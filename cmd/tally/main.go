// Command tally prints the time spent per project in a timesheet.
//
//	tally timesheet.txt
//	cat timesheet.txt | tally
//	tally --project site timesheet.txt
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/bsf-fixtures/go-cli/internal/tally"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// run is main's body with its process bound in as arguments, so every exit
// path can be reached by a test rather than only by running the binary.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("tally", flag.ContinueOnError)
	fs.SetOutput(stderr)
	project := fs.String("project", "", "print only this project's total")
	if err := fs.Parse(args); err != nil {
		// flag has already written the message and the usage to stderr.
		return 1
	}

	in := stdin
	if fs.NArg() > 0 {
		f, err := os.Open(fs.Arg(0))
		if err != nil {
			fmt.Fprintln(stderr, "tally:", err)
			return 1
		}
		defer f.Close()
		in = f
	}
	// The whole timesheet is read even when one project is named, so a
	// malformed line elsewhere is still an error rather than a quiet zero.
	entries, err := tally.Parse(in)
	if err != nil {
		fmt.Fprintln(stderr, "tally:", err)
		return 1
	}
	if *project != "" {
		t := tally.ProjectTotal(entries, *project)
		fmt.Fprintf(stdout, "%-20s %s\n", t.Project, tally.Format(t.Minutes))
		return 0
	}
	for _, t := range tally.Totals(entries) {
		fmt.Fprintf(stdout, "%-20s %s\n", t.Project, tally.Format(t.Minutes))
	}
	return 0
}
