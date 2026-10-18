# HOLDEXT: cases

3 cases reach 8 of the 15 rules of HOLDEXT's decision tables. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's.

## [1000-INITIALIZE](../../../blueprints/pangaea/rules/HOLDEXT.md#1000-initialize)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-FILE-STATUSES = ALL '0'` N |
| 2 | not reached: `WS-FILE-STATUSES = ALL '0'` Y, `READ CONTROL-FILE AT END` Y |
| 3 | [case 01](01/) |

## [1200-READ-HOLD](../../../blueprints/pangaea/rules/HOLDEXT.md#1200-read-hold)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |

## [2000-EXTRACT-HOLD](../../../blueprints/pangaea/rules/HOLDEXT.md#2000-extract-hold)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `HD-BUSINESS-DATE = BC-BUSINESS-DATE` N, `HD-REVERSED` Y |
| 2 | not reached: `HD-BUSINESS-DATE = BC-BUSINESS-DATE` N, `HD-REVERSED` N, `HD-PURCHASE` Y, `HD-FEE > 0` Y |
| 3 | not reached: `HD-BUSINESS-DATE = BC-BUSINESS-DATE` N, `HD-REVERSED` N, `HD-PURCHASE` Y, `HD-FEE > 0` N |
| 4 | not reached: `HD-BUSINESS-DATE = BC-BUSINESS-DATE` N, `HD-REVERSED` N, `HD-PURCHASE` N, `HD-FEE > 0` Y |
| 5 | not reached: `HD-BUSINESS-DATE = BC-BUSINESS-DATE` N, `HD-REVERSED` N, `HD-PURCHASE` N, `HD-FEE > 0` N |
| 6 | [case 02](02/) |
| 7 | [case 01](01/) |
| 8 | [case 03](03/) |
| 9 | [case 01](01/) |
| 10 | [case 03](03/) |
