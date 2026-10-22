# GENGO

```text
GENGO: 元号テーブル、日付変換用

各元号の記号、名称、西暦での開始日を古い順に持つ。
改元はここに1行加え、これを COPY するプログラムを
再コンパイルする。

変更履歴
日付       担当 変更
1987-04-01 KT   GL-0001 初版。
1989-01-16 KT   GL-0044 平成。
2019-04-22 YN   GL-0702 令和。
```

Copybook [GENGO.cpy](../../../gondwana/copybooks/GENGO.cpy), used by [GLC020](../programs/GLC020.md).

## Contents

- [ｹﾞﾝｺﾞｳ-ﾃｲｽﾞｳ](#ｹﾞﾝｺﾞｳ-ﾃｲｽﾞｳ)
- [ｹﾞﾝｺﾞｳ-ﾃｰﾌﾞﾙ](#ｹﾞﾝｺﾞｳ-ﾃｰﾌﾞﾙ)
- [ｹﾞﾝｺﾞｳ-ｹﾝｽｳ](#ｹﾞﾝｺﾞｳ-ｹﾝｽｳ)

## ｹﾞﾝｺﾞｳ-ﾃｲｽﾞｳ

65 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | ｹﾞﾝｺﾞｳ-ﾃｲｽﾞｳ | group | 1 | 65 | | | | [GENGO.cpy:14](../../../gondwana/copybooks/GENGO.cpy#L14) |
| 05 | ├ FILLER | filler | 1 | 13 | X(13) | | VALUE 'M明治18680908' | [GENGO.cpy:15](../../../gondwana/copybooks/GENGO.cpy#L15) |
| 05 | ├ FILLER | filler | 14 | 13 | X(13) | | VALUE 'T大正19120730' | [GENGO.cpy:16](../../../gondwana/copybooks/GENGO.cpy#L16) |
| 05 | ├ FILLER | filler | 27 | 13 | X(13) | | VALUE 'S昭和19261225' | [GENGO.cpy:17](../../../gondwana/copybooks/GENGO.cpy#L17) |
| 05 | ├ FILLER | filler | 40 | 13 | X(13) | | VALUE 'H平成19890108' | [GENGO.cpy:18](../../../gondwana/copybooks/GENGO.cpy#L18) |
| 05 | └ FILLER | filler | 53 | 13 | X(13) | | VALUE 'R令和20190501' | [GENGO.cpy:19](../../../gondwana/copybooks/GENGO.cpy#L19) |

## ｹﾞﾝｺﾞｳ-ﾃｰﾌﾞﾙ

65 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | ｹﾞﾝｺﾞｳ-ﾃｰﾌﾞﾙ | group | 1 | 65 | | | REDEFINES ｹﾞﾝｺﾞｳ-ﾃｲｽﾞｳ | [GENGO.cpy:20](../../../gondwana/copybooks/GENGO.cpy#L20) |
| 05 | └ ｹﾞﾝｺﾞｳ-ｴﾝﾄﾘ | group | 1 | 13 | | | OCCURS 5; INDEXED BY ｹﾞﾝｺﾞｳ-IX | [GENGO.cpy:21](../../../gondwana/copybooks/GENGO.cpy#L21) |
| 10 |   ├ ｹﾞﾝｺﾞｳ-ｷｺﾞｳ | text | 1 | 1 | X | | | [GENGO.cpy:23](../../../gondwana/copybooks/GENGO.cpy#L23) |
| 10 |   ├ ｹﾞﾝｺﾞｳ-ﾒｲｼｮｳ | text | 2 | 4 | X(4) | | | [GENGO.cpy:24](../../../gondwana/copybooks/GENGO.cpy#L24) |
| 10 |   ├ ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ | date | 6 | 8 | 9(8) | | redefined by ｹﾞﾝｺﾞｳ-ｶｲｼﾈﾝ | [GENGO.cpy:25](../../../gondwana/copybooks/GENGO.cpy#L25) |
| 10 |   ╘ ｹﾞﾝｺﾞｳ-ｶｲｼﾈﾝ | group | 6 | 8 | | | REDEFINES ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ | [GENGO.cpy:26](../../../gondwana/copybooks/GENGO.cpy#L26) |
| 15 |     ├ ｹﾞﾝｺﾞｳ-ｶｲｼ-YYYY | number | 6 | 4 | 9(4) | | | [GENGO.cpy:27](../../../gondwana/copybooks/GENGO.cpy#L27) |
| 15 |     └ FILLER | filler | 10 | 4 | 9(4) | | | [GENGO.cpy:28](../../../gondwana/copybooks/GENGO.cpy#L28) |

## ｹﾞﾝｺﾞｳ-ｹﾝｽｳ

1 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | ｹﾞﾝｺﾞｳ-ｹﾝｽｳ | number | 1 | 1 | 9 | | VALUE 5 | [GENGO.cpy:29](../../../gondwana/copybooks/GENGO.cpy#L29) |
