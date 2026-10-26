# MATCH: decision tables

The business rules of [MATCH](../programs/MATCH.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

## Contents

- [MAIN](#main)
- [READ-PRESENTMENT](#read-presentment)
- [READ-AUTHORIZATION](#read-authorization)
- [MATCH-NEXT](#match-next)
- [MARK-PRESENTMENT](#mark-presentment)
- [WRITE-PRESENTMENT](#write-presentment)

## MAIN

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-NOT-MATCHED > ZERO` | Y | N |
| → `OPEN INPUT PRES-IN AUTH-IN OUTPUT PRES-OUT REPORT-OUT` | 1 | 1 |
| → `PERFORM READ-PRESENTMENT` | 2 | 2 |
| → `PERFORM READ-AUTHORIZATION` | 3 | 3 |
| → `PERFORM MATCH-NEXT UNTIL WS-PRES-KEY = HIGH-VALUES AND WS-AUTH-KEY = HIGH-VALUES` | 4 | 4 |
| → `PERFORM WRITE-REPORT` | 5 | 5 |
| → `CLOSE PRES-IN AUTH-IN PRES-OUT REPORT-OUT` | 6 | 6 |
| → `STOP RUN RETURNING 4` | 7 | |
| → `STOP RUN RETURNING 0` | | 7 |

- **Rule 1.** When `WS-NOT-MATCHED` is more than zero, it opens `PRES-IN` and `AUTH-IN` for input and `PRES-OUT` and `REPORT-OUT` for output; performs `READ-PRESENTMENT`; performs `READ-AUTHORIZATION`; performs `MATCH-NEXT` until `WS-PRES-KEY` is high values and `WS-AUTH-KEY` is high values; performs `WRITE-REPORT`; closes `PRES-IN`, `AUTH-IN`, `PRES-OUT`, and `REPORT-OUT`; and stops the run.
- **Rule 2.** When `WS-NOT-MATCHED` is at most zero, it opens `PRES-IN` and `AUTH-IN` for input and `PRES-OUT` and `REPORT-OUT` for output; performs `READ-PRESENTMENT`; performs `READ-AUTHORIZATION`; performs `MATCH-NEXT` until `WS-PRES-KEY` is high values and `WS-AUTH-KEY` is high values; performs `WRITE-REPORT`; closes `PRES-IN`, `AUTH-IN`, `PRES-OUT`, and `REPORT-OUT`; and stops the run.

Source: [match.cbl:57](../../../laurasia/programs/match.cbl#L57).

## READ-PRESENTMENT

| | 1 | 2 |
| --- | :-: | :-: |
| `READ PRES-IN INTO PRES-RECORD AT END` | Y | N |
| → `READ PRES-IN INTO PRES-RECORD` | 1 | 1 |
| → `MOVE HIGH-VALUES TO WS-PRES-KEY` | 2 | |
| → `MOVE PS-AUTH-CODE TO WS-PRES-KEY` | | 2 |

- **Rule 1.** When reading `PRES-IN` finds no more records, it reads `PRES-IN` into `PRES-RECORD` and sets `WS-PRES-KEY` to high values.
- **Rule 2.** When reading `PRES-IN` finds a record, it reads `PRES-IN` into `PRES-RECORD` and sets `WS-PRES-KEY` to `PS-AUTH-CODE`.

Source: [match.cbl:70](../../../laurasia/programs/match.cbl#L70).

## READ-AUTHORIZATION

| | 1 | 2 |
| --- | :-: | :-: |
| `READ AUTH-IN AT END` | Y | N |
| → `READ AUTH-IN` | 1 | 1 |
| → `MOVE HIGH-VALUES TO WS-AUTH-KEY` | 2 | |
| → `MOVE 'N' TO WS-AUTH-CLAIMED` | 3 | 3 |
| → `MOVE AL-AUTH-CODE TO WS-AUTH-KEY` | | 2 |

- **Rule 1.** When reading `AUTH-IN` finds no more records, it reads `AUTH-IN`, sets `WS-AUTH-KEY` to high values, and sets `WS-AUTH-CLAIMED` to `'N'`.
- **Rule 2.** When reading `AUTH-IN` finds a record, it reads `AUTH-IN`, sets `WS-AUTH-KEY` to `AL-AUTH-CODE`, and sets `WS-AUTH-CLAIMED` to `'N'`.

Source: [match.cbl:76](../../../laurasia/programs/match.cbl#L76).

## MATCH-NEXT

| | 1 | 2 | 3 | 4 |
| --- | :-: | :-: | :-: | :-: |
| `WS-PRES-KEY < WS-AUTH-KEY` | Y | N | N | N |
| `WS-PRES-KEY > WS-AUTH-KEY` | | Y | Y | N |
| `AUTH-CLAIMED` | | N | Y | |
| → `SET PS-NO-AUTH TO TRUE` | 1 | | | |
| → `PERFORM WRITE-PRESENTMENT` | 2 | | | 3 |
| → `PERFORM READ-PRESENTMENT` | 3 | | | 4 |
| → `ADD 1 TO WS-LAPSED` | | 1 | | |
| → `PERFORM READ-AUTHORIZATION` | | 2 | 1 | |
| → `PERFORM MARK-PRESENTMENT` | | | | 1 |
| → `SET AUTH-CLAIMED TO TRUE` | | | | 2 |

- **Rule 1.** When `WS-PRES-KEY` is less than `WS-AUTH-KEY`, it sets `PS-MATCH` to `PS-NO-AUTH`, performs `WRITE-PRESENTMENT`, and performs `READ-PRESENTMENT`.
- **Rule 2.** When `WS-PRES-KEY` is at least `WS-AUTH-KEY`, `WS-PRES-KEY` is more than `WS-AUTH-KEY`, and `WS-AUTH-CLAIMED` is not `'Y'` (not `AUTH-CLAIMED`), it adds 1 to `WS-LAPSED` and performs `READ-AUTHORIZATION`.
- **Rule 3.** When `WS-PRES-KEY` is at least `WS-AUTH-KEY`, `WS-PRES-KEY` is more than `WS-AUTH-KEY`, and `WS-AUTH-CLAIMED` is `'Y'` (`AUTH-CLAIMED`), it performs `READ-AUTHORIZATION`.
- **Rule 4.** When `WS-PRES-KEY` is at least `WS-AUTH-KEY` and `WS-PRES-KEY` is at most `WS-AUTH-KEY`, it performs `MARK-PRESENTMENT`, sets `WS-AUTH-CLAIMED` to `AUTH-CLAIMED`, performs `WRITE-PRESENTMENT`, and performs `READ-PRESENTMENT`.

Source: [match.cbl:86](../../../laurasia/programs/match.cbl#L86).

## MARK-PRESENTMENT

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `AL-DECLINED` | Y | N | N |
| `PS-AMOUNT > AL-AMOUNT * TOLERANCE` | | Y | N |
| → `SET PS-DECLINED TO TRUE` | 1 | | |
| → `SET PS-OVER-AUTH TO TRUE` | | 1 | |
| → `SET PS-MATCHED TO TRUE` | | | 1 |

- **Rule 1.** When `AL-DECISION` is `'D'` (`AL-DECLINED`), it sets `PS-MATCH` to `PS-DECLINED`.
- **Rule 2.** When `AL-DECISION` is not `'D'` (not `AL-DECLINED`) and `PS-AMOUNT` is more than `AL-AMOUNT * TOLERANCE`, it sets `PS-MATCH` to `PS-OVER-AUTH`.
- **Rule 3.** When `AL-DECISION` is not `'D'` (not `AL-DECLINED`) and `PS-AMOUNT` is at most `AL-AMOUNT * TOLERANCE`, it sets `PS-MATCH` to `PS-MATCHED`.

Source: [match.cbl:104](../../../laurasia/programs/match.cbl#L104).

## WRITE-PRESENTMENT

| | 1 | 2 |
| --- | :-: | :-: |
| `PS-MATCHED` | Y | N |
| → `WRITE PRES-RECORD` | 1 | 1 |
| → `ADD 1 TO WS-MATCHED` | 2 | |
| → `ADD 1 TO WS-NOT-MATCHED` | | 2 |

- **Rule 1.** When `PS-MATCH` is `'M'` (`PS-MATCHED`), it writes `PRES-RECORD` and adds 1 to `WS-MATCHED`.
- **Rule 2.** When `PS-MATCH` is not `'M'` (not `PS-MATCHED`), it writes `PRES-RECORD` and adds 1 to `WS-NOT-MATCHED`.

Source: [match.cbl:114](../../../laurasia/programs/match.cbl#L114).
