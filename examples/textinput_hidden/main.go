package main

import (
	"fmt"
	"os"

	"github.com/erikgeiser/promptkit/textinput"
)

func main() {
	sp := textinput.New()
	sp.Prompt = "Choose a password:"
	sp.Placeholder = "minimum 10 characters"
	sp.Validate = func(s string) bool { return len(s) >= 10 } // nolint:gomnd
	sp.Hidden = true

	name, err := sp.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	// do something with the result
	_ = name
}
