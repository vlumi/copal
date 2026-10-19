# CYCLE: cases

7 cases reach 6 of the 10 rules of CYCLE's decision tables, and of the programs it calls, 3 of the 3 of INTCALC, 4 of the 4 of MINPAY. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's.

## CYCLE

### [MAIN](../../../blueprints/laurasia/rules/CYCLE.md#main)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ACCOUNTS-STATUS = '00'` N |
| 2 | [case 01](01/) |

### [MAIN, in its loop](../../../blueprints/laurasia/rules/CYCLE.md#main-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |

### [CLOSE-CYCLE](../../../blueprints/laurasia/rules/CYCLE.md#close-cycle)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `NOT CA-WAS-PAID-IN-FULL AND CA-STATEMENT-BALANCE > ZERO` Y, `WS-INTEREST > ZERO` Y, `REWRITE CARD-ACCOUNT INVALID KEY` Y |
| 2 | [case 01](01/) |
| 3 | not reached: `NOT CA-WAS-PAID-IN-FULL AND CA-STATEMENT-BALANCE > ZERO` Y, `WS-INTEREST > ZERO` N, `REWRITE CARD-ACCOUNT INVALID KEY` Y |
| 4 | [case 04](04/) |
| 5 | not reached: `NOT CA-WAS-PAID-IN-FULL AND CA-STATEMENT-BALANCE > ZERO` N, `WS-INTEREST > ZERO` N, `REWRITE CARD-ACCOUNT INVALID KEY` Y |
| 6 | [case 02](02/) |

## INTCALC

### [MAIN](../../../blueprints/laurasia/rules/INTCALC.md#main)

| Rule | Reached by |
| --: | --- |
| 1 | [case 06](06/) |
| 2 | [case 07](07/) |
| 3 | [case 01](01/) |

## MINPAY

### [MAIN](../../../blueprints/laurasia/rules/MINPAY.md#main)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 01](01/) |
| 3 | [case 05](05/) |
| 4 | [case 03](03/) |
