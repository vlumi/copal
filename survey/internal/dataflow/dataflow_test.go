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

// Affects is Trace turned around: values are followed forward by bytes,
// through groups and redefinitions, to the items they reach and the bytes
// of the files they are written to.
func TestAffects(t *testing.T) {
	p, err := cobol.ReadProgram("testdata/ORIGINS.cbl", nil)
	if err != nil {
		t.Fatal(err)
	}
	records, err := layout.Program(p)
	if err != nil {
		t.Fatal(err)
	}
	flow := Analyze(p, records)
	for item, want := range map[string]string{
		"IN-KIND":   "[WS-KIND] [] [MOVE IN-RECORD -> WS-HOLD]",
		"IN-ACCT":   "[WS-BODY WS-OUT-ACCT OUT-RECORD] [{OUT-FILE 1 5}] [MOVE IN-RECORD -> WS-HOLD MOVE IN-ACCT -> WS-OUT-ACCT WRITE WS-OUT -> OUT-RECORD WRITE OUT-RECORD -> OUT-FILE]",
		"IN-AMOUNT": "[WS-BODY] [] [MOVE IN-RECORD -> WS-HOLD]",
		"WS-OUT":    "[OUT-RECORD] [{OUT-FILE 1 8}] [WRITE WS-OUT -> OUT-RECORD WRITE OUT-RECORD -> OUT-FILE]",
		"WS-HOLD":   "[] [] []",
	} {
		e := flow.Affects(item)
		var edges []string
		for _, edge := range e.Edges {
			edges = append(edges, fmt.Sprintf("%s %s -> %s", edge.Verb, edge.From, edge.To))
		}
		if got := fmt.Sprintf("%v %v %v", e.Items, e.Files, edges); got != want {
			t.Errorf("affects of %s:\n%s\nwant\n%s", item, got, want)
		}
	}
}

// Items share bytes when one is in the other or redefines it.
func TestShare(t *testing.T) {
	p, err := cobol.ReadProgram("testdata/ORIGINS.cbl", nil)
	if err != nil {
		t.Fatal(err)
	}
	records, err := layout.Program(p)
	if err != nil {
		t.Fatal(err)
	}
	flow := Analyze(p, records)
	for _, c := range []struct {
		a, b string
		want bool
	}{{"IN-RECORD", "IN-ACCT", true}, {"IN-BODY", "IN-AMOUNT", true}, {"IN-KIND", "IN-ACCT", false}, {"WS-HOLD", "IN-RECORD", false}} {
		if got := flow.Share(c.a, c.b); got != c.want {
			t.Errorf("%s and %s share bytes: %v", c.a, c.b, got)
		}
	}
}

// A run of an item's bytes is traced and followed on its own: the bytes
// of IN-BODY that IN-ACCT redefines come from the file alone, and go on to
// WS-OUT-ACCT as well as into WS-HOLD.
func TestBytes(t *testing.T) {
	p, err := cobol.ReadProgram("testdata/ORIGINS.cbl", nil)
	if err != nil {
		t.Fatal(err)
	}
	records, err := layout.Program(p)
	if err != nil {
		t.Fatal(err)
	}
	flow := Analyze(p, records)
	if got := fmt.Sprint(flow.TraceBytes("IN-BODY", 1, 5).Origins); got != "[file IN-FILE IN-ACCT]" {
		t.Errorf("origins of IN-BODY's first bytes: %s", got)
	}
	if got := fmt.Sprint(flow.AffectsBytes("IN-BODY", 1, 5).Items); got != "[WS-BODY WS-OUT-ACCT OUT-RECORD]" {
		t.Errorf("IN-BODY's first bytes reach %s", got)
	}
	if e := flow.AffectsBytes("IN-BODY", 7, 5); e.Items != nil {
		t.Errorf("bytes beyond the item reach %s", e.Items)
	}
}
