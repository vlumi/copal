# GLSUM: cases

4 cases reach 16 of the 27 rules of GLSUM's decision tables. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's.

## [0000-MAIN](../../../blueprints/pangaea/rules/GLSUM.md#0000-main)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-CLOSED-MONTH-NUMBER = 0` N, `WS-TRANSIT-NET = 0` N |
| 2 | [case 01](01/) |
| 3 | not reached: `WS-CLOSED-MONTH-NUMBER = 0` Y, `WS-TRANSIT-NET = 0` N |
| 4 | [case 03](03/) |

## [1000-INITIALIZE](../../../blueprints/pangaea/rules/GLSUM.md#1000-initialize)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-FILE-STATUSES = ALL '0'` N |
| 2 | not reached: `WS-FILE-STATUSES = ALL '0'` Y, `READ ACCR-IN AT END` N, `WS-CM-MONTH = 12` Y |
| 3 | [case 01](01/) |
| 4 | not reached: `WS-FILE-STATUSES = ALL '0'` Y, `READ ACCR-IN AT END` Y, `WS-CM-MONTH = 12` Y |
| 5 | [case 03](03/) |

## [2000-POST-MONTH](../../../blueprints/pangaea/rules/GLSUM.md#2000-post-month)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-AUDIT-STATUS = '10'` N |
| 2 | [case 01](01/) |

## [2000-POST-MONTH, in its loop](../../../blueprints/pangaea/rules/GLSUM.md#2000-post-month-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 01](01/) |

## [2100-POST-TO-LEDGER](../../../blueprints/pangaea/rules/GLSUM.md#2100-post-to-ledger)

| Rule | Reached by |
| --: | --- |
| 1 | [case 04](04/) |
| 2 | [case 02](02/) |
| 3 | not reached: `GS-CODE(GS-CX) = AU-TRAN-CODE` N, `AU-AMOUNT < 0` Y |
| 4 | not reached: `GS-CODE(GS-CX) = AU-TRAN-CODE` N, `AU-AMOUNT < 0` N |

## [2200-FIND-ACCOUNT](../../../blueprints/pangaea/rules/GLSUM.md#2200-find-account)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | not reached: `WS-LE-ACCOUNT(WS-LX) = WS-WANTED-ACCOUNT` N |

## [3000-WRITE-LEDGER](../../../blueprints/pangaea/rules/GLSUM.md#3000-write-ledger)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-DEBIT-TOTAL = WS-CREDIT-TOTAL` N |
| 2 | [case 01](01/) |

## [3000-WRITE-LEDGER, in its loop](../../../blueprints/pangaea/rules/GLSUM.md#3000-write-ledger-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |
| 3 | [case 01](01/) |

## [8000-WRITE-REPORT](../../../blueprints/pangaea/rules/GLSUM.md#8000-write-report)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | [case 01](01/) |
| 3 | not reached: `WS-CLOSED-MONTH-NUMBER = 0` N, `WS-TRANSIT-NET = 0` N |
