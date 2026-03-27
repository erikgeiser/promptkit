package textinput_test

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/erikgeiser/promptkit"
	"github.com/erikgeiser/promptkit/test"
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/muesli/termenv"
)

var (
	keyShiftEnter = tea.KeyPressMsg{Code: tea.KeyEnter, Mod: tea.ModShift}
)

func TestAreaEnterText(t *testing.T) {
	t.Parallel()

	m := textinput.NewAreaModel(textinput.NewArea("foo:"))
	m.Placeholder = "placeholder"
	m.ColorProfile = termenv.TrueColor

	input := "bar"

	test.Run(t, m, test.MsgsFromText(input)...)
	assertNoAreaError(t, m)
	test.AssertGoldenView(t, m, "area_input.golden")

	value := getAreaValue(t, m)
	if value != input {
		t.Errorf("unexpected value: %q, expected %q", value, input)
	}

	view := m.View().Content
	strippedView := test.StripANSI(view)

	if strings.Contains(strippedView, m.Placeholder) {
		t.Errorf("placeholder %q is rendered after text input:\n%s",
			m.Placeholder, test.Indent(view))
	}

	test.Update(t, m, keyShiftEnter)
	test.AssertGoldenView(t, m, "area_input_confirmed.golden")
}

func TestAreaMultiline(t *testing.T) {
	t.Parallel()

	m := textinput.NewAreaModel(textinput.NewArea("notes:"))
	m.ColorProfile = termenv.TrueColor

	msgs := append(test.MsgsFromText("line one"), keyEnter)
	msgs = append(msgs, test.MsgsFromText("line two")...)
	test.Run(t, m, msgs...)
	assertNoAreaError(t, m)
	test.AssertGoldenView(t, m, "area_multiline.golden")

	value := getAreaValue(t, m)
	if !strings.Contains(value, "line one") || !strings.Contains(value, "line two") {
		t.Errorf("value %q does not contain both lines", value)
	}

	test.Update(t, m, keyShiftEnter)
	test.AssertGoldenView(t, m, "area_multiline_confirmed.golden")
}

func TestAreaPlaceholder(t *testing.T) {
	t.Parallel()

	placeholder := "enter some text"

	m := textinput.NewAreaModel(textinput.NewArea("Text:"))
	m.Placeholder = placeholder
	m.ColorProfile = termenv.TrueColor

	test.Run(t, m)
	assertNoAreaError(t, m)
	test.AssertGoldenView(t, m, "area_placeholder.golden")

	view := m.View().Content
	strippedView := test.StripANSI(view)

	if !strings.Contains(strippedView, placeholder) {
		t.Errorf("placeholder %q was not rendered:\n%s", placeholder, test.Indent(view))
	}

	test.Update(t, m, keyShiftEnter)
	test.AssertGoldenView(t, m, "area_placeholder_confirmed.golden")
}

func TestAreaInitialValue(t *testing.T) {
	t.Parallel()

	initialValue := "some text"

	m := textinput.NewAreaModel(textinput.NewArea("question?"))
	m.InitialValue = initialValue
	m.Placeholder = "placeholder"
	m.ColorProfile = termenv.TrueColor

	test.Run(t, m)
	assertNoAreaError(t, m)
	test.AssertGoldenView(t, m, "area_initial_value.golden")

	view := m.View().Content
	strippedView := test.StripANSI(view)

	if strings.Contains(strippedView, m.Placeholder) {
		t.Errorf("placeholder %q was rendered:\n%s", m.Placeholder, test.Indent(view))
	}

	value := getAreaValue(t, m)
	if value != initialValue {
		t.Errorf("value %q is not initial value %q", value, initialValue)
	}

	test.Update(t, m, keyShiftEnter)
	test.AssertGoldenView(t, m, "area_initial_value_confirmed.golden")
}

func TestAreaAbort(t *testing.T) {
	t.Parallel()

	m := textinput.NewAreaModel(textinput.NewArea("Question?"))
	m.Validate = nil
	m.ColorProfile = termenv.TrueColor

	test.Run(t, m, keyCtrlC)

	if m.Err == nil {
		t.Fatalf("aborting did not produce an error")
	}

	if !errors.Is(m.Err, promptkit.ErrAborted) {
		t.Fatalf("aborting produced %q instead of %q", m.Err, promptkit.ErrAborted)
	}

	test.AssertGoldenView(t, m, "area_abort.golden")
}

func TestAreaSubmit(t *testing.T) {
	t.Parallel()

	m := textinput.NewAreaModel(textinput.NewArea("foo:"))
	m.ResultTemplate = `result: {{ .FinalValue }}`
	m.Validate = nil
	m.ColorProfile = termenv.TrueColor

	test.Run(t, m)
	assertNoAreaError(t, m)

	cmd := test.Update(t, m, keyShiftEnter)
	if cmd == nil || cmd() != tea.Quit() {
		t.Errorf("alt+enter did not produce quit signal")
	}

	test.AssertGoldenView(t, m, "area_submit.golden")
}

func TestAreaValidate(t *testing.T) {
	t.Parallel()

	m := textinput.NewAreaModel(textinput.NewArea("foo:"))
	m.ColorProfile = termenv.TrueColor

	test.Run(t, m)
	assertNoAreaError(t, m)

	cmd := test.Update(t, m, keyShiftEnter)
	if cmd != nil {
		t.Errorf("alt+enter on empty input that does not validate did not produce a no-op")
	}

	test.Update(t, m, test.KeyMsg('x'))

	cmd = test.Update(t, m, keyShiftEnter)
	if cmd == nil || cmd() != tea.Quit() {
		t.Errorf("alt+enter on input that validates did not produce quit signal")
	}
}

func TestAreaTemplate(t *testing.T) {
	t.Parallel()

	separator := "|"

	m := textinput.NewAreaModel(textinput.NewArea("name?"))
	m.Template = `{{ print .Prompt Separator .Input}}`
	m.ResultTemplate = `my text is {{ .FinalValue }}`
	m.ExtendedTemplateFuncs["Separator"] = func() string { return separator }
	m.ColorProfile = termenv.TrueColor

	test.Run(t, m, test.MsgsFromText("hello")...)
	assertNoAreaError(t, m)
	test.AssertGoldenView(t, m, "area_template.golden")

	view := m.View().Content
	if !strings.Contains(test.StripANSI(view), separator) {
		t.Errorf("separator was not rendered:\n%s", test.Indent(view))
	}

	test.Update(t, m, keyShiftEnter)
	test.AssertGoldenView(t, m, "area_template_confirmed.golden")
}

func getAreaValue(tb testing.TB, m *textinput.AreaModel) string {
	tb.Helper()

	v, err := m.Value()
	if err != nil {
		tb.Fatalf("value: %v", err)
	}

	return v
}

func assertNoAreaError(tb testing.TB, m *textinput.AreaModel) {
	tb.Helper()

	if m.Err != nil {
		tb.Fatalf("model contains error: %v", m.Err)
	}
}
