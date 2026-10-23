package generate

import (
	"hash/fnv"
)

// Results makes up the collection agent's results for a transfer request:
// the request's records back with the banks' names filled in and each
// debit's result, most paid, a few not for lack of funds, a closed
// account, the depositor's refusal or a mandate the bank has no record
// of, and the trailer counting them.
func (g *Generator) Results(date Date, request [][]byte) ([][]byte, error) {
	r := g.rng("results", date.Number())
	l := g.Layouts.Transfer
	var results [][]byte
	var paidCount, paidAmount, unpaidCount, unpaidAmount int64
	for _, req := range request {
		rec := append([]byte(nil), req...)
		kind, err := l.Number(rec, "ZG-ﾃﾞｰﾀｸﾌﾞﾝ")
		if err != nil {
			return nil, err
		}
		switch kind {
		case 1:
			bankCode, _ := l.Number(rec, "ZG-H-ｷﾞﾝｺｳｺｰﾄﾞ")
			branchCode, _ := l.Number(rec, "ZG-H-ｼﾃﾝｺｰﾄﾞ")
			bankName, branchName := bankNamed(int(bankCode), int(branchCode))
			err = set(l, rec, map[string]any{"ZG-H-ｷﾞﾝｺｳﾒｲ": bankName, "ZG-H-ｼﾃﾝﾒｲ": branchName})
		case 2:
			bankCode, _ := l.Number(rec, "ZG-D-ｷﾞﾝｺｳｺｰﾄﾞ")
			branchCode, _ := l.Number(rec, "ZG-D-ｼﾃﾝｺｰﾄﾞ")
			amount, _ := l.Number(rec, "ZG-D-ｷﾝｶﾞｸ")
			customer, _ := l.Text(rec, "ZG-D-ｺｷｬｸﾊﾞﾝｺﾞｳ")
			bankName, branchName := bankNamed(int(bankCode), int(branchCode))
			result := "0"
			if chance(r, failureRate(customer)) {
				result = []string{"1", "1", "1", "1", "1", "1", "2", "3", "3", "4"}[r.IntN(10)]
				unpaidCount, unpaidAmount = unpaidCount+1, unpaidAmount+amount
			} else {
				paidCount, paidAmount = paidCount+1, paidAmount+amount
			}
			err = set(l, rec, map[string]any{"ZG-D-ｷﾞﾝｺｳﾒｲ": bankName, "ZG-D-ｼﾃﾝﾒｲ": branchName, "ZG-D-ﾌﾘｶｴｹｯｶ": result})
		case 8:
			err = set(l, rec, map[string]any{
				"ZG-T-ﾌﾘｶｴｽﾞﾐｹﾝｽｳ": paidCount, "ZG-T-ﾌﾘｶｴｽﾞﾐｷﾝｶﾞｸ": paidAmount,
				"ZG-T-ﾌﾘｶｴﾌﾉｳｹﾝｽｳ": unpaidCount, "ZG-T-ﾌﾘｶｴﾌﾉｳｷﾝｶﾞｸ": unpaidAmount,
			})
		}
		if err != nil {
			return nil, err
		}
		results = append(results, rec)
	}
	return results, nil
}

// failureRate is the chance a policyholder's debit fails: a few of them
// fail every other month, the rest rarely.
func failureRate(customer string) float64 {
	h := fnv.New32a()
	h.Write([]byte(customer))
	if h.Sum32()%100 < 3 {
		return 0.5
	}
	return 0.015
}
