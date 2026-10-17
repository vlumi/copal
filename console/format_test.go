package console

import "testing"

func TestMoney(t *testing.T) {
	for cents, want := range map[int64]string{
		0: "0.00", 5: "0.05", -5: "-0.05", 152000: "1,520.00", -3614490600: "-36,144,906.00", 99999: "999.99",
	} {
		if got := money(cents); got != want {
			t.Errorf("money(%d) = %q, want %q", cents, got, want)
		}
	}
}

func TestDecimal(t *testing.T) {
	if got := decimal(20260901, 0); got != "20260901" {
		t.Errorf("no decimals: %q", got)
	}
	if got := decimal(-15000, 2); got != "-150.00" {
		t.Errorf("two decimals: %q", got)
	}
}
