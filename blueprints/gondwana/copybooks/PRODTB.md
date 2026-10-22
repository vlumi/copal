# PRODTB

```text
PRODTB: 商品テーブル レコード (60バイト)

PRODTB.TXT の1行が1商品: 商品の種類、加入できる年齢、
保険金額の範囲、支払内容。入院日額率は保険金額1,000円
あたりの入院日額、免責日数は契約日からその日数の間は
入院給付金を支払わない日数。

変更履歴
日付       担当 変更
1987-04-01 KT   GL-0001 初版。
2008-06-02 YN   GL-0455 医療保険: 入院日額率、免責日数。
```

Copybook [PRODTB.cpy](../../../gondwana/copybooks/PRODTB.cpy), used by [GLD010](../programs/GLD010.md), [GLD020](../programs/GLD020.md), [GLD030](../programs/GLD030.md), [GLM020](../programs/GLM020.md), [GLY010](../programs/GLY010.md).

## PR-ｼｮｳﾋﾝ-ﾚｺｰﾄﾞ

60 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | PR-ｼｮｳﾋﾝ-ﾚｺｰﾄﾞ | group | 1 | 60 | | | | [PRODTB.cpy:14](../../../gondwana/copybooks/PRODTB.cpy#L14) |
| 05 | ├ PR-ｺｰﾄﾞ | text | 1 | 2 | X(2) | | | [PRODTB.cpy:15](../../../gondwana/copybooks/PRODTB.cpy#L15) |
| 05 | ├ PR-ﾒｲｼｮｳ | text | 3 | 20 | X(20) | | | [PRODTB.cpy:16](../../../gondwana/copybooks/PRODTB.cpy#L16) |
| 05 | ├ PR-ｼｭﾙｲ | code | 23 | 1 | 9 | | | [PRODTB.cpy:17](../../../gondwana/copybooks/PRODTB.cpy#L17) |
| 88 | │ ├ PR-ｼｭｳｼﾝ | | | | | | VALUE 1 | [PRODTB.cpy:18](../../../gondwana/copybooks/PRODTB.cpy#L18) |
| 88 | │ ├ PR-ｺｳｼﾝﾃｲｷ | | | | | | VALUE 2 | [PRODTB.cpy:19](../../../gondwana/copybooks/PRODTB.cpy#L19) |
| 88 | │ ├ PR-ﾖｳﾛｳ | | | | | | VALUE 3 | [PRODTB.cpy:20](../../../gondwana/copybooks/PRODTB.cpy#L20) |
| 88 | │ └ PR-ｲﾘｮｳ | | | | | | VALUE 4 | [PRODTB.cpy:21](../../../gondwana/copybooks/PRODTB.cpy#L21) |
| 05 | ├ PR-ｻｲﾃｲﾈﾝﾚｲ | number | 24 | 2 | 99 | | | [PRODTB.cpy:22](../../../gondwana/copybooks/PRODTB.cpy#L22) |
| 05 | ├ PR-ｻｲｺｳﾈﾝﾚｲ | number | 26 | 2 | 99 | | | [PRODTB.cpy:23](../../../gondwana/copybooks/PRODTB.cpy#L23) |
| 05 | ├ PR-ﾎｹﾝｷｶﾝ | number | 28 | 2 | 99 | | | [PRODTB.cpy:24](../../../gondwana/copybooks/PRODTB.cpy#L24) |
| 05 | ├ PR-ｻｲﾃｲﾎｹﾝｷﾝｶﾞｸ | amount | 30 | 9 | 9(9) | | | [PRODTB.cpy:25](../../../gondwana/copybooks/PRODTB.cpy#L25) |
| 05 | ├ PR-ｻｲｺｳﾎｹﾝｷﾝｶﾞｸ | amount | 39 | 9 | 9(9) | | | [PRODTB.cpy:26](../../../gondwana/copybooks/PRODTB.cpy#L26) |
| 05 | ├ PR-ﾆｭｳｲﾝﾆﾁｶﾞｸﾘﾂ | rate | 48 | 4 | 9V999 | | | [PRODTB.cpy:27](../../../gondwana/copybooks/PRODTB.cpy#L27) |
| 05 | ├ PR-ﾒﾝｾｷﾆｯｽｳ | number | 52 | 3 | 9(3) | | | [PRODTB.cpy:28](../../../gondwana/copybooks/PRODTB.cpy#L28) |
| 05 | ├ PR-ﾊｲﾄｳﾘﾂ | rate | 55 | 5 | 9V9(4) | | | [PRODTB.cpy:29](../../../gondwana/copybooks/PRODTB.cpy#L29) |
| 05 | └ FILLER | filler | 60 | 1 | X(1) | | | [PRODTB.cpy:30](../../../gondwana/copybooks/PRODTB.cpy#L30) |
