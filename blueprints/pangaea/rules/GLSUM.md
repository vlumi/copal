# GLSUM: decision tables

The business rules of [GLSUM](../programs/GLSUM.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

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

| | 1 | 2 | 3 | 4 |
| --- | :-: | :-: | :-: | :-: |
| `WS-CLOSED-MONTH-NUMBER = 0` | N | N | Y | Y |
| `WS-TRANSIT-NET = 0` | N | Y | N | Y |
| → `PERFORM 1000-INITIALIZE` | 1 | 1 | 1 | 1 |
| → `PERFORM 2000-POST-MONTH` | 2 | 2 | | |
| → `PERFORM 3000-WRITE-LEDGER` | 3 | 3 | | |
| → `PERFORM 8000-WRITE-REPORT` | 4 | 4 | 2 | 2 |
| → `PERFORM 9000-CLOSE-FILES` | 5 | 5 | 3 | 3 |
| → `MOVE 4 TO RETURN-CODE` | 6 | | 4 | |
| → `STOP RUN` | 7 | 7 | 5 | 5 |
| → `MOVE 0 TO RETURN-CODE` | | 6 | | 4 |

Source: [GLSUM.cbl:161](../../../pangaea/programs/GLSUM.cbl#L161).

## 1000-INITIALIZE

| | 1 | 2 | 3 | 4 | 5 |
| --- | :-: | :-: | :-: | :-: | :-: |
| `WS-FILE-STATUSES = ALL '0'` | N | Y | Y | Y | Y |
| `READ ACCR-IN AT END` | | N | N | Y | Y |
| `WS-CM-MONTH = 12` | | Y | N | Y | N |
| → `OPEN INPUT ACCR-IN AUDIT-IN OUTPUT GL-OUT REPORT-OUT` | 1 | 1 | 1 | 1 | 1 |
| → `DISPLAY 'GLSUM: CANNOT OPEN FILES ' WS-FILE-STATUSES` | 2 | | | | |
| → `MOVE 16 TO RETURN-CODE` | 3 | | | | |
| → `STOP RUN` | 4 | | | | |
| → `MOVE 0 TO WS-CLOSED-MONTH-NUMBER` | | 2 | 2 | 2 | 2 |
| → `READ ACCR-IN` | | 3 | 3 | 3 | 3 |
| → `MOVE AC-DUE-MONTH TO WS-CLOSED-MONTH-NUMBER` | | 4 | 4 | | |
| → `MOVE WS-CM-YEAR TO WS-NM-YEAR` | | 5 | 5 | 4 | 4 |
| → `ADD 1 TO WS-NM-YEAR` | | 6 | | 5 | |
| → `MOVE 1 TO WS-NM-MONTH` | | 7 | | 6 | |
| → `PERFORM VARYING WS-TX FROM 1 BY 1 UNTIL WS-TX > WS-LEDGER-COUNT` | | 8 | 7 | 7 | 6 |
| → `COMPUTE WS-NM-MONTH = WS-CM-MONTH + 1` | | | 6 | | 5 |

Source: [GLSUM.cbl:178](../../../pangaea/programs/GLSUM.cbl#L178).

## 2000-POST-MONTH

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-AUDIT-STATUS = '10'` | N | Y |
| → `READ AUDIT-IN` | 1 | 1 |
| → `PERFORM UNTIL WS-AUDIT-STATUS NOT = '00'` | 2 | 2 |
| → `MOVE 'CANNOT READ THE AUDIT TRAILS' TO WS-ABORT-MESSAGE` | 3 | |
| → `PERFORM 9900-ABORT` | 4 | |

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

Source: [GLSUM.cbl:251](../../../pangaea/programs/GLSUM.cbl#L251).

## 3000-WRITE-LEDGER

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-DEBIT-TOTAL = WS-CREDIT-TOTAL` | N | Y |
| → `PERFORM VARYING WS-LX FROM 1 BY 1 UNTIL WS-LX > WS-LEDGER-COUNT` | 1 | 1 |
| → `MOVE 'LEDGER OUT OF BALANCE' TO WS-ABORT-MESSAGE` | 2 | |
| → `PERFORM 9900-ABORT` | 3 | |

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

Source: [GLSUM.cbl:291](../../../pangaea/programs/GLSUM.cbl#L291).
