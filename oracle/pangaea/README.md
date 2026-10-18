# pangaea: cases

For each main program, input data that reaches the rules of its decision tables, and of the programs it calls, with what the program wrote and returned: the cases a rewrite has to match. `survey oracle` found them, and replays them to check them.

| Program | Rules reached | Cases |
| --- | --: | --: |
| [ACCTLOAD](ACCTLOAD/) | 6 of 9 | 1 |
| [ACCTMNT](ACCTMNT/) | 33 of 42 | 9 |
| [AUTHSRV](AUTHSRV/) | 49 of 88 | 25 |
| [BANKDATE](BANKDATE/) | 2 of 3 | 2 |
| [CARDLOAD](CARDLOAD/) | 7 of 10 | 2 |
| [GLSUM](GLSUM/) | 16 of 27 | 4 |
| [HOLDEXT](HOLDEXT/) | 8 of 15 | 3 |
| [INTACCR](INTACCR/) | 23 of 33 | 6 |
| [MONCHG](MONCHG/) | 16 of 25 | 9 |
| [STMTPRT](STMTPRT/) | 27 of 34 | 9 |
| [TXPOST](TXPOST/) | 26 of 30 | 10 |
| [TXSORT](TXSORT/) | 1 of 2 | 1 |
| [TXVALID](TXVALID/) | 31 of 66 | 19 |

And for each job, a batch: the job run through its steps on inputs from its programs' cleanest cases, with the datasets it wrote and its output.

| Job | Steps |
| --- | --- |
| [START](jobs/START/) | LOAD RC=0000, CARDS RC=0000 |
| [ONLINE](jobs/ONLINE/) | AUTH RC=0000, EXTRACT RC=0000 |
| [DAILY](jobs/DAILY/) | VALID RC=0000, VALIDC RC=0000, ACCRUE RC=0000, CHARGE RC=0000, POSTMO RC=0000, SORT RC=0000, POST RC=0000, SORTCR RC=0000, POSTCR RC=0000, SORTRT RC=0000, POSTRT RC=0000, MAINT RC=0000, LOAD RC=0000, NEXTDAY RC=0000, CARDS RC=0000 |
| [MONTHLY](jobs/MONTHLY/) | STMT RC=0000, LEDGER RC=0000 |
