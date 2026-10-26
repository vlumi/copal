// Command bankdate is the bank's BANKDATE rewritten in Go: it reads the
// bank control record from the file DD_BANKIN names and writes the next
// business day's, the next day from Monday to Friday, to the one
// DD_BANKOUT names. It returns 0, or 16 when a file cannot be opened or
// there is no control record, with the same lines displayed, as the COBOL
// does. It is here to show the survey's conformance check on a
// replacement.
package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

const recordLength = 80

// epoch is day 1 of COBOL's INTEGER-OF-DATE, a Monday, so a day number
// modulo 7 is 1 on Mondays, 6 on Saturdays and 0 on Sundays.
var epoch = time.Date(1601, 1, 1, 0, 0, 0, 0, time.UTC)

func main() {
	os.Exit(run())
}

func run() int {
	in, inStatus := openInput(os.Getenv("DD_BANKIN"))
	out, outStatus := openOutput(os.Getenv("DD_BANKOUT"))
	if inStatus != "00" || outStatus != "00" {
		fmt.Printf("BANKDATE: CANNOT OPEN FILES %s %s\n", inStatus, outStatus)
		return 16
	}
	defer in.Close()
	defer out.Close()
	record := make([]byte, recordLength)
	if _, err := io.ReadFull(in, record); err != nil {
		fmt.Println("BANKDATE: NO BANK CONTROL RECORD")
		return 16
	}
	copy(record, nextBusinessDay(string(record[:8])))
	if _, err := out.Write(record); err != nil {
		return 16
	}
	return 0
}

// openInput opens a file for reading, with the file status the open gives.
func openInput(path string) (*os.File, string) {
	f, err := os.Open(path)
	switch {
	case os.IsNotExist(err):
		return nil, "35"
	case err != nil:
		return nil, "30"
	}
	return f, "00"
}

// openOutput creates a file for writing, with the file status the open
// gives.
func openOutput(path string) (*os.File, string) {
	f, err := os.Create(path)
	if err != nil {
		return nil, "30"
	}
	return f, "00"
}

// nextBusinessDay moves a date, YYYYMMDD, to the next day that is not a
// weekend, by day numbers as the COBOL counts them: a date that is not
// one counts as day 0.
func nextBusinessDay(date string) string {
	day := 0
	if t, err := time.Parse("20060102", date); err == nil {
		day = int((t.Unix()-epoch.Unix())/86400) + 1
	}
	day++
	for weekday := day % 7; weekday == 6 || weekday == 0; weekday = day % 7 {
		day++
	}
	return epoch.AddDate(0, 0, day-1).Format("20060102")
}
