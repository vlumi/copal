# AUDITREC

```text
AUDITREC: AUDIT TRAIL RECORD, 100 BYTES, ONE PER POSTING
THE AMOUNT IS SIGNED: CREDITS POSITIVE, DEBITS NEGATIVE.
```

Copybook [AUDITREC.cpy](../../../pangaea/copybooks/AUDITREC.cpy), used by [GLSUM](../programs/GLSUM.md), [STMTPRT](../programs/STMTPRT.md), [TXPOST](../programs/TXPOST.md).

## AUDIT-RECORD

100 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | AUDIT-RECORD | group | 1 | 100 | | | | [AUDITREC.cpy:5](../../../pangaea/copybooks/AUDITREC.cpy#L5) |
| 05 | ├ AU-BUSINESS-DATE | date | 1 | 8 | 9(8) | | | [AUDITREC.cpy:6](../../../pangaea/copybooks/AUDITREC.cpy#L6) |
| 05 | ├ AU-TRAN-ID | number | 9 | 10 | 9(10) | | | [AUDITREC.cpy:7](../../../pangaea/copybooks/AUDITREC.cpy#L7) |
| 05 | ├ AU-ACCT-NO | number | 19 | 10 | 9(10) | | | [AUDITREC.cpy:8](../../../pangaea/copybooks/AUDITREC.cpy#L8) |
| 05 | ├ AU-TRAN-CODE | text | 29 | 2 | XX | | | [AUDITREC.cpy:9](../../../pangaea/copybooks/AUDITREC.cpy#L9) |
| 05 | ├ AU-AMOUNT | amount | 31 | 11 | S9(9)V99 | | | [AUDITREC.cpy:10](../../../pangaea/copybooks/AUDITREC.cpy#L10) |
| 05 | ├ AU-BALANCE-BEFORE | amount | 42 | 13 | S9(11)V99 | | | [AUDITREC.cpy:11](../../../pangaea/copybooks/AUDITREC.cpy#L11) |
| 05 | ├ AU-BALANCE-AFTER | amount | 55 | 13 | S9(11)V99 | | | [AUDITREC.cpy:12](../../../pangaea/copybooks/AUDITREC.cpy#L12) |
| 05 | └ FILLER | filler | 68 | 33 | X(33) | | | [AUDITREC.cpy:13](../../../pangaea/copybooks/AUDITREC.cpy#L13) |
