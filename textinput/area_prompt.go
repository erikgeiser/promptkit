package textinput

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
	// DefaultAreaTemplate defines the default appearance of the text area and
	// can be copied as a starting point for a custom template.
	DefaultAreaTemplate = `
	{{- if .Prompt }}{{ Bold .Prompt }}
{{ end -}}
{{- .Input -}}
{{- if .ValidationError }} {{ Foreground "1" (Bold "✘") }}
{{- else }} {{ Foreground "2" (Bold "✔") }}
{{- end }}
{{ Faint (printf "%s to submit • %s for newline" (index .KeyMap.Submit 0) (index .KeyMap.InsertNewline 0)) -}}
`

	// DefaultAreaResultTemplate defines the default appearance with which the
	// final result of the prompt is presented.
	DefaultAreaResultTemplate = `
	{{- print .Prompt " " (Foreground "32" .FinalValue) "\n" -}}
	`

	// DefaultAreaHeight is the default number of visible rows in the text area.
	DefaultAreaHeight = 6
)

// TextArea represents a configurable multi-line text area prompt.
type TextArea struct {
	// Prompt holds the prompt text or question that is printed above the
	// text area in the default template (if not empty).
	Prompt string

	// Placeholder holds the text that is displayed in the text area when no
	// input has been entered yet.
	Placeholder string

	// InitialValue sets the initial content of the text area, as if it had
	// been entered by the user. This can be used to provide an editable
	// default value.
	InitialValue string

	// Validate is a function that validates whether the current input is
	// valid. If it is not, the data cannot be submitted. By default, Validate
	// ensures that the input is not empty. If Validate is set to nil, no
	// validation is performed.
	Validate func(string) error

	// CharLimit is the maximum number of characters this text area will
	// accept. If 0 or less, there is no limit.
	CharLimit int

	// Height is the number of visible rows in the text area. If 0 or less,
	// the text area automatically expands and contracts to fit the content.
	// Defaults to 0 (auto-expand). Set to DefaultAreaHeight or another
	// positive value for a fixed-height text area.
	Height int

	// ShowLineNumbers specifies whether line numbers are displayed. Defaults
	// to false.
	ShowLineNumbers bool

	// Template holds the display template. A custom template can be used to
	// completely customize the appearance of the text area. If empty,
	// DefaultAreaTemplate is used. The following variables and functions are
	// available:
	//
	//  * Prompt string: The configured prompt.
	//  * InitialValue string: The configured initial value.
	//  * Placeholder string: The configured placeholder.
	//  * Input string: The rendered text area widget.
	//  * ValidationError error: The error returned by Validate.
	//  * TerminalWidth int: The width of the terminal.
	//  * KeyMap *AreaKeyMap: The configured key map.
	//  * promptkit.UtilFuncMap: Handy helper functions.
	//  * termenv TemplateFuncs (see https://github.com/muesli/termenv).
	//  * The functions specified in ExtendedTemplateFuncs.
	Template string

	// ResultTemplate is rendered once input has been confirmed. It is
	// intended to permanently indicate the result of the prompt. This
	// template is only rendered in RunPrompt() and NOT when the text area is
	// used as a model. The following variables and functions are available:
	//
	//  * FinalValue string: The submitted text.
	//  * Prompt string: The configured prompt.
	//  * InitialValue string: The configured initial value.
	//  * Placeholder string: The configured placeholder.
	//  * TerminalWidth int: The width of the terminal.
	//  * KeyMap *AreaKeyMap: The configured key map.
	//  * promptkit.UtilFuncMap: Handy helper functions.
	//  * termenv TemplateFuncs (see https://github.com/muesli/termenv).
	//  * The functions specified in ExtendedTemplateFuncs.
	ResultTemplate string

	// ExtendedTemplateFuncs can be used to add additional functions to the
	// evaluation scope of the templates.
	ExtendedTemplateFuncs template.FuncMap

	// InputTextStyle is the style applied to the text in the text area.
	InputTextStyle lipgloss.Style

	// InputPlaceholderStyle is the style applied to the placeholder text.
	InputPlaceholderStyle lipgloss.Style

	// KeyMap determines with which keys the text area is controlled. By
	// default, NewDefaultAreaKeyMap() is used.
	KeyMap *AreaKeyMap

	// WrapMode decides which way the prompt view is wrapped if it does not
	// fit the terminal. It can be a WrapMode provided by promptkit or a
	// custom function. By default it is promptkit.WordWrap. It can also be
	// nil which disables wrapping and likely causes output glitches.
	WrapMode promptkit.WrapMode

	// Output is the output writer, by default os.Stdout is used.
	Output io.Writer
	// Input is the input reader, by default os.Stdin is used.
	Input io.Reader

	// ColorProfile determines how colors are rendered. By default, the
	// terminal is queried.
	ColorProfile termenv.Profile
}

// NewArea creates a new text area prompt. See the TextArea properties for more
// documentation.
func NewArea(prompt string) *TextArea {
	return &TextArea{
		Prompt:                prompt,
		Template:              DefaultAreaTemplate,
		ResultTemplate:        DefaultAreaResultTemplate,
		KeyMap:                NewDefaultAreaKeyMap(),
		InputPlaceholderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		Validate: ValidateNotEmpty,
		ExtendedTemplateFuncs: template.FuncMap{},
		WrapMode:              promptkit.WordWrap,
		Output:                os.Stdout,
		Input:                 os.Stdin,
	}
}

// RunPrompt executes the text area prompt.
func (t *TextArea) RunPrompt() (string, error) {
	err := validateAreaKeyMap(t.KeyMap)
	if err != nil {
		return "", fmt.Errorf("insufficient key map: %w", err)
	}

	m := NewAreaModel(t)

	p := tea.NewProgram(m, tea.WithOutput(t.Output), tea.WithInput(t.Input))

	_, err = p.Run()
	if err != nil {
		return "", fmt.Errorf("running prompt: %w", err)
	}

	return m.Value()
}
