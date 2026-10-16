# STMTPRT: decision tables

The business rules of [STMTPRT](../programs/STMTPRT.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

## Contents

- [0000-MAIN](#0000-main)
- [1000-INITIALIZE](#1000-initialize)
- [2000-SELECT-POSTINGS](#2000-select-postings)
- [2000-SELECT-POSTINGS, in its loop](#2000-select-postings-in-its-loop)
- [3000-PRINT-STATEMENTS](#3000-print-statements)
- [3000-PRINT-STATEMENTS, in its loop](#3000-print-statements-in-its-loop)
- [3100-RETURN-POSTING](#3100-return-posting)
- [3200-READ-MASTER](#3200-read-master)
- [3300-PRINT-STATEMENT](#3300-print-statement)
- [3400-PRINT-HEADING](#3400-print-heading)
- [3500-PRINT-POSTING](#3500-print-posting)
- [8000-WRITE-REPORT](#8000-write-report)

## 0000-MAIN

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `WS-CLOSED-MONTH-NUMBER = 0` | N | N | Y |
| `SORT-RETURN = 0` | N | Y | |
| → `PERFORM 1000-INITIALIZE` | 1 | 1 | 1 |
| → `SORT STMT-SORT ON ASCENDING KEY SR-ACCT-NO SR-DATE SR-SEQ INPUT PROCEDURE IS 2000-SELECT-POSTINGS OUTPUT PROCEDURE IS 3000-PRINT-STATEMENTS` | 2 | 2 | |
| → `MOVE 'SORT FAILED' TO WS-ABORT-MESSAGE` | 3 | | |
| → `PERFORM 9900-ABORT` | 4 | | |
| → `PERFORM 8000-WRITE-REPORT` | 5 | 3 | 2 |
| → `PERFORM 9000-CLOSE-FILES` | 6 | 4 | 3 |
| → `MOVE 0 TO RETURN-CODE` | 7 | 5 | 4 |
| → `STOP RUN` | 8 | 6 | 5 |

Source: [STMTPRT.cbl:208](../../../pangaea/programs/STMTPRT.cbl#L208).

## 1000-INITIALIZE

| | 1 | 2 | 3 | 4 | 5 |
| --- | :-: | :-: | :-: | :-: | :-: |
| `WS-FILE-STATUSES = ALL '0'` | N | Y | Y | Y | Y |
| `READ ACCR-IN AT END` | | N | N | Y | Y |
| `WS-CM-MONTH = 12` | | Y | N | Y | N |
| → `OPEN INPUT ACCR-IN ACCT-IN AUDIT-IN OUTPUT STMT-OUT REPORT-OUT` | 1 | 1 | 1 | 1 | 1 |
| → `DISPLAY 'STMTPRT: CANNOT OPEN FILES ' WS-FILE-STATUSES` | 2 | | | | |
| → `MOVE 16 TO RETURN-CODE` | 3 | | | | |
| → `STOP RUN` | 4 | | | | |
| → `MOVE 0 TO WS-CLOSED-MONTH-NUMBER` | | 2 | 2 | 2 | 2 |
| → `READ ACCR-IN` | | 3 | 3 | 3 | 3 |
| → `MOVE AC-DUE-MONTH TO WS-CLOSED-MONTH-NUMBER` | | 4 | 4 | | |
| → `MOVE WS-CM-YEAR TO WS-NM-YEAR` | | 5 | 5 | 4 | 4 |
| → `ADD 1 TO WS-NM-YEAR` | | 6 | | 5 | |
| → `MOVE 1 TO WS-NM-MONTH` | | 7 | | 6 | |
| → `COMPUTE WS-NM-MONTH = WS-CM-MONTH + 1` | | | 6 | | 5 |

Source: [STMTPRT.cbl:226](../../../pangaea/programs/STMTPRT.cbl#L226).

## 2000-SELECT-POSTINGS

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-AUDIT-STATUS = '10'` | N | Y |
| → `READ AUDIT-IN` | 1 | 1 |
| → `PERFORM UNTIL WS-AUDIT-STATUS NOT = '00'` | 2 | 2 |
| → `MOVE 'CANNOT READ THE AUDIT TRAILS' TO WS-ABORT-MESSAGE` | 3 | |
| → `PERFORM 9900-ABORT` | 4 | |

Source: [STMTPRT.cbl:249](../../../pangaea/programs/STMTPRT.cbl#L249).

## 2000-SELECT-POSTINGS, in its loop

Each pass of `PERFORM UNTIL WS-AUDIT-STATUS NOT = '00'`.

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-POSTING-MONTH = WS-CLOSED-MONTH OR (WS-POSTING-MONTH = WS-NEXT-MONTH AND WS-MONTH-END-CODE)` | Y | N |
| → `ADD 1 TO WS-SEQ` | 1 | 1 |
| → `MOVE AU-BUSINESS-DATE(1:6) TO WS-POSTING-MONTH` | 2 | 2 |
| → `MOVE AU-TRAN-CODE TO WS-POSTING-CODE` | 3 | 3 |
| → `PERFORM 2100-RELEASE-POSTING` | 4 | |
| → `READ AUDIT-IN` | 5 | 4 |

Source: [STMTPRT.cbl:251](../../../pangaea/programs/STMTPRT.cbl#L251).

## 3000-PRINT-STATEMENTS

| | 1 | 2 |
| --- | :-: | :-: |
| `POSTINGS-AT-END` | N | Y |
| → `PERFORM 3100-RETURN-POSTING` | 1 | 1 |
| → `PERFORM 3200-READ-MASTER` | 2 | 2 |
| → `PERFORM UNTIL MASTER-AT-END` | 3 | 3 |
| → `MOVE 'POSTINGS FOR NO ACCOUNT' TO WS-ABORT-MESSAGE` | 4 | |
| → `PERFORM 9900-ABORT` | 5 | |

Source: [STMTPRT.cbl:282](../../../pangaea/programs/STMTPRT.cbl#L282).

## 3000-PRINT-STATEMENTS, in its loop

Each pass of `PERFORM UNTIL MASTER-AT-END`.

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-POSTING-KEY < AM-ACCT-NO` | Y | N |
| → `MOVE 'POSTINGS FOR NO ACCOUNT' TO WS-ABORT-MESSAGE` | 1 | |
| → `PERFORM 9900-ABORT` | 2 | |
| → `PERFORM 3300-PRINT-STATEMENT` | 3 | 1 |
| → `PERFORM 3200-READ-MASTER` | 4 | 2 |

Source: [STMTPRT.cbl:285](../../../pangaea/programs/STMTPRT.cbl#L285).

## 3100-RETURN-POSTING

| | 1 | 2 |
| --- | :-: | :-: |
| `RETURN STMT-SORT AT END` | Y | N |
| → `RETURN STMT-SORT` | 1 | 1 |
| → `SET POSTINGS-AT-END TO TRUE` | 2 | |
| → `MOVE HIGH-VALUES TO WS-POSTING-KEY` | 3 | |
| → `MOVE SR-ACCT-NO TO WS-POSTING-KEY` | | 2 |

Source: [STMTPRT.cbl:299](../../../pangaea/programs/STMTPRT.cbl#L299).

## 3200-READ-MASTER

| | 1 | 2 |
| --- | :-: | :-: |
| `READ ACCT-IN AT END` | Y | N |
| → `READ ACCT-IN` | 1 | 1 |
| → `SET MASTER-AT-END TO TRUE` | 2 | |
| → `ADD 1 TO WS-ACCOUNT-COUNT` | | 2 |

Source: [STMTPRT.cbl:309](../../../pangaea/programs/STMTPRT.cbl#L309).

## 3300-PRINT-STATEMENT

| | 1 | 2 | 3 | 4 |
| --- | :-: | :-: | :-: | :-: |
| `AM-CLOSED AND AM-BALANCE = 0 AND WS-POSTING-KEY NOT = AM-ACCT-NO` | Y | Y | N | N |
| `WS-POSTING-KEY = AM-ACCT-NO` | Y | N | Y | N |
| → `EXIT PARAGRAPH` | 1 | 1 | | |
| → `MOVE AM-BALANCE TO WS-OPENING` | 2 | 2 | 1 | 1 |
| → `MOVE SR-BALANCE-BEFORE TO WS-OPENING` | 3 | | 2 | |
| → `MOVE WS-OPENING TO WS-CLOSING` | 4 | 3 | 3 | 2 |
| → `MOVE 0 TO WS-CREDITS WS-DEBITS` | 5 | 4 | 4 | 3 |
| → `PERFORM 3400-PRINT-HEADING` | 6 | 5 | 5 | 4 |
| → `PERFORM UNTIL WS-POSTING-KEY NOT = AM-ACCT-NO` | 7 | 6 | 6 | 5 |
| → `PERFORM 3600-PRINT-CLOSING` | 8 | 7 | 7 | 6 |
| → `ADD 1 TO WS-STATEMENT-COUNT` | 9 | 8 | 8 | 7 |

Source: [STMTPRT.cbl:319](../../../pangaea/programs/STMTPRT.cbl#L319).

## 3400-PRINT-HEADING

| | 1 | 2 | 3 | 4 |
| --- | :-: | :-: | :-: | :-: |
| `AM-SAVINGS AND AM-FROZEN` | Y | N | N | N |
| `AM-SAVINGS` | | Y | N | N |
| `AM-FROZEN` | | | Y | N |
| → `WRITE STMT-LINE FROM WS-RULE-LINE` | 1 | 1 | 1 | 1 |
| → `MOVE SPACES TO WS-TL-PERIOD` | 2 | 2 | 2 | 2 |
| → `STRING 'STATEMENT FOR ' DELIMITED BY SIZE WS-MONTH-NAME(WS-CM-MONTH) DELIMITED BY SPACE ' ' WS-CM-YEAR DELIMITED BY SIZE INTO WS-TL-PERIOD` | 3 | 3 | 3 | 3 |
| → `WRITE STMT-LINE FROM WS-TITLE-LINE` | 4 | 4 | 4 | 4 |
| → `MOVE AM-ACCT-NO TO WS-AL-ACCT-NO` | 5 | 5 | 5 | 5 |
| → `MOVE AM-SHORT-NAME TO WS-AL-NAME` | 6 | 6 | 6 | 6 |
| → `MOVE 'SAVINGS ACCOUNT, FROZEN' TO WS-AL-KIND` | 7 | | | |
| → `WRITE STMT-LINE FROM WS-ACCOUNT-LINE` | 8 | 8 | 8 | 8 |
| → `MOVE SPACES TO STMT-LINE` | 9 | 9 | 9 | 9 |
| → `WRITE STMT-LINE` | 10 | 10 | 10 | 10 |
| → `WRITE STMT-LINE FROM WS-HEADING-LINE` | 11 | 11 | 11 | 11 |
| → `MOVE 'OPENING BALANCE' TO WS-BL-LABEL` | 12 | 12 | 12 | 12 |
| → `MOVE WS-OPENING TO WS-BL-BALANCE` | 13 | 13 | 13 | 13 |
| → `WRITE STMT-LINE FROM WS-BALANCE-LINE` | 14 | 14 | 14 | 14 |
| → `MOVE 'SAVINGS ACCOUNT' TO WS-AL-KIND` | | 7 | | |
| → `MOVE 'CURRENT ACCOUNT, FROZEN' TO WS-AL-KIND` | | | 7 | |
| → `MOVE 'CURRENT ACCOUNT' TO WS-AL-KIND` | | | | 7 |

Source: [STMTPRT.cbl:339](../../../pangaea/programs/STMTPRT.cbl#L339).

## 3500-PRINT-POSTING

| | 1 | 2 | 3 | 4 |
| --- | :-: | :-: | :-: | :-: |
| `ST-CODE(ST-CX) = SR-TRAN-CODE` | Y | Y | N | N |
| `SR-AMOUNT < 0` | Y | N | Y | N |
| → `MOVE SR-DATE TO WS-DATE` | 1 | 1 | 1 | 1 |
| → `MOVE WS-DATE-YEAR TO WS-PL-YEAR` | 2 | 2 | 2 | 2 |
| → `MOVE WS-DATE-MONTH TO WS-PL-MONTH` | 3 | 3 | 3 | 3 |
| → `MOVE WS-DATE-DAY TO WS-PL-DAY` | 4 | 4 | 4 | 4 |
| → `MOVE SR-TRAN-ID TO WS-PL-TRAN-ID` | 5 | 5 | 5 | 5 |
| → `SET ST-CX TO 1` | 6 | 6 | 6 | 6 |
| → `SEARCH ST-CODE-ENTRY` | 7 | 7 | 7 | 7 |
| → `MOVE ST-CODE-TEXT(ST-CX) TO WS-PL-DESCRIPTION` | 8 | 8 | | |
| → `MOVE SR-AMOUNT TO WS-PL-AMOUNT` | 9 | 9 | 9 | 9 |
| → `MOVE SR-BALANCE-AFTER TO WS-PL-BALANCE` | 10 | 10 | 10 | 10 |
| → `WRITE STMT-LINE FROM WS-POSTING-LINE` | 11 | 11 | 11 | 11 |
| → `SUBTRACT SR-AMOUNT FROM WS-DEBITS` | 12 | | 12 | |
| → `MOVE SR-BALANCE-AFTER TO WS-CLOSING` | 13 | 13 | 13 | 13 |
| → `ADD 1 TO WS-PRINTED-COUNT` | 14 | 14 | 14 | 14 |
| → `ADD SR-AMOUNT TO WS-CREDITS` | | 12 | | 12 |
| → `MOVE SR-TRAN-CODE TO WS-PL-DESCRIPTION` | | | 8 | 8 |

Source: [STMTPRT.cbl:368](../../../pangaea/programs/STMTPRT.cbl#L368).

## 8000-WRITE-REPORT

| | 1 | 2 | 3 | 4 |
| --- | :-: | :-: | :-: | :-: |
| `WS-CLOSED-MONTH-NUMBER = 0` | Y | Y | N | N |
| `WS-PRINTED-COUNT = WS-SELECTED-COUNT` | N | Y | N | Y |
| → `MOVE 'STMTPRT  MONTHLY STATEMENTS' TO REPORT-LINE` | 1 | 1 | 1 | 1 |
| → `WRITE REPORT-LINE` | 2, 4, 6, 10 | 2, 4, 6, 10 | 2, 4, 7 | 2, 4, 7 |
| → `MOVE SPACES TO REPORT-LINE` | 3, 8 | 3, 8 | 3, 5 | 3, 5 |
| → `MOVE 'NO MONTH CLOSED' TO REPORT-LINE` | 5 | 5 | | |
| → `EXIT PARAGRAPH` | 7 | 7 | | |
| → `STRING 'MONTH CLOSED                  ' WS-CLOSED-MONTH DELIMITED BY SIZE INTO REPORT-LINE` | 9 | 9 | 6 | 6 |
| → `MOVE 'ACCOUNTS READ' TO WS-RC-LABEL` | 11 | 11 | 8 | 8 |
| → `MOVE WS-ACCOUNT-COUNT TO WS-RC-VALUE` | 12 | 12 | 9 | 9 |
| → `WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE` | 13, 16, 19, 22 | 13, 16, 19, 22 | 10, 13, 16, 19 | 10, 13, 16, 19 |
| → `MOVE 'STATEMENTS PRINTED' TO WS-RC-LABEL` | 14 | 14 | 11 | 11 |
| → `MOVE WS-STATEMENT-COUNT TO WS-RC-VALUE` | 15 | 15 | 12 | 12 |
| → `MOVE 'POSTINGS SELECTED' TO WS-RC-LABEL` | 17 | 17 | 14 | 14 |
| → `MOVE WS-SELECTED-COUNT TO WS-RC-VALUE` | 18 | 18 | 15 | 15 |
| → `MOVE 'POSTINGS PRINTED' TO WS-RC-LABEL` | 20 | 20 | 17 | 17 |
| → `MOVE WS-PRINTED-COUNT TO WS-RC-VALUE` | 21 | 21 | 18 | 18 |
| → `MOVE 'CREDITS' TO WS-RA-LABEL` | 23 | 23 | 20 | 20 |
| → `MOVE WS-CREDIT-TOTAL TO WS-RA-VALUE` | 24 | 24 | 21 | 21 |
| → `WRITE REPORT-LINE FROM WS-REPORT-AMOUNT-LINE` | 25, 28 | 25, 28 | 22, 25 | 22, 25 |
| → `MOVE 'DEBITS' TO WS-RA-LABEL` | 26 | 26 | 23 | 23 |
| → `MOVE WS-DEBIT-TOTAL TO WS-RA-VALUE` | 27 | 27 | 24 | 24 |
| → `MOVE 'POSTINGS SELECTED BUT NOT PRINTED' TO WS-ABORT-MESSAGE` | 29 | | 26 | |
| → `PERFORM 9900-ABORT` | 30 | | 27 | |

Source: [STMTPRT.cbl:407](../../../pangaea/programs/STMTPRT.cbl#L407).
