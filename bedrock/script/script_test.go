package script

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"testing"
)

const night = `#!/bin/sh
# night.sh: one night.
#
# usage: night.sh DATE

DATE="$1"
DATA="$LCS_DATA"
WORK="$DATA/work/$DATE"
mkdir -p "$WORK"

# STEP LOAD: the day's file
# into the work area
export DD_IN="$DATA/inbound/IN.$DATE.txt"
export DD_OUT="$WORK/loaded.txt"
"$LCS_BIN/load" "$DATE" > "$DATA/reports/load.out" 2>&1
rc=$?
echo "LOAD load $rc" >> "$DATA/logs/night.log"
if [ "$rc" -ge 8 ]; then exit "$rc"; fi

# STEP BYKEY: in key order
sort -s -t "|" -k1.3,1.8 "$WORK/loaded.txt" > "$WORK/sorted.txt"
rc=$?
if [ "$rc" -ne 0 ]; then exit 16; fi

# STEP MORE: the rest of the night
"$LCS_SCRIPTS/more.sh" "$DATE"

mv "$WORK/sorted.txt" "$DATA/master/keyed.txt"
`

func TestParse(t *testing.T) {
	s, err := Parse("night", strings.NewReader(night))
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(s.Comment, []string{"night.sh: one night.", "", "usage: night.sh DATE"}) {
		t.Errorf("comment %q", s.Comment)
	}
	var got []string
	for _, st := range s.Steps {
		got = append(got, fmt.Sprintf("%s %s %s %v %v %q %q %v %v %v %d", st.Name, st.Kind, st.Program, st.Args, st.Input, st.Output, st.Stdout, st.Keys, st.StopsIf, st.Comment, len(st.Bindings)))
	}
	want := []string{
		`LOAD program load [$DATE] [] "" "$DATA/reports/load.out" [] &{-ge 8} [the day's file into the work area] 2`,
		`BYKEY sort sort [] [$DATA/work/$DATE/loaded.txt] "$DATA/work/$DATE/sorted.txt" "" [{3 8}] &{-ne 0} [in key order] 0`,
		`MORE script more [$DATE] [] "" "" [] <nil> [the rest of the night] 0`,
		`MV move mv [] [$DATA/work/$DATE/sorted.txt] "$DATA/master/keyed.txt" "" [] <nil> [] 0`,
	}
	if !slices.Equal(got, want) {
		t.Errorf("steps\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	if in := s.Steps[0].Bindings["IN"]; in != "$DATA/inbound/IN.$DATE.txt" {
		t.Errorf("DD_IN bound to %q", in)
	}
	if !maps.Equal(s.Parameters, map[string]string{"DATE": "1", "DATA": "LCS_DATA"}) {
		t.Errorf("parameters %v", s.Parameters)
	}
	values := map[string]string{"1": "20261019", "LCS_DATA": "/data", "LCS_BIN": "/bin"}
	value := func(parameter string) string { return values[parameter] }
	for text, want := range map[string]string{
		s.Steps[0].Bindings["IN"]: "/data/inbound/IN.20261019.txt",
		s.Steps[0].Command:        "/bin/load",
		"${DATA}/x.$UNSET":        "/data/x.",
	} {
		if got := s.Expand(text, value); got != want {
			t.Errorf("%s expands to %q, want %q", text, got, want)
		}
	}
}

func TestParseErrors(t *testing.T) {
	for name, source := range map[string]string{
		"a test of something else":  "if [ -f x ]; then exit 1; fi\n",
		"a test before any step":    "if [ \"$rc\" -ge 8 ]; then exit 8; fi\n",
		"a sort by fields":          "sort -k2 in > out\n",
		"a sort without its output": "sort in\n",
		"a quote not closed":        "echo \"x\n",
	} {
		if _, err := Parse("bad", strings.NewReader(source)); err == nil {
			t.Errorf("%s: parsed without error", name)
		}
	}
}

func TestHolds(t *testing.T) {
	ge8, ne0 := Test{Op: "-ge", Value: 8}, Test{Op: "-ne", Value: 0}
	if !ge8.Holds(8) || ge8.Holds(4) || !ne0.Holds(2) || ne0.Holds(0) {
		t.Error("a test of a return code does not hold as it should")
	}
}
