# HOLDEXT: decision tables

The business rules of [HOLDEXT](../programs/HOLDEXT.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

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

- **Rule 1.** When `WS-FILE-STATUSES` is not `ALL '0'`, it opens `CONTROL-FILE` and `HOLD-FILE` for input and `TRAN-OUT` and `REPORT-OUT` for output; displays `'HOLDEXT: CANNOT OPEN FILES '` `WS-FILE-STATUSES`; sets `RETURN-CODE` to 16; and stops the run.
- **Rule 2.** When `WS-FILE-STATUSES` is `ALL '0'` and reading `CONTROL-FILE` finds no more records, it opens `CONTROL-FILE` and `HOLD-FILE` for input and `TRAN-OUT` and `REPORT-OUT` for output; reads `CONTROL-FILE`; sets `WS-ABORT-MESSAGE` to `'NO BANK CONTROL RECORD'`; and performs `9900-ABORT`.
- **Rule 3.** When `WS-FILE-STATUSES` is `ALL '0'` and reading `CONTROL-FILE` finds a record, it opens `CONTROL-FILE` and `HOLD-FILE` for input and `TRAN-OUT` and `REPORT-OUT` for output; reads `CONTROL-FILE`; sets `TRAN-RECORD` to spaces; sets `TR-REC-TYPE` to `TR-HEADER`; sets `TH-BUSINESS-DATE` to `BC-BUSINESS-DATE`; sets `TH-SOURCE` to `TH-FROM-CARDS`; sets `TH-FILE-SEQ` to 1; writes `TRAN-OUT-RECORD` from `TRAN-RECORD`; and performs `1200-READ-HOLD`.

Source: [HOLDEXT.cbl:103](../../../pangaea/programs/HOLDEXT.cbl#L103).

## 1200-READ-HOLD

| | 1 | 2 |
| --- | :-: | :-: |
| `READ HOLD-FILE AT END` | Y | N |
| → `READ HOLD-FILE` | 1 | 1 |
| → `SET END-OF-HOLDS TO TRUE` | 2 | |
| → `ADD 1 TO WS-HOLD-COUNT` | | 2 |

- **Rule 1.** When reading `HOLD-FILE` finds no more records, it reads `HOLD-FILE` and sets `WS-HOLDS-END-SW` to `END-OF-HOLDS`.
- **Rule 2.** When reading `HOLD-FILE` finds a record, it reads `HOLD-FILE` and adds 1 to `WS-HOLD-COUNT`.

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

- **Rule 1.** When `HD-BUSINESS-DATE` is not `BC-BUSINESS-DATE`, it sets `WS-ABORT-MESSAGE` to `'A HOLD OF ANOTHER BUSINESS DAY'` and performs `9900-ABORT`.
- **Rule 2.** When `HD-BUSINESS-DATE` is `BC-BUSINESS-DATE` and `HD-STATUS` is `'R'` (`HD-REVERSED`), it adds 1 to `WS-REVERSED-COUNT` and performs `1200-READ-HOLD`.
- **Rule 3.** When `HD-BUSINESS-DATE` is `BC-BUSINESS-DATE`, `HD-STATUS` is not `'R'` (not `HD-REVERSED`), `HD-KIND` is `'P'` (`HD-PURCHASE`), and `HD-FEE` is more than 0, it sets `TRAN-RECORD` to spaces, sets `TR-REC-TYPE` to `TR-DETAIL`, sets `WS-TI-HOLD-NO` to `HD-HOLD-NO`, sets `TD-ACCT-NO` to `HD-ACCT-NO`, sets `TD-DESCRIPTION` to `HD-TERMINAL`, sets `TD-COUNTER-ACCT-NO` to 0, sets `TD-TRAN-CODE` to `TD-CARD-PURCHASE`, sets `TD-AMOUNT` to `HD-AMOUNT`, sets `WS-TI-PART` to 1, performs `2900-WRITE-DETAIL`, sets `TD-TRAN-CODE` to `TD-ATM-FEE`, sets `TD-AMOUNT` to `HD-FEE`, sets `TD-DESCRIPTION` to `'ATM FEE'`, sets `WS-TI-PART` to 2, performs `2900-WRITE-DETAIL`, and performs `1200-READ-HOLD`.
- **Rule 4.** When `HD-BUSINESS-DATE` is `BC-BUSINESS-DATE`, `HD-STATUS` is not `'R'` (not `HD-REVERSED`), `HD-KIND` is `'P'` (`HD-PURCHASE`), and `HD-FEE` is at most 0, it sets `TRAN-RECORD` to spaces, sets `TR-REC-TYPE` to `TR-DETAIL`, sets `WS-TI-HOLD-NO` to `HD-HOLD-NO`, sets `TD-ACCT-NO` to `HD-ACCT-NO`, sets `TD-DESCRIPTION` to `HD-TERMINAL`, sets `TD-COUNTER-ACCT-NO` to 0, sets `TD-TRAN-CODE` to `TD-CARD-PURCHASE`, sets `TD-AMOUNT` to `HD-AMOUNT`, sets `WS-TI-PART` to 1, performs `2900-WRITE-DETAIL`, and performs `1200-READ-HOLD`.
- **Rule 5.** When `HD-BUSINESS-DATE` is `BC-BUSINESS-DATE`, `HD-STATUS` is not `'R'` (not `HD-REVERSED`), `HD-KIND` is not `'P'` (not `HD-PURCHASE`), and `HD-FEE` is more than 0, it sets `TRAN-RECORD` to spaces, sets `TR-REC-TYPE` to `TR-DETAIL`, sets `WS-TI-HOLD-NO` to `HD-HOLD-NO`, sets `TD-ACCT-NO` to `HD-ACCT-NO`, sets `TD-DESCRIPTION` to `HD-TERMINAL`, sets `TD-COUNTER-ACCT-NO` to 0, sets `TD-TRAN-CODE` to `TD-ATM-WITHDRAWAL`, sets `TD-AMOUNT` to `HD-AMOUNT`, sets `WS-TI-PART` to 1, performs `2900-WRITE-DETAIL`, sets `TD-TRAN-CODE` to `TD-ATM-FEE`, sets `TD-AMOUNT` to `HD-FEE`, sets `TD-DESCRIPTION` to `'ATM FEE'`, sets `WS-TI-PART` to 2, performs `2900-WRITE-DETAIL`, and performs `1200-READ-HOLD`.
- **Rule 6.** When `HD-BUSINESS-DATE` is `BC-BUSINESS-DATE`, `HD-STATUS` is not `'R'` (not `HD-REVERSED`), `HD-KIND` is not `'P'` (not `HD-PURCHASE`), and `HD-FEE` is at most 0, it sets `TRAN-RECORD` to spaces, sets `TR-REC-TYPE` to `TR-DETAIL`, sets `WS-TI-HOLD-NO` to `HD-HOLD-NO`, sets `TD-ACCT-NO` to `HD-ACCT-NO`, sets `TD-DESCRIPTION` to `HD-TERMINAL`, sets `TD-COUNTER-ACCT-NO` to 0, sets `TD-TRAN-CODE` to `TD-ATM-WITHDRAWAL`, sets `TD-AMOUNT` to `HD-AMOUNT`, sets `WS-TI-PART` to 1, performs `2900-WRITE-DETAIL`, and performs `1200-READ-HOLD`.

Source: [HOLDEXT.cbl:134](../../../pangaea/programs/HOLDEXT.cbl#L134).
