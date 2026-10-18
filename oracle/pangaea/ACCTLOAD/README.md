# ACCTLOAD: cases

3 cases reach 8 of the 9 rules of ACCTLOAD's decision tables. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's.

## [1000-INITIALIZE](../../../blueprints/pangaea/rules/ACCTLOAD.md#1000-initialize)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-FILE-STATUSES = ALL '0'` N |
| 2 | [case 02](02/) |
| 3 | [case 01](01/) |

## [1200-READ-ACCOUNT](../../../blueprints/pangaea/rules/ACCTLOAD.md#1200-read-account)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |

## [2000-LOAD-ACCOUNT](../../../blueprints/pangaea/rules/ACCTLOAD.md#2000-load-account)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |

## [2100-WRITE-ONLINE](../../../blueprints/pangaea/rules/ACCTLOAD.md#2100-write-online)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | [case 01](01/) |
