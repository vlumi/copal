# TXVALID: cases

20 cases reach 32 of the 35 rules of TXVALID's decision tables. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it.

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
| 1 | [case 18](18/) |
| 2 | [case 13](13/) |
| 3 | [case 01](01/) |

## [2000-READ-HEADER](../../../blueprints/pangaea/rules/TXVALID.md#2000-read-header)

| Rule | Reached by |
| --: | --- |
| 1 | [case 14](14/) |
| 2 | not reached: `READ TRAN-IN AT END` N, `FILE-IS-GOOD` N |
| 3 | [case 03](03/) |
| 4 | [case 04](04/) |
| 5 | [case 06](06/) |
| 6 | [case 01](01/) |

## [3000-READ-DETAILS, in its loop](../../../blueprints/pangaea/rules/TXVALID.md#3000-read-details-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 02](02/) |
| 3 | not reached: `TR-DETAIL` N, `FILE-IS-GOOD` Y, `TR-HEADER` N |

## [3100-READ-NEXT](../../../blueprints/pangaea/rules/TXVALID.md#3100-read-next)

| Rule | Reached by |
| --: | --- |
| 1 | [case 17](17/) |
| 2 | [case 01](01/) |

## [3200-VALIDATE-DETAIL](../../../blueprints/pangaea/rules/TXVALID.md#3200-validate-detail)

| Rule | Reached by |
| --: | --- |
| 1 | [case 05](05/) |
| 2 | [case 08](08/) |
| 3 | [case 10](10/) |
| 4 | [case 01](01/) |
| 5 | [case 09](09/) |
| 6 | [case 19](19/) |
| 7 | [case 07](07/) |
| 8 | [case 20](20/) |

## [4000-CHECK-TRAILER](../../../blueprints/pangaea/rules/TXVALID.md#4000-check-trailer)

| Rule | Reached by |
| --: | --- |
| 1 | [case 16](16/) |
| 2 | [case 11](11/) |
| 3 | [case 12](12/) |
| 4 | [case 15](15/) |
| 5 | [case 01](01/) |
| 6 | not reached: `TT-DETAIL-COUNT = WS-DETAIL-COUNT` Y, `FILE-IS-GOOD` N, `TT-AMOUNT-TOTAL = WS-AMOUNT-TOTAL` Y, `TT-ACCT-HASH = WS-ACCT-HASH` Y |

## [8000-WRITE-REPORT](../../../blueprints/pangaea/rules/TXVALID.md#8000-write-report)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 02](02/) |
