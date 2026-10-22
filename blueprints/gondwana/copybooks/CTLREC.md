# CTLREC

```text
CTLREC: システム制御レコード (64バイト、バイナリ)

1レコードのみ、GLSCTL.BIN。システムの業務日付を両暦で、
最終請求年月、次に採番する支払番号と仕訳番号、および
最終アンロード時の件数を持つ。カウンタは NetCOBOL が
PC上で配置する通りの内部バイナリ、境界合わせ、
リトルエンディアン。

変更履歴
日付       担当 変更
1999-11-08 MS   GL-0233 初版、PC システム移行時。
```

Copybook [CTLREC.cpy](../../../gondwana/copybooks/CTLREC.cpy), used by [GLD010](../programs/GLD010.md), [GLD020](../programs/GLD020.md), [GLD030](../programs/GLD030.md), [GLD040](../programs/GLD040.md), [GLD050](../programs/GLD050.md), [GLD090](../programs/GLD090.md), [GLM010](../programs/GLM010.md), [GLM020](../programs/GLM020.md), [GLM030](../programs/GLM030.md), [GLO010](../programs/GLO010.md), [GLY010](../programs/GLY010.md).

## CT-ｾｲｷﾞｮ-ﾚｺｰﾄﾞ

64 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | CT-ｾｲｷﾞｮ-ﾚｺｰﾄﾞ | group | 1 | 64 | | | | [CTLREC.cpy:14](../../../gondwana/copybooks/CTLREC.cpy#L14) |
| 05 | ├ CT-ｷﾞｮｳﾑﾋﾂﾞｹ | date | 1 | 8 | 9(8) | | | [CTLREC.cpy:15](../../../gondwana/copybooks/CTLREC.cpy#L15) |
| 05 | ├ CT-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ | date | 9 | 7 | X(7) | | | [CTLREC.cpy:16](../../../gondwana/copybooks/CTLREC.cpy#L16) |
| 05 | ├ FILLER | filler | 16 | 1 | X | | | [CTLREC.cpy:17](../../../gondwana/copybooks/CTLREC.cpy#L17) |
| 05 | ├ CT-ｻｲｼｭｳｾｲｷｭｳﾈﾝｹﾞﾂ | date | 17 | 6 | 9(6) | | | [CTLREC.cpy:18](../../../gondwana/copybooks/CTLREC.cpy#L18) |
| 05 | ├ FILLER | filler | 23 | 2 | X(2) | | | [CTLREC.cpy:19](../../../gondwana/copybooks/CTLREC.cpy#L19) |
| 05 | ├ CT-ｼﾞｶｲｼﾊﾗｲﾊﾞﾝｺﾞｳ | number | 25 | 4 | S9(9) | COMP-5 | | [CTLREC.cpy:20](../../../gondwana/copybooks/CTLREC.cpy#L20) |
| 05 | ├ CT-ｼﾞｶｲｼﾜｹﾊﾞﾝｺﾞｳ | number | 29 | 4 | S9(9) | COMP-5 | | [CTLREC.cpy:21](../../../gondwana/copybooks/CTLREC.cpy#L21) |
| 05 | ├ CT-ﾕｳｺｳｹﾝｽｳ | number | 33 | 4 | S9(9) | COMP-5 | | [CTLREC.cpy:22](../../../gondwana/copybooks/CTLREC.cpy#L22) |
| 05 | ├ CT-ｿｳｹﾝｽｳ | number | 37 | 4 | S9(9) | COMP-5 | | [CTLREC.cpy:23](../../../gondwana/copybooks/CTLREC.cpy#L23) |
| 05 | ├ CT-ﾕｳｺｳﾎｹﾝｷﾝｶﾞｸ | amount | 41 | 8 | S9(18) | COMP-5 | | [CTLREC.cpy:24](../../../gondwana/copybooks/CTLREC.cpy#L24) |
| 05 | └ FILLER | filler | 49 | 16 | X(16) | | | [CTLREC.cpy:25](../../../gondwana/copybooks/CTLREC.cpy#L25) |
