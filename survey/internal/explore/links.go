package explore

import (
	"net/url"
	"strings"
)

// The explorer's pages, by what they show.

func programLink(name string) string { return "/programs/" + url.PathEscape(name) }

func paragraphLink(program, paragraph string) string {
	return programLink(program) + "/paragraphs/" + url.PathEscape(paragraph)
}

func fieldLink(program, field string) string {
	return programLink(program) + "/fields/" + url.PathEscape(field)
}

func copybookLink(name string) string { return "/copybooks/" + url.PathEscape(name) }

func jobLink(name string) string { return "/jobs/" + url.PathEscape(name) }

func stepLink(s StepRef) string { return jobLink(s.Job) + "/steps/" + url.PathEscape(s.Step) }

func datasetLink(base string) string { return "/datasets/" + url.PathEscape(base) }

// sourceLink links a source reference, path or path:line, to the source
// page, at the line.
func sourceLink(source string) string {
	file, line, ok := strings.Cut(source, ":")
	link := "/source/" + (&url.URL{Path: file}).EscapedPath()
	if ok {
		link += "#L" + line
	}
	return link
}
