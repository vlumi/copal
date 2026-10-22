      *----------------------------------------------------------------
      * RATLOD: 料率テーブルの読み込み手続き
      *
      * PROCEDURE DIVISION に COPY する。RATETB-F を開き、PRDLOD で
      * 商品テーブルを読み込んだ後に ﾘｮｳﾘﾂ-ﾖﾐｺﾐ-RTN を PERFORM
      * する。読めなければ ｲｼﾞｮｳ-RTN へ。読み込む先は TBLWS。
      *
      * 変更履歴
      * 日付       担当 変更
      * 1987-04-01 KT   GL-0001 初版。
      *----------------------------------------------------------------
      *----------------------------------------------------------------
      * ﾘｮｳﾘﾂ-ﾖﾐｺﾐ-RTN: 料率を商品・性別・年齢のテーブルに読み込む。
      * 商品テーブルにない商品の料率は読み飛ばす。
      *----------------------------------------------------------------
       ﾘｮｳﾘﾂ-ﾖﾐｺﾐ-RTN.
           MOVE ZEROS TO WS-ﾘｮｳﾘﾂ-TBL
           READ RATETB-F
           PERFORM UNTIL NOT WS-RATETB-OK
               SET WS-ｼｮｳﾋﾝ-IX TO 1
               SEARCH WS-ｼｮｳﾋﾝ
                   WHEN WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = RT-ｼｮｳﾋﾝ
                       COMPUTE WS-ﾈﾝﾚｲ-ｿｴｼﾞ = RT-ﾈﾝﾚｲ + 1
                       MOVE RT-ﾘｮｳﾘﾂ TO WS-ﾘｮｳﾘﾂ(WS-ｼｮｳﾋﾝ-IX,
                                                 RT-ｾｲﾍﾞﾂ,
                                                 WS-ﾈﾝﾚｲ-ｿｴｼﾞ)
               END-SEARCH
               READ RATETB-F
           END-PERFORM
           IF NOT WS-RATETB-ｵﾜﾘ
               MOVE "RATETB" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-RATETB-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.
