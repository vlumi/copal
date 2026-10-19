# AUTHSRV: cases

27 cases reach 27 of the 35 rules of AUTHSRV's decision tables, and of the programs it calls, 9 of the 9 of AUTHPUR, 6 of the 6 of AUTHREV, 9 of the 15 of ATMWDL, 3 of the 3 of ATMBAL. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's.

## AUTHSRV

### [1000-INITIALIZE](../../../blueprints/pangaea/rules/AUTHSRV.md#1000-initialize)

| Rule | Reached by |
| --: | --- |
| 1 | [case 13](13/) |
| 2 | [case 10](10/) |
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
| 2 | [case 02](02/) |
| 3 | [case 01](01/) |

### [2000-SERVE-REQUEST](../../../blueprints/pangaea/rules/AUTHSRV.md#2000-serve-request)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |

### [2100-AUTHORIZE](../../../blueprints/pangaea/rules/AUTHSRV.md#2100-authorize)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 15](15/) |
| 3 | not reached: `RQ-PURCHASE OR RQ-REVERSAL OR RQ-WITHDRAWAL OR RQ-BALANCE-ENQUIRY` Y, `RS-RESULT = SPACES AND RQ-REVERSAL` Y, `RS-RESULT = SPACES` N, `RS-APPROVED` Y, `RQ-BALANCE-ENQUIRY` N |
| 4 | [case 09](09/) |
| 5 | [case 08](08/) |
| 6 | [case 01](01/) |
| 7 | not reached: `RQ-PURCHASE OR RQ-REVERSAL OR RQ-WITHDRAWAL OR RQ-BALANCE-ENQUIRY` Y, `RS-RESULT = SPACES AND RQ-REVERSAL` N, `RS-RESULT = SPACES` Y, `RQ-PURCHASE` N, `RQ-REVERSAL` N, `RS-APPROVED` Y, `RQ-BALANCE-ENQUIRY` N, `RQ-WITHDRAWAL` Y |
| 8 | [case 09](09/) |
| 9 | [case 09](09/) |
| 10 | [case 22](22/) |
| 11 | not reached: `RQ-PURCHASE OR RQ-REVERSAL OR RQ-WITHDRAWAL OR RQ-BALANCE-ENQUIRY` Y, `RS-RESULT = SPACES AND RQ-REVERSAL` N, `RS-RESULT = SPACES` N, `RS-APPROVED` Y, `RQ-BALANCE-ENQUIRY` N |
| 12 | not reached: `RQ-PURCHASE OR RQ-REVERSAL OR RQ-WITHDRAWAL OR RQ-BALANCE-ENQUIRY` Y, `RS-RESULT = SPACES AND RQ-REVERSAL` N, `RS-RESULT = SPACES` N, `RS-APPROVED` Y, `RQ-BALANCE-ENQUIRY` Y |
| 13 | [case 01](01/) |
| 14 | [case 03](03/) |

### [2200-READ-CARD-AND-ACCOUNT](../../../blueprints/pangaea/rules/AUTHSRV.md#2200-read-card-and-account)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 11](11/) |
| 3 | [case 01](01/) |

### [2300-READ-HOLD](../../../blueprints/pangaea/rules/AUTHSRV.md#2300-read-hold)

| Rule | Reached by |
| --: | --- |
| 1 | [case 09](09/) |
| 2 | [case 02](02/) |

### [2400-SAVE-CHANGES](../../../blueprints/pangaea/rules/AUTHSRV.md#2400-save-changes)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `REWRITE ACCT-ONLINE-RECORD INVALID KEY` Y |
| 2 | not reached: `REWRITE ACCT-ONLINE-RECORD INVALID KEY` N, `REWRITE CARD-RECORD INVALID KEY` Y |
| 3 | not reached: `REWRITE ACCT-ONLINE-RECORD INVALID KEY` N, `REWRITE CARD-RECORD INVALID KEY` N, `RQ-REVERSAL` Y, `REWRITE HOLD-RECORD INVALID KEY` Y |
| 4 | [case 02](02/) |
| 5 | not reached: `REWRITE ACCT-ONLINE-RECORD INVALID KEY` N, `REWRITE CARD-RECORD INVALID KEY` N, `RQ-REVERSAL` N, `WRITE HOLD-RECORD INVALID KEY` Y |
| 6 | [case 08](08/) |

## AUTHPUR

### [0000-MAIN](../../../blueprints/pangaea/rules/AUTHPUR.md#0000-main)

| Rule | Reached by |
| --: | --- |
| 1 | [case 04](04/) |
| 2 | [case 05](05/) |
| 3 | [case 06](06/) |
| 4 | [case 01](01/) |
| 5 | [case 07](07/) |
| 6 | [case 08](08/) |

### [1000-START-CARD-TOTALS](../../../blueprints/pangaea/rules/AUTHPUR.md#1000-start-card-totals)

| Rule | Reached by |
| --: | --- |
| 1 | [case 08](08/) |
| 2 | [case 24](24/) |
| 3 | [case 01](01/) |

## AUTHREV

### [0000-MAIN](../../../blueprints/pangaea/rules/AUTHREV.md#0000-main)

| Rule | Reached by |
| --: | --- |
| 1 | [case 15](15/) |
| 2 | [case 16](16/) |
| 3 | [case 17](17/) |
| 4 | [case 02](02/) |

### [2000-APPROVE](../../../blueprints/pangaea/rules/AUTHREV.md#2000-approve)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 14](14/) |

## ATMWDL

### [0000-MAIN](../../../blueprints/pangaea/rules/ATMWDL.md#0000-main)

| Rule | Reached by |
| --: | --- |
| 1 | [case 18](18/) |
| 2 | [case 19](19/) |
| 3 | [case 09](09/) |
| 4 | not reached: `CR-MONTH-WITHDRAWALS < WS-FREE-WITHDRAWALS` Y, `CR-BLOCKED` N, `AO-OPEN` Y, `RQ-AMOUNT = 0 OR FUNCTION MOD(RQ-AMOUNT WS-NOTE) NOT = 0` N, `WS-DAY-CASH > CR-CASH-LIMIT` Y |
| 5 | not reached: `CR-MONTH-WITHDRAWALS < WS-FREE-WITHDRAWALS` Y, `CR-BLOCKED` N, `AO-OPEN` Y, `RQ-AMOUNT = 0 OR FUNCTION MOD(RQ-AMOUNT WS-NOTE) NOT = 0` N, `WS-DAY-CASH > CR-CASH-LIMIT` N, `WS-TO-HOLD > WS-AVAILABLE` Y |
| 6 | not reached: `CR-MONTH-WITHDRAWALS < WS-FREE-WITHDRAWALS` Y, `CR-BLOCKED` N, `AO-OPEN` Y, `RQ-AMOUNT = 0 OR FUNCTION MOD(RQ-AMOUNT WS-NOTE) NOT = 0` N, `WS-DAY-CASH > CR-CASH-LIMIT` N, `WS-TO-HOLD > WS-AVAILABLE` N |
| 7 | [case 25](25/) |
| 8 | [case 26](26/) |
| 9 | [case 21](21/) |
| 10 | not reached: `CR-MONTH-WITHDRAWALS < WS-FREE-WITHDRAWALS` N, `CR-BLOCKED` N, `AO-OPEN` Y, `RQ-AMOUNT = 0 OR FUNCTION MOD(RQ-AMOUNT WS-NOTE) NOT = 0` N, `WS-DAY-CASH > CR-CASH-LIMIT` Y |
| 11 | not reached: `CR-MONTH-WITHDRAWALS < WS-FREE-WITHDRAWALS` N, `CR-BLOCKED` N, `AO-OPEN` Y, `RQ-AMOUNT = 0 OR FUNCTION MOD(RQ-AMOUNT WS-NOTE) NOT = 0` N, `WS-DAY-CASH > CR-CASH-LIMIT` N, `WS-TO-HOLD > WS-AVAILABLE` Y |
| 12 | not reached: `CR-MONTH-WITHDRAWALS < WS-FREE-WITHDRAWALS` N, `CR-BLOCKED` N, `AO-OPEN` Y, `RQ-AMOUNT = 0 OR FUNCTION MOD(RQ-AMOUNT WS-NOTE) NOT = 0` N, `WS-DAY-CASH > CR-CASH-LIMIT` N, `WS-TO-HOLD > WS-AVAILABLE` N |

### [1000-START-CARD-TOTALS](../../../blueprints/pangaea/rules/ATMWDL.md#1000-start-card-totals)

| Rule | Reached by |
| --: | --- |
| 1 | [case 20](20/) |
| 2 | [case 27](27/) |
| 3 | [case 09](09/) |

## ATMBAL

### [0000-MAIN](../../../blueprints/pangaea/rules/ATMBAL.md#0000-main)

| Rule | Reached by |
| --: | --- |
| 1 | [case 22](22/) |
| 2 | [case 23](23/) |
| 3 | [case 09](09/) |
