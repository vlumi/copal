package console

import (
	"strconv"
	"strings"

	"golang.org/x/text/encoding/japanese"
)

// yen writes an amount in whole yen with thousands separators: ¥1,520,000.
func yen(n int64) string {
	if n < 0 {
		return "-¥" + grouped(strconv.FormatInt(-n, 10))
	}
	return "¥" + grouped(strconv.FormatInt(n, 10))
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

// utf8 reads Shift_JIS as UTF-8, the half-width katakana of the names and
// the kanji of the data alike; bytes it cannot read come out as the
// replacement character.
func utf8(sjis string) string {
	text, err := japanese.ShiftJIS.NewDecoder().String(sjis)
	if err != nil {
		return strings.ToValidUTF8(sjis, "�")
	}
	return text
}
