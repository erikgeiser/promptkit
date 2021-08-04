package textinput_test

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/erikgeiser/promptkit"
	"github.com/erikgeiser/promptkit/test"
	"github.com/erikgeiser/promptkit/textinput"
)

func TestEnterText(t *testing.T) {
	t.Parallel()

	m := textinput.NewModel(textinput.New("foo:"))
	m.Placeholder = "placeholder"

	input := "bar"

	test.Run(t, m, test.MsgsFromText(input)...)
	assertNoError(t, m)
	test.AssertGoldenView(t, m, "input.golden")

	value := getValue(t, m)
	if value != input {
		t.Errorf("unexpected value: %q, expected %q", value, input)
	}

	view := m.View()
	strippedView := test.StripANSI(view)

	if strings.Contains(strippedView, m.Placeholder) {
		t.Errorf("placeholder %q is rendered after text input:\n%s",
			m.Placeholder, test.Indent(view))
	}
}

func TestHidden(t *testing.T) {
	t.Parallel()

	m := textinput.NewModel(textinput.New("password?"))
	m.Hidden = true
	m.HideMask = 'X'

	input := "hunter2"

	test.Run(t, m, test.MsgsFromText(input)...)
	assertNoError(t, m)
	test.AssertGoldenView(t, m, "hidden.golden")

	view := m.View()
	strippedView := test.StripANSI(view)

	if strings.Contains(strippedView, input) {
		t.Errorf("hidden view contains input %q:\n%s", input, test.Indent(view))
	}

	if !strings.Contains(view, strings.Repeat(string(m.HideMask), len(input))) {
		t.Errorf("hidden view does not contain masked input:\n%s", test.Indent(view))
	}

	value := getValue(t, m)
	if value != input {
		t.Errorf("unexpected value: %q, expected %q", value, input)
	}
}

func TestPlaceholder(t *testing.T) {
	t.Parallel()

	placeholder := "enter some text"

	m := textinput.NewModel(textinput.New("Text:"))
	m.Placeholder = placeholder

	test.Run(t, m)
	assertNoError(t, m)
	test.AssertGoldenView(t, m, "placeholder.golden")

	view := m.View()
	strippedView := test.StripANSI(view)

	if !strings.Contains(strippedView, placeholder) {
		t.Errorf("placeholder %q was not rendered:\n%s", placeholder, test.Indent(view))
	}

	value := getValue(t, m)
	if value != "" {
		t.Errorf("value not empty: %s", value)
	}
}

func TestInitialValue(t *testing.T) {
	t.Parallel()

	initialValue := "some text"

	m := textinput.NewModel(textinput.New("question?"))
	m.InitialValue = initialValue
	m.Placeholder = "placeholder"

	test.Run(t, m)
	assertNoError(t, m)
	test.AssertGoldenView(t, m, "initial_value.golden")

	view := m.View()
	strippedView := test.StripANSI(view)

	if strings.Contains(strippedView, m.Placeholder) {
		t.Errorf("placeholder %q was rendered:\n%s", m.Placeholder, test.Indent(view))
	}

	value := getValue(t, m)
	if value != initialValue {
		t.Errorf("value %q is not initial value %q", value, initialValue)
	}
}

func TestModifiedInitialValue(t *testing.T) {
	t.Parallel()

	initialValue := "some test"
	modifiedInitialValue := "some text"

	m := textinput.NewModel(textinput.New("Text:"))
	m.InitialValue = initialValue

	test.Run(t, m, tea.KeyLeft, tea.KeyBackspace, test.KeyMsg('x'))
	assertNoError(t, m)
	test.AssertGoldenView(t, m, "modified_initial_value.golden")

	view := m.View()
	strippedView := test.StripANSI(view)
	value := getValue(t, m)

	if value == initialValue {
		t.Errorf("value %q is still initial value and was not modified to %q",
			value, modifiedInitialValue)
	}

	if strings.Contains(strippedView, initialValue) {
		t.Errorf("view still contains initial value:\n%s", test.Indent(view))
	}

	if value != modifiedInitialValue {
		t.Errorf("value %q is not modified initial value %q",
			value, modifiedInitialValue)
	}

	if !strings.Contains(strippedView, modifiedInitialValue) {
		t.Errorf("view does not contain modified initial value %q:\n%s",
			modifiedInitialValue, test.Indent(view))
	}
}

func TestTemplate(t *testing.T) {
	t.Parallel()

	separator := "|"

	m := textinput.NewModel(textinput.New("password?"))
	m.Template = `{{ print .Prompt Separator .Input}}`
	m.ExtendedTemplateScope["Separator"] = func() string { return separator }

	test.Run(t, m, tea.KeyLeft, tea.KeyBackspace, test.KeyMsg('s'))
	assertNoError(t, m)
	test.AssertGoldenView(t, m, "template.golden")

	view := m.View()
	if !strings.Contains(test.StripANSI(view), separator) {
		t.Errorf("sparator was not rendered:\n%s", test.Indent(view))
	}
}

func TestAbort(t *testing.T) {
	t.Parallel()

	m := textinput.NewModel(textinput.New("Question?"))
	m.Validate = nil

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

	m := textinput.NewModel(textinput.New("foo:"))
	m.Validate = nil

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

func TestValidate(t *testing.T) {
	t.Parallel()

	m := textinput.NewModel(textinput.New("foo:"))

	test.Run(t, m)
	assertNoError(t, m)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Errorf("enter on input that does not validate did not produce a no-op")
	}

	_, _ = m.Update(test.KeyMsg('x'))

	_, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil || cmd() != tea.Quit() {
		t.Errorf("enter on input that validates did not produce quit signal")
	}
}

func getValue(tb testing.TB, m *textinput.Model) string {
	tb.Helper()

	v, err := m.Value()
	if err != nil {
		tb.Fatalf("value: %v", err)
	}

	return v
}

func assertNoError(tb testing.TB, m *textinput.Model) {
	tb.Helper()

	if m.Err != nil {
		tb.Fatalf("model contains error: %v", m.Err)
	}
}
