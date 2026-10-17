package blueprint

// LayoutRow is one line of a record's layout as a table shows it: a field,
// or a condition name under its field, with its place in the record's
// tree drawn as a prefix of branches, so that a reader sees which fields
// are in which group without comparing level numbers.
type LayoutRow struct {
	Level  int
	Name   string
	Prefix string // the branches before the name, such as "│ ├ "
	Field  *Field // nil for a condition name
	Values []string
	Source string
}

// Tree branches: a double one for a field that redefines the one above
// it, so that the views of the same bytes stand out. The blanks are
// no-break spaces, so that a Markdown table cell keeps them.
const (
	branchMid       = "├\u00a0"
	branchLast      = "└\u00a0"
	branchRedefines = "╞\u00a0"
	branchRedefLast = "╘\u00a0"
	branchDown      = "│\u00a0"
	branchNone      = "\u00a0\u00a0"
)

// LayoutRows lays a record out as rows, each field's condition names
// before its fields, the way the source declares them.
func LayoutRows(record Field) []LayoutRow {
	var rows []LayoutRow
	var walk func(f *Field, ancestors string, last bool, depth int)
	walk = func(f *Field, ancestors string, last bool, depth int) {
		prefix := ancestors
		if depth > 0 {
			prefix += branch(last, f.Redefines != "")
		}
		rows = append(rows, LayoutRow{Level: f.Level, Name: f.Name, Prefix: prefix, Field: f, Source: f.Source})
		below := ancestors
		if depth > 0 {
			below += continuation(last)
		}
		items := len(f.Conditions) + len(f.Fields)
		for i, c := range f.Conditions {
			rows = append(rows, LayoutRow{Level: 88, Name: c.Name, Prefix: below + branch(i == items-1, false), Values: c.Values, Source: c.Source})
		}
		for i := range f.Fields {
			walk(&f.Fields[i], below, len(f.Conditions)+i == items-1, depth+1)
		}
	}
	walk(&record, "", true, 0)
	return rows
}

func branch(last, redefines bool) string {
	switch {
	case last && redefines:
		return branchRedefLast
	case redefines:
		return branchRedefines
	case last:
		return branchLast
	}
	return branchMid
}

func continuation(last bool) string {
	if last {
		return branchNone
	}
	return branchDown
}
