package selection

import (
	"fmt"
	"reflect"
)

type Choice struct {
	Index  int
	String string
	Value  interface{}
}

func NewChoice(item interface{}) *Choice {
	choice := &Choice{Index: 0, Value: item}

	switch i := item.(type) {
	case Choice:
		choice.String = i.String
		choice.Value = i.Value
	case *Choice:
		choice.String = i.String
		choice.Value = i.Value
	case string:
		choice.String = i
	case fmt.Stringer:
		choice.String = i.String()
	default:
		choice.String = fmt.Sprintf("%#v", i)
	}

	return choice
}

func StringChoices(choiceStrings []string) []*Choice {
	choices := make([]*Choice, 0, len(choiceStrings))

	for _, c := range choiceStrings {
		choices = append(choices, NewChoice(c))
	}

	return choices
}

func StringerChoices(choiceStrings []fmt.Stringer) []*Choice {
	choices := make([]*Choice, 0, len(choiceStrings))

	for _, c := range choiceStrings {
		choices = append(choices, NewChoice(c))
	}

	return choices
}

func SliceChoices(sliceChoices interface{}) []*Choice {
	switch reflect.TypeOf(sliceChoices).Kind() {
	case reflect.Slice:
		slice := reflect.ValueOf(sliceChoices)

		choices := make([]*Choice, 0, slice.Len())

		for i := 0; i < slice.Len(); i++ {
			choices = append(choices, NewChoice(slice.Index(i).Interface()))
		}

		return choices
	default:
		panic("SliceChoices argument is not a slice")
	}
}

func InterfaceChoices(choiceInterfaces []interface{}) []*Choice {
	choices := make([]*Choice, 0, len(choiceInterfaces))

	for _, c := range choiceInterfaces {
		choices = append(choices, NewChoice(c))
	}

	return choices
}
