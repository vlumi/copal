# MINPAY: decision tables

The business rules of [MINPAY](../programs/MINPAY.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

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

- **Rule 1.** When `LK-BALANCE` is at most zero, it sets `LK-MINIMUM` to zero and returns to its caller.
- **Rule 2.** When `LK-BALANCE` is more than zero and `LK-BALANCE` is at most `PAYMENT-FLOOR`, it sets `LK-MINIMUM` to `LK-BALANCE` and returns to its caller.
- **Rule 3.** When `LK-BALANCE` is more than zero, `LK-BALANCE` is more than `PAYMENT-FLOOR`, and `LK-MINIMUM` is less than `PAYMENT-FLOOR`, it computes `LK-MINIMUM` as `LK-BALANCE / 50`, sets `LK-MINIMUM` to `PAYMENT-FLOOR`, and returns to its caller.
- **Rule 4.** When `LK-BALANCE` is more than zero, `LK-BALANCE` is more than `PAYMENT-FLOOR`, and `LK-MINIMUM` is at least `PAYMENT-FLOOR`, it computes `LK-MINIMUM` as `LK-BALANCE / 50` and returns to its caller.

Source: [cycle.cbl:139](../../../laurasia/programs/cycle.cbl#L139).
