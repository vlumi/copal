package generate

import (
	"slices"
	"testing"
	"time"
)

func TestEraDates(t *testing.T) {
	for _, c := range []struct {
		date Date
		era  string
	}{
		{NewDate(2027, time.March, 1), "R090301"},
		{NewDate(2019, time.May, 1), "R010501"},
		{NewDate(2019, time.April, 30), "H310430"},
		{NewDate(1989, time.January, 8), "H010108"},
		{NewDate(1989, time.January, 7), "S640107"},
		{NewDate(1962, time.April, 1), "S370401"},
	} {
		if got := c.date.Era(); got != c.era {
			t.Errorf("%s in its era is %s, want %s", c.date, got, c.era)
		}
		back, err := DateFromEra(c.era)
		if err != nil || back != c.date {
			t.Errorf("%s reads as %s, %v; want %s", c.era, back, err, c.date)
		}
	}
	for _, s := range []string{"R000501", "H310501", "S650101", "X010101", "R09030", "R090230"} {
		if d, err := DateFromEra(s); err == nil {
			t.Errorf("%s read as %s", s, d)
		}
	}
}

func TestCalendar(t *testing.T) {
	wednesday := NewDate(2026, time.April, 1)
	if wednesday.Weekday() != time.Wednesday || !wednesday.BusinessDay() || wednesday.Number() != 20260401 || wednesday.Stamp() != "20260401" || wednesday.YearMonth() != 202604 {
		t.Errorf("%s: %s, business %v, %d %s %d", wednesday, wednesday.Weekday(), wednesday.BusinessDay(), wednesday.Number(), wednesday.Stamp(), wednesday.YearMonth())
	}
	friday := NewDate(2026, time.April, 3)
	if next := friday.NextBusinessDay(); next.String() != "2026-04-06" || friday.AddDays(1).BusinessDayOnOrAfter() != next {
		t.Errorf("the business day after %s is %s", friday, next)
	}
	if !NewDate(2026, time.April, 30).LastBusinessDayOfMonth() || !NewDate(2026, time.May, 29).LastBusinessDayOfMonth() || NewDate(2026, time.May, 28).LastBusinessDayOfMonth() {
		t.Error("the last business days of April and May 2026 are the 30th and the 29th")
	}
	if d := NewDate(2026, time.January, 31).AddMonths(1); d.String() != "2026-02-28" {
		t.Errorf("a month after January 31 is %s", d)
	}
	if d := NewDate(2007, time.March, 15).AddYears(20); d.String() != "2027-03-15" {
		t.Errorf("twenty years after 2007-03-15 is %s", d)
	}
	if n := MonthsBetween(NewDate(1990, time.July, 15), NewDate(2027, time.March, 1)); n != 439 {
		t.Errorf("%d months from 1990-07-15 to 2027-03-01, want 439", n)
	}
	if n := DaysBetween(NewDate(2026, time.December, 31), NewDate(2027, time.January, 2)); n != 2 {
		t.Errorf("%d days over the new year, want 2", n)
	}
	if d, err := DateFromNumber(20260230); err == nil {
		t.Errorf("20260230 read as %s", d)
	}
}

// The holidays of the fiscal year the insurer runs: those of a fixed
// day, of a Monday, the equinoxes, the Wednesday kept for a Sunday's
// Constitution Day, the Tuesday between two holidays, and the year's end.
func TestHolidays(t *testing.T) {
	var got []string
	for d := NewDate(2026, time.April, 1); d.Before(NewDate(2027, time.April, 1)); d = d.AddDays(1) {
		if d.Holiday() && d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
			got = append(got, d.String())
		}
	}
	want := []string{"2026-04-29", "2026-05-04", "2026-05-05", "2026-05-06", "2026-07-20", "2026-08-11",
		"2026-09-21", "2026-09-22", "2026-09-23", "2026-10-12", "2026-11-03", "2026-11-23", "2026-12-31",
		"2027-01-01", "2027-01-11", "2027-02-11", "2027-02-23", "2027-03-22"}
	if !slices.Equal(got, want) {
		t.Errorf("holidays on weekdays\n%q\nwant\n%q", got, want)
	}
	if next := NewDate(2026, time.December, 30).NextBusinessDay(); next.String() != "2027-01-04" {
		t.Errorf("the business day after December 30 is %s", next)
	}
	if !NewDate(2026, time.December, 30).LastBusinessDayOfMonth() {
		t.Error("the last business day of December 2026 is the 30th")
	}
}
