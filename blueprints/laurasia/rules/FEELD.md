# FEELD: decision tables

The business rules of [FEELD](../programs/FEELD.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

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

- **Rule 1.** When `WS-FEE-STATUS` is not `'00'`, it opens `FEE-SCHEDULE` for input and `FEE-TABLE` for output; and stops the run.
- **Rule 2.** When `WS-FEE-STATUS` is `'00'`, it opens `FEE-SCHEDULE` for input and `FEE-TABLE` for output; repeats its loop until `WS-AT-END` is `'Y'` (`SCHEDULE-AT-END`); closes `FEE-SCHEDULE` and `FEE-TABLE`; and stops the run.

Source: [feeld.cbl:41](../../../laurasia/programs/feeld.cbl#L41).

## MAIN, in its loop

Each pass of `PERFORM UNTIL SCHEDULE-AT-END`.

| | 1 | 2 |
| --- | :-: | :-: |
| `READ FEE-SCHEDULE AT END` | Y | N |
| → `READ FEE-SCHEDULE` | 1 | 1 |
| → `SET SCHEDULE-AT-END TO TRUE` | 2 | |
| → `PERFORM LOAD-FEE` | | 2 |

- **Rule 1.** When reading `FEE-SCHEDULE` finds no more records, it reads `FEE-SCHEDULE` and sets `WS-AT-END` to `SCHEDULE-AT-END`.
- **Rule 2.** When reading `FEE-SCHEDULE` finds a record, it reads `FEE-SCHEDULE` and performs `LOAD-FEE`.

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

- **Rule 1.** When `WS-SLOT` is zero, it calls `FEESLOT`, sets `WS-RESULT` to 8, and leaves the section.
- **Rule 2.** When `WS-SLOT` is not zero and writing `FEE-SLOT-RECORD` fails on its key, it calls `FEESLOT`, sets `FS-RATE` to `FD-RATE`, sets `FS-FIXED-CENTS` to `FD-FIXED`, writes `FEE-SLOT-RECORD`, and sets `WS-RESULT` to 8.
- **Rule 3.** When `WS-SLOT` is not zero and writing `FEE-SLOT-RECORD` succeeds by its key, it calls `FEESLOT`, sets `FS-RATE` to `FD-RATE`, sets `FS-FIXED-CENTS` to `FD-FIXED`, and writes `FEE-SLOT-RECORD`.

Source: [feeld.cbl:55](../../../laurasia/programs/feeld.cbl#L55).
