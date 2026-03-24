// Package main demonstrates how promptkit/selection.MultiSelection can be customized.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/erikgeiser/promptkit/selection"
	"github.com/muesli/termenv"
)

func main() {
	const (
		customTemplate = `
{{- if .Prompt -}}
  {{ Bold .Prompt }}
{{ end -}}
{{ if .IsFiltered }}
  {{- print .FilterPrompt " " .FilterInput }}
{{ end -}}
{{- range .OffscreenMarkedChoices }}
  {{- print "  " (Foreground "32" "[x] ") (Marked .) "\n" }}
{{- end }}
{{- range $i, $choice := .Choices }}
  {{- if IsScrollUpHintPosition $i }}
    {{- print "⇡ " -}}
  {{- else if IsScrollDownHintPosition $i -}}
    {{- print "⇣ " -}}
  {{- else -}}
    {{- print "  " -}}
  {{- end -}}

  {{- if eq $.SelectedIndex $i }}
    {{- if IsMarked $choice }}
      {{- print (Foreground "32" (Bold "▸ ")) (Foreground "32" (Bold "[x] ")) (CursorMarked $choice) "\n" }}
    {{- else }}
      {{- print (Foreground "32" (Bold "▸ ")) "[ ] " (Cursor $choice) "\n" }}
    {{- end }}
  {{- else if IsMarked $choice }}
    {{- print "  " (Foreground "32" "[x] ") (Marked $choice) "\n" }}
  {{- else }}
    {{- print "  " "[ ] " (Unmarked $choice) "\n" }}
  {{- end }}
{{- end}}
{{- if .NMarkedChoices }}
  {{ Faint (print .NMarkedChoices " selected") }}
{{ end }}`

		resultTemplate = `
		{{- print .Prompt " " (Foreground "32" .FinalChoicesStr) "\n" -}}
		`
	)

	type language struct {
		Name     string
		Paradigm string
	}

	choices := []language{
		{Name: "Go", Paradigm: "compiled"},
		{Name: "Rust", Paradigm: "compiled"},
		{Name: "Python", Paradigm: "interpreted"},
		{Name: "TypeScript", Paradigm: "interpreted"},
		{Name: "Haskell", Paradigm: "functional"},
		{Name: "Clojure", Paradigm: "functional"},
	}

	output := termenv.NewOutput(os.Stdout)
	green := output.String().Foreground(termenv.ANSI256Color(32)) //nolint:gomnd

	sp := selection.NewMulti("Which languages do you use?", choices)
	sp.FilterPrompt = "Filter by paradigm:"
	sp.FilterPlaceholder = "compiled / interpreted / functional"
	sp.PageSize = 4
	sp.LoopCursor = true
	sp.MinSelections = 1
	sp.MaxSelections = 3
	sp.Filter = func(filter string, choice *selection.Choice[language]) bool {
		return strings.HasPrefix(choice.Value.Paradigm, filter) ||
			strings.Contains(strings.ToLower(choice.Value.Name), strings.ToLower(filter))
	}
	sp.Template = customTemplate
	sp.ResultTemplate = resultTemplate
	sp.CursorMarkedChoiceStyle = func(c *selection.Choice[language]) string {
		return green.Bold().Styled(c.Value.Name)
	}
	sp.MarkedChoiceStyle = func(c *selection.Choice[language]) string {
		return green.Styled(c.Value.Name) + " " +
			output.String(c.Value.Paradigm).Faint().String()
	}
	sp.UnmarkedChoiceStyle = func(c *selection.Choice[language]) string {
		return c.Value.Name + " " +
			output.String(c.Value.Paradigm).Faint().String()
	}
	sp.FinalChoiceStyle = func(c *selection.Choice[language]) string {
		return c.Value.Name
	}
	sp.ExtendedTemplateFuncs = map[string]any{
		"Faint": func(s string) string { return output.String(s).Faint().String() },
	}

	choices2, err := sp.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	// do something with the final choices
	_ = choices2
}
