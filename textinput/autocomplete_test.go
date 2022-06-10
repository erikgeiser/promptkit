package textinput

import (
	"strconv"
	"testing"
)

func TestCommonPrefix(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		inputs []string

		prefix string
	}{
		{[]string{}, ""},
		{[]string{"abc", "abcd"}, "abc"},
		{[]string{"foobar", "fo", "foobarbaz"}, "fo"},
	}

	for i, testCase := range testCases { // nolint:paralleltest
		i, testCase := i, testCase

		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			prefix := commonPrefix(testCase.inputs)
			if prefix != testCase.prefix {
				t.Fatalf("got common prefix %q for input %#v instead of %q",
					prefix, testCase.inputs, testCase.prefix)
			}
		})
	}
}
