# STMTPRT: cases

12 cases reach 31 of the 34 rules of STMTPRT's decision tables. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile.

## [0000-MAIN](../../../blueprints/pangaea/rules/STMTPRT.md#0000-main)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-CLOSED-MONTH-NUMBER = 0` N, `SORT-RETURN = 0` N |
| 2 | [case 01](01/) |
| 3 | [case 09](09/) |

## [1000-INITIALIZE](../../../blueprints/pangaea/rules/STMTPRT.md#1000-initialize)

| Rule | Reached by |
| --: | --- |
| 1 | [case 10](10/) |
| 2 | [case 02](02/) |
| 3 | [case 01](01/) |
| 4 | not reached: `WS-FILE-STATUSES = ALL '0'` Y, `READ ACCR-IN AT END` Y, `WS-CM-MONTH = 12` Y |
| 5 | [case 09](09/) |

## [2000-SELECT-POSTINGS](../../../blueprints/pangaea/rules/STMTPRT.md#2000-select-postings)

| Rule | Reached by |
| --: | --- |
| 1 | [case 11](11/) |
| 2 | [case 01](01/) |

## [2000-SELECT-POSTINGS, in its loop](../../../blueprints/pangaea/rules/STMTPRT.md#2000-select-postings-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 02](02/) |

## [3000-PRINT-STATEMENTS](../../../blueprints/pangaea/rules/STMTPRT.md#3000-print-statements)

| Rule | Reached by |
| --: | --- |
| 1 | [case 05](05/) |
| 2 | [case 01](01/) |

## [3000-PRINT-STATEMENTS, in its loop](../../../blueprints/pangaea/rules/STMTPRT.md#3000-print-statements-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 04](04/) |
| 2 | [case 01](01/) |

## [3100-RETURN-POSTING](../../../blueprints/pangaea/rules/STMTPRT.md#3100-return-posting)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |

## [3200-READ-MASTER](../../../blueprints/pangaea/rules/STMTPRT.md#3200-read-master)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |

## [3300-PRINT-STATEMENT](../../../blueprints/pangaea/rules/STMTPRT.md#3300-print-statement)

| Rule | Reached by |
| --: | --- |
| 1 | [case 07](07/) |
| 2 | [case 01](01/) |
| 3 | [case 01](01/) |

## [3400-PRINT-HEADING](../../../blueprints/pangaea/rules/STMTPRT.md#3400-print-heading)

| Rule | Reached by |
| --: | --- |
| 1 | [case 06](06/) |
| 2 | [case 01](01/) |
| 3 | [case 06](06/) |
| 4 | [case 01](01/) |

## [3500-PRINT-POSTING](../../../blueprints/pangaea/rules/STMTPRT.md#3500-print-posting)

| Rule | Reached by |
| --: | --- |
| 1 | [case 12](12/) |
| 2 | [case 03](03/) |
| 3 | [case 08](08/) |
| 4 | [case 01](01/) |

## [8000-WRITE-REPORT](../../../blueprints/pangaea/rules/STMTPRT.md#8000-write-report)

| Rule | Reached by |
| --: | --- |
| 1 | [case 09](09/) |
| 2 | not reached, totals the program adds up itself found to differ: `WS-CLOSED-MONTH-NUMBER = 0` N, `WS-PRINTED-COUNT = WS-SELECTED-COUNT` N |
| 3 | [case 01](01/) |
