# ATMWDL: decision tables

The business rules of [ATMWDL](../programs/ATMWDL.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

## 0000-MAIN

| | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 |
| --- | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: |
| `CR-MONTH-WITHDRAWALS < WS-FREE-WITHDRAWALS` | Y | Y | Y | Y | Y | Y | N | N | N | N | N | N |
| `CR-BLOCKED` | Y | N | N | N | N | N | Y | N | N | N | N | N |
| `AO-OPEN` | | N | Y | Y | Y | Y | | N | Y | Y | Y | Y |
| `RQ-AMOUNT = 0 OR FUNCTION MOD(RQ-AMOUNT WS-NOTE) NOT = 0` | | | Y | N | N | N | | | Y | N | N | N |
| `WS-DAY-CASH > CR-CASH-LIMIT` | | | | Y | N | N | | | | Y | N | N |
| `WS-TO-HOLD > WS-AVAILABLE` | | | | | Y | N | | | | | Y | N |
| → `PERFORM 1000-START-CARD-TOTALS` | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 |
| → `ADD AO-BALANCE AO-OVERDRAFT-LIMIT GIVING WS-AVAILABLE` | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 |
| → `SUBTRACT AO-HOLDS FROM WS-AVAILABLE` | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 |
| → `COMPUTE WS-DAY-CASH = CR-DAY-CASH + RQ-AMOUNT` | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 |
| → `MOVE 0 TO WS-FEE` | 5 | 5 | 5 | 5 | 5 | 5 | | | | | | |
| → `COMPUTE WS-TO-HOLD = RQ-AMOUNT + WS-FEE` | 6 | 6 | 6 | 6 | 6 | 6 | 6 | 6 | 6 | 6 | 6 | 6 |
| → `SET RS-CARD-BLOCKED TO TRUE` | 7 | | | | | | 7 | | | | | |
| → `MOVE 'CARD BLOCKED' TO RS-MESSAGE` | 8 | | | | | | 8 | | | | | |
| → `GOBACK` | 9 | 9 | 9 | 9 | 9 | 8 | 9 | 9 | 9 | 9 | 9 | 8 |
| → `SET RS-ACCOUNT-NOT-OPEN TO TRUE` | | 7 | | | | | | 7 | | | | |
| → `MOVE 'ACCOUNT NOT OPEN' TO RS-MESSAGE` | | 8 | | | | | | 8 | | | | |
| → `SET RS-INVALID-AMOUNT TO TRUE` | | | 7 | | | | | | 7 | | | |
| → `MOVE 'AMOUNT NOT IN NOTES OF 20' TO RS-MESSAGE` | | | 8 | | | | | | 8 | | | |
| → `SET RS-OVER-LIMIT TO TRUE` | | | | 7 | | | | | | 7 | | |
| → `MOVE 'OVER DAILY CASH LIMIT' TO RS-MESSAGE` | | | | 8 | | | | | | 8 | | |
| → `SET RS-INSUFFICIENT-FUNDS TO TRUE` | | | | | 7 | | | | | | 7 | |
| → `MOVE 'INSUFFICIENT FUNDS' TO RS-MESSAGE` | | | | | 8 | | | | | | 8 | |
| → `PERFORM 2000-APPROVE` | | | | | | 7 | | | | | | 7 |
| → `MOVE WS-ATM-FEE TO WS-FEE` | | | | | | | 5 | 5 | 5 | 5 | 5 | 5 |

Source: [ATMWDL.cbl:47](../../../pangaea/programs/ATMWDL.cbl#L47).

## 1000-START-CARD-TOTALS

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `CR-TOTALS-DATE = BC-BUSINESS-DATE` | N | N | Y |
| `CR-TOTALS-DATE(1:6) = BC-BUSINESS-DATE(1:6)` | N | Y | |
| → `MOVE 0 TO CR-MONTH-WITHDRAWALS` | 1 | | |
| → `MOVE 0 TO CR-DAY-PURCHASES CR-DAY-CASH` | 2 | 1 | |
| → `MOVE BC-BUSINESS-DATE TO CR-TOTALS-DATE` | 3 | 2 | |

Source: [ATMWDL.cbl:82](../../../pangaea/programs/ATMWDL.cbl#L82).
