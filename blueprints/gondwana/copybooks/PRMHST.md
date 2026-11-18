# PRMHST

```text
PRMHST: 収納履歴 レコード (60バイト)

収納した保険料1件につき1レコード。キーは証券番号と契約内の
連番。第1回保険料、口座振替、復活時の窓口収納、自動振替貸付
による充当を残す。書き込みは GLC050 が行う。
収納日は西暦 YYYYMMDD: 保険料控除証明書が暦年で集計する。
回数は充当した保険料の回数、払込済年月は収納後のもの。

変更履歴
日付       担当 変更
2026-11-19 YN   GL-0827 初版、保険料控除証明書のため。
```

Copybook [PRMHST.cpy](../../../gondwana/copybooks/PRMHST.cpy), used by [GLC050](../programs/GLC050.md), [GLY030](../programs/GLY030.md).

## SH-ｼｭｳﾉｳ-ﾚｺｰﾄﾞ

60 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | SH-ｼｭｳﾉｳ-ﾚｺｰﾄﾞ | group | 1 | 60 | | | | [PRMHST.cpy:14](../../../gondwana/copybooks/PRMHST.cpy#L14) |
| 05 | ├ SH-ｷｰ | group | 1 | 14 | | | | [PRMHST.cpy:15](../../../gondwana/copybooks/PRMHST.cpy#L15) |
| 10 | │ ├ SH-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ | text | 1 | 10 | X(10) | | | [PRMHST.cpy:16](../../../gondwana/copybooks/PRMHST.cpy#L16) |
| 10 | │ └ SH-ﾚﾝﾊﾞﾝ | number | 11 | 4 | 9(4) | | | [PRMHST.cpy:17](../../../gondwana/copybooks/PRMHST.cpy#L17) |
| 05 | ├ SH-ｼｭｳﾉｳﾋﾞ | date | 15 | 8 | 9(8) | | | [PRMHST.cpy:18](../../../gondwana/copybooks/PRMHST.cpy#L18) |
| 05 | ├ SH-ｼｭﾍﾞﾂ | code | 23 | 1 | X | | | [PRMHST.cpy:19](../../../gondwana/copybooks/PRMHST.cpy#L19) |
| 88 | │ ├ SH-ｼｮｶｲ | | | | | | VALUE '1' | [PRMHST.cpy:20](../../../gondwana/copybooks/PRMHST.cpy#L20) |
| 88 | │ ├ SH-ﾌﾘｶｴ | | | | | | VALUE '2' | [PRMHST.cpy:21](../../../gondwana/copybooks/PRMHST.cpy#L21) |
| 88 | │ ├ SH-ﾏﾄﾞｸﾞﾁ | | | | | | VALUE '3' | [PRMHST.cpy:22](../../../gondwana/copybooks/PRMHST.cpy#L22) |
| 88 | │ └ SH-ｶｼﾂｹ | | | | | | VALUE '4' | [PRMHST.cpy:23](../../../gondwana/copybooks/PRMHST.cpy#L23) |
| 05 | ├ SH-ｷﾝｶﾞｸ | amount | 24 | 9 | 9(9) | | | [PRMHST.cpy:24](../../../gondwana/copybooks/PRMHST.cpy#L24) |
| 05 | ├ SH-ｶｲｽｳ | number | 33 | 2 | 99 | | | [PRMHST.cpy:25](../../../gondwana/copybooks/PRMHST.cpy#L25) |
| 05 | ├ SH-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ | date | 35 | 6 | 9(6) | | | [PRMHST.cpy:26](../../../gondwana/copybooks/PRMHST.cpy#L26) |
| 05 | └ FILLER | filler | 41 | 20 | X(20) | | | [PRMHST.cpy:27](../../../gondwana/copybooks/PRMHST.cpy#L27) |
