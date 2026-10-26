# nightly: a batch

The script run with sh after the scripts before it in the batch, on what they wrote, and on inputs taken from the cases of its steps' programs for the rest: each file it reads that no script before it, and no step of it before, wrote, with the records of the cleanest case of the first program that reads it, the first that ended with the lowest return code, and each of its parameters with the command line of the cleanest case of the first program it is passed to. The files it wrote and its output are recorded beside.

The script ended with exit status 0.

| Step | Program | Outcome | Rules reached |
| --- | --- | --- | --: |
| CLRIN | [CLRIN](../../CLRIN/) | RC=0004 | 7 of 23 |
| CLRONUS | [CLRIN](../../CLRIN/) | RC=0004 | 7 of 23 |
| PRESSORT | sort | RC=0000 | |
| AUTHSORT | sort | RC=0000 | |
| MATCH | [MATCH](../../MATCH/) | RC=0000 | 5 of 15 |
| FEELD | [FEELD](../../FEELD/) | RC=0000 | 5 of 23 |
| IXFEE | [IXFEE](../../IXFEE/) | RC=0000 | 2 of 29 |
| POST | [POST](../../POST/) | RC=0000 | 2 of 21 |
| PAYMENTS | payments | RC=0000 | |
| PAYMENTS.PAYIN | [PAYIN](../../PAYIN/) | RC=0004 | 5 of 17 |
| DISPUTES | disputes | RC=0000 | |
| DISPUTES.DSPSORT | sort | RC=0000 | |
| DISPUTES.DISPUTE | [DISPUTE](../../DISPUTE/) | RC=0004 | 15 of 21 |
| DISPUTES.DSPNEW | mv | RC=0000 | |
| SETLSORT | sort | RC=0000 | |
| ADJSORT | sort | RC=0000 | |
| SETTLE | [SETTLE](../../SETTLE/) | RC=0000 | 5 of 14 |
| ACCTUL | [ACCTUL](../../ACCTUL/) | RC=0000 | 1 of 2 |
