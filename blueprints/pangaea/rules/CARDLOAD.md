# CARDLOAD: decision tables

The business rules of [CARDLOAD](../programs/CARDLOAD.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

## Contents

- [1000-INITIALIZE](#1000-initialize)
- [1200-READ-CARD](#1200-read-card)
- [2000-LOAD-CARD](#2000-load-card)
- [2100-WRITE-ONLINE](#2100-write-online)

## 1000-INITIALIZE

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `WS-FILE-STATUSES = ALL '0'` | N | Y | Y |
| `READ CONTROL-FILE AT END` | | Y | N |
| → `OPEN INPUT CONTROL-FILE CARD-IN OUTPUT ONLINE-OUT REPORT-OUT` | 1 | 1 | 1 |
| → `DISPLAY 'CARDLOAD: CANNOT OPEN FILES ' WS-FILE-STATUSES` | 2 | | |
| → `MOVE 16 TO RETURN-CODE` | 3 | | |
| → `STOP RUN` | 4 | | |
| → `READ CONTROL-FILE` | | 2 | 2 |
| → `MOVE 'NO BANK CONTROL RECORD' TO WS-ABORT-MESSAGE` | | 3 | |
| → `PERFORM 9900-ABORT` | | 4 | |
| → `MOVE BC-BUSINESS-DATE TO WS-BUSINESS-DATE` | | 5 | 3 |
| → `PERFORM 1200-READ-CARD` | | 6 | 4 |

Source: [CARDLOAD.cbl:92](../../../pangaea/programs/CARDLOAD.cbl#L92).

## 1200-READ-CARD

| | 1 | 2 |
| --- | :-: | :-: |
| `READ CARD-IN AT END` | Y | N |
| → `READ CARD-IN` | 1 | 1 |
| → `SET MASTER-AT-END TO TRUE` | 2 | |
| → `ADD 1 TO WS-READ-COUNT` | | 2 |

Source: [CARDLOAD.cbl:108](../../../pangaea/programs/CARDLOAD.cbl#L108).

## 2000-LOAD-CARD

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `CR-EXPIRY-MONTH < WS-BUSINESS-MONTH` | Y | N | N |
| `CR-BLOCKED` | | Y | N |
| → `ADD 1 TO WS-EXPIRED-COUNT` | 1 | | |
| → `PERFORM 1200-READ-CARD` | 2 | 3 | 2 |
| → `ADD 1 TO WS-BLOCKED-COUNT` | | 1 | |
| → `PERFORM 2100-WRITE-ONLINE` | | 2 | 1 |

Source: [CARDLOAD.cbl:117](../../../pangaea/programs/CARDLOAD.cbl#L117).

## 2100-WRITE-ONLINE

| | 1 | 2 |
| --- | :-: | :-: |
| `WRITE ONLINE-CARD FROM CARD-RECORD INVALID KEY` | Y | N |
| → `WRITE ONLINE-CARD FROM CARD-RECORD` | 1 | 1 |
| → `MOVE 'MASTER NOT IN CARD NUMBER ORDER' TO WS-ABORT-MESSAGE` | 2 | |
| → `PERFORM 9900-ABORT` | 3 | |
| → `ADD 1 TO WS-LOADED-COUNT` | 4 | 2 |

Source: [CARDLOAD.cbl:128](../../../pangaea/programs/CARDLOAD.cbl#L128).
