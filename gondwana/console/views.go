package console

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/vlumi/cobolith/layout"
	"github.com/vlumi/cobolith/record"
	"github.com/vlumi/copal/gondwana/generate"
	"github.com/vlumi/copal/gondwana/insurer"
)

// errNotFound is for a policy or a day the insurer does not have.
var errNotFound = errors.New("not found")

// Policy is a policy of the master, its names in UTF-8, and its record.
type Policy struct {
	generate.Policy
	Record []byte
}

// statuses are a policy's states by their codes.
var statuses = map[int]string{1: "in force", 2: "in grace", 3: "lapsed", 4: "surrendered", 5: "matured", 6: "paid out"}

// StatusName is the policy's state in words.
func (p Policy) StatusName() string { return statuses[p.Status] }

// ModeName is how the policy's premium is paid, in words.
func (p Policy) ModeName() string {
	if p.Mode == 2 {
		return "annual"
	}
	return "monthly"
}

// PaidToMonth is the month the premiums are paid to, as yyyy-mm.
func (p Policy) PaidToMonth() string {
	s := strconv.FormatInt(p.PaidTo, 10)
	if len(s) != 6 {
		return s
	}
	return s[:4] + "-" + s[4:]
}

// Overview is the insurer as it stands: its policies by state and their
// premiums, and the days it has run, newest first.
type Overview struct {
	Policies int
	ByStatus []StatusCount
	Premiums int64 // of the policies in force or in grace, a month's or a year's each
	Days     []insurer.Day
}

// StatusCount is how many policies are in a state.
type StatusCount struct {
	Status string
	N      int
}

func (s *Server) overview(h insurer.History, _ *http.Request) (any, error) {
	policies, err := s.readPolicies()
	if err != nil {
		return nil, err
	}
	o := Overview{Policies: len(policies)}
	counts := map[int]int{}
	for _, p := range policies {
		counts[p.Status]++
		if p.Status <= 2 {
			o.Premiums += p.Premium
		}
	}
	for status := 1; status <= len(statuses); status++ {
		if counts[status] > 0 {
			o.ByStatus = append(o.ByStatus, StatusCount{Status: statuses[status], N: counts[status]})
		}
	}
	for _, d := range slices.Backward(h.Days) {
		o.Days = append(o.Days, d)
	}
	return o, nil
}

// readPolicies reads the policy master as the last unload left it, the
// names in UTF-8.
func (s *Server) readPolicies() ([]Policy, error) {
	content, err := s.insurer.Policies()
	if err != nil {
		return nil, err
	}
	records := generate.Records(s.layouts.Policy, content)
	read, err := (&generate.Generator{Layouts: s.layouts}).ReadPolicies(records)
	if err != nil {
		return nil, err
	}
	policies := make([]Policy, len(read))
	for i, p := range read {
		p.Name, p.Kana, p.Address, p.Beneficiary = utf8(p.Name), utf8(p.Kana), utf8(p.Address), utf8(p.Beneficiary)
		policies[i] = Policy{Policy: p, Record: records[i]}
	}
	return policies, nil
}

// shown is the most policies a search lists.
const shown = 200

// PolicyList is the policies a search found, at most shown of them.
type PolicyList struct {
	Query    string
	Found    int
	Total    int
	Policies []Policy
}

func (s *Server) policies(_ insurer.History, r *http.Request) (any, error) {
	policies, err := s.readPolicies()
	if err != nil {
		return nil, err
	}
	query := strings.TrimSpace(r.FormValue("q"))
	var found []Policy
	for _, p := range policies {
		if query == "" || strings.HasPrefix(p.No, query) || strings.Contains(p.Name, query) || strings.Contains(p.Kana, query) {
			found = append(found, p)
		}
	}
	return PolicyList{Query: query, Found: len(found), Total: len(policies), Policies: found[:min(len(found), shown)]}, nil
}

// PolicyPage is a policy's master record field by field.
type PolicyPage struct {
	Policy Policy
	Fields []FieldValue
}

// FieldValue is a field of a record: its bytes, and what they mean, in
// UTF-8.
type FieldValue struct {
	Field      *layout.Field
	Raw, Value string
}

func (s *Server) policy(_ insurer.History, r *http.Request) (any, error) {
	policies, err := s.readPolicies()
	if err != nil {
		return nil, err
	}
	i := slices.IndexFunc(policies, func(p Policy) bool { return p.No == r.PathValue("no") })
	if i < 0 {
		return nil, fmt.Errorf("policy %s: %w", r.PathValue("no"), errNotFound)
	}
	return PolicyPage{Policy: policies[i], Fields: fields(s.layouts.Policy, policies[i].Record)}, nil
}

// fields lays a record out field by field, numbers as numbers and text in
// UTF-8.
func fields(l *record.Layout, r []byte) []FieldValue {
	var values []FieldValue
	for _, f := range l.Fields() {
		v := FieldValue{Field: f, Raw: utf8(string(r[f.Start-1 : f.Start-1+f.Length]))}
		if n, err := l.Number(r, f.Name); err == nil && f.Category == layout.Numeric {
			v.Value = strconv.FormatInt(n, 10)
		} else if f.Name != "FILLER" {
			v.Value = strings.TrimRight(v.Raw, " ")
		}
		values = append(values, v)
	}
	return values
}

// DayPage is a day the insurer ran: its batch files and their steps, what
// it amounted to, and the files it read and wrote, a small one's text
// shown.
type DayPage struct {
	Day   insurer.Day
	Files []DayFile
}

// DayFile is a file of a day, with its text in UTF-8, or Large when it is
// too large to show.
type DayFile struct {
	insurer.File
	Text  string
	Large bool
}

// largeFile is the size from which a file is linked to, not shown.
const largeFile = 64 << 10

func (s *Server) day(h insurer.History, r *http.Request) (any, error) {
	i := slices.IndexFunc(h.Days, func(d insurer.Day) bool { return d.Date.String() == r.PathValue("date") })
	if i < 0 {
		return nil, fmt.Errorf("day %s: %w", r.PathValue("date"), errNotFound)
	}
	page := DayPage{Day: h.Days[i]}
	files, err := s.insurer.Files(page.Day.Date)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		df := DayFile{File: f, Large: f.Size >= largeFile}
		if !df.Large {
			content, err := s.insurer.ReadFile(f.Path)
			if err != nil {
				return nil, err
			}
			df.Text = s.text(f.Path, content)
		}
		page.Files = append(page.Files, df)
	}
	return page, nil
}

// text is a file's content in UTF-8: the banks' fixed-length files, which
// have no line ends, a record a line.
func (s *Server) text(path string, content []byte) string {
	if strings.HasSuffix(path, ".DAT") {
		var lines []string
		for _, r := range generate.Fixed(s.layouts.Transfer, content) {
			lines = append(lines, string(r))
		}
		content = []byte(strings.Join(lines, "\n") + "\n")
	}
	return utf8(string(content))
}
