# TXVALID: decision tables

The business rules of [TXVALID](../programs/TXVALID.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

## Contents

- [0000-MAIN](#0000-main)
- [1000-OPEN-FILES](#1000-open-files)
- [2000-READ-HEADER](#2000-read-header)
- [3000-READ-DETAILS, in its loop](#3000-read-details-in-its-loop)
- [3100-READ-NEXT](#3100-read-next)
- [3200-VALIDATE-DETAIL](#3200-validate-detail)
- [4000-CHECK-TRAILER](#4000-check-trailer)
- [8000-WRITE-REPORT](#8000-write-report)

## 0000-MAIN

| | 1 | 2 | 3 | 4 | 5 |
| --- | :-: | :-: | :-: | :-: | :-: |
| `FILE-IS-GOOD` | Y | Y | Y | Y | N |
| `FILE-IS-GOOD, again` | Y | Y | Y | N | |
| `FILE-IS-GOOD, again (3)` | N | Y | Y | | |
| `WS-REJECT-COUNT > 0` | | Y | N | | |
| → `PERFORM 1000-OPEN-FILES` | 1 | 1 | 1 | 1 | 1 |
| → `PERFORM 2000-READ-HEADER` | 2 | 2 | 2 | 2 | 2 |
| → `PERFORM 3000-READ-DETAILS` | 3 | 3 | 3 | 3 | |
| → `PERFORM 4000-CHECK-TRAILER` | 4 | 4 | 4 | | |
| → `PERFORM 8000-WRITE-REPORT` | 5 | 5 | 5 | 4 | 3 |
| → `PERFORM 9000-CLOSE-FILES` | 6 | 6 | 6 | 5 | 4 |
| → `MOVE 12 TO RETURN-CODE` | 7 | | | 6 | 5 |
| → `STOP RUN` | 8 | 8 | 8 | 7 | 6 |
| → `MOVE 4 TO RETURN-CODE` | | 7 | | | |
| → `MOVE 0 TO RETURN-CODE` | | | 7 | | |

Source: [TXVALID.cbl:113](../../../pangaea/programs/TXVALID.cbl#L113).

## 1000-OPEN-FILES

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `WS-CONTROL-STATUS NOT = '00' OR WS-TRAN-IN-STATUS NOT = '00' OR WS-TRAN-OUT-STATUS NOT = '00' OR NOT REJECTS-OPENED OR WS-REPORT-STATUS NOT = '00'` | Y | N | N |
| `READ CONTROL-FILE AT END` | | Y | N |
| → `OPEN INPUT CONTROL-FILE TRAN-IN OUTPUT TRAN-OUT REPORT-OUT EXTEND REJECT-OUT` | 1 | 1 | 1 |
| → `DISPLAY 'TXVALID: CANNOT OPEN FILES ' WS-FILE-STATUSES` | 2 | | |
| → `MOVE 16 TO RETURN-CODE` | 3 | 4 | |
| → `STOP RUN` | 4 | 5 | |
| → `READ CONTROL-FILE` | | 2 | 2 |
| → `DISPLAY 'TXVALID: NO BANK CONTROL RECORD'` | | 3 | |
| → `MOVE BC-BUSINESS-DATE TO WS-BUSINESS-DATE` | | | 3 |
| → `MOVE SPACES TO REJECT-RECORD` | | | 4 |
| → `STRING WS-BUSINESS-DATE(1:4) '-' WS-BUSINESS-DATE(5:2) '-' WS-BUSINESS-DATE(7:2) DELIMITED BY SIZE INTO RJ-REJECT-DATE` | | | 5 |

Source: [TXVALID.cbl:136](../../../pangaea/programs/TXVALID.cbl#L136).

## 2000-READ-HEADER

| | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 | 16 |
| --- | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: |
| `READ TRAN-IN AT END` | Y | Y | Y | Y | Y | Y | Y | Y | N | N | N | N | N | N | N | N |
| `FILE-IS-GOOD` | N | Y | Y | Y | Y | Y | Y | Y | N | Y | Y | Y | Y | Y | Y | Y |
| `TR-HEADER` | | N | N | Y | Y | Y | Y | Y | | N | N | Y | Y | Y | Y | Y |
| `FILE-IS-GOOD, again` | | Y | N | Y | N | Y | N | | | Y | N | Y | N | Y | N | |
| `TH-BUSINESS-DATE = WS-BUSINESS-DATE` | | | | N | N | Y | Y | Y | | | | N | N | Y | Y | Y |
| `NOT TH-FROM-BRANCHES AND NOT TH-FROM-CARDS` | | | | | | Y | Y | N | | | | | | Y | Y | N |
| → `READ TRAN-IN` | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 |
| → `MOVE 'EMPTY FILE' TO WS-FILE-ERROR` | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | | | | | | | | |
| → `CONTINUE` | 3 | | | | | | | | 2 | | | | | | | |
| → `MOVE 'FIRST RECORD IS NOT A HEADER' TO WS-FILE-ERROR` | | 3 | 3 | | | | | | | 2 | 2 | | | | | |
| → `MOVE TH-SOURCE TO WS-SOURCE` | | 4 | | 4 | | 4 | | 3 | | 3 | | 3 | | 3 | | 2 |
| → `MOVE 'HEADER IS NOT FOR THE BUSINESS DATE' TO WS-FILE-ERROR` | | | | 3 | 3 | | | | | | | 2 | 2 | | | |
| → `MOVE 'HEADER IS FROM AN UNKNOWN SOURCE' TO WS-FILE-ERROR` | | | | | | 3 | 3 | | | | | | | 2 | 2 | |

Source: [TXVALID.cbl:163](../../../pangaea/programs/TXVALID.cbl#L163).

## 3000-READ-DETAILS, in its loop

Each pass of `PERFORM UNTIL TR-TRAILER OR NOT FILE-IS-GOOD`.

| | 1 | 2 | 3 | 4 | 5 | 6 |
| --- | :-: | :-: | :-: | :-: | :-: | :-: |
| `TR-DETAIL` | Y | Y | N | N | N | N |
| `FILE-IS-GOOD` | Y | N | Y | N | Y | N |
| `TR-HEADER` | | | Y | Y | N | N |
| → `PERFORM 3200-VALIDATE-DETAIL` | 1 | 1 | | | | |
| → `PERFORM 3100-READ-NEXT` | 2 | | 2 | | 3 | |
| → `MOVE 'SECOND HEADER' TO WS-FILE-ERROR` | | | 1 | 1 | | |
| → `MOVE 'TYPE' TO RJ-REASON` | | | | | 1 | 1 |
| → `PERFORM 3300-REJECT-RECORD` | | | | | 2 | 2 |

Source: [TXVALID.cbl:187](../../../pangaea/programs/TXVALID.cbl#L187).

## 3100-READ-NEXT

| | 1 | 2 |
| --- | :-: | :-: |
| `READ TRAN-IN AT END` | Y | N |
| → `READ TRAN-IN` | 1 | 1 |
| → `MOVE 'NO TRAILER' TO WS-FILE-ERROR` | 2 | |

Source: [TXVALID.cbl:203](../../../pangaea/programs/TXVALID.cbl#L203).

## 3200-VALIDATE-DETAIL

| | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 | 16 | 17 | 18 | 19 | 20 |
| --- | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: |
| `TD-AMOUNT IS NUMERIC` | Y | Y | Y | Y | Y | Y | Y | Y | Y | Y | N | N | N | N | N | N | N | N | N | N |
| `TD-ACCT-NO IS NUMERIC` | Y | Y | Y | Y | Y | Y | Y | Y | N | N | Y | Y | Y | Y | Y | Y | Y | Y | N | N |
| `(FROM-BRANCHES AND NOT TD-KNOWN-CODE) OR (FROM-CARDS AND NOT TD-CARD-CODE)` | Y | Y | N | N | N | N | N | N | | | Y | Y | N | N | N | N | N | N | | |
| `DETAIL-IS-VALID` | Y | N | Y | N | Y | N | Y | N | Y | N | Y | N | Y | N | Y | N | Y | N | Y | N |
| `TD-AMOUNT IS NOT NUMERIC OR TD-AMOUNT = ZERO` | | | Y | Y | N | N | N | N | | | | | Y | Y | N | N | N | N | | |
| `TD-TRANSFER AND (TD-COUNTER-ACCT-NO IS NOT NUMERIC OR TD-COUNTER-ACCT-NO = TD-ACCT-NO)` | | | | | Y | Y | N | N | | | | | | | Y | Y | N | N | | |
| → `ADD 1 TO WS-DETAIL-COUNT` | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 |
| → `ADD TD-AMOUNT TO WS-AMOUNT-TOTAL` | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | | | | | | | | | | |
| → `ADD TD-ACCT-NO TO WS-ACCT-HASH` | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | | | 2 | 2 | 2 | 2 | 2 | 2 | 2 | 2 | | |
| → `SET DETAIL-IS-VALID TO TRUE` | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 4 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 3 | 2 | 2 |
| → `SET RJ-BAD-CODE TO TRUE` | 5 | 5 | | | | | | | | | 4 | 4 | | | | | | | | |
| → `SET DETAIL-IS-REJECTED TO TRUE` | 6 | 6 | 6 | 6 | 6 | 6 | | | 5 | 5 | 5 | 5 | 5 | 5 | 5 | 5 | | | 4 | 4 |
| → `WRITE TRAN-OUT-RECORD FROM TRAN-RECORD` | 7 | | 7 | | 7 | | 5 | | 6 | | 6 | | 6 | | 6 | | 4 | | 5 | |
| → `ADD 1 TO WS-VALID-COUNT` | 8 | | 8 | | 8 | | 6 | | 7 | | 7 | | 7 | | 7 | | 5 | | 6 | |
| → `PERFORM 3300-REJECT-RECORD` | | 7 | | 7 | | 7 | | 5 | | 6 | | 6 | | 6 | | 6 | | 4 | | 5 |
| → `SET RJ-BAD-AMOUNT TO TRUE` | | | 5 | 5 | | | | | | | | | 4 | 4 | | | | | | |
| → `SET RJ-BAD-COUNTER TO TRUE` | | | | | 5 | 5 | | | | | | | | | 4 | 4 | | | | |
| → `SET RJ-BAD-ACCOUNT TO TRUE` | | | | | | | | | 4 | 4 | | | | | | | | | 3 | 3 |

Source: [TXVALID.cbl:210](../../../pangaea/programs/TXVALID.cbl#L210).

## 4000-CHECK-TRAILER

| | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 |
| --- | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: |
| `TT-DETAIL-COUNT = WS-DETAIL-COUNT` | N | N | N | Y | Y | Y | Y | Y | Y | Y | Y | Y |
| `FILE-IS-GOOD` | Y | Y | N | Y | Y | N | Y | Y | N | Y | Y | N |
| `READ TRAN-IN AT END` | N | Y | | N | Y | | N | Y | | N | Y | |
| `TT-AMOUNT-TOTAL = WS-AMOUNT-TOTAL` | | | | N | N | N | Y | Y | Y | Y | Y | Y |
| `TT-ACCT-HASH = WS-ACCT-HASH` | | | | | | | N | N | N | Y | Y | Y |
| → `MOVE 'TRAILER DETAIL COUNT DIFFERS' TO WS-FILE-ERROR` | 1 | 1 | 1 | | | | | | | | | |
| → `READ TRAN-IN` | 2 | 2 | | 2 | 2 | | 2 | 2 | | 1 | 1 | |
| → `MOVE 'RECORDS AFTER THE TRAILER' TO WS-FILE-ERROR` | 3 | | | 3 | | | 3 | | | 2 | | |
| → `MOVE 'TRAILER AMOUNT TOTAL DIFFERS' TO WS-FILE-ERROR` | | | | 1 | 1 | 1 | | | | | | |
| → `MOVE 'TRAILER ACCOUNT HASH DIFFERS' TO WS-FILE-ERROR` | | | | | | | 1 | 1 | 1 | | | |

Source: [TXVALID.cbl:252](../../../pangaea/programs/TXVALID.cbl#L252).

## 8000-WRITE-REPORT

| | 1 | 2 |
| --- | :-: | :-: |
| `FILE-IS-GOOD` | Y | N |
| → `MOVE SPACES TO REPORT-LINE` | 1, 4, 18, 20 | 1, 4, 18, 20 |
| → `STRING 'TXVALID  TRANSACTION VALIDATION  BUSINESS DATE ' WS-BUSINESS-DATE DELIMITED BY SIZE INTO REPORT-LINE` | 2 | 2 |
| → `WRITE REPORT-LINE` | 3, 5, 19, 22 | 3, 5, 19, 22 |
| → `MOVE 'DETAILS READ' TO WS-RC-LABEL` | 6 | 6 |
| → `MOVE WS-DETAIL-COUNT TO WS-RC-VALUE` | 7 | 7 |
| → `WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE` | 8, 11, 14 | 8, 11, 14 |
| → `MOVE 'DETAILS VALID' TO WS-RC-LABEL` | 9 | 9 |
| → `MOVE WS-VALID-COUNT TO WS-RC-VALUE` | 10 | 10 |
| → `MOVE 'RECORDS REJECTED' TO WS-RC-LABEL` | 12 | 12 |
| → `MOVE WS-REJECT-COUNT TO WS-RC-VALUE` | 13 | 13 |
| → `MOVE 'AMOUNT TOTAL' TO WS-RA-LABEL` | 15 | 15 |
| → `MOVE WS-AMOUNT-TOTAL TO WS-RA-VALUE` | 16 | 16 |
| → `WRITE REPORT-LINE FROM WS-REPORT-AMOUNT-LINE` | 17 | 17 |
| → `MOVE 'FILE ACCEPTED' TO REPORT-LINE` | 21 | |
| → `STRING 'FILE REJECTED: ' WS-FILE-ERROR DELIMITED BY SIZE INTO REPORT-LINE` | | 21 |

Source: [TXVALID.cbl:273](../../../pangaea/programs/TXVALID.cbl#L273).
