# INTCALC: decision tables

The business rules of [INTCALC](../programs/INTCALC.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

## MAIN

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `LK-GOLD` | Y | N | N |
| `LK-BUSINESS` | | Y | N |
| → `MOVE 0.1690 TO WS-YEARLY-RATE` | 1 | | |
| → `COMPUTE WS-CENTS ROUNDED = LK-BALANCE * WS-YEARLY-RATE / 12 * 100` | 2 | 2 | 2 |
| → `MOVE WS-CENTS TO RETURN-CODE` | 3 | 3 | 3 |
| → `GOBACK` | 4 | 4 | 4 |
| → `MOVE 0.1490 TO WS-YEARLY-RATE` | | 1 | |
| → `MOVE 0.1990 TO WS-YEARLY-RATE` | | | 1 |

- **Rule 1.** When `LK-PRODUCT` is `'GD'` (`LK-GOLD`), it sets `WS-YEARLY-RATE` to 0.1690, computes `WS-CENTS` as `LK-BALANCE * WS-YEARLY-RATE / 12 * 100`, sets `RETURN-CODE` to `WS-CENTS`, and returns to its caller.
- **Rule 2.** When `LK-PRODUCT` is not `'GD'` (not `LK-GOLD`) and `LK-PRODUCT` is `'BZ'` (`LK-BUSINESS`), it sets `WS-YEARLY-RATE` to 0.1490, computes `WS-CENTS` as `LK-BALANCE * WS-YEARLY-RATE / 12 * 100`, sets `RETURN-CODE` to `WS-CENTS`, and returns to its caller.
- **Rule 3.** When `LK-PRODUCT` is not `'GD'` (not `LK-GOLD`) and `LK-PRODUCT` is not `'BZ'` (not `LK-BUSINESS`), it sets `WS-YEARLY-RATE` to 0.1990, computes `WS-CENTS` as `LK-BALANCE * WS-YEARLY-RATE / 12 * 100`, sets `RETURN-CODE` to `WS-CENTS`, and returns to its caller.

Source: [intcalc.cbl:23](../../../laurasia/programs/intcalc.cbl#L23).
