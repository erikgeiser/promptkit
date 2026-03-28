// Package main demonstrates how promptkit/textinput textarea can be customized.
package main

import (
	"fmt"
	"os"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/erikgeiser/promptkit/textinput"
)

func main() {
	const customTemplate = `
	{{- Bold .Prompt }}
{{ .Input -}}
{{- if .ValidationError }}
{{ Foreground "1" (Bold "✘ ") }}{{ Foreground "1" .ValidationError.Error }}{{ end }}
{{ Foreground "8" "Use ctrl+d to submit • Use ctrl+c to abort" }}
`

	const customResultTemplate = `
	{{- Bold "Commit message: " }}{{ Foreground "32" .FinalValue }}
`

	input := textinput.NewArea("Commit message:")
	input.Placeholder = "feat: add amazing new feature"
	input.Height = 3
	input.Template = customTemplate
	input.ResultTemplate = customResultTemplate
	input.KeyMap = textinput.NewDefaultAreaKeyMap()
	input.KeyMap.Submit = []string{"ctrl+d"}
	input.Validate = func(s string) error {
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("commit message cannot be empty")
		}

		if len(s) > 50 {
			return fmt.Errorf("commit message is too long (%d/50)", len(s))
		}

		return nil
	}
	input.InputPlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)

	message, err := input.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	// do something with the result
	_ = message
}
