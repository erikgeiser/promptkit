// Package main demonstrates how promptkit can be used as a bubbletea widget.
package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/erikgeiser/promptkit"
	"github.com/erikgeiser/promptkit/confirmation"
)

type confirmationPrompt struct {
	prompt       string
	selection    bool
	confirmation *confirmation.Model
	err          error
}

func newConfirmationPrompt(prompt string) *confirmationPrompt {
	return &confirmationPrompt{prompt: prompt}
}

var _ tea.Model = &confirmationPrompt{}

func (c *confirmationPrompt) Init() tea.Cmd {
	// Here we show how to specify your own keyMap since the default uses left
	// and right arrow keys which you may already be using.
	keyMap := &confirmation.KeyMap{
		Yes:    []string{"y", "Y"},
		No:     []string{"n", "N"},
		Toggle: []string{"tab"},
		Submit: []string{"enter"},
		Abort:  []string{"ctrl+c"},
	}

	conf := &confirmation.Confirmation{
		Prompt:         c.prompt,
		DefaultValue:   confirmation.Undecided,
		Template:       confirmation.DefaultTemplate,
		ResultTemplate: confirmation.DefaultResultTemplate,
		KeyMap:         keyMap,
		WrapMode:       promptkit.Truncate,
	}
	c.confirmation = confirmation.NewModel(conf)

	return c.confirmation.Init()
}

func (c *confirmationPrompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return c, nil
	}

	switch {
	case keyMsg.String() == "enter":
		selection, err := c.confirmation.Value()
		if err != nil {
			c.err = err
			return c, tea.Quit
		}

		c.selection = selection
	case keyMsg.String() == "esc":
		return c, tea.Quit
	default:
		_, cmd := c.confirmation.Update(msg)
		return c, cmd
	}

	return c, nil
}

func (c *confirmationPrompt) View() string {
	if c.err != nil {
		return fmt.Sprintf("Error: %v", c.err)
	}

	var b strings.Builder
	b.WriteString(c.confirmation.View())
	b.WriteString("\n=== You Chose: ===\n")
	b.WriteString(fmt.Sprintf("%t", c.selection))

	return b.String()
}

func main() {
	model := newConfirmationPrompt("Would you like to go shopping?")

	p := tea.NewProgram(model)

	_, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)

		os.Exit(1)
	}
}
