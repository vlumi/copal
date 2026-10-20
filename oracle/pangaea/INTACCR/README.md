# INTACCR: cases

13 cases reach 33 of the 33 rules of INTACCR's decision tables. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile.

## [0000-MAIN](../../../blueprints/pangaea/rules/INTACCR.md#0000-main)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | [case 01](01/) |

## [1000-INITIALIZE](../../../blueprints/pangaea/rules/INTACCR.md#1000-initialize)

| Rule | Reached by |
| --: | --- |
| 1 | [case 10](10/) |
| 2 | [case 08](08/) |
| 3 | [case 01](01/) |

## [1100-LOAD-RATES](../../../blueprints/pangaea/rules/INTACCR.md#1100-load-rates)

| Rule | Reached by |
| --: | --- |
| 1 | [case 11](11/) |
| 2 | [case 01](01/) |

## [1100-LOAD-RATES, in its loop](../../../blueprints/pangaea/rules/INTACCR.md#1100-load-rates-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 09](09/) |
| 2 | [case 01](01/) |

## [1200-READ-MASTER](../../../blueprints/pangaea/rules/INTACCR.md#1200-read-master)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |

## [1300-READ-ACCRUAL](../../../blueprints/pangaea/rules/INTACCR.md#1300-read-accrual)

| Rule | Reached by |
| --: | --- |
| 1 | [case 04](04/) |
| 2 | [case 01](01/) |

## [1400-TAKE-ACCRUAL](../../../blueprints/pangaea/rules/INTACCR.md#1400-take-accrual)

| Rule | Reached by |
| --: | --- |
| 1 | [case 04](04/) |
| 2 | [case 01](01/) |
| 3 | [case 02](02/) |

## [2000-ACCRUE-ACCOUNT](../../../blueprints/pangaea/rules/INTACCR.md#2000-accrue-account)

| Rule | Reached by |
| --: | --- |
| 1 | [case 12](12/) |
| 2 | [case 01](01/) |

## [2100-FIND-RATE](../../../blueprints/pangaea/rules/INTACCR.md#2100-find-rate)

| Rule | Reached by |
| --: | --- |
| 1 | [case 04](04/) |
| 2 | [case 01](01/) |

## [2100-FIND-RATE, in its loop](../../../blueprints/pangaea/rules/INTACCR.md#2100-find-rate-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 13](13/) |
| 2 | [case 06](06/) |
| 3 | [case 01](01/) |
| 4 | [case 05](05/) |
| 5 | [case 04](04/) |

## [2200-FIND-MONTH-START](../../../blueprints/pangaea/rules/INTACCR.md#2200-find-month-start)

| Rule | Reached by |
| --: | --- |
| 1 | [case 07](07/) |
| 2 | [case 01](01/) |

## [2300-ACCRUE-DAYS](../../../blueprints/pangaea/rules/INTACCR.md#2300-accrue-days)

| Rule | Reached by |
| --: | --- |
| 1 | [case 13](13/) |
| 2 | [case 01](01/) |

## [2400-CLOSE-MONTH](../../../blueprints/pangaea/rules/INTACCR.md#2400-close-month)

| Rule | Reached by |
| --: | --- |
| 1 | [case 12](12/) |
| 2 | [case 13](13/) |

## [8000-WRITE-REPORT](../../../blueprints/pangaea/rules/INTACCR.md#8000-write-report)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 12](12/) |
