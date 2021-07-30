package main

import (
	"fmt"
	"os"

	"github.com/erikgeiser/promptkit/textinput"
)

func main() {
	input := textinput.New()
	input.Prompt = "Choose a password:"
	input.Placeholder = "minimum 10 characters"
	input.Validate = func(s string) bool { return len(s) >= 10 } // nolint:gomnd
	input.Hidden = true

	name, err := input.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	// do something with the result
	_ = name
}
