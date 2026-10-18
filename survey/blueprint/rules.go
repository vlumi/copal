package blueprint

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/vlumi/copal/survey/internal/rules"
)

// Rules is the blueprint of a program's business rules: its decision
// tables, one for each paragraph, and each inline loop, that decides
// anything.
type Rules struct {
	Program string          `json:"program"`
	Tables  []DecisionTable `json:"tables"`
}

// DecisionTable is a paragraph's paths through its decisions: for each
// rule, the outcome of each condition, Y, N or empty when not tested, and
// its actions, in order, as indexes into the actions. ConditionRefs are
// the data items each condition decides on, a condition name counting as
// the field it is declared on.
type DecisionTable struct {
	Paragraph     string         `json:"paragraph"`
	Loop          string         `json:"loop,omitempty"`
	Source        string         `json:"source"`
	Conditions    []string       `json:"conditions"`
	ConditionRefs [][]string     `json:"conditionRefs"`
	Actions       []string       `json:"actions"`
	Rules         []DecisionRule `json:"rules"`
}

// DecisionRule is one path through a decision table.
type DecisionRule struct {
	Entries []string `json:"entries"`
	Actions []int    `json:"actions"`
}

// newRules converts a program's decision tables, resolving the condition
// names their conditions refer to with the fields they are declared on.
func newRules(program string, tables []rules.Table, records []Record) Rules {
	fields := conditionFields(records)
	r := Rules{Program: program}
	for _, t := range tables {
		dt := DecisionTable{Paragraph: t.Paragraph, Loop: t.Loop, Source: t.Pos.String(), Conditions: t.Conditions, Actions: t.Actions}
		for _, refs := range t.Refs {
			var items []string
			for _, ref := range refs {
				if field, ok := fields[ref]; ok {
					ref = field
				}
				addOnce(&items, ref)
			}
			dt.ConditionRefs = append(dt.ConditionRefs, items)
		}
		for _, rule := range t.Rules {
			dt.Rules = append(dt.Rules, DecisionRule{Entries: rule.Entries, Actions: rule.Actions})
		}
		r.Tables = append(r.Tables, dt)
	}
	return r
}

// Decision is a decision table with conditions on a data item: its
// paragraph, and loop if it is one's, and those conditions.
type Decision struct {
	Paragraph  string   `json:"paragraph"`
	Loop       string   `json:"loop,omitempty"`
	Source     string   `json:"source"`
	Conditions []string `json:"conditions"`
}

// Decisions lists the program's decision tables that decide on a data
// item, with the conditions on it.
func (p Program) Decisions(item string) []Decision {
	var decisions []Decision
	for _, t := range p.Rules.Tables {
		d := Decision{Paragraph: t.Paragraph, Loop: t.Loop, Source: t.Source}
		for i, refs := range t.ConditionRefs {
			if slices.Contains(refs, item) {
				d.Conditions = append(d.Conditions, t.Conditions[i])
			}
		}
		if len(d.Conditions) > 0 {
			decisions = append(decisions, d)
		}
	}
	return decisions
}

// conditionFields maps each condition name the records declare to its
// field.
func conditionFields(records []Record) map[string]string {
	fields := map[string]string{}
	for _, r := range records {
		walk(r.Field, func(f Field) {
			for _, c := range f.Conditions {
				fields[c.Name] = f.Name
			}
		})
	}
	return fields
}

// Title is the heading of a decision table: its paragraph, and for an
// inline loop's, that it is in its loop.
func (t DecisionTable) Title() string {
	if t.Loop != "" {
		return t.Paragraph + ", in its loop"
	}
	return t.Paragraph
}

func rulesMarkdown(dir string, r Rules) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s: decision tables\n\n", r.Program)
	fmt.Fprintf(&b, "The business rules of [%s](../programs/%s.md), paragraph by paragraph: each column is one path through "+
		"the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered "+
		"in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest "+
		"of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.\n", r.Program, r.Program)
	for _, t := range r.Tables {
		fmt.Fprintf(&b, "\n## %s\n\n", t.Title())
		if t.Loop != "" {
			fmt.Fprintf(&b, "Each pass of `%s`.\n\n", t.Loop)
		}
		header := []string{""}
		align := []string{"---"}
		for i := range t.Rules {
			header = append(header, strconv.Itoa(i+1))
			align = append(align, ":-:")
		}
		b.WriteString(row(header...))
		b.WriteString(row(align...))
		for i, c := range t.Conditions {
			cells := []string{"`" + c + "`"}
			for _, rule := range t.Rules {
				cells = append(cells, rule.Entries[i])
			}
			b.WriteString(row(cells...))
		}
		for i, a := range t.Actions {
			cells := []string{"→ `" + a + "`"}
			for _, rule := range t.Rules {
				var orders []string
				for order, action := range rule.Actions {
					if action == i {
						orders = append(orders, strconv.Itoa(order+1))
					}
				}
				cells = append(cells, strings.Join(orders, ", "))
			}
			b.WriteString(row(cells...))
		}
		fmt.Fprintf(&b, "\nSource: %s.\n", sourceLink(dir, t.Source, true))
	}
	return b.String()
}
