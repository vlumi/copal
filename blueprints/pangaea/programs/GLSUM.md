# GLSUM

```text
GLSUM: SUMMARIZE THE MONTH FOR THE GENERAL LEDGER

ON THE FIRST BUSINESS DAY OF A MONTH, AFTER THE DAILY BATCH,
POSTS EVERY POSTING OF THE MONTH THAT CLOSED TO THE GENERAL
LEDGER, AS STMTPRT PICKS THEM: ONE SIDE TO CUSTOMER DEPOSITS,
THE OTHER TO THE ACCOUNT THE CODE TABLE NAMES FOR ITS
TRANSACTION CODE. A CREDIT TO A CUSTOMER CREDITS DEPOSITS, A
DEBIT DEBITS THEM. WRITES EACH LEDGER ACCOUNT'S DEBITS AND
CREDITS FOR THE MONTH, AND REPORTS THEM.

EVERY TRANSFER'S LEGS POST ON ONE DAY, SO TRANSFERS IN TRANSIT
MUST NET TO ZERO OVER A MONTH. IF THEY DO NOT, A LEG WENT
MISSING, AND THE RUN ENDS WITH RETURN CODE 4.

RETURN CODES: 0 SUMMARIZED OR NO MONTH CLOSED,
              4 TRANSFERS IN TRANSIT DO NOT NET TO ZERO,
              16 AN UNKNOWN CODE OR ACCOUNT, OUT OF BALANCE,
              OR I/O.

MODIFICATION HISTORY
DATE       BY  CHANGE
2004-03-15 TL  CR-0390 ORIGINAL VERSION.
```

Program [GLSUM.cbl](../../../pangaea/programs/GLSUM.cbl): 4 files, 9 paragraphs.

- AUTHOR: T. LINDQVIST, BATCH SYSTEMS.
- INSTALLATION: FIRST BANK OF PANGAEA, DATA CENTER.
- DATE-WRITTEN: 2004-03-15.

## Files

| File | Kind | Assigned to | Organization | Status field | Records |
| --- | --- | --- | --- | --- | --- |
| ACCR-IN | FD | ACCRIN | SEQUENTIAL | WS-ACCR-IN-STATUS | ACCRUAL-RECORD |
| AUDIT-IN | FD | AUDIT | SEQUENTIAL | WS-AUDIT-STATUS | AUDIT-RECORD |
| GL-OUT | FD | GLOUT | SEQUENTIAL | WS-GL-STATUS | GL-SUMMARY-RECORD |
| REPORT-OUT | FD | GLRPT | LINE SEQUENTIAL | WS-REPORT-STATUS | REPORT-LINE |

## Flow

![The paragraphs of GLSUM and how control moves between them](GLSUM.svg)

Solid arrows are PERFORMs, dashed ones GO TOs, and dotted ones control running on into the next paragraph inside a PERFORM THRU range.

## Paragraphs

| Paragraph | Description | Source | Performs | Goes to | Performed by | Reached by |
| --- | --- | --- | --- | --- | --- | --- |
| 0000-MAIN | FIND THE MONTH THAT CLOSED, IF ONE DID, POST ITS POSTINGS TO THE LEDGER AND WRITE THE SUMMARY. | [GLSUM.cbl:156](../../../pangaea/programs/GLSUM.cbl#L156) | 1000-INITIALIZE, 2000-POST-MONTH, 3000-WRITE-LEDGER, 8000-WRITE-REPORT, 9000-CLOSE-FILES | | | |
| 1000-INITIALIZE | OPEN THE FILES, TAKE THE CLOSED MONTH FROM THE ACCRUALS AND START THE LEDGER AT ZERO. | [GLSUM.cbl:173](../../../pangaea/programs/GLSUM.cbl#L173) | | | 0000-MAIN | |
| 2000-POST-MONTH | POST THE POSTINGS OF THE CLOSED MONTH, AND THE INTEREST AND FEES POSTED FOR IT ON THE FIRST BUSINESS DAY OF THE NEXT. | [GLSUM.cbl:199](../../../pangaea/programs/GLSUM.cbl#L199) | 2100-POST-TO-LEDGER, 9900-ABORT | | 0000-MAIN | |
| 2100-POST-TO-LEDGER | POST ONE POSTING: DEPOSITS ON THE CUSTOMER'S SIDE, THE CODE'S LEDGER ACCOUNT ON THE OTHER. | [GLSUM.cbl:218](../../../pangaea/programs/GLSUM.cbl#L218) | 9900-ABORT, 2200-FIND-ACCOUNT | | 2000-POST-MONTH | |
| 2200-FIND-ACCOUNT | FIND THE WANTED LEDGER ACCOUNT, AND POINT AT ITS TOTALS. | [GLSUM.cbl:245](../../../pangaea/programs/GLSUM.cbl#L245) | 9900-ABORT | | 2100-POST-TO-LEDGER | |
| 3000-WRITE-LEDGER | WRITE A SUMMARY RECORD FOR EACH LEDGER ACCOUNT, AND CHECK THAT THE LEDGER BALANCES. | [GLSUM.cbl:257](../../../pangaea/programs/GLSUM.cbl#L257) | 9900-ABORT | | 0000-MAIN | |
| 8000-WRITE-REPORT | WRITE THE GENERAL LEDGER REPORT. | [GLSUM.cbl:284](../../../pangaea/programs/GLSUM.cbl#L284) | | | 0000-MAIN | |
| 9000-CLOSE-FILES | CLOSE ALL FILES. | [GLSUM.cbl:336](../../../pangaea/programs/GLSUM.cbl#L336) | | | 0000-MAIN, 9900-ABORT | |
| 9900-ABORT | STOP THE RUN WITH RETURN CODE 16, SAYING WHY. | [GLSUM.cbl:340](../../../pangaea/programs/GLSUM.cbl#L340) | 9000-CLOSE-FILES | | 2000-POST-MONTH, 2100-POST-TO-LEDGER, 2200-FIND-ACCOUNT, 3000-WRITE-LEDGER | |

## Records

### ACCRUAL-RECORD

File ACCR-IN. 40 bytes, layout in copybook [ACCRREC](../copybooks/ACCRREC.md).

### AUDIT-RECORD

File AUDIT-IN. 100 bytes, layout in copybook [AUDITREC](../copybooks/AUDITREC.md).

### GL-SUMMARY-RECORD

File GL-OUT. 60 bytes, layout in copybook [GLREC](../copybooks/GLREC.md).

### REPORT-LINE

File REPORT-OUT. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | REPORT-LINE | 1 | 80 | X(80) | | | [GLSUM.cbl:56](../../../pangaea/programs/GLSUM.cbl#L56) |

### WS-FILE-STATUSES

Working storage. 8 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-FILE-STATUSES | 1 | 8 | | | | [GLSUM.cbl:59](../../../pangaea/programs/GLSUM.cbl#L59) |
| 05 | WS-ACCR-IN-STATUS | 1 | 2 | XX | | | [GLSUM.cbl:60](../../../pangaea/programs/GLSUM.cbl#L60) |
| 05 | WS-AUDIT-STATUS | 3 | 2 | XX | | | [GLSUM.cbl:61](../../../pangaea/programs/GLSUM.cbl#L61) |
| 05 | WS-GL-STATUS | 5 | 2 | XX | | | [GLSUM.cbl:62](../../../pangaea/programs/GLSUM.cbl#L62) |
| 05 | WS-REPORT-STATUS | 7 | 2 | XX | | | [GLSUM.cbl:63](../../../pangaea/programs/GLSUM.cbl#L63) |

### WS-ABORT-MESSAGE

Working storage. 40 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-ABORT-MESSAGE | 1 | 40 | X(40) | | | [GLSUM.cbl:65](../../../pangaea/programs/GLSUM.cbl#L65) |

### WS-CLOSED-MONTH

Working storage. 6 bytes. THE MONTH THAT CLOSED AND THE ONE AFTER IT, AS IN STMTPRT.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-CLOSED-MONTH | 1 | 6 | | | | [GLSUM.cbl:68](../../../pangaea/programs/GLSUM.cbl#L68) |
| 05 | WS-CM-YEAR | 1 | 4 | 9(4) | | | [GLSUM.cbl:69](../../../pangaea/programs/GLSUM.cbl#L69) |
| 05 | WS-CM-MONTH | 5 | 2 | 99 | | | [GLSUM.cbl:70](../../../pangaea/programs/GLSUM.cbl#L70) |

### WS-CLOSED-MONTH-NUMBER

Working storage. 6 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-CLOSED-MONTH-NUMBER | 1 | 6 | 9(6) | | REDEFINES WS-CLOSED-MONTH | [GLSUM.cbl:71](../../../pangaea/programs/GLSUM.cbl#L71) |

### WS-NEXT-MONTH

Working storage. 6 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-NEXT-MONTH | 1 | 6 | | | | [GLSUM.cbl:73](../../../pangaea/programs/GLSUM.cbl#L73) |
| 05 | WS-NM-YEAR | 1 | 4 | 9(4) | | | [GLSUM.cbl:74](../../../pangaea/programs/GLSUM.cbl#L74) |
| 05 | WS-NM-MONTH | 5 | 2 | 99 | | | [GLSUM.cbl:75](../../../pangaea/programs/GLSUM.cbl#L75) |

### WS-POSTING-MONTH

Working storage. 6 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-POSTING-MONTH | 1 | 6 | X(6) | | | [GLSUM.cbl:76](../../../pangaea/programs/GLSUM.cbl#L76) |

### WS-POSTING-CODE

Working storage. 2 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-POSTING-CODE | 1 | 2 | XX | | | [GLSUM.cbl:77](../../../pangaea/programs/GLSUM.cbl#L77) |
| 88 | WS-MONTH-END-CODE | | | | | VALUE 'IN', 'OD', 'MF' | [GLSUM.cbl:78](../../../pangaea/programs/GLSUM.cbl#L78) |

### GS-CODE-VALUES

Working storage. 234 bytes, layout in copybook [CODETAB](../copybooks/CODETAB.md).

### GS-CODE-TABLE

Working storage. 234 bytes, layout in copybook [CODETAB](../copybooks/CODETAB.md).

### WS-LEDGER-VALUES

Working storage. 168 bytes. THE LEDGER ACCOUNTS, AND WHAT THE MONTH POSTED TO EACH.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-LEDGER-VALUES | 1 | 168 | | | | [GLSUM.cbl:84](../../../pangaea/programs/GLSUM.cbl#L84) |
| 05 | FILLER | 1 | 4 | 9(4) | | VALUE 1000 | [GLSUM.cbl:85](../../../pangaea/programs/GLSUM.cbl#L85) |
| 05 | FILLER | 5 | 24 | X(24) | | VALUE 'CASH' | [GLSUM.cbl:86](../../../pangaea/programs/GLSUM.cbl#L86) |
| 05 | FILLER | 29 | 4 | 9(4) | | VALUE 2000 | [GLSUM.cbl:87](../../../pangaea/programs/GLSUM.cbl#L87) |
| 05 | FILLER | 33 | 24 | X(24) | | VALUE 'CUSTOMER DEPOSITS' | [GLSUM.cbl:88](../../../pangaea/programs/GLSUM.cbl#L88) |
| 05 | FILLER | 57 | 4 | 9(4) | | VALUE 2900 | [GLSUM.cbl:90](../../../pangaea/programs/GLSUM.cbl#L90) |
| 05 | FILLER | 61 | 24 | X(24) | | VALUE 'TRANSFERS IN TRANSIT' | [GLSUM.cbl:91](../../../pangaea/programs/GLSUM.cbl#L91) |
| 05 | FILLER | 85 | 4 | 9(4) | | VALUE 4000 | [GLSUM.cbl:93](../../../pangaea/programs/GLSUM.cbl#L93) |
| 05 | FILLER | 89 | 24 | X(24) | | VALUE 'INTEREST INCOME' | [GLSUM.cbl:94](../../../pangaea/programs/GLSUM.cbl#L94) |
| 05 | FILLER | 113 | 4 | 9(4) | | VALUE 4100 | [GLSUM.cbl:96](../../../pangaea/programs/GLSUM.cbl#L96) |
| 05 | FILLER | 117 | 24 | X(24) | | VALUE 'FEE INCOME' | [GLSUM.cbl:97](../../../pangaea/programs/GLSUM.cbl#L97) |
| 05 | FILLER | 141 | 4 | 9(4) | | VALUE 5000 | [GLSUM.cbl:99](../../../pangaea/programs/GLSUM.cbl#L99) |
| 05 | FILLER | 145 | 24 | X(24) | | VALUE 'INTEREST EXPENSE' | [GLSUM.cbl:100](../../../pangaea/programs/GLSUM.cbl#L100) |

### WS-LEDGER

Working storage. 168 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-LEDGER | 1 | 168 | | | REDEFINES WS-LEDGER-VALUES | [GLSUM.cbl:102](../../../pangaea/programs/GLSUM.cbl#L102) |
| 05 | WS-LEDGER-ENTRY | 1 | 28 | | | OCCURS 6; INDEXED BY WS-LX | [GLSUM.cbl:103](../../../pangaea/programs/GLSUM.cbl#L103) |
| 10 | WS-LE-ACCOUNT | 1 | 4 | 9(4) | | | [GLSUM.cbl:105](../../../pangaea/programs/GLSUM.cbl#L105) |
| 10 | WS-LE-NAME | 5 | 24 | X(24) | | | [GLSUM.cbl:106](../../../pangaea/programs/GLSUM.cbl#L106) |

### WS-LEDGER-TOTALS

Working storage. 96 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-LEDGER-TOTALS | 1 | 96 | | | | [GLSUM.cbl:107](../../../pangaea/programs/GLSUM.cbl#L107) |
| 05 | WS-LEDGER-TOTAL | 1 | 16 | | | OCCURS 6; INDEXED BY WS-TX | [GLSUM.cbl:108](../../../pangaea/programs/GLSUM.cbl#L108) |
| 10 | WS-LT-DEBITS | 1 | 8 | S9(13)V99 | COMP-3 | | [GLSUM.cbl:110](../../../pangaea/programs/GLSUM.cbl#L110) |
| 10 | WS-LT-CREDITS | 9 | 8 | S9(13)V99 | COMP-3 | | [GLSUM.cbl:111](../../../pangaea/programs/GLSUM.cbl#L111) |

### WS-DEPOSITS-ACCOUNT

Working storage. 4 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-DEPOSITS-ACCOUNT | 1 | 4 | 9(4) | | VALUE 2000 | [GLSUM.cbl:113](../../../pangaea/programs/GLSUM.cbl#L113) |

### WS-TRANSIT-ACCOUNT

Working storage. 4 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TRANSIT-ACCOUNT | 1 | 4 | 9(4) | | VALUE 2900 | [GLSUM.cbl:114](../../../pangaea/programs/GLSUM.cbl#L114) |

### WS-POSTING

Working storage. 23 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-POSTING | 1 | 23 | | | | [GLSUM.cbl:115](../../../pangaea/programs/GLSUM.cbl#L115) |
| 05 | WS-AMOUNT | 1 | 11 | S9(9)V99 | | | [GLSUM.cbl:116](../../../pangaea/programs/GLSUM.cbl#L116) |
| 05 | WS-DEBIT-ACCOUNT | 12 | 4 | 9(4) | | | [GLSUM.cbl:117](../../../pangaea/programs/GLSUM.cbl#L117) |
| 05 | WS-CREDIT-ACCOUNT | 16 | 4 | 9(4) | | | [GLSUM.cbl:118](../../../pangaea/programs/GLSUM.cbl#L118) |
| 05 | WS-WANTED-ACCOUNT | 20 | 4 | 9(4) | | | [GLSUM.cbl:119](../../../pangaea/programs/GLSUM.cbl#L119) |

### WS-TOTALS

Working storage. 68 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TOTALS | 1 | 68 | | | | [GLSUM.cbl:121](../../../pangaea/programs/GLSUM.cbl#L121) |
| 05 | WS-POSTING-COUNT | 1 | 8 | 9(8) | | VALUE 0 | [GLSUM.cbl:122](../../../pangaea/programs/GLSUM.cbl#L122) |
| 05 | WS-DEBIT-TOTAL | 9 | 15 | S9(13)V99 | | VALUE 0 | [GLSUM.cbl:123](../../../pangaea/programs/GLSUM.cbl#L123) |
| 05 | WS-CREDIT-TOTAL | 24 | 15 | S9(13)V99 | | VALUE 0 | [GLSUM.cbl:124](../../../pangaea/programs/GLSUM.cbl#L124) |
| 05 | WS-DEPOSITS-CHANGE | 39 | 15 | S9(13)V99 | | VALUE 0 | [GLSUM.cbl:125](../../../pangaea/programs/GLSUM.cbl#L125) |
| 05 | WS-TRANSIT-NET | 54 | 15 | S9(13)V99 | | VALUE 0 | [GLSUM.cbl:126](../../../pangaea/programs/GLSUM.cbl#L126) |

### WS-HEADING-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-HEADING-LINE | 1 | 80 | | | | [GLSUM.cbl:128](../../../pangaea/programs/GLSUM.cbl#L128) |
| 05 | FILLER | 1 | 30 | X(30) | | VALUE 'ACCOUNT' | [GLSUM.cbl:129](../../../pangaea/programs/GLSUM.cbl#L129) |
| 05 | FILLER | 31 | 18 | X(18) | | VALUE '            DEBITS' | [GLSUM.cbl:130](../../../pangaea/programs/GLSUM.cbl#L130) |
| 05 | FILLER | 49 | 2 | X(2) | | VALUE SPACES | [GLSUM.cbl:132](../../../pangaea/programs/GLSUM.cbl#L132) |
| 05 | FILLER | 51 | 18 | X(18) | | VALUE '           CREDITS' | [GLSUM.cbl:133](../../../pangaea/programs/GLSUM.cbl#L133) |
| 05 | FILLER | 69 | 12 | X(12) | | VALUE SPACES | [GLSUM.cbl:135](../../../pangaea/programs/GLSUM.cbl#L135) |

### WS-LEDGER-LINE

Working storage. 82 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-LEDGER-LINE | 1 | 82 | | | | [GLSUM.cbl:136](../../../pangaea/programs/GLSUM.cbl#L136) |
| 05 | WS-LL-ACCOUNT | 1 | 4 | X(4) | | | [GLSUM.cbl:137](../../../pangaea/programs/GLSUM.cbl#L137) |
| 05 | FILLER | 5 | 1 | X(1) | | VALUE SPACES | [GLSUM.cbl:138](../../../pangaea/programs/GLSUM.cbl#L138) |
| 05 | WS-LL-NAME | 6 | 25 | X(25) | | | [GLSUM.cbl:139](../../../pangaea/programs/GLSUM.cbl#L139) |
| 05 | WS-LL-DEBITS | 31 | 20 | Z,ZZZ,ZZZ,ZZZ,ZZ9.99 | | | [GLSUM.cbl:140](../../../pangaea/programs/GLSUM.cbl#L140) |
| 05 | FILLER | 51 | 2 | X(2) | | VALUE SPACES | [GLSUM.cbl:141](../../../pangaea/programs/GLSUM.cbl#L141) |
| 05 | WS-LL-CREDITS | 53 | 20 | Z,ZZZ,ZZZ,ZZZ,ZZ9.99 | | | [GLSUM.cbl:142](../../../pangaea/programs/GLSUM.cbl#L142) |
| 05 | FILLER | 73 | 10 | X(10) | | VALUE SPACES | [GLSUM.cbl:143](../../../pangaea/programs/GLSUM.cbl#L143) |

### WS-REPORT-COUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-COUNT-LINE | 1 | 80 | | | | [GLSUM.cbl:144](../../../pangaea/programs/GLSUM.cbl#L144) |
| 05 | WS-RC-LABEL | 1 | 30 | X(30) | | | [GLSUM.cbl:145](../../../pangaea/programs/GLSUM.cbl#L145) |
| 05 | WS-RC-VALUE | 31 | 10 | ZZ,ZZZ,ZZ9 | | | [GLSUM.cbl:146](../../../pangaea/programs/GLSUM.cbl#L146) |
| 05 | FILLER | 41 | 40 | X(40) | | VALUE SPACES | [GLSUM.cbl:147](../../../pangaea/programs/GLSUM.cbl#L147) |

### WS-REPORT-AMOUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-AMOUNT-LINE | 1 | 80 | | | | [GLSUM.cbl:148](../../../pangaea/programs/GLSUM.cbl#L148) |
| 05 | WS-RA-LABEL | 1 | 30 | X(30) | | | [GLSUM.cbl:149](../../../pangaea/programs/GLSUM.cbl#L149) |
| 05 | WS-RA-VALUE | 31 | 21 | Z,ZZZ,ZZZ,ZZZ,ZZ9.99- | | | [GLSUM.cbl:150](../../../pangaea/programs/GLSUM.cbl#L150) |
| 05 | FILLER | 52 | 29 | X(29) | | VALUE SPACES | [GLSUM.cbl:151](../../../pangaea/programs/GLSUM.cbl#L151) |
