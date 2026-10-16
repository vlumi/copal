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
IN STEP WITH THE MASTER, ONE FOR EACH ACCOUNT. AN ACCOUNT
OPENED SINCE THE LAST RUN HAS NONE YET AND STARTS ACCRUING
TODAY.

RETURN CODES: 0 ACCRUED,
              16 AN ACCRUAL FOR NO ACCOUNT, NO RATE, OR I/O.

MODIFICATION HISTORY
DATE       BY  CHANGE
2001-09-03 TL  CR-0298 ORIGINAL VERSION.
2002-02-11 TL  CR-0311 ACCOUNTS OPENED SINCE THE LAST RUN.
```

Program [INTACCR.cbl](../../../pangaea/programs/INTACCR.cbl): 6 files, 15 paragraphs.

- AUTHOR: T. LINDQVIST, BATCH SYSTEMS.
- INSTALLATION: FIRST BANK OF PANGAEA, DATA CENTER.
- DATE-WRITTEN: 2001-09-03.

## Files

| File | Kind | Assigned to | Organization | Status field | Opened | Records |
| --- | --- | --- | --- | --- | --- | --- |
| CONTROL-FILE | FD | BANKCTL | SEQUENTIAL | WS-CONTROL-STATUS | INPUT | BANK-CONTROL-RECORD |
| RATE-FILE | FD | RATES | SEQUENTIAL | WS-RATE-STATUS | INPUT | RATE-RECORD |
| ACCT-IN | FD | ACCTIN | SEQUENTIAL | WS-ACCT-IN-STATUS | INPUT | ACCT-MASTER-RECORD |
| ACCR-IN | FD | ACCRIN | SEQUENTIAL | WS-ACCR-IN-STATUS | INPUT | ACCR-IN-RECORD |
| ACCR-OUT | FD | ACCROUT | SEQUENTIAL | WS-ACCR-OUT-STATUS | OUTPUT | ACCR-OUT-RECORD |
| REPORT-OUT | FD | ACCRRPT | LINE SEQUENTIAL | WS-REPORT-STATUS | OUTPUT | REPORT-LINE |

## Flow

![The paragraphs of INTACCR and how control moves between them](INTACCR.svg)

Solid arrows are PERFORMs, dashed ones GO TOs, and dotted ones control running on into the next paragraph inside a PERFORM THRU range.

## Paragraphs

| Paragraph | Description | Source | Performs | Goes to | Performed by | Reached by |
| --- | --- | --- | --- | --- | --- | --- |
| 0000-MAIN | ACCRUE EVERY ACCOUNT, THEN REPORT. | [INTACCR.cbl:155](../../../pangaea/programs/INTACCR.cbl#L155) | 1000-INITIALIZE, 2000-ACCRUE-ACCOUNT, 9900-ABORT, 8000-WRITE-REPORT, 9000-CLOSE-FILES | | | |
| 1000-INITIALIZE | OPEN THE FILES, LOAD THE RATES AND READ THE FIRST ACCOUNT. | [INTACCR.cbl:169](../../../pangaea/programs/INTACCR.cbl#L169) | 9900-ABORT, 1100-LOAD-RATES, 1200-READ-MASTER, 1300-READ-ACCRUAL | | 0000-MAIN | |
| 1100-LOAD-RATES | READ THE RATE FILE INTO THE RATE TABLE. | [INTACCR.cbl:189](../../../pangaea/programs/INTACCR.cbl#L189) | 9900-ABORT | | 1000-INITIALIZE | |
| 1200-READ-MASTER | READ THE NEXT ACCOUNT. | [INTACCR.cbl:209](../../../pangaea/programs/INTACCR.cbl#L209) | | | 1000-INITIALIZE, 2000-ACCRUE-ACCOUNT | |
| 1300-READ-ACCRUAL | READ THE NEXT ACCRUAL. | [INTACCR.cbl:216](../../../pangaea/programs/INTACCR.cbl#L216) | | | 1000-INITIALIZE, 1400-TAKE-ACCRUAL | |
| 1400-TAKE-ACCRUAL | TAKE THE ACCOUNT'S ACCRUAL, OR START ONE FOR AN ACCOUNT OPENED SINCE THE LAST RUN: NOTHING ACCRUED, THROUGH TODAY. | [INTACCR.cbl:226](../../../pangaea/programs/INTACCR.cbl#L226) | 1300-READ-ACCRUAL, 9900-ABORT | | 2000-ACCRUE-ACCOUNT | |
| 2000-ACCRUE-ACCOUNT | ACCRUE ONE ACCOUNT UP TO THE BUSINESS DATE. IF A MONTH CLOSED ON THE WAY, ACCRUE ITS DAYS FIRST AND MAKE ITS INTEREST DUE. | [INTACCR.cbl:244](../../../pangaea/programs/INTACCR.cbl#L244) | 1400-TAKE-ACCRUAL, 2100-FIND-RATE, 2200-FIND-MONTH-START, 2300-ACCRUE-DAYS, 2400-CLOSE-MONTH, 1200-READ-MASTER | | 0000-MAIN | |
| 2100-FIND-RATE | FIND THE ACCOUNT'S RATE: BELOW ZERO THE DEBIT RATE OF ITS TYPE, OTHERWISE THE CREDIT RATE OF THE HIGHEST TIER IT REACHES. | [INTACCR.cbl:266](../../../pangaea/programs/INTACCR.cbl#L266) | 2110-TAKE-RATE, 9900-ABORT | | 2000-ACCRUE-ACCOUNT | |
| 2110-TAKE-RATE | TAKE THE RATE IN HAND AS THE ACCOUNT'S, SO FAR. | [INTACCR.cbl:291](../../../pangaea/programs/INTACCR.cbl#L291) | | | 2100-FIND-RATE | |
| 2200-FIND-MONTH-START | THE FIRST DAY OF THE MONTH AFTER THE THROUGH DATE. | [INTACCR.cbl:297](../../../pangaea/programs/INTACCR.cbl#L297) | | | 2000-ACCRUE-ACCOUNT | |
| 2300-ACCRUE-DAYS | ACCRUE THE OPENING BALANCE AT THE ACCOUNT'S RATE FOR THE DAYS. | [INTACCR.cbl:311](../../../pangaea/programs/INTACCR.cbl#L311) | | | 2000-ACCRUE-ACCOUNT | |
| 2400-CLOSE-MONTH | MAKE THE CLOSED MONTH'S INTEREST DUE, TO THE CENT, AND KEEP THE FRACTION ACCRUED. | [INTACCR.cbl:326](../../../pangaea/programs/INTACCR.cbl#L326) | | | 2000-ACCRUE-ACCOUNT | |
| 8000-WRITE-REPORT | WRITE THE ACCRUAL REPORT. | [INTACCR.cbl:340](../../../pangaea/programs/INTACCR.cbl#L340) | | | 0000-MAIN | |
| 9000-CLOSE-FILES | CLOSE ALL FILES. | [INTACCR.cbl:382](../../../pangaea/programs/INTACCR.cbl#L382) | | | 0000-MAIN, 9900-ABORT | |
| 9900-ABORT | STOP THE RUN WITH RETURN CODE 16, SAYING WHY. | [INTACCR.cbl:387](../../../pangaea/programs/INTACCR.cbl#L387) | 9000-CLOSE-FILES | | 0000-MAIN, 1000-INITIALIZE, 1100-LOAD-RATES, 1400-TAKE-ACCRUAL, 2100-FIND-RATE | |

## Data flow

| Paragraph | Reads | Writes |
| --- | --- | --- |
| 0000-MAIN | WS-MASTER-END-SW, WS-ACCRUAL-KEY | WS-ABORT-MESSAGE, RETURN-CODE |
| 1000-INITIALIZE | WS-FILE-STATUSES, BC-BUSINESS-DATE, WS-BUSINESS-DATE | RETURN-CODE, BANK-CONTROL-RECORD, WS-ABORT-MESSAGE, WS-BUSINESS-DATE, WS-TODAY |
| 1100-LOAD-RATES | WS-RATE-STATUS, WS-RATE-COUNT, RT-ACCT-TYPE, RT-KIND, RT-FLOOR, RT-ANNUAL-RATE | RATE-RECORD, WS-ABORT-MESSAGE, WS-RATE-COUNT, WS-RT-ACCT-TYPE, WS-RT-KIND, WS-RT-FLOOR, WS-RT-ANNUAL-RATE |
| 1200-READ-MASTER | | ACCT-MASTER-RECORD, WS-MASTER-END-SW |
| 1300-READ-ACCRUAL | AI-ACCT-NO | ACCR-IN-RECORD, WS-ACCRUAL-KEY |
| 1400-TAKE-ACCRUAL | WS-ACCRUAL-KEY, AM-ACCT-NO, ACCR-IN-RECORD, WS-BUSINESS-DATE, WS-STARTED-COUNT | ACCRUAL-RECORD, AC-ACCT-NO, AC-THROUGH-DATE, AC-ACCRUED, WS-STARTED-COUNT, WS-ABORT-MESSAGE |
| 2000-ACCRUE-ACCOUNT | AC-THROUGH-DATE, WS-MONTH-START-DAY, WS-TODAY, WS-FROM-DAY, WS-BUSINESS-DATE, ACCRUAL-RECORD, WS-ACCOUNT-COUNT | AC-DUE, AC-DUE-MONTH, WS-FROM-DAY, WS-DAYS, AC-THROUGH-DATE, ACCR-OUT-RECORD, WS-ACCOUNT-COUNT |
| 2100-FIND-RATE | WS-RX, WS-RATE-COUNT, WS-RT-ACCT-TYPE, AM-ACCT-TYPE, AM-BALANCE, WS-RT-KIND, WS-RT-FLOOR, WS-RATE-FOUND-SW, WS-BEST-FLOOR | WS-RATE-FOUND-SW, WS-RX, WS-ABORT-MESSAGE |
| 2110-TAKE-RATE | WS-RT-ANNUAL-RATE, WS-RX, WS-RT-FLOOR | WS-RATE-FOUND-SW, WS-ACCOUNT-RATE, WS-BEST-FLOOR |
| 2200-FIND-MONTH-START | AC-THROUGH-DATE, WS-TH-YEAR, WS-TH-MONTH, WS-MS-YEAR, WS-MONTH-START-DATE | WS-THROUGH-DATE, WS-MS-YEAR, WS-MS-MONTH, WS-MS-DAY, WS-MONTH-START-DAY |
| 2300-ACCRUE-DAYS | AM-BALANCE, WS-ACCOUNT-RATE, WS-DAYS, WS-INTEREST-PRODUCT, WS-INCREMENT, AC-ACCRUED, WS-DEBIT-ACCRUED, WS-CREDIT-ACCRUED | WS-INTEREST-PRODUCT, WS-INCREMENT, AC-ACCRUED, WS-DEBIT-ACCRUED, WS-CREDIT-ACCRUED |
| 2400-CLOSE-MONTH | AC-ACCRUED, AC-DUE, WS-TH-YEAR, WS-TH-MONTH, WS-DUE-MONTH, AC-DUE-MONTH, WS-DEBIT-DUE, WS-CREDIT-DUE | AC-DUE, AC-ACCRUED, WS-DM-YEAR, WS-DM-MONTH, AC-DUE-MONTH, WS-CLOSED-MONTH, WS-DEBIT-DUE, WS-CREDIT-DUE |
| 8000-WRITE-REPORT | WS-BUSINESS-DATE, WS-RATE-COUNT, WS-REPORT-COUNT-LINE, WS-ACCOUNT-COUNT, WS-STARTED-COUNT, WS-CREDIT-ACCRUED, WS-REPORT-AMOUNT-LINE, WS-DEBIT-ACCRUED, WS-CLOSED-MONTH, WS-CREDIT-DUE, WS-DEBIT-DUE | REPORT-LINE, WS-RC-LABEL, WS-RC-VALUE, WS-RA-LABEL, WS-RA-VALUE |
| 9000-CLOSE-FILES | | |
| 9900-ABORT | WS-ABORT-MESSAGE | RETURN-CODE |

## Lineage

Where each data item's values come from, and which records go to which files.

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| AC-ACCRUED | WS-INCREMENT | ADD | 2300-ACCRUE-DAYS | [INTACCR.cbl:317](../../../pangaea/programs/INTACCR.cbl#L317) |
| AC-ACCRUED | AC-DUE | SUBTRACT | 2400-CLOSE-MONTH | [INTACCR.cbl:328](../../../pangaea/programs/INTACCR.cbl#L328) |
| AC-ACCT-NO | AM-ACCT-NO | MOVE | 1400-TAKE-ACCRUAL | [INTACCR.cbl:233](../../../pangaea/programs/INTACCR.cbl#L233) |
| AC-DUE | AC-ACCRUED | COMPUTE | 2400-CLOSE-MONTH | [INTACCR.cbl:327](../../../pangaea/programs/INTACCR.cbl#L327) |
| AC-DUE-MONTH | WS-DUE-MONTH | MOVE | 2400-CLOSE-MONTH | [INTACCR.cbl:331](../../../pangaea/programs/INTACCR.cbl#L331) |
| AC-THROUGH-DATE | WS-BUSINESS-DATE | MOVE | 1400-TAKE-ACCRUAL | [INTACCR.cbl:234](../../../pangaea/programs/INTACCR.cbl#L234) |
| AC-THROUGH-DATE | WS-BUSINESS-DATE | MOVE | 2000-ACCRUE-ACCOUNT | [INTACCR.cbl:259](../../../pangaea/programs/INTACCR.cbl#L259) |
| ACCR-IN-RECORD | file ACCR-IN | READ | 1300-READ-ACCRUAL | [INTACCR.cbl:217](../../../pangaea/programs/INTACCR.cbl#L217) |
| file ACCR-OUT | ACCR-OUT-RECORD | WRITE | 2000-ACCRUE-ACCOUNT | [INTACCR.cbl:260](../../../pangaea/programs/INTACCR.cbl#L260) |
| ACCR-OUT-RECORD | ACCRUAL-RECORD | WRITE | 2000-ACCRUE-ACCOUNT | [INTACCR.cbl:260](../../../pangaea/programs/INTACCR.cbl#L260) |
| ACCRUAL-RECORD | ACCR-IN-RECORD | MOVE | 1400-TAKE-ACCRUAL | [INTACCR.cbl:229](../../../pangaea/programs/INTACCR.cbl#L229) |
| ACCT-MASTER-RECORD | file ACCT-IN | READ | 1200-READ-MASTER | [INTACCR.cbl:210](../../../pangaea/programs/INTACCR.cbl#L210) |
| BANK-CONTROL-RECORD | file CONTROL-FILE | READ | 1000-INITIALIZE | [INTACCR.cbl:177](../../../pangaea/programs/INTACCR.cbl#L177) |
| RATE-RECORD | file RATE-FILE | READ | 1100-LOAD-RATES | [INTACCR.cbl:190](../../../pangaea/programs/INTACCR.cbl#L190) |
| RATE-RECORD | file RATE-FILE | READ | 1100-LOAD-RATES | [INTACCR.cbl:201](../../../pangaea/programs/INTACCR.cbl#L201) |
| REPORT-LINE | WS-BUSINESS-DATE | STRING | 8000-WRITE-REPORT | [INTACCR.cbl:342](../../../pangaea/programs/INTACCR.cbl#L342) |
| REPORT-LINE | WS-REPORT-COUNT-LINE | WRITE | 8000-WRITE-REPORT | [INTACCR.cbl:350](../../../pangaea/programs/INTACCR.cbl#L350) |
| REPORT-LINE | WS-REPORT-COUNT-LINE | WRITE | 8000-WRITE-REPORT | [INTACCR.cbl:353](../../../pangaea/programs/INTACCR.cbl#L353) |
| REPORT-LINE | WS-REPORT-COUNT-LINE | WRITE | 8000-WRITE-REPORT | [INTACCR.cbl:356](../../../pangaea/programs/INTACCR.cbl#L356) |
| REPORT-LINE | WS-REPORT-AMOUNT-LINE | WRITE | 8000-WRITE-REPORT | [INTACCR.cbl:359](../../../pangaea/programs/INTACCR.cbl#L359) |
| REPORT-LINE | WS-REPORT-AMOUNT-LINE | WRITE | 8000-WRITE-REPORT | [INTACCR.cbl:362](../../../pangaea/programs/INTACCR.cbl#L362) |
| REPORT-LINE | WS-CLOSED-MONTH | STRING | 8000-WRITE-REPORT | [INTACCR.cbl:370](../../../pangaea/programs/INTACCR.cbl#L370) |
| REPORT-LINE | WS-REPORT-AMOUNT-LINE | WRITE | 8000-WRITE-REPORT | [INTACCR.cbl:375](../../../pangaea/programs/INTACCR.cbl#L375) |
| REPORT-LINE | WS-REPORT-AMOUNT-LINE | WRITE | 8000-WRITE-REPORT | [INTACCR.cbl:378](../../../pangaea/programs/INTACCR.cbl#L378) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [INTACCR.cbl:345](../../../pangaea/programs/INTACCR.cbl#L345) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [INTACCR.cbl:347](../../../pangaea/programs/INTACCR.cbl#L347) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [INTACCR.cbl:350](../../../pangaea/programs/INTACCR.cbl#L350) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [INTACCR.cbl:353](../../../pangaea/programs/INTACCR.cbl#L353) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [INTACCR.cbl:356](../../../pangaea/programs/INTACCR.cbl#L356) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [INTACCR.cbl:359](../../../pangaea/programs/INTACCR.cbl#L359) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [INTACCR.cbl:362](../../../pangaea/programs/INTACCR.cbl#L362) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [INTACCR.cbl:364](../../../pangaea/programs/INTACCR.cbl#L364) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [INTACCR.cbl:367](../../../pangaea/programs/INTACCR.cbl#L367) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [INTACCR.cbl:372](../../../pangaea/programs/INTACCR.cbl#L372) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [INTACCR.cbl:375](../../../pangaea/programs/INTACCR.cbl#L375) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [INTACCR.cbl:378](../../../pangaea/programs/INTACCR.cbl#L378) |
| WS-ACCOUNT-RATE | WS-RT-ANNUAL-RATE | MOVE | 2110-TAKE-RATE | [INTACCR.cbl:293](../../../pangaea/programs/INTACCR.cbl#L293) |
| WS-ACCRUAL-KEY | AI-ACCT-NO | MOVE | 1300-READ-ACCRUAL | [INTACCR.cbl:221](../../../pangaea/programs/INTACCR.cbl#L221) |
| WS-BEST-FLOOR | WS-RT-FLOOR | MOVE | 2110-TAKE-RATE | [INTACCR.cbl:294](../../../pangaea/programs/INTACCR.cbl#L294) |
| WS-BUSINESS-DATE | BC-BUSINESS-DATE | MOVE | 1000-INITIALIZE | [INTACCR.cbl:182](../../../pangaea/programs/INTACCR.cbl#L182) |
| WS-CLOSED-MONTH | AC-DUE-MONTH | MOVE | 2400-CLOSE-MONTH | [INTACCR.cbl:332](../../../pangaea/programs/INTACCR.cbl#L332) |
| WS-CREDIT-ACCRUED | WS-INCREMENT | ADD | 2300-ACCRUE-DAYS | [INTACCR.cbl:321](../../../pangaea/programs/INTACCR.cbl#L321) |
| WS-CREDIT-DUE | AC-DUE | ADD | 2400-CLOSE-MONTH | [INTACCR.cbl:336](../../../pangaea/programs/INTACCR.cbl#L336) |
| WS-DAYS | WS-MONTH-START-DAY | COMPUTE | 2000-ACCRUE-ACCOUNT | [INTACCR.cbl:252](../../../pangaea/programs/INTACCR.cbl#L252) |
| WS-DAYS | WS-FROM-DAY | COMPUTE | 2000-ACCRUE-ACCOUNT | [INTACCR.cbl:252](../../../pangaea/programs/INTACCR.cbl#L252) |
| WS-DAYS | WS-TODAY | COMPUTE | 2000-ACCRUE-ACCOUNT | [INTACCR.cbl:257](../../../pangaea/programs/INTACCR.cbl#L257) |
| WS-DAYS | WS-FROM-DAY | COMPUTE | 2000-ACCRUE-ACCOUNT | [INTACCR.cbl:257](../../../pangaea/programs/INTACCR.cbl#L257) |
| WS-DEBIT-ACCRUED | WS-INCREMENT | ADD | 2300-ACCRUE-DAYS | [INTACCR.cbl:319](../../../pangaea/programs/INTACCR.cbl#L319) |
| WS-DEBIT-DUE | AC-DUE | ADD | 2400-CLOSE-MONTH | [INTACCR.cbl:334](../../../pangaea/programs/INTACCR.cbl#L334) |
| WS-DM-MONTH | WS-TH-MONTH | MOVE | 2400-CLOSE-MONTH | [INTACCR.cbl:330](../../../pangaea/programs/INTACCR.cbl#L330) |
| WS-DM-YEAR | WS-TH-YEAR | MOVE | 2400-CLOSE-MONTH | [INTACCR.cbl:329](../../../pangaea/programs/INTACCR.cbl#L329) |
| WS-FROM-DAY | AC-THROUGH-DATE | COMPUTE | 2000-ACCRUE-ACCOUNT | [INTACCR.cbl:248](../../../pangaea/programs/INTACCR.cbl#L248) |
| WS-FROM-DAY | WS-MONTH-START-DAY | MOVE | 2000-ACCRUE-ACCOUNT | [INTACCR.cbl:255](../../../pangaea/programs/INTACCR.cbl#L255) |
| WS-INCREMENT | WS-INTEREST-PRODUCT | DIVIDE | 2300-ACCRUE-DAYS | [INTACCR.cbl:315](../../../pangaea/programs/INTACCR.cbl#L315) |
| WS-INTEREST-PRODUCT | AM-BALANCE | MULTIPLY | 2300-ACCRUE-DAYS | [INTACCR.cbl:312](../../../pangaea/programs/INTACCR.cbl#L312) |
| WS-INTEREST-PRODUCT | WS-ACCOUNT-RATE | MULTIPLY | 2300-ACCRUE-DAYS | [INTACCR.cbl:312](../../../pangaea/programs/INTACCR.cbl#L312) |
| WS-INTEREST-PRODUCT | WS-DAYS | MULTIPLY | 2300-ACCRUE-DAYS | [INTACCR.cbl:314](../../../pangaea/programs/INTACCR.cbl#L314) |
| WS-MONTH-START-DAY | WS-MONTH-START-DATE | COMPUTE | 2200-FIND-MONTH-START | [INTACCR.cbl:307](../../../pangaea/programs/INTACCR.cbl#L307) |
| WS-MS-MONTH | WS-TH-MONTH | COMPUTE | 2200-FIND-MONTH-START | [INTACCR.cbl:304](../../../pangaea/programs/INTACCR.cbl#L304) |
| WS-MS-YEAR | WS-TH-YEAR | MOVE | 2200-FIND-MONTH-START | [INTACCR.cbl:299](../../../pangaea/programs/INTACCR.cbl#L299) |
| WS-RA-VALUE | WS-CREDIT-ACCRUED | MOVE | 8000-WRITE-REPORT | [INTACCR.cbl:358](../../../pangaea/programs/INTACCR.cbl#L358) |
| WS-RA-VALUE | WS-DEBIT-ACCRUED | MOVE | 8000-WRITE-REPORT | [INTACCR.cbl:361](../../../pangaea/programs/INTACCR.cbl#L361) |
| WS-RA-VALUE | WS-CREDIT-DUE | MOVE | 8000-WRITE-REPORT | [INTACCR.cbl:374](../../../pangaea/programs/INTACCR.cbl#L374) |
| WS-RA-VALUE | WS-DEBIT-DUE | MOVE | 8000-WRITE-REPORT | [INTACCR.cbl:377](../../../pangaea/programs/INTACCR.cbl#L377) |
| WS-RC-VALUE | WS-RATE-COUNT | MOVE | 8000-WRITE-REPORT | [INTACCR.cbl:349](../../../pangaea/programs/INTACCR.cbl#L349) |
| WS-RC-VALUE | WS-ACCOUNT-COUNT | MOVE | 8000-WRITE-REPORT | [INTACCR.cbl:352](../../../pangaea/programs/INTACCR.cbl#L352) |
| WS-RC-VALUE | WS-STARTED-COUNT | MOVE | 8000-WRITE-REPORT | [INTACCR.cbl:355](../../../pangaea/programs/INTACCR.cbl#L355) |
| WS-RT-ACCT-TYPE | RT-ACCT-TYPE | MOVE | 1100-LOAD-RATES | [INTACCR.cbl:197](../../../pangaea/programs/INTACCR.cbl#L197) |
| WS-RT-ANNUAL-RATE | RT-ANNUAL-RATE | MOVE | 1100-LOAD-RATES | [INTACCR.cbl:200](../../../pangaea/programs/INTACCR.cbl#L200) |
| WS-RT-FLOOR | RT-FLOOR | MOVE | 1100-LOAD-RATES | [INTACCR.cbl:199](../../../pangaea/programs/INTACCR.cbl#L199) |
| WS-RT-KIND | RT-KIND | MOVE | 1100-LOAD-RATES | [INTACCR.cbl:198](../../../pangaea/programs/INTACCR.cbl#L198) |
| WS-THROUGH-DATE | AC-THROUGH-DATE | MOVE | 2200-FIND-MONTH-START | [INTACCR.cbl:298](../../../pangaea/programs/INTACCR.cbl#L298) |
| WS-TODAY | WS-BUSINESS-DATE | COMPUTE | 1000-INITIALIZE | [INTACCR.cbl:183](../../../pangaea/programs/INTACCR.cbl#L183) |

## Records

### BANK-CONTROL-RECORD

File CONTROL-FILE. 80 bytes, layout in copybook [BANKCTL](../copybooks/BANKCTL.md).

### RATE-RECORD

File RATE-FILE. 20 bytes, layout in copybook [RATEREC](../copybooks/RATEREC.md).

### ACCT-MASTER-RECORD

File ACCT-IN. 100 bytes, layout in copybook [ACCTMAST](../copybooks/ACCTMAST.md).

### ACCR-IN-RECORD

File ACCR-IN. 40 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | ACCR-IN-RECORD | 1 | 40 | | | | [INTACCR.cbl:63](../../../pangaea/programs/INTACCR.cbl#L63) |
| 05 | AI-ACCT-NO | 1 | 10 | X(10) | | | [INTACCR.cbl:64](../../../pangaea/programs/INTACCR.cbl#L64) |
| 05 | FILLER | 11 | 30 | X(30) | | | [INTACCR.cbl:65](../../../pangaea/programs/INTACCR.cbl#L65) |

### ACCR-OUT-RECORD

File ACCR-OUT. 40 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | ACCR-OUT-RECORD | 1 | 40 | X(40) | | | [INTACCR.cbl:67](../../../pangaea/programs/INTACCR.cbl#L67) |

### REPORT-LINE

File REPORT-OUT. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | REPORT-LINE | 1 | 80 | X(80) | | | [INTACCR.cbl:69](../../../pangaea/programs/INTACCR.cbl#L69) |

### WS-FILE-STATUSES

Working storage. 12 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-FILE-STATUSES | 1 | 12 | | | | [INTACCR.cbl:72](../../../pangaea/programs/INTACCR.cbl#L72) |
| 05 | WS-CONTROL-STATUS | 1 | 2 | XX | | | [INTACCR.cbl:73](../../../pangaea/programs/INTACCR.cbl#L73) |
| 05 | WS-RATE-STATUS | 3 | 2 | XX | | | [INTACCR.cbl:74](../../../pangaea/programs/INTACCR.cbl#L74) |
| 05 | WS-ACCT-IN-STATUS | 5 | 2 | XX | | | [INTACCR.cbl:75](../../../pangaea/programs/INTACCR.cbl#L75) |
| 05 | WS-ACCR-IN-STATUS | 7 | 2 | XX | | | [INTACCR.cbl:76](../../../pangaea/programs/INTACCR.cbl#L76) |
| 05 | WS-ACCR-OUT-STATUS | 9 | 2 | XX | | | [INTACCR.cbl:77](../../../pangaea/programs/INTACCR.cbl#L77) |
| 05 | WS-REPORT-STATUS | 11 | 2 | XX | | | [INTACCR.cbl:78](../../../pangaea/programs/INTACCR.cbl#L78) |

### WS-BUSINESS-DATE

Working storage. 8 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-BUSINESS-DATE | 1 | 8 | 9(8) | | | [INTACCR.cbl:80](../../../pangaea/programs/INTACCR.cbl#L80) |

### WS-ABORT-MESSAGE

Working storage. 40 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-ABORT-MESSAGE | 1 | 40 | X(40) | | | [INTACCR.cbl:81](../../../pangaea/programs/INTACCR.cbl#L81) |

### WS-END-SWITCHES

Working storage. 1 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-END-SWITCHES | 1 | 1 | | | | [INTACCR.cbl:82](../../../pangaea/programs/INTACCR.cbl#L82) |
| 05 | WS-MASTER-END-SW | 1 | 1 | X | | VALUE 'N' | [INTACCR.cbl:83](../../../pangaea/programs/INTACCR.cbl#L83) |
| 88 | MASTER-AT-END | | | | | VALUE 'Y' | [INTACCR.cbl:84](../../../pangaea/programs/INTACCR.cbl#L84) |

### WS-ACCRUAL-KEY

Working storage. 10 bytes. THE ACCOUNT NUMBER OF THE NEXT ACCRUAL: HIGH-VALUES AT THE END.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-ACCRUAL-KEY | 1 | 10 | X(10) | | | [INTACCR.cbl:86](../../../pangaea/programs/INTACCR.cbl#L86) |

### ACCRUAL-RECORD

Working storage. 40 bytes, layout in copybook [ACCRREC](../copybooks/ACCRREC.md).

### WS-RATE-COUNT

Working storage. 2 bytes. THE RATE TABLE, AS MANY RATES AS THE RATE FILE HOLDS.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-RATE-COUNT | 1 | 2 | 99 | | VALUE 0 | [INTACCR.cbl:92](../../../pangaea/programs/INTACCR.cbl#L92) |

### WS-RATE-TABLE

Working storage. 260 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-RATE-TABLE | 1 | 260 | | | | [INTACCR.cbl:93](../../../pangaea/programs/INTACCR.cbl#L93) |
| 05 | WS-RATE | 1 | 13 | | | OCCURS 1 TO 20 DEPENDING ON WS-RATE-COUNT | [INTACCR.cbl:94](../../../pangaea/programs/INTACCR.cbl#L94) |
| 10 | WS-RT-ACCT-TYPE | 1 | 1 | X | | | [INTACCR.cbl:96](../../../pangaea/programs/INTACCR.cbl#L96) |
| 10 | WS-RT-KIND | 2 | 1 | X | | | [INTACCR.cbl:97](../../../pangaea/programs/INTACCR.cbl#L97) |
| 88 | WS-RT-CREDIT | | | | | VALUE 'C' | [INTACCR.cbl:98](../../../pangaea/programs/INTACCR.cbl#L98) |
| 88 | WS-RT-DEBIT | | | | | VALUE 'D' | [INTACCR.cbl:99](../../../pangaea/programs/INTACCR.cbl#L99) |
| 10 | WS-RT-FLOOR | 3 | 7 | S9(11)V99 | COMP-3 | | [INTACCR.cbl:100](../../../pangaea/programs/INTACCR.cbl#L100) |
| 10 | WS-RT-ANNUAL-RATE | 10 | 4 | 9V9(6) | COMP-3 | | [INTACCR.cbl:101](../../../pangaea/programs/INTACCR.cbl#L101) |

### WS-RX

Working storage. 2 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-RX | 1 | 2 | 99 | | | [INTACCR.cbl:102](../../../pangaea/programs/INTACCR.cbl#L102) |

### WS-ACCOUNT-RATE

Working storage. 4 bytes. THE RATE FOR THE ACCOUNT IN HAND, AND THE TIER IT CAME FROM.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-ACCOUNT-RATE | 1 | 4 | 9V9(6) | COMP-3 | | [INTACCR.cbl:105](../../../pangaea/programs/INTACCR.cbl#L105) |

### WS-RATE-FOUND-SW

Working storage. 1 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-RATE-FOUND-SW | 1 | 1 | X | | | [INTACCR.cbl:106](../../../pangaea/programs/INTACCR.cbl#L106) |
| 88 | RATE-FOUND | | | | | VALUE 'Y' | [INTACCR.cbl:107](../../../pangaea/programs/INTACCR.cbl#L107) |

### WS-BEST-FLOOR

Working storage. 7 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-BEST-FLOOR | 1 | 7 | S9(11)V99 | COMP-3 | | [INTACCR.cbl:108](../../../pangaea/programs/INTACCR.cbl#L108) |

### WS-TODAY

Working storage. 7 bytes. DAYS COUNT FROM 1601-01-01, AS FUNCTION INTEGER-OF-DATE GIVES.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TODAY | 1 | 7 | 9(7) | | | [INTACCR.cbl:111](../../../pangaea/programs/INTACCR.cbl#L111) |

### WS-FROM-DAY

Working storage. 7 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-FROM-DAY | 1 | 7 | 9(7) | | | [INTACCR.cbl:112](../../../pangaea/programs/INTACCR.cbl#L112) |

### WS-MONTH-START-DAY

Working storage. 7 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-MONTH-START-DAY | 1 | 7 | 9(7) | | | [INTACCR.cbl:113](../../../pangaea/programs/INTACCR.cbl#L113) |

### WS-DAYS

Working storage. 3 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-DAYS | 1 | 3 | 9(3) | | | [INTACCR.cbl:114](../../../pangaea/programs/INTACCR.cbl#L114) |

### WS-THROUGH-DATE

Working storage. 8 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-THROUGH-DATE | 1 | 8 | | | | [INTACCR.cbl:115](../../../pangaea/programs/INTACCR.cbl#L115) |
| 05 | WS-TH-YEAR | 1 | 4 | 9(4) | | | [INTACCR.cbl:116](../../../pangaea/programs/INTACCR.cbl#L116) |
| 05 | WS-TH-MONTH | 5 | 2 | 99 | | | [INTACCR.cbl:117](../../../pangaea/programs/INTACCR.cbl#L117) |
| 05 | WS-TH-DAY | 7 | 2 | 99 | | | [INTACCR.cbl:118](../../../pangaea/programs/INTACCR.cbl#L118) |

### WS-MONTH-START

Working storage. 8 bytes. THE FIRST DAY OF THE MONTH AFTER THE LAST ACCRUAL.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-MONTH-START | 1 | 8 | | | | [INTACCR.cbl:120](../../../pangaea/programs/INTACCR.cbl#L120) |
| 05 | WS-MS-YEAR | 1 | 4 | 9(4) | | | [INTACCR.cbl:121](../../../pangaea/programs/INTACCR.cbl#L121) |
| 05 | WS-MS-MONTH | 5 | 2 | 99 | | | [INTACCR.cbl:122](../../../pangaea/programs/INTACCR.cbl#L122) |
| 05 | WS-MS-DAY | 7 | 2 | 99 | | | [INTACCR.cbl:123](../../../pangaea/programs/INTACCR.cbl#L123) |

### WS-MONTH-START-DATE

Working storage. 8 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-MONTH-START-DATE | 1 | 8 | 9(8) | | REDEFINES WS-MONTH-START | [INTACCR.cbl:124](../../../pangaea/programs/INTACCR.cbl#L124) |

### WS-INTEREST-PRODUCT

Working storage. 11 bytes. THE BALANCE TIMES THE RATE TIMES THE DAYS, EXACT, BEFORE IT IS DIVIDED BY THE DAYS OF A YEAR AND ROUNDED.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-INTEREST-PRODUCT | 1 | 11 | S9(13)V9(8) | COMP-3 | | [INTACCR.cbl:129](../../../pangaea/programs/INTACCR.cbl#L129) |

### WS-INCREMENT

Working storage. 7 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-INCREMENT | 1 | 7 | S9(7)V9(6) | COMP-3 | | [INTACCR.cbl:130](../../../pangaea/programs/INTACCR.cbl#L130) |

### WS-DUE-MONTH

Working storage. 6 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-DUE-MONTH | 1 | 6 | | | | [INTACCR.cbl:131](../../../pangaea/programs/INTACCR.cbl#L131) |
| 05 | WS-DM-YEAR | 1 | 4 | 9(4) | | | [INTACCR.cbl:132](../../../pangaea/programs/INTACCR.cbl#L132) |
| 05 | WS-DM-MONTH | 5 | 2 | 99 | | | [INTACCR.cbl:133](../../../pangaea/programs/INTACCR.cbl#L133) |

### WS-TOTALS

Working storage. 52 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TOTALS | 1 | 52 | | | | [INTACCR.cbl:135](../../../pangaea/programs/INTACCR.cbl#L135) |
| 05 | WS-ACCOUNT-COUNT | 1 | 8 | 9(8) | | VALUE 0 | [INTACCR.cbl:136](../../../pangaea/programs/INTACCR.cbl#L136) |
| 05 | WS-STARTED-COUNT | 9 | 8 | 9(8) | | VALUE 0 | [INTACCR.cbl:137](../../../pangaea/programs/INTACCR.cbl#L137) |
| 05 | WS-CREDIT-ACCRUED | 17 | 8 | S9(9)V9(6) | COMP-3 | VALUE 0 | [INTACCR.cbl:138](../../../pangaea/programs/INTACCR.cbl#L138) |
| 05 | WS-DEBIT-ACCRUED | 25 | 8 | S9(9)V9(6) | COMP-3 | VALUE 0 | [INTACCR.cbl:139](../../../pangaea/programs/INTACCR.cbl#L139) |
| 05 | WS-CREDIT-DUE | 33 | 7 | S9(11)V99 | COMP-3 | VALUE 0 | [INTACCR.cbl:140](../../../pangaea/programs/INTACCR.cbl#L140) |
| 05 | WS-DEBIT-DUE | 40 | 7 | S9(11)V99 | COMP-3 | VALUE 0 | [INTACCR.cbl:141](../../../pangaea/programs/INTACCR.cbl#L141) |
| 05 | WS-CLOSED-MONTH | 47 | 6 | 9(6) | | VALUE 0 | [INTACCR.cbl:142](../../../pangaea/programs/INTACCR.cbl#L142) |

### WS-REPORT-COUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-COUNT-LINE | 1 | 80 | | | | [INTACCR.cbl:144](../../../pangaea/programs/INTACCR.cbl#L144) |
| 05 | WS-RC-LABEL | 1 | 30 | X(30) | | | [INTACCR.cbl:145](../../../pangaea/programs/INTACCR.cbl#L145) |
| 05 | WS-RC-VALUE | 31 | 10 | ZZ,ZZZ,ZZ9 | | | [INTACCR.cbl:146](../../../pangaea/programs/INTACCR.cbl#L146) |
| 05 | FILLER | 41 | 40 | X(40) | | VALUE SPACES | [INTACCR.cbl:147](../../../pangaea/programs/INTACCR.cbl#L147) |

### WS-REPORT-AMOUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-AMOUNT-LINE | 1 | 80 | | | | [INTACCR.cbl:148](../../../pangaea/programs/INTACCR.cbl#L148) |
| 05 | WS-RA-LABEL | 1 | 30 | X(30) | | | [INTACCR.cbl:149](../../../pangaea/programs/INTACCR.cbl#L149) |
| 05 | WS-RA-VALUE | 31 | 19 | ZZZ,ZZZ,ZZ9.999999- | | | [INTACCR.cbl:150](../../../pangaea/programs/INTACCR.cbl#L150) |
| 05 | FILLER | 50 | 31 | X(31) | | VALUE SPACES | [INTACCR.cbl:151](../../../pangaea/programs/INTACCR.cbl#L151) |
