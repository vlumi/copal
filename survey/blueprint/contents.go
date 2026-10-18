package blueprint

import (
	"fmt"
	"strings"
	"unicode"
)

// withContents adds a table of contents to a Markdown document, before its
// first section: a link to each section and subsection, by the anchors
// GitHub gives headings. A document with fewer than three sections and
// subsections gets none.
func withContents(md string) string {
	lines := strings.Split(md, "\n")
	var entries []string
	first := -1
	used := map[string]int{}
	fenced := false
	for i, line := range lines {
		if strings.HasPrefix(line, "```") {
			fenced = !fenced
		}
		level := strings.IndexFunc(line, func(r rune) bool { return r != '#' })
		if fenced || level < 2 || level > 3 || !strings.HasPrefix(line[level:], " ") {
			continue
		}
		if first < 0 {
			first = i
		}
		title := strings.TrimSpace(line[level:])
		anchor := HeadingAnchor(title)
		if n := used[anchor]; n > 0 {
			used[anchor]++
			anchor = fmt.Sprintf("%s-%d", anchor, n)
		} else {
			used[anchor] = 1
		}
		entries = append(entries, fmt.Sprintf("%s- [%s](#%s)", strings.Repeat("  ", level-2), title, anchor))
	}
	if len(entries) < 3 {
		return md
	}
	contents := "## Contents\n\n" + strings.Join(entries, "\n") + "\n"
	return strings.Join(lines[:first], "\n") + "\n" + contents + "\n" + strings.Join(lines[first:], "\n")
}

// HeadingAnchor is the anchor GitHub gives a heading: lower case, spaces as
// hyphens, and punctuation other than hyphens and underscores dropped.
func HeadingAnchor(title string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(title) {
		switch {
		case r == ' ':
			b.WriteRune('-')
		case r == '-' || r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
		}
	}
	return b.String()
}
