package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunTotalsStdinIntoAFixedWidthColumn(t *testing.T) {
	var stdout, stderr strings.Builder
	code := run(nil, strings.NewReader("2026-11-20 site 1h30m\n2026-11-20 docs 45m\n# note\n2026-11-21 site 2h\n"), &stdout, &stderr)
	want := "site                 3h30m\ndocs                 0h45m\n"
	if code != 0 || stdout.String() != want || stderr.String() != "" {
		t.Fatalf("got code %d, stdout %q, stderr %q, want code 0, stdout %q, stderr \"\"", code, stdout.String(), stderr.String(), want)
	}
}

func TestRunReadsTheFileNamedByItsArgumentInsteadOfStdin(t *testing.T) {
	path := filepath.Join(t.TempDir(), "timesheet.txt")
	if err := os.WriteFile(path, []byte("2026-11-20 site 1h30m\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr strings.Builder
	code := run([]string{path}, strings.NewReader("2026-11-20 ignored 9h\n"), &stdout, &stderr)
	want := "site                 1h30m\n"
	if code != 0 || stdout.String() != want || stderr.String() != "" {
		t.Fatalf("got code %d, stdout %q, stderr %q, want code 0, stdout %q, stderr \"\"", code, stdout.String(), stderr.String(), want)
	}
}

func TestRunReportsAMalformedLineOnStderrAndTotalsNothing(t *testing.T) {
	var stdout, stderr strings.Builder
	code := run(nil, strings.NewReader("2026-11-20 site 1h\n2026-11-21 site ninety\n"), &stdout, &stderr)
	if code != 1 || stdout.String() != "" || !strings.HasPrefix(stderr.String(), "tally: line 2: ") {
		t.Fatalf("got code %d, stdout %q, stderr %q, want code 1, no stdout, a stderr line starting \"tally: line 2: \"", code, stdout.String(), stderr.String())
	}
}

func TestRunReportsAMissingFileOnStderrAndExitsOne(t *testing.T) {
	path := filepath.Join(t.TempDir(), "absent.txt")
	var stdout, stderr strings.Builder
	code := run([]string{path}, strings.NewReader(""), &stdout, &stderr)
	if code != 1 || stdout.String() != "" || !strings.HasPrefix(stderr.String(), "tally: open "+path+": ") {
		t.Fatalf("got code %d, stdout %q, stderr %q, want code 1, no stdout, a stderr line starting \"tally: open %s: \"", code, stdout.String(), stderr.String(), path)
	}
}
