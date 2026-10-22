      *----------------------------------------------------------------
      * GLO010: システム開設: 契約マスタの初期登録
      *
      * システム開設時に1回だけ、コマンド行の業務日付 YYYYMMDD で
      * 実行する。契約抽出ファイル (1行1マスタレコード) を読み、
      * 証券番号をキーとする索引マスタに登録する。次に制御レコード
      * を、両暦の日付、採番カウンタ 1、登録した件数で書く。
      * 抽出に二重にある証券番号は1件だけ登録し、数える。
      *
      * 復帰コード: 0 登録済、4 抽出に重複あり、
      *             8 日付が日付でない、16 ファイル異常。
      *
      * 変更履歴
      * 日付       担当 変更
      * 1999-11-08 MS   GL-0233 初版、PC システム移行時。
      *----------------------------------------------------------------
       IDENTIFICATION DIVISION.
       PROGRAM-ID. GLO010.
       AUTHOR. ゴンドワナ生命 システム部.
       DATE-WRITTEN. 1999-11-08.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT EXTRACT-F ASSIGN TO "EXTRACT"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-ﾁｭｳｼｭﾂ-ST.
           SELECT POLMST-F ASSIGN TO "POLMST"
               ORGANIZATION IS INDEXED
               ACCESS MODE IS RANDOM
               RECORD KEY IS PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
               FILE STATUS IS WS-ｹｲﾔｸ-ST.
           SELECT GLSCTL-F ASSIGN TO "GLSCTL"
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-ｾｲｷﾞｮ-ST.
           SELECT REPORT-F ASSIGN TO "GLO010R"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-ﾎｳｺｸ-ST.

       DATA DIVISION.
       FILE SECTION.
       FD  EXTRACT-F.
       01  ﾁｭｳｼｭﾂ-ﾚｺｰﾄﾞ                PIC X(300).
       FD  POLMST-F.
       COPY POLMST.
       FD  GLSCTL-F
           RECORD CONTAINS 64 CHARACTERS.
       COPY CTLREC.
       FD  REPORT-F.
       01  ﾎｳｺｸ-ｷﾞｮｳ                   PIC X(80).

       WORKING-STORAGE SECTION.
       01  WS-ｽﾃｰﾀｽ.
           05  WS-ﾁｭｳｼｭﾂ-ST            PIC XX.
               88  WS-ﾁｭｳｼｭﾂ-OK                VALUE "00".
               88  WS-ﾁｭｳｼｭﾂ-ｵﾜﾘ               VALUE "10".
           05  WS-ｹｲﾔｸ-ST              PIC XX.
               88  WS-ｹｲﾔｸ-OK                  VALUE "00".
           05  WS-ｾｲｷﾞｮ-ST             PIC XX.
               88  WS-ｾｲｷﾞｮ-OK                 VALUE "00".
           05  WS-ﾎｳｺｸ-ST              PIC XX.
               88  WS-ﾎｳｺｸ-OK                  VALUE "00".
       01  WS-ｺﾏﾝﾄﾞｷﾞｮｳ                PIC X(80).
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ                PIC 9(8).
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ            PIC X(7).
       01  WS-ﾎｳｺｳ                     PIC X VALUE "W".
       01  WS-ﾍﾝｶﾝｹｯｶ                  PIC 9.
       01  WS-ｹﾝｽｳ.
           05  WS-ﾖﾐｺﾐｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ﾄｳﾛｸｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ﾁｮｳﾌｸｹﾝｽｳ            PIC 9(7) VALUE 0.
           05  WS-ﾕｳｺｳｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ﾕｳﾖｹﾝｽｳ              PIC 9(7) VALUE 0.
           05  WS-ｼｯｺｳｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ｼｮｳﾒﾂｹﾝｽｳ            PIC 9(7) VALUE 0.
           05  WS-ﾕｳｺｳﾎｹﾝｷﾝｶﾞｸ         PIC 9(15) VALUE 0.
       01  WS-ﾌｯｷｺｰﾄﾞ                  PIC 99 VALUE 0.
       01  WS-ﾎｳｺｸ.
           05  WS-ﾋｮｳﾀﾞｲ-ｷﾞｮｳ.
               10  FILLER                  PIC X(28)
                   VALUE "ゴンドワナ生命  システム開設".
               10  FILLER                  PIC X(4) VALUE SPACES.
               10  WS-ﾋｮｳﾀﾞｲ-ﾋﾂﾞｹ          PIC X(7).
           05  WS-ｹﾝｽｳ-ｷﾞｮｳ.
               10  WS-ｹﾝｽｳ-ﾗﾍﾞﾙ            PIC X(16).
               10  WS-ｹﾝｽｳ-ｱﾀｲ             PIC Z,ZZZ,ZZ9.
           05  WS-ｷﾝｶﾞｸ-ｷﾞｮｳ.
               10  FILLER                  PIC X(16)
                   VALUE "有効保険金額".
               10  WS-ｷﾝｶﾞｸ-ｱﾀｲ            PIC ZZZ,ZZZ,ZZZ,ZZ9.
               10  FILLER                  PIC X(2) VALUE "円".

       PROCEDURE DIVISION.
      *----------------------------------------------------------------
      * ﾒｲﾝ-RTN: 日付、登録、制御レコード、報告。
      *----------------------------------------------------------------
       ﾒｲﾝ-RTN.
           PERFORM ｼｮｷ-RTN
           PERFORM ﾄｳﾛｸ-RTN UNTIL WS-ﾁｭｳｼｭﾂ-ｵﾜﾘ
           PERFORM ｾｲｷﾞｮ-RTN
           PERFORM ﾎｳｺｸ-RTN
           PERFORM ｼｭｳﾘｮｳ-RTN
           MOVE WS-ﾌｯｷｺｰﾄﾞ TO RETURN-CODE
           STOP RUN.

      *----------------------------------------------------------------
      * ｼｮｷ-RTN: コマンド行の業務日付を両暦にし、ファイルを開く。
      *----------------------------------------------------------------
       ｼｮｷ-RTN.
           ACCEPT WS-ｺﾏﾝﾄﾞｷﾞｮｳ FROM COMMAND-LINE
           MOVE WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8) TO WS-ｷﾞｮｳﾑﾋﾂﾞｹ
           CALL "GLC020" USING WS-ﾎｳｺｳ WS-ｷﾞｮｳﾑﾋﾂﾞｹ
                               WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ WS-ﾍﾝｶﾝｹｯｶ
           IF WS-ﾍﾝｶﾝｹｯｶ NOT = 0
               DISPLAY "GLO010 日付不正: " WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8)
               MOVE 8 TO RETURN-CODE
               STOP RUN
           END-IF
           OPEN INPUT EXTRACT-F
           IF NOT WS-ﾁｭｳｼｭﾂ-OK
               MOVE "EXTRACT" TO WS-ｺﾏﾝﾄﾞｷﾞｮｳ
               MOVE WS-ﾁｭｳｼｭﾂ-ST TO WS-ｺﾏﾝﾄﾞｷﾞｮｳ(9:2)
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           OPEN OUTPUT POLMST-F
           IF NOT WS-ｹｲﾔｸ-OK
               MOVE "POLMST" TO WS-ｺﾏﾝﾄﾞｷﾞｮｳ
               MOVE WS-ｹｲﾔｸ-ST TO WS-ｺﾏﾝﾄﾞｷﾞｮｳ(9:2)
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           OPEN OUTPUT GLSCTL-F
           IF NOT WS-ｾｲｷﾞｮ-OK
               MOVE "GLSCTL" TO WS-ｺﾏﾝﾄﾞｷﾞｮｳ
               MOVE WS-ｾｲｷﾞｮ-ST TO WS-ｺﾏﾝﾄﾞｷﾞｮｳ(9:2)
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           OPEN OUTPUT REPORT-F
           IF NOT WS-ﾎｳｺｸ-OK
               MOVE "GLO010R" TO WS-ｺﾏﾝﾄﾞｷﾞｮｳ
               MOVE WS-ﾎｳｺｸ-ST TO WS-ｺﾏﾝﾄﾞｷﾞｮｳ(9:2)
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           PERFORM ﾁｭｳｼｭﾂ-ﾖﾐｺﾐ-RTN.

      *----------------------------------------------------------------
      * ﾁｭｳｼｭﾂ-ﾖﾐｺﾐ-RTN: 抽出ファイルの次の行。
      *----------------------------------------------------------------
       ﾁｭｳｼｭﾂ-ﾖﾐｺﾐ-RTN.
           READ EXTRACT-F
           IF WS-ﾁｭｳｼｭﾂ-OK
               ADD 1 TO WS-ﾖﾐｺﾐｹﾝｽｳ
           ELSE
               IF NOT WS-ﾁｭｳｼｭﾂ-ｵﾜﾘ
                   MOVE "EXTRACT" TO WS-ｺﾏﾝﾄﾞｷﾞｮｳ
                   MOVE WS-ﾁｭｳｼｭﾂ-ST TO WS-ｺﾏﾝﾄﾞｷﾞｮｳ(9:2)
                   PERFORM ｲｼﾞｮｳ-RTN
               END-IF
           END-IF.

      *----------------------------------------------------------------
      * ﾄｳﾛｸ-RTN: 契約1件をマスタに登録し、状態別に数える。
      *----------------------------------------------------------------
       ﾄｳﾛｸ-RTN.
           MOVE ﾁｭｳｼｭﾂ-ﾚｺｰﾄﾞ TO PM-ｹｲﾔｸ-ﾚｺｰﾄﾞ
           WRITE PM-ｹｲﾔｸ-ﾚｺｰﾄﾞ
               INVALID KEY
                   ADD 1 TO WS-ﾁｮｳﾌｸｹﾝｽｳ
                   MOVE 4 TO WS-ﾌｯｷｺｰﾄﾞ
               NOT INVALID KEY
                   ADD 1 TO WS-ﾄｳﾛｸｹﾝｽｳ
                   PERFORM ｼｭｳｹｲ-RTN
           END-WRITE
           PERFORM ﾁｭｳｼｭﾂ-ﾖﾐｺﾐ-RTN.

      *----------------------------------------------------------------
      * ｼｭｳｹｲ-RTN: 契約を状態別に数える。
      *----------------------------------------------------------------
       ｼｭｳｹｲ-RTN.
           EVALUATE TRUE
               WHEN PM-ﾕｳｺｳ
                   ADD 1 TO WS-ﾕｳｺｳｹﾝｽｳ
                   ADD PM-ﾎｹﾝｷﾝｶﾞｸ TO WS-ﾕｳｺｳﾎｹﾝｷﾝｶﾞｸ
               WHEN PM-ﾕｳﾖﾁｭｳ
                   ADD 1 TO WS-ﾕｳﾖｹﾝｽｳ
                   ADD PM-ﾎｹﾝｷﾝｶﾞｸ TO WS-ﾕｳｺｳﾎｹﾝｷﾝｶﾞｸ
               WHEN PM-ｼｯｺｳ
                   ADD 1 TO WS-ｼｯｺｳｹﾝｽｳ
               WHEN OTHER
                   ADD 1 TO WS-ｼｮｳﾒﾂｹﾝｽｳ
           END-EVALUATE.

      *----------------------------------------------------------------
      * ｾｲｷﾞｮ-RTN: 制御レコードを1回書く。
      *----------------------------------------------------------------
       ｾｲｷﾞｮ-RTN.
           MOVE SPACES TO CT-ｾｲｷﾞｮ-ﾚｺｰﾄﾞ
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ TO CT-ｷﾞｮｳﾑﾋﾂﾞｹ
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO CT-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ
           MOVE 0 TO CT-ｻｲｼｭｳｾｲｷｭｳﾈﾝｹﾞﾂ
           MOVE 1 TO CT-ｼﾞｶｲｼﾊﾗｲﾊﾞﾝｺﾞｳ
           MOVE 1 TO CT-ｼﾞｶｲｼﾜｹﾊﾞﾝｺﾞｳ
           COMPUTE CT-ﾕｳｺｳｹﾝｽｳ = WS-ﾕｳｺｳｹﾝｽｳ + WS-ﾕｳﾖｹﾝｽｳ
           MOVE WS-ﾄｳﾛｸｹﾝｽｳ TO CT-ｿｳｹﾝｽｳ
           MOVE WS-ﾕｳｺｳﾎｹﾝｷﾝｶﾞｸ TO CT-ﾕｳｺｳﾎｹﾝｷﾝｶﾞｸ
           WRITE CT-ｾｲｷﾞｮ-ﾚｺｰﾄﾞ
           IF NOT WS-ｾｲｷﾞｮ-OK
               MOVE "GLSCTL" TO WS-ｺﾏﾝﾄﾞｷﾞｮｳ
               MOVE WS-ｾｲｷﾞｮ-ST TO WS-ｺﾏﾝﾄﾞｷﾞｮｳ(9:2)
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ﾎｳｺｸ-RTN: 登録結果の報告。
      *----------------------------------------------------------------
       ﾎｳｺｸ-RTN.
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO WS-ﾋｮｳﾀﾞｲ-ﾋﾂﾞｹ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ﾋｮｳﾀﾞｲ-ｷﾞｮｳ
           MOVE SPACES TO ﾎｳｺｸ-ｷﾞｮｳ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ
           MOVE "抽出件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾖﾐｺﾐｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "登録件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾄｳﾛｸｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "重複件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾁｮｳﾌｸｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "有効" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾕｳｺｳｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "猶予中" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾕｳﾖｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "失効" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｼｯｺｳｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "消滅" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｼｮｳﾒﾂｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE WS-ﾕｳｺｳﾎｹﾝｷﾝｶﾞｸ TO WS-ｷﾝｶﾞｸ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｷﾝｶﾞｸ-ｷﾞｮｳ.

      *----------------------------------------------------------------
      * ｼｭｳﾘｮｳ-RTN: ファイルを閉じ、件数を表示する。
      *----------------------------------------------------------------
       ｼｭｳﾘｮｳ-RTN.
           CLOSE EXTRACT-F POLMST-F GLSCTL-F REPORT-F
           DISPLAY "GLO010 登録 " WS-ﾄｳﾛｸｹﾝｽｳ
                   " 重複 " WS-ﾁｮｳﾌｸｹﾝｽｳ.

      *----------------------------------------------------------------
      * ｲｼﾞｮｳ-RTN: ファイル異常は 16 で終わる。ファイル名と状態は
      * コマンド行域にある。
      *----------------------------------------------------------------
       ｲｼﾞｮｳ-RTN.
           DISPLAY "GLO010 ファイル異常 " WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8)
                   " 状態 " WS-ｺﾏﾝﾄﾞｷﾞｮｳ(9:2)
           MOVE 16 TO RETURN-CODE
           STOP RUN.
