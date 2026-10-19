# MONCHG: cases

11 cases reach 18 of the 21 rules of MONCHG's decision tables. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it.

## [1000-INITIALIZE](../../../blueprints/pangaea/rules/MONCHG.md#1000-initialize)

| Rule | Reached by |
| --: | --- |
| 1 | [case 11](11/) |
| 2 | [case 08](08/) |
| 3 | [case 01](01/) |

## [1200-READ-ACCOUNT](../../../blueprints/pangaea/rules/MONCHG.md#1200-read-account)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | not reached: `READ ACCT-IN AT END` Y, `READ ACCR-IN AT END` N, `MASTER-AT-END AND ACCRUALS-AT-END` Y |
| 3 | [case 09](09/) |
| 4 | not reached: `READ ACCT-IN AT END` N, `READ ACCR-IN AT END` Y, `MASTER-AT-END AND ACCRUALS-AT-END` Y |
| 5 | [case 10](10/) |
| 6 | not reached: `READ ACCT-IN AT END` N, `READ ACCR-IN AT END` N, `MASTER-AT-END AND ACCRUALS-AT-END` Y |
| 7 | [case 02](02/) |
| 8 | [case 01](01/) |

## [2000-CHARGE-ACCOUNT](../../../blueprints/pangaea/rules/MONCHG.md#2000-charge-account)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 03](03/) |

## [2100-SETTLE-INTEREST](../../../blueprints/pangaea/rules/MONCHG.md#2100-settle-interest)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 05](05/) |
| 3 | [case 04](04/) |

## [2200-CHARGE-FEE](../../../blueprints/pangaea/rules/MONCHG.md#2200-charge-fee)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 07](07/) |
| 3 | [case 06](06/) |

## [8000-WRITE-REPORT](../../../blueprints/pangaea/rules/MONCHG.md#8000-write-report)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | [case 01](01/) |
