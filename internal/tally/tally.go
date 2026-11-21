// Package tally reads a timesheet and totals the minutes per project.
//
// A timesheet is one entry per line: "2026-11-20 project 1h30m". Blank lines
// and lines starting with "#" are skipped. A malformed line is an error that
// names its line number, never a silent zero: a timesheet that loses an hour
// without saying so is worse than one that refuses to total.
package tally

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

// Entry is one line of a timesheet.
type Entry struct {
	Day     string
	Project string
	Minutes int
}

// Total is the minutes spent on one project.
type Total struct {
	Project string
	Minutes int
}

// Parse reads every entry, or returns the first malformed line.
func Parse(r io.Reader) ([]Entry, error) {
	var out []Entry
	sc := bufio.NewScanner(r)
	n := 0
	for sc.Scan() {
		n++
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		f := strings.Fields(line)
		if len(f) != 3 {
			return nil, fmt.Errorf("line %d: want \"DAY PROJECT DURATION\", got %q", n, line)
		}
		if _, err := time.Parse("2006-01-02", f[0]); err != nil {
			return nil, fmt.Errorf("line %d: %q is not a YYYY-MM-DD day", n, f[0])
		}
		d, err := time.ParseDuration(f[2])
		if err != nil || d <= 0 {
			return nil, fmt.Errorf("line %d: %q is not a positive duration like 1h30m", n, f[2])
		}
		out = append(out, Entry{Day: f[0], Project: f[1], Minutes: int(d.Minutes())})
	}
	return out, sc.Err()
}

// Totals sums minutes per project, largest first, ties by name.
func Totals(entries []Entry) []Total {
	by := map[string]int{}
	for _, e := range entries {
		by[e.Project] += e.Minutes
	}
	out := make([]Total, 0, len(by))
	for p, m := range by {
		out = append(out, Total{Project: p, Minutes: m})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Minutes != out[j].Minutes {
			return out[i].Minutes > out[j].Minutes
		}
		return out[i].Project < out[j].Project
	})
	return out
}

// Format renders minutes as "1h05m".
func Format(minutes int) string {
	return fmt.Sprintf("%dh%02dm", minutes/60, minutes%60)
}
