# PAYIN: decision tables

The business rules of [PAYIN](../programs/PAYIN.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

## Contents

- [MAIN](#main)
- [MAIN, in its loop](#main-in-its-loop)
- [APPLY-RECORD](#apply-record)
- [VALIDATE-RECORD](#validate-record)
- [APPLY-PAYMENT](#apply-payment)

## MAIN

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `WS-ACCOUNTS-STATUS = '00'` | N | Y | Y |
| `WS-REJECTED > ZERO` | | Y | N |
| → `ACCEPT WS-RUN-DATE FROM COMMAND-LINE` | 1 | 1 | 1 |
| → `OPEN INPUT FEED-IN I-O ACCOUNTS OUTPUT POSTINGS-OUT REJECT-OUT REPORT-OUT` | 2 | 2 | 2 |
| → `STOP RUN RETURNING 16` | 3 | | |
| → `PERFORM UNTIL FEED-AT-END` | | 3 | 3 |
| → `PERFORM WRITE-REPORT` | | 4 | 4 |
| → `CLOSE FEED-IN ACCOUNTS POSTINGS-OUT REJECT-OUT REPORT-OUT` | | 5 | 5 |
| → `STOP RUN RETURNING 4` | | 6 | |
| → `STOP RUN RETURNING 0` | | | 6 |

Source: [payin.cbl:72](../../../laurasia/programs/payin.cbl#L72).

## MAIN, in its loop

Each pass of `PERFORM UNTIL FEED-AT-END`.

| | 1 | 2 |
| --- | :-: | :-: |
| `READ FEED-IN AT END` | Y | N |
| → `READ FEED-IN` | 1 | 1 |
| → `SET FEED-AT-END TO TRUE` | 2 | |
| → `PERFORM APPLY-RECORD` | | 2 |

Source: [payin.cbl:80](../../../laurasia/programs/payin.cbl#L80).

## APPLY-RECORD

| | 1 | 2 | 3 | 4 | 5 | 6 |
| --- | :-: | :-: | :-: | :-: | :-: | :-: |
| `RECORD-IS-GOOD` | N | Y | Y | Y | Y | Y |
| `READ ACCOUNTS INVALID KEY` | | Y | N | N | N | N |
| `PF-PAYMENT` | | | Y | Y | N | N |
| `REWRITE CARD-ACCOUNT INVALID KEY` | | | Y | N | Y | N |
| → `PERFORM VALIDATE-RECORD` | 1 | 1 | 1 | 1 | 1 | 1 |
| → `PERFORM REJECT-RECORD` | 2 | 5 | | | | |
| → `EXIT SECTION` | 3 | 6 | | | | |
| → `MOVE PF-CARD-NO TO CA-CARD-NO` | | 2 | 2 | 2 | 2 | 2 |
| → `READ ACCOUNTS` | | 3 | 3 | 3 | 3 | 3 |
| → `MOVE 'CARD' TO WS-REASON` | | 4 | | | | |
| → `PERFORM APPLY-PAYMENT` | | | 4 | 4 | | |
| → `REWRITE CARD-ACCOUNT` | | | 5 | 5 | 5 | 5 |
| → `STOP RUN RETURNING 16` | | | 6 | | 6 | |
| → `PERFORM APPLY-STATUS` | | | | | 4 | 4 |

Source: [payin.cbl:93](../../../laurasia/programs/payin.cbl#L93).

## VALIDATE-RECORD

| | 1 | 2 | 3 | 4 |
| --- | :-: | :-: | :-: | :-: |
| `NOT PF-PAYMENT AND NOT PF-STATUS-CHANGE` | Y | N | N | N |
| `PF-PAYMENT AND PF-AMOUNT NOT > ZERO` | | Y | N | N |
| `PF-STATUS-CHANGE AND PF-STATUS NOT = 'A' AND PF-STATUS NOT = 'B' AND PF-STATUS NOT = 'C'` | | | Y | N |
| → `MOVE 'KIND' TO WS-REASON` | 1 | | | |
| → `MOVE 'AMNT' TO WS-REASON` | | 1 | | |
| → `MOVE 'STAT' TO WS-REASON` | | | 1 | |
| → `MOVE SPACES TO WS-REASON` | | | | 1 |

Source: [payin.cbl:115](../../../laurasia/programs/payin.cbl#L115).

## APPLY-PAYMENT

| | 1 | 2 |
| --- | :-: | :-: |
| `CA-BALANCE NOT > ZERO OR PF-AMOUNT NOT < CA-STATEMENT-BALANCE` | Y | N |
| → `SUBTRACT PF-AMOUNT FROM CA-BALANCE` | 1 | 1 |
| → `SET CA-WAS-PAID-IN-FULL TO TRUE` | 2 | |
| → `MOVE CA-CARD-NO TO PO-CARD-NO` | 3 | 2 |
| → `MOVE CA-ACCOUNT-NO TO PO-ACCOUNT-NO` | 4 | 3 |
| → `SET PO-PAYMENT TO TRUE` | 5 | 4 |
| → `COMPUTE PO-AMOUNT = 0 - PF-AMOUNT` | 6 | 5 |
| → `MOVE PF-REFERENCE TO PO-REFERENCE` | 7 | 6 |
| → `MOVE WS-RUN-DATE TO PO-DATE` | 8 | 7 |
| → `WRITE POSTING` | 9 | 8 |
| → `ADD 1 TO WS-PAYMENTS` | 10 | 9 |
| → `ADD PF-AMOUNT TO WS-PAID` | 11 | 10 |

Source: [payin.cbl:129](../../../laurasia/programs/payin.cbl#L129).
