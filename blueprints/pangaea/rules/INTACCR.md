# INTACCR: decision tables

The business rules of [INTACCR](../programs/INTACCR.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

## Contents

- [0000-MAIN](#0000-main)
- [1000-INITIALIZE](#1000-initialize)
- [1100-LOAD-RATES](#1100-load-rates)
- [1100-LOAD-RATES, in its loop](#1100-load-rates-in-its-loop)
- [1200-READ-MASTER](#1200-read-master)
- [1300-READ-ACCRUAL](#1300-read-accrual)
- [1400-TAKE-ACCRUAL](#1400-take-accrual)
- [2000-ACCRUE-ACCOUNT](#2000-accrue-account)
- [2100-FIND-RATE](#2100-find-rate)
- [2100-FIND-RATE, in its loop](#2100-find-rate-in-its-loop)
- [2200-FIND-MONTH-START](#2200-find-month-start)
- [2300-ACCRUE-DAYS](#2300-accrue-days)
- [2400-CLOSE-MONTH](#2400-close-month)
- [8000-WRITE-REPORT](#8000-write-report)

## 0000-MAIN

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-ACCRUAL-KEY = HIGH-VALUES` | N | Y |
| → `PERFORM 1000-INITIALIZE` | 1 | 1 |
| → `PERFORM 2000-ACCRUE-ACCOUNT UNTIL MASTER-AT-END` | 2 | 2 |
| → `MOVE 'AN ACCRUAL FOR NO ACCOUNT' TO WS-ABORT-MESSAGE` | 3 | |
| → `PERFORM 9900-ABORT` | 4 | |
| → `PERFORM 8000-WRITE-REPORT` | | 3 |
| → `PERFORM 9000-CLOSE-FILES` | | 4 |
| → `MOVE 0 TO RETURN-CODE` | | 5 |
| → `STOP RUN` | | 6 |

- **Rule 1.** When `WS-ACCRUAL-KEY` is not high values, it performs `1000-INITIALIZE`, performs `2000-ACCRUE-ACCOUNT` until `WS-MASTER-END-SW` is `'Y'` (`MASTER-AT-END`), sets `WS-ABORT-MESSAGE` to `'AN ACCRUAL FOR NO ACCOUNT'`, and performs `9900-ABORT`.
- **Rule 2.** When `WS-ACCRUAL-KEY` is high values, it performs `1000-INITIALIZE`, performs `2000-ACCRUE-ACCOUNT` until `WS-MASTER-END-SW` is `'Y'` (`MASTER-AT-END`), performs `8000-WRITE-REPORT`, performs `9000-CLOSE-FILES`, sets `RETURN-CODE` to 0, and stops the run.

Source: [INTACCR.cbl:156](../../../pangaea/programs/INTACCR.cbl#L156).

## 1000-INITIALIZE

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `WS-FILE-STATUSES = ALL '0'` | N | Y | Y |
| `READ CONTROL-FILE AT END` | | Y | N |
| → `OPEN INPUT CONTROL-FILE RATE-FILE ACCT-IN ACCR-IN OUTPUT ACCR-OUT REPORT-OUT` | 1 | 1 | 1 |
| → `DISPLAY 'INTACCR: CANNOT OPEN FILES ' WS-FILE-STATUSES` | 2 | | |
| → `MOVE 16 TO RETURN-CODE` | 3 | | |
| → `STOP RUN` | 4 | | |
| → `READ CONTROL-FILE` | | 2 | 2 |
| → `MOVE 'NO BANK CONTROL RECORD' TO WS-ABORT-MESSAGE` | | 3 | |
| → `PERFORM 9900-ABORT` | | 4 | |
| → `MOVE BC-BUSINESS-DATE TO WS-BUSINESS-DATE` | | | 3 |
| → `COMPUTE WS-TODAY = FUNCTION INTEGER-OF-DATE(WS-BUSINESS-DATE)` | | | 4 |
| → `PERFORM 1100-LOAD-RATES` | | | 5 |
| → `PERFORM 1200-READ-MASTER` | | | 6 |
| → `PERFORM 1300-READ-ACCRUAL` | | | 7 |

- **Rule 1.** When `WS-FILE-STATUSES` is not `ALL '0'`, it opens `CONTROL-FILE`, `RATE-FILE`, `ACCT-IN`, and `ACCR-IN` for input and `ACCR-OUT` and `REPORT-OUT` for output; displays `'INTACCR: CANNOT OPEN FILES '` `WS-FILE-STATUSES`; sets `RETURN-CODE` to 16; and stops the run.
- **Rule 2.** When `WS-FILE-STATUSES` is `ALL '0'` and reading `CONTROL-FILE` finds no more records, it opens `CONTROL-FILE`, `RATE-FILE`, `ACCT-IN`, and `ACCR-IN` for input and `ACCR-OUT` and `REPORT-OUT` for output; reads `CONTROL-FILE`; sets `WS-ABORT-MESSAGE` to `'NO BANK CONTROL RECORD'`; and performs `9900-ABORT`.
- **Rule 3.** When `WS-FILE-STATUSES` is `ALL '0'` and reading `CONTROL-FILE` finds a record, it opens `CONTROL-FILE`, `RATE-FILE`, `ACCT-IN`, and `ACCR-IN` for input and `ACCR-OUT` and `REPORT-OUT` for output; reads `CONTROL-FILE`; sets `WS-BUSINESS-DATE` to `BC-BUSINESS-DATE`; computes `WS-TODAY` as `FUNCTION INTEGER-OF-DATE(WS-BUSINESS-DATE)`; performs `1100-LOAD-RATES`; performs `1200-READ-MASTER`; and performs `1300-READ-ACCRUAL`.

Source: [INTACCR.cbl:170](../../../pangaea/programs/INTACCR.cbl#L170).

## 1100-LOAD-RATES

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-RATE-STATUS = '10'` | N | Y |
| → `READ RATE-FILE` | 1 | 1 |
| → `PERFORM UNTIL WS-RATE-STATUS NOT = '00'` | 2 | 2 |
| → `MOVE 'CANNOT READ THE RATES' TO WS-ABORT-MESSAGE` | 3 | |
| → `PERFORM 9900-ABORT` | 4 | |

- **Rule 1.** When `WS-RATE-STATUS` is not `'10'`, it reads `RATE-FILE`, repeats its loop until `WS-RATE-STATUS` is not `'00'`, sets `WS-ABORT-MESSAGE` to `'CANNOT READ THE RATES'`, and performs `9900-ABORT`.
- **Rule 2.** When `WS-RATE-STATUS` is `'10'`, it reads `RATE-FILE` and repeats its loop until `WS-RATE-STATUS` is not `'00'`.

Source: [INTACCR.cbl:190](../../../pangaea/programs/INTACCR.cbl#L190).

## 1100-LOAD-RATES, in its loop

Each pass of `PERFORM UNTIL WS-RATE-STATUS NOT = '00'`.

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-RATE-COUNT = 20` | Y | N |
| → `MOVE 'MORE THAN 20 RATES' TO WS-ABORT-MESSAGE` | 1 | |
| → `PERFORM 9900-ABORT` | 2 | |
| → `ADD 1 TO WS-RATE-COUNT` | | 1 |
| → `MOVE RT-ACCT-TYPE TO WS-RT-ACCT-TYPE(WS-RATE-COUNT)` | | 2 |
| → `MOVE RT-KIND TO WS-RT-KIND(WS-RATE-COUNT)` | | 3 |
| → `MOVE RT-FLOOR TO WS-RT-FLOOR(WS-RATE-COUNT)` | | 4 |
| → `MOVE RT-ANNUAL-RATE TO WS-RT-ANNUAL-RATE(WS-RATE-COUNT)` | | 5 |
| → `READ RATE-FILE` | | 6 |

- **Rule 1.** When `WS-RATE-COUNT` is 20, it sets `WS-ABORT-MESSAGE` to `'MORE THAN 20 RATES'` and performs `9900-ABORT`.
- **Rule 2.** When `WS-RATE-COUNT` is not 20, it adds 1 to `WS-RATE-COUNT`, sets `WS-RT-ACCT-TYPE(WS-RATE-COUNT)` to `RT-ACCT-TYPE`, sets `WS-RT-KIND(WS-RATE-COUNT)` to `RT-KIND`, sets `WS-RT-FLOOR(WS-RATE-COUNT)` to `RT-FLOOR`, sets `WS-RT-ANNUAL-RATE(WS-RATE-COUNT)` to `RT-ANNUAL-RATE`, and reads `RATE-FILE`.

Source: [INTACCR.cbl:192](../../../pangaea/programs/INTACCR.cbl#L192).

## 1200-READ-MASTER

| | 1 | 2 |
| --- | :-: | :-: |
| `READ ACCT-IN AT END` | Y | N |
| → `READ ACCT-IN` | 1 | 1 |
| → `SET MASTER-AT-END TO TRUE` | 2 | |

- **Rule 1.** When reading `ACCT-IN` finds no more records, it reads `ACCT-IN` and sets `WS-MASTER-END-SW` to `MASTER-AT-END`.
- **Rule 2.** When reading `ACCT-IN` finds a record, it reads `ACCT-IN`.

Source: [INTACCR.cbl:210](../../../pangaea/programs/INTACCR.cbl#L210).

## 1300-READ-ACCRUAL

| | 1 | 2 |
| --- | :-: | :-: |
| `READ ACCR-IN AT END` | Y | N |
| → `READ ACCR-IN` | 1 | 1 |
| → `MOVE HIGH-VALUES TO WS-ACCRUAL-KEY` | 2 | |
| → `MOVE AI-ACCT-NO TO WS-ACCRUAL-KEY` | | 2 |

- **Rule 1.** When reading `ACCR-IN` finds no more records, it reads `ACCR-IN` and sets `WS-ACCRUAL-KEY` to high values.
- **Rule 2.** When reading `ACCR-IN` finds a record, it reads `ACCR-IN` and sets `WS-ACCRUAL-KEY` to `AI-ACCT-NO`.

Source: [INTACCR.cbl:217](../../../pangaea/programs/INTACCR.cbl#L217).

## 1400-TAKE-ACCRUAL

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `WS-ACCRUAL-KEY = AM-ACCT-NO` | Y | N | N |
| `WS-ACCRUAL-KEY > AM-ACCT-NO` | | Y | N |
| → `MOVE ACCR-IN-RECORD TO ACCRUAL-RECORD` | 1 | | |
| → `PERFORM 1300-READ-ACCRUAL` | 2 | | |
| → `MOVE SPACES TO ACCRUAL-RECORD` | | 1 | |
| → `MOVE AM-ACCT-NO TO AC-ACCT-NO` | | 2 | |
| → `MOVE WS-BUSINESS-DATE TO AC-THROUGH-DATE` | | 3 | |
| → `MOVE 0 TO AC-ACCRUED` | | 4 | |
| → `ADD 1 TO WS-STARTED-COUNT` | | 5 | |
| → `MOVE 'AN ACCRUAL FOR NO ACCOUNT' TO WS-ABORT-MESSAGE` | | | 1 |
| → `PERFORM 9900-ABORT` | | | 2 |

- **Rule 1.** When `WS-ACCRUAL-KEY` is `AM-ACCT-NO`, it sets `ACCRUAL-RECORD` to `ACCR-IN-RECORD` and performs `1300-READ-ACCRUAL`.
- **Rule 2.** When `WS-ACCRUAL-KEY` is not `AM-ACCT-NO` and `WS-ACCRUAL-KEY` is more than `AM-ACCT-NO`, it sets `ACCRUAL-RECORD` to spaces, sets `AC-ACCT-NO` to `AM-ACCT-NO`, sets `AC-THROUGH-DATE` to `WS-BUSINESS-DATE`, sets `AC-ACCRUED` to 0, and adds 1 to `WS-STARTED-COUNT`.
- **Rule 3.** When `WS-ACCRUAL-KEY` is not `AM-ACCT-NO` and `WS-ACCRUAL-KEY` is at most `AM-ACCT-NO`, it sets `WS-ABORT-MESSAGE` to `'AN ACCRUAL FOR NO ACCOUNT'` and performs `9900-ABORT`.

Source: [INTACCR.cbl:227](../../../pangaea/programs/INTACCR.cbl#L227).

## 2000-ACCRUE-ACCOUNT

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-MONTH-START-DAY > WS-TODAY` | N | Y |
| → `PERFORM 1400-TAKE-ACCRUAL` | 1 | 1 |
| → `PERFORM 2100-FIND-RATE` | 2 | 2 |
| → `MOVE 0 TO AC-DUE AC-DUE-MONTH` | 3 | 3 |
| → `COMPUTE WS-FROM-DAY = FUNCTION INTEGER-OF-DATE(AC-THROUGH-DATE)` | 4 | 4 |
| → `PERFORM 2200-FIND-MONTH-START` | 5 | 5 |
| → `COMPUTE WS-DAYS = WS-MONTH-START-DAY - WS-FROM-DAY` | 6 | |
| → `PERFORM 2300-ACCRUE-DAYS` | 7, 11 | 7 |
| → `PERFORM 2400-CLOSE-MONTH` | 8 | |
| → `MOVE WS-MONTH-START-DAY TO WS-FROM-DAY` | 9 | |
| → `COMPUTE WS-DAYS = WS-TODAY - WS-FROM-DAY` | 10 | 6 |
| → `MOVE WS-BUSINESS-DATE TO AC-THROUGH-DATE` | 12 | 8 |
| → `WRITE ACCR-OUT-RECORD FROM ACCRUAL-RECORD` | 13 | 9 |
| → `ADD 1 TO WS-ACCOUNT-COUNT` | 14 | 10 |
| → `PERFORM 1200-READ-MASTER` | 15 | 11 |

- **Rule 1.** When `WS-MONTH-START-DAY` is at most `WS-TODAY`, it performs `1400-TAKE-ACCRUAL`; performs `2100-FIND-RATE`; sets `AC-DUE` and `AC-DUE-MONTH` to 0; computes `WS-FROM-DAY` as `FUNCTION INTEGER-OF-DATE(AC-THROUGH-DATE)`; performs `2200-FIND-MONTH-START`; computes `WS-DAYS` as `WS-MONTH-START-DAY - WS-FROM-DAY`; performs `2300-ACCRUE-DAYS`; performs `2400-CLOSE-MONTH`; sets `WS-FROM-DAY` to `WS-MONTH-START-DAY`; computes `WS-DAYS` as `WS-TODAY - WS-FROM-DAY`; performs `2300-ACCRUE-DAYS`; sets `AC-THROUGH-DATE` to `WS-BUSINESS-DATE`; writes `ACCR-OUT-RECORD` from `ACCRUAL-RECORD`; adds 1 to `WS-ACCOUNT-COUNT`; and performs `1200-READ-MASTER`.
- **Rule 2.** When `WS-MONTH-START-DAY` is more than `WS-TODAY`, it performs `1400-TAKE-ACCRUAL`; performs `2100-FIND-RATE`; sets `AC-DUE` and `AC-DUE-MONTH` to 0; computes `WS-FROM-DAY` as `FUNCTION INTEGER-OF-DATE(AC-THROUGH-DATE)`; performs `2200-FIND-MONTH-START`; computes `WS-DAYS` as `WS-TODAY - WS-FROM-DAY`; performs `2300-ACCRUE-DAYS`; sets `AC-THROUGH-DATE` to `WS-BUSINESS-DATE`; writes `ACCR-OUT-RECORD` from `ACCRUAL-RECORD`; adds 1 to `WS-ACCOUNT-COUNT`; and performs `1200-READ-MASTER`.

Source: [INTACCR.cbl:245](../../../pangaea/programs/INTACCR.cbl#L245).

## 2100-FIND-RATE

| | 1 | 2 |
| --- | :-: | :-: |
| `RATE-FOUND` | N | Y |
| → `MOVE 'N' TO WS-RATE-FOUND-SW` | 1 | 1 |
| → `PERFORM VARYING WS-RX FROM 1 BY 1 UNTIL WS-RX > WS-RATE-COUNT` | 2 | 2 |
| → `MOVE 'NO RATE FOR AN ACCOUNT TYPE' TO WS-ABORT-MESSAGE` | 3 | |
| → `PERFORM 9900-ABORT` | 4 | |

- **Rule 1.** When `WS-RATE-FOUND-SW` is not `'Y'` (not `RATE-FOUND`), it sets `WS-RATE-FOUND-SW` to `'N'`, repeats its loop varying `WS-RX` from 1 by 1 until `WS-RX` is more than `WS-RATE-COUNT`, sets `WS-ABORT-MESSAGE` to `'NO RATE FOR AN ACCOUNT TYPE'`, and performs `9900-ABORT`.
- **Rule 2.** When `WS-RATE-FOUND-SW` is `'Y'` (`RATE-FOUND`), it sets `WS-RATE-FOUND-SW` to `'N'` and repeats its loop varying `WS-RX` from 1 by 1 until `WS-RX` is more than `WS-RATE-COUNT`.

Source: [INTACCR.cbl:267](../../../pangaea/programs/INTACCR.cbl#L267).

## 2100-FIND-RATE, in its loop

Each pass of `PERFORM VARYING WS-RX FROM 1 BY 1 UNTIL WS-RX > WS-RATE-COUNT`.

| | 1 | 2 | 3 | 4 | 5 |
| --- | :-: | :-: | :-: | :-: | :-: |
| `WS-RT-ACCT-TYPE(WS-RX) = AM-ACCT-TYPE` | Y | Y | Y | Y | N |
| `AM-BALANCE < 0` | Y | Y | N | N | |
| `WS-RT-DEBIT(WS-RX)` | Y | N | | | |
| `WS-RT-CREDIT(WS-RX) AND WS-RT-FLOOR(WS-RX) NOT > AM-BALANCE AND (NOT RATE-FOUND OR WS-RT-FLOOR(WS-RX) > WS-BEST-FLOOR)` | | | Y | N | |
| → `PERFORM 2110-TAKE-RATE` | 1 | | 1 | | |

- **Rule 1.** When `WS-RT-ACCT-TYPE(WS-RX)` is `AM-ACCT-TYPE`, `AM-BALANCE` is less than 0, and `WS-RT-KIND` is `'D'` (`WS-RT-DEBIT(WS-RX)`), it performs `2110-TAKE-RATE`.
- **Rule 2.** When `WS-RT-ACCT-TYPE(WS-RX)` is `AM-ACCT-TYPE`, `AM-BALANCE` is less than 0, and `WS-RT-KIND` is not `'D'` (not `WS-RT-DEBIT(WS-RX)`), it does nothing more.
- **Rule 3.** When `WS-RT-ACCT-TYPE(WS-RX)` is `AM-ACCT-TYPE`; `AM-BALANCE` is at least 0; and `WS-RT-KIND` is `'C'` (`WS-RT-CREDIT(WS-RX)`), `WS-RT-FLOOR(WS-RX)` is at most `AM-BALANCE`, and either `WS-RATE-FOUND-SW` is not `'Y'` (not `RATE-FOUND`) or `WS-RT-FLOOR(WS-RX)` is more than `WS-BEST-FLOOR`, it performs `2110-TAKE-RATE`.
- **Rule 4.** When `WS-RT-ACCT-TYPE(WS-RX)` is `AM-ACCT-TYPE`; `AM-BALANCE` is at least 0; and `WS-RT-KIND` is not `'C'` (not `WS-RT-CREDIT(WS-RX)`), `WS-RT-FLOOR(WS-RX)` is more than `AM-BALANCE`, or both `WS-RATE-FOUND-SW` is `'Y'` (`RATE-FOUND`) and `WS-RT-FLOOR(WS-RX)` is at most `WS-BEST-FLOOR`, it does nothing more.
- **Rule 5.** When `WS-RT-ACCT-TYPE(WS-RX)` is not `AM-ACCT-TYPE`, it does nothing more.

Source: [INTACCR.cbl:269](../../../pangaea/programs/INTACCR.cbl#L269).

## 2200-FIND-MONTH-START

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-TH-MONTH = 12` | Y | N |
| → `MOVE AC-THROUGH-DATE TO WS-THROUGH-DATE` | 1 | 1 |
| → `MOVE WS-TH-YEAR TO WS-MS-YEAR` | 2 | 2 |
| → `ADD 1 TO WS-MS-YEAR` | 3 | |
| → `MOVE 1 TO WS-MS-MONTH` | 4 | |
| → `MOVE 1 TO WS-MS-DAY` | 5 | 4 |
| → `COMPUTE WS-MONTH-START-DAY = FUNCTION INTEGER-OF-DATE(WS-MONTH-START-DATE)` | 6 | 5 |
| → `COMPUTE WS-MS-MONTH = WS-TH-MONTH + 1` | | 3 |

- **Rule 1.** When `WS-TH-MONTH` is 12, it sets `WS-THROUGH-DATE` to `AC-THROUGH-DATE`, sets `WS-MS-YEAR` to `WS-TH-YEAR`, adds 1 to `WS-MS-YEAR`, sets `WS-MS-MONTH` to 1, sets `WS-MS-DAY` to 1, and computes `WS-MONTH-START-DAY` as `FUNCTION INTEGER-OF-DATE(WS-MONTH-START-DATE)`.
- **Rule 2.** When `WS-TH-MONTH` is not 12, it sets `WS-THROUGH-DATE` to `AC-THROUGH-DATE`, sets `WS-MS-YEAR` to `WS-TH-YEAR`, computes `WS-MS-MONTH` as `WS-TH-MONTH + 1`, sets `WS-MS-DAY` to 1, and computes `WS-MONTH-START-DAY` as `FUNCTION INTEGER-OF-DATE(WS-MONTH-START-DATE)`.

Source: [INTACCR.cbl:298](../../../pangaea/programs/INTACCR.cbl#L298).

## 2300-ACCRUE-DAYS

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-INCREMENT < 0` | Y | N |
| → `MULTIPLY AM-BALANCE BY WS-ACCOUNT-RATE GIVING WS-INTEREST-PRODUCT` | 1 | 1 |
| → `MULTIPLY WS-DAYS BY WS-INTEREST-PRODUCT` | 2 | 2 |
| → `DIVIDE WS-INTEREST-PRODUCT BY 365 GIVING WS-INCREMENT ROUNDED` | 3 | 3 |
| → `ADD WS-INCREMENT TO AC-ACCRUED` | 4 | 4 |
| → `ADD WS-INCREMENT TO WS-DEBIT-ACCRUED` | 5 | |
| → `ADD WS-INCREMENT TO WS-CREDIT-ACCRUED` | | 5 |

- **Rule 1.** When `WS-INCREMENT` is less than 0, it runs `MULTIPLY AM-BALANCE BY WS-ACCOUNT-RATE GIVING WS-INTEREST-PRODUCT`, runs `MULTIPLY WS-DAYS BY WS-INTEREST-PRODUCT`, runs `DIVIDE WS-INTEREST-PRODUCT BY 365 GIVING WS-INCREMENT ROUNDED`, adds `WS-INCREMENT` to `AC-ACCRUED`, and adds `WS-INCREMENT` to `WS-DEBIT-ACCRUED`.
- **Rule 2.** When `WS-INCREMENT` is at least 0, it runs `MULTIPLY AM-BALANCE BY WS-ACCOUNT-RATE GIVING WS-INTEREST-PRODUCT`, runs `MULTIPLY WS-DAYS BY WS-INTEREST-PRODUCT`, runs `DIVIDE WS-INTEREST-PRODUCT BY 365 GIVING WS-INCREMENT ROUNDED`, adds `WS-INCREMENT` to `AC-ACCRUED`, and adds `WS-INCREMENT` to `WS-CREDIT-ACCRUED`.

Source: [INTACCR.cbl:312](../../../pangaea/programs/INTACCR.cbl#L312).

## 2400-CLOSE-MONTH

| | 1 | 2 |
| --- | :-: | :-: |
| `AC-DUE < 0` | Y | N |
| → `COMPUTE AC-DUE ROUNDED = AC-ACCRUED` | 1 | 1 |
| → `SUBTRACT AC-DUE FROM AC-ACCRUED` | 2 | 2 |
| → `MOVE WS-TH-YEAR TO WS-DM-YEAR` | 3 | 3 |
| → `MOVE WS-TH-MONTH TO WS-DM-MONTH` | 4 | 4 |
| → `MOVE WS-DUE-MONTH TO AC-DUE-MONTH` | 5 | 5 |
| → `MOVE AC-DUE-MONTH TO WS-CLOSED-MONTH` | 6 | 6 |
| → `ADD AC-DUE TO WS-DEBIT-DUE` | 7 | |
| → `ADD AC-DUE TO WS-CREDIT-DUE` | | 7 |

- **Rule 1.** When `AC-DUE` is less than 0, it computes `AC-DUE` as `AC-ACCRUED`, subtracts `AC-DUE` from `AC-ACCRUED`, sets `WS-DM-YEAR` to `WS-TH-YEAR`, sets `WS-DM-MONTH` to `WS-TH-MONTH`, sets `AC-DUE-MONTH` to `WS-DUE-MONTH`, sets `WS-CLOSED-MONTH` to `AC-DUE-MONTH`, and adds `AC-DUE` to `WS-DEBIT-DUE`.
- **Rule 2.** When `AC-DUE` is at least 0, it computes `AC-DUE` as `AC-ACCRUED`, subtracts `AC-DUE` from `AC-ACCRUED`, sets `WS-DM-YEAR` to `WS-TH-YEAR`, sets `WS-DM-MONTH` to `WS-TH-MONTH`, sets `AC-DUE-MONTH` to `WS-DUE-MONTH`, sets `WS-CLOSED-MONTH` to `AC-DUE-MONTH`, and adds `AC-DUE` to `WS-CREDIT-DUE`.

Source: [INTACCR.cbl:327](../../../pangaea/programs/INTACCR.cbl#L327).

## 8000-WRITE-REPORT

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-CLOSED-MONTH = 0` | Y | N |
| → `MOVE SPACES TO REPORT-LINE` | 1, 4, 21 | 1, 4, 21, 23 |
| → `STRING 'INTACCR  INTEREST ACCRUAL        BUSINESS DATE ' WS-BUSINESS-DATE DELIMITED BY SIZE INTO REPORT-LINE` | 2 | 2 |
| → `WRITE REPORT-LINE` | 3, 5, 22, 24 | 3, 5, 22, 25 |
| → `MOVE 'RATES' TO WS-RC-LABEL` | 6 | 6 |
| → `MOVE WS-RATE-COUNT TO WS-RC-VALUE` | 7 | 7 |
| → `WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE` | 8, 11, 14 | 8, 11, 14 |
| → `MOVE 'ACCOUNTS ACCRUED' TO WS-RC-LABEL` | 9 | 9 |
| → `MOVE WS-ACCOUNT-COUNT TO WS-RC-VALUE` | 10 | 10 |
| → `MOVE 'ACCRUALS STARTED' TO WS-RC-LABEL` | 12 | 12 |
| → `MOVE WS-STARTED-COUNT TO WS-RC-VALUE` | 13 | 13 |
| → `MOVE 'CREDIT INTEREST ACCRUED' TO WS-RA-LABEL` | 15 | 15 |
| → `MOVE WS-CREDIT-ACCRUED TO WS-RA-VALUE` | 16 | 16 |
| → `WRITE REPORT-LINE FROM WS-REPORT-AMOUNT-LINE` | 17, 20 | 17, 20, 28, 31 |
| → `MOVE 'DEBIT INTEREST ACCRUED' TO WS-RA-LABEL` | 18 | 18 |
| → `MOVE WS-DEBIT-ACCRUED TO WS-RA-VALUE` | 19 | 19 |
| → `MOVE 'NO MONTH CLOSED' TO REPORT-LINE` | 23 | |
| → `STRING 'MONTH CLOSED                  ' WS-CLOSED-MONTH DELIMITED BY SIZE INTO REPORT-LINE` | | 24 |
| → `MOVE 'CREDIT INTEREST DUE' TO WS-RA-LABEL` | | 26 |
| → `MOVE WS-CREDIT-DUE TO WS-RA-VALUE` | | 27 |
| → `MOVE 'DEBIT INTEREST DUE' TO WS-RA-LABEL` | | 29 |
| → `MOVE WS-DEBIT-DUE TO WS-RA-VALUE` | | 30 |

- **Rule 1.** When `WS-CLOSED-MONTH` is 0, it sets `REPORT-LINE` to spaces, builds `REPORT-LINE` from `'INTACCR  INTEREST ACCRUAL        BUSINESS DATE '` `WS-BUSINESS-DATE`, writes `REPORT-LINE`, sets `REPORT-LINE` to spaces, writes `REPORT-LINE`, sets `WS-RC-LABEL` to `'RATES'`, sets `WS-RC-VALUE` to `WS-RATE-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `WS-RC-LABEL` to `'ACCOUNTS ACCRUED'`, sets `WS-RC-VALUE` to `WS-ACCOUNT-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `WS-RC-LABEL` to `'ACCRUALS STARTED'`, sets `WS-RC-VALUE` to `WS-STARTED-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `WS-RA-LABEL` to `'CREDIT INTEREST ACCRUED'`, sets `WS-RA-VALUE` to `WS-CREDIT-ACCRUED`, writes `REPORT-LINE` from `WS-REPORT-AMOUNT-LINE`, sets `WS-RA-LABEL` to `'DEBIT INTEREST ACCRUED'`, sets `WS-RA-VALUE` to `WS-DEBIT-ACCRUED`, writes `REPORT-LINE` from `WS-REPORT-AMOUNT-LINE`, sets `REPORT-LINE` to spaces, writes `REPORT-LINE`, sets `REPORT-LINE` to `'NO MONTH CLOSED'`, and writes `REPORT-LINE`.
- **Rule 2.** When `WS-CLOSED-MONTH` is not 0, it sets `REPORT-LINE` to spaces, builds `REPORT-LINE` from `'INTACCR  INTEREST ACCRUAL        BUSINESS DATE '` `WS-BUSINESS-DATE`, writes `REPORT-LINE`, sets `REPORT-LINE` to spaces, writes `REPORT-LINE`, sets `WS-RC-LABEL` to `'RATES'`, sets `WS-RC-VALUE` to `WS-RATE-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `WS-RC-LABEL` to `'ACCOUNTS ACCRUED'`, sets `WS-RC-VALUE` to `WS-ACCOUNT-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `WS-RC-LABEL` to `'ACCRUALS STARTED'`, sets `WS-RC-VALUE` to `WS-STARTED-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `WS-RA-LABEL` to `'CREDIT INTEREST ACCRUED'`, sets `WS-RA-VALUE` to `WS-CREDIT-ACCRUED`, writes `REPORT-LINE` from `WS-REPORT-AMOUNT-LINE`, sets `WS-RA-LABEL` to `'DEBIT INTEREST ACCRUED'`, sets `WS-RA-VALUE` to `WS-DEBIT-ACCRUED`, writes `REPORT-LINE` from `WS-REPORT-AMOUNT-LINE`, sets `REPORT-LINE` to spaces, writes `REPORT-LINE`, sets `REPORT-LINE` to spaces, builds `REPORT-LINE` from `'MONTH CLOSED                  '` `WS-CLOSED-MONTH`, writes `REPORT-LINE`, sets `WS-RA-LABEL` to `'CREDIT INTEREST DUE'`, sets `WS-RA-VALUE` to `WS-CREDIT-DUE`, writes `REPORT-LINE` from `WS-REPORT-AMOUNT-LINE`, sets `WS-RA-LABEL` to `'DEBIT INTEREST DUE'`, sets `WS-RA-VALUE` to `WS-DEBIT-DUE`, and writes `REPORT-LINE` from `WS-REPORT-AMOUNT-LINE`.

Source: [INTACCR.cbl:341](../../../pangaea/programs/INTACCR.cbl#L341).
