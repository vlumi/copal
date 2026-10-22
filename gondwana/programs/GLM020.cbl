      *----------------------------------------------------------------
      * GLM020: 責任準備金計算
      *
      * 月末。保障中 (有効・猶予中) の契約ごとに、責任準備金係数
      * (商品・経過年数、50年まで) を保険金額1,000円あたりに掛けて
      * 月末の責任準備金を求め、前月末との差額を繰入 (借方) または
      * 戻入 (貸方) として仕訳し、マスタに書き戻す。保障の終わった
      * 契約に残る準備金は全額戻入する。
      * 報告書は商品順に報告書作成機能 (REPORT WRITER) で: 商品ごとの
      * 見出しと計、総合計、頁番号。
      *
      * 復帰コード: 0 処理済、8 日付が日付でない、16 ファイル異常。
      *
      * 変更履歴
      * 日付       担当 変更
      * 1987-04-01 KT   GL-0001 初版。
      * 1995-10-02 KT   GL-0188 報告書作成機能で商品別に。
      *----------------------------------------------------------------
       IDENTIFICATION DIVISION.
       PROGRAM-ID. GLM020.
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
           SELECT RSVTB-F ASSIGN TO "RSVTB"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-RSVTB-ST.
           SELECT PRODTB-F ASSIGN TO "PRODTB"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-PRODTB-ST.
           SELECT SORT-F ASSIGN TO "SORTWK".
           SELECT GLSCTL-F ASSIGN TO "GLSCTL"
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-GLSCTL-ST.
           SELECT OPTIONAL GLJRNL-F ASSIGN TO "GLJRNL"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-GLJRNL-ST.
           SELECT REPORT-F ASSIGN TO "GLM020R"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-ST.

       DATA DIVISION.
       FILE SECTION.
       FD  POLMST-F.
       COPY POLMST.
       FD  RSVTB-F.
       COPY RSVTB.
       FD  PRODTB-F.
       COPY PRODTB.
       SD  SORT-F.
       01  SR-ｾｲﾚﾂ-ﾚｺｰﾄﾞ.
           05  SR-ｼｮｳﾋﾝ                PIC X(2).
           05  SR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ          PIC X(10).
           05  SR-ｼﾒｲ                  PIC X(20).
           05  SR-ｹｲｶﾈﾝｽｳ              PIC 99.
           05  SR-ﾎｹﾝｷﾝｶﾞｸ             PIC 9(9).
           05  SR-ｼﾞｭﾝﾋﾞｷﾝ             PIC 9(9).
           05  SR-ｿﾞｳｹﾞﾝ               PIC S9(9).
       FD  GLSCTL-F
           RECORD CONTAINS 64 CHARACTERS.
       COPY CTLREC.
       FD  GLJRNL-F.
       COPY GLJRNL.
       FD  REPORT-F
           REPORT IS RESERVE-RPT.

       WORKING-STORAGE SECTION.
       01  WS-ｽﾃｰﾀｽ.
           05  WS-POLMST-ST            PIC XX.
               88  WS-POLMST-OK                VALUE "00".
               88  WS-POLMST-ｵﾜﾘ               VALUE "10".
           05  WS-RSVTB-ST             PIC XX.
               88  WS-RSVTB-OK                 VALUE "00".
               88  WS-RSVTB-ｵﾜﾘ                VALUE "10".
           05  WS-PRODTB-ST            PIC XX.
               88  WS-PRODTB-OK                VALUE "00".
               88  WS-PRODTB-ｵﾜﾘ               VALUE "10".
           05  WS-GLSCTL-ST            PIC XX.
               88  WS-GLSCTL-OK                VALUE "00".
           05  WS-GLJRNL-ST            PIC XX.
               88  WS-GLJRNL-OK                VALUE "00" "05".
           05  WS-REPORT-ST            PIC XX.
               88  WS-REPORT-OK                VALUE "00".
       01  WS-ｺﾏﾝﾄﾞｷﾞｮｳ                PIC X(80).
       01  WS-ｴﾗｰﾌｧｲﾙ                  PIC X(8).
       01  WS-ｴﾗｰｽﾃｰﾀｽ                 PIC XX.
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ                PIC 9(8).
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ            PIC X(7).
       01  WS-ﾎｳｺｳ                     PIC X VALUE "W".
       01  WS-ﾍﾝｶﾝｹｯｶ                  PIC 9.
       01  WS-ｼｮｳﾋﾝﾒｲ                  PIC X(20).
       01  WS-ﾈﾝｽｳ                     PIC 99.
       COPY TBLWS.
       01  WS-ｼﾞｭﾝﾋﾞｷﾝ-TBL.
           05  WS-ｼﾞｭﾝﾋﾞｷﾝ-ｼｮｳﾋﾝ OCCURS 10 TIMES.
               10  WS-ｼﾞｭﾝﾋﾞｷﾝ-ｹｲｽｳ OCCURS 50 TIMES
                                           PIC 9(4)V9(3).
       01  WS-ｹﾝｽｳ.
           05  WS-ｹｲｻﾝｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ﾓﾄﾞｼｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ｼﾞｭﾝﾋﾞｷﾝｺﾞｳｹｲ        PIC 9(13) VALUE 0.
           05  WS-ｿﾞｳｹﾞﾝｺﾞｳｹｲ          PIC S9(13) VALUE 0.

       REPORT SECTION.
       RD  RESERVE-RPT
           CONTROLS ARE FINAL SR-ｼｮｳﾋﾝ
           PAGE LIMIT IS 60 LINES
           HEADING 1
           FIRST DETAIL 4
           LAST DETAIL 56
           FOOTING 58.
       01  TYPE IS PAGE HEADING.
           05  LINE 1.
               10  COLUMN 1   PIC X(32)
                   VALUE "ゴンドワナ生命  責任準備金計算表".
               10  COLUMN 36  PIC X(7)  SOURCE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ.
               10  COLUMN 70  PIC X(2)  VALUE "頁".
               10  COLUMN 73  PIC ZZ9   SOURCE PAGE-COUNTER.
           05  LINE 3.
               10  COLUMN 1   PIC X(8)  VALUE "証券番号".
               10  COLUMN 13  PIC X(4)  VALUE "氏名".
               10  COLUMN 35  PIC X(4)  VALUE "年数".
               10  COLUMN 45  PIC X(8)  VALUE "保険金額".
               10  COLUMN 58  PIC X(10) VALUE "責任準備金".
               10  COLUMN 72  PIC X(4)  VALUE "増減".
       01  TYPE IS CONTROL HEADING SR-ｼｮｳﾋﾝ.
           05  LINE PLUS 1.
               10  COLUMN 1   PIC X(4)  VALUE "商品".
               10  COLUMN 6   PIC X(2)  SOURCE SR-ｼｮｳﾋﾝ.
               10  COLUMN 10  PIC X(20) SOURCE WS-ｼｮｳﾋﾝﾒｲ.
       01  RESERVE-DETAIL TYPE IS DETAIL.
           05  LINE PLUS 1.
               10  COLUMN 1   PIC X(10) SOURCE SR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ.
               10  COLUMN 13  PIC X(20) SOURCE SR-ｼﾒｲ.
               10  COLUMN 36  PIC Z9    SOURCE SR-ｹｲｶﾈﾝｽｳ.
               10  COLUMN 41  PIC ZZZ,ZZZ,ZZ9 SOURCE SR-ﾎｹﾝｷﾝｶﾞｸ.
               10  COLUMN 56  PIC ZZZ,ZZZ,ZZ9 SOURCE SR-ｼﾞｭﾝﾋﾞｷﾝ.
               10  COLUMN 68  PIC ----,---,--9 SOURCE SR-ｿﾞｳｹﾞﾝ.
       01  TYPE IS CONTROL FOOTING SR-ｼｮｳﾋﾝ.
           05  LINE PLUS 1.
               10  COLUMN 1   PIC X(8)  VALUE "商品計".
               10  COLUMN 50  PIC Z,ZZZ,ZZZ,ZZZ,ZZ9 SUM SR-ｼﾞｭﾝﾋﾞｷﾝ.
               10  COLUMN 68  PIC ----,---,--9 SUM SR-ｿﾞｳｹﾞﾝ.
           05  LINE PLUS 1.
               10  COLUMN 1   PIC X(1)  VALUE " ".
       01  TYPE IS CONTROL FOOTING FINAL.
           05  LINE PLUS 1.
               10  COLUMN 1   PIC X(8)  VALUE "合計".
               10  COLUMN 50  PIC Z,ZZZ,ZZZ,ZZZ,ZZ9 SUM SR-ｼﾞｭﾝﾋﾞｷﾝ.
               10  COLUMN 68  PIC ----,---,--9 SUM SR-ｿﾞｳｹﾞﾝ.

       PROCEDURE DIVISION.
      *----------------------------------------------------------------
      * ﾒｲﾝ-RTN: 準備、計算と整列と報告、終了。
      *----------------------------------------------------------------
       ﾒｲﾝ-RTN.
           PERFORM ｼｮｷ-RTN
           SORT SORT-F
               ON ASCENDING KEY SR-ｼｮｳﾋﾝ SR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
               INPUT PROCEDURE IS ｹｲｻﾝ-RTN
               OUTPUT PROCEDURE IS ﾎｳｺｸ-RTN
           PERFORM ｼｭｳﾘｮｳ-RTN
           MOVE 0 TO RETURN-CODE
           STOP RUN.

      *----------------------------------------------------------------
      * ｼｮｷ-RTN: 業務日付、ファイル、テーブル、制御レコード。
      *----------------------------------------------------------------
       ｼｮｷ-RTN.
           ACCEPT WS-ｺﾏﾝﾄﾞｷﾞｮｳ FROM COMMAND-LINE
           MOVE WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8) TO WS-ｷﾞｮｳﾑﾋﾂﾞｹ
           CALL "GLC020" USING WS-ﾎｳｺｳ WS-ｷﾞｮｳﾑﾋﾂﾞｹ
                               WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ WS-ﾍﾝｶﾝｹｯｶ
           IF WS-ﾍﾝｶﾝｹｯｶ NOT = 0
               DISPLAY "GLM020 日付不正: " WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8)
               MOVE 8 TO RETURN-CODE
               STOP RUN
           END-IF
           OPEN INPUT RSVTB-F PRODTB-F
           OPEN I-O POLMST-F GLSCTL-F
           OPEN EXTEND GLJRNL-F
           OPEN OUTPUT REPORT-F
           EVALUATE TRUE
               WHEN NOT WS-RSVTB-OK
                   MOVE "RSVTB" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-RSVTB-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-PRODTB-OK
                   MOVE "PRODTB" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-PRODTB-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-POLMST-OK
                   MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-GLSCTL-OK
                   MOVE "GLSCTL" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-GLSCTL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-GLJRNL-OK
                   MOVE "GLJRNL" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-GLJRNL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-REPORT-OK
                   MOVE "GLM020R" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-REPORT-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN OTHER
                   MOVE SPACES TO WS-ｴﾗｰﾌｧｲﾙ
           END-EVALUATE
           IF WS-ｴﾗｰﾌｧｲﾙ NOT = SPACES
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           MOVE SPACES TO JR-ｼﾜｹ-ﾚｺｰﾄﾞ
           PERFORM ｼｮｳﾋﾝ-ﾖﾐｺﾐ-RTN
           PERFORM ｼﾞｭﾝﾋﾞｷﾝ-ﾖﾐｺﾐ-RTN
           READ GLSCTL-F
           IF NOT WS-GLSCTL-OK
               MOVE "GLSCTL" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-GLSCTL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ｼﾞｭﾝﾋﾞｷﾝ-ﾖﾐｺﾐ-RTN: 係数を商品・経過年数のテーブルに読み込む。
      *----------------------------------------------------------------
       ｼﾞｭﾝﾋﾞｷﾝ-ﾖﾐｺﾐ-RTN.
           MOVE ZEROS TO WS-ｼﾞｭﾝﾋﾞｷﾝ-TBL
           READ RSVTB-F
           PERFORM UNTIL NOT WS-RSVTB-OK
               SET WS-ｼｮｳﾋﾝ-IX TO 1
               SEARCH WS-ｼｮｳﾋﾝ
                   WHEN WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = RS-ｼｮｳﾋﾝ
                       MOVE RS-ｹｲｽｳ TO WS-ｼﾞｭﾝﾋﾞｷﾝ-ｹｲｽｳ(WS-ｼｮｳﾋﾝ-IX,
                                                        RS-ｹｲｶﾈﾝｽｳ)
               END-SEARCH
               READ RSVTB-F
           END-PERFORM
           IF NOT WS-RSVTB-ｵﾜﾘ
               MOVE "RSVTB" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-RSVTB-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ｹｲｻﾝ-RTN: マスタ全件。保障中は準備金を求め、終わった契約は
      * 戻し、差額を仕訳し、書き戻し、整列に渡す。
      *----------------------------------------------------------------
       ｹｲｻﾝ-RTN.
           READ POLMST-F NEXT
           PERFORM UNTIL NOT WS-POLMST-OK
               IF PM-ﾎｼｮｳﾁｭｳ
                   PERFORM ｾｷﾆﾝｼﾞｭﾝﾋﾞｷﾝ-RTN
                   PERFORM ｷｼﾞｭﾝ-RTN
               ELSE
                   IF PM-ｾｷﾆﾝｼﾞｭﾝﾋﾞｷﾝ > 0
                       MOVE 0 TO SR-ｼﾞｭﾝﾋﾞｷﾝ
                       ADD 1 TO WS-ﾓﾄﾞｼｹﾝｽｳ
                       PERFORM ｷｼﾞｭﾝ-RTN
                   END-IF
               END-IF
               READ POLMST-F NEXT
           END-PERFORM
           IF NOT WS-POLMST-ｵﾜﾘ
               MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ｾｷﾆﾝｼﾞｭﾝﾋﾞｷﾝ-RTN: 係数から月末の責任準備金。
      *----------------------------------------------------------------
       ｾｷﾆﾝｼﾞｭﾝﾋﾞｷﾝ-RTN.
           SET WS-ｼｮｳﾋﾝ-IX TO 1
           SEARCH WS-ｼｮｳﾋﾝ
               WHEN WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = PM-ｼｮｳﾋﾝ
                   CONTINUE
           END-SEARCH
           IF PM-ｹｲｶﾈﾝｽｳ > 50
               MOVE 50 TO WS-ﾈﾝｽｳ
           ELSE
               MOVE PM-ｹｲｶﾈﾝｽｳ TO WS-ﾈﾝｽｳ
           END-IF
           COMPUTE SR-ｼﾞｭﾝﾋﾞｷﾝ = PM-ﾎｹﾝｷﾝｶﾞｸ / 1000
               * WS-ｼﾞｭﾝﾋﾞｷﾝ-ｹｲｽｳ(WS-ｼｮｳﾋﾝ-IX, WS-ﾈﾝｽｳ)
           ADD 1 TO WS-ｹｲｻﾝｹﾝｽｳ.

      *----------------------------------------------------------------
      * ｷｼﾞｭﾝ-RTN: 差額の仕訳、マスタの更新、整列レコード。
      *----------------------------------------------------------------
       ｷｼﾞｭﾝ-RTN.
           COMPUTE SR-ｿﾞｳｹﾞﾝ = SR-ｼﾞｭﾝﾋﾞｷﾝ - PM-ｾｷﾆﾝｼﾞｭﾝﾋﾞｷﾝ
           IF SR-ｿﾞｳｹﾞﾝ NOT = 0
               PERFORM ｼﾜｹ-RTN
           END-IF
           MOVE SR-ｼﾞｭﾝﾋﾞｷﾝ TO PM-ｾｷﾆﾝｼﾞｭﾝﾋﾞｷﾝ
           REWRITE PM-ｹｲﾔｸ-ﾚｺｰﾄﾞ
           IF NOT WS-POLMST-OK
               MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           MOVE PM-ｼｮｳﾋﾝ TO SR-ｼｮｳﾋﾝ
           MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO SR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           MOVE PM-ｼﾒｲ TO SR-ｼﾒｲ
           MOVE PM-ｹｲｶﾈﾝｽｳ TO SR-ｹｲｶﾈﾝｽｳ
           MOVE PM-ﾎｹﾝｷﾝｶﾞｸ TO SR-ﾎｹﾝｷﾝｶﾞｸ
           RELEASE SR-ｾｲﾚﾂ-ﾚｺｰﾄﾞ
           ADD SR-ｼﾞｭﾝﾋﾞｷﾝ TO WS-ｼﾞｭﾝﾋﾞｷﾝｺﾞｳｹｲ
           ADD SR-ｿﾞｳｹﾞﾝ TO WS-ｿﾞｳｹﾞﾝｺﾞｳｹｲ.

      *----------------------------------------------------------------
      * ｼﾜｹ-RTN: 繰入は借方、戻入は貸方。
      *----------------------------------------------------------------
       ｼﾜｹ-RTN.
           MOVE CT-ｼﾞｶｲｼﾜｹﾊﾞﾝｺﾞｳ TO JR-ｼﾜｹﾊﾞﾝｺﾞｳ
           ADD 1 TO CT-ｼﾞｶｲｼﾜｹﾊﾞﾝｺﾞｳ
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO JR-ﾋﾂﾞｹ
           MOVE "GLM020" TO JR-ﾌﾟﾛｸﾞﾗﾑ
           MOVE 5210 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
           IF SR-ｿﾞｳｹﾞﾝ > 0
               SET JR-ｶﾘｶﾀ TO TRUE
               MOVE SR-ｿﾞｳｹﾞﾝ TO JR-ｷﾝｶﾞｸ
               MOVE "責任準備金繰入" TO JR-ﾃｷﾖｳ
           ELSE
               SET JR-ｶｼｶﾀ TO TRUE
               COMPUTE JR-ｷﾝｶﾞｸ = 0 - SR-ｿﾞｳｹﾞﾝ
               MOVE "責任準備金戻入" TO JR-ﾃｷﾖｳ
           END-IF
           MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO JR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           WRITE JR-ｼﾜｹ-ﾚｺｰﾄﾞ
           IF NOT WS-GLJRNL-OK
               MOVE "GLJRNL" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-GLJRNL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ﾎｳｺｸ-RTN: 整列済みの1件ごとに明細を生成する。商品名は
      * 見出しのために商品が変わるたびに引く。
      *----------------------------------------------------------------
       ﾎｳｺｸ-RTN.
           INITIATE RESERVE-RPT
           RETURN SORT-F
               AT END
                   CONTINUE
               NOT AT END
                   PERFORM ﾒｲｻｲ-RTN
                   PERFORM UNTIL SR-ｼｮｳﾋﾝ = SPACES
                       RETURN SORT-F
                           AT END
                               MOVE SPACES TO SR-ｼｮｳﾋﾝ
                           NOT AT END
                               PERFORM ﾒｲｻｲ-RTN
                       END-RETURN
                   END-PERFORM
           END-RETURN
           TERMINATE RESERVE-RPT.

      *----------------------------------------------------------------
      * ﾒｲｻｲ-RTN: 商品名を引き、1行生成する。
      *----------------------------------------------------------------
       ﾒｲｻｲ-RTN.
           SET WS-ｼｮｳﾋﾝ-IX TO 1
           SEARCH WS-ｼｮｳﾋﾝ
               AT END
                   MOVE SPACES TO WS-ｼｮｳﾋﾝﾒｲ
               WHEN WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = SR-ｼｮｳﾋﾝ
                   MOVE WS-ｼｮｳﾋﾝ-ﾒｲｼｮｳ(WS-ｼｮｳﾋﾝ-IX) TO WS-ｼｮｳﾋﾝﾒｲ
           END-SEARCH
           GENERATE RESERVE-DETAIL.

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
           CLOSE POLMST-F RSVTB-F PRODTB-F GLSCTL-F GLJRNL-F
                 REPORT-F
           DISPLAY "GLM020 計算 " WS-ｹｲｻﾝｹﾝｽｳ " 戻入 " WS-ﾓﾄﾞｼｹﾝｽｳ
                   " 準備金 " WS-ｼﾞｭﾝﾋﾞｷﾝｺﾞｳｹｲ.

      *----------------------------------------------------------------
      * ｲｼﾞｮｳ-RTN: ファイル異常は 16 で終わる。
      *----------------------------------------------------------------
       ｲｼﾞｮｳ-RTN.
           DISPLAY "GLM020 ファイル異常 " WS-ｴﾗｰﾌｧｲﾙ
                   " 状態 " WS-ｴﾗｰｽﾃｰﾀｽ
           MOVE 16 TO RETURN-CODE
           STOP RUN.

       COPY PRDLOD.
