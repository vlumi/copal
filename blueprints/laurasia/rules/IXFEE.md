# IXFEE: decision tables

The business rules of [IXFEE](../programs/IXFEE.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

## Contents

- [MAIN](#main)
- [MAIN, in its loop](#main-in-its-loop)
- [PRICE-PRESENTMENT](#price-presentment)
- [FIND-GROUP](#find-group)

## MAIN

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `WS-FEE-STATUS = '00'` | N | Y | Y |
| `WS-DEFAULT-COUNT > ZERO` | | Y | N |
| → `OPEN INPUT PRES-IN FEE-TABLE OUTPUT PRES-OUT REPORT-OUT` | 1 | 1 | 1 |
| → `STOP RUN RETURNING 16` | 2 | | |
| → `PERFORM UNTIL PRESENTMENTS-AT-END` | | 2 | 2 |
| → `PERFORM WRITE-REPORT` | | 3 | 3 |
| → `CLOSE PRES-IN FEE-TABLE PRES-OUT REPORT-OUT` | | 4 | 4 |
| → `STOP RUN RETURNING 4` | | 5 | |
| → `STOP RUN RETURNING 0` | | | 5 |

- **Rule 1.** When `WS-FEE-STATUS` is not `'00'`, it opens `PRES-IN` and `FEE-TABLE` for input and `PRES-OUT` and `REPORT-OUT` for output; and stops the run.
- **Rule 2.** When `WS-FEE-STATUS` is `'00'` and `WS-DEFAULT-COUNT` is more than zero, it opens `PRES-IN` and `FEE-TABLE` for input and `PRES-OUT` and `REPORT-OUT` for output; repeats its loop until `WS-AT-END` is `'Y'` (`PRESENTMENTS-AT-END`); performs `WRITE-REPORT`; closes `PRES-IN`, `FEE-TABLE`, `PRES-OUT`, and `REPORT-OUT`; and stops the run.
- **Rule 3.** When `WS-FEE-STATUS` is `'00'` and `WS-DEFAULT-COUNT` is at most zero, it opens `PRES-IN` and `FEE-TABLE` for input and `PRES-OUT` and `REPORT-OUT` for output; repeats its loop until `WS-AT-END` is `'Y'` (`PRESENTMENTS-AT-END`); performs `WRITE-REPORT`; closes `PRES-IN`, `FEE-TABLE`, `PRES-OUT`, and `REPORT-OUT`; and stops the run.

Source: [ixfee.cbl:69](../../../laurasia/programs/ixfee.cbl#L69).

## MAIN, in its loop

Each pass of `PERFORM UNTIL PRESENTMENTS-AT-END`.

| | 1 | 2 |
| --- | :-: | :-: |
| `READ PRES-IN INTO PRES-RECORD AT END` | Y | N |
| → `READ PRES-IN INTO PRES-RECORD` | 1 | 1 |
| → `SET PRESENTMENTS-AT-END TO TRUE` | 2 | |
| → `PERFORM PRICE-PRESENTMENT` | | 2 |

- **Rule 1.** When reading `PRES-IN` finds no more records, it reads `PRES-IN` into `PRES-RECORD` and sets `WS-AT-END` to `PRESENTMENTS-AT-END`.
- **Rule 2.** When reading `PRES-IN` finds a record, it reads `PRES-IN` into `PRES-RECORD` and performs `PRICE-PRESENTMENT`.

Source: [ixfee.cbl:75](../../../laurasia/programs/ixfee.cbl#L75).

## PRICE-PRESENTMENT

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `PS-MATCHED` | Y | Y | N |
| `READ FEE-TABLE INVALID KEY` | Y | N | |
| → `MOVE ZERO TO PS-FEE` | 1 | 1 | 1 |
| → `PERFORM FIND-GROUP` | 2 | 2 | |
| → `CALL 'feeslot' USING PS-PRODUCT PS-REGION WS-GROUP RETURNING WS-SLOT` | 3 | 3 | |
| → `READ FEE-TABLE` | 4 | 4 | |
| → `ADD 1 TO WS-DEFAULT-COUNT` | 5 | | |
| → `COMPUTE PS-FEE ROUNDED = PS-AMOUNT * DEFAULT-RATE + DEFAULT-FIXED-CENTS / 100` | 6 | | |
| → `ADD PS-FEE TO WS-FEES` | 7 | 6 | |
| → `WRITE PRES-RECORD` | 8 | 7 | 2 |
| → `COMPUTE PS-FEE ROUNDED = PS-AMOUNT * FS-RATE + FS-FIXED-CENTS / 100` | | 5 | |

- **Rule 1.** When `PS-MATCH` is `'M'` (`PS-MATCHED`) and reading `FEE-TABLE` fails on its key, it sets `PS-FEE` to zero, performs `FIND-GROUP`, calls `FEESLOT`, reads `FEE-TABLE`, adds 1 to `WS-DEFAULT-COUNT`, computes `PS-FEE` as `PS-AMOUNT * DEFAULT-RATE + DEFAULT-FIXED-CENTS / 100`, adds `PS-FEE` to `WS-FEES`, and writes `PRES-RECORD`.
- **Rule 2.** When `PS-MATCH` is `'M'` (`PS-MATCHED`) and reading `FEE-TABLE` succeeds by its key, it sets `PS-FEE` to zero, performs `FIND-GROUP`, calls `FEESLOT`, reads `FEE-TABLE`, computes `PS-FEE` as `PS-AMOUNT * FS-RATE + FS-FIXED-CENTS / 100`, adds `PS-FEE` to `WS-FEES`, and writes `PRES-RECORD`.
- **Rule 3.** When `PS-MATCH` is not `'M'` (not `PS-MATCHED`), it sets `PS-FEE` to zero and writes `PRES-RECORD`.

Source: [ixfee.cbl:88](../../../laurasia/programs/ixfee.cbl#L88).

## FIND-GROUP

| | 1 | 2 | 3 | 4 | 5 |
| --- | :-: | :-: | :-: | :-: | :-: |
| `PS-MCC >= 5411 AND PS-MCC <= 5499` | Y | N | N | N | N |
| `PS-MCC >= 5811 AND PS-MCC <= 5814` | | Y | N | N | N |
| `(PS-MCC >= 3000 AND PS-MCC <= 3999) OR (PS-MCC >= 4111 AND PS-MCC <= 4789)` | | | Y | N | N |
| `PS-MCC = 5541 OR PS-MCC = 5542 OR PS-MCC = 5983` | | | | Y | N |
| → `SET GROUP-GROCERY TO TRUE` | 1 | | | | |
| → `SET GROUP-RESTAURANTS TO TRUE` | | 1 | | | |
| → `SET GROUP-TRAVEL TO TRUE` | | | 1 | | |
| → `SET GROUP-FUEL TO TRUE` | | | | 1 | |
| → `SET GROUP-OTHER TO TRUE` | | | | | 1 |

- **Rule 1.** When `PS-MCC` is at least 5411 and `PS-MCC` is at most 5499, it sets `WS-GROUP` to `GROUP-GROCERY`.
- **Rule 2.** When `PS-MCC` is less than 5411 or `PS-MCC` is more than 5499; and `PS-MCC` is at least 5811 and `PS-MCC` is at most 5814, it sets `WS-GROUP` to `GROUP-RESTAURANTS`.
- **Rule 3.** When `PS-MCC` is less than 5411 or `PS-MCC` is more than 5499; `PS-MCC` is less than 5811 or `PS-MCC` is more than 5814; and both `PS-MCC` is at least 3000 and `PS-MCC` is at most 3999 or both `PS-MCC` is at least 4111 and `PS-MCC` is at most 4789, it sets `WS-GROUP` to `GROUP-TRAVEL`.
- **Rule 4.** When `PS-MCC` is less than 5411 or `PS-MCC` is more than 5499; `PS-MCC` is less than 5811 or `PS-MCC` is more than 5814; either `PS-MCC` is less than 3000 or `PS-MCC` is more than 3999 and either `PS-MCC` is less than 4111 or `PS-MCC` is more than 4789; and `PS-MCC` is 5541, `PS-MCC` is 5542, or `PS-MCC` is 5983, it sets `WS-GROUP` to `GROUP-FUEL`.
- **Rule 5.** When `PS-MCC` is less than 5411 or `PS-MCC` is more than 5499; `PS-MCC` is less than 5811 or `PS-MCC` is more than 5814; either `PS-MCC` is less than 3000 or `PS-MCC` is more than 3999 and either `PS-MCC` is less than 4111 or `PS-MCC` is more than 4789; and `PS-MCC` is not 5541, `PS-MCC` is not 5542, and `PS-MCC` is not 5983, it sets `WS-GROUP` to `GROUP-OTHER`.

Source: [ixfee.cbl:105](../../../laurasia/programs/ixfee.cbl#L105).
