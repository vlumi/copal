package batchfile

import "strings"

// expand expands the variables in text, %NAME%, %1 to %9 and %ERRORLEVEL%,
// through value, which gives a variable's value by its name in upper case
// and whether it has one; a variable with none stays as written, as
// cmd.exe leaves it. Backslashes become slashes, since the paths are run
// on a machine that has only those.
func expand(text string, value func(name string) (string, bool)) string {
	var b strings.Builder
	for i := 0; i < len(text); i++ {
		if text[i] == '%' {
			if j := strings.IndexByte(text[i+1:], '%'); j > 0 {
				name := strings.ToUpper(text[i+1 : i+1+j])
				if v, ok := value(name); ok {
					b.WriteString(v)
					i += j + 1
					continue
				}
			}
			if i+1 < len(text) && text[i+1] >= '1' && text[i+1] <= '9' {
				if v, ok := value(text[i+1 : i+2]); ok {
					b.WriteString(v)
					i++
					continue
				}
			}
		}
		if text[i] == '\\' {
			b.WriteByte('/')
			continue
		}
		b.WriteByte(text[i])
	}
	return b.String()
}

// programName is the name of the program a path runs: its file name
// without its directory or .EXE.
func programName(path string) string {
	name := path[strings.LastIndexAny(path, `/\`)+1:]
	if strings.HasSuffix(strings.ToUpper(name), ".EXE") {
		name = name[:len(name)-4]
	}
	return name
}

// batchName is the name of the batch file a path calls: its file name
// without its directory or .BAT.
func batchName(path string) string {
	name := path[strings.LastIndexAny(path, `/\`)+1:]
	if strings.HasSuffix(strings.ToUpper(name), ".BAT") {
		name = name[:len(name)-4]
	}
	return name
}
