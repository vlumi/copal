# RSVTB

```text
RSVTB: 責任準備金係数 レコード (13バイト)

RSVTB.TXT の1行が商品・性別・加入年齢区分・経過年数ごとの、
その年度末に積み立てる保険金額1,000円あたりの責任準備金。
年齢区分は解約返戻金係数と同じ: 20歳未満が 1、以後10歳
ごと、60歳以上が 6。経過年数は 1～50。特約の商品もある。

変更履歴
日付       担当 変更
1987-04-01 KT   GL-0001 初版。
2026-11-19 YN   GL-0829 性別と加入年齢区分。
```

Copybook [RSVTB.cpy](../../../gondwana/copybooks/RSVTB.cpy), used by [GLM020](../programs/GLM020.md).

## RS-ｼﾞｭﾝﾋﾞｷﾝ-ﾚｺｰﾄﾞ

13 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | RS-ｼﾞｭﾝﾋﾞｷﾝ-ﾚｺｰﾄﾞ | group | 1 | 13 | | | | [RSVTB.cpy:14](../../../gondwana/copybooks/RSVTB.cpy#L14) |
| 05 | ├ RS-ｼｮｳﾋﾝ | text | 1 | 2 | X(2) | | | [RSVTB.cpy:15](../../../gondwana/copybooks/RSVTB.cpy#L15) |
| 05 | ├ RS-ｾｲﾍﾞﾂ | number | 3 | 1 | 9 | | | [RSVTB.cpy:16](../../../gondwana/copybooks/RSVTB.cpy#L16) |
| 05 | ├ RS-ﾈﾝﾚｲｸﾌﾞﾝ | number | 4 | 1 | 9 | | | [RSVTB.cpy:17](../../../gondwana/copybooks/RSVTB.cpy#L17) |
| 05 | ├ RS-ｹｲｶﾈﾝｽｳ | number | 5 | 2 | 99 | | | [RSVTB.cpy:18](../../../gondwana/copybooks/RSVTB.cpy#L18) |
| 05 | └ RS-ｹｲｽｳ | rate | 7 | 7 | 9(4)V9(3) | | | [RSVTB.cpy:19](../../../gondwana/copybooks/RSVTB.cpy#L19) |
