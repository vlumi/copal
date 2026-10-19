# MATCH: decision tables

The business rules of [MATCH](../programs/MATCH.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

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

Source: [match.cbl:57](../../../laurasia/programs/match.cbl#L57).

## READ-PRESENTMENT

| | 1 | 2 |
| --- | :-: | :-: |
| `READ PRES-IN INTO PRES-RECORD AT END` | Y | N |
| → `READ PRES-IN INTO PRES-RECORD` | 1 | 1 |
| → `MOVE HIGH-VALUES TO WS-PRES-KEY` | 2 | |
| → `MOVE PS-AUTH-CODE TO WS-PRES-KEY` | | 2 |

Source: [match.cbl:70](../../../laurasia/programs/match.cbl#L70).

## READ-AUTHORIZATION

| | 1 | 2 |
| --- | :-: | :-: |
| `READ AUTH-IN AT END` | Y | N |
| → `READ AUTH-IN` | 1 | 1 |
| → `MOVE HIGH-VALUES TO WS-AUTH-KEY` | 2 | |
| → `MOVE 'N' TO WS-AUTH-CLAIMED` | 3 | 3 |
| → `MOVE AL-AUTH-CODE TO WS-AUTH-KEY` | | 2 |

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

Source: [match.cbl:86](../../../laurasia/programs/match.cbl#L86).

## MARK-PRESENTMENT

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `AL-DECLINED` | Y | N | N |
| `PS-AMOUNT > AL-AMOUNT * TOLERANCE` | | Y | N |
| → `SET PS-DECLINED TO TRUE` | 1 | | |
| → `SET PS-OVER-AUTH TO TRUE` | | 1 | |
| → `SET PS-MATCHED TO TRUE` | | | 1 |

Source: [match.cbl:104](../../../laurasia/programs/match.cbl#L104).

## WRITE-PRESENTMENT

| | 1 | 2 |
| --- | :-: | :-: |
| `PS-MATCHED` | Y | N |
| → `WRITE PRES-RECORD` | 1 | 1 |
| → `ADD 1 TO WS-MATCHED` | 2 | |
| → `ADD 1 TO WS-NOT-MATCHED` | | 2 |

Source: [match.cbl:114](../../../laurasia/programs/match.cbl#L114).
