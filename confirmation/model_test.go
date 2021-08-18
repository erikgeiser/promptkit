package confirmation_test

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/erikgeiser/promptkit"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/test"
	"github.com/muesli/termenv"
)

func TestDefaultYes(t *testing.T) {
	t.Parallel()

	c := confirmation.New("ready?", confirmation.Yes)
	c.ColorProfile = termenv.TrueColor
	m := confirmation.NewModel(c)

	test.Run(t, m, tea.KeyEnter)
	assertNoError(t, m)

	value := getValue(t, m)
	if !value {
		t.Errorf("default Yes produced a No")
	}

	test.AssertGoldenView(t, m, "default_yes.golden")
}

func TestDefaultNo(t *testing.T) {
	t.Parallel()

	c := confirmation.New("ready?", confirmation.No)
	c.ColorProfile = termenv.TrueColor
	m := confirmation.NewModel(c)

	test.Run(t, m, tea.KeyEnter)
	assertNoError(t, m)

	value := getValue(t, m)
	if value {
		t.Errorf("default No produced a Yes")
	}

	test.AssertGoldenView(t, m, "default_no.golden")
}

func TestDefaultUndecided(t *testing.T) {
	t.Parallel()

	c := confirmation.New("ready?", confirmation.Undecided)
	c.ColorProfile = termenv.TrueColor
	m := confirmation.NewModel(c)

	test.Run(t, m)
	assertNoError(t, m)

	cmd := test.Update(t, m, tea.KeyEnter)
	if cmd != nil {
		t.Errorf("enter when undecided not produce a no-op but a %v", cmd)
	}

	v, err := m.Value()
	if err == nil {
		t.Errorf("getting value before deciding did not return an error but %v", v)
	}

	test.AssertGoldenView(t, m, "default_undecided.golden")
}

func TestDefaultNilIsUndecided(t *testing.T) {
	t.Parallel()

	c := confirmation.New("ready?", nil)
	c.ColorProfile = termenv.TrueColor
	m := confirmation.NewModel(c)

	test.Run(t, m)
	assertNoError(t, m)

	cmd := test.Update(t, m, tea.KeyEnter)
	if cmd != nil {
		t.Errorf("enter when undecided not produce a no-op but a %v", cmd)
	}

	v, err := m.Value()
	if err == nil {
		t.Errorf("getting value before deciding did not return an error but %v", v)
	}

	test.AssertGoldenView(t, m, "default_nil.golden")
}

func TestImmediatelyChooseYes(t *testing.T) {
	t.Parallel()

	c := confirmation.New("ready?", confirmation.Undecided)
	c.ColorProfile = termenv.TrueColor
	m := confirmation.NewModel(c)

	test.Run(t, m)
	assertNoError(t, m)

	cmd := test.Update(t, m, test.KeyMsg('y'))
	if cmd == nil || cmd() != tea.Quit() {
		t.Errorf("direct answer selection did not result in quit message but in %v", cmd)
	}

	if !getValue(t, m) {
		t.Errorf("value is not Yes after entering y")
	}

	test.AssertGoldenView(t, m, "choose_yes.golden")
}

func TestImmediatelyChooseNo(t *testing.T) {
	t.Parallel()

	c := confirmation.New("ready?", confirmation.Undecided)
	c.ColorProfile = termenv.TrueColor
	m := confirmation.NewModel(c)

	test.Run(t, m)
	assertNoError(t, m)

	cmd := test.Update(t, m, test.KeyMsg('n'))
	if cmd == nil || cmd() != tea.Quit() {
		t.Errorf("direct answer selection did not result in quit message but in %v", cmd)
	}

	if getValue(t, m) {
		t.Errorf("value is not No after entering n")
	}

	test.AssertGoldenView(t, m, "choose_no.golden")
}

func TestToggle(t *testing.T) {
	t.Parallel()

	c := confirmation.New("ready?", confirmation.Undecided)
	c.ColorProfile = termenv.TrueColor
	m := confirmation.NewModel(c)

	test.Run(t, m)
	assertNoError(t, m)

	v, err := m.Value()
	if err == nil {
		t.Fatalf("getting value before decision did not produce error but %v", v)
	}

	test.AssertGoldenView(t, m, "toggle_before.golden")

	test.Update(t, m, tea.KeyTab)

	if !getValue(t, m) {
		t.Fatalf("toggle did not transition from Undecided to Yes")
	}

	test.AssertGoldenView(t, m, "toggle_once.golden")

	test.Update(t, m, tea.KeyTab)

	if getValue(t, m) {
		t.Fatalf("toggle did not transition from Yes to No")
	}

	test.AssertGoldenView(t, m, "toggle_twice.golden")

	test.Update(t, m, tea.KeyTab)

	if !getValue(t, m) {
		t.Fatalf("toggle did not transition from No to Yes")
	}

	test.Update(t, m, tea.KeyEnter)
	test.AssertGoldenView(t, m, "toggle_confirmed.golden")
}

func TestSelectYes(t *testing.T) {
	t.Parallel()

	c := confirmation.New("ready?", confirmation.Undecided)
	c.ColorProfile = termenv.TrueColor
	m := confirmation.NewModel(c)

	test.Run(t, m, tea.KeyLeft)
	assertNoError(t, m)

	if !getValue(t, m) {
		t.Fatalf("key left did not select yes")
	}

	test.AssertGoldenView(t, m, "select_yes.golden")
	test.Update(t, m, tea.KeyEnter)
	test.AssertGoldenView(t, m, "select_yes_confirmed.golden")
}

func TestSelectNo(t *testing.T) {
	t.Parallel()

	c := confirmation.New("ready?", confirmation.Undecided)
	c.ColorProfile = termenv.TrueColor
	m := confirmation.NewModel(c)

	test.Run(t, m, tea.KeyRight)
	assertNoError(t, m)

	if getValue(t, m) {
		t.Fatalf("key left did not select yes")
	}

	test.AssertGoldenView(t, m, "select_no.golden")
	test.Update(t, m, tea.KeyEnter)
	test.AssertGoldenView(t, m, "select_no_confirmed.golden")
}

func TestAbort(t *testing.T) {
	t.Parallel()

	c := confirmation.New("ready?", confirmation.Undecided)
	c.ColorProfile = termenv.TrueColor
	m := confirmation.NewModel(c)

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

	c := confirmation.New("ready?", confirmation.Yes)
	c.ColorProfile = termenv.TrueColor
	m := confirmation.NewModel(c)

	test.Run(t, m)
	assertNoError(t, m)

	cmd := test.Update(t, m, tea.KeyEnter)
	if cmd == nil || cmd() != tea.Quit() {
		t.Errorf("enter did not produce quit signal")
	}

	test.AssertGoldenView(t, m, "submit.golden")
}

func TestTemplateYN(t *testing.T) {
	t.Parallel()

	c := confirmation.New("yes or no?", confirmation.Undecided)
	c.ColorProfile = termenv.TrueColor
	c.Template = confirmation.TemplateYN
	c.ResultTemplate = confirmation.ResultTemplateYN
	m := confirmation.NewModel(c)

	test.Run(t, m)
	assertNoError(t, m)

	test.AssertGoldenView(t, m, "templateyn_undecided.golden")

	test.Update(t, m, tea.KeyRight)
	assertNoError(t, m)

	test.AssertGoldenView(t, m, "templateyn_no.golden")

	test.Update(t, m, tea.KeyLeft)
	assertNoError(t, m)

	test.AssertGoldenView(t, m, "templateyn_yes.golden")

	cmd := test.Update(t, m, tea.KeyEnter)
	if cmd == nil || cmd() != tea.Quit() {
		t.Errorf("enter did not produce quit signal")
	}

	assertNoError(t, m)

	test.AssertGoldenView(t, m, "templateyn_result.golden")
}

func getValue(tb testing.TB, m *confirmation.Model) bool {
	tb.Helper()

	v, err := m.Value()
	if err != nil {
		tb.Fatalf("value: %v", err)
	}

	return v
}

func assertNoError(tb testing.TB, m *confirmation.Model) {
	tb.Helper()

	if m.Err != nil {
		tb.Fatalf("model contains error: %v", m.Err)
	}
}
