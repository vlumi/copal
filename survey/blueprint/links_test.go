package blueprint

import "testing"

func TestIsHighest(t *testing.T) {
	for literal, want := range map[string]bool{
		"HIGH-VALUES": true, "HIGH-VALUE": true, `ALL "9"`: true, "ALL '9'": true, "99999999": true,
		"ZERO": false, "'9A'": false, "ALL '0'": false, "": false, "1": false,
	} {
		if got := isHighest(literal); got != want {
			t.Errorf("isHighest(%q) = %v, want %v", literal, got, want)
		}
	}
}
