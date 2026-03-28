// Package main demonstrates how promptkit/textinput is used.
package main

import (
	"fmt"
	"os"

	"github.com/erikgeiser/promptkit/textinput"
)

func main() {
	input := textinput.New("Choose a language?")
	input.Placeholder = "Use tab for auto completion"
	input.AutoComplete = textinput.AutoCompleteFromSlice([]string{
		"Go",
		"Python",
		"Rust",
	})

	lang, err := input.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	// do something with the result
	_ = lang
}
