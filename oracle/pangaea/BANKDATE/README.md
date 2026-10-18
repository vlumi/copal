# BANKDATE: cases

2 cases reach 2 of the 3 rules of BANKDATE's decision tables. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's.

## [0000-MAIN](../../../blueprints/pangaea/rules/BANKDATE.md#0000-main)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-IN-STATUS NOT = '00' OR WS-OUT-STATUS NOT = '00'` Y |
| 2 | [case 02](02/) |
| 3 | [case 01](01/) |
