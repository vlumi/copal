package console

import (
	"fmt"
	"strconv"
	"strings"
)

// decimal writes n, in units of the given number of decimals, with them:
// 152000 with 2 is 1520.00.
func decimal(n int64, scale int) string {
	if scale == 0 {
		return strconv.FormatInt(n, 10)
	}
	digits := fmt.Sprintf("%0*d", scale+1, abs(n))
	sign := ""
	if n < 0 {
		sign = "-"
	}
	return sign + digits[:len(digits)-scale] + "." + digits[len(digits)-scale:]
}

// money writes cents as an amount with thousands separators: -1,520.00.
func money(cents int64) string {
	whole, fraction, _ := strings.Cut(decimal(cents, 2), ".")
	return grouped(whole) + "." + fraction
}

// count writes a number with thousands separators.
func count(n int) string { return grouped(strconv.Itoa(n)) }

// recordCount writes a number of records, one in the singular.
func recordCount(n int) string {
	if n == 1 {
		return "1 record"
	}
	return count(n) + " records"
}

func grouped(digits string) string {
	sign := ""
	if rest, ok := strings.CutPrefix(digits, "-"); ok {
		sign, digits = "-", rest
	}
	var b strings.Builder
	for i, c := range digits {
		if i > 0 && (len(digits)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(c)
	}
	return sign + b.String()
}

// date writes a yyyymmdd number as yyyy-mm-dd.
func date(n int64) string {
	s := fmt.Sprintf("%08d", n)
	return s[:4] + "-" + s[4:6] + "-" + s[6:]
}
