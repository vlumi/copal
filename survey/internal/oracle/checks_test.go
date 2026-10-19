package oracle

import (
	"testing"

	"github.com/vlumi/copal/survey/internal/rules"
)

func TestFailsWrite(t *testing.T) {
	table := rules.Table{Conditions: []string{
		"READ MASTER INVALID KEY",
		"REWRITE MASTER-RECORD INVALID KEY",
		"WRITE LOG-RECORD INVALID KEY, again",
	}}
	for _, c := range []struct {
		entries []string
		want    bool
	}{
		{[]string{"Y", "", ""}, false},
		{[]string{"N", "N", ""}, false},
		{[]string{"N", "Y", ""}, true},
		{[]string{"N", "N", "Y"}, true},
	} {
		if got := failsWrite(table, rules.Rule{Entries: c.entries}); got != c.want {
			t.Errorf("rule %v fails a write: %v, want %v", c.entries, got, c.want)
		}
	}
}
