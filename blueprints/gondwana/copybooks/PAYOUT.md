# PAYOUT

```text
PAYOUT: 支払 レコード (100バイト)

システムが指図する支払1件が PAYOUT.YYYYMMDD.TXT の1行、
支払担当向け: 死亡保険金、入院給付金、解約返戻金、
満期保険金。支払番号は制御レコードから採番する。

変更履歴
日付       担当 変更
1987-04-01 KT   GL-0001 初版。
```

Copybook [PAYOUT.cpy](../../../gondwana/copybooks/PAYOUT.cpy), used by [GLD020](../programs/GLD020.md), [GLD030](../programs/GLD030.md).

## PY-ｼﾊﾗｲ-ﾚｺｰﾄﾞ

100 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | PY-ｼﾊﾗｲ-ﾚｺｰﾄﾞ | group | 1 | 100 | | | | [PAYOUT.cpy:12](../../../gondwana/copybooks/PAYOUT.cpy#L12) |
| 05 | ├ PY-ｼﾊﾗｲﾊﾞﾝｺﾞｳ | number | 1 | 9 | 9(9) | | | [PAYOUT.cpy:13](../../../gondwana/copybooks/PAYOUT.cpy#L13) |
| 05 | ├ PY-ｼﾊﾗｲﾋﾞ | date | 10 | 7 | X(7) | | | [PAYOUT.cpy:14](../../../gondwana/copybooks/PAYOUT.cpy#L14) |
| 05 | ├ PY-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ | text | 17 | 10 | X(10) | | | [PAYOUT.cpy:15](../../../gondwana/copybooks/PAYOUT.cpy#L15) |
| 05 | ├ PY-ｼｭﾍﾞﾂ | code | 27 | 1 | X | | | [PAYOUT.cpy:16](../../../gondwana/copybooks/PAYOUT.cpy#L16) |
| 88 | │ ├ PY-ｼﾎﾞｳ | | | | | | VALUE 'D' | [PAYOUT.cpy:17](../../../gondwana/copybooks/PAYOUT.cpy#L17) |
| 88 | │ ├ PY-ﾆｭｳｲﾝ | | | | | | VALUE 'H' | [PAYOUT.cpy:18](../../../gondwana/copybooks/PAYOUT.cpy#L18) |
| 88 | │ ├ PY-ｶｲﾔｸ | | | | | | VALUE 'S' | [PAYOUT.cpy:19](../../../gondwana/copybooks/PAYOUT.cpy#L19) |
| 88 | │ └ PY-ﾏﾝｷ | | | | | | VALUE 'M' | [PAYOUT.cpy:20](../../../gondwana/copybooks/PAYOUT.cpy#L20) |
| 05 | ├ PY-ｷﾝｶﾞｸ | amount | 28 | 9 | 9(9) | | | [PAYOUT.cpy:21](../../../gondwana/copybooks/PAYOUT.cpy#L21) |
| 05 | ├ PY-ｳｹﾄﾘﾆﾝｶﾅ | text | 37 | 30 | X(30) | | | [PAYOUT.cpy:22](../../../gondwana/copybooks/PAYOUT.cpy#L22) |
| 05 | ├ PY-ｾｲｷｭｳﾊﾞﾝｺﾞｳ | text | 67 | 8 | X(8) | | | [PAYOUT.cpy:23](../../../gondwana/copybooks/PAYOUT.cpy#L23) |
| 05 | └ FILLER | filler | 75 | 26 | X(26) | | | [PAYOUT.cpy:24](../../../gondwana/copybooks/PAYOUT.cpy#L24) |
