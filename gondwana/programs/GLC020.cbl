      *----------------------------------------------------------------
      * GLC020: 和暦・西暦 日付変換
      *
      * 方向 (W 和暦へ、S 西暦へ)、西暦日付 YYYYMMDD、和暦日付
      * GYYMMDD、結果を受けて呼ばれる。
      * 和暦へ: その日以前に始まった最後の元号を取り、元年を 1 と
      * して年を数える。
      * 西暦へ: 元号の開始年に年を足して 1 を引き、その元号の
      * 範囲内にあることを確かめる。昭和64年は 1989-01-07 まで、
      * 平成31年は 2019-04-30 まで。元号は GENGO にある。
      *
      * 結果: 0 変換済、4 いずれの元号の日付でもない。
      *
      * 変更履歴
      * 日付       担当 変更
      * 1987-04-01 KT   GL-0001 初版。
      * 1989-01-16 KT   GL-0044 元号をコードでなくテーブルから。
      * 2019-04-22 YN   GL-0702 元号の末日も検査する。
      *----------------------------------------------------------------
       IDENTIFICATION DIVISION.
       PROGRAM-ID. GLC020.
       AUTHOR. ゴンドワナ生命 システム部.
       DATE-WRITTEN. 1987-04-01.

       DATA DIVISION.
       WORKING-STORAGE SECTION.
       COPY GENGO.
       01  WS-ﾋﾂﾞｹ.
           05  WS-YYYYMMDD             PIC 9(8).
           05  WS-ﾋﾂﾞｹ-ﾌﾞﾝｶｲ REDEFINES WS-YYYYMMDD.
               10  WS-YYYY                 PIC 9(4).
               10  WS-MM                   PIC 99.
               10  WS-DD                   PIC 99.
       01  WS-ﾊｯｹﾝ                     PIC X.
           88  WS-ｹﾞﾝｺﾞｳ-ｱﾘ                    VALUE "Y".
           88  WS-ｹﾞﾝｺﾞｳ-ﾅｼ                    VALUE "N".

       LINKAGE SECTION.
       01  LK-ﾎｳｺｳ                     PIC X.
           88  LK-ﾜﾚｷ-ﾍ                        VALUE "W".
           88  LK-ｾｲﾚｷ-ﾍ                       VALUE "S".
       01  LK-ｾｲﾚｷ                     PIC 9(8).
       01  LK-ﾜﾚｷ.
           05  LK-ｹﾞﾝｺﾞｳ               PIC X.
           05  LK-ﾈﾝ                   PIC 99.
           05  LK-MMDD                 PIC 9(4).
       01  LK-ｹｯｶ                      PIC 9.
           88  LK-ﾍﾝｶﾝｽﾞﾐ                      VALUE 0.
           88  LK-ﾋﾂﾞｹ-ﾌｾｲ                     VALUE 4.

       PROCEDURE DIVISION USING LK-ﾎｳｺｳ LK-ｾｲﾚｷ LK-ﾜﾚｷ LK-ｹｯｶ.
      *----------------------------------------------------------------
      * ﾒｲﾝ-RTN: どちらかの方向に1回変換する。
      *----------------------------------------------------------------
       ﾒｲﾝ-RTN.
           SET LK-ﾍﾝｶﾝｽﾞﾐ TO TRUE
           EVALUATE TRUE
               WHEN LK-ﾜﾚｷ-ﾍ
                   PERFORM ﾜﾚｷ-ﾍﾝｶﾝ-RTN
               WHEN LK-ｾｲﾚｷ-ﾍ
                   PERFORM ｾｲﾚｷ-ﾍﾝｶﾝ-RTN
               WHEN OTHER
                   SET LK-ﾋﾂﾞｹ-ﾌｾｲ TO TRUE
           END-EVALUATE
           GOBACK.

      *----------------------------------------------------------------
      * ﾜﾚｷ-ﾍﾝｶﾝ-RTN: その日までに始まっていた最後の元号。
      *----------------------------------------------------------------
       ﾜﾚｷ-ﾍﾝｶﾝ-RTN.
           MOVE LK-ｾｲﾚｷ TO WS-YYYYMMDD
           SET WS-ｹﾞﾝｺﾞｳ-ﾅｼ TO TRUE
           PERFORM VARYING ｹﾞﾝｺﾞｳ-IX FROM ｹﾞﾝｺﾞｳ-ｹﾝｽｳ BY -1
                   UNTIL ｹﾞﾝｺﾞｳ-IX < 1 OR WS-ｹﾞﾝｺﾞｳ-ｱﾘ
               IF ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX) NOT > WS-YYYYMMDD
                   SET WS-ｹﾞﾝｺﾞｳ-ｱﾘ TO TRUE
                   MOVE ｹﾞﾝｺﾞｳ-ｷｺﾞｳ(ｹﾞﾝｺﾞｳ-IX) TO LK-ｹﾞﾝｺﾞｳ
                   COMPUTE LK-ﾈﾝ = WS-YYYY
                                 - ｹﾞﾝｺﾞｳ-ｶｲｼ-YYYY(ｹﾞﾝｺﾞｳ-IX) + 1
                   MOVE WS-ﾋﾂﾞｹ-ﾌﾞﾝｶｲ(5:4) TO LK-MMDD
               END-IF
           END-PERFORM
           IF WS-ｹﾞﾝｺﾞｳ-ﾅｼ
               SET LK-ﾋﾂﾞｹ-ﾌｾｲ TO TRUE
           END-IF.

      *----------------------------------------------------------------
      * ｾｲﾚｷ-ﾍﾝｶﾝ-RTN: 記号で元号を引き、その範囲内の日付にする。
      *----------------------------------------------------------------
       ｾｲﾚｷ-ﾍﾝｶﾝ-RTN.
           SET WS-ｹﾞﾝｺﾞｳ-ﾅｼ TO TRUE
           SET ｹﾞﾝｺﾞｳ-IX TO 1
           SEARCH ｹﾞﾝｺﾞｳ-ｴﾝﾄﾘ
               AT END
                   SET LK-ﾋﾂﾞｹ-ﾌｾｲ TO TRUE
               WHEN ｹﾞﾝｺﾞｳ-ｷｺﾞｳ(ｹﾞﾝｺﾞｳ-IX) = LK-ｹﾞﾝｺﾞｳ
                   SET WS-ｹﾞﾝｺﾞｳ-ｱﾘ TO TRUE
           END-SEARCH
           IF WS-ｹﾞﾝｺﾞｳ-ｱﾘ
               PERFORM ﾊﾝｲ-ｹﾝｻ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ﾊﾝｲ-ｹﾝｻ-RTN: 西暦日付を組み立て、元号の範囲内で、月日が
      * 月日であることを確かめる。
      *----------------------------------------------------------------
       ﾊﾝｲ-ｹﾝｻ-RTN.
           COMPUTE WS-YYYY = ｹﾞﾝｺﾞｳ-ｶｲｼ-YYYY(ｹﾞﾝｺﾞｳ-IX) + LK-ﾈﾝ - 1
           MOVE LK-MMDD TO WS-ﾋﾂﾞｹ-ﾌﾞﾝｶｲ(5:4)
           EVALUATE TRUE
               WHEN LK-ﾈﾝ = 0
                   SET LK-ﾋﾂﾞｹ-ﾌｾｲ TO TRUE
               WHEN WS-MM < 1 OR WS-MM > 12
                   SET LK-ﾋﾂﾞｹ-ﾌｾｲ TO TRUE
               WHEN WS-DD < 1 OR WS-DD > 31
                   SET LK-ﾋﾂﾞｹ-ﾌｾｲ TO TRUE
               WHEN WS-YYYYMMDD < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX)
                   SET LK-ﾋﾂﾞｹ-ﾌｾｲ TO TRUE
               WHEN ｹﾞﾝｺﾞｳ-IX < ｹﾞﾝｺﾞｳ-ｹﾝｽｳ
                   SET ｹﾞﾝｺﾞｳ-IX UP BY 1
                   IF WS-YYYYMMDD NOT < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX)
                       SET LK-ﾋﾂﾞｹ-ﾌｾｲ TO TRUE
                   END-IF
           END-EVALUATE
           IF LK-ﾍﾝｶﾝｽﾞﾐ
               MOVE WS-YYYYMMDD TO LK-ｾｲﾚｷ
           END-IF.
