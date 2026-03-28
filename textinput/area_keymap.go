package textinput

import (
	"fmt"
	"slices"

	tea "charm.land/bubbletea/v2"
)

// NewDefaultAreaKeyMap returns an AreaKeyMap with sensible default key mappings
// that can also be used as a starting point for customization.
func NewDefaultAreaKeyMap() *AreaKeyMap {
	return &AreaKeyMap{
		MoveBackward:           []string{"left", "ctrl+b"},
		MoveForward:            []string{"right", "ctrl+f"},
		MoveWordBackward:       []string{"alt+left", "alt+b"},
		MoveWordForward:        []string{"alt+right", "alt+f"},
		MoveUp:                 []string{"up", "ctrl+p"},
		MoveDown:               []string{"down", "ctrl+n"},
		JumpToLineBeginning:    []string{"home", "ctrl+a"},
		JumpToLineEnd:          []string{"end", "ctrl+e"},
		JumpToBeginning:        []string{"ctrl+home"},
		JumpToEnd:              []string{"ctrl+end"},
		DeleteBeforeCursor:     []string{"backspace"},
		DeleteWordBeforeCursor: []string{"ctrl+w", "alt+backspace"},
		DeleteUnderCursor:      []string{"delete"},
		DeleteWordAfterCursor:  []string{"alt+delete", "alt+d"},
		DeleteAllAfterCursor:   []string{"ctrl+k"},
		DeleteAllBeforeCursor:  []string{"ctrl+u"},
		InsertNewline:          []string{"enter"},
		Paste:                  []string{"ctrl+v"},
		Submit:                 []string{"shift+enter", "ctrl+d"},
		Abort:                  []string{"ctrl+c"},
	}
}

// upstreamAreaKeyMap lists keys handled natively by the underlying
// bubbles/textarea component. These are blocked from being passed through if
// they do not match any action in the AreaKeyMap, to avoid conflicts.
var upstreamAreaKeyMap = &AreaKeyMap{
	MoveBackward:           []string{"left", "ctrl+b"},
	MoveForward:            []string{"right", "ctrl+f"},
	MoveWordBackward:       []string{"alt+left", "alt+b"},
	MoveWordForward:        []string{"alt+right", "alt+f"},
	MoveUp:                 []string{"up", "ctrl+p"},
	MoveDown:               []string{"down", "ctrl+n"},
	JumpToLineBeginning:    []string{"home", "ctrl+a"},
	JumpToLineEnd:          []string{"end", "ctrl+e"},
	JumpToBeginning:        []string{"ctrl+home", "alt+<"},
	JumpToEnd:              []string{"ctrl+end", "alt+>"},
	DeleteBeforeCursor:     []string{"backspace", "ctrl+h"},
	DeleteWordBeforeCursor: []string{"ctrl+w", "alt+backspace"},
	DeleteUnderCursor:      []string{"delete", "ctrl+d"},
	DeleteWordAfterCursor:  []string{"alt+delete", "alt+d"},
	DeleteAllAfterCursor:   []string{"ctrl+k"},
	DeleteAllBeforeCursor:  []string{"ctrl+u"},
	InsertNewline:          []string{"enter", "ctrl+m"},
	Paste:                  []string{"ctrl+v"},
}

// AreaKeyMap defines the keys that trigger certain actions in the text area.
type AreaKeyMap struct {
	MoveBackward           []string
	MoveForward            []string
	MoveWordBackward       []string
	MoveWordForward        []string
	MoveUp                 []string
	MoveDown               []string
	JumpToLineBeginning    []string
	JumpToLineEnd          []string
	JumpToBeginning        []string
	JumpToEnd              []string
	DeleteBeforeCursor     []string
	DeleteWordBeforeCursor []string
	DeleteUnderCursor      []string
	DeleteWordAfterCursor  []string
	DeleteAllAfterCursor   []string
	DeleteAllBeforeCursor  []string
	InsertNewline          []string
	Paste                  []string
	Submit                 []string
	Abort                  []string
}

func areaKeyMatches(key tea.KeyPressMsg, mapping []string) bool {
	return slices.Contains(mapping, key.String())
}

func areaKeyMatchesUpstream(key tea.KeyPressMsg) bool {
	return areaKeyMatches(key, allAreaKeys(upstreamAreaKeyMap))
}

// validateAreaKeyMap returns an error if the key map is missing required bindings.
func validateAreaKeyMap(km *AreaKeyMap) error {
	if len(km.Submit) == 0 {
		return fmt.Errorf("no submit key")
	}

	if len(km.Abort) == 0 {
		return fmt.Errorf("no abort key")
	}

	return nil
}

func allAreaKeys(km *AreaKeyMap) (keys []string) {
	keys = append(keys, km.MoveBackward...)
	keys = append(keys, km.MoveForward...)
	keys = append(keys, km.MoveWordBackward...)
	keys = append(keys, km.MoveWordForward...)
	keys = append(keys, km.MoveUp...)
	keys = append(keys, km.MoveDown...)
	keys = append(keys, km.JumpToLineBeginning...)
	keys = append(keys, km.JumpToLineEnd...)
	keys = append(keys, km.JumpToBeginning...)
	keys = append(keys, km.JumpToEnd...)
	keys = append(keys, km.DeleteBeforeCursor...)
	keys = append(keys, km.DeleteWordBeforeCursor...)
	keys = append(keys, km.DeleteUnderCursor...)
	keys = append(keys, km.DeleteWordAfterCursor...)
	keys = append(keys, km.DeleteAllAfterCursor...)
	keys = append(keys, km.DeleteAllBeforeCursor...)
	keys = append(keys, km.InsertNewline...)
	keys = append(keys, km.Paste...)
	keys = append(keys, km.Submit...)
	keys = append(keys, km.Abort...)

	return keys
}
