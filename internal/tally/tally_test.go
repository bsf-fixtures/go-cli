package tally

import (
	"strings"
	"testing"
)

func TestParseSkipsBlankAndCommentLines(t *testing.T) {
	got, err := Parse(strings.NewReader("# week 47\n\n2026-11-20 site 1h30m\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Minutes != 90 || got[0].Project != "site" {
		t.Fatalf("got %+v", got)
	}
}

func TestParseNamesTheMalformedLineInsteadOfCountingZero(t *testing.T) {
	_, err := Parse(strings.NewReader("2026-11-20 site 1h\n2026-11-21 site ninety\n"))
	if err == nil || !strings.Contains(err.Error(), "line 2") {
		t.Fatalf("want an error naming line 2, got %v", err)
	}
}

func TestTotalsAreLargestFirstWithTiesByName(t *testing.T) {
	got := Totals([]Entry{{Project: "b", Minutes: 30}, {Project: "a", Minutes: 30}, {Project: "c", Minutes: 90}, {Project: "a", Minutes: 15}})
	want := []Total{{"c", 90}, {"a", 45}, {"b", 30}}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %+v, want %+v", got, want)
		}
	}
}

func TestFormatPadsMinutes(t *testing.T) {
	if Format(65) != "1h05m" {
		t.Fatalf("got %s", Format(65))
	}
}
