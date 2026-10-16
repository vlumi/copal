package dataflow

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/vlumi/copal/cobol"
)

func TestAnalyze(t *testing.T) {
	p, err := cobol.ReadProgram("testdata/FLOW.cbl", nil)
	if err != nil {
		t.Fatal(err)
	}
	flow, err := Analyze(p)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, para := range flow.Paragraphs {
		got = append(got, fmt.Sprintf("%s reads %s writes %s", para.Name, strings.Join(para.Reads, " "), strings.Join(para.Writes, " ")))
	}
	want := []string{
		"MAIN-PARA reads IN-CODE IN-AMOUNT WS-TOTAL WS-COUNT writes IN-RECORD WS-TOTAL WS-COUNT",
		"WRITE-PARA reads WS-TOTAL WS-COUNT WS-LINE writes WS-LINE-TOTAL IN-CODE OUT-RECORD",
	}
	if !slices.Equal(got, want) {
		t.Errorf("paragraphs\n%q\nwant\n%q", got, want)
	}
	got = nil
	for _, e := range flow.Lineage {
		from, to := e.From, e.To
		if e.FromFile {
			from = "file " + from
		}
		if e.ToFile {
			to = "file " + to
		}
		got = append(got, fmt.Sprintf("%s %s -> %s in %s", e.Verb, from, to, e.Paragraph))
	}
	want = []string{
		"READ file IN-FILE -> IN-RECORD in MAIN-PARA",
		"ADD IN-AMOUNT -> WS-TOTAL in MAIN-PARA",
		"COMPUTE WS-TOTAL -> WS-LINE-TOTAL in WRITE-PARA",
		"COMPUTE WS-COUNT -> WS-LINE-TOTAL in WRITE-PARA",
		"WRITE WS-LINE -> OUT-RECORD in WRITE-PARA",
		"WRITE OUT-RECORD -> file OUT-FILE in WRITE-PARA",
	}
	if !slices.Equal(got, want) {
		t.Errorf("lineage\n%q\nwant\n%q", got, want)
	}
}
