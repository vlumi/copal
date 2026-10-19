# CARDLOAD: cases

5 cases reach 10 of the 10 rules of CARDLOAD's decision tables. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it.

## [1000-INITIALIZE](../../../blueprints/pangaea/rules/CARDLOAD.md#1000-initialize)

| Rule | Reached by |
| --: | --- |
| 1 | [case 05](05/) |
| 2 | [case 03](03/) |
| 3 | [case 01](01/) |

## [1200-READ-CARD](../../../blueprints/pangaea/rules/CARDLOAD.md#1200-read-card)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |

## [2000-LOAD-CARD](../../../blueprints/pangaea/rules/CARDLOAD.md#2000-load-card)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 01](01/) |
| 3 | [case 01](01/) |

## [2100-WRITE-ONLINE](../../../blueprints/pangaea/rules/CARDLOAD.md#2100-write-online)

| Rule | Reached by |
| --: | --- |
| 1 | [case 04](04/) |
| 2 | [case 01](01/) |
