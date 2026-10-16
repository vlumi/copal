# ACCTMAST

```text
ACCTMAST: ACCOUNT MASTER RECORD, 100 BYTES
IN ACCOUNT NUMBER ORDER. THE BALANCE MAY GO BELOW ZERO DOWN
TO THE OVERDRAFT LIMIT, AND BELOW THAT ONLY BY FEES.
```

Copybook [ACCTMAST.cpy](../../../pangaea/copybooks/ACCTMAST.cpy), used by [TXPOST](../programs/TXPOST.md).

## ACCT-MASTER-RECORD

100 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | ACCT-MASTER-RECORD | 1 | 100 | | | | [ACCTMAST.cpy:6](../../../pangaea/copybooks/ACCTMAST.cpy#L6) |
| 05 | AM-ACCT-NO | 1 | 10 | 9(10) | | | [ACCTMAST.cpy:7](../../../pangaea/copybooks/ACCTMAST.cpy#L7) |
| 05 | AM-CUST-NO | 11 | 8 | 9(8) | | | [ACCTMAST.cpy:8](../../../pangaea/copybooks/ACCTMAST.cpy#L8) |
| 05 | AM-SHORT-NAME | 19 | 20 | X(20) | | | [ACCTMAST.cpy:9](../../../pangaea/copybooks/ACCTMAST.cpy#L9) |
| 05 | AM-ACCT-TYPE | 39 | 1 | X | | | [ACCTMAST.cpy:10](../../../pangaea/copybooks/ACCTMAST.cpy#L10) |
| 88 | AM-CURRENT | | | | | VALUE 'C' | [ACCTMAST.cpy:11](../../../pangaea/copybooks/ACCTMAST.cpy#L11) |
| 88 | AM-SAVINGS | | | | | VALUE 'S' | [ACCTMAST.cpy:12](../../../pangaea/copybooks/ACCTMAST.cpy#L12) |
| 05 | AM-STATUS | 40 | 1 | X | | | [ACCTMAST.cpy:13](../../../pangaea/copybooks/ACCTMAST.cpy#L13) |
| 88 | AM-OPEN | | | | | VALUE 'O' | [ACCTMAST.cpy:14](../../../pangaea/copybooks/ACCTMAST.cpy#L14) |
| 88 | AM-FROZEN | | | | | VALUE 'F' | [ACCTMAST.cpy:15](../../../pangaea/copybooks/ACCTMAST.cpy#L15) |
| 88 | AM-CLOSED | | | | | VALUE 'X' | [ACCTMAST.cpy:16](../../../pangaea/copybooks/ACCTMAST.cpy#L16) |
| 05 | AM-OPEN-DATE | 41 | 8 | 9(8) | | | [ACCTMAST.cpy:17](../../../pangaea/copybooks/ACCTMAST.cpy#L17) |
| 05 | AM-LAST-ACTIVITY-DATE | 49 | 8 | 9(8) | | | [ACCTMAST.cpy:18](../../../pangaea/copybooks/ACCTMAST.cpy#L18) |
| 05 | AM-BALANCE | 57 | 13 | S9(11)V99 | | | [ACCTMAST.cpy:19](../../../pangaea/copybooks/ACCTMAST.cpy#L19) |
| 05 | AM-OVERDRAFT-LIMIT | 70 | 9 | 9(7)V99 | | | [ACCTMAST.cpy:20](../../../pangaea/copybooks/ACCTMAST.cpy#L20) |
| 05 | FILLER | 79 | 22 | X(22) | | | [ACCTMAST.cpy:21](../../../pangaea/copybooks/ACCTMAST.cpy#L21) |
