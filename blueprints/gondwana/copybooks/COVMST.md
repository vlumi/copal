# COVMST

```text
COVMST: 保障マスタ レコード (80バイト)

契約の保障1件につき1レコード。キーは証券番号と保障番号で、
保障番号 01 が主契約、02 以降が特約。主契約はどの契約にも
あり、その商品・保険金額・保険期間・満期日は契約マスタと同じ。
保険金額は、入院特約と医療保険では入院日額の基になる額。
保険料は契約の払込方法での1回分で、有効な保障の保険料の合計
が契約マスタの保険料である。
状態は保障自身のもの: 契約が失効しても保障は有効のままで、
契約の状態が先に立つ。満了は保険期間の終わり、解約は特約の
解約か契約の解約、消滅は死亡保険金の支払による。
加入年齢は保険料を決めた年齢で、更新のたびにその時点の
保険年齢になる。
日付は和暦 GYYMMDD。消滅日は有効でなくなった日で、有効な
間は空白。

変更履歴
日付       担当 変更
2026-11-18 YN   GL-0825 初版、特約の導入。
```

Copybook [COVMST.cpy](../../../gondwana/copybooks/COVMST.cpy), used by [GLD010](../programs/GLD010.md), [GLD020](../programs/GLD020.md), [GLD030](../programs/GLD030.md), [GLD040](../programs/GLD040.md), [GLD090](../programs/GLD090.md), [GLM020](../programs/GLM020.md), [GLO010](../programs/GLO010.md), [GLY020](../programs/GLY020.md).

## HS-ﾎｼｮｳ-ﾚｺｰﾄﾞ

80 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | HS-ﾎｼｮｳ-ﾚｺｰﾄﾞ | group | 1 | 80 | | | | [COVMST.cpy:22](../../../gondwana/copybooks/COVMST.cpy#L22) |
| 05 | ├ HS-ｷｰ | group | 1 | 12 | | | | [COVMST.cpy:23](../../../gondwana/copybooks/COVMST.cpy#L23) |
| 10 | │ ├ HS-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ | text | 1 | 10 | X(10) | | | [COVMST.cpy:24](../../../gondwana/copybooks/COVMST.cpy#L24) |
| 10 | │ └ HS-ﾎｼｮｳﾊﾞﾝｺﾞｳ | code | 11 | 2 | 99 | | | [COVMST.cpy:25](../../../gondwana/copybooks/COVMST.cpy#L25) |
| 88 | │   └ HS-ｼｭｹｲﾔｸ | | | | | | VALUE 1 | [COVMST.cpy:26](../../../gondwana/copybooks/COVMST.cpy#L26) |
| 05 | ├ HS-ｼｮｳﾋﾝ | text | 13 | 2 | X(2) | | | [COVMST.cpy:27](../../../gondwana/copybooks/COVMST.cpy#L27) |
| 05 | ├ HS-ｼﾞｮｳﾀｲ | code | 15 | 1 | X | | | [COVMST.cpy:28](../../../gondwana/copybooks/COVMST.cpy#L28) |
| 88 | │ ├ HS-ﾕｳｺｳ | | | | | | VALUE '1' | [COVMST.cpy:29](../../../gondwana/copybooks/COVMST.cpy#L29) |
| 88 | │ ├ HS-ﾏﾝﾘｮｳ | | | | | | VALUE '2' | [COVMST.cpy:30](../../../gondwana/copybooks/COVMST.cpy#L30) |
| 88 | │ ├ HS-ｶｲﾔｸ | | | | | | VALUE '3' | [COVMST.cpy:31](../../../gondwana/copybooks/COVMST.cpy#L31) |
| 88 | │ └ HS-ｼｮｳﾒﾂ | | | | | | VALUE '4' | [COVMST.cpy:32](../../../gondwana/copybooks/COVMST.cpy#L32) |
| 05 | ├ HS-ﾎｹﾝｷﾝｶﾞｸ | amount | 16 | 9 | 9(9) | | | [COVMST.cpy:33](../../../gondwana/copybooks/COVMST.cpy#L33) |
| 05 | ├ HS-ﾎｹﾝﾘｮｳ | amount | 25 | 7 | 9(7) | | | [COVMST.cpy:34](../../../gondwana/copybooks/COVMST.cpy#L34) |
| 05 | ├ HS-ｹｲﾔｸﾋﾞ | date | 32 | 7 | X(7) | | | [COVMST.cpy:35](../../../gondwana/copybooks/COVMST.cpy#L35) |
| 05 | ├ HS-ﾎｹﾝｷｶﾝ | number | 39 | 2 | 99 | | | [COVMST.cpy:36](../../../gondwana/copybooks/COVMST.cpy#L36) |
| 05 | ├ HS-ﾏﾝｷﾋﾞ | date | 41 | 7 | X(7) | | | [COVMST.cpy:37](../../../gondwana/copybooks/COVMST.cpy#L37) |
| 05 | ├ HS-ｶﾆｭｳﾈﾝﾚｲ | number | 48 | 2 | 99 | | | [COVMST.cpy:38](../../../gondwana/copybooks/COVMST.cpy#L38) |
| 05 | ├ HS-ｼｮｳﾒﾂﾋﾞ | date | 50 | 7 | X(7) | | | [COVMST.cpy:39](../../../gondwana/copybooks/COVMST.cpy#L39) |
| 05 | └ FILLER | filler | 57 | 24 | X(24) | | | [COVMST.cpy:40](../../../gondwana/copybooks/COVMST.cpy#L40) |
