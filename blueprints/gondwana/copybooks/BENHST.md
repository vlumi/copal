# BENHST

```text
BENHST: 給付履歴 レコード (60バイト)

支払った入院給付金と手術給付金1件が1レコード。索引ファイル
BENHST.IDX では証券番号と契約内の連番をキーとし、日次の
アンロードでは同じレコードが1行である。入院の通算日数と、
同じ入院・手術の二重請求の検査に使う。手術は日数欄に倍率。

変更履歴
日付       担当 変更
2026-11-16 YN   GL-0822 初版。
```

Copybook [BENHST.cpy](../../../gondwana/copybooks/BENHST.cpy), used by [GLD030](../programs/GLD030.md), [GLD090](../programs/GLD090.md).

## BH-ﾘﾚｷ-ﾚｺｰﾄﾞ

60 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | BH-ﾘﾚｷ-ﾚｺｰﾄﾞ | group | 1 | 60 | | | | [BENHST.cpy:13](../../../gondwana/copybooks/BENHST.cpy#L13) |
| 05 | ├ BH-ｷｰ | group | 1 | 13 | | | | [BENHST.cpy:14](../../../gondwana/copybooks/BENHST.cpy#L14) |
| 10 | │ ├ BH-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ | text | 1 | 10 | X(10) | | | [BENHST.cpy:15](../../../gondwana/copybooks/BENHST.cpy#L15) |
| 10 | │ └ BH-ﾚﾝﾊﾞﾝ | number | 11 | 3 | 9(3) | | | [BENHST.cpy:16](../../../gondwana/copybooks/BENHST.cpy#L16) |
| 05 | ├ BH-ｼｭﾍﾞﾂ | code | 14 | 1 | X | | | [BENHST.cpy:17](../../../gondwana/copybooks/BENHST.cpy#L17) |
| 88 | │ ├ BH-ﾆｭｳｲﾝ | | | | | | VALUE 'H' | [BENHST.cpy:18](../../../gondwana/copybooks/BENHST.cpy#L18) |
| 88 | │ └ BH-ｼｭｼﾞｭﾂ | | | | | | VALUE 'O' | [BENHST.cpy:19](../../../gondwana/copybooks/BENHST.cpy#L19) |
| 05 | ├ BH-ｼﾞﾕｳﾊｯｾｲﾋﾞ | date | 15 | 7 | X(7) | | | [BENHST.cpy:20](../../../gondwana/copybooks/BENHST.cpy#L20) |
| 05 | ├ BH-ﾆｯｽｳ | number | 22 | 3 | 9(3) | | | [BENHST.cpy:21](../../../gondwana/copybooks/BENHST.cpy#L21) |
| 05 | ├ BH-ｷﾝｶﾞｸ | amount | 25 | 9 | 9(9) | | | [BENHST.cpy:22](../../../gondwana/copybooks/BENHST.cpy#L22) |
| 05 | ├ BH-ｼﾊﾗｲﾋﾞ | date | 34 | 7 | X(7) | | | [BENHST.cpy:23](../../../gondwana/copybooks/BENHST.cpy#L23) |
| 05 | ├ BH-ｾｲｷｭｳﾊﾞﾝｺﾞｳ | text | 41 | 8 | X(8) | | | [BENHST.cpy:24](../../../gondwana/copybooks/BENHST.cpy#L24) |
| 05 | └ FILLER | filler | 49 | 12 | X(12) | | | [BENHST.cpy:25](../../../gondwana/copybooks/BENHST.cpy#L25) |
