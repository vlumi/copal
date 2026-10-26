# MONCHG: decision tables

The business rules of [MONCHG](../programs/MONCHG.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

## Contents

- [1000-INITIALIZE](#1000-initialize)
- [1200-READ-ACCOUNT](#1200-read-account)
- [2000-CHARGE-ACCOUNT](#2000-charge-account)
- [2100-SETTLE-INTEREST](#2100-settle-interest)
- [2200-CHARGE-FEE](#2200-charge-fee)
- [8000-WRITE-REPORT](#8000-write-report)

## 1000-INITIALIZE

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `WS-FILE-STATUSES = ALL '0'` | N | Y | Y |
| `READ CONTROL-FILE AT END` | | Y | N |
| → `OPEN INPUT CONTROL-FILE ACCT-IN ACCR-IN OUTPUT TRAN-OUT REPORT-OUT` | 1 | 1 | 1 |
| → `DISPLAY 'MONCHG: CANNOT OPEN FILES ' WS-FILE-STATUSES` | 2 | | |
| → `MOVE 16 TO RETURN-CODE` | 3 | | |
| → `STOP RUN` | 4 | | |
| → `READ CONTROL-FILE` | | 2 | 2 |
| → `MOVE 'NO BANK CONTROL RECORD' TO WS-ABORT-MESSAGE` | | 3 | |
| → `PERFORM 9900-ABORT` | | 4 | |
| → `MOVE BC-BUSINESS-DATE TO WS-BUSINESS-DATE` | | | 3 |
| → `MOVE SPACES TO TRAN-RECORD` | | | 4 |
| → `PERFORM 1200-READ-ACCOUNT` | | | 5 |

- **Rule 1.** When `WS-FILE-STATUSES` is not `ALL '0'`, it opens `CONTROL-FILE`, `ACCT-IN`, and `ACCR-IN` for input and `TRAN-OUT` and `REPORT-OUT` for output; displays `'MONCHG: CANNOT OPEN FILES '` `WS-FILE-STATUSES`; sets `RETURN-CODE` to 16; and stops the run.
- **Rule 2.** When `WS-FILE-STATUSES` is `ALL '0'` and reading `CONTROL-FILE` finds no more records, it opens `CONTROL-FILE`, `ACCT-IN`, and `ACCR-IN` for input and `TRAN-OUT` and `REPORT-OUT` for output; reads `CONTROL-FILE`; sets `WS-ABORT-MESSAGE` to `'NO BANK CONTROL RECORD'`; and performs `9900-ABORT`.
- **Rule 3.** When `WS-FILE-STATUSES` is `ALL '0'` and reading `CONTROL-FILE` finds a record, it opens `CONTROL-FILE`, `ACCT-IN`, and `ACCR-IN` for input and `TRAN-OUT` and `REPORT-OUT` for output; reads `CONTROL-FILE`; sets `WS-BUSINESS-DATE` to `BC-BUSINESS-DATE`; sets `TRAN-RECORD` to spaces; and performs `1200-READ-ACCOUNT`.

Source: [MONCHG.cbl:122](../../../pangaea/programs/MONCHG.cbl#L122).

## 1200-READ-ACCOUNT

| | 1 | 2 | 3 | 4 | 5 | 6 | 7 |
| --- | :-: | :-: | :-: | :-: | :-: | :-: | :-: |
| `READ ACCT-IN AT END` | Y | Y | Y | N | N | N | N |
| `READ ACCR-IN AT END` | Y | N | N | Y | Y | N | N |
| `MASTER-AT-END AND ACCRUALS-AT-END` | Y | Y | N | Y | N | N | N |
| `MASTER-AT-END OR ACCRUALS-AT-END OR AC-ACCT-NO NOT = AM-ACCT-NO` | | | Y | | Y | Y | N |
| → `READ ACCT-IN` | 1 | 1 | 1 | 1 | 1 | 1 | 1 |
| → `SET MASTER-AT-END TO TRUE` | 2 | 2 | 2 | | | | |
| → `READ ACCR-IN` | 3 | 3 | 3 | 2 | 2 | 2 | 2 |
| → `SET ACCRUALS-AT-END TO TRUE` | 4 | | | 3 | 3 | | |
| → `CONTINUE` | 5 | 4 | | 4 | | | |
| → `MOVE 'ACCRUALS OUT OF STEP WITH THE MASTER' TO WS-ABORT-MESSAGE` | | | 4 | | 4 | 3 | |
| → `PERFORM 9900-ABORT` | | | 5 | | 5 | 4 | |

- **Rule 1.** When reading `ACCT-IN` finds no more records; reading `ACCR-IN` finds no more records; and `WS-MASTER-END-SW` is `'Y'` (`MASTER-AT-END`) and `WS-ACCRUAL-END-SW` is `'Y'` (`ACCRUALS-AT-END`), it reads `ACCT-IN`, sets `WS-MASTER-END-SW` to `MASTER-AT-END`, reads `ACCR-IN`, and sets `WS-ACCRUAL-END-SW` to `ACCRUALS-AT-END`.
- **Rule 2.** When reading `ACCT-IN` finds no more records; reading `ACCR-IN` finds a record; and `WS-MASTER-END-SW` is `'Y'` (`MASTER-AT-END`) and `WS-ACCRUAL-END-SW` is `'Y'` (`ACCRUALS-AT-END`), it reads `ACCT-IN`, sets `WS-MASTER-END-SW` to `MASTER-AT-END`, and reads `ACCR-IN`.
- **Rule 3.** When reading `ACCT-IN` finds no more records; reading `ACCR-IN` finds a record; `WS-MASTER-END-SW` is not `'Y'` (not `MASTER-AT-END`) or `WS-ACCRUAL-END-SW` is not `'Y'` (not `ACCRUALS-AT-END`); and `WS-MASTER-END-SW` is `'Y'` (`MASTER-AT-END`), `WS-ACCRUAL-END-SW` is `'Y'` (`ACCRUALS-AT-END`), or `AC-ACCT-NO` is not `AM-ACCT-NO`, it reads `ACCT-IN`, sets `WS-MASTER-END-SW` to `MASTER-AT-END`, reads `ACCR-IN`, sets `WS-ABORT-MESSAGE` to `'ACCRUALS OUT OF STEP WITH THE MASTER'`, and performs `9900-ABORT`.
- **Rule 4.** When reading `ACCT-IN` finds a record; reading `ACCR-IN` finds no more records; and `WS-MASTER-END-SW` is `'Y'` (`MASTER-AT-END`) and `WS-ACCRUAL-END-SW` is `'Y'` (`ACCRUALS-AT-END`), it reads `ACCT-IN`, reads `ACCR-IN`, and sets `WS-ACCRUAL-END-SW` to `ACCRUALS-AT-END`.
- **Rule 5.** When reading `ACCT-IN` finds a record; reading `ACCR-IN` finds no more records; `WS-MASTER-END-SW` is not `'Y'` (not `MASTER-AT-END`) or `WS-ACCRUAL-END-SW` is not `'Y'` (not `ACCRUALS-AT-END`); and `WS-MASTER-END-SW` is `'Y'` (`MASTER-AT-END`), `WS-ACCRUAL-END-SW` is `'Y'` (`ACCRUALS-AT-END`), or `AC-ACCT-NO` is not `AM-ACCT-NO`, it reads `ACCT-IN`, reads `ACCR-IN`, sets `WS-ACCRUAL-END-SW` to `ACCRUALS-AT-END`, sets `WS-ABORT-MESSAGE` to `'ACCRUALS OUT OF STEP WITH THE MASTER'`, and performs `9900-ABORT`.
- **Rule 6.** When reading `ACCT-IN` finds a record; reading `ACCR-IN` finds a record; `WS-MASTER-END-SW` is not `'Y'` (not `MASTER-AT-END`) or `WS-ACCRUAL-END-SW` is not `'Y'` (not `ACCRUALS-AT-END`); and `WS-MASTER-END-SW` is `'Y'` (`MASTER-AT-END`), `WS-ACCRUAL-END-SW` is `'Y'` (`ACCRUALS-AT-END`), or `AC-ACCT-NO` is not `AM-ACCT-NO`, it reads `ACCT-IN`, reads `ACCR-IN`, sets `WS-ABORT-MESSAGE` to `'ACCRUALS OUT OF STEP WITH THE MASTER'`, and performs `9900-ABORT`.
- **Rule 7.** When reading `ACCT-IN` finds a record; reading `ACCR-IN` finds a record; `WS-MASTER-END-SW` is not `'Y'` (not `MASTER-AT-END`) or `WS-ACCRUAL-END-SW` is not `'Y'` (not `ACCRUALS-AT-END`); and `WS-MASTER-END-SW` is not `'Y'` (not `MASTER-AT-END`), `WS-ACCRUAL-END-SW` is not `'Y'` (not `ACCRUALS-AT-END`), and `AC-ACCT-NO` is `AM-ACCT-NO`, it reads `ACCT-IN` and reads `ACCR-IN`.

Source: [MONCHG.cbl:140](../../../pangaea/programs/MONCHG.cbl#L140).

## 2000-CHARGE-ACCOUNT

| | 1 | 2 |
| --- | :-: | :-: |
| `AC-DUE-MONTH = 0` | N | Y |
| → `MOVE AC-DUE-MONTH TO WS-CLOSED-MONTH` | 1 | |
| → `MOVE WS-CM-YYMM TO WS-TI-YYMM` | 2 | |
| → `PERFORM 2100-SETTLE-INTEREST` | 3 | |
| → `PERFORM 2200-CHARGE-FEE` | 4 | |
| → `PERFORM 1200-READ-ACCOUNT` | 5 | 1 |

- **Rule 1.** When `AC-DUE-MONTH` is not 0, it sets `WS-CLOSED-MONTH` to `AC-DUE-MONTH`, sets `WS-TI-YYMM` to `WS-CM-YYMM`, performs `2100-SETTLE-INTEREST`, performs `2200-CHARGE-FEE`, and performs `1200-READ-ACCOUNT`.
- **Rule 2.** When `AC-DUE-MONTH` is 0, it performs `1200-READ-ACCOUNT`.

Source: [MONCHG.cbl:161](../../../pangaea/programs/MONCHG.cbl#L161).

## 2100-SETTLE-INTEREST

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `AC-DUE > 0` | Y | N | N |
| `AC-DUE < 0` | | Y | N |
| → `SET TD-INTEREST TO TRUE` | 1 | | |
| → `MOVE AC-DUE TO TD-AMOUNT` | 2 | | |
| → `MOVE 'CREDIT INTEREST' TO TD-DESCRIPTION` | 3 | | |
| → `PERFORM 2900-WRITE-TRAN` | 4 | 4 | |
| → `ADD 1 TO WS-INTEREST-COUNT` | 5 | | |
| → `ADD TD-AMOUNT TO WS-INTEREST-TOTAL` | 6 | | |
| → `SET TD-OVERDRAFT-INTEREST TO TRUE` | | 1 | |
| → `COMPUTE TD-AMOUNT = 0 - AC-DUE` | | 2 | |
| → `MOVE 'OVERDRAFT INTEREST' TO TD-DESCRIPTION` | | 3 | |
| → `ADD 1 TO WS-OVERDRAFT-COUNT` | | 5 | |
| → `ADD TD-AMOUNT TO WS-OVERDRAFT-TOTAL` | | 6 | |
| → `CONTINUE` | | | 1 |

- **Rule 1.** When `AC-DUE` is more than 0, it sets `TD-TRAN-CODE` to `TD-INTEREST`, sets `TD-AMOUNT` to `AC-DUE`, sets `TD-DESCRIPTION` to `'CREDIT INTEREST'`, performs `2900-WRITE-TRAN`, adds 1 to `WS-INTEREST-COUNT`, and adds `TD-AMOUNT` to `WS-INTEREST-TOTAL`.
- **Rule 2.** When `AC-DUE` is at most 0 and `AC-DUE` is less than 0, it sets `TD-TRAN-CODE` to `TD-OVERDRAFT-INTEREST`, computes `TD-AMOUNT` as `0 - AC-DUE`, sets `TD-DESCRIPTION` to `'OVERDRAFT INTEREST'`, performs `2900-WRITE-TRAN`, adds 1 to `WS-OVERDRAFT-COUNT`, and adds `TD-AMOUNT` to `WS-OVERDRAFT-TOTAL`.
- **Rule 3.** When `AC-DUE` is at most 0 and `AC-DUE` is at least 0, it does nothing more.

Source: [MONCHG.cbl:172](../../../pangaea/programs/MONCHG.cbl#L172).

## 2200-CHARGE-FEE

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `AM-ACCT-TYPE = 'C' AND AM-STATUS = 'O'` | Y | Y | N |
| `AM-BALANCE < WS-FEE-WAIVER-BALANCE` | Y | N | |
| → `SET TD-MAINTENANCE-FEE TO TRUE` | 1 | | |
| → `MOVE WS-MAINTENANCE-FEE TO TD-AMOUNT` | 2 | | |
| → `MOVE 'MAINTENANCE FEE' TO TD-DESCRIPTION` | 3 | | |
| → `PERFORM 2900-WRITE-TRAN` | 4 | | |
| → `ADD 1 TO WS-FEE-COUNT` | 5 | | |
| → `ADD TD-AMOUNT TO WS-FEE-TOTAL` | 6 | | |
| → `ADD 1 TO WS-WAIVED-COUNT` | | 1 | |
| → `CONTINUE` | | | 1 |

- **Rule 1.** When `AM-ACCT-TYPE` is `'C'` and `AM-STATUS` is `'O'`; and `AM-BALANCE` is less than `WS-FEE-WAIVER-BALANCE`, it sets `TD-TRAN-CODE` to `TD-MAINTENANCE-FEE`, sets `TD-AMOUNT` to `WS-MAINTENANCE-FEE`, sets `TD-DESCRIPTION` to `'MAINTENANCE FEE'`, performs `2900-WRITE-TRAN`, adds 1 to `WS-FEE-COUNT`, and adds `TD-AMOUNT` to `WS-FEE-TOTAL`.
- **Rule 2.** When `AM-ACCT-TYPE` is `'C'` and `AM-STATUS` is `'O'`; and `AM-BALANCE` is at least `WS-FEE-WAIVER-BALANCE`, it adds 1 to `WS-WAIVED-COUNT`.
- **Rule 3.** When `AM-ACCT-TYPE` is not `'C'` or `AM-STATUS` is not `'O'`, it does nothing more.

Source: [MONCHG.cbl:194](../../../pangaea/programs/MONCHG.cbl#L194).

## 8000-WRITE-REPORT

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-CLOSED-MONTH = 0` | Y | N |
| → `MOVE SPACES TO REPORT-LINE` | 1, 4 | 1, 4, 6, 9 |
| → `STRING 'MONCHG   MONTH-END CHARGES       BUSINESS DATE ' WS-BUSINESS-DATE DELIMITED BY SIZE INTO REPORT-LINE` | 2 | 2 |
| → `WRITE REPORT-LINE` | 3, 5, 7 | 3, 5, 8, 10 |
| → `MOVE 'NO MONTH CLOSED' TO REPORT-LINE` | 6 | |
| → `STRING 'MONTH CLOSED                  ' WS-CLOSED-MONTH DELIMITED BY SIZE INTO REPORT-LINE` | | 7 |
| → `MOVE 'INTEREST CREDITS' TO WS-RC-LABEL` | | 11 |
| → `MOVE WS-INTEREST-COUNT TO WS-RC-VALUE` | | 12 |
| → `WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE` | | 13, 19, 25, 31 |
| → `MOVE 'INTEREST CREDITED' TO WS-RA-LABEL` | | 14 |
| → `MOVE WS-INTEREST-TOTAL TO WS-RA-VALUE` | | 15 |
| → `WRITE REPORT-LINE FROM WS-REPORT-AMOUNT-LINE` | | 16, 22, 28 |
| → `MOVE 'OVERDRAFT INTEREST CHARGES' TO WS-RC-LABEL` | | 17 |
| → `MOVE WS-OVERDRAFT-COUNT TO WS-RC-VALUE` | | 18 |
| → `MOVE 'OVERDRAFT INTEREST CHARGED' TO WS-RA-LABEL` | | 20 |
| → `MOVE WS-OVERDRAFT-TOTAL TO WS-RA-VALUE` | | 21 |
| → `MOVE 'MAINTENANCE FEES' TO WS-RC-LABEL` | | 23 |
| → `MOVE WS-FEE-COUNT TO WS-RC-VALUE` | | 24 |
| → `MOVE 'MAINTENANCE FEES CHARGED' TO WS-RA-LABEL` | | 26 |
| → `MOVE WS-FEE-TOTAL TO WS-RA-VALUE` | | 27 |
| → `MOVE 'MAINTENANCE FEES WAIVED' TO WS-RC-LABEL` | | 29 |
| → `MOVE WS-WAIVED-COUNT TO WS-RC-VALUE` | | 30 |

- **Rule 1.** When `WS-CLOSED-MONTH` is 0, it sets `REPORT-LINE` to spaces, builds `REPORT-LINE` from `'MONCHG   MONTH-END CHARGES       BUSINESS DATE '` `WS-BUSINESS-DATE`, writes `REPORT-LINE`, sets `REPORT-LINE` to spaces, writes `REPORT-LINE`, sets `REPORT-LINE` to `'NO MONTH CLOSED'`, and writes `REPORT-LINE`.
- **Rule 2.** When `WS-CLOSED-MONTH` is not 0, it sets `REPORT-LINE` to spaces, builds `REPORT-LINE` from `'MONCHG   MONTH-END CHARGES       BUSINESS DATE '` `WS-BUSINESS-DATE`, writes `REPORT-LINE`, sets `REPORT-LINE` to spaces, writes `REPORT-LINE`, sets `REPORT-LINE` to spaces, builds `REPORT-LINE` from `'MONTH CLOSED                  '` `WS-CLOSED-MONTH`, writes `REPORT-LINE`, sets `REPORT-LINE` to spaces, writes `REPORT-LINE`, sets `WS-RC-LABEL` to `'INTEREST CREDITS'`, sets `WS-RC-VALUE` to `WS-INTEREST-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `WS-RA-LABEL` to `'INTEREST CREDITED'`, sets `WS-RA-VALUE` to `WS-INTEREST-TOTAL`, writes `REPORT-LINE` from `WS-REPORT-AMOUNT-LINE`, sets `WS-RC-LABEL` to `'OVERDRAFT INTEREST CHARGES'`, sets `WS-RC-VALUE` to `WS-OVERDRAFT-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `WS-RA-LABEL` to `'OVERDRAFT INTEREST CHARGED'`, sets `WS-RA-VALUE` to `WS-OVERDRAFT-TOTAL`, writes `REPORT-LINE` from `WS-REPORT-AMOUNT-LINE`, sets `WS-RC-LABEL` to `'MAINTENANCE FEES'`, sets `WS-RC-VALUE` to `WS-FEE-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `WS-RA-LABEL` to `'MAINTENANCE FEES CHARGED'`, sets `WS-RA-VALUE` to `WS-FEE-TOTAL`, writes `REPORT-LINE` from `WS-REPORT-AMOUNT-LINE`, sets `WS-RC-LABEL` to `'MAINTENANCE FEES WAIVED'`, sets `WS-RC-VALUE` to `WS-WAIVED-COUNT`, and writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`.

Source: [MONCHG.cbl:220](../../../pangaea/programs/MONCHG.cbl#L220).
