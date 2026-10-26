# GLSUM: decision tables

The business rules of [GLSUM](../programs/GLSUM.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

## Contents

- [0000-MAIN](#0000-main)
- [1000-INITIALIZE](#1000-initialize)
- [2000-POST-MONTH](#2000-post-month)
- [2000-POST-MONTH, in its loop](#2000-post-month-in-its-loop)
- [2100-POST-TO-LEDGER](#2100-post-to-ledger)
- [2200-FIND-ACCOUNT](#2200-find-account)
- [3000-WRITE-LEDGER](#3000-write-ledger)
- [3000-WRITE-LEDGER, in its loop](#3000-write-ledger-in-its-loop)
- [8000-WRITE-REPORT](#8000-write-report)

## 0000-MAIN

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `WS-CLOSED-MONTH-NUMBER = 0` | N | N | Y |
| `WS-TRANSIT-NET = 0` | N | Y | Y |
| → `PERFORM 1000-INITIALIZE` | 1 | 1 | 1 |
| → `PERFORM 2000-POST-MONTH` | 2 | 2 | |
| → `PERFORM 3000-WRITE-LEDGER` | 3 | 3 | |
| → `PERFORM 8000-WRITE-REPORT` | 4 | 4 | 2 |
| → `PERFORM 9000-CLOSE-FILES` | 5 | 5 | 3 |
| → `MOVE 4 TO RETURN-CODE` | 6 | | |
| → `STOP RUN` | 7 | 7 | 5 |
| → `MOVE 0 TO RETURN-CODE` | | 6 | 4 |

- **Rule 1.** When `WS-CLOSED-MONTH-NUMBER` is not 0 and `WS-TRANSIT-NET` is not 0, it performs `1000-INITIALIZE`, performs `2000-POST-MONTH`, performs `3000-WRITE-LEDGER`, performs `8000-WRITE-REPORT`, performs `9000-CLOSE-FILES`, sets `RETURN-CODE` to 4, and stops the run.
- **Rule 2.** When `WS-CLOSED-MONTH-NUMBER` is not 0 and `WS-TRANSIT-NET` is 0, it performs `1000-INITIALIZE`, performs `2000-POST-MONTH`, performs `3000-WRITE-LEDGER`, performs `8000-WRITE-REPORT`, performs `9000-CLOSE-FILES`, sets `RETURN-CODE` to 0, and stops the run.
- **Rule 3.** When `WS-CLOSED-MONTH-NUMBER` is 0 and `WS-TRANSIT-NET` is 0, it performs `1000-INITIALIZE`, performs `8000-WRITE-REPORT`, performs `9000-CLOSE-FILES`, sets `RETURN-CODE` to 0, and stops the run.

Source: [GLSUM.cbl:161](../../../pangaea/programs/GLSUM.cbl#L161).

## 1000-INITIALIZE

| | 1 | 2 | 3 | 4 |
| --- | :-: | :-: | :-: | :-: |
| `WS-FILE-STATUSES = ALL '0'` | N | Y | Y | Y |
| `READ ACCR-IN AT END` | | N | N | Y |
| `WS-CM-MONTH = 12` | | Y | N | N |
| → `OPEN INPUT ACCR-IN AUDIT-IN OUTPUT GL-OUT REPORT-OUT` | 1 | 1 | 1 | 1 |
| → `DISPLAY 'GLSUM: CANNOT OPEN FILES ' WS-FILE-STATUSES` | 2 | | | |
| → `MOVE 16 TO RETURN-CODE` | 3 | | | |
| → `STOP RUN` | 4 | | | |
| → `MOVE 0 TO WS-CLOSED-MONTH-NUMBER` | | 2 | 2 | 2 |
| → `READ ACCR-IN` | | 3 | 3 | 3 |
| → `MOVE AC-DUE-MONTH TO WS-CLOSED-MONTH-NUMBER` | | 4 | 4 | |
| → `MOVE WS-CM-YEAR TO WS-NM-YEAR` | | 5 | 5 | 4 |
| → `ADD 1 TO WS-NM-YEAR` | | 6 | | |
| → `MOVE 1 TO WS-NM-MONTH` | | 7 | | |
| → `PERFORM VARYING WS-TX FROM 1 BY 1 UNTIL WS-TX > WS-LEDGER-COUNT` | | 8 | 7 | 6 |
| → `COMPUTE WS-NM-MONTH = WS-CM-MONTH + 1` | | | 6 | 5 |

- **Rule 1.** When `WS-FILE-STATUSES` is not `ALL '0'`, it opens `ACCR-IN` and `AUDIT-IN` for input and `GL-OUT` and `REPORT-OUT` for output; displays `'GLSUM: CANNOT OPEN FILES '` `WS-FILE-STATUSES`; sets `RETURN-CODE` to 16; and stops the run.
- **Rule 2.** When `WS-FILE-STATUSES` is `ALL '0'`, reading `ACCR-IN` finds a record, and `WS-CM-MONTH` is 12, it opens `ACCR-IN` and `AUDIT-IN` for input and `GL-OUT` and `REPORT-OUT` for output; sets `WS-CLOSED-MONTH-NUMBER` to 0; reads `ACCR-IN`; sets `WS-CLOSED-MONTH-NUMBER` to `AC-DUE-MONTH`; sets `WS-NM-YEAR` to `WS-CM-YEAR`; adds 1 to `WS-NM-YEAR`; sets `WS-NM-MONTH` to 1; and repeats its loop varying `WS-TX` from 1 by 1 until `WS-TX` is more than `WS-LEDGER-COUNT`.
- **Rule 3.** When `WS-FILE-STATUSES` is `ALL '0'`, reading `ACCR-IN` finds a record, and `WS-CM-MONTH` is not 12, it opens `ACCR-IN` and `AUDIT-IN` for input and `GL-OUT` and `REPORT-OUT` for output; sets `WS-CLOSED-MONTH-NUMBER` to 0; reads `ACCR-IN`; sets `WS-CLOSED-MONTH-NUMBER` to `AC-DUE-MONTH`; sets `WS-NM-YEAR` to `WS-CM-YEAR`; computes `WS-NM-MONTH` as `WS-CM-MONTH + 1`; and repeats its loop varying `WS-TX` from 1 by 1 until `WS-TX` is more than `WS-LEDGER-COUNT`.
- **Rule 4.** When `WS-FILE-STATUSES` is `ALL '0'`, reading `ACCR-IN` finds no more records, and `WS-CM-MONTH` is not 12, it opens `ACCR-IN` and `AUDIT-IN` for input and `GL-OUT` and `REPORT-OUT` for output; sets `WS-CLOSED-MONTH-NUMBER` to 0; reads `ACCR-IN`; sets `WS-NM-YEAR` to `WS-CM-YEAR`; computes `WS-NM-MONTH` as `WS-CM-MONTH + 1`; and repeats its loop varying `WS-TX` from 1 by 1 until `WS-TX` is more than `WS-LEDGER-COUNT`.

Source: [GLSUM.cbl:178](../../../pangaea/programs/GLSUM.cbl#L178).

## 2000-POST-MONTH

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-AUDIT-STATUS = '10'` | N | Y |
| → `READ AUDIT-IN` | 1 | 1 |
| → `PERFORM UNTIL WS-AUDIT-STATUS NOT = '00'` | 2 | 2 |
| → `MOVE 'CANNOT READ THE AUDIT TRAILS' TO WS-ABORT-MESSAGE` | 3 | |
| → `PERFORM 9900-ABORT` | 4 | |

- **Rule 1.** When `WS-AUDIT-STATUS` is not `'10'`, it reads `AUDIT-IN`, repeats its loop until `WS-AUDIT-STATUS` is not `'00'`, sets `WS-ABORT-MESSAGE` to `'CANNOT READ THE AUDIT TRAILS'`, and performs `9900-ABORT`.
- **Rule 2.** When `WS-AUDIT-STATUS` is `'10'`, it reads `AUDIT-IN` and repeats its loop until `WS-AUDIT-STATUS` is not `'00'`.

Source: [GLSUM.cbl:205](../../../pangaea/programs/GLSUM.cbl#L205).

## 2000-POST-MONTH, in its loop

Each pass of `PERFORM UNTIL WS-AUDIT-STATUS NOT = '00'`.

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-POSTING-MONTH = WS-CLOSED-MONTH OR (WS-POSTING-MONTH = WS-NEXT-MONTH AND WS-MONTH-END-CODE)` | Y | N |
| → `MOVE AU-BUSINESS-DATE(1:6) TO WS-POSTING-MONTH` | 1 | 1 |
| → `MOVE AU-TRAN-CODE TO WS-POSTING-CODE` | 2 | 2 |
| → `PERFORM 2100-POST-TO-LEDGER` | 3 | |
| → `READ AUDIT-IN` | 4 | 3 |

- **Rule 1.** When `WS-POSTING-MONTH` is `WS-CLOSED-MONTH` or both `WS-POSTING-MONTH` is `WS-NEXT-MONTH` and `WS-POSTING-CODE` is `'IN'`, `'OD'`, or `'MF'` (`WS-MONTH-END-CODE`), it sets `WS-POSTING-MONTH` to `AU-BUSINESS-DATE(1:6)`, sets `WS-POSTING-CODE` to `AU-TRAN-CODE`, performs `2100-POST-TO-LEDGER`, and reads `AUDIT-IN`.
- **Rule 2.** When `WS-POSTING-MONTH` is not `WS-CLOSED-MONTH` and either `WS-POSTING-MONTH` is not `WS-NEXT-MONTH` or `WS-POSTING-CODE` is none of `'IN'`, `'OD'`, and `'MF'` (not `WS-MONTH-END-CODE`), it sets `WS-POSTING-MONTH` to `AU-BUSINESS-DATE(1:6)`, sets `WS-POSTING-CODE` to `AU-TRAN-CODE`, and reads `AUDIT-IN`.

Source: [GLSUM.cbl:207](../../../pangaea/programs/GLSUM.cbl#L207).

## 2100-POST-TO-LEDGER

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `GS-CODE(GS-CX) = AU-TRAN-CODE` | Y | Y | N |
| `AU-AMOUNT < 0` | Y | N | |
| → `SET GS-CX TO 1` | 1 | 1 | 1 |
| → `SEARCH GS-CODE-ENTRY` | 2 | 2 | 2 |
| → `CONTINUE` | 3 | 3 | |
| → `COMPUTE WS-AMOUNT = 0 - AU-AMOUNT` | 4 | | |
| → `MOVE WS-DEPOSITS-ACCOUNT TO WS-DEBIT-ACCOUNT` | 5 | | |
| → `MOVE GS-CODE-GL-ACCOUNT(GS-CX) TO WS-CREDIT-ACCOUNT` | 6 | | |
| → `MOVE WS-DEBIT-ACCOUNT TO WS-WANTED-ACCOUNT` | 7 | 7 | |
| → `PERFORM 2200-FIND-ACCOUNT` | 8, 11 | 8, 11 | |
| → `ADD WS-AMOUNT TO WS-LT-DEBITS(WS-TX)` | 9 | 9 | |
| → `MOVE WS-CREDIT-ACCOUNT TO WS-WANTED-ACCOUNT` | 10 | 10 | |
| → `ADD WS-AMOUNT TO WS-LT-CREDITS(WS-TX)` | 12 | 12 | |
| → `ADD 1 TO WS-POSTING-COUNT` | 13 | 13 | |
| → `MOVE AU-AMOUNT TO WS-AMOUNT` | | 4 | |
| → `MOVE GS-CODE-GL-ACCOUNT(GS-CX) TO WS-DEBIT-ACCOUNT` | | 5 | |
| → `MOVE WS-DEPOSITS-ACCOUNT TO WS-CREDIT-ACCOUNT` | | 6 | |
| → `MOVE 'UNKNOWN TRANSACTION CODE' TO WS-ABORT-MESSAGE` | | | 3 |
| → `PERFORM 9900-ABORT` | | | 4 |

- **Rule 1.** When `GS-CODE(GS-CX)` is `AU-TRAN-CODE` and `AU-AMOUNT` is less than 0, it sets `GS-CX` to 1, searches `GS-CODE-ENTRY`, computes `WS-AMOUNT` as `0 - AU-AMOUNT`, sets `WS-DEBIT-ACCOUNT` to `WS-DEPOSITS-ACCOUNT`, sets `WS-CREDIT-ACCOUNT` to `GS-CODE-GL-ACCOUNT(GS-CX)`, sets `WS-WANTED-ACCOUNT` to `WS-DEBIT-ACCOUNT`, performs `2200-FIND-ACCOUNT`, adds `WS-AMOUNT` to `WS-LT-DEBITS(WS-TX)`, sets `WS-WANTED-ACCOUNT` to `WS-CREDIT-ACCOUNT`, performs `2200-FIND-ACCOUNT`, adds `WS-AMOUNT` to `WS-LT-CREDITS(WS-TX)`, and adds 1 to `WS-POSTING-COUNT`.
- **Rule 2.** When `GS-CODE(GS-CX)` is `AU-TRAN-CODE` and `AU-AMOUNT` is at least 0, it sets `GS-CX` to 1, searches `GS-CODE-ENTRY`, sets `WS-AMOUNT` to `AU-AMOUNT`, sets `WS-DEBIT-ACCOUNT` to `GS-CODE-GL-ACCOUNT(GS-CX)`, sets `WS-CREDIT-ACCOUNT` to `WS-DEPOSITS-ACCOUNT`, sets `WS-WANTED-ACCOUNT` to `WS-DEBIT-ACCOUNT`, performs `2200-FIND-ACCOUNT`, adds `WS-AMOUNT` to `WS-LT-DEBITS(WS-TX)`, sets `WS-WANTED-ACCOUNT` to `WS-CREDIT-ACCOUNT`, performs `2200-FIND-ACCOUNT`, adds `WS-AMOUNT` to `WS-LT-CREDITS(WS-TX)`, and adds 1 to `WS-POSTING-COUNT`.
- **Rule 3.** When `GS-CODE(GS-CX)` is not `AU-TRAN-CODE`, it sets `GS-CX` to 1, searches `GS-CODE-ENTRY`, sets `WS-ABORT-MESSAGE` to `'UNKNOWN TRANSACTION CODE'`, and performs `9900-ABORT`.

Source: [GLSUM.cbl:224](../../../pangaea/programs/GLSUM.cbl#L224).

## 2200-FIND-ACCOUNT

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-LE-ACCOUNT(WS-LX) = WS-WANTED-ACCOUNT` | Y | N |
| → `SET WS-LX TO 1` | 1 | 1 |
| → `SEARCH WS-LEDGER-ENTRY` | 2 | 2 |
| → `SET WS-TX TO WS-LX` | 3 | |
| → `MOVE 'UNKNOWN LEDGER ACCOUNT' TO WS-ABORT-MESSAGE` | | 3 |
| → `PERFORM 9900-ABORT` | | 4 |

- **Rule 1.** When `WS-LE-ACCOUNT(WS-LX)` is `WS-WANTED-ACCOUNT`, it sets `WS-LX` to 1, searches `WS-LEDGER-ENTRY`, and sets `WS-TX` to `WS-LX`.
- **Rule 2.** When `WS-LE-ACCOUNT(WS-LX)` is not `WS-WANTED-ACCOUNT`, it sets `WS-LX` to 1, searches `WS-LEDGER-ENTRY`, sets `WS-ABORT-MESSAGE` to `'UNKNOWN LEDGER ACCOUNT'`, and performs `9900-ABORT`.

Source: [GLSUM.cbl:251](../../../pangaea/programs/GLSUM.cbl#L251).

## 3000-WRITE-LEDGER

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-DEBIT-TOTAL = WS-CREDIT-TOTAL` | N | Y |
| → `PERFORM VARYING WS-LX FROM 1 BY 1 UNTIL WS-LX > WS-LEDGER-COUNT` | 1 | 1 |
| → `MOVE 'LEDGER OUT OF BALANCE' TO WS-ABORT-MESSAGE` | 2 | |
| → `PERFORM 9900-ABORT` | 3 | |

- **Rule 1.** When `WS-DEBIT-TOTAL` is not `WS-CREDIT-TOTAL`, it repeats its loop varying `WS-LX` from 1 by 1 until `WS-LX` is more than `WS-LEDGER-COUNT`, sets `WS-ABORT-MESSAGE` to `'LEDGER OUT OF BALANCE'`, and performs `9900-ABORT`.
- **Rule 2.** When `WS-DEBIT-TOTAL` is `WS-CREDIT-TOTAL`, it repeats its loop varying `WS-LX` from 1 by 1 until `WS-LX` is more than `WS-LEDGER-COUNT`.

Source: [GLSUM.cbl:263](../../../pangaea/programs/GLSUM.cbl#L263).

## 3000-WRITE-LEDGER, in its loop

Each pass of `PERFORM VARYING WS-LX FROM 1 BY 1 UNTIL WS-LX > WS-LEDGER-COUNT`.

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `WS-LE-ACCOUNT(WS-LX) = WS-DEPOSITS-ACCOUNT` | Y | N | N |
| `WS-LE-ACCOUNT(WS-LX) = WS-TRANSIT-ACCOUNT` | | Y | N |
| → `SET WS-TX TO WS-LX` | 1 | 1 | 1 |
| → `MOVE SPACES TO GL-SUMMARY-RECORD` | 2 | 2 | 2 |
| → `MOVE WS-CLOSED-MONTH-NUMBER TO GL-MONTH` | 3 | 3 | 3 |
| → `MOVE WS-LE-ACCOUNT(WS-LX) TO GL-ACCOUNT` | 4 | 4 | 4 |
| → `MOVE WS-LE-NAME(WS-LX) TO GL-ACCOUNT-NAME` | 5 | 5 | 5 |
| → `MOVE WS-LT-DEBITS(WS-TX) TO GL-DEBITS` | 6 | 6 | 6 |
| → `MOVE WS-LT-CREDITS(WS-TX) TO GL-CREDITS` | 7 | 7 | 7 |
| → `WRITE GL-SUMMARY-RECORD` | 8 | 8 | 8 |
| → `ADD WS-LT-DEBITS(WS-TX) TO WS-DEBIT-TOTAL` | 9 | 9 | 9 |
| → `ADD WS-LT-CREDITS(WS-TX) TO WS-CREDIT-TOTAL` | 10 | 10 | 10 |
| → `COMPUTE WS-DEPOSITS-CHANGE = WS-LT-CREDITS(WS-TX) - WS-LT-DEBITS(WS-TX)` | 11 | | |
| → `COMPUTE WS-TRANSIT-NET = WS-LT-DEBITS(WS-TX) - WS-LT-CREDITS(WS-TX)` | | 11 | |

- **Rule 1.** When `WS-LE-ACCOUNT(WS-LX)` is `WS-DEPOSITS-ACCOUNT`, it sets `WS-TX` to `WS-LX`, sets `GL-SUMMARY-RECORD` to spaces, sets `GL-MONTH` to `WS-CLOSED-MONTH-NUMBER`, sets `GL-ACCOUNT` to `WS-LE-ACCOUNT(WS-LX)`, sets `GL-ACCOUNT-NAME` to `WS-LE-NAME(WS-LX)`, sets `GL-DEBITS` to `WS-LT-DEBITS(WS-TX)`, sets `GL-CREDITS` to `WS-LT-CREDITS(WS-TX)`, writes `GL-SUMMARY-RECORD`, adds `WS-LT-DEBITS(WS-TX)` to `WS-DEBIT-TOTAL`, adds `WS-LT-CREDITS(WS-TX)` to `WS-CREDIT-TOTAL`, and computes `WS-DEPOSITS-CHANGE` as `WS-LT-CREDITS(WS-TX) - WS-LT-DEBITS(WS-TX)`.
- **Rule 2.** When `WS-LE-ACCOUNT(WS-LX)` is not `WS-DEPOSITS-ACCOUNT` and `WS-LE-ACCOUNT(WS-LX)` is `WS-TRANSIT-ACCOUNT`, it sets `WS-TX` to `WS-LX`, sets `GL-SUMMARY-RECORD` to spaces, sets `GL-MONTH` to `WS-CLOSED-MONTH-NUMBER`, sets `GL-ACCOUNT` to `WS-LE-ACCOUNT(WS-LX)`, sets `GL-ACCOUNT-NAME` to `WS-LE-NAME(WS-LX)`, sets `GL-DEBITS` to `WS-LT-DEBITS(WS-TX)`, sets `GL-CREDITS` to `WS-LT-CREDITS(WS-TX)`, writes `GL-SUMMARY-RECORD`, adds `WS-LT-DEBITS(WS-TX)` to `WS-DEBIT-TOTAL`, adds `WS-LT-CREDITS(WS-TX)` to `WS-CREDIT-TOTAL`, and computes `WS-TRANSIT-NET` as `WS-LT-DEBITS(WS-TX) - WS-LT-CREDITS(WS-TX)`.
- **Rule 3.** When `WS-LE-ACCOUNT(WS-LX)` is not `WS-DEPOSITS-ACCOUNT` and `WS-LE-ACCOUNT(WS-LX)` is not `WS-TRANSIT-ACCOUNT`, it sets `WS-TX` to `WS-LX`, sets `GL-SUMMARY-RECORD` to spaces, sets `GL-MONTH` to `WS-CLOSED-MONTH-NUMBER`, sets `GL-ACCOUNT` to `WS-LE-ACCOUNT(WS-LX)`, sets `GL-ACCOUNT-NAME` to `WS-LE-NAME(WS-LX)`, sets `GL-DEBITS` to `WS-LT-DEBITS(WS-TX)`, sets `GL-CREDITS` to `WS-LT-CREDITS(WS-TX)`, writes `GL-SUMMARY-RECORD`, adds `WS-LT-DEBITS(WS-TX)` to `WS-DEBIT-TOTAL`, and adds `WS-LT-CREDITS(WS-TX)` to `WS-CREDIT-TOTAL`.

Source: [GLSUM.cbl:264](../../../pangaea/programs/GLSUM.cbl#L264).

## 8000-WRITE-REPORT

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `WS-CLOSED-MONTH-NUMBER = 0` | Y | N | N |
| `WS-TRANSIT-NET = 0` | | Y | N |
| → `MOVE 'GLSUM    GENERAL LEDGER SUMMARY' TO REPORT-LINE` | 1 | 1 | 1 |
| → `WRITE REPORT-LINE` | 2, 4, 6 | 2, 4, 7, 12, 21, 29, 31 | 2, 4, 7, 12, 21, 29, 31 |
| → `MOVE SPACES TO REPORT-LINE` | 3 | 3, 5, 11, 20, 28 | 3, 5, 11, 20, 28 |
| → `MOVE 'NO MONTH CLOSED' TO REPORT-LINE` | 5 | | |
| → `EXIT PARAGRAPH` | 7 | | |
| → `STRING 'MONTH CLOSED                  ' WS-CLOSED-MONTH DELIMITED BY SIZE INTO REPORT-LINE` | | 6 | 6 |
| → `MOVE 'POSTINGS' TO WS-RC-LABEL` | | 8 | 8 |
| → `MOVE WS-POSTING-COUNT TO WS-RC-VALUE` | | 9 | 9 |
| → `WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE` | | 10 | 10 |
| → `WRITE REPORT-LINE FROM WS-HEADING-LINE` | | 13 | 13 |
| → `PERFORM VARYING WS-LX FROM 1 BY 1 UNTIL WS-LX > WS-LEDGER-COUNT` | | 14 | 14 |
| → `MOVE SPACES TO WS-LL-ACCOUNT` | | 15 | 15 |
| → `MOVE 'TOTAL' TO WS-LL-NAME` | | 16 | 16 |
| → `MOVE WS-DEBIT-TOTAL TO WS-LL-DEBITS` | | 17 | 17 |
| → `MOVE WS-CREDIT-TOTAL TO WS-LL-CREDITS` | | 18 | 18 |
| → `WRITE REPORT-LINE FROM WS-LEDGER-LINE` | | 19 | 19 |
| → `MOVE 'CUSTOMER DEPOSITS CHANGE' TO WS-RA-LABEL` | | 22 | 22 |
| → `MOVE WS-DEPOSITS-CHANGE TO WS-RA-VALUE` | | 23 | 23 |
| → `WRITE REPORT-LINE FROM WS-REPORT-AMOUNT-LINE` | | 24, 27 | 24, 27 |
| → `MOVE 'TRANSFERS IN TRANSIT NET' TO WS-RA-LABEL` | | 25 | 25 |
| → `MOVE WS-TRANSIT-NET TO WS-RA-VALUE` | | 26 | 26 |
| → `MOVE 'LEDGER IN BALANCE' TO REPORT-LINE` | | 30 | |
| → `MOVE 'TRANSFERS IN TRANSIT DO NOT NET TO ZERO' TO REPORT-LINE` | | | 30 |

- **Rule 1.** When `WS-CLOSED-MONTH-NUMBER` is 0, it sets `REPORT-LINE` to `'GLSUM    GENERAL LEDGER SUMMARY'`, writes `REPORT-LINE`, sets `REPORT-LINE` to spaces, writes `REPORT-LINE`, sets `REPORT-LINE` to `'NO MONTH CLOSED'`, writes `REPORT-LINE`, and leaves the paragraph.
- **Rule 2.** When `WS-CLOSED-MONTH-NUMBER` is not 0 and `WS-TRANSIT-NET` is 0, it sets `REPORT-LINE` to `'GLSUM    GENERAL LEDGER SUMMARY'`, writes `REPORT-LINE`, sets `REPORT-LINE` to spaces, writes `REPORT-LINE`, sets `REPORT-LINE` to spaces, builds `REPORT-LINE` from `'MONTH CLOSED                  '` `WS-CLOSED-MONTH`, writes `REPORT-LINE`, sets `WS-RC-LABEL` to `'POSTINGS'`, sets `WS-RC-VALUE` to `WS-POSTING-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `REPORT-LINE` to spaces, writes `REPORT-LINE`, writes `REPORT-LINE` from `WS-HEADING-LINE`, repeats its loop varying `WS-LX` from 1 by 1 until `WS-LX` is more than `WS-LEDGER-COUNT`, sets `WS-LL-ACCOUNT` to spaces, sets `WS-LL-NAME` to `'TOTAL'`, sets `WS-LL-DEBITS` to `WS-DEBIT-TOTAL`, sets `WS-LL-CREDITS` to `WS-CREDIT-TOTAL`, writes `REPORT-LINE` from `WS-LEDGER-LINE`, sets `REPORT-LINE` to spaces, writes `REPORT-LINE`, sets `WS-RA-LABEL` to `'CUSTOMER DEPOSITS CHANGE'`, sets `WS-RA-VALUE` to `WS-DEPOSITS-CHANGE`, writes `REPORT-LINE` from `WS-REPORT-AMOUNT-LINE`, sets `WS-RA-LABEL` to `'TRANSFERS IN TRANSIT NET'`, sets `WS-RA-VALUE` to `WS-TRANSIT-NET`, writes `REPORT-LINE` from `WS-REPORT-AMOUNT-LINE`, sets `REPORT-LINE` to spaces, writes `REPORT-LINE`, sets `REPORT-LINE` to `'LEDGER IN BALANCE'`, and writes `REPORT-LINE`.
- **Rule 3.** When `WS-CLOSED-MONTH-NUMBER` is not 0 and `WS-TRANSIT-NET` is not 0, it sets `REPORT-LINE` to `'GLSUM    GENERAL LEDGER SUMMARY'`, writes `REPORT-LINE`, sets `REPORT-LINE` to spaces, writes `REPORT-LINE`, sets `REPORT-LINE` to spaces, builds `REPORT-LINE` from `'MONTH CLOSED                  '` `WS-CLOSED-MONTH`, writes `REPORT-LINE`, sets `WS-RC-LABEL` to `'POSTINGS'`, sets `WS-RC-VALUE` to `WS-POSTING-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `REPORT-LINE` to spaces, writes `REPORT-LINE`, writes `REPORT-LINE` from `WS-HEADING-LINE`, repeats its loop varying `WS-LX` from 1 by 1 until `WS-LX` is more than `WS-LEDGER-COUNT`, sets `WS-LL-ACCOUNT` to spaces, sets `WS-LL-NAME` to `'TOTAL'`, sets `WS-LL-DEBITS` to `WS-DEBIT-TOTAL`, sets `WS-LL-CREDITS` to `WS-CREDIT-TOTAL`, writes `REPORT-LINE` from `WS-LEDGER-LINE`, sets `REPORT-LINE` to spaces, writes `REPORT-LINE`, sets `WS-RA-LABEL` to `'CUSTOMER DEPOSITS CHANGE'`, sets `WS-RA-VALUE` to `WS-DEPOSITS-CHANGE`, writes `REPORT-LINE` from `WS-REPORT-AMOUNT-LINE`, sets `WS-RA-LABEL` to `'TRANSFERS IN TRANSIT NET'`, sets `WS-RA-VALUE` to `WS-TRANSIT-NET`, writes `REPORT-LINE` from `WS-REPORT-AMOUNT-LINE`, sets `REPORT-LINE` to spaces, writes `REPORT-LINE`, sets `REPORT-LINE` to `'TRANSFERS IN TRANSIT DO NOT NET TO ZERO'`, and writes `REPORT-LINE`.

Source: [GLSUM.cbl:291](../../../pangaea/programs/GLSUM.cbl#L291).
