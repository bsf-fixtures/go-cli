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
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// run chooses the reader, prints the totals, and returns the exit status
// rather than calling os.Exit itself, so that the command's behaviour — which
// input it reads, what it writes to stderr, and when it fails — can be tested
// without starting a process.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	in := stdin
	if len(args) > 0 {
		f, err := os.Open(args[0])
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
