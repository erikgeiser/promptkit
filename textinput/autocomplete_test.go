//nolint:goconst
package textinput

import (
	"strconv"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/erikgeiser/promptkit/test"
)

func TestAutocomplete(t *testing.T) {
	t.Parallel()

	defaultValue := "default"
	commonPrefix := "abcd"

	m := NewModel(New("foo:"))
	m.AutoComplete = AutoCompleteFromSliceWithDefault(
		[]string{commonPrefix + "ef", commonPrefix + "gh"}, defaultValue,
	)

	test.Run(t, m)
	assertNoError(t, m)

	test.Update(t, m, tea.KeyTab)

	if !m.autoCompleteTriggered {
		t.Fatalf("auto-complete trigger indication missing")
	}

	if m.autoCompleteIndecisive {
		t.Fatalf("auto-complete indecisive indication wrongfully triggered")
	}

	v := getValue(t, m)
	if v != defaultValue {
		t.Fatalf("completing default value resulted in %q instead of %q",
			v, defaultValue)
	}

	test.Update(t, m, tea.KeyEsc) // reset

	test.Update(t, m, test.KeyMsg('a'))

	if m.autoCompleteTriggered {
		t.Fatalf("auto-complete trigger was not reset")
	}

	test.Update(t, m, tea.KeyTab)

	if !m.autoCompleteTriggered {
		t.Fatalf("auto-complete trigger indication missing")
	}

	if !m.autoCompleteIndecisive {
		t.Fatalf("auto-complete indecisive indication did not trigger")
	}

	v = getValue(t, m)
	if v != commonPrefix {
		t.Fatalf("completing common prefix resulted in %q instead of %q",
			v, commonPrefix)
	}

	test.Update(t, m, test.KeyMsg('g'))

	if m.autoCompleteTriggered {
		t.Fatalf("auto-complete trigger was not reset")
	}

	test.Update(t, m, tea.KeyTab)

	if !m.autoCompleteTriggered {
		t.Fatalf("auto-complete trigger indication missing")
	}

	if m.autoCompleteIndecisive {
		t.Fatalf("auto-complete indecisive indication triggered wrongfully")
	}

	v = getValue(t, m)
	if v != commonPrefix+"gh" {
		t.Fatalf("completing full suggestion resulted in %q instead of %q",
			v, commonPrefix+"gh")
	}
}

func TestAutoCompleteFromSlice(t *testing.T) {
	t.Parallel()

	choices := []string{
		"abcde", "abxyz", "foo", "aBCD",
	}

	testCases := []struct {
		input   string
		outputs []string
	}{
		{"", choices},
		{"a", []string{"abcde", "abxyz", "aBCD"}},
		{"ab", []string{"abcde", "abxyz", "aBCD"}},
		{"aB", []string{"abcde", "abxyz", "aBCD"}},
		{"fo", []string{"foo"}},
		{"Fo", []string{"foo"}},
	}

	for i, testCase := range testCases {
		i, testCase := i, testCase

		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			assertSameContents(t,
				AutoCompleteFromSlice(choices)(testCase.input), testCase.outputs)
		})
	}
}

func TestAutoCompleteFromSliceWithDefault(t *testing.T) {
	t.Parallel()

	defaultValue := "default"
	choices := []string{
		"abcde", "abxyz", "foo", "aBCD",
	}

	testCases := []struct {
		input   string
		outputs []string
	}{
		{"", []string{defaultValue}},
		{"a", []string{"abcde", "abxyz", "aBCD"}},
		{"ab", []string{"abcde", "abxyz", "aBCD"}},
		{"aB", []string{"abcde", "abxyz", "aBCD"}},
		{"fo", []string{"foo"}},
		{"Fo", []string{"foo"}},
	}

	for i, testCase := range testCases {
		i, testCase := i, testCase

		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			assertSameContents(t,
				AutoCompleteFromSliceWithDefault(choices, defaultValue)(testCase.input),
				testCase.outputs)
		})
	}
}

func TestCaseSensitiveAutoCompleteFromSlice(t *testing.T) {
	t.Parallel()

	choices := []string{
		"abcde", "abxyz", "foo", "aBCD",
	}

	testCases := []struct {
		input   string
		outputs []string
	}{
		{"", choices},
		{"a", []string{"abcde", "abxyz", "aBCD"}},
		{"ab", []string{"abcde", "abxyz"}},
		{"aB", []string{"aBCD"}},
		{"fo", []string{"foo"}},
		{"Fo", nil},
	}

	for i, testCase := range testCases {
		i, testCase := i, testCase

		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			assertSameContents(t,
				CaseSensitiveAutoCompleteFromSlice(choices)(testCase.input),
				testCase.outputs)
		})
	}
}

func TestCaseSensitiveAutoCompleteFromSliceWithDefault(t *testing.T) {
	t.Parallel()

	defaultValue := "default"
	choices := []string{
		"abcde", "abxyz", "foo", "aBCD",
	}

	testCases := []struct {
		input   string
		outputs []string
	}{
		{"", []string{defaultValue}},
		{"a", []string{"abcde", "abxyz", "aBCD"}},
		{"ab", []string{"abcde", "abxyz"}},
		{"aB", []string{"aBCD"}},
		{"fo", []string{"foo"}},
		{"Fo", nil},
	}

	for i, testCase := range testCases {
		i, testCase := i, testCase

		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			assertSameContents(t,
				CaseSensitiveAutoCompleteFromSliceWithDefault(choices,
					defaultValue)(testCase.input),
				testCase.outputs)
		})
	}
}

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

	for i, testCase := range testCases {
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

func getValue(tb testing.TB, m *Model) string {
	tb.Helper()

	v, err := m.Value()
	if err != nil {
		tb.Fatalf("value: %v", err)
	}

	return v
}

func assertNoError(tb testing.TB, m *Model) {
	tb.Helper()

	if m.Err != nil {
		tb.Fatalf("model contains error: %v", m.Err)
	}
}

func assertSameContents(tb testing.TB, got []string, want []string) {
	tb.Helper()

	if len(want) != len(got) {
		tb.Fatalf("different length: want %#v (%d), got %#v (%d)",
			want, len(want), got, len(got))
	}

	for _, w := range want {
		found := false

		for _, g := range got {
			if g == w {
				found = true

				break
			}
		}

		if !found {
			tb.Fatalf("entry %q missing in %#v", w, got)
		}
	}
}
