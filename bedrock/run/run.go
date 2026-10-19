package run

import (
	"context"
	"errors"
	"io"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"time"
)

// Step is one run of a compiled program: the directory of the compiled
// programs, where the runtime also finds the modules the program calls,
// the arguments on its command line, its DD names bound to files, what it
// reads as SYSIN on standard input,
// where its standard output and error go, for a program compiled with
// trace code, the file the trace is written to, and, if set, how long it
// may run before it is killed.
type Step struct {
	Programs  string
	Program   string
	Args      []string
	DDs       map[string]string
	Stdin     io.Reader
	Stdout    io.Writer
	TraceFile string
	Timeout   time.Duration
}

// Outcome is how a run ended: its return code, or an abend, when the
// program was killed or crashed rather than ending, and whether it was
// killed for running out of time.
type Outcome struct {
	RC       int
	Abended  bool
	TimedOut bool
}

// Run runs the program and waits for it to end. The error is for a
// program that could not be run at all, not for its return code.
func (s Step) Run() (Outcome, error) {
	ctx := context.Background()
	if s.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.Timeout)
		defer cancel()
	}
	cmd := exec.CommandContext(ctx, filepath.Join(s.Programs, s.Program), s.Args...)
	cmd.Env = append(os.Environ(), "COB_LIBRARY_PATH="+s.Programs)
	for _, dd := range slices.Sorted(maps.Keys(s.DDs)) {
		cmd.Env = append(cmd.Env, "DD_"+dd+"="+s.DDs[dd])
	}
	if s.TraceFile != "" {
		cmd.Env = append(cmd.Env, "COB_SET_TRACE=Y", "COB_TRACE_FILE="+s.TraceFile)
	}
	cmd.Stdin, cmd.Stdout, cmd.Stderr = s.Stdin, s.Stdout, s.Stdout
	err := cmd.Run()
	var exit *exec.ExitError
	switch {
	case errors.As(err, &exit) && exit.Exited():
		return Outcome{RC: exit.ExitCode()}, nil
	case errors.As(err, &exit):
		return Outcome{Abended: true, TimedOut: ctx.Err() != nil}, nil
	case err != nil:
		return Outcome{}, err
	}
	return Outcome{}, nil
}
