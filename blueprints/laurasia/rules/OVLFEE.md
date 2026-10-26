# OVLFEE: decision tables

The business rules of [OVLFEE](../programs/OVLFEE.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

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

- **Rule 1.** When `LK-PRODUCT` is `'GD'` (`LK-GOLD`) and `LK-OVER-AMOUNT` is at most `GOLD-ALLOWANCE`, it sets `LK-FEE` to zero and returns to its caller.
- **Rule 2.** When `LK-PRODUCT` is not `'GD'` (not `LK-GOLD`) or `LK-OVER-AMOUNT` is more than `GOLD-ALLOWANCE`; and `LK-FEE` is less than `LEAST-FEE`, it computes `LK-FEE` as `LK-OVER-AMOUNT / 10`, sets `LK-FEE` to `LEAST-FEE`, and returns to its caller.
- **Rule 3.** When `LK-PRODUCT` is not `'GD'` (not `LK-GOLD`) or `LK-OVER-AMOUNT` is more than `GOLD-ALLOWANCE`; `LK-FEE` is at least `LEAST-FEE`; and `LK-FEE` is more than `MOST-FEE`, it computes `LK-FEE` as `LK-OVER-AMOUNT / 10`, sets `LK-FEE` to `MOST-FEE`, and returns to its caller.
- **Rule 4.** When `LK-PRODUCT` is not `'GD'` (not `LK-GOLD`) or `LK-OVER-AMOUNT` is more than `GOLD-ALLOWANCE`; `LK-FEE` is at least `LEAST-FEE`; and `LK-FEE` is at most `MOST-FEE`, it computes `LK-FEE` as `LK-OVER-AMOUNT / 10` and returns to its caller.

Source: [ovlfee.cbl:23](../../../laurasia/programs/ovlfee.cbl#L23).
