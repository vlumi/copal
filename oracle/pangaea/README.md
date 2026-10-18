# pangaea: cases

For each main program, input data that reaches the rules of its decision tables, and of the programs it calls, with what the program wrote and returned: the cases a rewrite has to match. `survey oracle` found them, and replays them to check them.

| Program | Rules reached | Cases |
| --- | --: | --: |
| [ACCTLOAD](ACCTLOAD/) | 8 of 9 | 3 |
| [ACCTMNT](ACCTMNT/) | 36 of 41 | 12 |
| [AUTHSRV](AUTHSRV/) | 53 of 79 | 26 |
| [BANKDATE](BANKDATE/) | 2 of 3 | 2 |
| [CARDLOAD](CARDLOAD/) | 9 of 10 | 4 |
| [GLSUM](GLSUM/) | 17 of 26 | 5 |
| [HOLDEXT](HOLDEXT/) | 10 of 11 | 5 |
| [INTACCR](INTACCR/) | 25 of 33 | 7 |
| [MONCHG](MONCHG/) | 17 of 21 | 10 |
| [STMTPRT](STMTPRT/) | 28 of 34 | 10 |
| [TXPOST](TXPOST/) | 29 of 30 | 12 |
| [TXSORT](TXSORT/) | 1 of 2 | 1 |
| [TXVALID](TXVALID/) | 31 of 39 | 19 |

And for each job, a batch: the job run through its steps on inputs from its programs' cleanest cases, with the datasets it wrote and its output.

| Job | Steps |
| --- | --- |
| [START](jobs/START/) | LOAD RC=0000, CARDS RC=0000 |
| [ONLINE](jobs/ONLINE/) | AUTH RC=0000, EXTRACT RC=0000 |
| [DAILY](jobs/DAILY/) | VALID RC=0000, VALIDC RC=0000, ACCRUE RC=0000, CHARGE RC=0000, POSTMO RC=0000, SORT RC=0000, POST RC=0000, SORTCR RC=0000, POSTCR RC=0000, SORTRT RC=0000, POSTRT RC=0000, MAINT RC=0000, LOAD RC=0000, NEXTDAY RC=0000, CARDS RC=0000 |
| [MONTHLY](jobs/MONTHLY/) | STMT RC=0000, LEDGER RC=0000 |
