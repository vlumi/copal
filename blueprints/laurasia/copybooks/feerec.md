# feerec

```text
feerec: an interchange fee in the fee table, a relative file whose slot
is worked out from the card product, the region and the merchant group.

history
2012-01-09 ah  lcs-0140 original version.
```

Copybook [feerec.cpy](../../../laurasia/copybooks/feerec.cpy), used by [FEELD](../programs/FEELD.md), [IXFEE](../programs/IXFEE.md).

## FEE-SLOT-RECORD

8 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | FEE-SLOT-RECORD | group | 1 | 8 | | | | [feerec.cpy:7](../../../laurasia/copybooks/feerec.cpy#L7) |
| 05 | ├ FS-RATE | amount | 1 | 5 | 9V9(4) | | | [feerec.cpy:8](../../../laurasia/copybooks/feerec.cpy#L8) |
| 05 | └ FS-FIXED-CENTS | number | 6 | 3 | 9(5) | COMP-5 | | [feerec.cpy:9](../../../laurasia/copybooks/feerec.cpy#L9) |
