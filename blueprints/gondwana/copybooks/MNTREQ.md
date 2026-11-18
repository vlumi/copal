# MNTREQ

```text
MNTREQ: 契約変更依頼 レコード (180バイト)

支社が入力した変更1件が MNTREQ.YYYYMMDD.TXT の1行。
種別により後続のどの項目が埋まるかが決まる: 住所、
口座、氏名とそのカナ、受取人。復活はいずれも空で、
窓口で収納した未納保険料の額を持つ。特約解約は解約する
特約の保障番号を持つ。

変更履歴
日付       担当 変更
1987-04-01 KT   GL-0001 初版。
2003-10-06 MS   GL-0310 口座。
2011-02-14 YN   GL-0520 窓口からの復活。
2026-11-18 YN   GL-0825 特約解約、保障番号 (予備から)。
```

Copybook [MNTREQ.cpy](../../../gondwana/copybooks/MNTREQ.cpy), used by [GLD040](../programs/GLD040.md).

## MR-ﾍﾝｺｳ-ﾚｺｰﾄﾞ

180 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | MR-ﾍﾝｺｳ-ﾚｺｰﾄﾞ | group | 1 | 180 | | | | [MNTREQ.cpy:17](../../../gondwana/copybooks/MNTREQ.cpy#L17) |
| 05 | ├ MR-ｲﾗｲﾊﾞﾝｺﾞｳ | text | 1 | 8 | X(8) | | | [MNTREQ.cpy:18](../../../gondwana/copybooks/MNTREQ.cpy#L18) |
| 05 | ├ MR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ | text | 9 | 10 | X(10) | | | [MNTREQ.cpy:19](../../../gondwana/copybooks/MNTREQ.cpy#L19) |
| 05 | ├ MR-ｼｭﾍﾞﾂ | code | 19 | 1 | X | | | [MNTREQ.cpy:20](../../../gondwana/copybooks/MNTREQ.cpy#L20) |
| 88 | │ ├ MR-ｼﾞｭｳｼｮﾍﾝｺｳ | | | | | | VALUE 'A' | [MNTREQ.cpy:21](../../../gondwana/copybooks/MNTREQ.cpy#L21) |
| 88 | │ ├ MR-ｺｳｻﾞﾍﾝｺｳ | | | | | | VALUE 'B' | [MNTREQ.cpy:22](../../../gondwana/copybooks/MNTREQ.cpy#L22) |
| 88 | │ ├ MR-ｼﾒｲﾍﾝｺｳ | | | | | | VALUE 'N' | [MNTREQ.cpy:23](../../../gondwana/copybooks/MNTREQ.cpy#L23) |
| 88 | │ ├ MR-ｳｹﾄﾘﾆﾝﾍﾝｺｳ | | | | | | VALUE 'K' | [MNTREQ.cpy:24](../../../gondwana/copybooks/MNTREQ.cpy#L24) |
| 88 | │ ├ MR-ﾌｯｶﾂ | | | | | | VALUE 'R' | [MNTREQ.cpy:25](../../../gondwana/copybooks/MNTREQ.cpy#L25) |
| 88 | │ └ MR-ﾄｸﾔｸｶｲﾔｸ | | | | | | VALUE 'T' | [MNTREQ.cpy:26](../../../gondwana/copybooks/MNTREQ.cpy#L26) |
| 05 | ├ MR-ｼﾞｭｳｼｮ | text | 20 | 40 | X(40) | | | [MNTREQ.cpy:27](../../../gondwana/copybooks/MNTREQ.cpy#L27) |
| 05 | ├ MR-ｺｳｻﾞ | group | 60 | 45 | | | | [MNTREQ.cpy:28](../../../gondwana/copybooks/MNTREQ.cpy#L28) |
| 10 | │ ├ MR-ｷﾞﾝｺｳｺｰﾄﾞ | number | 60 | 4 | 9(4) | | | [MNTREQ.cpy:29](../../../gondwana/copybooks/MNTREQ.cpy#L29) |
| 10 | │ ├ MR-ｼﾃﾝｺｰﾄﾞ | number | 64 | 3 | 9(3) | | | [MNTREQ.cpy:30](../../../gondwana/copybooks/MNTREQ.cpy#L30) |
| 10 | │ ├ MR-ﾖｷﾝｼｭﾓｸ | number | 67 | 1 | 9 | | | [MNTREQ.cpy:31](../../../gondwana/copybooks/MNTREQ.cpy#L31) |
| 10 | │ ├ MR-ｺｳｻﾞﾊﾞﾝｺﾞｳ | number | 68 | 7 | 9(7) | | | [MNTREQ.cpy:32](../../../gondwana/copybooks/MNTREQ.cpy#L32) |
| 10 | │ └ MR-ｺｳｻﾞﾒｲｷﾞ | text | 75 | 30 | X(30) | | | [MNTREQ.cpy:33](../../../gondwana/copybooks/MNTREQ.cpy#L33) |
| 05 | ├ MR-ｼﾒｲ | text | 105 | 20 | X(20) | | | [MNTREQ.cpy:34](../../../gondwana/copybooks/MNTREQ.cpy#L34) |
| 05 | ├ MR-ｼﾒｲｶﾅ | text | 125 | 20 | X(20) | | | [MNTREQ.cpy:35](../../../gondwana/copybooks/MNTREQ.cpy#L35) |
| 05 | ├ MR-ｳｹﾄﾘﾆﾝ | text | 145 | 20 | X(20) | | | [MNTREQ.cpy:36](../../../gondwana/copybooks/MNTREQ.cpy#L36) |
| 05 | ├ MR-ﾆｭｳｷﾝｶﾞｸ | amount | 165 | 7 | 9(7) | | | [MNTREQ.cpy:37](../../../gondwana/copybooks/MNTREQ.cpy#L37) |
| 05 | ├ MR-ﾎｼｮｳﾊﾞﾝｺﾞｳ | number | 172 | 2 | 99 | | | [MNTREQ.cpy:38](../../../gondwana/copybooks/MNTREQ.cpy#L38) |
| 05 | └ FILLER | filler | 174 | 7 | X(7) | | | [MNTREQ.cpy:39](../../../gondwana/copybooks/MNTREQ.cpy#L39) |
