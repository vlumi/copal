# PAYIN: cases

6 cases reach 15 of the 17 rules of PAYIN's decision tables. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile.

## [MAIN](../../../blueprints/laurasia/rules/PAYIN.md#main)

| Rule | Reached by |
| --: | --- |
| 1 | [case 05](05/) |
| 2 | [case 01](01/) |
| 3 | [case 04](04/) |

## [MAIN, in its loop](../../../blueprints/laurasia/rules/PAYIN.md#main-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |

## [APPLY-RECORD](../../../blueprints/laurasia/rules/PAYIN.md#apply-record)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 02](02/) |
| 3 | not reached, a write or rewrite by key found to fail: `RECORD-IS-GOOD` Y, `READ ACCOUNTS INVALID KEY` N, `PF-PAYMENT` Y, `REWRITE CARD-ACCOUNT INVALID KEY` Y |
| 4 | [case 04](04/) |
| 5 | not reached, a write or rewrite by key found to fail: `RECORD-IS-GOOD` Y, `READ ACCOUNTS INVALID KEY` N, `PF-PAYMENT` N, `REWRITE CARD-ACCOUNT INVALID KEY` Y |
| 6 | [case 04](04/) |

## [VALIDATE-RECORD](../../../blueprints/laurasia/rules/PAYIN.md#validate-record)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 02](02/) |
| 3 | [case 03](03/) |
| 4 | [case 02](02/) |

## [APPLY-PAYMENT](../../../blueprints/laurasia/rules/PAYIN.md#apply-payment)

| Rule | Reached by |
| --: | --- |
| 1 | [case 04](04/) |
| 2 | [case 06](06/) |
