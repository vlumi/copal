# DAILY

```text
DAILY BATCH: VALIDATE THE DAY'S TRANSACTIONS FROM THE BRANCHES
AND THE CARDS, ACCRUE INTEREST AND POST WHAT A CLOSED MONTH
OWES, THEN SORT AND POST THE DAY'S TRANSACTIONS, THE TRANSFER
CREDITS AND RETURNS THEY LEAD TO, APPLY THE DAY'S ACCOUNT
MAINTENANCE, LOAD THE ONLINE ACCOUNT FILE, MOVE TO THE NEXT
BUSINESS DAY, AND LOAD THE ONLINE CARD FILE FOR IT.

READS   BANK.CONTROL(0)  ACCT.MASTER(0)  TRAN.DAILY(0)
        TRAN.CARD(0)     ACCT.ACCRUAL(0)  ACCT.MAINT(0)
        CARD.MASTER(0)   BANK.RATES
WRITES  BANK.CONTROL(+1) ACCT.MASTER(+1) ACCT.ACCRUAL(+1)
        TRAN.AUDIT(+1)   TRAN.REJECTS(+1) MAINT.REJECTS(+1)
        ACCT.ONLINE(+1)  CARD.ONLINE(+1)  HOLD.ONLINE(+1)

A STEP RUNS ONLY IF NO EARLIER STEP ENDED WITH RC ABOVE 4, SO
A RUN THAT FAILS LEAVES THE BUSINESS DATE WHERE IT WAS.
```

Job [DAILY.jcl](../../../pangaea/jobs/DAILY.jcl): 15 steps.

## Contents

- [Network](#network)
- [Steps](#steps)
  - [VALID](#valid)
  - [VALIDC](#validc)
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
| VALID | [TXVALID](../programs/TXVALID.md) | BANK.CONTROL(0), TRAN.DAILY(0) | &&VALID, TRAN.REJECTS(+1) | | VALIDATE THE BRANCHES' TRANSACTION FILE | [DAILY.jcl:23](../../../pangaea/jobs/DAILY.jcl#L23) |
| VALIDC | [TXVALID](../programs/TXVALID.md) | BANK.CONTROL(0), TRAN.CARD(0) | &&CARDVAL, TRAN.REJECTS(+1) | COND=(4,LT) | VALIDATE THE CARD FILE, THE ONLINE DAY'S AUTHORIZED CARD AND ATM TRANSACTIONS | [DAILY.jcl:35](../../../pangaea/jobs/DAILY.jcl#L35) |
| ACCRUE | [INTACCR](../programs/INTACCR.md) | BANK.CONTROL(0), BANK.RATES, ACCT.MASTER(0), ACCT.ACCRUAL(0) | ACCT.ACCRUAL(+1) | COND=(4,LT) | ACCRUE INTEREST UP TO TODAY ON THE OPENING BALANCES. WHEN A MONTH HAS CLOSED, POST ITS INTEREST AND FEES. THE CHARGES COME IN MASTER ORDER, SO THEY ARE POSTED WITHOUT A SORT. | [DAILY.jcl:47](../../../pangaea/jobs/DAILY.jcl#L47) |
| CHARGE | [MONCHG](../programs/MONCHG.md) | BANK.CONTROL(0), ACCT.MASTER(0), ACCT.ACCRUAL(+1) | &&MONTHLY | COND=(4,LT) | | [DAILY.jcl:55](../../../pangaea/jobs/DAILY.jcl#L55) |
| POSTMO | [TXPOST](../programs/TXPOST.md) | BANK.CONTROL(0), ACCT.MASTER(0), &&MONTHLY | &&CHARGED, TRAN.AUDIT(+1), TRAN.REJECTS(+1) | COND=(4,LT) | | [DAILY.jcl:62](../../../pangaea/jobs/DAILY.jcl#L62) |
| SORT | [TXSORT](../programs/TXSORT.md) | &&VALID, &&CARDVAL | &&SORTED | COND=(4,LT) | SORT THE DAY'S TRANSACTIONS, FROM THE BRANCHES AND THE CARDS, INTO ACCOUNT ORDER AND POST THEM | [DAILY.jcl:77](../../../pangaea/jobs/DAILY.jcl#L77) |
| POST | [TXPOST](../programs/TXPOST.md) | BANK.CONTROL(0), &&CHARGED, &&SORTED | &&POSTED, TRAN.AUDIT(+1), TRAN.REJECTS(+1), &&CREDITS | COND=(4,LT) | | [DAILY.jcl:82](../../../pangaea/jobs/DAILY.jcl#L82) |
| SORTCR | [TXSORT](../programs/TXSORT.md) | &&CREDITS | &&CREDSORT | COND=(4,LT) | SORT AND POST THE CREDITS OF THE TRANSFERS THAT POSTED | [DAILY.jcl:96](../../../pangaea/jobs/DAILY.jcl#L96) |
| POSTCR | [TXPOST](../programs/TXPOST.md) | BANK.CONTROL(0), &&POSTED, &&CREDSORT | &&CREDITED, TRAN.AUDIT(+1), TRAN.REJECTS(+1), &&RETURNS | COND=(4,LT) | | [DAILY.jcl:100](../../../pangaea/jobs/DAILY.jcl#L100) |
| SORTRT | [TXSORT](../programs/TXSORT.md) | &&RETURNS | &&RETSORT | COND=(4,LT) | SORT AND POST THE RETURNS OF THE CREDITS THAT COULD NOT POST | [DAILY.jcl:114](../../../pangaea/jobs/DAILY.jcl#L114) |
| POSTRT | [TXPOST](../programs/TXPOST.md) | BANK.CONTROL(0), &&CREDITED, &&RETSORT | &&RETURNED, TRAN.AUDIT(+1), TRAN.REJECTS(+1) | COND=(4,LT) | | [DAILY.jcl:118](../../../pangaea/jobs/DAILY.jcl#L118) |
| MAINT | [ACCTMNT](../programs/ACCTMNT.md) | BANK.CONTROL(0), &&RETURNED, ACCT.MAINT(0) | ACCT.MASTER(+1), MAINT.REJECTS(+1) | COND=(4,LT) | APPLY THE DAY'S ACCOUNT MAINTENANCE, WRITING THE NEW MASTER | [DAILY.jcl:131](../../../pangaea/jobs/DAILY.jcl#L131) |
| LOAD | [ACCTLOAD](../programs/ACCTLOAD.md) | BANK.CONTROL(0), ACCT.MASTER(+1) | ACCT.ONLINE(+1), HOLD.ONLINE(+1) | COND=(4,LT) | LOAD THE NEW MASTER INTO THE ONLINE ACCOUNT FILE, WHICH THE CARD AND ATM AUTHORIZATIONS READ DURING THE NEXT BUSINESS DAY, AND START THAT DAY'S HOLD FILE | [DAILY.jcl:145](../../../pangaea/jobs/DAILY.jcl#L145) |
| NEXTDAY | [BANKDATE](../programs/BANKDATE.md) | BANK.CONTROL(0) | BANK.CONTROL(+1) | COND=(4,LT) | MOVE TO THE NEXT BUSINESS DAY | [DAILY.jcl:156](../../../pangaea/jobs/DAILY.jcl#L156) |
| CARDS | [CARDLOAD](../programs/CARDLOAD.md) | BANK.CONTROL(+1), CARD.MASTER(0) | CARD.ONLINE(+1) | COND=(4,LT) | LOAD THE CARDS INTO THE ONLINE CARD FILE FOR THE NEXT BUSINESS DAY, LEAVING OUT THOSE THAT HAVE EXPIRED BY THEN | [DAILY.jcl:164](../../../pangaea/jobs/DAILY.jcl#L164) |

### VALID

| DD | Dataset | DISP | File | Copybook | Use |
| --- | --- | --- | --- | --- | --- |
| BANKCTL | BANK.CONTROL(0) | SHR | CONTROL-FILE | [BANKCTL](../copybooks/BANKCTL.md) | read |
| TRANIN | TRAN.DAILY(0) | SHR | TRAN-IN | [TRANREC](../copybooks/TRANREC.md) | read |
| TRANOUT | &&VALID | NEW | TRAN-OUT | [TRANREC](../copybooks/TRANREC.md) | write |
| REJECTS | TRAN.REJECTS(+1) | NEW | REJECT-OUT | [REJREC](../copybooks/REJREC.md) | append |
| VALRPT | job output | | REPORT-OUT | | write |

### VALIDC

| DD | Dataset | DISP | File | Copybook | Use |
| --- | --- | --- | --- | --- | --- |
| BANKCTL | BANK.CONTROL(0) | SHR | CONTROL-FILE | [BANKCTL](../copybooks/BANKCTL.md) | read |
| TRANIN | TRAN.CARD(0) | SHR | TRAN-IN | [TRANREC](../copybooks/TRANREC.md) | read |
| TRANOUT | &&CARDVAL | NEW | TRAN-OUT | [TRANREC](../copybooks/TRANREC.md) | write |
| REJECTS | TRAN.REJECTS(+1) | MOD | REJECT-OUT | [REJREC](../copybooks/REJREC.md) | append |
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
| TRANIN | &&VALID, &&CARDVAL | OLD | TRAN-IN | | read |
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
| BANK.CONTROL(0) | [BANKCTL](../copybooks/BANKCTL.md) | | VALID, VALIDC, ACCRUE, CHARGE, POSTMO, POST, POSTCR, POSTRT, MAINT, LOAD, NEXTDAY |
| TRAN.DAILY(0) | [TRANREC](../copybooks/TRANREC.md) | | VALID |
| &&VALID | [TRANREC](../copybooks/TRANREC.md) | VALID | SORT |
| TRAN.REJECTS(+1) | [REJREC](../copybooks/REJREC.md) | VALID, VALIDC, POSTMO, POST, POSTCR, POSTRT | |
| TRAN.CARD(0) | [TRANREC](../copybooks/TRANREC.md) | | VALIDC |
| &&CARDVAL | [TRANREC](../copybooks/TRANREC.md) | VALIDC | SORT |
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

- Step VALID: TXVALID appends to DD REJECTS, which is not DISP=MOD.
