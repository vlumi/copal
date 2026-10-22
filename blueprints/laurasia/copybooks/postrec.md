# postrec

```text
postrec: a posting to a card's account, a record a line.

history
2011-04-04 rk  lcs-0102 original version.
2022-09-05 ah  lcs-0497 payments, as credits.
```

Copybook [postrec.cpy](../../../laurasia/copybooks/postrec.cpy), used by [CYCLE](../programs/CYCLE.md), [PAYIN](../programs/PAYIN.md), [POST](../programs/POST.md).

## POSTING

51 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | POSTING | group | 1 | 51 | | | | [postrec.cpy:7](../../../laurasia/copybooks/postrec.cpy#L7) |
| 05 | ├ PO-CARD-NO | text | 1 | 16 | X(16) | | | [postrec.cpy:8](../../../laurasia/copybooks/postrec.cpy#L8) |
| 05 | ├ PO-ACCOUNT-NO | text | 17 | 10 | X(10) | | | [postrec.cpy:9](../../../laurasia/copybooks/postrec.cpy#L9) |
| 05 | ├ PO-KIND | code | 27 | 1 | X | | | [postrec.cpy:10](../../../laurasia/copybooks/postrec.cpy#L10) |
| 88 | │ ├ PO-PURCHASE | | | | | | VALUE 'P' | [postrec.cpy:11](../../../laurasia/copybooks/postrec.cpy#L11) |
| 88 | │ ├ PO-OVER-LIMIT-FEE | | | | | | VALUE 'F' | [postrec.cpy:12](../../../laurasia/copybooks/postrec.cpy#L12) |
| 88 | │ ├ PO-INTEREST | | | | | | VALUE 'I' | [postrec.cpy:13](../../../laurasia/copybooks/postrec.cpy#L13) |
| 88 | │ └ PO-PAYMENT | | | | | | VALUE 'C' | [postrec.cpy:14](../../../laurasia/copybooks/postrec.cpy#L14) |
| 05 | ├ PO-AMOUNT | amount | 28 | 10 | S9(7)V99 | | SIGN LEADING SEPARATE | [postrec.cpy:15](../../../laurasia/copybooks/postrec.cpy#L15) |
| 05 | ├ PO-REFERENCE | text | 38 | 6 | X(6) | | | [postrec.cpy:16](../../../laurasia/copybooks/postrec.cpy#L16) |
| 05 | └ PO-DATE | date | 44 | 8 | 9(8) | | | [postrec.cpy:17](../../../laurasia/copybooks/postrec.cpy#L17) |
