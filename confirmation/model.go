package confirmation

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/erikgeiser/promptkit"
	"github.com/muesli/termenv"
)

// Model implements the bubbletea.Model for a confirmation prompt.
type Model struct {
	*Confirmation

	Err error

	tmpl             *template.Template
	confirmationTmpl *template.Template

	value Value

	quitting bool

	width int
}

// ensure that the Model interface is implemented.
var _ tea.Model = &Model{}

// NewModel returns a new model based on the provided confirmation prompt.
func NewModel(confirmation *Confirmation) *Model {
	return &Model{
		Confirmation: confirmation,
		value:        confirmation.DefaultValue,
	}
}

// Init initializes the confirmation prompt model.
func (m *Model) Init() tea.Cmd {
	m.tmpl, m.Err = m.initTemplate()
	if m.Err != nil {
		return tea.Quit
	}

	m.confirmationTmpl, m.Err = m.initConfirmationTemplate()
	if m.Err != nil {
		return tea.Quit
	}

	termenv.Reset()

	return textinput.Blink
}

func (m *Model) initTemplate() (*template.Template, error) {
	tmpl := template.New("")
	tmpl.Funcs(termenv.TemplateFuncs(termenv.ColorProfile()))
	tmpl.Funcs(promptkit.UtilFuncMap())
	tmpl.Funcs(m.ExtendedTemplateScope)

	return tmpl.Parse(m.Template)
}

func (m *Model) initConfirmationTemplate() (*template.Template, error) {
	if m.ConfirmationTemplate == "" {
		return nil, nil
	}

	tmpl := template.New("confirmed")
	tmpl.Funcs(termenv.TemplateFuncs(termenv.ColorProfile()))
	tmpl.Funcs(promptkit.UtilFuncMap())
	tmpl.Funcs(m.ExtendedTemplateScope)

	return tmpl.Parse(m.ConfirmationTemplate)
}

// Update updates the model based on the received message.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.Err != nil {
		return m, tea.Quit
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case keyMatches(msg, m.KeyMap.Submit):
			if m.value != Undecided {
				m.quitting = true

				return m, tea.Quit
			}
		case keyMatches(msg, m.KeyMap.Abort):
			m.Err = promptkit.ErrAborted
			m.quitting = true

			return m, tea.Quit
		case keyMatches(msg, m.KeyMap.Yes):
			m.value = Yes
			m.quitting = true

			return m, tea.Quit
		case keyMatches(msg, m.KeyMap.No):
			m.value = No
			m.quitting = true

			return m, tea.Quit
		case keyMatches(msg, m.KeyMap.SelectYes):
			m.value = Yes
		case keyMatches(msg, m.KeyMap.SelectNo):
			m.value = No
		case keyMatches(msg, m.KeyMap.Toggle):
			switch m.value {
			case Yes:
				m.value = No
			case No, Undecided:
				m.value = Yes
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case error:
		m.Err = msg

		return m, tea.Quit
	}

	return m, cmd
}

// View renders the confirmation prompt.
func (m *Model) View() string {
	defer termenv.Reset()

	// avoid panics if Quit is sent during Init
	if m.quitting {
		view, err := m.confirmationView()
		if err != nil {
			m.Err = err

			return ""
		}

		return promptkit.Wrap(view, m.width)
	}

	// avoid panics if Quit is sent during Init
	if m.tmpl == nil {
		return ""
	}

	viewBuffer := &bytes.Buffer{}

	err := m.tmpl.Execute(viewBuffer, map[string]interface{}{
		"Prompt":           m.Prompt,
		"YesSelected":      m.value == Yes,
		"NoSelected":       m.value == No,
		"Undecided":        m.value == Undecided,
		"DefaultYes":       m.DefaultValue == Yes,
		"DefaultNo":        m.DefaultValue == No,
		"DefaultUndecided": m.DefaultValue == Undecided,
		"TerminalWidth":    m.width,
	})
	if err != nil {
		m.Err = err

		return "Template Error: " + err.Error()
	}

	return promptkit.Wrap(viewBuffer.String(), m.width)
}

func (m *Model) confirmationView() (string, error) {
	viewBuffer := &bytes.Buffer{}

	if m.ConfirmationTemplate == "" {
		return "", nil
	}

	if m.confirmationTmpl == nil {
		return "", fmt.Errorf("rendering confirmation without loaded template")
	}

	value, err := m.Value()
	if err != nil {
		return "", fmt.Errorf("obtaining value for confirmation: %w", err)
	}

	err = m.confirmationTmpl.Execute(viewBuffer, map[string]interface{}{
		"FinalValue":       value,
		"FinalValueString": fmt.Sprintf("%v", value),
		"Prompt":           m.Prompt,
		"DefaultYes":       m.DefaultValue == Yes,
		"DefaultNo":        m.DefaultValue == No,
		"DefaultUndecided": m.DefaultValue == Undecided,
		"TerminalWidth":    m.width,
	})
	if err != nil {
		return "", fmt.Errorf("execute confirmation template: %w", err)
	}

	return viewBuffer.String(), nil
}

// Value returns the current value and error.
func (m *Model) Value() (bool, error) {
	if m.Err != nil {
		return false, m.Err
	}

	if m.value == Undecided {
		return false, fmt.Errorf("no decision was made")
	}

	return *m.value, m.Err
}
