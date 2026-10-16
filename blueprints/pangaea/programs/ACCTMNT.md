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
1994-11-14 PK  CR-0201 OVERDRAFT LIMIT CHANGES.
1999-01-11 SL  Y2K: 8-DIGIT OPEN DATES.
```

Program [ACCTMNT.cbl](../../../pangaea/programs/ACCTMNT.cbl): 6 files, 21 paragraphs.

- AUTHOR: J. ARKWRIGHT, BATCH SYSTEMS.
- INSTALLATION: FIRST BANK OF PANGAEA, DATA CENTER.
- DATE-WRITTEN: 1988-05-09.

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

Solid arrows are PERFORMs, dashed ones GO TOs, and dotted ones control running on into the next paragraph inside a PERFORM THRU range.

## Paragraphs

| Paragraph | Description | Source | Performs | Goes to | Performed by | Reached by |
| --- | --- | --- | --- | --- | --- | --- |
| 0000-MAIN | MATCH UNTIL BOTH FILES END, THEN RECONCILE AND REPORT. | [ACCTMNT.cbl:128](../../../pangaea/programs/ACCTMNT.cbl#L128) | 1000-INITIALIZE, 2000-NEXT-ACCOUNT THRU 2999-EXIT, 9900-ABORT, 8000-WRITE-REPORT, 9000-CLOSE-FILES | | | |
| 1000-INITIALIZE | OPEN THE FILES AND READ THE FIRST RECORD OF EACH. | [ACCTMNT.cbl:145](../../../pangaea/programs/ACCTMNT.cbl#L145) | 9900-ABORT, 1200-READ-MASTER THRU 1200-EXIT, 1300-READ-REQUEST THRU 1300-EXIT | | 0000-MAIN | |
| 1200-READ-MASTER | READ THE NEXT OLD MASTER RECORD, CHECKING ACCOUNT ORDER. | [ACCTMNT.cbl:162](../../../pangaea/programs/ACCTMNT.cbl#L162) | 9900-ABORT | 1200-EXIT, 1200-EXIT (falls into) | 1000-INITIALIZE, 2000-NEXT-ACCOUNT | |
| 1200-EXIT | | [ACCTMNT.cbl:173](../../../pangaea/programs/ACCTMNT.cbl#L173) | | | 1000-INITIALIZE, 2000-NEXT-ACCOUNT | 1200-READ-MASTER |
| 1300-READ-REQUEST | READ THE NEXT REQUEST, CHECKING ITS ORDER. | [ACCTMNT.cbl:177](../../../pangaea/programs/ACCTMNT.cbl#L177) | 9900-ABORT | 1300-EXIT, 1300-EXIT (falls into) | 1000-INITIALIZE, 2100-APPLY-REQUESTS | |
| 1300-EXIT | | [ACCTMNT.cbl:189](../../../pangaea/programs/ACCTMNT.cbl#L189) | | | 1000-INITIALIZE, 2100-APPLY-REQUESTS | 1300-READ-REQUEST |
| 2000-NEXT-ACCOUNT | TAKE THE LOWER KEY AS THE ACTIVE ACCOUNT, AND HOLD IT IF THE MASTER HAS IT. AN ACCOUNT ONLY THE REQUESTS NAME IS NOT HELD UNTIL AN OPENING MAKES IT. | [ACCTMNT.cbl:195](../../../pangaea/programs/ACCTMNT.cbl#L195) | 1200-READ-MASTER THRU 1200-EXIT | 2100-APPLY-REQUESTS (falls into) | 0000-MAIN | |
| 2100-APPLY-REQUESTS | APPLY THE ACTIVE ACCOUNT'S REQUESTS, ONE AFTER THE OTHER. | [ACCTMNT.cbl:207](../../../pangaea/programs/ACCTMNT.cbl#L207) | 2200-APPLY-REQUEST THRU 2299-EXIT, 1300-READ-REQUEST THRU 1300-EXIT | 2800-WRITE-HELD, 2100-APPLY-REQUESTS | 0000-MAIN | 2000-NEXT-ACCOUNT, 2100-APPLY-REQUESTS |
| 2200-APPLY-REQUEST | APPLY ONE REQUEST BY ITS TYPE. | [ACCTMNT.cbl:215](../../../pangaea/programs/ACCTMNT.cbl#L215) | | 2210-OPEN, 2290-REJECT, 2220-CLOSE, 2230-FREEZE, 2240-UNFREEZE, 2250-LIMIT | 0000-MAIN, 2100-APPLY-REQUESTS | |
| 2210-OPEN | OPEN A NEW ACCOUNT, EMPTY, FROM TODAY. | [ACCTMNT.cbl:236](../../../pangaea/programs/ACCTMNT.cbl#L236) | | 2290-REJECT, 2299-EXIT | 0000-MAIN, 2100-APPLY-REQUESTS | 2200-APPLY-REQUEST |
| 2220-CLOSE | CLOSE AN ACCOUNT, IF IT IS EMPTY. | [ACCTMNT.cbl:256](../../../pangaea/programs/ACCTMNT.cbl#L256) | | 2290-REJECT, 2299-EXIT | 0000-MAIN, 2100-APPLY-REQUESTS | 2200-APPLY-REQUEST |
| 2230-FREEZE | FREEZE AN OPEN ACCOUNT. | [ACCTMNT.cbl:266](../../../pangaea/programs/ACCTMNT.cbl#L266) | | 2290-REJECT, 2299-EXIT | 0000-MAIN, 2100-APPLY-REQUESTS | 2200-APPLY-REQUEST |
| 2240-UNFREEZE | UNFREEZE A FROZEN ACCOUNT. | [ACCTMNT.cbl:275](../../../pangaea/programs/ACCTMNT.cbl#L275) | | 2290-REJECT, 2299-EXIT | 0000-MAIN, 2100-APPLY-REQUESTS | 2200-APPLY-REQUEST |
| 2250-LIMIT | CHANGE A CURRENT ACCOUNT'S OVERDRAFT LIMIT. | [ACCTMNT.cbl:284](../../../pangaea/programs/ACCTMNT.cbl#L284) | | 2290-REJECT, 2299-EXIT | 0000-MAIN, 2100-APPLY-REQUESTS | 2200-APPLY-REQUEST |
| 2290-REJECT | WRITE THE REQUEST TO THE REJECTS WITH ITS REASON. | [ACCTMNT.cbl:293](../../../pangaea/programs/ACCTMNT.cbl#L293) | | 2299-EXIT (falls into) | 0000-MAIN, 2100-APPLY-REQUESTS | 2200-APPLY-REQUEST, 2210-OPEN, 2220-CLOSE, 2230-FREEZE, 2240-UNFREEZE, 2250-LIMIT |
| 2299-EXIT | | [ACCTMNT.cbl:298](../../../pangaea/programs/ACCTMNT.cbl#L298) | | | 0000-MAIN, 2100-APPLY-REQUESTS | 2210-OPEN, 2220-CLOSE, 2230-FREEZE, 2240-UNFREEZE, 2250-LIMIT, 2290-REJECT |
| 2800-WRITE-HELD | WRITE THE HELD ACCOUNT, ITS REQUESTS APPLIED, TO THE NEW MASTER. | [ACCTMNT.cbl:302](../../../pangaea/programs/ACCTMNT.cbl#L302) | | 2999-EXIT (falls into) | 0000-MAIN | 2100-APPLY-REQUESTS |
| 2999-EXIT | | [ACCTMNT.cbl:307](../../../pangaea/programs/ACCTMNT.cbl#L307) | | | 0000-MAIN | 2800-WRITE-HELD |
| 8000-WRITE-REPORT | WRITE THE MAINTENANCE REPORT. | [ACCTMNT.cbl:311](../../../pangaea/programs/ACCTMNT.cbl#L311) | | | 0000-MAIN | |
| 9000-CLOSE-FILES | CLOSE ALL FILES. | [ACCTMNT.cbl:357](../../../pangaea/programs/ACCTMNT.cbl#L357) | | | 0000-MAIN, 9900-ABORT | |
| 9900-ABORT | STOP THE RUN WITH RETURN CODE 16, SAYING WHY. | [ACCTMNT.cbl:362](../../../pangaea/programs/ACCTMNT.cbl#L362) | 9000-CLOSE-FILES | | 0000-MAIN, 1000-INITIALIZE, 1200-READ-MASTER, 1300-READ-REQUEST | |

## Records

### BANK-CONTROL-RECORD

File CONTROL-FILE. 80 bytes, layout in copybook [BANKCTL](../copybooks/BANKCTL.md).

### ACCT-IN-RECORD

File ACCT-IN. 100 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | ACCT-IN-RECORD | 1 | 100 | | | | [ACCTMNT.cbl:61](../../../pangaea/programs/ACCTMNT.cbl#L61) |
| 05 | AI-ACCT-NO | 1 | 10 | X(10) | | | [ACCTMNT.cbl:62](../../../pangaea/programs/ACCTMNT.cbl#L62) |
| 05 | FILLER | 11 | 90 | X(90) | | | [ACCTMNT.cbl:63](../../../pangaea/programs/ACCTMNT.cbl#L63) |

### MAINT-REQUEST

File MAINT-IN. 80 bytes, layout in copybook [MAINTREC](../copybooks/MAINTREC.md).

### ACCT-OUT-RECORD

File ACCT-OUT. 100 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | ACCT-OUT-RECORD | 1 | 100 | X(100) | | | [ACCTMNT.cbl:67](../../../pangaea/programs/ACCTMNT.cbl#L67) |

### REJECT-OUT-RECORD

File REJECT-OUT. 100 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | REJECT-OUT-RECORD | 1 | 100 | X(100) | | | [ACCTMNT.cbl:69](../../../pangaea/programs/ACCTMNT.cbl#L69) |

### REPORT-LINE

File REPORT-OUT. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | REPORT-LINE | 1 | 80 | X(80) | | | [ACCTMNT.cbl:71](../../../pangaea/programs/ACCTMNT.cbl#L71) |

### WS-FILE-STATUSES

Working storage. 12 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-FILE-STATUSES | 1 | 12 | | | | [ACCTMNT.cbl:74](../../../pangaea/programs/ACCTMNT.cbl#L74) |
| 05 | WS-CONTROL-STATUS | 1 | 2 | XX | | | [ACCTMNT.cbl:75](../../../pangaea/programs/ACCTMNT.cbl#L75) |
| 05 | WS-ACCT-IN-STATUS | 3 | 2 | XX | | | [ACCTMNT.cbl:76](../../../pangaea/programs/ACCTMNT.cbl#L76) |
| 05 | WS-MAINT-IN-STATUS | 5 | 2 | XX | | | [ACCTMNT.cbl:77](../../../pangaea/programs/ACCTMNT.cbl#L77) |
| 05 | WS-ACCT-OUT-STATUS | 7 | 2 | XX | | | [ACCTMNT.cbl:78](../../../pangaea/programs/ACCTMNT.cbl#L78) |
| 05 | WS-REJECT-STATUS | 9 | 2 | XX | | | [ACCTMNT.cbl:79](../../../pangaea/programs/ACCTMNT.cbl#L79) |
| 05 | WS-REPORT-STATUS | 11 | 2 | XX | | | [ACCTMNT.cbl:80](../../../pangaea/programs/ACCTMNT.cbl#L80) |

### WS-BUSINESS-DATE

Working storage. 8 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-BUSINESS-DATE | 1 | 8 | 9(8) | | | [ACCTMNT.cbl:82](../../../pangaea/programs/ACCTMNT.cbl#L82) |

### WS-ABORT-MESSAGE

Working storage. 40 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-ABORT-MESSAGE | 1 | 40 | X(40) | | | [ACCTMNT.cbl:83](../../../pangaea/programs/ACCTMNT.cbl#L83) |

### WS-MASTER-KEY

Working storage. 10 bytes. MATCH KEYS: HIGH-VALUES ONCE A FILE IS AT ITS END. THE ACTIVE KEY IS THE LOWER OF THE TWO: THE ACCOUNT BEING WORKED ON.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-MASTER-KEY | 1 | 10 | X(10) | | | [ACCTMNT.cbl:87](../../../pangaea/programs/ACCTMNT.cbl#L87) |

### WS-REQUEST-KEY

Working storage. 10 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REQUEST-KEY | 1 | 10 | X(10) | | | [ACCTMNT.cbl:88](../../../pangaea/programs/ACCTMNT.cbl#L88) |

### WS-ACTIVE-KEY

Working storage. 10 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-ACTIVE-KEY | 1 | 10 | X(10) | | | [ACCTMNT.cbl:89](../../../pangaea/programs/ACCTMNT.cbl#L89) |

### WS-PREV-MASTER-KEY

Working storage. 10 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-PREV-MASTER-KEY | 1 | 10 | X(10) | | VALUE LOW-VALUES | [ACCTMNT.cbl:90](../../../pangaea/programs/ACCTMNT.cbl#L90) |

### WS-REQUEST-ORDER-KEY

Working storage. 20 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REQUEST-ORDER-KEY | 1 | 20 | | | | [ACCTMNT.cbl:91](../../../pangaea/programs/ACCTMNT.cbl#L91) |
| 05 | WS-ROK-ACCT-NO | 1 | 10 | X(10) | | | [ACCTMNT.cbl:92](../../../pangaea/programs/ACCTMNT.cbl#L92) |
| 05 | WS-ROK-REQUEST-ID | 11 | 10 | X(10) | | | [ACCTMNT.cbl:93](../../../pangaea/programs/ACCTMNT.cbl#L93) |

### WS-PREV-REQUEST-ORDER-KEY

Working storage. 20 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-PREV-REQUEST-ORDER-KEY | 1 | 20 | X(20) | | VALUE LOW-VALUES | [ACCTMNT.cbl:94](../../../pangaea/programs/ACCTMNT.cbl#L94) |

### WS-HELD-SW

Working storage. 1 bytes. THE ACCOUNT HELD WHILE ITS REQUESTS APPLY.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-HELD-SW | 1 | 1 | X | | | [ACCTMNT.cbl:97](../../../pangaea/programs/ACCTMNT.cbl#L97) |
| 88 | ACCOUNT-HELD | | | | | VALUE 'Y' | [ACCTMNT.cbl:98](../../../pangaea/programs/ACCTMNT.cbl#L98) |
| 88 | NOTHING-HELD | | | | | VALUE 'N' | [ACCTMNT.cbl:99](../../../pangaea/programs/ACCTMNT.cbl#L99) |

### ACCT-MASTER-RECORD

Working storage. 100 bytes, layout in copybook [ACCTMAST](../copybooks/ACCTMAST.md).

### REJECT-RECORD

Working storage. 100 bytes, layout in copybook [REJREC](../copybooks/REJREC.md).

### WS-TOTALS

Working storage. 102 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TOTALS | 1 | 102 | | | | [ACCTMNT.cbl:104](../../../pangaea/programs/ACCTMNT.cbl#L104) |
| 05 | WS-MASTERS-READ | 1 | 8 | 9(8) | | VALUE 0 | [ACCTMNT.cbl:105](../../../pangaea/programs/ACCTMNT.cbl#L105) |
| 05 | WS-MASTERS-WRITTEN | 9 | 8 | 9(8) | | VALUE 0 | [ACCTMNT.cbl:106](../../../pangaea/programs/ACCTMNT.cbl#L106) |
| 05 | WS-REQUEST-COUNT | 17 | 8 | 9(8) | | VALUE 0 | [ACCTMNT.cbl:107](../../../pangaea/programs/ACCTMNT.cbl#L107) |
| 05 | WS-OPENED-COUNT | 25 | 8 | 9(8) | | VALUE 0 | [ACCTMNT.cbl:108](../../../pangaea/programs/ACCTMNT.cbl#L108) |
| 05 | WS-CLOSED-COUNT | 33 | 8 | 9(8) | | VALUE 0 | [ACCTMNT.cbl:109](../../../pangaea/programs/ACCTMNT.cbl#L109) |
| 05 | WS-FROZEN-COUNT | 41 | 8 | 9(8) | | VALUE 0 | [ACCTMNT.cbl:110](../../../pangaea/programs/ACCTMNT.cbl#L110) |
| 05 | WS-UNFROZEN-COUNT | 49 | 8 | 9(8) | | VALUE 0 | [ACCTMNT.cbl:111](../../../pangaea/programs/ACCTMNT.cbl#L111) |
| 05 | WS-LIMIT-COUNT | 57 | 8 | 9(8) | | VALUE 0 | [ACCTMNT.cbl:112](../../../pangaea/programs/ACCTMNT.cbl#L112) |
| 05 | WS-REJECT-COUNT | 65 | 8 | 9(8) | | VALUE 0 | [ACCTMNT.cbl:113](../../../pangaea/programs/ACCTMNT.cbl#L113) |
| 05 | WS-BALANCES-IN | 73 | 15 | S9(13)V99 | | VALUE 0 | [ACCTMNT.cbl:114](../../../pangaea/programs/ACCTMNT.cbl#L114) |
| 05 | WS-BALANCES-OUT | 88 | 15 | S9(13)V99 | | VALUE 0 | [ACCTMNT.cbl:115](../../../pangaea/programs/ACCTMNT.cbl#L115) |

### WS-REPORT-COUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-COUNT-LINE | 1 | 80 | | | | [ACCTMNT.cbl:117](../../../pangaea/programs/ACCTMNT.cbl#L117) |
| 05 | WS-RC-LABEL | 1 | 30 | X(30) | | | [ACCTMNT.cbl:118](../../../pangaea/programs/ACCTMNT.cbl#L118) |
| 05 | WS-RC-VALUE | 31 | 10 | ZZ,ZZZ,ZZ9 | | | [ACCTMNT.cbl:119](../../../pangaea/programs/ACCTMNT.cbl#L119) |
| 05 | FILLER | 41 | 40 | X(40) | | VALUE SPACES | [ACCTMNT.cbl:120](../../../pangaea/programs/ACCTMNT.cbl#L120) |

### WS-REPORT-AMOUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-AMOUNT-LINE | 1 | 80 | | | | [ACCTMNT.cbl:121](../../../pangaea/programs/ACCTMNT.cbl#L121) |
| 05 | WS-RA-LABEL | 1 | 30 | X(30) | | | [ACCTMNT.cbl:122](../../../pangaea/programs/ACCTMNT.cbl#L122) |
| 05 | WS-RA-VALUE | 31 | 21 | Z,ZZZ,ZZZ,ZZZ,ZZ9.99- | | | [ACCTMNT.cbl:123](../../../pangaea/programs/ACCTMNT.cbl#L123) |
| 05 | FILLER | 52 | 29 | X(29) | | VALUE SPACES | [ACCTMNT.cbl:124](../../../pangaea/programs/ACCTMNT.cbl#L124) |
