package runner

import (
	"reflect"
	"strings"
	"testing"
)

// The job log and the allocation messages a run writes read back as the
// run's results and allocations.
func TestReadLogs(t *testing.T) {
	r := newRunner(t, map[string]string{
		"COPY": `cat "$DD_IN" > "$DD_OUT"; exit 4`,
	})
	for dsn, text := range map[string]string{"A": "AB\n", "B": "CD\nEF\n"} {
		if err := r.Datasets.Load(dsn, strings.NewReader(text)); err != nil {
			t.Fatal(err)
		}
	}
	job := parse(t, `//J JOB
//S1 EXEC PGM=COPY
//IN DD DSN=A,DISP=SHR
// DD DSN=B,DISP=SHR
//OUT DD DSN=C,DISP=NEW,DCB=(RECFM=FB,LRECL=2)
//RPT DD SYSOUT=*
//NONE DD DUMMY
//S2 EXEC PGM=COPY,COND=(4,LE)
//IN DD DSN=C,DISP=SHR
//OUT DD DSN=D,DISP=NEW,DCB=(RECFM=FB,LRECL=2)
`)
	results, err := r.Run(job)
	if err != nil {
		t.Fatal(err)
	}
	name, logged, err := ReadJobLog(r.Sysout)
	if err != nil {
		t.Fatal(err)
	}
	if name != "J" || !reflect.DeepEqual(logged, results) {
		t.Errorf("job log: %s %+v, want J %+v", name, logged, results)
	}
	allocations, err := ReadMessages(r.Sysout)
	if err != nil {
		t.Fatal(err)
	}
	want := []Allocation{
		{Step: "S1", DD: "IN", Dataset: "A", Records: 1, Cataloged: true},
		{Step: "S1", DD: "IN", Dataset: "B", Records: 2, Cataloged: true},
		{Step: "S1", DD: "OUT", Dataset: "C", Records: 3, Cataloged: true},
		{Step: "S1", DD: "RPT", Sysout: true},
		{Step: "S1", DD: "NONE", Dummy: true},
	}
	if !reflect.DeepEqual(allocations, want) {
		t.Errorf("allocations\n%+v\nwant\n%+v", allocations, want)
	}
}

func TestParseAllocationErrors(t *testing.T) {
	for _, c := range []struct{ line, step, previous string }{
		{"  IN       A                            1 RECORDS", "", ""},
		{"           A                            1 RECORDS", "S1", ""},
		{"  IN       A                            many RECORDS", "S1", ""},
		{"  IN       A                            1", "S1", ""},
	} {
		if _, err := parseAllocation(c.line, c.step, c.previous); err == nil {
			t.Errorf("%q read without error", c.line)
		}
	}
}
