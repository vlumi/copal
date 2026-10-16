# ATMBAL: decision tables

The business rules of [ATMBAL](../programs/ATMBAL.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

## 0000-MAIN

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `CR-BLOCKED` | Y | N | N |
| `AO-OPEN` | | N | Y |
| → `SET RS-CARD-BLOCKED TO TRUE` | 1 | | |
| → `MOVE 'CARD BLOCKED' TO RS-MESSAGE` | 2 | | |
| → `GOBACK` | 3 | 3 | 7 |
| → `SET RS-ACCOUNT-NOT-OPEN TO TRUE` | | 1 | |
| → `MOVE 'ACCOUNT NOT OPEN' TO RS-MESSAGE` | | 2 | |
| → `SET RS-APPROVED TO TRUE` | | | 1 |
| → `MOVE 'BALANCE' TO RS-MESSAGE` | | | 2 |
| → `MOVE AO-BALANCE TO RS-BALANCE` | | | 3 |
| → `ADD AO-BALANCE AO-OVERDRAFT-LIMIT GIVING WS-AVAILABLE` | | | 4 |
| → `SUBTRACT AO-HOLDS FROM WS-AVAILABLE` | | | 5 |
| → `MOVE WS-AVAILABLE TO RS-AVAILABLE` | | | 6 |

Source: [ATMBAL.cbl:34](../../../pangaea/programs/ATMBAL.cbl#L34).
