# INTACCR

```text
INTACCR: ACCRUE INTEREST

ACCRUES EVERY ACCOUNT'S INTEREST FOR THE CALENDAR DAYS FROM
ITS LAST ACCRUAL UP TO THE BUSINESS DATE, ON THE OPENING
BALANCE: THE BALANCE THE LAST BUSINESS DAY CLOSED WITH, WHICH
STOOD OVER THE DAYS IN BETWEEN. THE RATE IS THE ACCOUNT
TYPE'S HIGHEST CREDIT TIER THE BALANCE REACHES, OR ITS DEBIT
RATE BELOW ZERO. A YEAR HAS 365 DAYS.

WHEN THE DAYS CROSS INTO A NEW MONTH, THE DAYS OF THE OLD
MONTH ARE ACCRUED FIRST AND ITS INTEREST BECOMES DUE, ROUNDED
TO THE CENT, TO BE CREDITED OR CHARGED. THE ACCRUALS ARE READ
IN STEP WITH THE MASTER, ONE FOR EACH ACCOUNT.

RETURN CODES: 0 ACCRUED,
              16 FILES OUT OF STEP, NO RATE, OR I/O.

MODIFICATION HISTORY
DATE       BY  CHANGE
2001-09-03 TL  CR-0298 ORIGINAL VERSION.
```

Program [INTACCR.cbl](../../../pangaea/programs/INTACCR.cbl): 6 files, 13 paragraphs.

- AUTHOR: T. LINDQVIST, BATCH SYSTEMS.
- INSTALLATION: FIRST BANK OF PANGAEA, DATA CENTER.
- DATE-WRITTEN: 2001-09-03.

## Files

| File | Kind | Assigned to | Organization | Status field | Records |
| --- | --- | --- | --- | --- | --- |
| CONTROL-FILE | FD | BANKCTL | SEQUENTIAL | WS-CONTROL-STATUS | BANK-CONTROL-RECORD |
| RATE-FILE | FD | RATES | SEQUENTIAL | WS-RATE-STATUS | RATE-RECORD |
| ACCT-IN | FD | ACCTIN | SEQUENTIAL | WS-ACCT-IN-STATUS | ACCT-MASTER-RECORD |
| ACCR-IN | FD | ACCRIN | SEQUENTIAL | WS-ACCR-IN-STATUS | ACCRUAL-RECORD |
| ACCR-OUT | FD | ACCROUT | SEQUENTIAL | WS-ACCR-OUT-STATUS | ACCR-OUT-RECORD |
| REPORT-OUT | FD | ACCRRPT | LINE SEQUENTIAL | WS-REPORT-STATUS | REPORT-LINE |

## Flow

![The paragraphs of INTACCR and the PERFORMs between them](INTACCR.svg)

## Paragraphs

| Paragraph | Description | Source | Performs | Performed by |
| --- | --- | --- | --- | --- |
| 0000-MAIN | ACCRUE EVERY ACCOUNT, THEN REPORT. | [INTACCR.cbl:146](../../../pangaea/programs/INTACCR.cbl#L146) | 1000-INITIALIZE, 2000-ACCRUE-ACCOUNT, 8000-WRITE-REPORT, 9000-CLOSE-FILES | |
| 1000-INITIALIZE | OPEN THE FILES, LOAD THE RATES AND READ THE FIRST ACCOUNT. | [INTACCR.cbl:156](../../../pangaea/programs/INTACCR.cbl#L156) | 9900-ABORT, 1100-LOAD-RATES, 1200-READ-ACCOUNT | 0000-MAIN |
| 1100-LOAD-RATES | READ THE RATE FILE INTO THE RATE TABLE. | [INTACCR.cbl:175](../../../pangaea/programs/INTACCR.cbl#L175) | 9900-ABORT | 1000-INITIALIZE |
| 1200-READ-ACCOUNT | READ THE NEXT ACCOUNT AND ITS ACCRUAL, WHICH MUST MATCH. | [INTACCR.cbl:195](../../../pangaea/programs/INTACCR.cbl#L195) | 9900-ABORT | 1000-INITIALIZE, 2000-ACCRUE-ACCOUNT |
| 2000-ACCRUE-ACCOUNT | ACCRUE ONE ACCOUNT UP TO THE BUSINESS DATE. IF A MONTH CLOSED ON THE WAY, ACCRUE ITS DAYS FIRST AND MAKE ITS INTEREST DUE. | [INTACCR.cbl:217](../../../pangaea/programs/INTACCR.cbl#L217) | 2100-FIND-RATE, 2200-FIND-MONTH-START, 2300-ACCRUE-DAYS, 2400-CLOSE-MONTH, 1200-READ-ACCOUNT | 0000-MAIN |
| 2100-FIND-RATE | FIND THE ACCOUNT'S RATE: BELOW ZERO THE DEBIT RATE OF ITS TYPE, OTHERWISE THE CREDIT RATE OF THE HIGHEST TIER IT REACHES. | [INTACCR.cbl:238](../../../pangaea/programs/INTACCR.cbl#L238) | 2110-TAKE-RATE, 9900-ABORT | 2000-ACCRUE-ACCOUNT |
| 2110-TAKE-RATE | TAKE THE RATE IN HAND AS THE ACCOUNT'S, SO FAR. | [INTACCR.cbl:263](../../../pangaea/programs/INTACCR.cbl#L263) | | 2100-FIND-RATE |
| 2200-FIND-MONTH-START | THE FIRST DAY OF THE MONTH AFTER THE THROUGH DATE. | [INTACCR.cbl:269](../../../pangaea/programs/INTACCR.cbl#L269) | | 2000-ACCRUE-ACCOUNT |
| 2300-ACCRUE-DAYS | ACCRUE THE OPENING BALANCE AT THE ACCOUNT'S RATE FOR THE DAYS. | [INTACCR.cbl:283](../../../pangaea/programs/INTACCR.cbl#L283) | | 2000-ACCRUE-ACCOUNT |
| 2400-CLOSE-MONTH | MAKE THE CLOSED MONTH'S INTEREST DUE, TO THE CENT, AND KEEP THE FRACTION ACCRUED. | [INTACCR.cbl:298](../../../pangaea/programs/INTACCR.cbl#L298) | | 2000-ACCRUE-ACCOUNT |
| 8000-WRITE-REPORT | WRITE THE ACCRUAL REPORT. | [INTACCR.cbl:312](../../../pangaea/programs/INTACCR.cbl#L312) | | 0000-MAIN |
| 9000-CLOSE-FILES | CLOSE ALL FILES. | [INTACCR.cbl:351](../../../pangaea/programs/INTACCR.cbl#L351) | | 0000-MAIN, 9900-ABORT |
| 9900-ABORT | STOP THE RUN WITH RETURN CODE 16, SAYING WHY. | [INTACCR.cbl:356](../../../pangaea/programs/INTACCR.cbl#L356) | 9000-CLOSE-FILES | 1000-INITIALIZE, 1100-LOAD-RATES, 1200-READ-ACCOUNT, 2100-FIND-RATE |

## Records

### BANK-CONTROL-RECORD

File CONTROL-FILE. 80 bytes, layout in copybook [BANKCTL](../copybooks/BANKCTL.md).

### RATE-RECORD

File RATE-FILE. 20 bytes, layout in copybook [RATEREC](../copybooks/RATEREC.md).

### ACCT-MASTER-RECORD

File ACCT-IN. 100 bytes, layout in copybook [ACCTMAST](../copybooks/ACCTMAST.md).

### ACCRUAL-RECORD

File ACCR-IN. 40 bytes, layout in copybook [ACCRREC](../copybooks/ACCRREC.md).

### ACCR-OUT-RECORD

File ACCR-OUT. 40 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | ACCR-OUT-RECORD | 1 | 40 | X(40) | | | [INTACCR.cbl:62](../../../pangaea/programs/INTACCR.cbl#L62) |

### REPORT-LINE

File REPORT-OUT. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | REPORT-LINE | 1 | 80 | X(80) | | | [INTACCR.cbl:64](../../../pangaea/programs/INTACCR.cbl#L64) |

### WS-FILE-STATUSES

Working storage. 12 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-FILE-STATUSES | 1 | 12 | | | | [INTACCR.cbl:67](../../../pangaea/programs/INTACCR.cbl#L67) |
| 05 | WS-CONTROL-STATUS | 1 | 2 | XX | | | [INTACCR.cbl:68](../../../pangaea/programs/INTACCR.cbl#L68) |
| 05 | WS-RATE-STATUS | 3 | 2 | XX | | | [INTACCR.cbl:69](../../../pangaea/programs/INTACCR.cbl#L69) |
| 05 | WS-ACCT-IN-STATUS | 5 | 2 | XX | | | [INTACCR.cbl:70](../../../pangaea/programs/INTACCR.cbl#L70) |
| 05 | WS-ACCR-IN-STATUS | 7 | 2 | XX | | | [INTACCR.cbl:71](../../../pangaea/programs/INTACCR.cbl#L71) |
| 05 | WS-ACCR-OUT-STATUS | 9 | 2 | XX | | | [INTACCR.cbl:72](../../../pangaea/programs/INTACCR.cbl#L72) |
| 05 | WS-REPORT-STATUS | 11 | 2 | XX | | | [INTACCR.cbl:73](../../../pangaea/programs/INTACCR.cbl#L73) |

### WS-BUSINESS-DATE

Working storage. 8 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-BUSINESS-DATE | 1 | 8 | 9(8) | | | [INTACCR.cbl:75](../../../pangaea/programs/INTACCR.cbl#L75) |

### WS-ABORT-MESSAGE

Working storage. 40 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-ABORT-MESSAGE | 1 | 40 | X(40) | | | [INTACCR.cbl:76](../../../pangaea/programs/INTACCR.cbl#L76) |

### WS-END-SWITCHES

Working storage. 2 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-END-SWITCHES | 1 | 2 | | | | [INTACCR.cbl:77](../../../pangaea/programs/INTACCR.cbl#L77) |
| 05 | WS-MASTER-END-SW | 1 | 1 | X | | VALUE 'N' | [INTACCR.cbl:78](../../../pangaea/programs/INTACCR.cbl#L78) |
| 88 | MASTER-AT-END | | | | | VALUE 'Y' | [INTACCR.cbl:79](../../../pangaea/programs/INTACCR.cbl#L79) |
| 05 | WS-ACCRUAL-END-SW | 2 | 1 | X | | VALUE 'N' | [INTACCR.cbl:80](../../../pangaea/programs/INTACCR.cbl#L80) |
| 88 | ACCRUALS-AT-END | | | | | VALUE 'Y' | [INTACCR.cbl:81](../../../pangaea/programs/INTACCR.cbl#L81) |

### WS-RATE-COUNT

Working storage. 2 bytes. THE RATE TABLE, AS MANY RATES AS THE RATE FILE HOLDS.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-RATE-COUNT | 1 | 2 | 99 | | VALUE 0 | [INTACCR.cbl:84](../../../pangaea/programs/INTACCR.cbl#L84) |

### WS-RATE-TABLE

Working storage. 260 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-RATE-TABLE | 1 | 260 | | | | [INTACCR.cbl:85](../../../pangaea/programs/INTACCR.cbl#L85) |
| 05 | WS-RATE | 1 | 13 | | | OCCURS 1 TO 20 DEPENDING ON WS-RATE-COUNT | [INTACCR.cbl:86](../../../pangaea/programs/INTACCR.cbl#L86) |
| 10 | WS-RT-ACCT-TYPE | 1 | 1 | X | | | [INTACCR.cbl:88](../../../pangaea/programs/INTACCR.cbl#L88) |
| 10 | WS-RT-KIND | 2 | 1 | X | | | [INTACCR.cbl:89](../../../pangaea/programs/INTACCR.cbl#L89) |
| 88 | WS-RT-CREDIT | | | | | VALUE 'C' | [INTACCR.cbl:90](../../../pangaea/programs/INTACCR.cbl#L90) |
| 88 | WS-RT-DEBIT | | | | | VALUE 'D' | [INTACCR.cbl:91](../../../pangaea/programs/INTACCR.cbl#L91) |
| 10 | WS-RT-FLOOR | 3 | 7 | S9(11)V99 | COMP-3 | | [INTACCR.cbl:92](../../../pangaea/programs/INTACCR.cbl#L92) |
| 10 | WS-RT-ANNUAL-RATE | 10 | 4 | 9V9(6) | COMP-3 | | [INTACCR.cbl:93](../../../pangaea/programs/INTACCR.cbl#L93) |

### WS-RX

Working storage. 2 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-RX | 1 | 2 | 99 | | | [INTACCR.cbl:94](../../../pangaea/programs/INTACCR.cbl#L94) |

### WS-ACCOUNT-RATE

Working storage. 4 bytes. THE RATE FOR THE ACCOUNT IN HAND, AND THE TIER IT CAME FROM.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-ACCOUNT-RATE | 1 | 4 | 9V9(6) | COMP-3 | | [INTACCR.cbl:97](../../../pangaea/programs/INTACCR.cbl#L97) |

### WS-RATE-FOUND-SW

Working storage. 1 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-RATE-FOUND-SW | 1 | 1 | X | | | [INTACCR.cbl:98](../../../pangaea/programs/INTACCR.cbl#L98) |
| 88 | RATE-FOUND | | | | | VALUE 'Y' | [INTACCR.cbl:99](../../../pangaea/programs/INTACCR.cbl#L99) |

### WS-BEST-FLOOR

Working storage. 7 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-BEST-FLOOR | 1 | 7 | S9(11)V99 | COMP-3 | | [INTACCR.cbl:100](../../../pangaea/programs/INTACCR.cbl#L100) |

### WS-TODAY

Working storage. 7 bytes. DAYS COUNT FROM 1601-01-01, AS FUNCTION INTEGER-OF-DATE GIVES.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TODAY | 1 | 7 | 9(7) | | | [INTACCR.cbl:103](../../../pangaea/programs/INTACCR.cbl#L103) |

### WS-FROM-DAY

Working storage. 7 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-FROM-DAY | 1 | 7 | 9(7) | | | [INTACCR.cbl:104](../../../pangaea/programs/INTACCR.cbl#L104) |

### WS-MONTH-START-DAY

Working storage. 7 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-MONTH-START-DAY | 1 | 7 | 9(7) | | | [INTACCR.cbl:105](../../../pangaea/programs/INTACCR.cbl#L105) |

### WS-DAYS

Working storage. 3 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-DAYS | 1 | 3 | 9(3) | | | [INTACCR.cbl:106](../../../pangaea/programs/INTACCR.cbl#L106) |

### WS-THROUGH-DATE

Working storage. 8 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-THROUGH-DATE | 1 | 8 | | | | [INTACCR.cbl:107](../../../pangaea/programs/INTACCR.cbl#L107) |
| 05 | WS-TH-YEAR | 1 | 4 | 9(4) | | | [INTACCR.cbl:108](../../../pangaea/programs/INTACCR.cbl#L108) |
| 05 | WS-TH-MONTH | 5 | 2 | 99 | | | [INTACCR.cbl:109](../../../pangaea/programs/INTACCR.cbl#L109) |
| 05 | WS-TH-DAY | 7 | 2 | 99 | | | [INTACCR.cbl:110](../../../pangaea/programs/INTACCR.cbl#L110) |

### WS-MONTH-START

Working storage. 8 bytes. THE FIRST DAY OF THE MONTH AFTER THE LAST ACCRUAL.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-MONTH-START | 1 | 8 | | | | [INTACCR.cbl:112](../../../pangaea/programs/INTACCR.cbl#L112) |
| 05 | WS-MS-YEAR | 1 | 4 | 9(4) | | | [INTACCR.cbl:113](../../../pangaea/programs/INTACCR.cbl#L113) |
| 05 | WS-MS-MONTH | 5 | 2 | 99 | | | [INTACCR.cbl:114](../../../pangaea/programs/INTACCR.cbl#L114) |
| 05 | WS-MS-DAY | 7 | 2 | 99 | | | [INTACCR.cbl:115](../../../pangaea/programs/INTACCR.cbl#L115) |

### WS-MONTH-START-DATE

Working storage. 8 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-MONTH-START-DATE | 1 | 8 | 9(8) | | REDEFINES WS-MONTH-START | [INTACCR.cbl:116](../../../pangaea/programs/INTACCR.cbl#L116) |

### WS-INTEREST-PRODUCT

Working storage. 11 bytes. THE BALANCE TIMES THE RATE TIMES THE DAYS, EXACT, BEFORE IT IS DIVIDED BY THE DAYS OF A YEAR AND ROUNDED.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-INTEREST-PRODUCT | 1 | 11 | S9(13)V9(8) | COMP-3 | | [INTACCR.cbl:121](../../../pangaea/programs/INTACCR.cbl#L121) |

### WS-INCREMENT

Working storage. 7 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-INCREMENT | 1 | 7 | S9(7)V9(6) | COMP-3 | | [INTACCR.cbl:122](../../../pangaea/programs/INTACCR.cbl#L122) |

### WS-DUE-MONTH

Working storage. 6 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-DUE-MONTH | 1 | 6 | | | | [INTACCR.cbl:123](../../../pangaea/programs/INTACCR.cbl#L123) |
| 05 | WS-DM-YEAR | 1 | 4 | 9(4) | | | [INTACCR.cbl:124](../../../pangaea/programs/INTACCR.cbl#L124) |
| 05 | WS-DM-MONTH | 5 | 2 | 99 | | | [INTACCR.cbl:125](../../../pangaea/programs/INTACCR.cbl#L125) |

### WS-TOTALS

Working storage. 44 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TOTALS | 1 | 44 | | | | [INTACCR.cbl:127](../../../pangaea/programs/INTACCR.cbl#L127) |
| 05 | WS-ACCOUNT-COUNT | 1 | 8 | 9(8) | | VALUE 0 | [INTACCR.cbl:128](../../../pangaea/programs/INTACCR.cbl#L128) |
| 05 | WS-CREDIT-ACCRUED | 9 | 8 | S9(9)V9(6) | COMP-3 | VALUE 0 | [INTACCR.cbl:129](../../../pangaea/programs/INTACCR.cbl#L129) |
| 05 | WS-DEBIT-ACCRUED | 17 | 8 | S9(9)V9(6) | COMP-3 | VALUE 0 | [INTACCR.cbl:130](../../../pangaea/programs/INTACCR.cbl#L130) |
| 05 | WS-CREDIT-DUE | 25 | 7 | S9(11)V99 | COMP-3 | VALUE 0 | [INTACCR.cbl:131](../../../pangaea/programs/INTACCR.cbl#L131) |
| 05 | WS-DEBIT-DUE | 32 | 7 | S9(11)V99 | COMP-3 | VALUE 0 | [INTACCR.cbl:132](../../../pangaea/programs/INTACCR.cbl#L132) |
| 05 | WS-CLOSED-MONTH | 39 | 6 | 9(6) | | VALUE 0 | [INTACCR.cbl:133](../../../pangaea/programs/INTACCR.cbl#L133) |

### WS-REPORT-COUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-COUNT-LINE | 1 | 80 | | | | [INTACCR.cbl:135](../../../pangaea/programs/INTACCR.cbl#L135) |
| 05 | WS-RC-LABEL | 1 | 30 | X(30) | | | [INTACCR.cbl:136](../../../pangaea/programs/INTACCR.cbl#L136) |
| 05 | WS-RC-VALUE | 31 | 10 | ZZ,ZZZ,ZZ9 | | | [INTACCR.cbl:137](../../../pangaea/programs/INTACCR.cbl#L137) |
| 05 | FILLER | 41 | 40 | X(40) | | VALUE SPACES | [INTACCR.cbl:138](../../../pangaea/programs/INTACCR.cbl#L138) |

### WS-REPORT-AMOUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-AMOUNT-LINE | 1 | 80 | | | | [INTACCR.cbl:139](../../../pangaea/programs/INTACCR.cbl#L139) |
| 05 | WS-RA-LABEL | 1 | 30 | X(30) | | | [INTACCR.cbl:140](../../../pangaea/programs/INTACCR.cbl#L140) |
| 05 | WS-RA-VALUE | 31 | 19 | ZZZ,ZZZ,ZZ9.999999- | | | [INTACCR.cbl:141](../../../pangaea/programs/INTACCR.cbl#L141) |
| 05 | FILLER | 50 | 31 | X(31) | | VALUE SPACES | [INTACCR.cbl:142](../../../pangaea/programs/INTACCR.cbl#L142) |
