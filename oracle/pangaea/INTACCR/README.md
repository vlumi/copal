# INTACCR: cases

6 cases reach 23 of the 33 rules of INTACCR's decision tables. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's.

## [0000-MAIN](../../../blueprints/pangaea/rules/INTACCR.md#0000-main)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ACCRUAL-KEY = HIGH-VALUES` N |
| 2 | [case 04](04/) |

## [1000-INITIALIZE](../../../blueprints/pangaea/rules/INTACCR.md#1000-initialize)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-FILE-STATUSES = ALL '0'` N |
| 2 | not reached: `WS-FILE-STATUSES = ALL '0'` Y, `READ CONTROL-FILE AT END` Y |
| 3 | [case 01](01/) |

## [1100-LOAD-RATES](../../../blueprints/pangaea/rules/INTACCR.md#1100-load-rates)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-RATE-STATUS = '10'` N |
| 2 | [case 01](01/) |

## [1100-LOAD-RATES, in its loop](../../../blueprints/pangaea/rules/INTACCR.md#1100-load-rates-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-RATE-COUNT = 20` Y |
| 2 | [case 01](01/) |

## [1200-READ-MASTER](../../../blueprints/pangaea/rules/INTACCR.md#1200-read-master)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
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
| 1 | [case 04](04/) |
| 2 | [case 03](03/) |

## [2100-FIND-RATE](../../../blueprints/pangaea/rules/INTACCR.md#2100-find-rate)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 03](03/) |

## [2100-FIND-RATE, in its loop](../../../blueprints/pangaea/rules/INTACCR.md#2100-find-rate-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-RT-ACCT-TYPE(WS-RX) = AM-ACCT-TYPE` Y, `AM-BALANCE < 0` Y, `WS-RT-DEBIT(WS-RX)` Y |
| 2 | [case 06](06/) |
| 3 | [case 03](03/) |
| 4 | [case 05](05/) |
| 5 | [case 01](01/) |

## [2200-FIND-MONTH-START](../../../blueprints/pangaea/rules/INTACCR.md#2200-find-month-start)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-TH-MONTH = 12` Y |
| 2 | [case 03](03/) |

## [2300-ACCRUE-DAYS](../../../blueprints/pangaea/rules/INTACCR.md#2300-accrue-days)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-INCREMENT < 0` Y |
| 2 | [case 03](03/) |

## [2400-CLOSE-MONTH](../../../blueprints/pangaea/rules/INTACCR.md#2400-close-month)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `AC-DUE < 0` Y |
| 2 | [case 04](04/) |

## [8000-WRITE-REPORT](../../../blueprints/pangaea/rules/INTACCR.md#8000-write-report)

| Rule | Reached by |
| --: | --- |
| 1 | [case 04](04/) |
| 2 | not reached: `WS-CLOSED-MONTH = 0` N |
