      *----------------------------------------------------------------
      * GLD060: 失効処理
      *
      * 日次。猶予中の契約のうち、払込猶予期間を過ぎたものを失効
      * させる。猶予期間は猶予開始日 (振替不能となった日) の属する
      * 月の翌月末まで: 業務月が猶予開始月の翌月より後なら失効。
      * 失効した契約は失効日を業務日付とし、1年以内なら GLD040 の
      * 復活で戻せる。
      *
      * 復帰コード: 0 処理済、8 日付が日付でない、16 ファイル異常。
      *
      * 変更履歴
      * 日付       担当 変更
      * 1987-04-01 KT   GL-0001 初版。
      * 2003-10-06 MS   GL-0310 猶予開始日を振替不能日から。
      *----------------------------------------------------------------
       IDENTIFICATION DIVISION.
       PROGRAM-ID. GLD060.
       AUTHOR. ゴンドワナ生命 システム部.
       DATE-WRITTEN. 1987-04-01.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT POLMST-F ASSIGN TO "POLMST"
               ORGANIZATION IS INDEXED
               ACCESS MODE IS DYNAMIC
               RECORD KEY IS PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
               FILE STATUS IS WS-POLMST-ST.
           SELECT REPORT-F ASSIGN TO "GLD060R"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-ST.

       DATA DIVISION.
       FILE SECTION.
       FD  POLMST-F.
       COPY POLMST.
       FD  REPORT-F.
       01  ﾎｳｺｸ-ｷﾞｮｳ                   PIC X(100).

       WORKING-STORAGE SECTION.
       01  WS-ｽﾃｰﾀｽ.
           05  WS-POLMST-ST            PIC XX.
               88  WS-POLMST-OK                VALUE "00".
               88  WS-POLMST-ｵﾜﾘ               VALUE "10".
           05  WS-REPORT-ST            PIC XX.
               88  WS-REPORT-OK                VALUE "00".
       01  WS-ｺﾏﾝﾄﾞｷﾞｮｳ                PIC X(80).
       01  WS-ｴﾗｰﾌｧｲﾙ                  PIC X(8).
       01  WS-ｴﾗｰｽﾃｰﾀｽ                 PIC XX.
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ                PIC 9(8).
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾌﾞﾝｶｲ REDEFINES WS-ｷﾞｮｳﾑﾋﾂﾞｹ.
           05  WS-ｷﾞｮｳﾑ-YYYYMM         PIC 9(6).
           05  WS-ｷﾞｮｳﾑ-DD             PIC 99.
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ            PIC X(7).
       01  WS-ﾎｳｺｳ                     PIC X.
       01  WS-ﾍﾝｶﾝｹｯｶ                  PIC 9.
       01  WS-ﾕｳﾖｶｲｼﾋﾞ                 PIC 9(8).
       01  WS-ﾕｳﾖ-ﾌﾞﾝｶｲ REDEFINES WS-ﾕｳﾖｶｲｼﾋﾞ.
           05  WS-ﾕｳﾖ-YYYY             PIC 9(4).
           05  WS-ﾕｳﾖ-MM               PIC 99.
           05  WS-ﾕｳﾖ-DD               PIC 99.
       01  WS-ｷｹﾞﾝ.
           05  WS-ｷｹﾞﾝ-YYYY            PIC 9(4).
           05  WS-ｷｹﾞﾝ-MM              PIC 99.
       01  WS-ｷｹﾞﾝ-YYYYMM REDEFINES WS-ｷｹﾞﾝ PIC 9(6).
       01  WS-ｹﾝｽｳ.
           05  WS-ﾕｳﾖｹﾝｽｳ              PIC 9(7) VALUE 0.
           05  WS-ｼｯｺｳｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ｼｯｺｳｷﾝｶﾞｸ            PIC 9(13) VALUE 0.
       01  WS-ﾎｳｺｸ.
           05  WS-ﾋｮｳﾀﾞｲ-ｷﾞｮｳ.
               10  FILLER                  PIC X(28)
                   VALUE "ゴンドワナ生命  失効処理".
               10  FILLER                  PIC X(4) VALUE SPACES.
               10  WS-ﾋｮｳﾀﾞｲ-ﾋﾂﾞｹ          PIC X(7).
           05  WS-ﾒｲｻｲ-ｷﾞｮｳ.
               10  WS-ﾒｲｻｲ-ｼｮｳｹﾝ           PIC X(10).
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾒｲｻｲ-ｼﾒｲ             PIC X(20).
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  FILLER                  PIC X(8) VALUE "猶予開始".
               10  WS-ﾒｲｻｲ-ﾕｳﾖ             PIC X(7).
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾒｲｻｲ-ｹｯｶ             PIC X(12).
           05  WS-ｹﾝｽｳ-ｷﾞｮｳ.
               10  WS-ｹﾝｽｳ-ﾗﾍﾞﾙ            PIC X(16).
               10  WS-ｹﾝｽｳ-ｱﾀｲ             PIC ZZZ,ZZZ,ZZ9.

       PROCEDURE DIVISION.
      *----------------------------------------------------------------
      * ﾒｲﾝ-RTN: 準備、マスタ全件の走査、報告、終了。
      *----------------------------------------------------------------
       ﾒｲﾝ-RTN.
           PERFORM ｼｮｷ-RTN
           PERFORM ｼｮﾘ-RTN UNTIL WS-POLMST-ｵﾜﾘ
           PERFORM ﾎｳｺｸ-RTN
           PERFORM ｼｭｳﾘｮｳ-RTN
           MOVE 0 TO RETURN-CODE
           STOP RUN.

      *----------------------------------------------------------------
      * ｼｮｷ-RTN: 業務日付とファイル。
      *----------------------------------------------------------------
       ｼｮｷ-RTN.
           ACCEPT WS-ｺﾏﾝﾄﾞｷﾞｮｳ FROM COMMAND-LINE
           MOVE WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8) TO WS-ｷﾞｮｳﾑﾋﾂﾞｹ
           MOVE "W" TO WS-ﾎｳｺｳ
           CALL "GLC020" USING WS-ﾎｳｺｳ WS-ｷﾞｮｳﾑﾋﾂﾞｹ
                               WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ WS-ﾍﾝｶﾝｹｯｶ
           IF WS-ﾍﾝｶﾝｹｯｶ NOT = 0
               DISPLAY "GLD060 日付不正: " WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8)
               MOVE 8 TO RETURN-CODE
               STOP RUN
           END-IF
           OPEN I-O POLMST-F
           IF NOT WS-POLMST-OK
               MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           OPEN OUTPUT REPORT-F
           IF NOT WS-REPORT-OK
               MOVE "GLD060R" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-REPORT-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO WS-ﾋｮｳﾀﾞｲ-ﾋﾂﾞｹ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ﾋｮｳﾀﾞｲ-ｷﾞｮｳ
           MOVE SPACES TO ﾎｳｺｸ-ｷﾞｮｳ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ
           PERFORM POLMST-ﾖﾐｺﾐ-RTN.

      *----------------------------------------------------------------
      * POLMST-ﾖﾐｺﾐ-RTN: マスタの次の契約。
      *----------------------------------------------------------------
       POLMST-ﾖﾐｺﾐ-RTN.
           READ POLMST-F NEXT
           IF NOT WS-POLMST-OK AND NOT WS-POLMST-ｵﾜﾘ
               MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ｼｮﾘ-RTN: 猶予中の契約の猶予期限を求め、過ぎていれば失効。
      *----------------------------------------------------------------
       ｼｮﾘ-RTN.
           IF PM-ﾕｳﾖﾁｭｳ
               ADD 1 TO WS-ﾕｳﾖｹﾝｽｳ
               PERFORM ｷｹﾞﾝ-RTN
               IF WS-ｷﾞｮｳﾑ-YYYYMM > WS-ｷｹﾞﾝ-YYYYMM
                   PERFORM ｼｯｺｳ-RTN
               ELSE
                   MOVE "猶予中" TO WS-ﾒｲｻｲ-ｹｯｶ
               END-IF
               PERFORM ﾒｲｻｲ-RTN
           END-IF
           PERFORM POLMST-ﾖﾐｺﾐ-RTN.

      *----------------------------------------------------------------
      * ｷｹﾞﾝ-RTN: 猶予開始月の翌月が猶予期限の月。
      *----------------------------------------------------------------
       ｷｹﾞﾝ-RTN.
           MOVE "S" TO WS-ﾎｳｺｳ
           CALL "GLC020" USING WS-ﾎｳｺｳ WS-ﾕｳﾖｶｲｼﾋﾞ
                               PM-ﾕｳﾖｶｲｼﾋﾞ WS-ﾍﾝｶﾝｹｯｶ
           MOVE WS-ﾕｳﾖ-YYYY TO WS-ｷｹﾞﾝ-YYYY
           MOVE WS-ﾕｳﾖ-MM TO WS-ｷｹﾞﾝ-MM
           IF WS-ｷｹﾞﾝ-MM = 12
               MOVE 1 TO WS-ｷｹﾞﾝ-MM
               ADD 1 TO WS-ｷｹﾞﾝ-YYYY
           ELSE
               ADD 1 TO WS-ｷｹﾞﾝ-MM
           END-IF.

      *----------------------------------------------------------------
      * ｼｯｺｳ-RTN: 失効にして書き戻す。
      *----------------------------------------------------------------
       ｼｯｺｳ-RTN.
           SET PM-ｼｯｺｳ TO TRUE
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO PM-ｼｯｺｳﾋﾞ PM-ｼﾞｮｳﾀｲﾋﾂﾞｹ
           REWRITE PM-ｹｲﾔｸ-ﾚｺｰﾄﾞ
           IF NOT WS-POLMST-OK
               MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           ADD 1 TO WS-ｼｯｺｳｹﾝｽｳ
           ADD PM-ﾎｹﾝｷﾝｶﾞｸ TO WS-ｼｯｺｳｷﾝｶﾞｸ
           MOVE "失効" TO WS-ﾒｲｻｲ-ｹｯｶ.

      *----------------------------------------------------------------
      * ﾒｲｻｲ-RTN: 1件の明細行。
      *----------------------------------------------------------------
       ﾒｲｻｲ-RTN.
           MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO WS-ﾒｲｻｲ-ｼｮｳｹﾝ
           MOVE PM-ｼﾒｲ TO WS-ﾒｲｻｲ-ｼﾒｲ
           MOVE PM-ﾕｳﾖｶｲｼﾋﾞ TO WS-ﾒｲｻｲ-ﾕｳﾖ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ﾒｲｻｲ-ｷﾞｮｳ.

      *----------------------------------------------------------------
      * ﾎｳｺｸ-RTN: 件数。
      *----------------------------------------------------------------
       ﾎｳｺｸ-RTN.
           MOVE SPACES TO ﾎｳｺｸ-ｷﾞｮｳ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ
           MOVE "猶予中件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾕｳﾖｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "失効件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｼｯｺｳｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "失効保険金額" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｼｯｺｳｷﾝｶﾞｸ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ.

      *----------------------------------------------------------------
      * ｼｭｳﾘｮｳ-RTN: ファイルを閉じ、件数を表示する。
      *----------------------------------------------------------------
       ｼｭｳﾘｮｳ-RTN.
           CLOSE POLMST-F REPORT-F
           DISPLAY "GLD060 猶予中 " WS-ﾕｳﾖｹﾝｽｳ " 失効 " WS-ｼｯｺｳｹﾝｽｳ.

      *----------------------------------------------------------------
      * ｲｼﾞｮｳ-RTN: ファイル異常は 16 で終わる。
      *----------------------------------------------------------------
       ｲｼﾞｮｳ-RTN.
           DISPLAY "GLD060 ファイル異常 " WS-ｴﾗｰﾌｧｲﾙ
                   " 状態 " WS-ｴﾗｰｽﾃｰﾀｽ
           MOVE 16 TO RETURN-CODE
           STOP RUN.
