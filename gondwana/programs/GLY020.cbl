      *----------------------------------------------------------------
      * GLY020: ご契約内容のお知らせ
      *
      * 年次、契約者配当 (GLY010) の後。保障中 (有効・猶予中) の契約
      * 1件につき1頁のお知らせを、レポートライタで作る。頁の見出しに
      * 契約者の住所と氏名、契約の内容 (証券番号、契約日、払込方法、
      * 払込済年月、保険料、配当金累計、貸付残高、解約返戻金) を置き、
      * 明細に有効な保障 (保障マスタ COVMST の主契約と特約) を1行
      * ずつ、最後に保障の保険料の合計を置く。解約返戻金は GLC030。
      * お知らせは印刷会社に渡す NOTICE に書く。
      *
      * 復帰コード: 0 作成済、8 日付が日付でない、16 ファイル異常。
      *
      * 変更履歴
      * 日付       担当 変更
      * 2026-11-19 YN   GL-0827 初版。
      *----------------------------------------------------------------
       IDENTIFICATION DIVISION.
       PROGRAM-ID. GLY020.
       AUTHOR. ゴンドワナ生命 システム部.
       DATE-WRITTEN. 2026-11-19.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT POLMST-F ASSIGN TO "POLMST"
               ORGANIZATION IS INDEXED
               ACCESS MODE IS SEQUENTIAL
               RECORD KEY IS PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
               FILE STATUS IS WS-POLMST-ST.
           SELECT COVMST-F ASSIGN TO "COVMST"
               ORGANIZATION IS INDEXED
               ACCESS MODE IS DYNAMIC
               RECORD KEY IS HS-ｷｰ
               FILE STATUS IS WS-COVMST-ST.
           SELECT PRODTB-F ASSIGN TO "PRODTB"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-PRODTB-ST.
           SELECT NOTICE-F ASSIGN TO "NOTICE"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-NOTICE-ST.

       DATA DIVISION.
       FILE SECTION.
       FD  POLMST-F.
       COPY POLMST.
       FD  COVMST-F.
       COPY COVMST.
       FD  PRODTB-F.
       COPY PRODTB.
       FD  NOTICE-F
           REPORT IS NOTICE-RPT.

       WORKING-STORAGE SECTION.
       01  WS-ｽﾃｰﾀｽ.
           05  WS-POLMST-ST            PIC XX.
               88  WS-POLMST-OK                VALUE "00".
               88  WS-POLMST-ｵﾜﾘ               VALUE "10".
           05  WS-COVMST-ST            PIC XX.
               88  WS-COVMST-OK                VALUE "00".
               88  WS-COVMST-ｵﾜﾘ               VALUE "10" "23".
           05  WS-PRODTB-ST            PIC XX.
               88  WS-PRODTB-OK                VALUE "00".
               88  WS-PRODTB-ｵﾜﾘ               VALUE "10".
           05  WS-NOTICE-ST            PIC XX.
               88  WS-NOTICE-OK                VALUE "00".
           05  WS-CSVTB-ST             PIC XX.
               88  WS-CSVTB-OK                 VALUE "00".
       01  WS-ｺﾏﾝﾄﾞｷﾞｮｳ                PIC X(80).
       01  WS-ｴﾗｰﾌｧｲﾙ                  PIC X(8).
       01  WS-ｴﾗｰｽﾃｰﾀｽ                 PIC XX.
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ                PIC 9(8).
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ            PIC X(7).
       01  WS-ﾎｳｺｳ                     PIC X.
       01  WS-ﾍﾝｶﾝｹｯｶ                  PIC 9.
       01  WS-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ              PIC X(10).
       01  WS-ﾊﾗｲｺﾐﾒｲ                  PIC X(6).
       01  WS-ﾍﾝﾚｲｷﾝ                   PIC 9(9).
       01  WS-ｼｮｳﾋﾝﾒｲ                  PIC X(20).
       01  WS-ﾂｳﾁｹﾝｽｳ                  PIC 9(7) VALUE 0.
       COPY TBLWS.

       REPORT SECTION.
       RD  NOTICE-RPT
           CONTROL IS WS-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           PAGE LIMIT IS 30 LINES
           HEADING 1
           FIRST DETAIL 3
           LAST DETAIL 26
           FOOTING 28.
       01  TYPE IS PAGE HEADING.
           05  LINE 1.
               10  COLUMN 1   PIC X(36)
                   VALUE "ゴンドワナ生命  ご契約内容のお知らせ".
               10  COLUMN 40  PIC X(7)  SOURCE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ.
               10  COLUMN 48  PIC X(4)  VALUE "現在".
       01  TYPE IS CONTROL HEADING WS-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ.
           05  LINE PLUS 1.
               10  COLUMN 1   PIC X(40) SOURCE PM-ｼﾞｭｳｼｮ.
           05  LINE PLUS 1.
               10  COLUMN 1   PIC X(20) SOURCE PM-ｼﾒｲ.
               10  COLUMN 22  PIC X(2)  VALUE "様".
           05  LINE PLUS 2.
               10  COLUMN 1   PIC X(8)  VALUE "証券番号".
               10  COLUMN 12  PIC X(10) SOURCE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ.
               10  COLUMN 26  PIC X(6)  VALUE "契約日".
               10  COLUMN 36  PIC X(7)  SOURCE PM-ｹｲﾔｸﾋﾞ.
               10  COLUMN 46  PIC X(8)  VALUE "払込方法".
               10  COLUMN 56  PIC X(6)  SOURCE WS-ﾊﾗｲｺﾐﾒｲ.
           05  LINE PLUS 1.
               10  COLUMN 1   PIC X(10) VALUE "払込済年月".
               10  COLUMN 12  PIC 9(6)  SOURCE PM-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ.
               10  COLUMN 26  PIC X(6)  VALUE "保険料".
               10  COLUMN 34  PIC Z,ZZZ,ZZ9 SOURCE PM-ﾎｹﾝﾘｮｳ.
           05  LINE PLUS 1.
               10  COLUMN 1   PIC X(10) VALUE "配当金累計".
               10  COLUMN 12  PIC ZZZ,ZZZ,ZZ9 SOURCE PM-ﾊｲﾄｳﾙｲｹｲ.
               10  COLUMN 26  PIC X(8)  VALUE "貸付残高".
               10  COLUMN 36  PIC ZZZ,ZZZ,ZZ9 SOURCE PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ.
           05  LINE PLUS 1.
               10  COLUMN 1   PIC X(10) VALUE "解約返戻金".
               10  COLUMN 12  PIC ZZZ,ZZZ,ZZ9 SOURCE WS-ﾍﾝﾚｲｷﾝ.
           05  LINE PLUS 2.
               10  COLUMN 1   PIC X(4)  VALUE "保障".
               10  COLUMN 6   PIC X(4)  VALUE "商品".
               10  COLUMN 33  PIC X(8)  VALUE "保険金額".
               10  COLUMN 47  PIC X(6)  VALUE "保険料".
               10  COLUMN 56  PIC X(6)  VALUE "満期日".
       01  NOTICE-DETAIL TYPE IS DETAIL.
           05  LINE PLUS 1.
               10  COLUMN 2   PIC 99    SOURCE HS-ﾎｼｮｳﾊﾞﾝｺﾞｳ.
               10  COLUMN 6   PIC X(2)  SOURCE HS-ｼｮｳﾋﾝ.
               10  COLUMN 9   PIC X(20) SOURCE WS-ｼｮｳﾋﾝﾒｲ.
               10  COLUMN 30  PIC ZZZ,ZZZ,ZZ9 SOURCE HS-ﾎｹﾝｷﾝｶﾞｸ.
               10  COLUMN 44  PIC Z,ZZZ,ZZ9 SOURCE HS-ﾎｹﾝﾘｮｳ.
               10  COLUMN 56  PIC X(7)  SOURCE HS-ﾏﾝｷﾋﾞ.
       01  TYPE IS CONTROL FOOTING WS-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           NEXT GROUP IS NEXT PAGE.
           05  LINE PLUS 1.
               10  COLUMN 1   PIC X(10) VALUE "保険料合計".
               10  COLUMN 42  PIC ZZZ,ZZZ,ZZ9 SUM HS-ﾎｹﾝﾘｮｳ.

       PROCEDURE DIVISION.
      *----------------------------------------------------------------
      * ﾒｲﾝ-RTN: 準備、契約ごとのお知らせ、終了。
      *----------------------------------------------------------------
       ﾒｲﾝ-RTN.
           PERFORM ｼｮｷ-RTN
           INITIATE NOTICE-RPT
           PERFORM ｹｲﾔｸ-RTN UNTIL WS-POLMST-ｵﾜﾘ
           TERMINATE NOTICE-RPT
           PERFORM ｼｭｳﾘｮｳ-RTN
           MOVE 0 TO RETURN-CODE
           STOP RUN.

      *----------------------------------------------------------------
      * ｼｮｷ-RTN: 業務日付、ファイル、商品。
      *----------------------------------------------------------------
       ｼｮｷ-RTN.
           ACCEPT WS-ｺﾏﾝﾄﾞｷﾞｮｳ FROM COMMAND-LINE
           MOVE WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8) TO WS-ｷﾞｮｳﾑﾋﾂﾞｹ
           MOVE "W" TO WS-ﾎｳｺｳ
           CALL "GLC020" USING WS-ﾎｳｺｳ WS-ｷﾞｮｳﾑﾋﾂﾞｹ
                               WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ WS-ﾍﾝｶﾝｹｯｶ
           IF WS-ﾍﾝｶﾝｹｯｶ NOT = 0
               DISPLAY "GLY020 日付不正: " WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8)
               MOVE 8 TO RETURN-CODE
               STOP RUN
           END-IF
           OPEN INPUT POLMST-F COVMST-F PRODTB-F
           OPEN OUTPUT NOTICE-F
           EVALUATE TRUE
               WHEN NOT WS-POLMST-OK
                   MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-COVMST-OK
                   MOVE "COVMST" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-COVMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-PRODTB-OK
                   MOVE "PRODTB" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-PRODTB-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-NOTICE-OK
                   MOVE "NOTICE" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-NOTICE-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN OTHER
                   MOVE SPACES TO WS-ｴﾗｰﾌｧｲﾙ
           END-EVALUATE
           IF WS-ｴﾗｰﾌｧｲﾙ NOT = SPACES
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           PERFORM ｼｮｳﾋﾝ-ﾖﾐｺﾐ-RTN
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
      * ｹｲﾔｸ-RTN: 保障中の契約の見出しの値を整え、保障をたどる。
      *----------------------------------------------------------------
       ｹｲﾔｸ-RTN.
           IF PM-ﾎｼｮｳﾁｭｳ
               MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO WS-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
               EVALUATE TRUE
                   WHEN PM-ﾈﾝﾊﾞﾗｲ
                       MOVE "年払" TO WS-ﾊﾗｲｺﾐﾒｲ
                   WHEN PM-ﾊﾝﾄｼﾊﾞﾗｲ
                       MOVE "半年払" TO WS-ﾊﾗｲｺﾐﾒｲ
                   WHEN OTHER
                       MOVE "月払" TO WS-ﾊﾗｲｺﾐﾒｲ
               END-EVALUATE
               CALL "GLC030" USING PM-ｼｮｳﾋﾝ PM-ｶﾆｭｳﾈﾝﾚｲ PM-ｹｲｶﾈﾝｽｳ
                                   PM-ﾎｹﾝｷﾝｶﾞｸ WS-ﾍﾝﾚｲｷﾝ WS-CSVTB-ST
               IF NOT WS-CSVTB-OK
                   MOVE "CSVTB" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-CSVTB-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
                   PERFORM ｲｼﾞｮｳ-RTN
               END-IF
               PERFORM ﾎｼｮｳ-ｿｳｻ-RTN
               ADD 1 TO WS-ﾂｳﾁｹﾝｽｳ
           END-IF
           PERFORM POLMST-ﾖﾐｺﾐ-RTN.

      *----------------------------------------------------------------
      * ﾎｼｮｳ-ｿｳｻ-RTN: 契約の保障を保障番号順にたどり、有効な保障を
      * 1行ずつ出す。
      *----------------------------------------------------------------
       ﾎｼｮｳ-ｿｳｻ-RTN.
           MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO HS-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           MOVE 0 TO HS-ﾎｼｮｳﾊﾞﾝｺﾞｳ
           START COVMST-F KEY NOT < HS-ｷｰ
           IF WS-COVMST-OK
               READ COVMST-F NEXT
           END-IF
           PERFORM UNTIL NOT WS-COVMST-OK
                   OR HS-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ NOT = PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
               IF HS-ﾕｳｺｳ
                   PERFORM ﾒｲｻｲ-RTN
               END-IF
               READ COVMST-F NEXT
           END-PERFORM
           IF NOT WS-COVMST-OK AND NOT WS-COVMST-ｵﾜﾘ
               MOVE "COVMST" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-COVMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ﾒｲｻｲ-RTN: 保障1件の商品名を引いて明細を出す。
      *----------------------------------------------------------------
       ﾒｲｻｲ-RTN.
           MOVE SPACES TO WS-ｼｮｳﾋﾝﾒｲ
           SET WS-ｼｮｳﾋﾝ-IX TO 1
           SEARCH WS-ｼｮｳﾋﾝ
               WHEN WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = HS-ｼｮｳﾋﾝ
                   MOVE WS-ｼｮｳﾋﾝ-ﾒｲｼｮｳ(WS-ｼｮｳﾋﾝ-IX) TO WS-ｼｮｳﾋﾝﾒｲ
           END-SEARCH
           GENERATE NOTICE-DETAIL.

      *----------------------------------------------------------------
      * ｼｭｳﾘｮｳ-RTN: ファイルを閉じ、件数を表示する。
      *----------------------------------------------------------------
       ｼｭｳﾘｮｳ-RTN.
           CLOSE POLMST-F COVMST-F PRODTB-F NOTICE-F
           DISPLAY "GLY020 お知らせ " WS-ﾂｳﾁｹﾝｽｳ.

      *----------------------------------------------------------------
      * ｲｼﾞｮｳ-RTN: ファイル異常は 16 で終わる。
      *----------------------------------------------------------------
       ｲｼﾞｮｳ-RTN.
           DISPLAY "GLY020 ファイル異常 " WS-ｴﾗｰﾌｧｲﾙ
                   " 状態 " WS-ｴﾗｰｽﾃｰﾀｽ
           MOVE 16 TO RETURN-CODE
           STOP RUN.

       COPY PRDLOD.
