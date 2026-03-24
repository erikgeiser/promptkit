// Package main demonstrates how promptkit/selection.MultiSelection is used.
package main

import (
	"fmt"
	"os"

	"github.com/erikgeiser/promptkit/selection"
)

func main() {
	sp := selection.NewMulti("What do you want to eat?",
		[]string{"Pizza", "Burger", "Sushi", "Salad", "Pasta", "Tacos"})
	sp.PageSize = 4

	choices, err := sp.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	// do something with the final choices
	_ = choices
}
