package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/erikgeiser/promptkit/selection"
)

func main() {
	sp := selection.NewModel(selection.Choices([]string{"Horse", "Car", "Plane", "Bike"}))
	sp.Label = "What do you pick?"
	sp.Filter = func(filter string, choice *selection.Choice) bool {
		return strings.Contains(strings.ToLower(choice.String), strings.ToLower(filter))
	}
	sp.PageSize = 3

	choice, err := sp.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Your choice: %v\n", choice.Value)
}
