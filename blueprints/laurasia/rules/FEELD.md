# FEELD: decision tables

The business rules of [FEELD](../programs/FEELD.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

## Contents

- [MAIN](#main)
- [MAIN, in its loop](#main-in-its-loop)
- [LOAD-FEE](#load-fee)

## MAIN

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-FEE-STATUS = '00'` | N | Y |
| → `OPEN INPUT FEE-SCHEDULE OUTPUT FEE-TABLE` | 1 | 1 |
| → `STOP RUN RETURNING 16` | 2 | |
| → `PERFORM UNTIL SCHEDULE-AT-END` | | 2 |
| → `CLOSE FEE-SCHEDULE FEE-TABLE` | | 3 |
| → `STOP RUN RETURNING WS-RESULT` | | 4 |

Source: [feeld.cbl:41](../../../laurasia/programs/feeld.cbl#L41).

## MAIN, in its loop

Each pass of `PERFORM UNTIL SCHEDULE-AT-END`.

| | 1 | 2 |
| --- | :-: | :-: |
| `READ FEE-SCHEDULE AT END` | Y | N |
| → `READ FEE-SCHEDULE` | 1 | 1 |
| → `SET SCHEDULE-AT-END TO TRUE` | 2 | |
| → `PERFORM LOAD-FEE` | | 2 |

Source: [feeld.cbl:46](../../../laurasia/programs/feeld.cbl#L46).

## LOAD-FEE

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `WS-SLOT = ZERO` | Y | N | N |
| `WRITE FEE-SLOT-RECORD INVALID KEY` | | Y | N |
| → `CALL 'feeslot' USING FD-PRODUCT FD-REGION FD-GROUP RETURNING WS-SLOT` | 1 | 1 | 1 |
| → `MOVE 8 TO WS-RESULT` | 2 | 5 | |
| → `EXIT SECTION` | 3 | | |
| → `MOVE FD-RATE TO FS-RATE` | | 2 | 2 |
| → `MOVE FD-FIXED TO FS-FIXED-CENTS` | | 3 | 3 |
| → `WRITE FEE-SLOT-RECORD` | | 4 | 4 |

Source: [feeld.cbl:55](../../../laurasia/programs/feeld.cbl#L55).
