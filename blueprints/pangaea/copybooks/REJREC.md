# REJREC

```text
REJREC: REJECTED TRANSACTION, 100 BYTES
THE TRANSACTION AS RECEIVED, WITH THE STEP THAT REJECTED IT
AND WHY. ACCOUNT MAINTENANCE REJECTS ITS REQUESTS IN THE SAME
FORM, THE REQUEST IN PLACE OF THE TRANSACTION.

MODIFICATION HISTORY
DATE       BY  CHANGE
1987-03-12 JA  ORIGINAL VERSION.
1988-05-09 JA  CR-0031 ACCOUNT MAINTENANCE REASONS.
```

Copybook [REJREC.cpy](../../../pangaea/copybooks/REJREC.cpy), used by [ACCTMNT](../programs/ACCTMNT.md), [TXPOST](../programs/TXPOST.md), [TXVALID](../programs/TXVALID.md).

## REJECT-RECORD

100 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | REJECT-RECORD | 1 | 100 | | | | [REJREC.cpy:12](../../../pangaea/copybooks/REJREC.cpy#L12) |
| 05 | RJ-PROGRAM | 1 | 8 | X(8) | | | [REJREC.cpy:13](../../../pangaea/copybooks/REJREC.cpy#L13) |
| 05 | RJ-REASON | 9 | 4 | X(4) | | | [REJREC.cpy:14](../../../pangaea/copybooks/REJREC.cpy#L14) |
| 88 | RJ-BAD-ACCOUNT | | | | | VALUE 'ACCT' | [REJREC.cpy:15](../../../pangaea/copybooks/REJREC.cpy#L15) |
| 88 | RJ-BAD-COUNTER | | | | | VALUE 'CTRA' | [REJREC.cpy:16](../../../pangaea/copybooks/REJREC.cpy#L16) |
| 88 | RJ-BAD-CODE | | | | | VALUE 'CODE' | [REJREC.cpy:17](../../../pangaea/copybooks/REJREC.cpy#L17) |
| 88 | RJ-BAD-AMOUNT | | | | | VALUE 'AMNT' | [REJREC.cpy:18](../../../pangaea/copybooks/REJREC.cpy#L18) |
| 88 | RJ-BAD-TYPE | | | | | VALUE 'TYPE' | [REJREC.cpy:19](../../../pangaea/copybooks/REJREC.cpy#L19) |
| 88 | RJ-NO-ACCOUNT | | | | | VALUE 'UNKN' | [REJREC.cpy:20](../../../pangaea/copybooks/REJREC.cpy#L20) |
| 88 | RJ-CLOSED | | | | | VALUE 'CLSD' | [REJREC.cpy:21](../../../pangaea/copybooks/REJREC.cpy#L21) |
| 88 | RJ-FROZEN | | | | | VALUE 'FRZN' | [REJREC.cpy:22](../../../pangaea/copybooks/REJREC.cpy#L22) |
| 88 | RJ-NO-FUNDS | | | | | VALUE 'FUND' | [REJREC.cpy:23](../../../pangaea/copybooks/REJREC.cpy#L23) |
| 88 | RJ-EXISTS | | | | | VALUE 'EXST' | [REJREC.cpy:24](../../../pangaea/copybooks/REJREC.cpy#L24) |
| 88 | RJ-HAS-BALANCE | | | | | VALUE 'BALN' | [REJREC.cpy:25](../../../pangaea/copybooks/REJREC.cpy#L25) |
| 88 | RJ-BAD-STATUS | | | | | VALUE 'STAT' | [REJREC.cpy:26](../../../pangaea/copybooks/REJREC.cpy#L26) |
| 05 | RJ-TRAN | 13 | 80 | X(80) | | | [REJREC.cpy:27](../../../pangaea/copybooks/REJREC.cpy#L27) |
| 05 | FILLER | 93 | 8 | X(8) | | | [REJREC.cpy:28](../../../pangaea/copybooks/REJREC.cpy#L28) |
