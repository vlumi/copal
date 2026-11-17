      *----------------------------------------------------------------
      * GLD070: 契約者貸付
      *
      * 日次。支社からの契約者貸付依頼ファイル LONREQ を読み、貸付
      * と返済を契約マスタの貸付残高に反映する。
      * どちらも、まず残高に業務日付までの利息 (GLC040) を繰り入れ
      * る。
      * 貸付: 契約が有効で、貸付後の残高が解約返戻金 (GLC030) の
      * 8割以内であること。貸付金を契約者の口座に支払う。
      * 返済: 契約が保障中 (有効・猶予中) で、返済額が残高以内で
      * あること。
      * 残高が 0 になれば利息起算日を空白にし、残れば業務日付と
      * する。支払番号と仕訳番号は制御レコードから採番する。
      *
      * 復帰コード: 0 全件反映、4 却下あり、8 日付が日付でない、
      *             16 ファイル異常。
      *
      * 変更履歴
      * 日付       担当 変更
      * 2012-09-03 YN   GL-0561 初版、支社の貸付台帳をシステム化。
      *----------------------------------------------------------------
       IDENTIFICATION DIVISION.
       PROGRAM-ID. GLD070.
       AUTHOR. ゴンドワナ生命 システム部.
       DATE-WRITTEN. 2012-09-03.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT LONREQ-F ASSIGN TO "LONREQ"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-LONREQ-ST.
           SELECT POLMST-F ASSIGN TO "POLMST"
               ORGANIZATION IS INDEXED
               ACCESS MODE IS RANDOM
               RECORD KEY IS PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
               FILE STATUS IS WS-POLMST-ST.
           SELECT GLSCTL-F ASSIGN TO "GLSCTL"
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-GLSCTL-ST.
           SELECT OPTIONAL PAYOUT-F ASSIGN TO "PAYOUT"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-PAYOUT-ST.
           SELECT OPTIONAL GLJRNL-F ASSIGN TO "GLJRNL"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-GLJRNL-ST.
           SELECT REJECT-F ASSIGN TO "LONREJ"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REJECT-ST.
           SELECT REPORT-F ASSIGN TO "GLD070R"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-ST.

       DATA DIVISION.
       FILE SECTION.
       FD  LONREQ-F.
       COPY LONREQ.
       FD  POLMST-F.
       COPY POLMST.
       FD  GLSCTL-F
           RECORD CONTAINS 64 CHARACTERS.
       COPY CTLREC.
       FD  PAYOUT-F.
       COPY PAYOUT.
       FD  GLJRNL-F.
       COPY GLJRNL.
       FD  REJECT-F.
       01  ｷｬｯｶ-ﾚｺｰﾄﾞ.
           05  ｷｬｯｶ-ﾃﾞｰﾀ               PIC X(60).
           05  ｷｬｯｶ-ﾘﾕｳ                PIC 99.
           05  ｷｬｯｶ-ﾒｯｾｰｼﾞ             PIC X(20).
       FD  REPORT-F.
       01  ﾎｳｺｸ-ｷﾞｮｳ                   PIC X(100).

       WORKING-STORAGE SECTION.
       01  WS-ｽﾃｰﾀｽ.
           05  WS-LONREQ-ST            PIC XX.
               88  WS-LONREQ-OK                VALUE "00".
               88  WS-LONREQ-ｵﾜﾘ               VALUE "10".
           05  WS-POLMST-ST            PIC XX.
               88  WS-POLMST-OK                VALUE "00".
               88  WS-POLMST-ﾅｼ                VALUE "23".
           05  WS-GLSCTL-ST            PIC XX.
               88  WS-GLSCTL-OK                VALUE "00".
           05  WS-PAYOUT-ST            PIC XX.
               88  WS-PAYOUT-OK                VALUE "00" "05".
           05  WS-GLJRNL-ST            PIC XX.
               88  WS-GLJRNL-OK                VALUE "00" "05".
           05  WS-REJECT-ST            PIC XX.
               88  WS-REJECT-OK                VALUE "00".
           05  WS-REPORT-ST            PIC XX.
               88  WS-REPORT-OK                VALUE "00".
           05  WS-CSVTB-ST             PIC XX.
               88  WS-CSVTB-OK                 VALUE "00".
       01  WS-ｺﾏﾝﾄﾞｷﾞｮｳ                PIC X(80).
       01  WS-ｴﾗｰﾌｧｲﾙ                  PIC X(8).
       01  WS-ｴﾗｰｽﾃｰﾀｽ                 PIC XX.
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ                PIC 9(8).
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ            PIC X(7).
       01  WS-ﾎｳｺｳ                     PIC X.
       01  WS-ﾍﾝｶﾝｹｯｶ                  PIC 9.
       01  WS-ｹﾞﾝﾄﾞﾘﾂ                  PIC V99 VALUE 0.80.
       01  WS-ﾍﾝﾚｲｷﾝ                   PIC 9(9).
       01  WS-ｹﾞﾝﾄﾞｶﾞｸ                 PIC 9(9).
       01  WS-ﾘｿｸ                      PIC 9(9).
       01  WS-ｻﾞﾝﾀﾞｶ                   PIC 9(10).
       01  WS-ｹﾝｽｳ.
           05  WS-ﾖﾐｺﾐｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ｶｼﾂｹｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ｶｼﾂｹｷﾝｶﾞｸ            PIC 9(13) VALUE 0.
           05  WS-ﾍﾝｻｲｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ﾍﾝｻｲｷﾝｶﾞｸ            PIC 9(13) VALUE 0.
           05  WS-ﾘｿｸｷﾝｶﾞｸ             PIC 9(13) VALUE 0.
           05  WS-ｷｬｯｶｹﾝｽｳ             PIC 9(7) VALUE 0.
       01  WS-ﾘﾕｳ                      PIC 99.
           88  WS-ﾘﾕｳ-ﾅｼ                       VALUE 0.
           88  WS-ｹｲﾔｸ-ﾅｼ                      VALUE 1.
           88  WS-ｼｭﾍﾞﾂ-ﾌｾｲ                    VALUE 2.
           88  WS-ｷﾝｶﾞｸ-ﾌｾｲ                    VALUE 3.
           88  WS-ｼﾞｮｳﾀｲ-ﾌｶ                    VALUE 4.
           88  WS-ｹﾞﾝﾄﾞ-ﾁｮｳｶ                   VALUE 5.
           88  WS-ﾍﾝｻｲ-ﾁｮｳｶ                    VALUE 6.
       01  WS-ﾘﾕｳ-ﾒｯｾｰｼﾞ.
           05  FILLER PIC X(20) VALUE "契約なし".
           05  FILLER PIC X(20) VALUE "種別不正".
           05  FILLER PIC X(20) VALUE "金額不正".
           05  FILLER PIC X(20) VALUE "対象外の状態".
           05  FILLER PIC X(20) VALUE "貸付限度超過".
           05  FILLER PIC X(20) VALUE "返済額超過".
       01  WS-ﾘﾕｳ-ﾒｯｾｰｼﾞ-TBL REDEFINES WS-ﾘﾕｳ-ﾒｯｾｰｼﾞ.
           05  WS-ﾒｯｾｰｼﾞ OCCURS 6 TIMES PIC X(20).
       01  WS-ｼｭﾍﾞﾂﾒｲ                  PIC X(6).
       01  WS-ﾌｯｷｺｰﾄﾞ                  PIC 99 VALUE 0.
       01  WS-ﾎｳｺｸ.
           05  WS-ﾋｮｳﾀﾞｲ-ｷﾞｮｳ.
               10  FILLER                  PIC X(28)
                   VALUE "ゴンドワナ生命  契約者貸付".
               10  FILLER                  PIC X(4) VALUE SPACES.
               10  WS-ﾋｮｳﾀﾞｲ-ﾋﾂﾞｹ          PIC X(7).
           05  WS-ﾒｲｻｲ-ｷﾞｮｳ.
               10  WS-ﾒｲｻｲ-ｲﾗｲ             PIC X(8).
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾒｲｻｲ-ｼｮｳｹﾝ           PIC X(10).
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾒｲｻｲ-ｼｭﾍﾞﾂ           PIC X(6).
               10  WS-ﾒｲｻｲ-ｷﾝｶﾞｸ           PIC ZZZ,ZZZ,ZZ9.
               10  FILLER                  PIC X(6) VALUE "  残高".
               10  WS-ﾒｲｻｲ-ｻﾞﾝﾀﾞｶ          PIC ZZZ,ZZZ,ZZ9.
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾒｲｻｲ-ｹｯｶ             PIC X(20).
           05  WS-ｹﾝｽｳ-ｷﾞｮｳ.
               10  WS-ｹﾝｽｳ-ﾗﾍﾞﾙ            PIC X(16).
               10  WS-ｹﾝｽｳ-ｱﾀｲ             PIC ZZZ,ZZZ,ZZ9.

       PROCEDURE DIVISION.
      *----------------------------------------------------------------
      * ﾒｲﾝ-RTN: 準備、1件ずつの反映、報告、終了。
      *----------------------------------------------------------------
       ﾒｲﾝ-RTN.
           PERFORM ｼｮｷ-RTN
           PERFORM ｼｮﾘ-RTN UNTIL WS-LONREQ-ｵﾜﾘ
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
               DISPLAY "GLD070 日付不正: " WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8)
               MOVE 8 TO RETURN-CODE
               STOP RUN
           END-IF
           OPEN INPUT LONREQ-F
           OPEN I-O POLMST-F GLSCTL-F
           OPEN EXTEND PAYOUT-F GLJRNL-F
           OPEN OUTPUT REJECT-F REPORT-F
           PERFORM ｵｰﾌﾟﾝ-ｹﾝｻ-RTN
           MOVE SPACES TO PY-ｼﾊﾗｲ-ﾚｺｰﾄﾞ JR-ｼﾜｹ-ﾚｺｰﾄﾞ
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
           PERFORM LONREQ-ﾖﾐｺﾐ-RTN.

      *----------------------------------------------------------------
      * ｵｰﾌﾟﾝ-ｹﾝｻ-RTN: 開けなかったファイルがあれば異常終了。
      *----------------------------------------------------------------
       ｵｰﾌﾟﾝ-ｹﾝｻ-RTN.
           EVALUATE TRUE
               WHEN NOT WS-LONREQ-OK
                   MOVE "LONREQ" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-LONREQ-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-POLMST-OK
                   MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-GLSCTL-OK
                   MOVE "GLSCTL" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-GLSCTL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-PAYOUT-OK
                   MOVE "PAYOUT" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-PAYOUT-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-GLJRNL-OK
                   MOVE "GLJRNL" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-GLJRNL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-REJECT-OK
                   MOVE "LONREJ" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-REJECT-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-REPORT-OK
                   MOVE "GLD070R" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-REPORT-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN OTHER
                   MOVE SPACES TO WS-ｴﾗｰﾌｧｲﾙ
           END-EVALUATE
           IF WS-ｴﾗｰﾌｧｲﾙ NOT = SPACES
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * LONREQ-ﾖﾐｺﾐ-RTN: 依頼ファイルの次の行。
      *----------------------------------------------------------------
       LONREQ-ﾖﾐｺﾐ-RTN.
           READ LONREQ-F
           IF WS-LONREQ-OK
               ADD 1 TO WS-ﾖﾐｺﾐｹﾝｽｳ
           ELSE
               IF NOT WS-LONREQ-ｵﾜﾘ
                   MOVE "LONREQ" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-LONREQ-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
                   PERFORM ｲｼﾞｮｳ-RTN
               END-IF
           END-IF.

      *----------------------------------------------------------------
      * ｼｮﾘ-RTN: 契約を読み、検査して反映か却下。
      *----------------------------------------------------------------
       ｼｮﾘ-RTN.
           SET WS-ﾘﾕｳ-ﾅｼ TO TRUE
           MOVE 0 TO WS-ｻﾞﾝﾀﾞｶ
           MOVE LR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           READ POLMST-F
           IF WS-POLMST-ﾅｼ
               SET WS-ｹｲﾔｸ-ﾅｼ TO TRUE
           ELSE
               IF NOT WS-POLMST-OK
                   MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
                   PERFORM ｲｼﾞｮｳ-RTN
               END-IF
           END-IF
           EVALUATE TRUE
               WHEN LR-ｶｼﾂｹ
                   MOVE "貸付" TO WS-ｼｭﾍﾞﾂﾒｲ
               WHEN LR-ﾍﾝｻｲ
                   MOVE "返済" TO WS-ｼｭﾍﾞﾂﾒｲ
               WHEN OTHER
                   MOVE "不明" TO WS-ｼｭﾍﾞﾂﾒｲ
           END-EVALUATE
           IF WS-ﾘﾕｳ-ﾅｼ
               PERFORM ｹﾝｻ-RTN
           END-IF
           IF WS-ﾘﾕｳ-ﾅｼ
               PERFORM ﾊﾝｴｲ-RTN
           ELSE
               PERFORM ｷｬｯｶ-RTN
           END-IF
           PERFORM LONREQ-ﾖﾐｺﾐ-RTN.

      *----------------------------------------------------------------
      * ｹﾝｻ-RTN: 種別、金額、契約の状態、次いで利息を繰り入れた残高
      * に対する限度。
      *----------------------------------------------------------------
       ｹﾝｻ-RTN.
           EVALUATE TRUE
               WHEN NOT LR-ｶｼﾂｹ AND NOT LR-ﾍﾝｻｲ
                   SET WS-ｼｭﾍﾞﾂ-ﾌｾｲ TO TRUE
               WHEN LR-ｷﾝｶﾞｸ NOT NUMERIC OR LR-ｷﾝｶﾞｸ = 0
                   SET WS-ｷﾝｶﾞｸ-ﾌｾｲ TO TRUE
               WHEN LR-ｶｼﾂｹ AND NOT PM-ﾕｳｺｳ
                   SET WS-ｼﾞｮｳﾀｲ-ﾌｶ TO TRUE
               WHEN LR-ﾍﾝｻｲ AND NOT PM-ﾎｼｮｳﾁｭｳ
                   SET WS-ｼﾞｮｳﾀｲ-ﾌｶ TO TRUE
               WHEN OTHER
                   PERFORM ｹﾞﾝﾄﾞ-RTN
           END-EVALUATE.

      *----------------------------------------------------------------
      * ｹﾞﾝﾄﾞ-RTN: 業務日付までの利息を加えた残高を求め、貸付は
      * 解約返戻金の8割、返済はその残高を限度とする。
      *----------------------------------------------------------------
       ｹﾞﾝﾄﾞ-RTN.
           CALL "GLC040" USING PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ PM-ｶｼﾂｹｷｻﾝﾋﾞ
                               WS-ｷﾞｮｳﾑﾋﾂﾞｹ WS-ﾘｿｸ
           COMPUTE WS-ｻﾞﾝﾀﾞｶ = PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ + WS-ﾘｿｸ
           IF LR-ｶｼﾂｹ
               CALL "GLC030" USING PM-ｼｮｳﾋﾝ PM-ｶﾆｭｳﾈﾝﾚｲ PM-ｹｲｶﾈﾝｽｳ
                                   PM-ﾎｹﾝｷﾝｶﾞｸ WS-ﾍﾝﾚｲｷﾝ WS-CSVTB-ST
               IF NOT WS-CSVTB-OK
                   MOVE "CSVTB" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-CSVTB-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
                   PERFORM ｲｼﾞｮｳ-RTN
               END-IF
               COMPUTE WS-ｹﾞﾝﾄﾞｶﾞｸ = WS-ﾍﾝﾚｲｷﾝ * WS-ｹﾞﾝﾄﾞﾘﾂ
               IF WS-ｻﾞﾝﾀﾞｶ + LR-ｷﾝｶﾞｸ > WS-ｹﾞﾝﾄﾞｶﾞｸ
                   SET WS-ｹﾞﾝﾄﾞ-ﾁｮｳｶ TO TRUE
               END-IF
           ELSE
               IF LR-ｷﾝｶﾞｸ > WS-ｻﾞﾝﾀﾞｶ
                   SET WS-ﾍﾝｻｲ-ﾁｮｳｶ TO TRUE
               END-IF
           END-IF.

      *----------------------------------------------------------------
      * ﾊﾝｴｲ-RTN: 利息を繰り入れ、貸付か返済を残高に反映して書き
      * 戻す。
      *----------------------------------------------------------------
       ﾊﾝｴｲ-RTN.
           IF WS-ﾘｿｸ > 0
               MOVE 4310 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
               SET JR-ｶｼｶﾀ TO TRUE
               MOVE WS-ﾘｿｸ TO JR-ｷﾝｶﾞｸ
               MOVE "貸付金利息" TO JR-ﾃｷﾖｳ
               PERFORM ｼﾜｹ-RTN
               ADD WS-ﾘｿｸ TO WS-ﾘｿｸｷﾝｶﾞｸ
           END-IF
           IF LR-ｶｼﾂｹ
               PERFORM ｶｼﾂｹ-RTN
           ELSE
               PERFORM ﾍﾝｻｲ-RTN
           END-IF
           MOVE WS-ｻﾞﾝﾀﾞｶ TO PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ
           IF PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0
               MOVE SPACES TO PM-ｶｼﾂｹｷｻﾝﾋﾞ
           ELSE
               MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO PM-ｶｼﾂｹｷｻﾝﾋﾞ
           END-IF
           REWRITE PM-ｹｲﾔｸ-ﾚｺｰﾄﾞ
           IF NOT WS-POLMST-OK
               MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           MOVE "反映" TO WS-ﾒｲｻｲ-ｹｯｶ
           PERFORM ﾒｲｻｲ-RTN.

      *----------------------------------------------------------------
      * ｶｼﾂｹ-RTN: 貸付金を残高に加え、契約者の口座への支払を書く。
      *----------------------------------------------------------------
       ｶｼﾂｹ-RTN.
           ADD LR-ｷﾝｶﾞｸ TO WS-ｻﾞﾝﾀﾞｶ
           MOVE CT-ｼﾞｶｲｼﾊﾗｲﾊﾞﾝｺﾞｳ TO PY-ｼﾊﾗｲﾊﾞﾝｺﾞｳ
           ADD 1 TO CT-ｼﾞｶｲｼﾊﾗｲﾊﾞﾝｺﾞｳ
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO PY-ｼﾊﾗｲﾋﾞ
           MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO PY-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           SET PY-ｶｼﾂｹ TO TRUE
           MOVE LR-ｷﾝｶﾞｸ TO PY-ｷﾝｶﾞｸ
           MOVE PM-ｺｳｻﾞﾒｲｷﾞ TO PY-ｳｹﾄﾘﾆﾝｶﾅ
           MOVE LR-ｲﾗｲﾊﾞﾝｺﾞｳ TO PY-ｾｲｷｭｳﾊﾞﾝｺﾞｳ
           WRITE PY-ｼﾊﾗｲ-ﾚｺｰﾄﾞ
           IF NOT WS-PAYOUT-OK
               MOVE "PAYOUT" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-PAYOUT-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           MOVE 1210 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
           SET JR-ｶﾘｶﾀ TO TRUE
           MOVE LR-ｷﾝｶﾞｸ TO JR-ｷﾝｶﾞｸ
           MOVE "契約者貸付" TO JR-ﾃｷﾖｳ
           PERFORM ｼﾜｹ-RTN
           ADD 1 TO WS-ｶｼﾂｹｹﾝｽｳ
           ADD LR-ｷﾝｶﾞｸ TO WS-ｶｼﾂｹｷﾝｶﾞｸ.

      *----------------------------------------------------------------
      * ﾍﾝｻｲ-RTN: 返済額を残高から引く。
      *----------------------------------------------------------------
       ﾍﾝｻｲ-RTN.
           SUBTRACT LR-ｷﾝｶﾞｸ FROM WS-ｻﾞﾝﾀﾞｶ
           MOVE 1210 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
           SET JR-ｶｼｶﾀ TO TRUE
           MOVE LR-ｷﾝｶﾞｸ TO JR-ｷﾝｶﾞｸ
           MOVE "貸付金返済" TO JR-ﾃｷﾖｳ
           PERFORM ｼﾜｹ-RTN
           ADD 1 TO WS-ﾍﾝｻｲｹﾝｽｳ
           ADD LR-ｷﾝｶﾞｸ TO WS-ﾍﾝｻｲｷﾝｶﾞｸ.

      *----------------------------------------------------------------
      * ｼﾜｹ-RTN: 仕訳1行。科目、貸借、金額、摘要は呼ぶ側が入れる。
      *----------------------------------------------------------------
       ｼﾜｹ-RTN.
           MOVE CT-ｼﾞｶｲｼﾜｹﾊﾞﾝｺﾞｳ TO JR-ｼﾜｹﾊﾞﾝｺﾞｳ
           ADD 1 TO CT-ｼﾞｶｲｼﾜｹﾊﾞﾝｺﾞｳ
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO JR-ﾋﾂﾞｹ
           MOVE "GLD070" TO JR-ﾌﾟﾛｸﾞﾗﾑ
           MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO JR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           WRITE JR-ｼﾜｹ-ﾚｺｰﾄﾞ
           IF NOT WS-GLJRNL-OK
               MOVE "GLJRNL" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-GLJRNL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ｷｬｯｶ-RTN: 却下ファイルに理由を付けて書き、明細を報告する。
      *----------------------------------------------------------------
       ｷｬｯｶ-RTN.
           MOVE LR-ｶｼﾂｹ-ﾚｺｰﾄﾞ TO ｷｬｯｶ-ﾃﾞｰﾀ
           MOVE WS-ﾘﾕｳ TO ｷｬｯｶ-ﾘﾕｳ
           MOVE WS-ﾒｯｾｰｼﾞ(WS-ﾘﾕｳ) TO ｷｬｯｶ-ﾒｯｾｰｼﾞ
           WRITE ｷｬｯｶ-ﾚｺｰﾄﾞ
           ADD 1 TO WS-ｷｬｯｶｹﾝｽｳ
           MOVE 4 TO WS-ﾌｯｷｺｰﾄﾞ
           MOVE WS-ﾒｯｾｰｼﾞ(WS-ﾘﾕｳ) TO WS-ﾒｲｻｲ-ｹｯｶ
           PERFORM ﾒｲｻｲ-RTN.

      *----------------------------------------------------------------
      * ﾒｲｻｲ-RTN: 1件の明細行。
      *----------------------------------------------------------------
       ﾒｲｻｲ-RTN.
           MOVE LR-ｲﾗｲﾊﾞﾝｺﾞｳ TO WS-ﾒｲｻｲ-ｲﾗｲ
           MOVE LR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO WS-ﾒｲｻｲ-ｼｮｳｹﾝ
           MOVE WS-ｼｭﾍﾞﾂﾒｲ TO WS-ﾒｲｻｲ-ｼｭﾍﾞﾂ
           IF LR-ｷﾝｶﾞｸ NUMERIC
               MOVE LR-ｷﾝｶﾞｸ TO WS-ﾒｲｻｲ-ｷﾝｶﾞｸ
           ELSE
               MOVE 0 TO WS-ﾒｲｻｲ-ｷﾝｶﾞｸ
           END-IF
           MOVE WS-ｻﾞﾝﾀﾞｶ TO WS-ﾒｲｻｲ-ｻﾞﾝﾀﾞｶ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ﾒｲｻｲ-ｷﾞｮｳ.

      *----------------------------------------------------------------
      * ﾎｳｺｸ-RTN: 件数と合計。
      *----------------------------------------------------------------
       ﾎｳｺｸ-RTN.
           MOVE SPACES TO ﾎｳｺｸ-ｷﾞｮｳ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ
           MOVE "読込件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾖﾐｺﾐｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "貸付件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｶｼﾂｹｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "貸付金額" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｶｼﾂｹｷﾝｶﾞｸ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "返済件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾍﾝｻｲｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "返済金額" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾍﾝｻｲｷﾝｶﾞｸ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "繰入利息" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾘｿｸｷﾝｶﾞｸ TO WS-ｹﾝｽｳ-ｱﾀｲ
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
           CLOSE LONREQ-F POLMST-F GLSCTL-F PAYOUT-F GLJRNL-F
                 REJECT-F REPORT-F
           DISPLAY "GLD070 貸付 " WS-ｶｼﾂｹｹﾝｽｳ
                   " 返済 " WS-ﾍﾝｻｲｹﾝｽｳ
                   " 却下 " WS-ｷｬｯｶｹﾝｽｳ.

      *----------------------------------------------------------------
      * ｲｼﾞｮｳ-RTN: ファイル異常は 16 で終わる。
      *----------------------------------------------------------------
       ｲｼﾞｮｳ-RTN.
           DISPLAY "GLD070 ファイル異常 " WS-ｴﾗｰﾌｧｲﾙ
                   " 状態 " WS-ｴﾗｰｽﾃｰﾀｽ
           MOVE 16 TO RETURN-CODE
           STOP RUN.
