# POLMST

```text
POLMST: 契約マスタ レコード (300バイト)

契約1件につき1レコード。索引マスタでは証券番号をキーとし、
開設時の抽出ファイルと日次のアンロードでは同じレコードが
1行である。
日付は和暦 GYYMMDD: 元号記号 (S 昭和, H 平成, R 令和) と
元号年。例外は払込済年月で、収納代行の振替ファイルに
合わせて西暦 YYYYMM とする。
金額は円単位、小数なし。
氏名はシフトJIS: 漢字氏名は1文字2バイト、カナ氏名は
半角カナ1文字1バイト。

変更履歴
日付       担当 変更
1987-04-01 KT   GL-0001 初版 (昭和62年)。
1989-01-16 KT   GL-0044 年の前に元号記号を置く。
2003-10-06 MS   GL-0310 口座振替用の口座。
2019-04-22 YN   GL-0702 令和: レイアウト変更なし。
```

Copybook [POLMST.cpy](../../../gondwana/copybooks/POLMST.cpy), used by [GLD010](../programs/GLD010.md), [GLD020](../programs/GLD020.md), [GLD030](../programs/GLD030.md), [GLD040](../programs/GLD040.md), [GLD050](../programs/GLD050.md), [GLD060](../programs/GLD060.md), [GLD090](../programs/GLD090.md), [GLM010](../programs/GLM010.md), [GLM020](../programs/GLM020.md), [GLM030](../programs/GLM030.md), [GLO010](../programs/GLO010.md), [GLY010](../programs/GLY010.md).

## PM-ｹｲﾔｸ-ﾚｺｰﾄﾞ

300 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | PM-ｹｲﾔｸ-ﾚｺｰﾄﾞ | group | 1 | 300 | | | | [POLMST.cpy:21](../../../gondwana/copybooks/POLMST.cpy#L21) |
| 05 | ├ PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ | key | 1 | 10 | X(10) | | | [POLMST.cpy:22](../../../gondwana/copybooks/POLMST.cpy#L22) |
| 05 | ├ PM-ｼﾞｮｳﾀｲ | code | 11 | 1 | X | | | [POLMST.cpy:23](../../../gondwana/copybooks/POLMST.cpy#L23) |
| 88 | │ ├ PM-ﾕｳｺｳ | | | | | | VALUE '1' | [POLMST.cpy:24](../../../gondwana/copybooks/POLMST.cpy#L24) |
| 88 | │ ├ PM-ﾕｳﾖﾁｭｳ | | | | | | VALUE '2' | [POLMST.cpy:25](../../../gondwana/copybooks/POLMST.cpy#L25) |
| 88 | │ ├ PM-ｼｯｺｳ | | | | | | VALUE '3' | [POLMST.cpy:26](../../../gondwana/copybooks/POLMST.cpy#L26) |
| 88 | │ ├ PM-ｶｲﾔｸ | | | | | | VALUE '4' | [POLMST.cpy:27](../../../gondwana/copybooks/POLMST.cpy#L27) |
| 88 | │ ├ PM-ﾏﾝｷ | | | | | | VALUE '5' | [POLMST.cpy:28](../../../gondwana/copybooks/POLMST.cpy#L28) |
| 88 | │ ├ PM-ｼﾊﾗｲｽﾞﾐ | | | | | | VALUE '6' | [POLMST.cpy:29](../../../gondwana/copybooks/POLMST.cpy#L29) |
| 88 | │ ├ PM-ﾊﾗｲｺﾐﾁｭｳ | | | | | | VALUE '1', '2' | [POLMST.cpy:30](../../../gondwana/copybooks/POLMST.cpy#L30) |
| 88 | │ └ PM-ﾎｼｮｳﾁｭｳ | | | | | | VALUE '1', '2' | [POLMST.cpy:31](../../../gondwana/copybooks/POLMST.cpy#L31) |
| 05 | ├ PM-ｼｮｳﾋﾝ | code | 12 | 2 | X(2) | | | [POLMST.cpy:32](../../../gondwana/copybooks/POLMST.cpy#L32) |
| 88 | │ ├ PM-ｼｭｳｼﾝ | | | | | | VALUE 'WL' | [POLMST.cpy:33](../../../gondwana/copybooks/POLMST.cpy#L33) |
| 88 | │ ├ PM-ﾃｲｷ | | | | | | VALUE 'TM' | [POLMST.cpy:34](../../../gondwana/copybooks/POLMST.cpy#L34) |
| 88 | │ ├ PM-ﾖｳﾛｳ | | | | | | VALUE 'EN' | [POLMST.cpy:35](../../../gondwana/copybooks/POLMST.cpy#L35) |
| 88 | │ └ PM-ｲﾘｮｳ | | | | | | VALUE 'MD' | [POLMST.cpy:36](../../../gondwana/copybooks/POLMST.cpy#L36) |
| 05 | ├ PM-ﾎｹﾝｷﾝｶﾞｸ | amount | 14 | 9 | 9(9) | | | [POLMST.cpy:37](../../../gondwana/copybooks/POLMST.cpy#L37) |
| 05 | ├ PM-ﾎｹﾝﾘｮｳ | amount | 23 | 7 | 9(7) | | | [POLMST.cpy:38](../../../gondwana/copybooks/POLMST.cpy#L38) |
| 05 | ├ PM-ﾊﾗｲｺﾐﾎｳﾎｳ | code | 30 | 1 | X | | | [POLMST.cpy:39](../../../gondwana/copybooks/POLMST.cpy#L39) |
| 88 | │ ├ PM-ﾂｷﾊﾞﾗｲ | | | | | | VALUE '1' | [POLMST.cpy:40](../../../gondwana/copybooks/POLMST.cpy#L40) |
| 88 | │ └ PM-ﾈﾝﾊﾞﾗｲ | | | | | | VALUE '2' | [POLMST.cpy:41](../../../gondwana/copybooks/POLMST.cpy#L41) |
| 05 | ├ PM-ｹｲﾔｸﾋﾞ | date | 31 | 7 | X(7) | | | [POLMST.cpy:42](../../../gondwana/copybooks/POLMST.cpy#L42) |
| 05 | ├ PM-ﾎｹﾝｷｶﾝ | number | 38 | 2 | 99 | | | [POLMST.cpy:43](../../../gondwana/copybooks/POLMST.cpy#L43) |
| 05 | ├ PM-ﾏﾝｷﾋﾞ | date | 40 | 7 | X(7) | | | [POLMST.cpy:44](../../../gondwana/copybooks/POLMST.cpy#L44) |
| 05 | ├ PM-ｹｲｶﾈﾝｽｳ | number | 47 | 2 | 99 | | | [POLMST.cpy:45](../../../gondwana/copybooks/POLMST.cpy#L45) |
| 05 | ├ PM-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ | date | 49 | 6 | 9(6) | | | [POLMST.cpy:46](../../../gondwana/copybooks/POLMST.cpy#L46) |
| 05 | ├ PM-ﾐﾉｳｶｲｽｳ | number | 55 | 2 | 99 | | | [POLMST.cpy:47](../../../gondwana/copybooks/POLMST.cpy#L47) |
| 05 | ├ PM-ﾕｳﾖｶｲｼﾋﾞ | date | 57 | 7 | X(7) | | | [POLMST.cpy:48](../../../gondwana/copybooks/POLMST.cpy#L48) |
| 05 | ├ PM-ｼｯｺｳﾋﾞ | date | 64 | 7 | X(7) | | | [POLMST.cpy:49](../../../gondwana/copybooks/POLMST.cpy#L49) |
| 05 | ├ PM-ｼﾞｮｳﾀｲﾋﾂﾞｹ | date | 71 | 7 | X(7) | | | [POLMST.cpy:50](../../../gondwana/copybooks/POLMST.cpy#L50) |
| 05 | ├ PM-ﾋﾎｹﾝｼｬ | group | 78 | 50 | | | | [POLMST.cpy:51](../../../gondwana/copybooks/POLMST.cpy#L51) |
| 10 | │ ├ PM-ｼﾒｲ | text | 78 | 20 | X(20) | | | [POLMST.cpy:52](../../../gondwana/copybooks/POLMST.cpy#L52) |
| 10 | │ ├ PM-ｼﾒｲｶﾅ | text | 98 | 20 | X(20) | | | [POLMST.cpy:53](../../../gondwana/copybooks/POLMST.cpy#L53) |
| 10 | │ ├ PM-ｾｲﾈﾝｶﾞｯﾋﾟ | date | 118 | 7 | X(7) | | | [POLMST.cpy:54](../../../gondwana/copybooks/POLMST.cpy#L54) |
| 10 | │ ├ PM-ｾｲﾍﾞﾂ | code | 125 | 1 | 9 | | | [POLMST.cpy:55](../../../gondwana/copybooks/POLMST.cpy#L55) |
| 88 | │ │ ├ PM-ｵﾄｺ | | | | | | VALUE 1 | [POLMST.cpy:56](../../../gondwana/copybooks/POLMST.cpy#L56) |
| 88 | │ │ └ PM-ｵﾝﾅ | | | | | | VALUE 2 | [POLMST.cpy:57](../../../gondwana/copybooks/POLMST.cpy#L57) |
| 10 | │ └ PM-ｶﾆｭｳﾈﾝﾚｲ | number | 126 | 2 | 99 | | | [POLMST.cpy:58](../../../gondwana/copybooks/POLMST.cpy#L58) |
| 05 | ├ PM-ｹｲﾔｸｼｬ | group | 128 | 60 | | | | [POLMST.cpy:59](../../../gondwana/copybooks/POLMST.cpy#L59) |
| 10 | │ ├ PM-ｼﾞｭｳｼｮ | text | 128 | 40 | X(40) | | | [POLMST.cpy:60](../../../gondwana/copybooks/POLMST.cpy#L60) |
| 10 | │ └ PM-ｳｹﾄﾘﾆﾝ | text | 168 | 20 | X(20) | | | [POLMST.cpy:61](../../../gondwana/copybooks/POLMST.cpy#L61) |
| 05 | ├ PM-ｺｳｻﾞ | group | 188 | 45 | | | | [POLMST.cpy:62](../../../gondwana/copybooks/POLMST.cpy#L62) |
| 10 | │ ├ PM-ｷﾞﾝｺｳｺｰﾄﾞ | number | 188 | 4 | 9(4) | | | [POLMST.cpy:63](../../../gondwana/copybooks/POLMST.cpy#L63) |
| 10 | │ ├ PM-ｼﾃﾝｺｰﾄﾞ | number | 192 | 3 | 9(3) | | | [POLMST.cpy:64](../../../gondwana/copybooks/POLMST.cpy#L64) |
| 10 | │ ├ PM-ﾖｷﾝｼｭﾓｸ | number | 195 | 1 | 9 | | | [POLMST.cpy:65](../../../gondwana/copybooks/POLMST.cpy#L65) |
| 10 | │ ├ PM-ｺｳｻﾞﾊﾞﾝｺﾞｳ | number | 196 | 7 | 9(7) | | | [POLMST.cpy:66](../../../gondwana/copybooks/POLMST.cpy#L66) |
| 10 | │ └ PM-ｺｳｻﾞﾒｲｷﾞ | text | 203 | 30 | X(30) | | | [POLMST.cpy:67](../../../gondwana/copybooks/POLMST.cpy#L67) |
| 05 | ├ PM-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ | text | 233 | 6 | X(6) | | | [POLMST.cpy:68](../../../gondwana/copybooks/POLMST.cpy#L68) |
| 05 | ├ PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ | amount | 239 | 9 | 9(9) | | | [POLMST.cpy:69](../../../gondwana/copybooks/POLMST.cpy#L69) |
| 05 | ├ PM-ｾｷﾆﾝｼﾞｭﾝﾋﾞｷﾝ | amount | 248 | 9 | 9(9) | | | [POLMST.cpy:70](../../../gondwana/copybooks/POLMST.cpy#L70) |
| 05 | ├ PM-ﾊｲﾄｳﾙｲｹｲ | amount | 257 | 9 | 9(9) | | | [POLMST.cpy:71](../../../gondwana/copybooks/POLMST.cpy#L71) |
| 05 | ├ PM-ｻｲｼｭｳﾆｭｳｷﾝﾈﾝｹﾞﾂ | date | 266 | 6 | 9(6) | | | [POLMST.cpy:72](../../../gondwana/copybooks/POLMST.cpy#L72) |
| 05 | ├ PM-ｻｲｼｭｳﾆｭｳｷﾝｶﾞｸ | amount | 272 | 7 | 9(7) | | | [POLMST.cpy:73](../../../gondwana/copybooks/POLMST.cpy#L73) |
| 05 | └ FILLER | filler | 279 | 22 | X(22) | | | [POLMST.cpy:74](../../../gondwana/copybooks/POLMST.cpy#L74) |
