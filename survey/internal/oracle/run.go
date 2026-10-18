package oracle

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

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

// result is what a run reached, and how it ended.
type result struct {
	reached map[rule]bool
	outcome run.Outcome
}

// run runs the program on a case, in a directory of its own, and replays
// its trace against the decision tables.
func (r *runner) run(c Case, dir string) (result, error) {
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
		case writes(f):
			err = os.WriteFile(path, nil, 0o644)
		}
		if err != nil {
			return result{}, err
		}
	}
	trace := filepath.Join(dir, "trace")
	outcome, err := run.Step{Programs: r.programs, Program: r.program.source.ID, DDs: dds, Stdout: io.Discard,
		Stdin: bytes.NewReader(lines(c[sysin])), TraceFile: trace, Timeout: timeout}.Run()
	if err != nil {
		return result{}, err
	}
	data, err := os.ReadFile(trace)
	if err != nil && !os.IsNotExist(err) {
		return result{}, err
	}
	events, err := coverage.Read(bytes.NewReader(data))
	if err != nil {
		return result{}, err
	}
	reached := map[rule]bool{}
	for u, unit := range r.program.units {
		for i, t := range coverage.Replay(unit.source, unit.Tables, events) {
			for j, n := range t.Reached {
				if n > 0 {
					reached[rule{u, i, j}] = true
				}
			}
		}
	}
	return result{reached: reached, outcome: outcome}, nil
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

// loaderSource is a program that writes the records of a sequential file
// into an indexed one, by the key at the given bytes.
func loaderSource(f blueprint.File, length, keyStart, keyLength int) string {
	fields := []string{}
	if keyStart > 1 {
		fields = append(fields, fmt.Sprintf("           05  FILLER                  PIC X(%d).", keyStart-1))
	}
	fields = append(fields, fmt.Sprintf("           05  LOAD-KEY                PIC X(%d).", keyLength))
	if rest := length - keyStart - keyLength + 1; rest > 0 {
		fields = append(fields, fmt.Sprintf("           05  FILLER                  PIC X(%d).", rest))
	}
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
`, loaderName(f), length, strings.Join(fields, "\n"))
}

// newRunner compiles the system's programs with trace code, and a loader
// for each indexed file the program reads.
func newRunner(p *program, sources, copybookDir, dir string) (*runner, error) {
	r := &runner{program: p, programs: filepath.Join(dir, "programs"), loaders: filepath.Join(dir, "loaders")}
	if _, err := (run.Build{Sources: sources, Copybooks: copybookDir, Out: r.programs, Trace: true}).Programs(); err != nil {
		return nil, err
	}
	loaderSources := filepath.Join(dir, "loader-sources")
	if err := os.MkdirAll(loaderSources, 0o755); err != nil {
		return nil, err
	}
	for _, in := range p.inputs {
		if in.file.Organization != "INDEXED" {
			continue
		}
		key := find(in.root, in.file.RecordKey)
		if key == nil {
			return nil, fmt.Errorf("%s: no record key %s", in.file.Name, in.file.RecordKey)
		}
		source := loaderSource(in.file, in.root.Size(), key.Start-in.root.Start+1, key.Size())
		if err := os.WriteFile(filepath.Join(loaderSources, loaderName(in.file)+".cbl"), []byte(source), 0o644); err != nil {
			return nil, err
		}
	}
	if _, err := (run.Build{Sources: loaderSources, Copybooks: loaderSources, Out: r.loaders}).Programs(); err != nil {
		return nil, err
	}
	return r, nil
}
