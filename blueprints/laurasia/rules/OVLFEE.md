# OVLFEE: decision tables

The business rules of [OVLFEE](../programs/OVLFEE.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

## MAIN

| | 1 | 2 | 3 | 4 |
| --- | :-: | :-: | :-: | :-: |
| `LK-GOLD AND LK-OVER-AMOUNT NOT > GOLD-ALLOWANCE` | Y | N | N | N |
| `LK-FEE < LEAST-FEE` | | Y | N | N |
| `LK-FEE > MOST-FEE` | | | Y | N |
| → `MOVE ZERO TO LK-FEE` | 1 | | | |
| → `GOBACK` | 2 | 3 | 3 | 2 |
| → `COMPUTE LK-FEE ROUNDED = LK-OVER-AMOUNT / 10` | | 1 | 1 | 1 |
| → `MOVE LEAST-FEE TO LK-FEE` | | 2 | | |
| → `MOVE MOST-FEE TO LK-FEE` | | | 2 | |

Source: [ovlfee.cbl:23](../../../laurasia/programs/ovlfee.cbl#L23).
