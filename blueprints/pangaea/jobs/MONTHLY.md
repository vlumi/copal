# MONTHLY

```text
MONTH-END: THE STATEMENTS AND THE GENERAL LEDGER SUMMARY OF
THE MONTH THAT CLOSED. THE SCHEDULER RUNS IT ON THE FIRST
BUSINESS DAY OF EVERY MONTH, AFTER THE DAILY BATCH HAS POSTED
THE MONTH'S INTEREST AND FEES. ON THE FIRST DAY OF A NEW BANK
NO MONTH HAS CLOSED, AND IT PRINTS NOTHING.

READS   ACCT.ACCRUAL(0)  ACCT.MASTER(-1)  TRAN.AUDIT
WRITES  GL.SUMMARY(+1)
```

Job [MONTHLY.jcl](../../../pangaea/jobs/MONTHLY.jcl): 2 steps.

## Contents

- [Network](#network)
- [Steps](#steps)
  - [STMT](#stmt)
  - [LEDGER](#ledger)
- [Datasets](#datasets)
- [Findings](#findings)

## Network

![The steps of MONTHLY and the datasets they read and write](MONTHLY.svg)

## Steps

| Step | Program | Reads | Writes | Bypassed by | Description | Source |
| --- | --- | --- | --- | --- | --- | --- |
| STMT | [STMTPRT](../programs/STMTPRT.md) | ACCT.ACCRUAL(0), ACCT.MASTER(-1), TRAN.AUDIT | | | PRINT THE STATEMENTS FROM EVERY AUDIT TRAIL AND THE MASTER THE CLOSED MONTH ENDED WITH, BEFORE THE DAILY BATCH OF TODAY | [MONTHLY.jcl:16](../../../pangaea/jobs/MONTHLY.jcl#L16) |
| LEDGER | [GLSUM](../programs/GLSUM.md) | ACCT.ACCRUAL(0), TRAN.AUDIT | GL.SUMMARY(+1) | COND=(4,LT) | SUMMARIZE THE MONTH'S POSTINGS FOR THE GENERAL LEDGER | [MONTHLY.jcl:25](../../../pangaea/jobs/MONTHLY.jcl#L25) |

### STMT

| DD | Dataset | DISP | File | Copybook | Use |
| --- | --- | --- | --- | --- | --- |
| ACCRIN | ACCT.ACCRUAL(0) | SHR | ACCR-IN | [ACCRREC](../copybooks/ACCRREC.md) | read |
| ACCTIN | ACCT.MASTER(-1) | SHR | ACCT-IN | [ACCTMAST](../copybooks/ACCTMAST.md) | read |
| AUDIT | TRAN.AUDIT | SHR | AUDIT-IN | [AUDITREC](../copybooks/AUDITREC.md) | read |
| STMTOUT | job output | | STMT-OUT | | write |
| STMTRPT | job output | | REPORT-OUT | | write |

### LEDGER

| DD | Dataset | DISP | File | Copybook | Use |
| --- | --- | --- | --- | --- | --- |
| ACCRIN | ACCT.ACCRUAL(0) | SHR | ACCR-IN | [ACCRREC](../copybooks/ACCRREC.md) | read |
| AUDIT | TRAN.AUDIT | SHR | AUDIT-IN | [AUDITREC](../copybooks/AUDITREC.md) | read |
| GLOUT | GL.SUMMARY(+1) | NEW | GL-OUT | [GLREC](../copybooks/GLREC.md) | write |
| GLRPT | job output | | REPORT-OUT | | write |

## Datasets

Each dataset the job names, the copybook that describes its records, from the programs that read or write it in any job, and the steps that write and read it.

| Dataset | Copybook | Written by | Read by |
| --- | --- | --- | --- |
| ACCT.ACCRUAL(0) | [ACCRREC](../copybooks/ACCRREC.md) | | STMT, LEDGER |
| ACCT.MASTER(-1) | [ACCTMAST](../copybooks/ACCTMAST.md) | | STMT |
| TRAN.AUDIT | [AUDITREC](../copybooks/AUDITREC.md) | | STMT, LEDGER |
| GL.SUMMARY(+1) | [GLREC](../copybooks/GLREC.md) | LEDGER | |

## Findings

None: every DD is assigned a file, every file a DD, and every disposition fits its use.
