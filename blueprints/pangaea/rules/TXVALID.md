# TXVALID: decision tables

The business rules of [TXVALID](../programs/TXVALID.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

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

- **Rule 1.** When `WS-FILE-ERROR` is spaces (`FILE-IS-GOOD`); `WS-FILE-ERROR` is spaces (`FILE-IS-GOOD`), tested again; and `WS-FILE-ERROR` is not spaces (not `FILE-IS-GOOD`), tested again, it performs `1000-OPEN-FILES`, performs `2000-READ-HEADER`, performs `3000-READ-DETAILS`, performs `4000-CHECK-TRAILER`, performs `8000-WRITE-REPORT`, performs `9000-CLOSE-FILES`, sets `RETURN-CODE` to 12, and stops the run.
- **Rule 2.** When `WS-FILE-ERROR` is spaces (`FILE-IS-GOOD`); `WS-FILE-ERROR` is spaces (`FILE-IS-GOOD`), tested again; `WS-FILE-ERROR` is spaces (`FILE-IS-GOOD`), tested again; and `WS-REJECT-COUNT` is more than 0, it performs `1000-OPEN-FILES`, performs `2000-READ-HEADER`, performs `3000-READ-DETAILS`, performs `4000-CHECK-TRAILER`, performs `8000-WRITE-REPORT`, performs `9000-CLOSE-FILES`, sets `RETURN-CODE` to 4, and stops the run.
- **Rule 3.** When `WS-FILE-ERROR` is spaces (`FILE-IS-GOOD`); `WS-FILE-ERROR` is spaces (`FILE-IS-GOOD`), tested again; `WS-FILE-ERROR` is spaces (`FILE-IS-GOOD`), tested again; and `WS-REJECT-COUNT` is at most 0, it performs `1000-OPEN-FILES`, performs `2000-READ-HEADER`, performs `3000-READ-DETAILS`, performs `4000-CHECK-TRAILER`, performs `8000-WRITE-REPORT`, performs `9000-CLOSE-FILES`, sets `RETURN-CODE` to 0, and stops the run.
- **Rule 4.** When `WS-FILE-ERROR` is spaces (`FILE-IS-GOOD`); and `WS-FILE-ERROR` is not spaces (not `FILE-IS-GOOD`), tested again, it performs `1000-OPEN-FILES`, performs `2000-READ-HEADER`, performs `3000-READ-DETAILS`, performs `8000-WRITE-REPORT`, performs `9000-CLOSE-FILES`, sets `RETURN-CODE` to 12, and stops the run.
- **Rule 5.** When `WS-FILE-ERROR` is not spaces (not `FILE-IS-GOOD`), it performs `1000-OPEN-FILES`, performs `2000-READ-HEADER`, performs `8000-WRITE-REPORT`, performs `9000-CLOSE-FILES`, sets `RETURN-CODE` to 12, and stops the run.

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

- **Rule 1.** When `WS-CONTROL-STATUS` is not `'00'`, `WS-TRAN-IN-STATUS` is not `'00'`, `WS-TRAN-OUT-STATUS` is not `'00'`, `WS-REJECT-STATUS` is neither `'00'` nor `'05'` (not `REJECTS-OPENED`), or `WS-REPORT-STATUS` is not `'00'`, it opens `CONTROL-FILE` and `TRAN-IN` for input, `TRAN-OUT` and `REPORT-OUT` for output, and `REJECT-OUT` to extend; displays `'TXVALID: CANNOT OPEN FILES '` `WS-FILE-STATUSES`; sets `RETURN-CODE` to 16; and stops the run.
- **Rule 2.** When `WS-CONTROL-STATUS` is `'00'`, `WS-TRAN-IN-STATUS` is `'00'`, `WS-TRAN-OUT-STATUS` is `'00'`, `WS-REJECT-STATUS` is `'00'` or `'05'` (`REJECTS-OPENED`), and `WS-REPORT-STATUS` is `'00'`; and reading `CONTROL-FILE` finds no more records, it opens `CONTROL-FILE` and `TRAN-IN` for input, `TRAN-OUT` and `REPORT-OUT` for output, and `REJECT-OUT` to extend; reads `CONTROL-FILE`; displays `'TXVALID: NO BANK CONTROL RECORD'`; sets `RETURN-CODE` to 16; and stops the run.
- **Rule 3.** When `WS-CONTROL-STATUS` is `'00'`, `WS-TRAN-IN-STATUS` is `'00'`, `WS-TRAN-OUT-STATUS` is `'00'`, `WS-REJECT-STATUS` is `'00'` or `'05'` (`REJECTS-OPENED`), and `WS-REPORT-STATUS` is `'00'`; and reading `CONTROL-FILE` finds a record, it opens `CONTROL-FILE` and `TRAN-IN` for input, `TRAN-OUT` and `REPORT-OUT` for output, and `REJECT-OUT` to extend; reads `CONTROL-FILE`; sets `WS-BUSINESS-DATE` to `BC-BUSINESS-DATE`; sets `REJECT-RECORD` to spaces; and builds `RJ-REJECT-DATE` from `WS-BUSINESS-DATE(1:4)` `'-'` `WS-BUSINESS-DATE(5:2)` `'-'` `WS-BUSINESS-DATE(7:2)`.

Source: [TXVALID.cbl:136](../../../pangaea/programs/TXVALID.cbl#L136).

## 2000-READ-HEADER

| | 1 | 2 | 3 | 4 | 5 |
| --- | :-: | :-: | :-: | :-: | :-: |
| `READ TRAN-IN AT END` | Y | N | N | N | N |
| `FILE-IS-GOOD` | N | Y | Y | Y | Y |
| `TR-HEADER` | | N | Y | Y | Y |
| `FILE-IS-GOOD, again` | | N | N | N | |
| `TH-BUSINESS-DATE = WS-BUSINESS-DATE` | | | N | Y | Y |
| `NOT TH-FROM-BRANCHES AND NOT TH-FROM-CARDS` | | | | Y | N |
| → `READ TRAN-IN` | 1 | 1 | 1 | 1 | 1 |
| → `MOVE 'EMPTY FILE' TO WS-FILE-ERROR` | 2 | | | | |
| → `CONTINUE` | 3 | | | | |
| → `MOVE 'FIRST RECORD IS NOT A HEADER' TO WS-FILE-ERROR` | | 2 | | | |
| → `MOVE 'HEADER IS NOT FOR THE BUSINESS DATE' TO WS-FILE-ERROR` | | | 2 | | |
| → `MOVE 'HEADER IS FROM AN UNKNOWN SOURCE' TO WS-FILE-ERROR` | | | | 2 | |
| → `MOVE TH-SOURCE TO WS-SOURCE` | | | | | 2 |

- **Rule 1.** When reading `TRAN-IN` finds no more records and `WS-FILE-ERROR` is not spaces (not `FILE-IS-GOOD`), it reads `TRAN-IN` and sets `WS-FILE-ERROR` to `'EMPTY FILE'`.
- **Rule 2.** When reading `TRAN-IN` finds a record; `WS-FILE-ERROR` is spaces (`FILE-IS-GOOD`); `TR-REC-TYPE` is not `'H'` (not `TR-HEADER`); and `WS-FILE-ERROR` is not spaces (not `FILE-IS-GOOD`), tested again, it reads `TRAN-IN` and sets `WS-FILE-ERROR` to `'FIRST RECORD IS NOT A HEADER'`.
- **Rule 3.** When reading `TRAN-IN` finds a record; `WS-FILE-ERROR` is spaces (`FILE-IS-GOOD`); `TR-REC-TYPE` is `'H'` (`TR-HEADER`); `WS-FILE-ERROR` is not spaces (not `FILE-IS-GOOD`), tested again; and `TH-BUSINESS-DATE` is not `WS-BUSINESS-DATE`, it reads `TRAN-IN` and sets `WS-FILE-ERROR` to `'HEADER IS NOT FOR THE BUSINESS DATE'`.
- **Rule 4.** When reading `TRAN-IN` finds a record; `WS-FILE-ERROR` is spaces (`FILE-IS-GOOD`); `TR-REC-TYPE` is `'H'` (`TR-HEADER`); `WS-FILE-ERROR` is not spaces (not `FILE-IS-GOOD`), tested again; `TH-BUSINESS-DATE` is `WS-BUSINESS-DATE`; and `TH-SOURCE` is not `'BRANCHES'` (not `TH-FROM-BRANCHES`) and `TH-SOURCE` is not `'CARDS'` (not `TH-FROM-CARDS`), it reads `TRAN-IN` and sets `WS-FILE-ERROR` to `'HEADER IS FROM AN UNKNOWN SOURCE'`.
- **Rule 5.** When reading `TRAN-IN` finds a record; `WS-FILE-ERROR` is spaces (`FILE-IS-GOOD`); `TR-REC-TYPE` is `'H'` (`TR-HEADER`); `TH-BUSINESS-DATE` is `WS-BUSINESS-DATE`; and `TH-SOURCE` is `'BRANCHES'` (`TH-FROM-BRANCHES`) or `TH-SOURCE` is `'CARDS'` (`TH-FROM-CARDS`), it reads `TRAN-IN` and sets `WS-SOURCE` to `TH-SOURCE`.

Source: [TXVALID.cbl:163](../../../pangaea/programs/TXVALID.cbl#L163).

## 3000-READ-DETAILS, in its loop

Each pass of `PERFORM UNTIL TR-TRAILER OR NOT FILE-IS-GOOD`.

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `TR-DETAIL` | Y | N | N |
| `FILE-IS-GOOD` | Y | N | Y |
| `TR-HEADER` | | Y | N |
| → `PERFORM 3200-VALIDATE-DETAIL` | 1 | | |
| → `PERFORM 3100-READ-NEXT` | 2 | | 3 |
| → `MOVE 'SECOND HEADER' TO WS-FILE-ERROR` | | 1 | |
| → `MOVE 'TYPE' TO RJ-REASON` | | | 1 |
| → `PERFORM 3300-REJECT-RECORD` | | | 2 |

- **Rule 1.** When `TR-REC-TYPE` is `'D'` (`TR-DETAIL`) and `WS-FILE-ERROR` is spaces (`FILE-IS-GOOD`), it performs `3200-VALIDATE-DETAIL` and performs `3100-READ-NEXT`.
- **Rule 2.** When `TR-REC-TYPE` is not `'D'` (not `TR-DETAIL`), `WS-FILE-ERROR` is not spaces (not `FILE-IS-GOOD`), and `TR-REC-TYPE` is `'H'` (`TR-HEADER`), it sets `WS-FILE-ERROR` to `'SECOND HEADER'`.
- **Rule 3.** When `TR-REC-TYPE` is not `'D'` (not `TR-DETAIL`), `WS-FILE-ERROR` is spaces (`FILE-IS-GOOD`), and `TR-REC-TYPE` is not `'H'` (not `TR-HEADER`), it sets `RJ-REASON` to `'TYPE'`, performs `3300-REJECT-RECORD`, and performs `3100-READ-NEXT`.

Source: [TXVALID.cbl:187](../../../pangaea/programs/TXVALID.cbl#L187).

## 3100-READ-NEXT

| | 1 | 2 |
| --- | :-: | :-: |
| `READ TRAN-IN AT END` | Y | N |
| → `READ TRAN-IN` | 1 | 1 |
| → `MOVE 'NO TRAILER' TO WS-FILE-ERROR` | 2 | |

- **Rule 1.** When reading `TRAN-IN` finds no more records, it reads `TRAN-IN` and sets `WS-FILE-ERROR` to `'NO TRAILER'`.
- **Rule 2.** When reading `TRAN-IN` finds a record, it reads `TRAN-IN`.

Source: [TXVALID.cbl:203](../../../pangaea/programs/TXVALID.cbl#L203).

## 3200-VALIDATE-DETAIL

| | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 |
| --- | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: |
| `TD-AMOUNT IS NUMERIC` | Y | Y | Y | Y | Y | N | N | N |
| `TD-ACCT-NO IS NUMERIC` | Y | Y | Y | Y | N | Y | Y | N |
| `(FROM-BRANCHES AND NOT TD-KNOWN-CODE) OR (FROM-CARDS AND NOT TD-CARD-CODE)` | Y | N | N | N | | Y | N | |
| `DETAIL-IS-VALID` | N | N | N | Y | N | N | N | N |
| `TD-AMOUNT IS NOT NUMERIC OR TD-AMOUNT = ZERO` | | Y | N | N | | | Y | |
| `TD-TRANSFER AND (TD-COUNTER-ACCT-NO IS NOT NUMERIC OR TD-COUNTER-ACCT-NO = TD-ACCT-NO)` | | | Y | N | | | | |
| → `ADD 1 TO WS-DETAIL-COUNT` | 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 |
| → `ADD TD-AMOUNT TO WS-AMOUNT-TOTAL` | 2 | 2 | 2 | 2 | 2 | | | |
| → `ADD TD-ACCT-NO TO WS-ACCT-HASH` | 3 | 3 | 3 | 3 | | 2 | 2 | |
| → `SET DETAIL-IS-VALID TO TRUE` | 4 | 4 | 4 | 4 | 3 | 3 | 3 | 2 |
| → `SET RJ-BAD-CODE TO TRUE` | 5 | | | | | 4 | | |
| → `SET DETAIL-IS-REJECTED TO TRUE` | 6 | 6 | 6 | | 5 | 5 | 5 | 4 |
| → `PERFORM 3300-REJECT-RECORD` | 7 | 7 | 7 | | 6 | 6 | 6 | 5 |
| → `SET RJ-BAD-AMOUNT TO TRUE` | | 5 | | | | | 4 | |
| → `SET RJ-BAD-COUNTER TO TRUE` | | | 5 | | | | | |
| → `WRITE TRAN-OUT-RECORD FROM TRAN-RECORD` | | | | 5 | | | | |
| → `ADD 1 TO WS-VALID-COUNT` | | | | 6 | | | | |
| → `SET RJ-BAD-ACCOUNT TO TRUE` | | | | | 4 | | | 3 |

- **Rule 1.** When `TD-AMOUNT` is numeric; `TD-ACCT-NO` is numeric; both `WS-SOURCE` is `'BRANCHES'` (`FROM-BRANCHES`) and `TD-TRAN-CODE` is none of `'DP'`, `'WD'`, `'FE'`, and `'TR'` (not `TD-KNOWN-CODE`) or both `WS-SOURCE` is `'CARDS'` (`FROM-CARDS`) and `TD-TRAN-CODE` is none of `'CP'`, `'AW'`, and `'AF'` (not `TD-CARD-CODE`); and `WS-DETAIL-STATE` is not `'Y'` (not `DETAIL-IS-VALID`), it adds 1 to `WS-DETAIL-COUNT`, adds `TD-AMOUNT` to `WS-AMOUNT-TOTAL`, adds `TD-ACCT-NO` to `WS-ACCT-HASH`, sets `WS-DETAIL-STATE` to `DETAIL-IS-VALID`, sets `RJ-REASON` to `RJ-BAD-CODE`, sets `WS-DETAIL-STATE` to `DETAIL-IS-REJECTED`, and performs `3300-REJECT-RECORD`.
- **Rule 2.** When `TD-AMOUNT` is numeric; `TD-ACCT-NO` is numeric; either `WS-SOURCE` is not `'BRANCHES'` (not `FROM-BRANCHES`) or `TD-TRAN-CODE` is `'DP'`, `'WD'`, `'FE'`, or `'TR'` (`TD-KNOWN-CODE`) and either `WS-SOURCE` is not `'CARDS'` (not `FROM-CARDS`) or `TD-TRAN-CODE` is `'CP'`, `'AW'`, or `'AF'` (`TD-CARD-CODE`); `WS-DETAIL-STATE` is not `'Y'` (not `DETAIL-IS-VALID`); and `TD-AMOUNT` is not numeric or `TD-AMOUNT` is zero, it adds 1 to `WS-DETAIL-COUNT`, adds `TD-AMOUNT` to `WS-AMOUNT-TOTAL`, adds `TD-ACCT-NO` to `WS-ACCT-HASH`, sets `WS-DETAIL-STATE` to `DETAIL-IS-VALID`, sets `RJ-REASON` to `RJ-BAD-AMOUNT`, sets `WS-DETAIL-STATE` to `DETAIL-IS-REJECTED`, and performs `3300-REJECT-RECORD`.
- **Rule 3.** When `TD-AMOUNT` is numeric; `TD-ACCT-NO` is numeric; either `WS-SOURCE` is not `'BRANCHES'` (not `FROM-BRANCHES`) or `TD-TRAN-CODE` is `'DP'`, `'WD'`, `'FE'`, or `'TR'` (`TD-KNOWN-CODE`) and either `WS-SOURCE` is not `'CARDS'` (not `FROM-CARDS`) or `TD-TRAN-CODE` is `'CP'`, `'AW'`, or `'AF'` (`TD-CARD-CODE`); `WS-DETAIL-STATE` is not `'Y'` (not `DETAIL-IS-VALID`); `TD-AMOUNT` is numeric and `TD-AMOUNT` is not zero; and `TD-TRAN-CODE` is `'TR'` (`TD-TRANSFER`) and either `TD-COUNTER-ACCT-NO` is not numeric or `TD-COUNTER-ACCT-NO` is `TD-ACCT-NO`, it adds 1 to `WS-DETAIL-COUNT`, adds `TD-AMOUNT` to `WS-AMOUNT-TOTAL`, adds `TD-ACCT-NO` to `WS-ACCT-HASH`, sets `WS-DETAIL-STATE` to `DETAIL-IS-VALID`, sets `RJ-REASON` to `RJ-BAD-COUNTER`, sets `WS-DETAIL-STATE` to `DETAIL-IS-REJECTED`, and performs `3300-REJECT-RECORD`.
- **Rule 4.** When `TD-AMOUNT` is numeric; `TD-ACCT-NO` is numeric; either `WS-SOURCE` is not `'BRANCHES'` (not `FROM-BRANCHES`) or `TD-TRAN-CODE` is `'DP'`, `'WD'`, `'FE'`, or `'TR'` (`TD-KNOWN-CODE`) and either `WS-SOURCE` is not `'CARDS'` (not `FROM-CARDS`) or `TD-TRAN-CODE` is `'CP'`, `'AW'`, or `'AF'` (`TD-CARD-CODE`); `WS-DETAIL-STATE` is `'Y'` (`DETAIL-IS-VALID`); `TD-AMOUNT` is numeric and `TD-AMOUNT` is not zero; and `TD-TRAN-CODE` is not `'TR'` (not `TD-TRANSFER`) or both `TD-COUNTER-ACCT-NO` is numeric and `TD-COUNTER-ACCT-NO` is not `TD-ACCT-NO`, it adds 1 to `WS-DETAIL-COUNT`, adds `TD-AMOUNT` to `WS-AMOUNT-TOTAL`, adds `TD-ACCT-NO` to `WS-ACCT-HASH`, sets `WS-DETAIL-STATE` to `DETAIL-IS-VALID`, writes `TRAN-OUT-RECORD` from `TRAN-RECORD`, and adds 1 to `WS-VALID-COUNT`.
- **Rule 5.** When `TD-AMOUNT` is numeric, `TD-ACCT-NO` is not numeric, and `WS-DETAIL-STATE` is not `'Y'` (not `DETAIL-IS-VALID`), it adds 1 to `WS-DETAIL-COUNT`, adds `TD-AMOUNT` to `WS-AMOUNT-TOTAL`, sets `WS-DETAIL-STATE` to `DETAIL-IS-VALID`, sets `RJ-REASON` to `RJ-BAD-ACCOUNT`, sets `WS-DETAIL-STATE` to `DETAIL-IS-REJECTED`, and performs `3300-REJECT-RECORD`.
- **Rule 6.** When `TD-AMOUNT` is not numeric; `TD-ACCT-NO` is numeric; both `WS-SOURCE` is `'BRANCHES'` (`FROM-BRANCHES`) and `TD-TRAN-CODE` is none of `'DP'`, `'WD'`, `'FE'`, and `'TR'` (not `TD-KNOWN-CODE`) or both `WS-SOURCE` is `'CARDS'` (`FROM-CARDS`) and `TD-TRAN-CODE` is none of `'CP'`, `'AW'`, and `'AF'` (not `TD-CARD-CODE`); and `WS-DETAIL-STATE` is not `'Y'` (not `DETAIL-IS-VALID`), it adds 1 to `WS-DETAIL-COUNT`, adds `TD-ACCT-NO` to `WS-ACCT-HASH`, sets `WS-DETAIL-STATE` to `DETAIL-IS-VALID`, sets `RJ-REASON` to `RJ-BAD-CODE`, sets `WS-DETAIL-STATE` to `DETAIL-IS-REJECTED`, and performs `3300-REJECT-RECORD`.
- **Rule 7.** When `TD-AMOUNT` is not numeric; `TD-ACCT-NO` is numeric; either `WS-SOURCE` is not `'BRANCHES'` (not `FROM-BRANCHES`) or `TD-TRAN-CODE` is `'DP'`, `'WD'`, `'FE'`, or `'TR'` (`TD-KNOWN-CODE`) and either `WS-SOURCE` is not `'CARDS'` (not `FROM-CARDS`) or `TD-TRAN-CODE` is `'CP'`, `'AW'`, or `'AF'` (`TD-CARD-CODE`); `WS-DETAIL-STATE` is not `'Y'` (not `DETAIL-IS-VALID`); and `TD-AMOUNT` is not numeric or `TD-AMOUNT` is zero, it adds 1 to `WS-DETAIL-COUNT`, adds `TD-ACCT-NO` to `WS-ACCT-HASH`, sets `WS-DETAIL-STATE` to `DETAIL-IS-VALID`, sets `RJ-REASON` to `RJ-BAD-AMOUNT`, sets `WS-DETAIL-STATE` to `DETAIL-IS-REJECTED`, and performs `3300-REJECT-RECORD`.
- **Rule 8.** When `TD-AMOUNT` is not numeric, `TD-ACCT-NO` is not numeric, and `WS-DETAIL-STATE` is not `'Y'` (not `DETAIL-IS-VALID`), it adds 1 to `WS-DETAIL-COUNT`, sets `WS-DETAIL-STATE` to `DETAIL-IS-VALID`, sets `RJ-REASON` to `RJ-BAD-ACCOUNT`, sets `WS-DETAIL-STATE` to `DETAIL-IS-REJECTED`, and performs `3300-REJECT-RECORD`.

Source: [TXVALID.cbl:210](../../../pangaea/programs/TXVALID.cbl#L210).

## 4000-CHECK-TRAILER

| | 1 | 2 | 3 | 4 | 5 |
| --- | :-: | :-: | :-: | :-: | :-: |
| `TT-DETAIL-COUNT = WS-DETAIL-COUNT` | N | Y | Y | Y | Y |
| `FILE-IS-GOOD` | N | N | N | Y | Y |
| `TT-AMOUNT-TOTAL = WS-AMOUNT-TOTAL` | | N | Y | Y | Y |
| `TT-ACCT-HASH = WS-ACCT-HASH` | | | N | Y | Y |
| `READ TRAN-IN AT END` | | | | N | Y |
| → `MOVE 'TRAILER DETAIL COUNT DIFFERS' TO WS-FILE-ERROR` | 1 | | | | |
| → `MOVE 'TRAILER AMOUNT TOTAL DIFFERS' TO WS-FILE-ERROR` | | 1 | | | |
| → `MOVE 'TRAILER ACCOUNT HASH DIFFERS' TO WS-FILE-ERROR` | | | 1 | | |
| → `READ TRAN-IN` | | | | 1 | 1 |
| → `MOVE 'RECORDS AFTER THE TRAILER' TO WS-FILE-ERROR` | | | | 2 | |

- **Rule 1.** When `TT-DETAIL-COUNT` is not `WS-DETAIL-COUNT` and `WS-FILE-ERROR` is not spaces (not `FILE-IS-GOOD`), it sets `WS-FILE-ERROR` to `'TRAILER DETAIL COUNT DIFFERS'`.
- **Rule 2.** When `TT-DETAIL-COUNT` is `WS-DETAIL-COUNT`, `WS-FILE-ERROR` is not spaces (not `FILE-IS-GOOD`), and `TT-AMOUNT-TOTAL` is not `WS-AMOUNT-TOTAL`, it sets `WS-FILE-ERROR` to `'TRAILER AMOUNT TOTAL DIFFERS'`.
- **Rule 3.** When `TT-DETAIL-COUNT` is `WS-DETAIL-COUNT`, `WS-FILE-ERROR` is not spaces (not `FILE-IS-GOOD`), `TT-AMOUNT-TOTAL` is `WS-AMOUNT-TOTAL`, and `TT-ACCT-HASH` is not `WS-ACCT-HASH`, it sets `WS-FILE-ERROR` to `'TRAILER ACCOUNT HASH DIFFERS'`.
- **Rule 4.** When `TT-DETAIL-COUNT` is `WS-DETAIL-COUNT`, `WS-FILE-ERROR` is spaces (`FILE-IS-GOOD`), `TT-AMOUNT-TOTAL` is `WS-AMOUNT-TOTAL`, `TT-ACCT-HASH` is `WS-ACCT-HASH`, and reading `TRAN-IN` finds a record, it reads `TRAN-IN` and sets `WS-FILE-ERROR` to `'RECORDS AFTER THE TRAILER'`.
- **Rule 5.** When `TT-DETAIL-COUNT` is `WS-DETAIL-COUNT`, `WS-FILE-ERROR` is spaces (`FILE-IS-GOOD`), `TT-AMOUNT-TOTAL` is `WS-AMOUNT-TOTAL`, `TT-ACCT-HASH` is `WS-ACCT-HASH`, and reading `TRAN-IN` finds no more records, it reads `TRAN-IN`.

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

- **Rule 1.** When `WS-FILE-ERROR` is spaces (`FILE-IS-GOOD`), it sets `REPORT-LINE` to spaces, builds `REPORT-LINE` from `'TXVALID  TRANSACTION VALIDATION  BUSINESS DATE '` `WS-BUSINESS-DATE`, writes `REPORT-LINE`, sets `REPORT-LINE` to spaces, writes `REPORT-LINE`, sets `WS-RC-LABEL` to `'DETAILS READ'`, sets `WS-RC-VALUE` to `WS-DETAIL-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `WS-RC-LABEL` to `'DETAILS VALID'`, sets `WS-RC-VALUE` to `WS-VALID-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `WS-RC-LABEL` to `'RECORDS REJECTED'`, sets `WS-RC-VALUE` to `WS-REJECT-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `WS-RA-LABEL` to `'AMOUNT TOTAL'`, sets `WS-RA-VALUE` to `WS-AMOUNT-TOTAL`, writes `REPORT-LINE` from `WS-REPORT-AMOUNT-LINE`, sets `REPORT-LINE` to spaces, writes `REPORT-LINE`, sets `REPORT-LINE` to spaces, sets `REPORT-LINE` to `'FILE ACCEPTED'`, and writes `REPORT-LINE`.
- **Rule 2.** When `WS-FILE-ERROR` is not spaces (not `FILE-IS-GOOD`), it sets `REPORT-LINE` to spaces, builds `REPORT-LINE` from `'TXVALID  TRANSACTION VALIDATION  BUSINESS DATE '` `WS-BUSINESS-DATE`, writes `REPORT-LINE`, sets `REPORT-LINE` to spaces, writes `REPORT-LINE`, sets `WS-RC-LABEL` to `'DETAILS READ'`, sets `WS-RC-VALUE` to `WS-DETAIL-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `WS-RC-LABEL` to `'DETAILS VALID'`, sets `WS-RC-VALUE` to `WS-VALID-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `WS-RC-LABEL` to `'RECORDS REJECTED'`, sets `WS-RC-VALUE` to `WS-REJECT-COUNT`, writes `REPORT-LINE` from `WS-REPORT-COUNT-LINE`, sets `WS-RA-LABEL` to `'AMOUNT TOTAL'`, sets `WS-RA-VALUE` to `WS-AMOUNT-TOTAL`, writes `REPORT-LINE` from `WS-REPORT-AMOUNT-LINE`, sets `REPORT-LINE` to spaces, writes `REPORT-LINE`, sets `REPORT-LINE` to spaces, builds `REPORT-LINE` from `'FILE REJECTED: '` `WS-FILE-ERROR`, and writes `REPORT-LINE`.

Source: [TXVALID.cbl:273](../../../pangaea/programs/TXVALID.cbl#L273).
