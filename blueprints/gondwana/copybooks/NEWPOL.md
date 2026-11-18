# NEWPOL

```text
NEWPOL: 新契約 レコード (220バイト)

営業システムで成立した申込1件が NEWPOL.YYYYMMDD.TXT の
1行。営業システムは本システムより新しく、日付を西暦
YYYYMMDD で書く。マスタには和暦で持つ。
特約は2件まで、商品と保険金額。なければ商品は空白。

変更履歴
日付       担当 変更
2003-10-06 MS   GL-0310 初版、営業システム連携。
2026-11-18 YN   GL-0825 特約 (予備を含めて20バイト延長)。
```

Copybook [NEWPOL.cpy](../../../gondwana/copybooks/NEWPOL.cpy), used by [GLD010](../programs/GLD010.md).

## NP-ｼﾝｹｲﾔｸ-ﾚｺｰﾄﾞ

220 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | NP-ｼﾝｹｲﾔｸ-ﾚｺｰﾄﾞ | group | 1 | 220 | | | | [NEWPOL.cpy:14](../../../gondwana/copybooks/NEWPOL.cpy#L14) |
| 05 | ├ NP-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ | text | 1 | 10 | X(10) | | | [NEWPOL.cpy:15](../../../gondwana/copybooks/NEWPOL.cpy#L15) |
| 05 | ├ NP-ｼｮｳﾋﾝ | text | 11 | 2 | X(2) | | | [NEWPOL.cpy:16](../../../gondwana/copybooks/NEWPOL.cpy#L16) |
| 05 | ├ NP-ﾎｹﾝｷﾝｶﾞｸ | amount | 13 | 9 | 9(9) | | | [NEWPOL.cpy:17](../../../gondwana/copybooks/NEWPOL.cpy#L17) |
| 05 | ├ NP-ﾊﾗｲｺﾐﾎｳﾎｳ | text | 22 | 1 | X | | | [NEWPOL.cpy:18](../../../gondwana/copybooks/NEWPOL.cpy#L18) |
| 05 | ├ NP-ｹｲﾔｸﾋﾞ | date | 23 | 8 | 9(8) | | | [NEWPOL.cpy:19](../../../gondwana/copybooks/NEWPOL.cpy#L19) |
| 05 | ├ NP-ｼﾒｲ | text | 31 | 20 | X(20) | | | [NEWPOL.cpy:20](../../../gondwana/copybooks/NEWPOL.cpy#L20) |
| 05 | ├ NP-ｼﾒｲｶﾅ | text | 51 | 20 | X(20) | | | [NEWPOL.cpy:21](../../../gondwana/copybooks/NEWPOL.cpy#L21) |
| 05 | ├ NP-ｾｲﾈﾝｶﾞｯﾋﾟ | date | 71 | 8 | 9(8) | | | [NEWPOL.cpy:22](../../../gondwana/copybooks/NEWPOL.cpy#L22) |
| 05 | ├ NP-ｾｲﾍﾞﾂ | number | 79 | 1 | 9 | | | [NEWPOL.cpy:23](../../../gondwana/copybooks/NEWPOL.cpy#L23) |
| 05 | ├ NP-ｼﾞｭｳｼｮ | text | 80 | 40 | X(40) | | | [NEWPOL.cpy:24](../../../gondwana/copybooks/NEWPOL.cpy#L24) |
| 05 | ├ NP-ｳｹﾄﾘﾆﾝ | text | 120 | 20 | X(20) | | | [NEWPOL.cpy:25](../../../gondwana/copybooks/NEWPOL.cpy#L25) |
| 05 | ├ NP-ｷﾞﾝｺｳｺｰﾄﾞ | number | 140 | 4 | 9(4) | | | [NEWPOL.cpy:26](../../../gondwana/copybooks/NEWPOL.cpy#L26) |
| 05 | ├ NP-ｼﾃﾝｺｰﾄﾞ | number | 144 | 3 | 9(3) | | | [NEWPOL.cpy:27](../../../gondwana/copybooks/NEWPOL.cpy#L27) |
| 05 | ├ NP-ﾖｷﾝｼｭﾓｸ | number | 147 | 1 | 9 | | | [NEWPOL.cpy:28](../../../gondwana/copybooks/NEWPOL.cpy#L28) |
| 05 | ├ NP-ｺｳｻﾞﾊﾞﾝｺﾞｳ | number | 148 | 7 | 9(7) | | | [NEWPOL.cpy:29](../../../gondwana/copybooks/NEWPOL.cpy#L29) |
| 05 | ├ NP-ｺｳｻﾞﾒｲｷﾞ | text | 155 | 30 | X(30) | | | [NEWPOL.cpy:30](../../../gondwana/copybooks/NEWPOL.cpy#L30) |
| 05 | ├ NP-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ | text | 185 | 6 | X(6) | | | [NEWPOL.cpy:31](../../../gondwana/copybooks/NEWPOL.cpy#L31) |
| 05 | ├ NP-ﾄｸﾔｸ | group | 191 | 11 | | | OCCURS 2 | [NEWPOL.cpy:32](../../../gondwana/copybooks/NEWPOL.cpy#L32) |
| 10 | │ ├ NP-ﾄｸﾔｸ-ｼｮｳﾋﾝ | text | 191 | 2 | X(2) | | | [NEWPOL.cpy:33](../../../gondwana/copybooks/NEWPOL.cpy#L33) |
| 10 | │ └ NP-ﾄｸﾔｸ-ｷﾝｶﾞｸ | amount | 193 | 9 | 9(9) | | | [NEWPOL.cpy:34](../../../gondwana/copybooks/NEWPOL.cpy#L34) |
| 05 | └ FILLER | filler | 213 | 8 | X(8) | | | [NEWPOL.cpy:35](../../../gondwana/copybooks/NEWPOL.cpy#L35) |
