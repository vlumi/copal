# STMTPRT: decision tables

The business rules of [STMTPRT](../programs/STMTPRT.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

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
| → `PERFORM 8000-WRITE-REPORT` | | 3 | 2 |
| → `PERFORM 9000-CLOSE-FILES` | | 4 | 3 |
| → `MOVE 0 TO RETURN-CODE` | | 5 | 4 |
| → `STOP RUN` | | 6 | 5 |

- **Rule 1.** When `WS-CLOSED-MONTH-NUMBER` is not 0 and `SORT-RETURN` is not 0, it performs `1000-INITIALIZE`, runs `SORT STMT-SORT ON ASCENDING KEY SR-ACCT-NO SR-DATE SR-SEQ INPUT PROCEDURE IS 2000-SELECT-POSTINGS OUTPUT PROCEDURE IS 3000-PRINT-STATEMENTS`, sets `WS-ABORT-MESSAGE` to `'SORT FAILED'`, and performs `9900-ABORT`.
- **Rule 2.** When `WS-CLOSED-MONTH-NUMBER` is not 0 and `SORT-RETURN` is 0, it performs `1000-INITIALIZE`, runs `SORT STMT-SORT ON ASCENDING KEY SR-ACCT-NO SR-DATE SR-SEQ INPUT PROCEDURE IS 2000-SELECT-POSTINGS OUTPUT PROCEDURE IS 3000-PRINT-STATEMENTS`, performs `8000-WRITE-REPORT`, performs `9000-CLOSE-FILES`, sets `RETURN-CODE` to 0, and stops the run.
- **Rule 3.** When `WS-CLOSED-MONTH-NUMBER` is 0, it performs `1000-INITIALIZE`, performs `8000-WRITE-REPORT`, performs `9000-CLOSE-FILES`, sets `RETURN-CODE` to 0, and stops the run.

Source: [STMTPRT.cbl:208](../../../pangaea/programs/STMTPRT.cbl#L208).

## 1000-INITIALIZE

| | 1 | 2 | 3 | 4 |
| --- | :-: | :-: | :-: | :-: |
| `WS-FILE-STATUSES = ALL '0'` | N | Y | Y | Y |
| `READ ACCR-IN AT END` | | N | N | Y |
| `WS-CM-MONTH = 12` | | Y | N | N |
| → `OPEN INPUT ACCR-IN ACCT-IN AUDIT-IN OUTPUT STMT-OUT REPORT-OUT` | 1 | 1 | 1 | 1 |
| → `DISPLAY 'STMTPRT: CANNOT OPEN FILES ' WS-FILE-STATUSES` | 2 | | | |
| → `MOVE 16 TO RETURN-CODE` | 3 | | | |
| → `STOP RUN` | 4 | | | |
| → `MOVE 0 TO WS-CLOSED-MONTH-NUMBER` | | 2 | 2 | 2 |
| → `READ ACCR-IN` | | 3 | 3 | 3 |
| → `MOVE AC-DUE-MONTH TO WS-CLOSED-MONTH-NUMBER` | | 4 | 4 | |
| → `MOVE WS-CM-YEAR TO WS-NM-YEAR` | | 5 | 5 | 4 |
| → `ADD 1 TO WS-NM-YEAR` | | 6 | | |
| → `MOVE 1 TO WS-NM-MONTH` | | 7 | | |
| → `COMPUTE WS-NM-MONTH = WS-CM-MONTH + 1` | | | 6 | 5 |

- **Rule 1.** When `WS-FILE-STATUSES` is not `ALL '0'`, it opens `ACCR-IN`, `ACCT-IN`, and `AUDIT-IN` for input and `STMT-OUT` and `REPORT-OUT` for output; displays `'STMTPRT: CANNOT OPEN FILES '` `WS-FILE-STATUSES`; sets `RETURN-CODE` to 16; and stops the run.
- **Rule 2.** When `WS-FILE-STATUSES` is `ALL '0'`, reading `ACCR-IN` finds a record, and `WS-CM-MONTH` is 12, it opens `ACCR-IN`, `ACCT-IN`, and `AUDIT-IN` for input and `STMT-OUT` and `REPORT-OUT` for output; sets `WS-CLOSED-MONTH-NUMBER` to 0; reads `ACCR-IN`; sets `WS-CLOSED-MONTH-NUMBER` to `AC-DUE-MONTH`; sets `WS-NM-YEAR` to `WS-CM-YEAR`; adds 1 to `WS-NM-YEAR`; and sets `WS-NM-MONTH` to 1.
- **Rule 3.** When `WS-FILE-STATUSES` is `ALL '0'`, reading `ACCR-IN` finds a record, and `WS-CM-MONTH` is not 12, it opens `ACCR-IN`, `ACCT-IN`, and `AUDIT-IN` for input and `STMT-OUT` and `REPORT-OUT` for output; sets `WS-CLOSED-MONTH-NUMBER` to 0; reads `ACCR-IN`; sets `WS-CLOSED-MONTH-NUMBER` to `AC-DUE-MONTH`; sets `WS-NM-YEAR` to `WS-CM-YEAR`; and computes `WS-NM-MONTH` as `WS-CM-MONTH + 1`.
- **Rule 4.** When `WS-FILE-STATUSES` is `ALL '0'`, reading `ACCR-IN` finds no more records, and `WS-CM-MONTH` is not 12, it opens `ACCR-IN`, `ACCT-IN`, and `AUDIT-IN` for input and `STMT-OUT` and `REPORT-OUT` for output; sets `WS-CLOSED-MONTH-NUMBER` to 0; reads `ACCR-IN`; sets `WS-NM-YEAR` to `WS-CM-YEAR`; and computes `WS-NM-MONTH` as `WS-CM-MONTH + 1`.

Source: [STMTPRT.cbl:226](../../../pangaea/programs/STMTPRT.cbl#L226).

## 2000-SELECT-POSTINGS

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-AUDIT-STATUS = '10'` | N | Y |
| → `READ AUDIT-IN` | 1 | 1 |
| → `PERFORM UNTIL WS-AUDIT-STATUS NOT = '00'` | 2 | 2 |
| → `MOVE 'CANNOT READ THE AUDIT TRAILS' TO WS-ABORT-MESSAGE` | 3 | |
| → `PERFORM 9900-ABORT` | 4 | |

- **Rule 1.** When `WS-AUDIT-STATUS` is not `'10'`, it reads `AUDIT-IN`, repeats its loop until `WS-AUDIT-STATUS` is not `'00'`, sets `WS-ABORT-MESSAGE` to `'CANNOT READ THE AUDIT TRAILS'`, and performs `9900-ABORT`.
- **Rule 2.** When `WS-AUDIT-STATUS` is `'10'`, it reads `AUDIT-IN` and repeats its loop until `WS-AUDIT-STATUS` is not `'00'`.

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

- **Rule 1.** When `WS-POSTING-MONTH` is `WS-CLOSED-MONTH` or both `WS-POSTING-MONTH` is `WS-NEXT-MONTH` and `WS-POSTING-CODE` is `'IN'`, `'OD'`, or `'MF'` (`WS-MONTH-END-CODE`), it adds 1 to `WS-SEQ`, sets `WS-POSTING-MONTH` to `AU-BUSINESS-DATE(1:6)`, sets `WS-POSTING-CODE` to `AU-TRAN-CODE`, performs `2100-RELEASE-POSTING`, and reads `AUDIT-IN`.
- **Rule 2.** When `WS-POSTING-MONTH` is not `WS-CLOSED-MONTH` and either `WS-POSTING-MONTH` is not `WS-NEXT-MONTH` or `WS-POSTING-CODE` is none of `'IN'`, `'OD'`, and `'MF'` (not `WS-MONTH-END-CODE`), it adds 1 to `WS-SEQ`, sets `WS-POSTING-MONTH` to `AU-BUSINESS-DATE(1:6)`, sets `WS-POSTING-CODE` to `AU-TRAN-CODE`, and reads `AUDIT-IN`.

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

- **Rule 1.** When `WS-POSTING-END-SW` is not `'Y'` (not `POSTINGS-AT-END`), it performs `3100-RETURN-POSTING`, performs `3200-READ-MASTER`, repeats its loop until `WS-MASTER-END-SW` is `'Y'` (`MASTER-AT-END`), sets `WS-ABORT-MESSAGE` to `'POSTINGS FOR NO ACCOUNT'`, and performs `9900-ABORT`.
- **Rule 2.** When `WS-POSTING-END-SW` is `'Y'` (`POSTINGS-AT-END`), it performs `3100-RETURN-POSTING`, performs `3200-READ-MASTER`, and repeats its loop until `WS-MASTER-END-SW` is `'Y'` (`MASTER-AT-END`).

Source: [STMTPRT.cbl:282](../../../pangaea/programs/STMTPRT.cbl#L282).

## 3000-PRINT-STATEMENTS, in its loop

Each pass of `PERFORM UNTIL MASTER-AT-END`.

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-POSTING-KEY < AM-ACCT-NO` | Y | N |
| → `MOVE 'POSTINGS FOR NO ACCOUNT' TO WS-ABORT-MESSAGE` | 1 | |
| → `PERFORM 9900-ABORT` | 2 | |
| → `PERFORM 3300-PRINT-STATEMENT` | | 1 |
| → `PERFORM 3200-READ-MASTER` | | 2 |

- **Rule 1.** When `WS-POSTING-KEY` is less than `AM-ACCT-NO`, it sets `WS-ABORT-MESSAGE` to `'POSTINGS FOR NO ACCOUNT'` and performs `9900-ABORT`.
- **Rule 2.** When `WS-POSTING-KEY` is at least `AM-ACCT-NO`, it performs `3300-PRINT-STATEMENT` and performs `3200-READ-MASTER`.

Source: [STMTPRT.cbl:285](../../../pangaea/programs/STMTPRT.cbl#L285).

## 3100-RETURN-POSTING

| | 1 | 2 |
| --- | :-: | :-: |
| `RETURN STMT-SORT AT END` | Y | N |
| → `RETURN STMT-SORT` | 1 | 1 |
| → `SET POSTINGS-AT-END TO TRUE` | 2 | |
| → `MOVE HIGH-VALUES TO WS-POSTING-KEY` | 3 | |
| → `MOVE SR-ACCT-NO TO WS-POSTING-KEY` | | 2 |

- **Rule 1.** When taking the next record from `STMT-SORT` finds no more records, it takes the next record from `STMT-SORT`, sets `WS-POSTING-END-SW` to `POSTINGS-AT-END`, and sets `WS-POSTING-KEY` to high values.
- **Rule 2.** When taking the next record from `STMT-SORT` finds a record, it takes the next record from `STMT-SORT` and sets `WS-POSTING-KEY` to `SR-ACCT-NO`.

Source: [STMTPRT.cbl:299](../../../pangaea/programs/STMTPRT.cbl#L299).

## 3200-READ-MASTER

| | 1 | 2 |
| --- | :-: | :-: |
| `READ ACCT-IN AT END` | Y | N |
| → `READ ACCT-IN` | 1 | 1 |
| → `SET MASTER-AT-END TO TRUE` | 2 | |
| → `ADD 1 TO WS-ACCOUNT-COUNT` | | 2 |

- **Rule 1.** When reading `ACCT-IN` finds no more records, it reads `ACCT-IN` and sets `WS-MASTER-END-SW` to `MASTER-AT-END`.
- **Rule 2.** When reading `ACCT-IN` finds a record, it reads `ACCT-IN` and adds 1 to `WS-ACCOUNT-COUNT`.

Source: [STMTPRT.cbl:309](../../../pangaea/programs/STMTPRT.cbl#L309).

## 3300-PRINT-STATEMENT

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `AM-CLOSED AND AM-BALANCE = 0 AND WS-POSTING-KEY NOT = AM-ACCT-NO` | Y | N | N |
| `WS-POSTING-KEY = AM-ACCT-NO` | | Y | N |
| → `EXIT PARAGRAPH` | 1 | | |
| → `MOVE AM-BALANCE TO WS-OPENING` | | 1 | 1 |
| → `MOVE SR-BALANCE-BEFORE TO WS-OPENING` | | 2 | |
| → `MOVE WS-OPENING TO WS-CLOSING` | | 3 | 2 |
| → `MOVE 0 TO WS-CREDITS WS-DEBITS` | | 4 | 3 |
| → `PERFORM 3400-PRINT-HEADING` | | 5 | 4 |
| → `PERFORM UNTIL WS-POSTING-KEY NOT = AM-ACCT-NO` | | 6 | 5 |
| → `PERFORM 3600-PRINT-CLOSING` | | 7 | 6 |
| → `ADD 1 TO WS-STATEMENT-COUNT` | | 8 | 7 |

- **Rule 1.** When `AM-STATUS` is `'X'` (`AM-CLOSED`), `AM-BALANCE` is 0, and `WS-POSTING-KEY` is not `AM-ACCT-NO`, it leaves the paragraph.
- **Rule 2.** When `AM-STATUS` is not `'X'` (not `AM-CLOSED`), `AM-BALANCE` is not 0, or `WS-POSTING-KEY` is `AM-ACCT-NO`; and `WS-POSTING-KEY` is `AM-ACCT-NO`, it sets `WS-OPENING` to `AM-BALANCE`; sets `WS-OPENING` to `SR-BALANCE-BEFORE`; sets `WS-CLOSING` to `WS-OPENING`; sets `WS-CREDITS` and `WS-DEBITS` to 0; performs `3400-PRINT-HEADING`; repeats its loop until `WS-POSTING-KEY` is not `AM-ACCT-NO`; performs `3600-PRINT-CLOSING`; and adds 1 to `WS-STATEMENT-COUNT`.
- **Rule 3.** When `AM-STATUS` is not `'X'` (not `AM-CLOSED`), `AM-BALANCE` is not 0, or `WS-POSTING-KEY` is `AM-ACCT-NO`; and `WS-POSTING-KEY` is not `AM-ACCT-NO`, it sets `WS-OPENING` to `AM-BALANCE`; sets `WS-CLOSING` to `WS-OPENING`; sets `WS-CREDITS` and `WS-DEBITS` to 0; performs `3400-PRINT-HEADING`; repeats its loop until `WS-POSTING-KEY` is not `AM-ACCT-NO`; performs `3600-PRINT-CLOSING`; and adds 1 to `WS-STATEMENT-COUNT`.

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

- **Rule 1.** When `AM-ACCT-TYPE` is `'S'` (`AM-SAVINGS`) and `AM-STATUS` is `'F'` (`AM-FROZEN`), it writes `STMT-LINE` from `WS-RULE-LINE`; sets `WS-TL-PERIOD` to spaces; builds `WS-TL-PERIOD` from `'STATEMENT FOR '` `WS-MONTH-NAME(WS-CM-MONTH)` spaces `' '` `WS-CM-YEAR`; writes `STMT-LINE` from `WS-TITLE-LINE`; sets `WS-AL-ACCT-NO` to `AM-ACCT-NO`; sets `WS-AL-NAME` to `AM-SHORT-NAME`; sets `WS-AL-KIND` to `'SAVINGS ACCOUNT, FROZEN'`; writes `STMT-LINE` from `WS-ACCOUNT-LINE`; sets `STMT-LINE` to spaces; writes `STMT-LINE`; writes `STMT-LINE` from `WS-HEADING-LINE`; sets `WS-BL-LABEL` to `'OPENING BALANCE'`; sets `WS-BL-BALANCE` to `WS-OPENING`; and writes `STMT-LINE` from `WS-BALANCE-LINE`.
- **Rule 2.** When `AM-ACCT-TYPE` is not `'S'` (not `AM-SAVINGS`) or `AM-STATUS` is not `'F'` (not `AM-FROZEN`); and `AM-ACCT-TYPE` is `'S'` (`AM-SAVINGS`), it writes `STMT-LINE` from `WS-RULE-LINE`, sets `WS-TL-PERIOD` to spaces, builds `WS-TL-PERIOD` from `'STATEMENT FOR '` `WS-MONTH-NAME(WS-CM-MONTH)` spaces `' '` `WS-CM-YEAR`, writes `STMT-LINE` from `WS-TITLE-LINE`, sets `WS-AL-ACCT-NO` to `AM-ACCT-NO`, sets `WS-AL-NAME` to `AM-SHORT-NAME`, sets `WS-AL-KIND` to `'SAVINGS ACCOUNT'`, writes `STMT-LINE` from `WS-ACCOUNT-LINE`, sets `STMT-LINE` to spaces, writes `STMT-LINE`, writes `STMT-LINE` from `WS-HEADING-LINE`, sets `WS-BL-LABEL` to `'OPENING BALANCE'`, sets `WS-BL-BALANCE` to `WS-OPENING`, and writes `STMT-LINE` from `WS-BALANCE-LINE`.
- **Rule 3.** When `AM-ACCT-TYPE` is not `'S'` (not `AM-SAVINGS`) or `AM-STATUS` is not `'F'` (not `AM-FROZEN`); `AM-ACCT-TYPE` is not `'S'` (not `AM-SAVINGS`); and `AM-STATUS` is `'F'` (`AM-FROZEN`), it writes `STMT-LINE` from `WS-RULE-LINE`; sets `WS-TL-PERIOD` to spaces; builds `WS-TL-PERIOD` from `'STATEMENT FOR '` `WS-MONTH-NAME(WS-CM-MONTH)` spaces `' '` `WS-CM-YEAR`; writes `STMT-LINE` from `WS-TITLE-LINE`; sets `WS-AL-ACCT-NO` to `AM-ACCT-NO`; sets `WS-AL-NAME` to `AM-SHORT-NAME`; sets `WS-AL-KIND` to `'CURRENT ACCOUNT, FROZEN'`; writes `STMT-LINE` from `WS-ACCOUNT-LINE`; sets `STMT-LINE` to spaces; writes `STMT-LINE`; writes `STMT-LINE` from `WS-HEADING-LINE`; sets `WS-BL-LABEL` to `'OPENING BALANCE'`; sets `WS-BL-BALANCE` to `WS-OPENING`; and writes `STMT-LINE` from `WS-BALANCE-LINE`.
- **Rule 4.** When `AM-ACCT-TYPE` is not `'S'` (not `AM-SAVINGS`) or `AM-STATUS` is not `'F'` (not `AM-FROZEN`); `AM-ACCT-TYPE` is not `'S'` (not `AM-SAVINGS`); and `AM-STATUS` is not `'F'` (not `AM-FROZEN`), it writes `STMT-LINE` from `WS-RULE-LINE`, sets `WS-TL-PERIOD` to spaces, builds `WS-TL-PERIOD` from `'STATEMENT FOR '` `WS-MONTH-NAME(WS-CM-MONTH)` spaces `' '` `WS-CM-YEAR`, writes `STMT-LINE` from `WS-TITLE-LINE`, sets `WS-AL-ACCT-NO` to `AM-ACCT-NO`, sets `WS-AL-NAME` to `AM-SHORT-NAME`, sets `WS-AL-KIND` to `'CURRENT ACCOUNT'`, writes `STMT-LINE` from `WS-ACCOUNT-LINE`, sets `STMT-LINE` to spaces, writes `STMT-LINE`, writes `STMT-LINE` from `WS-HEADING-LINE`, sets `WS-BL-LABEL` to `'OPENING BALANCE'`, sets `WS-BL-BALANCE` to `WS-OPENING`, and writes `STMT-LINE` from `WS-BALANCE-LINE`.

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

- **Rule 1.** When `ST-CODE(ST-CX)` is `SR-TRAN-CODE` and `SR-AMOUNT` is less than 0, it sets `WS-DATE` to `SR-DATE`, sets `WS-PL-YEAR` to `WS-DATE-YEAR`, sets `WS-PL-MONTH` to `WS-DATE-MONTH`, sets `WS-PL-DAY` to `WS-DATE-DAY`, sets `WS-PL-TRAN-ID` to `SR-TRAN-ID`, sets `ST-CX` to 1, searches `ST-CODE-ENTRY`, sets `WS-PL-DESCRIPTION` to `ST-CODE-TEXT(ST-CX)`, sets `WS-PL-AMOUNT` to `SR-AMOUNT`, sets `WS-PL-BALANCE` to `SR-BALANCE-AFTER`, writes `STMT-LINE` from `WS-POSTING-LINE`, subtracts `SR-AMOUNT` from `WS-DEBITS`, sets `WS-CLOSING` to `SR-BALANCE-AFTER`, and adds 1 to `WS-PRINTED-COUNT`.
- **Rule 2.** When `ST-CODE(ST-CX)` is `SR-TRAN-CODE` and `SR-AMOUNT` is at least 0, it sets `WS-DATE` to `SR-DATE`, sets `WS-PL-YEAR` to `WS-DATE-YEAR`, sets `WS-PL-MONTH` to `WS-DATE-MONTH`, sets `WS-PL-DAY` to `WS-DATE-DAY`, sets `WS-PL-TRAN-ID` to `SR-TRAN-ID`, sets `ST-CX` to 1, searches `ST-CODE-ENTRY`, sets `WS-PL-DESCRIPTION` to `ST-CODE-TEXT(ST-CX)`, sets `WS-PL-AMOUNT` to `SR-AMOUNT`, sets `WS-PL-BALANCE` to `SR-BALANCE-AFTER`, writes `STMT-LINE` from `WS-POSTING-LINE`, adds `SR-AMOUNT` to `WS-CREDITS`, sets `WS-CLOSING` to `SR-BALANCE-AFTER`, and adds 1 to `WS-PRINTED-COUNT`.
- **Rule 3.** When `ST-CODE(ST-CX)` is not `SR-TRAN-CODE` and `SR-AMOUNT` is less than 0, it sets `WS-DATE` to `SR-DATE`, sets `WS-PL-YEAR` to `WS-DATE-YEAR`, sets `WS-PL-MONTH` to `WS-DATE-MONTH`, sets `WS-PL-DAY` to `WS-DATE-DAY`, sets `WS-PL-TRAN-ID` to `SR-TRAN-ID`, sets `ST-CX` to 1, searches `ST-CODE-ENTRY`, sets `WS-PL-DESCRIPTION` to `SR-TRAN-CODE`, sets `WS-PL-AMOUNT` to `SR-AMOUNT`, sets `WS-PL-BALANCE` to `SR-BALANCE-AFTER`, writes `STMT-LINE` from `WS-POSTING-LINE`, subtracts `SR-AMOUNT` from `WS-DEBITS`, sets `WS-CLOSING` to `SR-BALANCE-AFTER`, and adds 1 to `WS-PRINTED-COUNT`.
- **Rule 4.** When `ST-CODE(ST-CX)` is not `SR-TRAN-CODE` and `SR-AMOUNT` is at least 0, it sets `WS-DATE` to `SR-DATE`, sets `WS-PL-YEAR` to `WS-DATE-YEAR`, sets `WS-PL-MONTH` to `WS-DATE-MONTH`, sets `WS-PL-DAY` to `WS-DATE-DAY`, sets `WS-PL-TRAN-ID` to `SR-TRAN-ID`, sets `ST-CX` to 1, searches `ST-CODE-ENTRY`, sets `WS-PL-DESCRIPTION` to `SR-TRAN-CODE`, sets `WS-PL-AMOUNT` to `SR-AMOUNT`, sets `WS-PL-BALANCE` to `SR-BALANCE-AFTER`, writes `STMT-LINE` from `WS-POSTING-LINE`, adds `SR-AMOUNT` to `WS-CREDITS`, sets `WS-CLOSING` to `SR-BALANCE-AFTER`, and adds 1 to `WS-PRINTED-COUNT`.

Source: [STMTPRT.cbl:368](../../../pangaea/programs/STMTPRT.cbl#L368).

## 8000-WRITE-REPORT

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `WS-CLOSED-MONTH-NUMBER = 0` | Y | N | N |
| `WS-PRINTED-COUNT = WS-SELECTED-COUNT` | | N | Y |
| → `MOVE 'STMTPRT  MONTHLY STATEMENTS' TO REPORT-LINE` | 1 | 1 | 1 |
| → `WRITE REPORT-LINE` | 2, 4, 6 | 2, 4, 7 | 2, 4, 7 |
| → `MOVE SPACES TO REPORT-LINE` | 3 | 3, 5 | 3, 5 |
| → `MOVE 'NO MONTH CLOSED' TO REPORT-LINE` | 5 | | |
| → `EXIT PARAGRAPH` | 7 | | |
| → `STRING 'MONTH CLOSED                  ' WS-CLOSED-MONTH DELIMITED BY SIZE INTO REPORT-LINE` | | 6 | 6 |
| → `MOVE 'ACCOUNTS READ' TO WS-RC-LABEL` | | 8 | 8 |
| → `MOVE WS-ACCOUNT-COUNT TO WS-RC-VALUE` | | 9 | 9 |
| → `WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE` | | 10, 13, 16, 19 | 10, 13, 16, 19 |
| → `MOVE 'STATEMENTS PRINTED' TO WS-RC-LABEL` | | 11 | 11 |
| → `MOVE WS-STATEMENT-COUNT TO WS-RC-VALUE` | | 12 | 12 |
| → `MOVE 'POSTINGS SELECTED' TO WS-RC-LABEL` | | 14 | 14 |
| → `MOVE WS-SELECTED-COUNT TO WS-RC-VALUE` | | 15 | 15 |
| → `MOVE 'POSTINGS PRINTED' TO WS-RC-LABEL` | | 17 | 17 |
| → `MOVE WS-PRINTED-COUNT TO WS-RC-VALUE` | | 18 | 18 |
| → `MOVE 'CREDITS' TO WS-RA-LABEL` | | 20 | 20 |
| → `MOVE WS-CREDIT-TOTAL TO WS-RA-VALUE` | | 21 | 21 |
| → `WRITE REPORT-LINE FROM WS-REPORT-AMOUNT-LINE` | | 22, 25 | 22, 25 |
| → `MOVE 'DEBITS' TO WS-RA-LABEL` | | 23 | 23 |
| → `MOVE WS-DEBIT-TOTAL TO WS-RA-VALUE` | | 24 | 24 |
| → `MOVE 'POSTINGS SELECTED BUT NOT PRINTED' TO WS-ABORT-MESSAGE` | | 26 | |
| → `PERFORM 9900-ABORT` | | 27 | |

- **Rule 1.** When `WS-CLOSED-MONTH-NUMBER` is 0, it sets `REPORT-LINE` to `'STMTPRT  MONTHLY STATEMENTS'`, writes `REPORT-LINE`, sets `REPORT-LINE` to spaces, writes `REPORT-LINE`, sets `REPORT-LINE` to `'NO MONTH CLOSED'`, writes `REPORT-LINE`, and leaves the paragraph.
- **Rule 2.** When `WS-CLOSED-MONTH-NUMBER` is not 0 and `WS-PRINTED-COUNT` is not `WS-SELECTED-COUNT`, it sets `REPORT-LINE` to `'STMTPRT  MONTHLY STATEMENTS'`, writes `REPORT-LINE`, sets `REPORT-LINE` to spaces, writes `REPORT-LINE`, sets `REPORT-LINE` to spaces, builds `REPORT-LINE` from `'MONTH CLOSED                  '` `WS-CLOSED-MONTH`, writes `REPORT-LINE`, sets `WS-RC-LABEL` to `'ACCOUNTS READ'`, sets `WS-RC-VALUE` to `WS-ACCOUNT-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `WS-RC-LABEL` to `'STATEMENTS PRINTED'`, sets `WS-RC-VALUE` to `WS-STATEMENT-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `WS-RC-LABEL` to `'POSTINGS SELECTED'`, sets `WS-RC-VALUE` to `WS-SELECTED-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `WS-RC-LABEL` to `'POSTINGS PRINTED'`, sets `WS-RC-VALUE` to `WS-PRINTED-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `WS-RA-LABEL` to `'CREDITS'`, sets `WS-RA-VALUE` to `WS-CREDIT-TOTAL`, writes `REPORT-LINE` from `WS-REPORT-AMOUNT-LINE`, sets `WS-RA-LABEL` to `'DEBITS'`, sets `WS-RA-VALUE` to `WS-DEBIT-TOTAL`, writes `REPORT-LINE` from `WS-REPORT-AMOUNT-LINE`, sets `WS-ABORT-MESSAGE` to `'POSTINGS SELECTED BUT NOT PRINTED'`, and performs `9900-ABORT`.
- **Rule 3.** When `WS-CLOSED-MONTH-NUMBER` is not 0 and `WS-PRINTED-COUNT` is `WS-SELECTED-COUNT`, it sets `REPORT-LINE` to `'STMTPRT  MONTHLY STATEMENTS'`, writes `REPORT-LINE`, sets `REPORT-LINE` to spaces, writes `REPORT-LINE`, sets `REPORT-LINE` to spaces, builds `REPORT-LINE` from `'MONTH CLOSED                  '` `WS-CLOSED-MONTH`, writes `REPORT-LINE`, sets `WS-RC-LABEL` to `'ACCOUNTS READ'`, sets `WS-RC-VALUE` to `WS-ACCOUNT-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `WS-RC-LABEL` to `'STATEMENTS PRINTED'`, sets `WS-RC-VALUE` to `WS-STATEMENT-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `WS-RC-LABEL` to `'POSTINGS SELECTED'`, sets `WS-RC-VALUE` to `WS-SELECTED-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `WS-RC-LABEL` to `'POSTINGS PRINTED'`, sets `WS-RC-VALUE` to `WS-PRINTED-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `WS-RA-LABEL` to `'CREDITS'`, sets `WS-RA-VALUE` to `WS-CREDIT-TOTAL`, writes `REPORT-LINE` from `WS-REPORT-AMOUNT-LINE`, sets `WS-RA-LABEL` to `'DEBITS'`, sets `WS-RA-VALUE` to `WS-DEBIT-TOTAL`, and writes `REPORT-LINE` from `WS-REPORT-AMOUNT-LINE`.

Source: [STMTPRT.cbl:407](../../../pangaea/programs/STMTPRT.cbl#L407).
