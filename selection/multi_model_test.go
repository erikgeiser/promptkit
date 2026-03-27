package selection_test

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/erikgeiser/promptkit"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/test"
	"github.com/muesli/termenv"
)

var keySpace = tea.KeyPressMsg{Code: ' '}

func TestMultiToggle(t *testing.T) {
	t.Parallel()

	m := selection.NewMultiModel(selection.NewMulti("foo:", []string{"a", "b", "c"}))
	m.ColorProfile = termenv.TrueColor

	test.Run(t, m, keyDown, keySpace)
	assertMultiNoError(t, m)
	test.AssertGoldenView(t, m, "multi_toggle.golden")

	values := getValues(t, m)
	if len(values) != 1 || values[0] != "b" {
		t.Errorf("unexpected values: %v, expected [b]", values)
	}
}

func TestMultiToggleMultiple(t *testing.T) {
	t.Parallel()

	m := selection.NewMultiModel(selection.NewMulti("foo:", []string{"a", "b", "c"}))
	m.ColorProfile = termenv.TrueColor

	test.Run(t, m, keySpace, keyDown, keySpace)
	assertMultiNoError(t, m)
	test.AssertGoldenView(t, m, "multi_toggle_multiple.golden")

	values := getValues(t, m)
	if len(values) != 2 || values[0] != "a" || values[1] != "b" {
		t.Errorf("unexpected values: %v, expected [a b]", values)
	}
}

func TestMultiUntoggle(t *testing.T) {
	t.Parallel()

	m := selection.NewMultiModel(selection.NewMulti("foo:", []string{"a", "b", "c"}))
	m.ColorProfile = termenv.TrueColor

	// toggle a, move down, toggle b, move back up, untoggle a
	test.Run(t, m, keySpace, keyDown, keySpace, keyUp, keySpace)
	assertMultiNoError(t, m)

	values := getValues(t, m)
	if len(values) != 1 || values[0] != "b" {
		t.Errorf("unexpected values after untoggle: %v, expected [b]", values)
	}
}

func TestMultiConfirm(t *testing.T) {
	t.Parallel()

	m := selection.NewMultiModel(selection.NewMulti("foo:", []string{"a", "b", "c"}))
	m.ColorProfile = termenv.TrueColor

	test.Run(t, m, keySpace)
	assertMultiNoError(t, m)

	cmd := test.Update(t, m, keyEnter)
	if cmd == nil || cmd() != tea.Quit() {
		t.Errorf("enter after marking did not produce quit signal")
	}

	test.AssertGoldenView(t, m, "multi_confirm.golden")
}

func TestMultiMinSelectionsBlocked(t *testing.T) {
	t.Parallel()

	m := selection.NewMultiModel(selection.NewMulti("foo:", []string{"a", "b", "c"}))
	m.ColorProfile = termenv.TrueColor
	m.MinSelections = 2

	test.Run(t, m, keySpace)
	assertMultiNoError(t, m)

	// only 1 marked, MinSelections=2 — Enter should be a no-op
	cmd := test.Update(t, m, keyEnter)
	if cmd != nil && cmd() == tea.Quit() {
		t.Errorf("enter with insufficient selections should not quit")
	}

	// mark second item and confirm
	test.Update(t, m, keyDown)
	test.Update(t, m, keySpace)

	cmd = test.Update(t, m, keyEnter)
	if cmd == nil || cmd() != tea.Quit() {
		t.Errorf("enter after reaching MinSelections did not produce quit signal")
	}
}

func TestMultiMaxSelections(t *testing.T) {
	t.Parallel()

	m := selection.NewMultiModel(selection.NewMulti("foo:", []string{"a", "b", "c"}))
	m.ColorProfile = termenv.TrueColor
	m.MaxSelections = 2

	test.Run(t, m, keySpace, keyDown, keySpace, keyDown, keySpace)
	assertMultiNoError(t, m)

	// only a and b should be marked; c should be rejected by MaxSelections
	values := getValues(t, m)
	if len(values) != 2 {
		t.Errorf("unexpected number of values: %d, expected 2 (MaxSelections=2)", len(values))
	}
}

func TestMultiAbort(t *testing.T) {
	t.Parallel()

	m := selection.NewMultiModel(selection.NewMulti("foo:", []string{"a", "b", "c"}))
	m.ColorProfile = termenv.TrueColor

	test.Run(t, m, keyCtrlC)

	if m.Err == nil {
		t.Fatalf("aborting did not produce an error")
	}

	if !errors.Is(m.Err, promptkit.ErrAborted) {
		t.Fatalf("aborting produced %q instead of %q", m.Err, promptkit.ErrAborted)
	}
}

func TestMultiFilter(t *testing.T) {
	t.Parallel()

	m := selection.NewMultiModel(selection.NewMulti("foo:", []string{
		"AAA", "BBB", "CCC1", "CCC2", "DDD",
	}))
	m.PageSize = 2
	m.ColorProfile = termenv.TrueColor

	// filter to CCC items, mark both
	inputs := append(test.MsgsFromText("CC"), keySpace, keyDown, keySpace)
	test.Run(t, m, inputs...)
	assertMultiNoError(t, m)

	view := m.View().Content
	strippedView := test.StripANSI(view)

	if !strings.Contains(strippedView, "CCC1") || !strings.Contains(strippedView, "CCC2") {
		t.Errorf("filtered view should show CCC1 and CCC2:\n%s", view)
	}

	// clear filter — CCC1 and CCC2 should still be marked
	test.Update(t, m, keyEsc)

	values := getValues(t, m)
	if len(values) != 2 {
		t.Errorf("selections should persist after clearing filter, got %v", values)
	}
}

func TestMultiSelectionsPersistThroughFilter(t *testing.T) {
	t.Parallel()

	m := selection.NewMultiModel(selection.NewMulti("foo:", []string{"alpha", "beta", "gamma"}))
	m.ColorProfile = termenv.TrueColor

	// mark alpha
	test.Run(t, m, keySpace)
	assertMultiNoError(t, m)

	// filter to "beta" — alpha is marked so it appears in the offscreen section above the list
	test.Update(t, m, tea.KeyPressMsg{Code: 'b', Text: "b"})
	test.Update(t, m, tea.KeyPressMsg{Code: 'e', Text: "e"})

	view := m.View().Content
	if !strings.Contains(test.StripANSI(view), "alpha") {
		t.Errorf("marked item alpha should appear in offscreen section during filter:\n%s", view)
	}

	// cursor is on beta (position 0 of filtered list), mark it
	test.Update(t, m, keySpace)

	// clear filter
	test.Update(t, m, keyEsc)

	// both alpha and beta should be marked
	values := getValues(t, m)
	if len(values) != 2 {
		t.Errorf("selections should persist through filter, got %v", values)
	}
}

func TestMultiMarkedItemsVisibleWhenFiltered(t *testing.T) {
	t.Parallel()

	m := selection.NewMultiModel(selection.NewMulti("foo:", []string{"apple", "apricot", "banana"}))
	m.ColorProfile = termenv.TrueColor

	// mark apple (first item)
	test.Run(t, m, keySpace)
	assertMultiNoError(t, m)

	// filter to "ban" — only banana matches; apple is marked so it appears in the offscreen section
	for _, msg := range test.MsgsFromText("ban") {
		test.Update(t, m, msg)
	}

	view := m.View().Content
	strippedView := test.StripANSI(view)

	if !strings.Contains(strippedView, "apple") {
		t.Errorf("marked item apple should appear in offscreen section when filtered:\n%s", view)
	}

	if !strings.Contains(strippedView, "banana") {
		t.Errorf("matching item banana should be visible when filtered:\n%s", view)
	}

	if strings.Contains(strippedView, "apricot") {
		t.Errorf("unmarked non-matching item apricot should be hidden when filtered:\n%s", view)
	}
}

func TestMultiLoopCursor(t *testing.T) {
	t.Parallel()

	m := selection.NewMultiModel(selection.NewMulti("foo:", []string{"a", "b", "c"}))
	m.ColorProfile = termenv.TrueColor
	m.LoopCursor = true

	test.Run(t, m, keyUp)
	assertMultiNoError(t, m)

	// should have looped to last element
	values := getValues(t, m)
	if len(values) != 0 {
		t.Errorf("looping should not mark items: %v", values)
	}

	// cursor should be at last element — mark it
	test.Update(t, m, keySpace)

	values = getValues(t, m)
	if len(values) != 1 || values[0] != "c" {
		t.Errorf("expected [c] after marking last element, got %v", values)
	}
}

func TestMultiPaginate(t *testing.T) {
	t.Parallel()

	m := selection.NewMultiModel(selection.NewMulti("foo:", []string{
		"First1", "First2", "Second1", "Second2",
	}))
	m.PageSize = 2
	m.ColorProfile = termenv.TrueColor

	// mark First1, page down (cursor lands on First2), move down to Second1, mark it
	test.Run(t, m, keySpace, keyPgDown, keyDown, keySpace)
	assertMultiNoError(t, m)

	values := getValues(t, m)
	if len(values) != 2 || values[0] != "First1" || values[1] != "Second1" {
		t.Errorf("unexpected values across pages: %v", values)
	}
}

func TestMultiResultOrder(t *testing.T) {
	t.Parallel()

	m := selection.NewMultiModel(selection.NewMulti("foo:", []string{"a", "b", "c", "d"}))
	m.ColorProfile = termenv.TrueColor

	// mark in reverse order: d, c, b
	test.Run(t, m,
		keyDown, keyDown, keyDown, keySpace, // mark d
		keyUp, keySpace, // mark c
		keyUp, keySpace, // mark b
	)
	assertMultiNoError(t, m)

	values := getValues(t, m)

	// result should be in original order regardless of toggle order
	if len(values) != 3 || values[0] != "b" || values[1] != "c" || values[2] != "d" {
		t.Errorf("values not in original order: %v, expected [b c d]", values)
	}
}

func TestMultiSelectionAnnotation(t *testing.T) {
	t.Parallel()

	choices := []string{"a", "b", "c", "d"}

	cases := []struct {
		name          string
		minSelections int
		maxSelections int
		keys          []tea.Msg // keys to send after init
		golden        string
	}{
		{
			name:          "min_and_max_none_selected",
			minSelections: 2,
			maxSelections: 4,
			keys:          nil,
			golden:        "annotation_min_max_none_selected.golden",
		},
		{
			name:          "min_and_max_within_range",
			minSelections: 2,
			maxSelections: 4,
			keys:          []tea.Msg{keySpace, keyDown, keySpace},
			golden:        "annotation_min_max_within_range.golden",
		},
		{
			name:          "min_and_max_below_min",
			minSelections: 2,
			maxSelections: 4,
			keys:          []tea.Msg{keySpace},
			golden:        "annotation_min_max_below_min.golden",
		},
		{
			name:          "only_min_below",
			minSelections: 2,
			maxSelections: 0,
			keys:          []tea.Msg{keySpace},
			golden:        "annotation_only_min_below.golden",
		},
		{
			name:          "only_min_satisfied",
			minSelections: 2,
			maxSelections: 0,
			keys:          []tea.Msg{keySpace, keyDown, keySpace},
			golden:        "annotation_only_min_satisfied.golden",
		},
		{
			name:          "only_max_within",
			minSelections: 0,
			maxSelections: 3,
			keys:          []tea.Msg{keySpace, keyDown, keySpace},
			golden:        "annotation_only_max_within.golden",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := selection.NewMultiModel(selection.NewMulti("Pick:", choices))
			m.ColorProfile = termenv.TrueColor
			m.MinSelections = tc.minSelections
			m.MaxSelections = tc.maxSelections

			test.Run(t, m, tc.keys...)
			assertMultiNoError(t, m)
			test.AssertGoldenView(t, m, tc.golden)
		})
	}
}

func TestMultiZeroMinSelections(t *testing.T) {
	t.Parallel()

	m := selection.NewMultiModel(selection.NewMulti("foo:", []string{"a", "b", "c"}))
	m.ColorProfile = termenv.TrueColor
	m.MinSelections = 0

	test.Run(t, m)
	assertMultiNoError(t, m)

	// should be able to confirm with nothing marked
	cmd := test.Update(t, m, keyEnter)
	if cmd == nil || cmd() != tea.Quit() {
		t.Errorf("enter with MinSelections=0 and no marks should quit")
	}

	values := getValues(t, m)
	if len(values) != 0 {
		t.Errorf("expected empty values, got %v", values)
	}
}

var keyEsc = tea.KeyPressMsg{Code: tea.KeyEsc}

func getValues[T any](tb testing.TB, m *selection.MultiModel[T]) []T {
	tb.Helper()

	v, err := m.Values()
	if err != nil {
		tb.Fatalf("values: %v", err)
	}

	return v
}

func assertMultiNoError[T any](tb testing.TB, m *selection.MultiModel[T]) {
	tb.Helper()

	if m.Err != nil {
		tb.Fatalf("model contains error: %v", m.Err)
	}
}
