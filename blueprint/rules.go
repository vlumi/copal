package blueprint

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vlumi/copal/rules"
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
// its actions, in order, as indexes into the actions.
type DecisionTable struct {
	Paragraph  string         `json:"paragraph"`
	Loop       string         `json:"loop,omitempty"`
	Source     string         `json:"source"`
	Conditions []string       `json:"conditions"`
	Actions    []string       `json:"actions"`
	Rules      []DecisionRule `json:"rules"`
}

// DecisionRule is one path through a decision table.
type DecisionRule struct {
	Entries []string `json:"entries"`
	Actions []int    `json:"actions"`
}

func newRules(program string, tables []rules.Table) Rules {
	r := Rules{Program: program}
	for _, t := range tables {
		dt := DecisionTable{Paragraph: t.Paragraph, Loop: t.Loop, Source: t.Pos.String(), Conditions: t.Conditions, Actions: t.Actions}
		for _, rule := range t.Rules {
			dt.Rules = append(dt.Rules, DecisionRule{Entries: rule.Entries, Actions: rule.Actions})
		}
		r.Tables = append(r.Tables, dt)
	}
	return r
}

func rulesMarkdown(dir string, r Rules) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s: decision tables\n\n", r.Program)
	fmt.Fprintf(&b, "The business rules of [%s](../programs/%s.md), paragraph by paragraph: each column is one path through "+
		"the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered "+
		"in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest "+
		"of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.\n", r.Program, r.Program)
	for _, t := range r.Tables {
		title := t.Paragraph
		if t.Loop != "" {
			title += ", in its loop"
		}
		fmt.Fprintf(&b, "\n## %s\n\n", title)
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
