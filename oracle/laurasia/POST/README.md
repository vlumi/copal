# POST: cases

9 cases reach 13 of the 17 rules of POST's decision tables, and of the programs it calls, 0 of the 0 of OVLBIZ, 4 of the 4 of OVLFEE. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's.

## POST

### [MAIN](../../../blueprints/laurasia/rules/POST.md#main)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ACCOUNTS-STATUS = '00'` N |
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
| 2 | [case 04](04/) |
| 3 | [case 05](05/) |
| 4 | not reached: `READ ACCOUNTS INVALID KEY` N, `CA-BLOCKED` N, `CA-CLOSED` N, `CA-BALANCE > CA-LIMIT` Y, `REWRITE CARD-ACCOUNT INVALID KEY` Y |
| 5 | [case 02](02/) |
| 6 | not reached: `READ ACCOUNTS INVALID KEY` N, `CA-BLOCKED` N, `CA-CLOSED` N, `CA-BALANCE > CA-LIMIT` N, `REWRITE CARD-ACCOUNT INVALID KEY` Y |
| 7 | [case 06](06/) |

### [CHARGE-OVER-LIMIT](../../../blueprints/laurasia/rules/POST.md#charge-over-limit)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `CA-BUSINESS` Y, `WS-FEE > ZERO` Y |
| 2 | [case 08](08/) |
| 3 | [case 02](02/) |
| 4 | [case 07](07/) |

## OVLBIZ

## OVLFEE

### [MAIN](../../../blueprints/laurasia/rules/OVLFEE.md#main)

| Rule | Reached by |
| --: | --- |
| 1 | [case 07](07/) |
| 2 | [case 02](02/) |
| 3 | [case 09](09/) |
| 4 | [case 03](03/) |
