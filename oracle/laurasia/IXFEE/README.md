# IXFEE: cases

15 cases reach 13 of the 13 rules of IXFEE's decision tables, and of the programs it calls, 10 of the 16 of FEESLOT. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's.

## IXFEE

### [MAIN](../../../blueprints/laurasia/rules/IXFEE.md#main)

| Rule | Reached by |
| --: | --- |
| 1 | [case 11](11/) |
| 2 | [case 01](01/) |
| 3 | [case 02](02/) |

### [MAIN, in its loop](../../../blueprints/laurasia/rules/IXFEE.md#main-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |

### [PRICE-PRESENTMENT](../../../blueprints/laurasia/rules/IXFEE.md#price-presentment)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 02](02/) |
| 3 | [case 01](01/) |

### [FIND-GROUP](../../../blueprints/laurasia/rules/IXFEE.md#find-group)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 03](03/) |
| 3 | [case 04](04/) |
| 4 | [case 05](05/) |
| 5 | [case 01](01/) |

## FEESLOT

### [MAIN](../../../blueprints/laurasia/rules/FEESLOT.md#main)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `LK-PRODUCT = 'CL'` Y, `LK-REGION = 'D'` Y, `LK-GROUP < 1 OR LK-GROUP > 5` Y |
| 2 | [case 01](01/) |
| 3 | not reached: `LK-PRODUCT = 'CL'` Y, `LK-REGION = 'D'` N, `LK-GROUP < 1 OR LK-GROUP > 5` Y, `LK-REGION = 'I'` Y |
| 4 | [case 10](10/) |
| 5 | [case 09](09/) |
| 6 | not reached: `LK-PRODUCT = 'CL'` N, `LK-REGION = 'D'` Y, `LK-GROUP < 1 OR LK-GROUP > 5` Y, `LK-PRODUCT = 'GD'` Y |
| 7 | [case 07](07/) |
| 8 | not reached: `LK-PRODUCT = 'CL'` N, `LK-REGION = 'D'` N, `LK-GROUP < 1 OR LK-GROUP > 5` Y, `LK-REGION = 'I'` Y, `LK-PRODUCT = 'GD'` Y |
| 9 | [case 13](13/) |
| 10 | [case 12](12/) |
| 11 | not reached: `LK-PRODUCT = 'CL'` N, `LK-REGION = 'D'` Y, `LK-GROUP < 1 OR LK-GROUP > 5` Y, `LK-PRODUCT = 'GD'` N, `LK-PRODUCT = 'BZ'` Y |
| 12 | [case 08](08/) |
| 13 | not reached: `LK-PRODUCT = 'CL'` N, `LK-REGION = 'D'` N, `LK-GROUP < 1 OR LK-GROUP > 5` Y, `LK-REGION = 'I'` Y, `LK-PRODUCT = 'GD'` N, `LK-PRODUCT = 'BZ'` Y |
| 14 | [case 15](15/) |
| 15 | [case 14](14/) |
| 16 | [case 06](06/) |
