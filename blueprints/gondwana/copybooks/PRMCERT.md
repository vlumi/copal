# PRMCERT

```text
PRMCERT: 保険料控除証明 レコード (140バイト)

証明書1通が PRMCERT.YYYYMMDD.TXT の1行で、印刷会社に渡す。
年分は西暦。証明額はその年の1月1日から発行日までに収納した
保険料、予定額は12月まで払込が続いた場合の年間の額。
発行日と契約日は和暦 GYYMMDD。

変更履歴
日付       担当 変更
2026-11-19 YN   GL-0827 初版。
```

Copybook [PRMCERT.cpy](../../../gondwana/copybooks/PRMCERT.cpy), used by [GLY030](../programs/GLY030.md).

## SM-ｼｮｳﾒｲ-ﾚｺｰﾄﾞ

140 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | SM-ｼｮｳﾒｲ-ﾚｺｰﾄﾞ | group | 1 | 140 | | | | [PRMCERT.cpy:13](../../../gondwana/copybooks/PRMCERT.cpy#L13) |
| 05 | ├ SM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ | text | 1 | 10 | X(10) | | | [PRMCERT.cpy:14](../../../gondwana/copybooks/PRMCERT.cpy#L14) |
| 05 | ├ SM-ﾈﾝﾌﾞﾝ | number | 11 | 4 | 9(4) | | | [PRMCERT.cpy:15](../../../gondwana/copybooks/PRMCERT.cpy#L15) |
| 05 | ├ SM-ﾊｯｺｳﾋﾞ | date | 15 | 7 | X(7) | | | [PRMCERT.cpy:16](../../../gondwana/copybooks/PRMCERT.cpy#L16) |
| 05 | ├ SM-ｼﾒｲ | text | 22 | 20 | X(20) | | | [PRMCERT.cpy:17](../../../gondwana/copybooks/PRMCERT.cpy#L17) |
| 05 | ├ SM-ｼﾞｭｳｼｮ | text | 42 | 40 | X(40) | | | [PRMCERT.cpy:18](../../../gondwana/copybooks/PRMCERT.cpy#L18) |
| 05 | ├ SM-ｼｮｳﾋﾝ | text | 82 | 2 | X(2) | | | [PRMCERT.cpy:19](../../../gondwana/copybooks/PRMCERT.cpy#L19) |
| 05 | ├ SM-ｼｮｳﾋﾝﾒｲ | text | 84 | 20 | X(20) | | | [PRMCERT.cpy:20](../../../gondwana/copybooks/PRMCERT.cpy#L20) |
| 05 | ├ SM-ｹｲﾔｸﾋﾞ | date | 104 | 7 | X(7) | | | [PRMCERT.cpy:21](../../../gondwana/copybooks/PRMCERT.cpy#L21) |
| 05 | ├ SM-ﾊﾗｲｺﾐﾎｳﾎｳ | text | 111 | 1 | X | | | [PRMCERT.cpy:22](../../../gondwana/copybooks/PRMCERT.cpy#L22) |
| 05 | ├ SM-ｼｮｳﾒｲｶﾞｸ | amount | 112 | 9 | 9(9) | | | [PRMCERT.cpy:23](../../../gondwana/copybooks/PRMCERT.cpy#L23) |
| 05 | ├ SM-ﾖﾃｲｶﾞｸ | amount | 121 | 9 | 9(9) | | | [PRMCERT.cpy:24](../../../gondwana/copybooks/PRMCERT.cpy#L24) |
| 05 | └ FILLER | filler | 130 | 11 | X(11) | | | [PRMCERT.cpy:25](../../../gondwana/copybooks/PRMCERT.cpy#L25) |
