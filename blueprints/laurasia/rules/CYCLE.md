# CYCLE: decision tables

The business rules of [CYCLE](../programs/CYCLE.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

## Contents

- [MAIN](#main)
- [MAIN, in its loop](#main-in-its-loop)
- [CLOSE-CYCLE](#close-cycle)

## MAIN

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-ACCOUNTS-STATUS = '00'` | N | Y |
| → `ACCEPT WS-RUN-DATE FROM COMMAND-LINE` | 1 | 1 |
| → `OPEN I-O ACCOUNTS OUTPUT POSTINGS-OUT STATEMENTS-OUT` | 2 | 2 |
| → `STOP RUN RETURNING 16` | 3 | |
| → `READ ACCOUNTS NEXT RECORD` | | 3 |
| → `PERFORM UNTIL ACCOUNTS-AT-END` | | 4 |
| → `CLOSE ACCOUNTS POSTINGS-OUT STATEMENTS-OUT` | | 5 |
| → `STOP RUN RETURNING 0` | | 6 |

- **Rule 1.** When `WS-ACCOUNTS-STATUS` is not `'00'`, it accepts `WS-RUN-DATE` from its command line; opens `ACCOUNTS` to update and `POSTINGS-OUT` and `STATEMENTS-OUT` for output; and stops the run.
- **Rule 2.** When `WS-ACCOUNTS-STATUS` is `'00'`, it accepts `WS-RUN-DATE` from its command line; opens `ACCOUNTS` to update and `POSTINGS-OUT` and `STATEMENTS-OUT` for output; reads `ACCOUNTS`; repeats its loop until `WS-ACCOUNTS-STATUS` is `'10'` (`ACCOUNTS-AT-END`); closes `ACCOUNTS`, `POSTINGS-OUT`, and `STATEMENTS-OUT`; and stops the run.

Source: [cycle.cbl:63](../../../laurasia/programs/cycle.cbl#L63).

## MAIN, in its loop

Each pass of `PERFORM UNTIL ACCOUNTS-AT-END`.

| | 1 | 2 |
| --- | :-: | :-: |
| `CA-CYCLE-DAY = WS-RUN-DAY AND NOT CA-CLOSED` | Y | N |
| → `PERFORM CLOSE-CYCLE` | 1 | |
| → `READ ACCOUNTS NEXT RECORD` | 2 | 1 |

- **Rule 1.** When `CA-CYCLE-DAY` is `WS-RUN-DAY` and `CA-STATUS` is not `'C'` (not `CA-CLOSED`), it performs `CLOSE-CYCLE` and reads `ACCOUNTS`.
- **Rule 2.** When `CA-CYCLE-DAY` is not `WS-RUN-DAY` or `CA-STATUS` is `'C'` (`CA-CLOSED`), it reads `ACCOUNTS`.

Source: [cycle.cbl:71](../../../laurasia/programs/cycle.cbl#L71).

## CLOSE-CYCLE

| | 1 | 2 | 3 | 4 | 5 | 6 |
| --- | :-: | :-: | :-: | :-: | :-: | :-: |
| `NOT CA-WAS-PAID-IN-FULL AND CA-STATEMENT-BALANCE > ZERO` | Y | Y | Y | Y | N | N |
| `WS-INTEREST > ZERO` | Y | Y | N | N | N | N |
| `REWRITE CARD-ACCOUNT INVALID KEY` | Y | N | Y | N | Y | N |
| → `MOVE ZERO TO WS-INTEREST` | 1 | 1 | 1 | 1 | 1 | 1 |
| → `CALL 'intcalc' USING CA-PRODUCT CA-STATEMENT-BALANCE RETURNING WS-INTEREST-CENTS` | 2 | 2 | 2 | 2 | | |
| → `COMPUTE WS-INTEREST = WS-INTEREST-CENTS / 100` | 3 | 3 | 3 | 3 | | |
| → `ADD WS-INTEREST TO CA-BALANCE` | 4 | 4 | | | | |
| → `PERFORM POST-INTEREST` | 5 | 5 | | | | |
| → `CALL 'minpay' USING CA-BALANCE WS-MINIMUM` | 6 | 6 | 4 | 4 | 2 | 2 |
| → `PERFORM WRITE-STATEMENT` | 7 | 7 | 5 | 5 | 3 | 3 |
| → `MOVE CA-BALANCE TO CA-STATEMENT-BALANCE` | 8 | 8 | 6 | 6 | 4 | 4 |
| → `MOVE 'N' TO CA-PAID-IN-FULL` | 9 | 9 | 7 | 7 | 5 | 5 |
| → `REWRITE CARD-ACCOUNT` | 10 | 10 | 8 | 8 | 6 | 6 |
| → `STOP RUN RETURNING 16` | 11 | | 9 | | 7 | |
| → `ADD 1 TO WS-CLOSED-COUNT` | | 11 | | 9 | | 7 |

- **Rule 1.** When `CA-PAID-IN-FULL` is not `'Y'` (not `CA-WAS-PAID-IN-FULL`) and `CA-STATEMENT-BALANCE` is more than zero; `WS-INTEREST` is more than zero; and rewriting `CARD-ACCOUNT` fails on its key, it sets `WS-INTEREST` to zero, calls `INTCALC`, computes `WS-INTEREST` as `WS-INTEREST-CENTS / 100`, adds `WS-INTEREST` to `CA-BALANCE`, performs `POST-INTEREST`, calls `MINPAY`, performs `WRITE-STATEMENT`, sets `CA-STATEMENT-BALANCE` to `CA-BALANCE`, sets `CA-PAID-IN-FULL` to `'N'`, rewrites `CARD-ACCOUNT`, and stops the run.
- **Rule 2.** When `CA-PAID-IN-FULL` is not `'Y'` (not `CA-WAS-PAID-IN-FULL`) and `CA-STATEMENT-BALANCE` is more than zero; `WS-INTEREST` is more than zero; and rewriting `CARD-ACCOUNT` succeeds by its key, it sets `WS-INTEREST` to zero, calls `INTCALC`, computes `WS-INTEREST` as `WS-INTEREST-CENTS / 100`, adds `WS-INTEREST` to `CA-BALANCE`, performs `POST-INTEREST`, calls `MINPAY`, performs `WRITE-STATEMENT`, sets `CA-STATEMENT-BALANCE` to `CA-BALANCE`, sets `CA-PAID-IN-FULL` to `'N'`, rewrites `CARD-ACCOUNT`, and adds 1 to `WS-CLOSED-COUNT`.
- **Rule 3.** When `CA-PAID-IN-FULL` is not `'Y'` (not `CA-WAS-PAID-IN-FULL`) and `CA-STATEMENT-BALANCE` is more than zero; `WS-INTEREST` is at most zero; and rewriting `CARD-ACCOUNT` fails on its key, it sets `WS-INTEREST` to zero, calls `INTCALC`, computes `WS-INTEREST` as `WS-INTEREST-CENTS / 100`, calls `MINPAY`, performs `WRITE-STATEMENT`, sets `CA-STATEMENT-BALANCE` to `CA-BALANCE`, sets `CA-PAID-IN-FULL` to `'N'`, rewrites `CARD-ACCOUNT`, and stops the run.
- **Rule 4.** When `CA-PAID-IN-FULL` is not `'Y'` (not `CA-WAS-PAID-IN-FULL`) and `CA-STATEMENT-BALANCE` is more than zero; `WS-INTEREST` is at most zero; and rewriting `CARD-ACCOUNT` succeeds by its key, it sets `WS-INTEREST` to zero, calls `INTCALC`, computes `WS-INTEREST` as `WS-INTEREST-CENTS / 100`, calls `MINPAY`, performs `WRITE-STATEMENT`, sets `CA-STATEMENT-BALANCE` to `CA-BALANCE`, sets `CA-PAID-IN-FULL` to `'N'`, rewrites `CARD-ACCOUNT`, and adds 1 to `WS-CLOSED-COUNT`.
- **Rule 5.** When `CA-PAID-IN-FULL` is `'Y'` (`CA-WAS-PAID-IN-FULL`) or `CA-STATEMENT-BALANCE` is at most zero; `WS-INTEREST` is at most zero; and rewriting `CARD-ACCOUNT` fails on its key, it sets `WS-INTEREST` to zero, calls `MINPAY`, performs `WRITE-STATEMENT`, sets `CA-STATEMENT-BALANCE` to `CA-BALANCE`, sets `CA-PAID-IN-FULL` to `'N'`, rewrites `CARD-ACCOUNT`, and stops the run.
- **Rule 6.** When `CA-PAID-IN-FULL` is `'Y'` (`CA-WAS-PAID-IN-FULL`) or `CA-STATEMENT-BALANCE` is at most zero; `WS-INTEREST` is at most zero; and rewriting `CARD-ACCOUNT` succeeds by its key, it sets `WS-INTEREST` to zero, calls `MINPAY`, performs `WRITE-STATEMENT`, sets `CA-STATEMENT-BALANCE` to `CA-BALANCE`, sets `CA-PAID-IN-FULL` to `'N'`, rewrites `CARD-ACCOUNT`, and adds 1 to `WS-CLOSED-COUNT`.

Source: [cycle.cbl:80](../../../laurasia/programs/cycle.cbl#L80).
