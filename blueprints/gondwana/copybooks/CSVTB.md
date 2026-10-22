# CSVTB

```text
CSVTB: 解約返戻金係数 レコード (12バイト)

CSVTB.TXT の1行が商品・加入年齢区分・経過年数ごとの、
保険金額1,000円あたりの解約返戻金。年齢区分は 1 が20歳
未満、2 が20歳代、3 が30歳代、4 が40歳代、5 が50歳代、
6 が60歳以上。経過年数は 1～40。

変更履歴
日付       担当 変更
1987-04-01 KT   GL-0001 初版。
```

Copybook [CSVTB.cpy](../../../gondwana/copybooks/CSVTB.cpy), used by [GLD030](../programs/GLD030.md).

## CV-ｶｲﾔｸﾍﾝﾚｲｷﾝ-ﾚｺｰﾄﾞ

12 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | CV-ｶｲﾔｸﾍﾝﾚｲｷﾝ-ﾚｺｰﾄﾞ | group | 1 | 12 | | | | [CSVTB.cpy:13](../../../gondwana/copybooks/CSVTB.cpy#L13) |
| 05 | ├ CV-ｼｮｳﾋﾝ | text | 1 | 2 | X(2) | | | [CSVTB.cpy:14](../../../gondwana/copybooks/CSVTB.cpy#L14) |
| 05 | ├ CV-ﾈﾝﾚｲｸﾌﾞﾝ | number | 3 | 1 | 9 | | | [CSVTB.cpy:15](../../../gondwana/copybooks/CSVTB.cpy#L15) |
| 05 | ├ CV-ｹｲｶﾈﾝｽｳ | number | 4 | 2 | 99 | | | [CSVTB.cpy:16](../../../gondwana/copybooks/CSVTB.cpy#L16) |
| 05 | └ CV-ｹｲｽｳ | amount | 6 | 7 | 9(4)V9(3) | | | [CSVTB.cpy:17](../../../gondwana/copybooks/CSVTB.cpy#L17) |
