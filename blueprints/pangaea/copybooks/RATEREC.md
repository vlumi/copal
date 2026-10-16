# RATEREC

```text
RATEREC: INTEREST RATE TABLE RECORD, 20 BYTES
ONE RECORD PER RATE: CREDIT RATES IN TIERS BY BALANCE, ONE
DEBIT RATE FOR BALANCES BELOW ZERO, PER ACCOUNT TYPE. A
BALANCE EARNS THE RATE OF THE HIGHEST TIER IT REACHES, ON THE
WHOLE BALANCE. RATES ARE ANNUAL FRACTIONS: 0.015 IS 1.5%.

MODIFICATION HISTORY
DATE       BY  CHANGE
-1754      HAM ORIGINAL VERSION, ON CLAY, IN BABYLON: AT MOST
               20% A YEAR ON SILVER, 33 1/3% ON GRAIN.
2001-09-03 TL  CR-0298 IN COBOL, WITH TIERS. GRAIN DROPPED.
```

Copybook [RATEREC.cpy](../../../pangaea/copybooks/RATEREC.cpy), used by [INTACCR](../programs/INTACCR.md).

## RATE-RECORD

20 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | RATE-RECORD | 1 | 20 | | | | [RATEREC.cpy:14](../../../pangaea/copybooks/RATEREC.cpy#L14) |
| 05 | RT-ACCT-TYPE | 1 | 1 | X | | | [RATEREC.cpy:15](../../../pangaea/copybooks/RATEREC.cpy#L15) |
| 05 | RT-KIND | 2 | 1 | X | | | [RATEREC.cpy:16](../../../pangaea/copybooks/RATEREC.cpy#L16) |
| 88 | RT-CREDIT | | | | | VALUE 'C' | [RATEREC.cpy:17](../../../pangaea/copybooks/RATEREC.cpy#L17) |
| 88 | RT-DEBIT | | | | | VALUE 'D' | [RATEREC.cpy:18](../../../pangaea/copybooks/RATEREC.cpy#L18) |
| 05 | RT-FLOOR | 3 | 7 | S9(11)V99 | COMP-3 | | [RATEREC.cpy:19](../../../pangaea/copybooks/RATEREC.cpy#L19) |
| 05 | RT-ANNUAL-RATE | 10 | 4 | 9V9(6) | COMP-3 | | [RATEREC.cpy:20](../../../pangaea/copybooks/RATEREC.cpy#L20) |
| 05 | FILLER | 14 | 7 | X(7) | | | [RATEREC.cpy:21](../../../pangaea/copybooks/RATEREC.cpy#L21) |
