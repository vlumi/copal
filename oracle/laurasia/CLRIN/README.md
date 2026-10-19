# CLRIN: cases

13 cases reach 22 of the 23 rules of CLRIN's decision tables. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile.

## [MAIN](../../../blueprints/laurasia/rules/CLRIN.md#main)

| Rule | Reached by |
| --: | --- |
| 1 | [case 04](04/) |
| 2 | [case 03](03/) |
| 3 | [case 01](01/) |
| 4 | not reached: `IO-FAILED` N, `CLEARING-AT-END OR NOT IS-FILE-HEADER OR FH-FILE-DATE NOT = WS-RUN-DATE` N, `IO-FAILED, again` Y |
| 5 | [case 02](02/) |

## [READ-CLEARING](../../../blueprints/laurasia/rules/CLRIN.md#read-clearing)

| Rule | Reached by |
| --: | --- |
| 1 | [case 05](05/) |
| 2 | [case 01](01/) |

## [TAKE-BATCH](../../../blueprints/laurasia/rules/CLRIN.md#take-batch)

| Rule | Reached by |
| --: | --- |
| 1 | [case 05](05/) |
| 2 | [case 02](02/) |
| 3 | [case 06](06/) |
| 4 | [case 09](09/) |

## [TAKE-PRESENTMENT](../../../blueprints/laurasia/rules/CLRIN.md#take-presentment)

| Rule | Reached by |
| --: | --- |
| 1 | [case 13](13/) |
| 2 | [case 02](02/) |

## [VALIDATE-PRESENTMENT](../../../blueprints/laurasia/rules/CLRIN.md#validate-presentment)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 10](10/) |
| 3 | [case 11](11/) |
| 4 | [case 07](07/) |
| 5 | [case 12](12/) |
| 6 | [case 13](13/) |

## [REJECT-RECORD](../../../blueprints/laurasia/rules/CLRIN.md#reject-record)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 05](05/) |

## [CHECK-FILE-TRAILER](../../../blueprints/laurasia/rules/CLRIN.md#check-file-trailer)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 08](08/) |
