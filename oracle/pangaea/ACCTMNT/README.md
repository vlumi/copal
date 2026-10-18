# ACCTMNT: cases

9 cases reach 33 of the 42 rules of ACCTMNT's decision tables. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's.

## [0000-MAIN](../../../blueprints/pangaea/rules/ACCTMNT.md#0000-main)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-BALANCES-OUT = WS-BALANCES-IN` N, `WS-REJECT-COUNT > 0` Y |
| 2 | not reached: `WS-BALANCES-OUT = WS-BALANCES-IN` N, `WS-REJECT-COUNT > 0` N |
| 3 | [case 01](01/) |
| 4 | [case 09](09/) |

## [1000-INITIALIZE](../../../blueprints/pangaea/rules/ACCTMNT.md#1000-initialize)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-FILE-STATUSES = ALL '0'` N |
| 2 | not reached: `WS-FILE-STATUSES = ALL '0'` Y, `READ CONTROL-FILE AT END` Y |
| 3 | [case 01](01/) |

## [1200-READ-MASTER](../../../blueprints/pangaea/rules/ACCTMNT.md#1200-read-master)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | not reached: `READ ACCT-IN AT END` N, `WS-MASTER-KEY > WS-PREV-MASTER-KEY` N |
| 3 | [case 01](01/) |

## [1300-READ-REQUEST](../../../blueprints/pangaea/rules/ACCTMNT.md#1300-read-request)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | not reached: `READ MAINT-IN AT END` N, `WS-REQUEST-ORDER-KEY > WS-PREV-REQUEST-ORDER-KEY` N |
| 3 | [case 01](01/) |

## [2000-NEXT-ACCOUNT](../../../blueprints/pangaea/rules/ACCTMNT.md#2000-next-account)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-REQUEST-KEY < WS-MASTER-KEY` Y, `WS-MASTER-KEY = WS-ACTIVE-KEY` Y |
| 2 | [case 01](01/) |
| 3 | [case 01](01/) |
| 4 | not reached: `WS-REQUEST-KEY < WS-MASTER-KEY` N, `WS-MASTER-KEY = WS-ACTIVE-KEY` N |

## [2100-APPLY-REQUESTS](../../../blueprints/pangaea/rules/ACCTMNT.md#2100-apply-requests)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |

## [2200-APPLY-REQUEST](../../../blueprints/pangaea/rules/ACCTMNT.md#2200-apply-request)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |
| 3 | [case 03](03/) |
| 4 | [case 03](03/) |
| 5 | [case 04](04/) |
| 6 | [case 04](04/) |
| 7 | [case 04](04/) |
| 8 | [case 08](08/) |

## [2210-OPEN](../../../blueprints/pangaea/rules/ACCTMNT.md#2210-open)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | [case 02](02/) |
| 3 | [case 04](04/) |
| 4 | [case 05](05/) |
| 5 | [case 01](01/) |

## [2220-CLOSE](../../../blueprints/pangaea/rules/ACCTMNT.md#2220-close)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `AM-BALANCE = 0` N |
| 2 | [case 03](03/) |

## [2230-FREEZE](../../../blueprints/pangaea/rules/ACCTMNT.md#2230-freeze)

| Rule | Reached by |
| --: | --- |
| 1 | [case 06](06/) |
| 2 | [case 04](04/) |

## [2240-UNFREEZE](../../../blueprints/pangaea/rules/ACCTMNT.md#2240-unfreeze)

| Rule | Reached by |
| --: | --- |
| 1 | [case 07](07/) |
| 2 | [case 04](04/) |

## [2250-LIMIT](../../../blueprints/pangaea/rules/ACCTMNT.md#2250-limit)

| Rule | Reached by |
| --: | --- |
| 1 | [case 04](04/) |
| 2 | [case 09](09/) |

## [2800-WRITE-HELD](../../../blueprints/pangaea/rules/ACCTMNT.md#2800-write-held)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |
