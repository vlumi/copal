# TXPOST: cases

13 cases reach 30 of the 30 rules of TXPOST's decision tables. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile.

## [0000-MAIN](../../../blueprints/pangaea/rules/TXPOST.md#0000-main)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 03](03/) |

## [1000-INITIALIZE](../../../blueprints/pangaea/rules/TXPOST.md#1000-initialize)

| Rule | Reached by |
| --: | --- |
| 1 | [case 10](10/) |
| 2 | [case 08](08/) |
| 3 | [case 01](01/) |

## [1200-READ-MASTER](../../../blueprints/pangaea/rules/TXPOST.md#1200-read-master)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 01](01/) |
| 3 | [case 01](01/) |

## [1300-READ-TRAN](../../../blueprints/pangaea/rules/TXPOST.md#1300-read-tran)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 09](09/) |
| 3 | [case 01](01/) |

## [2000-MATCH](../../../blueprints/pangaea/rules/TXPOST.md#2000-match)

| Rule | Reached by |
| --: | --- |
| 1 | [case 11](11/) |
| 2 | [case 02](02/) |
| 3 | [case 01](01/) |

## [3000-APPLY-TRAN](../../../blueprints/pangaea/rules/TXPOST.md#3000-apply-tran)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 09](09/) |
| 3 | [case 04](04/) |
| 4 | [case 06](06/) |
| 5 | [case 12](12/) |
| 6 | [case 03](03/) |
| 7 | [case 07](07/) |

## [3200-POST-DEBIT](../../../blueprints/pangaea/rules/TXPOST.md#3200-post-debit)

| Rule | Reached by |
| --: | --- |
| 1 | [case 13](13/) |
| 2 | [case 05](05/) |
| 3 | [case 03](03/) |

## [3500-WRITE-LEG](../../../blueprints/pangaea/rules/TXPOST.md#3500-write-leg)

| Rule | Reached by |
| --: | --- |
| 1 | [case 05](05/) |
| 2 | [case 11](11/) |

## [3900-REJECT-TRAN](../../../blueprints/pangaea/rules/TXPOST.md#3900-reject-tran)

| Rule | Reached by |
| --: | --- |
| 1 | [case 11](11/) |
| 2 | [case 02](02/) |

## [5000-RECONCILE](../../../blueprints/pangaea/rules/TXPOST.md#5000-reconcile)

| Rule | Reached by |
| --: | --- |
| 1 | [case 11](11/) |
| 2 | [case 02](02/) |
