# CLRIN: decision tables

The business rules of [CLRIN](../programs/CLRIN.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

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

- **Rule 1.** When `WS-IO-FAILED` is `'Y'` (`IO-FAILED`), it accepts `WS-RUN-DATE` from its command line; opens `CLEARING-IN` for input and `PRES-OUT`, `REJECT-OUT`, and `REPORT-OUT` for output; and stops the run.
- **Rule 2.** When `WS-IO-FAILED` is not `'Y'` (not `IO-FAILED`); `WS-CLEARING-STATUS` is `'10'` (`CLEARING-AT-END`), `CLR-RECORD-TYPE` is not `'FH'` (not `IS-FILE-HEADER`), or `FH-FILE-DATE` is not `WS-RUN-DATE`; and `WS-IO-FAILED` is `'Y'` (`IO-FAILED`), tested again, it accepts `WS-RUN-DATE` from its command line; opens `CLEARING-IN` for input and `PRES-OUT`, `REJECT-OUT`, and `REPORT-OUT` for output; performs `READ-CLEARING`; sets `WS-RESULT` to 12; performs `WRITE-REPORT`; closes `CLEARING-IN`, `PRES-OUT`, `REJECT-OUT`, and `REPORT-OUT`; sets `WS-RESULT` to 16; and stops the run.
- **Rule 3.** When `WS-IO-FAILED` is not `'Y'` (not `IO-FAILED`); `WS-CLEARING-STATUS` is `'10'` (`CLEARING-AT-END`), `CLR-RECORD-TYPE` is not `'FH'` (not `IS-FILE-HEADER`), or `FH-FILE-DATE` is not `WS-RUN-DATE`; and `WS-IO-FAILED` is not `'Y'` (not `IO-FAILED`), tested again, it accepts `WS-RUN-DATE` from its command line; opens `CLEARING-IN` for input and `PRES-OUT`, `REJECT-OUT`, and `REPORT-OUT` for output; performs `READ-CLEARING`; sets `WS-RESULT` to 12; performs `WRITE-REPORT`; closes `CLEARING-IN`, `PRES-OUT`, `REJECT-OUT`, and `REPORT-OUT`; and stops the run.
- **Rule 4.** When `WS-IO-FAILED` is not `'Y'` (not `IO-FAILED`); `WS-CLEARING-STATUS` is not `'10'` (not `CLEARING-AT-END`), `CLR-RECORD-TYPE` is `'FH'` (`IS-FILE-HEADER`), and `FH-FILE-DATE` is `WS-RUN-DATE`; and `WS-IO-FAILED` is `'Y'` (`IO-FAILED`), tested again, it accepts `WS-RUN-DATE` from its command line; opens `CLEARING-IN` for input and `PRES-OUT`, `REJECT-OUT`, and `REPORT-OUT` for output; performs `READ-CLEARING`; performs `READ-CLEARING`; performs `TAKE-BATCH` until `WS-CLEARING-STATUS` is `'10'` (`CLEARING-AT-END`), `CLR-RECORD-TYPE` is `'FT'` (`IS-FILE-TRAILER`), or `WS-IO-FAILED` is `'Y'` (`IO-FAILED`); performs `CHECK-FILE-TRAILER`; performs `WRITE-REPORT`; closes `CLEARING-IN`, `PRES-OUT`, `REJECT-OUT`, and `REPORT-OUT`; sets `WS-RESULT` to 16; and stops the run.
- **Rule 5.** When `WS-IO-FAILED` is not `'Y'` (not `IO-FAILED`); `WS-CLEARING-STATUS` is not `'10'` (not `CLEARING-AT-END`), `CLR-RECORD-TYPE` is `'FH'` (`IS-FILE-HEADER`), and `FH-FILE-DATE` is `WS-RUN-DATE`; and `WS-IO-FAILED` is not `'Y'` (not `IO-FAILED`), tested again, it accepts `WS-RUN-DATE` from its command line; opens `CLEARING-IN` for input and `PRES-OUT`, `REJECT-OUT`, and `REPORT-OUT` for output; performs `READ-CLEARING`; performs `READ-CLEARING`; performs `TAKE-BATCH` until `WS-CLEARING-STATUS` is `'10'` (`CLEARING-AT-END`), `CLR-RECORD-TYPE` is `'FT'` (`IS-FILE-TRAILER`), or `WS-IO-FAILED` is `'Y'` (`IO-FAILED`); performs `CHECK-FILE-TRAILER`; performs `WRITE-REPORT`; closes `CLEARING-IN`, `PRES-OUT`, `REJECT-OUT`, and `REPORT-OUT`; and stops the run.

Source: [clrin.cbl:94](../../../laurasia/programs/clrin.cbl#L94).

## READ-CLEARING

| | 1 | 2 |
| --- | :-: | :-: |
| `READ CLEARING-IN AT END` | Y | N |
| → `READ CLEARING-IN` | 1 | 1 |
| → `SET CLEARING-AT-END TO TRUE` | 2 | |

- **Rule 1.** When reading `CLEARING-IN` finds no more records, it reads `CLEARING-IN` and sets `WS-CLEARING-STATUS` to `CLEARING-AT-END`.
- **Rule 2.** When reading `CLEARING-IN` finds a record, it reads `CLEARING-IN`.

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

- **Rule 1.** When `CLR-RECORD-TYPE` is not `'BH'` (not `IS-BATCH-HEADER`), it sets `WS-REASON` to `'ORDR'`, performs `REJECT-RECORD`, performs `READ-CLEARING`, and leaves the section.
- **Rule 2.** When `CLR-RECORD-TYPE` is `'BH'` (`IS-BATCH-HEADER`); `CLR-RECORD-TYPE` is `'BT'` (`IS-BATCH-TRAILER`), `BT-COUNT` is `WS-BATCH-COUNT`, and `BT-AMOUNT` is `WS-BATCH-AMOUNT`; and `CLR-RECORD-TYPE` is `'BT'` (`IS-BATCH-TRAILER`), it sets `WS-BATCH-NO` to `BH-BATCH-NO`; sets `WS-BATCH-COUNT`, `WS-BATCH-AMOUNT`, and `WS-HELD-COUNT` to zero; adds 1 to `WS-BATCHES`; performs `READ-CLEARING`; performs `TAKE-PRESENTMENT` until `WS-CLEARING-STATUS` is `'10'` (`CLEARING-AT-END`) or `CLR-RECORD-TYPE` is not `'PR'` (not `IS-PRESENTMENT`); sets `WS-RECORD-READ` to `CLR-PRESENTMENT`; performs `RELEASE-BATCH`; sets `CLR-PRESENTMENT` to `WS-RECORD-READ`; and performs `READ-CLEARING`.
- **Rule 3.** When `CLR-RECORD-TYPE` is `'BH'` (`IS-BATCH-HEADER`); `CLR-RECORD-TYPE` is not `'BT'` (not `IS-BATCH-TRAILER`), `BT-COUNT` is not `WS-BATCH-COUNT`, or `BT-AMOUNT` is not `WS-BATCH-AMOUNT`; and `CLR-RECORD-TYPE` is `'BT'` (`IS-BATCH-TRAILER`), it sets `WS-BATCH-NO` to `BH-BATCH-NO`; sets `WS-BATCH-COUNT`, `WS-BATCH-AMOUNT`, and `WS-HELD-COUNT` to zero; adds 1 to `WS-BATCHES`; performs `READ-CLEARING`; performs `TAKE-PRESENTMENT` until `WS-CLEARING-STATUS` is `'10'` (`CLEARING-AT-END`) or `CLR-RECORD-TYPE` is not `'PR'` (not `IS-PRESENTMENT`); sets `WS-RECORD-READ` to `CLR-PRESENTMENT`; performs `REJECT-BATCH`; sets `CLR-PRESENTMENT` to `WS-RECORD-READ`; and performs `READ-CLEARING`.
- **Rule 4.** When `CLR-RECORD-TYPE` is `'BH'` (`IS-BATCH-HEADER`); `CLR-RECORD-TYPE` is not `'BT'` (not `IS-BATCH-TRAILER`), `BT-COUNT` is not `WS-BATCH-COUNT`, or `BT-AMOUNT` is not `WS-BATCH-AMOUNT`; and `CLR-RECORD-TYPE` is not `'BT'` (not `IS-BATCH-TRAILER`), it sets `WS-BATCH-NO` to `BH-BATCH-NO`; sets `WS-BATCH-COUNT`, `WS-BATCH-AMOUNT`, and `WS-HELD-COUNT` to zero; adds 1 to `WS-BATCHES`; performs `READ-CLEARING`; performs `TAKE-PRESENTMENT` until `WS-CLEARING-STATUS` is `'10'` (`CLEARING-AT-END`) or `CLR-RECORD-TYPE` is not `'PR'` (not `IS-PRESENTMENT`); sets `WS-RECORD-READ` to `CLR-PRESENTMENT`; performs `REJECT-BATCH`; and sets `CLR-PRESENTMENT` to `WS-RECORD-READ`.

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

- **Rule 1.** When `WS-REASON` is spaces (`PRESENTMENT-IS-GOOD`) and `WS-HELD-COUNT` is less than `MAX-BATCH`, it adds 1 to `WS-BATCH-COUNT` and `WS-RECORDS`; adds `PR-AMOUNT` to `WS-BATCH-AMOUNT` and `WS-AMOUNT`; performs `VALIDATE-PRESENTMENT`; adds 1 to `WS-HELD-COUNT`; sets `WS-HELD-PRESENTMENT(WS-HELD-COUNT)` to `CLR-PRESENTMENT`; and performs `READ-CLEARING`.
- **Rule 2.** When `WS-REASON` is not spaces (not `PRESENTMENT-IS-GOOD`) or `WS-HELD-COUNT` is at least `MAX-BATCH`, it adds 1 to `WS-BATCH-COUNT` and `WS-RECORDS`; adds `PR-AMOUNT` to `WS-BATCH-AMOUNT` and `WS-AMOUNT`; performs `VALIDATE-PRESENTMENT`; performs `REJECT-RECORD`; and performs `READ-CLEARING`.

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

- **Rule 1.** When `PR-CARD-NO` is not numeric, it sets `WS-REASON` to `'CARD'`.
- **Rule 2.** When `PR-CARD-NO` is numeric and `PR-AMOUNT` is at most zero, it sets `WS-REASON` to `'AMNT'`.
- **Rule 3.** When `PR-CARD-NO` is numeric, `PR-AMOUNT` is more than zero, and `PR-MCC` is not numeric, it sets `WS-REASON` to `'MCC '`.
- **Rule 4.** When `PR-CARD-NO` is numeric; `PR-AMOUNT` is more than zero; `PR-MCC` is numeric; and `PR-REGION` is not `'D'` and `PR-REGION` is not `'I'`, it sets `WS-REASON` to `'REGN'`.
- **Rule 5.** When `PR-CARD-NO` is numeric; `PR-AMOUNT` is more than zero; `PR-MCC` is numeric; `PR-REGION` is `'D'` or `PR-REGION` is `'I'`; and `PR-PRODUCT` is not `'CL'`, `PR-PRODUCT` is not `'GD'`, and `PR-PRODUCT` is not `'BZ'`, it sets `WS-REASON` to `'PROD'`.
- **Rule 6.** When `PR-CARD-NO` is numeric; `PR-AMOUNT` is more than zero; `PR-MCC` is numeric; `PR-REGION` is `'D'` or `PR-REGION` is `'I'`; and `PR-PRODUCT` is `'CL'`, `PR-PRODUCT` is `'GD'`, or `PR-PRODUCT` is `'BZ'`, it sets `WS-REASON` to spaces.

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

- **Rule 1.** When `WS-RESULT` is less than 4, it sets `REJECT-LINE` to spaces, sets `RJ-REASON` to `WS-REASON`, sets `RJ-BATCH-NO` to `WS-BATCH-NO`, sets `RJ-RECORD` to `CLR-PRESENTMENT`, writes `REJECT-LINE`, adds 1 to `WS-REJECTED`, and sets `WS-RESULT` to 4.
- **Rule 2.** When `WS-RESULT` is at least 4, it sets `REJECT-LINE` to spaces, sets `RJ-REASON` to `WS-REASON`, sets `RJ-BATCH-NO` to `WS-BATCH-NO`, sets `RJ-RECORD` to `CLR-PRESENTMENT`, writes `REJECT-LINE`, and adds 1 to `WS-REJECTED`.

Source: [clrin.cbl:200](../../../laurasia/programs/clrin.cbl#L200).

## CHECK-FILE-TRAILER

| | 1 | 2 |
| --- | :-: | :-: |
| `NOT IS-FILE-TRAILER OR FT-BATCH-COUNT NOT = WS-BATCHES OR FT-RECORD-COUNT NOT = WS-RECORDS OR FT-AMOUNT NOT = WS-AMOUNT` | Y | N |
| → `MOVE 8 TO WS-RESULT` | 1 | |

- **Rule 1.** When `CLR-RECORD-TYPE` is not `'FT'` (not `IS-FILE-TRAILER`), `FT-BATCH-COUNT` is not `WS-BATCHES`, `FT-RECORD-COUNT` is not `WS-RECORDS`, or `FT-AMOUNT` is not `WS-AMOUNT`, it sets `WS-RESULT` to 8.
- **Rule 2.** When `CLR-RECORD-TYPE` is `'FT'` (`IS-FILE-TRAILER`), `FT-BATCH-COUNT` is `WS-BATCHES`, `FT-RECORD-COUNT` is `WS-RECORDS`, and `FT-AMOUNT` is `WS-AMOUNT`, it does nothing more.

Source: [clrin.cbl:211](../../../laurasia/programs/clrin.cbl#L211).
