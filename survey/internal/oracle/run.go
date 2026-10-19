package oracle

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/vlumi/copal/bedrock/cobol"
	"github.com/vlumi/copal/bedrock/layout"
	"github.com/vlumi/copal/bedrock/run"
	"github.com/vlumi/copal/survey/blueprint"
	"github.com/vlumi/copal/survey/internal/coverage"
)

// Case is the input of one run of a program: the records of each input,
// by the name of its file, or SYSIN for the lines the program accepts.
type Case map[string][][]byte

// rule is a rule of one of the decision tables of a unit, by their
// indexes.
type rule struct{ unit, table, rule int }

// timeout is how long a run may take before the search gives up on it.
const timeout = 10 * time.Second

// runner runs a program on cases: the system's programs compiled with
// trace code, and the loaders of the indexed files it reads.
type runner struct {
	program  *program
	programs string
	loaders  string
}

// result is what a run reached, and how it ended; and, for a run that
// records its outputs, what it wrote to each file it writes or updates,
// and to its standard output.
type result struct {
	reached map[rule]bool
	outcome run.Outcome
	outputs map[string][]byte
	stdout  []byte
}

// run runs the program on a case, in a directory of its own, and replays
// its trace against the decision tables, recording its outputs if asked.
func (r *runner) run(c Case, dir string, record bool) (result, error) {
	dds := map[string]string{}
	for _, f := range r.program.blueprint.Files {
		path := filepath.Join(dir, f.Name)
		dds[f.Assign] = path
		var err error
		switch {
		case reads(f) && f.Organization == "INDEXED":
			err = r.load(f, c[f.Name], path, dir)
		case reads(f):
			err = os.WriteFile(path, sequential(f, c[f.Name]), 0o644)
		case writes(f) && f.Organization != "INDEXED":
			err = os.WriteFile(path, nil, 0o644)
		}
		if err != nil {
			return result{}, err
		}
	}
	trace := filepath.Join(dir, "trace")
	var stdout bytes.Buffer
	var args []string
	if line := c[commandLine]; len(line) > 0 {
		args = []string{strings.TrimRight(string(line[0]), " ")}
	}
	outcome, err := run.Step{Programs: r.programs, Program: r.program.source.ID, Args: args, DDs: dds, Stdout: &stdout,
		Stdin: bytes.NewReader(lines(c[sysin])), TraceFile: trace, Timeout: timeout}.Run()
	if err != nil {
		return result{}, err
	}
	res := result{outcome: outcome, stdout: stdout.Bytes()}
	if record {
		if res.outputs, err = r.outputs(dir); err != nil {
			return result{}, err
		}
	}
	data, err := os.ReadFile(trace)
	if err != nil && !os.IsNotExist(err) {
		return result{}, err
	}
	events, err := coverage.Read(bytes.NewReader(data))
	if err != nil {
		return result{}, err
	}
	res.reached = map[rule]bool{}
	for u, unit := range r.program.units {
		for i, t := range coverage.Replay(unit.source, unit.Tables, events) {
			for j, n := range t.Reached {
				if n > 0 {
					res.reached[rule{u, i, j}] = true
				}
			}
		}
	}
	return res, nil
}

// outputs reads what a run wrote to the files the program writes or
// updates, a record a line: an indexed file through its unloader, in key
// order.
func (r *runner) outputs(dir string) (map[string][]byte, error) {
	outputs := map[string][]byte{}
	for _, f := range r.program.blueprint.Files {
		if !writes(f) && !slices.Contains(f.Opened, "I-O") {
			continue
		}
		path := filepath.Join(dir, f.Name)
		if f.Organization == "INDEXED" {
			unloaded := path + ".unload"
			outcome, err := run.Step{Programs: r.loaders, Program: unloaderName(f), DDs: map[string]string{"UNLOADIN": path, "UNLOADOUT": unloaded},
				Stdout: io.Discard, Timeout: timeout}.Run()
			if err == nil && (outcome.RC != 0 || outcome.Abended) {
				err = fmt.Errorf("unloading %s: %+v", f.Name, outcome)
			}
			if err != nil {
				return nil, err
			}
			path = unloaded
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		outputs[f.Name] = asLines(data, f.Organization == "LINE SEQUENTIAL", r.recordLength(f))
	}
	return outputs, nil
}

func (r *runner) recordLength(f blueprint.File) int {
	if root := fileRecord(r.program.records, f.Name); root != nil {
		return root.Size()
	}
	return 0
}

// asLines writes a file's records a line each: as they are, with trailing
// spaces dropped, or quoted, with escapes, for a record that holds bytes
// that are not printable.
func asLines(data []byte, lineSequential bool, length int) []byte {
	var records [][]byte
	switch {
	case lineSequential:
		records = bytes.Split(bytes.TrimSuffix(data, []byte("\n")), []byte("\n"))
		if len(data) == 0 {
			records = nil
		}
	case length > 0:
		for len(data) >= length {
			records, data = append(records, data[:length]), data[length:]
		}
	}
	var b bytes.Buffer
	for _, record := range records {
		line := string(bytes.TrimRight(record, " "))
		if !printable(line) {
			line = strconv.Quote(line)
		}
		b.WriteString(line + "\n")
	}
	return b.Bytes()
}

func printable(s string) bool {
	for _, r := range s {
		if r < ' ' || r > '~' {
			return false
		}
	}
	return !strings.HasPrefix(s, `"`)
}

// sequential writes records as a sequential file holds them: back to back,
// or a line each for a line sequential file.
func sequential(f blueprint.File, records [][]byte) []byte {
	if f.Organization == "LINE SEQUENTIAL" {
		return lines(records)
	}
	return bytes.Join(records, nil)
}

func lines(records [][]byte) []byte {
	var b bytes.Buffer
	for _, r := range records {
		b.Write(bytes.TrimRight(r, " "))
		b.WriteByte('\n')
	}
	return b.Bytes()
}

// load writes an indexed file through its loader, which reads the records
// from a sequential file.
func (r *runner) load(f blueprint.File, records [][]byte, path, dir string) error {
	in := path + ".load"
	if err := os.WriteFile(in, bytes.Join(records, nil), 0o644); err != nil {
		return err
	}
	outcome, err := run.Step{Programs: r.loaders, Program: loaderName(f), DDs: map[string]string{"LOADIN": in, "LOADOUT": path},
		Stdout: io.Discard, Timeout: timeout}.Run()
	if err == nil && (outcome.RC != 0 || outcome.Abended) {
		err = fmt.Errorf("loading %s: %+v", f.Name, outcome)
	}
	return err
}

func loaderName(f blueprint.File) string {
	return "LOAD" + strings.ReplaceAll(f.Name, "-", "")
}

func unloaderName(f blueprint.File) string {
	return "UNLD" + strings.ReplaceAll(f.Name, "-", "")
}

// unloaderSource is a program that writes the records of an indexed file,
// in key order, to a sequential one.
func unloaderSource(f blueprint.File, length, keyStart, keyLength int) string {
	return fmt.Sprintf(`       IDENTIFICATION DIVISION.
       PROGRAM-ID. %s.
       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT UNLOAD-IN ASSIGN TO UNLOADIN
               ORGANIZATION IS INDEXED
               ACCESS MODE IS SEQUENTIAL
               RECORD KEY IS UNLOAD-KEY.
           SELECT UNLOAD-OUT ASSIGN TO UNLOADOUT
               ORGANIZATION IS SEQUENTIAL.
       DATA DIVISION.
       FILE SECTION.
       FD  UNLOAD-IN.
       01  UNLOAD-RECORD.
%s
       FD  UNLOAD-OUT.
       01  UNLOAD-OUT-RECORD           PIC X(%d).
       WORKING-STORAGE SECTION.
       01  UNLOAD-EOF                  PIC X          VALUE 'N'.
       PROCEDURE DIVISION.
       UNLOAD-MAIN.
           OPEN INPUT UNLOAD-IN OUTPUT UNLOAD-OUT
           PERFORM UNTIL UNLOAD-EOF = 'Y'
               READ UNLOAD-IN NEXT RECORD
                   AT END
                       MOVE 'Y' TO UNLOAD-EOF
                   NOT AT END
                       WRITE UNLOAD-OUT-RECORD FROM UNLOAD-RECORD
               END-READ
           END-PERFORM
           CLOSE UNLOAD-IN UNLOAD-OUT
           STOP RUN.
`, unloaderName(f), keyFields("UNLOAD-KEY", length, keyStart, keyLength), length)
}

// keyFields lays out a record of the given length as its key, at the
// given bytes, and the fillers around it.
func keyFields(key string, length, keyStart, keyLength int) string {
	var fields []string
	if keyStart > 1 {
		fields = append(fields, fmt.Sprintf("           05  FILLER                  PIC X(%d).", keyStart-1))
	}
	fields = append(fields, fmt.Sprintf("           05  %-23s PIC X(%d).", key, keyLength))
	if rest := length - keyStart - keyLength + 1; rest > 0 {
		fields = append(fields, fmt.Sprintf("           05  FILLER                  PIC X(%d).", rest))
	}
	return strings.Join(fields, "\n")
}

// loaderSource is a program that writes the records of a sequential file
// into an indexed one, by the key at the given bytes.
func loaderSource(f blueprint.File, length, keyStart, keyLength int) string {
	return fmt.Sprintf(`       IDENTIFICATION DIVISION.
       PROGRAM-ID. %s.
       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT LOAD-IN ASSIGN TO LOADIN
               ORGANIZATION IS SEQUENTIAL.
           SELECT LOAD-OUT ASSIGN TO LOADOUT
               ORGANIZATION IS INDEXED
               ACCESS MODE IS RANDOM
               RECORD KEY IS LOAD-KEY.
       DATA DIVISION.
       FILE SECTION.
       FD  LOAD-IN.
       01  LOAD-IN-RECORD              PIC X(%d).
       FD  LOAD-OUT.
       01  LOAD-RECORD.
%s
       WORKING-STORAGE SECTION.
       01  LOAD-EOF                    PIC X          VALUE 'N'.
       PROCEDURE DIVISION.
       LOAD-MAIN.
           OPEN INPUT LOAD-IN OUTPUT LOAD-OUT
           PERFORM UNTIL LOAD-EOF = 'Y'
               READ LOAD-IN
                   AT END
                       MOVE 'Y' TO LOAD-EOF
                   NOT AT END
                       WRITE LOAD-RECORD FROM LOAD-IN-RECORD
                           INVALID KEY
                               CONTINUE
                       END-WRITE
               END-READ
           END-PERFORM
           CLOSE LOAD-IN LOAD-OUT
           STOP RUN.
`, loaderName(f), length, keyFields("LOAD-KEY", length, keyStart, keyLength))
}

// newRunner prepares the runs of a program in a scratch directory: the
// system's programs compiled there with trace code, once for every
// program the scratch directory serves, and for each indexed file the
// program uses, a loader and an unloader.
func newRunner(p *program, sources, copybookDir, scratch string) (*runner, error) {
	dir := filepath.Join(scratch, p.blueprint.Name)
	r := &runner{program: p, programs: filepath.Join(scratch, "programs"), loaders: filepath.Join(dir, "loaders")}
	if _, err := os.Stat(r.programs); os.IsNotExist(err) {
		if _, err := (run.Build{Dialect: p.dialect, Sources: sources, Copybooks: copybookDir, Out: r.programs, Trace: true}).Programs(); err != nil {
			return nil, err
		}
	}
	loaderSources := filepath.Join(dir, "loader-sources")
	if err := os.MkdirAll(loaderSources, 0o755); err != nil {
		return nil, err
	}
	for _, f := range p.blueprint.Files {
		if f.Organization != "INDEXED" {
			continue
		}
		root := fileRecord(p.records, f.Name)
		var key *layout.Field
		if root != nil {
			key = find(root, f.RecordKey)
		}
		if key == nil {
			return nil, fmt.Errorf("%s: no record key %s", f.Name, f.RecordKey)
		}
		length, keyStart := root.Size(), key.Start-root.Start+1
		for name, source := range map[string]string{
			loaderName(f):   loaderSource(f, length, keyStart, key.Size()),
			unloaderName(f): unloaderSource(f, length, keyStart, key.Size()),
		} {
			if err := os.WriteFile(filepath.Join(loaderSources, name+".cbl"), []byte(source), 0o644); err != nil {
				return nil, err
			}
		}
	}
	// The loaders are the oracle's own, in IBM's fixed format.
	if _, err := (run.Build{Dialect: cobol.IBM, Sources: loaderSources, Copybooks: loaderSources, Out: r.loaders}).Programs(); err != nil {
		return nil, err
	}
	return r, nil
}
