// Package main demonstrates how promptkit/confirmation can be customized.
package main

import (
	"fmt"
	"os"

	"github.com/erikgeiser/promptkit/confirmation"
)

func main() {
	input := confirmation.New("Do you want to try out promptkit?",
		confirmation.NewValue(true))
	input.Template = confirmation.TemplateYN
	input.ResultTemplate = confirmation.ResultTemplateYN
	input.KeyMap.SelectYes = append(input.KeyMap.SelectYes, "+")
	input.KeyMap.SelectNo = append(input.KeyMap.SelectNo, "-")

	ready, err := input.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	// do something with the result
	_ = ready
}
