# FEESLOT: decision tables

The business rules of [FEESLOT](../programs/FEESLOT.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

## MAIN

| | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 | 16 |
| --- | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: |
| `LK-PRODUCT = 'CL'` | Y | Y | Y | Y | Y | N | N | N | N | N | N | N | N | N | N | N |
| `LK-REGION = 'D'` | Y | Y | N | N | N | Y | Y | N | N | N | Y | Y | N | N | N | |
| `LK-GROUP < 1 OR LK-GROUP > 5` | Y | N | Y | N | | Y | N | Y | N | | Y | N | Y | N | | |
| `LK-REGION = 'I'` | | | Y | Y | N | | | Y | Y | N | | | Y | Y | N | |
| `LK-PRODUCT = 'GD'` | | | | | | Y | Y | Y | Y | Y | N | N | N | N | N | N |
| `LK-PRODUCT = 'BZ'` | | | | | | | | | | | Y | Y | Y | Y | Y | N |
| → `MOVE ZERO TO RETURN-CODE` | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 |
| → `MOVE 1 TO WS-PRODUCT-INDEX` | 2 | 2 | 2 | 2 | 2 | | | | | | | | | | | |
| → `MOVE 1 TO WS-REGION-INDEX` | 3 | 3 | | | | 3 | 3 | | | | 3 | 3 | | | | |
| → `GOBACK` | 4 | 6 | 4 | 6 | 3 | 4 | 6 | 4 | 6 | 3 | 4 | 6 | 4 | 6 | 3 | 2 |
| → `COMPUTE WS-SLOT = ((WS-PRODUCT-INDEX - 1) * 2 + WS-REGION-INDEX - 1) * 5 + LK-GROUP` | | 4 | | 4 | | | 4 | | 4 | | | 4 | | 4 | | |
| → `MOVE WS-SLOT TO RETURN-CODE` | | 5 | | 5 | | | 5 | | 5 | | | 5 | | 5 | | |
| → `MOVE 2 TO WS-REGION-INDEX` | | | 3 | 3 | | | | 3 | 3 | | | | 3 | 3 | | |
| → `MOVE 2 TO WS-PRODUCT-INDEX` | | | | | | 2 | 2 | 2 | 2 | 2 | | | | | | |
| → `MOVE 3 TO WS-PRODUCT-INDEX` | | | | | | | | | | | 2 | 2 | 2 | 2 | 2 | |

- **Rule 1.** When `LK-PRODUCT` is `'CL'`; `LK-REGION` is `'D'`; and `LK-GROUP` is less than 1 or `LK-GROUP` is more than 5, it sets `RETURN-CODE` to zero, sets `WS-PRODUCT-INDEX` to 1, sets `WS-REGION-INDEX` to 1, and returns to its caller.
- **Rule 2.** When `LK-PRODUCT` is `'CL'`; `LK-REGION` is `'D'`; and `LK-GROUP` is at least 1 and `LK-GROUP` is at most 5, it sets `RETURN-CODE` to zero, sets `WS-PRODUCT-INDEX` to 1, sets `WS-REGION-INDEX` to 1, computes `WS-SLOT` as `((WS-PRODUCT-INDEX - 1) * 2 + WS-REGION-INDEX - 1) * 5 + LK-GROUP`, sets `RETURN-CODE` to `WS-SLOT`, and returns to its caller.
- **Rule 3.** When `LK-PRODUCT` is `'CL'`; `LK-REGION` is not `'D'`; `LK-GROUP` is less than 1 or `LK-GROUP` is more than 5; and `LK-REGION` is `'I'`, it sets `RETURN-CODE` to zero, sets `WS-PRODUCT-INDEX` to 1, sets `WS-REGION-INDEX` to 2, and returns to its caller.
- **Rule 4.** When `LK-PRODUCT` is `'CL'`; `LK-REGION` is not `'D'`; `LK-GROUP` is at least 1 and `LK-GROUP` is at most 5; and `LK-REGION` is `'I'`, it sets `RETURN-CODE` to zero, sets `WS-PRODUCT-INDEX` to 1, sets `WS-REGION-INDEX` to 2, computes `WS-SLOT` as `((WS-PRODUCT-INDEX - 1) * 2 + WS-REGION-INDEX - 1) * 5 + LK-GROUP`, sets `RETURN-CODE` to `WS-SLOT`, and returns to its caller.
- **Rule 5.** When `LK-PRODUCT` is `'CL'`, `LK-REGION` is not `'D'`, and `LK-REGION` is not `'I'`, it sets `RETURN-CODE` to zero, sets `WS-PRODUCT-INDEX` to 1, and returns to its caller.
- **Rule 6.** When `LK-PRODUCT` is not `'CL'`; `LK-REGION` is `'D'`; `LK-GROUP` is less than 1 or `LK-GROUP` is more than 5; and `LK-PRODUCT` is `'GD'`, it sets `RETURN-CODE` to zero, sets `WS-PRODUCT-INDEX` to 2, sets `WS-REGION-INDEX` to 1, and returns to its caller.
- **Rule 7.** When `LK-PRODUCT` is not `'CL'`; `LK-REGION` is `'D'`; `LK-GROUP` is at least 1 and `LK-GROUP` is at most 5; and `LK-PRODUCT` is `'GD'`, it sets `RETURN-CODE` to zero, sets `WS-PRODUCT-INDEX` to 2, sets `WS-REGION-INDEX` to 1, computes `WS-SLOT` as `((WS-PRODUCT-INDEX - 1) * 2 + WS-REGION-INDEX - 1) * 5 + LK-GROUP`, sets `RETURN-CODE` to `WS-SLOT`, and returns to its caller.
- **Rule 8.** When `LK-PRODUCT` is not `'CL'`; `LK-REGION` is not `'D'`; `LK-GROUP` is less than 1 or `LK-GROUP` is more than 5; `LK-REGION` is `'I'`; and `LK-PRODUCT` is `'GD'`, it sets `RETURN-CODE` to zero, sets `WS-PRODUCT-INDEX` to 2, sets `WS-REGION-INDEX` to 2, and returns to its caller.
- **Rule 9.** When `LK-PRODUCT` is not `'CL'`; `LK-REGION` is not `'D'`; `LK-GROUP` is at least 1 and `LK-GROUP` is at most 5; `LK-REGION` is `'I'`; and `LK-PRODUCT` is `'GD'`, it sets `RETURN-CODE` to zero, sets `WS-PRODUCT-INDEX` to 2, sets `WS-REGION-INDEX` to 2, computes `WS-SLOT` as `((WS-PRODUCT-INDEX - 1) * 2 + WS-REGION-INDEX - 1) * 5 + LK-GROUP`, sets `RETURN-CODE` to `WS-SLOT`, and returns to its caller.
- **Rule 10.** When `LK-PRODUCT` is not `'CL'`, `LK-REGION` is not `'D'`, `LK-REGION` is not `'I'`, and `LK-PRODUCT` is `'GD'`, it sets `RETURN-CODE` to zero, sets `WS-PRODUCT-INDEX` to 2, and returns to its caller.
- **Rule 11.** When `LK-PRODUCT` is not `'CL'`; `LK-REGION` is `'D'`; `LK-GROUP` is less than 1 or `LK-GROUP` is more than 5; `LK-PRODUCT` is not `'GD'`; and `LK-PRODUCT` is `'BZ'`, it sets `RETURN-CODE` to zero, sets `WS-PRODUCT-INDEX` to 3, sets `WS-REGION-INDEX` to 1, and returns to its caller.
- **Rule 12.** When `LK-PRODUCT` is not `'CL'`; `LK-REGION` is `'D'`; `LK-GROUP` is at least 1 and `LK-GROUP` is at most 5; `LK-PRODUCT` is not `'GD'`; and `LK-PRODUCT` is `'BZ'`, it sets `RETURN-CODE` to zero, sets `WS-PRODUCT-INDEX` to 3, sets `WS-REGION-INDEX` to 1, computes `WS-SLOT` as `((WS-PRODUCT-INDEX - 1) * 2 + WS-REGION-INDEX - 1) * 5 + LK-GROUP`, sets `RETURN-CODE` to `WS-SLOT`, and returns to its caller.
- **Rule 13.** When `LK-PRODUCT` is not `'CL'`; `LK-REGION` is not `'D'`; `LK-GROUP` is less than 1 or `LK-GROUP` is more than 5; `LK-REGION` is `'I'`; `LK-PRODUCT` is not `'GD'`; and `LK-PRODUCT` is `'BZ'`, it sets `RETURN-CODE` to zero, sets `WS-PRODUCT-INDEX` to 3, sets `WS-REGION-INDEX` to 2, and returns to its caller.
- **Rule 14.** When `LK-PRODUCT` is not `'CL'`; `LK-REGION` is not `'D'`; `LK-GROUP` is at least 1 and `LK-GROUP` is at most 5; `LK-REGION` is `'I'`; `LK-PRODUCT` is not `'GD'`; and `LK-PRODUCT` is `'BZ'`, it sets `RETURN-CODE` to zero, sets `WS-PRODUCT-INDEX` to 3, sets `WS-REGION-INDEX` to 2, computes `WS-SLOT` as `((WS-PRODUCT-INDEX - 1) * 2 + WS-REGION-INDEX - 1) * 5 + LK-GROUP`, sets `RETURN-CODE` to `WS-SLOT`, and returns to its caller.
- **Rule 15.** When `LK-PRODUCT` is not `'CL'`, `LK-REGION` is not `'D'`, `LK-REGION` is not `'I'`, `LK-PRODUCT` is not `'GD'`, and `LK-PRODUCT` is `'BZ'`, it sets `RETURN-CODE` to zero, sets `WS-PRODUCT-INDEX` to 3, and returns to its caller.
- **Rule 16.** When `LK-PRODUCT` is not `'CL'`, `LK-PRODUCT` is not `'GD'`, and `LK-PRODUCT` is not `'BZ'`, it sets `RETURN-CODE` to zero and returns to its caller.

Source: [feeslot.cbl:23](../../../laurasia/programs/feeslot.cbl#L23).
