/*
Package selection implements a selection prompt that allows users to to select
one of the pre-defined choices. It also offers customizable appreance and key
map as well as optional support for pagination, filtering.
*/
package selection

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	// DefaultTemplate defines the default appearance of the selection and can
	// be copied as a starting point for a custom template.
	DefaultTemplate = `
{{- if .Prompt -}}
  {{ Bold .Prompt }}
{{ end -}}
{{ if .IsFiltered }}
  {{- print .FilterPrompt " " .FilterInput }}
{{ end }}

{{- range  $i, $choice := .Choices }}
  {{- if IsScrollUpHintPosition $i }}
    {{- "⇡ " -}}
  {{- else if IsScrollDownHintPosition $i -}}
    {{- "⇣ " -}} 
  {{- else -}}
    {{- "  " -}}
  {{- end -}} 

  {{- if eq $.SelectedIndex $i }}
   {{- Foreground "32" (Bold (print "▸ " $choice.String "\n")) }}
  {{- else }}
    {{- print "  " $choice.String "\n"}}
  {{- end }}
{{- end}}`

	// DefaultConfirmationTemplate defines the default appearance with which the
	// finale result of the selection is presented.
	DefaultConfirmationTemplate = `
	{{- print .Prompt " " (Foreground "32"  .FinalChoice.String) "\n" -}}
	`

	// DefaultFilterPrompt is the default prompt for the filter input when
	// filtering is enabled.
	DefaultFilterPrompt = "Filter:"

	// DefaultFilterPlaceholder is printed by default when no filter text was
	// entered yet.
	DefaultFilterPlaceholder = "Type to filter choices"
)

// Selection represents a configurable selection prompt.
type Selection struct {
	// Choices represent all selectable choices of the selection. Slices of
	// arbitrary types can be converted to a slice of choices using the helper
	// selection.Choices.
	Choices []*Choice

	// Prompt holds the the prompt text or question that is to be answered by
	// one of the choices.
	Prompt string

	// FilterPrompt is the prompt for the filter if filtering is enabled.
	FilterPrompt string

	// Filter is a function that decides whether a given choice should be
	// displayed based on the text entered by the user into the filter input
	// field. If Filter is nil, filtering will be disabled. By default the
	// filter FilterContainsCaseInsensitive is used.
	Filter func(filterText string, choice *Choice) bool

	// FilterPlaceholder holds the text that is displayed in the filter input
	// field when no text was entered by the user yet. If empty, the
	// DefaultFilterPlaceholder is used. If Filter is nil, filtering is disabled
	// and FilterPlaceholder does nothing.
	FilterPlaceholder string

	// PageSize is the number of choices that are displayed at once. If PageSize
	// is smaller than the number of choices, pagination is enabled. If PageSize
	// is 0, pagenation is always disabled.
	PageSize int

	// Template holds the display template. A custom template can be used to
	// completely customize the appearance of the selection prompt. If empty,
	// the DefaultTemplate is used. The following variables and functions are
	// available:
	//
	//  * Prompt string: The configured prompt.
	//  * IsFiltered bool: Whether or not filtering is enabled.
	//  * FilterPrompt string: The configured filter prompt.
	//  * FilterInput string: The view of the filter input model.
	//  * Choices []*Choice: The choices on the current page.
	//  * NChoices int: The number of choices on the current page.
	//  * SelectedIndex int: The index that is currently selected.
	//  * PageSize int: The configured page size.
	//  * IsPaged bool: Whether pagination is currently active.
	//  * AllChoices []*Choice: All configured choices.
	//  * NAllChoices int: The number of configured choices.
	//  * TerminalWidth int: The width of the terminal.
	//  * IsScrollDownHintPosition(idx int) bool: Returns whether
	//    the scroll down hint shoud be displayed at the given index.
	//  * IsScrollUpHintPosition(idx int) bool: Returns whether the
	//    scroll up hint shoud be displayed at the given index).
	//  * promptkit.UtilFuncMap: Handy helper functions.
	//  * termenv TemplateFuncs (see https://github.com/muesli/termenv).
	//  * The functions specified in ExtendedTemplateScope.
	Template string

	// ConfirmationTemplate is rendered as soon as a choice has been selected.
	// It is intended to permanently indicate the result of the prompt when the
	// selection itself has disappeared. This template is only rendered in the
	// Run() method and NOT when the selection prompt is used as a model. The
	// following variables and functions are available:
	//
	//  * FinalChoice: The choice that was selected by the user.
	//  * Prompt string: The configured prompt.
	//  * AllChoices []*Choice: All configured choices.
	//  * NAllChoices int: The number of configured choices.
	//  * TerminalWidth int: The width of the terminal.
	//  * promptkit.UtilFuncMap: Handy helper functions.
	//  * termenv TemplateFuncs (see https://github.com/muesli/termenv).
	//  * The functions specified in ExtendedTemplateScope.
	ConfirmationTemplate string

	// ExtendedTemplateScope can be used to add additional functions to the
	// evaluation scope of the templates.
	ExtendedTemplateScope template.FuncMap

	// Styles of the filter input field. These will be applied as inline styles.
	//
	// For an introduction to styling with Lip Gloss see:
	// https://github.com/charmbracelet/lipgloss
	FilterInputTextStyle        lipgloss.Style
	FilterInputBackgroundStyle  lipgloss.Style
	FilterInputPlaceholderStyle lipgloss.Style
	FilterInputCursorStyle      lipgloss.Style

	// KeyMap determines with which keys the selection prompt is controlled. By
	// default, DefaultKeyMap is used.
	KeyMap *KeyMap

	// Output is the output writer, by default os.Stdout is used.
	Output io.Writer
	// Input is the input reader, by default, os.Stdin is used.
	Input io.Reader
}

// New creates a new selection prompt.
func New(prompt string, choices []*Choice) *Selection {
	return &Selection{
		Choices:                     choices,
		Prompt:                      prompt,
		FilterPrompt:                DefaultFilterPrompt,
		Template:                    DefaultTemplate,
		ConfirmationTemplate:        DefaultConfirmationTemplate,
		Filter:                      FilterContainsCaseInsensitive,
		FilterInputPlaceholderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		KeyMap:                      NewDefaultKeyMap(),
		FilterPlaceholder:           DefaultFilterPlaceholder,
		ExtendedTemplateScope:       template.FuncMap{},
		Output:                      os.Stdout,
		Input:                       os.Stdin,
	}
}

// RunPrompt executes the selection prompt.
func (s *Selection) RunPrompt() (*Choice, error) {
	err := validateKeyMap(s.KeyMap)
	if err != nil {
		return nil, fmt.Errorf("insufficient key map: %w", err)
	}

	m := NewModel(s)

	p := tea.NewProgram(m, tea.WithOutput(s.Output), tea.WithInput(s.Input))
	if err := p.Start(); err != nil {
		return nil, fmt.Errorf("running prompt: %w", err)
	}

	return m.Value()
}

// FilterContainsCaseInsensitive returns true if the string representation of
// the choice contains the filter string without regard for capitalization.
func FilterContainsCaseInsensitive(filter string, choice *Choice) bool {
	return strings.Contains(strings.ToLower(choice.String), strings.ToLower(filter))
}

// FilterContainsCaseSensitive returns true if the string representation of the
// choice contains the filter string respecting capitalization.
func FilterContainsCaseSensitive(filter string, choice *Choice) bool {
	return strings.Contains(choice.String, filter)
}
