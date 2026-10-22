# dispcase

```text
dispcase: a cardholder's dispute of a presentment, a case in the
dispute master, with where it stands.

history
2015-03-23 mt  lcs-0302 original version.
2019-11-04 mt  lcs-0455 pre-arbitration.
```

Copybook [dispcase.cpy](../../../laurasia/copybooks/dispcase.cpy), used by [DISPUTE](../programs/DISPUTE.md).

## DISPUTE-CASE

48 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | DISPUTE-CASE | group | 1 | 48 | | | | [dispcase.cpy:8](../../../laurasia/copybooks/dispcase.cpy#L8) |
| 05 | ├ DC-CASE-NO | number | 1 | 8 | 9(8) | | | [dispcase.cpy:9](../../../laurasia/copybooks/dispcase.cpy#L9) |
| 05 | ├ DC-STATUS | code | 9 | 2 | X(2) | | | [dispcase.cpy:10](../../../laurasia/copybooks/dispcase.cpy#L10) |
| 88 | │ ├ DC-OPEN | | | | | | VALUE 'OP' | [dispcase.cpy:11](../../../laurasia/copybooks/dispcase.cpy#L11) |
| 88 | │ ├ DC-CHARGED-BACK | | | | | | VALUE 'CB' | [dispcase.cpy:12](../../../laurasia/copybooks/dispcase.cpy#L12) |
| 88 | │ ├ DC-REPRESENTED | | | | | | VALUE 'RP' | [dispcase.cpy:13](../../../laurasia/copybooks/dispcase.cpy#L13) |
| 88 | │ ├ DC-PRE-ARBITRATION | | | | | | VALUE 'PA' | [dispcase.cpy:14](../../../laurasia/copybooks/dispcase.cpy#L14) |
| 88 | │ ├ DC-WON | | | | | | VALUE 'WN' | [dispcase.cpy:15](../../../laurasia/copybooks/dispcase.cpy#L15) |
| 88 | │ └ DC-LOST | | | | | | VALUE 'LS' | [dispcase.cpy:16](../../../laurasia/copybooks/dispcase.cpy#L16) |
| 05 | ├ DC-MERCHANT | text | 11 | 10 | X(10) | | | [dispcase.cpy:17](../../../laurasia/copybooks/dispcase.cpy#L17) |
| 05 | ├ DC-AMOUNT | amount | 21 | 12 | S9(9)V99 | | SIGN LEADING SEPARATE | [dispcase.cpy:18](../../../laurasia/copybooks/dispcase.cpy#L18) |
| 05 | ├ DC-OPENED | date | 33 | 8 | 9(8) | | | [dispcase.cpy:19](../../../laurasia/copybooks/dispcase.cpy#L19) |
| 05 | └ DC-UPDATED | date | 41 | 8 | 9(8) | | | [dispcase.cpy:20](../../../laurasia/copybooks/dispcase.cpy#L20) |
