package blueprint

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/vlumi/copal/survey/internal/diagram"
)

// StateMachine is how a code field moves between its states, the
// condition names declared on it, as the programs that use its entity
// move it.
type StateMachine struct {
	States      []State      `json:"states"`
	Transitions []Transition `json:"transitions,omitempty"`
	// Copies are the statements that move another field's value into it,
	// which are not transitions: the state they give is whatever the other
	// field holds.
	Copies []Copy `json:"copies,omitempty"`
}

// Copy is a statement of a program that moves another data item's value
// into a code field.
type Copy struct {
	Program   string `json:"program"`
	From      string `json:"from"`
	Paragraph string `json:"paragraph"`
	Source    string `json:"source"`
}

// State is a condition name of a code field and its values.
type State struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

// AnyState is where a transition starts that sets a state without testing
// the field first.
const AnyState = "any"

// Transition is a paragraph of a program setting a field to a state, To,
// on a path through it that holds the field in the state From, or in any
// state when the path does not establish one. Source is where the
// statement that sets the state is.
type Transition struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Program   string    `json:"program"`
	Paragraph string    `json:"paragraph"`
	Source    string    `json:"source"`
	When      []Outcome `json:"when,omitempty"`
}

// Outcome is a condition a rule tests, and whether the rule finds it true.
type Outcome struct {
	Condition string `json:"condition"`
	Holds     bool   `json:"holds"`
}

// codeField is a code field as a program knows it, its condition names
// mapped to the copybook's, which name the states.
type codeField struct {
	program *Program
	name    string
	states  map[string]string // the program's condition name to the state
	values  map[string]string // a value, as written, to its state
}

func newCodeField(p *Program, field, copybook Field) codeField {
	c := codeField{program: p, name: field.Name, states: map[string]string{}, values: map[string]string{}}
	for i, cond := range field.Conditions {
		state := copybook.Conditions[i].Name
		c.states[cond.Name] = state
		for _, v := range cond.Values {
			c.values[v] = state
		}
	}
	return c
}

// newStateMachine starts a code field's state machine with its states.
func newStateMachine(copybook Field) *StateMachine {
	m := &StateMachine{}
	for _, c := range copybook.Conditions {
		m.States = append(m.States, State{Name: c.Name, Values: c.Values})
	}
	return m
}

// add adds the transitions and copies of the field in one program.
//
// A decision table's rule holds the field in state A where it tests A's
// condition name, or FIELD = 'value' with A's value, and finds it true. Each
// of its actions that sets the field, by SET B TO TRUE or MOVE 'value' TO
// FIELD with B's value, is a transition from A to B. A rule that holds the
// field in no state starts from the states known where its paragraph
// starts: those the rules that perform the paragraph hold the field in,
// and theirs in turn, or any state where some way in knows none. A
// paragraph without a decision table sets states the same way, as its
// statements in the lineage say.
func (m *StateMachine) add(c codeField) {
	p := c.program
	entries := c.entries()
	decided := map[string]bool{}
	for _, t := range p.Rules.Tables {
		decided[t.Paragraph] = true
		for a, action := range t.Actions {
			to, ok := c.set(action)
			if !ok {
				continue
			}
			source := c.statementSource(t.Paragraph, to)
			for _, r := range taking(t, a) {
				when := c.reasons(t, r, a)
				for _, e := range c.startingFrom(t, r, entries[t.Paragraph]) {
					for _, f := range e.states {
						m.addTransition(Transition{From: f, To: to, Program: p.Name, Paragraph: t.Paragraph, Source: source, When: outcomes(e.when, when)})
					}
				}
			}
		}
	}
	for _, e := range p.Lineage {
		if e.To != c.name || e.ToFile {
			continue
		}
		if !e.FromLiteral {
			m.Copies = append(m.Copies, Copy{Program: p.Name, From: e.From, Paragraph: e.Paragraph, Source: e.Source})
			continue
		}
		if to, ok := c.literalState(e.From); ok && !decided[e.Paragraph] {
			for _, in := range startingFrom([]string{AnyState}, nil, entries[e.Paragraph]) {
				for _, f := range in.states {
					m.addTransition(Transition{From: f, To: to, Program: p.Name, Paragraph: e.Paragraph, Source: e.Source, When: in.when})
				}
			}
		}
	}
}

// entry is one way into a paragraph: the states a rule that performs it
// holds the field in, or those known where the rule's own paragraph
// starts, and the outcomes that decide the PERFORM.
type entry struct {
	states []string
	when   []Outcome
}

// startingFrom gives the states a rule starts from: those it holds the
// field in, or, where it holds none, each way into its paragraph, less
// the states the rule finds the field not in.
func (c codeField) startingFrom(t DecisionTable, r DecisionRule, entries []entry) []entry {
	return startingFrom(c.heldIn(t, r), c.excludedIn(t, r), entries)
}

func startingFrom(held, excluded []string, entries []entry) []entry {
	if len(held) != 1 || held[0] != AnyState || len(entries) == 0 {
		return []entry{{states: held}}
	}
	var from []entry
	for _, e := range entries {
		states := slices.DeleteFunc(slices.Clone(e.states), func(state string) bool { return slices.Contains(excluded, state) })
		if len(states) > 0 {
			from = append(from, entry{states: states, when: e.when})
		}
	}
	return from
}

// outcomes joins the outcomes a way into a paragraph is decided by and
// those a rule in it is, each condition once.
func outcomes(entry, own []Outcome) []Outcome {
	joined := slices.Clone(entry)
	for _, o := range own {
		if !slices.Contains(joined, o) {
			joined = append(joined, o)
		}
	}
	return joined
}

// entries lists the ways into each paragraph the program's rules perform,
// from the performing rules, and through them from the rules that perform
// their paragraphs, until the states known no longer change. A paragraph
// reached by other than a PERFORM, by GO TO or fall-through, is entered
// from any state.
func (c codeField) entries() map[string][]entry {
	p := c.program
	reached := map[string]bool{}
	for _, para := range p.Paragraphs {
		if len(para.ReachedBy) > 0 {
			reached[para.Name] = true
		}
	}
	entries := map[string][]entry{}
	for range 8 {
		next := map[string][]entry{}
		for _, t := range p.Rules.Tables {
			for a, action := range t.Actions {
				target, ok := performed(action)
				if !ok {
					continue
				}
				for _, r := range taking(t, a) {
					var ways []entry
					if reached[target] {
						ways = []entry{{states: []string{AnyState}}}
					} else {
						ways = c.startingFrom(t, r, entries[t.Paragraph])
					}
					when := c.reasons(t, r, a)
					for _, in := range ways {
						e := entry{states: in.states, when: outcomes(in.when, when)}
						if !slices.ContainsFunc(next[target], func(have entry) bool { return reflect.DeepEqual(have, e) }) {
							next[target] = append(next[target], e)
						}
					}
				}
			}
		}
		if reflect.DeepEqual(next, entries) {
			break
		}
		entries = next
	}
	return entries
}

// taking are the rules of a table that take an action.
func taking(t DecisionTable, action int) []DecisionRule {
	var rules []DecisionRule
	for _, r := range t.Rules {
		if slices.Contains(r.Actions, action) {
			rules = append(rules, r)
		}
	}
	return rules
}

// performed is the paragraph an action performs, if it is a PERFORM of
// one: the first after PERFORM, of a THRU range too.
func performed(action string) (string, bool) {
	rest, ok := strings.CutPrefix(action, "PERFORM ")
	if !ok {
		return "", false
	}
	name, _, _ := strings.Cut(rest, " ")
	switch name {
	case "", "UNTIL", "VARYING", "WITH", "TEST":
		return "", false
	}
	return name, true
}

func (m *StateMachine) addTransition(t Transition) {
	if !slices.ContainsFunc(m.Transitions, func(u Transition) bool { return reflect.DeepEqual(t, u) }) {
		m.Transitions = append(m.Transitions, t)
	}
}

// reasons are the outcomes of the conditions that decide an action of a
// rule, other than those on the field itself: why the rule takes it. A
// condition decides the action where the rules part on it, one that
// found the conditions before it alike finding it otherwise and not
// taking the action; one every rule finds the same way, one some rule
// finds otherwise and takes the action all the same, or one tested after
// the rules have parted, says nothing. Of a conjunction the rule finds
// true, with the field in it, the other parts are reasons, each true.
func (c codeField) reasons(t DecisionTable, r DecisionRule, action int) []Outcome {
	var when []Outcome
	for i, condition := range t.Conditions {
		if r.Entries[i] == "" || !decides(t, r, i, action) {
			continue
		}
		switch {
		case !slices.Contains(t.ConditionRefs[i], c.name):
			when = append(when, Outcome{Condition: condition, Holds: r.Entries[i] == "Y"})
		case r.Entries[i] == "Y":
			for _, part := range conjuncts(condition) {
				if _, held := c.state(part); !held {
					when = append(when, Outcome{Condition: part, Holds: true})
				}
			}
		}
	}
	return when
}

// decides reports whether a condition decides an action for a rule. The
// rules that part from it there, finding the conditions before alike, or
// not testing them, and this one otherwise, include one that does not
// take the action, and none that takes it and goes on as the rule does.
func decides(t DecisionTable, r DecisionRule, condition, action int) bool {
	parted := false
	for _, other := range t.Rules {
		if other.Entries[condition] == "" || other.Entries[condition] == r.Entries[condition] || !agree(r.Entries[:condition], other.Entries[:condition]) {
			continue
		}
		if slices.Contains(other.Actions, action) && slices.Equal(r.Entries[condition+1:], other.Entries[condition+1:]) {
			return false
		}
		parted = parted || !slices.Contains(other.Actions, action)
	}
	return parted
}

// agree reports whether two rules find the conditions they both test the
// same way.
func agree(a, b []string) bool {
	for i := range a {
		if a[i] != "" && b[i] != "" && a[i] != b[i] {
			return false
		}
	}
	return true
}

// heldIn lists the states a rule holds the field in, or any state: those
// of the conditions on the field it finds true, alone or as a part of a
// conjunction.
func (c codeField) heldIn(t DecisionTable, r DecisionRule) []string {
	var states []string
	for i, condition := range t.Conditions {
		if r.Entries[i] != "Y" || !slices.Contains(t.ConditionRefs[i], c.name) {
			continue
		}
		for _, part := range conjuncts(condition) {
			if state, ok := c.state(part); ok {
				states = append(states, state)
			}
		}
	}
	if len(states) == 0 {
		return []string{AnyState}
	}
	return states
}

// excludedIn are the states a rule finds the field not in: those of the
// conditions it finds false that are one state's.
func (c codeField) excludedIn(t DecisionTable, r DecisionRule) []string {
	var states []string
	for i, condition := range t.Conditions {
		if r.Entries[i] != "N" || !slices.Contains(t.ConditionRefs[i], c.name) {
			continue
		}
		if state, ok := c.state(condition); ok {
			states = append(states, state)
		}
	}
	return states
}

// state is the state a condition holds the field in: its condition name,
// or FIELD = 'value' with the value of one.
func (c codeField) state(condition string) (string, bool) {
	if state, ok := c.states[condition]; ok {
		return state, true
	}
	if value, ok := strings.CutPrefix(condition, c.name+" = "); ok && c.values[value] != "" {
		return c.values[value], true
	}
	return "", false
}

// conjuncts splits a condition into the parts an AND joins at its top
// level, A AND B AND C, which all hold when it does; a condition with OR
// at its top level is one part, as it does not say which held.
func conjuncts(condition string) []string {
	var parts []string
	depth, start := 0, 0
	var quote byte
	for i := 0; i < len(condition); i++ {
		switch {
		case quote != 0:
			if condition[i] == quote {
				quote = 0
			}
		case condition[i] == '\'' || condition[i] == '"':
			quote = condition[i]
		case condition[i] == '(':
			depth++
		case condition[i] == ')':
			depth--
		case depth > 0:
		case strings.HasPrefix(condition[i:], " OR "):
			return []string{condition}
		case strings.HasPrefix(condition[i:], " AND "):
			parts = append(parts, condition[start:i])
			start = i + len(" AND ")
		}
	}
	return append(parts, condition[start:])
}

// set is the state an action sets the field to, if it sets one.
func (c codeField) set(action string) (string, bool) {
	if name, ok := strings.CutPrefix(action, "SET "); ok {
		name, ok = strings.CutSuffix(name, " TO TRUE")
		state, known := c.states[name]
		return state, ok && known
	}
	if value, ok := strings.CutPrefix(action, "MOVE "); ok {
		value, ok = strings.CutSuffix(value, " TO "+c.name)
		state, known := c.values[value]
		return state, ok && known
	}
	return "", false
}

// literalState is the state a literal of the lineage sets: "condition
// NAME" for a SET, or a value moved.
func (c codeField) literalState(literal string) (string, bool) {
	if name, ok := strings.CutPrefix(literal, "condition "); ok {
		state, known := c.states[name]
		return state, known
	}
	state, known := c.values[literal]
	return state, known
}

// statementSource is where a paragraph sets the field to a state.
func (c codeField) statementSource(paragraph, state string) string {
	for _, e := range c.program.Lineage {
		if s, ok := c.literalState(e.From); ok && e.FromLiteral && e.To == c.name && e.Paragraph == paragraph && s == state {
			return e.Source
		}
	}
	return ""
}

// Label is a transition's program and paragraph, as an edge shows it.
func (t Transition) Label() string {
	return fmt.Sprintf("%s %s", t.Program, t.Paragraph)
}

// StateGraph draws a state machine: its states, an "any" node for the
// transitions that start from any state, and an edge for each transition
// between two of them, labeled with every program and paragraph that
// makes it.
func StateGraph(m StateMachine) diagram.Graph {
	g := diagram.Graph{Legend: []diagram.Key{
		diagram.NodeKey(diagram.Node{}, "state, its condition name and values"),
		diagram.NodeKey(diagram.Node{Dashed: true}, "any state: set without a test"),
		diagram.EdgeKey(diagram.Solid, "set by, program and paragraph"),
	}}
	if slices.ContainsFunc(m.Transitions, func(t Transition) bool { return t.From == AnyState }) {
		g.Nodes = append(g.Nodes, diagram.Node{ID: AnyState, Label: "any state", Dashed: true})
	}
	for _, s := range m.States {
		g.Nodes = append(g.Nodes, diagram.Node{ID: s.Name, Label: s.Name + "\n" + strings.Join(s.Values, ", ")})
	}
	labels := map[[2]string][]string{}
	var pairs [][2]string
	for _, t := range m.Transitions {
		pair := [2]string{t.From, t.To}
		if _, ok := labels[pair]; !ok {
			pairs = append(pairs, pair)
		}
		if !slices.Contains(labels[pair], t.Label()) {
			labels[pair] = append(labels[pair], t.Label())
		}
	}
	for _, pair := range pairs {
		g.Edges = append(g.Edges, diagram.Edge{From: pair[0], To: pair[1], Label: strings.Join(labels[pair], "\n")})
	}
	return g
}
