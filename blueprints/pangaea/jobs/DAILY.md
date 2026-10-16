# DAILY

```text
DAILY BATCH: VALIDATE THE DAY'S TRANSACTIONS, ACCRUE INTEREST
AND POST WHAT A CLOSED MONTH OWES, THEN SORT AND POST THE
DAY'S TRANSACTIONS, THE TRANSFER CREDITS AND RETURNS THEY LEAD
TO, APPLY THE DAY'S ACCOUNT MAINTENANCE, AND MOVE TO THE NEXT
BUSINESS DAY.

READS   BANK.CONTROL(0)  ACCT.MASTER(0)  TRAN.DAILY(0)
        ACCT.ACCRUAL(0)  ACCT.MAINT(0)    BANK.RATES
WRITES  BANK.CONTROL(+1) ACCT.MASTER(+1) ACCT.ACCRUAL(+1)
        TRAN.AUDIT(+1)   TRAN.REJECTS(+1) MAINT.REJECTS(+1)

A STEP RUNS ONLY IF NO EARLIER STEP ENDED WITH RC ABOVE 4, SO
A RUN THAT FAILS LEAVES THE BUSINESS DATE WHERE IT WAS.
```

Job [DAILY.jcl](../../../pangaea/jobs/DAILY.jcl): 12 steps.

## Contents

- [Network](#network)
- [Steps](#steps)
  - [VALID](#valid)
  - [ACCRUE](#accrue)
  - [CHARGE](#charge)
  - [POSTMO](#postmo)
  - [SORT](#sort)
  - [POST](#post)
  - [SORTCR](#sortcr)
  - [POSTCR](#postcr)
  - [SORTRT](#sortrt)
  - [POSTRT](#postrt)
  - [MAINT](#maint)
  - [NEXTDAY](#nextday)
- [Datasets](#datasets)
- [Findings](#findings)

## Network

![The steps of DAILY and the datasets they read and write](DAILY.svg)

## Steps

| Step | Program | Reads | Writes | Bypassed by | Description | Source |
| --- | --- | --- | --- | --- | --- | --- |
| VALID | [TXVALID](../programs/TXVALID.md) | BANK.CONTROL(0), TRAN.DAILY(0) | &&VALID, TRAN.REJECTS(+1) | | VALIDATE THE DAY'S TRANSACTION FILE | [DAILY.jcl:20](../../../pangaea/jobs/DAILY.jcl#L20) |
| ACCRUE | [INTACCR](../programs/INTACCR.md) | BANK.CONTROL(0), BANK.RATES, ACCT.MASTER(0), ACCT.ACCRUAL(0) | ACCT.ACCRUAL(+1) | COND=(4,LT) | ACCRUE INTEREST UP TO TODAY ON THE OPENING BALANCES. WHEN A MONTH HAS CLOSED, POST ITS INTEREST AND FEES. THE CHARGES COME IN MASTER ORDER, SO THEY ARE POSTED WITHOUT A SORT. | [DAILY.jcl:33](../../../pangaea/jobs/DAILY.jcl#L33) |
| CHARGE | [MONCHG](../programs/MONCHG.md) | BANK.CONTROL(0), ACCT.MASTER(0), ACCT.ACCRUAL(+1) | &&MONTHLY | COND=(4,LT) | | [DAILY.jcl:41](../../../pangaea/jobs/DAILY.jcl#L41) |
| POSTMO | [TXPOST](../programs/TXPOST.md) | BANK.CONTROL(0), ACCT.MASTER(0), &&MONTHLY | &&CHARGED, TRAN.AUDIT(+1), TRAN.REJECTS(+1) | COND=(4,LT) | | [DAILY.jcl:48](../../../pangaea/jobs/DAILY.jcl#L48) |
| SORT | [TXSORT](../programs/TXSORT.md) | &&VALID | &&SORTED | COND=(4,LT) | SORT THE DAY'S TRANSACTIONS INTO ACCOUNT ORDER AND POST THEM | [DAILY.jcl:62](../../../pangaea/jobs/DAILY.jcl#L62) |
| POST | [TXPOST](../programs/TXPOST.md) | BANK.CONTROL(0), &&CHARGED, &&SORTED | &&POSTED, TRAN.AUDIT(+1), TRAN.REJECTS(+1), &&CREDITS | COND=(4,LT) | | [DAILY.jcl:66](../../../pangaea/jobs/DAILY.jcl#L66) |
| SORTCR | [TXSORT](../programs/TXSORT.md) | &&CREDITS | &&CREDSORT | COND=(4,LT) | SORT AND POST THE CREDITS OF THE TRANSFERS THAT POSTED | [DAILY.jcl:80](../../../pangaea/jobs/DAILY.jcl#L80) |
| POSTCR | [TXPOST](../programs/TXPOST.md) | BANK.CONTROL(0), &&POSTED, &&CREDSORT | &&CREDITED, TRAN.AUDIT(+1), TRAN.REJECTS(+1), &&RETURNS | COND=(4,LT) | | [DAILY.jcl:84](../../../pangaea/jobs/DAILY.jcl#L84) |
| SORTRT | [TXSORT](../programs/TXSORT.md) | &&RETURNS | &&RETSORT | COND=(4,LT) | SORT AND POST THE RETURNS OF THE CREDITS THAT COULD NOT POST | [DAILY.jcl:98](../../../pangaea/jobs/DAILY.jcl#L98) |
| POSTRT | [TXPOST](../programs/TXPOST.md) | BANK.CONTROL(0), &&CREDITED, &&RETSORT | &&RETURNED, TRAN.AUDIT(+1), TRAN.REJECTS(+1) | COND=(4,LT) | | [DAILY.jcl:102](../../../pangaea/jobs/DAILY.jcl#L102) |
| MAINT | [ACCTMNT](../programs/ACCTMNT.md) | BANK.CONTROL(0), &&RETURNED, ACCT.MAINT(0) | ACCT.MASTER(+1), MAINT.REJECTS(+1) | COND=(4,LT) | APPLY THE DAY'S ACCOUNT MAINTENANCE, WRITING THE NEW MASTER | [DAILY.jcl:115](../../../pangaea/jobs/DAILY.jcl#L115) |
| NEXTDAY | [BANKDATE](../programs/BANKDATE.md) | BANK.CONTROL(0) | BANK.CONTROL(+1) | COND=(4,LT) | MOVE TO THE NEXT BUSINESS DAY | [DAILY.jcl:127](../../../pangaea/jobs/DAILY.jcl#L127) |

### VALID

| DD | Dataset | DISP | File | Copybook | Use |
| --- | --- | --- | --- | --- | --- |
| BANKCTL | BANK.CONTROL(0) | SHR | CONTROL-FILE | [BANKCTL](../copybooks/BANKCTL.md) | read |
| TRANIN | TRAN.DAILY(0) | SHR | TRAN-IN | [TRANREC](../copybooks/TRANREC.md) | read |
| TRANOUT | &&VALID | NEW | TRAN-OUT | [TRANREC](../copybooks/TRANREC.md) | write |
| REJECTS | TRAN.REJECTS(+1) | NEW | REJECT-OUT | [REJREC](../copybooks/REJREC.md) | write |
| VALRPT | job output | | REPORT-OUT | | write |

### ACCRUE

| DD | Dataset | DISP | File | Copybook | Use |
| --- | --- | --- | --- | --- | --- |
| BANKCTL | BANK.CONTROL(0) | SHR | CONTROL-FILE | [BANKCTL](../copybooks/BANKCTL.md) | read |
| RATES | BANK.RATES | SHR | RATE-FILE | [RATEREC](../copybooks/RATEREC.md) | read |
| ACCTIN | ACCT.MASTER(0) | SHR | ACCT-IN | [ACCTMAST](../copybooks/ACCTMAST.md) | read |
| ACCRIN | ACCT.ACCRUAL(0) | SHR | ACCR-IN | | read |
| ACCROUT | ACCT.ACCRUAL(+1) | NEW | ACCR-OUT | [ACCRREC](../copybooks/ACCRREC.md) | write |
| ACCRRPT | job output | | REPORT-OUT | | write |

### CHARGE

| DD | Dataset | DISP | File | Copybook | Use |
| --- | --- | --- | --- | --- | --- |
| BANKCTL | BANK.CONTROL(0) | SHR | CONTROL-FILE | [BANKCTL](../copybooks/BANKCTL.md) | read |
| ACCTIN | ACCT.MASTER(0) | SHR | ACCT-IN | [ACCTMAST](../copybooks/ACCTMAST.md) | read |
| ACCRIN | ACCT.ACCRUAL(+1) | SHR | ACCR-IN | [ACCRREC](../copybooks/ACCRREC.md) | read |
| TRANOUT | &&MONTHLY | NEW | TRAN-OUT | [TRANREC](../copybooks/TRANREC.md) | write |
| MONRPT | job output | | REPORT-OUT | | write |

### POSTMO

| DD | Dataset | DISP | File | Copybook | Use |
| --- | --- | --- | --- | --- | --- |
| BANKCTL | BANK.CONTROL(0) | SHR | CONTROL-FILE | [BANKCTL](../copybooks/BANKCTL.md) | read |
| ACCTIN | ACCT.MASTER(0) | SHR | ACCT-IN | [ACCTMAST](../copybooks/ACCTMAST.md) | read |
| TRANIN | &&MONTHLY | OLD | TRAN-IN | [TRANREC](../copybooks/TRANREC.md) | read |
| ACCTOUT | &&CHARGED | NEW | ACCT-OUT | [ACCTMAST](../copybooks/ACCTMAST.md) | write |
| AUDIT | TRAN.AUDIT(+1) | MOD | AUDIT-OUT | [AUDITREC](../copybooks/AUDITREC.md) | append |
| REJECTS | TRAN.REJECTS(+1) | MOD | REJECT-OUT | [REJREC](../copybooks/REJREC.md) | append |
| LEGS | dummy | | LEGS-OUT | [TRANREC](../copybooks/TRANREC.md) | write |
| POSTRPT | job output | | REPORT-OUT | | write |

### SORT

| DD | Dataset | DISP | File | Copybook | Use |
| --- | --- | --- | --- | --- | --- |
| TRANIN | &&VALID | OLD | TRAN-IN | | read |
| TRANOUT | &&SORTED | NEW | TRAN-OUT | | write |

### POST

| DD | Dataset | DISP | File | Copybook | Use |
| --- | --- | --- | --- | --- | --- |
| BANKCTL | BANK.CONTROL(0) | SHR | CONTROL-FILE | [BANKCTL](../copybooks/BANKCTL.md) | read |
| ACCTIN | &&CHARGED | OLD | ACCT-IN | [ACCTMAST](../copybooks/ACCTMAST.md) | read |
| TRANIN | &&SORTED | OLD | TRAN-IN | [TRANREC](../copybooks/TRANREC.md) | read |
| ACCTOUT | &&POSTED | NEW | ACCT-OUT | [ACCTMAST](../copybooks/ACCTMAST.md) | write |
| AUDIT | TRAN.AUDIT(+1) | MOD | AUDIT-OUT | [AUDITREC](../copybooks/AUDITREC.md) | append |
| REJECTS | TRAN.REJECTS(+1) | MOD | REJECT-OUT | [REJREC](../copybooks/REJREC.md) | append |
| LEGS | &&CREDITS | NEW | LEGS-OUT | [TRANREC](../copybooks/TRANREC.md) | write |
| POSTRPT | job output | | REPORT-OUT | | write |

### SORTCR

| DD | Dataset | DISP | File | Copybook | Use |
| --- | --- | --- | --- | --- | --- |
| TRANIN | &&CREDITS | OLD | TRAN-IN | | read |
| TRANOUT | &&CREDSORT | NEW | TRAN-OUT | | write |

### POSTCR

| DD | Dataset | DISP | File | Copybook | Use |
| --- | --- | --- | --- | --- | --- |
| BANKCTL | BANK.CONTROL(0) | SHR | CONTROL-FILE | [BANKCTL](../copybooks/BANKCTL.md) | read |
| ACCTIN | &&POSTED | OLD | ACCT-IN | [ACCTMAST](../copybooks/ACCTMAST.md) | read |
| TRANIN | &&CREDSORT | OLD | TRAN-IN | [TRANREC](../copybooks/TRANREC.md) | read |
| ACCTOUT | &&CREDITED | NEW | ACCT-OUT | [ACCTMAST](../copybooks/ACCTMAST.md) | write |
| AUDIT | TRAN.AUDIT(+1) | MOD | AUDIT-OUT | [AUDITREC](../copybooks/AUDITREC.md) | append |
| REJECTS | TRAN.REJECTS(+1) | MOD | REJECT-OUT | [REJREC](../copybooks/REJREC.md) | append |
| LEGS | &&RETURNS | NEW | LEGS-OUT | [TRANREC](../copybooks/TRANREC.md) | write |
| POSTRPT | job output | | REPORT-OUT | | write |

### SORTRT

| DD | Dataset | DISP | File | Copybook | Use |
| --- | --- | --- | --- | --- | --- |
| TRANIN | &&RETURNS | OLD | TRAN-IN | | read |
| TRANOUT | &&RETSORT | NEW | TRAN-OUT | | write |

### POSTRT

| DD | Dataset | DISP | File | Copybook | Use |
| --- | --- | --- | --- | --- | --- |
| BANKCTL | BANK.CONTROL(0) | SHR | CONTROL-FILE | [BANKCTL](../copybooks/BANKCTL.md) | read |
| ACCTIN | &&CREDITED | OLD | ACCT-IN | [ACCTMAST](../copybooks/ACCTMAST.md) | read |
| TRANIN | &&RETSORT | OLD | TRAN-IN | [TRANREC](../copybooks/TRANREC.md) | read |
| ACCTOUT | &&RETURNED | NEW | ACCT-OUT | [ACCTMAST](../copybooks/ACCTMAST.md) | write |
| AUDIT | TRAN.AUDIT(+1) | MOD | AUDIT-OUT | [AUDITREC](../copybooks/AUDITREC.md) | append |
| REJECTS | TRAN.REJECTS(+1) | MOD | REJECT-OUT | [REJREC](../copybooks/REJREC.md) | append |
| LEGS | dummy | | LEGS-OUT | [TRANREC](../copybooks/TRANREC.md) | write |
| POSTRPT | job output | | REPORT-OUT | | write |

### MAINT

| DD | Dataset | DISP | File | Copybook | Use |
| --- | --- | --- | --- | --- | --- |
| BANKCTL | BANK.CONTROL(0) | SHR | CONTROL-FILE | [BANKCTL](../copybooks/BANKCTL.md) | read |
| ACCTIN | &&RETURNED | OLD | ACCT-IN | | read |
| MAINTIN | ACCT.MAINT(0) | SHR | MAINT-IN | [MAINTREC](../copybooks/MAINTREC.md) | read |
| ACCTOUT | ACCT.MASTER(+1) | NEW | ACCT-OUT | [ACCTMAST](../copybooks/ACCTMAST.md) | write |
| REJECTS | MAINT.REJECTS(+1) | NEW | REJECT-OUT | [REJREC](../copybooks/REJREC.md) | write |
| MNTRPT | job output | | REPORT-OUT | | write |

### NEXTDAY

| DD | Dataset | DISP | File | Copybook | Use |
| --- | --- | --- | --- | --- | --- |
| BANKIN | BANK.CONTROL(0) | SHR | CONTROL-IN | [BANKCTL](../copybooks/BANKCTL.md) | read |
| BANKOUT | BANK.CONTROL(+1) | NEW | CONTROL-OUT | [BANKCTL](../copybooks/BANKCTL.md) | write |

## Datasets

Each dataset the job names, the copybook that describes its records, from the programs that read or write it in any job, and the steps that write and read it.

| Dataset | Copybook | Written by | Read by |
| --- | --- | --- | --- |
| BANK.CONTROL(0) | [BANKCTL](../copybooks/BANKCTL.md) | | VALID, ACCRUE, CHARGE, POSTMO, POST, POSTCR, POSTRT, MAINT, NEXTDAY |
| TRAN.DAILY(0) | [TRANREC](../copybooks/TRANREC.md) | | VALID |
| &&VALID | [TRANREC](../copybooks/TRANREC.md) | VALID | SORT |
| TRAN.REJECTS(+1) | [REJREC](../copybooks/REJREC.md) | VALID, POSTMO, POST, POSTCR, POSTRT | |
| BANK.RATES | [RATEREC](../copybooks/RATEREC.md) | | ACCRUE |
| ACCT.MASTER(0) | [ACCTMAST](../copybooks/ACCTMAST.md) | | ACCRUE, CHARGE, POSTMO |
| ACCT.ACCRUAL(0) | [ACCRREC](../copybooks/ACCRREC.md) | | ACCRUE |
| ACCT.ACCRUAL(+1) | [ACCRREC](../copybooks/ACCRREC.md) | ACCRUE | CHARGE |
| &&MONTHLY | [TRANREC](../copybooks/TRANREC.md) | CHARGE | POSTMO |
| &&CHARGED | [ACCTMAST](../copybooks/ACCTMAST.md) | POSTMO | POST |
| TRAN.AUDIT(+1) | [AUDITREC](../copybooks/AUDITREC.md) | POSTMO, POST, POSTCR, POSTRT | |
| &&SORTED | [TRANREC](../copybooks/TRANREC.md) | SORT | POST |
| &&POSTED | [ACCTMAST](../copybooks/ACCTMAST.md) | POST | POSTCR |
| &&CREDITS | [TRANREC](../copybooks/TRANREC.md) | POST | SORTCR |
| &&CREDSORT | [TRANREC](../copybooks/TRANREC.md) | SORTCR | POSTCR |
| &&CREDITED | [ACCTMAST](../copybooks/ACCTMAST.md) | POSTCR | POSTRT |
| &&RETURNS | [TRANREC](../copybooks/TRANREC.md) | POSTCR | SORTRT |
| &&RETSORT | [TRANREC](../copybooks/TRANREC.md) | SORTRT | POSTRT |
| &&RETURNED | [ACCTMAST](../copybooks/ACCTMAST.md) | POSTRT | MAINT |
| ACCT.MAINT(0) | [MAINTREC](../copybooks/MAINTREC.md) | | MAINT |
| ACCT.MASTER(+1) | [ACCTMAST](../copybooks/ACCTMAST.md) | MAINT | |
| MAINT.REJECTS(+1) | [REJREC](../copybooks/REJREC.md) | MAINT | |
| BANK.CONTROL(+1) | [BANKCTL](../copybooks/BANKCTL.md) | NEXTDAY | |

## Findings

None: every DD is assigned a file, every file a DD, and every disposition fits its use.
