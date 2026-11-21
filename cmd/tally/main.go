// Command tally prints the time spent per project in a timesheet.
//
//	tally timesheet.txt
//	cat timesheet.txt | tally
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/bsf-fixtures/go-cli/internal/tally"
)

func main() {
	os.Exit(run(os.Args, os.Stdin, os.Stdout, os.Stderr))
}

// run is main's wiring with the process edges passed in, so the choice
// between stdin and a file argument, the output column and the exit status
// can be tested without starting a process. args is os.Args, program name
// included; the result is the exit status.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	in := stdin
	if len(args) > 1 {
		f, err := os.Open(args[1])
		if err != nil {
			fmt.Fprintln(stderr, "tally:", err)
			return 1
		}
		defer f.Close()
		in = f
	}
	entries, err := tally.Parse(in)
	if err != nil {
		fmt.Fprintln(stderr, "tally:", err)
		return 1
	}
	for _, t := range tally.Totals(entries) {
		fmt.Fprintf(stdout, "%-20s %s\n", t.Project, tally.Format(t.Minutes))
	}
	return 0
}
