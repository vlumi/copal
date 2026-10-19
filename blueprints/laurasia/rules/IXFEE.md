# IXFEE: decision tables

The business rules of [IXFEE](../programs/IXFEE.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

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

Source: [ixfee.cbl:69](../../../laurasia/programs/ixfee.cbl#L69).

## MAIN, in its loop

Each pass of `PERFORM UNTIL PRESENTMENTS-AT-END`.

| | 1 | 2 |
| --- | :-: | :-: |
| `READ PRES-IN INTO PRES-RECORD AT END` | Y | N |
| → `READ PRES-IN INTO PRES-RECORD` | 1 | 1 |
| → `SET PRESENTMENTS-AT-END TO TRUE` | 2 | |
| → `PERFORM PRICE-PRESENTMENT` | | 2 |

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

Source: [ixfee.cbl:105](../../../laurasia/programs/ixfee.cbl#L105).
