// Package main demonstrates how promptkit/textinput is used.
package main

import (
	"fmt"
	"os"

	"github.com/erikgeiser/promptkit/textinput"
)

func main() {
	input := textinput.New("What is your name?")
	input.InitialValue = os.Getenv("USER")
	input.Placeholder = "Your name cannot be empty"

	name, err := input.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	// do something with the result
	_ = name
}
