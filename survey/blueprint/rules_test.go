package blueprint

import (
	"reflect"
	"testing"

	"github.com/vlumi/copal/survey/internal/rules"
)

// A condition name decides on its field, and a program's decisions on a
// field are the tables with a condition on it, with those conditions.
func TestDecisions(t *testing.T) {
	records := []Record{{Field: Field{Name: "R", Fields: []Field{
		{Name: "STATUS", Conditions: []Condition{{Name: "IS-OPEN"}, {Name: "IS-SHUT"}}},
		{Name: "AMOUNT"},
	}}}}
	p := Program{Name: "P", Records: records, Rules: newRules("P", []rules.Table{
		{Paragraph: "A", Conditions: []string{"IS-OPEN", "AMOUNT > 0", "NOT-OURS"}, Refs: [][]string{{"IS-OPEN"}, {"AMOUNT"}, {"OTHER"}}},
		{Paragraph: "B", Loop: "PERFORM UNTIL IS-SHUT", Conditions: []string{"IS-SHUT OR STATUS = 'Z'"}, Refs: [][]string{{"IS-SHUT", "STATUS"}}},
		{Paragraph: "C", Conditions: []string{"AMOUNT = 0"}, Refs: [][]string{{"AMOUNT"}}},
	}, records)}
	if got, want := p.Rules.Tables[1].ConditionRefs, [][]string{{"STATUS"}}; !reflect.DeepEqual(got, want) {
		t.Errorf("refs %q, want %q", got, want)
	}
	want := []Decision{
		{Paragraph: "A", Source: ":0", Conditions: []string{"IS-OPEN"}},
		{Paragraph: "B", Loop: "PERFORM UNTIL IS-SHUT", Source: ":0", Conditions: []string{"IS-SHUT OR STATUS = 'Z'"}},
	}
	if got := p.Decisions("STATUS"); !reflect.DeepEqual(got, want) {
		t.Errorf("decisions %+v, want %+v", got, want)
	}
}
