package oracle

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Coverage is a program's cases for tools: how each run ended, and for
// each decision table of the program and of those it calls, the case that
// reaches each rule first.
type Coverage struct {
	Program string         `json:"program"`
	Cases   []CaseOutcome  `json:"cases"`
	Units   []UnitCoverage `json:"units"`
}

// CaseOutcome is how a case's run ended.
type CaseOutcome struct {
	ReturnCode int  `json:"returnCode"`
	Abended    bool `json:"abended,omitempty"`
	TimedOut   bool `json:"timedOut,omitempty"`
}

// UnitCoverage is the coverage of a program's decision tables.
type UnitCoverage struct {
	Program string          `json:"program"`
	Tables  []TableCoverage `json:"tables"`
}

// TableCoverage is a decision table's coverage: for each rule, the number
// of the first case that reaches it, from 1, or 0 for none.
type TableCoverage struct {
	Paragraph string `json:"paragraph"`
	Loop      string `json:"loop,omitempty"`
	ReachedBy []int  `json:"reachedBy"`
}

// Reached counts the rules of a table some case reaches.
func (t TableCoverage) Reached() int {
	n := 0
	for _, c := range t.ReachedBy {
		if c > 0 {
			n++
		}
	}
	return n
}

func (r *Result) coverageJSON() []byte {
	c := Coverage{Program: r.Program}
	for _, run := range r.runs {
		c.Cases = append(c.Cases, CaseOutcome{ReturnCode: run.outcome.RC, Abended: run.outcome.Abended, TimedOut: run.outcome.TimedOut})
	}
	for u, unit := range r.Units {
		uc := UnitCoverage{Program: unit.Program}
		for t, table := range unit.Tables {
			tc := TableCoverage{Paragraph: table.Paragraph, Loop: table.Loop, ReachedBy: make([]int, len(table.Rules))}
			for i := range table.Rules {
				if first, ok := r.First[rule{u, t, i}]; ok {
					tc.ReachedBy[i] = first + 1
				}
			}
			uc.Tables = append(uc.Tables, tc)
		}
		c.Units = append(c.Units, uc)
	}
	data, _ := json.MarshalIndent(c, "", "  ")
	return append(data, '\n')
}

// ReadCoverage reads the coverage of each program whose cases are in a
// directory, by program.
func ReadCoverage(dir string) (map[string]*Coverage, error) {
	if _, err := os.Stat(dir); err != nil {
		return nil, err
	}
	paths, err := filepath.Glob(filepath.Join(dir, "*", "coverage.json"))
	if err != nil {
		return nil, err
	}
	coverage := map[string]*Coverage{}
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var c Coverage
		if err := json.Unmarshal(data, &c); err != nil {
			return nil, err
		}
		coverage[c.Program] = &c
	}
	return coverage, nil
}
