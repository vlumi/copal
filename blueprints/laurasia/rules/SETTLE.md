# SETTLE: decision tables

The business rules of [SETTLE](../programs/SETTLE.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

## Contents

- [MAIN](#main)
- [READ-PRESENTMENT](#read-presentment)
- [READ-ADJUSTMENT](#read-adjustment)
- [SETTLE-MERCHANT](#settle-merchant)
- [SETTLE-MERCHANT, in its loop](#settle-merchant-in-its-loop)
- [WRITE-PAYOUT](#write-payout)

## MAIN

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-OWING-COUNT > ZERO` | Y | N |
| → `OPEN INPUT PRES-IN ADJUST-IN OUTPUT PAYOUT-OUT REPORT-OUT` | 1 | 1 |
| → `PERFORM READ-PRESENTMENT` | 2 | 2 |
| → `PERFORM READ-ADJUSTMENT` | 3 | 3 |
| → `PERFORM SETTLE-MERCHANT UNTIL WS-PRES-KEY = HIGH-VALUES AND WS-ADJUST-KEY = HIGH-VALUES` | 4 | 4 |
| → `PERFORM WRITE-REPORT` | 5 | 5 |
| → `CLOSE PRES-IN ADJUST-IN PAYOUT-OUT REPORT-OUT` | 6 | 6 |
| → `STOP RUN RETURNING 4` | 7 | |
| → `STOP RUN RETURNING 0` | | 7 |

- **Rule 1.** When `WS-OWING-COUNT` is more than zero, it opens `PRES-IN` and `ADJUST-IN` for input and `PAYOUT-OUT` and `REPORT-OUT` for output; performs `READ-PRESENTMENT`; performs `READ-ADJUSTMENT`; performs `SETTLE-MERCHANT` until `WS-PRES-KEY` is high values and `WS-ADJUST-KEY` is high values; performs `WRITE-REPORT`; closes `PRES-IN`, `ADJUST-IN`, `PAYOUT-OUT`, and `REPORT-OUT`; and stops the run.
- **Rule 2.** When `WS-OWING-COUNT` is at most zero, it opens `PRES-IN` and `ADJUST-IN` for input and `PAYOUT-OUT` and `REPORT-OUT` for output; performs `READ-PRESENTMENT`; performs `READ-ADJUSTMENT`; performs `SETTLE-MERCHANT` until `WS-PRES-KEY` is high values and `WS-ADJUST-KEY` is high values; performs `WRITE-REPORT`; closes `PRES-IN`, `ADJUST-IN`, `PAYOUT-OUT`, and `REPORT-OUT`; and stops the run.

Source: [settle.cbl:72](../../../laurasia/programs/settle.cbl#L72).

## READ-PRESENTMENT

| | 1 | 2 |
| --- | :-: | :-: |
| `READ PRES-IN AT END` | Y | N |
| → `READ PRES-IN` | 1 | 1 |
| → `MOVE HIGH-VALUES TO WS-PRES-KEY` | 2 | |
| → `MOVE PS-MERCHANT TO WS-PRES-KEY` | | 2 |

- **Rule 1.** When reading `PRES-IN` finds no more records, it reads `PRES-IN` and sets `WS-PRES-KEY` to high values.
- **Rule 2.** When reading `PRES-IN` finds a record, it reads `PRES-IN` and sets `WS-PRES-KEY` to `PS-MERCHANT`.

Source: [settle.cbl:85](../../../laurasia/programs/settle.cbl#L85).

## READ-ADJUSTMENT

| | 1 | 2 |
| --- | :-: | :-: |
| `READ ADJUST-IN AT END` | Y | N |
| → `READ ADJUST-IN` | 1 | 1 |
| → `MOVE HIGH-VALUES TO WS-ADJUST-KEY` | 2 | |
| → `MOVE AD-MERCHANT TO WS-ADJUST-KEY` | | 2 |

- **Rule 1.** When reading `ADJUST-IN` finds no more records, it reads `ADJUST-IN` and sets `WS-ADJUST-KEY` to high values.
- **Rule 2.** When reading `ADJUST-IN` finds a record, it reads `ADJUST-IN` and sets `WS-ADJUST-KEY` to `AD-MERCHANT`.

Source: [settle.cbl:91](../../../laurasia/programs/settle.cbl#L91).

## SETTLE-MERCHANT

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-PRES-KEY < WS-ADJUST-KEY` | Y | N |
| → `MOVE WS-PRES-KEY TO WS-MERCHANT` | 1 | |
| → `MOVE ZERO TO WS-GROSS WS-INTERCHANGE WS-ADJUSTMENTS` | 2 | 2 |
| → `PERFORM UNTIL WS-PRES-KEY NOT = WS-MERCHANT` | 3 | 3 |
| → `PERFORM UNTIL WS-ADJUST-KEY NOT = WS-MERCHANT` | 4 | 4 |
| → `PERFORM WRITE-PAYOUT` | 5 | 5 |
| → `MOVE WS-ADJUST-KEY TO WS-MERCHANT` | | 1 |

- **Rule 1.** When `WS-PRES-KEY` is less than `WS-ADJUST-KEY`, it sets `WS-MERCHANT` to `WS-PRES-KEY`; sets `WS-GROSS`, `WS-INTERCHANGE`, and `WS-ADJUSTMENTS` to zero; repeats its loop until `WS-PRES-KEY` is not `WS-MERCHANT`; repeats its loop until `WS-ADJUST-KEY` is not `WS-MERCHANT`; and performs `WRITE-PAYOUT`.
- **Rule 2.** When `WS-PRES-KEY` is at least `WS-ADJUST-KEY`, it sets `WS-MERCHANT` to `WS-ADJUST-KEY`; sets `WS-GROSS`, `WS-INTERCHANGE`, and `WS-ADJUSTMENTS` to zero; repeats its loop until `WS-PRES-KEY` is not `WS-MERCHANT`; repeats its loop until `WS-ADJUST-KEY` is not `WS-MERCHANT`; and performs `WRITE-PAYOUT`.

Source: [settle.cbl:99](../../../laurasia/programs/settle.cbl#L99).

## SETTLE-MERCHANT, in its loop

Each pass of `PERFORM UNTIL WS-PRES-KEY NOT = WS-MERCHANT`.

| | 1 | 2 |
| --- | :-: | :-: |
| `PS-MATCHED` | Y | N |
| → `ADD PS-AMOUNT TO WS-GROSS` | 1 | |
| → `ADD PS-FEE TO WS-INTERCHANGE` | 2 | |
| → `PERFORM READ-PRESENTMENT` | 3 | 1 |

- **Rule 1.** When `PS-MATCH` is `'M'` (`PS-MATCHED`), it adds `PS-AMOUNT` to `WS-GROSS`, adds `PS-FEE` to `WS-INTERCHANGE`, and performs `READ-PRESENTMENT`.
- **Rule 2.** When `PS-MATCH` is not `'M'` (not `PS-MATCHED`), it performs `READ-PRESENTMENT`.

Source: [settle.cbl:106](../../../laurasia/programs/settle.cbl#L106).

## WRITE-PAYOUT

| | 1 | 2 |
| --- | :-: | :-: |
| `PA-NET < ZERO` | Y | N |
| → `MOVE WS-MERCHANT TO PA-MERCHANT` | 1 | 1 |
| → `MOVE WS-GROSS TO PA-GROSS` | 2 | 2 |
| → `MOVE WS-INTERCHANGE TO PA-INTERCHANGE` | 3 | 3 |
| → `MOVE WS-ADJUSTMENTS TO PA-ADJUSTMENTS` | 4 | 4 |
| → `COMPUTE PA-NET = WS-GROSS - WS-INTERCHANGE + WS-ADJUSTMENTS` | 5 | 5 |
| → `SET PA-OWES TO TRUE` | 6 | |
| → `ADD 1 TO WS-OWING-COUNT` | 7 | |
| → `WRITE PAYOUT` | 8 | 9 |
| → `SET PA-PAID TO TRUE` | | 6 |
| → `ADD 1 TO WS-PAID-COUNT` | | 7 |
| → `ADD PA-NET TO WS-PAID-TOTAL` | | 8 |

- **Rule 1.** When `PA-NET` is less than zero, it sets `PA-MERCHANT` to `WS-MERCHANT`, sets `PA-GROSS` to `WS-GROSS`, sets `PA-INTERCHANGE` to `WS-INTERCHANGE`, sets `PA-ADJUSTMENTS` to `WS-ADJUSTMENTS`, computes `PA-NET` as `WS-GROSS - WS-INTERCHANGE + WS-ADJUSTMENTS`, sets `PA-STATUS` to `PA-OWES`, adds 1 to `WS-OWING-COUNT`, and writes `PAYOUT`.
- **Rule 2.** When `PA-NET` is at least zero, it sets `PA-MERCHANT` to `WS-MERCHANT`, sets `PA-GROSS` to `WS-GROSS`, sets `PA-INTERCHANGE` to `WS-INTERCHANGE`, sets `PA-ADJUSTMENTS` to `WS-ADJUSTMENTS`, computes `PA-NET` as `WS-GROSS - WS-INTERCHANGE + WS-ADJUSTMENTS`, sets `PA-STATUS` to `PA-PAID`, adds 1 to `WS-PAID-COUNT`, adds `PA-NET` to `WS-PAID-TOTAL`, and writes `PAYOUT`.

Source: [settle.cbl:119](../../../laurasia/programs/settle.cbl#L119).
