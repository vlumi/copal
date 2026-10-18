# TXVALID: cases

19 cases reach 31 of the 66 rules of TXVALID's decision tables. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's.

## [0000-MAIN](../../../blueprints/pangaea/rules/TXVALID.md#0000-main)

| Rule | Reached by |
| --: | --- |
| 1 | [case 11](11/) |
| 2 | [case 05](05/) |
| 3 | [case 01](01/) |
| 4 | [case 02](02/) |
| 5 | [case 03](03/) |

## [1000-OPEN-FILES](../../../blueprints/pangaea/rules/TXVALID.md#1000-open-files)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-CONTROL-STATUS NOT = '00' OR WS-TRAN-IN-STATUS NOT = '00' OR WS-TRAN-OUT-STATUS NOT = '00' OR NOT REJECTS-OPENED OR WS-REPORT-STATUS NOT = '00'` Y |
| 2 | [case 13](13/) |
| 3 | [case 01](01/) |

## [2000-READ-HEADER](../../../blueprints/pangaea/rules/TXVALID.md#2000-read-header)

| Rule | Reached by |
| --: | --- |
| 1 | [case 14](14/) |
| 2 | not reached: `READ TRAN-IN AT END` Y, `FILE-IS-GOOD` Y, `TR-HEADER` N, `FILE-IS-GOOD, again` Y |
| 3 | not reached: `READ TRAN-IN AT END` Y, `FILE-IS-GOOD` Y, `TR-HEADER` N, `FILE-IS-GOOD, again` N |
| 4 | not reached: `READ TRAN-IN AT END` Y, `FILE-IS-GOOD` Y, `TR-HEADER` Y, `FILE-IS-GOOD, again` Y, `TH-BUSINESS-DATE = WS-BUSINESS-DATE` N |
| 5 | not reached: `READ TRAN-IN AT END` Y, `FILE-IS-GOOD` Y, `TR-HEADER` Y, `FILE-IS-GOOD, again` N, `TH-BUSINESS-DATE = WS-BUSINESS-DATE` N |
| 6 | not reached: `READ TRAN-IN AT END` Y, `FILE-IS-GOOD` Y, `TR-HEADER` Y, `FILE-IS-GOOD, again` Y, `TH-BUSINESS-DATE = WS-BUSINESS-DATE` Y, `NOT TH-FROM-BRANCHES AND NOT TH-FROM-CARDS` Y |
| 7 | not reached: `READ TRAN-IN AT END` Y, `FILE-IS-GOOD` Y, `TR-HEADER` Y, `FILE-IS-GOOD, again` N, `TH-BUSINESS-DATE = WS-BUSINESS-DATE` Y, `NOT TH-FROM-BRANCHES AND NOT TH-FROM-CARDS` Y |
| 8 | not reached: `READ TRAN-IN AT END` Y, `FILE-IS-GOOD` Y, `TR-HEADER` Y, `TH-BUSINESS-DATE = WS-BUSINESS-DATE` Y, `NOT TH-FROM-BRANCHES AND NOT TH-FROM-CARDS` N |
| 9 | not reached: `READ TRAN-IN AT END` N, `FILE-IS-GOOD` N |
| 10 | not reached: `READ TRAN-IN AT END` N, `FILE-IS-GOOD` Y, `TR-HEADER` N, `FILE-IS-GOOD, again` Y |
| 11 | [case 03](03/) |
| 12 | not reached: `READ TRAN-IN AT END` N, `FILE-IS-GOOD` Y, `TR-HEADER` Y, `FILE-IS-GOOD, again` Y, `TH-BUSINESS-DATE = WS-BUSINESS-DATE` N |
| 13 | [case 04](04/) |
| 14 | not reached: `READ TRAN-IN AT END` N, `FILE-IS-GOOD` Y, `TR-HEADER` Y, `FILE-IS-GOOD, again` Y, `TH-BUSINESS-DATE = WS-BUSINESS-DATE` Y, `NOT TH-FROM-BRANCHES AND NOT TH-FROM-CARDS` Y |
| 15 | [case 06](06/) |
| 16 | [case 01](01/) |

## [3000-READ-DETAILS, in its loop](../../../blueprints/pangaea/rules/TXVALID.md#3000-read-details-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | not reached: `TR-DETAIL` Y, `FILE-IS-GOOD` N |
| 3 | not reached: `TR-DETAIL` N, `FILE-IS-GOOD` Y, `TR-HEADER` Y |
| 4 | [case 02](02/) |
| 5 | not reached: `TR-DETAIL` N, `FILE-IS-GOOD` Y, `TR-HEADER` N |
| 6 | not reached: `TR-DETAIL` N, `FILE-IS-GOOD` N, `TR-HEADER` N |

## [3100-READ-NEXT](../../../blueprints/pangaea/rules/TXVALID.md#3100-read-next)

| Rule | Reached by |
| --: | --- |
| 1 | [case 17](17/) |
| 2 | [case 01](01/) |

## [3200-VALIDATE-DETAIL](../../../blueprints/pangaea/rules/TXVALID.md#3200-validate-detail)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `TD-AMOUNT IS NUMERIC` Y, `TD-ACCT-NO IS NUMERIC` Y, `(FROM-BRANCHES AND NOT TD-KNOWN-CODE) OR (FROM-CARDS AND NOT TD-CARD-CODE)` Y, `DETAIL-IS-VALID` Y |
| 2 | [case 05](05/) |
| 3 | not reached: `TD-AMOUNT IS NUMERIC` Y, `TD-ACCT-NO IS NUMERIC` Y, `(FROM-BRANCHES AND NOT TD-KNOWN-CODE) OR (FROM-CARDS AND NOT TD-CARD-CODE)` N, `DETAIL-IS-VALID` Y, `TD-AMOUNT IS NOT NUMERIC OR TD-AMOUNT = ZERO` Y |
| 4 | [case 08](08/) |
| 5 | not reached: `TD-AMOUNT IS NUMERIC` Y, `TD-ACCT-NO IS NUMERIC` Y, `(FROM-BRANCHES AND NOT TD-KNOWN-CODE) OR (FROM-CARDS AND NOT TD-CARD-CODE)` N, `DETAIL-IS-VALID` Y, `TD-AMOUNT IS NOT NUMERIC OR TD-AMOUNT = ZERO` N, `TD-TRANSFER AND (TD-COUNTER-ACCT-NO IS NOT NUMERIC OR TD-COUNTER-ACCT-NO = TD-ACCT-NO)` Y |
| 6 | [case 10](10/) |
| 7 | [case 01](01/) |
| 8 | not reached: `TD-AMOUNT IS NUMERIC` Y, `TD-ACCT-NO IS NUMERIC` Y, `(FROM-BRANCHES AND NOT TD-KNOWN-CODE) OR (FROM-CARDS AND NOT TD-CARD-CODE)` N, `DETAIL-IS-VALID` N, `TD-AMOUNT IS NOT NUMERIC OR TD-AMOUNT = ZERO` N, `TD-TRANSFER AND (TD-COUNTER-ACCT-NO IS NOT NUMERIC OR TD-COUNTER-ACCT-NO = TD-ACCT-NO)` N |
| 9 | not reached: `TD-AMOUNT IS NUMERIC` Y, `TD-ACCT-NO IS NUMERIC` N, `DETAIL-IS-VALID` Y |
| 10 | [case 09](09/) |
| 11 | not reached: `TD-AMOUNT IS NUMERIC` N, `TD-ACCT-NO IS NUMERIC` Y, `(FROM-BRANCHES AND NOT TD-KNOWN-CODE) OR (FROM-CARDS AND NOT TD-CARD-CODE)` Y, `DETAIL-IS-VALID` Y |
| 12 | [case 18](18/) |
| 13 | not reached: `TD-AMOUNT IS NUMERIC` N, `TD-ACCT-NO IS NUMERIC` Y, `(FROM-BRANCHES AND NOT TD-KNOWN-CODE) OR (FROM-CARDS AND NOT TD-CARD-CODE)` N, `DETAIL-IS-VALID` Y, `TD-AMOUNT IS NOT NUMERIC OR TD-AMOUNT = ZERO` Y |
| 14 | [case 07](07/) |
| 15 | not reached: `TD-AMOUNT IS NUMERIC` N, `TD-ACCT-NO IS NUMERIC` Y, `(FROM-BRANCHES AND NOT TD-KNOWN-CODE) OR (FROM-CARDS AND NOT TD-CARD-CODE)` N, `DETAIL-IS-VALID` Y, `TD-AMOUNT IS NOT NUMERIC OR TD-AMOUNT = ZERO` N, `TD-TRANSFER AND (TD-COUNTER-ACCT-NO IS NOT NUMERIC OR TD-COUNTER-ACCT-NO = TD-ACCT-NO)` Y |
| 16 | not reached: `TD-AMOUNT IS NUMERIC` N, `TD-ACCT-NO IS NUMERIC` Y, `(FROM-BRANCHES AND NOT TD-KNOWN-CODE) OR (FROM-CARDS AND NOT TD-CARD-CODE)` N, `DETAIL-IS-VALID` N, `TD-AMOUNT IS NOT NUMERIC OR TD-AMOUNT = ZERO` N, `TD-TRANSFER AND (TD-COUNTER-ACCT-NO IS NOT NUMERIC OR TD-COUNTER-ACCT-NO = TD-ACCT-NO)` Y |
| 17 | not reached: `TD-AMOUNT IS NUMERIC` N, `TD-ACCT-NO IS NUMERIC` Y, `(FROM-BRANCHES AND NOT TD-KNOWN-CODE) OR (FROM-CARDS AND NOT TD-CARD-CODE)` N, `DETAIL-IS-VALID` Y, `TD-AMOUNT IS NOT NUMERIC OR TD-AMOUNT = ZERO` N, `TD-TRANSFER AND (TD-COUNTER-ACCT-NO IS NOT NUMERIC OR TD-COUNTER-ACCT-NO = TD-ACCT-NO)` N |
| 18 | not reached: `TD-AMOUNT IS NUMERIC` N, `TD-ACCT-NO IS NUMERIC` Y, `(FROM-BRANCHES AND NOT TD-KNOWN-CODE) OR (FROM-CARDS AND NOT TD-CARD-CODE)` N, `DETAIL-IS-VALID` N, `TD-AMOUNT IS NOT NUMERIC OR TD-AMOUNT = ZERO` N, `TD-TRANSFER AND (TD-COUNTER-ACCT-NO IS NOT NUMERIC OR TD-COUNTER-ACCT-NO = TD-ACCT-NO)` N |
| 19 | not reached: `TD-AMOUNT IS NUMERIC` N, `TD-ACCT-NO IS NUMERIC` N, `DETAIL-IS-VALID` Y |
| 20 | [case 19](19/) |

## [4000-CHECK-TRAILER](../../../blueprints/pangaea/rules/TXVALID.md#4000-check-trailer)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `TT-DETAIL-COUNT = WS-DETAIL-COUNT` N, `FILE-IS-GOOD` Y, `READ TRAN-IN AT END` N |
| 2 | not reached: `TT-DETAIL-COUNT = WS-DETAIL-COUNT` N, `FILE-IS-GOOD` Y, `READ TRAN-IN AT END` Y |
| 3 | [case 16](16/) |
| 4 | not reached: `TT-DETAIL-COUNT = WS-DETAIL-COUNT` Y, `FILE-IS-GOOD` Y, `READ TRAN-IN AT END` N, `TT-AMOUNT-TOTAL = WS-AMOUNT-TOTAL` N |
| 5 | not reached: `TT-DETAIL-COUNT = WS-DETAIL-COUNT` Y, `FILE-IS-GOOD` Y, `READ TRAN-IN AT END` Y, `TT-AMOUNT-TOTAL = WS-AMOUNT-TOTAL` N |
| 6 | [case 11](11/) |
| 7 | not reached: `TT-DETAIL-COUNT = WS-DETAIL-COUNT` Y, `FILE-IS-GOOD` Y, `READ TRAN-IN AT END` N, `TT-AMOUNT-TOTAL = WS-AMOUNT-TOTAL` Y, `TT-ACCT-HASH = WS-ACCT-HASH` N |
| 8 | not reached: `TT-DETAIL-COUNT = WS-DETAIL-COUNT` Y, `FILE-IS-GOOD` Y, `READ TRAN-IN AT END` Y, `TT-AMOUNT-TOTAL = WS-AMOUNT-TOTAL` Y, `TT-ACCT-HASH = WS-ACCT-HASH` N |
| 9 | [case 12](12/) |
| 10 | [case 15](15/) |
| 11 | [case 01](01/) |
| 12 | not reached: `TT-DETAIL-COUNT = WS-DETAIL-COUNT` Y, `FILE-IS-GOOD` N, `TT-AMOUNT-TOTAL = WS-AMOUNT-TOTAL` Y, `TT-ACCT-HASH = WS-ACCT-HASH` Y |

## [8000-WRITE-REPORT](../../../blueprints/pangaea/rules/TXVALID.md#8000-write-report)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 02](02/) |
