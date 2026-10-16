# TXPOST

```text
TXPOST: POST THE DAY'S TRANSACTIONS TO THE ACCOUNT MASTER

MATCHES THE SORTED TRANSACTIONS AGAINST THE OLD MASTER, BOTH
IN ACCOUNT NUMBER ORDER, AND WRITES THE NEW MASTER: EVERY OLD
RECORD, WITH ITS TRANSACTIONS APPLIED. EACH POSTING GOES TO
THE AUDIT TRAIL; A TRANSACTION THAT CANNOT BE POSTED GOES TO
THE REJECTS.

RULES: NOTHING POSTS TO A CLOSED ACCOUNT. DEPOSITS AND FEES
POST TO A FROZEN ACCOUNT, WITHDRAWALS DO NOT. A WITHDRAWAL
MAY TAKE THE BALANCE DOWN TO THE OVERDRAFT LIMIT, NO FURTHER.
A FEE ALWAYS POSTS, EVEN PAST THE LIMIT.

RETURN CODES: 0 ALL POSTED, 4 SOME REJECTED,
              16 INPUT OUT OF ORDER, OUT OF BALANCE OR I/O.

MODIFICATION HISTORY
DATE       BY  CHANGE
1987-04-02 JA  ORIGINAL VERSION.
1994-11-14 PK  CR-0201 OVERDRAFT LIMITS.
1999-01-11 SL  Y2K: 8-DIGIT DATES IN THE AUDIT TRAIL.
2003-05-26 MV  CR-0352 FEES POST EVEN PAST THE OVERDRAFT
               LIMIT.
2011-08-15 MV  CR-0471 RECONCILE THE NEW MASTER BEFORE
               ENDING THE RUN.
```

Program [TXPOST.cbl](../../../pangaea/programs/TXPOST.cbl): 7 files, 16 paragraphs.

- AUTHOR: J. ARKWRIGHT, BATCH SYSTEMS.
- INSTALLATION: FIRST BANK OF PANGAEA, DATA CENTRE.
- DATE-WRITTEN: 1987-04-02.

## Files

| File | Kind | Assigned to | Organization | Status field | Records |
| --- | --- | --- | --- | --- | --- |
| CONTROL-FILE | FD | BANKCTL | SEQUENTIAL | WS-CONTROL-STATUS | BANK-CONTROL-RECORD |
| ACCT-IN | FD | ACCTIN | SEQUENTIAL | WS-ACCT-IN-STATUS | ACCT-MASTER-RECORD |
| TRAN-IN | FD | TRANIN | SEQUENTIAL | WS-TRAN-IN-STATUS | TRAN-RECORD |
| ACCT-OUT | FD | ACCTOUT | SEQUENTIAL | WS-ACCT-OUT-STATUS | ACCT-OUT-RECORD |
| AUDIT-OUT | FD | AUDIT | SEQUENTIAL | WS-AUDIT-STATUS | AUDIT-OUT-RECORD |
| REJECT-OUT | FD | REJECTS | SEQUENTIAL | WS-REJECT-STATUS | REJECT-OUT-RECORD |
| REPORT-OUT | FD | POSTRPT | LINE SEQUENTIAL | WS-REPORT-STATUS | REPORT-LINE |

## Flow

![The paragraphs of TXPOST and the PERFORMs between them](TXPOST.svg)

## Paragraphs

| Paragraph | Description | Source | Performs | Performed by |
| --- | --- | --- | --- | --- |
| 0000-MAIN | MATCH UNTIL BOTH FILES END, THEN RECONCILE AND REPORT. | [TXPOST.cbl:131](../../../pangaea/programs/TXPOST.cbl#L131) | 1000-INITIALIZE, 2000-MATCH, 5000-RECONCILE, 8000-WRITE-REPORT, 9000-CLOSE-FILES | |
| 1000-INITIALIZE | OPEN THE FILES AND READ THE FIRST RECORD OF EACH. | [TXPOST.cbl:147](../../../pangaea/programs/TXPOST.cbl#L147) | 9900-ABORT, 1200-READ-MASTER, 1300-READ-TRAN | 0000-MAIN |
| 1200-READ-MASTER | READ THE NEXT OLD MASTER RECORD, CHECKING ACCOUNT ORDER. | [TXPOST.cbl:172](../../../pangaea/programs/TXPOST.cbl#L172) | 9900-ABORT | 1000-INITIALIZE, 2000-MATCH |
| 1300-READ-TRAN | READ THE NEXT TRANSACTION, CHECKING ITS ORDER. | [TXPOST.cbl:189](../../../pangaea/programs/TXPOST.cbl#L189) | 9900-ABORT | 1000-INITIALIZE, 2000-MATCH |
| 2000-MATCH | ONE STEP OF THE MATCH, BY WHICHEVER KEY IS LOWER: REJECT A TRANSACTION FOR A MISSING ACCOUNT, APPLY ONE TO ITS ACCOUNT, OR WRITE A FINISHED ACCOUNT. | [TXPOST.cbl:207](../../../pangaea/programs/TXPOST.cbl#L207) | 3900-REJECT-TRAN, 1300-READ-TRAN, 3000-APPLY-TRAN, 4000-WRITE-MASTER, 1200-READ-MASTER | 0000-MAIN |
| 3000-APPLY-TRAN | APPLY ONE TRANSACTION TO ITS ACCOUNT BY THE POSTING RULES. | [TXPOST.cbl:222](../../../pangaea/programs/TXPOST.cbl#L222) | 3900-REJECT-TRAN, 3100-POST-DEPOSIT, 3300-POST-FEE, 3200-POST-WITHDRAWAL | 2000-MATCH |
| 3100-POST-DEPOSIT | CREDIT A DEPOSIT. | [TXPOST.cbl:240](../../../pangaea/programs/TXPOST.cbl#L240) | 3800-WRITE-AUDIT | 3000-APPLY-TRAN |
| 3200-POST-WITHDRAWAL | DEBIT A WITHDRAWAL IF IT STAYS WITHIN THE OVERDRAFT LIMIT. | [TXPOST.cbl:247](../../../pangaea/programs/TXPOST.cbl#L247) | 3900-REJECT-TRAN, 3800-WRITE-AUDIT | 3000-APPLY-TRAN |
| 3300-POST-FEE | DEBIT A FEE, WHATEVER THE BALANCE. | [TXPOST.cbl:261](../../../pangaea/programs/TXPOST.cbl#L261) | 3800-WRITE-AUDIT | 3000-APPLY-TRAN |
| 3800-WRITE-AUDIT | RECORD A POSTING IN THE AUDIT TRAIL. | [TXPOST.cbl:268](../../../pangaea/programs/TXPOST.cbl#L268) | | 3100-POST-DEPOSIT, 3200-POST-WITHDRAWAL, 3300-POST-FEE |
| 3900-REJECT-TRAN | WRITE THE TRANSACTION TO THE REJECTS WITH ITS REASON. | [TXPOST.cbl:281](../../../pangaea/programs/TXPOST.cbl#L281) | | 2000-MATCH, 3000-APPLY-TRAN, 3200-POST-WITHDRAWAL |
| 4000-WRITE-MASTER | WRITE AN ACCOUNT, ITS TRANSACTIONS APPLIED, TO THE NEW MASTER. | [TXPOST.cbl:288](../../../pangaea/programs/TXPOST.cbl#L288) | | 2000-MATCH |
| 5000-RECONCILE | THE NEW MASTER MUST HOLD THE OLD BALANCES PLUS WHAT WAS POSTED. | [TXPOST.cbl:294](../../../pangaea/programs/TXPOST.cbl#L294) | 9900-ABORT | 0000-MAIN |
| 8000-WRITE-REPORT | WRITE THE POSTING REPORT. | [TXPOST.cbl:306](../../../pangaea/programs/TXPOST.cbl#L306) | | 0000-MAIN |
| 9000-CLOSE-FILES | CLOSE ALL FILES. | [TXPOST.cbl:349](../../../pangaea/programs/TXPOST.cbl#L349) | | 0000-MAIN, 9900-ABORT |
| 9900-ABORT | STOP THE RUN WITH RETURN CODE 16, SAYING WHY. | [TXPOST.cbl:354](../../../pangaea/programs/TXPOST.cbl#L354) | 9000-CLOSE-FILES | 1000-INITIALIZE, 1200-READ-MASTER, 1300-READ-TRAN, 5000-RECONCILE |

## Records

### BANK-CONTROL-RECORD

File CONTROL-FILE. 80 bytes, layout in copybook [BANKCTL](../copybooks/BANKCTL.md).

### ACCT-MASTER-RECORD

File ACCT-IN. 100 bytes, layout in copybook [ACCTMAST](../copybooks/ACCTMAST.md).

### TRAN-RECORD

File TRAN-IN. 80 bytes, layout in copybook [TRANREC](../copybooks/TRANREC.md).

### ACCT-OUT-RECORD

File ACCT-OUT. 100 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | ACCT-OUT-RECORD | 1 | 100 | X(100) | | | [TXPOST.cbl:68](../../../pangaea/programs/TXPOST.cbl#L68) |

### AUDIT-OUT-RECORD

File AUDIT-OUT. 100 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | AUDIT-OUT-RECORD | 1 | 100 | X(100) | | | [TXPOST.cbl:70](../../../pangaea/programs/TXPOST.cbl#L70) |

### REJECT-OUT-RECORD

File REJECT-OUT. 100 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | REJECT-OUT-RECORD | 1 | 100 | X(100) | | | [TXPOST.cbl:72](../../../pangaea/programs/TXPOST.cbl#L72) |

### REPORT-LINE

File REPORT-OUT. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | REPORT-LINE | 1 | 80 | X(80) | | | [TXPOST.cbl:74](../../../pangaea/programs/TXPOST.cbl#L74) |

### WS-FILE-STATUSES

Working storage. 14 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-FILE-STATUSES | 1 | 14 | | | | [TXPOST.cbl:77](../../../pangaea/programs/TXPOST.cbl#L77) |
| 05 | WS-CONTROL-STATUS | 1 | 2 | XX | | | [TXPOST.cbl:78](../../../pangaea/programs/TXPOST.cbl#L78) |
| 05 | WS-ACCT-IN-STATUS | 3 | 2 | XX | | | [TXPOST.cbl:79](../../../pangaea/programs/TXPOST.cbl#L79) |
| 05 | WS-TRAN-IN-STATUS | 5 | 2 | XX | | | [TXPOST.cbl:80](../../../pangaea/programs/TXPOST.cbl#L80) |
| 05 | WS-ACCT-OUT-STATUS | 7 | 2 | XX | | | [TXPOST.cbl:81](../../../pangaea/programs/TXPOST.cbl#L81) |
| 05 | WS-AUDIT-STATUS | 9 | 2 | XX | | | [TXPOST.cbl:82](../../../pangaea/programs/TXPOST.cbl#L82) |
| 05 | WS-REJECT-STATUS | 11 | 2 | XX | | | [TXPOST.cbl:83](../../../pangaea/programs/TXPOST.cbl#L83) |
| 05 | WS-REPORT-STATUS | 13 | 2 | XX | | | [TXPOST.cbl:84](../../../pangaea/programs/TXPOST.cbl#L84) |

### WS-BUSINESS-DATE

Working storage. 8 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-BUSINESS-DATE | 1 | 8 | 9(8) | | | [TXPOST.cbl:86](../../../pangaea/programs/TXPOST.cbl#L86) |

### WS-ABORT-MESSAGE

Working storage. 40 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-ABORT-MESSAGE | 1 | 40 | X(40) | | | [TXPOST.cbl:87](../../../pangaea/programs/TXPOST.cbl#L87) |

### WS-MASTER-KEY

Working storage. 10 bytes. MATCH KEYS: HIGH-VALUES ONCE A FILE IS AT ITS END.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-MASTER-KEY | 1 | 10 | X(10) | | | [TXPOST.cbl:90](../../../pangaea/programs/TXPOST.cbl#L90) |

### WS-TRAN-KEY

Working storage. 10 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TRAN-KEY | 1 | 10 | X(10) | | | [TXPOST.cbl:91](../../../pangaea/programs/TXPOST.cbl#L91) |

### WS-PREV-MASTER-KEY

Working storage. 10 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-PREV-MASTER-KEY | 1 | 10 | X(10) | | VALUE LOW-VALUES | [TXPOST.cbl:92](../../../pangaea/programs/TXPOST.cbl#L92) |

### WS-TRAN-ORDER-KEY

Working storage. 20 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TRAN-ORDER-KEY | 1 | 20 | | | | [TXPOST.cbl:93](../../../pangaea/programs/TXPOST.cbl#L93) |
| 05 | WS-TOK-ACCT-NO | 1 | 10 | X(10) | | | [TXPOST.cbl:94](../../../pangaea/programs/TXPOST.cbl#L94) |
| 05 | WS-TOK-TRAN-ID | 11 | 10 | X(10) | | | [TXPOST.cbl:95](../../../pangaea/programs/TXPOST.cbl#L95) |

### WS-PREV-TRAN-ORDER-KEY

Working storage. 20 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-PREV-TRAN-ORDER-KEY | 1 | 20 | X(20) | | VALUE LOW-VALUES | [TXPOST.cbl:96](../../../pangaea/programs/TXPOST.cbl#L96) |

### WS-POSTING

Working storage. 50 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-POSTING | 1 | 50 | | | | [TXPOST.cbl:98](../../../pangaea/programs/TXPOST.cbl#L98) |
| 05 | WS-BALANCE-BEFORE | 1 | 13 | S9(11)V99 | | | [TXPOST.cbl:99](../../../pangaea/programs/TXPOST.cbl#L99) |
| 05 | WS-NEW-BALANCE | 14 | 13 | S9(11)V99 | | | [TXPOST.cbl:100](../../../pangaea/programs/TXPOST.cbl#L100) |
| 05 | WS-BALANCE-FLOOR | 27 | 13 | S9(11)V99 | | | [TXPOST.cbl:101](../../../pangaea/programs/TXPOST.cbl#L101) |
| 05 | WS-SIGNED-AMOUNT | 40 | 11 | S9(9)V99 | | | [TXPOST.cbl:102](../../../pangaea/programs/TXPOST.cbl#L102) |

### WS-TOTALS

Working storage. 122 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TOTALS | 1 | 122 | | | | [TXPOST.cbl:104](../../../pangaea/programs/TXPOST.cbl#L104) |
| 05 | WS-MASTERS-READ | 1 | 8 | 9(8) | | VALUE 0 | [TXPOST.cbl:105](../../../pangaea/programs/TXPOST.cbl#L105) |
| 05 | WS-MASTERS-WRITTEN | 9 | 8 | 9(8) | | VALUE 0 | [TXPOST.cbl:106](../../../pangaea/programs/TXPOST.cbl#L106) |
| 05 | WS-POSTED-COUNT | 17 | 8 | 9(8) | | VALUE 0 | [TXPOST.cbl:107](../../../pangaea/programs/TXPOST.cbl#L107) |
| 05 | WS-REJECT-COUNT | 25 | 8 | 9(8) | | VALUE 0 | [TXPOST.cbl:108](../../../pangaea/programs/TXPOST.cbl#L108) |
| 05 | WS-OPENING-TOTAL | 33 | 15 | S9(13)V99 | | VALUE 0 | [TXPOST.cbl:109](../../../pangaea/programs/TXPOST.cbl#L109) |
| 05 | WS-DEPOSIT-TOTAL | 48 | 15 | 9(13)V99 | | VALUE 0 | [TXPOST.cbl:110](../../../pangaea/programs/TXPOST.cbl#L110) |
| 05 | WS-WITHDRAWAL-TOTAL | 63 | 15 | 9(13)V99 | | VALUE 0 | [TXPOST.cbl:111](../../../pangaea/programs/TXPOST.cbl#L111) |
| 05 | WS-FEE-TOTAL | 78 | 15 | 9(13)V99 | | VALUE 0 | [TXPOST.cbl:112](../../../pangaea/programs/TXPOST.cbl#L112) |
| 05 | WS-CLOSING-TOTAL | 93 | 15 | S9(13)V99 | | VALUE 0 | [TXPOST.cbl:113](../../../pangaea/programs/TXPOST.cbl#L113) |
| 05 | WS-EXPECTED-TOTAL | 108 | 15 | S9(13)V99 | | VALUE 0 | [TXPOST.cbl:114](../../../pangaea/programs/TXPOST.cbl#L114) |

### AUDIT-RECORD

Working storage. 100 bytes, layout in copybook [AUDITREC](../copybooks/AUDITREC.md).

### REJECT-RECORD

Working storage. 100 bytes, layout in copybook [REJREC](../copybooks/REJREC.md).

### WS-REPORT-COUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-COUNT-LINE | 1 | 80 | | | | [TXPOST.cbl:120](../../../pangaea/programs/TXPOST.cbl#L120) |
| 05 | WS-RC-LABEL | 1 | 30 | X(30) | | | [TXPOST.cbl:121](../../../pangaea/programs/TXPOST.cbl#L121) |
| 05 | WS-RC-VALUE | 31 | 10 | ZZ,ZZZ,ZZ9 | | | [TXPOST.cbl:122](../../../pangaea/programs/TXPOST.cbl#L122) |
| 05 | FILLER | 41 | 40 | X(40) | | VALUE SPACES | [TXPOST.cbl:123](../../../pangaea/programs/TXPOST.cbl#L123) |

### WS-REPORT-AMOUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-AMOUNT-LINE | 1 | 80 | | | | [TXPOST.cbl:124](../../../pangaea/programs/TXPOST.cbl#L124) |
| 05 | WS-RA-LABEL | 1 | 30 | X(30) | | | [TXPOST.cbl:125](../../../pangaea/programs/TXPOST.cbl#L125) |
| 05 | WS-RA-VALUE | 31 | 21 | Z,ZZZ,ZZZ,ZZZ,ZZ9.99- | | | [TXPOST.cbl:126](../../../pangaea/programs/TXPOST.cbl#L126) |
| 05 | FILLER | 52 | 29 | X(29) | | VALUE SPACES | [TXPOST.cbl:127](../../../pangaea/programs/TXPOST.cbl#L127) |
