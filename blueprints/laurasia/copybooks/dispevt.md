# dispevt

```text
dispevt: an event in a dispute's life, from the cardholder's claim to
the network's ruling.

history
2015-03-23 mt  lcs-0302 original version.
2019-11-04 mt  lcs-0455 pre-arbitration.
```

Copybook [dispevt.cpy](../../../laurasia/copybooks/dispevt.cpy), used by [DISPUTE](../programs/DISPUTE.md).

## DISPUTE-EVENT

40 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | DISPUTE-EVENT | group | 1 | 40 | | | | [dispevt.cpy:8](../../../laurasia/copybooks/dispevt.cpy#L8) |
| 05 | ├ DE-CASE-NO | number | 1 | 8 | 9(8) | | | [dispevt.cpy:9](../../../laurasia/copybooks/dispevt.cpy#L9) |
| 05 | ├ DE-EVENT | code | 9 | 2 | X(2) | | | [dispevt.cpy:10](../../../laurasia/copybooks/dispevt.cpy#L10) |
| 88 | │ ├ DE-OPENED | | | | | | VALUE 'OP' | [dispevt.cpy:11](../../../laurasia/copybooks/dispevt.cpy#L11) |
| 88 | │ ├ DE-CHARGEBACK | | | | | | VALUE 'CB' | [dispevt.cpy:12](../../../laurasia/copybooks/dispevt.cpy#L12) |
| 88 | │ ├ DE-REPRESENTMENT | | | | | | VALUE 'RP' | [dispevt.cpy:13](../../../laurasia/copybooks/dispevt.cpy#L13) |
| 88 | │ ├ DE-PRE-ARBITRATION | | | | | | VALUE 'PA' | [dispevt.cpy:14](../../../laurasia/copybooks/dispevt.cpy#L14) |
| 88 | │ ├ DE-MERCHANT-ACCEPTS | | | | | | VALUE 'AC' | [dispevt.cpy:15](../../../laurasia/copybooks/dispevt.cpy#L15) |
| 88 | │ ├ DE-ISSUER-DROPS | | | | | | VALUE 'DR' | [dispevt.cpy:16](../../../laurasia/copybooks/dispevt.cpy#L16) |
| 88 | │ ├ DE-RULED-FOR-ISSUER | | | | | | VALUE 'RI' | [dispevt.cpy:17](../../../laurasia/copybooks/dispevt.cpy#L17) |
| 88 | │ └ DE-RULED-FOR-MERCHANT | | | | | | VALUE 'RM' | [dispevt.cpy:18](../../../laurasia/copybooks/dispevt.cpy#L18) |
| 05 | ├ DE-MERCHANT | text | 11 | 10 | X(10) | | | [dispevt.cpy:19](../../../laurasia/copybooks/dispevt.cpy#L19) |
| 05 | ├ DE-AMOUNT | amount | 21 | 12 | S9(9)V99 | | SIGN LEADING SEPARATE | [dispevt.cpy:20](../../../laurasia/copybooks/dispevt.cpy#L20) |
| 05 | └ DE-DATE | date | 33 | 8 | 9(8) | | | [dispevt.cpy:21](../../../laurasia/copybooks/dispevt.cpy#L21) |
