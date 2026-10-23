package generate

import (
	"fmt"
	"time"
)

// Date is a calendar date, as the insurer's files hold it: in the Western
// calendar, yyyymmdd, in what the sales system and the banks send and in
// the business date on the command line, and in the era, Gyymmdd, in the
// master and the branches' requests.
type Date struct{ t time.Time }

// NewDate is the given day.
func NewDate(year int, month time.Month, day int) Date {
	return Date{time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

// DateFromNumber reads a yyyymmdd number.
func DateFromNumber(n int64) (Date, error) {
	d := NewDate(int(n/10000), time.Month(n/100%100), int(n%100))
	if d.Number() != n {
		return Date{}, fmt.Errorf("%d is not a date", n)
	}
	return d, nil
}

// era is a Japanese era: its letter and the day it began.
type era struct {
	letter byte
	start  Date
}

// eras are the eras the master's dates can be in, newest first, as the
// programs' GENGO table has them.
var eras = []era{
	{'R', NewDate(2019, time.May, 1)},
	{'H', NewDate(1989, time.January, 8)},
	{'S', NewDate(1926, time.December, 25)},
	{'T', NewDate(1912, time.July, 30)},
	{'M', NewDate(1868, time.September, 8)},
}

// DateFromEra reads an era date, Gyymmdd: the era's letter, its year
// counting the year it began as 1, then the month and the day.
func DateFromEra(s string) (Date, error) {
	if len(s) != 7 {
		return Date{}, fmt.Errorf("%q is not an era date", s)
	}
	var year, month, day int
	if _, err := fmt.Sscanf(s[1:], "%2d%2d%2d", &year, &month, &day); err != nil {
		return Date{}, fmt.Errorf("%q is not an era date", s)
	}
	for i, e := range eras {
		if e.letter != s[0] {
			continue
		}
		d := NewDate(e.start.Year()+year-1, time.Month(month), day)
		if d.Before(e.start) || i > 0 && !d.Before(eras[i-1].start) || d.Month() != time.Month(month) || d.Day() != day {
			return Date{}, fmt.Errorf("%s is not a date of its era", s)
		}
		return d, nil
	}
	return Date{}, fmt.Errorf("%q is not an era date", s)
}

// Era is the date in its era, Gyymmdd, as the master holds it.
func (d Date) Era() string {
	for _, e := range eras {
		if !d.Before(e.start) {
			return fmt.Sprintf("%c%02d%02d%02d", e.letter, d.Year()-e.start.Year()+1, d.Month(), d.Day())
		}
	}
	return ""
}

// Number is the date as yyyymmdd.
func (d Date) Number() int64 {
	return int64(d.Year())*10000 + int64(d.Month())*100 + int64(d.Day())
}

// Stamp is the date as the files are named: yyyymmdd.
func (d Date) Stamp() string { return fmt.Sprintf("%08d", d.Number()) }

// YearMonth is the date's month as yyyymm, as the paid-to month is held.
func (d Date) YearMonth() int64 { return d.Number() / 100 }

func (d Date) String() string { return d.t.Format("2006-01-02") }

// IsZero reports whether d is the zero Date, no date at all.
func (d Date) IsZero() bool { return d.t.IsZero() }

// Year, Month and Day are the date's parts.
func (d Date) Year() int         { return d.t.Year() }
func (d Date) Month() time.Month { return d.t.Month() }
func (d Date) Day() int          { return d.t.Day() }

// Weekday is the day of the week.
func (d Date) Weekday() time.Weekday { return d.t.Weekday() }

// AddDays is the date n days later.
func (d Date) AddDays(n int) Date { return Date{d.t.AddDate(0, 0, n)} }

// AddMonths is the date n months later, on the same day of the month, or
// the last day of the month when that one has fewer days.
func (d Date) AddMonths(n int) Date {
	first := time.Date(d.Year(), d.Month()+time.Month(n), 1, 0, 0, 0, 0, time.UTC)
	last := first.AddDate(0, 1, -1).Day()
	return Date{time.Date(first.Year(), first.Month(), min(d.Day(), last), 0, 0, 0, 0, time.UTC)}
}

// AddYears is the date n years later, as the programs add a term to a
// contract date: the year moved on, the month and day kept.
func (d Date) AddYears(n int) Date { return d.AddMonths(12 * n) }

// After and Before compare dates.
func (d Date) After(e Date) bool  { return d.t.After(e.t) }
func (d Date) Before(e Date) bool { return d.t.Before(e.t) }

// MonthsBetween is the number of whole months from d to e, counting a
// month only once e's day of the month has come, as the insurance age
// counts them.
func MonthsBetween(d, e Date) int {
	months := (e.Year()-d.Year())*12 + int(e.Month()) - int(d.Month())
	if e.Day() < d.Day() {
		months--
	}
	return max(months, 0)
}

// DaysBetween is the number of days from d to e.
func DaysBetween(d, e Date) int {
	return int(e.t.Sub(d.t).Hours() / 24)
}

// BusinessDay reports whether the insurer runs on d: Monday to Friday.
func (d Date) BusinessDay() bool {
	w := d.t.Weekday()
	return w != time.Saturday && w != time.Sunday
}

// NextBusinessDay is the first business day after d.
func (d Date) NextBusinessDay() Date {
	for d = d.AddDays(1); !d.BusinessDay(); d = d.AddDays(1) {
	}
	return d
}

// BusinessDayOnOrAfter is d when it is a business day, else the next.
func (d Date) BusinessDayOnOrAfter() Date {
	if d.BusinessDay() {
		return d
	}
	return d.NextBusinessDay()
}

// LastBusinessDayOfMonth reports whether d is the last business day of
// its month.
func (d Date) LastBusinessDayOfMonth() bool {
	return d.BusinessDay() && d.NextBusinessDay().Month() != d.Month()
}

// MonthStart is the first day of d's month.
func (d Date) MonthStart() Date { return NewDate(d.Year(), d.Month(), 1) }
