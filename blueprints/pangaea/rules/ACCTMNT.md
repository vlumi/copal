# ACCTMNT: decision tables

The business rules of [ACCTMNT](../programs/ACCTMNT.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

## Contents

- [0000-MAIN](#0000-main)
- [1000-INITIALIZE](#1000-initialize)
- [1200-READ-MASTER](#1200-read-master)
- [1300-READ-REQUEST](#1300-read-request)
- [2000-NEXT-ACCOUNT](#2000-next-account)
- [2100-APPLY-REQUESTS](#2100-apply-requests)
- [2200-APPLY-REQUEST](#2200-apply-request)
- [2210-OPEN](#2210-open)
- [2220-CLOSE](#2220-close)
- [2230-FREEZE](#2230-freeze)
- [2240-UNFREEZE](#2240-unfreeze)
- [2250-LIMIT](#2250-limit)
- [2800-WRITE-HELD](#2800-write-held)

## 0000-MAIN

| | 1 | 2 | 3 | 4 |
| --- | :-: | :-: | :-: | :-: |
| `WS-BALANCES-OUT = WS-BALANCES-IN` | N | N | Y | Y |
| `WS-REJECT-COUNT > 0` | Y | N | Y | N |
| → `PERFORM 1000-INITIALIZE` | 1 | 1 | 1 | 1 |
| → `PERFORM 2000-NEXT-ACCOUNT THRU 2999-EXIT UNTIL WS-MASTER-KEY = HIGH-VALUES AND WS-REQUEST-KEY = HIGH-VALUES` | 2 | 2 | 2 | 2 |
| → `MOVE 'NEW MASTER OUT OF BALANCE' TO WS-ABORT-MESSAGE` | 3 | 3 | | |
| → `PERFORM 9900-ABORT` | 4 | 4 | | |
| → `PERFORM 8000-WRITE-REPORT` | 5 | 5 | 3 | 3 |
| → `PERFORM 9000-CLOSE-FILES` | 6 | 6 | 4 | 4 |
| → `MOVE 4 TO RETURN-CODE` | 7 | | 5 | |
| → `STOP RUN` | 8 | 8 | 6 | 6 |
| → `MOVE 0 TO RETURN-CODE` | | 7 | | 5 |

Source: [ACCTMNT.cbl:128](../../../pangaea/programs/ACCTMNT.cbl#L128).

## 1000-INITIALIZE

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `WS-FILE-STATUSES = ALL '0'` | N | Y | Y |
| `READ CONTROL-FILE AT END` | | Y | N |
| → `OPEN INPUT CONTROL-FILE ACCT-IN MAINT-IN OUTPUT ACCT-OUT REJECT-OUT REPORT-OUT` | 1 | 1 | 1 |
| → `DISPLAY 'ACCTMNT: CANNOT OPEN FILES ' WS-FILE-STATUSES` | 2 | | |
| → `MOVE 16 TO RETURN-CODE` | 3 | | |
| → `STOP RUN` | 4 | | |
| → `READ CONTROL-FILE` | | 2 | 2 |
| → `MOVE 'NO BANK CONTROL RECORD' TO WS-ABORT-MESSAGE` | | 3 | |
| → `PERFORM 9900-ABORT` | | 4 | |
| → `MOVE BC-BUSINESS-DATE TO WS-BUSINESS-DATE` | | 5 | 3 |
| → `MOVE SPACES TO REJECT-RECORD` | | 6 | 4 |
| → `PERFORM 1200-READ-MASTER THRU 1200-EXIT` | | 7 | 5 |
| → `PERFORM 1300-READ-REQUEST THRU 1300-EXIT` | | 8 | 6 |

Source: [ACCTMNT.cbl:145](../../../pangaea/programs/ACCTMNT.cbl#L145).

## 1200-READ-MASTER

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `READ ACCT-IN AT END` | Y | N | N |
| `WS-MASTER-KEY > WS-PREV-MASTER-KEY` | | N | Y |
| → `READ ACCT-IN` | 1 | 1 | 1 |
| → `MOVE HIGH-VALUES TO WS-MASTER-KEY` | 2 | | |
| → `GO TO 1200-EXIT` | 3 | | |
| → `MOVE AI-ACCT-NO TO WS-MASTER-KEY` | | 2 | 2 |
| → `MOVE 'MASTER NOT IN ACCOUNT ORDER' TO WS-ABORT-MESSAGE` | | 3 | |
| → `PERFORM 9900-ABORT` | | 4 | |
| → `MOVE WS-MASTER-KEY TO WS-PREV-MASTER-KEY` | | 5 | 3 |
| → `ADD 1 TO WS-MASTERS-READ` | | 6 | 4 |

Source: [ACCTMNT.cbl:162](../../../pangaea/programs/ACCTMNT.cbl#L162).

## 1300-READ-REQUEST

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `READ MAINT-IN AT END` | Y | N | N |
| `WS-REQUEST-ORDER-KEY > WS-PREV-REQUEST-ORDER-KEY` | | N | Y |
| → `READ MAINT-IN` | 1 | 1 | 1 |
| → `MOVE HIGH-VALUES TO WS-REQUEST-KEY` | 2 | | |
| → `GO TO 1300-EXIT` | 3 | | |
| → `MOVE MR-ACCT-NO TO WS-REQUEST-KEY WS-ROK-ACCT-NO` | | 2 | 2 |
| → `MOVE MR-REQUEST-ID TO WS-ROK-REQUEST-ID` | | 3 | 3 |
| → `MOVE 'REQUESTS NOT IN ACCOUNT ORDER' TO WS-ABORT-MESSAGE` | | 4 | |
| → `PERFORM 9900-ABORT` | | 5 | |
| → `MOVE WS-REQUEST-ORDER-KEY TO WS-PREV-REQUEST-ORDER-KEY` | | 6 | 4 |
| → `ADD 1 TO WS-REQUEST-COUNT` | | 7 | 5 |

Source: [ACCTMNT.cbl:177](../../../pangaea/programs/ACCTMNT.cbl#L177).

## 2000-NEXT-ACCOUNT

| | 1 | 2 | 3 | 4 |
| --- | :-: | :-: | :-: | :-: |
| `WS-REQUEST-KEY < WS-MASTER-KEY` | Y | Y | N | N |
| `WS-MASTER-KEY = WS-ACTIVE-KEY` | Y | N | Y | N |
| → `MOVE WS-MASTER-KEY TO WS-ACTIVE-KEY` | 1 | 1 | 1 | 1 |
| → `MOVE WS-REQUEST-KEY TO WS-ACTIVE-KEY` | 2 | 2 | | |
| → `SET NOTHING-HELD TO TRUE` | 3 | 3 | 2 | 2 |
| → `MOVE ACCT-IN-RECORD TO ACCT-MASTER-RECORD` | 4 | | 3 | |
| → `ADD AM-BALANCE TO WS-BALANCES-IN` | 5 | | 4 | |
| → `SET ACCOUNT-HELD TO TRUE` | 6 | | 5 | |
| → `PERFORM 1200-READ-MASTER THRU 1200-EXIT` | 7 | | 6 | |

Source: [ACCTMNT.cbl:195](../../../pangaea/programs/ACCTMNT.cbl#L195).

## 2100-APPLY-REQUESTS

| | 1 | 2 |
| --- | :-: | :-: |
| `WS-REQUEST-KEY = WS-ACTIVE-KEY` | N | Y |
| → `GO TO 2800-WRITE-HELD` | 1 | |
| → `PERFORM 2200-APPLY-REQUEST THRU 2299-EXIT` | | 1 |
| → `PERFORM 1300-READ-REQUEST THRU 1300-EXIT` | | 2 |
| → `GO TO 2100-APPLY-REQUESTS` | | 3 |

Source: [ACCTMNT.cbl:207](../../../pangaea/programs/ACCTMNT.cbl#L207).

## 2200-APPLY-REQUEST

| | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 |
| --- | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: |
| `MR-OPEN` | Y | N | N | N | N | N | N | N |
| `NOTHING-HELD` | | Y | N | N | N | N | N | N |
| `AM-CLOSED` | | | Y | N | N | N | N | N |
| `MR-CLOSE` | | | | Y | N | N | N | N |
| `MR-FREEZE` | | | | | Y | N | N | N |
| `MR-UNFREEZE` | | | | | | Y | N | N |
| `MR-LIMIT` | | | | | | | Y | N |
| → `GO TO 2210-OPEN` | 1 | | | | | | | |
| → `SET RJ-NO-ACCOUNT TO TRUE` | | 1 | | | | | | |
| → `GO TO 2290-REJECT` | | 2 | 2 | | | | | 2 |
| → `SET RJ-CLOSED TO TRUE` | | | 1 | | | | | |
| → `GO TO 2220-CLOSE` | | | | 1 | | | | |
| → `GO TO 2230-FREEZE` | | | | | 1 | | | |
| → `GO TO 2240-UNFREEZE` | | | | | | 1 | | |
| → `GO TO 2250-LIMIT` | | | | | | | 1 | |
| → `SET RJ-BAD-TYPE TO TRUE` | | | | | | | | 1 |

Source: [ACCTMNT.cbl:215](../../../pangaea/programs/ACCTMNT.cbl#L215).

## 2210-OPEN

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `ACCOUNT-HELD` | Y | N | N |
| `AM-SAVINGS` | | Y | N |
| → `SET RJ-EXISTS TO TRUE` | 1 | | |
| → `GO TO 2290-REJECT` | 2 | | |
| → `MOVE SPACES TO ACCT-MASTER-RECORD` | | 1 | 1 |
| → `MOVE MR-ACCT-NO TO AM-ACCT-NO` | | 2 | 2 |
| → `MOVE MR-CUST-NO TO AM-CUST-NO` | | 3 | 3 |
| → `MOVE MR-SHORT-NAME TO AM-SHORT-NAME` | | 4 | 4 |
| → `MOVE MR-ACCT-TYPE TO AM-ACCT-TYPE` | | 5 | 5 |
| → `SET AM-OPEN TO TRUE` | | 6 | 6 |
| → `MOVE WS-BUSINESS-DATE TO AM-OPEN-DATE AM-LAST-ACTIVITY-DATE` | | 7 | 7 |
| → `MOVE 0 TO AM-BALANCE` | | 8 | 8 |
| → `MOVE MR-OVERDRAFT-LIMIT TO AM-OVERDRAFT-LIMIT` | | 9 | 9 |
| → `MOVE 0 TO AM-OVERDRAFT-LIMIT` | | 10 | |
| → `SET ACCOUNT-HELD TO TRUE` | | 11 | 10 |
| → `ADD 1 TO WS-OPENED-COUNT` | | 12 | 11 |
| → `GO TO 2299-EXIT` | | 13 | 12 |

Source: [ACCTMNT.cbl:236](../../../pangaea/programs/ACCTMNT.cbl#L236).

## 2220-CLOSE

| | 1 | 2 |
| --- | :-: | :-: |
| `AM-BALANCE = 0` | N | Y |
| → `SET RJ-HAS-BALANCE TO TRUE` | 1 | |
| → `GO TO 2290-REJECT` | 2 | |
| → `SET AM-CLOSED TO TRUE` | | 1 |
| → `MOVE 0 TO AM-OVERDRAFT-LIMIT` | | 2 |
| → `ADD 1 TO WS-CLOSED-COUNT` | | 3 |
| → `GO TO 2299-EXIT` | | 4 |

Source: [ACCTMNT.cbl:256](../../../pangaea/programs/ACCTMNT.cbl#L256).

## 2230-FREEZE

| | 1 | 2 |
| --- | :-: | :-: |
| `AM-OPEN` | N | Y |
| → `SET RJ-BAD-STATUS TO TRUE` | 1 | |
| → `GO TO 2290-REJECT` | 2 | |
| → `SET AM-FROZEN TO TRUE` | | 1 |
| → `ADD 1 TO WS-FROZEN-COUNT` | | 2 |
| → `GO TO 2299-EXIT` | | 3 |

Source: [ACCTMNT.cbl:266](../../../pangaea/programs/ACCTMNT.cbl#L266).

## 2240-UNFREEZE

| | 1 | 2 |
| --- | :-: | :-: |
| `AM-FROZEN` | N | Y |
| → `SET RJ-BAD-STATUS TO TRUE` | 1 | |
| → `GO TO 2290-REJECT` | 2 | |
| → `SET AM-OPEN TO TRUE` | | 1 |
| → `ADD 1 TO WS-UNFROZEN-COUNT` | | 2 |
| → `GO TO 2299-EXIT` | | 3 |

Source: [ACCTMNT.cbl:275](../../../pangaea/programs/ACCTMNT.cbl#L275).

## 2250-LIMIT

| | 1 | 2 |
| --- | :-: | :-: |
| `AM-CURRENT` | N | Y |
| → `SET RJ-BAD-TYPE TO TRUE` | 1 | |
| → `GO TO 2290-REJECT` | 2 | |
| → `MOVE MR-OVERDRAFT-LIMIT TO AM-OVERDRAFT-LIMIT` | | 1 |
| → `ADD 1 TO WS-LIMIT-COUNT` | | 2 |
| → `GO TO 2299-EXIT` | | 3 |

Source: [ACCTMNT.cbl:284](../../../pangaea/programs/ACCTMNT.cbl#L284).

## 2800-WRITE-HELD

| | 1 | 2 |
| --- | :-: | :-: |
| `ACCOUNT-HELD` | Y | N |
| → `WRITE ACCT-OUT-RECORD FROM ACCT-MASTER-RECORD` | 1 | |
| → `ADD 1 TO WS-MASTERS-WRITTEN` | 2 | |
| → `ADD AM-BALANCE TO WS-BALANCES-OUT` | 3 | |

Source: [ACCTMNT.cbl:302](../../../pangaea/programs/ACCTMNT.cbl#L302).
