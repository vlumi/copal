# laurasia: cases

For each main program, input data that reaches the rules of its decision tables, and of the programs it calls, with what the program wrote and returned: the cases a rewrite has to match. `survey oracle` found them, and replays them to check them.

| Program | Rules reached | Cases |
| --- | --: | --: |
| [ACCTLD](ACCTLD/) | 4 of 5 | 1 |
| [ACCTUL](ACCTUL/) | 1 of 2 | 1 |
| [CLRIN](CLRIN/) | 20 of 24 | 11 |
| [CYCLE](CYCLE/) | 13 of 17 | 7 |
| [DISPUTE](DISPUTE/) | 21 of 22 | 4 |
| [FEELD](FEELD/) | 22 of 23 | 16 |
| [IXFEE](IXFEE/) | 22 of 29 | 14 |
| [MATCH](MATCH/) | 15 of 15 | 4 |
| [POST](POST/) | 17 of 21 | 9 |
| [SETTLE](SETTLE/) | 12 of 12 | 2 |

And for each job, a batch: the job run through its steps on inputs from its programs' cleanest cases, with the datasets it wrote and its output.

| Job | Steps |
| --- | --- |
| [open](jobs/open/) | ACCTLD RC=0008 |
| [nightly](jobs/nightly/) | CLRIN RC=0004, PRESSORT RC=0000, AUTHSORT RC=0000, MATCH RC=0000, FEELD RC=0000, IXFEE RC=0000, POST RC=0000, DISPUTES RC=0000, DISPUTES.DSPSORT RC=0000, DISPUTES.DISPUTE RC=0004, DISPUTES.DSPNEW RC=0000, SETLSORT RC=0000, ADJSORT RC=0000, SETTLE RC=0004, ACCTUL RC=0000 |
| [cycle](jobs/cycle/) | CYCLE RC=0000, ACCTUL RC=0000 |
