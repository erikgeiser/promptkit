package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/erikgeiser/promptkit/selection"
)

func main() {
	sp := &selection.Prompt{
		Label:   "What do you pick?",
		Choices: selection.StringChoices([]string{"Horse", "Car", "Plane", "Bike"}),
		Filter: func(filter string, choice *selection.Choice) bool {
			return strings.Contains(strings.ToLower(choice.String), strings.ToLower(filter))
		},
		PageSize: 3,
	}

	choice, err := sp.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Your choice: %v\n", choice.Value)
}
