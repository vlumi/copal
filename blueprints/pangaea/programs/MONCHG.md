# MONCHG

```text
MONCHG: MONTH-END INTEREST AND FEES

ON THE DAY INTACCR CLOSES A MONTH, WRITES THE MONTH'S CHARGES
AND CREDITS AS TRANSACTIONS FOR TXPOST: THE INTEREST DUE,
CREDITED (IN) OR CHARGED AS OVERDRAFT INTEREST (OD), AND THE
MAINTENANCE FEE (MF) OF AN OPEN CURRENT ACCOUNT WHOSE BALANCE
CLOSED THE MONTH BELOW THE WAIVER BALANCE. ON OTHER DAYS IT
WRITES NOTHING.

THE MASTER AND THE ACCRUALS ARE READ IN STEP, SO THE
TRANSACTIONS COME OUT IN ACCOUNT NUMBER ORDER, AND TXPOST
POSTS THEM WITHOUT A SORT. TRANSACTION IDS ARE 00, THE CLOSED
MONTH AS YYMM, AND A SEQUENCE NUMBER.

RETURN CODES: 0 WRITTEN, 16 FILES OUT OF STEP OR I/O.

MODIFICATION HISTORY
DATE       BY  CHANGE
2001-09-03 TL  CR-0298 ORIGINAL VERSION.
```

Program [MONCHG.cbl](../../../pangaea/programs/MONCHG.cbl): 5 files, 10 paragraphs. Its business rules: [6 decision tables](../rules/MONCHG.md).

- AUTHOR: T. LINDQVIST, BATCH SYSTEMS.
- INSTALLATION: FIRST BANK OF PANGAEA, DATA CENTER.
- DATE-WRITTEN: 2001-09-03.

## Contents

- [Files](#files)
- [Flow](#flow)
- [Paragraphs](#paragraphs)
- [Data flow](#data-flow)
- [Lineage](#lineage)
- [Records](#records)
  - [BANK-CONTROL-RECORD](#bank-control-record)
  - [ACCT-MASTER-RECORD](#acct-master-record)
  - [ACCRUAL-RECORD](#accrual-record)
  - [TRAN-OUT-RECORD](#tran-out-record)
  - [REPORT-LINE](#report-line)
  - [WS-FILE-STATUSES](#ws-file-statuses)
  - [WS-BUSINESS-DATE](#ws-business-date)
  - [WS-ABORT-MESSAGE](#ws-abort-message)
  - [WS-END-SWITCHES](#ws-end-switches)
  - [WS-MAINTENANCE-FEE](#ws-maintenance-fee)
  - [WS-FEE-WAIVER-BALANCE](#ws-fee-waiver-balance)
  - [WS-CLOSED-MONTH](#ws-closed-month)
  - [WS-CLOSED-MONTH-PARTS](#ws-closed-month-parts)
  - [WS-TRAN-ID](#ws-tran-id)
  - [TRAN-RECORD](#tran-record)
  - [WS-TOTALS](#ws-totals)
  - [WS-REPORT-COUNT-LINE](#ws-report-count-line)
  - [WS-REPORT-AMOUNT-LINE](#ws-report-amount-line)

## Files

| File | Kind | Assigned to | Organization | Status field | Opened | Records |
| --- | --- | --- | --- | --- | --- | --- |
| CONTROL-FILE | FD | BANKCTL | SEQUENTIAL | WS-CONTROL-STATUS | INPUT | BANK-CONTROL-RECORD |
| ACCT-IN | FD | ACCTIN | SEQUENTIAL | WS-ACCT-IN-STATUS | INPUT | ACCT-MASTER-RECORD |
| ACCR-IN | FD | ACCRIN | SEQUENTIAL | WS-ACCR-IN-STATUS | INPUT | ACCRUAL-RECORD |
| TRAN-OUT | FD | TRANOUT | SEQUENTIAL | WS-TRAN-OUT-STATUS | OUTPUT | TRAN-OUT-RECORD |
| REPORT-OUT | FD | MONRPT | LINE SEQUENTIAL | WS-REPORT-STATUS | OUTPUT | REPORT-LINE |

## Flow

![The paragraphs of MONCHG and how control moves between them](MONCHG.svg)

Solid arrows are PERFORMs, dashed ones GO TOs, and dotted ones control running on into the next paragraph inside a PERFORM THRU range.

## Paragraphs

| Paragraph | Description | Source | Performs | Goes to | Performed by | Reached by |
| --- | --- | --- | --- | --- | --- | --- |
| 0000-MAIN | CHARGE EVERY ACCOUNT, THEN REPORT. | [MONCHG.cbl:111](../../../pangaea/programs/MONCHG.cbl#L111) | 1000-INITIALIZE, 2000-CHARGE-ACCOUNT, 8000-WRITE-REPORT, 9000-CLOSE-FILES | | | |
| 1000-INITIALIZE | OPEN THE FILES AND READ THE FIRST ACCOUNT. THE TRANSACTION RECORD STARTS BLANK: NO COUNTER ACCOUNT. | [MONCHG.cbl:122](../../../pangaea/programs/MONCHG.cbl#L122) | 9900-ABORT, 1200-READ-ACCOUNT | | 0000-MAIN | |
| 1200-READ-ACCOUNT | READ THE NEXT ACCOUNT AND ITS ACCRUAL, WHICH MUST MATCH. | [MONCHG.cbl:140](../../../pangaea/programs/MONCHG.cbl#L140) | 9900-ABORT | | 1000-INITIALIZE, 2000-CHARGE-ACCOUNT | |
| 2000-CHARGE-ACCOUNT | WRITE AN ACCOUNT'S INTEREST AND FEE, IF ITS MONTH CLOSED. | [MONCHG.cbl:161](../../../pangaea/programs/MONCHG.cbl#L161) | 2100-SETTLE-INTEREST, 2200-CHARGE-FEE, 1200-READ-ACCOUNT | | 0000-MAIN | |
| 2100-SETTLE-INTEREST | CREDIT THE INTEREST DUE, OR CHARGE IT IF THE ACCOUNT WAS OVERDRAWN. | [MONCHG.cbl:172](../../../pangaea/programs/MONCHG.cbl#L172) | 2900-WRITE-TRAN | | 2000-CHARGE-ACCOUNT | |
| 2200-CHARGE-FEE | CHARGE THE MAINTENANCE FEE: OPEN CURRENT ACCOUNTS ONLY, AND NOT IF THE MONTH CLOSED AT OR ABOVE THE WAIVER BALANCE. | [MONCHG.cbl:194](../../../pangaea/programs/MONCHG.cbl#L194) | 2900-WRITE-TRAN | | 2000-CHARGE-ACCOUNT | |
| 2900-WRITE-TRAN | WRITE THE TRANSACTION FOR THE ACCOUNT, WITH THE NEXT ID. | [MONCHG.cbl:212](../../../pangaea/programs/MONCHG.cbl#L212) | | | 2100-SETTLE-INTEREST, 2200-CHARGE-FEE | |
| 8000-WRITE-REPORT | WRITE THE MONTH-END CHARGES REPORT. | [MONCHG.cbl:220](../../../pangaea/programs/MONCHG.cbl#L220) | | | 0000-MAIN | |
| 9000-CLOSE-FILES | CLOSE ALL FILES. | [MONCHG.cbl:262](../../../pangaea/programs/MONCHG.cbl#L262) | | | 0000-MAIN, 9900-ABORT | |
| 9900-ABORT | STOP THE RUN WITH RETURN CODE 16, SAYING WHY. | [MONCHG.cbl:266](../../../pangaea/programs/MONCHG.cbl#L266) | 9000-CLOSE-FILES | | 1000-INITIALIZE, 1200-READ-ACCOUNT | |

## Data flow

| Paragraph | Reads | Writes |
| --- | --- | --- |
| 0000-MAIN | WS-MASTER-END-SW, WS-ACCRUAL-END-SW | RETURN-CODE |
| 1000-INITIALIZE | WS-FILE-STATUSES, BC-BUSINESS-DATE | RETURN-CODE, BANK-CONTROL-RECORD, WS-ABORT-MESSAGE, WS-BUSINESS-DATE, TRAN-RECORD |
| 1200-READ-ACCOUNT | WS-MASTER-END-SW, WS-ACCRUAL-END-SW, AC-ACCT-NO, AM-ACCT-NO | ACCT-MASTER-RECORD, WS-MASTER-END-SW, ACCRUAL-RECORD, WS-ACCRUAL-END-SW, WS-ABORT-MESSAGE |
| 2000-CHARGE-ACCOUNT | AC-DUE-MONTH, WS-CM-YYMM | WS-CLOSED-MONTH, WS-TI-YYMM |
| 2100-SETTLE-INTEREST | AC-DUE, WS-INTEREST-COUNT, TD-AMOUNT, WS-INTEREST-TOTAL, WS-OVERDRAFT-COUNT, WS-OVERDRAFT-TOTAL | TD-TRAN-CODE, TD-AMOUNT, TD-DESCRIPTION, WS-INTEREST-COUNT, WS-INTEREST-TOTAL, WS-OVERDRAFT-COUNT, WS-OVERDRAFT-TOTAL |
| 2200-CHARGE-FEE | AM-ACCT-TYPE, AM-STATUS, AM-BALANCE, WS-FEE-WAIVER-BALANCE, WS-MAINTENANCE-FEE, WS-FEE-COUNT, TD-AMOUNT, WS-FEE-TOTAL, WS-WAIVED-COUNT | TD-TRAN-CODE, TD-AMOUNT, TD-DESCRIPTION, WS-FEE-COUNT, WS-FEE-TOTAL, WS-WAIVED-COUNT |
| 2900-WRITE-TRAN | WS-TI-SEQ, WS-TRAN-ID, AM-ACCT-NO, TRAN-RECORD | TR-REC-TYPE, WS-TI-SEQ, TD-TRAN-ID, TD-ACCT-NO, TRAN-OUT-RECORD |
| 8000-WRITE-REPORT | WS-BUSINESS-DATE, WS-CLOSED-MONTH, WS-INTEREST-COUNT, WS-REPORT-COUNT-LINE, WS-INTEREST-TOTAL, WS-REPORT-AMOUNT-LINE, WS-OVERDRAFT-COUNT, WS-OVERDRAFT-TOTAL, WS-FEE-COUNT, WS-FEE-TOTAL, WS-WAIVED-COUNT | REPORT-LINE, WS-RC-LABEL, WS-RC-VALUE, WS-RA-LABEL, WS-RA-VALUE |
| 9000-CLOSE-FILES | | |
| 9900-ABORT | WS-ABORT-MESSAGE | RETURN-CODE |

## Lineage

Where each data item's values come from, and which records go to which files.

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| ACCRUAL-RECORD | file ACCR-IN | READ | 1200-READ-ACCOUNT | [MONCHG.cbl:145](../../../pangaea/programs/MONCHG.cbl#L145) |
| ACCT-MASTER-RECORD | file ACCT-IN | READ | 1200-READ-ACCOUNT | [MONCHG.cbl:141](../../../pangaea/programs/MONCHG.cbl#L141) |
| BANK-CONTROL-RECORD | file CONTROL-FILE | READ | 1000-INITIALIZE | [MONCHG.cbl:130](../../../pangaea/programs/MONCHG.cbl#L130) |
| REPORT-LINE | WS-BUSINESS-DATE | STRING | 8000-WRITE-REPORT | [MONCHG.cbl:222](../../../pangaea/programs/MONCHG.cbl#L222) |
| REPORT-LINE | WS-CLOSED-MONTH | STRING | 8000-WRITE-REPORT | [MONCHG.cbl:233](../../../pangaea/programs/MONCHG.cbl#L233) |
| REPORT-LINE | WS-REPORT-COUNT-LINE | WRITE | 8000-WRITE-REPORT | [MONCHG.cbl:240](../../../pangaea/programs/MONCHG.cbl#L240) |
| REPORT-LINE | WS-REPORT-AMOUNT-LINE | WRITE | 8000-WRITE-REPORT | [MONCHG.cbl:243](../../../pangaea/programs/MONCHG.cbl#L243) |
| REPORT-LINE | WS-REPORT-COUNT-LINE | WRITE | 8000-WRITE-REPORT | [MONCHG.cbl:246](../../../pangaea/programs/MONCHG.cbl#L246) |
| REPORT-LINE | WS-REPORT-AMOUNT-LINE | WRITE | 8000-WRITE-REPORT | [MONCHG.cbl:249](../../../pangaea/programs/MONCHG.cbl#L249) |
| REPORT-LINE | WS-REPORT-COUNT-LINE | WRITE | 8000-WRITE-REPORT | [MONCHG.cbl:252](../../../pangaea/programs/MONCHG.cbl#L252) |
| REPORT-LINE | WS-REPORT-AMOUNT-LINE | WRITE | 8000-WRITE-REPORT | [MONCHG.cbl:255](../../../pangaea/programs/MONCHG.cbl#L255) |
| REPORT-LINE | WS-REPORT-COUNT-LINE | WRITE | 8000-WRITE-REPORT | [MONCHG.cbl:258](../../../pangaea/programs/MONCHG.cbl#L258) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [MONCHG.cbl:225](../../../pangaea/programs/MONCHG.cbl#L225) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [MONCHG.cbl:227](../../../pangaea/programs/MONCHG.cbl#L227) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [MONCHG.cbl:230](../../../pangaea/programs/MONCHG.cbl#L230) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [MONCHG.cbl:235](../../../pangaea/programs/MONCHG.cbl#L235) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [MONCHG.cbl:237](../../../pangaea/programs/MONCHG.cbl#L237) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [MONCHG.cbl:240](../../../pangaea/programs/MONCHG.cbl#L240) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [MONCHG.cbl:243](../../../pangaea/programs/MONCHG.cbl#L243) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [MONCHG.cbl:246](../../../pangaea/programs/MONCHG.cbl#L246) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [MONCHG.cbl:249](../../../pangaea/programs/MONCHG.cbl#L249) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [MONCHG.cbl:252](../../../pangaea/programs/MONCHG.cbl#L252) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [MONCHG.cbl:255](../../../pangaea/programs/MONCHG.cbl#L255) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [MONCHG.cbl:258](../../../pangaea/programs/MONCHG.cbl#L258) |
| TD-ACCT-NO | AM-ACCT-NO | MOVE | 2900-WRITE-TRAN | [MONCHG.cbl:216](../../../pangaea/programs/MONCHG.cbl#L216) |
| TD-AMOUNT | AC-DUE | MOVE | 2100-SETTLE-INTEREST | [MONCHG.cbl:176](../../../pangaea/programs/MONCHG.cbl#L176) |
| TD-AMOUNT | AC-DUE | COMPUTE | 2100-SETTLE-INTEREST | [MONCHG.cbl:183](../../../pangaea/programs/MONCHG.cbl#L183) |
| TD-AMOUNT | WS-MAINTENANCE-FEE | MOVE | 2200-CHARGE-FEE | [MONCHG.cbl:199](../../../pangaea/programs/MONCHG.cbl#L199) |
| TD-TRAN-ID | WS-TRAN-ID | MOVE | 2900-WRITE-TRAN | [MONCHG.cbl:215](../../../pangaea/programs/MONCHG.cbl#L215) |
| file TRAN-OUT | TRAN-OUT-RECORD | WRITE | 2900-WRITE-TRAN | [MONCHG.cbl:217](../../../pangaea/programs/MONCHG.cbl#L217) |
| TRAN-OUT-RECORD | TRAN-RECORD | WRITE | 2900-WRITE-TRAN | [MONCHG.cbl:217](../../../pangaea/programs/MONCHG.cbl#L217) |
| WS-BUSINESS-DATE | BC-BUSINESS-DATE | MOVE | 1000-INITIALIZE | [MONCHG.cbl:135](../../../pangaea/programs/MONCHG.cbl#L135) |
| WS-CLOSED-MONTH | AC-DUE-MONTH | MOVE | 2000-CHARGE-ACCOUNT | [MONCHG.cbl:163](../../../pangaea/programs/MONCHG.cbl#L163) |
| WS-FEE-TOTAL | TD-AMOUNT | ADD | 2200-CHARGE-FEE | [MONCHG.cbl:203](../../../pangaea/programs/MONCHG.cbl#L203) |
| WS-INTEREST-TOTAL | TD-AMOUNT | ADD | 2100-SETTLE-INTEREST | [MONCHG.cbl:180](../../../pangaea/programs/MONCHG.cbl#L180) |
| WS-OVERDRAFT-TOTAL | TD-AMOUNT | ADD | 2100-SETTLE-INTEREST | [MONCHG.cbl:187](../../../pangaea/programs/MONCHG.cbl#L187) |
| WS-RA-VALUE | WS-INTEREST-TOTAL | MOVE | 8000-WRITE-REPORT | [MONCHG.cbl:242](../../../pangaea/programs/MONCHG.cbl#L242) |
| WS-RA-VALUE | WS-OVERDRAFT-TOTAL | MOVE | 8000-WRITE-REPORT | [MONCHG.cbl:248](../../../pangaea/programs/MONCHG.cbl#L248) |
| WS-RA-VALUE | WS-FEE-TOTAL | MOVE | 8000-WRITE-REPORT | [MONCHG.cbl:254](../../../pangaea/programs/MONCHG.cbl#L254) |
| WS-RC-VALUE | WS-INTEREST-COUNT | MOVE | 8000-WRITE-REPORT | [MONCHG.cbl:239](../../../pangaea/programs/MONCHG.cbl#L239) |
| WS-RC-VALUE | WS-OVERDRAFT-COUNT | MOVE | 8000-WRITE-REPORT | [MONCHG.cbl:245](../../../pangaea/programs/MONCHG.cbl#L245) |
| WS-RC-VALUE | WS-FEE-COUNT | MOVE | 8000-WRITE-REPORT | [MONCHG.cbl:251](../../../pangaea/programs/MONCHG.cbl#L251) |
| WS-RC-VALUE | WS-WAIVED-COUNT | MOVE | 8000-WRITE-REPORT | [MONCHG.cbl:257](../../../pangaea/programs/MONCHG.cbl#L257) |
| WS-TI-YYMM | WS-CM-YYMM | MOVE | 2000-CHARGE-ACCOUNT | [MONCHG.cbl:164](../../../pangaea/programs/MONCHG.cbl#L164) |

## Records

### BANK-CONTROL-RECORD

File CONTROL-FILE. 80 bytes, layout in copybook [BANKCTL](../copybooks/BANKCTL.md).

### ACCT-MASTER-RECORD

File ACCT-IN. 100 bytes, layout in copybook [ACCTMAST](../copybooks/ACCTMAST.md).

### ACCRUAL-RECORD

File ACCR-IN. 40 bytes, layout in copybook [ACCRREC](../copybooks/ACCRREC.md).

### TRAN-OUT-RECORD

File TRAN-OUT. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | TRAN-OUT-RECORD | 1 | 80 | X(80) | | | [MONCHG.cbl:56](../../../pangaea/programs/MONCHG.cbl#L56) |

### REPORT-LINE

File REPORT-OUT. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | REPORT-LINE | 1 | 80 | X(80) | | | [MONCHG.cbl:58](../../../pangaea/programs/MONCHG.cbl#L58) |

### WS-FILE-STATUSES

Working storage. 10 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-FILE-STATUSES | 1 | 10 | | | | [MONCHG.cbl:61](../../../pangaea/programs/MONCHG.cbl#L61) |
| 05 | WS-CONTROL-STATUS | 1 | 2 | XX | | | [MONCHG.cbl:62](../../../pangaea/programs/MONCHG.cbl#L62) |
| 05 | WS-ACCT-IN-STATUS | 3 | 2 | XX | | | [MONCHG.cbl:63](../../../pangaea/programs/MONCHG.cbl#L63) |
| 05 | WS-ACCR-IN-STATUS | 5 | 2 | XX | | | [MONCHG.cbl:64](../../../pangaea/programs/MONCHG.cbl#L64) |
| 05 | WS-TRAN-OUT-STATUS | 7 | 2 | XX | | | [MONCHG.cbl:65](../../../pangaea/programs/MONCHG.cbl#L65) |
| 05 | WS-REPORT-STATUS | 9 | 2 | XX | | | [MONCHG.cbl:66](../../../pangaea/programs/MONCHG.cbl#L66) |

### WS-BUSINESS-DATE

Working storage. 8 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-BUSINESS-DATE | 1 | 8 | 9(8) | | | [MONCHG.cbl:68](../../../pangaea/programs/MONCHG.cbl#L68) |

### WS-ABORT-MESSAGE

Working storage. 40 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-ABORT-MESSAGE | 1 | 40 | X(40) | | | [MONCHG.cbl:69](../../../pangaea/programs/MONCHG.cbl#L69) |

### WS-END-SWITCHES

Working storage. 2 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-END-SWITCHES | 1 | 2 | | | | [MONCHG.cbl:70](../../../pangaea/programs/MONCHG.cbl#L70) |
| 05 | WS-MASTER-END-SW | 1 | 1 | X | | VALUE 'N' | [MONCHG.cbl:71](../../../pangaea/programs/MONCHG.cbl#L71) |
| 88 | MASTER-AT-END | | | | | VALUE 'Y' | [MONCHG.cbl:72](../../../pangaea/programs/MONCHG.cbl#L72) |
| 05 | WS-ACCRUAL-END-SW | 2 | 1 | X | | VALUE 'N' | [MONCHG.cbl:73](../../../pangaea/programs/MONCHG.cbl#L73) |
| 88 | ACCRUALS-AT-END | | | | | VALUE 'Y' | [MONCHG.cbl:74](../../../pangaea/programs/MONCHG.cbl#L74) |

### WS-MAINTENANCE-FEE

Working storage. 5 bytes. THE MAINTENANCE FEE, AND THE BALANCE THAT WAIVES IT.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-MAINTENANCE-FEE | 1 | 5 | 9(3)V99 | | VALUE 2.50 | [MONCHG.cbl:77](../../../pangaea/programs/MONCHG.cbl#L77) |

### WS-FEE-WAIVER-BALANCE

Working storage. 9 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-FEE-WAIVER-BALANCE | 1 | 9 | 9(7)V99 | | VALUE 1000.00 | [MONCHG.cbl:78](../../../pangaea/programs/MONCHG.cbl#L78) |

### WS-CLOSED-MONTH

Working storage. 6 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-CLOSED-MONTH | 1 | 6 | 9(6) | | VALUE 0 | [MONCHG.cbl:80](../../../pangaea/programs/MONCHG.cbl#L80) |

### WS-CLOSED-MONTH-PARTS

Working storage. 6 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-CLOSED-MONTH-PARTS | 1 | 6 | | | REDEFINES WS-CLOSED-MONTH | [MONCHG.cbl:81](../../../pangaea/programs/MONCHG.cbl#L81) |
| 05 | WS-CM-CENTURY | 1 | 2 | 99 | | | [MONCHG.cbl:82](../../../pangaea/programs/MONCHG.cbl#L82) |
| 05 | WS-CM-YYMM | 3 | 4 | 9(4) | | | [MONCHG.cbl:83](../../../pangaea/programs/MONCHG.cbl#L83) |

### WS-TRAN-ID

Working storage. 10 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TRAN-ID | 1 | 10 | | | | [MONCHG.cbl:84](../../../pangaea/programs/MONCHG.cbl#L84) |
| 05 | FILLER | 1 | 2 | 99 | | VALUE 0 | [MONCHG.cbl:85](../../../pangaea/programs/MONCHG.cbl#L85) |
| 05 | WS-TI-YYMM | 3 | 4 | 9(4) | | | [MONCHG.cbl:86](../../../pangaea/programs/MONCHG.cbl#L86) |
| 05 | WS-TI-SEQ | 7 | 4 | 9(4) | | VALUE 0 | [MONCHG.cbl:87](../../../pangaea/programs/MONCHG.cbl#L87) |

### TRAN-RECORD

Working storage. 80 bytes, layout in copybook [TRANREC](../copybooks/TRANREC.md).

### WS-TOTALS

Working storage. 71 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TOTALS | 1 | 71 | | | | [MONCHG.cbl:91](../../../pangaea/programs/MONCHG.cbl#L91) |
| 05 | WS-INTEREST-COUNT | 1 | 8 | 9(8) | | VALUE 0 | [MONCHG.cbl:92](../../../pangaea/programs/MONCHG.cbl#L92) |
| 05 | WS-INTEREST-TOTAL | 9 | 13 | 9(11)V99 | | VALUE 0 | [MONCHG.cbl:93](../../../pangaea/programs/MONCHG.cbl#L93) |
| 05 | WS-OVERDRAFT-COUNT | 22 | 8 | 9(8) | | VALUE 0 | [MONCHG.cbl:94](../../../pangaea/programs/MONCHG.cbl#L94) |
| 05 | WS-OVERDRAFT-TOTAL | 30 | 13 | 9(11)V99 | | VALUE 0 | [MONCHG.cbl:95](../../../pangaea/programs/MONCHG.cbl#L95) |
| 05 | WS-FEE-COUNT | 43 | 8 | 9(8) | | VALUE 0 | [MONCHG.cbl:96](../../../pangaea/programs/MONCHG.cbl#L96) |
| 05 | WS-FEE-TOTAL | 51 | 13 | 9(11)V99 | | VALUE 0 | [MONCHG.cbl:97](../../../pangaea/programs/MONCHG.cbl#L97) |
| 05 | WS-WAIVED-COUNT | 64 | 8 | 9(8) | | VALUE 0 | [MONCHG.cbl:98](../../../pangaea/programs/MONCHG.cbl#L98) |

### WS-REPORT-COUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-COUNT-LINE | 1 | 80 | | | | [MONCHG.cbl:100](../../../pangaea/programs/MONCHG.cbl#L100) |
| 05 | WS-RC-LABEL | 1 | 30 | X(30) | | | [MONCHG.cbl:101](../../../pangaea/programs/MONCHG.cbl#L101) |
| 05 | WS-RC-VALUE | 31 | 10 | ZZ,ZZZ,ZZ9 | | | [MONCHG.cbl:102](../../../pangaea/programs/MONCHG.cbl#L102) |
| 05 | FILLER | 41 | 40 | X(40) | | VALUE SPACES | [MONCHG.cbl:103](../../../pangaea/programs/MONCHG.cbl#L103) |

### WS-REPORT-AMOUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-AMOUNT-LINE | 1 | 80 | | | | [MONCHG.cbl:104](../../../pangaea/programs/MONCHG.cbl#L104) |
| 05 | WS-RA-LABEL | 1 | 30 | X(30) | | | [MONCHG.cbl:105](../../../pangaea/programs/MONCHG.cbl#L105) |
| 05 | WS-RA-VALUE | 31 | 21 | Z,ZZZ,ZZZ,ZZZ,ZZ9.99- | | | [MONCHG.cbl:106](../../../pangaea/programs/MONCHG.cbl#L106) |
| 05 | FILLER | 52 | 29 | X(29) | | VALUE SPACES | [MONCHG.cbl:107](../../../pangaea/programs/MONCHG.cbl#L107) |
