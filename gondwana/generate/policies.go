package generate

// Policy is what the generator knows of a policy: its record in the
// master, as the extract or the daily unload has it.
type Policy struct {
	No            string
	Status        int // 1 in force, 2 in grace, 3 lapsed, 4 surrendered, 5 matured, 6 paid out
	Product       string
	Sum           int64
	Premium       int64
	Mode          int // 1 monthly, 2 annual
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
	Reserve       int64
	Dividends     int64
	LastPaidMonth int64
	LastPaid      int64
}

// Covered reports whether the policy is in force or in grace: insured,
// and paying premiums.
func (p Policy) Covered() bool { return p.Status == 1 || p.Status == 2 }

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
		"PM-ｻｲｼｭｳﾆｭｳｷﾝﾈﾝｹﾞﾂ": p.LastPaidMonth, "PM-ｻｲｼｭｳﾆｭｳｷﾝｶﾞｸ": p.LastPaid,
	})
}

// ReadPolicies reads policies from master records, as the extract or the
// daily unload has them.
func (g *Generator) ReadPolicies(records [][]byte) ([]Policy, error) {
	policies := make([]Policy, 0, len(records))
	for _, r := range records {
		var p Policy
		var status, mode, contract, maturity, grace, lapse, statusDate, birth string
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
			"PM-ｻｲｼｭｳﾆｭｳｷﾝﾈﾝｹﾞﾂ": &p.LastPaidMonth, "PM-ｻｲｼｭｳﾆｭｳｷﾝｶﾞｸ": &p.LastPaid,
		}); err != nil {
			return nil, err
		}
		p.Status, p.Mode = int(status[0]-'0'), int(mode[0]-'0')
		var err error
		for _, d := range []struct {
			target *Date
			text   string
		}{{&p.Contract, contract}, {&p.Maturity, maturity}, {&p.GraceStart, grace}, {&p.LapseDate, lapse}, {&p.StatusDate, statusDate}, {&p.Birth, birth}} {
			if *d.target, err = dateOrZero(d.text); err != nil {
				return nil, err
			}
		}
		policies = append(policies, p)
	}
	return policies, nil
}
