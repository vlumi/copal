      *----------------------------------------------------------------
      * GLD020: 契約応当日処理
      *
      * 日次。保障中 (有効・猶予中) の契約のうち、契約日の月日が
      * 業務日付の月日と同じものの経過年数を1つ進める。その日が
      * 満期日でもあれば、養老保険は満期保険金 (保険金額に配当累計
      * を加えた額) を支払って満期とし、定期保険は更新する: その
      * 時点の保険年齢で保険料を引き直し、満期日を保険期間分
      * 先に送る。保険年齢が商品の最高年齢を超えていれば更新せず
      * 期間満了とする (支払なし)。終身・医療は経過年数のみ。
      * 支払番号と仕訳番号は制御レコードから採番する。
      *
      * 復帰コード: 0 処理済、8 日付が日付でない、16 ファイル異常。
      *
      * 変更履歴
      * 日付       担当 変更
      * 1987-04-01 KT   GL-0001 初版。
      * 1992-03-02 KT   GL-0117 定期保険の更新。
      * 2001-07-02 MS   GL-0271 満期保険金に配当累計を加える。
      *----------------------------------------------------------------
       IDENTIFICATION DIVISION.
       PROGRAM-ID. GLD020.
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
           SELECT PRODTB-F ASSIGN TO "PRODTB"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-PRODTB-ST.
           SELECT RATETB-F ASSIGN TO "RATETB"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-RATETB-ST.
           SELECT GLSCTL-F ASSIGN TO "GLSCTL"
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-GLSCTL-ST.
           SELECT OPTIONAL PAYOUT-F ASSIGN TO "PAYOUT"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-PAYOUT-ST.
           SELECT OPTIONAL GLJRNL-F ASSIGN TO "GLJRNL"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-GLJRNL-ST.
           SELECT REPORT-F ASSIGN TO "GLD020R"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-ST.

       DATA DIVISION.
       FILE SECTION.
       FD  POLMST-F.
       COPY POLMST.
       FD  PRODTB-F.
       COPY PRODTB.
       FD  RATETB-F.
       COPY RATETB.
       FD  GLSCTL-F
           RECORD CONTAINS 64 CHARACTERS.
       COPY CTLREC.
       FD  PAYOUT-F.
       COPY PAYOUT.
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
           05  WS-RATETB-ST            PIC XX.
               88  WS-RATETB-OK                VALUE "00".
               88  WS-RATETB-ｵﾜﾘ               VALUE "10".
           05  WS-GLSCTL-ST            PIC XX.
               88  WS-GLSCTL-OK                VALUE "00".
           05  WS-PAYOUT-ST            PIC XX.
               88  WS-PAYOUT-OK                VALUE "00" "05".
           05  WS-GLJRNL-ST            PIC XX.
               88  WS-GLJRNL-OK                VALUE "00" "05".
           05  WS-REPORT-ST            PIC XX.
               88  WS-REPORT-OK                VALUE "00".
       01  WS-ｺﾏﾝﾄﾞｷﾞｮｳ                PIC X(80).
       01  WS-ｴﾗｰﾌｧｲﾙ                  PIC X(8).
       01  WS-ｴﾗｰｽﾃｰﾀｽ                 PIC XX.
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ                PIC 9(8).
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾌﾞﾝｶｲ REDEFINES WS-ｷﾞｮｳﾑﾋﾂﾞｹ.
           05  WS-ｷﾞｮｳﾑ-YYYY           PIC 9(4).
           05  WS-ｷﾞｮｳﾑ-MMDD           PIC 9(4).
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ            PIC X(7).
       01  WS-ﾎｳｺｳ                     PIC X.
       01  WS-ﾍﾝｶﾝｹｯｶ                  PIC 9.
       01  WS-ｾｲﾈﾝｶﾞｯﾋﾟ                PIC 9(8).
       01  WS-ﾈﾝﾚｲ                     PIC 99.
       01  WS-ｹﾞﾂｶﾞｸ                   PIC 9(9).
       01  WS-ﾏﾝｷﾋﾞ                    PIC 9(8).
       01  WS-ﾏﾝｷﾋﾞ-ﾌﾞﾝｶｲ REDEFINES WS-ﾏﾝｷﾋﾞ.
           05  WS-ﾏﾝｷ-YYYY             PIC 9(4).
           05  WS-ﾏﾝｷ-MMDD             PIC 9(4).
       01  WS-ﾏﾝｷﾋﾞ-ﾜﾚｷ                PIC X(7).
       COPY TBLWS.
       01  WS-ｹﾝｽｳ.
           05  WS-ｵｳﾄｳｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ｺｳｼﾝｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ﾏﾝﾘｮｳｹﾝｽｳ            PIC 9(7) VALUE 0.
           05  WS-ﾏﾝｷｹﾝｽｳ              PIC 9(7) VALUE 0.
           05  WS-ﾏﾝｷｷﾝｶﾞｸ             PIC 9(13) VALUE 0.
       01  WS-ﾎｳｺｸ.
           05  WS-ﾋｮｳﾀﾞｲ-ｷﾞｮｳ.
               10  FILLER                  PIC X(30)
                   VALUE "ゴンドワナ生命  契約応当日処理".
               10  FILLER                  PIC X(2) VALUE SPACES.
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
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾒｲｻｲ-ｹｯｶ             PIC X(12).
               10  WS-ﾒｲｻｲ-ｷﾝｶﾞｸ           PIC ZZZ,ZZZ,ZZ9.
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
      * ｼｮｷ-RTN: 業務日付、ファイル、テーブル、制御レコード。
      *----------------------------------------------------------------
       ｼｮｷ-RTN.
           ACCEPT WS-ｺﾏﾝﾄﾞｷﾞｮｳ FROM COMMAND-LINE
           MOVE WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8) TO WS-ｷﾞｮｳﾑﾋﾂﾞｹ
           MOVE "W" TO WS-ﾎｳｺｳ
           CALL "GLC020" USING WS-ﾎｳｺｳ WS-ｷﾞｮｳﾑﾋﾂﾞｹ
                               WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ WS-ﾍﾝｶﾝｹｯｶ
           IF WS-ﾍﾝｶﾝｹｯｶ NOT = 0
               DISPLAY "GLD020 日付不正: " WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8)
               MOVE 8 TO RETURN-CODE
               STOP RUN
           END-IF
           OPEN INPUT PRODTB-F RATETB-F
           OPEN I-O POLMST-F GLSCTL-F
           OPEN EXTEND PAYOUT-F GLJRNL-F
           OPEN OUTPUT REPORT-F
           PERFORM ｵｰﾌﾟﾝ-ｹﾝｻ-RTN
           MOVE SPACES TO PY-ｼﾊﾗｲ-ﾚｺｰﾄﾞ JR-ｼﾜｹ-ﾚｺｰﾄﾞ
           PERFORM ｼｮｳﾋﾝ-ﾖﾐｺﾐ-RTN
           PERFORM ﾘｮｳﾘﾂ-ﾖﾐｺﾐ-RTN
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
      * ｵｰﾌﾟﾝ-ｹﾝｻ-RTN: 開けなかったファイルがあれば異常終了。
      *----------------------------------------------------------------
       ｵｰﾌﾟﾝ-ｹﾝｻ-RTN.
           EVALUATE TRUE
               WHEN NOT WS-PRODTB-OK
                   MOVE "PRODTB" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-PRODTB-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-RATETB-OK
                   MOVE "RATETB" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-RATETB-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
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
               WHEN NOT WS-REPORT-OK
                   MOVE "GLD020R" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-REPORT-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN OTHER
                   MOVE SPACES TO WS-ｴﾗｰﾌｧｲﾙ
           END-EVALUATE
           IF WS-ｴﾗｰﾌｧｲﾙ NOT = SPACES
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

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
      * ｼｮﾘ-RTN: 応当日の契約なら処理し、書き戻す。当日契約した
      * ものは応当日ではない。
      *----------------------------------------------------------------
       ｼｮﾘ-RTN.
           IF PM-ﾎｼｮｳﾁｭｳ
                   AND PM-ｹｲﾔｸﾋﾞ(4:4) = WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ(4:4)
                   AND PM-ｹｲﾔｸﾋﾞ NOT = WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ
               PERFORM ｵｳﾄｳ-RTN
               REWRITE PM-ｹｲﾔｸ-ﾚｺｰﾄﾞ
               IF NOT WS-POLMST-OK
                   MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
                   PERFORM ｲｼﾞｮｳ-RTN
               END-IF
           END-IF
           PERFORM POLMST-ﾖﾐｺﾐ-RTN.

      *----------------------------------------------------------------
      * ｵｳﾄｳ-RTN: 経過年数を進め、満期日なら商品により満期か更新。
      *----------------------------------------------------------------
       ｵｳﾄｳ-RTN.
           ADD 1 TO WS-ｵｳﾄｳｹﾝｽｳ
           ADD 1 TO PM-ｹｲｶﾈﾝｽｳ
           MOVE 0 TO WS-ﾒｲｻｲ-ｷﾝｶﾞｸ
           MOVE "年度更新" TO WS-ﾒｲｻｲ-ｹｯｶ
           EVALUATE TRUE
               WHEN PM-ﾖｳﾛｳ AND PM-ﾏﾝｷﾋﾞ = WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ
                   PERFORM ﾏﾝｷ-RTN
               WHEN PM-ﾃｲｷ AND PM-ﾏﾝｷﾋﾞ = WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ
                   PERFORM ｺｳｼﾝ-RTN
           END-EVALUATE
           PERFORM ﾒｲｻｲ-RTN.

      *----------------------------------------------------------------
      * ﾏﾝｷ-RTN: 満期保険金に配当累計を加えて支払い、満期とする。
      *----------------------------------------------------------------
       ﾏﾝｷ-RTN.
           ADD 1 TO WS-ﾏﾝｷｹﾝｽｳ
           MOVE CT-ｼﾞｶｲｼﾊﾗｲﾊﾞﾝｺﾞｳ TO PY-ｼﾊﾗｲﾊﾞﾝｺﾞｳ
           ADD 1 TO CT-ｼﾞｶｲｼﾊﾗｲﾊﾞﾝｺﾞｳ
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO PY-ｼﾊﾗｲﾋﾞ
           MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO PY-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           SET PY-ﾏﾝｷ TO TRUE
           COMPUTE PY-ｷﾝｶﾞｸ = PM-ﾎｹﾝｷﾝｶﾞｸ + PM-ﾊｲﾄｳﾙｲｹｲ
           MOVE PM-ｺｳｻﾞﾒｲｷﾞ TO PY-ｳｹﾄﾘﾆﾝｶﾅ
           MOVE SPACES TO PY-ｾｲｷｭｳﾊﾞﾝｺﾞｳ
           WRITE PY-ｼﾊﾗｲ-ﾚｺｰﾄﾞ
           MOVE 5140 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
           MOVE PM-ﾎｹﾝｷﾝｶﾞｸ TO JR-ｷﾝｶﾞｸ
           MOVE "満期保険金" TO JR-ﾃｷﾖｳ
           PERFORM ｼﾜｹ-RTN
           IF PM-ﾊｲﾄｳﾙｲｹｲ > 0
               MOVE 2310 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
               MOVE PM-ﾊｲﾄｳﾙｲｹｲ TO JR-ｷﾝｶﾞｸ
               MOVE "配当金支払" TO JR-ﾃｷﾖｳ
               PERFORM ｼﾜｹ-RTN
           END-IF
           ADD PY-ｷﾝｶﾞｸ TO WS-ﾏﾝｷｷﾝｶﾞｸ
           MOVE PY-ｷﾝｶﾞｸ TO WS-ﾒｲｻｲ-ｷﾝｶﾞｸ
           MOVE 0 TO PM-ﾊｲﾄｳﾙｲｹｲ
           SET PM-ﾏﾝｷ TO TRUE
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO PM-ｼﾞｮｳﾀｲﾋﾂﾞｹ
           MOVE "満期" TO WS-ﾒｲｻｲ-ｹｯｶ.

      *----------------------------------------------------------------
      * ｺｳｼﾝ-RTN: 定期保険の更新。その時点の保険年齢が商品の最高
      * 年齢以内なら保険料を引き直して期間を延ばし、超えていれば
      * 期間満了。
      *----------------------------------------------------------------
       ｺｳｼﾝ-RTN.
           MOVE "S" TO WS-ﾎｳｺｳ
           CALL "GLC020" USING WS-ﾎｳｺｳ WS-ｾｲﾈﾝｶﾞｯﾋﾟ
                               PM-ｾｲﾈﾝｶﾞｯﾋﾟ WS-ﾍﾝｶﾝｹｯｶ
           CALL "GLC010" USING WS-ｾｲﾈﾝｶﾞｯﾋﾟ WS-ｷﾞｮｳﾑﾋﾂﾞｹ WS-ﾈﾝﾚｲ
           SET WS-ｼｮｳﾋﾝ-IX TO 1
           SEARCH WS-ｼｮｳﾋﾝ
               WHEN WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = PM-ｼｮｳﾋﾝ
                   CONTINUE
           END-SEARCH
           IF WS-ﾈﾝﾚｲ > WS-ｼｮｳﾋﾝ-ｻｲｺｳﾈﾝﾚｲ(WS-ｼｮｳﾋﾝ-IX)
               ADD 1 TO WS-ﾏﾝﾘｮｳｹﾝｽｳ
               SET PM-ﾏﾝｷ TO TRUE
               MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO PM-ｼﾞｮｳﾀｲﾋﾂﾞｹ
               MOVE "期間満了" TO WS-ﾒｲｻｲ-ｹｯｶ
           ELSE
               ADD 1 TO WS-ｺｳｼﾝｹﾝｽｳ
               COMPUTE WS-ﾈﾝﾚｲ-ｿｴｼﾞ = WS-ﾈﾝﾚｲ + 1
               COMPUTE WS-ｹﾞﾂｶﾞｸ = PM-ﾎｹﾝｷﾝｶﾞｸ / 1000
                   * WS-ﾘｮｳﾘﾂ(WS-ｼｮｳﾋﾝ-IX, PM-ｾｲﾍﾞﾂ, WS-ﾈﾝﾚｲ-ｿｴｼﾞ)
               IF PM-ﾈﾝﾊﾞﾗｲ
                   COMPUTE PM-ﾎｹﾝﾘｮｳ = WS-ｹﾞﾂｶﾞｸ * 12 * 0.95
               ELSE
                   MOVE WS-ｹﾞﾂｶﾞｸ TO PM-ﾎｹﾝﾘｮｳ
               END-IF
               MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ TO WS-ﾏﾝｷﾋﾞ
               ADD PM-ﾎｹﾝｷｶﾝ TO WS-ﾏﾝｷ-YYYY
               MOVE "W" TO WS-ﾎｳｺｳ
               CALL "GLC020" USING WS-ﾎｳｺｳ WS-ﾏﾝｷﾋﾞ
                                   WS-ﾏﾝｷﾋﾞ-ﾜﾚｷ WS-ﾍﾝｶﾝｹｯｶ
               MOVE WS-ﾏﾝｷﾋﾞ-ﾜﾚｷ TO PM-ﾏﾝｷﾋﾞ
               MOVE PM-ﾎｹﾝﾘｮｳ TO WS-ﾒｲｻｲ-ｷﾝｶﾞｸ
               MOVE "更新" TO WS-ﾒｲｻｲ-ｹｯｶ
           END-IF.

      *----------------------------------------------------------------
      * ｼﾜｹ-RTN: 借方1行。科目、金額、摘要は呼ぶ側が入れる。
      *----------------------------------------------------------------
       ｼﾜｹ-RTN.
           MOVE CT-ｼﾞｶｲｼﾜｹﾊﾞﾝｺﾞｳ TO JR-ｼﾜｹﾊﾞﾝｺﾞｳ
           ADD 1 TO CT-ｼﾞｶｲｼﾜｹﾊﾞﾝｺﾞｳ
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO JR-ﾋﾂﾞｹ
           MOVE "GLD020" TO JR-ﾌﾟﾛｸﾞﾗﾑ
           SET JR-ｶﾘｶﾀ TO TRUE
           MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO JR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           WRITE JR-ｼﾜｹ-ﾚｺｰﾄﾞ
           IF NOT WS-GLJRNL-OK
               MOVE "GLJRNL" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-GLJRNL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ﾒｲｻｲ-RTN: 1件の明細行。
      *----------------------------------------------------------------
       ﾒｲｻｲ-RTN.
           MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO WS-ﾒｲｻｲ-ｼｮｳｹﾝ
           MOVE PM-ｼｮｳﾋﾝ TO WS-ﾒｲｻｲ-ｼｮｳﾋﾝ
           MOVE PM-ｼﾒｲ TO WS-ﾒｲｻｲ-ｼﾒｲ
           MOVE PM-ｹｲｶﾈﾝｽｳ TO WS-ﾒｲｻｲ-ﾈﾝｽｳ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ﾒｲｻｲ-ｷﾞｮｳ.

      *----------------------------------------------------------------
      * ﾎｳｺｸ-RTN: 件数。
      *----------------------------------------------------------------
       ﾎｳｺｸ-RTN.
           MOVE SPACES TO ﾎｳｺｸ-ｷﾞｮｳ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ
           MOVE "応当件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｵｳﾄｳｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "更新件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｺｳｼﾝｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "期間満了件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾏﾝﾘｮｳｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "満期件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾏﾝｷｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "満期保険金" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾏﾝｷｷﾝｶﾞｸ TO WS-ｹﾝｽｳ-ｱﾀｲ
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
           CLOSE POLMST-F PRODTB-F RATETB-F GLSCTL-F PAYOUT-F
                 GLJRNL-F REPORT-F
           DISPLAY "GLD020 応当 " WS-ｵｳﾄｳｹﾝｽｳ " 更新 " WS-ｺｳｼﾝｹﾝｽｳ
                   " 満期 " WS-ﾏﾝｷｹﾝｽｳ.

      *----------------------------------------------------------------
      * ｲｼﾞｮｳ-RTN: ファイル異常は 16 で終わる。
      *----------------------------------------------------------------
       ｲｼﾞｮｳ-RTN.
           DISPLAY "GLD020 ファイル異常 " WS-ｴﾗｰﾌｧｲﾙ
                   " 状態 " WS-ｴﾗｰｽﾃｰﾀｽ
           MOVE 16 TO RETURN-CODE
           STOP RUN.

       COPY PRDLOD.

       COPY RATLOD.
