# GLJRNL

```text
GLJRNL: 仕訳 レコード (80バイト)

総勘定元帳向けの仕訳1件が GLJRNL.YYYYMMDD.TXT の1行:
起票したプログラム、勘定科目、貸借と金額、元の契約。
勘定科目は 4110 保険料収入、5110 死亡保険金、5120 入院
給付金、5130 解約返戻金、5140 満期保険金、5210 責任準備金
繰入、5310 契約者配当金、5410 代理店手数料。

変更履歴
日付       担当 変更
1987-04-01 KT   GL-0001 初版。
```

Copybook [GLJRNL.cpy](../../../gondwana/copybooks/GLJRNL.cpy), used by [GLD010](../programs/GLD010.md), [GLD020](../programs/GLD020.md), [GLD030](../programs/GLD030.md), [GLD040](../programs/GLD040.md), [GLD050](../programs/GLD050.md), [GLM020](../programs/GLM020.md), [GLM030](../programs/GLM030.md), [GLY010](../programs/GLY010.md).

## JR-ｼﾜｹ-ﾚｺｰﾄﾞ

80 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | JR-ｼﾜｹ-ﾚｺｰﾄﾞ | group | 1 | 80 | | | | [GLJRNL.cpy:14](../../../gondwana/copybooks/GLJRNL.cpy#L14) |
| 05 | ├ JR-ｼﾜｹﾊﾞﾝｺﾞｳ | number | 1 | 9 | 9(9) | | | [GLJRNL.cpy:15](../../../gondwana/copybooks/GLJRNL.cpy#L15) |
| 05 | ├ JR-ﾋﾂﾞｹ | date | 10 | 7 | X(7) | | | [GLJRNL.cpy:16](../../../gondwana/copybooks/GLJRNL.cpy#L16) |
| 05 | ├ JR-ﾌﾟﾛｸﾞﾗﾑ | text | 17 | 6 | X(6) | | | [GLJRNL.cpy:17](../../../gondwana/copybooks/GLJRNL.cpy#L17) |
| 05 | ├ JR-ｶﾝｼﾞｮｳｶﾓｸ | number | 23 | 4 | 9(4) | | | [GLJRNL.cpy:18](../../../gondwana/copybooks/GLJRNL.cpy#L18) |
| 05 | ├ JR-ﾀｲｼｬｸ | code | 27 | 1 | X | | | [GLJRNL.cpy:19](../../../gondwana/copybooks/GLJRNL.cpy#L19) |
| 88 | │ ├ JR-ｶﾘｶﾀ | | | | | | VALUE 'D' | [GLJRNL.cpy:20](../../../gondwana/copybooks/GLJRNL.cpy#L20) |
| 88 | │ └ JR-ｶｼｶﾀ | | | | | | VALUE 'C' | [GLJRNL.cpy:21](../../../gondwana/copybooks/GLJRNL.cpy#L21) |
| 05 | ├ JR-ｷﾝｶﾞｸ | amount | 28 | 11 | 9(11) | | | [GLJRNL.cpy:22](../../../gondwana/copybooks/GLJRNL.cpy#L22) |
| 05 | ├ JR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ | text | 39 | 10 | X(10) | | | [GLJRNL.cpy:23](../../../gondwana/copybooks/GLJRNL.cpy#L23) |
| 05 | ├ JR-ﾃｷﾖｳ | text | 49 | 20 | X(20) | | | [GLJRNL.cpy:24](../../../gondwana/copybooks/GLJRNL.cpy#L24) |
| 05 | └ FILLER | filler | 69 | 12 | X(12) | | | [GLJRNL.cpy:25](../../../gondwana/copybooks/GLJRNL.cpy#L25) |
