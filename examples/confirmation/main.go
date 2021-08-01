package main

import (
	"fmt"
	"os"

	"github.com/erikgeiser/promptkit/confirmation"
)

func main() {
	input := confirmation.New()
	input.Prompt = "Are you ready?"
	input.DefaultValue = confirmation.Undecided

	ready, err := input.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	// do something with the result
	_ = ready
}
