package console

import (
	"slices"
	"testing"

	"github.com/vlumi/copal/bedrock/runner"
)

var accounts = []Account{
	{No: 1000000026, Customer: 1, Name: "YUNNANOZOON FELIX", Status: "O", Balance: 100},
	{No: 1000000067, Customer: 2, Name: "KOMODO ADA", Status: "O", Balance: -50},
	{No: 1000002600, Customer: 26, Name: "FERN FELIX", Status: "F", Balance: 0},
}

func numbers(found []Account) []int64 {
	var nos []int64
	for _, a := range found {
		nos = append(nos, a.No)
	}
	return nos
}

func TestSearch(t *testing.T) {
	for query, want := range map[string][]int64{
		"":           {1000000026, 1000000067, 1000002600},
		"felix":      {1000000026, 1000002600},
		"100000006":  {1000000067},
		" 26 ":       {1000002600},
		"no such":    nil,
		"1000000026": {1000000026},
	} {
		if got := numbers(Search(accounts, query)); !slices.Equal(got, want) {
			t.Errorf("search %q: %v, want %v", query, got, want)
		}
	}
}

func TestWorst(t *testing.T) {
	for _, c := range []struct {
		steps []runner.Result
		want  string
	}{
		{nil, ""},
		{[]runner.Result{{RC: 0}, {RC: 4}, {RC: 0}}, "RC=0004"},
		{[]runner.Result{{RC: 12}, {Bypassed: true}}, "BYPASSED"},
		{[]runner.Result{{Abended: true}, {RC: 16}}, "ABENDED"},
	} {
		if got := (DaySummary{Steps: c.steps}).Worst(); got != c.want {
			t.Errorf("%+v: worst %q, want %q", c.steps, got, c.want)
		}
	}
}

func TestComparison(t *testing.T) {
	after := []Account{
		{No: 1000000026, Status: "O", Balance: 70},
		{No: 1000000067, Status: "O", Balance: -50},
		{No: 1000002600, Status: "X", Balance: 0},
		{No: 1000002700, Status: "O", Balance: 0},
	}
	c := Comparison{Before: balances(accounts), After: balances(after), Posted: -30}
	c.changes(accounts, after)
	if !c.Reconciled() {
		t.Errorf("balances %d to %d, posted %d: not reconciled", c.Before.Total, c.After.Total, c.Posted)
	}
	if len(c.Changes) != 1 || c.Changes[0].Delta() != -30 {
		t.Errorf("changes %+v, want 1000000026 down 30", c.Changes)
	}
	if len(c.StatusChanges) != 1 || c.StatusChanges[0].After.Status != "X" {
		t.Errorf("status changes %+v, want 1000002600 closed", c.StatusChanges)
	}
	if got := numbers(c.Opened); !slices.Equal(got, []int64{1000002700}) {
		t.Errorf("opened %v, want 1000002700", got)
	}
	c.Posted = 0
	if c.Reconciled() {
		t.Error("reconciled with nothing posted")
	}
}
