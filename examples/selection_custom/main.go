// Package main demonstrates how promptkit/selection can be customized.
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
{{ end }}

{{- range  $i, $choice := .Choices }}
  {{- if IsScrollUpHintPosition $i }}
    {{- print "⇡ " -}}
  {{- else if IsScrollDownHintPosition $i -}}
    {{- print "⇣ " -}} 
  {{- else -}}
    {{- print "  " -}}
  {{- end -}} 

  {{- if eq $.SelectedIndex $i }}
   {{- print "[" (Foreground "32" (Bold "x")) "] " (Selected $choice) "\n" }}
  {{- else }}
    {{- print "[ ] " (Unselected $choice) "\n" }}
  {{- end }}
{{- end}}`
		resultTemplate = `
		{{- print .Prompt " " (Foreground "32"  (name .FinalChoice)) "\n" -}}
		`
	)

	type article struct {
		ID   string
		Name string
	}

	choices := []article{
		{ID: "123", Name: "Article A"},
		{ID: "321", Name: "Article B"},
		{ID: "345", Name: "Article C"},
		{ID: "456", Name: "Article D"},
		{ID: "444", Name: "Article E"},
	}

	blue := termenv.String().Foreground(termenv.ANSI256Color(32)) //nolint:gomnd

	sp := selection.New("Choose an article!", choices)
	sp.FilterPrompt = "Filter by ID:"
	sp.FilterPlaceholder = "Type to filter"
	sp.PageSize = 3
	sp.LoopCursor = true
	sp.Filter = func(filter string, choice *selection.Choice[article]) bool {
		return strings.HasPrefix(choice.Value.ID, filter)
	}
	sp.Template = customTemplate
	sp.ResultTemplate = resultTemplate
	sp.SelectedChoiceStyle = func(c *selection.Choice[article]) string {
		return (blue.Bold().Styled(c.Value.Name) + " " +
			termenv.String("("+c.Value.ID+")").Faint().String())
	}
	sp.UnselectedChoiceStyle = func(c *selection.Choice[article]) string {
		return c.Value.Name + " " + termenv.String("("+c.Value.ID+")").Faint().String()
	}
	sp.ExtendedTemplateFuncs = map[string]interface{}{
		"name": func(c *selection.Choice[article]) string { return c.Value.Name },
	}

	choice, err := sp.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// do something with the final choice
	_ = choice
}
