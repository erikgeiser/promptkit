package selection

import (
	"fmt"
	"io"
	"os"
	"text/template"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/erikgeiser/promptkit"
	"github.com/muesli/termenv"
)

const (
	// DefaultMultiTemplate defines the default appearance of the multi-selection
	// prompt and can be copied as a starting point for a custom template.
	DefaultMultiTemplate = `
{{- if .Prompt -}}
  {{ Bold .Prompt }}{{ if and (gt .MinSelections 0) (gt .MaxSelections 0) -}}
    {{- if not (and (ge .NMarkedChoices .MinSelections)
                    (le .NMarkedChoices .MaxSelections)) -}}
      {{- print " " }}{{ Faint (print "Select " .MinSelections "-" .MaxSelections " items") -}}
    {{- end -}}
  {{- else if gt .MinSelections 0 -}}
    {{- if lt .NMarkedChoices .MinSelections }} {{ Faint (print "Select at least " .MinSelections
      " items") -}}
    {{- end -}}
  {{- else if gt .MaxSelections 0 -}}
    {{- if gt .NMarkedChoices .MaxSelections }} {{ Faint (print "Select up to " .MaxSelections " items") -}}
    {{- end -}}
  {{- end }}
{{ end -}}
{{ if .IsFiltered }}
  {{- print .FilterPrompt " " .FilterInput }}
{{ end -}}
{{- range $i, $choice := .Choices }}
  {{- if IsScrollUpHintPosition $i }}
    {{- "⇡ " -}}
  {{- else if IsScrollDownHintPosition $i -}}
    {{- "⇣ " -}}
  {{- else -}}
    {{- "  " -}}
  {{- end -}}
  {{- if eq $.SelectedIndex $i }}
    {{- if IsMarked $choice }}
      {{- print (Foreground "32" (Bold "▸ ")) "☑ " (CursorMarked $choice) "\n" }}
    {{- else }}
      {{- print (Foreground "32" (Bold "▸ ")) "☐ " (Cursor $choice) "\n" }}
    {{- end }}
  {{- else if IsMarked $choice }}
    {{- print "  ☑ " (Marked $choice) "\n" }}
  {{- else }}
    {{- print "  ☐ " (Unmarked $choice) "\n" }}
  {{- end }}
{{- end}}
{{- range .OffscreenMarkedChoices }}
  {{- print (Faint "  ☑ ") (Faint (Marked .)) "\n" }}
{{- end }}`

	// DefaultMultiResultTemplate defines the default appearance with which the
	// final result of the multi-selection prompt is presented.
	DefaultMultiResultTemplate = `
	{{- print .Prompt " " .FinalChoicesStr "\n" -}}
	`
)

// DefaultCursorMarkedChoiceStyle is the default style for a marked choice at
// the cursor position.
func DefaultCursorMarkedChoiceStyle[T any](c *Choice[T]) string {
	return termenv.String(c.String).Foreground(accentColor).Bold().String()
}

// DefaultMarkedChoiceStyle is the default style for marked choices.
func DefaultMarkedChoiceStyle[T any](c *Choice[T]) string {
	return termenv.String(c.String).Foreground(accentColor).String()
}

// MultiSelection represents a configurable multi-selection prompt.
type MultiSelection[T any] struct {
	// choices represent all selectable choices. Slices of arbitrary types can
	// be converted to a slice of choices using asChoices.
	choices []*Choice[T]

	// Prompt holds the prompt text or question that is to be answered by
	// marking one or more of the choices.
	Prompt string

	// FilterPrompt is the prompt for the filter if filtering is enabled.
	FilterPrompt string

	// Filter is a function that decides whether a given choice should be
	// displayed based on the text entered by the user into the filter input
	// field. If Filter is nil, filtering will be disabled. By default the
	// filter FilterContainsCaseInsensitive is used.
	Filter func(filterText string, choice *Choice[T]) bool

	// FilterPlaceholder holds the text that is displayed in the filter input
	// field when no text was entered by the user yet. If empty, the
	// DefaultFilterPlaceholder is used. If Filter is nil, filtering is
	// disabled and FilterPlaceholder does nothing.
	FilterPlaceholder string

	// PageSize is the number of choices that are displayed at once. If
	// PageSize is smaller than the number of choices, pagination is enabled.
	// If PageSize is 0, pagination is disabled. Regardless of the value of
	// PageSize, pagination is always enabled when the prompt does not fit the
	// terminal.
	PageSize int

	// LoopCursor enables the cursor to loop around to the first choice when
	// navigating down from the last choice and the other way around.
	LoopCursor bool

	// MinSelections is the minimum number of choices that must be marked
	// before the prompt can be confirmed. A value of 0 allows confirming with
	// no choices marked. The default is 1.
	MinSelections int

	// MaxSelections is the maximum number of choices that can be marked. Once
	// the limit is reached, toggling unmarked choices has no effect. A value
	// of 0 means there is no limit.
	MaxSelections int

	// Template holds the display template. A custom template can be used to
	// completely customize the appearance of the multi-selection prompt. If
	// empty, DefaultMultiTemplate is used. The following variables and
	// functions are available:
	//
	//  * Prompt string: The configured prompt.
	//  * IsFiltered bool: Whether or not filtering is enabled.
	//  * FilterPrompt string: The configured filter prompt.
	//  * FilterInput string: The view of the filter input model.
	//  * Choices []*Choice: The choices on the current page.
	//  * NChoices int: The number of choices on the current page.
	//  * SelectedIndex int: The index of the cursor in the current page.
	//  * PageSize int: The configured page size.
	//  * IsPaged bool: Whether pagination is currently active.
	//  * AllChoices []*Choice: All configured choices.
	//  * NAllChoices int: The number of configured choices.
	//  * MarkedChoices []*Choice: All currently marked choices in original order.
	//  * NMarkedChoices int: The number of currently marked choices.
	//  * OffscreenMarkedChoices []*Choice: Marked choices not visible in the current
	//    page because they are scrolled out of view or filtered out.
	//  * NOffscreenMarkedChoices int: The number of offscreen marked choices.
	//  * TerminalWidth int: The width of the terminal.
	//  * IsMarked(*Choice) bool: Whether the given choice is marked.
	//  * Cursor(*Choice) string: The configured CursorChoiceStyle.
	//  * CursorMarked(*Choice) string: The configured CursorMarkedChoiceStyle.
	//  * Marked(*Choice) string: The configured MarkedChoiceStyle.
	//  * Unmarked(*Choice) string: The configured UnmarkedChoiceStyle.
	//  * IsScrollDownHintPosition(idx int) bool: Returns whether the scroll
	//    down hint should be displayed at the given index.
	//  * IsScrollUpHintPosition(idx int) bool: Returns whether the scroll up
	//    hint should be displayed at the given index.
	//  * promptkit.UtilFuncMap: Handy helper functions.
	//  * termenv TemplateFuncs (see https://github.com/muesli/termenv).
	//  * The functions specified in ExtendedTemplateFuncs.
	Template string

	// ResultTemplate is rendered as soon as the choices have been confirmed.
	// It is intended to permanently indicate the result of the prompt when the
	// selection itself has disappeared. This template is only rendered in the
	// RunPrompt() method and NOT when the prompt is used as a model. The
	// following variables and functions are available:
	//
	//  * FinalChoices []*Choice: The marked choices in original order.
	//  * FinalChoicesStr string: The marked choices rendered with FinalChoiceStyle
	//    and joined by ", ".
	//  * Prompt string: The configured prompt.
	//  * AllChoices []*Choice: All configured choices.
	//  * NAllChoices int: The number of configured choices.
	//  * TerminalWidth int: The width of the terminal.
	//  * Final(*Choice) string: The configured FinalChoiceStyle.
	//  * promptkit.UtilFuncMap: Handy helper functions.
	//  * termenv TemplateFuncs (see https://github.com/muesli/termenv).
	//  * The functions specified in ExtendedTemplateFuncs.
	ResultTemplate string

	// ExtendedTemplateFuncs can be used to add additional functions to the
	// evaluation scope of the templates.
	ExtendedTemplateFuncs template.FuncMap

	// Styles of the filter input field. These will be applied as inline
	// styles.
	//
	// For an introduction to styling with Lip Gloss see:
	// https://github.com/charmbracelet/lipgloss
	FilterInputTextStyle        lipgloss.Style
	FilterInputBackgroundStyle  lipgloss.Style // Deprecated: This property is not used anymore.
	FilterInputPlaceholderStyle lipgloss.Style
	FilterInputCursorStyle      lipgloss.Style

	// CursorChoiceStyle allows customizing the appearance of the choice at the
	// cursor position when it is not marked. If nil, no style is applied and
	// the plain string representation is used. This style is available as the
	// template function Cursor.
	CursorChoiceStyle func(*Choice[T]) string

	// CursorMarkedChoiceStyle allows customizing the appearance of the choice
	// at the cursor position when it is marked. By default
	// DefaultCursorMarkedChoiceStyle is used. If nil, no style is applied.
	// This style is available as the template function CursorMarked.
	CursorMarkedChoiceStyle func(*Choice[T]) string

	// MarkedChoiceStyle allows customizing the appearance of marked choices
	// that are not at the cursor position. By default DefaultMarkedChoiceStyle
	// is used. If nil, no style is applied. This style is available as the
	// template function Marked.
	MarkedChoiceStyle func(*Choice[T]) string

	// UnmarkedChoiceStyle allows customizing the appearance of choices that
	// are neither marked nor at the cursor position. If nil, no style is
	// applied and the plain string representation is used. This style is
	// available as the template function Unmarked.
	UnmarkedChoiceStyle func(*Choice[T]) string

	// FinalChoiceStyle allows customizing the appearance of the confirmed
	// choices in the result template. By default DefaultMarkedChoiceStyle is
	// used. If nil, no style is applied. This style is available as the
	// template function Final.
	FinalChoiceStyle func(*Choice[T]) string

	// KeyMap determines with which keys the multi-selection prompt is
	// controlled. By default, DefaultMultiKeyMap is used.
	KeyMap *MultiKeyMap

	// WrapMode decides which way the prompt view is wrapped if it does not fit
	// the terminal. It can be a WrapMode provided by promptkit or a custom
	// function. By default it is promptkit.Truncate. It can also be nil which
	// disables wrapping and likely causes output glitches.
	WrapMode promptkit.WrapMode

	// Output is the output writer, by default os.Stdout is used.
	Output io.Writer

	// Input is the input reader, by default, os.Stdin is used.
	Input io.Reader

	// ColorProfile determines how colors are rendered. By default, the
	// terminal is queried.
	ColorProfile termenv.Profile
}

// NewMulti creates a new multi-selection prompt. See the MultiSelection
// properties for more documentation.
func NewMulti[T any](prompt string, choices []T) *MultiSelection[T] {
	return &MultiSelection[T]{
		choices:                     asChoices(choices),
		Prompt:                      prompt,
		FilterPrompt:                DefaultFilterPrompt,
		Template:                    DefaultMultiTemplate,
		ResultTemplate:              DefaultMultiResultTemplate,
		Filter:                      FilterContainsCaseInsensitive[T],
		FilterInputPlaceholderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		CursorMarkedChoiceStyle:     DefaultCursorMarkedChoiceStyle[T],
		MarkedChoiceStyle:           DefaultMarkedChoiceStyle[T],
		FinalChoiceStyle:            DefaultMarkedChoiceStyle[T],
		KeyMap:                      NewDefaultMultiKeyMap(),
		FilterPlaceholder:           DefaultFilterPlaceholder,
		ExtendedTemplateFuncs:       template.FuncMap{},
		WrapMode:                    promptkit.Truncate,
		MinSelections:               1,
		Output:                      os.Stdout,
		Input:                       os.Stdin,
	}
}

// RunPrompt executes the multi-selection prompt.
func (s *MultiSelection[T]) RunPrompt() ([]T, error) {
	err := validateMultiKeyMap(s.KeyMap)
	if err != nil {
		return nil, fmt.Errorf("insufficient key map: %w", err)
	}

	m := NewMultiModel(s)

	p := tea.NewProgram(m, tea.WithOutput(s.Output), tea.WithInput(s.Input))

	_, err = p.Run()
	if err != nil {
		return nil, fmt.Errorf("running prompt: %w", err)
	}

	return m.Values()
}
