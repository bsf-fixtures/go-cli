package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const twoProjects = "2026-11-20 site 1h30m\n2026-11-20 docs 45m\n# note\n2026-11-21 site 2h\n"

func TestRunWithProjectFlagPrintsOnlyThatProject(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"--project", "site"}, strings.NewReader(twoProjects), &stdout, &stderr)
	want := "site                 3h30m\n"
	if code != 0 || stdout.String() != want {
		t.Fatalf("got code %d stdout %q stderr %q, want code 0 stdout %q", code, stdout.String(), stderr.String(), want)
	}
	if strings.Contains(stdout.String(), "docs") {
		t.Fatalf("got %q, want no mention of another project", stdout.String())
	}
}

func TestRunPrintsZeroForAProjectWithNoEntries(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"--project", "nope"}, strings.NewReader(twoProjects), &stdout, &stderr)
	want := "nope                 0h00m\n"
	if code != 0 || stdout.String() != want {
		t.Fatalf("got code %d stdout %q stderr %q, want code 0 stdout %q", code, stdout.String(), stderr.String(), want)
	}
}

func TestRunWithoutProjectFlagPrintsEveryProjectLargestFirst(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run(nil, strings.NewReader(twoProjects), &stdout, &stderr)
	want := "site                 3h30m\ndocs                 0h45m\n"
	if code != 0 || stdout.String() != want {
		t.Fatalf("got code %d stdout %q stderr %q, want code 0 stdout %q", code, stdout.String(), stderr.String(), want)
	}
}

func TestRunReportsAMalformedLineOnStderrAndExitsOne(t *testing.T) {
	// The two cases vary only in whether a project is named: naming one must
	// not suppress a bad line belonging to a different project.
	for _, args := range [][]string{nil, {"--project", "site"}} {
		var stdout, stderr bytes.Buffer
		code := run(args, strings.NewReader("2026-11-20 site 1h\n2026-11-21 docs ninety\n"), &stdout, &stderr)
		if code != 1 || stdout.String() != "" {
			t.Fatalf("args %v: got code %d stdout %q, want code 1 and empty stdout", args, code, stdout.String())
		}
		if !strings.Contains(stderr.String(), "tally:") || !strings.Contains(stderr.String(), "line 2") {
			t.Fatalf("args %v: got stderr %q, want it to name tally: and line 2", args, stderr.String())
		}
	}
}

func TestRunReadsTheFileNamedByThePositionalArgument(t *testing.T) {
	path := filepath.Join(t.TempDir(), "timesheet.txt")
	if err := os.WriteFile(path, []byte("2026-11-20 docs 45m\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{path}, strings.NewReader(twoProjects), &stdout, &stderr)
	want := "docs                 0h45m\n"
	if code != 0 || stdout.String() != want {
		t.Fatalf("got code %d stdout %q stderr %q, want code 0 stdout %q", code, stdout.String(), stderr.String(), want)
	}

	stdout.Reset()
	stderr.Reset()
	code = run([]string{"/nope.txt"}, strings.NewReader(twoProjects), &stdout, &stderr)
	if code != 1 || !strings.Contains(stderr.String(), "tally: open") {
		t.Fatalf("got code %d stderr %q, want code 1 and a \"tally: open\" line", code, stderr.String())
	}
}
