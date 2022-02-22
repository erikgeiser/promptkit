package main

import (
	"fmt"
	"os"

	"github.com/erikgeiser/promptkit/textinput"
)

func main() {
	input := textinput.New("Choose a password:")
	input.Placeholder = "pick a strong password"
	input.Validate = func(s string) error {
		if len(s) < 10 {
			return fmt.Errorf("needs %d more characters", 10-len(s))
		}
		if s == "1234567890" {
			return fmt.Errorf("too easy")
		}
		return nil
	}
	input.Hidden = true

	name, err := input.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	// do something with the result
	_ = name
}
