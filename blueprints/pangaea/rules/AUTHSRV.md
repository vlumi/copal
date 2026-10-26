# AUTHSRV: decision tables

The business rules of [AUTHSRV](../programs/AUTHSRV.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

## Contents

- [1000-INITIALIZE](#1000-initialize)
- [1100-FIND-LAST-HOLD](#1100-find-last-hold)
- [1100-FIND-LAST-HOLD, in its loop](#1100-find-last-hold-in-its-loop)
- [2000-SERVE-REQUEST](#2000-serve-request)
- [2100-AUTHORIZE](#2100-authorize)
- [2200-READ-CARD-AND-ACCOUNT](#2200-read-card-and-account)
- [2300-READ-HOLD](#2300-read-hold)
- [2400-SAVE-CHANGES](#2400-save-changes)

## 1000-INITIALIZE

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `WS-FILE-STATUSES = ALL '0'` | N | Y | Y |
| `READ CONTROL-FILE AT END` | | Y | N |
| → `OPEN INPUT CONTROL-FILE I-O ACCT-FILE CARD-FILE HOLD-FILE OUTPUT REPORT-OUT` | 1 | 1 | 1 |
| → `DISPLAY 'AUTHSRV: CANNOT OPEN FILES ' WS-FILE-STATUSES UPON SYSERR` | 2 | | |
| → `MOVE 16 TO RETURN-CODE` | 3 | | |
| → `STOP RUN` | 4 | | |
| → `READ CONTROL-FILE` | | 2 | 2 |
| → `MOVE 'NO BANK CONTROL RECORD' TO WS-ABORT-MESSAGE` | | 3 | |
| → `PERFORM 9900-ABORT` | | 4 | |
| → `PERFORM 1100-FIND-LAST-HOLD` | | | 3 |

- **Rule 1.** When `WS-FILE-STATUSES` is not `ALL '0'`, it opens `CONTROL-FILE` for input, `ACCT-FILE`, `CARD-FILE`, and `HOLD-FILE` to update, and `REPORT-OUT` for output; displays `'AUTHSRV: CANNOT OPEN FILES '` `WS-FILE-STATUSES`; sets `RETURN-CODE` to 16; and stops the run.
- **Rule 2.** When `WS-FILE-STATUSES` is `ALL '0'` and reading `CONTROL-FILE` finds no more records, it opens `CONTROL-FILE` for input, `ACCT-FILE`, `CARD-FILE`, and `HOLD-FILE` to update, and `REPORT-OUT` for output; reads `CONTROL-FILE`; sets `WS-ABORT-MESSAGE` to `'NO BANK CONTROL RECORD'`; and performs `9900-ABORT`.
- **Rule 3.** When `WS-FILE-STATUSES` is `ALL '0'` and reading `CONTROL-FILE` finds a record, it opens `CONTROL-FILE` for input, `ACCT-FILE`, `CARD-FILE`, and `HOLD-FILE` to update, and `REPORT-OUT` for output; reads `CONTROL-FILE`; and performs `1100-FIND-LAST-HOLD`.

Source: [AUTHSRV.cbl:113](../../../pangaea/programs/AUTHSRV.cbl#L113).

## 1100-FIND-LAST-HOLD

| | 1 | 2 |
| --- | :-: | :-: |
| `START HOLD-FILE KEY IS NOT LESS THAN HD-KEY INVALID KEY` | Y | N |
| → `MOVE LOW-VALUES TO HD-KEY` | 1 | 1 |
| → `START HOLD-FILE KEY IS NOT LESS THAN HD-KEY` | 2 | 2 |
| → `SET END-OF-HOLDS TO TRUE` | 3 | |
| → `PERFORM UNTIL END-OF-HOLDS` | 4 | 3 |

- **Rule 1.** When positioning `HOLD-FILE` fails on its key, it sets `HD-KEY` to low values, runs `START HOLD-FILE KEY IS NOT LESS THAN HD-KEY`, sets `WS-HOLDS-END-SW` to `END-OF-HOLDS`, and repeats its loop until `WS-HOLDS-END-SW` is `'Y'` (`END-OF-HOLDS`).
- **Rule 2.** When positioning `HOLD-FILE` succeeds by its key, it sets `HD-KEY` to low values, runs `START HOLD-FILE KEY IS NOT LESS THAN HD-KEY`, and repeats its loop until `WS-HOLDS-END-SW` is `'Y'` (`END-OF-HOLDS`).

Source: [AUTHSRV.cbl:131](../../../pangaea/programs/AUTHSRV.cbl#L131).

## 1100-FIND-LAST-HOLD, in its loop

Each pass of `PERFORM UNTIL END-OF-HOLDS`.

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `READ HOLD-FILE NEXT RECORD AT END` | Y | N | N |
| `HD-HOLD-NO > WS-LAST-HOLD-NO` | | Y | N |
| → `READ HOLD-FILE NEXT RECORD` | 1 | 1 | 1 |
| → `SET END-OF-HOLDS TO TRUE` | 2 | | |
| → `MOVE HD-HOLD-NO TO WS-LAST-HOLD-NO` | | 2 | |

- **Rule 1.** When reading `HOLD-FILE` finds no more records, it reads `HOLD-FILE` and sets `WS-HOLDS-END-SW` to `END-OF-HOLDS`.
- **Rule 2.** When reading `HOLD-FILE` finds a record and `HD-HOLD-NO` is more than `WS-LAST-HOLD-NO`, it reads `HOLD-FILE` and sets `WS-LAST-HOLD-NO` to `HD-HOLD-NO`.
- **Rule 3.** When reading `HOLD-FILE` finds a record and `HD-HOLD-NO` is at most `WS-LAST-HOLD-NO`, it reads `HOLD-FILE`.

Source: [AUTHSRV.cbl:137](../../../pangaea/programs/AUTHSRV.cbl#L137).

## 2000-SERVE-REQUEST

| | 1 | 2 |
| --- | :-: | :-: |
| `ACCEPT AUTH-REQUEST ON EXCEPTION` | Y | N |
| → `ACCEPT AUTH-REQUEST` | 1 | 1 |
| → `SET END-OF-INPUT TO TRUE` | 2 | |
| → `PERFORM 2100-AUTHORIZE` | | 2 |
| → `DISPLAY AUTH-RESPONSE` | | 3 |

- **Rule 1.** When `ACCEPT AUTH-REQUEST` raises on exception, it accepts `AUTH-REQUEST` from its input and sets `WS-INPUT-END-SW` to `END-OF-INPUT`.
- **Rule 2.** When `ACCEPT AUTH-REQUEST` does not raise on exception, it accepts `AUTH-REQUEST` from its input, performs `2100-AUTHORIZE`, and displays `AUTH-RESPONSE`.

Source: [AUTHSRV.cbl:149](../../../pangaea/programs/AUTHSRV.cbl#L149).

## 2100-AUTHORIZE

| | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 |
| --- | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: |
| `RQ-PURCHASE OR RQ-REVERSAL OR RQ-WITHDRAWAL OR RQ-BALANCE-ENQUIRY` | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | N |
| `RS-RESULT = SPACES AND RQ-REVERSAL` | Y | Y | Y | N | N | N | N | N | N | N | N |
| `RS-RESULT = SPACES` | Y | Y | N | Y | Y | Y | Y | Y | Y | N | N |
| `RQ-PURCHASE` | N | N | | Y | Y | N | N | N | N | | |
| `RQ-REVERSAL` | Y | Y | | | | N | N | N | N | | |
| `RS-APPROVED` | Y | N | N | Y | N | Y | N | Y | N | N | N |
| `RQ-BALANCE-ENQUIRY` | N | | | N | | N | | Y | Y | | |
| `RQ-WITHDRAWAL` | | | | | | Y | Y | N | N | | |
| → `ADD 1 TO WS-REQUEST-COUNT` | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 |
| → `MOVE SPACES TO AUTH-RESPONSE HOLD-RECORD` | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 |
| → `PERFORM 2200-READ-CARD-AND-ACCOUNT` | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | |
| → `PERFORM 2300-READ-HOLD` | 4 | 4 | 4 | | | | | | | | |
| → `CALL 'AUTHREV' USING AUTH-REQUEST AUTH-RESPONSE BANK-CONTROL-RECORD ACCT-ONLINE-RECORD CARD-RECORD HOLD-RECORD` | 5 | 5 | | | | | | | | | |
| → `ADD 1 TO WS-APPROVED-COUNT` | 6 | | | 5 | | 5 | | 5 | | | |
| → `PERFORM 2400-SAVE-CHANGES` | 7 | | | 6 | | 6 | | | | | |
| → `ADD 1 TO WS-DECLINED-COUNT` | | 6 | 5 | | 5 | | 5 | | 5 | 4 | 5 |
| → `CALL 'AUTHPUR' USING AUTH-REQUEST AUTH-RESPONSE BANK-CONTROL-RECORD ACCT-ONLINE-RECORD CARD-RECORD HOLD-RECORD` | | | | 4 | 4 | | | | | | |
| → `CALL 'ATMWDL' USING AUTH-REQUEST AUTH-RESPONSE BANK-CONTROL-RECORD ACCT-ONLINE-RECORD CARD-RECORD HOLD-RECORD` | | | | | | 4 | 4 | | | | |
| → `CALL 'ATMBAL' USING AUTH-REQUEST AUTH-RESPONSE BANK-CONTROL-RECORD ACCT-ONLINE-RECORD CARD-RECORD HOLD-RECORD` | | | | | | | | 4 | 4 | | |
| → `SET RS-INVALID-REQUEST TO TRUE` | | | | | | | | | | | 3 |
| → `MOVE 'UNKNOWN REQUEST' TO RS-MESSAGE` | | | | | | | | | | | 4 |

- **Rule 1.** When `RQ-FUNCTION` is `'PURC'` (`RQ-PURCHASE`), `RQ-FUNCTION` is `'REVS'` (`RQ-REVERSAL`), `RQ-FUNCTION` is `'WDRL'` (`RQ-WITHDRAWAL`), or `RQ-FUNCTION` is `'BALQ'` (`RQ-BALANCE-ENQUIRY`); `RS-RESULT` is spaces and `RQ-FUNCTION` is `'REVS'` (`RQ-REVERSAL`); `RS-RESULT` is spaces; `RQ-FUNCTION` is not `'PURC'` (not `RQ-PURCHASE`); `RQ-FUNCTION` is `'REVS'` (`RQ-REVERSAL`); `RS-RESULT` is `'00'` (`RS-APPROVED`); and `RQ-FUNCTION` is not `'BALQ'` (not `RQ-BALANCE-ENQUIRY`), it adds 1 to `WS-REQUEST-COUNT`; sets `AUTH-RESPONSE` and `HOLD-RECORD` to spaces; performs `2200-READ-CARD-AND-ACCOUNT`; performs `2300-READ-HOLD`; calls `AUTHREV`; adds 1 to `WS-APPROVED-COUNT`; and performs `2400-SAVE-CHANGES`.
- **Rule 2.** When `RQ-FUNCTION` is `'PURC'` (`RQ-PURCHASE`), `RQ-FUNCTION` is `'REVS'` (`RQ-REVERSAL`), `RQ-FUNCTION` is `'WDRL'` (`RQ-WITHDRAWAL`), or `RQ-FUNCTION` is `'BALQ'` (`RQ-BALANCE-ENQUIRY`); `RS-RESULT` is spaces and `RQ-FUNCTION` is `'REVS'` (`RQ-REVERSAL`); `RS-RESULT` is spaces; `RQ-FUNCTION` is not `'PURC'` (not `RQ-PURCHASE`); `RQ-FUNCTION` is `'REVS'` (`RQ-REVERSAL`); and `RS-RESULT` is not `'00'` (not `RS-APPROVED`), it adds 1 to `WS-REQUEST-COUNT`; sets `AUTH-RESPONSE` and `HOLD-RECORD` to spaces; performs `2200-READ-CARD-AND-ACCOUNT`; performs `2300-READ-HOLD`; calls `AUTHREV`; and adds 1 to `WS-DECLINED-COUNT`.
- **Rule 3.** When `RQ-FUNCTION` is `'PURC'` (`RQ-PURCHASE`), `RQ-FUNCTION` is `'REVS'` (`RQ-REVERSAL`), `RQ-FUNCTION` is `'WDRL'` (`RQ-WITHDRAWAL`), or `RQ-FUNCTION` is `'BALQ'` (`RQ-BALANCE-ENQUIRY`); `RS-RESULT` is spaces and `RQ-FUNCTION` is `'REVS'` (`RQ-REVERSAL`); `RS-RESULT` is not spaces; and `RS-RESULT` is not `'00'` (not `RS-APPROVED`), it adds 1 to `WS-REQUEST-COUNT`; sets `AUTH-RESPONSE` and `HOLD-RECORD` to spaces; performs `2200-READ-CARD-AND-ACCOUNT`; performs `2300-READ-HOLD`; and adds 1 to `WS-DECLINED-COUNT`.
- **Rule 4.** When `RQ-FUNCTION` is `'PURC'` (`RQ-PURCHASE`), `RQ-FUNCTION` is `'REVS'` (`RQ-REVERSAL`), `RQ-FUNCTION` is `'WDRL'` (`RQ-WITHDRAWAL`), or `RQ-FUNCTION` is `'BALQ'` (`RQ-BALANCE-ENQUIRY`); `RS-RESULT` is not spaces or `RQ-FUNCTION` is not `'REVS'` (not `RQ-REVERSAL`); `RS-RESULT` is spaces; `RQ-FUNCTION` is `'PURC'` (`RQ-PURCHASE`); `RS-RESULT` is `'00'` (`RS-APPROVED`); and `RQ-FUNCTION` is not `'BALQ'` (not `RQ-BALANCE-ENQUIRY`), it adds 1 to `WS-REQUEST-COUNT`; sets `AUTH-RESPONSE` and `HOLD-RECORD` to spaces; performs `2200-READ-CARD-AND-ACCOUNT`; calls `AUTHPUR`; adds 1 to `WS-APPROVED-COUNT`; and performs `2400-SAVE-CHANGES`.
- **Rule 5.** When `RQ-FUNCTION` is `'PURC'` (`RQ-PURCHASE`), `RQ-FUNCTION` is `'REVS'` (`RQ-REVERSAL`), `RQ-FUNCTION` is `'WDRL'` (`RQ-WITHDRAWAL`), or `RQ-FUNCTION` is `'BALQ'` (`RQ-BALANCE-ENQUIRY`); `RS-RESULT` is not spaces or `RQ-FUNCTION` is not `'REVS'` (not `RQ-REVERSAL`); `RS-RESULT` is spaces; `RQ-FUNCTION` is `'PURC'` (`RQ-PURCHASE`); and `RS-RESULT` is not `'00'` (not `RS-APPROVED`), it adds 1 to `WS-REQUEST-COUNT`; sets `AUTH-RESPONSE` and `HOLD-RECORD` to spaces; performs `2200-READ-CARD-AND-ACCOUNT`; calls `AUTHPUR`; and adds 1 to `WS-DECLINED-COUNT`.
- **Rule 6.** When `RQ-FUNCTION` is `'PURC'` (`RQ-PURCHASE`), `RQ-FUNCTION` is `'REVS'` (`RQ-REVERSAL`), `RQ-FUNCTION` is `'WDRL'` (`RQ-WITHDRAWAL`), or `RQ-FUNCTION` is `'BALQ'` (`RQ-BALANCE-ENQUIRY`); `RS-RESULT` is not spaces or `RQ-FUNCTION` is not `'REVS'` (not `RQ-REVERSAL`); `RS-RESULT` is spaces; `RQ-FUNCTION` is not `'PURC'` (not `RQ-PURCHASE`); `RQ-FUNCTION` is not `'REVS'` (not `RQ-REVERSAL`); `RS-RESULT` is `'00'` (`RS-APPROVED`); `RQ-FUNCTION` is not `'BALQ'` (not `RQ-BALANCE-ENQUIRY`); and `RQ-FUNCTION` is `'WDRL'` (`RQ-WITHDRAWAL`), it adds 1 to `WS-REQUEST-COUNT`; sets `AUTH-RESPONSE` and `HOLD-RECORD` to spaces; performs `2200-READ-CARD-AND-ACCOUNT`; calls `ATMWDL`; adds 1 to `WS-APPROVED-COUNT`; and performs `2400-SAVE-CHANGES`.
- **Rule 7.** When `RQ-FUNCTION` is `'PURC'` (`RQ-PURCHASE`), `RQ-FUNCTION` is `'REVS'` (`RQ-REVERSAL`), `RQ-FUNCTION` is `'WDRL'` (`RQ-WITHDRAWAL`), or `RQ-FUNCTION` is `'BALQ'` (`RQ-BALANCE-ENQUIRY`); `RS-RESULT` is not spaces or `RQ-FUNCTION` is not `'REVS'` (not `RQ-REVERSAL`); `RS-RESULT` is spaces; `RQ-FUNCTION` is not `'PURC'` (not `RQ-PURCHASE`); `RQ-FUNCTION` is not `'REVS'` (not `RQ-REVERSAL`); `RS-RESULT` is not `'00'` (not `RS-APPROVED`); and `RQ-FUNCTION` is `'WDRL'` (`RQ-WITHDRAWAL`), it adds 1 to `WS-REQUEST-COUNT`; sets `AUTH-RESPONSE` and `HOLD-RECORD` to spaces; performs `2200-READ-CARD-AND-ACCOUNT`; calls `ATMWDL`; and adds 1 to `WS-DECLINED-COUNT`.
- **Rule 8.** When `RQ-FUNCTION` is `'PURC'` (`RQ-PURCHASE`), `RQ-FUNCTION` is `'REVS'` (`RQ-REVERSAL`), `RQ-FUNCTION` is `'WDRL'` (`RQ-WITHDRAWAL`), or `RQ-FUNCTION` is `'BALQ'` (`RQ-BALANCE-ENQUIRY`); `RS-RESULT` is not spaces or `RQ-FUNCTION` is not `'REVS'` (not `RQ-REVERSAL`); `RS-RESULT` is spaces; `RQ-FUNCTION` is not `'PURC'` (not `RQ-PURCHASE`); `RQ-FUNCTION` is not `'REVS'` (not `RQ-REVERSAL`); `RS-RESULT` is `'00'` (`RS-APPROVED`); `RQ-FUNCTION` is `'BALQ'` (`RQ-BALANCE-ENQUIRY`); and `RQ-FUNCTION` is not `'WDRL'` (not `RQ-WITHDRAWAL`), it adds 1 to `WS-REQUEST-COUNT`; sets `AUTH-RESPONSE` and `HOLD-RECORD` to spaces; performs `2200-READ-CARD-AND-ACCOUNT`; calls `ATMBAL`; and adds 1 to `WS-APPROVED-COUNT`.
- **Rule 9.** When `RQ-FUNCTION` is `'PURC'` (`RQ-PURCHASE`), `RQ-FUNCTION` is `'REVS'` (`RQ-REVERSAL`), `RQ-FUNCTION` is `'WDRL'` (`RQ-WITHDRAWAL`), or `RQ-FUNCTION` is `'BALQ'` (`RQ-BALANCE-ENQUIRY`); `RS-RESULT` is not spaces or `RQ-FUNCTION` is not `'REVS'` (not `RQ-REVERSAL`); `RS-RESULT` is spaces; `RQ-FUNCTION` is not `'PURC'` (not `RQ-PURCHASE`); `RQ-FUNCTION` is not `'REVS'` (not `RQ-REVERSAL`); `RS-RESULT` is not `'00'` (not `RS-APPROVED`); `RQ-FUNCTION` is `'BALQ'` (`RQ-BALANCE-ENQUIRY`); and `RQ-FUNCTION` is not `'WDRL'` (not `RQ-WITHDRAWAL`), it adds 1 to `WS-REQUEST-COUNT`; sets `AUTH-RESPONSE` and `HOLD-RECORD` to spaces; performs `2200-READ-CARD-AND-ACCOUNT`; calls `ATMBAL`; and adds 1 to `WS-DECLINED-COUNT`.
- **Rule 10.** When `RQ-FUNCTION` is `'PURC'` (`RQ-PURCHASE`), `RQ-FUNCTION` is `'REVS'` (`RQ-REVERSAL`), `RQ-FUNCTION` is `'WDRL'` (`RQ-WITHDRAWAL`), or `RQ-FUNCTION` is `'BALQ'` (`RQ-BALANCE-ENQUIRY`); `RS-RESULT` is not spaces or `RQ-FUNCTION` is not `'REVS'` (not `RQ-REVERSAL`); `RS-RESULT` is not spaces; and `RS-RESULT` is not `'00'` (not `RS-APPROVED`), it adds 1 to `WS-REQUEST-COUNT`; sets `AUTH-RESPONSE` and `HOLD-RECORD` to spaces; performs `2200-READ-CARD-AND-ACCOUNT`; and adds 1 to `WS-DECLINED-COUNT`.
- **Rule 11.** When `RQ-FUNCTION` is not `'PURC'` (not `RQ-PURCHASE`), `RQ-FUNCTION` is not `'REVS'` (not `RQ-REVERSAL`), `RQ-FUNCTION` is not `'WDRL'` (not `RQ-WITHDRAWAL`), and `RQ-FUNCTION` is not `'BALQ'` (not `RQ-BALANCE-ENQUIRY`); `RS-RESULT` is not spaces or `RQ-FUNCTION` is not `'REVS'` (not `RQ-REVERSAL`); `RS-RESULT` is not spaces; and `RS-RESULT` is not `'00'` (not `RS-APPROVED`), it adds 1 to `WS-REQUEST-COUNT`; sets `AUTH-RESPONSE` and `HOLD-RECORD` to spaces; sets `RS-RESULT` to `RS-INVALID-REQUEST`; sets `RS-MESSAGE` to `'UNKNOWN REQUEST'`; and adds 1 to `WS-DECLINED-COUNT`.

Source: [AUTHSRV.cbl:160](../../../pangaea/programs/AUTHSRV.cbl#L160).

## 2200-READ-CARD-AND-ACCOUNT

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `READ CARD-FILE INVALID KEY` | Y | N | N |
| `READ ACCT-FILE INVALID KEY` | | Y | N |
| → `MOVE RQ-CARD-NO TO CR-CARD-NO` | 1 | 1 | 1 |
| → `READ CARD-FILE` | 2 | 2 | 2 |
| → `SET RS-NO-SUCH-CARD TO TRUE` | 3 | | |
| → `MOVE 'NO SUCH CARD' TO RS-MESSAGE` | 4 | | |
| → `MOVE CR-ACCT-NO TO AO-ACCT-NO` | | 3 | 3 |
| → `READ ACCT-FILE` | | 4 | 4 |
| → `SET RS-ACCOUNT-NOT-OPEN TO TRUE` | | 5 | |
| → `MOVE 'NO SUCH ACCOUNT' TO RS-MESSAGE` | | 6 | |

- **Rule 1.** When reading `CARD-FILE` fails on its key, it sets `CR-CARD-NO` to `RQ-CARD-NO`, reads `CARD-FILE`, sets `RS-RESULT` to `RS-NO-SUCH-CARD`, and sets `RS-MESSAGE` to `'NO SUCH CARD'`.
- **Rule 2.** When reading `CARD-FILE` succeeds by its key and reading `ACCT-FILE` fails on its key, it sets `CR-CARD-NO` to `RQ-CARD-NO`, reads `CARD-FILE`, sets `AO-ACCT-NO` to `CR-ACCT-NO`, reads `ACCT-FILE`, sets `RS-RESULT` to `RS-ACCOUNT-NOT-OPEN`, and sets `RS-MESSAGE` to `'NO SUCH ACCOUNT'`.
- **Rule 3.** When reading `CARD-FILE` succeeds by its key and reading `ACCT-FILE` succeeds by its key, it sets `CR-CARD-NO` to `RQ-CARD-NO`, reads `CARD-FILE`, sets `AO-ACCT-NO` to `CR-ACCT-NO`, and reads `ACCT-FILE`.

Source: [AUTHSRV.cbl:203](../../../pangaea/programs/AUTHSRV.cbl#L203).

## 2300-READ-HOLD

| | 1 | 2 |
| --- | :-: | :-: |
| `READ HOLD-FILE INVALID KEY` | Y | N |
| → `MOVE CR-ACCT-NO TO HD-ACCT-NO` | 1 | 1 |
| → `MOVE RQ-AUTH-CODE TO HD-HOLD-NO` | 2 | 2 |
| → `READ HOLD-FILE` | 3 | 3 |
| → `SET RS-NO-SUCH-AUTH TO TRUE` | 4 | |
| → `MOVE 'NO SUCH AUTHORIZATION' TO RS-MESSAGE` | 5 | |

- **Rule 1.** When reading `HOLD-FILE` fails on its key, it sets `HD-ACCT-NO` to `CR-ACCT-NO`, sets `HD-HOLD-NO` to `RQ-AUTH-CODE`, reads `HOLD-FILE`, sets `RS-RESULT` to `RS-NO-SUCH-AUTH`, and sets `RS-MESSAGE` to `'NO SUCH AUTHORIZATION'`.
- **Rule 2.** When reading `HOLD-FILE` succeeds by its key, it sets `HD-ACCT-NO` to `CR-ACCT-NO`, sets `HD-HOLD-NO` to `RQ-AUTH-CODE`, and reads `HOLD-FILE`.

Source: [AUTHSRV.cbl:219](../../../pangaea/programs/AUTHSRV.cbl#L219).

## 2400-SAVE-CHANGES

| | 1 | 2 | 3 | 4 | 5 | 6 |
| --- | :-: | :-: | :-: | :-: | :-: | :-: |
| `REWRITE ACCT-ONLINE-RECORD INVALID KEY` | Y | N | N | N | N | N |
| `REWRITE CARD-RECORD INVALID KEY` | | Y | N | N | N | N |
| `RQ-REVERSAL` | | | Y | Y | N | N |
| `REWRITE HOLD-RECORD INVALID KEY` | | | Y | N | | |
| `WRITE HOLD-RECORD INVALID KEY` | | | | | Y | N |
| → `REWRITE ACCT-ONLINE-RECORD` | 1 | 1 | 1 | 1 | 1 | 1 |
| → `MOVE 'CANNOT REWRITE THE ACCOUNT' TO WS-ABORT-MESSAGE` | 2 | | | | | |
| → `PERFORM 9900-ABORT` | 3 | 4 | 5 | | 7 | |
| → `REWRITE CARD-RECORD` | | 2 | 2 | 2 | 2 | 2 |
| → `MOVE 'CANNOT REWRITE THE CARD' TO WS-ABORT-MESSAGE` | | 3 | | | | |
| → `REWRITE HOLD-RECORD` | | | 3 | 3 | | |
| → `MOVE 'CANNOT REWRITE THE HOLD' TO WS-ABORT-MESSAGE` | | | 4 | | | |
| → `ADD HD-AMOUNT TO WS-RELEASED-TOTAL` | | | | 4 | | |
| → `ADD 1 TO WS-LAST-HOLD-NO` | | | | | 3 | 3 |
| → `MOVE WS-LAST-HOLD-NO TO HD-HOLD-NO RS-AUTH-CODE` | | | | | 4 | 4 |
| → `WRITE HOLD-RECORD` | | | | | 5 | 5 |
| → `MOVE 'CANNOT WRITE THE HOLD' TO WS-ABORT-MESSAGE` | | | | | 6 | |
| → `ADD HD-AMOUNT TO WS-HELD-TOTAL` | | | | | | 6 |
| → `ADD HD-FEE TO WS-HELD-TOTAL` | | | | | | 7 |

- **Rule 1.** When rewriting `ACCT-ONLINE-RECORD` fails on its key, it rewrites `ACCT-ONLINE-RECORD`, sets `WS-ABORT-MESSAGE` to `'CANNOT REWRITE THE ACCOUNT'`, and performs `9900-ABORT`.
- **Rule 2.** When rewriting `ACCT-ONLINE-RECORD` succeeds by its key and rewriting `CARD-RECORD` fails on its key, it rewrites `ACCT-ONLINE-RECORD`, rewrites `CARD-RECORD`, sets `WS-ABORT-MESSAGE` to `'CANNOT REWRITE THE CARD'`, and performs `9900-ABORT`.
- **Rule 3.** When rewriting `ACCT-ONLINE-RECORD` succeeds by its key, rewriting `CARD-RECORD` succeeds by its key, `RQ-FUNCTION` is `'REVS'` (`RQ-REVERSAL`), and rewriting `HOLD-RECORD` fails on its key, it rewrites `ACCT-ONLINE-RECORD`, rewrites `CARD-RECORD`, rewrites `HOLD-RECORD`, sets `WS-ABORT-MESSAGE` to `'CANNOT REWRITE THE HOLD'`, and performs `9900-ABORT`.
- **Rule 4.** When rewriting `ACCT-ONLINE-RECORD` succeeds by its key, rewriting `CARD-RECORD` succeeds by its key, `RQ-FUNCTION` is `'REVS'` (`RQ-REVERSAL`), and rewriting `HOLD-RECORD` succeeds by its key, it rewrites `ACCT-ONLINE-RECORD`, rewrites `CARD-RECORD`, rewrites `HOLD-RECORD`, and adds `HD-AMOUNT` to `WS-RELEASED-TOTAL`.
- **Rule 5.** When rewriting `ACCT-ONLINE-RECORD` succeeds by its key, rewriting `CARD-RECORD` succeeds by its key, `RQ-FUNCTION` is not `'REVS'` (not `RQ-REVERSAL`), and writing `HOLD-RECORD` fails on its key, it rewrites `ACCT-ONLINE-RECORD`; rewrites `CARD-RECORD`; adds 1 to `WS-LAST-HOLD-NO`; sets `HD-HOLD-NO` and `RS-AUTH-CODE` to `WS-LAST-HOLD-NO`; writes `HOLD-RECORD`; sets `WS-ABORT-MESSAGE` to `'CANNOT WRITE THE HOLD'`; and performs `9900-ABORT`.
- **Rule 6.** When rewriting `ACCT-ONLINE-RECORD` succeeds by its key, rewriting `CARD-RECORD` succeeds by its key, `RQ-FUNCTION` is not `'REVS'` (not `RQ-REVERSAL`), and writing `HOLD-RECORD` succeeds by its key, it rewrites `ACCT-ONLINE-RECORD`; rewrites `CARD-RECORD`; adds 1 to `WS-LAST-HOLD-NO`; sets `HD-HOLD-NO` and `RS-AUTH-CODE` to `WS-LAST-HOLD-NO`; writes `HOLD-RECORD`; adds `HD-AMOUNT` to `WS-HELD-TOTAL`; and adds `HD-FEE` to `WS-HELD-TOTAL`.

Source: [AUTHSRV.cbl:230](../../../pangaea/programs/AUTHSRV.cbl#L230).
