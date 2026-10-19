# POST: decision tables

The business rules of [POST](../programs/POST.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

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

Source: [post.cbl:131](../../../laurasia/programs/post.cbl#L131).
