package generate

import "time"

// Holiday reports whether d is a day the insurer and the banks are
// closed on a weekday's account: a national holiday, the day one that
// fell on a Sunday is kept instead, a day between two holidays, or one
// of the year's end and start, December 31 to January 3.
func (d Date) Holiday() bool {
	if national(d) || yearEnd(d) {
		return true
	}
	if national(d.AddDays(-1)) && national(d.AddDays(1)) {
		return true
	}
	for kept := d.AddDays(-1); national(kept); kept = kept.AddDays(-1) {
		if kept.Weekday() == time.Sunday {
			return true
		}
	}
	return false
}

// yearEnd reports whether d is one of the days the banks close over the
// new year.
func yearEnd(d Date) bool {
	return d.Month() == time.December && d.Day() == 31 || d.Month() == time.January && d.Day() <= 3
}

// national reports whether d is a national holiday as the law has had
// them since 2020: those of a fixed day, those of a Monday of their
// month, and the two equinoxes.
func national(d Date) bool {
	monday := 0
	if d.Weekday() == time.Monday {
		monday = (d.Day() + 6) / 7
	}
	switch d.Month() {
	case time.January:
		return d.Day() == 1 || monday == 2
	case time.February:
		return d.Day() == 11 || d.Day() == 23
	case time.March:
		return d.Day() == equinox(d.Year(), 20.8431)
	case time.April:
		return d.Day() == 29
	case time.May:
		return d.Day() >= 3 && d.Day() <= 5
	case time.July:
		return monday == 3
	case time.August:
		return d.Day() == 11
	case time.September:
		return monday == 3 || d.Day() == equinox(d.Year(), 23.2488)
	case time.October:
		return monday == 2
	case time.November:
		return d.Day() == 3 || d.Day() == 23
	}
	return false
}

// equinox is the day of March or September an equinox falls on in a
// year from 1980 to 2099, from the day it fell on in 1980.
func equinox(year int, in1980 float64) int {
	years := year - 1980
	return int(in1980+0.242194*float64(years)) - years/4
}
