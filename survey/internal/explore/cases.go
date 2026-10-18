package explore

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/vlumi/copal/survey/internal/oracle"
)

// CasesView is the coverage of a program's decision tables by the cases
// of a main program: its own, or its caller's.
type CasesView struct {
	Main           string
	Cases          int
	Unit           oracle.UnitCoverage
	Reached, Total int
}

// casesOf finds the coverage of a program's tables: by its own cases, or
// by those of a main program that calls it.
func (s *Server) casesOf(program string) *CasesView {
	for _, main := range s.mainsCovering(program) {
		c := s.cases[main]
		if c == nil {
			continue
		}
		for _, u := range c.Units {
			if u.Program != program {
				continue
			}
			v := &CasesView{Main: main, Cases: len(c.Cases), Unit: u}
			for _, t := range u.Tables {
				v.Reached += t.Reached()
				v.Total += len(t.ReachedBy)
			}
			return v
		}
	}
	return nil
}

// mainsCovering lists the main programs whose cases may cover a program:
// the program itself first, then the others with cases, by name.
func (s *Server) mainsCovering(program string) []string {
	mains := []string{program}
	for _, p := range s.index.system.Programs {
		if p.Name != program && s.cases[p.Name] != nil {
			mains = append(mains, p.Name)
		}
	}
	return mains
}

// CaseView is one case of a main program: how its run ended, and its
// input and output files, as they were recorded.
type CaseView struct {
	Program string
	Number  int
	Outcome oracle.CaseOutcome
	Inputs  []CaseFile
	Outputs []CaseFile
}

// CaseFile is a file of a case, by name.
type CaseFile struct {
	Name string
	Text string
}

func (s *Server) caseView(r *http.Request) (any, error) {
	p, err := s.lookupProgram(r)
	if err != nil {
		return nil, err
	}
	c, ok := s.cases[p.Name]
	var n int
	if _, err := fmt.Sscanf(r.PathValue("case"), "%d", &n); !ok || err != nil || n < 1 || n > len(c.Cases) {
		return nil, fmt.Errorf("case %s of %s: %w", r.PathValue("case"), p.Name, errNotFound)
	}
	v := CaseView{Program: p.Name, Number: n, Outcome: c.Cases[n-1]}
	dir := filepath.Join(s.casesDir, p.Name, fmt.Sprintf("%02d", n))
	if v.Inputs, err = caseFiles(filepath.Join(dir, "input")); err != nil {
		return nil, err
	}
	v.Outputs, err = caseFiles(filepath.Join(dir, "output"))
	return v, err
}

func caseFiles(dir string) ([]CaseFile, error) {
	paths, err := filepath.Glob(filepath.Join(dir, "*.txt"))
	if err != nil {
		return nil, err
	}
	var files []CaseFile
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		files = append(files, CaseFile{Name: strings.TrimSuffix(filepath.Base(path), ".txt"), Text: string(data)})
	}
	return files, nil
}

func caseLink(program string, n int) string {
	return fmt.Sprintf("/programs/%s/cases/%02d", program, n)
}
