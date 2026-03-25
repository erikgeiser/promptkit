// Package main demonstrates how promptkit/textinput is used.
package main

import (
	"fmt"
	"os"

	"github.com/erikgeiser/promptkit/textinput"
)

func main() {
	firstNamePrompt := textinput.New("What is your first name?")
	firstNamePrompt.InitialValue = os.Getenv("USER")
	firstNamePrompt.Placeholder = "Your first name cannot be empty"

	firstName, err := firstNamePrompt.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	lastNamePrompt := textinput.New("What is your last name?")
	lastNamePrompt.Placeholder = "Your last name cannot be empty"

	lastName, err := lastNamePrompt.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	fmt.Println("Hi", firstName, lastName)
}
