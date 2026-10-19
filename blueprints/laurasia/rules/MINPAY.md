# MINPAY: decision tables

The business rules of [MINPAY](../programs/MINPAY.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

## MAIN

| | 1 | 2 | 3 | 4 |
| --- | :-: | :-: | :-: | :-: |
| `LK-BALANCE > ZERO` | N | Y | Y | Y |
| `LK-BALANCE > PAYMENT-FLOOR` | | N | Y | Y |
| `LK-MINIMUM < PAYMENT-FLOOR` | | | Y | N |
| → `MOVE ZERO TO LK-MINIMUM` | 1 | | | |
| → `GOBACK` | 2 | 2 | 3 | 2 |
| → `MOVE LK-BALANCE TO LK-MINIMUM` | | 1 | | |
| → `COMPUTE LK-MINIMUM ROUNDED = LK-BALANCE / 50` | | | 1 | 1 |
| → `MOVE PAYMENT-FLOOR TO LK-MINIMUM` | | | 2 | |

Source: [cycle.cbl:137](../../../laurasia/programs/cycle.cbl#L137).
