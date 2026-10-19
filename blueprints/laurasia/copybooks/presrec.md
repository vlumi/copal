# presrec

```text
presrec: a presentment on its way through the night, a record a line.
clrin fills in where it came from, match its authorization's outcome,
ixfee the interchange fee.

history
2011-04-04 rk  lcs-0102 original version.
2014-02-17 rk  lcs-0251 match status and interchange fee.
```

Copybook [presrec.cpy](../../../laurasia/copybooks/presrec.cpy), used by [CLRIN](../programs/CLRIN.md), [IXFEE](../programs/IXFEE.md), [MATCH](../programs/MATCH.md), [POST](../programs/POST.md), [SETTLE](../programs/SETTLE.md).

## PRES-RECORD

76 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | PRES-RECORD | group | 1 | 76 | | | | [presrec.cpy:9](../../../laurasia/copybooks/presrec.cpy#L9) |
| 05 | ├ PS-CARD-NO | text | 1 | 16 | X(16) | | | [presrec.cpy:10](../../../laurasia/copybooks/presrec.cpy#L10) |
| 05 | ├ PS-AMOUNT | amount | 17 | 12 | S9(9)V99 | | SIGN LEADING SEPARATE | [presrec.cpy:11](../../../laurasia/copybooks/presrec.cpy#L11) |
| 05 | ├ PS-AUTH-CODE | text | 29 | 6 | X(6) | | | [presrec.cpy:12](../../../laurasia/copybooks/presrec.cpy#L12) |
| 05 | ├ PS-MERCHANT | text | 35 | 10 | X(10) | | | [presrec.cpy:13](../../../laurasia/copybooks/presrec.cpy#L13) |
| 05 | ├ PS-MCC | number | 45 | 4 | 9(4) | | | [presrec.cpy:14](../../../laurasia/copybooks/presrec.cpy#L14) |
| 05 | ├ PS-REGION | code | 49 | 1 | X | | | [presrec.cpy:15](../../../laurasia/copybooks/presrec.cpy#L15) |
| 88 | │ ├ PS-DOMESTIC | | | | | | VALUE 'D' | [presrec.cpy:16](../../../laurasia/copybooks/presrec.cpy#L16) |
| 88 | │ └ PS-INTERNATIONAL | | | | | | VALUE 'I' | [presrec.cpy:17](../../../laurasia/copybooks/presrec.cpy#L17) |
| 05 | ├ PS-PRODUCT | code | 50 | 2 | X(2) | | | [presrec.cpy:18](../../../laurasia/copybooks/presrec.cpy#L18) |
| 88 | │ ├ PS-CLASSIC | | | | | | VALUE 'CL' | [presrec.cpy:19](../../../laurasia/copybooks/presrec.cpy#L19) |
| 88 | │ ├ PS-GOLD | | | | | | VALUE 'GD' | [presrec.cpy:20](../../../laurasia/copybooks/presrec.cpy#L20) |
| 88 | │ └ PS-BUSINESS | | | | | | VALUE 'BZ' | [presrec.cpy:21](../../../laurasia/copybooks/presrec.cpy#L21) |
| 05 | ├ PS-TXN-DATE | date | 52 | 8 | 9(8) | | | [presrec.cpy:22](../../../laurasia/copybooks/presrec.cpy#L22) |
| 05 | ├ PS-BATCH-NO | number | 60 | 6 | 9(6) | | | [presrec.cpy:23](../../../laurasia/copybooks/presrec.cpy#L23) |
| 05 | ├ PS-MATCH | code | 66 | 1 | X | | | [presrec.cpy:24](../../../laurasia/copybooks/presrec.cpy#L24) |
| 88 | │ ├ PS-UNMATCHED | | | | | | VALUE ' ' | [presrec.cpy:25](../../../laurasia/copybooks/presrec.cpy#L25) |
| 88 | │ ├ PS-MATCHED | | | | | | VALUE 'M' | [presrec.cpy:26](../../../laurasia/copybooks/presrec.cpy#L26) |
| 88 | │ ├ PS-DECLINED | | | | | | VALUE 'D' | [presrec.cpy:27](../../../laurasia/copybooks/presrec.cpy#L27) |
| 88 | │ ├ PS-NO-AUTH | | | | | | VALUE 'U' | [presrec.cpy:28](../../../laurasia/copybooks/presrec.cpy#L28) |
| 88 | │ └ PS-OVER-AUTH | | | | | | VALUE 'X' | [presrec.cpy:29](../../../laurasia/copybooks/presrec.cpy#L29) |
| 05 | └ PS-FEE | amount | 67 | 10 | S9(7)V99 | | SIGN LEADING SEPARATE | [presrec.cpy:30](../../../laurasia/copybooks/presrec.cpy#L30) |
