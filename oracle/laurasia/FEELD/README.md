# FEELD: cases

16 cases reach 6 of the 7 rules of FEELD's decision tables, and of the programs it calls, 16 of the 16 of FEESLOT. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's.

## FEELD

### [MAIN](../../../blueprints/laurasia/rules/FEELD.md#main)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-FEE-STATUS = '00'` N |
| 2 | [case 01](01/) |

### [MAIN, in its loop](../../../blueprints/laurasia/rules/FEELD.md#main-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |

### [LOAD-FEE](../../../blueprints/laurasia/rules/FEELD.md#load-fee)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 11](11/) |
| 3 | [case 05](05/) |

## FEESLOT

### [MAIN](../../../blueprints/laurasia/rules/FEESLOT.md#main)

| Rule | Reached by |
| --: | --- |
| 1 | [case 10](10/) |
| 2 | [case 11](11/) |
| 3 | [case 12](12/) |
| 4 | [case 05](05/) |
| 5 | [case 02](02/) |
| 6 | [case 13](13/) |
| 7 | [case 06](06/) |
| 8 | [case 14](14/) |
| 9 | [case 07](07/) |
| 10 | [case 03](03/) |
| 11 | [case 15](15/) |
| 12 | [case 08](08/) |
| 13 | [case 16](16/) |
| 14 | [case 09](09/) |
| 15 | [case 04](04/) |
| 16 | [case 01](01/) |
