// Package main demonstrates how promptkit/confirmation is used.
package main

import (
	"fmt"
	"os"

	"github.com/erikgeiser/promptkit/confirmation"
)

func main() {
	input := confirmation.New("Are you ready?", confirmation.Undecided)

	ready, err := input.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	// do something with the result
	_ = ready
}
