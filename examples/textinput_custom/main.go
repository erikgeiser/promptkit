package main

import (
	"fmt"
	"net"
	"os"

	"github.com/erikgeiser/promptkit/textinput"
)

func main() {
	const customTemplate = `
	{{- Bold .Prompt }} {{ .Input -}}
	{{- if not .Valid }} {{ Foreground "1" "✘" }}
	{{- else }} 🖥️{{- end -}}
	`

	input := textinput.New()
	input.Prompt = "Enter an IP address:"
	input.Placeholder = "e.g. 127.0.0.1"
	input.Validate = func(input string) bool { return net.ParseIP(input) != nil }
	input.Template = customTemplate

	ip, err := input.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	// do something with the result
	_ = ip
}
