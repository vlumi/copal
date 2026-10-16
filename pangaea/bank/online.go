package bank

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"sync"

	"github.com/vlumi/copal/pangaea/dataset"
	"github.com/vlumi/copal/pangaea/generate"
	"github.com/vlumi/copal/pangaea/runner"
)

// ErrOnlineDayRun is for opening the online day of a business day whose
// online day has already run: its batch is next.
var ErrOnlineDayRun = errors.New("the business day's online day has run; its batch is next")

// ErrOnlineDayClosed is for a request to an online day that has ended.
var ErrOnlineDayClosed = errors.New("the online day has ended")

// OnlineDay is a business day's online day: the ONLINE job, running with
// its in-stream data open, so that it authorizes requests one at a time
// until it is closed, when it extracts the day's holds for the batch.
type OnlineDay struct {
	Date generate.Date

	mu        sync.Mutex
	requests  *io.PipeWriter
	responses *bufio.Scanner
	done      chan struct{}
	results   []runner.Result
	err       error
}

// OpenDay starts the online day of the bank's business date.
func (b Bank) OpenDay() (*OnlineDay, error) {
	g, store, err := b.open()
	if err != nil {
		return nil, err
	}
	if onlineDayRun(store) {
		return nil, ErrOnlineDayRun
	}
	date, err := b.businessDate(g, store)
	if err != nil {
		return nil, err
	}
	job, err := b.job("ONLINE")
	if err != nil {
		return nil, err
	}
	inStream, requests := io.Pipe()
	responses, output := io.Pipe()
	d := &OnlineDay{Date: date, requests: requests, responses: bufio.NewScanner(responses), done: make(chan struct{})}
	r := b.runner(store, filepath.Join(b.sysout(date), job.Name))
	r.InStream, r.Output = inStream, output
	go func() {
		d.results, d.err = runJob(r, job)
		_ = output.Close()
		_ = inStream.Close()
		close(d.done)
	}()
	return d, nil
}

// Authorize sends a request, an AUTHMSG request record, and returns the
// response record.
func (d *OnlineDay) Authorize(request []byte) ([]byte, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	select {
	case <-d.done:
		return nil, ErrOnlineDayClosed
	default:
	}
	if _, err := d.requests.Write(append(slices.Clone(request), '\n')); err != nil {
		return nil, ErrOnlineDayClosed
	}
	if !d.responses.Scan() {
		return nil, fmt.Errorf("%w: no response", ErrOnlineDayClosed)
	}
	return slices.Clone(d.responses.Bytes()), nil
}

// Close ends the online day, as the batch window opens, and returns how
// the ONLINE job's steps ended.
func (d *OnlineDay) Close() ([]runner.Result, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	_ = d.requests.Close()
	<-d.done
	return d.results, d.err
}

// onlineDayRun tells whether the business day's online day has run: day n
// writes generation n of the card file, and starts from generation n of
// the control record.
func onlineDayRun(store *dataset.Store) bool {
	return store.LatestGeneration("TRAN.CARD") == store.LatestGeneration("BANK.CONTROL")
}
