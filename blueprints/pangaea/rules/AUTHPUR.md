# AUTHPUR: decision tables

The business rules of [AUTHPUR](../programs/AUTHPUR.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

## 0000-MAIN

| | 1 | 2 | 3 | 4 | 5 | 6 |
| --- | :-: | :-: | :-: | :-: | :-: | :-: |
| `CR-BLOCKED` | Y | N | N | N | N | N |
| `AO-OPEN` | | N | Y | Y | Y | Y |
| `RQ-AMOUNT = 0` | | | Y | N | N | N |
| `WS-DAY-PURCHASES > CR-PURCHASE-LIMIT` | | | | Y | N | N |
| `RQ-AMOUNT > WS-AVAILABLE` | | | | | Y | N |
| → `PERFORM 1000-START-CARD-TOTALS` | 1 | 1 | 1 | 1 | 1 | 1 |
| → `ADD AO-BALANCE AO-OVERDRAFT-LIMIT GIVING WS-AVAILABLE` | 2 | 2 | 2 | 2 | 2 | 2 |
| → `SUBTRACT AO-HOLDS FROM WS-AVAILABLE` | 3 | 3 | 3 | 3 | 3 | 3 |
| → `COMPUTE WS-DAY-PURCHASES = CR-DAY-PURCHASES + RQ-AMOUNT` | 4 | 4 | 4 | 4 | 4 | 4 |
| → `SET RS-CARD-BLOCKED TO TRUE` | 5 | | | | | |
| → `MOVE 'CARD BLOCKED' TO RS-MESSAGE` | 6 | | | | | |
| → `GOBACK` | 7 | 7 | 7 | 7 | 7 | 6 |
| → `SET RS-ACCOUNT-NOT-OPEN TO TRUE` | | 5 | | | | |
| → `MOVE 'ACCOUNT NOT OPEN' TO RS-MESSAGE` | | 6 | | | | |
| → `SET RS-INVALID-AMOUNT TO TRUE` | | | 5 | | | |
| → `MOVE 'NO AMOUNT' TO RS-MESSAGE` | | | 6 | | | |
| → `SET RS-OVER-LIMIT TO TRUE` | | | | 5 | | |
| → `MOVE 'OVER DAILY PURCHASE LIMIT' TO RS-MESSAGE` | | | | 6 | | |
| → `SET RS-INSUFFICIENT-FUNDS TO TRUE` | | | | | 5 | |
| → `MOVE 'INSUFFICIENT FUNDS' TO RS-MESSAGE` | | | | | 6 | |
| → `PERFORM 2000-APPROVE` | | | | | | 5 |

- **Rule 1.** When `CR-STATUS` is `'B'` (`CR-BLOCKED`), it performs `1000-START-CARD-TOTALS`, runs `ADD AO-BALANCE AO-OVERDRAFT-LIMIT GIVING WS-AVAILABLE`, subtracts `AO-HOLDS` from `WS-AVAILABLE`, computes `WS-DAY-PURCHASES` as `CR-DAY-PURCHASES + RQ-AMOUNT`, sets `RS-RESULT` to `RS-CARD-BLOCKED`, sets `RS-MESSAGE` to `'CARD BLOCKED'`, and returns to its caller.
- **Rule 2.** When `CR-STATUS` is not `'B'` (not `CR-BLOCKED`) and `AO-STATUS` is not `'O'` (not `AO-OPEN`), it performs `1000-START-CARD-TOTALS`, runs `ADD AO-BALANCE AO-OVERDRAFT-LIMIT GIVING WS-AVAILABLE`, subtracts `AO-HOLDS` from `WS-AVAILABLE`, computes `WS-DAY-PURCHASES` as `CR-DAY-PURCHASES + RQ-AMOUNT`, sets `RS-RESULT` to `RS-ACCOUNT-NOT-OPEN`, sets `RS-MESSAGE` to `'ACCOUNT NOT OPEN'`, and returns to its caller.
- **Rule 3.** When `CR-STATUS` is not `'B'` (not `CR-BLOCKED`), `AO-STATUS` is `'O'` (`AO-OPEN`), and `RQ-AMOUNT` is 0, it performs `1000-START-CARD-TOTALS`, runs `ADD AO-BALANCE AO-OVERDRAFT-LIMIT GIVING WS-AVAILABLE`, subtracts `AO-HOLDS` from `WS-AVAILABLE`, computes `WS-DAY-PURCHASES` as `CR-DAY-PURCHASES + RQ-AMOUNT`, sets `RS-RESULT` to `RS-INVALID-AMOUNT`, sets `RS-MESSAGE` to `'NO AMOUNT'`, and returns to its caller.
- **Rule 4.** When `CR-STATUS` is not `'B'` (not `CR-BLOCKED`), `AO-STATUS` is `'O'` (`AO-OPEN`), `RQ-AMOUNT` is not 0, and `WS-DAY-PURCHASES` is more than `CR-PURCHASE-LIMIT`, it performs `1000-START-CARD-TOTALS`, runs `ADD AO-BALANCE AO-OVERDRAFT-LIMIT GIVING WS-AVAILABLE`, subtracts `AO-HOLDS` from `WS-AVAILABLE`, computes `WS-DAY-PURCHASES` as `CR-DAY-PURCHASES + RQ-AMOUNT`, sets `RS-RESULT` to `RS-OVER-LIMIT`, sets `RS-MESSAGE` to `'OVER DAILY PURCHASE LIMIT'`, and returns to its caller.
- **Rule 5.** When `CR-STATUS` is not `'B'` (not `CR-BLOCKED`), `AO-STATUS` is `'O'` (`AO-OPEN`), `RQ-AMOUNT` is not 0, `WS-DAY-PURCHASES` is at most `CR-PURCHASE-LIMIT`, and `RQ-AMOUNT` is more than `WS-AVAILABLE`, it performs `1000-START-CARD-TOTALS`, runs `ADD AO-BALANCE AO-OVERDRAFT-LIMIT GIVING WS-AVAILABLE`, subtracts `AO-HOLDS` from `WS-AVAILABLE`, computes `WS-DAY-PURCHASES` as `CR-DAY-PURCHASES + RQ-AMOUNT`, sets `RS-RESULT` to `RS-INSUFFICIENT-FUNDS`, sets `RS-MESSAGE` to `'INSUFFICIENT FUNDS'`, and returns to its caller.
- **Rule 6.** When `CR-STATUS` is not `'B'` (not `CR-BLOCKED`), `AO-STATUS` is `'O'` (`AO-OPEN`), `RQ-AMOUNT` is not 0, `WS-DAY-PURCHASES` is at most `CR-PURCHASE-LIMIT`, and `RQ-AMOUNT` is at most `WS-AVAILABLE`, it performs `1000-START-CARD-TOTALS`, runs `ADD AO-BALANCE AO-OVERDRAFT-LIMIT GIVING WS-AVAILABLE`, subtracts `AO-HOLDS` from `WS-AVAILABLE`, computes `WS-DAY-PURCHASES` as `CR-DAY-PURCHASES + RQ-AMOUNT`, performs `2000-APPROVE`, and returns to its caller.

Source: [AUTHPUR.cbl:38](../../../pangaea/programs/AUTHPUR.cbl#L38).

## 1000-START-CARD-TOTALS

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `CR-TOTALS-DATE = BC-BUSINESS-DATE` | N | N | Y |
| `CR-TOTALS-DATE(1:6) = BC-BUSINESS-DATE(1:6)` | N | Y | |
| → `MOVE 0 TO CR-MONTH-WITHDRAWALS` | 1 | | |
| → `MOVE 0 TO CR-DAY-PURCHASES CR-DAY-CASH` | 2 | 1 | |
| → `MOVE BC-BUSINESS-DATE TO CR-TOTALS-DATE` | 3 | 2 | |

- **Rule 1.** When `CR-TOTALS-DATE` is not `BC-BUSINESS-DATE` and `CR-TOTALS-DATE(1:6)` is not `BC-BUSINESS-DATE(1:6)`, it sets `CR-MONTH-WITHDRAWALS` to 0; sets `CR-DAY-PURCHASES` and `CR-DAY-CASH` to 0; and sets `CR-TOTALS-DATE` to `BC-BUSINESS-DATE`.
- **Rule 2.** When `CR-TOTALS-DATE` is not `BC-BUSINESS-DATE` and `CR-TOTALS-DATE(1:6)` is `BC-BUSINESS-DATE(1:6)`, it sets `CR-DAY-PURCHASES` and `CR-DAY-CASH` to 0; and sets `CR-TOTALS-DATE` to `BC-BUSINESS-DATE`.
- **Rule 3.** When `CR-TOTALS-DATE` is `BC-BUSINESS-DATE`, it does nothing more.

Source: [AUTHPUR.cbl:65](../../../pangaea/programs/AUTHPUR.cbl#L65).
