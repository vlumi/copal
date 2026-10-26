# AUTHREV: decision tables

The business rules of [AUTHREV](../programs/AUTHREV.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

## 0000-MAIN

| | 1 | 2 | 3 | 4 |
| --- | :-: | :-: | :-: | :-: |
| `HD-CARD-NO = CR-CARD-NO` | N | Y | Y | Y |
| `HD-PURCHASE` | | N | Y | Y |
| `HD-REVERSED` | | | Y | N |
| → `SET RS-NO-SUCH-AUTH TO TRUE` | 1 | | 1 | |
| → `MOVE 'AUTHORIZATION OF ANOTHER CARD' TO RS-MESSAGE` | 2 | | | |
| → `GOBACK` | 3 | 3 | 3 | 2 |
| → `SET RS-INVALID-REQUEST TO TRUE` | | 1 | | |
| → `MOVE 'ONLY PURCHASES ARE REVERSED' TO RS-MESSAGE` | | 2 | | |
| → `MOVE 'ALREADY REVERSED' TO RS-MESSAGE` | | | 2 | |
| → `PERFORM 2000-APPROVE` | | | | 1 |

- **Rule 1.** When `HD-CARD-NO` is not `CR-CARD-NO`, it sets `RS-RESULT` to `RS-NO-SUCH-AUTH`, sets `RS-MESSAGE` to `'AUTHORIZATION OF ANOTHER CARD'`, and returns to its caller.
- **Rule 2.** When `HD-CARD-NO` is `CR-CARD-NO` and `HD-KIND` is not `'P'` (not `HD-PURCHASE`), it sets `RS-RESULT` to `RS-INVALID-REQUEST`, sets `RS-MESSAGE` to `'ONLY PURCHASES ARE REVERSED'`, and returns to its caller.
- **Rule 3.** When `HD-CARD-NO` is `CR-CARD-NO`, `HD-KIND` is `'P'` (`HD-PURCHASE`), and `HD-STATUS` is `'R'` (`HD-REVERSED`), it sets `RS-RESULT` to `RS-NO-SUCH-AUTH`, sets `RS-MESSAGE` to `'ALREADY REVERSED'`, and returns to its caller.
- **Rule 4.** When `HD-CARD-NO` is `CR-CARD-NO`, `HD-KIND` is `'P'` (`HD-PURCHASE`), and `HD-STATUS` is not `'R'` (not `HD-REVERSED`), it performs `2000-APPROVE` and returns to its caller.

Source: [AUTHREV.cbl:37](../../../pangaea/programs/AUTHREV.cbl#L37).

## 2000-APPROVE

| | 1 | 2 |
| --- | :-: | :-: |
| `CR-TOTALS-DATE = HD-BUSINESS-DATE` | Y | N |
| → `SUBTRACT HD-AMOUNT FROM AO-HOLDS` | 1 | 1 |
| → `SUBTRACT HD-AMOUNT FROM CR-DAY-PURCHASES` | 2 | |
| → `SET HD-REVERSED TO TRUE` | 3 | 2 |
| → `SET RS-APPROVED TO TRUE` | 4 | 3 |
| → `MOVE 'REVERSED' TO RS-MESSAGE` | 5 | 4 |
| → `MOVE HD-HOLD-NO TO RS-AUTH-CODE` | 6 | 5 |
| → `MOVE HD-AMOUNT TO RS-AMOUNT` | 7 | 6 |
| → `MOVE 0 TO RS-FEE` | 8 | 7 |
| → `MOVE AO-BALANCE TO RS-BALANCE` | 9 | 8 |
| → `ADD AO-BALANCE AO-OVERDRAFT-LIMIT GIVING WS-AVAILABLE` | 10 | 9 |
| → `SUBTRACT AO-HOLDS FROM WS-AVAILABLE` | 11 | 10 |
| → `MOVE WS-AVAILABLE TO RS-AVAILABLE` | 12 | 11 |

- **Rule 1.** When `CR-TOTALS-DATE` is `HD-BUSINESS-DATE`, it subtracts `HD-AMOUNT` from `AO-HOLDS`, subtracts `HD-AMOUNT` from `CR-DAY-PURCHASES`, sets `HD-STATUS` to `HD-REVERSED`, sets `RS-RESULT` to `RS-APPROVED`, sets `RS-MESSAGE` to `'REVERSED'`, sets `RS-AUTH-CODE` to `HD-HOLD-NO`, sets `RS-AMOUNT` to `HD-AMOUNT`, sets `RS-FEE` to 0, sets `RS-BALANCE` to `AO-BALANCE`, runs `ADD AO-BALANCE AO-OVERDRAFT-LIMIT GIVING WS-AVAILABLE`, subtracts `AO-HOLDS` from `WS-AVAILABLE`, and sets `RS-AVAILABLE` to `WS-AVAILABLE`.
- **Rule 2.** When `CR-TOTALS-DATE` is not `HD-BUSINESS-DATE`, it subtracts `HD-AMOUNT` from `AO-HOLDS`, sets `HD-STATUS` to `HD-REVERSED`, sets `RS-RESULT` to `RS-APPROVED`, sets `RS-MESSAGE` to `'REVERSED'`, sets `RS-AUTH-CODE` to `HD-HOLD-NO`, sets `RS-AMOUNT` to `HD-AMOUNT`, sets `RS-FEE` to 0, sets `RS-BALANCE` to `AO-BALANCE`, runs `ADD AO-BALANCE AO-OVERDRAFT-LIMIT GIVING WS-AVAILABLE`, subtracts `AO-HOLDS` from `WS-AVAILABLE`, and sets `RS-AVAILABLE` to `WS-AVAILABLE`.

Source: [AUTHREV.cbl:55](../../../pangaea/programs/AUTHREV.cbl#L55).
