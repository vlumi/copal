      *----------------------------------------------------------------
      * GLY010: 契約者配当
      *
      * 年次、年度末 (3月末)。保障中 (有効・猶予中) で経過年数が
      * 3年以上の契約に、商品の配当率を保険金額に掛けた配当金を
      * 割り当て、配当累計に積み立てる (円未満切捨て)。配当率のない
      * 商品 (定期、医療) には配当しない。積み立てた配当は満期・死亡
      * のとき保険金に加えて支払う。仕訳は契約者配当金 (借方)。
      *
      * 復帰コード: 0 処理済、8 日付が日付でない、16 ファイル異常。
      *
      * 変更履歴
      * 日付       担当 変更
      * 1988-03-07 KT   GL-0021 初版、昭和62年度決算から。
      *----------------------------------------------------------------
       IDENTIFICATION DIVISION.
       PROGRAM-ID. GLY010.
       AUTHOR. ゴンドワナ生命 システム部.
       DATE-WRITTEN. 1988-03-07.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT POLMST-F ASSIGN TO "POLMST"
               ORGANIZATION IS INDEXED
               ACCESS MODE IS DYNAMIC
               RECORD KEY IS PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
               FILE STATUS IS WS-POLMST-ST.
           SELECT PRODTB-F ASSIGN TO "PRODTB"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-PRODTB-ST.
           SELECT GLSCTL-F ASSIGN TO "GLSCTL"
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-GLSCTL-ST.
           SELECT OPTIONAL GLJRNL-F ASSIGN TO "GLJRNL"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-GLJRNL-ST.
           SELECT REPORT-F ASSIGN TO "GLY010R"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-ST.

       DATA DIVISION.
       FILE SECTION.
       FD  POLMST-F.
       COPY POLMST.
       FD  PRODTB-F.
       COPY PRODTB.
       FD  GLSCTL-F
           RECORD CONTAINS 64 CHARACTERS.
       COPY CTLREC.
       FD  GLJRNL-F.
       COPY GLJRNL.
       FD  REPORT-F.
       01  ﾎｳｺｸ-ｷﾞｮｳ                   PIC X(100).

       WORKING-STORAGE SECTION.
       01  WS-ｽﾃｰﾀｽ.
           05  WS-POLMST-ST            PIC XX.
               88  WS-POLMST-OK                VALUE "00".
               88  WS-POLMST-ｵﾜﾘ               VALUE "10".
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
       01  WS-ﾊｲﾄｳｷﾝ                   PIC 9(9).
       COPY TBLWS.
       01  WS-ｹﾝｽｳ.
           05  WS-ﾀｲｼｮｳｹﾝｽｳ            PIC 9(7) VALUE 0.
           05  WS-ﾊｲﾄｳｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ﾊｲﾄｳｺﾞｳｹｲ            PIC 9(13) VALUE 0.
       01  WS-ﾎｳｺｸ.
           05  WS-ﾋｮｳﾀﾞｲ-ｷﾞｮｳ.
               10  FILLER                  PIC X(28)
                   VALUE "ゴンドワナ生命  契約者配当".
               10  FILLER                  PIC X(4) VALUE SPACES.
               10  WS-ﾋｮｳﾀﾞｲ-ﾋﾂﾞｹ          PIC X(7).
           05  WS-ﾒｲｻｲ-ｷﾞｮｳ.
               10  WS-ﾒｲｻｲ-ｼｮｳｹﾝ           PIC X(10).
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾒｲｻｲ-ｼｮｳﾋﾝ           PIC X(2).
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾒｲｻｲ-ｼﾒｲ             PIC X(20).
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾒｲｻｲ-ﾈﾝｽｳ            PIC Z9.
               10  FILLER                  PIC X(4) VALUE "年目".
               10  WS-ﾒｲｻｲ-ｷﾝｶﾞｸ           PIC ZZZ,ZZZ,ZZ9.
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  FILLER                  PIC X(4) VALUE "累計".
               10  WS-ﾒｲｻｲ-ﾙｲｹｲ            PIC ZZZ,ZZZ,ZZ9.
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
      * ｼｮｷ-RTN: 業務日付、ファイル、商品テーブル、制御レコード。
      *----------------------------------------------------------------
       ｼｮｷ-RTN.
           ACCEPT WS-ｺﾏﾝﾄﾞｷﾞｮｳ FROM COMMAND-LINE
           MOVE WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8) TO WS-ｷﾞｮｳﾑﾋﾂﾞｹ
           CALL "GLC020" USING WS-ﾎｳｺｳ WS-ｷﾞｮｳﾑﾋﾂﾞｹ
                               WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ WS-ﾍﾝｶﾝｹｯｶ
           IF WS-ﾍﾝｶﾝｹｯｶ NOT = 0
               DISPLAY "GLY010 日付不正: " WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8)
               MOVE 8 TO RETURN-CODE
               STOP RUN
           END-IF
           OPEN INPUT PRODTB-F
           OPEN I-O POLMST-F GLSCTL-F
           OPEN EXTEND GLJRNL-F
           OPEN OUTPUT REPORT-F
           EVALUATE TRUE
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
                   MOVE "GLY010R" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-REPORT-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN OTHER
                   MOVE SPACES TO WS-ｴﾗｰﾌｧｲﾙ
           END-EVALUATE
           IF WS-ｴﾗｰﾌｧｲﾙ NOT = SPACES
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           MOVE SPACES TO JR-ｼﾜｹ-ﾚｺｰﾄﾞ
           PERFORM ｼｮｳﾋﾝ-ﾖﾐｺﾐ-RTN
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
      * ｼｮﾘ-RTN: 保障中で3年目以上なら配当を求め、あれば積み立てる。
      *----------------------------------------------------------------
       ｼｮﾘ-RTN.
           IF PM-ﾎｼｮｳﾁｭｳ AND PM-ｹｲｶﾈﾝｽｳ NOT < 3
               ADD 1 TO WS-ﾀｲｼｮｳｹﾝｽｳ
               PERFORM ﾊｲﾄｳ-RTN
               IF WS-ﾊｲﾄｳｷﾝ > 0
                   PERFORM ﾂﾐﾀﾃ-RTN
               END-IF
           END-IF
           PERFORM POLMST-ﾖﾐｺﾐ-RTN.

      *----------------------------------------------------------------
      * ﾊｲﾄｳ-RTN: 商品の配当率を保険金額に掛ける。
      *----------------------------------------------------------------
       ﾊｲﾄｳ-RTN.
           MOVE 0 TO WS-ﾊｲﾄｳｷﾝ
           SET WS-ｼｮｳﾋﾝ-IX TO 1
           SEARCH WS-ｼｮｳﾋﾝ
               WHEN WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = PM-ｼｮｳﾋﾝ
                   COMPUTE WS-ﾊｲﾄｳｷﾝ = PM-ﾎｹﾝｷﾝｶﾞｸ
                       * WS-ｼｮｳﾋﾝ-ﾊｲﾄｳﾘﾂ(WS-ｼｮｳﾋﾝ-IX)
           END-SEARCH.

      *----------------------------------------------------------------
      * ﾂﾐﾀﾃ-RTN: 配当累計に積み立て、仕訳し、明細を報告する。
      *----------------------------------------------------------------
       ﾂﾐﾀﾃ-RTN.
           ADD WS-ﾊｲﾄｳｷﾝ TO PM-ﾊｲﾄｳﾙｲｹｲ
           REWRITE PM-ｹｲﾔｸ-ﾚｺｰﾄﾞ
           IF NOT WS-POLMST-OK
               MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           MOVE CT-ｼﾞｶｲｼﾜｹﾊﾞﾝｺﾞｳ TO JR-ｼﾜｹﾊﾞﾝｺﾞｳ
           ADD 1 TO CT-ｼﾞｶｲｼﾜｹﾊﾞﾝｺﾞｳ
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO JR-ﾋﾂﾞｹ
           MOVE "GLY010" TO JR-ﾌﾟﾛｸﾞﾗﾑ
           MOVE 5310 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
           SET JR-ｶﾘｶﾀ TO TRUE
           MOVE WS-ﾊｲﾄｳｷﾝ TO JR-ｷﾝｶﾞｸ
           MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO JR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           MOVE "契約者配当金" TO JR-ﾃｷﾖｳ
           WRITE JR-ｼﾜｹ-ﾚｺｰﾄﾞ
           IF NOT WS-GLJRNL-OK
               MOVE "GLJRNL" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-GLJRNL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           ADD 1 TO WS-ﾊｲﾄｳｹﾝｽｳ
           ADD WS-ﾊｲﾄｳｷﾝ TO WS-ﾊｲﾄｳｺﾞｳｹｲ
           MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO WS-ﾒｲｻｲ-ｼｮｳｹﾝ
           MOVE PM-ｼｮｳﾋﾝ TO WS-ﾒｲｻｲ-ｼｮｳﾋﾝ
           MOVE PM-ｼﾒｲ TO WS-ﾒｲｻｲ-ｼﾒｲ
           MOVE PM-ｹｲｶﾈﾝｽｳ TO WS-ﾒｲｻｲ-ﾈﾝｽｳ
           MOVE WS-ﾊｲﾄｳｷﾝ TO WS-ﾒｲｻｲ-ｷﾝｶﾞｸ
           MOVE PM-ﾊｲﾄｳﾙｲｹｲ TO WS-ﾒｲｻｲ-ﾙｲｹｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ﾒｲｻｲ-ｷﾞｮｳ.

      *----------------------------------------------------------------
      * ﾎｳｺｸ-RTN: 件数と合計。
      *----------------------------------------------------------------
       ﾎｳｺｸ-RTN.
           MOVE SPACES TO ﾎｳｺｸ-ｷﾞｮｳ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ
           MOVE "対象件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾀｲｼｮｳｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "配当件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾊｲﾄｳｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "配当金合計" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾊｲﾄｳｺﾞｳｹｲ TO WS-ｹﾝｽｳ-ｱﾀｲ
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
           CLOSE POLMST-F PRODTB-F GLSCTL-F GLJRNL-F REPORT-F
           DISPLAY "GLY010 配当 " WS-ﾊｲﾄｳｹﾝｽｳ " 件 " WS-ﾊｲﾄｳｺﾞｳｹｲ.

      *----------------------------------------------------------------
      * ｲｼﾞｮｳ-RTN: ファイル異常は 16 で終わる。
      *----------------------------------------------------------------
       ｲｼﾞｮｳ-RTN.
           DISPLAY "GLY010 ファイル異常 " WS-ｴﾗｰﾌｧｲﾙ
                   " 状態 " WS-ｴﾗｰｽﾃｰﾀｽ
           MOVE 16 TO RETURN-CODE
           STOP RUN.

       COPY PRDLOD.
