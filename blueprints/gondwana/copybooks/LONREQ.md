# LONREQ

```text
LONREQ: 契約者貸付依頼 レコード (60バイト)

支社が入力した貸付または返済1件が LONREQ.YYYYMMDD.TXT の
1行。貸付は契約者の口座に送金し、返済は窓口で収納した額。

変更履歴
日付       担当 変更
2012-09-03 YN   GL-0561 初版、支社の貸付台帳をシステム化。
```

Copybook [LONREQ.cpy](../../../gondwana/copybooks/LONREQ.cpy), used by [GLD070](../programs/GLD070.md).

## LR-ｶｼﾂｹ-ﾚｺｰﾄﾞ

60 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | LR-ｶｼﾂｹ-ﾚｺｰﾄﾞ | group | 1 | 60 | | | | [LONREQ.cpy:11](../../../gondwana/copybooks/LONREQ.cpy#L11) |
| 05 | ├ LR-ｲﾗｲﾊﾞﾝｺﾞｳ | text | 1 | 8 | X(8) | | | [LONREQ.cpy:12](../../../gondwana/copybooks/LONREQ.cpy#L12) |
| 05 | ├ LR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ | text | 9 | 10 | X(10) | | | [LONREQ.cpy:13](../../../gondwana/copybooks/LONREQ.cpy#L13) |
| 05 | ├ LR-ｼｭﾍﾞﾂ | code | 19 | 1 | X | | | [LONREQ.cpy:14](../../../gondwana/copybooks/LONREQ.cpy#L14) |
| 88 | │ ├ LR-ｶｼﾂｹ | | | | | | VALUE 'L' | [LONREQ.cpy:15](../../../gondwana/copybooks/LONREQ.cpy#L15) |
| 88 | │ └ LR-ﾍﾝｻｲ | | | | | | VALUE 'P' | [LONREQ.cpy:16](../../../gondwana/copybooks/LONREQ.cpy#L16) |
| 05 | ├ LR-ｷﾝｶﾞｸ | amount | 20 | 9 | 9(9) | | | [LONREQ.cpy:17](../../../gondwana/copybooks/LONREQ.cpy#L17) |
| 05 | ├ LR-ｼﾃﾝ | text | 29 | 4 | X(4) | | | [LONREQ.cpy:18](../../../gondwana/copybooks/LONREQ.cpy#L18) |
| 05 | └ FILLER | filler | 33 | 28 | X(28) | | | [LONREQ.cpy:19](../../../gondwana/copybooks/LONREQ.cpy#L19) |
