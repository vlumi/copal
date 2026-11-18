      *----------------------------------------------------------------
      * GLC050: 収納履歴の書き込み
      *
      * 保険料を収納したプログラムから呼ばれる。機能 W は収納1件を
      * 収納履歴 PRMHST に、その契約の次の連番で書く。ファイルは
      * 最初の W で開き (なければ作る)、開いたままにする。機能 C は
      * 開いていれば閉じる: 呼ぶ側は終了時に必ず C で呼ぶ。
      * 状態は "00" で正常、それ以外はファイル状態をそのまま返す。
      *
      * 呼び出し: CALL "GLC050" USING 機能 収納 状態
      *   機能 PIC X      W 書き込み、C クローズ
      *   収納            証券番号 X(10)、収納日 9(8) 西暦、種別 X、
      *                   金額 9(9)、回数 99、払込済年月 9(6)
      *   状態 PIC XX
      *
      * 変更履歴
      * 日付       担当 変更
      * 2026-11-19 YN   GL-0827 初版。
      *----------------------------------------------------------------
       IDENTIFICATION DIVISION.
       PROGRAM-ID. GLC050.
       AUTHOR. ゴンドワナ生命 システム部.
       DATE-WRITTEN. 2026-11-19.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT OPTIONAL PRMHST-F ASSIGN TO "PRMHST"
               ORGANIZATION IS INDEXED
               ACCESS MODE IS DYNAMIC
               RECORD KEY IS SH-ｷｰ
               FILE STATUS IS WS-PRMHST-ST.

       DATA DIVISION.
       FILE SECTION.
       FD  PRMHST-F.
       COPY PRMHST.

       WORKING-STORAGE SECTION.
       01  WS-PRMHST-ST                PIC XX.
           88  WS-PRMHST-OK                    VALUE "00" "05".
           88  WS-PRMHST-ｵﾜﾘ                   VALUE "10" "23".
       01  WS-ｵｰﾌﾟﾝ                    PIC X VALUE "N".
           88  WS-ｵｰﾌﾟﾝｽﾞﾐ                     VALUE "Y".
           88  WS-ｵｰﾌﾟﾝﾏｴ                      VALUE "N".
       01  WS-ﾚﾝﾊﾞﾝ                    PIC 9(4).

       LINKAGE SECTION.
       01  LK-ｷﾉｳ                      PIC X.
           88  LK-ｶｷｺﾐ                         VALUE "W".
           88  LK-ｸﾛｰｽﾞ                        VALUE "C".
       01  LK-ｼｭｳﾉｳ.
           05  LK-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ          PIC X(10).
           05  LK-ｼｭｳﾉｳﾋﾞ              PIC 9(8).
           05  LK-ｼｭﾍﾞﾂ                PIC X.
           05  LK-ｷﾝｶﾞｸ                PIC 9(9).
           05  LK-ｶｲｽｳ                 PIC 99.
           05  LK-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ        PIC 9(6).
       01  LK-ｼﾞｮｳﾀｲ                   PIC XX.
           88  LK-ｾｲｼﾞｮｳ                       VALUE "00".

       PROCEDURE DIVISION USING LK-ｷﾉｳ LK-ｼｭｳﾉｳ LK-ｼﾞｮｳﾀｲ.
      *----------------------------------------------------------------
      * ﾒｲﾝ-RTN: 機能により書き込みかクローズ。
      *----------------------------------------------------------------
       ﾒｲﾝ-RTN.
           SET LK-ｾｲｼﾞｮｳ TO TRUE
           EVALUATE TRUE
               WHEN LK-ｶｷｺﾐ
                   PERFORM ｶｷｺﾐ-RTN
               WHEN LK-ｸﾛｰｽﾞ AND WS-ｵｰﾌﾟﾝｽﾞﾐ
                   CLOSE PRMHST-F
                   SET WS-ｵｰﾌﾟﾝﾏｴ TO TRUE
           END-EVALUATE
           GOBACK.

      *----------------------------------------------------------------
      * ｶｷｺﾐ-RTN: 開いていなければ開き、契約の最後の連番を求めて
      * 次の連番で書く。
      *----------------------------------------------------------------
       ｶｷｺﾐ-RTN.
           IF WS-ｵｰﾌﾟﾝﾏｴ
               OPEN I-O PRMHST-F
               IF WS-PRMHST-OK
                   SET WS-ｵｰﾌﾟﾝｽﾞﾐ TO TRUE
               END-IF
           END-IF
           IF WS-ｵｰﾌﾟﾝｽﾞﾐ
               PERFORM ﾚﾝﾊﾞﾝ-RTN
           END-IF
           IF WS-ｵｰﾌﾟﾝｽﾞﾐ AND (WS-PRMHST-OK OR WS-PRMHST-ｵﾜﾘ)
               MOVE SPACES TO SH-ｼｭｳﾉｳ-ﾚｺｰﾄﾞ
               MOVE LK-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO SH-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
               COMPUTE SH-ﾚﾝﾊﾞﾝ = WS-ﾚﾝﾊﾞﾝ + 1
               MOVE LK-ｼｭｳﾉｳﾋﾞ TO SH-ｼｭｳﾉｳﾋﾞ
               MOVE LK-ｼｭﾍﾞﾂ TO SH-ｼｭﾍﾞﾂ
               MOVE LK-ｷﾝｶﾞｸ TO SH-ｷﾝｶﾞｸ
               MOVE LK-ｶｲｽｳ TO SH-ｶｲｽｳ
               MOVE LK-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ TO SH-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ
               WRITE SH-ｼｭｳﾉｳ-ﾚｺｰﾄﾞ
           END-IF
           IF NOT WS-PRMHST-OK
               MOVE WS-PRMHST-ST TO LK-ｼﾞｮｳﾀｲ
           END-IF.

      *----------------------------------------------------------------
      * ﾚﾝﾊﾞﾝ-RTN: 契約の収納履歴を連番順にたどり、最後の連番を
      * 求める。なければ 0。
      *----------------------------------------------------------------
       ﾚﾝﾊﾞﾝ-RTN.
           MOVE 0 TO WS-ﾚﾝﾊﾞﾝ
           MOVE LK-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO SH-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           MOVE 0 TO SH-ﾚﾝﾊﾞﾝ
           START PRMHST-F KEY NOT < SH-ｷｰ
           IF WS-PRMHST-OK
               READ PRMHST-F NEXT
           END-IF
           PERFORM UNTIL NOT WS-PRMHST-OK
                   OR SH-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ NOT = LK-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
               MOVE SH-ﾚﾝﾊﾞﾝ TO WS-ﾚﾝﾊﾞﾝ
               READ PRMHST-F NEXT
           END-PERFORM.
