// Package test contains helper functions for prompt tests.
package test

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	// UpdateGoldenFiles specifies whether the golden files should be updated.
	UpdateGoldenFiles = flag.Bool("update", false, "update the golden files")

	// Inspect prints string mismatches escaped such that the differences can be
	// inspected in detail.
	Inspect = flag.Bool("inspect", false, "inspect strings in detail")
)

// MsgsFromText generates KeyMsg events from a given text.
func MsgsFromText(text string) []tea.Msg {
	msgs := make([]tea.Msg, 0, len(text))

	for _, c := range text {
		msgs = append(msgs, KeyMsg(c))
	}

	return msgs
}

// KeyMsg returns the KeyMsg that corresponds to the given rune.
func KeyMsg(r rune) tea.Msg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
}

// Run initializes the model and applies all events.
func Run(tb testing.TB, model tea.Model, events ...tea.Msg) {
	tb.Helper()

	model.Update(model.Init())

	for _, event := range events {
		Update(tb, model, event)
	}
}

// Update applies the event to an already initialized model. If the model was not
// initialized, use Run instead.
func Update(tb testing.TB, model tea.Model, event tea.Msg) tea.Cmd {
	tb.Helper()

	var cmd tea.Cmd

	switch e := event.(type) {
	case tea.KeyType:
		_, cmd = model.Update(tea.KeyMsg{Type: e})
	default:
		_, cmd = model.Update(event)
	}

	return cmd
}

// AssertGoldenView compares the view to an exected view in an updatable golden file.
func AssertGoldenView(tb testing.TB, m tea.Model, expectedViewFile string) {
	tb.Helper()

	view := m.View()
	goldenFilePath := filepath.Join("testdata", expectedViewFile)

	if _, err := os.Stat(goldenFilePath); errors.Is(err, os.ErrNotExist) || *UpdateGoldenFiles {
		err := os.WriteFile(goldenFilePath, []byte(view), 0o664) //nolint:gosec,gomnd
		if err != nil {
			tb.Fatalf("updating golden view: %v", err)
		}

		return
	}

	goldenViewFileContent, err := os.ReadFile(goldenFilePath)
	if err != nil {
		tb.Fatalf("reading golden view: %v", err)
	}

	expectedView := string(goldenViewFileContent)

	if view != expectedView {
		comparison := "Expected:\n%s\nGot:\n%s"
		if *Inspect {
			comparison = "Expected:\n%q\nGot:\n%q"
		}

		tb.Errorf("view mismatch in %s:\n"+comparison,
			expectedViewFile, Indent(expectedView), Indent(view))
	}
}

// Indent is intended to indent views for easier comparison in test error logs.
func Indent(text string) string {
	res := make([]byte, 0, len(text))

	newLine := true

	for _, c := range text {
		if newLine && c != '\n' {
			res = append(res, []byte("    ")...)
		}

		res = append(res, byte(c))
		newLine = c == '\n'
	}

	return string(res)
}

var ansiRE = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*" +
	"(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))")

// StripANSI removes all ANSI sequences from a string.
func StripANSI(str string) string {
	return ansiRE.ReplaceAllString(str, "")
}
