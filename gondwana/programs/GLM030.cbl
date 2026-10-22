      *----------------------------------------------------------------
      * GLM030: 代理店手数料計算
      *
      * 月末。当月に保険料の入金があった契約 (最終入金年月が業務月)
      * ごとに、入金額に経過年数による料率を掛けて代理店手数料を
      * 求める: 1年目 40%、2〜3年目 7%、4〜10年目 2%、11年目から
      * なし、円未満切捨て。代理店順に整列し、報告書作成機能で代理店
      * ごとの見出し (代理店マスタの名称、支社、入社年月) と計を
      * 印字する。代理店ごとの合計を手数料ファイルに書き、仕訳する。
      * 代理店マスタにない代理店は "(未登録)" として印字する。
      *
      * 復帰コード: 0 処理済、8 日付が日付でない、16 ファイル異常。
      *
      * 変更履歴
      * 日付       担当 変更
      * 1987-04-01 KT   GL-0001 初版。
      * 1995-10-02 KT   GL-0188 報告書作成機能で代理店別に。
      *----------------------------------------------------------------
       IDENTIFICATION DIVISION.
       PROGRAM-ID. GLM030.
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
           SELECT AGTMST-F ASSIGN TO "AGTMST"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-AGTMST-ST.
           SELECT SORT-F ASSIGN TO "SORTWK".
           SELECT GLSCTL-F ASSIGN TO "GLSCTL"
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-GLSCTL-ST.
           SELECT OPTIONAL GLJRNL-F ASSIGN TO "GLJRNL"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-GLJRNL-ST.
           SELECT COMMIS-F ASSIGN TO "COMMIS"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-COMMIS-ST.
           SELECT REPORT-F ASSIGN TO "GLM030R"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-ST.

       DATA DIVISION.
       FILE SECTION.
       FD  POLMST-F.
       COPY POLMST.
       FD  AGTMST-F.
       COPY AGTMST.
       SD  SORT-F.
       01  SR-ｾｲﾚﾂ-ﾚｺｰﾄﾞ.
           05  SR-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ           PIC X(6).
           05  SR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ          PIC X(10).
           05  SR-ｼｮｳﾋﾝ                PIC X(2).
           05  SR-ｹｲｶﾈﾝｽｳ              PIC 99.
           05  SR-ﾆｭｳｷﾝｶﾞｸ             PIC 9(7).
           05  SR-ﾘﾂ                   PIC 99.
           05  SR-ﾃｽｳﾘｮｳ               PIC 9(7).
       FD  GLSCTL-F
           RECORD CONTAINS 64 CHARACTERS.
       COPY CTLREC.
       FD  GLJRNL-F.
       COPY GLJRNL.
       FD  COMMIS-F.
       01  CM-ﾃｽｳﾘｮｳ-ﾚｺｰﾄﾞ.
           05  CM-ﾈﾝｹﾞﾂ                PIC 9(6).
           05  CM-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ           PIC X(6).
           05  CM-ｹﾝｽｳ                 PIC 9(5).
           05  CM-ﾆｭｳｷﾝｶﾞｸ             PIC 9(11).
           05  CM-ﾃｽｳﾘｮｳ               PIC 9(9).
       FD  REPORT-F
           REPORT IS COMMISSION-RPT.

       WORKING-STORAGE SECTION.
       01  WS-ｽﾃｰﾀｽ.
           05  WS-POLMST-ST            PIC XX.
               88  WS-POLMST-OK                VALUE "00".
               88  WS-POLMST-ｵﾜﾘ               VALUE "10".
           05  WS-AGTMST-ST            PIC XX.
               88  WS-AGTMST-OK                VALUE "00".
               88  WS-AGTMST-ｵﾜﾘ               VALUE "10".
           05  WS-GLSCTL-ST            PIC XX.
               88  WS-GLSCTL-OK                VALUE "00".
           05  WS-GLJRNL-ST            PIC XX.
               88  WS-GLJRNL-OK                VALUE "00" "05".
           05  WS-COMMIS-ST            PIC XX.
               88  WS-COMMIS-OK                VALUE "00".
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
       01  WS-ﾎｳｺｳ                     PIC X VALUE "W".
       01  WS-ﾍﾝｶﾝｹｯｶ                  PIC 9.
       01  WS-ﾀﾞｲﾘﾃﾝ-TBL.
           05  WS-ﾀﾞｲﾘﾃﾝｹﾝｽｳ           PIC 999 VALUE 0.
           05  WS-ﾀﾞｲﾘﾃﾝ OCCURS 200 TIMES INDEXED BY WS-ﾀﾞｲﾘﾃﾝ-IX.
               10  WS-ﾀﾞｲﾘﾃﾝ-ｺｰﾄﾞ          PIC X(6).
               10  WS-ﾀﾞｲﾘﾃﾝ-ｼﾒｲ           PIC X(20).
               10  WS-ﾀﾞｲﾘﾃﾝ-ｼﾃﾝ           PIC X(4).
               10  WS-ﾀﾞｲﾘﾃﾝ-ﾆｭｳｼｬ         PIC 9(4).
       01  WS-ﾐﾀﾞｼ.
           05  WS-ﾐﾀﾞｼ-ｼﾒｲ             PIC X(20).
           05  WS-ﾐﾀﾞｼ-ｼﾃﾝ             PIC X(4).
           05  WS-ﾐﾀﾞｼ-ﾆｭｳｼｬ           PIC X(5).
       01  WS-ｿｳﾀﾞｲﾘﾃﾝ                 PIC X(6).
       01  WS-ｹﾝｽｳ.
           05  WS-ﾀｲｼｮｳｹﾝｽｳ            PIC 9(7) VALUE 0.
           05  WS-ﾃｽｳﾘｮｳｺﾞｳｹｲ          PIC 9(13) VALUE 0.
           05  WS-ﾀﾞｲﾘﾃﾝ-ｹﾝｽｳ          PIC 9(5) VALUE 0.
           05  WS-ﾀﾞｲﾘﾃﾝ-ﾆｭｳｷﾝ         PIC 9(11) VALUE 0.
           05  WS-ﾀﾞｲﾘﾃﾝ-ﾃｽｳﾘｮｳ        PIC 9(9) VALUE 0.

       REPORT SECTION.
       RD  COMMISSION-RPT
           CONTROLS ARE FINAL SR-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ
           PAGE LIMIT IS 60 LINES
           HEADING 1
           FIRST DETAIL 4
           LAST DETAIL 56
           FOOTING 58.
       01  TYPE IS PAGE HEADING.
           05  LINE 1.
               10  COLUMN 1   PIC X(34)
                   VALUE "ゴンドワナ生命  代理店手数料計算表".
               10  COLUMN 38  PIC X(7)  SOURCE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ.
               10  COLUMN 70  PIC X(2)  VALUE "頁".
               10  COLUMN 73  PIC ZZ9   SOURCE PAGE-COUNTER.
           05  LINE 3.
               10  COLUMN 1   PIC X(8)  VALUE "証券番号".
               10  COLUMN 13  PIC X(4)  VALUE "商品".
               10  COLUMN 19  PIC X(4)  VALUE "年数".
               10  COLUMN 29  PIC X(6)  VALUE "入金額".
               10  COLUMN 38  PIC X(4)  VALUE "料率".
               10  COLUMN 48  PIC X(6)  VALUE "手数料".
       01  TYPE IS CONTROL HEADING SR-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ.
           05  LINE PLUS 1.
               10  COLUMN 1   PIC X(6)  VALUE "代理店".
               10  COLUMN 8   PIC X(6)  SOURCE SR-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ.
               10  COLUMN 16  PIC X(20) SOURCE WS-ﾐﾀﾞｼ-ｼﾒｲ.
               10  COLUMN 38  PIC X(4)  SOURCE WS-ﾐﾀﾞｼ-ｼﾃﾝ.
               10  COLUMN 44  PIC X(8)  VALUE "入社年月".
               10  COLUMN 53  PIC X(5)  SOURCE WS-ﾐﾀﾞｼ-ﾆｭｳｼｬ.
       01  COMMISSION-DETAIL TYPE IS DETAIL.
           05  LINE PLUS 1.
               10  COLUMN 1   PIC X(10) SOURCE SR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ.
               10  COLUMN 13  PIC X(2)  SOURCE SR-ｼｮｳﾋﾝ.
               10  COLUMN 20  PIC Z9    SOURCE SR-ｹｲｶﾈﾝｽｳ.
               10  COLUMN 25  PIC Z,ZZZ,ZZ9 SOURCE SR-ﾆｭｳｷﾝｶﾞｸ.
               10  COLUMN 38  PIC Z9    SOURCE SR-ﾘﾂ.
               10  COLUMN 40  PIC X(1)  VALUE "%".
               10  COLUMN 44  PIC Z,ZZZ,ZZ9 SOURCE SR-ﾃｽｳﾘｮｳ.
       01  TYPE IS CONTROL FOOTING SR-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ.
           05  LINE PLUS 1.
               10  COLUMN 1   PIC X(8)  VALUE "代理店計".
               10  COLUMN 21  PIC ZZZ,ZZZ,ZZ9 SUM SR-ﾆｭｳｷﾝｶﾞｸ.
               10  COLUMN 40  PIC ZZZ,ZZZ,ZZ9 SUM SR-ﾃｽｳﾘｮｳ.
           05  LINE PLUS 1.
               10  COLUMN 1   PIC X(1)  VALUE " ".
       01  TYPE IS CONTROL FOOTING FINAL.
           05  LINE PLUS 1.
               10  COLUMN 1   PIC X(8)  VALUE "合計".
               10  COLUMN 21  PIC ZZZ,ZZZ,ZZ9 SUM SR-ﾆｭｳｷﾝｶﾞｸ.
               10  COLUMN 40  PIC ZZZ,ZZZ,ZZ9 SUM SR-ﾃｽｳﾘｮｳ.

       PROCEDURE DIVISION.
      *----------------------------------------------------------------
      * ﾒｲﾝ-RTN: 準備、選択と整列と報告、終了。
      *----------------------------------------------------------------
       ﾒｲﾝ-RTN.
           PERFORM ｼｮｷ-RTN
           SORT SORT-F
               ON ASCENDING KEY SR-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ SR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
               INPUT PROCEDURE IS ｾﾝﾀｸ-RTN
               OUTPUT PROCEDURE IS ﾎｳｺｸ-RTN
           PERFORM ｼｭｳﾘｮｳ-RTN
           MOVE 0 TO RETURN-CODE
           STOP RUN.

      *----------------------------------------------------------------
      * ｼｮｷ-RTN: 業務日付、ファイル、代理店マスタ、制御レコード。
      *----------------------------------------------------------------
       ｼｮｷ-RTN.
           ACCEPT WS-ｺﾏﾝﾄﾞｷﾞｮｳ FROM COMMAND-LINE
           MOVE WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8) TO WS-ｷﾞｮｳﾑﾋﾂﾞｹ
           CALL "GLC020" USING WS-ﾎｳｺｳ WS-ｷﾞｮｳﾑﾋﾂﾞｹ
                               WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ WS-ﾍﾝｶﾝｹｯｶ
           IF WS-ﾍﾝｶﾝｹｯｶ NOT = 0
               DISPLAY "GLM030 日付不正: " WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8)
               MOVE 8 TO RETURN-CODE
               STOP RUN
           END-IF
           OPEN INPUT POLMST-F AGTMST-F
           OPEN I-O GLSCTL-F
           OPEN EXTEND GLJRNL-F
           OPEN OUTPUT COMMIS-F REPORT-F
           EVALUATE TRUE
               WHEN NOT WS-POLMST-OK
                   MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-AGTMST-OK
                   MOVE "AGTMST" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-AGTMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-GLSCTL-OK
                   MOVE "GLSCTL" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-GLSCTL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-GLJRNL-OK
                   MOVE "GLJRNL" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-GLJRNL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-COMMIS-OK
                   MOVE "COMMIS" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-COMMIS-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-REPORT-OK
                   MOVE "GLM030R" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-REPORT-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN OTHER
                   MOVE SPACES TO WS-ｴﾗｰﾌｧｲﾙ
           END-EVALUATE
           IF WS-ｴﾗｰﾌｧｲﾙ NOT = SPACES
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           MOVE SPACES TO JR-ｼﾜｹ-ﾚｺｰﾄﾞ
           PERFORM ﾀﾞｲﾘﾃﾝ-ﾖﾐｺﾐ-RTN
           READ GLSCTL-F
           IF NOT WS-GLSCTL-OK
               MOVE "GLSCTL" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-GLSCTL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ﾀﾞｲﾘﾃﾝ-ﾖﾐｺﾐ-RTN: 代理店マスタをテーブルに読み込む。
      *----------------------------------------------------------------
       ﾀﾞｲﾘﾃﾝ-ﾖﾐｺﾐ-RTN.
           READ AGTMST-F
           PERFORM UNTIL NOT WS-AGTMST-OK
               ADD 1 TO WS-ﾀﾞｲﾘﾃﾝｹﾝｽｳ
               SET WS-ﾀﾞｲﾘﾃﾝ-IX TO WS-ﾀﾞｲﾘﾃﾝｹﾝｽｳ
               MOVE AG-ｺｰﾄﾞ TO WS-ﾀﾞｲﾘﾃﾝ-ｺｰﾄﾞ(WS-ﾀﾞｲﾘﾃﾝ-IX)
               MOVE AG-ｼﾒｲ TO WS-ﾀﾞｲﾘﾃﾝ-ｼﾒｲ(WS-ﾀﾞｲﾘﾃﾝ-IX)
               MOVE AG-ｼﾃﾝ TO WS-ﾀﾞｲﾘﾃﾝ-ｼﾃﾝ(WS-ﾀﾞｲﾘﾃﾝ-IX)
               MOVE AG-ﾆｭｳｼｬﾈﾝｹﾞﾂ TO WS-ﾀﾞｲﾘﾃﾝ-ﾆｭｳｼｬ(WS-ﾀﾞｲﾘﾃﾝ-IX)
               READ AGTMST-F
           END-PERFORM
           IF NOT WS-AGTMST-ｵﾜﾘ
               MOVE "AGTMST" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-AGTMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ｾﾝﾀｸ-RTN: 当月入金のあった契約の手数料を求め、整列に渡す。
      *----------------------------------------------------------------
       ｾﾝﾀｸ-RTN.
           READ POLMST-F NEXT
           PERFORM UNTIL NOT WS-POLMST-OK
               IF PM-ｻｲｼｭｳﾆｭｳｷﾝﾈﾝｹﾞﾂ = WS-ｷﾞｮｳﾑ-YYYYMM
                       AND PM-ｻｲｼｭｳﾆｭｳｷﾝｶﾞｸ > 0
                   PERFORM ﾃｽｳﾘｮｳ-RTN
               END-IF
               READ POLMST-F NEXT
           END-PERFORM
           IF NOT WS-POLMST-ｵﾜﾘ
               MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ﾃｽｳﾘｮｳ-RTN: 経過年数による料率と手数料。
      *----------------------------------------------------------------
       ﾃｽｳﾘｮｳ-RTN.
           EVALUATE TRUE
               WHEN PM-ｹｲｶﾈﾝｽｳ = 1
                   MOVE 40 TO SR-ﾘﾂ
               WHEN PM-ｹｲｶﾈﾝｽｳ < 4
                   MOVE 7 TO SR-ﾘﾂ
               WHEN PM-ｹｲｶﾈﾝｽｳ < 11
                   MOVE 2 TO SR-ﾘﾂ
               WHEN OTHER
                   MOVE 0 TO SR-ﾘﾂ
           END-EVALUATE
           COMPUTE SR-ﾃｽｳﾘｮｳ = PM-ｻｲｼｭｳﾆｭｳｷﾝｶﾞｸ * SR-ﾘﾂ / 100
           MOVE PM-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ TO SR-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ
           MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO SR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           MOVE PM-ｼｮｳﾋﾝ TO SR-ｼｮｳﾋﾝ
           MOVE PM-ｹｲｶﾈﾝｽｳ TO SR-ｹｲｶﾈﾝｽｳ
           MOVE PM-ｻｲｼｭｳﾆｭｳｷﾝｶﾞｸ TO SR-ﾆｭｳｷﾝｶﾞｸ
           RELEASE SR-ｾｲﾚﾂ-ﾚｺｰﾄﾞ
           ADD 1 TO WS-ﾀｲｼｮｳｹﾝｽｳ
           ADD SR-ﾃｽｳﾘｮｳ TO WS-ﾃｽｳﾘｮｳｺﾞｳｹｲ.

      *----------------------------------------------------------------
      * ﾎｳｺｸ-RTN: 整列済みの1件ごとに明細を生成し、代理店が変わる
      * ごとに合計を書き出す。
      *----------------------------------------------------------------
       ﾎｳｺｸ-RTN.
           INITIATE COMMISSION-RPT
           MOVE SPACES TO WS-ｿｳﾀﾞｲﾘﾃﾝ
           RETURN SORT-F
               AT END
                   CONTINUE
               NOT AT END
                   PERFORM ﾒｲｻｲ-RTN
                   PERFORM UNTIL SR-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ = SPACES
                       RETURN SORT-F
                           AT END
                               MOVE SPACES TO SR-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ
                           NOT AT END
                               PERFORM ﾒｲｻｲ-RTN
                       END-RETURN
                   END-PERFORM
           END-RETURN
           IF WS-ｿｳﾀﾞｲﾘﾃﾝ NOT = SPACES
               PERFORM ﾀﾞｲﾘﾃﾝ-ｹｲ-RTN
           END-IF
           TERMINATE COMMISSION-RPT.

      *----------------------------------------------------------------
      * ﾒｲｻｲ-RTN: 代理店が変わったら前の代理店の計と見出しの準備、
      * 1行生成。
      *----------------------------------------------------------------
       ﾒｲｻｲ-RTN.
           IF SR-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ NOT = WS-ｿｳﾀﾞｲﾘﾃﾝ
               IF WS-ｿｳﾀﾞｲﾘﾃﾝ NOT = SPACES
                   PERFORM ﾀﾞｲﾘﾃﾝ-ｹｲ-RTN
               END-IF
               MOVE SR-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ TO WS-ｿｳﾀﾞｲﾘﾃﾝ
               PERFORM ﾐﾀﾞｼ-RTN
           END-IF
           GENERATE COMMISSION-DETAIL
           ADD 1 TO WS-ﾀﾞｲﾘﾃﾝ-ｹﾝｽｳ
           ADD SR-ﾆｭｳｷﾝｶﾞｸ TO WS-ﾀﾞｲﾘﾃﾝ-ﾆｭｳｷﾝ
           ADD SR-ﾃｽｳﾘｮｳ TO WS-ﾀﾞｲﾘﾃﾝ-ﾃｽｳﾘｮｳ.

      *----------------------------------------------------------------
      * ﾐﾀﾞｼ-RTN: 代理店マスタから見出しの名称、支社、入社年月。
      *----------------------------------------------------------------
       ﾐﾀﾞｼ-RTN.
           SET WS-ﾀﾞｲﾘﾃﾝ-IX TO 1
           SEARCH WS-ﾀﾞｲﾘﾃﾝ
               AT END
                   MOVE "(未登録)" TO WS-ﾐﾀﾞｼ-ｼﾒｲ
                   MOVE SPACES TO WS-ﾐﾀﾞｼ-ｼﾃﾝ WS-ﾐﾀﾞｼ-ﾆｭｳｼｬ
               WHEN WS-ﾀﾞｲﾘﾃﾝ-ｺｰﾄﾞ(WS-ﾀﾞｲﾘﾃﾝ-IX) = SR-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ
                   MOVE WS-ﾀﾞｲﾘﾃﾝ-ｼﾒｲ(WS-ﾀﾞｲﾘﾃﾝ-IX) TO WS-ﾐﾀﾞｼ-ｼﾒｲ
                   MOVE WS-ﾀﾞｲﾘﾃﾝ-ｼﾃﾝ(WS-ﾀﾞｲﾘﾃﾝ-IX) TO WS-ﾐﾀﾞｼ-ｼﾃﾝ
                   MOVE WS-ﾀﾞｲﾘﾃﾝ-ﾆｭｳｼｬ(WS-ﾀﾞｲﾘﾃﾝ-IX)(1:2)
                       TO WS-ﾐﾀﾞｼ-ﾆｭｳｼｬ(1:2)
                   MOVE "/" TO WS-ﾐﾀﾞｼ-ﾆｭｳｼｬ(3:1)
                   MOVE WS-ﾀﾞｲﾘﾃﾝ-ﾆｭｳｼｬ(WS-ﾀﾞｲﾘﾃﾝ-IX)(3:2)
                       TO WS-ﾐﾀﾞｼ-ﾆｭｳｼｬ(4:2)
           END-SEARCH.

      *----------------------------------------------------------------
      * ﾀﾞｲﾘﾃﾝ-ｹｲ-RTN: 代理店の合計を手数料ファイルに書き、仕訳し、
      * 数え直す。
      *----------------------------------------------------------------
       ﾀﾞｲﾘﾃﾝ-ｹｲ-RTN.
           MOVE WS-ｷﾞｮｳﾑ-YYYYMM TO CM-ﾈﾝｹﾞﾂ
           MOVE WS-ｿｳﾀﾞｲﾘﾃﾝ TO CM-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ
           MOVE WS-ﾀﾞｲﾘﾃﾝ-ｹﾝｽｳ TO CM-ｹﾝｽｳ
           MOVE WS-ﾀﾞｲﾘﾃﾝ-ﾆｭｳｷﾝ TO CM-ﾆｭｳｷﾝｶﾞｸ
           MOVE WS-ﾀﾞｲﾘﾃﾝ-ﾃｽｳﾘｮｳ TO CM-ﾃｽｳﾘｮｳ
           WRITE CM-ﾃｽｳﾘｮｳ-ﾚｺｰﾄﾞ
           IF NOT WS-COMMIS-OK
               MOVE "COMMIS" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-COMMIS-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           IF WS-ﾀﾞｲﾘﾃﾝ-ﾃｽｳﾘｮｳ > 0
               MOVE CT-ｼﾞｶｲｼﾜｹﾊﾞﾝｺﾞｳ TO JR-ｼﾜｹﾊﾞﾝｺﾞｳ
               ADD 1 TO CT-ｼﾞｶｲｼﾜｹﾊﾞﾝｺﾞｳ
               MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO JR-ﾋﾂﾞｹ
               MOVE "GLM030" TO JR-ﾌﾟﾛｸﾞﾗﾑ
               MOVE 5410 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
               SET JR-ｶﾘｶﾀ TO TRUE
               MOVE WS-ﾀﾞｲﾘﾃﾝ-ﾃｽｳﾘｮｳ TO JR-ｷﾝｶﾞｸ
               MOVE WS-ｿｳﾀﾞｲﾘﾃﾝ TO JR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
               MOVE "代理店手数料" TO JR-ﾃｷﾖｳ
               WRITE JR-ｼﾜｹ-ﾚｺｰﾄﾞ
               IF NOT WS-GLJRNL-OK
                   MOVE "GLJRNL" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-GLJRNL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
                   PERFORM ｲｼﾞｮｳ-RTN
               END-IF
           END-IF
           MOVE 0 TO WS-ﾀﾞｲﾘﾃﾝ-ｹﾝｽｳ WS-ﾀﾞｲﾘﾃﾝ-ﾆｭｳｷﾝ WS-ﾀﾞｲﾘﾃﾝ-ﾃｽｳﾘｮｳ.

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
           CLOSE POLMST-F AGTMST-F GLSCTL-F GLJRNL-F COMMIS-F
                 REPORT-F
           DISPLAY "GLM030 対象 " WS-ﾀｲｼｮｳｹﾝｽｳ " 手数料 "
                   WS-ﾃｽｳﾘｮｳｺﾞｳｹｲ.

      *----------------------------------------------------------------
      * ｲｼﾞｮｳ-RTN: ファイル異常は 16 で終わる。
      *----------------------------------------------------------------
       ｲｼﾞｮｳ-RTN.
           DISPLAY "GLM030 ファイル異常 " WS-ｴﾗｰﾌｧｲﾙ
                   " 状態 " WS-ｴﾗｰｽﾃｰﾀｽ
           MOVE 16 TO RETURN-CODE
           STOP RUN.
