package generate

import (
	"fmt"
	"time"
)

// Date is a calendar date, as the processor's files hold and name it:
// yyyymmdd.
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

// Number is the date as yyyymmdd.
func (d Date) Number() int64 {
	return int64(d.t.Year())*10000 + int64(d.t.Month())*100 + int64(d.t.Day())
}

// Stamp is the date as the files are named: yyyymmdd.
func (d Date) Stamp() string { return fmt.Sprintf("%08d", d.Number()) }

func (d Date) String() string { return d.t.Format("2006-01-02") }

// IsZero reports whether d is the zero Date, no date at all.
func (d Date) IsZero() bool { return d.t.IsZero() }

// Day is the day of the month.
func (d Date) Day() int { return d.t.Day() }

// Weekday is the day of the week.
func (d Date) Weekday() time.Weekday { return d.t.Weekday() }

// AddDays is the date n days later.
func (d Date) AddDays(n int) Date { return Date{d.t.AddDate(0, 0, n)} }

// After reports whether d is later than e.
func (d Date) After(e Date) bool { return d.t.After(e.t) }

// BusinessDay reports whether the networks clear on d: Monday to Friday.
// The switch authorizes every day, and so does the billing cycle run.
func (d Date) BusinessDay() bool {
	w := d.t.Weekday()
	return w != time.Saturday && w != time.Sunday
}

// PreviousBusinessDay is the last business day before d.
func (d Date) PreviousBusinessDay() Date {
	for d = d.AddDays(-1); !d.BusinessDay(); d = d.AddDays(-1) {
	}
	return d
}
