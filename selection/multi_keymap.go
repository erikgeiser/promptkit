package selection

import "fmt"

// NewDefaultMultiKeyMap returns a MultiKeyMap with sensible default key
// mappings that can also be used as a starting point for customization.
func NewDefaultMultiKeyMap() *MultiKeyMap {
	return &MultiKeyMap{
		KeyMap: *NewDefaultKeyMap(),
		Toggle: []string{"space"},
	}
}

// MultiKeyMap defines the keys that trigger actions in the multi-selection
// prompt. It extends KeyMap with a Toggle key for marking and unmarking
// choices.
type MultiKeyMap struct {
	KeyMap
	Toggle []string
}

func validateMultiKeyMap(km *MultiKeyMap) error {
	err := validateKeyMap(&km.KeyMap)
	if err != nil {
		return err
	}

	if len(km.Toggle) == 0 {
		return fmt.Errorf("no toggle key")
	}

	return nil
}
