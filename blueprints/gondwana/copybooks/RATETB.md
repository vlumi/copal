# RATETB

```text
RATETB: 保険料率 レコード (10バイト)

RATETB.TXT の1行が商品・性別・加入年齢ごとの、保険金額
1,000円あたりの月払保険料。プログラムは商品・性別・年齢
(0～80) のテーブルに読み込む。

変更履歴
日付       担当 変更
1987-04-01 KT   GL-0001 初版。
```

Copybook [RATETB.cpy](../../../gondwana/copybooks/RATETB.cpy), used by [GLD010](../programs/GLD010.md), [GLD020](../programs/GLD020.md).

## RT-ﾘｮｳﾘﾂ-ﾚｺｰﾄﾞ

10 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | RT-ﾘｮｳﾘﾂ-ﾚｺｰﾄﾞ | group | 1 | 10 | | | | [RATETB.cpy:12](../../../gondwana/copybooks/RATETB.cpy#L12) |
| 05 | ├ RT-ｼｮｳﾋﾝ | text | 1 | 2 | X(2) | | | [RATETB.cpy:13](../../../gondwana/copybooks/RATETB.cpy#L13) |
| 05 | ├ RT-ｾｲﾍﾞﾂ | number | 3 | 1 | 9 | | | [RATETB.cpy:14](../../../gondwana/copybooks/RATETB.cpy#L14) |
| 05 | ├ RT-ﾈﾝﾚｲ | number | 4 | 2 | 99 | | | [RATETB.cpy:15](../../../gondwana/copybooks/RATETB.cpy#L15) |
| 05 | └ RT-ﾘｮｳﾘﾂ | rate | 6 | 5 | 9(3)V99 | | | [RATETB.cpy:16](../../../gondwana/copybooks/RATETB.cpy#L16) |
