// Package main demonstrates how promptkit/selection.MultiSelection can be customized.
// This example showcases custom styling and templates to create a distinctive look.
package main

import (
	"fmt"
	"os"

	"github.com/erikgeiser/promptkit/selection"
	"github.com/muesli/termenv"
)

func main() {
	const (
		// Custom template with emoji checkboxes and a cleaner layout.
		customTemplate = `{{ Bold .Prompt }}
{{ if .IsFiltered }}{{ Faint (print "filter: " .FilterInput) }}{{ end }}
{{ range $i, $choice := .Choices -}}
	{{- if eq $.SelectedIndex $i }}▸ {{ end }}
	{{- if IsMarked $choice }}✓ {{ Marked $choice }}{{ else }}· {{ Unmarked $choice }}{{ end }}
{{ end }}
{{- if .NMarkedChoices }}{{ Faint (print .NMarkedChoices " selected, up to " $.MaxSelections) }}{{ end }}`

		resultTemplate = `{{ Bold .Prompt }} {{ .FinalChoicesStr }}`
	)

	type lang struct {
		name string
		icon string
	}

	choices := []lang{
		{name: "Python", icon: "🐍"},
		{name: "Go", icon: "🐹"},
		{name: "Rust", icon: "🦀"},
		{name: "TypeScript", icon: "📘"},
		{name: "Haskell", icon: "λ"},
	}

	output := termenv.NewOutput(os.Stdout)
	accent := output.String().Foreground(termenv.ANSI256Color(211)) //nolint:gomnd

	sp := selection.NewMulti("Pick your favorite languages:", choices)
	sp.MinSelections = 2
	sp.MaxSelections = 4
	sp.LoopCursor = true
	sp.PreSelected = selection.PreSelect(choices[1])

	// Custom styling with emoji icons.
	sp.UnmarkedChoiceStyle = func(c *selection.Choice[lang]) string {
		return c.Value.icon + "  " + c.Value.name
	}
	sp.CursorChoiceStyle = func(c *selection.Choice[lang]) string {
		return accent.Styled(c.Value.icon + "  " + c.Value.name)
	}
	sp.MarkedChoiceStyle = func(c *selection.Choice[lang]) string {
		return accent.Styled(c.Value.icon + "  " + c.Value.name)
	}
	sp.CursorMarkedChoiceStyle = func(c *selection.Choice[lang]) string {
		return accent.Bold().Styled(c.Value.icon + "  " + c.Value.name)
	}
	sp.FinalChoiceStyle = func(c *selection.Choice[lang]) string {
		return c.Value.icon + " " + c.Value.name
	}

	sp.Template = customTemplate
	sp.ResultTemplate = resultTemplate

	selected, err := sp.RunPrompt()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	_ = selected
}
