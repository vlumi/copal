# CLRIN: decision tables

The business rules of [CLRIN](../programs/CLRIN.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

## Contents

- [MAIN](#main)
- [READ-CLEARING](#read-clearing)
- [TAKE-BATCH](#take-batch)
- [TAKE-PRESENTMENT](#take-presentment)
- [VALIDATE-PRESENTMENT](#validate-presentment)
- [REJECT-RECORD](#reject-record)
- [CHECK-FILE-TRAILER](#check-file-trailer)

## MAIN

| | 1 | 2 | 3 | 4 | 5 |
| --- | :-: | :-: | :-: | :-: | :-: |
| `IO-FAILED` | Y | N | N | N | N |
| `CLEARING-AT-END OR NOT IS-FILE-HEADER OR FH-FILE-DATE NOT = WS-RUN-DATE` | | Y | Y | N | N |
| `IO-FAILED, again` | | Y | N | Y | N |
| → `ACCEPT WS-RUN-DATE FROM COMMAND-LINE` | 1 | 1 | 1 | 1 | 1 |
| → `OPEN INPUT CLEARING-IN OUTPUT PRES-OUT REJECT-OUT REPORT-OUT` | 2 | 2 | 2 | 2 | 2 |
| → `STOP RUN RETURNING 16` | 3 | | | | |
| → `PERFORM READ-CLEARING` | | 3 | 3 | 3, 4 | 3, 4 |
| → `MOVE 12 TO WS-RESULT` | | 4 | 4 | | |
| → `PERFORM WRITE-REPORT` | | 5 | 5 | 7 | 7 |
| → `CLOSE CLEARING-IN PRES-OUT REJECT-OUT REPORT-OUT` | | 6 | 6 | 8 | 8 |
| → `MOVE 16 TO WS-RESULT` | | 7 | | 9 | |
| → `STOP RUN RETURNING WS-RESULT` | | 8 | 7 | 10 | 9 |
| → `PERFORM TAKE-BATCH UNTIL CLEARING-AT-END OR IS-FILE-TRAILER OR IO-FAILED` | | | | 5 | 5 |
| → `PERFORM CHECK-FILE-TRAILER` | | | | 6 | 6 |

Source: [clrin.cbl:94](../../../laurasia/programs/clrin.cbl#L94).

## READ-CLEARING

| | 1 | 2 |
| --- | :-: | :-: |
| `READ CLEARING-IN AT END` | Y | N |
| → `READ CLEARING-IN` | 1 | 1 |
| → `SET CLEARING-AT-END TO TRUE` | 2 | |

Source: [clrin.cbl:116](../../../laurasia/programs/clrin.cbl#L116).

## TAKE-BATCH

| | 1 | 2 | 3 | 4 |
| --- | :-: | :-: | :-: | :-: |
| `IS-BATCH-HEADER` | N | Y | Y | Y |
| `IS-BATCH-TRAILER AND BT-COUNT = WS-BATCH-COUNT AND BT-AMOUNT = WS-BATCH-AMOUNT` | | Y | N | N |
| `IS-BATCH-TRAILER` | | Y | Y | N |
| → `MOVE 'ORDR' TO WS-REASON` | 1 | | | |
| → `PERFORM REJECT-RECORD` | 2 | | | |
| → `PERFORM READ-CLEARING` | 3 | 4, 9 | 4, 9 | 4 |
| → `EXIT SECTION` | 4 | | | |
| → `MOVE BH-BATCH-NO TO WS-BATCH-NO` | | 1 | 1 | 1 |
| → `MOVE ZERO TO WS-BATCH-COUNT WS-BATCH-AMOUNT WS-HELD-COUNT` | | 2 | 2 | 2 |
| → `ADD 1 TO WS-BATCHES` | | 3 | 3 | 3 |
| → `PERFORM TAKE-PRESENTMENT UNTIL CLEARING-AT-END OR NOT IS-PRESENTMENT` | | 5 | 5 | 5 |
| → `MOVE CLR-PRESENTMENT TO WS-RECORD-READ` | | 6 | 6 | 6 |
| → `PERFORM RELEASE-BATCH` | | 7 | | |
| → `MOVE WS-RECORD-READ TO CLR-PRESENTMENT` | | 8 | 8 | 8 |
| → `PERFORM REJECT-BATCH` | | | 7 | 7 |

Source: [clrin.cbl:124](../../../laurasia/programs/clrin.cbl#L124).

## TAKE-PRESENTMENT

| | 1 | 2 |
| --- | :-: | :-: |
| `PRESENTMENT-IS-GOOD AND WS-HELD-COUNT < MAX-BATCH` | Y | N |
| → `ADD 1 TO WS-BATCH-COUNT WS-RECORDS` | 1 | 1 |
| → `ADD PR-AMOUNT TO WS-BATCH-AMOUNT WS-AMOUNT` | 2 | 2 |
| → `PERFORM VALIDATE-PRESENTMENT` | 3 | 3 |
| → `ADD 1 TO WS-HELD-COUNT` | 4 | |
| → `MOVE CLR-PRESENTMENT TO WS-HELD-PRESENTMENT(WS-HELD-COUNT)` | 5 | |
| → `PERFORM READ-CLEARING` | 6 | 5 |
| → `PERFORM REJECT-RECORD` | | 4 |

Source: [clrin.cbl:147](../../../laurasia/programs/clrin.cbl#L147).

## VALIDATE-PRESENTMENT

| | 1 | 2 | 3 | 4 | 5 | 6 |
| --- | :-: | :-: | :-: | :-: | :-: | :-: |
| `PR-CARD-NO IS NUMERIC` | N | Y | Y | Y | Y | Y |
| `PR-AMOUNT > ZERO` | | N | Y | Y | Y | Y |
| `PR-MCC IS NUMERIC` | | | N | Y | Y | Y |
| `PR-REGION NOT = 'D' AND PR-REGION NOT = 'I'` | | | | Y | N | N |
| `PR-PRODUCT NOT = 'CL' AND PR-PRODUCT NOT = 'GD' AND PR-PRODUCT NOT = 'BZ'` | | | | | Y | N |
| → `MOVE 'CARD' TO WS-REASON` | 1 | | | | | |
| → `MOVE 'AMNT' TO WS-REASON` | | 1 | | | | |
| → `MOVE 'MCC ' TO WS-REASON` | | | 1 | | | |
| → `MOVE 'REGN' TO WS-REASON` | | | | 1 | | |
| → `MOVE 'PROD' TO WS-REASON` | | | | | 1 | |
| → `MOVE SPACES TO WS-REASON` | | | | | | 1 |

Source: [clrin.cbl:159](../../../laurasia/programs/clrin.cbl#L159).

## REJECT-RECORD

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-RESULT < 4` | Y | N |
| → `MOVE SPACES TO REJECT-LINE` | 1 | 1 |
| → `MOVE WS-REASON TO RJ-REASON` | 2 | 2 |
| → `MOVE WS-BATCH-NO TO RJ-BATCH-NO` | 3 | 3 |
| → `MOVE CLR-PRESENTMENT TO RJ-RECORD` | 4 | 4 |
| → `WRITE REJECT-LINE` | 5 | 5 |
| → `ADD 1 TO WS-REJECTED` | 6 | 6 |
| → `MOVE 4 TO WS-RESULT` | 7 | |

Source: [clrin.cbl:200](../../../laurasia/programs/clrin.cbl#L200).

## CHECK-FILE-TRAILER

| | 1 | 2 |
| --- | :-: | :-: |
| `NOT IS-FILE-TRAILER OR FT-BATCH-COUNT NOT = WS-BATCHES OR FT-RECORD-COUNT NOT = WS-RECORDS OR FT-AMOUNT NOT = WS-AMOUNT` | Y | N |
| → `MOVE 8 TO WS-RESULT` | 1 | |

Source: [clrin.cbl:211](../../../laurasia/programs/clrin.cbl#L211).
