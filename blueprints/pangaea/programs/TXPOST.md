# TXPOST

```text
TXPOST: POST TRANSACTIONS TO THE ACCOUNT MASTER

MATCHES THE SORTED TRANSACTIONS AGAINST THE OLD MASTER, BOTH
IN ACCOUNT NUMBER ORDER, AND WRITES THE NEW MASTER: EVERY OLD
RECORD, WITH ITS TRANSACTIONS APPLIED. EACH POSTING GOES TO
THE AUDIT TRAIL; A TRANSACTION THAT CANNOT BE POSTED GOES TO
THE REJECTS. BOTH ARE APPENDED TO, AS THE DAILY BATCH RUNS
TXPOST FOUR TIMES:
  0. THE MONTH-END INTEREST AND FEES FROM MONCHG, ON THE DAY
     A MONTH CLOSES; ON OTHER DAYS NONE.
  1. THE DAY'S TRANSACTIONS. A POSTED TRANSFER WRITES A
     TRANSFER CREDIT FOR THE RECEIVING ACCOUNT TO THE LEGS.
  2. THE TRANSFER CREDITS, SORTED. ONE FOR AN UNKNOWN OR
     CLOSED ACCOUNT WRITES A RETURNED TRANSFER TO THE LEGS.
  3. THE RETURNED TRANSFERS, SORTED.

RULES: NOTHING POSTS TO A CLOSED ACCOUNT. CREDITS AND FEES
POST TO A FROZEN ACCOUNT, WITHDRAWALS AND TRANSFERS DO NOT. A
DEBIT MAY TAKE THE BALANCE DOWN TO THE OVERDRAFT LIMIT, NO
FURTHER. A FEE ALWAYS POSTS, EVEN PAST THE LIMIT. INTEREST
POSTS AS A CREDIT, OVERDRAFT INTEREST AND MAINTENANCE FEES AS
FEES.

RETURN CODES: 0 ALL POSTED, 4 SOME REJECTED,
              16 INPUT OUT OF ORDER, OUT OF BALANCE OR I/O.

MODIFICATION HISTORY
DATE       BY  CHANGE
1987-04-02 JA  ORIGINAL VERSION.
1994-11-14 PK  CR-0201 OVERDRAFT LIMITS.
1996-04-15 PK  CR-0233 TRANSFERS, POSTED IN LEGS OVER THREE
               RUNS. AUDIT AND REJECTS ARE APPENDED TO.
1999-01-11 SL  Y2K: 8-DIGIT DATES IN THE AUDIT TRAIL.
2003-05-26 MV  CR-0352 FEES POST EVEN PAST THE OVERDRAFT
               LIMIT.
2001-09-03 TL  CR-0298 MONTH-END INTEREST AND FEES.
2011-08-15 MV  CR-0471 RECONCILE THE NEW MASTER BEFORE
               ENDING THE RUN.
```

Program [TXPOST.cbl](../../../pangaea/programs/TXPOST.cbl): 8 files, 17 paragraphs.

- AUTHOR: J. ARKWRIGHT, BATCH SYSTEMS.
- INSTALLATION: FIRST BANK OF PANGAEA, DATA CENTER.
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
| LEGS-OUT | FD | LEGS | SEQUENTIAL | WS-LEGS-STATUS | LEGS-OUT-RECORD |
| REPORT-OUT | FD | POSTRPT | LINE SEQUENTIAL | WS-REPORT-STATUS | REPORT-LINE |

## Flow

![The paragraphs of TXPOST and the PERFORMs between them](TXPOST.svg)

## Paragraphs

| Paragraph | Description | Source | Performs | Performed by |
| --- | --- | --- | --- | --- |
| 0000-MAIN | MATCH UNTIL BOTH FILES END, THEN RECONCILE AND REPORT. | [TXPOST.cbl:158](../../../pangaea/programs/TXPOST.cbl#L158) | 1000-INITIALIZE, 2000-MATCH, 5000-RECONCILE, 8000-WRITE-REPORT, 9000-CLOSE-FILES | |
| 1000-INITIALIZE | OPEN THE FILES AND READ THE FIRST RECORD OF EACH. THE AUDIT TRAIL AND THE REJECTS ARE EXTENDED, AND CREATED IF NEW. | [TXPOST.cbl:175](../../../pangaea/programs/TXPOST.cbl#L175) | 9900-ABORT, 1200-READ-MASTER, 1300-READ-TRAN | 0000-MAIN |
| 1200-READ-MASTER | READ THE NEXT OLD MASTER RECORD, CHECKING ACCOUNT ORDER. | [TXPOST.cbl:202](../../../pangaea/programs/TXPOST.cbl#L202) | 9900-ABORT | 1000-INITIALIZE, 2000-MATCH |
| 1300-READ-TRAN | READ THE NEXT TRANSACTION, CHECKING ITS ORDER. | [TXPOST.cbl:219](../../../pangaea/programs/TXPOST.cbl#L219) | 9900-ABORT | 1000-INITIALIZE, 2000-MATCH |
| 2000-MATCH | ONE STEP OF THE MATCH, BY WHICHEVER KEY IS LOWER: REJECT A TRANSACTION FOR A MISSING ACCOUNT, APPLY ONE TO ITS ACCOUNT, OR WRITE A FINISHED ACCOUNT. | [TXPOST.cbl:237](../../../pangaea/programs/TXPOST.cbl#L237) | 3900-REJECT-TRAN, 1300-READ-TRAN, 3000-APPLY-TRAN, 4000-WRITE-MASTER, 1200-READ-MASTER | 0000-MAIN |
| 3000-APPLY-TRAN | APPLY ONE TRANSACTION TO ITS ACCOUNT BY THE POSTING RULES. | [TXPOST.cbl:252](../../../pangaea/programs/TXPOST.cbl#L252) | 3900-REJECT-TRAN, 3100-POST-CREDIT, 3300-POST-FEE, 3200-POST-DEBIT | 2000-MATCH |
| 3100-POST-CREDIT | CREDIT A DEPOSIT, A TRANSFER CREDIT, A RETURNED TRANSFER OR INTEREST. | [TXPOST.cbl:277](../../../pangaea/programs/TXPOST.cbl#L277) | 3800-WRITE-AUDIT | 3000-APPLY-TRAN |
| 3200-POST-DEBIT | DEBIT A WITHDRAWAL OR TRANSFER IF IT STAYS WITHIN THE OVERDRAFT LIMIT. A POSTED TRANSFER SENDS ITS CREDIT ON. | [TXPOST.cbl:285](../../../pangaea/programs/TXPOST.cbl#L285) | 3900-REJECT-TRAN, 3800-WRITE-AUDIT, 3500-WRITE-LEG | 3000-APPLY-TRAN |
| 3300-POST-FEE | DEBIT A FEE, OVERDRAFT INTEREST OR A MAINTENANCE FEE, WHATEVER THE BALANCE. | [TXPOST.cbl:303](../../../pangaea/programs/TXPOST.cbl#L303) | 3800-WRITE-AUDIT | 3000-APPLY-TRAN |
| 3500-WRITE-LEG | WRITE THE NEXT LEG OF A TRANSFER FOR THE COUNTER ACCOUNT: A TRANSFER CREDIT AFTER A TRANSFER, A RETURN AFTER A CREDIT THAT COULD NOT POST. | [TXPOST.cbl:312](../../../pangaea/programs/TXPOST.cbl#L312) | | 3200-POST-DEBIT, 3900-REJECT-TRAN |
| 3800-WRITE-AUDIT | RECORD A POSTING IN THE AUDIT TRAIL. | [TXPOST.cbl:327](../../../pangaea/programs/TXPOST.cbl#L327) | | 3100-POST-CREDIT, 3200-POST-DEBIT, 3300-POST-FEE |
| 3900-REJECT-TRAN | WRITE THE TRANSACTION TO THE REJECTS WITH ITS REASON. A TRANSFER CREDIT THAT CANNOT POST GOES BACK TO ITS SENDER. | [TXPOST.cbl:341](../../../pangaea/programs/TXPOST.cbl#L341) | 3500-WRITE-LEG | 2000-MATCH, 3000-APPLY-TRAN, 3200-POST-DEBIT |
| 4000-WRITE-MASTER | WRITE AN ACCOUNT, ITS TRANSACTIONS APPLIED, TO THE NEW MASTER. | [TXPOST.cbl:351](../../../pangaea/programs/TXPOST.cbl#L351) | | 2000-MATCH |
| 5000-RECONCILE | THE NEW MASTER MUST HOLD THE OLD BALANCES PLUS WHAT WAS POSTED. | [TXPOST.cbl:357](../../../pangaea/programs/TXPOST.cbl#L357) | 9900-ABORT | 0000-MAIN |
| 8000-WRITE-REPORT | WRITE THE POSTING REPORT. | [TXPOST.cbl:369](../../../pangaea/programs/TXPOST.cbl#L369) | | 0000-MAIN |
| 9000-CLOSE-FILES | CLOSE ALL FILES. | [TXPOST.cbl:415](../../../pangaea/programs/TXPOST.cbl#L415) | | 0000-MAIN, 9900-ABORT |
| 9900-ABORT | STOP THE RUN WITH RETURN CODE 16, SAYING WHY. | [TXPOST.cbl:420](../../../pangaea/programs/TXPOST.cbl#L420) | 9000-CLOSE-FILES | 1000-INITIALIZE, 1200-READ-MASTER, 1300-READ-TRAN, 5000-RECONCILE |

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
| 01 | ACCT-OUT-RECORD | 1 | 100 | X(100) | | | [TXPOST.cbl:84](../../../pangaea/programs/TXPOST.cbl#L84) |

### AUDIT-OUT-RECORD

File AUDIT-OUT. 100 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | AUDIT-OUT-RECORD | 1 | 100 | X(100) | | | [TXPOST.cbl:86](../../../pangaea/programs/TXPOST.cbl#L86) |

### REJECT-OUT-RECORD

File REJECT-OUT. 100 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | REJECT-OUT-RECORD | 1 | 100 | X(100) | | | [TXPOST.cbl:88](../../../pangaea/programs/TXPOST.cbl#L88) |

### LEGS-OUT-RECORD

File LEGS-OUT. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | LEGS-OUT-RECORD | 1 | 80 | X(80) | | | [TXPOST.cbl:90](../../../pangaea/programs/TXPOST.cbl#L90) |

### REPORT-LINE

File REPORT-OUT. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | REPORT-LINE | 1 | 80 | X(80) | | | [TXPOST.cbl:92](../../../pangaea/programs/TXPOST.cbl#L92) |

### WS-FILE-STATUSES

Working storage. 16 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-FILE-STATUSES | 1 | 16 | | | | [TXPOST.cbl:95](../../../pangaea/programs/TXPOST.cbl#L95) |
| 05 | WS-CONTROL-STATUS | 1 | 2 | XX | | | [TXPOST.cbl:96](../../../pangaea/programs/TXPOST.cbl#L96) |
| 05 | WS-ACCT-IN-STATUS | 3 | 2 | XX | | | [TXPOST.cbl:97](../../../pangaea/programs/TXPOST.cbl#L97) |
| 05 | WS-TRAN-IN-STATUS | 5 | 2 | XX | | | [TXPOST.cbl:98](../../../pangaea/programs/TXPOST.cbl#L98) |
| 05 | WS-ACCT-OUT-STATUS | 7 | 2 | XX | | | [TXPOST.cbl:99](../../../pangaea/programs/TXPOST.cbl#L99) |
| 05 | WS-AUDIT-STATUS | 9 | 2 | XX | | | [TXPOST.cbl:100](../../../pangaea/programs/TXPOST.cbl#L100) |
| 88 | AUDIT-OPENED | | | | | VALUE '00', '05' | [TXPOST.cbl:101](../../../pangaea/programs/TXPOST.cbl#L101) |
| 05 | WS-REJECT-STATUS | 11 | 2 | XX | | | [TXPOST.cbl:102](../../../pangaea/programs/TXPOST.cbl#L102) |
| 88 | REJECTS-OPENED | | | | | VALUE '00', '05' | [TXPOST.cbl:103](../../../pangaea/programs/TXPOST.cbl#L103) |
| 05 | WS-LEGS-STATUS | 13 | 2 | XX | | | [TXPOST.cbl:104](../../../pangaea/programs/TXPOST.cbl#L104) |
| 05 | WS-REPORT-STATUS | 15 | 2 | XX | | | [TXPOST.cbl:105](../../../pangaea/programs/TXPOST.cbl#L105) |

### WS-BUSINESS-DATE

Working storage. 8 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-BUSINESS-DATE | 1 | 8 | 9(8) | | | [TXPOST.cbl:107](../../../pangaea/programs/TXPOST.cbl#L107) |

### WS-ABORT-MESSAGE

Working storage. 40 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-ABORT-MESSAGE | 1 | 40 | X(40) | | | [TXPOST.cbl:108](../../../pangaea/programs/TXPOST.cbl#L108) |

### WS-MASTER-KEY

Working storage. 10 bytes. MATCH KEYS: HIGH-VALUES ONCE A FILE IS AT ITS END.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-MASTER-KEY | 1 | 10 | X(10) | | | [TXPOST.cbl:111](../../../pangaea/programs/TXPOST.cbl#L111) |

### WS-TRAN-KEY

Working storage. 10 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TRAN-KEY | 1 | 10 | X(10) | | | [TXPOST.cbl:112](../../../pangaea/programs/TXPOST.cbl#L112) |

### WS-PREV-MASTER-KEY

Working storage. 10 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-PREV-MASTER-KEY | 1 | 10 | X(10) | | VALUE LOW-VALUES | [TXPOST.cbl:113](../../../pangaea/programs/TXPOST.cbl#L113) |

### WS-TRAN-ORDER-KEY

Working storage. 20 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TRAN-ORDER-KEY | 1 | 20 | | | | [TXPOST.cbl:114](../../../pangaea/programs/TXPOST.cbl#L114) |
| 05 | WS-TOK-ACCT-NO | 1 | 10 | X(10) | | | [TXPOST.cbl:115](../../../pangaea/programs/TXPOST.cbl#L115) |
| 05 | WS-TOK-TRAN-ID | 11 | 10 | X(10) | | | [TXPOST.cbl:116](../../../pangaea/programs/TXPOST.cbl#L116) |

### WS-PREV-TRAN-ORDER-KEY

Working storage. 20 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-PREV-TRAN-ORDER-KEY | 1 | 20 | X(20) | | VALUE LOW-VALUES | [TXPOST.cbl:117](../../../pangaea/programs/TXPOST.cbl#L117) |

### WS-POSTING

Working storage. 50 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-POSTING | 1 | 50 | | | | [TXPOST.cbl:119](../../../pangaea/programs/TXPOST.cbl#L119) |
| 05 | WS-BALANCE-BEFORE | 1 | 13 | S9(11)V99 | | | [TXPOST.cbl:120](../../../pangaea/programs/TXPOST.cbl#L120) |
| 05 | WS-NEW-BALANCE | 14 | 13 | S9(11)V99 | | | [TXPOST.cbl:121](../../../pangaea/programs/TXPOST.cbl#L121) |
| 05 | WS-BALANCE-FLOOR | 27 | 13 | S9(11)V99 | | | [TXPOST.cbl:122](../../../pangaea/programs/TXPOST.cbl#L122) |
| 05 | WS-SIGNED-AMOUNT | 40 | 11 | S9(9)V99 | | | [TXPOST.cbl:123](../../../pangaea/programs/TXPOST.cbl#L123) |

### WS-SAVED-TRAN

Working storage. 80 bytes. A LEG IS WRITTEN FROM THE TRANSACTION RECORD ITSELF, TURNED AROUND; THE TRANSACTION IS SAVED HERE AND PUT BACK AFTER.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-SAVED-TRAN | 1 | 80 | X(80) | | | [TXPOST.cbl:127](../../../pangaea/programs/TXPOST.cbl#L127) |

### WS-LEG-ACCT-NO

Working storage. 10 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-LEG-ACCT-NO | 1 | 10 | 9(10) | | | [TXPOST.cbl:128](../../../pangaea/programs/TXPOST.cbl#L128) |

### WS-TOTALS

Working storage. 130 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TOTALS | 1 | 130 | | | | [TXPOST.cbl:130](../../../pangaea/programs/TXPOST.cbl#L130) |
| 05 | WS-MASTERS-READ | 1 | 8 | 9(8) | | VALUE 0 | [TXPOST.cbl:131](../../../pangaea/programs/TXPOST.cbl#L131) |
| 05 | WS-MASTERS-WRITTEN | 9 | 8 | 9(8) | | VALUE 0 | [TXPOST.cbl:132](../../../pangaea/programs/TXPOST.cbl#L132) |
| 05 | WS-POSTED-COUNT | 17 | 8 | 9(8) | | VALUE 0 | [TXPOST.cbl:133](../../../pangaea/programs/TXPOST.cbl#L133) |
| 05 | WS-REJECT-COUNT | 25 | 8 | 9(8) | | VALUE 0 | [TXPOST.cbl:134](../../../pangaea/programs/TXPOST.cbl#L134) |
| 05 | WS-LEG-COUNT | 33 | 8 | 9(8) | | VALUE 0 | [TXPOST.cbl:135](../../../pangaea/programs/TXPOST.cbl#L135) |
| 05 | WS-OPENING-TOTAL | 41 | 15 | S9(13)V99 | | VALUE 0 | [TXPOST.cbl:136](../../../pangaea/programs/TXPOST.cbl#L136) |
| 05 | WS-CREDIT-TOTAL | 56 | 15 | 9(13)V99 | | VALUE 0 | [TXPOST.cbl:137](../../../pangaea/programs/TXPOST.cbl#L137) |
| 05 | WS-DEBIT-TOTAL | 71 | 15 | 9(13)V99 | | VALUE 0 | [TXPOST.cbl:138](../../../pangaea/programs/TXPOST.cbl#L138) |
| 05 | WS-FEE-TOTAL | 86 | 15 | 9(13)V99 | | VALUE 0 | [TXPOST.cbl:139](../../../pangaea/programs/TXPOST.cbl#L139) |
| 05 | WS-CLOSING-TOTAL | 101 | 15 | S9(13)V99 | | VALUE 0 | [TXPOST.cbl:140](../../../pangaea/programs/TXPOST.cbl#L140) |
| 05 | WS-EXPECTED-TOTAL | 116 | 15 | S9(13)V99 | | VALUE 0 | [TXPOST.cbl:141](../../../pangaea/programs/TXPOST.cbl#L141) |

### AUDIT-RECORD

Working storage. 100 bytes, layout in copybook [AUDITREC](../copybooks/AUDITREC.md).

### REJECT-RECORD

Working storage. 100 bytes, layout in copybook [REJREC](../copybooks/REJREC.md).

### WS-REPORT-COUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-COUNT-LINE | 1 | 80 | | | | [TXPOST.cbl:147](../../../pangaea/programs/TXPOST.cbl#L147) |
| 05 | WS-RC-LABEL | 1 | 30 | X(30) | | | [TXPOST.cbl:148](../../../pangaea/programs/TXPOST.cbl#L148) |
| 05 | WS-RC-VALUE | 31 | 10 | ZZ,ZZZ,ZZ9 | | | [TXPOST.cbl:149](../../../pangaea/programs/TXPOST.cbl#L149) |
| 05 | FILLER | 41 | 40 | X(40) | | VALUE SPACES | [TXPOST.cbl:150](../../../pangaea/programs/TXPOST.cbl#L150) |

### WS-REPORT-AMOUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-AMOUNT-LINE | 1 | 80 | | | | [TXPOST.cbl:151](../../../pangaea/programs/TXPOST.cbl#L151) |
| 05 | WS-RA-LABEL | 1 | 30 | X(30) | | | [TXPOST.cbl:152](../../../pangaea/programs/TXPOST.cbl#L152) |
| 05 | WS-RA-VALUE | 31 | 21 | Z,ZZZ,ZZZ,ZZZ,ZZ9.99- | | | [TXPOST.cbl:153](../../../pangaea/programs/TXPOST.cbl#L153) |
| 05 | FILLER | 52 | 29 | X(29) | | VALUE SPACES | [TXPOST.cbl:154](../../../pangaea/programs/TXPOST.cbl#L154) |
