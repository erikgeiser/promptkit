package confirmation_test

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/erikgeiser/promptkit"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/test"
)

func TestDefaultYes(t *testing.T) {
	t.Parallel()

	c := confirmation.New("ready?")
	c.DefaultValue = confirmation.Yes
	m := confirmation.NewModel(c)

	test.Run(t, m, tea.KeyEnter)
	assertNoError(t, m)

	value := getValue(t, m)
	if !value {
		t.Errorf("default Yes produced a No")
	}
}

func TestDefaultNo(t *testing.T) {
	t.Parallel()

	c := confirmation.New("ready?")
	c.DefaultValue = confirmation.No
	m := confirmation.NewModel(c)

	test.Run(t, m, tea.KeyEnter)
	assertNoError(t, m)

	value := getValue(t, m)
	if value {
		t.Errorf("default No produced a Yes")
	}
}

func TestDefaultUndecided(t *testing.T) {
	t.Parallel()

	c := confirmation.New("ready?")
	c.DefaultValue = confirmation.Undecided
	m := confirmation.NewModel(c)

	test.Run(t, m)
	assertNoError(t, m)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Errorf("enter when undecided not produce a no-op but a %v", cmd)
	}

	v, err := m.Value()
	if err == nil {
		t.Errorf("getting value before deciding did not return an error but %v", v)
	}
}

func TestImmediatelyChooseYes(t *testing.T) {
	t.Parallel()

	m := confirmation.NewModel(confirmation.New("ready?"))

	test.Run(t, m)
	assertNoError(t, m)

	_, cmd := m.Update(test.KeyMsg('y'))
	if cmd == nil || cmd() != tea.Quit() {
		t.Errorf("direct answer selection did not result in quit message but in %v", cmd)
	}

	if !getValue(t, m) {
		t.Errorf("value is not Yes after entering y")
	}
}

func TestImmediatelyChooseNo(t *testing.T) {
	t.Parallel()

	m := confirmation.NewModel(confirmation.New("ready?"))

	test.Run(t, m)
	assertNoError(t, m)

	_, cmd := m.Update(test.KeyMsg('n'))
	if cmd == nil || cmd() != tea.Quit() {
		t.Errorf("direct answer selection did not result in quit message but in %v", cmd)
	}

	if getValue(t, m) {
		t.Errorf("value is not No after entering n")
	}
}

func TestToggle(t *testing.T) {
	t.Parallel()

	m := confirmation.NewModel(confirmation.New("ready?"))

	test.Run(t, m)
	assertNoError(t, m)

	v, err := m.Value()
	if err == nil {
		t.Fatalf("getting value before decision did not produce error but %v", v)
	}

	test.AssertGoldenView(t, m, "toggle_before.golden")

	m.Update(tea.KeyMsg{Type: tea.KeyTab})

	if !getValue(t, m) {
		t.Fatalf("toggle did not transition from Undecided to Yes")
	}

	test.AssertGoldenView(t, m, "toggle_once.golden")

	m.Update(tea.KeyMsg{Type: tea.KeyTab})

	if getValue(t, m) {
		t.Fatalf("toggle did not transition from Yes to No")
	}

	test.AssertGoldenView(t, m, "toggle_twice.golden")

	m.Update(tea.KeyMsg{Type: tea.KeyTab})

	if !getValue(t, m) {
		t.Fatalf("toggle did not transition from No to Yes")
	}
}

func TestSelectYes(t *testing.T) {
	t.Parallel()

	m := confirmation.NewModel(confirmation.New("ready?"))

	test.Run(t, m, tea.KeyLeft)
	assertNoError(t, m)

	if !getValue(t, m) {
		t.Fatalf("key left did not select yes")
	}
}

func TestSelectNo(t *testing.T) {
	t.Parallel()

	m := confirmation.NewModel(confirmation.New("ready?"))

	test.Run(t, m, tea.KeyRight)
	assertNoError(t, m)

	if getValue(t, m) {
		t.Fatalf("key left did not select yes")
	}
}

func TestAbort(t *testing.T) {
	t.Parallel()

	m := confirmation.NewModel(confirmation.New("ready?"))

	test.Run(t, m, tea.KeyCtrlC)

	if m.Err == nil {
		t.Fatalf("aborting did not produce an error")
	}

	if !errors.Is(m.Err, promptkit.ErrAborted) {
		t.Fatalf("aborting produced %q instead of %q", m.Err, promptkit.ErrAborted)
	}
}

func TestSubmit(t *testing.T) {
	t.Parallel()

	c := confirmation.New("ready?")
	c.DefaultValue = confirmation.Yes
	m := confirmation.NewModel(c)

	test.Run(t, m)
	assertNoError(t, m)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil || cmd() != tea.Quit() {
		t.Errorf("enter did not produce quit signal")
	}

	if m.View() != "" {
		t.Errorf("view not empty after quitting")
	}
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
