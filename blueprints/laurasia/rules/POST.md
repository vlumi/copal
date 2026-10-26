# POST: decision tables

The business rules of [POST](../programs/POST.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

## Contents

- [MAIN](#main)
- [MAIN, in its loop](#main-in-its-loop)
- [POST-PRESENTMENT](#post-presentment)
- [CHARGE-OVER-LIMIT](#charge-over-limit)

## MAIN

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `WS-ACCOUNTS-STATUS = '00'` | N | Y | Y |
| `WS-REJECTED > ZERO` | | Y | N |
| → `ACCEPT WS-RUN-DATE FROM COMMAND-LINE` | 1 | 1 | 1 |
| → `OPEN INPUT PRES-IN I-O ACCOUNTS OUTPUT POSTINGS-OUT REJECT-OUT REPORT-OUT` | 2 | 2 | 2 |
| → `STOP RUN RETURNING 16` | 3 | | |
| → `PERFORM UNTIL PRESENTMENTS-AT-END` | | 3 | 3 |
| → `PERFORM WRITE-REPORT` | | 4 | 4 |
| → `CLOSE PRES-IN ACCOUNTS POSTINGS-OUT REJECT-OUT REPORT-OUT` | | 5 | 5 |
| → `STOP RUN RETURNING 4` | | 6 | |
| → `STOP RUN RETURNING 0` | | | 6 |

- **Rule 1.** When `WS-ACCOUNTS-STATUS` is not `'00'`, it accepts `WS-RUN-DATE` from its command line; opens `PRES-IN` for input, `ACCOUNTS` to update, and `POSTINGS-OUT`, `REJECT-OUT`, and `REPORT-OUT` for output; and stops the run.
- **Rule 2.** When `WS-ACCOUNTS-STATUS` is `'00'` and `WS-REJECTED` is more than zero, it accepts `WS-RUN-DATE` from its command line; opens `PRES-IN` for input, `ACCOUNTS` to update, and `POSTINGS-OUT`, `REJECT-OUT`, and `REPORT-OUT` for output; repeats its loop until `WS-AT-END` is `'Y'` (`PRESENTMENTS-AT-END`); performs `WRITE-REPORT`; closes `PRES-IN`, `ACCOUNTS`, `POSTINGS-OUT`, `REJECT-OUT`, and `REPORT-OUT`; and stops the run.
- **Rule 3.** When `WS-ACCOUNTS-STATUS` is `'00'` and `WS-REJECTED` is at most zero, it accepts `WS-RUN-DATE` from its command line; opens `PRES-IN` for input, `ACCOUNTS` to update, and `POSTINGS-OUT`, `REJECT-OUT`, and `REPORT-OUT` for output; repeats its loop until `WS-AT-END` is `'Y'` (`PRESENTMENTS-AT-END`); performs `WRITE-REPORT`; closes `PRES-IN`, `ACCOUNTS`, `POSTINGS-OUT`, `REJECT-OUT`, and `REPORT-OUT`; and stops the run.

Source: [post.cbl:75](../../../laurasia/programs/post.cbl#L75).

## MAIN, in its loop

Each pass of `PERFORM UNTIL PRESENTMENTS-AT-END`.

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `READ PRES-IN AT END` | Y | N | N |
| `PS-MATCHED` | | Y | N |
| → `READ PRES-IN` | 1 | 1 | 1 |
| → `SET PRESENTMENTS-AT-END TO TRUE` | 2 | | |
| → `PERFORM POST-PRESENTMENT` | | 2 | |

- **Rule 1.** When reading `PRES-IN` finds no more records, it reads `PRES-IN` and sets `WS-AT-END` to `PRESENTMENTS-AT-END`.
- **Rule 2.** When reading `PRES-IN` finds a record and `PS-MATCH` is `'M'` (`PS-MATCHED`), it reads `PRES-IN` and performs `POST-PRESENTMENT`.
- **Rule 3.** When reading `PRES-IN` finds a record and `PS-MATCH` is not `'M'` (not `PS-MATCHED`), it reads `PRES-IN`.

Source: [post.cbl:83](../../../laurasia/programs/post.cbl#L83).

## POST-PRESENTMENT

| | 1 | 2 | 3 | 4 | 5 | 6 | 7 |
| --- | :-: | :-: | :-: | :-: | :-: | :-: | :-: |
| `READ ACCOUNTS INVALID KEY` | Y | N | N | N | N | N | N |
| `CA-BLOCKED` | | Y | N | N | N | N | N |
| `CA-CLOSED` | | | Y | N | N | N | N |
| `CA-BALANCE > CA-LIMIT` | | | | Y | Y | N | N |
| `REWRITE CARD-ACCOUNT INVALID KEY` | | | | Y | N | Y | N |
| → `MOVE PS-CARD-NO TO CA-CARD-NO` | 1 | 1 | 1 | 1 | 1 | 1 | 1 |
| → `READ ACCOUNTS` | 2 | 2 | 2 | 2 | 2 | 2 | 2 |
| → `MOVE 'CARD' TO RJ-REASON` | 3 | | | | | | |
| → `PERFORM REJECT-PRESENTMENT` | 4 | 4 | 4 | | | | |
| → `EXIT SECTION` | 5 | 5 | 5 | | | | |
| → `MOVE 'BLKD' TO RJ-REASON` | | 3 | | | | | |
| → `MOVE 'CLSD' TO RJ-REASON` | | | 3 | | | | |
| → `ADD PS-AMOUNT TO CA-BALANCE` | | | | 3 | 3 | 3 | 3 |
| → `MOVE 'P' TO PO-KIND` | | | | 4 | 4 | 4 | 4 |
| → `MOVE PS-AMOUNT TO PO-AMOUNT` | | | | 5 | 5 | 5 | 5 |
| → `PERFORM WRITE-POSTING` | | | | 6 | 6 | 6 | 6 |
| → `ADD PS-AMOUNT TO WS-PURCHASES` | | | | 7 | 7 | 7 | 7 |
| → `PERFORM CHARGE-OVER-LIMIT` | | | | 8 | 8 | | |
| → `REWRITE CARD-ACCOUNT` | | | | 9 | 9 | 8 | 8 |
| → `STOP RUN RETURNING 16` | | | | 10 | | 9 | |

- **Rule 1.** When reading `ACCOUNTS` fails on its key, it sets `CA-CARD-NO` to `PS-CARD-NO`, reads `ACCOUNTS`, sets `RJ-REASON` to `'CARD'`, performs `REJECT-PRESENTMENT`, and leaves the section.
- **Rule 2.** When reading `ACCOUNTS` succeeds by its key and `CA-STATUS` is `'B'` (`CA-BLOCKED`), it sets `CA-CARD-NO` to `PS-CARD-NO`, reads `ACCOUNTS`, sets `RJ-REASON` to `'BLKD'`, performs `REJECT-PRESENTMENT`, and leaves the section.
- **Rule 3.** When reading `ACCOUNTS` succeeds by its key, `CA-STATUS` is not `'B'` (not `CA-BLOCKED`), and `CA-STATUS` is `'C'` (`CA-CLOSED`), it sets `CA-CARD-NO` to `PS-CARD-NO`, reads `ACCOUNTS`, sets `RJ-REASON` to `'CLSD'`, performs `REJECT-PRESENTMENT`, and leaves the section.
- **Rule 4.** When reading `ACCOUNTS` succeeds by its key, `CA-STATUS` is not `'B'` (not `CA-BLOCKED`), `CA-STATUS` is not `'C'` (not `CA-CLOSED`), `CA-BALANCE` is more than `CA-LIMIT`, and rewriting `CARD-ACCOUNT` fails on its key, it sets `CA-CARD-NO` to `PS-CARD-NO`, reads `ACCOUNTS`, adds `PS-AMOUNT` to `CA-BALANCE`, sets `PO-KIND` to `'P'`, sets `PO-AMOUNT` to `PS-AMOUNT`, performs `WRITE-POSTING`, adds `PS-AMOUNT` to `WS-PURCHASES`, performs `CHARGE-OVER-LIMIT`, rewrites `CARD-ACCOUNT`, and stops the run.
- **Rule 5.** When reading `ACCOUNTS` succeeds by its key, `CA-STATUS` is not `'B'` (not `CA-BLOCKED`), `CA-STATUS` is not `'C'` (not `CA-CLOSED`), `CA-BALANCE` is more than `CA-LIMIT`, and rewriting `CARD-ACCOUNT` succeeds by its key, it sets `CA-CARD-NO` to `PS-CARD-NO`, reads `ACCOUNTS`, adds `PS-AMOUNT` to `CA-BALANCE`, sets `PO-KIND` to `'P'`, sets `PO-AMOUNT` to `PS-AMOUNT`, performs `WRITE-POSTING`, adds `PS-AMOUNT` to `WS-PURCHASES`, performs `CHARGE-OVER-LIMIT`, and rewrites `CARD-ACCOUNT`.
- **Rule 6.** When reading `ACCOUNTS` succeeds by its key, `CA-STATUS` is not `'B'` (not `CA-BLOCKED`), `CA-STATUS` is not `'C'` (not `CA-CLOSED`), `CA-BALANCE` is at most `CA-LIMIT`, and rewriting `CARD-ACCOUNT` fails on its key, it sets `CA-CARD-NO` to `PS-CARD-NO`, reads `ACCOUNTS`, adds `PS-AMOUNT` to `CA-BALANCE`, sets `PO-KIND` to `'P'`, sets `PO-AMOUNT` to `PS-AMOUNT`, performs `WRITE-POSTING`, adds `PS-AMOUNT` to `WS-PURCHASES`, rewrites `CARD-ACCOUNT`, and stops the run.
- **Rule 7.** When reading `ACCOUNTS` succeeds by its key, `CA-STATUS` is not `'B'` (not `CA-BLOCKED`), `CA-STATUS` is not `'C'` (not `CA-CLOSED`), `CA-BALANCE` is at most `CA-LIMIT`, and rewriting `CARD-ACCOUNT` succeeds by its key, it sets `CA-CARD-NO` to `PS-CARD-NO`, reads `ACCOUNTS`, adds `PS-AMOUNT` to `CA-BALANCE`, sets `PO-KIND` to `'P'`, sets `PO-AMOUNT` to `PS-AMOUNT`, performs `WRITE-POSTING`, adds `PS-AMOUNT` to `WS-PURCHASES`, and rewrites `CARD-ACCOUNT`.

Source: [post.cbl:99](../../../laurasia/programs/post.cbl#L99).

## CHARGE-OVER-LIMIT

| | 1 | 2 | 3 | 4 |
| --- | :-: | :-: | :-: | :-: |
| `CA-BUSINESS` | Y | Y | N | N |
| `WS-FEE > ZERO` | Y | N | Y | N |
| → `MOVE 'ovlbiz' TO WS-FEE-PROGRAM` | 1 | 1 | | |
| → `COMPUTE WS-OVER-AMOUNT = CA-BALANCE - CA-LIMIT` | 2 | 2 | 2 | 2 |
| → `CALL WS-FEE-PROGRAM USING CA-PRODUCT WS-OVER-AMOUNT WS-FEE` | 3 | 3 | 3 | 3 |
| → `ADD WS-FEE TO CA-BALANCE WS-FEES` | 4 | | 4 | |
| → `MOVE 'F' TO PO-KIND` | 5 | | 5 | |
| → `MOVE WS-FEE TO PO-AMOUNT` | 6 | | 6 | |
| → `PERFORM WRITE-POSTING` | 7 | | 7 | |
| → `MOVE 'ovlfee' TO WS-FEE-PROGRAM` | | | 1 | 1 |

- **Rule 1.** When `CA-PRODUCT` is `'BZ'` (`CA-BUSINESS`) and `WS-FEE` is more than zero, it sets `WS-FEE-PROGRAM` to `'ovlbiz'`; computes `WS-OVER-AMOUNT` as `CA-BALANCE - CA-LIMIT`; calls the program `WS-FEE-PROGRAM` names; adds `WS-FEE` to `CA-BALANCE` and `WS-FEES`; sets `PO-KIND` to `'F'`; sets `PO-AMOUNT` to `WS-FEE`; and performs `WRITE-POSTING`.
- **Rule 2.** When `CA-PRODUCT` is `'BZ'` (`CA-BUSINESS`) and `WS-FEE` is at most zero, it sets `WS-FEE-PROGRAM` to `'ovlbiz'`, computes `WS-OVER-AMOUNT` as `CA-BALANCE - CA-LIMIT`, and calls the program `WS-FEE-PROGRAM` names.
- **Rule 3.** When `CA-PRODUCT` is not `'BZ'` (not `CA-BUSINESS`) and `WS-FEE` is more than zero, it sets `WS-FEE-PROGRAM` to `'ovlfee'`; computes `WS-OVER-AMOUNT` as `CA-BALANCE - CA-LIMIT`; calls the program `WS-FEE-PROGRAM` names; adds `WS-FEE` to `CA-BALANCE` and `WS-FEES`; sets `PO-KIND` to `'F'`; sets `PO-AMOUNT` to `WS-FEE`; and performs `WRITE-POSTING`.
- **Rule 4.** When `CA-PRODUCT` is not `'BZ'` (not `CA-BUSINESS`) and `WS-FEE` is at most zero, it sets `WS-FEE-PROGRAM` to `'ovlfee'`, computes `WS-OVER-AMOUNT` as `CA-BALANCE - CA-LIMIT`, and calls the program `WS-FEE-PROGRAM` names.

Source: [post.cbl:131](../../../laurasia/programs/post.cbl#L131).
