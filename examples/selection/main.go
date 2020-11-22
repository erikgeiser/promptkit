package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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

	p := tea.NewProgram(sp)
	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	choice, err := sp.Choice()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Your choice: %v\n", choice)
}
