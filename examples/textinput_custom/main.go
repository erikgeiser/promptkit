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

	sp := textinput.New()
	sp.Prompt = "Enter an IP address:"
	sp.Placeholder = "e.g. 127.0.0.1"
	sp.Validate = func(input string) bool { return net.ParseIP(input) != nil }
	sp.Template = customTemplate

	ip, err := sp.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	// do something with the result
	_ = ip
}
