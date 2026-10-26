# ACCTLD: decision tables

The business rules of [ACCTLD](../programs/ACCTLD.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

## MAIN

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-ACCOUNTS-STATUS = '00'` | N | Y |
| → `OPEN INPUT ACCOUNTS-IN OUTPUT ACCOUNTS` | 1 | 1 |
| → `STOP RUN RETURNING 16` | 2 | |
| → `PERFORM UNTIL ACCOUNTS-AT-END` | | 2 |
| → `CLOSE ACCOUNTS-IN ACCOUNTS` | | 3 |
| → `STOP RUN RETURNING WS-RESULT` | | 4 |

- **Rule 1.** When `WS-ACCOUNTS-STATUS` is not `'00'`, it opens `ACCOUNTS-IN` for input and `ACCOUNTS` for output; and stops the run.
- **Rule 2.** When `WS-ACCOUNTS-STATUS` is `'00'`, it opens `ACCOUNTS-IN` for input and `ACCOUNTS` for output; repeats its loop until `WS-AT-END` is `'Y'` (`ACCOUNTS-AT-END`); closes `ACCOUNTS-IN` and `ACCOUNTS`; and stops the run.

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

- **Rule 1.** When reading `ACCOUNTS-IN` finds no more records, it reads `ACCOUNTS-IN` and sets `WS-AT-END` to `ACCOUNTS-AT-END`.
- **Rule 2.** When reading `ACCOUNTS-IN` finds a record and writing `CARD-ACCOUNT` fails on its key, it reads `ACCOUNTS-IN`, writes `CARD-ACCOUNT` from `ACCOUNT-LINE`, and sets `WS-RESULT` to 8.
- **Rule 3.** When reading `ACCOUNTS-IN` finds a record and writing `CARD-ACCOUNT` succeeds by its key, it reads `ACCOUNTS-IN` and writes `CARD-ACCOUNT` from `ACCOUNT-LINE`.

Source: [acctld.cbl:41](../../../laurasia/programs/acctld.cbl#L41).
