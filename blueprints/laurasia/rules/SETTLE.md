# SETTLE: decision tables

The business rules of [SETTLE](../programs/SETTLE.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

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

Source: [settle.cbl:72](../../../laurasia/programs/settle.cbl#L72).

## READ-PRESENTMENT

| | 1 | 2 |
| --- | :-: | :-: |
| `READ PRES-IN AT END` | Y | N |
| → `READ PRES-IN` | 1 | 1 |
| → `MOVE HIGH-VALUES TO WS-PRES-KEY` | 2 | |
| → `MOVE PS-MERCHANT TO WS-PRES-KEY` | | 2 |

Source: [settle.cbl:85](../../../laurasia/programs/settle.cbl#L85).

## READ-ADJUSTMENT

| | 1 | 2 |
| --- | :-: | :-: |
| `READ ADJUST-IN AT END` | Y | N |
| → `READ ADJUST-IN` | 1 | 1 |
| → `MOVE HIGH-VALUES TO WS-ADJUST-KEY` | 2 | |
| → `MOVE AD-MERCHANT TO WS-ADJUST-KEY` | | 2 |

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

Source: [settle.cbl:99](../../../laurasia/programs/settle.cbl#L99).

## SETTLE-MERCHANT, in its loop

Each pass of `PERFORM UNTIL WS-PRES-KEY NOT = WS-MERCHANT`.

| | 1 | 2 |
| --- | :-: | :-: |
| `PS-MATCHED` | Y | N |
| → `ADD PS-AMOUNT TO WS-GROSS` | 1 | |
| → `ADD PS-FEE TO WS-INTERCHANGE` | 2 | |
| → `PERFORM READ-PRESENTMENT` | 3 | 1 |

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

Source: [settle.cbl:119](../../../laurasia/programs/settle.cbl#L119).
