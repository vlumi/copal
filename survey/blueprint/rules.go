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

// conditionNames resolves the condition names of a program's records: the
// field each is declared on, and its values.
func conditionNames(records []Record) rules.ConditionNames {
	names := map[string]Condition{}
	fields := map[string]string{}
	for _, r := range records {
		walk(r.Field, func(f Field) {
			for _, c := range f.Conditions {
				names[c.Name], fields[c.Name] = c, f.Name
			}
		})
	}
	return func(name string) (string, []string, bool) {
		c, ok := names[name]
		return fields[name], c.Values, ok
	}
}

// initialValues are the values a program's working storage starts with:
// of each elementary field of one VALUE, in no table and no redefinition,
// by a name no other field of the program has.
func initialValues(records []Record) map[string]string {
	values := map[string]string{}
	count := map[string]int{}
	var visit func(f Field, plain bool)
	visit = func(f Field, plain bool) {
		count[f.Name]++
		plain = plain && f.Occurs == 0 && f.Redefines == ""
		if plain && len(f.Fields) == 0 && len(f.Values) == 1 {
			values[f.Name] = f.Values[0]
		}
		for _, child := range f.Fields {
			visit(child, plain)
		}
	}
	for _, r := range records {
		visit(r.Field, r.Section == "working-storage")
	}
	for name := range values {
		if count[name] > 1 {
			delete(values, name)
		}
	}
	return values
}

// textFields finds the elementary alphanumeric fields each group of a
// program's records holds, in no table, by names no other field of the
// program has: those a MOVE of SPACES to the group blanks.
func textFields(records []Record) func(group string) []string {
	count := map[string]int{}
	for _, r := range records {
		walk(r.Field, func(f Field) { count[f.Name]++ })
	}
	fields := map[string][]string{}
	var visit func(f Field) []string
	visit = func(f Field) []string {
		var held []string
		for _, child := range f.Fields {
			held = append(held, visit(child)...)
		}
		if len(f.Fields) > 0 && count[f.Name] == 1 {
			fields[f.Name] = held
		}
		switch {
		case f.Occurs > 0:
			return nil
		case len(f.Fields) == 0 && f.Category == "alphanumeric" && count[f.Name] == 1:
			return []string{f.Name}
		}
		return held
	}
	for _, r := range records {
		visit(r.Field)
	}
	return func(group string) []string { return fields[group] }
}

// copies reports whether a MOVE from one field of a program's records to
// another keeps its value whole: two elementary fields, by names no other
// field has, laid out alike.
func copies(records []Record) func(from, to string) bool {
	count := map[string]int{}
	fields := map[string]Field{}
	for _, r := range records {
		walk(r.Field, func(f Field) {
			count[f.Name]++
			fields[f.Name] = f
		})
	}
	return func(from, to string) bool {
		a, b := fields[from], fields[to]
		return count[from] == 1 && count[to] == 1 && len(a.Fields) == 0 && len(b.Fields) == 0 &&
			a.Length == b.Length && a.Category == b.Category && a.Usage == b.Usage &&
			a.Digits == b.Digits && a.Scale == b.Scale && a.Signed == b.Signed && a.Occurs == 0 && b.Occurs == 0
	}
}

// digitFields finds, for a field of a program's records of unsigned
// digits in display, its digits and the other fields of digits over its
// bytes: those of a group it is in, of one it redefines, or of one that
// redefines either, record and all, in no table and by names no other
// field has.
func digitFields(records []Record) func(field string) (int, []rules.DigitField) {
	type place struct {
		record string
		start  int
		field  Field
	}
	count := map[string]int{}
	var places []place
	base := map[string]string{}
	var visit func(record string, f Field, offset int, inTable bool)
	visit = func(record string, f Field, offset int, inTable bool) {
		count[f.Name]++
		inTable = inTable || f.Occurs > 0
		if !inTable && len(f.Fields) == 0 && f.Category == "numeric" && f.Usage == "DISPLAY" && !f.Signed && f.Scale == 0 && f.Digits == f.Length {
			places = append(places, place{record, f.Start + offset, f})
		}
		for _, child := range f.Fields {
			visit(record, child, offset, inTable)
		}
	}
	for _, r := range records {
		record := r.Name
		if b, ok := base[r.Redefines]; ok {
			record = b
		}
		base[r.Name] = record
		visit(record, r.Field, 0, false)
	}
	return func(field string) (int, []rules.DigitField) {
		i := slices.IndexFunc(places, func(p place) bool { return p.field.Name == field })
		if i < 0 || count[field] != 1 {
			return 0, nil
		}
		at := places[i]
		var over []rules.DigitField
		for _, p := range places {
			if p.record == at.record && p.field.Name != field && count[p.field.Name] == 1 &&
				p.start >= at.start && p.start+p.field.Length <= at.start+at.field.Length {
				over = append(over, rules.DigitField{Name: p.field.Name, Offset: p.start - at.start, Length: p.field.Length})
			}
		}
		return at.field.Length, over
	}
}
