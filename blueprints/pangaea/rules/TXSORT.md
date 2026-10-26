# TXSORT: decision tables

The business rules of [TXSORT](../programs/TXSORT.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

## 0000-MAIN

| | 1 | 2 |
| --- | :-: | :-: |
| `SORT-RETURN = 0` | N | Y |
| → `SORT SORT-WORK ON ASCENDING KEY SR-ACCT-NO SR-TRAN-ID USING TRAN-IN GIVING TRAN-OUT` | 1 | 1 |
| → `DISPLAY 'TXSORT: SORT FAILED ' SORT-RETURN` | 2 | |
| → `MOVE 16 TO RETURN-CODE` | 3 | |
| → `STOP RUN` | 4 | 3 |
| → `MOVE 0 TO RETURN-CODE` | | 2 |

- **Rule 1.** When `SORT-RETURN` is not 0, it runs `SORT SORT-WORK ON ASCENDING KEY SR-ACCT-NO SR-TRAN-ID USING TRAN-IN GIVING TRAN-OUT`, displays `'TXSORT: SORT FAILED '` `SORT-RETURN`, sets `RETURN-CODE` to 16, and stops the run.
- **Rule 2.** When `SORT-RETURN` is 0, it runs `SORT SORT-WORK ON ASCENDING KEY SR-ACCT-NO SR-TRAN-ID USING TRAN-IN GIVING TRAN-OUT`, sets `RETURN-CODE` to 0, and stops the run.

Source: [TXSORT.cbl:48](../../../pangaea/programs/TXSORT.cbl#L48).
