# PAYIN: decision tables

The business rules of [PAYIN](../programs/PAYIN.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

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

- **Rule 1.** When `WS-ACCOUNTS-STATUS` is not `'00'`, it accepts `WS-RUN-DATE` from its command line; opens `FEED-IN` for input, `ACCOUNTS` to update, and `POSTINGS-OUT`, `REJECT-OUT`, and `REPORT-OUT` for output; and stops the run.
- **Rule 2.** When `WS-ACCOUNTS-STATUS` is `'00'` and `WS-REJECTED` is more than zero, it accepts `WS-RUN-DATE` from its command line; opens `FEED-IN` for input, `ACCOUNTS` to update, and `POSTINGS-OUT`, `REJECT-OUT`, and `REPORT-OUT` for output; repeats its loop until `WS-AT-END` is `'Y'` (`FEED-AT-END`); performs `WRITE-REPORT`; closes `FEED-IN`, `ACCOUNTS`, `POSTINGS-OUT`, `REJECT-OUT`, and `REPORT-OUT`; and stops the run.
- **Rule 3.** When `WS-ACCOUNTS-STATUS` is `'00'` and `WS-REJECTED` is at most zero, it accepts `WS-RUN-DATE` from its command line; opens `FEED-IN` for input, `ACCOUNTS` to update, and `POSTINGS-OUT`, `REJECT-OUT`, and `REPORT-OUT` for output; repeats its loop until `WS-AT-END` is `'Y'` (`FEED-AT-END`); performs `WRITE-REPORT`; closes `FEED-IN`, `ACCOUNTS`, `POSTINGS-OUT`, `REJECT-OUT`, and `REPORT-OUT`; and stops the run.

Source: [payin.cbl:72](../../../laurasia/programs/payin.cbl#L72).

## MAIN, in its loop

Each pass of `PERFORM UNTIL FEED-AT-END`.

| | 1 | 2 |
| --- | :-: | :-: |
| `READ FEED-IN AT END` | Y | N |
| → `READ FEED-IN` | 1 | 1 |
| → `SET FEED-AT-END TO TRUE` | 2 | |
| → `PERFORM APPLY-RECORD` | | 2 |

- **Rule 1.** When reading `FEED-IN` finds no more records, it reads `FEED-IN` and sets `WS-AT-END` to `FEED-AT-END`.
- **Rule 2.** When reading `FEED-IN` finds a record, it reads `FEED-IN` and performs `APPLY-RECORD`.

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

- **Rule 1.** When `WS-REASON` is not spaces (not `RECORD-IS-GOOD`), it performs `VALIDATE-RECORD`, performs `REJECT-RECORD`, and leaves the section.
- **Rule 2.** When `WS-REASON` is spaces (`RECORD-IS-GOOD`) and reading `ACCOUNTS` fails on its key, it performs `VALIDATE-RECORD`, sets `CA-CARD-NO` to `PF-CARD-NO`, reads `ACCOUNTS`, sets `WS-REASON` to `'CARD'`, performs `REJECT-RECORD`, and leaves the section.
- **Rule 3.** When `WS-REASON` is spaces (`RECORD-IS-GOOD`), reading `ACCOUNTS` succeeds by its key, `PF-KIND` is `'P'` (`PF-PAYMENT`), and rewriting `CARD-ACCOUNT` fails on its key, it performs `VALIDATE-RECORD`, sets `CA-CARD-NO` to `PF-CARD-NO`, reads `ACCOUNTS`, performs `APPLY-PAYMENT`, rewrites `CARD-ACCOUNT`, and stops the run.
- **Rule 4.** When `WS-REASON` is spaces (`RECORD-IS-GOOD`), reading `ACCOUNTS` succeeds by its key, `PF-KIND` is `'P'` (`PF-PAYMENT`), and rewriting `CARD-ACCOUNT` succeeds by its key, it performs `VALIDATE-RECORD`, sets `CA-CARD-NO` to `PF-CARD-NO`, reads `ACCOUNTS`, performs `APPLY-PAYMENT`, and rewrites `CARD-ACCOUNT`.
- **Rule 5.** When `WS-REASON` is spaces (`RECORD-IS-GOOD`), reading `ACCOUNTS` succeeds by its key, `PF-KIND` is not `'P'` (not `PF-PAYMENT`), and rewriting `CARD-ACCOUNT` fails on its key, it performs `VALIDATE-RECORD`, sets `CA-CARD-NO` to `PF-CARD-NO`, reads `ACCOUNTS`, performs `APPLY-STATUS`, rewrites `CARD-ACCOUNT`, and stops the run.
- **Rule 6.** When `WS-REASON` is spaces (`RECORD-IS-GOOD`), reading `ACCOUNTS` succeeds by its key, `PF-KIND` is not `'P'` (not `PF-PAYMENT`), and rewriting `CARD-ACCOUNT` succeeds by its key, it performs `VALIDATE-RECORD`, sets `CA-CARD-NO` to `PF-CARD-NO`, reads `ACCOUNTS`, performs `APPLY-STATUS`, and rewrites `CARD-ACCOUNT`.

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

- **Rule 1.** When `PF-KIND` is not `'P'` (not `PF-PAYMENT`) and `PF-KIND` is not `'S'` (not `PF-STATUS-CHANGE`), it sets `WS-REASON` to `'KIND'`.
- **Rule 2.** When `PF-KIND` is `'P'` (`PF-PAYMENT`) or `PF-KIND` is `'S'` (`PF-STATUS-CHANGE`); and `PF-KIND` is `'P'` (`PF-PAYMENT`) and `PF-AMOUNT` is at most zero, it sets `WS-REASON` to `'AMNT'`.
- **Rule 3.** When `PF-KIND` is `'P'` (`PF-PAYMENT`) or `PF-KIND` is `'S'` (`PF-STATUS-CHANGE`); `PF-KIND` is not `'P'` (not `PF-PAYMENT`) or `PF-AMOUNT` is more than zero; and `PF-KIND` is `'S'` (`PF-STATUS-CHANGE`), `PF-STATUS` is not `'A'`, `PF-STATUS` is not `'B'`, and `PF-STATUS` is not `'C'`, it sets `WS-REASON` to `'STAT'`.
- **Rule 4.** When `PF-KIND` is `'P'` (`PF-PAYMENT`) or `PF-KIND` is `'S'` (`PF-STATUS-CHANGE`); `PF-KIND` is not `'P'` (not `PF-PAYMENT`) or `PF-AMOUNT` is more than zero; and `PF-KIND` is not `'S'` (not `PF-STATUS-CHANGE`), `PF-STATUS` is `'A'`, `PF-STATUS` is `'B'`, or `PF-STATUS` is `'C'`, it sets `WS-REASON` to spaces.

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

- **Rule 1.** When `CA-BALANCE` is at most zero or `PF-AMOUNT` is at least `CA-STATEMENT-BALANCE`, it subtracts `PF-AMOUNT` from `CA-BALANCE`, sets `CA-PAID-IN-FULL` to `CA-WAS-PAID-IN-FULL`, sets `PO-CARD-NO` to `CA-CARD-NO`, sets `PO-ACCOUNT-NO` to `CA-ACCOUNT-NO`, sets `PO-KIND` to `PO-PAYMENT`, computes `PO-AMOUNT` as `0 - PF-AMOUNT`, sets `PO-REFERENCE` to `PF-REFERENCE`, sets `PO-DATE` to `WS-RUN-DATE`, writes `POSTING`, adds 1 to `WS-PAYMENTS`, and adds `PF-AMOUNT` to `WS-PAID`.
- **Rule 2.** When `CA-BALANCE` is more than zero and `PF-AMOUNT` is less than `CA-STATEMENT-BALANCE`, it subtracts `PF-AMOUNT` from `CA-BALANCE`, sets `PO-CARD-NO` to `CA-CARD-NO`, sets `PO-ACCOUNT-NO` to `CA-ACCOUNT-NO`, sets `PO-KIND` to `PO-PAYMENT`, computes `PO-AMOUNT` as `0 - PF-AMOUNT`, sets `PO-REFERENCE` to `PF-REFERENCE`, sets `PO-DATE` to `WS-RUN-DATE`, writes `POSTING`, adds 1 to `WS-PAYMENTS`, and adds `PF-AMOUNT` to `WS-PAID`.

Source: [payin.cbl:129](../../../laurasia/programs/payin.cbl#L129).
