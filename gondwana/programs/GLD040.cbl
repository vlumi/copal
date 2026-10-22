      *----------------------------------------------------------------
      * GLD040: 契約変更
      *
      * 日次。支社からの契約変更依頼ファイル MNTREQ を読み、契約
      * マスタに反映する。住所変更、口座変更、氏名変更 (漢字とカナ)、
      * 受取人変更は該当項目が埋まっていること。復活は契約が失効中
      * で、失効日から1年以内で、窓口で収納した額が未納保険料
      * (払込済年月の翌月から業務月までの月数に保険料を掛けた額)
      * 以上であること。復活した契約は有効に戻り、払込済年月を業務
      * 月とし、収納額を保険料収入に仕訳する。
      *
      * 復帰コード: 0 全件反映、4 却下あり、8 日付が日付でない、
      *             16 ファイル異常。
      *
      * 変更履歴
      * 日付       担当 変更
      * 1987-04-01 KT   GL-0001 初版。
      * 2003-10-06 MS   GL-0310 口座変更。
      * 2011-02-14 YN   GL-0520 窓口からの復活。
      *----------------------------------------------------------------
       IDENTIFICATION DIVISION.
       PROGRAM-ID. GLD040.
       AUTHOR. ゴンドワナ生命 システム部.
       DATE-WRITTEN. 1987-04-01.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT MNTREQ-F ASSIGN TO "MNTREQ"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-MNTREQ-ST.
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
           SELECT REJECT-F ASSIGN TO "MNTREJ"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REJECT-ST.
           SELECT REPORT-F ASSIGN TO "GLD040R"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-ST.

       DATA DIVISION.
       FILE SECTION.
       FD  MNTREQ-F.
       COPY MNTREQ.
       FD  POLMST-F.
       COPY POLMST.
       FD  GLSCTL-F
           RECORD CONTAINS 64 CHARACTERS.
       COPY CTLREC.
       FD  GLJRNL-F.
       COPY GLJRNL.
       FD  REJECT-F.
       01  ｷｬｯｶ-ﾚｺｰﾄﾞ.
           05  ｷｬｯｶ-ﾃﾞｰﾀ               PIC X(180).
           05  ｷｬｯｶ-ﾘﾕｳ                PIC 99.
           05  ｷｬｯｶ-ﾒｯｾｰｼﾞ             PIC X(20).
       FD  REPORT-F.
       01  ﾎｳｺｸ-ｷﾞｮｳ                   PIC X(100).

       WORKING-STORAGE SECTION.
       01  WS-ｽﾃｰﾀｽ.
           05  WS-MNTREQ-ST            PIC XX.
               88  WS-MNTREQ-OK                VALUE "00".
               88  WS-MNTREQ-ｵﾜﾘ               VALUE "10".
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
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾌﾞﾝｶｲ REDEFINES WS-ｷﾞｮｳﾑﾋﾂﾞｹ.
           05  WS-ｷﾞｮｳﾑ-YYYY           PIC 9(4).
           05  WS-ｷﾞｮｳﾑ-MM             PIC 99.
           05  WS-ｷﾞｮｳﾑ-DD             PIC 99.
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ            PIC X(7).
       01  WS-ﾎｳｺｳ                     PIC X.
       01  WS-ﾍﾝｶﾝｹｯｶ                  PIC 9.
       01  WS-ｼｯｺｳﾋﾞ                   PIC 9(8).
       01  WS-ｹｲｶﾆｯｽｳ                  PIC S9(7).
       01  WS-ﾊﾗｲｺﾐｽﾞﾐ-ﾌﾞﾝｶｲ.
           05  WS-ﾊﾗｲｺﾐ-YYYY           PIC 9(4).
           05  WS-ﾊﾗｲｺﾐ-MM             PIC 99.
       01  WS-ﾐﾉｳﾂｷｽｳ                  PIC S9(5).
       01  WS-ﾐﾉｳｶﾞｸ                   PIC 9(11).
       01  WS-ｹﾝｽｳ.
           05  WS-ﾖﾐｺﾐｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ﾊﾝｴｲｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ｷｬｯｶｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ﾌｯｶﾂｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ﾌｯｶﾂｷﾝｶﾞｸ            PIC 9(13) VALUE 0.
       01  WS-ﾘﾕｳ                      PIC 99.
           88  WS-ﾘﾕｳ-ﾅｼ                       VALUE 0.
           88  WS-ｹｲﾔｸ-ﾅｼ                      VALUE 1.
           88  WS-ｼﾞｭｳｼｮ-ﾅｼ                    VALUE 2.
           88  WS-ｺｳｻﾞ-ﾌｾｲ                     VALUE 3.
           88  WS-ｼﾒｲ-ﾅｼ                       VALUE 4.
           88  WS-ｳｹﾄﾘﾆﾝ-ﾅｼ                    VALUE 5.
           88  WS-ｼｯｺｳﾁｭｳ-ﾃﾞﾅｲ                 VALUE 6.
           88  WS-ﾌｯｶﾂｷｶﾝ-ｹｲｶ                  VALUE 7.
           88  WS-ﾐﾉｳｶﾞｸ-ﾌｿｸ                   VALUE 8.
           88  WS-ｼｭﾍﾞﾂ-ﾌｾｲ                    VALUE 9.
       01  WS-ﾘﾕｳ-ﾒｯｾｰｼﾞ.
           05  FILLER PIC X(20) VALUE "契約なし".
           05  FILLER PIC X(20) VALUE "住所なし".
           05  FILLER PIC X(20) VALUE "口座不正".
           05  FILLER PIC X(20) VALUE "氏名なし".
           05  FILLER PIC X(20) VALUE "受取人なし".
           05  FILLER PIC X(20) VALUE "失効中でない".
           05  FILLER PIC X(20) VALUE "復活期間経過".
           05  FILLER PIC X(20) VALUE "未納額不足".
           05  FILLER PIC X(20) VALUE "種別不正".
       01  WS-ﾘﾕｳ-ﾒｯｾｰｼﾞ-TBL REDEFINES WS-ﾘﾕｳ-ﾒｯｾｰｼﾞ.
           05  WS-ﾒｯｾｰｼﾞ OCCURS 9 TIMES PIC X(20).
       01  WS-ｼｭﾍﾞﾂﾒｲ                  PIC X(12).
       01  WS-ﾌｯｷｺｰﾄﾞ                  PIC 99 VALUE 0.
       01  WS-ﾎｳｺｸ.
           05  WS-ﾋｮｳﾀﾞｲ-ｷﾞｮｳ.
               10  FILLER                  PIC X(28)
                   VALUE "ゴンドワナ生命  契約変更".
               10  FILLER                  PIC X(4) VALUE SPACES.
               10  WS-ﾋｮｳﾀﾞｲ-ﾋﾂﾞｹ          PIC X(7).
           05  WS-ﾒｲｻｲ-ｷﾞｮｳ.
               10  WS-ﾒｲｻｲ-ｲﾗｲ             PIC X(8).
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾒｲｻｲ-ｼｮｳｹﾝ           PIC X(10).
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾒｲｻｲ-ｼｭﾍﾞﾂ           PIC X(12).
               10  WS-ﾒｲｻｲ-ｷﾝｶﾞｸ           PIC ZZZ,ZZZ,ZZ9.
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
           PERFORM ｼｮﾘ-RTN UNTIL WS-MNTREQ-ｵﾜﾘ
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
               DISPLAY "GLD040 日付不正: " WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8)
               MOVE 8 TO RETURN-CODE
               STOP RUN
           END-IF
           OPEN INPUT MNTREQ-F
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
           PERFORM MNTREQ-ﾖﾐｺﾐ-RTN.

      *----------------------------------------------------------------
      * ｵｰﾌﾟﾝ-ｹﾝｻ-RTN: 開けなかったファイルがあれば異常終了。
      *----------------------------------------------------------------
       ｵｰﾌﾟﾝ-ｹﾝｻ-RTN.
           EVALUATE TRUE
               WHEN NOT WS-MNTREQ-OK
                   MOVE "MNTREQ" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-MNTREQ-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
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
                   MOVE "MNTREJ" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-REJECT-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-REPORT-OK
                   MOVE "GLD040R" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-REPORT-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN OTHER
                   MOVE SPACES TO WS-ｴﾗｰﾌｧｲﾙ
           END-EVALUATE
           IF WS-ｴﾗｰﾌｧｲﾙ NOT = SPACES
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * MNTREQ-ﾖﾐｺﾐ-RTN: 依頼ファイルの次の行。
      *----------------------------------------------------------------
       MNTREQ-ﾖﾐｺﾐ-RTN.
           READ MNTREQ-F
           IF WS-MNTREQ-OK
               ADD 1 TO WS-ﾖﾐｺﾐｹﾝｽｳ
           ELSE
               IF NOT WS-MNTREQ-ｵﾜﾘ
                   MOVE "MNTREQ" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-MNTREQ-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
                   PERFORM ｲｼﾞｮｳ-RTN
               END-IF
           END-IF.

      *----------------------------------------------------------------
      * ｼｮﾘ-RTN: 契約を読み、種別ごとに検査して反映か却下。
      *----------------------------------------------------------------
       ｼｮﾘ-RTN.
           SET WS-ﾘﾕｳ-ﾅｼ TO TRUE
           MOVE 0 TO WS-ﾒｲｻｲ-ｷﾝｶﾞｸ
           MOVE MR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
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
               WHEN MR-ｼﾞｭｳｼｮﾍﾝｺｳ
                   MOVE "住所変更" TO WS-ｼｭﾍﾞﾂﾒｲ
               WHEN MR-ｺｳｻﾞﾍﾝｺｳ
                   MOVE "口座変更" TO WS-ｼｭﾍﾞﾂﾒｲ
               WHEN MR-ｼﾒｲﾍﾝｺｳ
                   MOVE "氏名変更" TO WS-ｼｭﾍﾞﾂﾒｲ
               WHEN MR-ｳｹﾄﾘﾆﾝﾍﾝｺｳ
                   MOVE "受取人変更" TO WS-ｼｭﾍﾞﾂﾒｲ
               WHEN MR-ﾌｯｶﾂ
                   MOVE "復活" TO WS-ｼｭﾍﾞﾂﾒｲ
               WHEN OTHER
                   MOVE "不明" TO WS-ｼｭﾍﾞﾂﾒｲ
           END-EVALUATE
           IF WS-ﾘﾕｳ-ﾅｼ
               PERFORM ﾍﾝｺｳ-RTN
           END-IF
           IF WS-ﾘﾕｳ-ﾅｼ
               REWRITE PM-ｹｲﾔｸ-ﾚｺｰﾄﾞ
               IF NOT WS-POLMST-OK
                   MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
                   PERFORM ｲｼﾞｮｳ-RTN
               END-IF
               ADD 1 TO WS-ﾊﾝｴｲｹﾝｽｳ
               MOVE "反映" TO WS-ﾒｲｻｲ-ｹｯｶ
               PERFORM ﾒｲｻｲ-RTN
           ELSE
               PERFORM ｷｬｯｶ-RTN
           END-IF
           PERFORM MNTREQ-ﾖﾐｺﾐ-RTN.

      *----------------------------------------------------------------
      * ﾍﾝｺｳ-RTN: 種別ごとの検査と反映。
      *----------------------------------------------------------------
       ﾍﾝｺｳ-RTN.
           EVALUATE TRUE
               WHEN MR-ｼﾞｭｳｼｮﾍﾝｺｳ
                   IF MR-ｼﾞｭｳｼｮ = SPACES
                       SET WS-ｼﾞｭｳｼｮ-ﾅｼ TO TRUE
                   ELSE
                       MOVE MR-ｼﾞｭｳｼｮ TO PM-ｼﾞｭｳｼｮ
                   END-IF
               WHEN MR-ｺｳｻﾞﾍﾝｺｳ
                   IF MR-ｷﾞﾝｺｳｺｰﾄﾞ = 0 OR MR-ｺｳｻﾞﾊﾞﾝｺﾞｳ = 0
                           OR MR-ｺｳｻﾞﾒｲｷﾞ = SPACES
                       SET WS-ｺｳｻﾞ-ﾌｾｲ TO TRUE
                   ELSE
                       MOVE MR-ｺｳｻﾞ TO PM-ｺｳｻﾞ
                   END-IF
               WHEN MR-ｼﾒｲﾍﾝｺｳ
                   IF MR-ｼﾒｲ = SPACES OR MR-ｼﾒｲｶﾅ = SPACES
                       SET WS-ｼﾒｲ-ﾅｼ TO TRUE
                   ELSE
                       MOVE MR-ｼﾒｲ TO PM-ｼﾒｲ
                       MOVE MR-ｼﾒｲｶﾅ TO PM-ｼﾒｲｶﾅ
                   END-IF
               WHEN MR-ｳｹﾄﾘﾆﾝﾍﾝｺｳ
                   IF MR-ｳｹﾄﾘﾆﾝ = SPACES
                       SET WS-ｳｹﾄﾘﾆﾝ-ﾅｼ TO TRUE
                   ELSE
                       MOVE MR-ｳｹﾄﾘﾆﾝ TO PM-ｳｹﾄﾘﾆﾝ
                   END-IF
               WHEN MR-ﾌｯｶﾂ
                   PERFORM ﾌｯｶﾂ-RTN
               WHEN OTHER
                   SET WS-ｼｭﾍﾞﾂ-ﾌｾｲ TO TRUE
           END-EVALUATE.

      *----------------------------------------------------------------
      * ﾌｯｶﾂ-RTN: 失効中で、失効日から1年以内で、未納額が収納されて
      * いれば有効に戻す。
      *----------------------------------------------------------------
       ﾌｯｶﾂ-RTN.
           IF NOT PM-ｼｯｺｳ
               SET WS-ｼｯｺｳﾁｭｳ-ﾃﾞﾅｲ TO TRUE
           ELSE
               MOVE "S" TO WS-ﾎｳｺｳ
               CALL "GLC020" USING WS-ﾎｳｺｳ WS-ｼｯｺｳﾋﾞ
                                   PM-ｼｯｺｳﾋﾞ WS-ﾍﾝｶﾝｹｯｶ
               COMPUTE WS-ｹｲｶﾆｯｽｳ
                   = FUNCTION INTEGER-OF-DATE(WS-ｷﾞｮｳﾑﾋﾂﾞｹ)
                   - FUNCTION INTEGER-OF-DATE(WS-ｼｯｺｳﾋﾞ)
               MOVE PM-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ TO WS-ﾊﾗｲｺﾐｽﾞﾐ-ﾌﾞﾝｶｲ
               COMPUTE WS-ﾐﾉｳﾂｷｽｳ = (WS-ｷﾞｮｳﾑ-YYYY - WS-ﾊﾗｲｺﾐ-YYYY) * 12
                                 + WS-ｷﾞｮｳﾑ-MM - WS-ﾊﾗｲｺﾐ-MM
               IF PM-ﾈﾝﾊﾞﾗｲ
                   COMPUTE WS-ﾐﾉｳﾂｷｽｳ = (WS-ﾐﾉｳﾂｷｽｳ + 11) / 12
               END-IF
               COMPUTE WS-ﾐﾉｳｶﾞｸ = WS-ﾐﾉｳﾂｷｽｳ * PM-ﾎｹﾝﾘｮｳ
               EVALUATE TRUE
                   WHEN WS-ｹｲｶﾆｯｽｳ > 365
                       SET WS-ﾌｯｶﾂｷｶﾝ-ｹｲｶ TO TRUE
                   WHEN MR-ﾆｭｳｷﾝｶﾞｸ < WS-ﾐﾉｳｶﾞｸ
                       SET WS-ﾐﾉｳｶﾞｸ-ﾌｿｸ TO TRUE
                   WHEN OTHER
                       PERFORM ﾌｯｶﾂ-ﾊﾝｴｲ-RTN
               END-EVALUATE
           END-IF.

      *----------------------------------------------------------------
      * ﾌｯｶﾂ-ﾊﾝｴｲ-RTN: 有効に戻し、払込済年月を業務月とし、収納額を
      * 仕訳する。
      *----------------------------------------------------------------
       ﾌｯｶﾂ-ﾊﾝｴｲ-RTN.
           SET PM-ﾕｳｺｳ TO TRUE
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ(1:6) TO PM-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ
           MOVE 0 TO PM-ﾐﾉｳｶｲｽｳ
           MOVE SPACES TO PM-ﾕｳﾖｶｲｼﾋﾞ PM-ｼｯｺｳﾋﾞ
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO PM-ｼﾞｮｳﾀｲﾋﾂﾞｹ
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ(1:6) TO PM-ｻｲｼｭｳﾆｭｳｷﾝﾈﾝｹﾞﾂ
           MOVE MR-ﾆｭｳｷﾝｶﾞｸ TO PM-ｻｲｼｭｳﾆｭｳｷﾝｶﾞｸ
           MOVE SPACES TO JR-ｼﾜｹ-ﾚｺｰﾄﾞ
           MOVE CT-ｼﾞｶｲｼﾜｹﾊﾞﾝｺﾞｳ TO JR-ｼﾜｹﾊﾞﾝｺﾞｳ
           ADD 1 TO CT-ｼﾞｶｲｼﾜｹﾊﾞﾝｺﾞｳ
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO JR-ﾋﾂﾞｹ
           MOVE "GLD040" TO JR-ﾌﾟﾛｸﾞﾗﾑ
           MOVE 4110 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
           SET JR-ｶｼｶﾀ TO TRUE
           MOVE MR-ﾆｭｳｷﾝｶﾞｸ TO JR-ｷﾝｶﾞｸ
           MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO JR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           MOVE "復活時保険料" TO JR-ﾃｷﾖｳ
           WRITE JR-ｼﾜｹ-ﾚｺｰﾄﾞ
           ADD 1 TO WS-ﾌｯｶﾂｹﾝｽｳ
           ADD MR-ﾆｭｳｷﾝｶﾞｸ TO WS-ﾌｯｶﾂｷﾝｶﾞｸ
           MOVE MR-ﾆｭｳｷﾝｶﾞｸ TO WS-ﾒｲｻｲ-ｷﾝｶﾞｸ.

      *----------------------------------------------------------------
      * ｷｬｯｶ-RTN: 却下ファイルに理由を付けて書き、明細を報告する。
      *----------------------------------------------------------------
       ｷｬｯｶ-RTN.
           MOVE MR-ﾍﾝｺｳ-ﾚｺｰﾄﾞ TO ｷｬｯｶ-ﾃﾞｰﾀ
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
           MOVE MR-ｲﾗｲﾊﾞﾝｺﾞｳ TO WS-ﾒｲｻｲ-ｲﾗｲ
           MOVE MR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO WS-ﾒｲｻｲ-ｼｮｳｹﾝ
           MOVE WS-ｼｭﾍﾞﾂﾒｲ TO WS-ﾒｲｻｲ-ｼｭﾍﾞﾂ
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
           MOVE "反映件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾊﾝｴｲｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "却下件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｷｬｯｶｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "復活件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾌｯｶﾂｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "復活時収納額" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾌｯｶﾂｷﾝｶﾞｸ TO WS-ｹﾝｽｳ-ｱﾀｲ
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
           CLOSE MNTREQ-F POLMST-F GLSCTL-F GLJRNL-F REJECT-F
                 REPORT-F
           DISPLAY "GLD040 反映 " WS-ﾊﾝｴｲｹﾝｽｳ
                   " 却下 " WS-ｷｬｯｶｹﾝｽｳ.

      *----------------------------------------------------------------
      * ｲｼﾞｮｳ-RTN: ファイル異常は 16 で終わる。
      *----------------------------------------------------------------
       ｲｼﾞｮｳ-RTN.
           DISPLAY "GLD040 ファイル異常 " WS-ｴﾗｰﾌｧｲﾙ
                   " 状態 " WS-ｴﾗｰｽﾃｰﾀｽ
           MOVE 16 TO RETURN-CODE
           STOP RUN.
