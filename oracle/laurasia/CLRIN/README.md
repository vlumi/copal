# CLRIN: cases

11 cases reach 20 of the 24 rules of CLRIN's decision tables. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's.

## [MAIN](../../../blueprints/laurasia/rules/CLRIN.md#main)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `IO-FAILED` Y |
| 2 | not reached: `IO-FAILED` N, `CLEARING-AT-END OR NOT IS-FILE-HEADER OR FH-FILE-DATE NOT = WS-RUN-DATE` Y, `IO-FAILED, again` Y |
| 3 | [case 01](01/) |
| 4 | not reached: `IO-FAILED` N, `CLEARING-AT-END OR NOT IS-FILE-HEADER OR FH-FILE-DATE NOT = WS-RUN-DATE` N, `IO-FAILED, again` Y |
| 5 | [case 02](02/) |

## [READ-CLEARING](../../../blueprints/laurasia/rules/CLRIN.md#read-clearing)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | [case 01](01/) |

## [TAKE-BATCH](../../../blueprints/laurasia/rules/CLRIN.md#take-batch)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | [case 02](02/) |
| 3 | not reached: `IS-BATCH-HEADER` Y, `IS-BATCH-TRAILER AND BT-COUNT = WS-BATCH-COUNT AND BT-AMOUNT = WS-BATCH-AMOUNT` Y, `IS-BATCH-TRAILER` N |
| 4 | [case 04](04/) |
| 5 | [case 07](07/) |

## [TAKE-PRESENTMENT](../../../blueprints/laurasia/rules/CLRIN.md#take-presentment)

| Rule | Reached by |
| --: | --- |
| 1 | [case 11](11/) |
| 2 | [case 02](02/) |

## [VALIDATE-PRESENTMENT](../../../blueprints/laurasia/rules/CLRIN.md#validate-presentment)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 08](08/) |
| 3 | [case 09](09/) |
| 4 | [case 05](05/) |
| 5 | [case 10](10/) |
| 6 | [case 11](11/) |

## [REJECT-RECORD](../../../blueprints/laurasia/rules/CLRIN.md#reject-record)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 03](03/) |

## [CHECK-FILE-TRAILER](../../../blueprints/laurasia/rules/CLRIN.md#check-file-trailer)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 06](06/) |
