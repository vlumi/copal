      *----------------------------------------------------------------
      * PRMLOG: 収納履歴に残す手続き
      *
      * PROCEDURE DIVISION に COPY する。作業領域は PRMCAL。
      * ｼｭｳﾉｳ-ｷﾛｸ-RTN は、種別、金額、回数を入れた後、契約マスタの
      * 払込済年月を進めてから PERFORM する。証券番号、業務日付、
      * 払込済年月を加えて GLC050 を呼ぶ。ｼｭｳﾉｳ-ｸﾛｰｽﾞ-RTN は終了時
      * に PERFORM する。書けなければ ｲｼﾞｮｳ-RTN へ。
      *
      * 変更履歴
      * 日付       担当 変更
      * 2026-11-19 YN   GL-0827 初版。
      *----------------------------------------------------------------
       ｼｭｳﾉｳ-ｷﾛｸ-RTN.
           MOVE "W" TO WS-ｼｭｳﾉｳ-ｷﾉｳ
           MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO WS-ｼｭｳﾉｳ-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ TO WS-ｼｭｳﾉｳ-ﾋﾂﾞｹ
           MOVE PM-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ TO WS-ｼｭｳﾉｳ-ﾈﾝｹﾞﾂ
           PERFORM ｼｭｳﾉｳ-ﾖﾋﾞﾀﾞｼ-RTN.

       ｼｭｳﾉｳ-ｸﾛｰｽﾞ-RTN.
           MOVE "C" TO WS-ｼｭｳﾉｳ-ｷﾉｳ
           PERFORM ｼｭｳﾉｳ-ﾖﾋﾞﾀﾞｼ-RTN.

       ｼｭｳﾉｳ-ﾖﾋﾞﾀﾞｼ-RTN.
           CALL "GLC050" USING WS-ｼｭｳﾉｳ-ｷﾉｳ WS-ｼｭｳﾉｳ WS-PRMHST-ST
           IF NOT WS-PRMHST-OK
               MOVE "PRMHST" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-PRMHST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.
