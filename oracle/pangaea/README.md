# pangaea: cases

For each main program, input data that reaches the rules of its decision tables, and of the programs it calls, with what the program wrote and returned: the cases a rewrite has to match. `survey oracle` found them, and replays them to check them.

| Program | Rules reached | Cases |
| --- | --: | --: |
| [ACCTLOAD](ACCTLOAD/) | 9 of 9 | 4 |
| [ACCTMNT](ACCTMNT/) | 37 of 41 | 13 |
| [AUTHSRV](AUTHSRV/) | 60 of 65 | 32 |
| [BANKDATE](BANKDATE/) | 3 of 3 | 3 |
| [CARDLOAD](CARDLOAD/) | 10 of 10 | 5 |
| [GLSUM](GLSUM/) | 20 of 25 | 7 |
| [HOLDEXT](HOLDEXT/) | 11 of 11 | 6 |
| [INTACCR](INTACCR/) | 27 of 33 | 9 |
| [MONCHG](MONCHG/) | 18 of 21 | 11 |
| [STMTPRT](STMTPRT/) | 31 of 34 | 12 |
| [TXPOST](TXPOST/) | 30 of 30 | 13 |
| [TXSORT](TXSORT/) | 2 of 2 | 2 |
| [TXVALID](TXVALID/) | 32 of 33 | 20 |

And for each job, a batch: the job run through its steps on inputs from its programs' cleanest cases, with the datasets it wrote and its output.

| Job | Steps |
| --- | --- |
| [START](jobs/START/) | LOAD RC=0000, CARDS RC=0000 |
| [ONLINE](jobs/ONLINE/) | AUTH RC=0000, EXTRACT RC=0000 |
| [DAILY](jobs/DAILY/) | VALID RC=0000, VALIDC RC=0000, ACCRUE RC=0000, CHARGE RC=0000, POSTMO RC=0000, SORT RC=0000, POST RC=0000, SORTCR RC=0000, POSTCR RC=0000, SORTRT RC=0000, POSTRT RC=0000, MAINT RC=0000, LOAD RC=0000, NEXTDAY RC=0000, CARDS RC=0000 |
| [MONTHLY](jobs/MONTHLY/) | STMT RC=0000, LEDGER RC=0000 |
