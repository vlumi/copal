# ACCTLOAD: decision tables

The business rules of [ACCTLOAD](../programs/ACCTLOAD.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

## Contents

- [1000-INITIALIZE](#1000-initialize)
- [1200-READ-ACCOUNT](#1200-read-account)
- [2000-LOAD-ACCOUNT](#2000-load-account)
- [2100-WRITE-ONLINE](#2100-write-online)

## 1000-INITIALIZE

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `WS-FILE-STATUSES = ALL '0'` | N | Y | Y |
| `READ CONTROL-FILE AT END` | | Y | N |
| → `OPEN INPUT CONTROL-FILE ACCT-IN OUTPUT ONLINE-OUT REPORT-OUT` | 1 | 1 | 1 |
| → `DISPLAY 'ACCTLOAD: CANNOT OPEN FILES ' WS-FILE-STATUSES` | 2 | | |
| → `MOVE 16 TO RETURN-CODE` | 3 | | |
| → `STOP RUN` | 4 | | |
| → `READ CONTROL-FILE` | | 2 | 2 |
| → `MOVE 'NO BANK CONTROL RECORD' TO WS-ABORT-MESSAGE` | | 3 | |
| → `PERFORM 9900-ABORT` | | 4 | |
| → `MOVE BC-BUSINESS-DATE TO WS-BUSINESS-DATE` | | 5 | 3 |
| → `PERFORM 1200-READ-ACCOUNT` | | 6 | 4 |

Source: [ACCTLOAD.cbl:90](../../../pangaea/programs/ACCTLOAD.cbl#L90).

## 1200-READ-ACCOUNT

| | 1 | 2 |
| --- | :-: | :-: |
| `READ ACCT-IN AT END` | Y | N |
| → `READ ACCT-IN` | 1 | 1 |
| → `SET MASTER-AT-END TO TRUE` | 2 | |
| → `ADD 1 TO WS-READ-COUNT` | | 2 |

Source: [ACCTLOAD.cbl:106](../../../pangaea/programs/ACCTLOAD.cbl#L106).

## 2000-LOAD-ACCOUNT

| | 1 | 2 |
| --- | :-: | :-: |
| `AM-CLOSED` | Y | N |
| → `ADD 1 TO WS-CLOSED-COUNT` | 1 | |
| → `PERFORM 1200-READ-ACCOUNT` | 2 | 2 |
| → `PERFORM 2100-WRITE-ONLINE` | | 1 |

Source: [ACCTLOAD.cbl:115](../../../pangaea/programs/ACCTLOAD.cbl#L115).

## 2100-WRITE-ONLINE

| | 1 | 2 |
| --- | :-: | :-: |
| `WRITE ACCT-ONLINE-RECORD INVALID KEY` | Y | N |
| → `MOVE SPACES TO ACCT-ONLINE-RECORD` | 1 | 1 |
| → `MOVE AM-ACCT-NO TO AO-ACCT-NO` | 2 | 2 |
| → `MOVE AM-CUST-NO TO AO-CUST-NO` | 3 | 3 |
| → `MOVE AM-SHORT-NAME TO AO-SHORT-NAME` | 4 | 4 |
| → `MOVE AM-ACCT-TYPE TO AO-ACCT-TYPE` | 5 | 5 |
| → `MOVE AM-STATUS TO AO-STATUS` | 6 | 6 |
| → `MOVE AM-BALANCE TO AO-BALANCE` | 7 | 7 |
| → `MOVE AM-OVERDRAFT-LIMIT TO AO-OVERDRAFT-LIMIT` | 8 | 8 |
| → `MOVE 0 TO AO-HOLDS` | 9 | 9 |
| → `MOVE WS-BUSINESS-DATE TO AO-AS-OF-DATE` | 10 | 10 |
| → `WRITE ACCT-ONLINE-RECORD` | 11 | 11 |
| → `MOVE 'MASTER NOT IN ACCOUNT NUMBER ORDER' TO WS-ABORT-MESSAGE` | 12 | |
| → `PERFORM 9900-ABORT` | 13 | |
| → `ADD 1 TO WS-LOADED-COUNT` | 14 | 12 |
| → `ADD AM-BALANCE TO WS-LOADED-BALANCE` | 15 | 13 |

Source: [ACCTLOAD.cbl:123](../../../pangaea/programs/ACCTLOAD.cbl#L123).
