package promptkit_test

import (
	"testing"

	"github.com/erikgeiser/promptkit"
	"github.com/erikgeiser/promptkit/test"
)

func TestWordWrap(t *testing.T) {
	t.Parallel()

	text := "ab cde fgh ijklmnopq rs"
	expected := "ab cde\nfgh\nijklmno\npq\nrs"
	assertEqual(t, expected, promptkit.WordWrap(text, 7))
}

func TestHardWrap(t *testing.T) {
	t.Parallel()

	text := "ab cde fgh ijklmnopq rs"
	expected := "ab cde \nfgh ijk\nlmnopq \nrs"
	assertEqual(t, expected, promptkit.HardWrap(text, 7))
}

func TestTruncate(t *testing.T) {
	t.Parallel()

	text := "0123456789\n0123\n0123456789\n"
	expected := "012345\n0123\n012345\n"
	assertEqual(t, expected, promptkit.Truncate(text, 6))
}

func assertEqual(tb testing.TB, expected string, got string) {
	tb.Helper()

	if expected == got {
		return
	}

	comparison := "Expected:\n%s\nGot:\n%s"
	if *test.Inspect {
		comparison = "Expected:\n%q\nGot:\n%q"
	}

	tb.Errorf("unexpected result:\n"+comparison, expected, got)
}
