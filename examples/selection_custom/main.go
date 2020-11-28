package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/erikgeiser/promptkit/selection"
)

// nolint:lll
const customTemplate = `
{{- if .Label -}}
  {{ Bold .Label }}
{{ end -}}
{{ if .Filter }}
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
    {{- print "[ ] " $choice.Value.Name " (" $choice.Value.ID ") " "\n"}}
  {{- end }}
{{- end}}`

type Article struct {
	ID   string
	Name string
}

func main() {
	choices := []Article{
		{ID: "123", Name: "Article A"},
		{ID: "234", Name: "Article B"},
		{ID: "345", Name: "Article C"},
		{ID: "456", Name: "Article D"},
		{ID: "567", Name: "Article E"},
	}

	sp := &selection.Prompt{
		Label:    "Choose an article!",
		Choices:  selection.SliceChoices(choices),
		Template: customTemplate,
		Filter: func(filter string, choice *selection.Choice) bool {
			article, _ := choice.Value.(Article)

			return strings.Contains(article.ID, filter)
		},
		PageSize: 3,
	}

	choice, err := sp.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Your choice: %v\n", choice.Value)
}
