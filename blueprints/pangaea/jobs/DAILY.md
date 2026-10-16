# DAILY

```text
DAILY BATCH: VALIDATE THE DAY'S TRANSACTIONS, ACCRUE INTEREST
AND POST WHAT A CLOSED MONTH OWES, THEN SORT AND POST THE
DAY'S TRANSACTIONS, THE TRANSFER CREDITS AND RETURNS THEY LEAD
TO, APPLY THE DAY'S ACCOUNT MAINTENANCE, LOAD THE ONLINE
ACCOUNT FILE, MOVE TO THE NEXT BUSINESS DAY, AND LOAD THE
ONLINE CARD FILE FOR IT.

READS   BANK.CONTROL(0)  ACCT.MASTER(0)  TRAN.DAILY(0)
        ACCT.ACCRUAL(0)  ACCT.MAINT(0)    BANK.RATES
        CARD.MASTER(0)
WRITES  BANK.CONTROL(+1) ACCT.MASTER(+1) ACCT.ACCRUAL(+1)
        TRAN.AUDIT(+1)   TRAN.REJECTS(+1) MAINT.REJECTS(+1)
        ACCT.ONLINE(+1)  CARD.ONLINE(+1)  HOLD.ONLINE(+1)

A STEP RUNS ONLY IF NO EARLIER STEP ENDED WITH RC ABOVE 4, SO
A RUN THAT FAILS LEAVES THE BUSINESS DATE WHERE IT WAS.
```

Job [DAILY.jcl](../../../pangaea/jobs/DAILY.jcl): 14 steps.

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
  - [LOAD](#load)
  - [NEXTDAY](#nextday)
  - [CARDS](#cards)
- [Datasets](#datasets)
- [Findings](#findings)

## Network

![The steps of DAILY and the datasets they read and write](DAILY.svg)

## Steps

| Step | Program | Reads | Writes | Bypassed by | Description | Source |
| --- | --- | --- | --- | --- | --- | --- |
| VALID | [TXVALID](../programs/TXVALID.md) | BANK.CONTROL(0), TRAN.DAILY(0) | &&VALID, TRAN.REJECTS(+1) | | VALIDATE THE DAY'S TRANSACTION FILE | [DAILY.jcl:23](../../../pangaea/jobs/DAILY.jcl#L23) |
| ACCRUE | [INTACCR](../programs/INTACCR.md) | BANK.CONTROL(0), BANK.RATES, ACCT.MASTER(0), ACCT.ACCRUAL(0) | ACCT.ACCRUAL(+1) | COND=(4,LT) | ACCRUE INTEREST UP TO TODAY ON THE OPENING BALANCES. WHEN A MONTH HAS CLOSED, POST ITS INTEREST AND FEES. THE CHARGES COME IN MASTER ORDER, SO THEY ARE POSTED WITHOUT A SORT. | [DAILY.jcl:36](../../../pangaea/jobs/DAILY.jcl#L36) |
| CHARGE | [MONCHG](../programs/MONCHG.md) | BANK.CONTROL(0), ACCT.MASTER(0), ACCT.ACCRUAL(+1) | &&MONTHLY | COND=(4,LT) | | [DAILY.jcl:44](../../../pangaea/jobs/DAILY.jcl#L44) |
| POSTMO | [TXPOST](../programs/TXPOST.md) | BANK.CONTROL(0), ACCT.MASTER(0), &&MONTHLY | &&CHARGED, TRAN.AUDIT(+1), TRAN.REJECTS(+1) | COND=(4,LT) | | [DAILY.jcl:51](../../../pangaea/jobs/DAILY.jcl#L51) |
| SORT | [TXSORT](../programs/TXSORT.md) | &&VALID | &&SORTED | COND=(4,LT) | SORT THE DAY'S TRANSACTIONS INTO ACCOUNT ORDER AND POST THEM | [DAILY.jcl:65](../../../pangaea/jobs/DAILY.jcl#L65) |
| POST | [TXPOST](../programs/TXPOST.md) | BANK.CONTROL(0), &&CHARGED, &&SORTED | &&POSTED, TRAN.AUDIT(+1), TRAN.REJECTS(+1), &&CREDITS | COND=(4,LT) | | [DAILY.jcl:69](../../../pangaea/jobs/DAILY.jcl#L69) |
| SORTCR | [TXSORT](../programs/TXSORT.md) | &&CREDITS | &&CREDSORT | COND=(4,LT) | SORT AND POST THE CREDITS OF THE TRANSFERS THAT POSTED | [DAILY.jcl:83](../../../pangaea/jobs/DAILY.jcl#L83) |
| POSTCR | [TXPOST](../programs/TXPOST.md) | BANK.CONTROL(0), &&POSTED, &&CREDSORT | &&CREDITED, TRAN.AUDIT(+1), TRAN.REJECTS(+1), &&RETURNS | COND=(4,LT) | | [DAILY.jcl:87](../../../pangaea/jobs/DAILY.jcl#L87) |
| SORTRT | [TXSORT](../programs/TXSORT.md) | &&RETURNS | &&RETSORT | COND=(4,LT) | SORT AND POST THE RETURNS OF THE CREDITS THAT COULD NOT POST | [DAILY.jcl:101](../../../pangaea/jobs/DAILY.jcl#L101) |
| POSTRT | [TXPOST](../programs/TXPOST.md) | BANK.CONTROL(0), &&CREDITED, &&RETSORT | &&RETURNED, TRAN.AUDIT(+1), TRAN.REJECTS(+1) | COND=(4,LT) | | [DAILY.jcl:105](../../../pangaea/jobs/DAILY.jcl#L105) |
| MAINT | [ACCTMNT](../programs/ACCTMNT.md) | BANK.CONTROL(0), &&RETURNED, ACCT.MAINT(0) | ACCT.MASTER(+1), MAINT.REJECTS(+1) | COND=(4,LT) | APPLY THE DAY'S ACCOUNT MAINTENANCE, WRITING THE NEW MASTER | [DAILY.jcl:118](../../../pangaea/jobs/DAILY.jcl#L118) |
| LOAD | [ACCTLOAD](../programs/ACCTLOAD.md) | BANK.CONTROL(0), ACCT.MASTER(+1) | ACCT.ONLINE(+1), HOLD.ONLINE(+1) | COND=(4,LT) | LOAD THE NEW MASTER INTO THE ONLINE ACCOUNT FILE, WHICH THE CARD AND ATM AUTHORIZATIONS READ DURING THE NEXT BUSINESS DAY, AND START THAT DAY'S HOLD FILE | [DAILY.jcl:132](../../../pangaea/jobs/DAILY.jcl#L132) |
| NEXTDAY | [BANKDATE](../programs/BANKDATE.md) | BANK.CONTROL(0) | BANK.CONTROL(+1) | COND=(4,LT) | MOVE TO THE NEXT BUSINESS DAY | [DAILY.jcl:143](../../../pangaea/jobs/DAILY.jcl#L143) |
| CARDS | [CARDLOAD](../programs/CARDLOAD.md) | BANK.CONTROL(+1), CARD.MASTER(0) | CARD.ONLINE(+1) | COND=(4,LT) | LOAD THE CARDS INTO THE ONLINE CARD FILE FOR THE NEXT BUSINESS DAY, LEAVING OUT THOSE THAT HAVE EXPIRED BY THEN | [DAILY.jcl:151](../../../pangaea/jobs/DAILY.jcl#L151) |

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

### LOAD

| DD | Dataset | DISP | File | Copybook | Use |
| --- | --- | --- | --- | --- | --- |
| BANKCTL | BANK.CONTROL(0) | SHR | CONTROL-FILE | [BANKCTL](../copybooks/BANKCTL.md) | read |
| ACCTIN | ACCT.MASTER(+1) | SHR | ACCT-IN | [ACCTMAST](../copybooks/ACCTMAST.md) | read |
| ACCTOUT | ACCT.ONLINE(+1) (indexed) | NEW | ONLINE-OUT | [ACCTONL](../copybooks/ACCTONL.md) | write |
| HOLDOUT | HOLD.ONLINE(+1) (indexed) | NEW | HOLD-OUT | [HOLDREC](../copybooks/HOLDREC.md) | write |
| LOADRPT | job output | | REPORT-OUT | | write |

### NEXTDAY

| DD | Dataset | DISP | File | Copybook | Use |
| --- | --- | --- | --- | --- | --- |
| BANKIN | BANK.CONTROL(0) | SHR | CONTROL-IN | [BANKCTL](../copybooks/BANKCTL.md) | read |
| BANKOUT | BANK.CONTROL(+1) | NEW | CONTROL-OUT | [BANKCTL](../copybooks/BANKCTL.md) | write |

### CARDS

| DD | Dataset | DISP | File | Copybook | Use |
| --- | --- | --- | --- | --- | --- |
| BANKCTL | BANK.CONTROL(+1) | SHR | CONTROL-FILE | [BANKCTL](../copybooks/BANKCTL.md) | read |
| CARDIN | CARD.MASTER(0) | SHR | CARD-IN | [CARDREC](../copybooks/CARDREC.md) | read |
| CARDOUT | CARD.ONLINE(+1) (indexed) | NEW | ONLINE-OUT | [CARDREC](../copybooks/CARDREC.md) | write |
| LOADRPT | job output | | REPORT-OUT | | write |

## Datasets

Each dataset the job names, the copybook that describes its records, from the programs that read or write it in any job, and the steps that write and read it.

| Dataset | Copybook | Written by | Read by |
| --- | --- | --- | --- |
| BANK.CONTROL(0) | [BANKCTL](../copybooks/BANKCTL.md) | | VALID, ACCRUE, CHARGE, POSTMO, POST, POSTCR, POSTRT, MAINT, LOAD, NEXTDAY |
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
| ACCT.MASTER(+1) | [ACCTMAST](../copybooks/ACCTMAST.md) | MAINT | LOAD |
| MAINT.REJECTS(+1) | [REJREC](../copybooks/REJREC.md) | MAINT | |
| ACCT.ONLINE(+1) | [ACCTONL](../copybooks/ACCTONL.md) | LOAD | |
| HOLD.ONLINE(+1) | [HOLDREC](../copybooks/HOLDREC.md) | LOAD | |
| BANK.CONTROL(+1) | [BANKCTL](../copybooks/BANKCTL.md) | NEXTDAY | CARDS |
| CARD.MASTER(0) | [CARDREC](../copybooks/CARDREC.md) | | CARDS |
| CARD.ONLINE(+1) | [CARDREC](../copybooks/CARDREC.md) | CARDS | |

## Findings

None: every DD is assigned a file, every file a DD, and every disposition fits its use.
