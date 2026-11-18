      *----------------------------------------------------------------
      * GLY030: 保険料控除証明書
      *
      * 年次、10月の最初の営業日。収納履歴 PRMHST を証券番号順に
      * 読み、コマンド行の日付の年の1月1日からその日までに収納した
      * 保険料を契約ごとに合計する。収納のあった契約に証明書1通を
      * PRMCERT に書く。証明額はその合計。予定額は、払込中 (有効・
      * 猶予中) の契約では、払込済年月の後、12月までに払込期月の来る
      * 保険料を証明額に加えた額: 月払は残りの月数分、年払は次の
      * 払込期月が年内なら1回分、半年払は年内の回数分。払込中で
      * ない契約の予定額は証明額に等しい。
      * 収納履歴にあって契約マスタにない契約は数えて飛ばす。
      *
      * 復帰コード: 0 作成済、4 契約のない収納履歴あり、
      *             8 日付が日付でない、16 ファイル異常。
      *
      * 変更履歴
      * 日付       担当 変更
      * 2026-11-19 YN   GL-0827 初版。
      *----------------------------------------------------------------
       IDENTIFICATION DIVISION.
       PROGRAM-ID. GLY030.
       AUTHOR. ゴンドワナ生命 システム部.
       DATE-WRITTEN. 2026-11-19.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT OPTIONAL PRMHST-F ASSIGN TO "PRMHST"
               ORGANIZATION IS INDEXED
               ACCESS MODE IS SEQUENTIAL
               RECORD KEY IS SH-ｷｰ
               FILE STATUS IS WS-PRMHST-ST.
           SELECT POLMST-F ASSIGN TO "POLMST"
               ORGANIZATION IS INDEXED
               ACCESS MODE IS RANDOM
               RECORD KEY IS PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
               FILE STATUS IS WS-POLMST-ST.
           SELECT PRODTB-F ASSIGN TO "PRODTB"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-PRODTB-ST.
           SELECT PRMCERT-F ASSIGN TO "PRMCERT"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-PRMCERT-ST.
           SELECT REPORT-F ASSIGN TO "GLY030R"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-ST.

       DATA DIVISION.
       FILE SECTION.
       FD  PRMHST-F.
       COPY PRMHST.
       FD  POLMST-F.
       COPY POLMST.
       FD  PRODTB-F.
       COPY PRODTB.
       FD  PRMCERT-F.
       COPY PRMCERT.
       FD  REPORT-F.
       01  ﾎｳｺｸ-ｷﾞｮｳ                   PIC X(100).

       WORKING-STORAGE SECTION.
       01  WS-ｽﾃｰﾀｽ.
           05  WS-PRMHST-ST            PIC XX.
               88  WS-PRMHST-OK                VALUE "00" "05".
               88  WS-PRMHST-ｵﾜﾘ               VALUE "10".
           05  WS-POLMST-ST            PIC XX.
               88  WS-POLMST-OK                VALUE "00".
               88  WS-POLMST-ﾅｼ                VALUE "23".
           05  WS-PRODTB-ST            PIC XX.
               88  WS-PRODTB-OK                VALUE "00".
               88  WS-PRODTB-ｵﾜﾘ               VALUE "10".
           05  WS-PRMCERT-ST           PIC XX.
               88  WS-PRMCERT-OK               VALUE "00".
           05  WS-REPORT-ST            PIC XX.
               88  WS-REPORT-OK                VALUE "00".
       01  WS-ｺﾏﾝﾄﾞｷﾞｮｳ                PIC X(80).
       01  WS-ｴﾗｰﾌｧｲﾙ                  PIC X(8).
       01  WS-ｴﾗｰｽﾃｰﾀｽ                 PIC XX.
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ                PIC 9(8).
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾌﾞﾝｶｲ REDEFINES WS-ｷﾞｮｳﾑﾋﾂﾞｹ.
           05  WS-ｷﾞｮｳﾑ-YYYY           PIC 9(4).
           05  WS-ｷﾞｮｳﾑ-MMDD           PIC 9(4).
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ            PIC X(7).
       01  WS-ﾎｳｺｳ                     PIC X.
       01  WS-ﾍﾝｶﾝｹｯｶ                  PIC 9.
       01  WS-ﾈﾝｼｮ                     PIC 9(8).
       01  WS-ﾈﾝﾏﾂ                     PIC 9(6).
       01  WS-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ              PIC X(10).
       01  WS-ｼｮｳﾒｲｶﾞｸ                 PIC 9(11).
       01  WS-ﾖﾃｲｶｲｽｳ                  PIC S9(5).
       01  WS-ﾊﾗｲｺﾐｽﾞﾐ.
           05  WS-ﾊﾗｲｺﾐ-YYYY           PIC 9(4).
           05  WS-ﾊﾗｲｺﾐ-MM             PIC 99.
       01  WS-ﾊﾗｲｺﾐ-YYYYMM REDEFINES WS-ﾊﾗｲｺﾐｽﾞﾐ PIC 9(6).
       COPY TBLWS.
       01  WS-ｹﾝｽｳ.
           05  WS-ﾘﾚｷｹﾝｽｳ              PIC 9(7) VALUE 0.
           05  WS-ｼｮｳﾒｲｹﾝｽｳ            PIC 9(7) VALUE 0.
           05  WS-ｹｲﾔｸﾅｼｹﾝｽｳ           PIC 9(7) VALUE 0.
           05  WS-ｼｮｳﾒｲｺﾞｳｹｲ           PIC 9(13) VALUE 0.
           05  WS-ﾖﾃｲｺﾞｳｹｲ             PIC 9(13) VALUE 0.
       01  WS-ﾌｯｷｺｰﾄﾞ                  PIC 99 VALUE 0.
       01  WS-ﾎｳｺｸ.
           05  WS-ﾋｮｳﾀﾞｲ-ｷﾞｮｳ.
               10  FILLER                  PIC X(32)
                   VALUE "ゴンドワナ生命  保険料控除証明書".
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾋｮｳﾀﾞｲ-ﾈﾝ            PIC 9(4).
               10  FILLER                  PIC X(6) VALUE "年分".
               10  WS-ﾋｮｳﾀﾞｲ-ﾋﾂﾞｹ          PIC X(7).
           05  WS-ﾒｲｻｲ-ｷﾞｮｳ.
               10  WS-ﾒｲｻｲ-ｼｮｳｹﾝ           PIC X(10).
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾒｲｻｲ-ｼｮｳﾋﾝ           PIC X(2).
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾒｲｻｲ-ｼﾒｲ             PIC X(20).
               10  WS-ﾒｲｻｲ-ｼｮｳﾒｲｶﾞｸ        PIC ZZZ,ZZZ,ZZ9.
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾒｲｻｲ-ﾖﾃｲｶﾞｸ          PIC ZZZ,ZZZ,ZZ9.
           05  WS-ｹﾝｽｳ-ｷﾞｮｳ.
               10  WS-ｹﾝｽｳ-ﾗﾍﾞﾙ            PIC X(16).
               10  WS-ｹﾝｽｳ-ｱﾀｲ             PIC Z,ZZZ,ZZZ,ZZZ,ZZ9.

       PROCEDURE DIVISION.
      *----------------------------------------------------------------
      * ﾒｲﾝ-RTN: 準備、契約ごとの集計と証明書、報告、終了。
      *----------------------------------------------------------------
       ﾒｲﾝ-RTN.
           PERFORM ｼｮｷ-RTN
           PERFORM ｹｲﾔｸ-RTN UNTIL WS-PRMHST-ｵﾜﾘ
           PERFORM ﾎｳｺｸ-RTN
           PERFORM ｼｭｳﾘｮｳ-RTN
           MOVE WS-ﾌｯｷｺｰﾄﾞ TO RETURN-CODE
           STOP RUN.

      *----------------------------------------------------------------
      * ｼｮｷ-RTN: 業務日付、その年の初めと終わり、ファイル、商品。
      *----------------------------------------------------------------
       ｼｮｷ-RTN.
           ACCEPT WS-ｺﾏﾝﾄﾞｷﾞｮｳ FROM COMMAND-LINE
           MOVE WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8) TO WS-ｷﾞｮｳﾑﾋﾂﾞｹ
           MOVE "W" TO WS-ﾎｳｺｳ
           CALL "GLC020" USING WS-ﾎｳｺｳ WS-ｷﾞｮｳﾑﾋﾂﾞｹ
                               WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ WS-ﾍﾝｶﾝｹｯｶ
           IF WS-ﾍﾝｶﾝｹｯｶ NOT = 0
               DISPLAY "GLY030 日付不正: " WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8)
               MOVE 8 TO RETURN-CODE
               STOP RUN
           END-IF
           COMPUTE WS-ﾈﾝｼｮ = WS-ｷﾞｮｳﾑ-YYYY * 10000 + 101
           COMPUTE WS-ﾈﾝﾏﾂ = WS-ｷﾞｮｳﾑ-YYYY * 100 + 12
           OPEN INPUT PRMHST-F POLMST-F PRODTB-F
           OPEN OUTPUT PRMCERT-F REPORT-F
           EVALUATE TRUE
               WHEN NOT WS-PRMHST-OK
                   MOVE "PRMHST" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-PRMHST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-POLMST-OK
                   MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-PRODTB-OK
                   MOVE "PRODTB" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-PRODTB-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-PRMCERT-OK
                   MOVE "PRMCERT" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-PRMCERT-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-REPORT-OK
                   MOVE "GLY030R" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-REPORT-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN OTHER
                   MOVE SPACES TO WS-ｴﾗｰﾌｧｲﾙ
           END-EVALUATE
           IF WS-ｴﾗｰﾌｧｲﾙ NOT = SPACES
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           PERFORM ｼｮｳﾋﾝ-ﾖﾐｺﾐ-RTN
           MOVE WS-ｷﾞｮｳﾑ-YYYY TO WS-ﾋｮｳﾀﾞｲ-ﾈﾝ
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO WS-ﾋｮｳﾀﾞｲ-ﾋﾂﾞｹ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ﾋｮｳﾀﾞｲ-ｷﾞｮｳ
           MOVE SPACES TO ﾎｳｺｸ-ｷﾞｮｳ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ
           PERFORM PRMHST-ﾖﾐｺﾐ-RTN.

      *----------------------------------------------------------------
      * PRMHST-ﾖﾐｺﾐ-RTN: 収納履歴の次の1件。
      *----------------------------------------------------------------
       PRMHST-ﾖﾐｺﾐ-RTN.
           READ PRMHST-F NEXT
           IF WS-PRMHST-OK
               ADD 1 TO WS-ﾘﾚｷｹﾝｽｳ
           ELSE
               IF NOT WS-PRMHST-ｵﾜﾘ
                   MOVE "PRMHST" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-PRMHST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
                   PERFORM ｲｼﾞｮｳ-RTN
               END-IF
           END-IF.

      *----------------------------------------------------------------
      * ｹｲﾔｸ-RTN: 証券番号が変わるまで、その年の業務日付までの収納を
      * 合計し、あれば証明書を作る。
      *----------------------------------------------------------------
       ｹｲﾔｸ-RTN.
           MOVE SH-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO WS-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           MOVE 0 TO WS-ｼｮｳﾒｲｶﾞｸ
           PERFORM UNTIL WS-PRMHST-ｵﾜﾘ
                   OR SH-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ NOT = WS-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
               IF SH-ｼｭｳﾉｳﾋﾞ NOT < WS-ﾈﾝｼｮ
                       AND SH-ｼｭｳﾉｳﾋﾞ NOT > WS-ｷﾞｮｳﾑﾋﾂﾞｹ
                   ADD SH-ｷﾝｶﾞｸ TO WS-ｼｮｳﾒｲｶﾞｸ
               END-IF
               PERFORM PRMHST-ﾖﾐｺﾐ-RTN
           END-PERFORM
           IF WS-ｼｮｳﾒｲｶﾞｸ > 0
               PERFORM ｼｮｳﾒｲ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ｼｮｳﾒｲ-RTN: 契約を読み、予定額を求めて証明書1通を書く。
      *----------------------------------------------------------------
       ｼｮｳﾒｲ-RTN.
           MOVE WS-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           READ POLMST-F
           EVALUATE TRUE
               WHEN WS-POLMST-ﾅｼ
                   ADD 1 TO WS-ｹｲﾔｸﾅｼｹﾝｽｳ
                   MOVE 4 TO WS-ﾌｯｷｺｰﾄﾞ
               WHEN NOT WS-POLMST-OK
                   MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
                   PERFORM ｲｼﾞｮｳ-RTN
               WHEN OTHER
                   PERFORM ﾖﾃｲ-RTN
                   PERFORM ｶｷｺﾐ-RTN
           END-EVALUATE.

      *----------------------------------------------------------------
      * ﾖﾃｲ-RTN: 払込中の契約の、払込済年月の後、12月までに払込期月
      * の来る保険料の回数。月払は残りの月数、年払は次の払込期月が
      * 年内なら 1、半年払は年内に来る回数。
      *----------------------------------------------------------------
       ﾖﾃｲ-RTN.
           MOVE 0 TO WS-ﾖﾃｲｶｲｽｳ
           MOVE PM-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ TO WS-ﾊﾗｲｺﾐｽﾞﾐ
           IF PM-ﾊﾗｲｺﾐﾁｭｳ
               EVALUATE TRUE
                   WHEN PM-ﾈﾝﾊﾞﾗｲ
                       ADD 1 TO WS-ﾊﾗｲｺﾐ-YYYY
                       PERFORM ｷｹﾞﾂ-RTN
                   WHEN PM-ﾊﾝﾄｼﾊﾞﾗｲ
                       PERFORM 2 TIMES
                           ADD 6 TO WS-ﾊﾗｲｺﾐ-MM
                           PERFORM ｷｹﾞﾂ-RTN
                       END-PERFORM
                   WHEN OTHER
                       COMPUTE WS-ﾖﾃｲｶｲｽｳ
                           = (WS-ｷﾞｮｳﾑ-YYYY - WS-ﾊﾗｲｺﾐ-YYYY) * 12
                           + 12 - WS-ﾊﾗｲｺﾐ-MM
                       IF WS-ﾖﾃｲｶｲｽｳ < 0
                           MOVE 0 TO WS-ﾖﾃｲｶｲｽｳ
                       END-IF
               END-EVALUATE
           END-IF.

      *----------------------------------------------------------------
      * ｷｹﾞﾂ-RTN: 進めた払込期月が年内なら1回と数える。
      *----------------------------------------------------------------
       ｷｹﾞﾂ-RTN.
           IF WS-ﾊﾗｲｺﾐ-MM > 12
               SUBTRACT 12 FROM WS-ﾊﾗｲｺﾐ-MM
               ADD 1 TO WS-ﾊﾗｲｺﾐ-YYYY
           END-IF
           IF WS-ﾊﾗｲｺﾐ-YYYYMM NOT > WS-ﾈﾝﾏﾂ
               ADD 1 TO WS-ﾖﾃｲｶｲｽｳ
           END-IF.

      *----------------------------------------------------------------
      * ｶｷｺﾐ-RTN: 証明書1通と明細行。
      *----------------------------------------------------------------
       ｶｷｺﾐ-RTN.
           MOVE SPACES TO SM-ｼｮｳﾒｲ-ﾚｺｰﾄﾞ
           MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO SM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           MOVE WS-ｷﾞｮｳﾑ-YYYY TO SM-ﾈﾝﾌﾞﾝ
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO SM-ﾊｯｺｳﾋﾞ
           MOVE PM-ｼﾒｲ TO SM-ｼﾒｲ
           MOVE PM-ｼﾞｭｳｼｮ TO SM-ｼﾞｭｳｼｮ
           MOVE PM-ｼｮｳﾋﾝ TO SM-ｼｮｳﾋﾝ
           SET WS-ｼｮｳﾋﾝ-IX TO 1
           SEARCH WS-ｼｮｳﾋﾝ
               WHEN WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = PM-ｼｮｳﾋﾝ
                   MOVE WS-ｼｮｳﾋﾝ-ﾒｲｼｮｳ(WS-ｼｮｳﾋﾝ-IX) TO SM-ｼｮｳﾋﾝﾒｲ
           END-SEARCH
           MOVE PM-ｹｲﾔｸﾋﾞ TO SM-ｹｲﾔｸﾋﾞ
           MOVE PM-ﾊﾗｲｺﾐﾎｳﾎｳ TO SM-ﾊﾗｲｺﾐﾎｳﾎｳ
           MOVE WS-ｼｮｳﾒｲｶﾞｸ TO SM-ｼｮｳﾒｲｶﾞｸ
           COMPUTE SM-ﾖﾃｲｶﾞｸ = WS-ｼｮｳﾒｲｶﾞｸ
                             + PM-ﾎｹﾝﾘｮｳ * WS-ﾖﾃｲｶｲｽｳ
           WRITE SM-ｼｮｳﾒｲ-ﾚｺｰﾄﾞ
           IF NOT WS-PRMCERT-OK
               MOVE "PRMCERT" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-PRMCERT-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           ADD 1 TO WS-ｼｮｳﾒｲｹﾝｽｳ
           ADD SM-ｼｮｳﾒｲｶﾞｸ TO WS-ｼｮｳﾒｲｺﾞｳｹｲ
           ADD SM-ﾖﾃｲｶﾞｸ TO WS-ﾖﾃｲｺﾞｳｹｲ
           MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO WS-ﾒｲｻｲ-ｼｮｳｹﾝ
           MOVE PM-ｼｮｳﾋﾝ TO WS-ﾒｲｻｲ-ｼｮｳﾋﾝ
           MOVE PM-ｼﾒｲ TO WS-ﾒｲｻｲ-ｼﾒｲ
           MOVE SM-ｼｮｳﾒｲｶﾞｸ TO WS-ﾒｲｻｲ-ｼｮｳﾒｲｶﾞｸ
           MOVE SM-ﾖﾃｲｶﾞｸ TO WS-ﾒｲｻｲ-ﾖﾃｲｶﾞｸ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ﾒｲｻｲ-ｷﾞｮｳ.

      *----------------------------------------------------------------
      * ﾎｳｺｸ-RTN: 件数と合計。
      *----------------------------------------------------------------
       ﾎｳｺｸ-RTN.
           MOVE SPACES TO ﾎｳｺｸ-ｷﾞｮｳ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ
           MOVE "収納履歴件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾘﾚｷｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "証明書件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｼｮｳﾒｲｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "契約なし" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｹｲﾔｸﾅｼｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "証明額合計" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｼｮｳﾒｲｺﾞｳｹｲ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "予定額合計" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾖﾃｲｺﾞｳｹｲ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ.

      *----------------------------------------------------------------
      * ｼｭｳﾘｮｳ-RTN: ファイルを閉じ、件数を表示する。
      *----------------------------------------------------------------
       ｼｭｳﾘｮｳ-RTN.
           CLOSE PRMHST-F POLMST-F PRODTB-F PRMCERT-F REPORT-F
           DISPLAY "GLY030 証明書 " WS-ｼｮｳﾒｲｹﾝｽｳ.

      *----------------------------------------------------------------
      * ｲｼﾞｮｳ-RTN: ファイル異常は 16 で終わる。
      *----------------------------------------------------------------
       ｲｼﾞｮｳ-RTN.
           DISPLAY "GLY030 ファイル異常 " WS-ｴﾗｰﾌｧｲﾙ
                   " 状態 " WS-ｴﾗｰｽﾃｰﾀｽ
           MOVE 16 TO RETURN-CODE
           STOP RUN.

       COPY PRDLOD.
