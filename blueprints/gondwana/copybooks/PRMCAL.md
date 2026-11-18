# PRMCAL

```text
PRMCAL: 収納履歴の書き込み GLC050 を呼ぶ作業領域

手続きは PRMLOG。呼ぶ側は種別、金額、回数を入れて
ｼｭｳﾉｳ-ｷﾛｸ-RTN を PERFORM する。

変更履歴
日付       担当 変更
2026-11-19 YN   GL-0827 初版。
```

Copybook [PRMCAL.cpy](../../../gondwana/copybooks/PRMCAL.cpy), used by [GLD010](../programs/GLD010.md), [GLD040](../programs/GLD040.md), [GLD050](../programs/GLD050.md), [GLD060](../programs/GLD060.md).

## Contents

- [WS-ｼｭｳﾉｳ-ｷﾉｳ](#ws-ｼｭｳﾉｳ-ｷﾉｳ)
- [WS-ｼｭｳﾉｳ](#ws-ｼｭｳﾉｳ)
- [WS-PRMHST-ST](#ws-prmhst-st)

## WS-ｼｭｳﾉｳ-ｷﾉｳ

1 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-ｼｭｳﾉｳ-ｷﾉｳ | text | 1 | 1 | X | | | [PRMCAL.cpy:11](../../../gondwana/copybooks/PRMCAL.cpy#L11) |

## WS-ｼｭｳﾉｳ

36 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-ｼｭｳﾉｳ | group | 1 | 36 | | | | [PRMCAL.cpy:12](../../../gondwana/copybooks/PRMCAL.cpy#L12) |
| 05 | ├ WS-ｼｭｳﾉｳ-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ | text | 1 | 10 | X(10) | | | [PRMCAL.cpy:13](../../../gondwana/copybooks/PRMCAL.cpy#L13) |
| 05 | ├ WS-ｼｭｳﾉｳ-ﾋﾂﾞｹ | date | 11 | 8 | 9(8) | | | [PRMCAL.cpy:14](../../../gondwana/copybooks/PRMCAL.cpy#L14) |
| 05 | ├ WS-ｼｭｳﾉｳ-ｼｭﾍﾞﾂ | code | 19 | 1 | X | | | [PRMCAL.cpy:15](../../../gondwana/copybooks/PRMCAL.cpy#L15) |
| 88 | │ ├ WS-ｼｭｳﾉｳ-ｼｮｶｲ | | | | | | VALUE '1' | [PRMCAL.cpy:16](../../../gondwana/copybooks/PRMCAL.cpy#L16) |
| 88 | │ ├ WS-ｼｭｳﾉｳ-ﾌﾘｶｴ | | | | | | VALUE '2' | [PRMCAL.cpy:17](../../../gondwana/copybooks/PRMCAL.cpy#L17) |
| 88 | │ ├ WS-ｼｭｳﾉｳ-ﾏﾄﾞｸﾞﾁ | | | | | | VALUE '3' | [PRMCAL.cpy:18](../../../gondwana/copybooks/PRMCAL.cpy#L18) |
| 88 | │ └ WS-ｼｭｳﾉｳ-ｶｼﾂｹ | | | | | | VALUE '4' | [PRMCAL.cpy:19](../../../gondwana/copybooks/PRMCAL.cpy#L19) |
| 05 | ├ WS-ｼｭｳﾉｳ-ｷﾝｶﾞｸ | amount | 20 | 9 | 9(9) | | | [PRMCAL.cpy:20](../../../gondwana/copybooks/PRMCAL.cpy#L20) |
| 05 | ├ WS-ｼｭｳﾉｳ-ｶｲｽｳ | number | 29 | 2 | 99 | | | [PRMCAL.cpy:21](../../../gondwana/copybooks/PRMCAL.cpy#L21) |
| 05 | └ WS-ｼｭｳﾉｳ-ﾈﾝｹﾞﾂ | date | 31 | 6 | 9(6) | | | [PRMCAL.cpy:22](../../../gondwana/copybooks/PRMCAL.cpy#L22) |

## WS-PRMHST-ST

2 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-PRMHST-ST | code | 1 | 2 | XX | | | [PRMCAL.cpy:23](../../../gondwana/copybooks/PRMCAL.cpy#L23) |
| 88 | └ WS-PRMHST-OK | | | | | | VALUE '00' | [PRMCAL.cpy:24](../../../gondwana/copybooks/PRMCAL.cpy#L24) |
