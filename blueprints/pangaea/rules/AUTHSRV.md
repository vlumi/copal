# AUTHSRV: decision tables

The business rules of [AUTHSRV](../programs/AUTHSRV.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

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

Source: [AUTHSRV.cbl:113](../../../pangaea/programs/AUTHSRV.cbl#L113).

## 1100-FIND-LAST-HOLD

| | 1 | 2 |
| --- | :-: | :-: |
| `START HOLD-FILE KEY IS NOT LESS THAN HD-KEY INVALID KEY` | Y | N |
| → `MOVE LOW-VALUES TO HD-KEY` | 1 | 1 |
| → `START HOLD-FILE KEY IS NOT LESS THAN HD-KEY` | 2 | 2 |
| → `SET END-OF-HOLDS TO TRUE` | 3 | |
| → `PERFORM UNTIL END-OF-HOLDS` | 4 | 3 |

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

Source: [AUTHSRV.cbl:137](../../../pangaea/programs/AUTHSRV.cbl#L137).

## 2000-SERVE-REQUEST

| | 1 | 2 |
| --- | :-: | :-: |
| `ACCEPT AUTH-REQUEST ON EXCEPTION` | Y | N |
| → `ACCEPT AUTH-REQUEST` | 1 | 1 |
| → `SET END-OF-INPUT TO TRUE` | 2 | |
| → `PERFORM 2100-AUTHORIZE` | | 2 |
| → `DISPLAY AUTH-RESPONSE` | | 3 |

Source: [AUTHSRV.cbl:149](../../../pangaea/programs/AUTHSRV.cbl#L149).

## 2100-AUTHORIZE

| | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 | 16 | 17 | 18 | 19 | 20 | 21 | 22 | 23 | 24 | 25 | 26 | 27 | 28 | 29 | 30 | 31 | 32 | 33 | 34 | 35 | 36 | 37 | 38 | 39 | 40 | 41 | 42 | 43 | 44 | 45 | 46 | 47 | 48 | 49 | 50 | 51 | 52 | 53 | 54 | 55 | 56 | 57 | 58 | 59 | 60 | 61 | 62 | 63 | 64 |
| --- | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: |
| `RQ-PURCHASE OR RQ-REVERSAL OR RQ-WITHDRAWAL OR RQ-BALANCE-ENQUIRY` | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | N | N | N | N | N | N | N | N | N | N | N | N | N | N | N | N | N | N | N | N | N | N | N | N | N | N | N | N | N | N | N | N |
| `RS-RESULT = SPACES AND RQ-REVERSAL` | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | N | N | N | N | N | N | N | N | N | N | N | N | N | N | N | N | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | N | N | N | N | N | N | N | N | N | N | N | N | N | N | N | N |
| `RS-RESULT = SPACES` | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | N | N | N | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | N | N | N | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | N | N | N | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | N | N | N |
| `RQ-PURCHASE` | Y | Y | Y | N | N | N | N | N | N | N | N | N | N | | | | Y | Y | Y | N | N | N | N | N | N | N | N | N | N | | | | Y | Y | Y | N | N | N | N | N | N | N | N | N | N | | | | Y | Y | Y | N | N | N | N | N | N | N | N | N | N | | | |
| `RS-APPROVED` | Y | Y | N | Y | Y | N | Y | Y | N | Y | N | Y | N | Y | Y | N | Y | Y | N | Y | Y | N | Y | Y | N | Y | N | Y | N | Y | Y | N | Y | Y | N | Y | Y | N | Y | Y | N | Y | N | Y | N | Y | Y | N | Y | Y | N | Y | Y | N | Y | Y | N | Y | N | Y | N | Y | Y | N |
| `RQ-BALANCE-ENQUIRY` | N | Y | | N | Y | | N | Y | | Y | Y | N | N | N | Y | | N | Y | | N | Y | | N | Y | | Y | Y | N | N | N | Y | | N | Y | | N | Y | | N | Y | | Y | Y | N | N | N | Y | | N | Y | | N | Y | | N | Y | | Y | Y | N | N | N | Y | |
| `RQ-REVERSAL` | | | | Y | Y | Y | N | N | N | N | N | N | N | | | | | | | Y | Y | Y | N | N | N | N | N | N | N | | | | | | | Y | Y | Y | N | N | N | N | N | N | N | | | | | | | Y | Y | Y | N | N | N | N | N | N | N | | | |
| `RQ-WITHDRAWAL` | | | | | | | Y | Y | Y | N | N | N | N | | | | | | | | | | Y | Y | Y | N | N | N | N | | | | | | | | | | Y | Y | Y | N | N | N | N | | | | | | | | | | Y | Y | Y | N | N | N | N | | | |
| → `ADD 1 TO WS-REQUEST-COUNT` | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 |
| → `MOVE SPACES TO AUTH-RESPONSE HOLD-RECORD` | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 |
| → `PERFORM 2200-READ-CARD-AND-ACCOUNT` | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | | | | | | | | | | | | | | | | | | | | | | | | | | | | | | | | |
| → `PERFORM 2300-READ-HOLD` | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | | | | | | | | | | | | | | | | | 5 | 5 | 5 | 5 | 5 | 5 | 5 | 5 | 5 | 5 | 5 | 5 | 5 | 5 | 5 | 5 | | | | | | | | | | | | | | | | |
| → `CALL 'AUTHPUR' USING AUTH-REQUEST AUTH-RESPONSE BANK-CONTROL-RECORD ACCT-ONLINE-RECORD CARD-RECORD HOLD-RECORD` | 5 | 5 | 5 | | | | | | | | | | | | | | 4 | 4 | 4 | | | | | | | | | | | | | | 6 | 6 | 6 | | | | | | | | | | | | | | 5 | 5 | 5 | | | | | | | | | | | | | |
| → `ADD 1 TO WS-APPROVED-COUNT` | 6 | 6 | | 6 | 6 | | 6 | 6 | | 6 | | 5 | | 5 | 5 | | 5 | 5 | | 5 | 5 | | 5 | 5 | | 5 | | 4 | | 4 | 4 | | 7 | 7 | | 7 | 7 | | 7 | 7 | | 7 | | 6 | | 6 | 6 | | 6 | 6 | | 6 | 6 | | 6 | 6 | | 6 | | 5 | | 5 | 5 | |
| → `PERFORM 2400-SAVE-CHANGES` | 7 | | | 7 | | | 7 | | | | | 6 | | 6 | | | 6 | | | 6 | | | 6 | | | | | 5 | | 5 | | | 8 | | | 8 | | | 8 | | | | | 7 | | 7 | | | 7 | | | 7 | | | 7 | | | | | 6 | | 6 | | |
| → `ADD 1 TO WS-DECLINED-COUNT` | | | 6 | | | 6 | | | 6 | | 6 | | 5 | | | 5 | | | 5 | | | 5 | | | 5 | | 5 | | 4 | | | 4 | | | 7 | | | 7 | | | 7 | | 7 | | 6 | | | 6 | | | 6 | | | 6 | | | 6 | | 6 | | 5 | | | 5 |
| → `CALL 'AUTHREV' USING AUTH-REQUEST AUTH-RESPONSE BANK-CONTROL-RECORD ACCT-ONLINE-RECORD CARD-RECORD HOLD-RECORD` | | | | 5 | 5 | 5 | | | | | | | | | | | | | | 4 | 4 | 4 | | | | | | | | | | | | | | 6 | 6 | 6 | | | | | | | | | | | | | | 5 | 5 | 5 | | | | | | | | | | |
| → `CALL 'ATMWDL' USING AUTH-REQUEST AUTH-RESPONSE BANK-CONTROL-RECORD ACCT-ONLINE-RECORD CARD-RECORD HOLD-RECORD` | | | | | | | 5 | 5 | 5 | | | | | | | | | | | | | | 4 | 4 | 4 | | | | | | | | | | | | | | 6 | 6 | 6 | | | | | | | | | | | | | | 5 | 5 | 5 | | | | | | | |
| → `CALL 'ATMBAL' USING AUTH-REQUEST AUTH-RESPONSE BANK-CONTROL-RECORD ACCT-ONLINE-RECORD CARD-RECORD HOLD-RECORD` | | | | | | | | | | 5 | 5 | | | | | | | | | | | | | | | 4 | 4 | | | | | | | | | | | | | | | 6 | 6 | | | | | | | | | | | | | | | 5 | 5 | | | | | |
| → `SET RS-INVALID-REQUEST TO TRUE` | | | | | | | | | | | | | | | | | | | | | | | | | | | | | | | | | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 |
| → `MOVE 'UNKNOWN REQUEST' TO RS-MESSAGE` | | | | | | | | | | | | | | | | | | | | | | | | | | | | | | | | | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 |

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

Source: [AUTHSRV.cbl:230](../../../pangaea/programs/AUTHSRV.cbl#L230).
