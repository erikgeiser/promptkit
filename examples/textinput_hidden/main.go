// Package main demonstrates how promptkit/textinput can be used to ask for
// passphrases.
package main

import (
	"fmt"
	"os"

	"github.com/erikgeiser/promptkit/textinput"
)

const minCharacters = 10

func main() {
	input := textinput.New("Choose a passphrase:")
	input.Placeholder = "make it strong!"
	input.Validate = func(s string) error {
		if s == "hunter2" {
			return fmt.Errorf("not that one")
		}

		if len(s) < minCharacters {
			return fmt.Errorf("at least %d more characters", minCharacters-len(s))
		}

		return nil
	}
	input.Hidden = true
	input.Template += `
	{{- if .ValidationError -}}
		{{- print " " (Foreground "1" .ValidationError.Error) -}}
	{{- end -}}`

	name, err := input.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	// do something with the result
	_ = name
}
