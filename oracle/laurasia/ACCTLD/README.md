# ACCTLD: cases

1 cases reach 4 of the 5 rules of ACCTLD's decision tables. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's.

## [MAIN](../../../blueprints/laurasia/rules/ACCTLD.md#main)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ACCOUNTS-STATUS = '00'` N |
| 2 | [case 01](01/) |

## [MAIN, in its loop](../../../blueprints/laurasia/rules/ACCTLD.md#main-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |
| 3 | [case 01](01/) |
