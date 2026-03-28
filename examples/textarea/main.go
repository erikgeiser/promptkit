// Package main demonstrates how promptkit/textinput is used for multi-line input.
package main

import (
	"fmt"
	"os"

	"github.com/erikgeiser/promptkit/textinput"
)

func main() {
	input := textinput.NewArea("Describe your issue:")
	input.Placeholder = "Please provide as much detail as possible"

	description, err := input.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	// do something with the result
	_ = description
}
