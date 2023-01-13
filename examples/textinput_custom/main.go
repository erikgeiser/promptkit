// Package main demonstrates how promptkit/textinput can be customized.
package main

import (
	"fmt"
	"net"
	"os"

	"github.com/erikgeiser/promptkit/textinput"
)

func main() {
	const customTemplate = `
	{{- "┏" }}━{{ Repeat "━" (Len .Prompt) }}━┯━{{ Repeat "━" 13 }}{{ "━━━━┓\n" }}
	{{- "┃" }} {{ Bold .Prompt }} │ {{ .Input -}}
	{{- Repeat " " (Max 0 (Sub 16 (Len .Input))) }}
	{{- if .ValidationError -}}
		{{- Foreground "1" (Bold "✘") -}}
	{{- else -}}
		{{- Foreground "2" (Bold "✔") -}}
	{{- end -}}┃
	{{- if .ValidationError -}}
		{{- (print " Error: " (Foreground "1" .ValidationError.Error)) -}}
	{{- end -}}
	{{- "\n┗" }}━{{ Repeat "━" (Len .Prompt) }}━┷━{{ Repeat "━" 13 }}{{ "━━━━┛\n" -}}
	{{- if .AutoCompleteIndecisive -}}
		{{ print "  Suggestions: " }}
		{{- range $suggestion := AutoCompleteSuggestions -}}
			{{- print $suggestion " " -}}
		{{- end -}}
	{{- end -}}
	`

	const customResultTemplate = `
	{{- Bold (print "🖥️  Connecting to " (Foreground "32" .FinalValue) ) -}}`

	input := textinput.New("Enter an IP address")
	input.Placeholder = "127.0.0.1"
	input.Validate = func(input string) error {
		if net.ParseIP(input) == nil {
			return fmt.Errorf("invalid IP address")
		}

		return nil
	}
	input.Template = customTemplate
	input.ResultTemplate = customResultTemplate
	input.CharLimit = 15
	input.AutoComplete = textinput.AutoCompleteFromSliceWithDefault([]string{
		"10.0.0.1",
		"10.0.0.2",
		"127.0.0.1",
		"fe80::1",
	}, input.Placeholder)

	ip, err := input.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	// do something with the result
	_ = ip
}
