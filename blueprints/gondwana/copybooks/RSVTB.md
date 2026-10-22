# RSVTB

```text
RSVTB: 責任準備金係数 レコード (11バイト)

RSVTB.TXT の1行が商品・経過年数ごとの、その年度末に
積み立てる保険金額1,000円あたりの責任準備金。経過年数は
1～50。

変更履歴
日付       担当 変更
1987-04-01 KT   GL-0001 初版。
```

Copybook [RSVTB.cpy](../../../gondwana/copybooks/RSVTB.cpy), used by [GLM020](../programs/GLM020.md).

## RS-ｼﾞｭﾝﾋﾞｷﾝ-ﾚｺｰﾄﾞ

11 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | RS-ｼﾞｭﾝﾋﾞｷﾝ-ﾚｺｰﾄﾞ | group | 1 | 11 | | | | [RSVTB.cpy:12](../../../gondwana/copybooks/RSVTB.cpy#L12) |
| 05 | ├ RS-ｼｮｳﾋﾝ | text | 1 | 2 | X(2) | | | [RSVTB.cpy:13](../../../gondwana/copybooks/RSVTB.cpy#L13) |
| 05 | ├ RS-ｹｲｶﾈﾝｽｳ | number | 3 | 2 | 99 | | | [RSVTB.cpy:14](../../../gondwana/copybooks/RSVTB.cpy#L14) |
| 05 | └ RS-ｹｲｽｳ | rate | 5 | 7 | 9(4)V9(3) | | | [RSVTB.cpy:15](../../../gondwana/copybooks/RSVTB.cpy#L15) |
