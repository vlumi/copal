# AGTMST

```text
AGTMST: 代理店マスタ レコード (60バイト)

AGTMST.TXT の1行が1代理店、手数料計算用。入社年月は
当時の代理店システムが書いた通りの YYMM で、元号も世紀
もない。

変更履歴
日付       担当 変更
1987-04-01 KT   GL-0001 初版。
```

Copybook [AGTMST.cpy](../../../gondwana/copybooks/AGTMST.cpy), used by [GLM030](../programs/GLM030.md).

## AG-ﾀﾞｲﾘﾃﾝ-ﾚｺｰﾄﾞ

60 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | AG-ﾀﾞｲﾘﾃﾝ-ﾚｺｰﾄﾞ | group | 1 | 60 | | | | [AGTMST.cpy:12](../../../gondwana/copybooks/AGTMST.cpy#L12) |
| 05 | ├ AG-ｺｰﾄﾞ | text | 1 | 6 | X(6) | | | [AGTMST.cpy:13](../../../gondwana/copybooks/AGTMST.cpy#L13) |
| 05 | ├ AG-ｼﾒｲ | text | 7 | 20 | X(20) | | | [AGTMST.cpy:14](../../../gondwana/copybooks/AGTMST.cpy#L14) |
| 05 | ├ AG-ｼﾃﾝ | text | 27 | 4 | X(4) | | | [AGTMST.cpy:15](../../../gondwana/copybooks/AGTMST.cpy#L15) |
| 05 | ├ AG-ﾆｭｳｼｬﾈﾝｹﾞﾂ | number | 31 | 4 | 9(4) | | | [AGTMST.cpy:16](../../../gondwana/copybooks/AGTMST.cpy#L16) |
| 05 | ├ AG-ｼｶｸ | number | 35 | 1 | 9 | | | [AGTMST.cpy:17](../../../gondwana/copybooks/AGTMST.cpy#L17) |
| 05 | └ FILLER | filler | 36 | 25 | X(25) | | | [AGTMST.cpy:18](../../../gondwana/copybooks/AGTMST.cpy#L18) |
