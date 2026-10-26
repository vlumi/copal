package blueprint

import (
	"regexp"
	"strings"

	"github.com/vlumi/cobolith/cobol"
	"github.com/vlumi/copal/survey/internal/rules"
)

// sentences writes a decision table's rules as sentences, from the
// conditions each tests and the statements it runs, with no more in them
// than the table has: what is tested and found, in words, a condition name
// with the field it is declared on and its values; then what runs, in
// order. A statement with no words of its own is quoted.
type sentences struct {
	names rules.ConditionNames
}

// rule writes one rule of a table as a sentence.
func (w sentences) rule(t rules.Table, r rules.Rule) string {
	var found []string
	for i, entry := range r.Entries {
		if entry != "" {
			found = append(found, w.test(t.Tests[i], entry == "Y"))
		}
	}
	var done []string
	for _, a := range r.Actions {
		if words := w.action(t.Statements[a]); words != "" {
			done = append(done, words)
		}
	}
	does := "it does nothing more"
	if len(done) > 0 {
		does = "it " + list(done)
	}
	if len(found) == 0 {
		return "Always, " + does + "."
	}
	return "When " + list(found) + ", " + does + "."
}

// test writes a condition, or a conditional phrase, as found true or
// false, and that it is tested again if it is.
func (w sentences) test(t rules.Test, outcome bool) string {
	var words string
	if t.Condition != nil {
		words = w.condition(t.Condition, outcome)
	} else {
		words = phrase(*t.Statement, t.Phrase, outcome)
	}
	if t.Time > 1 {
		words += ", tested again"
	}
	return words
}

// condition writes a condition as found true or false, turning NOT around
// and a compound found false into its parts found false.
func (w sentences) condition(c *cobol.Condition, outcome bool) string {
	switch c.Op {
	case "NOT":
		return w.condition(c.Parts[0], !outcome)
	case "AND", "OR":
		joined := joining(c.Op, outcome)
		var parts []string
		for _, part := range c.Parts {
			words := w.condition(part, outcome)
			if (part.Op == "AND" || part.Op == "OR") && joining(part.Op, outcome) != joined {
				words = map[string]string{"and": "both ", "or": "either "}[joining(part.Op, outcome)] + words
			}
			parts = append(parts, words)
		}
		return joinWith(parts, joined)
	case "IS":
		class := strings.ToLower(strings.ReplaceAll(c.Class, "-", " "))
		if outcome {
			return operand(*c.Left) + " is " + class
		}
		return operand(*c.Left) + " is not " + class
	case "NAME":
		return w.conditionName(*c.Left, outcome)
	case "TRUE":
		if outcome {
			return "always"
		}
		return "never"
	}
	return operand(*c.Left) + " " + relations[c.Op][boolIndex(outcome)] + " " + operand(*c.Right)
}

// relations are the words of a comparison, found true and found false.
var relations = map[string][2]string{
	"=":  {"is", "is not"},
	">":  {"is more than", "is at most"},
	"<":  {"is less than", "is at least"},
	">=": {"is at least", "is less than"},
	"<=": {"is at most", "is more than"},
}

func boolIndex(b bool) int {
	if b {
		return 0
	}
	return 1
}

// joining is the word a compound's parts join with as found: AND found
// false is one of its parts found false, OR found false all of them.
func joining(op string, outcome bool) string {
	if (op == "AND") == outcome {
		return "and"
	}
	return "or"
}

// conditionName writes a condition name as the field it is declared on
// being, or not being, the name's values.
func (w sentences) conditionName(name cobol.Operand, outcome bool) string {
	field, values, ok := w.names(name.Name)
	if !ok {
		if outcome {
			return operand(name) + " holds"
		}
		return operand(name) + " does not hold"
	}
	var shown []string
	for _, v := range values {
		shown = append(shown, value(v))
	}
	if outcome {
		return "`" + field + "` is " + joinWith(shown, "or") + " (" + operand(name) + ")"
	}
	var is string
	switch len(shown) {
	case 1:
		is = "is not " + shown[0]
	case 2:
		is = "is neither " + shown[0] + " nor " + shown[1]
	default:
		is = "is none of " + joinWith(shown, "and")
	}
	return "`" + field + "` " + is + " (not " + operand(name) + ")"
}

// value writes a value of a condition name, a range of them as one.
func value(v string) string {
	if low, high, ok := strings.Cut(v, " THRU "); ok {
		return "from " + literal(low) + " to " + literal(high)
	}
	return literal(v)
}

// phrase writes a statement's conditional phrase as raised or not: the end
// of a file, a key the file refuses, an exception.
func phrase(s cobol.Statement, kind string, raised bool) string {
	object := ""
	switch {
	case len(s.Files) > 0:
		object = "`" + s.Files[0].Name + "`"
	case len(s.Writes) > 0:
		object = operand(s.Writes[0])
	}
	doing := map[string]string{"READ": "reading", "RETURN": "taking the next record from", "WRITE": "writing",
		"REWRITE": "rewriting", "DELETE": "deleting", "START": "positioning"}[s.Verb]
	switch {
	case kind == "AT END" && doing != "" && raised:
		return doing + " " + object + " finds no more records"
	case kind == "AT END" && doing != "":
		return doing + " " + object + " finds a record"
	case kind == "INVALID KEY" && doing != "" && raised:
		return doing + " " + object + " fails on its key"
	case kind == "INVALID KEY" && doing != "":
		return doing + " " + object + " succeeds by its key"
	case raised:
		return "`" + s.Text + "` raises " + strings.ToLower(kind)
	}
	return "`" + s.Text + "` does not raise " + strings.ToLower(kind)
}

// action writes what a statement does, or nothing for one that does
// nothing, CONTINUE and a plain EXIT.
func (w sentences) action(s cobol.Statement) string {
	switch s.Verb {
	case "CONTINUE", "EXIT":
		return ""
	case "STOP":
		return "stops the run"
	case "GOBACK", "EXIT PROGRAM":
		return "returns to its caller"
	case "EXIT PARAGRAPH":
		return "leaves the paragraph"
	case "EXIT SECTION":
		return "leaves the section"
	case "MOVE":
		if strings.Contains(s.Text, " CORR") || len(s.Reads) != 1 {
			break
		}
		return "sets " + operands(s.Writes) + " to " + operand(s.Reads[0])
	case "SET":
		if words := w.set(s); words != "" {
			return words
		}
	case "ADD", "SUBTRACT":
		if words := arithmetic(s); words != "" {
			return words
		}
	case "COMPUTE":
		if len(s.Reads) == 1 {
			return "computes " + operands(s.Writes) + " as " + quoted(s.Reads[0].Text)
		}
	case "INITIALIZE":
		return "initializes " + operands(s.Writes)
	case "DISPLAY":
		return "displays " + spaced(s.Reads)
	case "ACCEPT":
		return accept(s)
	case "STRING":
		return "builds " + operands(s.Writes) + " from " + spaced(s.Reads)
	case "OPEN":
		return open(s)
	case "CLOSE":
		return "closes " + files(s.Files)
	case "READ":
		if len(s.Files) == 1 {
			return "reads " + files(s.Files) + into(s)
		}
	case "RETURN":
		if len(s.Files) == 1 {
			return "takes the next record from " + files(s.Files) + into(s)
		}
	case "WRITE", "REWRITE":
		if len(s.Writes) == 1 {
			return strings.ToLower(s.Verb) + "s " + operand(s.Writes[0]) + from(s)
		}
	case "PERFORM":
		return w.perform(s)
	case "GO":
		if len(s.Procedures) == 1 {
			return "goes to `" + s.Procedures[0].Name + "`"
		}
	case "CALL":
		switch {
		case s.Calls != "":
			return "calls `" + s.Calls + "`"
		case s.CallsThrough != "":
			return "calls the program `" + s.CallsThrough + "` names"
		}
	case "SEARCH":
		table := strings.TrimPrefix(strings.TrimPrefix(s.Text, "SEARCH "), "ALL ")
		if name, _, _ := strings.Cut(table, " "); name != "" {
			return "searches `" + name + "`"
		}
	}
	return "runs " + quoted(s.Text)
}

// set writes a SET of condition names to true as the fields they are
// declared on given their values, and a SET of an index to a value, or up
// or down by one.
func (w sentences) set(s cobol.Statement) string {
	if m := setIndex.FindStringSubmatch(s.Text); m != nil && len(s.Writes) == 1 && m[2] != "TRUE" && m[2] != "FALSE" {
		switch m[1] {
		case "TO":
			return "sets " + operand(s.Writes[0]) + " to " + literal(m[2])
		case "UP BY":
			return "steps " + operand(s.Writes[0]) + " up by " + literal(m[2])
		}
		return "steps " + operand(s.Writes[0]) + " down by " + literal(m[2])
	}
	if !strings.HasSuffix(s.Text, " TO TRUE") {
		return ""
	}
	var parts []string
	for _, name := range s.Writes {
		if field, _, ok := w.names(name.Name); ok {
			parts = append(parts, "sets `"+field+"` to "+operand(name))
		} else {
			parts = append(parts, "sets "+operand(name)+" true")
		}
	}
	return joinWith(parts, "and")
}

var setIndex = regexp.MustCompile(`^SET \S+ (TO|UP BY|DOWN BY) ([0-9]+|[A-Z0-9-]+)$`)

// arithmetic writes an ADD or SUBTRACT with no GIVING: what is added to,
// or taken from, the fields it changes.
func arithmetic(s cobol.Statement) string {
	if strings.Contains(s.Text, " GIVING ") {
		return ""
	}
	changed := map[string]bool{}
	for _, o := range s.Writes {
		changed[o.Text] = true
	}
	var amounts []cobol.Operand
	for _, o := range s.Reads {
		if !changed[o.Text] {
			amounts = append(amounts, o)
		}
	}
	if len(amounts) == 0 {
		return ""
	}
	if s.Verb == "ADD" {
		return "adds " + operands(amounts) + " to " + operands(s.Writes)
	}
	return "subtracts " + operands(amounts) + " from " + operands(s.Writes)
}

func accept(s cobol.Statement) string {
	if len(s.Writes) != 1 {
		return "runs " + quoted(s.Text)
	}
	switch s.Device {
	case "":
		return "accepts " + operand(s.Writes[0]) + " from its input"
	case "COMMAND-LINE":
		return "accepts " + operand(s.Writes[0]) + " from its command line"
	}
	return "accepts " + operand(s.Writes[0]) + " from " + strings.ToLower(strings.ReplaceAll(s.Device, "-", " "))
}

// open writes the files an OPEN opens by the mode it opens them in.
func open(s cobol.Statement) string {
	var parts []string
	var group []cobol.FileRef
	flush := func() {
		if len(group) > 0 {
			parts = append(parts, files(group)+" "+modes[group[0].Mode])
			group = nil
		}
	}
	for _, f := range s.Files {
		if len(group) > 0 && group[0].Mode != f.Mode {
			flush()
		}
		group = append(group, f)
	}
	flush()
	return "opens " + joinWith(parts, "and")
}

var modes = map[string]string{"INPUT": "for input", "OUTPUT": "for output", "EXTEND": "to extend", "I-O": "to update"}

var varying = regexp.MustCompile(`^VARYING (\S+) FROM (\S+) BY (\S+) UNTIL `)

// perform writes a PERFORM: what it performs, or that it repeats its own
// loop, and until when, or the rest of the statement quoted.
func (w sentences) perform(s cobol.Statement) string {
	words := "repeats its loop"
	rest := strings.TrimPrefix(s.Text, "PERFORM")
	if len(s.Procedures) == 1 && len(s.Branches) == 0 {
		p := s.Procedures[0]
		words = "performs `" + p.Name + "`"
		rest = strings.TrimPrefix(rest, " "+p.Name)
		if p.Thru != "" {
			words += " through `" + p.Thru + "`"
			_, rest, _ = strings.Cut(rest, " "+p.Thru)
		}
	}
	rest = strings.TrimSpace(rest)
	switch {
	case rest == "":
		return words
	case s.Condition != nil && strings.HasPrefix(rest, "UNTIL "):
		return words + " until " + w.condition(s.Condition, true)
	case s.Condition != nil && strings.HasPrefix(rest, "WITH TEST AFTER UNTIL "):
		return words + " until, after each pass, " + w.condition(s.Condition, true)
	}
	if m := varying.FindStringSubmatch(rest); m != nil && s.Condition != nil {
		return words + " varying `" + m[1] + "` from " + literal(m[2]) + " by " + literal(m[3]) + " until " + w.condition(s.Condition, true)
	}
	return words + ", " + quoted(rest)
}

func into(s cobol.Statement) string {
	if len(s.Writes) == 1 {
		return " into " + operand(s.Writes[0])
	}
	return ""
}

func from(s cobol.Statement) string {
	if len(s.Reads) == 1 {
		return " from " + operand(s.Reads[0])
	}
	return ""
}

func files(refs []cobol.FileRef) string {
	var names []string
	for _, f := range refs {
		names = append(names, "`"+f.Name+"`")
	}
	return joinWith(names, "and")
}

func operands(ops []cobol.Operand) string {
	var shown []string
	for _, o := range ops {
		shown = append(shown, operand(o))
	}
	return joinWith(shown, "and")
}

// spaced writes what a DISPLAY or STRING puts together, one after another.
func spaced(ops []cobol.Operand) string {
	var shown []string
	for _, o := range ops {
		shown = append(shown, operand(o))
	}
	return strings.Join(shown, " ")
}

// operand writes a data item or an expression quoted, a number as it is,
// and a figurative constant in words.
func operand(o cobol.Operand) string {
	if o.Name != "" {
		return quoted(o.Text)
	}
	return literal(o.Text)
}

var number = regexp.MustCompile(`^[+-]?[0-9]+(\.[0-9]+)?$`)

var figurative = map[string]string{
	"SPACE": "spaces", "SPACES": "spaces", "ZERO": "zero", "ZEROS": "zero", "ZEROES": "zero",
	"HIGH-VALUE": "high values", "HIGH-VALUES": "high values", "LOW-VALUE": "low values", "LOW-VALUES": "low values",
	"QUOTE": "quotes", "QUOTES": "quotes",
}

func literal(text string) string {
	if words, ok := figurative[text]; ok {
		return words
	}
	if number.MatchString(text) {
		return text
	}
	return quoted(text)
}

func quoted(text string) string {
	return "`" + text + "`"
}

// list joins the parts of a sentence, with a comma before the last when
// there are more than two; with semicolons when a part has a comma or a
// conjunction of its own, so that the parts stay apart.
func list(parts []string) string {
	for _, part := range parts {
		if len(parts) > 1 && (strings.Contains(part, ",") || strings.Contains(part, " and ") || strings.Contains(part, " or ")) {
			return strings.Join(parts[:len(parts)-1], "; ") + "; and " + parts[len(parts)-1]
		}
	}
	return joinWith(parts, "and")
}

func joinWith(parts []string, word string) string {
	switch len(parts) {
	case 0:
		return ""
	case 1:
		return parts[0]
	case 2:
		return parts[0] + " " + word + " " + parts[1]
	}
	return strings.Join(parts[:len(parts)-1], ", ") + ", " + word + " " + parts[len(parts)-1]
}
