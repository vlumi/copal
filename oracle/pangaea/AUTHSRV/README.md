# AUTHSRV: cases

26 cases reach 28 of the 32 rules of AUTHSRV's decision tables, and of the programs it calls, 9 of the 9 of AUTHPUR, 6 of the 6 of AUTHREV, 15 of the 15 of ATMWDL, 3 of the 3 of ATMBAL. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile. A rule that finds a file to have failed once open says so: it takes the file's status after a read, a write, a rewrite, a delete or a close to be none that the file's records make, a record or the end for a read in order, a record or none of that key for a read by key, which is the device's doing. And a rule of a program called that takes a parameter to be what no call of this program passes, a direction it never asks for, says so, as do the rules of the paragraphs performed only from such rules: no run of this program reaches them, whatever its input.

## AUTHSRV

### [1000-INITIALIZE](../../../blueprints/pangaea/rules/AUTHSRV.md#1000-initialize)

| Rule | Reached by |
| --: | --- |
| 1 | [case 13](13/) |
| 2 | [case 11](11/) |
| 3 | [case 01](01/) |

### [1100-FIND-LAST-HOLD](../../../blueprints/pangaea/rules/AUTHSRV.md#1100-find-last-hold)

| Rule | Reached by |
| --: | --- |
| 1 | [case 12](12/) |
| 2 | [case 01](01/) |

### [1100-FIND-LAST-HOLD, in its loop](../../../blueprints/pangaea/rules/AUTHSRV.md#1100-find-last-hold-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |
| 3 | [case 04](04/) |

### [2000-SERVE-REQUEST](../../../blueprints/pangaea/rules/AUTHSRV.md#2000-serve-request)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |

### [2100-AUTHORIZE](../../../blueprints/pangaea/rules/AUTHSRV.md#2100-authorize)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | [case 10](10/) |
| 3 | [case 09](09/) |
| 4 | [case 03](03/) |
| 5 | [case 01](01/) |
| 6 | [case 21](21/) |
| 7 | [case 03](03/) |
| 8 | [case 03](03/) |
| 9 | [case 09](09/) |
| 10 | [case 01](01/) |
| 11 | [case 05](05/) |

### [2200-READ-CARD-AND-ACCOUNT](../../../blueprints/pangaea/rules/AUTHSRV.md#2200-read-card-and-account)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 02](02/) |
| 3 | [case 01](01/) |

### [2300-READ-HOLD](../../../blueprints/pangaea/rules/AUTHSRV.md#2300-read-hold)

| Rule | Reached by |
| --: | --- |
| 1 | [case 09](09/) |
| 2 | [case 03](03/) |

### [2400-SAVE-CHANGES](../../../blueprints/pangaea/rules/AUTHSRV.md#2400-save-changes)

| Rule | Reached by |
| --: | --- |
| 1 | not reached, a write or rewrite by key found to fail: `REWRITE ACCT-ONLINE-RECORD INVALID KEY` Y |
| 2 | not reached, a write or rewrite by key found to fail: `REWRITE ACCT-ONLINE-RECORD INVALID KEY` N, `REWRITE CARD-RECORD INVALID KEY` Y |
| 3 | not reached, a write or rewrite by key found to fail: `REWRITE ACCT-ONLINE-RECORD INVALID KEY` N, `REWRITE CARD-RECORD INVALID KEY` N, `RQ-REVERSAL` Y, `REWRITE HOLD-RECORD INVALID KEY` Y |
| 4 | [case 03](03/) |
| 5 | not reached, a write or rewrite by key found to fail: `REWRITE ACCT-ONLINE-RECORD INVALID KEY` N, `REWRITE CARD-RECORD INVALID KEY` N, `RQ-REVERSAL` N, `WRITE HOLD-RECORD INVALID KEY` Y |
| 6 | [case 03](03/) |

## AUTHPUR

### [0000-MAIN](../../../blueprints/pangaea/rules/AUTHPUR.md#0000-main)

| Rule | Reached by |
| --: | --- |
| 1 | [case 09](09/) |
| 2 | [case 16](16/) |
| 3 | [case 01](01/) |
| 4 | [case 07](07/) |
| 5 | [case 06](06/) |
| 6 | [case 03](03/) |

### [1000-START-CARD-TOTALS](../../../blueprints/pangaea/rules/AUTHPUR.md#1000-start-card-totals)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 08](08/) |
| 3 | [case 07](07/) |

## AUTHREV

### [0000-MAIN](../../../blueprints/pangaea/rules/AUTHREV.md#0000-main)

| Rule | Reached by |
| --: | --- |
| 1 | [case 10](10/) |
| 2 | [case 18](18/) |
| 3 | [case 15](15/) |
| 4 | [case 03](03/) |

### [2000-APPROVE](../../../blueprints/pangaea/rules/AUTHREV.md#2000-approve)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | [case 14](14/) |

## ATMWDL

### [0000-MAIN](../../../blueprints/pangaea/rules/ATMWDL.md#0000-main)

| Rule | Reached by |
| --: | --- |
| 1 | [case 09](09/) |
| 2 | [case 16](16/) |
| 3 | [case 03](03/) |
| 4 | [case 17](17/) |
| 5 | [case 19](19/) |
| 6 | [case 21](21/) |
| 7 | [case 22](22/) |
| 8 | [case 23](23/) |
| 9 | [case 20](20/) |
| 10 | [case 24](24/) |
| 11 | [case 25](25/) |
| 12 | [case 26](26/) |

### [1000-START-CARD-TOTALS](../../../blueprints/pangaea/rules/ATMWDL.md#1000-start-card-totals)

| Rule | Reached by |
| --: | --- |
| 1 | [case 09](09/) |
| 2 | [case 20](20/) |
| 3 | [case 03](03/) |

## ATMBAL

### [0000-MAIN](../../../blueprints/pangaea/rules/ATMBAL.md#0000-main)

| Rule | Reached by |
| --: | --- |
| 1 | [case 09](09/) |
| 2 | [case 16](16/) |
| 3 | [case 03](03/) |
