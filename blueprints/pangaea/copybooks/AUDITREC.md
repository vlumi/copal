# AUDITREC

Copybook [AUDITREC.cpy](../../../pangaea/copybooks/AUDITREC.cpy), used by [TXPOST](../programs/TXPOST.md).

## AUDIT-RECORD

100 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | AUDIT-RECORD | 1 | 100 | | | | [AUDITREC.cpy:5](../../../pangaea/copybooks/AUDITREC.cpy#L5) |
| 05 | AU-BUSINESS-DATE | 1 | 8 | 9(8) | | | [AUDITREC.cpy:6](../../../pangaea/copybooks/AUDITREC.cpy#L6) |
| 05 | AU-TRAN-ID | 9 | 10 | 9(10) | | | [AUDITREC.cpy:7](../../../pangaea/copybooks/AUDITREC.cpy#L7) |
| 05 | AU-ACCT-NO | 19 | 10 | 9(10) | | | [AUDITREC.cpy:8](../../../pangaea/copybooks/AUDITREC.cpy#L8) |
| 05 | AU-TRAN-CODE | 29 | 2 | XX | | | [AUDITREC.cpy:9](../../../pangaea/copybooks/AUDITREC.cpy#L9) |
| 05 | AU-AMOUNT | 31 | 11 | S9(9)V99 | | | [AUDITREC.cpy:10](../../../pangaea/copybooks/AUDITREC.cpy#L10) |
| 05 | AU-BALANCE-BEFORE | 42 | 13 | S9(11)V99 | | | [AUDITREC.cpy:11](../../../pangaea/copybooks/AUDITREC.cpy#L11) |
| 05 | AU-BALANCE-AFTER | 55 | 13 | S9(11)V99 | | | [AUDITREC.cpy:12](../../../pangaea/copybooks/AUDITREC.cpy#L12) |
| 05 | FILLER | 68 | 33 | X(33) | | | [AUDITREC.cpy:13](../../../pangaea/copybooks/AUDITREC.cpy#L13) |
