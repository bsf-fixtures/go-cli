package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunTotalsStdinIntoAColumnLargestFirst(t *testing.T) {
	var out, errOut strings.Builder
	code := run([]string{"tally"}, strings.NewReader("2026-11-20 site 1h30m\n2026-11-20 docs 45m\n# note\n2026-11-21 site 2h\n"), &out, &errOut)
	want := "site                 3h30m\ndocs                 0h45m\n"
	if code != 0 || out.String() != want || errOut.String() != "" {
		t.Fatalf("got code %d, stdout %q, stderr %q, want code 0, stdout %q, stderr \"\"", code, out.String(), errOut.String(), want)
	}
}

func TestRunReadsTheFileArgumentInsteadOfStdin(t *testing.T) {
	path := filepath.Join(t.TempDir(), "timesheet.txt")
	if err := os.WriteFile(path, []byte("2026-11-20 site 1h30m\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errOut strings.Builder
	code := run([]string{"tally", path}, strings.NewReader("2026-11-20 ignored 9h\n"), &out, &errOut)
	want := "site                 1h30m\n"
	if code != 0 || out.String() != want || errOut.String() != "" {
		t.Fatalf("got code %d, stdout %q, stderr %q, want code 0, stdout %q, stderr \"\"", code, out.String(), errOut.String(), want)
	}
}

func TestRunReportsAMissingFileOnStderrAndExitsOne(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nope.txt")
	var out, errOut strings.Builder
	code := run([]string{"tally", path}, strings.NewReader(""), &out, &errOut)
	if code != 1 || out.String() != "" || !strings.HasPrefix(errOut.String(), "tally: open "+path+":") {
		t.Fatalf("got code %d, stdout %q, stderr %q, want code 1, empty stdout and a \"tally: open %s: ...\" diagnostic", code, out.String(), errOut.String(), path)
	}
}

func TestRunPrintsNoTotalsWhenALineIsMalformed(t *testing.T) {
	var out, errOut strings.Builder
	code := run([]string{"tally"}, strings.NewReader("2026-11-20 site 1h\n2026-11-21 site ninety\n"), &out, &errOut)
	if code != 1 || out.String() != "" || !strings.HasPrefix(errOut.String(), "tally: line 2: ") {
		t.Fatalf("got code %d, stdout %q, stderr %q, want code 1, empty stdout and a \"tally: line 2: ...\" diagnostic", code, out.String(), errOut.String())
	}
}
