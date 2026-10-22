# RATLOD

```text
RATLOD: 料率テーブルの読み込み手続き

PROCEDURE DIVISION に COPY する。RATETB-F を開き、PRDLOD で
商品テーブルを読み込んだ後に ﾘｮｳﾘﾂ-ﾖﾐｺﾐ-RTN を PERFORM
する。読めなければ ｲｼﾞｮｳ-RTN へ。読み込む先は TBLWS。

変更履歴
日付       担当 変更
1987-04-01 KT   GL-0001 初版。
```

Copybook [RATLOD.cpy](../../../gondwana/copybooks/RATLOD.cpy), used by [GLD010](../programs/GLD010.md), [GLD020](../programs/GLD020.md).

A copybook of procedure: its paragraphs, ﾘｮｳﾘﾂ-ﾖﾐｺﾐ-RTN, are copied into the procedure division of each program that uses it, and are in each program's blueprint.
