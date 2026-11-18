      *----------------------------------------------------------------
      * GLD090: 契約マスタ アンロード
      *
      * 日次と月次の最後。契約マスタを証券番号順に1行1レコードで
      * アンロードし、下流 (営業システム、経理) に渡す。制御レコード
      * の業務日付を処理した日に進め、件数と有効保険金額を書き戻す。
      * 報告には制御レコードの内容を印字する。
      * 給付履歴 BENHST も、あれば証券番号・連番順にアンロードする。
      * 保障マスタ COVMST は証券番号・保障番号順にアンロードする。
      *
      * 復帰コード: 0 処理済、8 日付が日付でない、16 ファイル異常。
      *
      * 変更履歴
      * 日付       担当 変更
      * 1987-04-01 KT   GL-0001 初版。
      * 1999-11-08 MS   GL-0233 制御レコードの更新。
      * 2026-11-16 YN   GL-0822 給付履歴のアンロード。
      * 2026-11-18 YN   GL-0825 保障マスタのアンロード。
      *----------------------------------------------------------------
       IDENTIFICATION DIVISION.
       PROGRAM-ID. GLD090.
       AUTHOR. ゴンドワナ生命 システム部.
       DATE-WRITTEN. 1987-04-01.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT POLMST-F ASSIGN TO "POLMST"
               ORGANIZATION IS INDEXED
               ACCESS MODE IS SEQUENTIAL
               RECORD KEY IS PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
               FILE STATUS IS WS-POLMST-ST.
           SELECT UNLOAD-F ASSIGN TO "UNLOAD"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-UNLOAD-ST.
           SELECT OPTIONAL BENHST-F ASSIGN TO "BENHST"
               ORGANIZATION IS INDEXED
               ACCESS MODE IS SEQUENTIAL
               RECORD KEY IS BH-ｷｰ
               FILE STATUS IS WS-BENHST-ST.
           SELECT HSTUNL-F ASSIGN TO "HSTUNL"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-HSTUNL-ST.
           SELECT COVMST-F ASSIGN TO "COVMST"
               ORGANIZATION IS INDEXED
               ACCESS MODE IS SEQUENTIAL
               RECORD KEY IS HS-ｷｰ
               FILE STATUS IS WS-COVMST-ST.
           SELECT COVUNL-F ASSIGN TO "COVUNL"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-COVUNL-ST.
           SELECT GLSCTL-F ASSIGN TO "GLSCTL"
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-GLSCTL-ST.
           SELECT REPORT-F ASSIGN TO "GLD090R"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-ST.

       DATA DIVISION.
       FILE SECTION.
       FD  POLMST-F.
       COPY POLMST.
       FD  UNLOAD-F.
       01  ｱﾝﾛｰﾄﾞ-ﾚｺｰﾄﾞ                PIC X(300).
       FD  BENHST-F.
       COPY BENHST.
       FD  HSTUNL-F.
       01  ﾘﾚｷ-ｱﾝﾛｰﾄﾞ-ﾚｺｰﾄﾞ            PIC X(60).
       FD  COVMST-F.
       COPY COVMST.
       FD  COVUNL-F.
       01  ﾎｼｮｳ-ｱﾝﾛｰﾄﾞ-ﾚｺｰﾄﾞ           PIC X(80).
       FD  GLSCTL-F
           RECORD CONTAINS 64 CHARACTERS.
       COPY CTLREC.
       FD  REPORT-F.
       01  ﾎｳｺｸ-ｷﾞｮｳ                   PIC X(80).

       WORKING-STORAGE SECTION.
       01  WS-ｽﾃｰﾀｽ.
           05  WS-POLMST-ST            PIC XX.
               88  WS-POLMST-OK                VALUE "00".
               88  WS-POLMST-ｵﾜﾘ               VALUE "10".
           05  WS-UNLOAD-ST            PIC XX.
               88  WS-UNLOAD-OK                VALUE "00".
           05  WS-BENHST-ST            PIC XX.
               88  WS-BENHST-OK                VALUE "00" "05".
               88  WS-BENHST-ｵﾜﾘ               VALUE "10".
           05  WS-HSTUNL-ST            PIC XX.
               88  WS-HSTUNL-OK                VALUE "00".
           05  WS-COVMST-ST            PIC XX.
               88  WS-COVMST-OK                VALUE "00".
               88  WS-COVMST-ｵﾜﾘ               VALUE "10".
           05  WS-COVUNL-ST            PIC XX.
               88  WS-COVUNL-OK                VALUE "00".
           05  WS-GLSCTL-ST            PIC XX.
               88  WS-GLSCTL-OK                VALUE "00".
           05  WS-REPORT-ST            PIC XX.
               88  WS-REPORT-OK                VALUE "00".
       01  WS-ｺﾏﾝﾄﾞｷﾞｮｳ                PIC X(80).
       01  WS-ｴﾗｰﾌｧｲﾙ                  PIC X(8).
       01  WS-ｴﾗｰｽﾃｰﾀｽ                 PIC XX.
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ                PIC 9(8).
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ            PIC X(7).
       01  WS-ﾎｳｺｳ                     PIC X VALUE "W".
       01  WS-ﾍﾝｶﾝｹｯｶ                  PIC 9.
       01  WS-ｹﾝｽｳ.
           05  WS-ｿｳｹﾝｽｳ               PIC 9(7) VALUE 0.
           05  WS-ﾕｳｺｳｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ﾕｳｺｳﾎｹﾝｷﾝｶﾞｸ         PIC 9(15) VALUE 0.
           05  WS-ﾘﾚｷｹﾝｽｳ              PIC 9(7) VALUE 0.
           05  WS-ﾎｼｮｳｹﾝｽｳ             PIC 9(7) VALUE 0.
       01  WS-ﾎｳｺｸ.
           05  WS-ﾋｮｳﾀﾞｲ-ｷﾞｮｳ.
               10  FILLER                  PIC X(28)
                   VALUE "ゴンドワナ生命  アンロード".
               10  FILLER                  PIC X(4) VALUE SPACES.
               10  WS-ﾋｮｳﾀﾞｲ-ﾋﾂﾞｹ          PIC X(7).
           05  WS-ｹﾝｽｳ-ｷﾞｮｳ.
               10  WS-ｹﾝｽｳ-ﾗﾍﾞﾙ            PIC X(16).
               10  WS-ｹﾝｽｳ-ｱﾀｲ             PIC ZZZ,ZZZ,ZZZ,ZZ9.
           05  WS-ﾋﾂﾞｹ-ｷﾞｮｳ.
               10  WS-ﾋﾂﾞｹ-ﾗﾍﾞﾙ            PIC X(16).
               10  WS-ﾋﾂﾞｹ-ｱﾀｲ             PIC X(8).

       PROCEDURE DIVISION.
      *----------------------------------------------------------------
      * ﾒｲﾝ-RTN: 準備、全件のアンロード、制御レコード、報告、終了。
      *----------------------------------------------------------------
       ﾒｲﾝ-RTN.
           PERFORM ｼｮｷ-RTN
           PERFORM ｱﾝﾛｰﾄﾞ-RTN UNTIL WS-POLMST-ｵﾜﾘ
           PERFORM BENHST-ﾖﾐｺﾐ-RTN
           PERFORM ﾘﾚｷ-ｱﾝﾛｰﾄﾞ-RTN UNTIL WS-BENHST-ｵﾜﾘ
           PERFORM COVMST-ﾖﾐｺﾐ-RTN
           PERFORM ﾎｼｮｳ-ｱﾝﾛｰﾄﾞ-RTN UNTIL WS-COVMST-ｵﾜﾘ
           PERFORM ｾｲｷﾞｮ-RTN
           PERFORM ﾎｳｺｸ-RTN
           PERFORM ｼｭｳﾘｮｳ-RTN
           MOVE 0 TO RETURN-CODE
           STOP RUN.

      *----------------------------------------------------------------
      * ｼｮｷ-RTN: 業務日付、ファイル、制御レコード。
      *----------------------------------------------------------------
       ｼｮｷ-RTN.
           ACCEPT WS-ｺﾏﾝﾄﾞｷﾞｮｳ FROM COMMAND-LINE
           MOVE WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8) TO WS-ｷﾞｮｳﾑﾋﾂﾞｹ
           CALL "GLC020" USING WS-ﾎｳｺｳ WS-ｷﾞｮｳﾑﾋﾂﾞｹ
                               WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ WS-ﾍﾝｶﾝｹｯｶ
           IF WS-ﾍﾝｶﾝｹｯｶ NOT = 0
               DISPLAY "GLD090 日付不正: " WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8)
               MOVE 8 TO RETURN-CODE
               STOP RUN
           END-IF
           OPEN INPUT POLMST-F BENHST-F COVMST-F
           OPEN I-O GLSCTL-F
           OPEN OUTPUT UNLOAD-F HSTUNL-F COVUNL-F REPORT-F
           EVALUATE TRUE
               WHEN NOT WS-COVMST-OK
                   MOVE "COVMST" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-COVMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-COVUNL-OK
                   MOVE "COVUNL" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-COVUNL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-BENHST-OK
                   MOVE "BENHST" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-BENHST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-HSTUNL-OK
                   MOVE "HSTUNL" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-HSTUNL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-POLMST-OK
                   MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-GLSCTL-OK
                   MOVE "GLSCTL" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-GLSCTL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-UNLOAD-OK
                   MOVE "UNLOAD" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-UNLOAD-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-REPORT-OK
                   MOVE "GLD090R" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-REPORT-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN OTHER
                   MOVE SPACES TO WS-ｴﾗｰﾌｧｲﾙ
           END-EVALUATE
           IF WS-ｴﾗｰﾌｧｲﾙ NOT = SPACES
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           READ GLSCTL-F
           IF NOT WS-GLSCTL-OK
               MOVE "GLSCTL" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-GLSCTL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
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
      * ｱﾝﾛｰﾄﾞ-RTN: 1件を書き、数える。
      *----------------------------------------------------------------
       ｱﾝﾛｰﾄﾞ-RTN.
           WRITE ｱﾝﾛｰﾄﾞ-ﾚｺｰﾄﾞ FROM PM-ｹｲﾔｸ-ﾚｺｰﾄﾞ
           IF NOT WS-UNLOAD-OK
               MOVE "UNLOAD" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-UNLOAD-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           ADD 1 TO WS-ｿｳｹﾝｽｳ
           IF PM-ﾎｼｮｳﾁｭｳ
               ADD 1 TO WS-ﾕｳｺｳｹﾝｽｳ
               ADD PM-ﾎｹﾝｷﾝｶﾞｸ TO WS-ﾕｳｺｳﾎｹﾝｷﾝｶﾞｸ
           END-IF
           PERFORM POLMST-ﾖﾐｺﾐ-RTN.

      *----------------------------------------------------------------
      * BENHST-ﾖﾐｺﾐ-RTN: 給付履歴の次の1件。
      *----------------------------------------------------------------
       BENHST-ﾖﾐｺﾐ-RTN.
           READ BENHST-F NEXT
           IF NOT WS-BENHST-OK AND NOT WS-BENHST-ｵﾜﾘ
               MOVE "BENHST" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-BENHST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ﾘﾚｷ-ｱﾝﾛｰﾄﾞ-RTN: 給付履歴1件を書き、数える。
      *----------------------------------------------------------------
       ﾘﾚｷ-ｱﾝﾛｰﾄﾞ-RTN.
           WRITE ﾘﾚｷ-ｱﾝﾛｰﾄﾞ-ﾚｺｰﾄﾞ FROM BH-ﾘﾚｷ-ﾚｺｰﾄﾞ
           IF NOT WS-HSTUNL-OK
               MOVE "HSTUNL" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-HSTUNL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           ADD 1 TO WS-ﾘﾚｷｹﾝｽｳ
           PERFORM BENHST-ﾖﾐｺﾐ-RTN.

      *----------------------------------------------------------------
      * COVMST-ﾖﾐｺﾐ-RTN: 保障マスタの次の1件。
      *----------------------------------------------------------------
       COVMST-ﾖﾐｺﾐ-RTN.
           READ COVMST-F NEXT
           IF NOT WS-COVMST-OK AND NOT WS-COVMST-ｵﾜﾘ
               MOVE "COVMST" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-COVMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ﾎｼｮｳ-ｱﾝﾛｰﾄﾞ-RTN: 保障1件を書き、数える。
      *----------------------------------------------------------------
       ﾎｼｮｳ-ｱﾝﾛｰﾄﾞ-RTN.
           WRITE ﾎｼｮｳ-ｱﾝﾛｰﾄﾞ-ﾚｺｰﾄﾞ FROM HS-ﾎｼｮｳ-ﾚｺｰﾄﾞ
           IF NOT WS-COVUNL-OK
               MOVE "COVUNL" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-COVUNL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           ADD 1 TO WS-ﾎｼｮｳｹﾝｽｳ
           PERFORM COVMST-ﾖﾐｺﾐ-RTN.

      *----------------------------------------------------------------
      * ｾｲｷﾞｮ-RTN: 業務日付と件数を制御レコードに書き戻す。
      *----------------------------------------------------------------
       ｾｲｷﾞｮ-RTN.
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ TO CT-ｷﾞｮｳﾑﾋﾂﾞｹ
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO CT-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ
           MOVE WS-ﾕｳｺｳｹﾝｽｳ TO CT-ﾕｳｺｳｹﾝｽｳ
           MOVE WS-ｿｳｹﾝｽｳ TO CT-ｿｳｹﾝｽｳ
           MOVE WS-ﾕｳｺｳﾎｹﾝｷﾝｶﾞｸ TO CT-ﾕｳｺｳﾎｹﾝｷﾝｶﾞｸ
           REWRITE CT-ｾｲｷﾞｮ-ﾚｺｰﾄﾞ
           IF NOT WS-GLSCTL-OK
               MOVE "GLSCTL" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-GLSCTL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ﾎｳｺｸ-RTN: 件数と、書き戻した制御レコードの内容。
      *----------------------------------------------------------------
       ﾎｳｺｸ-RTN.
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO WS-ﾋｮｳﾀﾞｲ-ﾋﾂﾞｹ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ﾋｮｳﾀﾞｲ-ｷﾞｮｳ
           MOVE SPACES TO ﾎｳｺｸ-ｷﾞｮｳ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ
           MOVE "総件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｿｳｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "有効件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾕｳｺｳｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "有効保険金額" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾕｳｺｳﾎｹﾝｷﾝｶﾞｸ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "給付履歴件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾘﾚｷｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "保障件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾎｼｮｳｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE SPACES TO ﾎｳｺｸ-ｷﾞｮｳ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ
           MOVE "制御: 業務日付" TO WS-ﾋﾂﾞｹ-ﾗﾍﾞﾙ
           MOVE CT-ｷﾞｮｳﾑﾋﾂﾞｹ TO WS-ﾋﾂﾞｹ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ﾋﾂﾞｹ-ｷﾞｮｳ
           MOVE "制御: 和暦" TO WS-ﾋﾂﾞｹ-ﾗﾍﾞﾙ
           MOVE CT-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO WS-ﾋﾂﾞｹ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ﾋﾂﾞｹ-ｷﾞｮｳ
           MOVE "制御: 請求年月" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE CT-ｻｲｼｭｳｾｲｷｭｳﾈﾝｹﾞﾂ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "制御: 次回支払" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE CT-ｼﾞｶｲｼﾊﾗｲﾊﾞﾝｺﾞｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "制御: 次回仕訳" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE CT-ｼﾞｶｲｼﾜｹﾊﾞﾝｺﾞｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ.

      *----------------------------------------------------------------
      * ｼｭｳﾘｮｳ-RTN: ファイルを閉じ、件数を表示する。
      *----------------------------------------------------------------
       ｼｭｳﾘｮｳ-RTN.
           CLOSE POLMST-F BENHST-F UNLOAD-F HSTUNL-F COVMST-F
                 COVUNL-F GLSCTL-F REPORT-F
           DISPLAY "GLD090 アンロード " WS-ｿｳｹﾝｽｳ " 有効 " WS-ﾕｳｺｳｹﾝｽｳ.

      *----------------------------------------------------------------
      * ｲｼﾞｮｳ-RTN: ファイル異常は 16 で終わる。
      *----------------------------------------------------------------
       ｲｼﾞｮｳ-RTN.
           DISPLAY "GLD090 ファイル異常 " WS-ｴﾗｰﾌｧｲﾙ
                   " 状態 " WS-ｴﾗｰｽﾃｰﾀｽ
           MOVE 16 TO RETURN-CODE
           STOP RUN.
