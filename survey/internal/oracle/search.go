package oracle

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	"github.com/vlumi/copal/survey/blueprint"
)

// maxRounds bounds the search: each round tries every lever on each case
// the round before kept.
const maxRounds = 4

// Result is the cases of a program: the units whose decision tables they
// reach, the cases, and for each case, what its run reached, wrote and
// returned; and for each rule reached, the first case that reaches it.
type Result struct {
	Program string
	Units   []Unit
	Cases   []Case
	First   map[rule]int
	runs    []result
	program *program
	runner  *runner
}

// Rules counts the rules of a unit's tables, and those reached.
func (r *Result) Rules(unit int) (reached, total int) {
	for i, t := range r.Units[unit].Tables {
		for j := range t.Rules {
			total++
			if _, ok := r.First[rule{unit, i, j}]; ok {
				reached++
			}
		}
	}
	return reached, total
}

// Search finds cases that reach the rules of a main program's decision
// tables, and of the programs it calls, as few as reach them all, and
// records what each writes. The system's programs are compiled from the
// directory of the program's source, in a scratch directory that other
// searches of the same system may share.
func Search(system *blueprint.System, bp *blueprint.Program, copybookDir, scratch string) (*Result, error) {
	s, err := newSearcher(system, bp, copybookDir, scratch)
	if err != nil {
		return nil, err
	}
	found, err := s.search()
	if err != nil {
		return nil, err
	}
	return s.record(fewest(found))
}

// Replay runs a program on the cases its Files were written from, in the
// program's directory, and records what each reaches and writes.
func Replay(system *blueprint.System, bp *blueprint.Program, copybookDir, scratch, dir string) (*Result, error) {
	s, err := newSearcher(system, bp, copybookDir, scratch)
	if err != nil {
		return nil, err
	}
	cases, err := s.program.readCases(dir)
	if err != nil {
		return nil, err
	}
	return s.record(cases)
}

type searcher struct {
	program *program
	runner  *runner
	scratch string
	runs    int
}

func newSearcher(system *blueprint.System, bp *blueprint.Program, copybookDir, scratch string) (*searcher, error) {
	p, err := prepare(system, bp, copybookDir)
	if err != nil {
		return nil, err
	}
	r, err := newRunner(p, filepath.Dir(bp.Source), copybookDir, scratch)
	if err != nil {
		return nil, err
	}
	return &searcher{program: p, runner: r, scratch: filepath.Join(scratch, bp.Name)}, nil
}

// found is a case the search kept, and the rules its run reached.
type found struct {
	c       Case
	reached map[rule]bool
}

func (s *searcher) search() ([]found, error) {
	var kept []found
	reached := map[rule]bool{}
	keep := func(candidates []Case) ([]Case, error) {
		results, err := s.runAll(candidates, false)
		if err != nil {
			return nil, err
		}
		var fresh []Case
		for i, res := range results {
			added := false
			for r := range res.reached {
				if !reached[r] {
					reached[r], added = true, true
				}
			}
			if added {
				kept = append(kept, found{candidates[i], res.reached})
				fresh = append(fresh, candidates[i])
			}
		}
		return fresh, nil
	}
	levers := s.program.levers()
	frontier, err := keep([]Case{s.program.base()})
	if err != nil {
		return nil, err
	}
	for range maxRounds {
		var candidates []Case
		for _, c := range frontier {
			for _, l := range levers {
				if next, ok := l.apply(s.program, c); ok {
					candidates = append(candidates, next)
				}
			}
		}
		if frontier, err = keep(candidates); err != nil {
			return nil, err
		}
		if len(frontier) == 0 {
			break
		}
	}
	return kept, nil
}

// fewest picks, from the cases found, as few as reach every rule they
// reach between them: each time, the case that reaches the most rules no
// case picked before reaches, the earliest found on a tie. The cases
// picked keep the order they were found in.
func fewest(cases []found) []Case {
	covered := map[rule]bool{}
	picked := make([]bool, len(cases))
	for {
		best, gain := -1, 0
		for i, f := range cases {
			n := 0
			for r := range f.reached {
				if !covered[r] {
					n++
				}
			}
			if !picked[i] && n > gain {
				best, gain = i, n
			}
		}
		if best < 0 {
			break
		}
		picked[best] = true
		for r := range cases[best].reached {
			covered[r] = true
		}
	}
	var out []Case
	for i, f := range cases {
		if picked[i] {
			out = append(out, f.c)
		}
	}
	return out
}

// record runs the cases again, recording what each writes, and notes the
// first case that reaches each rule.
func (s *searcher) record(cases []Case) (*Result, error) {
	runs, err := s.runAll(cases, true)
	if err != nil {
		return nil, err
	}
	res := &Result{Program: s.program.blueprint.Name, Units: s.program.units, Cases: cases, First: map[rule]int{},
		runs: runs, program: s.program, runner: s.runner}
	for i, run := range runs {
		for r := range run.reached {
			if _, ok := res.First[r]; !ok {
				res.First[r] = i
			}
		}
	}
	return res, nil
}

func (s *searcher) runAll(cases []Case, record bool) ([]result, error) {
	results := make([]result, len(cases))
	errs := make([]error, len(cases))
	work := make(chan int)
	var wg sync.WaitGroup
	for range runtime.NumCPU() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range work {
				dir := filepath.Join(s.scratch, "runs", strconv.Itoa(s.runs+i))
				if errs[i] = os.MkdirAll(dir, 0o755); errs[i] == nil {
					results[i], errs[i] = s.runner.run(cases[i], dir, record)
					_ = os.RemoveAll(dir)
				}
			}
		}()
	}
	for i := range cases {
		work <- i
	}
	close(work)
	wg.Wait()
	s.runs += len(cases)
	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}
