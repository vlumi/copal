package blueprint

import "testing"

func TestWithContents(t *testing.T) {
	md := "# TXPOST\n\nIntro.\n\n## Files\n\ntext\n\n```text\n## not a heading\n```\n\n## Records\n\n### WS-TOTALS\n\n### WS-TOTALS\n"
	want := "# TXPOST\n\nIntro.\n\n## Contents\n\n- [Files](#files)\n- [Records](#records)\n  - [WS-TOTALS](#ws-totals)\n" +
		"  - [WS-TOTALS](#ws-totals-1)\n\n## Files\n\ntext\n\n```text\n## not a heading\n```\n\n## Records\n\n### WS-TOTALS\n\n### WS-TOTALS\n"
	if got := withContents(md); got != want {
		t.Errorf("got\n%s\nwant\n%s", got, want)
	}
	short := "# A\n\n## Only\n"
	if got := withContents(short); got != short {
		t.Errorf("a short document got contents:\n%s", got)
	}
}
