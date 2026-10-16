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
