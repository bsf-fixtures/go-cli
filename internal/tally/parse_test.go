package tally

import (
	"strings"
	"testing"
)

func TestParseRejectsALineThatIsNotThreeFields(t *testing.T) {
	_, err := Parse(strings.NewReader("2026-11-20 site 1h\n2026-11-21 site\n"))
	if err == nil || !strings.Contains(err.Error(), "line 2") || !strings.Contains(err.Error(), "DAY PROJECT DURATION") {
		t.Fatalf("want line 2 named with the shape it wanted, got %v", err)
	}
}

func TestParseRejectsADayThatIsNotYYYYMMDD(t *testing.T) {
	_, err := Parse(strings.NewReader("20-11-2026 site 1h\n"))
	if err == nil || !strings.Contains(err.Error(), "line 1") || !strings.Contains(err.Error(), "YYYY-MM-DD") {
		t.Fatalf("want line 1 named as a bad day, got %v", err)
	}
}

func TestParseRejectsANonPositiveDurationInsteadOfTotallingZero(t *testing.T) {
	// Zero and negative both reach the same guard by a different route: one
	// parses to nothing, the other to time taken back.
	for _, d := range []string{"0m", "-1h"} {
		got, err := Parse(strings.NewReader("2026-11-20 site " + d + "\n"))
		if err == nil || !strings.Contains(err.Error(), "line 1") {
			t.Fatalf("%q: want an error naming line 1, got %+v, %v", d, got, err)
		}
		if got != nil {
			t.Fatalf("%q: want no entries, got %+v", d, got)
		}
	}
}
