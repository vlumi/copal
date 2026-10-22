package batchfile

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

const daily = "@echo off\r\n" +
	"rem DAILY.BAT: one day.\r\n" +
	"rem\r\n" +
	"rem usage: DAILY.BAT DATE\r\n" +
	"\r\n" +
	"setlocal\r\n" +
	"set DATE=%1\r\n" +
	"set DATA=%GLS_DATA%\r\n" +
	"set WORK=%DATA%\\WORK\\%DATE%\r\n" +
	"set LOG=%DATA%\\LOGS\\DAILY.%DATE%.LOG\r\n" +
	"if not exist \"%WORK%\" mkdir \"%WORK%\"\r\n" +
	"\r\n" +
	"rem STEP LOAD: the day's file\r\n" +
	"rem into the work area\r\n" +
	"set DD_IN=%DATA%\\IN\\IN.%DATE%.TXT\r\n" +
	"set \"DD_OUT=%WORK%\\LOADED.TXT\"\r\n" +
	"\"%GLS_BIN%\\LOAD.EXE\" %DATE% > \"%DATA%\\REPORTS\\LOAD.OUT\" 2>&1\r\n" +
	">>\"%LOG%\" echo LOAD LOAD %ERRORLEVEL%\r\n" +
	"if errorlevel 8 exit /b %ERRORLEVEL%\r\n" +
	"\r\n" +
	"rem STEP MORE: the rest of the day\r\n" +
	"call \"%GLS_SCRIPTS%\\MORE.BAT\" %DATE%\r\n" +
	"if errorlevel 8 exit /b 16\r\n" +
	"\r\n" +
	"move /y \"%WORK%\\LOADED.TXT\" \"%DATA%\\MASTER\\KEYED.TXT\"\r\n" +
	"exit /b 0\r\n"

func TestScript(t *testing.T) {
	f, err := Parse("DAILY", strings.NewReader(daily))
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(f.Comment, []string{"DAILY.BAT: one day.", "", "usage: DAILY.BAT DATE"}) {
		t.Errorf("comment %q", f.Comment)
	}
	s := f.Script()
	var got []string
	for _, st := range s.Steps {
		got = append(got, fmt.Sprintf("%s %s %s %v %v %q %q %v %v %d", st.Name, st.Kind, st.Program, st.Args, st.Input, st.Output, st.Stdout, st.StopsIf, st.Comment, len(st.Bindings)))
	}
	want := []string{
		`LOAD program LOAD [$DATE] [] "" "$DATA/REPORTS/LOAD.OUT" &{-ge 8} [the day's file into the work area] 2`,
		`MORE script MORE [$DATE] [] "" "" &{-ge 8} [the rest of the day] 0`,
		`MOVE move move [] [$DATA/WORK/$DATE/LOADED.TXT] "$DATA/MASTER/KEYED.TXT" "" <nil> [] 0`,
	}
	if !slices.Equal(got, want) {
		t.Errorf("steps\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	if in := s.Steps[0].Bindings["IN"]; in != "$DATA/IN/IN.$DATE.TXT" {
		t.Errorf("DD_IN bound to %q", in)
	}
	if !maps.Equal(s.Parameters, map[string]string{"DATE": "1", "DATA": "GLS_DATA"}) {
		t.Errorf("parameters %v", s.Parameters)
	}
	values := map[string]string{"1": "20261001", "GLS_DATA": "/data", "GLS_BIN": "/bin"}
	value := func(parameter string) string { return values[parameter] }
	if got, want := s.Expand(s.Steps[0].Bindings["OUT"], value), "/data/WORK/20261001/LOADED.TXT"; got != want {
		t.Errorf("DD_OUT expands to %q, want %q", got, want)
	}
	if got, want := s.Expand(s.Steps[0].Command, value), "/bin/LOAD.EXE"; got != want {
		t.Errorf("the command expands to %q, want %q", got, want)
	}
}

func TestParseErrors(t *testing.T) {
	for name, source := range map[string]string{
		"a set without a value":     "set X\r\n",
		"an if of something else":   "if exist X call Y\r\n",
		"an echo to the console":    "echo hello\r\n",
		"a move without its target": "move X\r\n",
		"a quote not closed":        "call \"X\r\n",
	} {
		if _, err := Parse("bad", strings.NewReader(source)); err == nil {
			t.Errorf("%s: parsed without error", name)
		}
	}
}

// A run sets the variables, makes the directories, binds the DD names,
// runs the programs with their output redirected, echoes into the log,
// calls the other batch file, moves the file, and exits as the tests of
// errorlevel say.
func TestRun(t *testing.T) {
	bin, scripts, data := t.TempDir(), t.TempDir(), t.TempDir()
	program := "#!/bin/sh\necho \"LOAD $1 $DD_IN $(cat \"$DD_IN\")\"\necho loaded > \"$DD_OUT\"\nexit 4\n"
	if err := os.WriteFile(filepath.Join(bin, "LOAD"), []byte(program), 0o755); err != nil {
		t.Fatal(err)
	}
	more := "@echo off\r\nset LOG=%GLS_DATA%\\LOGS\\DAILY.%1.LOG\r\n>>\"%LOG%\" echo MORE %1 %2\r\nexit /b 0\r\n"
	for name, content := range map[string]string{"DAILY.BAT": daily, "MORE.BAT": more} {
		if err := os.WriteFile(filepath.Join(scripts, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, dir := range []string{"IN", "LOGS", "REPORTS", "MASTER"} {
		if err := os.Mkdir(filepath.Join(data, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(data, "IN", "IN.20261001.TXT"), []byte("records"), 0o644); err != nil {
		t.Fatal(err)
	}

	r := Runner{Env: map[string]string{"GLS_BIN": bin, "GLS_SCRIPTS": scripts, "GLS_DATA": data}}
	rc, err := r.Run(filepath.Join(scripts, "DAILY.BAT"), "20261001")
	if err != nil {
		t.Fatal(err)
	}
	if rc != 0 {
		t.Errorf("errorlevel %d, want 0", rc)
	}
	for name, want := range map[string]string{
		"LOGS/DAILY.20261001.LOG": "LOAD LOAD 4\nMORE 20261001 \n",
		"REPORTS/LOAD.OUT":        "LOAD 20261001 " + filepath.Join(data, "IN", "IN.20261001.TXT") + " records\n",
		"MASTER/KEYED.TXT":        "loaded\n",
	} {
		got, err := os.ReadFile(filepath.Join(data, name))
		if err != nil || string(got) != want {
			t.Errorf("%s is %q, %v, want %q", name, got, err, want)
		}
	}
	if _, err := os.Stat(filepath.Join(data, "WORK", "20261001", "LOADED.TXT")); !os.IsNotExist(err) {
		t.Errorf("the loaded file was not moved: %v", err)
	}

	// A program that ends with 8 stops the file, which exits with the code.
	if err := os.WriteFile(filepath.Join(bin, "LOAD"), []byte("#!/bin/sh\nexit 12\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if rc, err := r.Run(filepath.Join(scripts, "DAILY.BAT"), "20261002"); err != nil || rc != 12 {
		t.Errorf("errorlevel %d, %v, want 12 from the step that failed", rc, err)
	}
	if _, err := os.Stat(filepath.Join(data, "REPORTS", "MORE.OUT")); !os.IsNotExist(err) {
		t.Error("the file went on after the step that failed")
	}
}
