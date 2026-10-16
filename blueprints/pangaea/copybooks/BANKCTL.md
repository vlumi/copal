# BANKCTL

Copybook [BANKCTL.cpy](../../../pangaea/copybooks/BANKCTL.cpy), used by [TXPOST](../programs/TXPOST.md), [TXVALID](../programs/TXVALID.md).

## BANK-CONTROL-RECORD

80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | BANK-CONTROL-RECORD | 1 | 80 | | | | [BANKCTL.cpy:5](../../../pangaea/copybooks/BANKCTL.cpy#L5) |
| 05 | BC-BUSINESS-DATE | 1 | 8 | 9(8) | | | [BANKCTL.cpy:6](../../../pangaea/copybooks/BANKCTL.cpy#L6) |
| 05 | FILLER | 9 | 72 | X(72) | | | [BANKCTL.cpy:7](../../../pangaea/copybooks/BANKCTL.cpy#L7) |
