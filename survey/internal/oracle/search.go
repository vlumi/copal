package oracle

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/vlumi/copal/survey/blueprint"
)

// maxRounds bounds the search: each round tries every lever on each case
// the round before kept.
const maxRounds = 4

// Result is what the search found for a program: the cases it kept, in
// the order it found them, the units whose decision tables they reach,
// and for each rule reached, the case that reached it first.
type Result struct {
	Program string
	Cases   []Case
	Units   []Unit
	First   map[rule]int
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
// tables, and of the programs it calls. The system's programs are
// compiled from the directory of the program's source, and the work is
// done in a scratch directory.
func Search(system *blueprint.System, bp *blueprint.Program, copybookDir, scratch string) (*Result, error) {
	p, err := prepare(system, bp, copybookDir)
	if err != nil {
		return nil, err
	}
	r, err := newRunner(p, filepath.Dir(bp.Source), copybookDir, scratch)
	if err != nil {
		return nil, err
	}
	s := &searcher{program: p, runner: r, scratch: scratch,
		result: &Result{Program: bp.Name, Units: p.units, First: map[rule]int{}}}
	return s.result, s.search()
}

type searcher struct {
	program *program
	runner  *runner
	scratch string
	runs    int
	result  *Result
}

func (s *searcher) search() error {
	levers := s.program.levers()
	frontier := []Case{s.program.base()}
	if _, err := s.keep(frontier); err != nil {
		return err
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
		kept, err := s.keep(candidates)
		if err != nil {
			return err
		}
		if len(kept) == 0 {
			break
		}
		frontier = kept
	}
	return nil
}

// keep runs the candidates, in parallel, and keeps, in their order, those
// that reach a rule that no case before them reached.
func (s *searcher) keep(candidates []Case) ([]Case, error) {
	results, err := s.runAll(candidates)
	if err != nil {
		return nil, err
	}
	var kept []Case
	for i, res := range results {
		var fresh []rule
		for r := range res.reached {
			if _, ok := s.result.First[r]; !ok {
				fresh = append(fresh, r)
			}
		}
		if len(fresh) == 0 {
			continue
		}
		for _, r := range fresh {
			s.result.First[r] = len(s.result.Cases)
		}
		s.result.Cases = append(s.result.Cases, candidates[i])
		kept = append(kept, candidates[i])
	}
	return kept, nil
}

func (s *searcher) runAll(candidates []Case) ([]result, error) {
	results := make([]result, len(candidates))
	errs := make([]error, len(candidates))
	work := make(chan int)
	var wg sync.WaitGroup
	for range runtime.NumCPU() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range work {
				dir := filepath.Join(s.scratch, "runs", fmt.Sprint(s.runs+i))
				if errs[i] = os.MkdirAll(dir, 0o755); errs[i] == nil {
					results[i], errs[i] = s.runner.run(candidates[i], dir)
					_ = os.RemoveAll(dir)
				}
			}
		}()
	}
	for i := range candidates {
		work <- i
	}
	close(work)
	wg.Wait()
	s.runs += len(candidates)
	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}
