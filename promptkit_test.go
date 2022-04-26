package promptkit_test

import (
	"testing"

	"github.com/erikgeiser/promptkit"
	_ "github.com/erikgeiser/promptkit/test"
)

func TestWordWrap(t *testing.T) {
	t.Parallel()

	text := "ab cde fgh ijklmnopq rs"
	expected := "ab cde\nfgh\nijklmno\npq\nrs"
	got := promptkit.WordWrap(text, 7)

	if got != expected {
		t.Errorf("unexpected result:\nExpected:\n%s\nGot:\n%s", expected, got)
	}
}

func TestHardWrap(t *testing.T) {
	t.Parallel()

	text := "ab cde fgh ijklmnopq rs"
	expected := "ab cde \nfgh ijk\nlmnopq \nrs"
	got := promptkit.HardWrap(text, 7)

	if got != expected {
		t.Errorf("unexpected result:\nExpected:\n%s\nGot:\n%s", expected, got)
	}
}
