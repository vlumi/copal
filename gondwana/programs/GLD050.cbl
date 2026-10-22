      *----------------------------------------------------------------
      * GLD050: 口座振替結果処理
      *
      * 振替日の数日後。収納代行から戻った振替結果ファイル PRMRES
      * (銀行間の固定長レイアウト) を読み、契約マスタに反映する。
      * ヘッダが先頭で種別コード 91 であること、データはトレーラ
      * まで表に溜め、トレーラの件数と金額が合えば反映する。合わなけ
      * れば何も反映せず 8 で終わる。
      * 振替済 (結果 0): 払込済年月を1か月 (年払は12か月) 進め、未納
      * 回数を 0 にし、猶予中なら有効に戻し、保険料収入を仕訳する。
      * 振替不能: 未納回数を1つ増やし、有効なら猶予中にして猶予開始
      * 日を業務日付とする。失効は GLD060 が猶予期間の満了で行う。
      * 顧客番号の契約がない、払込中 (有効・猶予中) でない、金額が
      * 保険料の1回分でも未納分を合わせた2回分でもない記録は却下
      * する。
      *
      * 復帰コード: 0 全件反映、4 却下あり、8 ファイル不正、
      *             16 ファイル異常。
      *
      * 変更履歴
      * 日付       担当 変更
      * 2003-10-06 MS   GL-0310 初版、口座振替の開始。
      *----------------------------------------------------------------
       IDENTIFICATION DIVISION.
       PROGRAM-ID. GLD050.
       AUTHOR. ゴンドワナ生命 システム部.
       DATE-WRITTEN. 2003-10-06.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT PRMRES-F ASSIGN TO "PRMRES"
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-PRMRES-ST.
           SELECT POLMST-F ASSIGN TO "POLMST"
               ORGANIZATION IS INDEXED
               ACCESS MODE IS RANDOM
               RECORD KEY IS PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
               FILE STATUS IS WS-POLMST-ST.
           SELECT GLSCTL-F ASSIGN TO "GLSCTL"
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-GLSCTL-ST.
           SELECT OPTIONAL GLJRNL-F ASSIGN TO "GLJRNL"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-GLJRNL-ST.
           SELECT REJECT-F ASSIGN TO "RESREJ"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REJECT-ST.
           SELECT REPORT-F ASSIGN TO "GLD050R"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-ST.

       DATA DIVISION.
       FILE SECTION.
       FD  PRMRES-F
           RECORD CONTAINS 120 CHARACTERS.
       COPY ZENGIN.
       FD  POLMST-F.
       COPY POLMST.
       FD  GLSCTL-F
           RECORD CONTAINS 64 CHARACTERS.
       COPY CTLREC.
       FD  GLJRNL-F.
       COPY GLJRNL.
       FD  REJECT-F.
       01  ｷｬｯｶ-ﾚｺｰﾄﾞ.
           05  ｷｬｯｶ-ﾃﾞｰﾀ               PIC X(120).
           05  ｷｬｯｶ-ﾘﾕｳ                PIC 99.
           05  ｷｬｯｶ-ﾒｯｾｰｼﾞ             PIC X(20).
       FD  REPORT-F.
       01  ﾎｳｺｸ-ｷﾞｮｳ                   PIC X(100).

       WORKING-STORAGE SECTION.
       01  WS-ｽﾃｰﾀｽ.
           05  WS-PRMRES-ST            PIC XX.
               88  WS-PRMRES-OK                VALUE "00".
               88  WS-PRMRES-ｵﾜﾘ               VALUE "10".
           05  WS-POLMST-ST            PIC XX.
               88  WS-POLMST-OK                VALUE "00".
               88  WS-POLMST-ﾅｼ                VALUE "23".
           05  WS-GLSCTL-ST            PIC XX.
               88  WS-GLSCTL-OK                VALUE "00".
           05  WS-GLJRNL-ST            PIC XX.
               88  WS-GLJRNL-OK                VALUE "00" "05".
           05  WS-REJECT-ST            PIC XX.
               88  WS-REJECT-OK                VALUE "00".
           05  WS-REPORT-ST            PIC XX.
               88  WS-REPORT-OK                VALUE "00".
       01  WS-ｺﾏﾝﾄﾞｷﾞｮｳ                PIC X(80).
       01  WS-ｴﾗｰﾌｧｲﾙ                  PIC X(8).
       01  WS-ｴﾗｰｽﾃｰﾀｽ                 PIC XX.
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ                PIC 9(8).
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ            PIC X(7).
       01  WS-ﾎｳｺｳ                     PIC X.
       01  WS-ﾍﾝｶﾝｹｯｶ                  PIC 9.
       01  WS-ﾊﾗｲｺﾐｽﾞﾐ.
           05  WS-ﾊﾗｲｺﾐ-YYYY           PIC 9(4).
           05  WS-ﾊﾗｲｺﾐ-MM             PIC 99.
       01  WS-ｶｲｽｳ                     PIC 9.
      * 振替結果の表: トレーラで検算するまで溜める。
       01  WS-ｹｯｶ-TBL.
           05  WS-ｹｯｶｹﾝｽｳ              PIC 9(5) VALUE 0.
           05  WS-ｹｯｶ OCCURS 5000 TIMES.
               10  WS-ｹｯｶ-ﾚｺｰﾄﾞ            PIC X(120).
       01  WS-ｹｯｶ-IX                   PIC 9(5).
       01  WS-ﾌｧｲﾙ-ｼﾞｮｳﾀｲ              PIC X.
           88  WS-ﾍｯﾀﾞ-ﾏﾁ                      VALUE "H".
           88  WS-ﾃﾞｰﾀ-ﾏﾁ                      VALUE "D".
           88  WS-ｴﾝﾄﾞ-ﾏﾁ                      VALUE "E".
           88  WS-ﾌｧｲﾙ-ｶﾝﾘｮｳ                   VALUE "F".
           88  WS-ﾌｧｲﾙ-ﾌｾｲ                     VALUE "X".
       01  WS-ｹﾝｻ.
           05  WS-ｹﾝｻ-ｹﾝｽｳ             PIC 9(6) VALUE 0.
           05  WS-ｹﾝｻ-ｷﾝｶﾞｸ            PIC 9(12) VALUE 0.
           05  WS-ｹﾝｻ-ｽﾞﾐｹﾝｽｳ          PIC 9(6) VALUE 0.
           05  WS-ｹﾝｻ-ｽﾞﾐｷﾝｶﾞｸ         PIC 9(12) VALUE 0.
           05  WS-ｹﾝｻ-ﾌﾉｳｹﾝｽｳ          PIC 9(6) VALUE 0.
           05  WS-ｹﾝｻ-ﾌﾉｳｷﾝｶﾞｸ         PIC 9(12) VALUE 0.
       01  WS-ｹﾝｽｳ.
           05  WS-ﾊﾝｴｲｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ｽﾞﾐｹﾝｽｳ              PIC 9(7) VALUE 0.
           05  WS-ｽﾞﾐｷﾝｶﾞｸ             PIC 9(13) VALUE 0.
           05  WS-ﾌﾉｳｹﾝｽｳ              PIC 9(7) VALUE 0.
           05  WS-ﾌﾉｳｷﾝｶﾞｸ             PIC 9(13) VALUE 0.
           05  WS-ｷｬｯｶｹﾝｽｳ             PIC 9(7) VALUE 0.
       01  WS-ﾘﾕｳ                      PIC 99.
           88  WS-ﾘﾕｳ-ﾅｼ                       VALUE 0.
           88  WS-ｹｲﾔｸ-ﾅｼ                      VALUE 1.
           88  WS-ﾊﾗｲｺﾐﾁｭｳ-ﾃﾞﾅｲ                VALUE 2.
           88  WS-ｷﾝｶﾞｸ-ｿｳｲ                    VALUE 3.
       01  WS-ﾘﾕｳ-ﾒｯｾｰｼﾞ.
           05  FILLER PIC X(20) VALUE "契約なし".
           05  FILLER PIC X(20) VALUE "払込中でない".
           05  FILLER PIC X(20) VALUE "金額相違".
       01  WS-ﾘﾕｳ-ﾒｯｾｰｼﾞ-TBL REDEFINES WS-ﾘﾕｳ-ﾒｯｾｰｼﾞ.
           05  WS-ﾒｯｾｰｼﾞ OCCURS 3 TIMES PIC X(20).
       01  WS-ﾌｯｷｺｰﾄﾞ                  PIC 99 VALUE 0.
       01  WS-ﾎｳｺｸ.
           05  WS-ﾋｮｳﾀﾞｲ-ｷﾞｮｳ.
               10  FILLER                  PIC X(32)
                   VALUE "ゴンドワナ生命  口座振替結果処理".
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾋｮｳﾀﾞｲ-ﾋﾂﾞｹ          PIC X(7).
           05  WS-ﾒｲｻｲ-ｷﾞｮｳ.
               10  WS-ﾒｲｻｲ-ｼｮｳｹﾝ           PIC X(10).
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾒｲｻｲ-ﾒｲｷﾞ            PIC X(30).
               10  WS-ﾒｲｻｲ-ｷﾝｶﾞｸ           PIC ZZZ,ZZZ,ZZ9.
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾒｲｻｲ-ｹｯｶ             PIC X(20).
           05  WS-ｹﾝｽｳ-ｷﾞｮｳ.
               10  WS-ｹﾝｽｳ-ﾗﾍﾞﾙ            PIC X(16).
               10  WS-ｹﾝｽｳ-ｱﾀｲ             PIC ZZZ,ZZZ,ZZ9.
           05  WS-ﾌｾｲ-ｷﾞｮｳ.
               10  FILLER                  PIC X(24)
                   VALUE "ファイル不正: 反映せず".
               10  WS-ﾌｾｲ-ﾘﾕｳ              PIC X(30).

       PROCEDURE DIVISION.
      *----------------------------------------------------------------
      * ﾒｲﾝ-RTN: 準備、ファイルの読み込みと検算、反映、報告、終了。
      *----------------------------------------------------------------
       ﾒｲﾝ-RTN.
           PERFORM ｼｮｷ-RTN
           PERFORM ﾖﾐｺﾐ-RTN UNTIL WS-ﾌｧｲﾙ-ｶﾝﾘｮｳ OR WS-ﾌｧｲﾙ-ﾌｾｲ
           IF WS-ﾌｧｲﾙ-ｶﾝﾘｮｳ
               PERFORM ﾊﾝｴｲ-RTN VARYING WS-ｹｯｶ-IX FROM 1 BY 1
                   UNTIL WS-ｹｯｶ-IX > WS-ｹｯｶｹﾝｽｳ
           ELSE
               WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ﾌｾｲ-ｷﾞｮｳ
               MOVE 8 TO WS-ﾌｯｷｺｰﾄﾞ
           END-IF
           PERFORM ﾎｳｺｸ-RTN
           PERFORM ｼｭｳﾘｮｳ-RTN
           MOVE WS-ﾌｯｷｺｰﾄﾞ TO RETURN-CODE
           STOP RUN.

      *----------------------------------------------------------------
      * ｼｮｷ-RTN: 業務日付、ファイル、制御レコード。
      *----------------------------------------------------------------
       ｼｮｷ-RTN.
           ACCEPT WS-ｺﾏﾝﾄﾞｷﾞｮｳ FROM COMMAND-LINE
           MOVE WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8) TO WS-ｷﾞｮｳﾑﾋﾂﾞｹ
           MOVE "W" TO WS-ﾎｳｺｳ
           CALL "GLC020" USING WS-ﾎｳｺｳ WS-ｷﾞｮｳﾑﾋﾂﾞｹ
                               WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ WS-ﾍﾝｶﾝｹｯｶ
           IF WS-ﾍﾝｶﾝｹｯｶ NOT = 0
               DISPLAY "GLD050 日付不正: " WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8)
               MOVE 8 TO RETURN-CODE
               STOP RUN
           END-IF
           OPEN INPUT PRMRES-F
           OPEN I-O POLMST-F GLSCTL-F
           OPEN EXTEND GLJRNL-F
           OPEN OUTPUT REJECT-F REPORT-F
           PERFORM ｵｰﾌﾟﾝ-ｹﾝｻ-RTN
           READ GLSCTL-F
           IF NOT WS-GLSCTL-OK
               MOVE "GLSCTL" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-GLSCTL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO WS-ﾋｮｳﾀﾞｲ-ﾋﾂﾞｹ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ﾋｮｳﾀﾞｲ-ｷﾞｮｳ
           MOVE SPACES TO ﾎｳｺｸ-ｷﾞｮｳ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ
           SET WS-ﾍｯﾀﾞ-ﾏﾁ TO TRUE.

      *----------------------------------------------------------------
      * ｵｰﾌﾟﾝ-ｹﾝｻ-RTN: 開けなかったファイルがあれば異常終了。
      *----------------------------------------------------------------
       ｵｰﾌﾟﾝ-ｹﾝｻ-RTN.
           EVALUATE TRUE
               WHEN NOT WS-PRMRES-OK
                   MOVE "PRMRES" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-PRMRES-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-POLMST-OK
                   MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-GLSCTL-OK
                   MOVE "GLSCTL" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-GLSCTL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-GLJRNL-OK
                   MOVE "GLJRNL" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-GLJRNL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-REJECT-OK
                   MOVE "RESREJ" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-REJECT-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-REPORT-OK
                   MOVE "GLD050R" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-REPORT-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN OTHER
                   MOVE SPACES TO WS-ｴﾗｰﾌｧｲﾙ
           END-EVALUATE
           IF WS-ｴﾗｰﾌｧｲﾙ NOT = SPACES
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ﾖﾐｺﾐ-RTN: 次の記録を読み、順序とデータ区分で振り分ける。
      * ヘッダ、データ、トレーラ、エンドの順でなければファイル不正。
      *----------------------------------------------------------------
       ﾖﾐｺﾐ-RTN.
           READ PRMRES-F
           EVALUATE TRUE
               WHEN WS-PRMRES-ｵﾜﾘ
                   MOVE "エンド記録なし" TO WS-ﾌｾｲ-ﾘﾕｳ
                   SET WS-ﾌｧｲﾙ-ﾌｾｲ TO TRUE
               WHEN NOT WS-PRMRES-OK
                   MOVE "PRMRES" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-PRMRES-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
                   PERFORM ｲｼﾞｮｳ-RTN
               WHEN WS-ﾍｯﾀﾞ-ﾏﾁ AND ZG-ﾍｯﾀﾞ AND ZG-H-ｼｭﾍﾞﾂｺｰﾄﾞ = 91
                   SET WS-ﾃﾞｰﾀ-ﾏﾁ TO TRUE
               WHEN WS-ﾍｯﾀﾞ-ﾏﾁ
                   MOVE "ヘッダ不正" TO WS-ﾌｾｲ-ﾘﾕｳ
                   SET WS-ﾌｧｲﾙ-ﾌｾｲ TO TRUE
               WHEN WS-ﾃﾞｰﾀ-ﾏﾁ AND ZG-ﾃﾞｰﾀ
                   PERFORM ﾀﾒﾙ-RTN
               WHEN WS-ﾃﾞｰﾀ-ﾏﾁ AND ZG-ﾄﾚｰﾗ
                   PERFORM ｹﾝｻﾝ-RTN
               WHEN WS-ｴﾝﾄﾞ-ﾏﾁ AND ZG-ｴﾝﾄﾞ
                   SET WS-ﾌｧｲﾙ-ｶﾝﾘｮｳ TO TRUE
               WHEN OTHER
                   MOVE "記録の順序不正" TO WS-ﾌｾｲ-ﾘﾕｳ
                   SET WS-ﾌｧｲﾙ-ﾌｾｲ TO TRUE
           END-EVALUATE.

      *----------------------------------------------------------------
      * ﾀﾒﾙ-RTN: データ記録を表に溜め、検算の件数と金額を足す。
      *----------------------------------------------------------------
       ﾀﾒﾙ-RTN.
           ADD 1 TO WS-ｹｯｶｹﾝｽｳ
           MOVE ZG-ﾌﾘｶｴ-ﾚｺｰﾄﾞ TO WS-ｹｯｶ-ﾚｺｰﾄﾞ(WS-ｹｯｶｹﾝｽｳ)
           ADD 1 TO WS-ｹﾝｻ-ｹﾝｽｳ
           ADD ZG-D-ｷﾝｶﾞｸ TO WS-ｹﾝｻ-ｷﾝｶﾞｸ
           IF ZG-D-ﾌﾘｶｴｽﾞﾐ
               ADD 1 TO WS-ｹﾝｻ-ｽﾞﾐｹﾝｽｳ
               ADD ZG-D-ｷﾝｶﾞｸ TO WS-ｹﾝｻ-ｽﾞﾐｷﾝｶﾞｸ
           ELSE
               ADD 1 TO WS-ｹﾝｻ-ﾌﾉｳｹﾝｽｳ
               ADD ZG-D-ｷﾝｶﾞｸ TO WS-ｹﾝｻ-ﾌﾉｳｷﾝｶﾞｸ
           END-IF.

      *----------------------------------------------------------------
      * ｹﾝｻﾝ-RTN: トレーラの件数と金額を溜めたものと突き合わせる。
      *----------------------------------------------------------------
       ｹﾝｻﾝ-RTN.
           EVALUATE TRUE
               WHEN ZG-T-ｺﾞｳｹｲｹﾝｽｳ NOT = WS-ｹﾝｻ-ｹﾝｽｳ
                 OR ZG-T-ｺﾞｳｹｲｷﾝｶﾞｸ NOT = WS-ｹﾝｻ-ｷﾝｶﾞｸ
                   MOVE "トレーラの合計不一致" TO WS-ﾌｾｲ-ﾘﾕｳ
                   SET WS-ﾌｧｲﾙ-ﾌｾｲ TO TRUE
               WHEN ZG-T-ﾌﾘｶｴｽﾞﾐｹﾝｽｳ NOT = WS-ｹﾝｻ-ｽﾞﾐｹﾝｽｳ
                 OR ZG-T-ﾌﾘｶｴｽﾞﾐｷﾝｶﾞｸ NOT = WS-ｹﾝｻ-ｽﾞﾐｷﾝｶﾞｸ
                   MOVE "トレーラの振替済不一致" TO WS-ﾌｾｲ-ﾘﾕｳ
                   SET WS-ﾌｧｲﾙ-ﾌｾｲ TO TRUE
               WHEN ZG-T-ﾌﾘｶｴﾌﾉｳｹﾝｽｳ NOT = WS-ｹﾝｻ-ﾌﾉｳｹﾝｽｳ
                 OR ZG-T-ﾌﾘｶｴﾌﾉｳｷﾝｶﾞｸ NOT = WS-ｹﾝｻ-ﾌﾉｳｷﾝｶﾞｸ
                   MOVE "トレーラの振替不能不一致" TO WS-ﾌｾｲ-ﾘﾕｳ
                   SET WS-ﾌｧｲﾙ-ﾌｾｲ TO TRUE
               WHEN OTHER
                   SET WS-ｴﾝﾄﾞ-ﾏﾁ TO TRUE
           END-EVALUATE.

      *----------------------------------------------------------------
      * ﾊﾝｴｲ-RTN: 溜めた1件を契約に反映する。
      *----------------------------------------------------------------
       ﾊﾝｴｲ-RTN.
           MOVE WS-ｹｯｶ-ﾚｺｰﾄﾞ(WS-ｹｯｶ-IX) TO ZG-ﾌﾘｶｴ-ﾚｺｰﾄﾞ
           SET WS-ﾘﾕｳ-ﾅｼ TO TRUE
           MOVE ZG-D-ｺｷｬｸﾊﾞﾝｺﾞｳ(1:10) TO PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           READ POLMST-F
           EVALUATE TRUE
               WHEN WS-POLMST-ﾅｼ
                   SET WS-ｹｲﾔｸ-ﾅｼ TO TRUE
               WHEN NOT WS-POLMST-OK
                   MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
                   PERFORM ｲｼﾞｮｳ-RTN
               WHEN NOT PM-ﾊﾗｲｺﾐﾁｭｳ
                   SET WS-ﾊﾗｲｺﾐﾁｭｳ-ﾃﾞﾅｲ TO TRUE
               WHEN OTHER
                   PERFORM ｷﾝｶﾞｸ-ｹﾝｻ-RTN
           END-EVALUATE
           IF WS-ﾘﾕｳ-ﾅｼ
               IF ZG-D-ﾌﾘｶｴｽﾞﾐ
                   PERFORM ﾌﾘｶｴｽﾞﾐ-RTN
               ELSE
                   PERFORM ﾌﾘｶｴﾌﾉｳ-RTN
               END-IF
               REWRITE PM-ｹｲﾔｸ-ﾚｺｰﾄﾞ
               IF NOT WS-POLMST-OK
                   MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
                   PERFORM ｲｼﾞｮｳ-RTN
               END-IF
               ADD 1 TO WS-ﾊﾝｴｲｹﾝｽｳ
           ELSE
               PERFORM ｷｬｯｶ-RTN
           END-IF
           PERFORM ﾒｲｻｲ-RTN.

      *----------------------------------------------------------------
      * ｷﾝｶﾞｸ-ｹﾝｻ-RTN: 振替額は保険料の1回分か、未納分を合わせた
      * 2回分 (年払は1回分のみ)。それ以外は金額相違。
      *----------------------------------------------------------------
       ｷﾝｶﾞｸ-ｹﾝｻ-RTN.
           EVALUATE TRUE
               WHEN ZG-D-ｷﾝｶﾞｸ = PM-ﾎｹﾝﾘｮｳ
                   MOVE 1 TO WS-ｶｲｽｳ
               WHEN PM-ﾂｷﾊﾞﾗｲ AND ZG-D-ｷﾝｶﾞｸ = PM-ﾎｹﾝﾘｮｳ * 2
                   MOVE 2 TO WS-ｶｲｽｳ
               WHEN OTHER
                   SET WS-ｷﾝｶﾞｸ-ｿｳｲ TO TRUE
           END-EVALUATE.

      *----------------------------------------------------------------
      * ﾌﾘｶｴｽﾞﾐ-RTN: 払込済年月を回数分進め、猶予中なら有効に戻し、
      * 仕訳。
      *----------------------------------------------------------------
       ﾌﾘｶｴｽﾞﾐ-RTN.
           MOVE PM-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ TO WS-ﾊﾗｲｺﾐｽﾞﾐ
           IF PM-ﾈﾝﾊﾞﾗｲ
               ADD 1 TO WS-ﾊﾗｲｺﾐ-YYYY
           ELSE
               ADD WS-ｶｲｽｳ TO WS-ﾊﾗｲｺﾐ-MM
               IF WS-ﾊﾗｲｺﾐ-MM > 12
                   SUBTRACT 12 FROM WS-ﾊﾗｲｺﾐ-MM
                   ADD 1 TO WS-ﾊﾗｲｺﾐ-YYYY
               END-IF
           END-IF
           MOVE WS-ﾊﾗｲｺﾐｽﾞﾐ TO PM-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ
           MOVE 0 TO PM-ﾐﾉｳｶｲｽｳ
           IF PM-ﾕｳﾖﾁｭｳ
               SET PM-ﾕｳｺｳ TO TRUE
               MOVE SPACES TO PM-ﾕｳﾖｶｲｼﾋﾞ
               MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO PM-ｼﾞｮｳﾀｲﾋﾂﾞｹ
           END-IF
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ(1:6) TO PM-ｻｲｼｭｳﾆｭｳｷﾝﾈﾝｹﾞﾂ
           MOVE ZG-D-ｷﾝｶﾞｸ TO PM-ｻｲｼｭｳﾆｭｳｷﾝｶﾞｸ
           MOVE SPACES TO JR-ｼﾜｹ-ﾚｺｰﾄﾞ
           MOVE CT-ｼﾞｶｲｼﾜｹﾊﾞﾝｺﾞｳ TO JR-ｼﾜｹﾊﾞﾝｺﾞｳ
           ADD 1 TO CT-ｼﾞｶｲｼﾜｹﾊﾞﾝｺﾞｳ
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO JR-ﾋﾂﾞｹ
           MOVE "GLD050" TO JR-ﾌﾟﾛｸﾞﾗﾑ
           MOVE 4110 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
           SET JR-ｶｼｶﾀ TO TRUE
           MOVE ZG-D-ｷﾝｶﾞｸ TO JR-ｷﾝｶﾞｸ
           MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO JR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           MOVE "口座振替保険料" TO JR-ﾃｷﾖｳ
           WRITE JR-ｼﾜｹ-ﾚｺｰﾄﾞ
           ADD 1 TO WS-ｽﾞﾐｹﾝｽｳ
           ADD ZG-D-ｷﾝｶﾞｸ TO WS-ｽﾞﾐｷﾝｶﾞｸ
           MOVE "振替済" TO WS-ﾒｲｻｲ-ｹｯｶ.

      *----------------------------------------------------------------
      * ﾌﾘｶｴﾌﾉｳ-RTN: 未納回数を増やし、有効なら猶予中にする。
      *----------------------------------------------------------------
       ﾌﾘｶｴﾌﾉｳ-RTN.
           ADD 1 TO PM-ﾐﾉｳｶｲｽｳ
           IF PM-ﾕｳｺｳ
               SET PM-ﾕｳﾖﾁｭｳ TO TRUE
               MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO PM-ﾕｳﾖｶｲｼﾋﾞ
               MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO PM-ｼﾞｮｳﾀｲﾋﾂﾞｹ
           END-IF
           ADD 1 TO WS-ﾌﾉｳｹﾝｽｳ
           ADD ZG-D-ｷﾝｶﾞｸ TO WS-ﾌﾉｳｷﾝｶﾞｸ
           EVALUATE TRUE
               WHEN ZG-D-ｼｷﾝﾌｿｸ
                   MOVE "振替不能: 資金不足" TO WS-ﾒｲｻｲ-ｹｯｶ
               WHEN ZG-D-ﾄﾘﾋｷﾅｼ
                   MOVE "振替不能: 取引なし" TO WS-ﾒｲｻｲ-ｹｯｶ
               WHEN ZG-D-ﾖｷﾝｼｬﾂｺﾞｳ
                   MOVE "振替不能: 預金者都合" TO WS-ﾒｲｻｲ-ｹｯｶ
               WHEN ZG-D-ｲﾗｲｼｮﾅｼ
                   MOVE "振替不能: 依頼書なし" TO WS-ﾒｲｻｲ-ｹｯｶ
               WHEN OTHER
                   MOVE "振替不能: その他" TO WS-ﾒｲｻｲ-ｹｯｶ
           END-EVALUATE.

      *----------------------------------------------------------------
      * ｷｬｯｶ-RTN: 却下ファイルに理由を付けて書く。
      *----------------------------------------------------------------
       ｷｬｯｶ-RTN.
           MOVE ZG-ﾌﾘｶｴ-ﾚｺｰﾄﾞ TO ｷｬｯｶ-ﾃﾞｰﾀ
           MOVE WS-ﾘﾕｳ TO ｷｬｯｶ-ﾘﾕｳ
           MOVE WS-ﾒｯｾｰｼﾞ(WS-ﾘﾕｳ) TO ｷｬｯｶ-ﾒｯｾｰｼﾞ
           WRITE ｷｬｯｶ-ﾚｺｰﾄﾞ
           ADD 1 TO WS-ｷｬｯｶｹﾝｽｳ
           MOVE 4 TO WS-ﾌｯｷｺｰﾄﾞ
           MOVE WS-ﾒｯｾｰｼﾞ(WS-ﾘﾕｳ) TO WS-ﾒｲｻｲ-ｹｯｶ.

      *----------------------------------------------------------------
      * ﾒｲｻｲ-RTN: 1件の明細行。
      *----------------------------------------------------------------
       ﾒｲｻｲ-RTN.
           MOVE ZG-D-ｺｷｬｸﾊﾞﾝｺﾞｳ(1:10) TO WS-ﾒｲｻｲ-ｼｮｳｹﾝ
           MOVE ZG-D-ｺｳｻﾞﾒｲｷﾞ TO WS-ﾒｲｻｲ-ﾒｲｷﾞ
           MOVE ZG-D-ｷﾝｶﾞｸ TO WS-ﾒｲｻｲ-ｷﾝｶﾞｸ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ﾒｲｻｲ-ｷﾞｮｳ.

      *----------------------------------------------------------------
      * ﾎｳｺｸ-RTN: 件数と合計。
      *----------------------------------------------------------------
       ﾎｳｺｸ-RTN.
           MOVE SPACES TO ﾎｳｺｸ-ｷﾞｮｳ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ
           MOVE "読込件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｹｯｶｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "振替済件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｽﾞﾐｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "振替済金額" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｽﾞﾐｷﾝｶﾞｸ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "振替不能件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾌﾉｳｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "振替不能金額" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾌﾉｳｷﾝｶﾞｸ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "却下件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｷｬｯｶｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ.

      *----------------------------------------------------------------
      * ｼｭｳﾘｮｳ-RTN: 制御レコードを書き戻し、ファイルを閉じる。
      *----------------------------------------------------------------
       ｼｭｳﾘｮｳ-RTN.
           REWRITE CT-ｾｲｷﾞｮ-ﾚｺｰﾄﾞ
           IF NOT WS-GLSCTL-OK
               MOVE "GLSCTL" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-GLSCTL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           CLOSE PRMRES-F POLMST-F GLSCTL-F GLJRNL-F REJECT-F
                 REPORT-F
           DISPLAY "GLD050 振替済 " WS-ｽﾞﾐｹﾝｽｳ " 不能 " WS-ﾌﾉｳｹﾝｽｳ
                   " 却下 " WS-ｷｬｯｶｹﾝｽｳ.

      *----------------------------------------------------------------
      * ｲｼﾞｮｳ-RTN: ファイル異常は 16 で終わる。
      *----------------------------------------------------------------
       ｲｼﾞｮｳ-RTN.
           DISPLAY "GLD050 ファイル異常 " WS-ｴﾗｰﾌｧｲﾙ
                   " 状態 " WS-ｴﾗｰｽﾃｰﾀｽ
           MOVE 16 TO RETURN-CODE
           STOP RUN.
