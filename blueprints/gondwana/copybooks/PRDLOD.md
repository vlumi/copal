# PRDLOD

```text
PRDLOD: 商品テーブルの読み込み手続き

PROCEDURE DIVISION に COPY する。PRODTB-F を開いた後に
ｼｮｳﾋﾝ-ﾖﾐｺﾐ-RTN を PERFORM する。読めなければ ｲｼﾞｮｳ-RTN へ。
読み込む先は TBLWS。

変更履歴
日付       担当 変更
1987-04-01 KT   GL-0001 初版。
2008-06-02 YN   GL-0455 医療保険の項目。
```

Copybook [PRDLOD.cpy](../../../gondwana/copybooks/PRDLOD.cpy), used by [GLD010](../programs/GLD010.md), [GLD020](../programs/GLD020.md), [GLD030](../programs/GLD030.md), [GLM020](../programs/GLM020.md), [GLY010](../programs/GLY010.md).

A copybook of procedure: its paragraphs, ｼｮｳﾋﾝ-ﾖﾐｺﾐ-RTN, are copied into the procedure division of each program that uses it, and are in each program's blueprint.
