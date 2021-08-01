package selection_test

import (
	"flag"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/test"
)

var update = flag.Bool("update", false, "update the golden files")

func TestSelectSecond(t *testing.T) {
	t.Parallel()

	m := selection.NewModel(selection.New("foo:",
		selection.Choices([]string{"a", "b", "c"})))

	test.Run(t, m, tea.KeyDown)
	assertNoError(t, m)
	test.AssertGoldenView(t, m, "select_second.golden", *update)

	choice := getChoice(t, m)
	if choice.Value != "b" {
		t.Errorf("unexpected choice: %v, expected b", choice.Value)
	}
}

func TestPaginate(t *testing.T) {
	t.Parallel()

	m := selection.NewModel(selection.New("foo:",
		selection.Choices([]string{"First1", "First2", "Second1"})))
	m.PageSize = 2

	test.Run(t, m)
	assertNoError(t, m)
	test.AssertGoldenView(t, m, "paginate.golden", *update)

	view := m.View()
	strippedView := test.StripANSI(view)

	if strings.Contains(strippedView, "Second1") {
		t.Errorf("initial paginated view contains element of second page:\n%s", view)
	}
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
	test.AssertGoldenView(t, m, "paginate_push.golden", *update)

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
	test.AssertGoldenView(t, m, "paginate_scroll.golden", *update)

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
	test.AssertGoldenView(t, m, "paginate_last.golden", *update)

	choice := getChoice(t, m)
	if choice.Value != "Second2" {
		t.Errorf("unexpected selected element: %v", choice.Value)
	}
}

func getChoice(tb testing.TB, m *selection.Model) *selection.Choice {
	tb.Helper()

	v, err := m.Choice()
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
