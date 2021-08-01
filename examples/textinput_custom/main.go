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
	{{- if not .Valid -}}
		{{- Foreground "1" (Bold "✘") -}}
	{{- else -}}
		{{- Foreground "2" (Bold "✔") -}}
	{{- end -}}┃
	{{- "\n┗" }}━{{ Repeat "━" (Len .Prompt) }}━┷━{{ Repeat "━" 13 }}{{ "━━━━┛" -}}
	`

	const customConfirmationTemplate = `
	{{- Bold (print "🖥️  Connecting to " (Foreground "32" .FinalValue) "\n") -}}
	`

	input := textinput.New()
	input.Prompt = "Enter an IP address"
	input.Placeholder = "127.0.0.1"
	input.Validate = func(input string) bool { return net.ParseIP(input) != nil }
	input.Template = customTemplate
	input.ConfirmationTemplate = customConfirmationTemplate
	input.CharLimit = 15

	ip, err := input.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	// do something with the result
	_ = ip
}
