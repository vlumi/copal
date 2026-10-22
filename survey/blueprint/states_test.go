package blueprint

import (
	"reflect"
	"slices"
	"strings"
	"testing"
)

// A rule that finds the field in a state and sets another is a transition
// between them; one that sets a state without finding the field in one,
// or a paragraph that decides nothing, starts from any state; a value
// moved from another field is a copy, not a transition. The program's
// condition names are mapped to the copybook's, which name the states.
func TestStateMachine(t *testing.T) {
	copybook := Field{Name: "ST", Conditions: []Condition{{Name: "S-OPEN", Values: []string{"'O'"}}, {Name: "S-SHUT", Values: []string{"'X'"}}}}
	field := Field{Name: "XX-ST", Conditions: []Condition{{Name: "XX-OPEN", Values: []string{"'O'"}}, {Name: "XX-SHUT", Values: []string{"'X'"}}}}
	p := &Program{Name: "P",
		Rules: Rules{Tables: []DecisionTable{
			{Paragraph: "SHUT-IT", Conditions: []string{"XX-OPEN", "XX-ST = 'X'"}, ConditionRefs: [][]string{{"XX-ST"}, {"XX-ST"}},
				Actions: []string{"SET XX-SHUT TO TRUE", "MOVE 'O' TO XX-ST", "DISPLAY 'NO'"},
				Rules: []DecisionRule{{Entries: []string{"Y", ""}, Actions: []int{0}}, {Entries: []string{"N", "Y"}, Actions: []int{1}},
					{Entries: []string{"N", "N"}, Actions: []int{2}}}},
			{Paragraph: "LOG-IT", Conditions: []string{"WS-LOG"}, ConditionRefs: [][]string{{"WS-LOG"}},
				Actions: []string{"SET XX-OPEN TO TRUE"}, Rules: []DecisionRule{{Entries: []string{"Y"}, Actions: []int{0}}, {Entries: []string{"N"}}}},
		}},
		Lineage: []Lineage{
			{From: "condition XX-SHUT", FromLiteral: true, To: "XX-ST", Paragraph: "SHUT-IT", Verb: "SET", Source: "P.cbl:10"},
			{From: "'O'", FromLiteral: true, To: "XX-ST", Paragraph: "SHUT-IT", Verb: "MOVE", Source: "P.cbl:11"},
			{From: "condition XX-OPEN", FromLiteral: true, To: "XX-ST", Paragraph: "LOG-IT", Verb: "SET", Source: "P.cbl:20"},
			{From: "'O'", FromLiteral: true, To: "XX-ST", Paragraph: "START-IT", Verb: "MOVE", Source: "P.cbl:30"},
			{From: "OTHER-ST", To: "XX-ST", Paragraph: "COPY-IT", Verb: "MOVE", Source: "P.cbl:40"},
		},
	}
	m := newStateMachine(copybook)
	m.add(newCodeField(p, field, copybook))
	want := []Transition{
		{From: "S-OPEN", To: "S-SHUT", Program: "P", Paragraph: "SHUT-IT", Source: "P.cbl:10"},
		{From: "S-SHUT", To: "S-OPEN", Program: "P", Paragraph: "SHUT-IT", Source: "P.cbl:11"},
		{From: AnyState, To: "S-OPEN", Program: "P", Paragraph: "LOG-IT", Source: "P.cbl:20", When: []Outcome{{Condition: "WS-LOG", Holds: true}}},
		{From: AnyState, To: "S-OPEN", Program: "P", Paragraph: "START-IT", Source: "P.cbl:30"},
	}
	if !reflect.DeepEqual(m.Transitions, want) {
		t.Errorf("transitions\n%+v\nwant\n%+v", m.Transitions, want)
	}
	if len(m.Copies) != 1 || m.Copies[0].From != "OTHER-ST" {
		t.Errorf("copies %+v", m.Copies)
	}
	g := StateGraph(*m)
	if len(g.Nodes) != 3 || len(g.Edges) != 3 || g.Edges[2].Label != "P LOG-IT\nP START-IT" {
		t.Errorf("graph %+v", g)
	}
	if svg, err := g.SVG(); err != nil || !strings.Contains(string(svg), ">P SHUT-IT</text>") {
		t.Errorf("svg %s, %v", svg, err)
	}
}

// A rule that sets a state without finding the field in one starts from
// the states the rules performing its paragraph hold the field in, and
// those two performs up, with the outcomes that decided each PERFORM; a
// state the rule finds the field not in is left out. Of its own conditions
// only those the rules part on are reasons: one found the same way by a
// rule that takes the action all the same, with the same conditions after
// it, is none, so the rules that take an action whatever it says make one
// transition.
func TestStateMachineEntries(t *testing.T) {
	copybook := Field{Name: "ST", Conditions: []Condition{{Name: "S-OPEN", Values: []string{"'O'"}}, {Name: "S-SHUT", Values: []string{"'X'"}}, {Name: "S-GONE", Values: []string{"'G'"}}}}
	p := &Program{Name: "P",
		Paragraphs: []Paragraph{{Name: "MAIN"}, {Name: "CHECK"}, {Name: "SHUT"}, {Name: "LOG", ReachedBy: []string{"MAIN"}}},
		Rules: Rules{Tables: []DecisionTable{
			{Paragraph: "MAIN", Conditions: []string{"S-OPEN", "S-SHUT"}, ConditionRefs: [][]string{{"ST"}, {"ST"}},
				Actions: []string{"PERFORM CHECK", "PERFORM LOG"},
				Rules:   []DecisionRule{{Entries: []string{"Y", ""}, Actions: []int{0}}, {Entries: []string{"N", "Y"}, Actions: []int{0, 1}}, {Entries: []string{"N", "N"}}}},
			{Paragraph: "CHECK", Conditions: []string{"WS-DUE", "S-SHUT"}, ConditionRefs: [][]string{{"WS-DUE"}, {"ST"}},
				Actions: []string{"PERFORM SHUT"},
				Rules:   []DecisionRule{{Entries: []string{"Y", "N"}, Actions: []int{0}}, {Entries: []string{"Y", "Y"}}, {Entries: []string{"N", ""}}}},
			{Paragraph: "SHUT", Conditions: []string{"WS-FORCE", "WS-LATE"}, ConditionRefs: [][]string{{"WS-FORCE"}, {"WS-LATE"}},
				Actions: []string{"SET S-GONE TO TRUE", "DISPLAY 'LATE'"},
				Rules:   []DecisionRule{{Entries: []string{"Y", "Y"}, Actions: []int{0, 1}}, {Entries: []string{"Y", "N"}, Actions: []int{0}}, {Entries: []string{"N", ""}}}},
			{Paragraph: "LOG", Conditions: []string{"WS-LOG"}, ConditionRefs: [][]string{{"WS-LOG"}},
				Actions: []string{"SET S-OPEN TO TRUE"}, Rules: []DecisionRule{{Entries: []string{"Y"}, Actions: []int{0}}, {Entries: []string{"N"}}}},
		}},
		Lineage: []Lineage{
			{From: "condition S-GONE", FromLiteral: true, To: "ST", Paragraph: "SHUT", Verb: "SET", Source: "P.cbl:10"},
			{From: "condition S-OPEN", FromLiteral: true, To: "ST", Paragraph: "LOG", Verb: "SET", Source: "P.cbl:20"},
		},
	}
	m := newStateMachine(copybook)
	m.add(newCodeField(p, copybook, copybook))
	want := []Transition{
		{From: "S-OPEN", To: "S-GONE", Program: "P", Paragraph: "SHUT", Source: "P.cbl:10", When: []Outcome{{Condition: "WS-DUE", Holds: true}, {Condition: "WS-FORCE", Holds: true}}},
		{From: AnyState, To: "S-OPEN", Program: "P", Paragraph: "LOG", Source: "P.cbl:20", When: []Outcome{{Condition: "WS-LOG", Holds: true}}},
	}
	if !reflect.DeepEqual(m.Transitions, want) {
		t.Errorf("transitions\n%+v\nwant\n%+v", m.Transitions, want)
	}
}

func TestConjuncts(t *testing.T) {
	for condition, want := range map[string][]string{
		"HOLDING-CASE AND DC-OPEN AND DE-CHARGEBACK": {"HOLDING-CASE", "DC-OPEN", "DE-CHARGEBACK"},
		"A AND (B OR C)":      {"A", "(B OR C)"},
		"A OR B AND C":        {"A OR B AND C"},
		"A = 'X AND Y' AND B": {"A = 'X AND Y'", "B"},
		"AM-CLOSED":           {"AM-CLOSED"},
	} {
		if got := conjuncts(condition); !slices.Equal(got, want) {
			t.Errorf("conjuncts(%q) = %q, want %q", condition, got, want)
		}
	}
}
