# BANKDATE: cases

3 cases reach 3 of the 3 rules of BANKDATE's decision tables. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it.

## [0000-MAIN](../../../blueprints/pangaea/rules/BANKDATE.md#0000-main)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | [case 02](02/) |
| 3 | [case 01](01/) |
