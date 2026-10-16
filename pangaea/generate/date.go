package generate

import (
	"fmt"
	"time"
)

// Date is a calendar date, as the bank's records hold it: yyyymmdd.
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

func (d Date) String() string { return d.t.Format("2006-01-02") }

// IsZero reports whether d is the zero Date, no date at all.
func (d Date) IsZero() bool { return d.t.IsZero() }

// Day is the day of the month.
func (d Date) Day() int { return d.t.Day() }

// AddDays is the date n days later.
func (d Date) AddDays(n int) Date { return Date{d.t.AddDate(0, 0, n)} }

// BusinessDay reports whether the bank is open: Monday to Friday.
func (d Date) BusinessDay() bool {
	w := d.t.Weekday()
	return w != time.Saturday && w != time.Sunday
}

// sameMonth reports whether two dates are in the same month.
func (d Date) sameMonth(e Date) bool {
	return d.t.Year() == e.t.Year() && d.t.Month() == e.t.Month()
}

// FirstBusinessDayOfMonth reports whether no business day of the same
// month comes before d.
func (d Date) FirstBusinessDayOfMonth() bool {
	for e := d.AddDays(-1); e.sameMonth(d); e = e.AddDays(-1) {
		if e.BusinessDay() {
			return false
		}
	}
	return d.BusinessDay()
}

// Payday reports whether salaries are paid on d: the 25th, or the last
// business day before it when the 25th is not one.
func (d Date) Payday() bool {
	payday := NewDate(d.t.Year(), d.t.Month(), 25)
	for !payday.BusinessDay() {
		payday = payday.AddDays(-1)
	}
	return d == payday
}
