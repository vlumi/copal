# FEESLOT: decision tables

The business rules of [FEESLOT](../programs/FEESLOT.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

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

Source: [feeslot.cbl:23](../../../laurasia/programs/feeslot.cbl#L23).
