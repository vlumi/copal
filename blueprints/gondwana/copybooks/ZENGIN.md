# ZENGIN

```text
ZENGIN: 口座振替ファイル レコード (120バイト、固定長)

収納代行に送る振替依頼 PRMREQ.YYYYMM.DAT と、戻ってくる
振替結果 PRMRES.YYYYMMDD.DAT に共通の、銀行間で交換する
レイアウト: ヘッダ、データ、件数と金額のトレーラ、
エンドの4種で、先頭1バイトのデータ区分で区別する。
レコード区切りなし。名義は半角カナ。
振替結果は依頼では空白、結果で埋まる。

変更履歴
日付       担当 変更
2003-10-06 MS   GL-0310 初版。
```

Copybook [ZENGIN.cpy](../../../gondwana/copybooks/ZENGIN.cpy), used by [GLD050](../programs/GLD050.md), [GLM010](../programs/GLM010.md).

## ZG-ﾌﾘｶｴ-ﾚｺｰﾄﾞ

120 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | ZG-ﾌﾘｶｴ-ﾚｺｰﾄﾞ | group | 1 | 120 | | | | [ZENGIN.cpy:15](../../../gondwana/copybooks/ZENGIN.cpy#L15) |
| 05 | ├ ZG-ﾃﾞｰﾀｸﾌﾞﾝ | code | 1 | 1 | 9 | | | [ZENGIN.cpy:16](../../../gondwana/copybooks/ZENGIN.cpy#L16) |
| 88 | │ ├ ZG-ﾍｯﾀﾞ | | | | | | VALUE 1 | [ZENGIN.cpy:17](../../../gondwana/copybooks/ZENGIN.cpy#L17) |
| 88 | │ ├ ZG-ﾃﾞｰﾀ | | | | | | VALUE 2 | [ZENGIN.cpy:18](../../../gondwana/copybooks/ZENGIN.cpy#L18) |
| 88 | │ ├ ZG-ﾄﾚｰﾗ | | | | | | VALUE 8 | [ZENGIN.cpy:19](../../../gondwana/copybooks/ZENGIN.cpy#L19) |
| 88 | │ └ ZG-ｴﾝﾄﾞ | | | | | | VALUE 9 | [ZENGIN.cpy:20](../../../gondwana/copybooks/ZENGIN.cpy#L20) |
| 05 | ├ ZG-ﾍｯﾀﾞﾌﾞ | group | 2 | 119 | | | redefined by ZG-ﾃﾞｰﾀﾌﾞ, ZG-ﾄﾚｰﾗﾌﾞ | [ZENGIN.cpy:21](../../../gondwana/copybooks/ZENGIN.cpy#L21) |
| 10 | │ ├ ZG-H-ｼｭﾍﾞﾂｺｰﾄﾞ | number | 2 | 2 | 99 | | | [ZENGIN.cpy:22](../../../gondwana/copybooks/ZENGIN.cpy#L22) |
| 10 | │ ├ ZG-H-ｺｰﾄﾞｸﾌﾞﾝ | number | 4 | 1 | 9 | | | [ZENGIN.cpy:23](../../../gondwana/copybooks/ZENGIN.cpy#L23) |
| 10 | │ ├ ZG-H-ｲﾀｸｼｬｺｰﾄﾞ | number | 5 | 10 | 9(10) | | | [ZENGIN.cpy:24](../../../gondwana/copybooks/ZENGIN.cpy#L24) |
| 10 | │ ├ ZG-H-ｲﾀｸｼｬﾒｲ | text | 15 | 40 | X(40) | | | [ZENGIN.cpy:25](../../../gondwana/copybooks/ZENGIN.cpy#L25) |
| 10 | │ ├ ZG-H-ﾋｷｵﾄｼﾋﾞ | date | 55 | 4 | 9(4) | | | [ZENGIN.cpy:26](../../../gondwana/copybooks/ZENGIN.cpy#L26) |
| 10 | │ ├ ZG-H-ｷﾞﾝｺｳｺｰﾄﾞ | number | 59 | 4 | 9(4) | | | [ZENGIN.cpy:27](../../../gondwana/copybooks/ZENGIN.cpy#L27) |
| 10 | │ ├ ZG-H-ｷﾞﾝｺｳﾒｲ | text | 63 | 15 | X(15) | | | [ZENGIN.cpy:28](../../../gondwana/copybooks/ZENGIN.cpy#L28) |
| 10 | │ ├ ZG-H-ｼﾃﾝｺｰﾄﾞ | number | 78 | 3 | 9(3) | | | [ZENGIN.cpy:29](../../../gondwana/copybooks/ZENGIN.cpy#L29) |
| 10 | │ ├ ZG-H-ｼﾃﾝﾒｲ | text | 81 | 15 | X(15) | | | [ZENGIN.cpy:30](../../../gondwana/copybooks/ZENGIN.cpy#L30) |
| 10 | │ ├ ZG-H-ﾖｷﾝｼｭﾓｸ | number | 96 | 1 | 9 | | | [ZENGIN.cpy:31](../../../gondwana/copybooks/ZENGIN.cpy#L31) |
| 10 | │ ├ ZG-H-ｺｳｻﾞﾊﾞﾝｺﾞｳ | number | 97 | 7 | 9(7) | | | [ZENGIN.cpy:32](../../../gondwana/copybooks/ZENGIN.cpy#L32) |
| 10 | │ └ FILLER | filler | 104 | 17 | X(17) | | | [ZENGIN.cpy:33](../../../gondwana/copybooks/ZENGIN.cpy#L33) |
| 05 | ╞ ZG-ﾃﾞｰﾀﾌﾞ | group | 2 | 119 | | | REDEFINES ZG-ﾍｯﾀﾞﾌﾞ | [ZENGIN.cpy:34](../../../gondwana/copybooks/ZENGIN.cpy#L34) |
| 10 | │ ├ ZG-D-ｷﾞﾝｺｳｺｰﾄﾞ | number | 2 | 4 | 9(4) | | | [ZENGIN.cpy:35](../../../gondwana/copybooks/ZENGIN.cpy#L35) |
| 10 | │ ├ ZG-D-ｷﾞﾝｺｳﾒｲ | text | 6 | 15 | X(15) | | | [ZENGIN.cpy:36](../../../gondwana/copybooks/ZENGIN.cpy#L36) |
| 10 | │ ├ ZG-D-ｼﾃﾝｺｰﾄﾞ | number | 21 | 3 | 9(3) | | | [ZENGIN.cpy:37](../../../gondwana/copybooks/ZENGIN.cpy#L37) |
| 10 | │ ├ ZG-D-ｼﾃﾝﾒｲ | text | 24 | 15 | X(15) | | | [ZENGIN.cpy:38](../../../gondwana/copybooks/ZENGIN.cpy#L38) |
| 10 | │ ├ FILLER | filler | 39 | 4 | X(4) | | | [ZENGIN.cpy:39](../../../gondwana/copybooks/ZENGIN.cpy#L39) |
| 10 | │ ├ ZG-D-ﾖｷﾝｼｭﾓｸ | number | 43 | 1 | 9 | | | [ZENGIN.cpy:40](../../../gondwana/copybooks/ZENGIN.cpy#L40) |
| 10 | │ ├ ZG-D-ｺｳｻﾞﾊﾞﾝｺﾞｳ | number | 44 | 7 | 9(7) | | | [ZENGIN.cpy:41](../../../gondwana/copybooks/ZENGIN.cpy#L41) |
| 10 | │ ├ ZG-D-ｺｳｻﾞﾒｲｷﾞ | text | 51 | 30 | X(30) | | | [ZENGIN.cpy:42](../../../gondwana/copybooks/ZENGIN.cpy#L42) |
| 10 | │ ├ ZG-D-ｷﾝｶﾞｸ | amount | 81 | 10 | 9(10) | | | [ZENGIN.cpy:43](../../../gondwana/copybooks/ZENGIN.cpy#L43) |
| 10 | │ ├ ZG-D-ｼﾝｷｺｰﾄﾞ | number | 91 | 1 | 9 | | | [ZENGIN.cpy:44](../../../gondwana/copybooks/ZENGIN.cpy#L44) |
| 10 | │ ├ ZG-D-ｺｷｬｸﾊﾞﾝｺﾞｳ | text | 92 | 20 | X(20) | | | [ZENGIN.cpy:45](../../../gondwana/copybooks/ZENGIN.cpy#L45) |
| 10 | │ ├ ZG-D-ﾌﾘｶｴｹｯｶ | code | 112 | 1 | X | | | [ZENGIN.cpy:46](../../../gondwana/copybooks/ZENGIN.cpy#L46) |
| 88 | │ │ ├ ZG-D-ﾌﾘｶｴｽﾞﾐ | | | | | | VALUE '0' | [ZENGIN.cpy:47](../../../gondwana/copybooks/ZENGIN.cpy#L47) |
| 88 | │ │ ├ ZG-D-ｼｷﾝﾌｿｸ | | | | | | VALUE '1' | [ZENGIN.cpy:48](../../../gondwana/copybooks/ZENGIN.cpy#L48) |
| 88 | │ │ ├ ZG-D-ﾄﾘﾋｷﾅｼ | | | | | | VALUE '2' | [ZENGIN.cpy:49](../../../gondwana/copybooks/ZENGIN.cpy#L49) |
| 88 | │ │ ├ ZG-D-ﾖｷﾝｼｬﾂｺﾞｳ | | | | | | VALUE '3' | [ZENGIN.cpy:50](../../../gondwana/copybooks/ZENGIN.cpy#L50) |
| 88 | │ │ ├ ZG-D-ｲﾗｲｼｮﾅｼ | | | | | | VALUE '4' | [ZENGIN.cpy:51](../../../gondwana/copybooks/ZENGIN.cpy#L51) |
| 88 | │ │ └ ZG-D-ﾐｼｮﾘ | | | | | | VALUE ' ' | [ZENGIN.cpy:52](../../../gondwana/copybooks/ZENGIN.cpy#L52) |
| 10 | │ └ FILLER | filler | 113 | 8 | X(8) | | | [ZENGIN.cpy:53](../../../gondwana/copybooks/ZENGIN.cpy#L53) |
| 05 | ╘ ZG-ﾄﾚｰﾗﾌﾞ | group | 2 | 119 | | | REDEFINES ZG-ﾍｯﾀﾞﾌﾞ | [ZENGIN.cpy:54](../../../gondwana/copybooks/ZENGIN.cpy#L54) |
| 10 |   ├ ZG-T-ｺﾞｳｹｲｹﾝｽｳ | number | 2 | 6 | 9(6) | | | [ZENGIN.cpy:55](../../../gondwana/copybooks/ZENGIN.cpy#L55) |
| 10 |   ├ ZG-T-ｺﾞｳｹｲｷﾝｶﾞｸ | amount | 8 | 12 | 9(12) | | | [ZENGIN.cpy:56](../../../gondwana/copybooks/ZENGIN.cpy#L56) |
| 10 |   ├ ZG-T-ﾌﾘｶｴｽﾞﾐｹﾝｽｳ | number | 20 | 6 | 9(6) | | | [ZENGIN.cpy:57](../../../gondwana/copybooks/ZENGIN.cpy#L57) |
| 10 |   ├ ZG-T-ﾌﾘｶｴｽﾞﾐｷﾝｶﾞｸ | amount | 26 | 12 | 9(12) | | | [ZENGIN.cpy:58](../../../gondwana/copybooks/ZENGIN.cpy#L58) |
| 10 |   ├ ZG-T-ﾌﾘｶｴﾌﾉｳｹﾝｽｳ | number | 38 | 6 | 9(6) | | | [ZENGIN.cpy:59](../../../gondwana/copybooks/ZENGIN.cpy#L59) |
| 10 |   ├ ZG-T-ﾌﾘｶｴﾌﾉｳｷﾝｶﾞｸ | amount | 44 | 12 | 9(12) | | | [ZENGIN.cpy:60](../../../gondwana/copybooks/ZENGIN.cpy#L60) |
| 10 |   └ FILLER | filler | 56 | 65 | X(65) | | | [ZENGIN.cpy:61](../../../gondwana/copybooks/ZENGIN.cpy#L61) |
