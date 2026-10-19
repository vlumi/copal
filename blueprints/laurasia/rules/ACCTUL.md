# ACCTUL: decision tables

The business rules of [ACCTUL](../programs/ACCTUL.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

## MAIN

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-ACCOUNTS-STATUS = '00'` | N | Y |
| → `OPEN INPUT ACCOUNTS OUTPUT ACCOUNTS-OUT` | 1 | 1 |
| → `STOP RUN RETURNING 16` | 2 | |
| → `READ ACCOUNTS NEXT RECORD` | | 2 |
| → `PERFORM UNTIL ACCOUNTS-AT-END` | | 3 |
| → `CLOSE ACCOUNTS ACCOUNTS-OUT` | | 4 |
| → `STOP RUN RETURNING 0` | | 5 |

Source: [acctul.cbl:32](../../../laurasia/programs/acctul.cbl#L32).
