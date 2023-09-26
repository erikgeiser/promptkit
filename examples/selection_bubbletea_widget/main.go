// Package main demonstrates how promptkit can be used as a bubbletea widget.
package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/erikgeiser/promptkit/selection"
)

type shoppingCart struct {
	availableItems []string
	addedItems     map[string]int
	selection      *selection.Model[string]
	err            error
}

func newShoppingCart(items ...string) *shoppingCart {
	return &shoppingCart{availableItems: items, addedItems: make(map[string]int)}
}

var _ tea.Model = &shoppingCart{}

func (s *shoppingCart) Init() tea.Cmd {
	sel := selection.New("Add Items to Your Shopping Cart:", s.availableItems)
	sel.Filter = nil

	s.selection = selection.NewModel(sel)

	return s.selection.Init()
}

func (s *shoppingCart) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return s, nil
	}

	switch {
	case keyMsg.String() == "enter":
		c, err := s.selection.Value()
		if err != nil {
			s.err = err

			return s, tea.Quit
		}

		s.addedItems[c]++
	case keyMsg.String() == "esc":
		return s, tea.Quit
	default:
		_, cmd := s.selection.Update(msg)

		return s, cmd
	}

	return s, nil
}

func (s *shoppingCart) View() string {
	if s.err != nil {
		return fmt.Sprintf("Error: %v", s.err)
	}

	var b strings.Builder

	b.WriteString(s.selection.View())
	b.WriteString("=== Your Shopping Cart: ===\n")

	if len(s.addedItems) == 0 {
		b.WriteString("no items\n")

		return b.String()
	}

	for item, amount := range s.addedItems {
		fmt.Fprintf(&b, "%dx %s\n", amount, item)
	}

	return b.String()
}

func main() {
	model := newShoppingCart("Apples", "Milk", "Bread")

	p := tea.NewProgram(model)

	_, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)

		os.Exit(1)
	}
}
