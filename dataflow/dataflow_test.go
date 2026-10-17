package dataflow

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/vlumi/copal/bedrock/cobol"
	"github.com/vlumi/copal/bedrock/layout"
)

func TestAnalyze(t *testing.T) {
	p, err := cobol.ReadProgram("testdata/FLOW.cbl", nil)
	if err != nil {
		t.Fatal(err)
	}
	records, err := layout.Program(p)
	if err != nil {
		t.Fatal(err)
	}
	flow := Analyze(p, records)
	var got []string
	for _, para := range flow.Paragraphs {
		got = append(got, fmt.Sprintf("%s reads %s writes %s", para.Name, strings.Join(para.Reads, " "), strings.Join(para.Writes, " ")))
	}
	want := []string{
		"MAIN-PARA reads IN-CODE IN-AMOUNT WS-TOTAL WS-COUNT writes IN-RECORD WS-TOTAL WS-COUNT",
		"WRITE-PARA reads WS-TOTAL WS-COUNT WS-LINE writes WS-LINE-TOTAL IN-CODE OUT-RECORD IN-RECORD",
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
		if e.FromLiteral {
			from = "literal " + from
		}
		if e.ToFile {
			to = "file " + to
		}
		got = append(got, fmt.Sprintf("%s %s -> %s in %s", e.Verb, from, to, e.Paragraph))
	}
	want = []string{
		"READ file IN-FILE -> IN-RECORD in MAIN-PARA",
		"ADD IN-AMOUNT -> WS-TOTAL in MAIN-PARA",
		"ADD literal 1 -> WS-COUNT in MAIN-PARA",
		"COMPUTE WS-TOTAL -> WS-LINE-TOTAL in WRITE-PARA",
		"COMPUTE WS-COUNT -> WS-LINE-TOTAL in WRITE-PARA",
		"SET literal condition IN-DEPOSIT -> IN-CODE in WRITE-PARA",
		"WRITE WS-LINE -> OUT-RECORD in WRITE-PARA",
		"WRITE OUT-RECORD -> file OUT-FILE in WRITE-PARA",
		"REWRITE IN-RECORD -> file IN-FILE in WRITE-PARA",
	}
	if !slices.Equal(got, want) {
		t.Errorf("lineage\n%q\nwant\n%q", got, want)
	}
}

// Values are followed by bytes: a field of a group moved into a group comes
// from the same bytes of the source, which a field that redefines them
// writes too.
func TestOrigins(t *testing.T) {
	p, err := cobol.ReadProgram("testdata/ORIGINS.cbl", nil)
	if err != nil {
		t.Fatal(err)
	}
	records, err := layout.Program(p)
	if err != nil {
		t.Fatal(err)
	}
	flow := Analyze(p, records)
	for item, want := range map[string][]string{
		"WS-KIND":     {"file IN-FILE IN-KIND"},
		"WS-BODY":     {"file IN-FILE IN-BODY", "literal 5"},
		"WS-OUT-ACCT": {"file IN-FILE IN-ACCT"},
		"WS-OUT-TAG":  {"initial value of WS-OUT-TAG"},
		"WS-NOTHING":  nil,
	} {
		var got []string
		for _, o := range flow.Origins(item) {
			got = append(got, o.String())
		}
		if !slices.Equal(got, want) {
			t.Errorf("origins of %s: %q, want %q", item, got, want)
		}
	}
}

// A trace keeps the statements it follows, in program order.
func TestTrace(t *testing.T) {
	p, err := cobol.ReadProgram("testdata/ORIGINS.cbl", nil)
	if err != nil {
		t.Fatal(err)
	}
	records, err := layout.Program(p)
	if err != nil {
		t.Fatal(err)
	}
	flow := Analyze(p, records)
	for item, want := range map[string][]string{
		"WS-OUT-ACCT": {"READ IN-FILE -> IN-RECORD", "MOVE IN-ACCT -> WS-OUT-ACCT"},
		"WS-BODY":     {"READ IN-FILE -> IN-RECORD", "MOVE 5 -> IN-AMOUNT", "MOVE IN-RECORD -> WS-HOLD"},
		"WS-OUT-TAG":  nil,
	} {
		var got []string
		for _, e := range flow.Trace(item).Edges {
			got = append(got, fmt.Sprintf("%s %s -> %s", e.Verb, e.From, e.To))
		}
		if !slices.Equal(got, want) {
			t.Errorf("trace of %s: %q, want %q", item, got, want)
		}
	}
}
