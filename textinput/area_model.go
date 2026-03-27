package textinput

import (
	"bytes"
	"fmt"
	"text/template"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"github.com/erikgeiser/promptkit"
	"github.com/muesli/termenv"
)

// AreaModel implements the bubbletea.Model for a text area.
type AreaModel struct {
	*TextArea

	// Err holds errors that may occur during the execution of the text area.
	Err error

	// MaxWidth limits the width of the view using the TextArea's WrapMode.
	MaxWidth int

	input textarea.Model

	tmpl       *template.Template
	resultTmpl *template.Template

	quitting bool

	width int
}

// ensure that the Model interface is implemented.
var _ tea.Model = &AreaModel{}

// NewAreaModel returns a new model based on the provided text area.
func NewAreaModel(textArea *TextArea) *AreaModel {
	return &AreaModel{TextArea: textArea}
}

// Init initializes the text area model.
func (m *AreaModel) Init() tea.Cmd {
	m.tmpl, m.Err = m.initAreaTemplate()
	if m.Err != nil {
		return tea.Quit
	}

	m.resultTmpl, m.Err = m.initAreaResultTemplate()
	if m.Err != nil {
		return tea.Quit
	}

	m.input = m.initAreaInput()

	return tea.Batch(textarea.Blink, m.input.Focus())
}

func (m *AreaModel) initAreaTemplate() (*template.Template, error) {
	tmpl := template.New("view")
	tmpl.Funcs(termenv.TemplateFuncs(m.ColorProfile))
	tmpl.Funcs(promptkit.UtilFuncMap())
	tmpl.Funcs(m.ExtendedTemplateFuncs)

	return tmpl.Parse(m.Template)
}

func (m *AreaModel) initAreaResultTemplate() (*template.Template, error) {
	if m.ResultTemplate == "" {
		return nil, nil
	}

	tmpl := template.New("result")
	tmpl.Funcs(termenv.TemplateFuncs(m.ColorProfile))
	tmpl.Funcs(promptkit.UtilFuncMap())
	tmpl.Funcs(m.ExtendedTemplateFuncs)

	return tmpl.Parse(m.ResultTemplate)
}

func (m *AreaModel) initAreaInput() textarea.Model {
	input := textarea.New()
	input.Placeholder = m.Placeholder
	input.CharLimit = m.CharLimit
	input.ShowLineNumbers = m.ShowLineNumbers

	height := m.Height
	if height <= 0 {
		height = 1
	}

	input.SetHeight(height)

	styles := input.Styles()
	styles.Focused.Text = m.InputTextStyle
	styles.Focused.Placeholder = m.InputPlaceholderStyle
	input.SetStyles(styles)

	if m.InitialValue != "" {
		input.SetValue(m.InitialValue)
	}

	return input
}

// Update updates the model based on the received message.
func (m *AreaModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	if m.Err != nil {
		return m, tea.Quit
	}

	var cmd tea.Cmd

	switch msg := message.(type) {
	case tea.KeyPressMsg:
		switch {
		case areaKeyMatches(msg, m.KeyMap.Submit):
			if m.Validate == nil || m.Validate(m.input.Value()) == nil {
				m.quitting = true

				return m, tea.Quit
			}
		case areaKeyMatches(msg, m.KeyMap.Abort):
			m.Err = promptkit.ErrAborted
			m.quitting = true

			return m, tea.Quit
		case areaKeyMatches(msg, m.KeyMap.InsertNewline):
			message = tea.KeyPressMsg{Code: tea.KeyEnter}
		case areaKeyMatches(msg, m.KeyMap.MoveBackward):
			message = tea.KeyPressMsg{Code: tea.KeyLeft}
		case areaKeyMatches(msg, m.KeyMap.MoveForward):
			message = tea.KeyPressMsg{Code: tea.KeyRight}
		case areaKeyMatches(msg, m.KeyMap.MoveWordBackward):
			message = tea.KeyPressMsg{Code: tea.KeyLeft, Mod: tea.ModAlt}
		case areaKeyMatches(msg, m.KeyMap.MoveWordForward):
			message = tea.KeyPressMsg{Code: tea.KeyRight, Mod: tea.ModAlt}
		case areaKeyMatches(msg, m.KeyMap.MoveUp):
			message = tea.KeyPressMsg{Code: tea.KeyUp}
		case areaKeyMatches(msg, m.KeyMap.MoveDown):
			message = tea.KeyPressMsg{Code: tea.KeyDown}
		case areaKeyMatches(msg, m.KeyMap.JumpToLineBeginning):
			message = tea.KeyPressMsg{Code: tea.KeyHome}
		case areaKeyMatches(msg, m.KeyMap.JumpToLineEnd):
			message = tea.KeyPressMsg{Code: tea.KeyEnd}
		case areaKeyMatches(msg, m.KeyMap.JumpToBeginning):
			message = tea.KeyPressMsg{Code: tea.KeyHome, Mod: tea.ModCtrl}
		case areaKeyMatches(msg, m.KeyMap.JumpToEnd):
			message = tea.KeyPressMsg{Code: tea.KeyEnd, Mod: tea.ModCtrl}
		case areaKeyMatches(msg, m.KeyMap.DeleteBeforeCursor):
			message = tea.KeyPressMsg{Code: tea.KeyBackspace}
		case areaKeyMatches(msg, m.KeyMap.DeleteWordBeforeCursor):
			message = tea.KeyPressMsg{Code: 'w', Mod: tea.ModCtrl}
		case areaKeyMatches(msg, m.KeyMap.DeleteUnderCursor):
			message = tea.KeyPressMsg{Code: tea.KeyDelete}
		case areaKeyMatches(msg, m.KeyMap.DeleteWordAfterCursor):
			message = tea.KeyPressMsg{Code: 'd', Mod: tea.ModAlt}
		case areaKeyMatches(msg, m.KeyMap.DeleteAllAfterCursor):
			message = tea.KeyPressMsg{Code: 'k', Mod: tea.ModCtrl}
		case areaKeyMatches(msg, m.KeyMap.DeleteAllBeforeCursor):
			message = tea.KeyPressMsg{Code: 'u', Mod: tea.ModCtrl}
		case areaKeyMatches(msg, m.KeyMap.Paste):
			message = tea.KeyPressMsg{Code: 'v', Mod: tea.ModCtrl}
		case areaKeyMatchesUpstream(msg):
			return m, cmd // do not pass to bubbles/textarea
		default: // do nothing
		}

		m.prepareAutoResize()
		m.input, cmd = m.input.Update(message)
		m.autoResizeInput()

		return m, cmd
	case tea.WindowSizeMsg:
		m.width = zeroAwareMin(msg.Width, m.MaxWidth)
		m.input.SetWidth(m.width)
	case error:
		m.Err = msg

		return m, tea.Quit
	}

	m.prepareAutoResize()
	m.input, cmd = m.input.Update(message)
	m.autoResizeInput()

	return m, cmd
}

func (m *AreaModel) prepareAutoResize() {
	if m.Height <= 0 {
		// Set a large height before Update() so the viewport never scrolls,
		// preventing lines from disappearing above the visible area.
		m.input.SetHeight(9999)
	}
}

func (m *AreaModel) autoResizeInput() {
	if m.Height <= 0 {
		m.input.SetHeight(m.input.LineCount())
	}
}

// View renders the text area.
func (m *AreaModel) View() tea.View {
	return tea.NewView(m.areaView())
}

func (m *AreaModel) areaView() string {
	if m.quitting {
		view, err := m.areaResultView()
		if err != nil {
			m.Err = err

			return ""
		}

		return m.areaWrap(view)
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
		"Prompt":          m.Prompt,
		"InitialValue":    m.InitialValue,
		"Placeholder":     m.Placeholder,
		"Input":           m.input.View(),
		"ValidationError": validationErr,
		"TerminalWidth":   m.width,
		"KeyMap":          m.KeyMap,
	})
	if err != nil {
		m.Err = err

		return "Template Error: " + err.Error()
	}

	return m.areaWrap(viewBuffer.String())
}

func (m *AreaModel) areaResultView() (string, error) {
	if m.ResultTemplate == "" {
		return "", nil
	}

	if m.resultTmpl == nil {
		return "", fmt.Errorf("rendering result without loaded template")
	}

	value, err := m.Value()
	if err != nil {
		return "", err
	}

	viewBuffer := &bytes.Buffer{}

	err = m.resultTmpl.Execute(viewBuffer, map[string]any{
		"FinalValue":    value,
		"Prompt":        m.Prompt,
		"InitialValue":  m.InitialValue,
		"Placeholder":   m.Placeholder,
		"TerminalWidth": m.width,
	})
	if err != nil {
		return "", fmt.Errorf("execute result template: %w", err)
	}

	return viewBuffer.String(), nil
}

func (m *AreaModel) areaWrap(text string) string {
	if m.WrapMode == nil {
		return text
	}

	return m.WrapMode(text, m.width)
}

// Value returns the current value and error.
func (m *AreaModel) Value() (string, error) {
	return m.input.Value(), m.Err
}
