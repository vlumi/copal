# POST: cases

10 cases reach 14 of the 17 rules of POST's decision tables, and of the programs it calls, 0 of the 0 of OVLBIZ, 4 of the 4 of OVLFEE. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile. A rule that finds a file to have failed once open says so: it takes the file's status after a read, a write, a rewrite, a delete or a close to be none that the file's records make, a record or the end for a read in order, a record or none of that key for a read by key, which is the device's doing. And a rule of a program called that takes a parameter to be what no call of this program passes, a direction it never asks for, says so, as do the rules of the paragraphs performed only from such rules: no run of this program reaches them, whatever its input.

## POST

### [MAIN](../../../blueprints/laurasia/rules/POST.md#main)

| Rule | Reached by |
| --: | --- |
| 1 | [case 04](04/) |
| 2 | [case 01](01/) |
| 3 | [case 02](02/) |

### [MAIN, in its loop](../../../blueprints/laurasia/rules/POST.md#main-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |
| 3 | [case 01](01/) |

### [POST-PRESENTMENT](../../../blueprints/laurasia/rules/POST.md#post-presentment)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 06](06/) |
| 3 | [case 07](07/) |
| 4 | not reached, a write or rewrite by key found to fail: `READ ACCOUNTS INVALID KEY` N, `CA-BLOCKED` N, `CA-CLOSED` N, `CA-BALANCE > CA-LIMIT` Y, `REWRITE CARD-ACCOUNT INVALID KEY` Y |
| 5 | [case 02](02/) |
| 6 | not reached, a write or rewrite by key found to fail: `READ ACCOUNTS INVALID KEY` N, `CA-BLOCKED` N, `CA-CLOSED` N, `CA-BALANCE > CA-LIMIT` N, `REWRITE CARD-ACCOUNT INVALID KEY` Y |
| 7 | [case 08](08/) |

### [CHARGE-OVER-LIMIT](../../../blueprints/laurasia/rules/POST.md#charge-over-limit)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `CA-BUSINESS` Y, `WS-FEE > ZERO` Y |
| 2 | [case 03](03/) |
| 3 | [case 02](02/) |
| 4 | [case 09](09/) |

## OVLBIZ

## OVLFEE

### [MAIN](../../../blueprints/laurasia/rules/OVLFEE.md#main)

| Rule | Reached by |
| --: | --- |
| 1 | [case 09](09/) |
| 2 | [case 02](02/) |
| 3 | [case 10](10/) |
| 4 | [case 05](05/) |
