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
	var in io.Reader = os.Stdin
	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "tally:", err)
			os.Exit(1)
		}
		defer f.Close()
		in = f
	}
	entries, err := tally.Parse(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "tally:", err)
		os.Exit(1)
	}
	for _, t := range tally.Totals(entries) {
		fmt.Printf("%-20s %s\n", t.Project, tally.Format(t.Minutes))
	}
}
