package textinput_test

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/erikgeiser/promptkit"
	"github.com/erikgeiser/promptkit/test"
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/muesli/termenv"
)

func TestEnterText(t *testing.T) {
	t.Parallel()

	m := textinput.NewModel(textinput.New("foo:"))
	m.Placeholder = "placeholder"
	m.ColorProfile = termenv.TrueColor

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

	test.Update(t, m, tea.KeyEnter)
	test.AssertGoldenView(t, m, "input_confirmed.golden")
}

func TestHidden(t *testing.T) {
	t.Parallel()

	m := textinput.NewModel(textinput.New("password?"))
	m.Hidden = true
	m.HideMask = 'X'
	m.ColorProfile = termenv.TrueColor

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

	test.Update(t, m, tea.KeyEnter)
	test.AssertGoldenView(t, m, "hidden_confirmed.golden")
}

func TestPlaceholder(t *testing.T) {
	t.Parallel()

	placeholder := "enter some text"

	m := textinput.NewModel(textinput.New("Text:"))
	m.Placeholder = placeholder
	m.ColorProfile = termenv.TrueColor

	test.Run(t, m)
	assertNoError(t, m)
	test.AssertGoldenView(t, m, "placeholder.golden")

	view := m.View()
	strippedView := test.StripANSI(view)

	if !strings.Contains(strippedView, placeholder) {
		t.Errorf("placeholder %q was not rendered:\n%s", placeholder, test.Indent(view))
	}

	test.Update(t, m, tea.KeyEnter)
	test.AssertGoldenView(t, m, "placeholder_confirmed.golden")
}

func TestInitialValue(t *testing.T) {
	t.Parallel()

	initialValue := "some text"

	m := textinput.NewModel(textinput.New("question?"))
	m.InitialValue = initialValue
	m.Placeholder = "placeholder"
	m.ColorProfile = termenv.TrueColor

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

	test.Update(t, m, tea.KeyEnter)
	test.AssertGoldenView(t, m, "initial_value_confirmed.golden")
}

func TestModifiedInitialValue(t *testing.T) {
	t.Parallel()

	initialValue := "some test"
	modifiedInitialValue := "some text"

	m := textinput.NewModel(textinput.New("Text:"))
	m.InitialValue = initialValue
	m.ColorProfile = termenv.TrueColor

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

	test.Update(t, m, tea.KeyEnter)
	test.AssertGoldenView(t, m, "modified_initial_value_confirmed.golden")
}

func TestTemplate(t *testing.T) {
	t.Parallel()

	separator := "|"

	m := textinput.NewModel(textinput.New("name?"))
	m.Template = `{{ print .Prompt Separator .Input}}`
	m.ResultTemplate = `my name is {{ .FinalValue }}`
	m.ExtendedTemplateFuncs["Separator"] = func() string { return separator }
	m.ColorProfile = termenv.TrueColor

	test.Run(t, m, tea.KeyLeft, tea.KeyBackspace, test.KeyMsg('s'))
	assertNoError(t, m)
	test.AssertGoldenView(t, m, "template.golden")

	view := m.View()
	if !strings.Contains(test.StripANSI(view), separator) {
		t.Errorf("sparator was not rendered:\n%s", test.Indent(view))
	}

	test.Update(t, m, tea.KeyEnter)
	test.AssertGoldenView(t, m, "template_confirmed.golden")
}

func TestAbort(t *testing.T) {
	t.Parallel()

	m := textinput.NewModel(textinput.New("Question?"))
	m.Validate = nil
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

	m := textinput.NewModel(textinput.New("foo:"))
	m.ResultTemplate = `result: {{ .FinalValue }}`
	m.Validate = nil
	m.ColorProfile = termenv.TrueColor

	test.Run(t, m)
	assertNoError(t, m)

	cmd := test.Update(t, m, tea.KeyEnter)
	if cmd == nil || cmd() != tea.Quit() {
		t.Errorf("enter did not produce quit signal")
	}

	test.AssertGoldenView(t, m, "submit.golden")
}

func TestValidate(t *testing.T) {
	t.Parallel()

	m := textinput.NewModel(textinput.New("foo:"))
	m.ColorProfile = termenv.TrueColor

	test.Run(t, m)
	assertNoError(t, m)

	cmd := test.Update(t, m, tea.KeyEnter)
	if cmd != nil {
		t.Errorf("enter on input that does not validate did not produce a no-op")
	}

	test.Update(t, m, test.KeyMsg('x'))

	cmd = test.Update(t, m, tea.KeyEnter)
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
