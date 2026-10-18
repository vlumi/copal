# HOLDEXT: decision tables

The business rules of [HOLDEXT](../programs/HOLDEXT.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

## Contents

- [1000-INITIALIZE](#1000-initialize)
- [1200-READ-HOLD](#1200-read-hold)
- [2000-EXTRACT-HOLD](#2000-extract-hold)

## 1000-INITIALIZE

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `WS-FILE-STATUSES = ALL '0'` | N | Y | Y |
| `READ CONTROL-FILE AT END` | | Y | N |
| → `OPEN INPUT CONTROL-FILE HOLD-FILE OUTPUT TRAN-OUT REPORT-OUT` | 1 | 1 | 1 |
| → `DISPLAY 'HOLDEXT: CANNOT OPEN FILES ' WS-FILE-STATUSES` | 2 | | |
| → `MOVE 16 TO RETURN-CODE` | 3 | | |
| → `STOP RUN` | 4 | | |
| → `READ CONTROL-FILE` | | 2 | 2 |
| → `MOVE 'NO BANK CONTROL RECORD' TO WS-ABORT-MESSAGE` | | 3 | |
| → `PERFORM 9900-ABORT` | | 4 | |
| → `MOVE SPACES TO TRAN-RECORD` | | | 3 |
| → `SET TR-HEADER TO TRUE` | | | 4 |
| → `MOVE BC-BUSINESS-DATE TO TH-BUSINESS-DATE` | | | 5 |
| → `SET TH-FROM-CARDS TO TRUE` | | | 6 |
| → `MOVE 1 TO TH-FILE-SEQ` | | | 7 |
| → `WRITE TRAN-OUT-RECORD FROM TRAN-RECORD` | | | 8 |
| → `PERFORM 1200-READ-HOLD` | | | 9 |

Source: [HOLDEXT.cbl:103](../../../pangaea/programs/HOLDEXT.cbl#L103).

## 1200-READ-HOLD

| | 1 | 2 |
| --- | :-: | :-: |
| `READ HOLD-FILE AT END` | Y | N |
| → `READ HOLD-FILE` | 1 | 1 |
| → `SET END-OF-HOLDS TO TRUE` | 2 | |
| → `ADD 1 TO WS-HOLD-COUNT` | | 2 |

Source: [HOLDEXT.cbl:124](../../../pangaea/programs/HOLDEXT.cbl#L124).

## 2000-EXTRACT-HOLD

| | 1 | 2 | 3 | 4 | 5 | 6 |
| --- | :-: | :-: | :-: | :-: | :-: | :-: |
| `HD-BUSINESS-DATE = BC-BUSINESS-DATE` | N | Y | Y | Y | Y | Y |
| `HD-REVERSED` | | Y | N | N | N | N |
| `HD-PURCHASE` | | | Y | Y | N | N |
| `HD-FEE > 0` | | | Y | N | Y | N |
| → `MOVE 'A HOLD OF ANOTHER BUSINESS DAY' TO WS-ABORT-MESSAGE` | 1 | | | | | |
| → `PERFORM 9900-ABORT` | 2 | | | | | |
| → `ADD 1 TO WS-REVERSED-COUNT` | | 1 | | | | |
| → `PERFORM 1200-READ-HOLD` | | 2 | 16 | 11 | 16 | 11 |
| → `MOVE SPACES TO TRAN-RECORD` | | | 1 | 1 | 1 | 1 |
| → `SET TR-DETAIL TO TRUE` | | | 2 | 2 | 2 | 2 |
| → `MOVE HD-HOLD-NO TO WS-TI-HOLD-NO` | | | 3 | 3 | 3 | 3 |
| → `MOVE HD-ACCT-NO TO TD-ACCT-NO` | | | 4 | 4 | 4 | 4 |
| → `MOVE HD-TERMINAL TO TD-DESCRIPTION` | | | 5 | 5 | 5 | 5 |
| → `MOVE 0 TO TD-COUNTER-ACCT-NO` | | | 6 | 6 | 6 | 6 |
| → `SET TD-CARD-PURCHASE TO TRUE` | | | 7 | 7 | | |
| → `MOVE HD-AMOUNT TO TD-AMOUNT` | | | 8 | 8 | 8 | 8 |
| → `MOVE 1 TO WS-TI-PART` | | | 9 | 9 | 9 | 9 |
| → `PERFORM 2900-WRITE-DETAIL` | | | 10, 15 | 10 | 10, 15 | 10 |
| → `SET TD-ATM-FEE TO TRUE` | | | 11 | | 11 | |
| → `MOVE HD-FEE TO TD-AMOUNT` | | | 12 | | 12 | |
| → `MOVE 'ATM FEE' TO TD-DESCRIPTION` | | | 13 | | 13 | |
| → `MOVE 2 TO WS-TI-PART` | | | 14 | | 14 | |
| → `SET TD-ATM-WITHDRAWAL TO TRUE` | | | | | 7 | 7 |

Source: [HOLDEXT.cbl:134](../../../pangaea/programs/HOLDEXT.cbl#L134).
