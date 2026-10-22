# CYCLE: decision tables

The business rules of [CYCLE](../programs/CYCLE.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

## Contents

- [MAIN](#main)
- [MAIN, in its loop](#main-in-its-loop)
- [CLOSE-CYCLE](#close-cycle)

## MAIN

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-ACCOUNTS-STATUS = '00'` | N | Y |
| → `ACCEPT WS-RUN-DATE FROM COMMAND-LINE` | 1 | 1 |
| → `OPEN I-O ACCOUNTS OUTPUT POSTINGS-OUT STATEMENTS-OUT` | 2 | 2 |
| → `STOP RUN RETURNING 16` | 3 | |
| → `READ ACCOUNTS NEXT RECORD` | | 3 |
| → `PERFORM UNTIL ACCOUNTS-AT-END` | | 4 |
| → `CLOSE ACCOUNTS POSTINGS-OUT STATEMENTS-OUT` | | 5 |
| → `STOP RUN RETURNING 0` | | 6 |

Source: [cycle.cbl:63](../../../laurasia/programs/cycle.cbl#L63).

## MAIN, in its loop

Each pass of `PERFORM UNTIL ACCOUNTS-AT-END`.

| | 1 | 2 |
| --- | :-: | :-: |
| `CA-CYCLE-DAY = WS-RUN-DAY AND NOT CA-CLOSED` | Y | N |
| → `PERFORM CLOSE-CYCLE` | 1 | |
| → `READ ACCOUNTS NEXT RECORD` | 2 | 1 |

Source: [cycle.cbl:71](../../../laurasia/programs/cycle.cbl#L71).

## CLOSE-CYCLE

| | 1 | 2 | 3 | 4 | 5 | 6 |
| --- | :-: | :-: | :-: | :-: | :-: | :-: |
| `NOT CA-WAS-PAID-IN-FULL AND CA-STATEMENT-BALANCE > ZERO` | Y | Y | Y | Y | N | N |
| `WS-INTEREST > ZERO` | Y | Y | N | N | N | N |
| `REWRITE CARD-ACCOUNT INVALID KEY` | Y | N | Y | N | Y | N |
| → `MOVE ZERO TO WS-INTEREST` | 1 | 1 | 1 | 1 | 1 | 1 |
| → `CALL 'intcalc' USING CA-PRODUCT CA-STATEMENT-BALANCE RETURNING WS-INTEREST-CENTS` | 2 | 2 | 2 | 2 | | |
| → `COMPUTE WS-INTEREST = WS-INTEREST-CENTS / 100` | 3 | 3 | 3 | 3 | | |
| → `ADD WS-INTEREST TO CA-BALANCE` | 4 | 4 | | | | |
| → `PERFORM POST-INTEREST` | 5 | 5 | | | | |
| → `CALL 'minpay' USING CA-BALANCE WS-MINIMUM` | 6 | 6 | 4 | 4 | 2 | 2 |
| → `PERFORM WRITE-STATEMENT` | 7 | 7 | 5 | 5 | 3 | 3 |
| → `MOVE CA-BALANCE TO CA-STATEMENT-BALANCE` | 8 | 8 | 6 | 6 | 4 | 4 |
| → `MOVE 'N' TO CA-PAID-IN-FULL` | 9 | 9 | 7 | 7 | 5 | 5 |
| → `REWRITE CARD-ACCOUNT` | 10 | 10 | 8 | 8 | 6 | 6 |
| → `STOP RUN RETURNING 16` | 11 | | 9 | | 7 | |
| → `ADD 1 TO WS-CLOSED-COUNT` | | 11 | | 9 | | 7 |

Source: [cycle.cbl:80](../../../laurasia/programs/cycle.cbl#L80).
