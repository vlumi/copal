# TXPOST: decision tables

The business rules of [TXPOST](../programs/TXPOST.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

## Contents

- [0000-MAIN](#0000-main)
- [1000-INITIALIZE](#1000-initialize)
- [1200-READ-MASTER](#1200-read-master)
- [1300-READ-TRAN](#1300-read-tran)
- [2000-MATCH](#2000-match)
- [3000-APPLY-TRAN](#3000-apply-tran)
- [3200-POST-DEBIT](#3200-post-debit)
- [3500-WRITE-LEG](#3500-write-leg)
- [3900-REJECT-TRAN](#3900-reject-tran)
- [5000-RECONCILE](#5000-reconcile)

## 0000-MAIN

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-REJECT-COUNT > 0` | Y | N |
| → `PERFORM 1000-INITIALIZE` | 1 | 1 |
| → `PERFORM 2000-MATCH UNTIL WS-MASTER-KEY = HIGH-VALUES AND WS-TRAN-KEY = HIGH-VALUES` | 2 | 2 |
| → `PERFORM 5000-RECONCILE` | 3 | 3 |
| → `PERFORM 8000-WRITE-REPORT` | 4 | 4 |
| → `PERFORM 9000-CLOSE-FILES` | 5 | 5 |
| → `MOVE 4 TO RETURN-CODE` | 6 | |
| → `STOP RUN` | 7 | 7 |
| → `MOVE 0 TO RETURN-CODE` | | 6 |

- **Rule 1.** When `WS-REJECT-COUNT` is more than 0, it performs `1000-INITIALIZE`; performs `2000-MATCH` until `WS-MASTER-KEY` is high values and `WS-TRAN-KEY` is high values; performs `5000-RECONCILE`; performs `8000-WRITE-REPORT`; performs `9000-CLOSE-FILES`; sets `RETURN-CODE` to 4; and stops the run.
- **Rule 2.** When `WS-REJECT-COUNT` is at most 0, it performs `1000-INITIALIZE`; performs `2000-MATCH` until `WS-MASTER-KEY` is high values and `WS-TRAN-KEY` is high values; performs `5000-RECONCILE`; performs `8000-WRITE-REPORT`; performs `9000-CLOSE-FILES`; sets `RETURN-CODE` to 0; and stops the run.

Source: [TXPOST.cbl:164](../../../pangaea/programs/TXPOST.cbl#L164).

## 1000-INITIALIZE

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `WS-CONTROL-STATUS NOT = '00' OR WS-ACCT-IN-STATUS NOT = '00' OR WS-TRAN-IN-STATUS NOT = '00' OR WS-ACCT-OUT-STATUS NOT = '00' OR NOT AUDIT-OPENED OR NOT REJECTS-OPENED OR WS-LEGS-STATUS NOT = '00' OR WS-REPORT-STATUS NOT = '00'` | Y | N | N |
| `READ CONTROL-FILE AT END` | | Y | N |
| → `OPEN INPUT CONTROL-FILE ACCT-IN TRAN-IN OUTPUT ACCT-OUT LEGS-OUT REPORT-OUT EXTEND AUDIT-OUT REJECT-OUT` | 1 | 1 | 1 |
| → `DISPLAY 'TXPOST: CANNOT OPEN FILES ' WS-FILE-STATUSES` | 2 | | |
| → `MOVE 16 TO RETURN-CODE` | 3 | | |
| → `STOP RUN` | 4 | | |
| → `READ CONTROL-FILE` | | 2 | 2 |
| → `MOVE 'NO BANK CONTROL RECORD' TO WS-ABORT-MESSAGE` | | 3 | |
| → `PERFORM 9900-ABORT` | | 4 | |
| → `MOVE BC-BUSINESS-DATE TO WS-BUSINESS-DATE` | | | 3 |
| → `MOVE SPACES TO AUDIT-RECORD REJECT-RECORD` | | | 4 |
| → `STRING WS-BUSINESS-DATE(1:4) '-' WS-BUSINESS-DATE(5:2) '-' WS-BUSINESS-DATE(7:2) DELIMITED BY SIZE INTO RJ-REJECT-DATE` | | | 5 |
| → `PERFORM 1200-READ-MASTER` | | | 6 |
| → `PERFORM 1300-READ-TRAN` | | | 7 |

- **Rule 1.** When `WS-CONTROL-STATUS` is not `'00'`, `WS-ACCT-IN-STATUS` is not `'00'`, `WS-TRAN-IN-STATUS` is not `'00'`, `WS-ACCT-OUT-STATUS` is not `'00'`, `WS-AUDIT-STATUS` is neither `'00'` nor `'05'` (not `AUDIT-OPENED`), `WS-REJECT-STATUS` is neither `'00'` nor `'05'` (not `REJECTS-OPENED`), `WS-LEGS-STATUS` is not `'00'`, or `WS-REPORT-STATUS` is not `'00'`, it opens `CONTROL-FILE`, `ACCT-IN`, and `TRAN-IN` for input, `ACCT-OUT`, `LEGS-OUT`, and `REPORT-OUT` for output, and `AUDIT-OUT` and `REJECT-OUT` to extend; displays `'TXPOST: CANNOT OPEN FILES '` `WS-FILE-STATUSES`; sets `RETURN-CODE` to 16; and stops the run.
- **Rule 2.** When `WS-CONTROL-STATUS` is `'00'`, `WS-ACCT-IN-STATUS` is `'00'`, `WS-TRAN-IN-STATUS` is `'00'`, `WS-ACCT-OUT-STATUS` is `'00'`, `WS-AUDIT-STATUS` is `'00'` or `'05'` (`AUDIT-OPENED`), `WS-REJECT-STATUS` is `'00'` or `'05'` (`REJECTS-OPENED`), `WS-LEGS-STATUS` is `'00'`, and `WS-REPORT-STATUS` is `'00'`; and reading `CONTROL-FILE` finds no more records, it opens `CONTROL-FILE`, `ACCT-IN`, and `TRAN-IN` for input, `ACCT-OUT`, `LEGS-OUT`, and `REPORT-OUT` for output, and `AUDIT-OUT` and `REJECT-OUT` to extend; reads `CONTROL-FILE`; sets `WS-ABORT-MESSAGE` to `'NO BANK CONTROL RECORD'`; and performs `9900-ABORT`.
- **Rule 3.** When `WS-CONTROL-STATUS` is `'00'`, `WS-ACCT-IN-STATUS` is `'00'`, `WS-TRAN-IN-STATUS` is `'00'`, `WS-ACCT-OUT-STATUS` is `'00'`, `WS-AUDIT-STATUS` is `'00'` or `'05'` (`AUDIT-OPENED`), `WS-REJECT-STATUS` is `'00'` or `'05'` (`REJECTS-OPENED`), `WS-LEGS-STATUS` is `'00'`, and `WS-REPORT-STATUS` is `'00'`; and reading `CONTROL-FILE` finds a record, it opens `CONTROL-FILE`, `ACCT-IN`, and `TRAN-IN` for input, `ACCT-OUT`, `LEGS-OUT`, and `REPORT-OUT` for output, and `AUDIT-OUT` and `REJECT-OUT` to extend; reads `CONTROL-FILE`; sets `WS-BUSINESS-DATE` to `BC-BUSINESS-DATE`; sets `AUDIT-RECORD` and `REJECT-RECORD` to spaces; builds `RJ-REJECT-DATE` from `WS-BUSINESS-DATE(1:4)` `'-'` `WS-BUSINESS-DATE(5:2)` `'-'` `WS-BUSINESS-DATE(7:2)`; performs `1200-READ-MASTER`; and performs `1300-READ-TRAN`.

Source: [TXPOST.cbl:181](../../../pangaea/programs/TXPOST.cbl#L181).

## 1200-READ-MASTER

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `READ ACCT-IN AT END` | Y | N | N |
| `WS-MASTER-KEY > WS-PREV-MASTER-KEY` | | N | Y |
| → `READ ACCT-IN` | 1 | 1 | 1 |
| → `MOVE HIGH-VALUES TO WS-MASTER-KEY` | 2 | | |
| → `MOVE AM-ACCT-NO TO WS-MASTER-KEY` | | 2 | 2 |
| → `MOVE 'MASTER NOT IN ACCOUNT ORDER' TO WS-ABORT-MESSAGE` | | 3 | |
| → `PERFORM 9900-ABORT` | | 4 | |
| → `MOVE WS-MASTER-KEY TO WS-PREV-MASTER-KEY` | | | 3 |
| → `ADD 1 TO WS-MASTERS-READ` | | | 4 |
| → `ADD AM-BALANCE TO WS-OPENING-TOTAL` | | | 5 |

- **Rule 1.** When reading `ACCT-IN` finds no more records, it reads `ACCT-IN` and sets `WS-MASTER-KEY` to high values.
- **Rule 2.** When reading `ACCT-IN` finds a record and `WS-MASTER-KEY` is at most `WS-PREV-MASTER-KEY`, it reads `ACCT-IN`, sets `WS-MASTER-KEY` to `AM-ACCT-NO`, sets `WS-ABORT-MESSAGE` to `'MASTER NOT IN ACCOUNT ORDER'`, and performs `9900-ABORT`.
- **Rule 3.** When reading `ACCT-IN` finds a record and `WS-MASTER-KEY` is more than `WS-PREV-MASTER-KEY`, it reads `ACCT-IN`, sets `WS-MASTER-KEY` to `AM-ACCT-NO`, sets `WS-PREV-MASTER-KEY` to `WS-MASTER-KEY`, adds 1 to `WS-MASTERS-READ`, and adds `AM-BALANCE` to `WS-OPENING-TOTAL`.

Source: [TXPOST.cbl:211](../../../pangaea/programs/TXPOST.cbl#L211).

## 1300-READ-TRAN

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `READ TRAN-IN AT END` | Y | N | N |
| `WS-TRAN-ORDER-KEY > WS-PREV-TRAN-ORDER-KEY` | | N | Y |
| → `READ TRAN-IN` | 1 | 1 | 1 |
| → `MOVE HIGH-VALUES TO WS-TRAN-KEY` | 2 | | |
| → `MOVE TD-ACCT-NO TO WS-TRAN-KEY WS-TOK-ACCT-NO` | | 2 | 2 |
| → `MOVE TD-TRAN-ID TO WS-TOK-TRAN-ID` | | 3 | 3 |
| → `MOVE 'TRANSACTIONS NOT IN ACCOUNT ORDER' TO WS-ABORT-MESSAGE` | | 4 | |
| → `PERFORM 9900-ABORT` | | 5 | |
| → `MOVE WS-TRAN-ORDER-KEY TO WS-PREV-TRAN-ORDER-KEY` | | | 4 |

- **Rule 1.** When reading `TRAN-IN` finds no more records, it reads `TRAN-IN` and sets `WS-TRAN-KEY` to high values.
- **Rule 2.** When reading `TRAN-IN` finds a record and `WS-TRAN-ORDER-KEY` is at most `WS-PREV-TRAN-ORDER-KEY`, it reads `TRAN-IN`; sets `WS-TRAN-KEY` and `WS-TOK-ACCT-NO` to `TD-ACCT-NO`; sets `WS-TOK-TRAN-ID` to `TD-TRAN-ID`; sets `WS-ABORT-MESSAGE` to `'TRANSACTIONS NOT IN ACCOUNT ORDER'`; and performs `9900-ABORT`.
- **Rule 3.** When reading `TRAN-IN` finds a record and `WS-TRAN-ORDER-KEY` is more than `WS-PREV-TRAN-ORDER-KEY`, it reads `TRAN-IN`; sets `WS-TRAN-KEY` and `WS-TOK-ACCT-NO` to `TD-ACCT-NO`; sets `WS-TOK-TRAN-ID` to `TD-TRAN-ID`; and sets `WS-PREV-TRAN-ORDER-KEY` to `WS-TRAN-ORDER-KEY`.

Source: [TXPOST.cbl:228](../../../pangaea/programs/TXPOST.cbl#L228).

## 2000-MATCH

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `WS-TRAN-KEY < WS-MASTER-KEY` | Y | N | N |
| `WS-TRAN-KEY = WS-MASTER-KEY` | | Y | N |
| → `SET RJ-NO-ACCOUNT TO TRUE` | 1 | | |
| → `PERFORM 3900-REJECT-TRAN` | 2 | | |
| → `PERFORM 1300-READ-TRAN` | 3 | 2 | |
| → `PERFORM 3000-APPLY-TRAN` | | 1 | |
| → `PERFORM 4000-WRITE-MASTER` | | | 1 |
| → `PERFORM 1200-READ-MASTER` | | | 2 |

- **Rule 1.** When `WS-TRAN-KEY` is less than `WS-MASTER-KEY`, it sets `RJ-REASON` to `RJ-NO-ACCOUNT`, performs `3900-REJECT-TRAN`, and performs `1300-READ-TRAN`.
- **Rule 2.** When `WS-TRAN-KEY` is at least `WS-MASTER-KEY` and `WS-TRAN-KEY` is `WS-MASTER-KEY`, it performs `3000-APPLY-TRAN` and performs `1300-READ-TRAN`.
- **Rule 3.** When `WS-TRAN-KEY` is at least `WS-MASTER-KEY` and `WS-TRAN-KEY` is not `WS-MASTER-KEY`, it performs `4000-WRITE-MASTER` and performs `1200-READ-MASTER`.

Source: [TXPOST.cbl:246](../../../pangaea/programs/TXPOST.cbl#L246).

## 3000-APPLY-TRAN

| | 1 | 2 | 3 | 4 | 5 | 6 | 7 |
| --- | :-: | :-: | :-: | :-: | :-: | :-: | :-: |
| `AM-CLOSED` | Y | N | N | N | N | N | N |
| `TD-DEPOSIT OR TD-TRANSFER-CREDIT OR TD-TRANSFER-RETURN OR TD-INTEREST` | | Y | N | N | N | N | N |
| `TD-FEE OR TD-OVERDRAFT-INTEREST OR TD-MAINTENANCE-FEE OR TD-ATM-FEE` | | | Y | N | N | N | N |
| `TD-CARD-PURCHASE OR TD-ATM-WITHDRAWAL` | | | | Y | N | N | N |
| `AM-FROZEN` | | | | | Y | N | N |
| `TD-WITHDRAWAL OR TD-TRANSFER` | | | | | | Y | N |
| → `MOVE AM-BALANCE TO WS-BALANCE-BEFORE` | 1 | 1 | 1 | 1 | 1 | 1 | 1 |
| → `SET RJ-CLOSED TO TRUE` | 2 | | | | | | |
| → `PERFORM 3900-REJECT-TRAN` | 3 | | | | 3 | | |
| → `PERFORM 3100-POST-CREDIT` | | 2 | | | | | |
| → `PERFORM 3300-POST-FEE` | | | 2 | | | | |
| → `PERFORM 3400-POST-AUTHORIZED` | | | | 2 | | | |
| → `SET RJ-FROZEN TO TRUE` | | | | | 2 | | |
| → `PERFORM 3200-POST-DEBIT` | | | | | | 2 | |

- **Rule 1.** When `AM-STATUS` is `'X'` (`AM-CLOSED`), it sets `WS-BALANCE-BEFORE` to `AM-BALANCE`, sets `RJ-REASON` to `RJ-CLOSED`, and performs `3900-REJECT-TRAN`.
- **Rule 2.** When `AM-STATUS` is not `'X'` (not `AM-CLOSED`); and `TD-TRAN-CODE` is `'DP'` (`TD-DEPOSIT`), `TD-TRAN-CODE` is `'TC'` (`TD-TRANSFER-CREDIT`), `TD-TRAN-CODE` is `'RT'` (`TD-TRANSFER-RETURN`), or `TD-TRAN-CODE` is `'IN'` (`TD-INTEREST`), it sets `WS-BALANCE-BEFORE` to `AM-BALANCE` and performs `3100-POST-CREDIT`.
- **Rule 3.** When `AM-STATUS` is not `'X'` (not `AM-CLOSED`); `TD-TRAN-CODE` is not `'DP'` (not `TD-DEPOSIT`), `TD-TRAN-CODE` is not `'TC'` (not `TD-TRANSFER-CREDIT`), `TD-TRAN-CODE` is not `'RT'` (not `TD-TRANSFER-RETURN`), and `TD-TRAN-CODE` is not `'IN'` (not `TD-INTEREST`); and `TD-TRAN-CODE` is `'FE'` (`TD-FEE`), `TD-TRAN-CODE` is `'OD'` (`TD-OVERDRAFT-INTEREST`), `TD-TRAN-CODE` is `'MF'` (`TD-MAINTENANCE-FEE`), or `TD-TRAN-CODE` is `'AF'` (`TD-ATM-FEE`), it sets `WS-BALANCE-BEFORE` to `AM-BALANCE` and performs `3300-POST-FEE`.
- **Rule 4.** When `AM-STATUS` is not `'X'` (not `AM-CLOSED`); `TD-TRAN-CODE` is not `'DP'` (not `TD-DEPOSIT`), `TD-TRAN-CODE` is not `'TC'` (not `TD-TRANSFER-CREDIT`), `TD-TRAN-CODE` is not `'RT'` (not `TD-TRANSFER-RETURN`), and `TD-TRAN-CODE` is not `'IN'` (not `TD-INTEREST`); `TD-TRAN-CODE` is not `'FE'` (not `TD-FEE`), `TD-TRAN-CODE` is not `'OD'` (not `TD-OVERDRAFT-INTEREST`), `TD-TRAN-CODE` is not `'MF'` (not `TD-MAINTENANCE-FEE`), and `TD-TRAN-CODE` is not `'AF'` (not `TD-ATM-FEE`); and `TD-TRAN-CODE` is `'CP'` (`TD-CARD-PURCHASE`) or `TD-TRAN-CODE` is `'AW'` (`TD-ATM-WITHDRAWAL`), it sets `WS-BALANCE-BEFORE` to `AM-BALANCE` and performs `3400-POST-AUTHORIZED`.
- **Rule 5.** When `AM-STATUS` is not `'X'` (not `AM-CLOSED`); `TD-TRAN-CODE` is not `'DP'` (not `TD-DEPOSIT`), `TD-TRAN-CODE` is not `'TC'` (not `TD-TRANSFER-CREDIT`), `TD-TRAN-CODE` is not `'RT'` (not `TD-TRANSFER-RETURN`), and `TD-TRAN-CODE` is not `'IN'` (not `TD-INTEREST`); `TD-TRAN-CODE` is not `'FE'` (not `TD-FEE`), `TD-TRAN-CODE` is not `'OD'` (not `TD-OVERDRAFT-INTEREST`), `TD-TRAN-CODE` is not `'MF'` (not `TD-MAINTENANCE-FEE`), and `TD-TRAN-CODE` is not `'AF'` (not `TD-ATM-FEE`); `TD-TRAN-CODE` is not `'CP'` (not `TD-CARD-PURCHASE`) and `TD-TRAN-CODE` is not `'AW'` (not `TD-ATM-WITHDRAWAL`); and `AM-STATUS` is `'F'` (`AM-FROZEN`), it sets `WS-BALANCE-BEFORE` to `AM-BALANCE`, sets `RJ-REASON` to `RJ-FROZEN`, and performs `3900-REJECT-TRAN`.
- **Rule 6.** When `AM-STATUS` is not `'X'` (not `AM-CLOSED`); `TD-TRAN-CODE` is not `'DP'` (not `TD-DEPOSIT`), `TD-TRAN-CODE` is not `'TC'` (not `TD-TRANSFER-CREDIT`), `TD-TRAN-CODE` is not `'RT'` (not `TD-TRANSFER-RETURN`), and `TD-TRAN-CODE` is not `'IN'` (not `TD-INTEREST`); `TD-TRAN-CODE` is not `'FE'` (not `TD-FEE`), `TD-TRAN-CODE` is not `'OD'` (not `TD-OVERDRAFT-INTEREST`), `TD-TRAN-CODE` is not `'MF'` (not `TD-MAINTENANCE-FEE`), and `TD-TRAN-CODE` is not `'AF'` (not `TD-ATM-FEE`); `TD-TRAN-CODE` is not `'CP'` (not `TD-CARD-PURCHASE`) and `TD-TRAN-CODE` is not `'AW'` (not `TD-ATM-WITHDRAWAL`); `AM-STATUS` is not `'F'` (not `AM-FROZEN`); and `TD-TRAN-CODE` is `'WD'` (`TD-WITHDRAWAL`) or `TD-TRAN-CODE` is `'TR'` (`TD-TRANSFER`), it sets `WS-BALANCE-BEFORE` to `AM-BALANCE` and performs `3200-POST-DEBIT`.
- **Rule 7.** When `AM-STATUS` is not `'X'` (not `AM-CLOSED`); `TD-TRAN-CODE` is not `'DP'` (not `TD-DEPOSIT`), `TD-TRAN-CODE` is not `'TC'` (not `TD-TRANSFER-CREDIT`), `TD-TRAN-CODE` is not `'RT'` (not `TD-TRANSFER-RETURN`), and `TD-TRAN-CODE` is not `'IN'` (not `TD-INTEREST`); `TD-TRAN-CODE` is not `'FE'` (not `TD-FEE`), `TD-TRAN-CODE` is not `'OD'` (not `TD-OVERDRAFT-INTEREST`), `TD-TRAN-CODE` is not `'MF'` (not `TD-MAINTENANCE-FEE`), and `TD-TRAN-CODE` is not `'AF'` (not `TD-ATM-FEE`); `TD-TRAN-CODE` is not `'CP'` (not `TD-CARD-PURCHASE`) and `TD-TRAN-CODE` is not `'AW'` (not `TD-ATM-WITHDRAWAL`); `AM-STATUS` is not `'F'` (not `AM-FROZEN`); and `TD-TRAN-CODE` is not `'WD'` (not `TD-WITHDRAWAL`) and `TD-TRAN-CODE` is not `'TR'` (not `TD-TRANSFER`), it sets `WS-BALANCE-BEFORE` to `AM-BALANCE`.

Source: [TXPOST.cbl:261](../../../pangaea/programs/TXPOST.cbl#L261).

## 3200-POST-DEBIT

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `WS-NEW-BALANCE < WS-BALANCE-FLOOR` | Y | N | N |
| `TD-TRANSFER` | | Y | N |
| → `COMPUTE WS-NEW-BALANCE = AM-BALANCE - TD-AMOUNT` | 1 | 1 | 1 |
| → `COMPUTE WS-BALANCE-FLOOR = 0 - AM-OVERDRAFT-LIMIT` | 2 | 2 | 2 |
| → `SET RJ-NO-FUNDS TO TRUE` | 3 | | |
| → `PERFORM 3900-REJECT-TRAN` | 4 | | |
| → `MOVE WS-NEW-BALANCE TO AM-BALANCE` | | 3 | 3 |
| → `ADD TD-AMOUNT TO WS-DEBIT-TOTAL` | | 4 | 4 |
| → `COMPUTE WS-SIGNED-AMOUNT = 0 - TD-AMOUNT` | | 5 | 5 |
| → `PERFORM 3800-WRITE-AUDIT` | | 6 | 6 |
| → `PERFORM 3500-WRITE-LEG` | | 7 | |

- **Rule 1.** When `WS-NEW-BALANCE` is less than `WS-BALANCE-FLOOR`, it computes `WS-NEW-BALANCE` as `AM-BALANCE - TD-AMOUNT`, computes `WS-BALANCE-FLOOR` as `0 - AM-OVERDRAFT-LIMIT`, sets `RJ-REASON` to `RJ-NO-FUNDS`, and performs `3900-REJECT-TRAN`.
- **Rule 2.** When `WS-NEW-BALANCE` is at least `WS-BALANCE-FLOOR` and `TD-TRAN-CODE` is `'TR'` (`TD-TRANSFER`), it computes `WS-NEW-BALANCE` as `AM-BALANCE - TD-AMOUNT`, computes `WS-BALANCE-FLOOR` as `0 - AM-OVERDRAFT-LIMIT`, sets `AM-BALANCE` to `WS-NEW-BALANCE`, adds `TD-AMOUNT` to `WS-DEBIT-TOTAL`, computes `WS-SIGNED-AMOUNT` as `0 - TD-AMOUNT`, performs `3800-WRITE-AUDIT`, and performs `3500-WRITE-LEG`.
- **Rule 3.** When `WS-NEW-BALANCE` is at least `WS-BALANCE-FLOOR` and `TD-TRAN-CODE` is not `'TR'` (not `TD-TRANSFER`), it computes `WS-NEW-BALANCE` as `AM-BALANCE - TD-AMOUNT`, computes `WS-BALANCE-FLOOR` as `0 - AM-OVERDRAFT-LIMIT`, sets `AM-BALANCE` to `WS-NEW-BALANCE`, adds `TD-AMOUNT` to `WS-DEBIT-TOTAL`, computes `WS-SIGNED-AMOUNT` as `0 - TD-AMOUNT`, and performs `3800-WRITE-AUDIT`.

Source: [TXPOST.cbl:298](../../../pangaea/programs/TXPOST.cbl#L298).

## 3500-WRITE-LEG

| | 1 | 2 |
| --- | :-: | :-: |
| `TD-TRANSFER` | Y | N |
| → `MOVE TRAN-RECORD TO WS-SAVED-TRAN` | 1 | 1 |
| → `SET TD-TRANSFER-CREDIT TO TRUE` | 2 | |
| → `MOVE TD-ACCT-NO TO WS-LEG-ACCT-NO` | 3 | 3 |
| → `MOVE TD-COUNTER-ACCT-NO TO TD-ACCT-NO` | 4 | 4 |
| → `MOVE WS-LEG-ACCT-NO TO TD-COUNTER-ACCT-NO` | 5 | 5 |
| → `WRITE LEGS-OUT-RECORD FROM TRAN-RECORD` | 6 | 6 |
| → `ADD 1 TO WS-LEG-COUNT` | 7 | 7 |
| → `MOVE WS-SAVED-TRAN TO TRAN-RECORD` | 8 | 8 |
| → `SET TD-TRANSFER-RETURN TO TRUE` | | 2 |

- **Rule 1.** When `TD-TRAN-CODE` is `'TR'` (`TD-TRANSFER`), it sets `WS-SAVED-TRAN` to `TRAN-RECORD`, sets `TD-TRAN-CODE` to `TD-TRANSFER-CREDIT`, sets `WS-LEG-ACCT-NO` to `TD-ACCT-NO`, sets `TD-ACCT-NO` to `TD-COUNTER-ACCT-NO`, sets `TD-COUNTER-ACCT-NO` to `WS-LEG-ACCT-NO`, writes `LEGS-OUT-RECORD` from `TRAN-RECORD`, adds 1 to `WS-LEG-COUNT`, and sets `TRAN-RECORD` to `WS-SAVED-TRAN`.
- **Rule 2.** When `TD-TRAN-CODE` is not `'TR'` (not `TD-TRANSFER`), it sets `WS-SAVED-TRAN` to `TRAN-RECORD`, sets `TD-TRAN-CODE` to `TD-TRANSFER-RETURN`, sets `WS-LEG-ACCT-NO` to `TD-ACCT-NO`, sets `TD-ACCT-NO` to `TD-COUNTER-ACCT-NO`, sets `TD-COUNTER-ACCT-NO` to `WS-LEG-ACCT-NO`, writes `LEGS-OUT-RECORD` from `TRAN-RECORD`, adds 1 to `WS-LEG-COUNT`, and sets `TRAN-RECORD` to `WS-SAVED-TRAN`.

Source: [TXPOST.cbl:333](../../../pangaea/programs/TXPOST.cbl#L333).

## 3900-REJECT-TRAN

| | 1 | 2 |
| --- | :-: | :-: |
| `TD-TRANSFER-CREDIT` | Y | N |
| → `MOVE 'TXPOST' TO RJ-PROGRAM` | 1 | 1 |
| → `MOVE TRAN-RECORD TO RJ-TRAN` | 2 | 2 |
| → `WRITE REJECT-OUT-RECORD FROM REJECT-RECORD` | 3 | 3 |
| → `ADD 1 TO WS-REJECT-COUNT` | 4 | 4 |
| → `PERFORM 3500-WRITE-LEG` | 5 | |

- **Rule 1.** When `TD-TRAN-CODE` is `'TC'` (`TD-TRANSFER-CREDIT`), it sets `RJ-PROGRAM` to `'TXPOST'`, sets `RJ-TRAN` to `TRAN-RECORD`, writes `REJECT-OUT-RECORD` from `REJECT-RECORD`, adds 1 to `WS-REJECT-COUNT`, and performs `3500-WRITE-LEG`.
- **Rule 2.** When `TD-TRAN-CODE` is not `'TC'` (not `TD-TRANSFER-CREDIT`), it sets `RJ-PROGRAM` to `'TXPOST'`, sets `RJ-TRAN` to `TRAN-RECORD`, writes `REJECT-OUT-RECORD` from `REJECT-RECORD`, and adds 1 to `WS-REJECT-COUNT`.

Source: [TXPOST.cbl:362](../../../pangaea/programs/TXPOST.cbl#L362).

## 5000-RECONCILE

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-EXPECTED-TOTAL NOT = WS-CLOSING-TOTAL OR WS-MASTERS-WRITTEN NOT = WS-MASTERS-READ` | Y | N |
| → `MOVE WS-OPENING-TOTAL TO WS-EXPECTED-TOTAL` | 1 | 1 |
| → `ADD WS-CREDIT-TOTAL TO WS-EXPECTED-TOTAL` | 2 | 2 |
| → `SUBTRACT WS-DEBIT-TOTAL FROM WS-EXPECTED-TOTAL` | 3 | 3 |
| → `SUBTRACT WS-FEE-TOTAL FROM WS-EXPECTED-TOTAL` | 4 | 4 |
| → `MOVE 'NEW MASTER OUT OF BALANCE' TO WS-ABORT-MESSAGE` | 5 | |
| → `PERFORM 9900-ABORT` | 6 | |

- **Rule 1.** When `WS-EXPECTED-TOTAL` is not `WS-CLOSING-TOTAL` or `WS-MASTERS-WRITTEN` is not `WS-MASTERS-READ`, it sets `WS-EXPECTED-TOTAL` to `WS-OPENING-TOTAL`, adds `WS-CREDIT-TOTAL` to `WS-EXPECTED-TOTAL`, subtracts `WS-DEBIT-TOTAL` from `WS-EXPECTED-TOTAL`, subtracts `WS-FEE-TOTAL` from `WS-EXPECTED-TOTAL`, sets `WS-ABORT-MESSAGE` to `'NEW MASTER OUT OF BALANCE'`, and performs `9900-ABORT`.
- **Rule 2.** When `WS-EXPECTED-TOTAL` is `WS-CLOSING-TOTAL` and `WS-MASTERS-WRITTEN` is `WS-MASTERS-READ`, it sets `WS-EXPECTED-TOTAL` to `WS-OPENING-TOTAL`, adds `WS-CREDIT-TOTAL` to `WS-EXPECTED-TOTAL`, subtracts `WS-DEBIT-TOTAL` from `WS-EXPECTED-TOTAL`, and subtracts `WS-FEE-TOTAL` from `WS-EXPECTED-TOTAL`.

Source: [TXPOST.cbl:378](../../../pangaea/programs/TXPOST.cbl#L378).
