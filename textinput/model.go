package textinput

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/erikgeiser/promptkit"
	"github.com/muesli/termenv"
)

// Model implements the bubbletea.Model for a text input.
type Model struct {
	*TextInput

	// Err holds errors that may occur during the execution of
	// the textinput.
	Err error

	// MaxWidth limits the width of the view using the TextInput's WrapMode.
	MaxWidth int

	input textinput.Model

	tmpl       *template.Template
	resultTmpl *template.Template

	autoCompleteTriggered  bool
	autoCompleteIndecisive bool

	quitting bool

	width int
}

// ensure that the Model interface is implemented.
var _ tea.Model = &Model{}

// NewModel returns a new model based on the provided text input.
func NewModel(textInput *TextInput) *Model {
	return &Model{TextInput: textInput}
}

// Init initializes the text input model.
func (m *Model) Init() tea.Cmd {
	m.tmpl, m.Err = m.initTemplate()
	if m.Err != nil {
		return tea.Quit
	}

	m.resultTmpl, m.Err = m.initResultTemplate()
	if m.Err != nil {
		return tea.Quit
	}

	m.input = m.initInput()

	return textinput.Blink
}

func (m *Model) initTemplate() (*template.Template, error) {
	tmpl := template.New("view")
	tmpl.Funcs(termenv.TemplateFuncs(m.ColorProfile))
	tmpl.Funcs(promptkit.UtilFuncMap())
	tmpl.Funcs(m.ExtendedTemplateFuncs)
	tmpl.Funcs(template.FuncMap{
		"Mask": m.mask,
		"AutoCompleteSuggestions": func() []string {
			return m.AutoComplete(m.input.Value())
		},
	})

	return tmpl.Parse(m.Template)
}

func (m *Model) initResultTemplate() (*template.Template, error) {
	if m.ResultTemplate == "" {
		return nil, nil
	}

	tmpl := template.New("result")
	tmpl.Funcs(termenv.TemplateFuncs(m.ColorProfile))
	tmpl.Funcs(promptkit.UtilFuncMap())
	tmpl.Funcs(m.ExtendedTemplateFuncs)
	tmpl.Funcs(template.FuncMap{
		"Mask": m.mask,
	})

	return tmpl.Parse(m.ResultTemplate)
}

func (m *Model) initInput() textinput.Model {
	input := textinput.New()
	input.Prompt = ""
	input.Placeholder = m.Placeholder

	input.CharLimit = m.CharLimit
	if m.InputWidth > 0 {
		input.SetWidth(m.InputWidth)
	}

	styles := input.Styles()
	styles.Focused.Text = m.InputTextStyle
	styles.Focused.Placeholder = m.InputPlaceholderStyle
	input.SetStyles(styles)

	if m.Hidden {
		input.EchoMode = textinput.EchoPassword
		input.EchoCharacter = m.HideMask
	}

	input.SetValue(m.InitialValue)
	input.Focus()

	return input
}

// Update updates the model based on the received message.
func (m *Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	if m.Err != nil {
		return m, tea.Quit
	}

	var cmd tea.Cmd

	switch msg := message.(type) {
	case tea.KeyPressMsg:
		m.autoCompleteTriggered = false
		m.autoCompleteIndecisive = false

		switch {
		case keyMatches(msg, m.KeyMap.Submit):
			if m.Validate == nil || m.Validate(m.input.Value()) == nil {
				m.quitting = true

				return m, tea.Quit
			}
		case keyMatches(msg, m.KeyMap.AutoComplete):
			if m.AutoComplete != nil {
				m.input.SetValue(m.autoCompleteResult(m.input.Value()))
				m.input.CursorEnd()
			}
		case keyMatches(msg, m.KeyMap.Abort):
			m.Err = promptkit.ErrAborted
			m.quitting = true

			return m, tea.Quit
		case keyMatches(msg, m.KeyMap.Reset):
			m.input.SetValue(m.InitialValue)
			m.input.CursorStart()

			return m, cmd
		case keyMatches(msg, m.KeyMap.Clear):
			m.input.SetValue("")

			return m, cmd
		case keyMatches(msg, m.KeyMap.DeleteAllAfterCursor):
			message = tea.KeyPressMsg{Code: 'k', Mod: tea.ModCtrl}
		case keyMatches(msg, m.KeyMap.DeleteAllBeforeCursor):
			message = tea.KeyPressMsg{Code: 'u', Mod: tea.ModCtrl}
		case keyMatches(msg, m.KeyMap.DeleteWordBeforeCursor):
			message = tea.KeyPressMsg{Code: 'w', Mod: tea.ModCtrl}
		case keyMatches(msg, m.KeyMap.DeleteUnderCursor):
			message = tea.KeyPressMsg{Code: tea.KeyDelete}
		case keyMatches(msg, m.KeyMap.DeleteBeforeCursor):
			message = tea.KeyPressMsg{Code: tea.KeyBackspace}
		case keyMatches(msg, m.KeyMap.MoveBackward):
			message = tea.KeyPressMsg{Code: tea.KeyLeft}
		case keyMatches(msg, m.KeyMap.MoveForward):
			message = tea.KeyPressMsg{Code: tea.KeyRight}
		case keyMatches(msg, m.KeyMap.JumpToBeginning):
			message = tea.KeyPressMsg{Code: tea.KeyHome}
		case keyMatches(msg, m.KeyMap.JumpToEnd):
			message = tea.KeyPressMsg{Code: tea.KeyEnd}
		case keyMatches(msg, m.KeyMap.Paste):
			message = tea.KeyPressMsg{Code: 'v', Mod: tea.ModCtrl}
		case keyMatchesUpstreamKeyMap(msg):
			return m, cmd // do not pass to bubbles/textinput
		default: // do nothing
		}
	case tea.WindowSizeMsg:
		m.width = zeroAwareMin(msg.Width, m.MaxWidth)
	case error:
		m.Err = msg

		return m, tea.Quit
	}

	m.input, cmd = m.input.Update(message)

	return m, cmd
}

// View renders the text input.
func (m *Model) View() tea.View {
	return tea.NewView(m.view())
}

// inputView returns the rendered input, working around a bubbles v2 bug where
// placeholderView() only renders 1 character when Width=0 (because it allocates
// p := make([]rune, Width+1) which is length 1, triggering an early return).
// When no fixed width is configured we render the placeholder ourselves so the
// dynamic (non-scrolling) width is preserved.
func (m *Model) inputView() string {
	if m.InputWidth <= 0 && m.input.Value() == "" && m.Placeholder != "" {
		ph := []rune(m.Placeholder)

		cursor := m.InputPlaceholderStyle.Render(string(ph[:1]))
		if len(ph) > 1 {
			return cursor + m.InputPlaceholderStyle.Render(string(ph[1:]))
		}

		return cursor
	}

	return m.input.View()
}

func (m *Model) view() string {
	if m.quitting {
		view, err := m.resultView()
		if err != nil {
			m.Err = err

			return ""
		}

		return m.wrap(view)
	}

	// avoid panics if Quit is sent during Init
	if m.tmpl == nil {
		return ""
	}

	viewBuffer := &bytes.Buffer{}

	var validationErr error
	if m.Validate != nil {
		validationErr = m.Validate(m.input.Value())
	}

	err := m.tmpl.Execute(viewBuffer, map[string]any{
		"Prompt":                 m.Prompt,
		"InitialValue":           m.InitialValue,
		"Placeholder":            m.Placeholder,
		"Input":                  m.inputView(),
		"ValidationError":        validationErr,
		"TerminalWidth":          m.width,
		"AutoCompleteTriggered":  m.autoCompleteTriggered,
		"AutoCompleteIndecisive": m.autoCompleteIndecisive,
	})
	if err != nil {
		m.Err = err

		return "Template Error: " + err.Error()
	}

	return m.wrap(viewBuffer.String())
}

func (m *Model) resultView() (string, error) {
	viewBuffer := &bytes.Buffer{}

	if m.ResultTemplate == "" {
		return "", nil
	}

	if m.resultTmpl == nil {
		return "", fmt.Errorf("rendering confirmation without loaded template")
	}

	value, err := m.Value()
	if err != nil {
		return "", err
	}

	err = m.resultTmpl.Execute(viewBuffer, map[string]any{
		"FinalValue":    value,
		"Prompt":        m.Prompt,
		"InitialValue":  m.InitialValue,
		"Placeholder":   m.Placeholder,
		"Hidden":        m.Hidden,
		"TerminalWidth": m.width,
	})
	if err != nil {
		return "", fmt.Errorf("execute confirmation template: %w", err)
	}

	return viewBuffer.String(), nil
}

func (m *Model) wrap(text string) string {
	if m.WrapMode == nil {
		return text
	}

	return m.WrapMode(text, m.width)
}

// Value returns the current value and error.
func (m *Model) Value() (string, error) {
	return m.input.Value(), m.Err
}

// mask replaces each character with HideMask if Hidden is true.
func (m *Model) mask(s string) string {
	if !m.Hidden {
		return s
	}

	return strings.Repeat(string(m.HideMask), len(s))
}

func (m *Model) autoCompleteResult(input string) string {
	m.autoCompleteTriggered = true

	if m.AutoComplete == nil {
		return input
	}

	switch candidates := m.AutoComplete(input); len(candidates) {
	case 0:
		return input
	case 1:
		return candidates[0]
	default:
		m.autoCompleteIndecisive = true

		return commonPrefix(candidates)
	}
}

func zeroAwareMin(a int, b int) int {
	switch {
	case a == 0:
		return b
	case b == 0:
		return a
	case a > b:
		return b
	default:
		return a
	}
}
