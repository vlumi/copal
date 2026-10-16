# INTACCR: decision tables

The business rules of [INTACCR](../programs/INTACCR.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

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
| → `PERFORM 8000-WRITE-REPORT` | 5 | 3 |
| → `PERFORM 9000-CLOSE-FILES` | 6 | 4 |
| → `MOVE 0 TO RETURN-CODE` | 7 | 5 |
| → `STOP RUN` | 8 | 6 |

Source: [INTACCR.cbl:155](../../../pangaea/programs/INTACCR.cbl#L155).

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
| → `MOVE BC-BUSINESS-DATE TO WS-BUSINESS-DATE` | | 5 | 3 |
| → `COMPUTE WS-TODAY = FUNCTION INTEGER-OF-DATE(WS-BUSINESS-DATE)` | | 6 | 4 |
| → `PERFORM 1100-LOAD-RATES` | | 7 | 5 |
| → `PERFORM 1200-READ-MASTER` | | 8 | 6 |
| → `PERFORM 1300-READ-ACCRUAL` | | 9 | 7 |

Source: [INTACCR.cbl:169](../../../pangaea/programs/INTACCR.cbl#L169).

## 1100-LOAD-RATES

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-RATE-STATUS = '10'` | N | Y |
| → `READ RATE-FILE` | 1 | 1 |
| → `PERFORM UNTIL WS-RATE-STATUS NOT = '00'` | 2 | 2 |
| → `MOVE 'CANNOT READ THE RATES' TO WS-ABORT-MESSAGE` | 3 | |
| → `PERFORM 9900-ABORT` | 4 | |

Source: [INTACCR.cbl:189](../../../pangaea/programs/INTACCR.cbl#L189).

## 1100-LOAD-RATES, in its loop

Each pass of `PERFORM UNTIL WS-RATE-STATUS NOT = '00'`.

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-RATE-COUNT = 20` | Y | N |
| → `MOVE 'MORE THAN 20 RATES' TO WS-ABORT-MESSAGE` | 1 | |
| → `PERFORM 9900-ABORT` | 2 | |
| → `ADD 1 TO WS-RATE-COUNT` | 3 | 1 |
| → `MOVE RT-ACCT-TYPE TO WS-RT-ACCT-TYPE(WS-RATE-COUNT)` | 4 | 2 |
| → `MOVE RT-KIND TO WS-RT-KIND(WS-RATE-COUNT)` | 5 | 3 |
| → `MOVE RT-FLOOR TO WS-RT-FLOOR(WS-RATE-COUNT)` | 6 | 4 |
| → `MOVE RT-ANNUAL-RATE TO WS-RT-ANNUAL-RATE(WS-RATE-COUNT)` | 7 | 5 |
| → `READ RATE-FILE` | 8 | 6 |

Source: [INTACCR.cbl:191](../../../pangaea/programs/INTACCR.cbl#L191).

## 1200-READ-MASTER

| | 1 | 2 |
| --- | :-: | :-: |
| `READ ACCT-IN AT END` | Y | N |
| → `READ ACCT-IN` | 1 | 1 |
| → `SET MASTER-AT-END TO TRUE` | 2 | |

Source: [INTACCR.cbl:209](../../../pangaea/programs/INTACCR.cbl#L209).

## 1300-READ-ACCRUAL

| | 1 | 2 |
| --- | :-: | :-: |
| `READ ACCR-IN AT END` | Y | N |
| → `READ ACCR-IN` | 1 | 1 |
| → `MOVE HIGH-VALUES TO WS-ACCRUAL-KEY` | 2 | |
| → `MOVE AI-ACCT-NO TO WS-ACCRUAL-KEY` | | 2 |

Source: [INTACCR.cbl:216](../../../pangaea/programs/INTACCR.cbl#L216).

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

Source: [INTACCR.cbl:226](../../../pangaea/programs/INTACCR.cbl#L226).

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

Source: [INTACCR.cbl:244](../../../pangaea/programs/INTACCR.cbl#L244).

## 2100-FIND-RATE

| | 1 | 2 |
| --- | :-: | :-: |
| `RATE-FOUND` | N | Y |
| → `MOVE 'N' TO WS-RATE-FOUND-SW` | 1 | 1 |
| → `PERFORM VARYING WS-RX FROM 1 BY 1 UNTIL WS-RX > WS-RATE-COUNT` | 2 | 2 |
| → `MOVE 'NO RATE FOR AN ACCOUNT TYPE' TO WS-ABORT-MESSAGE` | 3 | |
| → `PERFORM 9900-ABORT` | 4 | |

Source: [INTACCR.cbl:266](../../../pangaea/programs/INTACCR.cbl#L266).

## 2100-FIND-RATE, in its loop

Each pass of `PERFORM VARYING WS-RX FROM 1 BY 1 UNTIL WS-RX > WS-RATE-COUNT`.

| | 1 | 2 | 3 | 4 | 5 |
| --- | :-: | :-: | :-: | :-: | :-: |
| `WS-RT-ACCT-TYPE(WS-RX) = AM-ACCT-TYPE` | Y | Y | Y | Y | N |
| `AM-BALANCE < 0` | Y | Y | N | N | |
| `WS-RT-DEBIT(WS-RX)` | Y | N | | | |
| `WS-RT-CREDIT(WS-RX) AND WS-RT-FLOOR(WS-RX) NOT > AM-BALANCE AND (NOT RATE-FOUND OR WS-RT-FLOOR(WS-RX) > WS-BEST-FLOOR)` | | | Y | N | |
| → `PERFORM 2110-TAKE-RATE` | 1 | | 1 | | |

Source: [INTACCR.cbl:268](../../../pangaea/programs/INTACCR.cbl#L268).

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

Source: [INTACCR.cbl:297](../../../pangaea/programs/INTACCR.cbl#L297).

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

Source: [INTACCR.cbl:311](../../../pangaea/programs/INTACCR.cbl#L311).

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

Source: [INTACCR.cbl:326](../../../pangaea/programs/INTACCR.cbl#L326).

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

Source: [INTACCR.cbl:340](../../../pangaea/programs/INTACCR.cbl#L340).
