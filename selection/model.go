package selection

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
	"github.com/muesli/termenv"
)

// Model implements the bubbletea.Model for a selection prompt.
type Model struct {
	*Selection

	// Err holds errors that may occur during the execution of
	// the selection prompt.
	Err error

	filterInput textinput.Model
	// currently displayed choices, after filtering and pagination
	currentChoices []*Choice
	// number of available choices after filtering
	availableChoices int
	// index of current selection in currentChoices slice
	currentIdx   int
	scrollOffset int
	width        int
	tmpl         *template.Template

	quitting bool
}

// ensure that the Model interface is implemented.
var _ tea.Model = &Model{}

// NewModel returns a new selection prompt model for the
// provided choices.
func NewModel(selection *Selection) *Model {
	return &Model{Selection: selection}
}

// Init initializes the selection prompt model.
func (m *Model) Init() tea.Cmd {
	m.reindexChoices()

	m.Err = m.validate()
	if m.Err != nil {
		return tea.Quit
	}

	m.tmpl, m.Err = m.initTemplate()
	if m.Err != nil {
		return tea.Quit
	}

	m.filterInput = m.initFilterInput()

	m.currentChoices, m.availableChoices = m.filteredAndPagedChoices()

	return textinput.Blink
}

func (m *Model) initTemplate() (*template.Template, error) {
	tmpl := template.New("")
	tmpl.Funcs(termenv.TemplateFuncs(termenv.ColorProfile()))
	tmpl.Funcs(m.ExtendedTemplateScope)
	tmpl.Funcs(template.FuncMap{
		"IsScrollDownHintPosition": func(idx int) bool {
			return m.canScrollDown() && (idx == len(m.currentChoices)-1)
		},
		"IsScrollUpHintPosition": func(idx int) bool {
			return m.canScrollUp() && idx == 0 && m.scrollOffset > 0
		},
	})

	return tmpl.Parse(m.Template)
}

func (m *Model) initFilterInput() textinput.Model {
	filterInput := textinput.NewModel()
	filterInput.Prompt = ""
	filterInput.TextStyle = m.FilterInputTextStyle
	filterInput.BackgroundStyle = m.FilterInputBackgroundStyle
	filterInput.PlaceholderStyle = m.FilterInputPlaceholderStyle
	filterInput.CursorStyle = m.FilterInputCursorStyle
	filterInput.Placeholder = m.FilterPlaceholder
	filterInput.Width = 80
	filterInput.Focus()

	return filterInput
}

// Choice returns the choice that is currently selected or the final
// choice after the prompt has concluded.
func (m *Model) Choice() (*Choice, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	if len(m.currentChoices) == 0 {
		return nil, fmt.Errorf("no choices")
	}

	if m.currentIdx < 0 || m.currentIdx >= len(m.currentChoices) {
		return nil, fmt.Errorf("choice index out of bounds")
	}

	return m.currentChoices[m.currentIdx], nil
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
		case keyMatches(msg, m.KeyMap.Abort):
			m.Err = fmt.Errorf("selection was aborted")

			return m, tea.Quit
		case keyMatches(msg, m.KeyMap.ClearFilter):
			m.filterInput.Reset()
			m.currentChoices, m.availableChoices = m.filteredAndPagedChoices()

			return m, nil
		case keyMatches(msg, m.KeyMap.Select):
			if len(m.currentChoices) == 0 {
				return m, nil
			}

			m.quitting = true

			return m, tea.Quit
		case keyMatches(msg, m.KeyMap.Down):
			m.cursorDown()

			return m, nil
		case keyMatches(msg, m.KeyMap.Up):
			m.cursorUp()

			return m, nil
		case keyMatches(msg, m.KeyMap.ScrollDown):
			m.scrollDown()

			return m, nil
		case keyMatches(msg, m.KeyMap.ScrollUp):
			m.scrollUp()

			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case error:
		m.Err = msg

		return m, tea.Quit
	}

	if m.Filter == nil {
		return m, cmd
	}

	previousFilter := m.filterInput.Value()

	m.filterInput, cmd = m.filterInput.Update(msg)

	if m.filterInput.Value() != previousFilter {
		m.currentIdx = 0
		m.scrollOffset = 0
		m.currentChoices, m.availableChoices = m.filteredAndPagedChoices()
	}

	return m, cmd
}

// View renders the selection prompt.
func (m *Model) View() string {
	if m.quitting {
		return ""
	}

	// avoid panics if Quit is sent during Init
	if m.tmpl == nil {
		m.Err = fmt.Errorf("rendering view without loaded template")

		return ""
	}

	viewBuffer := &bytes.Buffer{}

	err := m.tmpl.Execute(viewBuffer, map[string]interface{}{
		"Prompt":        m.Prompt,
		"IsFiltered":    m.Filter != nil,
		"FilterInput":   m.filterInput.View(),
		"Choices":       m.currentChoices,
		"NChoices":      len(m.currentChoices),
		"SelectedIndex": m.currentIdx,
		"PageSize":      m.PageSize,
		"IsPaged":       m.PageSize > 0 && len(m.currentChoices) > m.PageSize,
		"AllChoices":    m.Choices,
		"NAllChoices":   len(m.Choices),
	})
	if err != nil {
		m.Err = err

		return "Template Error: " + err.Error()
	}

	termenv.Reset()

	return wrap.String(wordwrap.String(viewBuffer.String(), m.width), m.width)
}

func (m Model) filteredAndPagedChoices() ([]*Choice, int) {
	choices := []*Choice{}

	var available, ignored int

	for _, choice := range m.Choices {
		if m.Filter != nil && !m.Filter(m.filterInput.Value(), choice) {
			continue
		}

		available++

		if m.PageSize > 0 && len(choices) >= m.PageSize {
			break
		}

		if (m.PageSize > 0) && (ignored < m.scrollOffset) {
			ignored++

			continue
		}

		choices = append(choices, choice)
	}

	return choices, available
}

func (m *Model) canScrollDown() bool {
	if m.PageSize <= 0 || m.availableChoices <= m.PageSize {
		return false
	}

	if m.scrollOffset+m.PageSize >= len(m.Choices) {
		return false
	}

	return true
}

func (m *Model) canScrollUp() bool {
	return m.scrollOffset > 0
}

func (m *Model) cursorDown() {
	if m.currentIdx == len(m.currentChoices)-1 && m.canScrollDown() {
		m.scrollDown()
	}

	m.currentIdx = min(len(m.currentChoices)-1, m.currentIdx+1)
}

func (m *Model) cursorUp() {
	if m.currentIdx == 0 && m.canScrollUp() {
		m.scrollUp()
	}

	m.currentIdx = max(0, m.currentIdx-1)
}

func (m *Model) scrollDown() {
	if m.PageSize <= 0 || m.scrollOffset+m.PageSize >= m.availableChoices {
		return
	}

	m.currentIdx = max(0, m.currentIdx-1)
	m.scrollOffset++
	m.currentChoices, m.availableChoices = m.filteredAndPagedChoices()
}

func (m *Model) scrollUp() {
	if m.PageSize <= 0 || m.scrollOffset <= 0 {
		return
	}

	m.currentIdx = min(len(m.currentChoices)-1, m.currentIdx+1)
	m.scrollOffset--
	m.currentChoices, m.availableChoices = m.filteredAndPagedChoices()
}

func (m *Model) reindexChoices() {
	for i, choice := range m.Choices {
		choice.Index = i
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
