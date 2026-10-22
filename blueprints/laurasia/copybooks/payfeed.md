# payfeed

```text
payfeed: the card system's daily feed, a record a line: a cardholder's
payment, with the bank's reference, or a change of a card's status.

history
2022-09-05 ah  lcs-0497 original version.
```

Copybook [payfeed.cpy](../../../laurasia/copybooks/payfeed.cpy), used by [PAYIN](../programs/PAYIN.md).

## FEED-RECORD

42 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | FEED-RECORD | group | 1 | 42 | | | | [payfeed.cpy:7](../../../laurasia/copybooks/payfeed.cpy#L7) |
| 05 | ├ PF-CARD-NO | text | 1 | 16 | X(16) | | | [payfeed.cpy:8](../../../laurasia/copybooks/payfeed.cpy#L8) |
| 05 | ├ PF-KIND | code | 17 | 1 | X | | | [payfeed.cpy:9](../../../laurasia/copybooks/payfeed.cpy#L9) |
| 88 | │ ├ PF-PAYMENT | | | | | | VALUE 'P' | [payfeed.cpy:10](../../../laurasia/copybooks/payfeed.cpy#L10) |
| 88 | │ └ PF-STATUS-CHANGE | | | | | | VALUE 'S' | [payfeed.cpy:11](../../../laurasia/copybooks/payfeed.cpy#L11) |
| 05 | ├ PF-AMOUNT | amount | 18 | 10 | S9(7)V99 | | SIGN LEADING SEPARATE | [payfeed.cpy:12](../../../laurasia/copybooks/payfeed.cpy#L12) |
| 05 | ├ PF-STATUS | text | 28 | 1 | X | | | [payfeed.cpy:13](../../../laurasia/copybooks/payfeed.cpy#L13) |
| 05 | ├ PF-REFERENCE | text | 29 | 6 | X(6) | | | [payfeed.cpy:14](../../../laurasia/copybooks/payfeed.cpy#L14) |
| 05 | └ PF-DATE | date | 35 | 8 | 9(8) | | | [payfeed.cpy:15](../../../laurasia/copybooks/payfeed.cpy#L15) |
