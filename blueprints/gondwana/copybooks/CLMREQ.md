# CLMREQ

```text
CLMREQ: 給付請求 レコード (80バイト)

支社が入力した保険金・給付金・解約の請求1件が
CLMREQ.YYYYMMDD.TXT の1行。支社は請求書を転記し、
請求書は和暦で記入されるため、事由発生日は GYYMMDD。
事由コードは 1 疾病、2 災害、3 自殺。古い請求は空白。
手術は日数欄に給付倍率 (10, 20, 40) を入れる。

変更履歴
日付       担当 変更
1987-04-01 KT   GL-0001 初版。
2008-06-02 YN   GL-0455 入院給付金、日数付き。
2026-11-16 YN   GL-0821 事由コード (予備から)。
2026-11-16 YN   GL-0822 手術給付金。
```

Copybook [CLMREQ.cpy](../../../gondwana/copybooks/CLMREQ.cpy), used by [GLD030](../programs/GLD030.md).

## CR-ｾｲｷｭｳ-ﾚｺｰﾄﾞ

80 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | CR-ｾｲｷｭｳ-ﾚｺｰﾄﾞ | group | 1 | 80 | | | | [CLMREQ.cpy:17](../../../gondwana/copybooks/CLMREQ.cpy#L17) |
| 05 | ├ CR-ｾｲｷｭｳﾊﾞﾝｺﾞｳ | text | 1 | 8 | X(8) | | | [CLMREQ.cpy:18](../../../gondwana/copybooks/CLMREQ.cpy#L18) |
| 05 | ├ CR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ | text | 9 | 10 | X(10) | | | [CLMREQ.cpy:19](../../../gondwana/copybooks/CLMREQ.cpy#L19) |
| 05 | ├ CR-ｼｭﾍﾞﾂ | code | 19 | 1 | X | | | [CLMREQ.cpy:20](../../../gondwana/copybooks/CLMREQ.cpy#L20) |
| 88 | │ ├ CR-ｼﾎﾞｳ | | | | | | VALUE 'D' | [CLMREQ.cpy:21](../../../gondwana/copybooks/CLMREQ.cpy#L21) |
| 88 | │ ├ CR-ﾆｭｳｲﾝ | | | | | | VALUE 'H' | [CLMREQ.cpy:22](../../../gondwana/copybooks/CLMREQ.cpy#L22) |
| 88 | │ ├ CR-ｶｲﾔｸ | | | | | | VALUE 'S' | [CLMREQ.cpy:23](../../../gondwana/copybooks/CLMREQ.cpy#L23) |
| 88 | │ └ CR-ｼｭｼﾞｭﾂ | | | | | | VALUE 'O' | [CLMREQ.cpy:24](../../../gondwana/copybooks/CLMREQ.cpy#L24) |
| 05 | ├ CR-ｼﾞﾕｳﾊｯｾｲﾋﾞ | date | 20 | 7 | X(7) | | | [CLMREQ.cpy:25](../../../gondwana/copybooks/CLMREQ.cpy#L25) |
| 05 | ├ CR-ﾆｯｽｳ | number | 27 | 3 | 9(3) | | | [CLMREQ.cpy:26](../../../gondwana/copybooks/CLMREQ.cpy#L26) |
| 05 | ├ CR-ｳｹﾄﾘﾆﾝｶﾅ | text | 30 | 30 | X(30) | | | [CLMREQ.cpy:27](../../../gondwana/copybooks/CLMREQ.cpy#L27) |
| 05 | ├ CR-ｼﾃﾝ | text | 60 | 4 | X(4) | | | [CLMREQ.cpy:28](../../../gondwana/copybooks/CLMREQ.cpy#L28) |
| 05 | ├ CR-ｼﾞﾕｳｺｰﾄﾞ | code | 64 | 1 | X | | | [CLMREQ.cpy:29](../../../gondwana/copybooks/CLMREQ.cpy#L29) |
| 88 | │ └ CR-ｼﾞｻﾂ | | | | | | VALUE '3' | [CLMREQ.cpy:30](../../../gondwana/copybooks/CLMREQ.cpy#L30) |
| 05 | └ FILLER | filler | 65 | 16 | X(16) | | | [CLMREQ.cpy:31](../../../gondwana/copybooks/CLMREQ.cpy#L31) |
