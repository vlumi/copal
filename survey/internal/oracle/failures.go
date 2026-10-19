package oracle

import (
	"os"
	"path/filepath"
	"strings"
)

// failure is how a case fails one of the files a program opens, the same
// way every run: a file it reads is missing, so that opening it fails, or
// a directory is where it should be, so that it opens but reading it
// fails; or the directory of a file it writes is missing, so that creating
// it fails.
type failure string

const (
	missing     failure = "missing"
	unreadable  failure = "unreadable"
	uncreatable failure = "uncreatable"
)

// failed is the name a case's failed file goes by among its inputs: one
// line, the file's name and how it fails; failedName is the file the line
// is written to, beside the inputs, in lower case so that it is no input's
// name.
const (
	failed     = "\x00failed"
	failedName = "failed.txt"
)

// failures lists how each of a program's files may fail: any one it reads
// may be missing, and a sequential one unreadable, and any one it only
// writes uncreatable. An indexed or relative file is read through its file
// handler, which may not read a directory the same way everywhere.
func (p *program) failures() []lever {
	var levers []lever
	for _, in := range p.inputs {
		if in.file.Name == "" {
			continue
		}
		levers = append(levers, lever{input: in.name(), change: failFile, failure: missing})
		if !keyed(in.file) {
			levers = append(levers, lever{input: in.name(), change: failFile, failure: unreadable})
		}
	}
	for _, f := range p.blueprint.Files {
		if writes(f) && !reads(f) {
			levers = append(levers, lever{input: f.Name, change: failFile, failure: uncreatable})
		}
	}
	return levers
}

// fail makes a case fail a file, if the case fails none yet: the file has
// no records, and the case records how it fails.
func (c Case) fail(file string, how failure) (Case, bool) {
	if _, _, ok := c.failedFile(); ok {
		return nil, false
	}
	next := c.clone()
	delete(next, file)
	next[failed] = [][]byte{[]byte(file + " " + string(how))}
	return next, true
}

// failedFile is the file a case fails, and how, if it fails one.
func (c Case) failedFile() (string, failure, bool) {
	if len(c[failed]) == 0 {
		return "", "", false
	}
	file, how, _ := strings.Cut(string(c[failed][0]), " ")
	return file, failure(how), true
}

// failsFile reports whether a case fails a file.
func (c Case) failsFile(file string) (failure, bool) {
	name, how, ok := c.failedFile()
	return how, ok && name == file
}

// prepareFailure puts in place what fails a file at a path, and gives the
// path to bind it to: nothing for a missing file, a directory for an
// unreadable one, and for an uncreatable one, a path in a directory that
// is not there.
func prepareFailure(how failure, path string) (string, error) {
	switch how {
	case unreadable:
		return path, os.Mkdir(path, 0o755)
	case uncreatable:
		return filepath.Join(path+".missing", filepath.Base(path)), nil
	}
	return path, nil
}

// describe says how a case fails its file, for its README.
func (how failure) describe(file string) string {
	return map[failure]string{
		missing:     file + " is missing, so opening it fails.",
		unreadable:  file + " is a directory, so it opens but reading it fails.",
		uncreatable: file + "'s directory is missing, so creating it fails.",
	}[how]
}
