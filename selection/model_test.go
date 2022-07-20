package selection_test

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/erikgeiser/promptkit"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/test"
	"github.com/muesli/termenv"
)

func TestSelectSecond(t *testing.T) {
	t.Parallel()

	s := selection.New("foo:", []string{"a", "b", "c"})
	s.ColorProfile = termenv.TrueColor
	m := selection.NewModel(s)

	test.Run(t, m, tea.KeyDown)
	assertNoError(t, m)
	test.AssertGoldenView(t, m, "select_second.golden")

	choice := getChoice(t, m)
	if choice != "b" {
		t.Errorf("unexpected choice: %v, expected b", choice)
	}

	test.Update(t, m, tea.KeyEnter)
	test.AssertGoldenView(t, m, "select_second_confirmed.golden")
}

func TestPaginate(t *testing.T) {
	t.Parallel()

	s := selection.New("foo:", []string{"First1", "First2", "Second1"})
	s.ColorProfile = termenv.TrueColor
	s.PageSize = 2
	m := selection.NewModel(s)

	test.Run(t, m)
	assertNoError(t, m)
	test.AssertGoldenView(t, m, "paginate.golden")

	view := m.View()
	strippedView := test.StripANSI(view)

	if strings.Contains(strippedView, "Second1") {
		t.Errorf("initial paginated view contains element of second page:\n%s", view)
	}

	test.Update(t, m, tea.KeyEnter)
	test.AssertGoldenView(t, m, "paginate_confirmed.golden")
}

func TestPaginatePush(t *testing.T) {
	t.Parallel()

	m := selection.NewModel(selection.New("foo:",
		[]string{
			"First1", "First2",
			"Second1", "Second2",
		}))
	m.PageSize = 2
	m.ColorProfile = termenv.TrueColor

	test.Run(t, m, tea.KeyDown, tea.KeyDown)
	assertNoError(t, m)
	test.AssertGoldenView(t, m, "paginate_push.golden")

	view := m.View()
	strippedView := test.StripANSI(view)

	if !strings.Contains(strippedView, "Second1") {
		t.Errorf("scolled view does not contain element of second page:\n%s", view)
	}

	if strings.Contains(strippedView, "Second2") {
		t.Errorf("scolled view contains \"Second2\" before scrolling that far")
	}

	if strings.Contains(strippedView, "First1") {
		t.Errorf("scolled view contains \"First1\" from first page")
	}

	test.Update(t, m, tea.KeyEnter)
	test.AssertGoldenView(t, m, "paginate_push_confirmed.golden")
}

func TestPaginateScroll(t *testing.T) {
	t.Parallel()

	m := selection.NewModel(selection.New("foo:", []string{
		"First1", "First2",
		"Second1", "Second2",
	}))
	m.PageSize = 2
	m.ColorProfile = termenv.TrueColor

	test.Run(t, m, tea.KeyPgDown)
	assertNoError(t, m)
	test.AssertGoldenView(t, m, "paginate_scroll.golden")

	view := m.View()
	strippedView := test.StripANSI(view)

	if !strings.Contains(strippedView, "Second1") {
		t.Errorf("scolled view does not contain element of second page:\n%s", view)
	}

	if strings.Contains(strippedView, "Second2") {
		t.Errorf("scolled view contains \"Second2\" before scrolling that far")
	}

	if strings.Contains(strippedView, "First1") {
		t.Errorf("scolled view contains \"First1\" from first page")
	}

	test.Update(t, m, tea.KeyEnter)
	test.AssertGoldenView(t, m, "paginate_scroll_confirmed.golden")
}

func TestPaginateLast(t *testing.T) {
	t.Parallel()

	m := selection.NewModel(selection.New("foo:",
		[]string{
			"First1", "First2",
			"Second1", "Second2",
		}))
	m.PageSize = 2
	m.ColorProfile = termenv.TrueColor

	test.Run(t, m, tea.KeyPgDown, tea.KeyPgDown, tea.KeyPgDown, tea.KeyPgDown,
		tea.KeyDown, tea.KeyDown, tea.KeyDown, tea.KeyDown, tea.KeyDown,
		tea.KeyPgDown, tea.KeyPgDown, tea.KeyPgDown, tea.KeyPgDown)
	assertNoError(t, m)
	test.AssertGoldenView(t, m, "paginate_last.golden")

	choice := getChoice(t, m)
	if choice != "Second2" {
		t.Errorf("unexpected selected element: %v", choice)
	}

	test.Update(t, m, tea.KeyEnter)
	test.AssertGoldenView(t, m, "paginate_last_confirmed.golden")
}

func TestFilter(t *testing.T) {
	t.Parallel()

	m := selection.NewModel(selection.New("foo:", []string{
		"AAA", "BBB", "CCC1", "CCC2", "DDD",
	}))
	m.PageSize = 2
	m.ColorProfile = termenv.TrueColor

	inputs := append(test.MsgsFromText("CC"), tea.KeyDown)
	test.Run(t, m, inputs...)
	assertNoError(t, m)
	test.AssertGoldenView(t, m, "filter.golden")

	choice := getChoice(t, m)
	if choice != "CCC2" {
		t.Errorf("unexpected selected element: %v", choice)
	}

	view := m.View()
	strippedView := test.StripANSI(view)

	if !strings.Contains(strippedView, "CCC1") {
		t.Errorf("filtered view does not contain first element that matches filter:\n%s",
			view)
	}

	if !strings.Contains(strippedView, "CCC2") {
		t.Errorf("filtered view does not contain first element that matches filter:\n%s",
			view)
	}

	if strings.Contains(strippedView, "AAA") || strings.Contains(strippedView, "BBB") ||
		strings.Contains(strippedView, "DDD") {
		t.Errorf("filtered contains elements that do not match filter:\n%s", view)
	}

	test.Update(t, m, tea.KeyEnter)
	test.AssertGoldenView(t, m, "filter_confirmed.golden")
}

func TestNoFilter(t *testing.T) {
	t.Parallel()

	m := selection.NewModel(selection.New("foo:", []string{
		"AAA", "BBB", "CCC", "DDD",
	}))
	m.Filter = nil
	m.PageSize = 2
	m.ColorProfile = termenv.TrueColor

	inputs := append(test.MsgsFromText("CC"), tea.KeyDown)
	test.Run(t, m, inputs...)
	assertNoError(t, m)
	test.AssertGoldenView(t, m, "no_filter.golden")

	choice := getChoice(t, m)
	if choice != "BBB" {
		t.Errorf("unexpected selected element: %v", choice)
	}

	view := m.View()
	strippedView := test.StripANSI(view)

	if !strings.Contains(strippedView, "AAA") {
		t.Errorf("filtered view does not contain first element that matches filter:\n%s",
			view)
	}

	if !strings.Contains(strippedView, "BBB") {
		t.Errorf("filtered view does not contain first element that matches filter:\n%s",
			view)
	}

	if strings.Contains(strippedView, "CCC") || strings.Contains(strippedView, "DDD") {
		t.Errorf("filtered contains elements that do not match filter:\n%s", view)
	}

	test.Update(t, m, tea.KeyEnter)
	test.AssertGoldenView(t, m, "no_filter_confirmed.golden")
}

func TestAbort(t *testing.T) {
	t.Parallel()

	m := selection.NewModel(selection.New("foo:", []string{
		"a", "b", "c",
	}))
	m.ColorProfile = termenv.TrueColor

	test.Run(t, m, tea.KeyCtrlC)

	if m.Err == nil {
		t.Fatalf("aborting did not produce an error")
	}

	if !errors.Is(m.Err, promptkit.ErrAborted) {
		t.Fatalf("aborting produced %q instead of %q", m.Err, promptkit.ErrAborted)
	}

	test.AssertGoldenView(t, m, "abort.golden")
}

func TestSubmit(t *testing.T) {
	t.Parallel()

	m := selection.NewModel(selection.New("foo:", []string{
		"a", "b", "c",
	}))
	m.ColorProfile = termenv.TrueColor

	test.Run(t, m)
	assertNoError(t, m)

	cmd := test.Update(t, m, tea.KeyEnter)
	if cmd == nil || cmd() != tea.Quit() {
		t.Errorf("enter did not produce quit signal")
	}

	test.AssertGoldenView(t, m, "submit.golden")
}

func TestLoopCursorTopToBottom(t *testing.T) {
	t.Parallel()

	choices := []string{
		"a", "b", "c", "d", "lastelement",
	}
	lastElement := choices[len(choices)-1]

	m := selection.NewModel(selection.New("foo:", choices))
	m.ColorProfile = termenv.TrueColor
	m.LoopCursor = true

	test.Run(t, m, tea.KeyUp)
	assertNoError(t, m)

	if getChoice(t, m) != lastElement {
		t.Fatalf("value did not loop to the last element %q but %q",
			lastElement, getChoice(t, m))
	}

	test.AssertGoldenView(t, m, "loop_top_to_bottom.golden")
}

func TestLoopCursorBottomToTop(t *testing.T) {
	t.Parallel()

	choices := []string{
		"firstelement", "b", "c", "d", "lastelement",
	}
	firstElement := choices[0]
	lastElement := choices[len(choices)-1]

	m := selection.NewModel(selection.New("foo:", choices))
	m.ColorProfile = termenv.TrueColor
	m.LoopCursor = true

	test.Run(t, m, tea.KeyDown, tea.KeyDown, tea.KeyDown, tea.KeyDown)
	assertNoError(t, m)

	if value := getChoice(t, m); value != lastElement {
		t.Fatalf("value did not loop to the last element before looping but %q",
			value)
	}

	test.Update(t, m, tea.KeyDown)

	if value := getChoice(t, m); value != firstElement {
		t.Fatalf("value did not loop to the first element %q but %q",
			firstElement, value)
	}

	test.AssertGoldenView(t, m, "loop_bottom_to_top.golden")
}

func TestLoopCursorTopToBottomPaged(t *testing.T) {
	t.Parallel()

	choices := []string{
		"a", "b", "c", "d", "lastelement",
	}
	lastElement := choices[len(choices)-1]

	m := selection.NewModel(selection.New("foo:", choices))
	m.ColorProfile = termenv.TrueColor
	m.PageSize = 3
	m.LoopCursor = true

	test.Run(t, m)
	assertNoError(t, m)

	if strings.Contains(m.View(), lastElement) {
		t.Fatalf("last element is already shown before looping")
	}

	test.Update(t, m, tea.KeyUp)

	if value := getChoice(t, m); value != lastElement {
		t.Fatalf("value did not loop to the last element %q but %q",
			lastElement, value)
	}

	test.AssertGoldenView(t, m, "loop_top_to_bottom_paged.golden")
}

func TestLoopCursorBottomToTopPaged(t *testing.T) {
	t.Parallel()

	choices := []string{
		"firstelement", "b", "c", "d", "lastelement",
	}
	firstElement := choices[0]
	lastElement := choices[len(choices)-1]

	m := selection.NewModel(selection.New("foo:", choices))
	m.ColorProfile = termenv.TrueColor
	m.PageSize = 3
	m.LoopCursor = true

	test.Run(t, m, tea.KeyDown, tea.KeyDown, tea.KeyDown, tea.KeyDown)
	assertNoError(t, m)

	if strings.Contains(m.View(), firstElement) {
		t.Fatalf("first element is already shown before looping")
	}

	if value := getChoice(t, m); value != lastElement {
		t.Fatalf("value did not loop to the last element before looping but %q",
			value)
	}

	test.Update(t, m, tea.KeyDown)

	if value := getChoice(t, m); value != firstElement {
		t.Fatalf("value did not loop to the first element %q but %q",
			firstElement, value)
	}

	test.AssertGoldenView(t, m, "loop_bottom_to_top_paged.golden")
}

func getChoice[T any](tb testing.TB, m *selection.Model[T]) T {
	tb.Helper()

	v, err := m.Value()
	if err != nil {
		tb.Fatalf("value: %v", err)
	}

	return v
}

func assertNoError[T any](tb testing.TB, m *selection.Model[T]) {
	tb.Helper()

	if m.Err != nil {
		tb.Fatalf("model contains error: %v", m.Err)
	}
}
