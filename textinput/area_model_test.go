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
	test.AssertGoldenView(t, m, "area_placeholder_rejected.golden")
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

func TestAreaValidationIndicator(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := textinput.NewAreaModel(textinput.NewArea("foo:"))
		m.ColorProfile = termenv.TrueColor

		test.Run(t, m, test.MsgsFromText("hello")...)
		assertNoAreaError(t, m)

		view := test.StripANSI(m.View().Content)

		if !strings.Contains(view, "✔") {
			t.Errorf("success indicator ✔ not rendered:\n%s", test.Indent(m.View().Content))
		}

		if strings.Contains(view, "✘") {
			t.Errorf("failure indicator ✘ unexpectedly rendered:\n%s", test.Indent(m.View().Content))
		}

		for _, line := range strings.Split(view, "\n") {
			if strings.Contains(line, "✔") && !strings.Contains(line, "Newline:") {
				t.Errorf("success indicator ✔ is not on the same line as the hints:\n%s", test.Indent(m.View().Content))
			}
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()

		m := textinput.NewAreaModel(textinput.NewArea("foo:"))
		m.ColorProfile = termenv.TrueColor

		test.Run(t, m) // no input — ValidateNotEmpty fails
		assertNoAreaError(t, m)

		view := test.StripANSI(m.View().Content)

		if !strings.Contains(view, "✘") {
			t.Errorf("failure indicator ✘ not rendered:\n%s", test.Indent(m.View().Content))
		}

		if strings.Contains(view, "✔") {
			t.Errorf("success indicator ✔ unexpectedly rendered:\n%s", test.Indent(m.View().Content))
		}

		for _, line := range strings.Split(view, "\n") {
			if strings.Contains(line, "✘") && !strings.Contains(line, "Newline:") {
				t.Errorf("failure indicator ✘ is not on the same line as the hints:\n%s", test.Indent(m.View().Content))
			}
		}
	})
}

func TestAreaScrollIndicators(t *testing.T) {
	t.Parallel()

	newScrollModel := func() *textinput.AreaModel {
		m := textinput.NewAreaModel(textinput.NewArea("notes:"))
		m.Height = 3
		m.Validate = nil
		m.ColorProfile = termenv.TrueColor

		return m
	}

	typelines := func(lines ...string) []tea.Msg {
		msgs := test.MsgsFromText(lines[0])
		for _, line := range lines[1:] {
			msgs = append(msgs, keyEnter)
			msgs = append(msgs, test.MsgsFromText(line)...)
		}

		return msgs
	}

	t.Run("up indicator when content is above", func(t *testing.T) {
		t.Parallel()

		m := newScrollModel()
		// 4 lines in a height-3 area — viewport scrolls 1 line down, leaving line 1 above.
		test.Run(t, m, typelines("line 1", "line 2", "line 3", "line 4")...)
		assertNoAreaError(t, m)

		view := test.StripANSI(m.View().Content)
		if !strings.Contains(view, "↑") {
			t.Errorf("up indicator not rendered when content is above:\n%s", test.Indent(m.View().Content))
		}

		if strings.Contains(view, "↓") {
			t.Errorf("down indicator unexpectedly rendered when cursor is on last line:\n%s", test.Indent(m.View().Content))
		}
	})

	t.Run("down indicator when content is below", func(t *testing.T) {
		t.Parallel()

		m := newScrollModel()
		// Type 5 lines, then move cursor back to the top so lines 4-5 are below.
		test.Run(t, m, typelines("line 1", "line 2", "line 3", "line 4", "line 5")...)
		assertNoAreaError(t, m)

		for range 4 {
			test.Update(t, m, tea.KeyPressMsg{Code: tea.KeyUp})
		}

		view := test.StripANSI(m.View().Content)
		if !strings.Contains(view, "↓") {
			t.Errorf("down indicator not rendered when content is below:\n%s", test.Indent(m.View().Content))
		}

		if strings.Contains(view, "↑") {
			t.Errorf("up indicator unexpectedly rendered when cursor is on first line:\n%s", test.Indent(m.View().Content))
		}
	})

	t.Run("no indicators when content fits", func(t *testing.T) {
		t.Parallel()

		m := newScrollModel()
		// Exactly height-many lines — no scrolling needed.
		test.Run(t, m, typelines("line 1", "line 2", "line 3")...)
		assertNoAreaError(t, m)

		view := test.StripANSI(m.View().Content)
		if strings.Contains(view, "↑") {
			t.Errorf("up indicator unexpectedly rendered when all content fits:\n%s", test.Indent(m.View().Content))
		}

		if strings.Contains(view, "↓") {
			t.Errorf("down indicator unexpectedly rendered when all content fits:\n%s", test.Indent(m.View().Content))
		}
	})

	t.Run("no indicators in auto-resize mode", func(t *testing.T) {
		t.Parallel()

		m := textinput.NewAreaModel(textinput.NewArea("notes:"))
		// Height = 0 (default) — auto-resize, no scroll indicators expected.
		m.Validate = nil
		m.ColorProfile = termenv.TrueColor

		test.Run(t, m, typelines("line 1", "line 2", "line 3", "line 4", "line 5")...)
		assertNoAreaError(t, m)

		view := test.StripANSI(m.View().Content)
		if strings.Contains(view, "↑") || strings.Contains(view, "↓") {
			t.Errorf("scroll indicators unexpectedly rendered in auto-resize mode:\n%s", test.Indent(m.View().Content))
		}
	})
}

func TestAreaCharLimit(t *testing.T) {
	t.Parallel()

	m := textinput.NewAreaModel(textinput.NewArea("input:"))
	m.CharLimit = 5
	m.Validate = nil
	m.ColorProfile = termenv.TrueColor

	test.Run(t, m, test.MsgsFromText("abcdefgh")...)
	assertNoAreaError(t, m)

	value := getAreaValue(t, m)
	if len(value) > 5 {
		t.Errorf("char limit not respected: got %q (len %d), expected max 5", value, len(value))
	}
}

func TestAreaFixedHeight(t *testing.T) {
	t.Parallel()

	m := textinput.NewAreaModel(textinput.NewArea("notes:"))
	m.Height = 3
	m.Validate = nil
	m.ColorProfile = termenv.TrueColor

	// Enter multiple lines to exceed the fixed height.
	msgs := test.MsgsFromText("line 1")
	msgs = append(msgs, keyEnter)
	msgs = append(msgs, test.MsgsFromText("line 2")...)
	msgs = append(msgs, keyEnter)
	msgs = append(msgs, test.MsgsFromText("line 3")...)
	msgs = append(msgs, keyEnter)
	msgs = append(msgs, test.MsgsFromText("line 4")...)
	test.Run(t, m, msgs...)
	assertNoAreaError(t, m)

	value := getAreaValue(t, m)

	lines := strings.Split(value, "\n")
	if len(lines) != 4 {
		t.Errorf("expected 4 lines, got %d: %q", len(lines), value)
	}
}

func TestAreaShowLineNumbers(t *testing.T) {
	t.Parallel()

	m := textinput.NewAreaModel(textinput.NewArea("code:"))
	m.ShowLineNumbers = true
	m.Validate = nil
	m.ColorProfile = termenv.TrueColor

	msgs := test.MsgsFromText("hello")
	msgs = append(msgs, keyEnter)
	msgs = append(msgs, test.MsgsFromText("world")...)
	test.Run(t, m, msgs...)
	assertNoAreaError(t, m)

	view := m.View().Content
	strippedView := test.StripANSI(view)

	// Line numbers should show "1" somewhere in the rendered view.
	if !strings.Contains(strippedView, "1") {
		t.Errorf("line numbers not visible in view:\n%s", test.Indent(view))
	}
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
