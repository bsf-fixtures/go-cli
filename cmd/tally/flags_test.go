package main

import (
	"bytes"
	"strings"
	"testing"
)

// The flag-parse failure path is the one exit in run that main_test.go does
// not reach. It matters for the same reason a malformed line does: a command
// line the user got wrong must stop, not total whatever it managed to read.

func TestRunRejectsAnUnknownFlagAndTotalsNothing(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"--nope", "--project", "site"}, strings.NewReader(twoProjects), &stdout, &stderr)
	if code != 1 || stdout.String() != "" {
		t.Fatalf("got code %d stdout %q, want code 1 and empty stdout", code, stdout.String())
	}
	if !strings.Contains(stderr.String(), "-nope") {
		t.Fatalf("got stderr %q, want it to name the flag it did not recognise", stderr.String())
	}
}

func TestRunRejectsAFlagMissingItsValueAndTotalsNothing(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"--project"}, strings.NewReader(twoProjects), &stdout, &stderr)
	if code != 1 || stdout.String() != "" {
		t.Fatalf("got code %d stdout %q, want code 1 and empty stdout", code, stdout.String())
	}
	if !strings.Contains(stderr.String(), "-project") {
		t.Fatalf("got stderr %q, want it to name the flag that needs an argument", stderr.String())
	}
}
