# SETTLE: cases

2 cases reach 12 of the 12 rules of SETTLE's decision tables. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile.

## [MAIN](../../../blueprints/laurasia/rules/SETTLE.md#main)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 01](01/) |

## [READ-PRESENTMENT](../../../blueprints/laurasia/rules/SETTLE.md#read-presentment)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |

## [READ-ADJUSTMENT](../../../blueprints/laurasia/rules/SETTLE.md#read-adjustment)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |

## [SETTLE-MERCHANT](../../../blueprints/laurasia/rules/SETTLE.md#settle-merchant)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |

## [SETTLE-MERCHANT, in its loop](../../../blueprints/laurasia/rules/SETTLE.md#settle-merchant-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |

## [WRITE-PAYOUT](../../../blueprints/laurasia/rules/SETTLE.md#write-payout)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 01](01/) |
