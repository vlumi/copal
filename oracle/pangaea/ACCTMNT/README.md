# ACCTMNT: cases

13 cases reach 38 of the 41 rules of ACCTMNT's decision tables. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile.

## [0000-MAIN](../../../blueprints/pangaea/rules/ACCTMNT.md#0000-main)

| Rule | Reached by |
| --: | --- |
| 1 | not reached, totals the program adds up itself found to differ: `WS-BALANCES-OUT = WS-BALANCES-IN` N |
| 2 | [case 06](06/) |
| 3 | [case 05](05/) |

## [1000-INITIALIZE](../../../blueprints/pangaea/rules/ACCTMNT.md#1000-initialize)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | [case 02](02/) |
| 3 | [case 01](01/) |

## [1200-READ-MASTER](../../../blueprints/pangaea/rules/ACCTMNT.md#1200-read-master)

| Rule | Reached by |
| --: | --- |
| 1 | [case 08](08/) |
| 2 | [case 01](01/) |
| 3 | [case 01](01/) |

## [1300-READ-REQUEST](../../../blueprints/pangaea/rules/ACCTMNT.md#1300-read-request)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 04](04/) |
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
| 2 | [case 06](06/) |
| 3 | [case 01](01/) |
| 4 | [case 01](01/) |
| 5 | [case 08](08/) |
| 6 | [case 08](08/) |
| 7 | [case 08](08/) |
| 8 | [case 10](10/) |

## [2210-OPEN](../../../blueprints/pangaea/rules/ACCTMNT.md#2210-open)

| Rule | Reached by |
| --: | --- |
| 1 | [case 08](08/) |
| 2 | [case 07](07/) |
| 3 | [case 01](01/) |
| 4 | [case 11](11/) |
| 5 | [case 06](06/) |

## [2220-CLOSE](../../../blueprints/pangaea/rules/ACCTMNT.md#2220-close)

| Rule | Reached by |
| --: | --- |
| 1 | [case 08](08/) |
| 2 | [case 01](01/) |

## [2230-FREEZE](../../../blueprints/pangaea/rules/ACCTMNT.md#2230-freeze)

| Rule | Reached by |
| --: | --- |
| 1 | [case 13](13/) |
| 2 | [case 08](08/) |

## [2240-UNFREEZE](../../../blueprints/pangaea/rules/ACCTMNT.md#2240-unfreeze)

| Rule | Reached by |
| --: | --- |
| 1 | [case 09](09/) |
| 2 | [case 08](08/) |

## [2250-LIMIT](../../../blueprints/pangaea/rules/ACCTMNT.md#2250-limit)

| Rule | Reached by |
| --: | --- |
| 1 | [case 08](08/) |
| 2 | [case 12](12/) |

## [2800-WRITE-HELD](../../../blueprints/pangaea/rules/ACCTMNT.md#2800-write-held)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 06](06/) |
