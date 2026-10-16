package bank

import (
	"path/filepath"

	"github.com/vlumi/copal/pangaea/dataset"
	"github.com/vlumi/copal/pangaea/generate"
)

// RunDay is a business day the bank has run, and where its files are. Day
// n is generation n of the transaction file, the audit trail and the
// rejects; it started from generation n of the control record and the
// master and wrote generation n+1 of both.
type RunDay struct {
	Date       generate.Date
	Generation int
	Sysout     string // the directory of the day's job log and reports
}

// Transactions is the dataset of the day's transaction file.
func (d RunDay) Transactions() string { return dataset.GenerationName("TRAN.DAILY", d.Generation) }

// Audit is the dataset of the day's audit trail.
func (d RunDay) Audit() string { return dataset.GenerationName("TRAN.AUDIT", d.Generation) }

// Rejects is the dataset of the day's rejected transactions.
func (d RunDay) Rejects() string { return dataset.GenerationName("TRAN.REJECTS", d.Generation) }

// MasterBefore is the account master the day started from.
func (d RunDay) MasterBefore() string { return dataset.GenerationName("ACCT.MASTER", d.Generation) }

// MasterAfter is the account master the day wrote.
func (d RunDay) MasterAfter() string { return dataset.GenerationName("ACCT.MASTER", d.Generation+1) }

// History is what the bank has done since its reset.
type History struct {
	Seed uint64
	// Date is the business day the bank runs next.
	Date generate.Date
	// Master is the current account master.
	Master string
	Days   []RunDay
}

// History reads the bank's history from its catalog. A day counts as run
// when its batch wrote the next day's control record, the last thing it
// does.
func (b Bank) History() (History, error) {
	seed, err := b.seed()
	if err != nil {
		return History{}, err
	}
	g, err := b.generator(seed)
	if err != nil {
		return History{}, err
	}
	store, err := dataset.Open(b.datasets())
	if err != nil {
		return History{}, err
	}
	h := History{Seed: seed, Master: latest(store, "ACCT.MASTER")}
	controls := store.LatestGeneration("BANK.CONTROL")
	for n := 1; n <= controls; n++ {
		date, err := b.controlDate(g, store, dataset.GenerationName("BANK.CONTROL", n))
		if err != nil {
			return History{}, err
		}
		if n == controls {
			h.Date = date
		} else {
			h.Days = append(h.Days, RunDay{Date: date, Generation: n, Sysout: b.sysout(date)})
		}
	}
	return h, nil
}

func (b Bank) sysout(date generate.Date) string { return filepath.Join(b.Dir, "sysout", date.String()) }
