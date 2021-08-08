package selection_test

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/erikgeiser/promptkit"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/test"
)

func TestSelectSecond(t *testing.T) {
	t.Parallel()

	m := selection.NewModel(selection.New("foo:",
		selection.Choices([]string{"a", "b", "c"})))

	test.Run(t, m, tea.KeyDown)
	assertNoError(t, m)
	test.AssertGoldenView(t, m, "select_second.golden")

	choice := getChoice(t, m)
	if choice.Value != "b" {
		t.Errorf("unexpected choice: %v, expected b", choice.Value)
	}

	test.Update(t, m, tea.KeyEnter)
	test.AssertGoldenView(t, m, "select_second_confirmed.golden")
}

func TestPaginate(t *testing.T) {
	t.Parallel()

	m := selection.NewModel(selection.New("foo:",
		selection.Choices([]string{"First1", "First2", "Second1"})))
	m.PageSize = 2

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
		selection.Choices([]string{
			"First1", "First2",
			"Second1", "Second2",
		})))
	m.PageSize = 2

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

	m := selection.NewModel(selection.New("foo:",
		selection.Choices([]string{
			"First1", "First2",
			"Second1", "Second2",
		})))
	m.PageSize = 2

	test.Run(t, m, tea.KeyRight)
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
		selection.Choices([]string{
			"First1", "First2",
			"Second1", "Second2",
		})))
	m.PageSize = 2

	test.Run(t, m, tea.KeyRight, tea.KeyRight, tea.KeyRight, tea.KeyRight,
		tea.KeyDown, tea.KeyDown, tea.KeyDown, tea.KeyDown, tea.KeyDown,
		tea.KeyRight, tea.KeyRight, tea.KeyRight, tea.KeyRight)
	assertNoError(t, m)
	test.AssertGoldenView(t, m, "paginate_last.golden")

	choice := getChoice(t, m)
	if choice.Value != "Second2" {
		t.Errorf("unexpected selected element: %v", choice.Value)
	}

	test.Update(t, m, tea.KeyEnter)
	test.AssertGoldenView(t, m, "paginate_last_confirmed.golden")
}

func TestFilter(t *testing.T) {
	t.Parallel()

	m := selection.NewModel(selection.New("foo:",
		selection.Choices([]string{
			"AAA", "BBB", "CCC1", "CCC2", "DDD",
		})))
	m.PageSize = 2

	inputs := append(test.MsgsFromText("CC"), tea.KeyDown)
	test.Run(t, m, inputs...)
	assertNoError(t, m)
	test.AssertGoldenView(t, m, "filter.golden")

	choice := getChoice(t, m)
	if choice.Value != "CCC2" {
		t.Errorf("unexpected selected element: %v", choice.Value)
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

	m := selection.NewModel(selection.New("foo:",
		selection.Choices([]string{
			"AAA", "BBB", "CCC", "DDD",
		})))
	m.Filter = nil
	m.PageSize = 2

	inputs := append(test.MsgsFromText("CC"), tea.KeyDown)
	test.Run(t, m, inputs...)
	assertNoError(t, m)
	test.AssertGoldenView(t, m, "no_filter.golden")

	choice := getChoice(t, m)
	if choice.Value != "BBB" {
		t.Errorf("unexpected selected element: %v", choice.Value)
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

	m := selection.NewModel(selection.New("foo:",
		selection.Choices([]string{
			"a", "b", "c",
		})))

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

	m := selection.NewModel(selection.New("foo:",
		selection.Choices([]string{
			"a", "b", "c",
		})))

	test.Run(t, m)
	assertNoError(t, m)

	cmd := test.Update(t, m, tea.KeyEnter)
	if cmd == nil || cmd() != tea.Quit() {
		t.Errorf("enter did not produce quit signal")
	}

	test.AssertGoldenView(t, m, "submit.golden")
}

func getChoice(tb testing.TB, m *selection.Model) *selection.Choice {
	tb.Helper()

	v, err := m.Value()
	if err != nil {
		tb.Fatalf("value: %v", err)
	}

	return v
}

func assertNoError(tb testing.TB, m *selection.Model) {
	tb.Helper()

	if m.Err != nil {
		tb.Fatalf("model contains error: %v", m.Err)
	}
}
