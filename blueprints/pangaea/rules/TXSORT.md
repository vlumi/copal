# TXSORT: decision tables

The business rules of [TXSORT](../programs/TXSORT.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

## 0000-MAIN

| | 1 | 2 |
| --- | :-: | :-: |
| `SORT-RETURN = 0` | N | Y |
| → `SORT SORT-WORK ON ASCENDING KEY SR-ACCT-NO SR-TRAN-ID USING TRAN-IN GIVING TRAN-OUT` | 1 | 1 |
| → `DISPLAY 'TXSORT: SORT FAILED ' SORT-RETURN` | 2 | |
| → `MOVE 16 TO RETURN-CODE` | 3 | |
| → `STOP RUN` | 4 | 3 |
| → `MOVE 0 TO RETURN-CODE` | | 2 |

Source: [TXSORT.cbl:48](../../../pangaea/programs/TXSORT.cbl#L48).
