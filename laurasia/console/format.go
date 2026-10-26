package console

import (
	"fmt"
	"strconv"
	"strings"
)

// money writes cents as an amount with thousands separators: -1,520.00.
func money(cents int64) string {
	sign := ""
	if cents < 0 {
		sign, cents = "-", -cents
	}
	return fmt.Sprintf("%s%s.%02d", sign, grouped(strconv.FormatInt(cents/100, 10)), cents%100)
}

// count writes a number with thousands separators.
func count(n int) string { return grouped(strconv.Itoa(n)) }

// size writes a file's size in bytes, or in kilobytes from ten of them.
func size(n int64) string {
	if n < 10<<10 {
		return count(int(n)) + " bytes"
	}
	return count(int(n>>10)) + " kB"
}

func grouped(digits string) string {
	var b strings.Builder
	for i, c := range digits {
		if i > 0 && (len(digits)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(c)
	}
	return b.String()
}
