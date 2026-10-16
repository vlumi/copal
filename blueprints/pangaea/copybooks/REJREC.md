# REJREC

Copybook [REJREC.cpy](../../../pangaea/copybooks/REJREC.cpy), used by [TXPOST](../programs/TXPOST.md), [TXVALID](../programs/TXVALID.md).

## REJECT-RECORD

100 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | REJECT-RECORD | 1 | 100 | | | | [REJREC.cpy:6](../../../pangaea/copybooks/REJREC.cpy#L6) |
| 05 | RJ-PROGRAM | 1 | 8 | X(8) | | | [REJREC.cpy:7](../../../pangaea/copybooks/REJREC.cpy#L7) |
| 05 | RJ-REASON | 9 | 4 | X(4) | | | [REJREC.cpy:8](../../../pangaea/copybooks/REJREC.cpy#L8) |
| 88 | RJ-BAD-ACCOUNT | | | | | VALUE 'ACCT' | [REJREC.cpy:9](../../../pangaea/copybooks/REJREC.cpy#L9) |
| 88 | RJ-BAD-CODE | | | | | VALUE 'CODE' | [REJREC.cpy:10](../../../pangaea/copybooks/REJREC.cpy#L10) |
| 88 | RJ-BAD-AMOUNT | | | | | VALUE 'AMNT' | [REJREC.cpy:11](../../../pangaea/copybooks/REJREC.cpy#L11) |
| 88 | RJ-BAD-TYPE | | | | | VALUE 'TYPE' | [REJREC.cpy:12](../../../pangaea/copybooks/REJREC.cpy#L12) |
| 88 | RJ-NO-ACCOUNT | | | | | VALUE 'UNKN' | [REJREC.cpy:13](../../../pangaea/copybooks/REJREC.cpy#L13) |
| 88 | RJ-CLOSED | | | | | VALUE 'CLSD' | [REJREC.cpy:14](../../../pangaea/copybooks/REJREC.cpy#L14) |
| 88 | RJ-FROZEN | | | | | VALUE 'FRZN' | [REJREC.cpy:15](../../../pangaea/copybooks/REJREC.cpy#L15) |
| 88 | RJ-NO-FUNDS | | | | | VALUE 'FUND' | [REJREC.cpy:16](../../../pangaea/copybooks/REJREC.cpy#L16) |
| 05 | RJ-TRAN | 13 | 80 | X(80) | | | [REJREC.cpy:17](../../../pangaea/copybooks/REJREC.cpy#L17) |
| 05 | FILLER | 93 | 8 | X(8) | | | [REJREC.cpy:18](../../../pangaea/copybooks/REJREC.cpy#L18) |
