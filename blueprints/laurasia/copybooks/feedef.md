# feedef

```text
feedef: an interchange fee as the fee schedule gives it, a line each:
for a card product, a region and a merchant group, a rate on the amount
and a fixed fee in cents.

history
2012-01-09 ah  lcs-0140 original version.
```

Copybook [feedef.cpy](../../../laurasia/copybooks/feedef.cpy), used by [FEELD](../programs/FEELD.md).

## FEE-DEFINITION

14 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | FEE-DEFINITION | group | 1 | 14 | | | | [feedef.cpy:8](../../../laurasia/copybooks/feedef.cpy#L8) |
| 05 | ├ FD-PRODUCT | text | 1 | 2 | X(2) | | | [feedef.cpy:9](../../../laurasia/copybooks/feedef.cpy#L9) |
| 05 | ├ FD-REGION | text | 3 | 1 | X | | | [feedef.cpy:10](../../../laurasia/copybooks/feedef.cpy#L10) |
| 05 | ├ FD-GROUP | number | 4 | 1 | 9 | | | [feedef.cpy:11](../../../laurasia/copybooks/feedef.cpy#L11) |
| 05 | ├ FD-RATE | amount | 5 | 5 | 9V9(4) | | | [feedef.cpy:12](../../../laurasia/copybooks/feedef.cpy#L12) |
| 05 | └ FD-FIXED | number | 10 | 5 | 9(5) | | | [feedef.cpy:13](../../../laurasia/copybooks/feedef.cpy#L13) |
