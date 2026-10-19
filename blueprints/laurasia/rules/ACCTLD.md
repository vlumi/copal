# ACCTLD: decision tables

The business rules of [ACCTLD](../programs/ACCTLD.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

## MAIN

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-ACCOUNTS-STATUS = '00'` | N | Y |
| → `OPEN INPUT ACCOUNTS-IN OUTPUT ACCOUNTS` | 1 | 1 |
| → `STOP RUN RETURNING 16` | 2 | |
| → `PERFORM UNTIL ACCOUNTS-AT-END` | | 2 |
| → `CLOSE ACCOUNTS-IN ACCOUNTS` | | 3 |
| → `STOP RUN RETURNING WS-RESULT` | | 4 |

Source: [acctld.cbl:36](../../../laurasia/programs/acctld.cbl#L36).

## MAIN, in its loop

Each pass of `PERFORM UNTIL ACCOUNTS-AT-END`.

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `READ ACCOUNTS-IN AT END` | Y | N | N |
| `WRITE CARD-ACCOUNT FROM ACCOUNT-LINE INVALID KEY` | | Y | N |
| → `READ ACCOUNTS-IN` | 1 | 1 | 1 |
| → `SET ACCOUNTS-AT-END TO TRUE` | 2 | | |
| → `WRITE CARD-ACCOUNT FROM ACCOUNT-LINE` | | 2 | 2 |
| → `MOVE 8 TO WS-RESULT` | | 3 | |

Source: [acctld.cbl:41](../../../laurasia/programs/acctld.cbl#L41).
