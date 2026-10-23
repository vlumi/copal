package generate

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

// name is a Japanese name as the files hold it: in kanji and in
// half-width katakana, both Shift_JIS.
type name struct{ kanji, kana string }

// familyNames are the family names the policyholders are drawn from.
var familyNames = encoded([]name{
	{"山田", "ﾔﾏﾀﾞ"}, {"佐藤", "ｻﾄｳ"}, {"鈴木", "ｽｽﾞｷ"}, {"高橋", "ﾀｶﾊｼ"}, {"田中", "ﾀﾅｶ"},
	{"伊藤", "ｲﾄｳ"}, {"渡辺", "ﾜﾀﾅﾍﾞ"}, {"中村", "ﾅｶﾑﾗ"}, {"小林", "ｺﾊﾞﾔｼ"}, {"加藤", "ｶﾄｳ"},
	{"吉田", "ﾖｼﾀﾞ"}, {"山本", "ﾔﾏﾓﾄ"}, {"松本", "ﾏﾂﾓﾄ"}, {"井上", "ｲﾉｳｴ"}, {"木村", "ｷﾑﾗ"},
	{"林", "ﾊﾔｼ"}, {"清水", "ｼﾐｽﾞ"}, {"山口", "ﾔﾏｸﾞﾁ"}, {"森", "ﾓﾘ"}, {"池田", "ｲｹﾀﾞ"},
	{"橋本", "ﾊｼﾓﾄ"}, {"石川", "ｲｼｶﾜ"}, {"前田", "ﾏｴﾀﾞ"}, {"藤田", "ﾌｼﾞﾀ"}, {"岡田", "ｵｶﾀﾞ"},
	{"後藤", "ｺﾞﾄｳ"}, {"長谷川", "ﾊｾｶﾞﾜ"}, {"村上", "ﾑﾗｶﾐ"}, {"近藤", "ｺﾝﾄﾞｳ"}, {"石井", "ｲｼｲ"},
	{"斎藤", "ｻｲﾄｳ"}, {"坂本", "ｻｶﾓﾄ"}, {"遠藤", "ｴﾝﾄﾞｳ"}, {"青木", "ｱｵｷ"}, {"藤井", "ﾌｼﾞｲ"},
	{"西村", "ﾆｼﾑﾗ"}, {"福田", "ﾌｸﾀﾞ"}, {"太田", "ｵｵﾀ"}, {"三浦", "ﾐｳﾗ"}, {"岡本", "ｵｶﾓﾄ"},
})

// givenNames are the given names, by sex: 1 for a man, 2 for a woman, as
// the records code it.
var givenNames = map[int][]name{
	1: encoded([]name{
		{"太郎", "ﾀﾛｳ"}, {"一郎", "ｲﾁﾛｳ"}, {"健", "ｹﾝ"}, {"修", "ｵｻﾑ"}, {"誠", "ﾏｺﾄ"},
		{"大輔", "ﾀﾞｲｽｹ"}, {"浩二", "ｺｳｼﾞ"}, {"清", "ｷﾖｼ"}, {"隆", "ﾀｶｼ"}, {"直人", "ﾅｵﾄ"},
		{"翔太", "ｼｮｳﾀ"}, {"亮", "ﾘｮｳ"}, {"誠一", "ｾｲｲﾁ"}, {"拓也", "ﾀｸﾔ"}, {"茂", "ｼｹﾞﾙ"},
		{"健太", "ｹﾝﾀ"}, {"剛", "ﾂﾖｼ"}, {"弘", "ﾋﾛｼ"}, {"竜", "ﾘｭｳ"}, {"和夫", "ｶｽﾞｵ"},
		{"勇", "ｲｻﾑ"}, {"大樹", "ﾀﾞｲｷ"}, {"蓮", "ﾚﾝ"}, {"悠斗", "ﾕｳﾄ"}, {"正雄", "ﾏｻｵ"},
		{"哲也", "ﾃﾂﾔ"}, {"慎一", "ｼﾝｲﾁ"}, {"光", "ﾋｶﾙ"}, {"優", "ﾕｳ"}, {"昭夫", "ｱｷｵ"},
	}),
	2: encoded([]name{
		{"花子", "ﾊﾅｺ"}, {"美咲", "ﾐｻｷ"}, {"結衣", "ﾕｲ"}, {"恵子", "ｹｲｺ"}, {"幸子", "ｻﾁｺ"},
		{"和子", "ｶｽﾞｺ"}, {"さくら", "ｻｸﾗ"}, {"由美", "ﾕﾐ"}, {"典子", "ﾉﾘｺ"}, {"愛", "ｱｲ"},
		{"莉子", "ﾘｺ"}, {"美香", "ﾐｶ"}, {"真理", "ﾏﾘ"}, {"陽菜", "ﾋﾅ"}, {"葵", "ｱｵｲ"},
		{"直美", "ﾅｵﾐ"}, {"麻衣", "ﾏｲ"}, {"恵", "ﾒｸﾞﾐ"}, {"聡美", "ｻﾄﾐ"}, {"久美子", "ｸﾐｺ"},
		{"裕子", "ﾕｳｺ"}, {"明美", "ｱｹﾐ"}, {"千代", "ﾁﾖ"}, {"凛", "ﾘﾝ"}, {"彩", "ｱﾔ"},
		{"春香", "ﾊﾙｶ"}, {"智子", "ﾄﾓｺ"}, {"玲奈", "ﾚｲﾅ"}, {"文子", "ﾌﾐｺ"}, {"静香", "ｼｽﾞｶ"},
	}),
}

// encoded encodes a table of names in Shift_JIS.
func encoded(names []name) []name {
	for i, n := range names {
		names[i] = name{sjis(n.kanji), sjis(n.kana)}
	}
	return names
}

// fullWidthSpace separates a family name from a given name, in Shift_JIS.
var fullWidthSpace = sjis("　")

// person is someone named: a family name, a given name and a sex.
type person struct {
	family, given name
	sex           int
}

// Kanji is the full name in kanji, as the master holds it.
func (p person) Kanji() string { return p.family.kanji + fullWidthSpace + p.given.kanji }

// Kana is the full name in half-width katakana, as the master and the
// banks' file hold it.
func (p person) Kana() string { return p.family.kana + " " + p.given.kana }

// somebody makes up a person.
func somebody(r *rand.Rand) person {
	sex := 1 + r.IntN(2)
	return person{family: familyNames[r.IntN(len(familyNames))], given: givenNames[sex][r.IntN(len(givenNames[sex]))], sex: sex}
}

// relative makes up a relative of a person, a beneficiary: the same
// family name, usually.
func relative(r *rand.Rand, p person) person {
	q := somebody(r)
	if r.IntN(100) < 85 {
		q.family = p.family
	}
	return q
}

// kanaOf finds the kana of a name in kanji, as the branches write a
// beneficiary's name on a benefit request from the policy; a name this
// package did not make up has none.
func kanaOf(kanji string) string {
	family, given, ok := strings.Cut(kanji, fullWidthSpace)
	if !ok {
		return ""
	}
	var kana []string
	for _, n := range familyNames {
		if n.kanji == family {
			kana = append(kana, n.kana)
		}
	}
	for _, names := range givenNames {
		for _, n := range names {
			if n.kanji == given {
				kana = append(kana, n.kana)
			}
		}
	}
	if len(kana) != 2 {
		return ""
	}
	return kana[0] + " " + kana[1]
}

// addresses are where the policyholders live: a prefecture, a city, a
// ward and a district, within the 40 bytes of the field.
var addresses = encodedAll([]string{
	"東京都千代田区丸の内一丁目", "東京都世田谷区成城二丁目", "東京都港区赤坂五丁目", "東京都杉並区高円寺南",
	"神奈川県横浜市西区みなとみらい", "神奈川県川崎市中原区小杉町", "埼玉県さいたま市浦和区高砂", "千葉県千葉市中央区新町",
	"大阪府大阪市北区梅田三丁目", "大阪府堺市堺区南瓦町", "京都府京都市中京区烏丸通", "兵庫県神戸市中央区三宮町",
	"愛知県名古屋市中区栄一丁目", "静岡県静岡市葵区呉服町", "新潟県新潟市中央区万代", "宮城県仙台市青葉区一番町",
	"北海道札幌市中央区北一条西", "広島県広島市中区紙屋町", "福岡県福岡市博多区博多駅前", "熊本県熊本市中央区手取本町",
})

func encodedAll(texts []string) []string {
	for i, t := range texts {
		texts[i] = sjis(t)
	}
	return texts
}

// bank is a bank the premiums are debited at, with its code and name as
// the collection agent writes them, and its branches.
type bank struct {
	code     int
	name     string
	branches []branch
}

type branch struct {
	code int
	name string
}

// banks are the banks of the policyholders' accounts.
var banks = []bank{
	{1, sjis("ﾊﾞﾝｺﾞﾅ"), []branch{{1, sjis("ﾎﾝﾃﾝ")}, {12, sjis("ｼﾝｼﾞｭｸ")}, {34, sjis("ｳﾒﾀﾞ")}, {56, sjis("ｻｶｴ")}}},
	{5, sjis("ﾛｰﾗｼｱ"), []branch{{1, sjis("ﾎﾝﾃﾝ")}, {12, sjis("ｼﾌﾞﾔ")}, {78, sjis("ﾊｶﾀ")}}},
	{9, sjis("ﾃﾁｽ"), []branch{{1, sjis("ﾎﾝﾃﾝ")}, {123, sjis("ｻｯﾎﾟﾛ")}, {45, sjis("ｾﾝﾀﾞｲ")}}},
	{33, sjis("ｲｱﾍﾟﾄｽ"), []branch{{7, sjis("ｺｳﾍﾞ")}, {1, sjis("ﾎﾝﾃﾝ")}, {89, sjis("ﾋﾛｼﾏ")}}},
}

// account is a policyholder's bank account for the premium debits.
type account struct {
	bank, branch, kind, number int
	holder                     string
}

// anAccount makes up an account for a person, at one of the banks.
func anAccount(r *rand.Rand, holder person) account {
	b := banks[r.IntN(len(banks))]
	kind := 1
	if r.IntN(100) < 5 {
		kind = 2
	}
	return account{bank: b.code, branch: b.branches[r.IntN(len(b.branches))].code, kind: kind, number: 1000000 + r.IntN(9000000), holder: holder.Kana()}
}

// bankNamed finds a bank and a branch by their codes, for the names the
// agent fills into the results; the agent leaves unknown ones blank.
func bankNamed(bankCode, branchCode int) (string, string) {
	for _, b := range banks {
		if b.code != bankCode {
			continue
		}
		for _, br := range b.branches {
			if br.code == branchCode {
				return b.name, br.name
			}
		}
		return b.name, ""
	}
	return "", ""
}

// agent is a sales office: its code, name, branch and the month it joined,
// as the agent master has them.
type agent struct {
	code, name, branch string
	joined             int // yymm, as the agent master has it
	grade              int
}

// agents are the sales offices, by the branch they belong to.
var agents = []agent{
	{"A00101", sjis("東京第一営業所"), "0101", 6204, 3}, {"A00102", sjis("東京第二営業所"), "0101", 310, 2},
	{"A00103", sjis("東京第三営業所"), "0101", 1504, 1}, {"A00201", sjis("大阪営業所"), "0201", 1504, 2},
	{"A00202", sjis("大阪南営業所"), "0201", 904, 2}, {"A00301", sjis("名古屋営業所"), "0301", 9011, 1},
	{"A00401", sjis("福岡営業所"), "0401", 2206, 1}, {"A00501", sjis("仙台営業所"), "0501", 710, 2},
	{"A00601", sjis("札幌営業所"), "0601", 1910, 1}, {"A00701", sjis("広島営業所"), "0701", 8604, 3},
}

// closedAgent is an office closed long ago, whose policies the master
// still names and the agent master no longer lists.
const closedAgent = "A00199"

// branchOf is the branch an agent belongs to, as a benefit request names
// the branch that took it.
func branchOf(code string) string {
	for _, a := range agents {
		if a.code == code {
			return a.branch
		}
	}
	return "0101"
}

// policyNumber is a policy number the sales system gives an application:
// the year, the day of the year and the application's number that day.
func policyNumber(date Date, n int) string {
	return fmt.Sprintf("%04d%03d%03d", date.Year(), date.t.YearDay(), n)
}

// requestNumber is a branch's number for a request: a letter for its
// kind, the year, the day of the year and the request's number that day.
func requestNumber(kind byte, date Date, n int) string {
	return fmt.Sprintf("%c%02d%03d%02d", kind, date.Year()%100, date.t.YearDay(), n)
}
