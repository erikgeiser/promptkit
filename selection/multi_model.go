package selection

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/erikgeiser/promptkit"
	"github.com/muesli/termenv"
	"golang.org/x/term"
)

// MultiModel implements the bubbletea.Model for a multi-selection prompt.
type MultiModel[T any] struct {
	*MultiSelection[T]

	// Err holds errors that may occur during the execution of the
	// multi-selection prompt.
	Err error

	// MaxWidth limits the width of the view using the MultiSelection's WrapMode.
	MaxWidth int

	filterInput textinput.Model
	// currently displayed choices, after filtering and pagination
	currentChoices []*Choice[T]
	// number of available choices after filtering
	availableChoices int
	// index of cursor in currentChoices slice
	currentIdx        int
	scrollOffset      int
	width             int
	height            int
	tmpl              *template.Template
	resultTmpl        *template.Template
	requestedPageSize int

	// selectedChoices tracks marked choices by their stable choice index.
	selectedChoices map[int]bool
	quitting        bool
}

// ensure that the Model interface is implemented.
var _ tea.Model = &MultiModel[any]{}

// NewMultiModel returns a new multi-selection prompt model for the provided
// choices.
func NewMultiModel[T any](ms *MultiSelection[T]) *MultiModel[T] {
	return &MultiModel[T]{MultiSelection: ms}
}

// Init initializes the multi-selection prompt model.
func (m *MultiModel[T]) Init() tea.Cmd {
	m.reindexChoices()

	if len(m.choices) == 0 {
		m.Err = fmt.Errorf("no choices provided")

		return tea.Quit
	}

	if m.Template == "" {
		m.Err = fmt.Errorf("empty template")

		return tea.Quit
	}

	m.tmpl, m.Err = m.initTemplate()
	if m.Err != nil {
		return tea.Quit
	}

	m.resultTmpl, m.Err = m.initResultTemplate()
	if m.Err != nil {
		return tea.Quit
	}

	m.filterInput = m.initFilterInput()
	m.selectedChoices = make(map[int]bool)
	m.currentChoices, m.availableChoices = m.filteredAndPagedChoices()
	m.requestedPageSize = m.PageSize

	outputFile, ok := m.Output.(*os.File)
	if ok {
		width, height, err := term.GetSize(int(outputFile.Fd()))
		if err == nil {
			m.resize(width, height)
		}
	}

	return textinput.Blink
}

func (m *MultiModel[T]) initTemplate() (*template.Template, error) {
	tmpl := template.New("view")
	tmpl.Funcs(termenv.TemplateFuncs(m.ColorProfile))
	tmpl.Funcs(m.ExtendedTemplateFuncs)
	tmpl.Funcs(promptkit.UtilFuncMap())
	tmpl.Funcs(template.FuncMap{
		"Faint": func(s string) string {
			return termenv.String(s).Faint().String()
		},
		"IsScrollDownHintPosition": func(idx int) bool {
			return m.canScrollDown() && (idx == len(m.currentChoices)-1)
		},
		"IsScrollUpHintPosition": func(idx int) bool {
			return m.canScrollUp() && idx == 0 && m.scrollOffset > 0
		},
		"IsMarked": func(c *Choice[T]) bool {
			return m.selectedChoices[c.idx]
		},
		"Cursor": func(c *Choice[T]) string {
			if m.CursorChoiceStyle == nil {
				return c.String
			}

			return m.CursorChoiceStyle(c)
		},
		"CursorMarked": func(c *Choice[T]) string {
			if m.CursorMarkedChoiceStyle == nil {
				return c.String
			}

			return m.CursorMarkedChoiceStyle(c)
		},
		"Marked": func(c *Choice[T]) string {
			if m.MarkedChoiceStyle == nil {
				return c.String
			}

			return m.MarkedChoiceStyle(c)
		},
		"Unmarked": func(c *Choice[T]) string {
			if m.UnmarkedChoiceStyle == nil {
				return c.String
			}

			return m.UnmarkedChoiceStyle(c)
		},
	})

	return tmpl.Parse(m.Template)
}

func (m *MultiModel[T]) initResultTemplate() (*template.Template, error) {
	if m.ResultTemplate == "" {
		return nil, nil //nolint:nilnil
	}

	tmpl := template.New("result")
	tmpl.Funcs(termenv.TemplateFuncs(m.ColorProfile))
	tmpl.Funcs(m.ExtendedTemplateFuncs)
	tmpl.Funcs(promptkit.UtilFuncMap())
	tmpl.Funcs(template.FuncMap{
		"Final": func(c *Choice[T]) string {
			if m.FinalChoiceStyle == nil {
				return c.String
			}

			return m.FinalChoiceStyle(c)
		},
	})

	return tmpl.Parse(m.ResultTemplate)
}

func (m *MultiModel[T]) initFilterInput() textinput.Model {
	filterInput := textinput.New()
	filterInput.Prompt = ""
	styles := filterInput.Styles()
	styles.Focused.Text = m.FilterInputTextStyle
	styles.Focused.Placeholder = m.FilterInputPlaceholderStyle
	filterInput.SetStyles(styles)
	filterInput.Placeholder = m.FilterPlaceholder
	filterInput.SetWidth(80)
	filterInput.Focus()

	return filterInput
}

// ValuesAsChoices returns the marked choices in original choice order.
func (m *MultiModel[T]) ValuesAsChoices() ([]*Choice[T], error) {
	if m.Err != nil {
		return nil, m.Err
	}

	choices := make([]*Choice[T], 0, len(m.selectedChoices))

	for _, c := range m.choices {
		if m.selectedChoices[c.idx] {
			choices = append(choices, c)
		}
	}

	return choices, nil
}

// Values returns the marked values in original choice order.
func (m *MultiModel[T]) Values() ([]T, error) {
	choices, err := m.ValuesAsChoices()
	if err != nil {
		return nil, err
	}

	values := make([]T, len(choices))

	for i, c := range choices {
		values[i] = c.Value
	}

	return values, nil
}

// Update updates the model based on the received message.
func (m *MultiModel[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.Err != nil {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case keyMatches(msg, m.KeyMap.Abort):
			m.Err = promptkit.ErrAborted
			m.quitting = true

			return m, tea.Quit
		case keyMatches(msg, m.KeyMap.Select):
			if len(m.selectedChoices) < m.MinSelections {
				return m, nil
			}

			if m.MaxSelections > 0 && len(m.selectedChoices) > m.MaxSelections {
				return m, nil
			}

			m.quitting = true

			return m, tea.Quit
		case keyMatches(msg, m.KeyMap.Toggle):
			m.toggle()
		case keyMatches(msg, m.KeyMap.ClearFilter):
			m.filterInput.Reset()
			m.currentChoices, m.availableChoices = m.filteredAndPagedChoices()
		case keyMatches(msg, m.KeyMap.Down):
			m.cursorDown()
		case keyMatches(msg, m.KeyMap.Up):
			m.cursorUp()
		case keyMatches(msg, m.KeyMap.ScrollDown):
			m.scrollDown()
		case keyMatches(msg, m.KeyMap.ScrollUp):
			m.scrollUp()
		default:
			return m.updateFilter(msg)
		}

		return m, nil
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)

		return m, nil
	case error:
		m.Err = msg

		return m, tea.Quit
	}

	var cmd tea.Cmd

	return m, cmd
}

func (m *MultiModel[T]) toggle() {
	if len(m.currentChoices) == 0 {
		return
	}

	choice := m.currentChoices[m.currentIdx]

	if m.selectedChoices[choice.idx] {
		delete(m.selectedChoices, choice.idx)

		return
	}

	if m.MaxSelections > 0 && len(m.selectedChoices) >= m.MaxSelections {
		return
	}

	m.selectedChoices[choice.idx] = true
}

func (m *MultiModel[T]) resize(width int, height int) {
	m.width = zeroAwareMin(width, m.MaxWidth)

	if m.height != height {
		m.height = height
		m.forceUpdatePageSizeForHeight()
	}
}

func (m *MultiModel[T]) forceUpdatePageSizeForHeight() { //nolint:dupl
	maxAcceptablePageSize := len(m.choices)
	if m.requestedPageSize != 0 {
		maxAcceptablePageSize = min(len(m.choices), m.requestedPageSize)
	}

	m.PageSize = maxAcceptablePageSize
	m.currentIdx = 0
	m.scrollOffset = 0
	m.currentChoices, m.availableChoices = m.filteredAndPagedChoices()

	if lipgloss.Height(m.view()) < m.height {
		return
	}

	for m.PageSize = 1; m.PageSize <= maxAcceptablePageSize; m.PageSize++ {
		m.currentChoices, m.availableChoices = m.filteredAndPagedChoices()

		if lipgloss.Height(m.view()) >= m.height {
			m.PageSize--
			m.currentChoices, m.availableChoices = m.filteredAndPagedChoices()

			return
		}
	}

	m.PageSize--
	m.currentChoices, m.availableChoices = m.filteredAndPagedChoices()
}

func (m *MultiModel[T]) updateFilter(msg tea.Msg) (*MultiModel[T], tea.Cmd) {
	if m.Filter == nil {
		return m, nil
	}

	previousFilter := m.filterInput.Value()

	var cmd tea.Cmd

	m.filterInput, cmd = m.filterInput.Update(msg)

	if m.filterInput.Value() != previousFilter {
		m.currentIdx = 0
		m.scrollOffset = 0
		m.currentChoices, m.availableChoices = m.filteredAndPagedChoices()
	}

	return m, cmd
}

// View renders the multi-selection prompt.
func (m *MultiModel[T]) View() tea.View {
	return tea.NewView(m.view())
}

func (m *MultiModel[T]) view() string {
	viewBuffer := &bytes.Buffer{}

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

	markedChoices := m.markedChoicesInOrder()
	offscreenChoices := m.offscreenSelectedChoices()

	err := m.tmpl.Execute(viewBuffer, map[string]any{
		"Prompt":                  m.Prompt,
		"IsFiltered":              m.Filter != nil,
		"FilterPrompt":            m.FilterPrompt,
		"FilterInput":             m.filterInput.View(),
		"Choices":                 m.currentChoices,
		"NChoices":                len(m.currentChoices),
		"SelectedIndex":           m.currentIdx,
		"PageSize":                m.PageSize,
		"IsPaged":                 m.PageSize > 0 && len(m.currentChoices) > m.PageSize,
		"AllChoices":              m.choices,
		"NAllChoices":             len(m.choices),
		"MarkedChoices":           markedChoices,
		"NMarkedChoices":          len(markedChoices),
		"OffscreenMarkedChoices":  offscreenChoices,
		"NOffscreenMarkedChoices": len(offscreenChoices),
		"TerminalWidth":           m.width,
		"MinSelections":           m.MinSelections,
		"MaxSelections":           m.MaxSelections,
	})
	if err != nil {
		m.Err = err

		return "Template Error: " + err.Error()
	}

	return m.wrap(viewBuffer.String())
}

func (m *MultiModel[T]) resultView() (string, error) {
	viewBuffer := &bytes.Buffer{}

	if m.ResultTemplate == "" {
		return "", nil
	}

	if m.resultTmpl == nil {
		return "", fmt.Errorf("rendering result without loaded template")
	}

	markedChoices := m.markedChoicesInOrder()

	err := m.resultTmpl.Execute(viewBuffer, map[string]any{
		"FinalChoices":    markedChoices,
		"FinalChoicesStr": m.finalChoicesStr(markedChoices),
		"Prompt":          m.Prompt,
		"AllChoices":      m.choices,
		"NAllChoices":     len(m.choices),
		"TerminalWidth":   m.width,
	})
	if err != nil {
		return "", fmt.Errorf("execute result template: %w", err)
	}

	return viewBuffer.String(), nil
}

func (m *MultiModel[T]) wrap(text string) string {
	if m.WrapMode == nil {
		return text
	}

	return m.WrapMode(text, m.width)
}

func (m *MultiModel[T]) filteredAndPagedChoices() ([]*Choice[T], int) {
	choices := []*Choice[T]{}

	var available, ignored int

	for _, choice := range m.choices {
		if m.Filter != nil && !m.Filter(m.filterInput.Value(), choice) {
			continue
		}

		available++

		if m.PageSize > 0 && (len(choices) >= m.PageSize || ignored < m.scrollOffset) {
			ignored++

			continue
		}

		choices = append(choices, choice)
	}

	return choices, available
}

func (m *MultiModel[T]) offscreenSelectedChoices() []*Choice[T] {
	onScreen := make(map[int]bool, len(m.currentChoices))
	for _, c := range m.currentChoices {
		onScreen[c.idx] = true
	}

	var result []*Choice[T]

	for _, choice := range m.choices {
		if m.selectedChoices[choice.idx] && !onScreen[choice.idx] {
			result = append(result, choice)
		}
	}

	return result
}

func (m *MultiModel[T]) canScrollDown() bool {
	if m.PageSize <= 0 || m.availableChoices <= m.PageSize {
		return false
	}

	if m.scrollOffset+m.PageSize >= len(m.choices) {
		return false
	}

	return true
}

func (m *MultiModel[T]) canScrollUp() bool {
	return m.scrollOffset > 0
}

func (m *MultiModel[T]) cursorDown() {
	if m.currentIdx == len(m.currentChoices)-1 {
		if m.canScrollDown() {
			m.scrollDown()
		} else if m.LoopCursor {
			m.scrollToTop()

			return
		}
	}

	m.currentIdx = min(len(m.currentChoices)-1, m.currentIdx+1)
}

func (m *MultiModel[T]) cursorUp() {
	if m.currentIdx == 0 {
		if m.canScrollUp() {
			m.scrollUp()
		} else if m.LoopCursor {
			m.scrollToBottom()

			return
		}
	}

	m.currentIdx = max(0, m.currentIdx-1)
}

func (m *MultiModel[T]) scrollDown() {
	if m.PageSize <= 0 || m.scrollOffset+m.PageSize >= m.availableChoices {
		return
	}

	m.currentIdx = max(0, m.currentIdx-1)
	m.scrollOffset++
	m.currentChoices, m.availableChoices = m.filteredAndPagedChoices()
}

func (m *MultiModel[T]) scrollToBottom() {
	m.currentIdx = len(m.currentChoices) - 1

	if m.PageSize <= 0 || m.availableChoices < m.PageSize {
		return
	}

	m.scrollOffset = m.availableChoices - m.PageSize
	m.currentChoices, m.availableChoices = m.filteredAndPagedChoices()
}

func (m *MultiModel[T]) scrollUp() {
	if m.PageSize <= 0 || m.scrollOffset <= 0 {
		return
	}

	m.currentIdx = min(len(m.currentChoices)-1, m.currentIdx+1)
	m.scrollOffset--
	m.currentChoices, m.availableChoices = m.filteredAndPagedChoices()
}

func (m *MultiModel[T]) scrollToTop() {
	m.currentIdx = 0

	if m.PageSize <= 0 || m.availableChoices < m.PageSize {
		return
	}

	m.scrollOffset = 0
	m.currentChoices, m.availableChoices = m.filteredAndPagedChoices()
}

func (m *MultiModel[T]) reindexChoices() {
	for i, choice := range m.choices {
		choice.idx = i
	}
}

func (m *MultiModel[T]) markedChoicesInOrder() []*Choice[T] {
	choices := make([]*Choice[T], 0, len(m.selectedChoices))

	for _, c := range m.choices {
		if m.selectedChoices[c.idx] {
			choices = append(choices, c)
		}
	}

	return choices
}

func (m *MultiModel[T]) finalChoicesStr(choices []*Choice[T]) string {
	strs := make([]string, len(choices))

	for i, c := range choices {
		if m.FinalChoiceStyle == nil {
			strs[i] = c.String
		} else {
			strs[i] = m.FinalChoiceStyle(c)
		}
	}

	return strings.Join(strs, ", ")
}
