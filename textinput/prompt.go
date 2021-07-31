/*
Package textinput implements prompt for a string input that can also be used for
secret strings such as passwords. It also offers customizable appreance as well
as optional support for input validation and a customizable key map.
*/
package textinput

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/erikgeiser/promptkit"
	"github.com/muesli/termenv"
)

const (
	// DefaultTemplate defines the default appearance of the text input and can
	// be copied as a starting point for a custom template.
	DefaultTemplate = `
	{{- Bold .Prompt }} {{ .Input -}}
	{{- if not .Valid }} {{ Foreground "1" "✘" }}
	{{- else }} {{ Foreground "2" "✔" }}
	{{- end -}}
	`

	// DefaultConfirmationTemplate defines the default appearance with which the
	// finale result of the prompt is presented.
	DefaultConfirmationTemplate = `
	{{- print (Bold (print .Prompt " " (Foreground "32"  (Mask .FinalValue)))) "\n" -}}
	`

	// DefaultMask specified the character with which the input is masked by
	// default if Hidden is true.
	DefaultMask = '●'
)

// TextInput represents a configurable selection prompt.
type TextInput struct {
	// Prompt holds the the prompt text or question that is printed above the
	// choices in the default template (if not empty).
	Prompt string

	// Placeholder holds the text that is displayed in the input field when the
	// input data is empty, e.g. when no text was entered yet.
	Placeholder string

	// InitialValue is similar to Placeholder, however, the actual input data is
	// set to InitialValue such that as if it was entered by the user. This can
	// be used to provide an editable default value.
	InitialValue string

	// Validate is a function that validates whether the current input data is
	// valid. If it is not, the data cannot be submitted. By default, Validate
	// ensures that the input data is not empty. If Validate is set to nil, no
	// validation is performed.
	Validate func(string) bool

	// Hidden specified whether or not the input data is considered secret and
	// should be masked. This is useful for password prompts.
	Hidden bool

	// HideMask specified the character with which the input data should be
	// masked when Hidden is set to true.
	HideMask rune

	// CharLimit is the maximum amount of characters this input element will
	// accept. If 0 or less, there's no limit.
	CharLimit int

	// Width is the maximum number of characters that can be displayed at once.
	// It essentially treats the text field like a horizontally scrolling
	// viewport. If 0 or less this setting is ignored.
	Width int

	// Template holds the display template. A custom template can be used to
	// completely customize the appearance of the text input. If empty,
	// the DefaultTemplate is used. The following variables and functions are
	// available:
	//
	//  * Prompt string: The configured prompt.
	//  * InitialValue string: The configured initial value of the input.
	//  * Placeholder string: The configured placeholder of the input.
	//  * Input string: The actual input field.
	//  * Valid bool: Whether or not the current value is valid according
	//    to the configured Validate function.
	//  * promptkit.UtilFuncMap: Handy helper functions.
	//  * termenv TemplateFuncs (see https://github.com/muesli/termenv).
	//  * The functions specified in ExtendedTemplateScope.
	Template string

	// ConfirmationTemplate is rendered as soon as a input has been confirmed.
	// It is intended to permanently indicate the result of the prompt when the
	// input itself has disappeared. This template is only rendered in the Run()
	// method and NOT when the text input is used as a model. The following
	// variables and functions are available:
	//
	//  * FinalChoice: The choice that was selected by the user.
	//  * Prompt string: The configured prompt.
	//  * InitialValue string: The configured initial value of the input.
	//  * Placeholder string: The configured placeholder of the input.
	//  * Mask(string) string: A function that replaces all characters of
	//    a string with the character specified in HideMask if Hidden is
	//    true and returns the input string if Hidden is false.
	//  * promptkit.UtilFuncMap: Handy helper functions.
	//  * termenv TemplateFuncs (see https://github.com/muesli/termenv).
	//  * The functions specified in ExtendedTemplateScope.
	ConfirmationTemplate string

	// ExtendedTemplateScope can be used to add additional functions to the
	// evaluation scope of the templates.
	ExtendedTemplateScope template.FuncMap

	// Styles of the actual input field. These will be applied as inline styles.
	//
	// For an introduction to styling with Lip Gloss see:
	// https://github.com/charmbracelet/lipgloss
	InputTextStyle        lipgloss.Style
	InputBackgroundStyle  lipgloss.Style
	InputPlaceholderStyle lipgloss.Style
	InputCursorStyle      lipgloss.Style

	// KeyMap determines with which keys the selection prompt is controlled. By
	// default, DefaultKeyMap is used.
	KeyMap *KeyMap
}

// New creates a new text input.
func New() *TextInput {
	return &TextInput{
		Template:              DefaultTemplate,
		ConfirmationTemplate:  DefaultConfirmationTemplate,
		KeyMap:                NewDefaultKeyMap(),
		InputPlaceholderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		Validate:              func(s string) bool { return len(s) > 0 },
		HideMask:              DefaultMask,
	}
}

func (t *TextInput) initConfirmationTemplate() (*template.Template, error) {
	if t.ConfirmationTemplate == "" {
		return nil, nil
	}

	tmpl := template.New("confirmed")
	tmpl.Funcs(termenv.TemplateFuncs(termenv.ColorProfile()))
	tmpl.Funcs(promptkit.UtilFuncMap())
	tmpl.Funcs(t.ExtendedTemplateScope)
	tmpl.Funcs(template.FuncMap{"Mask": t.mask})

	return tmpl.Parse(t.ConfirmationTemplate)
}

// RunPrompt executes the text input prompt.
func (t *TextInput) RunPrompt() (string, error) {
	tmpl, err := t.initConfirmationTemplate()
	if err != nil {
		return "", fmt.Errorf("initializing confirmation template: %w", err)
	}

	m := NewModel(t)

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		return "", err
	}

	value, err := m.Value()
	if err != nil {
		return "", err
	}

	if t.ConfirmationTemplate != "" {
		err = tmpl.Execute(os.Stdout, map[string]interface{}{
			"FinalValue":   value,
			"Prompt":       m.Prompt,
			"InitialValue": m.InitialValue,
			"Placeholder":  m.Placeholder,
			"Hidden":       m.Hidden,
		})
	}

	return value, err
}

// mask replaces each character with HideMask if Hidden is true.
func (t *TextInput) mask(s string) string {
	if !t.Hidden {
		return s
	}

	return strings.Repeat(string(t.HideMask), len(s))
}
