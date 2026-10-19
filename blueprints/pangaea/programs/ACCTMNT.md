# ACCTMNT

```text
ACCTMNT: ACCOUNT MAINTENANCE

APPLIES THE BRANCHES' MAINTENANCE REQUESTS TO THE MASTER AT
THE END OF THE DAY: OPENS ACCOUNTS, CLOSES THEM, FREEZES AND
UNFREEZES THEM, AND CHANGES OVERDRAFT LIMITS. THE REQUESTS
COME IN ACCOUNT NUMBER ORDER AND ARE MATCHED AGAINST THE
MASTER: EACH ACCOUNT IS HELD WHILE ITS REQUESTS APPLY, THEN
WRITTEN TO THE NEW MASTER. A REQUEST THAT CANNOT APPLY GOES TO
THE REJECTS WITH ITS REASON.

RULES: AN OPENING NEEDS AN ACCOUNT NUMBER NOT YET USED; ANY
OTHER REQUEST AN ACCOUNT THAT EXISTS AND IS NOT CLOSED. ONLY
AN EMPTY ACCOUNT CLOSES. ONLY AN OPEN ACCOUNT FREEZES, ONLY A
FROZEN ONE UNFREEZES. ONLY CURRENT ACCOUNTS HAVE OVERDRAFTS.
MAINTENANCE NEVER MOVES MONEY: THE NEW MASTER MUST HOLD THE
SAME BALANCES AS THE OLD.

RETURN CODES: 0 ALL APPLIED, 4 SOME REJECTED,
              16 INPUT OUT OF ORDER, OUT OF BALANCE OR I/O.

MODIFICATION HISTORY
DATE       BY  CHANGE
1988-05-09 JA  CR-0031 ORIGINAL VERSION.
1991-10-07 PK  CR-0129 AN OPENING IS DATED FROM ITS REQUEST'S
               EFFECTIVE DATE.
1994-11-14 PK  CR-0201 OVERDRAFT LIMIT CHANGES.
1999-01-11 SL  Y2K: 8-DIGIT OPEN DATES. THE EFFECTIVE DATE
               LEFT AS IS, WINDOWED ON READ: BELOW 50 IS 20YY.
2009-03-02 MV  CR-0444 THE DATE OF THE REJECTION ON EVERY
               REJECT.
```

Program [ACCTMNT.cbl](../../../pangaea/programs/ACCTMNT.cbl): 6 files, 21 paragraphs. Its business rules: [13 decision tables](../rules/ACCTMNT.md).

- AUTHOR: J. ARKWRIGHT, BATCH SYSTEMS.
- INSTALLATION: FIRST BANK OF PANGAEA, DATA CENTER.
- DATE-WRITTEN: 1988-05-09.

## Contents

- [Files](#files)
- [Flow](#flow)
- [Paragraphs](#paragraphs)
- [Data flow](#data-flow)
- [Lineage](#lineage)
  - [Into BANK-CONTROL-RECORD](#into-bank-control-record)
  - [Into ACCT-IN-RECORD](#into-acct-in-record)
  - [Into MAINT-REQUEST](#into-maint-request)
  - [Into ACCT-OUT-RECORD](#into-acct-out-record)
  - [Into REJECT-OUT-RECORD](#into-reject-out-record)
  - [Into REPORT-LINE](#into-report-line)
  - [Into WS-BUSINESS-DATE](#into-ws-business-date)
  - [Into WS-EFFECTIVE-DATE](#into-ws-effective-date)
  - [Into WS-ABORT-MESSAGE](#into-ws-abort-message)
  - [Into WS-MASTER-KEY](#into-ws-master-key)
  - [Into WS-REQUEST-KEY](#into-ws-request-key)
  - [Into WS-ACTIVE-KEY](#into-ws-active-key)
  - [Into WS-PREV-MASTER-KEY](#into-ws-prev-master-key)
  - [Into WS-REQUEST-ORDER-KEY](#into-ws-request-order-key)
  - [Into WS-PREV-REQUEST-ORDER-KEY](#into-ws-prev-request-order-key)
  - [Into WS-HELD-SW](#into-ws-held-sw)
  - [Into ACCT-MASTER-RECORD](#into-acct-master-record)
  - [Into REJECT-RECORD](#into-reject-record)
  - [Into WS-TOTALS](#into-ws-totals)
  - [Into WS-REPORT-COUNT-LINE](#into-ws-report-count-line)
  - [Into WS-REPORT-AMOUNT-LINE](#into-ws-report-amount-line)
  - [Out of the program](#out-of-the-program)
  - [Into other data items](#into-other-data-items)
- [Field origins](#field-origins)
  - [File REJECT-OUT](#file-reject-out)
  - [File ACCT-OUT](#file-acct-out)
  - [File REPORT-OUT](#file-report-out)
- [Links between records](#links-between-records)
- [Records](#records)
  - [BANK-CONTROL-RECORD](#bank-control-record)
  - [ACCT-IN-RECORD](#acct-in-record)
  - [MAINT-REQUEST](#maint-request)
  - [ACCT-OUT-RECORD](#acct-out-record)
  - [REJECT-OUT-RECORD](#reject-out-record)
  - [REPORT-LINE](#report-line)
  - [WS-FILE-STATUSES](#ws-file-statuses)
  - [WS-BUSINESS-DATE](#ws-business-date)
  - [WS-EFFECTIVE-DATE](#ws-effective-date)
  - [WS-ABORT-MESSAGE](#ws-abort-message)
  - [WS-MASTER-KEY](#ws-master-key)
  - [WS-REQUEST-KEY](#ws-request-key)
  - [WS-ACTIVE-KEY](#ws-active-key)
  - [WS-PREV-MASTER-KEY](#ws-prev-master-key)
  - [WS-REQUEST-ORDER-KEY](#ws-request-order-key)
  - [WS-PREV-REQUEST-ORDER-KEY](#ws-prev-request-order-key)
  - [WS-HELD-SW](#ws-held-sw)
  - [ACCT-MASTER-RECORD](#acct-master-record)
  - [REJECT-RECORD](#reject-record)
  - [WS-TOTALS](#ws-totals)
  - [WS-REPORT-COUNT-LINE](#ws-report-count-line)
  - [WS-REPORT-AMOUNT-LINE](#ws-report-amount-line)

## Files

| File | Kind | Assigned to | Organization | Status field | Opened | Records |
| --- | --- | --- | --- | --- | --- | --- |
| CONTROL-FILE | FD | BANKCTL | SEQUENTIAL | WS-CONTROL-STATUS | INPUT | BANK-CONTROL-RECORD |
| ACCT-IN | FD | ACCTIN | SEQUENTIAL | WS-ACCT-IN-STATUS | INPUT | ACCT-IN-RECORD |
| MAINT-IN | FD | MAINTIN | SEQUENTIAL | WS-MAINT-IN-STATUS | INPUT | MAINT-REQUEST |
| ACCT-OUT | FD | ACCTOUT | SEQUENTIAL | WS-ACCT-OUT-STATUS | OUTPUT | ACCT-OUT-RECORD |
| REJECT-OUT | FD | REJECTS | SEQUENTIAL | WS-REJECT-STATUS | OUTPUT | REJECT-OUT-RECORD |
| REPORT-OUT | FD | MNTRPT | LINE SEQUENTIAL | WS-REPORT-STATUS | OUTPUT | REPORT-LINE |

## Flow

![The paragraphs of ACCTMNT and how control moves between them](ACCTMNT.svg)

## Paragraphs

| Paragraph | Description | Source | Performs | Goes to | Performed by | Reached by |
| --- | --- | --- | --- | --- | --- | --- |
| 0000-MAIN | MATCH UNTIL BOTH FILES END, THEN RECONCILE AND REPORT. | [ACCTMNT.cbl:138](../../../pangaea/programs/ACCTMNT.cbl#L138) | 1000-INITIALIZE, 2000-NEXT-ACCOUNT THRU 2999-EXIT, 9900-ABORT, 8000-WRITE-REPORT, 9000-CLOSE-FILES | | | |
| 1000-INITIALIZE | OPEN THE FILES AND READ THE FIRST RECORD OF EACH. | [ACCTMNT.cbl:155](../../../pangaea/programs/ACCTMNT.cbl#L155) | 9900-ABORT, 1200-READ-MASTER THRU 1200-EXIT, 1300-READ-REQUEST THRU 1300-EXIT | | 0000-MAIN | |
| 1200-READ-MASTER | READ THE NEXT OLD MASTER RECORD, CHECKING ACCOUNT ORDER. | [ACCTMNT.cbl:175](../../../pangaea/programs/ACCTMNT.cbl#L175) | 9900-ABORT | 1200-EXIT, 1200-EXIT (falls into) | 1000-INITIALIZE, 2000-NEXT-ACCOUNT | |
| 1200-EXIT | | [ACCTMNT.cbl:186](../../../pangaea/programs/ACCTMNT.cbl#L186) | | | 1000-INITIALIZE, 2000-NEXT-ACCOUNT | 1200-READ-MASTER |
| 1300-READ-REQUEST | READ THE NEXT REQUEST, CHECKING ITS ORDER. | [ACCTMNT.cbl:190](../../../pangaea/programs/ACCTMNT.cbl#L190) | 9900-ABORT | 1300-EXIT, 1300-EXIT (falls into) | 1000-INITIALIZE, 2100-APPLY-REQUESTS | |
| 1300-EXIT | | [ACCTMNT.cbl:202](../../../pangaea/programs/ACCTMNT.cbl#L202) | | | 1000-INITIALIZE, 2100-APPLY-REQUESTS | 1300-READ-REQUEST |
| 2000-NEXT-ACCOUNT | TAKE THE LOWER KEY AS THE ACTIVE ACCOUNT, AND HOLD IT IF THE MASTER HAS IT. AN ACCOUNT ONLY THE REQUESTS NAME IS NOT HELD UNTIL AN OPENING MAKES IT. | [ACCTMNT.cbl:208](../../../pangaea/programs/ACCTMNT.cbl#L208) | 1200-READ-MASTER THRU 1200-EXIT | 2100-APPLY-REQUESTS (falls into) | 0000-MAIN | |
| 2100-APPLY-REQUESTS | APPLY THE ACTIVE ACCOUNT'S REQUESTS, ONE AFTER THE OTHER. | [ACCTMNT.cbl:220](../../../pangaea/programs/ACCTMNT.cbl#L220) | 2200-APPLY-REQUEST THRU 2299-EXIT, 1300-READ-REQUEST THRU 1300-EXIT | 2800-WRITE-HELD, 2100-APPLY-REQUESTS | 0000-MAIN | 2000-NEXT-ACCOUNT, 2100-APPLY-REQUESTS |
| 2200-APPLY-REQUEST | APPLY ONE REQUEST BY ITS TYPE. | [ACCTMNT.cbl:228](../../../pangaea/programs/ACCTMNT.cbl#L228) | | 2210-OPEN, 2290-REJECT, 2220-CLOSE, 2230-FREEZE, 2240-UNFREEZE, 2250-LIMIT | 0000-MAIN, 2100-APPLY-REQUESTS | |
| 2210-OPEN | OPEN A NEW ACCOUNT, EMPTY, FROM THE REQUEST'S EFFECTIVE DATE. | [ACCTMNT.cbl:249](../../../pangaea/programs/ACCTMNT.cbl#L249) | | 2290-REJECT, 2299-EXIT | 0000-MAIN, 2100-APPLY-REQUESTS | 2200-APPLY-REQUEST |
| 2220-CLOSE | CLOSE AN ACCOUNT, IF IT IS EMPTY. | [ACCTMNT.cbl:276](../../../pangaea/programs/ACCTMNT.cbl#L276) | | 2290-REJECT, 2299-EXIT | 0000-MAIN, 2100-APPLY-REQUESTS | 2200-APPLY-REQUEST |
| 2230-FREEZE | FREEZE AN OPEN ACCOUNT. | [ACCTMNT.cbl:286](../../../pangaea/programs/ACCTMNT.cbl#L286) | | 2290-REJECT, 2299-EXIT | 0000-MAIN, 2100-APPLY-REQUESTS | 2200-APPLY-REQUEST |
| 2240-UNFREEZE | UNFREEZE A FROZEN ACCOUNT. | [ACCTMNT.cbl:295](../../../pangaea/programs/ACCTMNT.cbl#L295) | | 2290-REJECT, 2299-EXIT | 0000-MAIN, 2100-APPLY-REQUESTS | 2200-APPLY-REQUEST |
| 2250-LIMIT | CHANGE A CURRENT ACCOUNT'S OVERDRAFT LIMIT. | [ACCTMNT.cbl:304](../../../pangaea/programs/ACCTMNT.cbl#L304) | | 2290-REJECT, 2299-EXIT | 0000-MAIN, 2100-APPLY-REQUESTS | 2200-APPLY-REQUEST |
| 2290-REJECT | WRITE THE REQUEST TO THE REJECTS WITH ITS REASON. | [ACCTMNT.cbl:313](../../../pangaea/programs/ACCTMNT.cbl#L313) | | 2299-EXIT (falls into) | 0000-MAIN, 2100-APPLY-REQUESTS | 2200-APPLY-REQUEST, 2210-OPEN, 2220-CLOSE, 2230-FREEZE, 2240-UNFREEZE, 2250-LIMIT |
| 2299-EXIT | | [ACCTMNT.cbl:318](../../../pangaea/programs/ACCTMNT.cbl#L318) | | | 0000-MAIN, 2100-APPLY-REQUESTS | 2210-OPEN, 2220-CLOSE, 2230-FREEZE, 2240-UNFREEZE, 2250-LIMIT, 2290-REJECT |
| 2800-WRITE-HELD | WRITE THE HELD ACCOUNT, ITS REQUESTS APPLIED, TO THE NEW MASTER. | [ACCTMNT.cbl:322](../../../pangaea/programs/ACCTMNT.cbl#L322) | | 2999-EXIT (falls into) | 0000-MAIN | 2100-APPLY-REQUESTS |
| 2999-EXIT | | [ACCTMNT.cbl:327](../../../pangaea/programs/ACCTMNT.cbl#L327) | | | 0000-MAIN | 2800-WRITE-HELD |
| 8000-WRITE-REPORT | WRITE THE MAINTENANCE REPORT. | [ACCTMNT.cbl:331](../../../pangaea/programs/ACCTMNT.cbl#L331) | | | 0000-MAIN | |
| 9000-CLOSE-FILES | CLOSE ALL FILES. | [ACCTMNT.cbl:377](../../../pangaea/programs/ACCTMNT.cbl#L377) | | | 0000-MAIN, 9900-ABORT | |
| 9900-ABORT | STOP THE RUN WITH RETURN CODE 16, SAYING WHY. | [ACCTMNT.cbl:382](../../../pangaea/programs/ACCTMNT.cbl#L382) | 9000-CLOSE-FILES | | 0000-MAIN, 1000-INITIALIZE, 1200-READ-MASTER, 1300-READ-REQUEST | |

## Data flow

| Paragraph | Reads | Writes |
| --- | --- | --- |
| 0000-MAIN | WS-MASTER-KEY, WS-REQUEST-KEY, WS-BALANCES-OUT, WS-BALANCES-IN, WS-REJECT-COUNT | WS-ABORT-MESSAGE, RETURN-CODE |
| 1000-INITIALIZE | WS-FILE-STATUSES, BC-BUSINESS-DATE, WS-BUSINESS-DATE | RETURN-CODE, BANK-CONTROL-RECORD, WS-ABORT-MESSAGE, WS-BUSINESS-DATE, REJECT-RECORD, RJ-REJECT-DATE |
| 1200-READ-MASTER | AI-ACCT-NO, WS-MASTER-KEY, WS-PREV-MASTER-KEY, WS-MASTERS-READ | ACCT-IN-RECORD, WS-MASTER-KEY, WS-ABORT-MESSAGE, WS-PREV-MASTER-KEY, WS-MASTERS-READ |
| 1200-EXIT | | |
| 1300-READ-REQUEST | MR-ACCT-NO, MR-REQUEST-ID, WS-REQUEST-ORDER-KEY, WS-PREV-REQUEST-ORDER-KEY, WS-REQUEST-COUNT | MAINT-REQUEST, WS-REQUEST-KEY, WS-ROK-ACCT-NO, WS-ROK-REQUEST-ID, WS-ABORT-MESSAGE, WS-PREV-REQUEST-ORDER-KEY, WS-REQUEST-COUNT |
| 1300-EXIT | | |
| 2000-NEXT-ACCOUNT | WS-MASTER-KEY, WS-REQUEST-KEY, WS-ACTIVE-KEY, ACCT-IN-RECORD, AM-BALANCE, WS-BALANCES-IN | WS-ACTIVE-KEY, WS-HELD-SW, ACCT-MASTER-RECORD, WS-BALANCES-IN |
| 2100-APPLY-REQUESTS | WS-REQUEST-KEY, WS-ACTIVE-KEY | |
| 2200-APPLY-REQUEST | MR-REQUEST-TYPE, WS-HELD-SW, AM-STATUS | RJ-REASON |
| 2210-OPEN | WS-HELD-SW, MR-ACCT-NO, MR-CUST-NO, MR-SHORT-NAME, MR-ACCT-TYPE, MR-EFFECTIVE-DATE, WS-EFF-YY, WS-EFFECTIVE-DATE, WS-BUSINESS-DATE, MR-OVERDRAFT-LIMIT, AM-ACCT-TYPE, WS-OPENED-COUNT | RJ-REASON, ACCT-MASTER-RECORD, AM-ACCT-NO, AM-CUST-NO, AM-SHORT-NAME, AM-ACCT-TYPE, AM-STATUS, WS-EFF-YYMMDD, WS-EFF-CENTURY, AM-OPEN-DATE, AM-LAST-ACTIVITY-DATE, AM-BALANCE, AM-OVERDRAFT-LIMIT, WS-HELD-SW, WS-OPENED-COUNT |
| 2220-CLOSE | AM-BALANCE, WS-CLOSED-COUNT | RJ-REASON, AM-STATUS, AM-OVERDRAFT-LIMIT, WS-CLOSED-COUNT |
| 2230-FREEZE | AM-STATUS, WS-FROZEN-COUNT | RJ-REASON, AM-STATUS, WS-FROZEN-COUNT |
| 2240-UNFREEZE | AM-STATUS, WS-UNFROZEN-COUNT | RJ-REASON, AM-STATUS, WS-UNFROZEN-COUNT |
| 2250-LIMIT | AM-ACCT-TYPE, MR-OVERDRAFT-LIMIT, WS-LIMIT-COUNT | RJ-REASON, AM-OVERDRAFT-LIMIT, WS-LIMIT-COUNT |
| 2290-REJECT | MAINT-REQUEST, REJECT-RECORD, WS-REJECT-COUNT | RJ-PROGRAM, RJ-TRAN, REJECT-OUT-RECORD, WS-REJECT-COUNT |
| 2299-EXIT | | |
| 2800-WRITE-HELD | WS-HELD-SW, ACCT-MASTER-RECORD, WS-MASTERS-WRITTEN, AM-BALANCE, WS-BALANCES-OUT | ACCT-OUT-RECORD, WS-MASTERS-WRITTEN, WS-BALANCES-OUT |
| 2999-EXIT | | |
| 8000-WRITE-REPORT | WS-BUSINESS-DATE, WS-MASTERS-READ, WS-REPORT-COUNT-LINE, WS-MASTERS-WRITTEN, WS-REQUEST-COUNT, WS-OPENED-COUNT, WS-CLOSED-COUNT, WS-FROZEN-COUNT, WS-UNFROZEN-COUNT, WS-LIMIT-COUNT, WS-REJECT-COUNT, WS-BALANCES-OUT, WS-REPORT-AMOUNT-LINE | REPORT-LINE, WS-RC-LABEL, WS-RC-VALUE, WS-RA-LABEL, WS-RA-VALUE |
| 9000-CLOSE-FILES | | |
| 9900-ABORT | WS-ABORT-MESSAGE | RETURN-CODE |

## Lineage

Where each data item's values come from, grouped by the record they go into. A data item links to the record it belongs to.

### Into BANK-CONTROL-RECORD

Record [BANK-CONTROL-RECORD](#bank-control-record), from copybook [BANKCTL](../copybooks/BANKCTL.md).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| BANK-CONTROL-RECORD | file CONTROL-FILE | READ | 1000-INITIALIZE | [ACCTMNT.cbl:162](../../../pangaea/programs/ACCTMNT.cbl#L162) |

### Into ACCT-IN-RECORD

Record [ACCT-IN-RECORD](#acct-in-record).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| ACCT-IN-RECORD | file ACCT-IN | READ | 1200-READ-MASTER | [ACCTMNT.cbl:176](../../../pangaea/programs/ACCTMNT.cbl#L176) |

### Into MAINT-REQUEST

Record [MAINT-REQUEST](#maint-request), from copybook [MAINTREC](../copybooks/MAINTREC.md).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| MAINT-REQUEST | file MAINT-IN | READ | 1300-READ-REQUEST | [ACCTMNT.cbl:191](../../../pangaea/programs/ACCTMNT.cbl#L191) |

### Into ACCT-OUT-RECORD

Record [ACCT-OUT-RECORD](#acct-out-record).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| ACCT-OUT-RECORD | [ACCT-MASTER-RECORD](#acct-master-record) | WRITE | 2800-WRITE-HELD | [ACCTMNT.cbl:324](../../../pangaea/programs/ACCTMNT.cbl#L324) |

### Into REJECT-OUT-RECORD

Record [REJECT-OUT-RECORD](#reject-out-record).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| REJECT-OUT-RECORD | [REJECT-RECORD](#reject-record) | WRITE | 2290-REJECT | [ACCTMNT.cbl:316](../../../pangaea/programs/ACCTMNT.cbl#L316) |

### Into REPORT-LINE

Record [REPORT-LINE](#report-line).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| REPORT-LINE | literal SPACES | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:332](../../../pangaea/programs/ACCTMNT.cbl#L332) |
| REPORT-LINE | literal 'ACCTMNT  ACCOUNT MAINTENANCE     BUSINESS DATE ' | STRING | 8000-WRITE-REPORT | [ACCTMNT.cbl:333](../../../pangaea/programs/ACCTMNT.cbl#L333) |
| REPORT-LINE | [WS-BUSINESS-DATE](#ws-business-date) | STRING | 8000-WRITE-REPORT | [ACCTMNT.cbl:333](../../../pangaea/programs/ACCTMNT.cbl#L333) |
| REPORT-LINE | literal SPACES | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:337](../../../pangaea/programs/ACCTMNT.cbl#L337) |
| REPORT-LINE | [WS-REPORT-COUNT-LINE](#ws-report-count-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:341](../../../pangaea/programs/ACCTMNT.cbl#L341) |
| REPORT-LINE | [WS-REPORT-COUNT-LINE](#ws-report-count-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:344](../../../pangaea/programs/ACCTMNT.cbl#L344) |
| REPORT-LINE | [WS-REPORT-COUNT-LINE](#ws-report-count-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:347](../../../pangaea/programs/ACCTMNT.cbl#L347) |
| REPORT-LINE | [WS-REPORT-COUNT-LINE](#ws-report-count-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:350](../../../pangaea/programs/ACCTMNT.cbl#L350) |
| REPORT-LINE | [WS-REPORT-COUNT-LINE](#ws-report-count-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:353](../../../pangaea/programs/ACCTMNT.cbl#L353) |
| REPORT-LINE | [WS-REPORT-COUNT-LINE](#ws-report-count-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:356](../../../pangaea/programs/ACCTMNT.cbl#L356) |
| REPORT-LINE | [WS-REPORT-COUNT-LINE](#ws-report-count-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:359](../../../pangaea/programs/ACCTMNT.cbl#L359) |
| REPORT-LINE | [WS-REPORT-COUNT-LINE](#ws-report-count-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:362](../../../pangaea/programs/ACCTMNT.cbl#L362) |
| REPORT-LINE | [WS-REPORT-COUNT-LINE](#ws-report-count-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:365](../../../pangaea/programs/ACCTMNT.cbl#L365) |
| REPORT-LINE | literal SPACES | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:366](../../../pangaea/programs/ACCTMNT.cbl#L366) |
| REPORT-LINE | [WS-REPORT-AMOUNT-LINE](#ws-report-amount-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:370](../../../pangaea/programs/ACCTMNT.cbl#L370) |
| REPORT-LINE | literal SPACES | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:371](../../../pangaea/programs/ACCTMNT.cbl#L371) |
| REPORT-LINE | literal 'NEW MASTER IN BALANCE' | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:373](../../../pangaea/programs/ACCTMNT.cbl#L373) |

### Into WS-BUSINESS-DATE

Record [WS-BUSINESS-DATE](#ws-business-date).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-BUSINESS-DATE | [BC-BUSINESS-DATE](#bank-control-record "in BANK-CONTROL-RECORD") | MOVE | 1000-INITIALIZE | [ACCTMNT.cbl:166](../../../pangaea/programs/ACCTMNT.cbl#L166) |

### Into WS-EFFECTIVE-DATE

Record [WS-EFFECTIVE-DATE](#ws-effective-date).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-EFF-CENTURY | literal 20 | MOVE | 2210-OPEN | [ACCTMNT.cbl:261](../../../pangaea/programs/ACCTMNT.cbl#L261) |
| WS-EFF-CENTURY | literal 19 | MOVE | 2210-OPEN | [ACCTMNT.cbl:263](../../../pangaea/programs/ACCTMNT.cbl#L263) |
| WS-EFF-YYMMDD | [MR-EFFECTIVE-DATE](#maint-request "in MAINT-REQUEST") | MOVE | 2210-OPEN | [ACCTMNT.cbl:259](../../../pangaea/programs/ACCTMNT.cbl#L259) |

### Into WS-ABORT-MESSAGE

Record [WS-ABORT-MESSAGE](#ws-abort-message).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-ABORT-MESSAGE | literal 'NEW MASTER OUT OF BALANCE' | MOVE | 0000-MAIN | [ACCTMNT.cbl:144](../../../pangaea/programs/ACCTMNT.cbl#L144) |
| WS-ABORT-MESSAGE | literal 'NO BANK CONTROL RECORD' | MOVE | 1000-INITIALIZE | [ACCTMNT.cbl:164](../../../pangaea/programs/ACCTMNT.cbl#L164) |
| WS-ABORT-MESSAGE | literal 'MASTER NOT IN ACCOUNT ORDER' | MOVE | 1200-READ-MASTER | [ACCTMNT.cbl:182](../../../pangaea/programs/ACCTMNT.cbl#L182) |
| WS-ABORT-MESSAGE | literal 'REQUESTS NOT IN ACCOUNT ORDER' | MOVE | 1300-READ-REQUEST | [ACCTMNT.cbl:198](../../../pangaea/programs/ACCTMNT.cbl#L198) |

### Into WS-MASTER-KEY

Record [WS-MASTER-KEY](#ws-master-key).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-MASTER-KEY | literal HIGH-VALUES | MOVE | 1200-READ-MASTER | [ACCTMNT.cbl:178](../../../pangaea/programs/ACCTMNT.cbl#L178) |
| WS-MASTER-KEY | [AI-ACCT-NO](#acct-in-record "in ACCT-IN-RECORD") | MOVE | 1200-READ-MASTER | [ACCTMNT.cbl:180](../../../pangaea/programs/ACCTMNT.cbl#L180) |

### Into WS-REQUEST-KEY

Record [WS-REQUEST-KEY](#ws-request-key).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-REQUEST-KEY | literal HIGH-VALUES | MOVE | 1300-READ-REQUEST | [ACCTMNT.cbl:193](../../../pangaea/programs/ACCTMNT.cbl#L193) |
| WS-REQUEST-KEY | [MR-ACCT-NO](#maint-request "in MAINT-REQUEST") | MOVE | 1300-READ-REQUEST | [ACCTMNT.cbl:195](../../../pangaea/programs/ACCTMNT.cbl#L195) |

### Into WS-ACTIVE-KEY

Record [WS-ACTIVE-KEY](#ws-active-key).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-ACTIVE-KEY | [WS-MASTER-KEY](#ws-master-key) | MOVE | 2000-NEXT-ACCOUNT | [ACCTMNT.cbl:209](../../../pangaea/programs/ACCTMNT.cbl#L209) |
| WS-ACTIVE-KEY | [WS-REQUEST-KEY](#ws-request-key) | MOVE | 2000-NEXT-ACCOUNT | [ACCTMNT.cbl:211](../../../pangaea/programs/ACCTMNT.cbl#L211) |

### Into WS-PREV-MASTER-KEY

Record [WS-PREV-MASTER-KEY](#ws-prev-master-key).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-PREV-MASTER-KEY | [WS-MASTER-KEY](#ws-master-key) | MOVE | 1200-READ-MASTER | [ACCTMNT.cbl:184](../../../pangaea/programs/ACCTMNT.cbl#L184) |

### Into WS-REQUEST-ORDER-KEY

Record [WS-REQUEST-ORDER-KEY](#ws-request-order-key).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-ROK-ACCT-NO | [MR-ACCT-NO](#maint-request "in MAINT-REQUEST") | MOVE | 1300-READ-REQUEST | [ACCTMNT.cbl:195](../../../pangaea/programs/ACCTMNT.cbl#L195) |
| WS-ROK-REQUEST-ID | [MR-REQUEST-ID](#maint-request "in MAINT-REQUEST") | MOVE | 1300-READ-REQUEST | [ACCTMNT.cbl:196](../../../pangaea/programs/ACCTMNT.cbl#L196) |

### Into WS-PREV-REQUEST-ORDER-KEY

Record [WS-PREV-REQUEST-ORDER-KEY](#ws-prev-request-order-key).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-PREV-REQUEST-ORDER-KEY | [WS-REQUEST-ORDER-KEY](#ws-request-order-key) | MOVE | 1300-READ-REQUEST | [ACCTMNT.cbl:200](../../../pangaea/programs/ACCTMNT.cbl#L200) |

### Into WS-HELD-SW

Record [WS-HELD-SW](#ws-held-sw).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-HELD-SW | literal condition NOTHING-HELD | SET | 2000-NEXT-ACCOUNT | [ACCTMNT.cbl:212](../../../pangaea/programs/ACCTMNT.cbl#L212) |
| WS-HELD-SW | literal condition ACCOUNT-HELD | SET | 2000-NEXT-ACCOUNT | [ACCTMNT.cbl:216](../../../pangaea/programs/ACCTMNT.cbl#L216) |
| WS-HELD-SW | literal condition ACCOUNT-HELD | SET | 2210-OPEN | [ACCTMNT.cbl:271](../../../pangaea/programs/ACCTMNT.cbl#L271) |

### Into ACCT-MASTER-RECORD

Record [ACCT-MASTER-RECORD](#acct-master-record), from copybook [ACCTMAST](../copybooks/ACCTMAST.md).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| ACCT-MASTER-RECORD | [ACCT-IN-RECORD](#acct-in-record) | MOVE | 2000-NEXT-ACCOUNT | [ACCTMNT.cbl:214](../../../pangaea/programs/ACCTMNT.cbl#L214) |
| ACCT-MASTER-RECORD | literal SPACES | MOVE | 2210-OPEN | [ACCTMNT.cbl:253](../../../pangaea/programs/ACCTMNT.cbl#L253) |
| AM-ACCT-NO | [MR-ACCT-NO](#maint-request "in MAINT-REQUEST") | MOVE | 2210-OPEN | [ACCTMNT.cbl:254](../../../pangaea/programs/ACCTMNT.cbl#L254) |
| AM-ACCT-TYPE | [MR-ACCT-TYPE](#maint-request "in MAINT-REQUEST") | MOVE | 2210-OPEN | [ACCTMNT.cbl:257](../../../pangaea/programs/ACCTMNT.cbl#L257) |
| AM-BALANCE | literal 0 | MOVE | 2210-OPEN | [ACCTMNT.cbl:267](../../../pangaea/programs/ACCTMNT.cbl#L267) |
| AM-CUST-NO | [MR-CUST-NO](#maint-request "in MAINT-REQUEST") | MOVE | 2210-OPEN | [ACCTMNT.cbl:255](../../../pangaea/programs/ACCTMNT.cbl#L255) |
| AM-LAST-ACTIVITY-DATE | [WS-BUSINESS-DATE](#ws-business-date) | MOVE | 2210-OPEN | [ACCTMNT.cbl:266](../../../pangaea/programs/ACCTMNT.cbl#L266) |
| AM-OPEN-DATE | [WS-EFFECTIVE-DATE](#ws-effective-date) | MOVE | 2210-OPEN | [ACCTMNT.cbl:265](../../../pangaea/programs/ACCTMNT.cbl#L265) |
| AM-OVERDRAFT-LIMIT | [MR-OVERDRAFT-LIMIT](#maint-request "in MAINT-REQUEST") | MOVE | 2210-OPEN | [ACCTMNT.cbl:268](../../../pangaea/programs/ACCTMNT.cbl#L268) |
| AM-OVERDRAFT-LIMIT | literal 0 | MOVE | 2210-OPEN | [ACCTMNT.cbl:270](../../../pangaea/programs/ACCTMNT.cbl#L270) |
| AM-OVERDRAFT-LIMIT | literal 0 | MOVE | 2220-CLOSE | [ACCTMNT.cbl:281](../../../pangaea/programs/ACCTMNT.cbl#L281) |
| AM-OVERDRAFT-LIMIT | [MR-OVERDRAFT-LIMIT](#maint-request "in MAINT-REQUEST") | MOVE | 2250-LIMIT | [ACCTMNT.cbl:308](../../../pangaea/programs/ACCTMNT.cbl#L308) |
| AM-SHORT-NAME | [MR-SHORT-NAME](#maint-request "in MAINT-REQUEST") | MOVE | 2210-OPEN | [ACCTMNT.cbl:256](../../../pangaea/programs/ACCTMNT.cbl#L256) |
| AM-STATUS | literal condition AM-OPEN | SET | 2210-OPEN | [ACCTMNT.cbl:258](../../../pangaea/programs/ACCTMNT.cbl#L258) |
| AM-STATUS | literal condition AM-CLOSED | SET | 2220-CLOSE | [ACCTMNT.cbl:280](../../../pangaea/programs/ACCTMNT.cbl#L280) |
| AM-STATUS | literal condition AM-FROZEN | SET | 2230-FREEZE | [ACCTMNT.cbl:290](../../../pangaea/programs/ACCTMNT.cbl#L290) |
| AM-STATUS | literal condition AM-OPEN | SET | 2240-UNFREEZE | [ACCTMNT.cbl:299](../../../pangaea/programs/ACCTMNT.cbl#L299) |

### Into REJECT-RECORD

Record [REJECT-RECORD](#reject-record), from copybook [REJREC](../copybooks/REJREC.md).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| REJECT-RECORD | literal SPACES | MOVE | 1000-INITIALIZE | [ACCTMNT.cbl:167](../../../pangaea/programs/ACCTMNT.cbl#L167) |
| RJ-PROGRAM | literal 'ACCTMNT' | MOVE | 2290-REJECT | [ACCTMNT.cbl:314](../../../pangaea/programs/ACCTMNT.cbl#L314) |
| RJ-REASON | literal condition RJ-NO-ACCOUNT | SET | 2200-APPLY-REQUEST | [ACCTMNT.cbl:232](../../../pangaea/programs/ACCTMNT.cbl#L232) |
| RJ-REASON | literal condition RJ-CLOSED | SET | 2200-APPLY-REQUEST | [ACCTMNT.cbl:235](../../../pangaea/programs/ACCTMNT.cbl#L235) |
| RJ-REASON | literal condition RJ-BAD-TYPE | SET | 2200-APPLY-REQUEST | [ACCTMNT.cbl:245](../../../pangaea/programs/ACCTMNT.cbl#L245) |
| RJ-REASON | literal condition RJ-EXISTS | SET | 2210-OPEN | [ACCTMNT.cbl:251](../../../pangaea/programs/ACCTMNT.cbl#L251) |
| RJ-REASON | literal condition RJ-HAS-BALANCE | SET | 2220-CLOSE | [ACCTMNT.cbl:278](../../../pangaea/programs/ACCTMNT.cbl#L278) |
| RJ-REASON | literal condition RJ-BAD-STATUS | SET | 2230-FREEZE | [ACCTMNT.cbl:288](../../../pangaea/programs/ACCTMNT.cbl#L288) |
| RJ-REASON | literal condition RJ-BAD-STATUS | SET | 2240-UNFREEZE | [ACCTMNT.cbl:297](../../../pangaea/programs/ACCTMNT.cbl#L297) |
| RJ-REASON | literal condition RJ-BAD-TYPE | SET | 2250-LIMIT | [ACCTMNT.cbl:306](../../../pangaea/programs/ACCTMNT.cbl#L306) |
| RJ-REJECT-DATE | [WS-BUSINESS-DATE](#ws-business-date) | STRING | 1000-INITIALIZE | [ACCTMNT.cbl:168](../../../pangaea/programs/ACCTMNT.cbl#L168) |
| RJ-REJECT-DATE | literal '-' | STRING | 1000-INITIALIZE | [ACCTMNT.cbl:168](../../../pangaea/programs/ACCTMNT.cbl#L168) |
| RJ-TRAN | [MAINT-REQUEST](#maint-request) | MOVE | 2290-REJECT | [ACCTMNT.cbl:315](../../../pangaea/programs/ACCTMNT.cbl#L315) |

### Into WS-TOTALS

Record [WS-TOTALS](#ws-totals).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-BALANCES-IN | [AM-BALANCE](#acct-master-record "in ACCT-MASTER-RECORD") | ADD | 2000-NEXT-ACCOUNT | [ACCTMNT.cbl:215](../../../pangaea/programs/ACCTMNT.cbl#L215) |
| WS-BALANCES-OUT | [AM-BALANCE](#acct-master-record "in ACCT-MASTER-RECORD") | ADD | 2800-WRITE-HELD | [ACCTMNT.cbl:326](../../../pangaea/programs/ACCTMNT.cbl#L326) |
| WS-CLOSED-COUNT | literal 1 | ADD | 2220-CLOSE | [ACCTMNT.cbl:282](../../../pangaea/programs/ACCTMNT.cbl#L282) |
| WS-FROZEN-COUNT | literal 1 | ADD | 2230-FREEZE | [ACCTMNT.cbl:291](../../../pangaea/programs/ACCTMNT.cbl#L291) |
| WS-LIMIT-COUNT | literal 1 | ADD | 2250-LIMIT | [ACCTMNT.cbl:309](../../../pangaea/programs/ACCTMNT.cbl#L309) |
| WS-MASTERS-READ | literal 1 | ADD | 1200-READ-MASTER | [ACCTMNT.cbl:185](../../../pangaea/programs/ACCTMNT.cbl#L185) |
| WS-MASTERS-WRITTEN | literal 1 | ADD | 2800-WRITE-HELD | [ACCTMNT.cbl:325](../../../pangaea/programs/ACCTMNT.cbl#L325) |
| WS-OPENED-COUNT | literal 1 | ADD | 2210-OPEN | [ACCTMNT.cbl:272](../../../pangaea/programs/ACCTMNT.cbl#L272) |
| WS-REJECT-COUNT | literal 1 | ADD | 2290-REJECT | [ACCTMNT.cbl:317](../../../pangaea/programs/ACCTMNT.cbl#L317) |
| WS-REQUEST-COUNT | literal 1 | ADD | 1300-READ-REQUEST | [ACCTMNT.cbl:201](../../../pangaea/programs/ACCTMNT.cbl#L201) |
| WS-UNFROZEN-COUNT | literal 1 | ADD | 2240-UNFREEZE | [ACCTMNT.cbl:300](../../../pangaea/programs/ACCTMNT.cbl#L300) |

### Into WS-REPORT-COUNT-LINE

Record [WS-REPORT-COUNT-LINE](#ws-report-count-line).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-RC-LABEL | literal 'ACCOUNTS READ' | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:339](../../../pangaea/programs/ACCTMNT.cbl#L339) |
| WS-RC-LABEL | literal 'ACCOUNTS WRITTEN' | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:342](../../../pangaea/programs/ACCTMNT.cbl#L342) |
| WS-RC-LABEL | literal 'REQUESTS READ' | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:345](../../../pangaea/programs/ACCTMNT.cbl#L345) |
| WS-RC-LABEL | literal 'ACCOUNTS OPENED' | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:348](../../../pangaea/programs/ACCTMNT.cbl#L348) |
| WS-RC-LABEL | literal 'ACCOUNTS CLOSED' | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:351](../../../pangaea/programs/ACCTMNT.cbl#L351) |
| WS-RC-LABEL | literal 'ACCOUNTS FROZEN' | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:354](../../../pangaea/programs/ACCTMNT.cbl#L354) |
| WS-RC-LABEL | literal 'ACCOUNTS UNFROZEN' | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:357](../../../pangaea/programs/ACCTMNT.cbl#L357) |
| WS-RC-LABEL | literal 'OVERDRAFT LIMITS CHANGED' | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:360](../../../pangaea/programs/ACCTMNT.cbl#L360) |
| WS-RC-LABEL | literal 'REQUESTS REJECTED' | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:363](../../../pangaea/programs/ACCTMNT.cbl#L363) |
| WS-RC-VALUE | [WS-MASTERS-READ](#ws-totals "in WS-TOTALS") | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:340](../../../pangaea/programs/ACCTMNT.cbl#L340) |
| WS-RC-VALUE | [WS-MASTERS-WRITTEN](#ws-totals "in WS-TOTALS") | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:343](../../../pangaea/programs/ACCTMNT.cbl#L343) |
| WS-RC-VALUE | [WS-REQUEST-COUNT](#ws-totals "in WS-TOTALS") | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:346](../../../pangaea/programs/ACCTMNT.cbl#L346) |
| WS-RC-VALUE | [WS-OPENED-COUNT](#ws-totals "in WS-TOTALS") | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:349](../../../pangaea/programs/ACCTMNT.cbl#L349) |
| WS-RC-VALUE | [WS-CLOSED-COUNT](#ws-totals "in WS-TOTALS") | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:352](../../../pangaea/programs/ACCTMNT.cbl#L352) |
| WS-RC-VALUE | [WS-FROZEN-COUNT](#ws-totals "in WS-TOTALS") | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:355](../../../pangaea/programs/ACCTMNT.cbl#L355) |
| WS-RC-VALUE | [WS-UNFROZEN-COUNT](#ws-totals "in WS-TOTALS") | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:358](../../../pangaea/programs/ACCTMNT.cbl#L358) |
| WS-RC-VALUE | [WS-LIMIT-COUNT](#ws-totals "in WS-TOTALS") | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:361](../../../pangaea/programs/ACCTMNT.cbl#L361) |
| WS-RC-VALUE | [WS-REJECT-COUNT](#ws-totals "in WS-TOTALS") | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:364](../../../pangaea/programs/ACCTMNT.cbl#L364) |

### Into WS-REPORT-AMOUNT-LINE

Record [WS-REPORT-AMOUNT-LINE](#ws-report-amount-line).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-RA-LABEL | literal 'BALANCES' | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:368](../../../pangaea/programs/ACCTMNT.cbl#L368) |
| WS-RA-VALUE | [WS-BALANCES-OUT](#ws-totals "in WS-TOTALS") | MOVE | 8000-WRITE-REPORT | [ACCTMNT.cbl:369](../../../pangaea/programs/ACCTMNT.cbl#L369) |

### Out of the program

Records written to files, items passed to the programs called, and items displayed to the job's output.

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| file ACCT-OUT | [ACCT-OUT-RECORD](#acct-out-record) | WRITE | 2800-WRITE-HELD | [ACCTMNT.cbl:324](../../../pangaea/programs/ACCTMNT.cbl#L324) |
| file REJECT-OUT | [REJECT-OUT-RECORD](#reject-out-record) | WRITE | 2290-REJECT | [ACCTMNT.cbl:316](../../../pangaea/programs/ACCTMNT.cbl#L316) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:336](../../../pangaea/programs/ACCTMNT.cbl#L336) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:338](../../../pangaea/programs/ACCTMNT.cbl#L338) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:341](../../../pangaea/programs/ACCTMNT.cbl#L341) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:344](../../../pangaea/programs/ACCTMNT.cbl#L344) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:347](../../../pangaea/programs/ACCTMNT.cbl#L347) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:350](../../../pangaea/programs/ACCTMNT.cbl#L350) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:353](../../../pangaea/programs/ACCTMNT.cbl#L353) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:356](../../../pangaea/programs/ACCTMNT.cbl#L356) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:359](../../../pangaea/programs/ACCTMNT.cbl#L359) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:362](../../../pangaea/programs/ACCTMNT.cbl#L362) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:365](../../../pangaea/programs/ACCTMNT.cbl#L365) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:367](../../../pangaea/programs/ACCTMNT.cbl#L367) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:370](../../../pangaea/programs/ACCTMNT.cbl#L370) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:372](../../../pangaea/programs/ACCTMNT.cbl#L372) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [ACCTMNT.cbl:374](../../../pangaea/programs/ACCTMNT.cbl#L374) |
| SYSOUT | [WS-FILE-STATUSES](#ws-file-statuses) | DISPLAY | 1000-INITIALIZE | [ACCTMNT.cbl:159](../../../pangaea/programs/ACCTMNT.cbl#L159) |
| SYSOUT | [WS-ABORT-MESSAGE](#ws-abort-message) | DISPLAY | 9900-ABORT | [ACCTMNT.cbl:383](../../../pangaea/programs/ACCTMNT.cbl#L383) |

### Into other data items

Data items the program does not declare, such as special registers and index names.

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| RETURN-CODE | literal 4 | MOVE | 0000-MAIN | [ACCTMNT.cbl:149](../../../pangaea/programs/ACCTMNT.cbl#L149) |
| RETURN-CODE | literal 0 | MOVE | 0000-MAIN | [ACCTMNT.cbl:151](../../../pangaea/programs/ACCTMNT.cbl#L151) |
| RETURN-CODE | literal 16 | MOVE | 1000-INITIALIZE | [ACCTMNT.cbl:160](../../../pangaea/programs/ACCTMNT.cbl#L160) |
| RETURN-CODE | literal 16 | MOVE | 9900-ABORT | [ACCTMNT.cbl:385](../../../pangaea/programs/ACCTMNT.cbl#L385) |

## Field origins

Where the values of each field the program writes to a file come from: the fields of the files it reads, literals, and initial values, followed through every statement that moves a value, and through groups and redefinitions by the bytes they share.

### File REJECT-OUT

Written from [REJECT-RECORD](#reject-record).

| Field | Position | Origins |
| --- | --: | --- |
| RJ-PROGRAM | 1-8 | literal 'ACCTMNT' |
| RJ-REASON | 9-12 | literal condition RJ-BAD-STATUS; literal condition RJ-BAD-TYPE; literal condition RJ-CLOSED; literal condition RJ-EXISTS; literal condition RJ-HAS-BALANCE; literal condition RJ-NO-ACCOUNT |
| RJ-TRAN | 13-92 | file MAINT-IN MAINT-REQUEST |
| RJ-REJECT-DATE | 101-110 | file CONTROL-FILE BC-BUSINESS-DATE; literal '-' |

### File ACCT-OUT

Written from [ACCT-MASTER-RECORD](#acct-master-record).

| Field | Position | Origins |
| --- | --: | --- |
| AM-ACCT-NO | 1-10 | file ACCT-IN AI-ACCT-NO; file MAINT-IN MR-ACCT-NO |
| AM-CUST-NO | 11-18 | file ACCT-IN ACCT-IN-RECORD 11-18; file MAINT-IN MR-CUST-NO |
| AM-SHORT-NAME | 19-38 | file ACCT-IN ACCT-IN-RECORD 19-38; file MAINT-IN MR-SHORT-NAME |
| AM-ACCT-TYPE | 39 | file ACCT-IN ACCT-IN-RECORD 39-39; file MAINT-IN MR-ACCT-TYPE |
| AM-STATUS | 40 | file ACCT-IN ACCT-IN-RECORD 40-40; literal condition AM-CLOSED; literal condition AM-FROZEN; literal condition AM-OPEN |
| AM-OPEN-DATE | 41-48 | file ACCT-IN ACCT-IN-RECORD 41-48; file MAINT-IN MR-EFFECTIVE-DATE; literal 19; literal 20 |
| AM-LAST-ACTIVITY-DATE | 49-56 | file ACCT-IN ACCT-IN-RECORD 49-56; file CONTROL-FILE BC-BUSINESS-DATE |
| AM-BALANCE | 57-69 | file ACCT-IN ACCT-IN-RECORD 57-69; literal 0 |
| AM-OVERDRAFT-LIMIT | 70-78 | file ACCT-IN ACCT-IN-RECORD 70-78; file MAINT-IN MR-OVERDRAFT-LIMIT; literal 0 |

### File REPORT-OUT

Written from [REPORT-LINE](#report-line), [WS-REPORT-COUNT-LINE](#ws-report-count-line) and [WS-REPORT-AMOUNT-LINE](#ws-report-amount-line).

#### Record REPORT-LINE

| Field | Position | Origins |
| --- | --: | --- |
| REPORT-LINE | 1-80 | file CONTROL-FILE BC-BUSINESS-DATE; literal 'ACCTMNT  ACCOUNT MAINTENANCE     BUSINESS DATE '; literal 'NEW MASTER IN BALANCE'; literal SPACES |

#### Record WS-REPORT-COUNT-LINE

| Field | Position | Origins |
| --- | --: | --- |
| WS-RC-LABEL | 1-30 | literal 'ACCOUNTS CLOSED'; literal 'ACCOUNTS FROZEN'; literal 'ACCOUNTS OPENED'; literal 'ACCOUNTS READ'; literal 'ACCOUNTS UNFROZEN'; literal 'ACCOUNTS WRITTEN'; literal 'OVERDRAFT LIMITS CHANGED'; literal 'REQUESTS READ'; literal 'REQUESTS REJECTED' |
| WS-RC-VALUE | 31-40 | initial value of WS-CLOSED-COUNT; initial value of WS-FROZEN-COUNT; initial value of WS-LIMIT-COUNT; initial value of WS-MASTERS-READ; initial value of WS-MASTERS-WRITTEN; initial value of WS-OPENED-COUNT; initial value of WS-REJECT-COUNT; initial value of WS-REQUEST-COUNT; initial value of WS-UNFROZEN-COUNT; literal 1 |

#### Record WS-REPORT-AMOUNT-LINE

| Field | Position | Origins |
| --- | --: | --- |
| WS-RA-LABEL | 1-30 | literal 'BALANCES' |
| WS-RA-VALUE | 31-51 | file ACCT-IN ACCT-IN-RECORD 57-69; initial value of WS-BALANCES-OUT; literal 0 |

## Links between records

Fields of one input the program relies on to find or match the records of another: moved into an indexed file's record key to read it by, compared for equality, or compared for order to merge the two in that order, each side's values from that field alone.

| Field | Of | Links to | Of | How | Source |
| --- | --- | --- | --- | --- | --- |
| MR-ACCT-NO | MAINT-IN | AI-ACCT-NO | ACCT-IN | merged in order with | [ACCTMNT.cbl:210](../../../pangaea/programs/ACCTMNT.cbl#L210) |

## Records

### BANK-CONTROL-RECORD

File CONTROL-FILE. 80 bytes, layout in copybook [BANKCTL](../copybooks/BANKCTL.md).

### ACCT-IN-RECORD

File ACCT-IN. 100 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | ACCT-IN-RECORD | group | 1 | 100 | | | | [ACCTMNT.cbl:66](../../../pangaea/programs/ACCTMNT.cbl#L66) |
| 05 | ├ AI-ACCT-NO | text | 1 | 10 | X(10) | | | [ACCTMNT.cbl:67](../../../pangaea/programs/ACCTMNT.cbl#L67) |
| 05 | └ FILLER | filler | 11 | 90 | X(90) | | | [ACCTMNT.cbl:68](../../../pangaea/programs/ACCTMNT.cbl#L68) |

### MAINT-REQUEST

File MAINT-IN. 80 bytes, layout in copybook [MAINTREC](../copybooks/MAINTREC.md).

### ACCT-OUT-RECORD

File ACCT-OUT. 100 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | ACCT-OUT-RECORD | text | 1 | 100 | X(100) | | | [ACCTMNT.cbl:72](../../../pangaea/programs/ACCTMNT.cbl#L72) |

### REJECT-OUT-RECORD

File REJECT-OUT. 110 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | REJECT-OUT-RECORD | text | 1 | 110 | X(110) | | | [ACCTMNT.cbl:74](../../../pangaea/programs/ACCTMNT.cbl#L74) |

### REPORT-LINE

File REPORT-OUT. 80 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | REPORT-LINE | text | 1 | 80 | X(80) | | | [ACCTMNT.cbl:76](../../../pangaea/programs/ACCTMNT.cbl#L76) |

### WS-FILE-STATUSES

Working storage. 12 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-FILE-STATUSES | group | 1 | 12 | | | | [ACCTMNT.cbl:79](../../../pangaea/programs/ACCTMNT.cbl#L79) |
| 05 | ├ WS-CONTROL-STATUS | text | 1 | 2 | XX | | | [ACCTMNT.cbl:80](../../../pangaea/programs/ACCTMNT.cbl#L80) |
| 05 | ├ WS-ACCT-IN-STATUS | text | 3 | 2 | XX | | | [ACCTMNT.cbl:81](../../../pangaea/programs/ACCTMNT.cbl#L81) |
| 05 | ├ WS-MAINT-IN-STATUS | text | 5 | 2 | XX | | | [ACCTMNT.cbl:82](../../../pangaea/programs/ACCTMNT.cbl#L82) |
| 05 | ├ WS-ACCT-OUT-STATUS | text | 7 | 2 | XX | | | [ACCTMNT.cbl:83](../../../pangaea/programs/ACCTMNT.cbl#L83) |
| 05 | ├ WS-REJECT-STATUS | text | 9 | 2 | XX | | | [ACCTMNT.cbl:84](../../../pangaea/programs/ACCTMNT.cbl#L84) |
| 05 | └ WS-REPORT-STATUS | text | 11 | 2 | XX | | | [ACCTMNT.cbl:85](../../../pangaea/programs/ACCTMNT.cbl#L85) |

### WS-BUSINESS-DATE

Working storage. 8 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-BUSINESS-DATE | date | 1 | 8 | 9(8) | | | [ACCTMNT.cbl:87](../../../pangaea/programs/ACCTMNT.cbl#L87) |

### WS-EFFECTIVE-DATE

Working storage. 8 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-EFFECTIVE-DATE | group | 1 | 8 | | | | [ACCTMNT.cbl:88](../../../pangaea/programs/ACCTMNT.cbl#L88) |
| 05 | ├ WS-EFF-CENTURY | number | 1 | 2 | 99 | | | [ACCTMNT.cbl:89](../../../pangaea/programs/ACCTMNT.cbl#L89) |
| 05 | └ WS-EFF-YYMMDD | group | 3 | 6 | | | | [ACCTMNT.cbl:90](../../../pangaea/programs/ACCTMNT.cbl#L90) |
| 10 |   ├ WS-EFF-YY | number | 3 | 2 | 99 | | | [ACCTMNT.cbl:91](../../../pangaea/programs/ACCTMNT.cbl#L91) |
| 10 |   └ WS-EFF-MMDD | number | 5 | 4 | 9(4) | | | [ACCTMNT.cbl:92](../../../pangaea/programs/ACCTMNT.cbl#L92) |

### WS-ABORT-MESSAGE

Working storage. 40 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-ABORT-MESSAGE | text | 1 | 40 | X(40) | | | [ACCTMNT.cbl:93](../../../pangaea/programs/ACCTMNT.cbl#L93) |

### WS-MASTER-KEY

Working storage. 10 bytes. MATCH KEYS: HIGH-VALUES ONCE A FILE IS AT ITS END. THE ACTIVE KEY IS THE LOWER OF THE TWO: THE ACCOUNT BEING WORKED ON.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-MASTER-KEY | text | 1 | 10 | X(10) | | | [ACCTMNT.cbl:97](../../../pangaea/programs/ACCTMNT.cbl#L97) |

### WS-REQUEST-KEY

Working storage. 10 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REQUEST-KEY | text | 1 | 10 | X(10) | | | [ACCTMNT.cbl:98](../../../pangaea/programs/ACCTMNT.cbl#L98) |

### WS-ACTIVE-KEY

Working storage. 10 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-ACTIVE-KEY | text | 1 | 10 | X(10) | | | [ACCTMNT.cbl:99](../../../pangaea/programs/ACCTMNT.cbl#L99) |

### WS-PREV-MASTER-KEY

Working storage. 10 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-PREV-MASTER-KEY | text | 1 | 10 | X(10) | | VALUE LOW-VALUES | [ACCTMNT.cbl:100](../../../pangaea/programs/ACCTMNT.cbl#L100) |

### WS-REQUEST-ORDER-KEY

Working storage. 20 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REQUEST-ORDER-KEY | group | 1 | 20 | | | | [ACCTMNT.cbl:101](../../../pangaea/programs/ACCTMNT.cbl#L101) |
| 05 | ├ WS-ROK-ACCT-NO | text | 1 | 10 | X(10) | | | [ACCTMNT.cbl:102](../../../pangaea/programs/ACCTMNT.cbl#L102) |
| 05 | └ WS-ROK-REQUEST-ID | text | 11 | 10 | X(10) | | | [ACCTMNT.cbl:103](../../../pangaea/programs/ACCTMNT.cbl#L103) |

### WS-PREV-REQUEST-ORDER-KEY

Working storage. 20 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-PREV-REQUEST-ORDER-KEY | text | 1 | 20 | X(20) | | VALUE LOW-VALUES | [ACCTMNT.cbl:104](../../../pangaea/programs/ACCTMNT.cbl#L104) |

### WS-HELD-SW

Working storage. 1 bytes. THE ACCOUNT HELD WHILE ITS REQUESTS APPLY.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-HELD-SW | code | 1 | 1 | X | | | [ACCTMNT.cbl:107](../../../pangaea/programs/ACCTMNT.cbl#L107) |
| 88 | ├ ACCOUNT-HELD | | | | | | VALUE 'Y' | [ACCTMNT.cbl:108](../../../pangaea/programs/ACCTMNT.cbl#L108) |
| 88 | └ NOTHING-HELD | | | | | | VALUE 'N' | [ACCTMNT.cbl:109](../../../pangaea/programs/ACCTMNT.cbl#L109) |

### ACCT-MASTER-RECORD

Working storage. 100 bytes, layout in copybook [ACCTMAST](../copybooks/ACCTMAST.md).

### REJECT-RECORD

Working storage. 110 bytes, layout in copybook [REJREC](../copybooks/REJREC.md).

### WS-TOTALS

Working storage. 102 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TOTALS | group | 1 | 102 | | | | [ACCTMNT.cbl:114](../../../pangaea/programs/ACCTMNT.cbl#L114) |
| 05 | ├ WS-MASTERS-READ | number | 1 | 8 | 9(8) | | VALUE 0 | [ACCTMNT.cbl:115](../../../pangaea/programs/ACCTMNT.cbl#L115) |
| 05 | ├ WS-MASTERS-WRITTEN | number | 9 | 8 | 9(8) | | VALUE 0 | [ACCTMNT.cbl:116](../../../pangaea/programs/ACCTMNT.cbl#L116) |
| 05 | ├ WS-REQUEST-COUNT | number | 17 | 8 | 9(8) | | VALUE 0 | [ACCTMNT.cbl:117](../../../pangaea/programs/ACCTMNT.cbl#L117) |
| 05 | ├ WS-OPENED-COUNT | number | 25 | 8 | 9(8) | | VALUE 0 | [ACCTMNT.cbl:118](../../../pangaea/programs/ACCTMNT.cbl#L118) |
| 05 | ├ WS-CLOSED-COUNT | number | 33 | 8 | 9(8) | | VALUE 0 | [ACCTMNT.cbl:119](../../../pangaea/programs/ACCTMNT.cbl#L119) |
| 05 | ├ WS-FROZEN-COUNT | number | 41 | 8 | 9(8) | | VALUE 0 | [ACCTMNT.cbl:120](../../../pangaea/programs/ACCTMNT.cbl#L120) |
| 05 | ├ WS-UNFROZEN-COUNT | number | 49 | 8 | 9(8) | | VALUE 0 | [ACCTMNT.cbl:121](../../../pangaea/programs/ACCTMNT.cbl#L121) |
| 05 | ├ WS-LIMIT-COUNT | number | 57 | 8 | 9(8) | | VALUE 0 | [ACCTMNT.cbl:122](../../../pangaea/programs/ACCTMNT.cbl#L122) |
| 05 | ├ WS-REJECT-COUNT | number | 65 | 8 | 9(8) | | VALUE 0 | [ACCTMNT.cbl:123](../../../pangaea/programs/ACCTMNT.cbl#L123) |
| 05 | ├ WS-BALANCES-IN | amount | 73 | 15 | S9(13)V99 | | VALUE 0 | [ACCTMNT.cbl:124](../../../pangaea/programs/ACCTMNT.cbl#L124) |
| 05 | └ WS-BALANCES-OUT | amount | 88 | 15 | S9(13)V99 | | VALUE 0 | [ACCTMNT.cbl:125](../../../pangaea/programs/ACCTMNT.cbl#L125) |

### WS-REPORT-COUNT-LINE

Working storage. 80 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-COUNT-LINE | group | 1 | 80 | | | | [ACCTMNT.cbl:127](../../../pangaea/programs/ACCTMNT.cbl#L127) |
| 05 | ├ WS-RC-LABEL | text | 1 | 30 | X(30) | | | [ACCTMNT.cbl:128](../../../pangaea/programs/ACCTMNT.cbl#L128) |
| 05 | ├ WS-RC-VALUE | edited | 31 | 10 | ZZ,ZZZ,ZZ9 | | | [ACCTMNT.cbl:129](../../../pangaea/programs/ACCTMNT.cbl#L129) |
| 05 | └ FILLER | filler | 41 | 40 | X(40) | | VALUE SPACES | [ACCTMNT.cbl:130](../../../pangaea/programs/ACCTMNT.cbl#L130) |

### WS-REPORT-AMOUNT-LINE

Working storage. 80 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-AMOUNT-LINE | group | 1 | 80 | | | | [ACCTMNT.cbl:131](../../../pangaea/programs/ACCTMNT.cbl#L131) |
| 05 | ├ WS-RA-LABEL | text | 1 | 30 | X(30) | | | [ACCTMNT.cbl:132](../../../pangaea/programs/ACCTMNT.cbl#L132) |
| 05 | ├ WS-RA-VALUE | edited | 31 | 21 | Z,ZZZ,ZZZ,ZZZ,ZZ9.99- | | | [ACCTMNT.cbl:133](../../../pangaea/programs/ACCTMNT.cbl#L133) |
| 05 | └ FILLER | filler | 52 | 29 | X(29) | | VALUE SPACES | [ACCTMNT.cbl:134](../../../pangaea/programs/ACCTMNT.cbl#L134) |
