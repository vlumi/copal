package layouts

import "testing"

func TestLoad(t *testing.T) {
	l, err := Load("../copybooks")
	if err != nil {
		t.Fatal(err)
	}
	for name, got := range map[string]int{"master": l.Master.Length, "transaction": l.Tran.Length, "accrual": l.Accrual.Length} {
		if want := map[string]int{"master": 100, "transaction": 80, "accrual": 40}[name]; got != want {
			t.Errorf("%s record %d bytes, want %d", name, got, want)
		}
	}
}
