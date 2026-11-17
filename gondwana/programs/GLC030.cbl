      *----------------------------------------------------------------
      * GLC030: 解約返戻金
      *
      * 商品、加入年齢、経過年数、保険金額、返戻金、状態を受けて
      * 呼ばれる。最初の呼び出しで解約返戻金係数 CSVTB を商品・加入
      * 年齢区分・経過年数のテーブルに読み込み、以後は引くだけ。
      * 返戻金は保険金額 1,000 円あたりの係数を掛けた額で、円未満は
      * 切り捨てる。年齢区分は 20歳未満が 1、以後10歳ごと、60歳
      * 以上が 6。経過年数は 40 年で頭打ち。係数表にない商品と
      * 経過年数 0 の契約は返戻金なし。
      *
      * 状態: CSVTB のファイル状態。00 なら返戻金を求めた。
      *
      * 変更履歴
      * 日付       担当 変更
      * 2012-09-03 YN   GL-0561 初版、GLD030 から切り出し。
      *----------------------------------------------------------------
       IDENTIFICATION DIVISION.
       PROGRAM-ID. GLC030.
       AUTHOR. ゴンドワナ生命 システム部.
       DATE-WRITTEN. 2012-09-03.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT CSVTB-F ASSIGN TO "CSVTB"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-CSVTB-ST.

       DATA DIVISION.
       FILE SECTION.
       FD  CSVTB-F.
       COPY CSVTB.

       WORKING-STORAGE SECTION.
       01  WS-CSVTB-ST                 PIC XX.
           88  WS-CSVTB-OK                     VALUE "00".
           88  WS-CSVTB-ｵﾜﾘ                    VALUE "10".
       01  WS-ﾖﾐｺﾐ                     PIC X VALUE "N".
           88  WS-ﾖﾐｺﾐｽﾞﾐ                      VALUE "Y".
       01  WS-ｶｲﾔｸ-TBL.
           05  WS-ｶｲﾔｸ-ｼｮｳﾋﾝ OCCURS 10 TIMES
                                       INDEXED BY WS-ｼｮｳﾋﾝ-IX.
               10  WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ           PIC X(2).
               10  WS-ｶｲﾔｸ-ｸﾌﾞﾝ OCCURS 6 TIMES.
                   15  WS-ｶｲﾔｸ-ｹｲｽｳ OCCURS 40 TIMES
                                           PIC 9(4)V9(3).
       01  WS-ﾈﾝﾚｲｸﾌﾞﾝ                 PIC 9.
       01  WS-ﾈﾝｽｳ                     PIC 99.

       LINKAGE SECTION.
       01  LK-ｼｮｳﾋﾝ                    PIC X(2).
       01  LK-ｶﾆｭｳﾈﾝﾚｲ                 PIC 99.
       01  LK-ｹｲｶﾈﾝｽｳ                  PIC 99.
       01  LK-ﾎｹﾝｷﾝｶﾞｸ                 PIC 9(9).
       01  LK-ﾍﾝﾚｲｷﾝ                   PIC 9(9).
       01  LK-ｼﾞｮｳﾀｲ                   PIC XX.
           88  LK-ﾓﾄﾒﾀ                         VALUE "00".

       PROCEDURE DIVISION USING LK-ｼｮｳﾋﾝ LK-ｶﾆｭｳﾈﾝﾚｲ LK-ｹｲｶﾈﾝｽｳ
                                LK-ﾎｹﾝｷﾝｶﾞｸ LK-ﾍﾝﾚｲｷﾝ LK-ｼﾞｮｳﾀｲ.
      *----------------------------------------------------------------
      * ﾒｲﾝ-RTN: 係数表がまだなら読み込み、返戻金を求める。
      *----------------------------------------------------------------
       ﾒｲﾝ-RTN.
           SET LK-ﾓﾄﾒﾀ TO TRUE
           MOVE 0 TO LK-ﾍﾝﾚｲｷﾝ
           IF NOT WS-ﾖﾐｺﾐｽﾞﾐ
               PERFORM ﾖﾐｺﾐ-RTN
           END-IF
           IF WS-ﾖﾐｺﾐｽﾞﾐ
               PERFORM ｹｲｻﾝ-RTN
           END-IF
           GOBACK.

      *----------------------------------------------------------------
      * ﾖﾐｺﾐ-RTN: 係数表を読み込む。読めなければ状態を返す。
      *----------------------------------------------------------------
       ﾖﾐｺﾐ-RTN.
           MOVE ZEROS TO WS-ｶｲﾔｸ-TBL
           OPEN INPUT CSVTB-F
           IF WS-CSVTB-OK
               READ CSVTB-F
               PERFORM UNTIL NOT WS-CSVTB-OK
                   PERFORM ﾄｳﾛｸ-RTN
                   READ CSVTB-F
               END-PERFORM
               IF WS-CSVTB-ｵﾜﾘ
                   SET WS-ﾖﾐｺﾐｽﾞﾐ TO TRUE
               ELSE
                   MOVE WS-CSVTB-ST TO LK-ｼﾞｮｳﾀｲ
               END-IF
               CLOSE CSVTB-F
           ELSE
               MOVE WS-CSVTB-ST TO LK-ｼﾞｮｳﾀｲ
           END-IF.

      *----------------------------------------------------------------
      * ﾄｳﾛｸ-RTN: 係数1件。その商品の行、なければ空いた行に置く。
      * 11番目の商品は置かない。
      *----------------------------------------------------------------
       ﾄｳﾛｸ-RTN.
           SET WS-ｼｮｳﾋﾝ-IX TO 1
           SEARCH WS-ｶｲﾔｸ-ｼｮｳﾋﾝ
               AT END
                   CONTINUE
               WHEN WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = CV-ｼｮｳﾋﾝ
                       OR WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = ZEROS
                   MOVE CV-ｼｮｳﾋﾝ TO WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX)
                   MOVE CV-ｹｲｽｳ TO WS-ｶｲﾔｸ-ｹｲｽｳ(WS-ｼｮｳﾋﾝ-IX,
                                                CV-ﾈﾝﾚｲｸﾌﾞﾝ,
                                                CV-ｹｲｶﾈﾝｽｳ)
           END-SEARCH.

      *----------------------------------------------------------------
      * ｹｲｻﾝ-RTN: 年齢区分と経過年数で係数を引き、保険金額に掛ける。
      *----------------------------------------------------------------
       ｹｲｻﾝ-RTN.
           EVALUATE TRUE
               WHEN LK-ｶﾆｭｳﾈﾝﾚｲ < 20
                   MOVE 1 TO WS-ﾈﾝﾚｲｸﾌﾞﾝ
               WHEN LK-ｶﾆｭｳﾈﾝﾚｲ < 60
                   COMPUTE WS-ﾈﾝﾚｲｸﾌﾞﾝ = LK-ｶﾆｭｳﾈﾝﾚｲ / 10
               WHEN OTHER
                   MOVE 6 TO WS-ﾈﾝﾚｲｸﾌﾞﾝ
           END-EVALUATE
           IF LK-ｹｲｶﾈﾝｽｳ > 40
               MOVE 40 TO WS-ﾈﾝｽｳ
           ELSE
               MOVE LK-ｹｲｶﾈﾝｽｳ TO WS-ﾈﾝｽｳ
           END-IF
           IF WS-ﾈﾝｽｳ > 0
               SET WS-ｼｮｳﾋﾝ-IX TO 1
               SEARCH WS-ｶｲﾔｸ-ｼｮｳﾋﾝ
                   AT END
                       CONTINUE
                   WHEN WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ
                       COMPUTE LK-ﾍﾝﾚｲｷﾝ = LK-ﾎｹﾝｷﾝｶﾞｸ / 1000
                           * WS-ｶｲﾔｸ-ｹｲｽｳ(WS-ｼｮｳﾋﾝ-IX,
                                          WS-ﾈﾝﾚｲｸﾌﾞﾝ, WS-ﾈﾝｽｳ)
               END-SEARCH
           END-IF.
