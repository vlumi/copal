# TBLWS

```text
TBLWS: 商品テーブルと料率テーブルの作業領域

PRODTB と RATETB を読み込む先。商品の並びが料率の第1添字で、
年齢 0～80 は添字 1～81。読み込みは PRDLOD と RATLOD の手続き。

変更履歴
日付       担当 変更
1987-04-01 KT   GL-0001 初版。
```

Copybook [TBLWS.cpy](../../../gondwana/copybooks/TBLWS.cpy), used by [GLD010](../programs/GLD010.md), [GLD020](../programs/GLD020.md), [GLD030](../programs/GLD030.md), [GLM020](../programs/GLM020.md), [GLY010](../programs/GLY010.md).

## Contents

- [WS-ｼｮｳﾋﾝ-TBL](#ws-ｼｮｳﾋﾝ-tbl)
- [WS-ﾘｮｳﾘﾂ-TBL](#ws-ﾘｮｳﾘﾂ-tbl)
- [WS-ﾈﾝﾚｲ-ｿｴｼﾞ](#ws-ﾈﾝﾚｲ-ｿｴｼﾞ)

## WS-ｼｮｳﾋﾝ-TBL

592 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-ｼｮｳﾋﾝ-TBL | group | 1 | 592 | | | | [TBLWS.cpy:11](../../../gondwana/copybooks/TBLWS.cpy#L11) |
| 05 | ├ WS-ｼｮｳﾋﾝｹﾝｽｳ | number | 1 | 2 | 99 | | VALUE 0 | [TBLWS.cpy:12](../../../gondwana/copybooks/TBLWS.cpy#L12) |
| 05 | └ WS-ｼｮｳﾋﾝ | group | 3 | 59 | | | OCCURS 10; INDEXED BY WS-ｼｮｳﾋﾝ-IX | [TBLWS.cpy:13](../../../gondwana/copybooks/TBLWS.cpy#L13) |
| 10 |   ├ WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ | text | 3 | 2 | X(2) | | | [TBLWS.cpy:14](../../../gondwana/copybooks/TBLWS.cpy#L14) |
| 10 |   ├ WS-ｼｮｳﾋﾝ-ﾒｲｼｮｳ | text | 5 | 20 | X(20) | | | [TBLWS.cpy:15](../../../gondwana/copybooks/TBLWS.cpy#L15) |
| 10 |   ├ WS-ｼｮｳﾋﾝ-ｼｭﾙｲ | code | 25 | 1 | 9 | | | [TBLWS.cpy:16](../../../gondwana/copybooks/TBLWS.cpy#L16) |
| 88 |   │ ├ WS-ｼｮｳﾋﾝ-ｼｭｳｼﾝ | | | | | | VALUE 1 | [TBLWS.cpy:17](../../../gondwana/copybooks/TBLWS.cpy#L17) |
| 88 |   │ ├ WS-ｼｮｳﾋﾝ-ｺｳｼﾝﾃｲｷ | | | | | | VALUE 2 | [TBLWS.cpy:18](../../../gondwana/copybooks/TBLWS.cpy#L18) |
| 88 |   │ ├ WS-ｼｮｳﾋﾝ-ﾖｳﾛｳ | | | | | | VALUE 3 | [TBLWS.cpy:19](../../../gondwana/copybooks/TBLWS.cpy#L19) |
| 88 |   │ └ WS-ｼｮｳﾋﾝ-ｲﾘｮｳ | | | | | | VALUE 4 | [TBLWS.cpy:20](../../../gondwana/copybooks/TBLWS.cpy#L20) |
| 10 |   ├ WS-ｼｮｳﾋﾝ-ｻｲﾃｲﾈﾝﾚｲ | number | 26 | 2 | 99 | | | [TBLWS.cpy:21](../../../gondwana/copybooks/TBLWS.cpy#L21) |
| 10 |   ├ WS-ｼｮｳﾋﾝ-ｻｲｺｳﾈﾝﾚｲ | number | 28 | 2 | 99 | | | [TBLWS.cpy:22](../../../gondwana/copybooks/TBLWS.cpy#L22) |
| 10 |   ├ WS-ｼｮｳﾋﾝ-ﾎｹﾝｷｶﾝ | number | 30 | 2 | 99 | | | [TBLWS.cpy:23](../../../gondwana/copybooks/TBLWS.cpy#L23) |
| 10 |   ├ WS-ｼｮｳﾋﾝ-ｻｲﾃｲｷﾝｶﾞｸ | number | 32 | 9 | 9(9) | | | [TBLWS.cpy:24](../../../gondwana/copybooks/TBLWS.cpy#L24) |
| 10 |   ├ WS-ｼｮｳﾋﾝ-ｻｲｺｳｷﾝｶﾞｸ | number | 41 | 9 | 9(9) | | | [TBLWS.cpy:25](../../../gondwana/copybooks/TBLWS.cpy#L25) |
| 10 |   ├ WS-ｼｮｳﾋﾝ-ﾆｭｳｲﾝﾆﾁｶﾞｸﾘﾂ | amount | 50 | 4 | 9V999 | | | [TBLWS.cpy:26](../../../gondwana/copybooks/TBLWS.cpy#L26) |
| 10 |   ├ WS-ｼｮｳﾋﾝ-ﾒﾝｾｷﾆｯｽｳ | number | 54 | 3 | 9(3) | | | [TBLWS.cpy:27](../../../gondwana/copybooks/TBLWS.cpy#L27) |
| 10 |   └ WS-ｼｮｳﾋﾝ-ﾊｲﾄｳﾘﾂ | amount | 57 | 5 | 9V9(4) | | | [TBLWS.cpy:28](../../../gondwana/copybooks/TBLWS.cpy#L28) |

## WS-ﾘｮｳﾘﾂ-TBL

8100 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-ﾘｮｳﾘﾂ-TBL | group | 1 | 8100 | | | | [TBLWS.cpy:29](../../../gondwana/copybooks/TBLWS.cpy#L29) |
| 05 | └ WS-ﾘｮｳﾘﾂ-ｼｮｳﾋﾝ | group | 1 | 810 | | | OCCURS 10 | [TBLWS.cpy:30](../../../gondwana/copybooks/TBLWS.cpy#L30) |
| 10 |   └ WS-ﾘｮｳﾘﾂ-ｾｲﾍﾞﾂ | group | 1 | 405 | | | OCCURS 2 | [TBLWS.cpy:31](../../../gondwana/copybooks/TBLWS.cpy#L31) |
| 15 |     └ WS-ﾘｮｳﾘﾂ | amount | 1 | 5 | 9(3)V99 | | OCCURS 81 | [TBLWS.cpy:32](../../../gondwana/copybooks/TBLWS.cpy#L32) |

## WS-ﾈﾝﾚｲ-ｿｴｼﾞ

2 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-ﾈﾝﾚｲ-ｿｴｼﾞ | number | 1 | 2 | 99 | | | [TBLWS.cpy:34](../../../gondwana/copybooks/TBLWS.cpy#L34) |
