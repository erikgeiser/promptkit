package textinput

import (
	"sort"
	"strings"
)

// AutoCompleteFromSlice creates a case-insensitive auto-complete function from
// a slice of choices.
func AutoCompleteFromSlice(choices []string) func(string) []string {
	return autoCompleteFromSlice(choices, false)
}

// AutoCompleteFromSliceWithDefault creates a case-insensitive auto-complete
// function from a slice of choices with a default completion value that is
// inserted if the function is called on an empty input.
func AutoCompleteFromSliceWithDefault(
	choices []string, defaultValue string,
) func(string) []string {
	autoComplete := autoCompleteFromSlice(choices, false)

	return func(s string) []string {
		if s == "" {
			return []string{defaultValue}
		}

		return autoComplete(s)
	}
}

// CaseSensitiveAutoCompleteFromSlice creates a case-sensitive auto-complete
// function from a slice of choices.
func CaseSensitiveAutoCompleteFromSlice(choices []string) func(string) []string {
	return autoCompleteFromSlice(choices, true)
}

// CaseSensitiveAutoCompleteFromSliceWithDefault creates a case-sensitive
// auto-complete function from a slice of choices with a default completion
// value that is inserted if the function is called on an empty input.
func CaseSensitiveAutoCompleteFromSliceWithDefault(
	choices []string, defaultValue string,
) func(string) []string {
	autoComplete := autoCompleteFromSlice(choices, true)

	return func(s string) []string {
		if s == "" {
			return []string{defaultValue}
		}

		return autoComplete(s)
	}
}

func autoCompleteFromSlice(choices []string, caseSensitive bool) func(string) []string {
	return func(value string) []string {
		v := value
		if !caseSensitive {
			v = strings.ToLower(value)
		}

		var candidates []string

		for _, choice := range choices {
			ch := choice
			if !caseSensitive {
				ch = strings.ToLower(choice)
			}

			if strings.HasPrefix(ch, v) {
				candidates = append(candidates, choice)
			}
		}

		return candidates
	}
}

func commonPrefix(candidates []string) string {
	if len(candidates) == 0 {
		return ""
	}

	commonPrefix := ""
	endPrefix := false

	if len(candidates) > 0 {
		sort.Strings(candidates)
		first := candidates[0]
		last := candidates[len(candidates)-1]

		for i := 0; i < len(first); i++ {
			if !endPrefix && string(last[i]) == string(first[i]) {
				commonPrefix += string(last[i])
			} else {
				endPrefix = true
			}
		}
	}

	return commonPrefix
}
