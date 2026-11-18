package generate

// Policy is what the generator knows of a policy: its record in the
// master, as the extract or the daily unload has it, and its coverages,
// the main contract first.
type Policy struct {
	No            string
	Status        int // 1 in force, 2 in grace, 3 lapsed, 4 surrendered, 5 matured, 6 paid out
	Product       string
	Sum           int64
	Premium       int64
	Mode          int // 1 monthly, 2 annual, 3 half-yearly
	Contract      Date
	Term          int
	Maturity      Date
	Year          int   // the policy year, from 1
	PaidTo        int64 // yyyymm
	Arrears       int
	GraceStart    Date
	LapseDate     Date
	StatusDate    Date
	Name          string // kanji, Shift_JIS
	Kana          string
	Birth         Date
	Sex           int
	EntryAge      int
	Address       string
	Beneficiary   string
	Account       account
	Agent         string
	Loan          int64
	LoanSince     Date // the day interest was last added to the loan
	Reinstated    Date // the day the policy was last reinstated
	Reserve       int64
	Dividends     int64
	LastPaidMonth int64
	LastPaid      int64
	Coverages     []Coverage
}

// Coverage is a coverage of a policy, as the coverage file has it: the
// main contract, numbered 1, or a rider.
type Coverage struct {
	No       int
	Product  string
	Status   int // 1 in force, 2 expired, 3 cancelled, 4 ended by a death
	Amount   int64
	Premium  int64
	Contract Date
	Term     int
	Maturity Date
	RatedAge int // the age its premium is rated at: the entry age, or the age at its last renewal
	Ended    Date
}

// InForce reports whether the coverage has not ended. A lapsed policy's
// coverages stay in force: the policy's state comes first.
func (c Coverage) InForce() bool { return c.Status == 1 }

// Riders are the policy's riders in force.
func (p Policy) Riders() []Coverage {
	var riders []Coverage
	for _, c := range p.Coverages {
		if c.No > 1 && c.InForce() {
			riders = append(riders, c)
		}
	}
	return riders
}

// Hospital reports whether a coverage in force pays hospitalizations: a
// medical policy's main contract, or a hospitalization rider.
func (p Policy) Hospital() bool {
	for _, c := range p.Coverages {
		if c.InForce() && productOf(c.Product).kind == 4 {
			return true
		}
	}
	return false
}

// Covered reports whether the policy is in force or in grace: insured,
// and paying premiums.
func (p Policy) Covered() bool { return p.Status == 1 || p.Status == 2 }

// cashValue is what the policy would pay on surrender before any loan is
// taken off it, as the cash value table has it.
func (p Policy) cashValue() int64 {
	for _, pr := range products {
		if pr.code == p.Product && p.Year > 0 {
			return p.Sum / 1000 * int64(pr.cashValue(ageBand(p.EntryAge), min(p.Year, 40))) / 1000
		}
	}
	return 0
}

// Lapsed reports whether the policy has lapsed and may be reinstated.
func (p Policy) Lapsed() bool { return p.Status == 3 }

// eraOrBlank is a date in its era, or blank for no date, as the master
// holds the dates a policy has not reached.
func eraOrBlank(d Date) string {
	if d.IsZero() {
		return ""
	}
	return d.Era()
}

// dateOrZero reads an era date that may be blank.
func dateOrZero(s string) (Date, error) {
	if s == "" {
		return Date{}, nil
	}
	return DateFromEra(s)
}

// record writes a policy as a record of the master.
func (g *Generator) record(p Policy) ([]byte, error) {
	return made(g.Layouts.Policy, map[string]any{
		"PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ": p.No, "PM-ｼﾞｮｳﾀｲ": string(rune('0' + p.Status)), "PM-ｼｮｳﾋﾝ": p.Product,
		"PM-ﾎｹﾝｷﾝｶﾞｸ": p.Sum, "PM-ﾎｹﾝﾘｮｳ": p.Premium, "PM-ﾊﾗｲｺﾐﾎｳﾎｳ": string(rune('0' + p.Mode)),
		"PM-ｹｲﾔｸﾋﾞ": p.Contract.Era(), "PM-ﾎｹﾝｷｶﾝ": p.Term, "PM-ﾏﾝｷﾋﾞ": eraOrBlank(p.Maturity),
		"PM-ｹｲｶﾈﾝｽｳ": p.Year, "PM-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ": p.PaidTo, "PM-ﾐﾉｳｶｲｽｳ": p.Arrears,
		"PM-ﾕｳﾖｶｲｼﾋﾞ": eraOrBlank(p.GraceStart), "PM-ｼｯｺｳﾋﾞ": eraOrBlank(p.LapseDate), "PM-ｼﾞｮｳﾀｲﾋﾂﾞｹ": p.StatusDate.Era(),
		"PM-ｼﾒｲ": p.Name, "PM-ｼﾒｲｶﾅ": p.Kana, "PM-ｾｲﾈﾝｶﾞｯﾋﾟ": p.Birth.Era(), "PM-ｾｲﾍﾞﾂ": p.Sex, "PM-ｶﾆｭｳﾈﾝﾚｲ": p.EntryAge,
		"PM-ｼﾞｭｳｼｮ": p.Address, "PM-ｳｹﾄﾘﾆﾝ": p.Beneficiary,
		"PM-ｷﾞﾝｺｳｺｰﾄﾞ": p.Account.bank, "PM-ｼﾃﾝｺｰﾄﾞ": p.Account.branch, "PM-ﾖｷﾝｼｭﾓｸ": p.Account.kind,
		"PM-ｺｳｻﾞﾊﾞﾝｺﾞｳ": p.Account.number, "PM-ｺｳｻﾞﾒｲｷﾞ": p.Account.holder, "PM-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ": p.Agent,
		"PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ": p.Loan, "PM-ｾｷﾆﾝｼﾞｭﾝﾋﾞｷﾝ": p.Reserve, "PM-ﾊｲﾄｳﾙｲｹｲ": p.Dividends,
		"PM-ｻｲｼｭｳﾆｭｳｷﾝﾈﾝｹﾞﾂ": p.LastPaidMonth, "PM-ｻｲｼｭｳﾆｭｳｷﾝｶﾞｸ": p.LastPaid, "PM-ｶｼﾂｹｷｻﾝﾋﾞ": eraOrBlank(p.LoanSince),
		"PM-ﾌｯｶﾂﾋﾞ": eraOrBlank(p.Reinstated),
	})
}

// coverageRecords writes a policy's coverages as records of the coverage
// file.
func (g *Generator) coverageRecords(p Policy) ([][]byte, error) {
	var recs [][]byte
	for _, c := range p.Coverages {
		rec, err := made(g.Layouts.Coverage, map[string]any{
			"HS-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ": p.No, "HS-ﾎｼｮｳﾊﾞﾝｺﾞｳ": c.No, "HS-ｼｮｳﾋﾝ": c.Product, "HS-ｼﾞｮｳﾀｲ": string(rune('0' + c.Status)),
			"HS-ﾎｹﾝｷﾝｶﾞｸ": c.Amount, "HS-ﾎｹﾝﾘｮｳ": c.Premium, "HS-ｹｲﾔｸﾋﾞ": c.Contract.Era(), "HS-ﾎｹﾝｷｶﾝ": c.Term,
			"HS-ﾏﾝｷﾋﾞ": eraOrBlank(c.Maturity), "HS-ｶﾆｭｳﾈﾝﾚｲ": c.RatedAge, "HS-ｼｮｳﾒﾂﾋﾞ": eraOrBlank(c.Ended),
		})
		if err != nil {
			return nil, err
		}
		recs = append(recs, rec)
	}
	return recs, nil
}

// ReadCoverages gives policies their coverages from records of the
// coverage file, as the extract or the daily unload has them.
func (g *Generator) ReadCoverages(policies []Policy, records [][]byte) error {
	byNo := make(map[string]*Policy, len(policies))
	for i := range policies {
		byNo[policies[i].No] = &policies[i]
	}
	for _, r := range records {
		var c Coverage
		var no, status, contract, maturity, ended string
		if err := get(g.Layouts.Coverage, r, map[string]any{
			"HS-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ": &no, "HS-ﾎｼｮｳﾊﾞﾝｺﾞｳ": &c.No, "HS-ｼｮｳﾋﾝ": &c.Product, "HS-ｼﾞｮｳﾀｲ": &status,
			"HS-ﾎｹﾝｷﾝｶﾞｸ": &c.Amount, "HS-ﾎｹﾝﾘｮｳ": &c.Premium, "HS-ｹｲﾔｸﾋﾞ": &contract, "HS-ﾎｹﾝｷｶﾝ": &c.Term,
			"HS-ﾏﾝｷﾋﾞ": &maturity, "HS-ｶﾆｭｳﾈﾝﾚｲ": &c.RatedAge, "HS-ｼｮｳﾒﾂﾋﾞ": &ended,
		}); err != nil {
			return err
		}
		c.Status = int(status[0] - '0')
		var err error
		for _, d := range []struct {
			target *Date
			text   string
		}{{&c.Contract, contract}, {&c.Maturity, maturity}, {&c.Ended, ended}} {
			if *d.target, err = dateOrZero(d.text); err != nil {
				return err
			}
		}
		if p, ok := byNo[no]; ok {
			p.Coverages = append(p.Coverages, c)
		}
	}
	return nil
}

// ReadPolicies reads policies from master records, as the extract or the
// daily unload has them.
func (g *Generator) ReadPolicies(records [][]byte) ([]Policy, error) {
	policies := make([]Policy, 0, len(records))
	for _, r := range records {
		var p Policy
		var status, mode, contract, maturity, grace, lapse, statusDate, birth, loanSince, reinstated string
		if err := get(g.Layouts.Policy, r, map[string]any{
			"PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ": &p.No, "PM-ｼﾞｮｳﾀｲ": &status, "PM-ｼｮｳﾋﾝ": &p.Product,
			"PM-ﾎｹﾝｷﾝｶﾞｸ": &p.Sum, "PM-ﾎｹﾝﾘｮｳ": &p.Premium, "PM-ﾊﾗｲｺﾐﾎｳﾎｳ": &mode,
			"PM-ｹｲﾔｸﾋﾞ": &contract, "PM-ﾎｹﾝｷｶﾝ": &p.Term, "PM-ﾏﾝｷﾋﾞ": &maturity,
			"PM-ｹｲｶﾈﾝｽｳ": &p.Year, "PM-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ": &p.PaidTo, "PM-ﾐﾉｳｶｲｽｳ": &p.Arrears,
			"PM-ﾕｳﾖｶｲｼﾋﾞ": &grace, "PM-ｼｯｺｳﾋﾞ": &lapse, "PM-ｼﾞｮｳﾀｲﾋﾂﾞｹ": &statusDate,
			"PM-ｼﾒｲ": &p.Name, "PM-ｼﾒｲｶﾅ": &p.Kana, "PM-ｾｲﾈﾝｶﾞｯﾋﾟ": &birth, "PM-ｾｲﾍﾞﾂ": &p.Sex, "PM-ｶﾆｭｳﾈﾝﾚｲ": &p.EntryAge,
			"PM-ｼﾞｭｳｼｮ": &p.Address, "PM-ｳｹﾄﾘﾆﾝ": &p.Beneficiary,
			"PM-ｷﾞﾝｺｳｺｰﾄﾞ": &p.Account.bank, "PM-ｼﾃﾝｺｰﾄﾞ": &p.Account.branch, "PM-ﾖｷﾝｼｭﾓｸ": &p.Account.kind,
			"PM-ｺｳｻﾞﾊﾞﾝｺﾞｳ": &p.Account.number, "PM-ｺｳｻﾞﾒｲｷﾞ": &p.Account.holder, "PM-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ": &p.Agent,
			"PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ": &p.Loan, "PM-ｾｷﾆﾝｼﾞｭﾝﾋﾞｷﾝ": &p.Reserve, "PM-ﾊｲﾄｳﾙｲｹｲ": &p.Dividends,
			"PM-ｻｲｼｭｳﾆｭｳｷﾝﾈﾝｹﾞﾂ": &p.LastPaidMonth, "PM-ｻｲｼｭｳﾆｭｳｷﾝｶﾞｸ": &p.LastPaid, "PM-ｶｼﾂｹｷｻﾝﾋﾞ": &loanSince,
			"PM-ﾌｯｶﾂﾋﾞ": &reinstated,
		}); err != nil {
			return nil, err
		}
		p.Status, p.Mode = int(status[0]-'0'), int(mode[0]-'0')
		var err error
		for _, d := range []struct {
			target *Date
			text   string
		}{{&p.Contract, contract}, {&p.Maturity, maturity}, {&p.GraceStart, grace}, {&p.LapseDate, lapse}, {&p.StatusDate, statusDate}, {&p.Birth, birth}, {&p.LoanSince, loanSince}, {&p.Reinstated, reinstated}} {
			if *d.target, err = dateOrZero(d.text); err != nil {
				return nil, err
			}
		}
		policies = append(policies, p)
	}
	return policies, nil
}
