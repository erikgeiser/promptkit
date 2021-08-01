package main

import (
	"fmt"
	"os"

	"github.com/erikgeiser/promptkit/selection"
)

func main() {
	sp := selection.New("What do you pick?",
		selection.Choices([]string{"Horse", "Car", "Plane", "Bike"}))
	sp.PageSize = 3

	choice, err := sp.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	// do something with the final choice
	_ = choice
}
