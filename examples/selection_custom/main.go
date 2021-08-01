package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/erikgeiser/promptkit/selection"
)

func main() {
	// nolint:lll
	const (
		customTemplate = `
{{- if .Prompt -}}
  {{ Bold .Prompt }}
{{ end -}}
{{ if .IsFiltered }}
  {{- print "Filter by ID: " .FilterInput }}
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
   {{- Foreground "32" (print (Bold (print "[x] " $choice.Value.Name)) (Faint (print " (" $choice.Value.ID ") " "\n"))) }}
  {{- else }}
    {{- print "[ ] " $choice.Value.Name (Faint (print " (" $choice.Value.ID ") ")) "\n"}}
  {{- end }}
{{- end}}`
		customConfirmationTempalte = `
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

	sp := selection.New("Choose an article!", selection.Choices(choices))
	sp.FilterPlaceholder = "Type to filter"
	sp.PageSize = 3
	sp.Filter = func(filter string, choice *selection.Choice) bool {
		chosenArticle, _ := choice.Value.(article)

		return strings.HasPrefix(chosenArticle.ID, filter)
	}
	sp.Template = customTemplate
	sp.ConfirmationTemplate = customConfirmationTempalte
	sp.ExtendedTemplateScope = map[string]interface{}{
		"name": func(c *selection.Choice) string { return c.Value.(article).Name },
	}

	choice, err := sp.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// do something with the final choice
	_ = choice
}
