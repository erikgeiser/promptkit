package main

import (
	"fmt"
	"os"

	"github.com/erikgeiser/promptkit/textinput"
)

func main() {
	sp := textinput.New()
	sp.Prompt = "What is your name?"
	sp.InitialValue = os.Getenv("USER")
	sp.Placeholder = "Your name cannot be empty"

	name, err := sp.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	// do something with the result
	_ = name
}
