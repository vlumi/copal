# CARDLOAD: decision tables

The business rules of [CARDLOAD](../programs/CARDLOAD.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

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
| → `MOVE BC-BUSINESS-DATE TO WS-BUSINESS-DATE` | | | 3 |
| → `PERFORM 1200-READ-CARD` | | | 4 |

- **Rule 1.** When `WS-FILE-STATUSES` is not `ALL '0'`, it opens `CONTROL-FILE` and `CARD-IN` for input and `ONLINE-OUT` and `REPORT-OUT` for output; displays `'CARDLOAD: CANNOT OPEN FILES '` `WS-FILE-STATUSES`; sets `RETURN-CODE` to 16; and stops the run.
- **Rule 2.** When `WS-FILE-STATUSES` is `ALL '0'` and reading `CONTROL-FILE` finds no more records, it opens `CONTROL-FILE` and `CARD-IN` for input and `ONLINE-OUT` and `REPORT-OUT` for output; reads `CONTROL-FILE`; sets `WS-ABORT-MESSAGE` to `'NO BANK CONTROL RECORD'`; and performs `9900-ABORT`.
- **Rule 3.** When `WS-FILE-STATUSES` is `ALL '0'` and reading `CONTROL-FILE` finds a record, it opens `CONTROL-FILE` and `CARD-IN` for input and `ONLINE-OUT` and `REPORT-OUT` for output; reads `CONTROL-FILE`; sets `WS-BUSINESS-DATE` to `BC-BUSINESS-DATE`; and performs `1200-READ-CARD`.

Source: [CARDLOAD.cbl:92](../../../pangaea/programs/CARDLOAD.cbl#L92).

## 1200-READ-CARD

| | 1 | 2 |
| --- | :-: | :-: |
| `READ CARD-IN AT END` | Y | N |
| → `READ CARD-IN` | 1 | 1 |
| → `SET MASTER-AT-END TO TRUE` | 2 | |
| → `ADD 1 TO WS-READ-COUNT` | | 2 |

- **Rule 1.** When reading `CARD-IN` finds no more records, it reads `CARD-IN` and sets `WS-MASTER-END-SW` to `MASTER-AT-END`.
- **Rule 2.** When reading `CARD-IN` finds a record, it reads `CARD-IN` and adds 1 to `WS-READ-COUNT`.

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

- **Rule 1.** When `CR-EXPIRY-MONTH` is less than `WS-BUSINESS-MONTH`, it adds 1 to `WS-EXPIRED-COUNT` and performs `1200-READ-CARD`.
- **Rule 2.** When `CR-EXPIRY-MONTH` is at least `WS-BUSINESS-MONTH` and `CR-STATUS` is `'B'` (`CR-BLOCKED`), it adds 1 to `WS-BLOCKED-COUNT`, performs `2100-WRITE-ONLINE`, and performs `1200-READ-CARD`.
- **Rule 3.** When `CR-EXPIRY-MONTH` is at least `WS-BUSINESS-MONTH` and `CR-STATUS` is not `'B'` (not `CR-BLOCKED`), it performs `2100-WRITE-ONLINE` and performs `1200-READ-CARD`.

Source: [CARDLOAD.cbl:117](../../../pangaea/programs/CARDLOAD.cbl#L117).

## 2100-WRITE-ONLINE

| | 1 | 2 |
| --- | :-: | :-: |
| `WRITE ONLINE-CARD FROM CARD-RECORD INVALID KEY` | Y | N |
| → `WRITE ONLINE-CARD FROM CARD-RECORD` | 1 | 1 |
| → `MOVE 'MASTER NOT IN CARD NUMBER ORDER' TO WS-ABORT-MESSAGE` | 2 | |
| → `PERFORM 9900-ABORT` | 3 | |
| → `ADD 1 TO WS-LOADED-COUNT` | | 2 |

- **Rule 1.** When writing `ONLINE-CARD` fails on its key, it writes `ONLINE-CARD` from `CARD-RECORD`, sets `WS-ABORT-MESSAGE` to `'MASTER NOT IN CARD NUMBER ORDER'`, and performs `9900-ABORT`.
- **Rule 2.** When writing `ONLINE-CARD` succeeds by its key, it writes `ONLINE-CARD` from `CARD-RECORD` and adds 1 to `WS-LOADED-COUNT`.

Source: [CARDLOAD.cbl:128](../../../pangaea/programs/CARDLOAD.cbl#L128).
