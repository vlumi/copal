# TXPOST: decision tables

The business rules of [TXPOST](../programs/TXPOST.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

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

Source: [TXPOST.cbl:378](../../../pangaea/programs/TXPOST.cbl#L378).
