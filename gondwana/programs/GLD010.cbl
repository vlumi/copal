      *----------------------------------------------------------------
      * GLD010: 新契約登録
      *
      * 日次。営業システムからの新契約ファイル NEWPOL を読み、
      * 1件ずつ検査して契約マスタに登録する。検査は商品、証券番号
      * の重複、払込方法、契約日と生年月日 (西暦で届くので和暦に
      * 変換できること)、性別、契約日時点の保険年齢が商品の範囲内
      * であること、保険金額が範囲内であること、料率があること。
      * 保険料は保険金額1,000円あたりの月払料率から求め、円未満
      * 切捨て。年払は月払の12倍から5%を引く。第1回保険料は申込時
      * に収納済みのため、払込済年月は契約月とし、保険料収入を
      * 仕訳する。満期日は契約日に保険期間を足した日 (終身・医療は
      * なし)。却下は理由コードを付けて却下ファイルへ。
      *
      * 復帰コード: 0 全件登録、4 却下あり、8 日付が日付でない、
      *             16 ファイル異常。
      *
      * 変更履歴
      * 日付       担当 変更
      * 1987-04-01 KT   GL-0001 初版、申込書のパンチ入力から。
      * 2003-10-06 MS   GL-0310 営業システムのファイルから。
      * 2008-06-02 YN   GL-0455 医療保険。
      *----------------------------------------------------------------
       IDENTIFICATION DIVISION.
       PROGRAM-ID. GLD010.
       AUTHOR. ゴンドワナ生命 システム部.
       DATE-WRITTEN. 1987-04-01.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT NEWPOL-F ASSIGN TO "NEWPOL"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-NEWPOL-ST.
           SELECT POLMST-F ASSIGN TO "POLMST"
               ORGANIZATION IS INDEXED
               ACCESS MODE IS RANDOM
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
           SELECT OPTIONAL GLJRNL-F ASSIGN TO "GLJRNL"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-GLJRNL-ST.
           SELECT REJECT-F ASSIGN TO "NEWREJ"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REJECT-ST.
           SELECT REPORT-F ASSIGN TO "GLD010R"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-ST.

       DATA DIVISION.
       FILE SECTION.
       FD  NEWPOL-F.
       COPY NEWPOL.
       FD  POLMST-F.
       COPY POLMST.
       FD  PRODTB-F.
       COPY PRODTB.
       FD  RATETB-F.
       COPY RATETB.
       FD  GLSCTL-F
           RECORD CONTAINS 64 CHARACTERS.
       COPY CTLREC.
       FD  GLJRNL-F.
       COPY GLJRNL.
       FD  REJECT-F.
       01  ｷｬｯｶ-ﾚｺｰﾄﾞ.
           05  ｷｬｯｶ-ﾃﾞｰﾀ               PIC X(200).
           05  ｷｬｯｶ-ﾘﾕｳ                PIC 99.
           05  ｷｬｯｶ-ﾒｯｾｰｼﾞ             PIC X(20).
       FD  REPORT-F.
       01  ﾎｳｺｸ-ｷﾞｮｳ                   PIC X(100).

       WORKING-STORAGE SECTION.
       01  WS-ｽﾃｰﾀｽ.
           05  WS-NEWPOL-ST            PIC XX.
               88  WS-NEWPOL-OK                VALUE "00".
               88  WS-NEWPOL-ｵﾜﾘ               VALUE "10".
           05  WS-POLMST-ST            PIC XX.
               88  WS-POLMST-OK                VALUE "00".
               88  WS-POLMST-ﾅｼ                VALUE "23".
           05  WS-PRODTB-ST            PIC XX.
               88  WS-PRODTB-OK                VALUE "00".
               88  WS-PRODTB-ｵﾜﾘ               VALUE "10".
           05  WS-RATETB-ST            PIC XX.
               88  WS-RATETB-OK                VALUE "00".
               88  WS-RATETB-ｵﾜﾘ               VALUE "10".
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
       01  WS-ｹｲﾔｸﾋﾞ-ﾜﾚｷ               PIC X(7).
       01  WS-ｾｲﾈﾝｶﾞｯﾋﾟ-ﾜﾚｷ            PIC X(7).
       01  WS-ﾏﾝｷﾋﾞ                    PIC 9(8).
       01  WS-ﾏﾝｷﾋﾞ-ﾌﾞﾝｶｲ REDEFINES WS-ﾏﾝｷﾋﾞ.
           05  WS-ﾏﾝｷ-YYYY             PIC 9(4).
           05  WS-ﾏﾝｷ-MMDD             PIC 9(4).
       01  WS-ﾏﾝｷﾋﾞ-ﾜﾚｷ                PIC X(7).
       01  WS-ﾈﾝﾚｲ                     PIC 99.
       01  WS-ｹﾞﾂｶﾞｸ                   PIC 9(9).
       COPY TBLWS.
       01  WS-ｹﾝｽｳ.
           05  WS-ﾖﾐｺﾐｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ﾄｳﾛｸｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ｷｬｯｶｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ﾄｳﾛｸｷﾝｶﾞｸ            PIC 9(13) VALUE 0.
           05  WS-ﾄｳﾛｸﾎｹﾝﾘｮｳ           PIC 9(11) VALUE 0.
       01  WS-ﾘﾕｳ                      PIC 99.
           88  WS-ﾘﾕｳ-ﾅｼ                       VALUE 0.
           88  WS-ｼｮｳﾋﾝ-ﾅｼ                     VALUE 1.
           88  WS-ｼｮｳｹﾝ-ﾁｮｳﾌｸ                  VALUE 2.
           88  WS-ﾊﾗｲｺﾐ-ﾌｾｲ                    VALUE 3.
           88  WS-ｹｲﾔｸﾋﾞ-ﾌｾｲ                   VALUE 4.
           88  WS-ｾｲﾈﾝｶﾞｯﾋﾟ-ﾌｾｲ                VALUE 5.
           88  WS-ｾｲﾍﾞﾂ-ﾌｾｲ                    VALUE 6.
           88  WS-ﾈﾝﾚｲ-ﾊﾝｲｶﾞｲ                  VALUE 7.
           88  WS-ｷﾝｶﾞｸ-ﾊﾝｲｶﾞｲ                 VALUE 8.
           88  WS-ﾘｮｳﾘﾂ-ﾅｼ                     VALUE 9.
       01  WS-ﾘﾕｳ-ﾒｯｾｰｼﾞ.
           05  FILLER PIC X(20) VALUE "商品なし".
           05  FILLER PIC X(20) VALUE "証券番号重複".
           05  FILLER PIC X(20) VALUE "払込方法不正".
           05  FILLER PIC X(20) VALUE "契約日不正".
           05  FILLER PIC X(20) VALUE "生年月日不正".
           05  FILLER PIC X(20) VALUE "性別不正".
           05  FILLER PIC X(20) VALUE "年齢範囲外".
           05  FILLER PIC X(20) VALUE "保険金額範囲外".
           05  FILLER PIC X(20) VALUE "料率なし".
       01  WS-ﾘﾕｳ-ﾒｯｾｰｼﾞ-TBL REDEFINES WS-ﾘﾕｳ-ﾒｯｾｰｼﾞ.
           05  WS-ﾒｯｾｰｼﾞ OCCURS 9 TIMES PIC X(20).
       01  WS-ﾌｯｷｺｰﾄﾞ                  PIC 99 VALUE 0.
       01  WS-ﾎｳｺｸ.
           05  WS-ﾋｮｳﾀﾞｲ-ｷﾞｮｳ.
               10  FILLER                  PIC X(28)
                   VALUE "ゴンドワナ生命  新契約登録".
               10  FILLER                  PIC X(4) VALUE SPACES.
               10  WS-ﾋｮｳﾀﾞｲ-ﾋﾂﾞｹ          PIC X(7).
           05  WS-ﾒｲｻｲ-ｷﾞｮｳ.
               10  WS-ﾒｲｻｲ-ｼｮｳｹﾝ           PIC X(10).
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾒｲｻｲ-ｼｮｳﾋﾝ           PIC X(2).
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾒｲｻｲ-ｼﾒｲ             PIC X(20).
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾒｲｻｲ-ﾈﾝﾚｲ            PIC Z9.
               10  FILLER                  PIC X(2) VALUE "歳".
               10  WS-ﾒｲｻｲ-ｷﾝｶﾞｸ           PIC ZZZ,ZZZ,ZZ9.
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾒｲｻｲ-ﾎｹﾝﾘｮｳ          PIC Z,ZZZ,ZZ9.
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾒｲｻｲ-ｹｯｶ             PIC X(20).
           05  WS-ｹﾝｽｳ-ｷﾞｮｳ.
               10  WS-ｹﾝｽｳ-ﾗﾍﾞﾙ            PIC X(16).
               10  WS-ｹﾝｽｳ-ｱﾀｲ             PIC ZZZ,ZZZ,ZZ9.

       PROCEDURE DIVISION.
      *----------------------------------------------------------------
      * ﾒｲﾝ-RTN: 準備、1件ずつの登録、報告、終了。
      *----------------------------------------------------------------
       ﾒｲﾝ-RTN.
           PERFORM ｼｮｷ-RTN
           PERFORM ｼｮﾘ-RTN UNTIL WS-NEWPOL-ｵﾜﾘ
           PERFORM ﾎｳｺｸ-RTN
           PERFORM ｼｭｳﾘｮｳ-RTN
           MOVE WS-ﾌｯｷｺｰﾄﾞ TO RETURN-CODE
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
               DISPLAY "GLD010 日付不正: " WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8)
               MOVE 8 TO RETURN-CODE
               STOP RUN
           END-IF
           OPEN INPUT NEWPOL-F PRODTB-F RATETB-F
           OPEN I-O POLMST-F GLSCTL-F
           OPEN EXTEND GLJRNL-F
           OPEN OUTPUT REJECT-F REPORT-F
           PERFORM ｵｰﾌﾟﾝ-ｹﾝｻ-RTN
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
           PERFORM NEWPOL-ﾖﾐｺﾐ-RTN.

      *----------------------------------------------------------------
      * ｵｰﾌﾟﾝ-ｹﾝｻ-RTN: 開けなかったファイルがあれば異常終了。
      *----------------------------------------------------------------
       ｵｰﾌﾟﾝ-ｹﾝｻ-RTN.
           EVALUATE TRUE
               WHEN NOT WS-NEWPOL-OK
                   MOVE "NEWPOL" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-NEWPOL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
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
               WHEN NOT WS-GLJRNL-OK
                   MOVE "GLJRNL" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-GLJRNL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-REJECT-OK
                   MOVE "NEWREJ" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-REJECT-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-REPORT-OK
                   MOVE "GLD010R" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-REPORT-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN OTHER
                   MOVE SPACES TO WS-ｴﾗｰﾌｧｲﾙ
           END-EVALUATE
           IF WS-ｴﾗｰﾌｧｲﾙ NOT = SPACES
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * NEWPOL-ﾖﾐｺﾐ-RTN: 新契約ファイルの次の行。
      *----------------------------------------------------------------
       NEWPOL-ﾖﾐｺﾐ-RTN.
           READ NEWPOL-F
           IF WS-NEWPOL-OK
               ADD 1 TO WS-ﾖﾐｺﾐｹﾝｽｳ
           ELSE
               IF NOT WS-NEWPOL-ｵﾜﾘ
                   MOVE "NEWPOL" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-NEWPOL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
                   PERFORM ｲｼﾞｮｳ-RTN
               END-IF
           END-IF.

      *----------------------------------------------------------------
      * ｼｮﾘ-RTN: 1件を検査し、通れば登録、通らなければ却下。
      *----------------------------------------------------------------
       ｼｮﾘ-RTN.
           PERFORM ｹﾝｻ-RTN
           IF WS-ﾘﾕｳ-ﾅｼ
               PERFORM ﾄｳﾛｸ-RTN
           ELSE
               PERFORM ｷｬｯｶ-RTN
           END-IF
           PERFORM NEWPOL-ﾖﾐｺﾐ-RTN.

      *----------------------------------------------------------------
      * ｹﾝｻ-RTN: 検査を順に行い、最初に引っかかった理由を残す。
      *----------------------------------------------------------------
       ｹﾝｻ-RTN.
           SET WS-ﾘﾕｳ-ﾅｼ TO TRUE
           SET WS-ｼｮｳﾋﾝ-IX TO 1
           SEARCH WS-ｼｮｳﾋﾝ
               AT END
                   SET WS-ｼｮｳﾋﾝ-ﾅｼ TO TRUE
               WHEN WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = NP-ｼｮｳﾋﾝ
                   CONTINUE
           END-SEARCH
           IF WS-ﾘﾕｳ-ﾅｼ
               MOVE NP-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
               READ POLMST-F
               IF WS-POLMST-OK
                   SET WS-ｼｮｳｹﾝ-ﾁｮｳﾌｸ TO TRUE
               ELSE
                   IF NOT WS-POLMST-ﾅｼ
                       MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
                       MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
                       PERFORM ｲｼﾞｮｳ-RTN
                   END-IF
               END-IF
           END-IF
           IF WS-ﾘﾕｳ-ﾅｼ AND NP-ﾊﾗｲｺﾐﾎｳﾎｳ NOT = "1" AND "2"
               SET WS-ﾊﾗｲｺﾐ-ﾌｾｲ TO TRUE
           END-IF
           IF WS-ﾘﾕｳ-ﾅｼ
               MOVE "W" TO WS-ﾎｳｺｳ
               CALL "GLC020" USING WS-ﾎｳｺｳ NP-ｹｲﾔｸﾋﾞ
                                   WS-ｹｲﾔｸﾋﾞ-ﾜﾚｷ WS-ﾍﾝｶﾝｹｯｶ
               IF WS-ﾍﾝｶﾝｹｯｶ NOT = 0
                   SET WS-ｹｲﾔｸﾋﾞ-ﾌｾｲ TO TRUE
               END-IF
           END-IF
           IF WS-ﾘﾕｳ-ﾅｼ
               CALL "GLC020" USING WS-ﾎｳｺｳ NP-ｾｲﾈﾝｶﾞｯﾋﾟ
                                   WS-ｾｲﾈﾝｶﾞｯﾋﾟ-ﾜﾚｷ WS-ﾍﾝｶﾝｹｯｶ
               IF WS-ﾍﾝｶﾝｹｯｶ NOT = 0
                   SET WS-ｾｲﾈﾝｶﾞｯﾋﾟ-ﾌｾｲ TO TRUE
               END-IF
           END-IF
           IF WS-ﾘﾕｳ-ﾅｼ AND NP-ｾｲﾍﾞﾂ NOT = 1 AND 2
               SET WS-ｾｲﾍﾞﾂ-ﾌｾｲ TO TRUE
           END-IF
           IF WS-ﾘﾕｳ-ﾅｼ
               PERFORM ﾈﾝﾚｲ-ｷﾝｶﾞｸ-ｹﾝｻ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ﾈﾝﾚｲ-ｷﾝｶﾞｸ-ｹﾝｻ-RTN: 契約日時点の保険年齢と保険金額が商品の
      * 範囲内で、料率があること。
      *----------------------------------------------------------------
       ﾈﾝﾚｲ-ｷﾝｶﾞｸ-ｹﾝｻ-RTN.
           CALL "GLC010" USING NP-ｾｲﾈﾝｶﾞｯﾋﾟ NP-ｹｲﾔｸﾋﾞ WS-ﾈﾝﾚｲ
           EVALUATE TRUE
               WHEN WS-ﾈﾝﾚｲ < WS-ｼｮｳﾋﾝ-ｻｲﾃｲﾈﾝﾚｲ(WS-ｼｮｳﾋﾝ-IX)
                 OR WS-ﾈﾝﾚｲ > WS-ｼｮｳﾋﾝ-ｻｲｺｳﾈﾝﾚｲ(WS-ｼｮｳﾋﾝ-IX)
                   SET WS-ﾈﾝﾚｲ-ﾊﾝｲｶﾞｲ TO TRUE
               WHEN NP-ﾎｹﾝｷﾝｶﾞｸ < WS-ｼｮｳﾋﾝ-ｻｲﾃｲｷﾝｶﾞｸ(WS-ｼｮｳﾋﾝ-IX)
                 OR NP-ﾎｹﾝｷﾝｶﾞｸ > WS-ｼｮｳﾋﾝ-ｻｲｺｳｷﾝｶﾞｸ(WS-ｼｮｳﾋﾝ-IX)
                   SET WS-ｷﾝｶﾞｸ-ﾊﾝｲｶﾞｲ TO TRUE
               WHEN OTHER
                   COMPUTE WS-ﾈﾝﾚｲ-ｿｴｼﾞ = WS-ﾈﾝﾚｲ + 1
                   IF WS-ﾘｮｳﾘﾂ(WS-ｼｮｳﾋﾝ-IX, NP-ｾｲﾍﾞﾂ, WS-ﾈﾝﾚｲ-ｿｴｼﾞ)
                           = 0
                       SET WS-ﾘｮｳﾘﾂ-ﾅｼ TO TRUE
                   END-IF
           END-EVALUATE.

      *----------------------------------------------------------------
      * ﾄｳﾛｸ-RTN: マスタレコードを組み立てて書き、第1回保険料を
      * 仕訳し、明細を報告する。
      *----------------------------------------------------------------
       ﾄｳﾛｸ-RTN.
           MOVE SPACES TO PM-ｹｲﾔｸ-ﾚｺｰﾄﾞ
           MOVE NP-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           SET PM-ﾕｳｺｳ TO TRUE
           MOVE NP-ｼｮｳﾋﾝ TO PM-ｼｮｳﾋﾝ
           MOVE NP-ﾎｹﾝｷﾝｶﾞｸ TO PM-ﾎｹﾝｷﾝｶﾞｸ
           PERFORM ﾎｹﾝﾘｮｳ-RTN
           MOVE NP-ﾊﾗｲｺﾐﾎｳﾎｳ TO PM-ﾊﾗｲｺﾐﾎｳﾎｳ
           MOVE WS-ｹｲﾔｸﾋﾞ-ﾜﾚｷ TO PM-ｹｲﾔｸﾋﾞ
           MOVE WS-ｼｮｳﾋﾝ-ﾎｹﾝｷｶﾝ(WS-ｼｮｳﾋﾝ-IX) TO PM-ﾎｹﾝｷｶﾝ
           PERFORM ﾏﾝｷﾋﾞ-RTN
           MOVE 1 TO PM-ｹｲｶﾈﾝｽｳ
           MOVE NP-ｹｲﾔｸﾋﾞ(1:6) TO PM-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ
           MOVE 0 TO PM-ﾐﾉｳｶｲｽｳ
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO PM-ｼﾞｮｳﾀｲﾋﾂﾞｹ
           MOVE NP-ｼﾒｲ TO PM-ｼﾒｲ
           MOVE NP-ｼﾒｲｶﾅ TO PM-ｼﾒｲｶﾅ
           MOVE WS-ｾｲﾈﾝｶﾞｯﾋﾟ-ﾜﾚｷ TO PM-ｾｲﾈﾝｶﾞｯﾋﾟ
           MOVE NP-ｾｲﾍﾞﾂ TO PM-ｾｲﾍﾞﾂ
           MOVE WS-ﾈﾝﾚｲ TO PM-ｶﾆｭｳﾈﾝﾚｲ
           MOVE NP-ｼﾞｭｳｼｮ TO PM-ｼﾞｭｳｼｮ
           MOVE NP-ｳｹﾄﾘﾆﾝ TO PM-ｳｹﾄﾘﾆﾝ
           MOVE NP-ｷﾞﾝｺｳｺｰﾄﾞ TO PM-ｷﾞﾝｺｳｺｰﾄﾞ
           MOVE NP-ｼﾃﾝｺｰﾄﾞ TO PM-ｼﾃﾝｺｰﾄﾞ
           MOVE NP-ﾖｷﾝｼｭﾓｸ TO PM-ﾖｷﾝｼｭﾓｸ
           MOVE NP-ｺｳｻﾞﾊﾞﾝｺﾞｳ TO PM-ｺｳｻﾞﾊﾞﾝｺﾞｳ
           MOVE NP-ｺｳｻﾞﾒｲｷﾞ TO PM-ｺｳｻﾞﾒｲｷﾞ
           MOVE NP-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ TO PM-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ
           MOVE 0 TO PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ PM-ｾｷﾆﾝｼﾞｭﾝﾋﾞｷﾝ PM-ﾊｲﾄｳﾙｲｹｲ
           MOVE NP-ｹｲﾔｸﾋﾞ(1:6) TO PM-ｻｲｼｭｳﾆｭｳｷﾝﾈﾝｹﾞﾂ
           MOVE PM-ﾎｹﾝﾘｮｳ TO PM-ｻｲｼｭｳﾆｭｳｷﾝｶﾞｸ
           WRITE PM-ｹｲﾔｸ-ﾚｺｰﾄﾞ
           IF NOT WS-POLMST-OK
               MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           PERFORM ｼﾜｹ-RTN
           ADD 1 TO WS-ﾄｳﾛｸｹﾝｽｳ
           ADD PM-ﾎｹﾝｷﾝｶﾞｸ TO WS-ﾄｳﾛｸｷﾝｶﾞｸ
           ADD PM-ﾎｹﾝﾘｮｳ TO WS-ﾄｳﾛｸﾎｹﾝﾘｮｳ
           MOVE "登録" TO WS-ﾒｲｻｲ-ｹｯｶ
           PERFORM ﾒｲｻｲ-RTN.

      *----------------------------------------------------------------
      * ﾎｹﾝﾘｮｳ-RTN: 保険金額1,000円あたりの月払料率から、円未満
      * 切捨て。年払は12か月分の5%引き。
      *----------------------------------------------------------------
       ﾎｹﾝﾘｮｳ-RTN.
           COMPUTE WS-ｹﾞﾂｶﾞｸ = NP-ﾎｹﾝｷﾝｶﾞｸ / 1000
               * WS-ﾘｮｳﾘﾂ(WS-ｼｮｳﾋﾝ-IX, NP-ｾｲﾍﾞﾂ, WS-ﾈﾝﾚｲ-ｿｴｼﾞ)
           IF NP-ﾊﾗｲｺﾐﾎｳﾎｳ = "2"
               COMPUTE PM-ﾎｹﾝﾘｮｳ = WS-ｹﾞﾂｶﾞｸ * 12 * 0.95
           ELSE
               MOVE WS-ｹﾞﾂｶﾞｸ TO PM-ﾎｹﾝﾘｮｳ
           END-IF.

      *----------------------------------------------------------------
      * ﾏﾝｷﾋﾞ-RTN: 契約日に保険期間を足した日を和暦で。期間のない
      * 商品 (終身、医療) は空白。
      *----------------------------------------------------------------
       ﾏﾝｷﾋﾞ-RTN.
           IF PM-ﾎｹﾝｷｶﾝ = 0
               MOVE SPACES TO PM-ﾏﾝｷﾋﾞ
           ELSE
               MOVE NP-ｹｲﾔｸﾋﾞ TO WS-ﾏﾝｷﾋﾞ
               ADD PM-ﾎｹﾝｷｶﾝ TO WS-ﾏﾝｷ-YYYY
               CALL "GLC020" USING WS-ﾎｳｺｳ WS-ﾏﾝｷﾋﾞ
                                   WS-ﾏﾝｷﾋﾞ-ﾜﾚｷ WS-ﾍﾝｶﾝｹｯｶ
               MOVE WS-ﾏﾝｷﾋﾞ-ﾜﾚｷ TO PM-ﾏﾝｷﾋﾞ
           END-IF.

      *----------------------------------------------------------------
      * ｼﾜｹ-RTN: 第1回保険料の保険料収入を貸方に仕訳する。
      *----------------------------------------------------------------
       ｼﾜｹ-RTN.
           MOVE SPACES TO JR-ｼﾜｹ-ﾚｺｰﾄﾞ
           MOVE CT-ｼﾞｶｲｼﾜｹﾊﾞﾝｺﾞｳ TO JR-ｼﾜｹﾊﾞﾝｺﾞｳ
           ADD 1 TO CT-ｼﾞｶｲｼﾜｹﾊﾞﾝｺﾞｳ
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO JR-ﾋﾂﾞｹ
           MOVE "GLD010" TO JR-ﾌﾟﾛｸﾞﾗﾑ
           MOVE 4110 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
           SET JR-ｶｼｶﾀ TO TRUE
           MOVE PM-ﾎｹﾝﾘｮｳ TO JR-ｷﾝｶﾞｸ
           MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO JR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           MOVE "第1回保険料" TO JR-ﾃｷﾖｳ
           WRITE JR-ｼﾜｹ-ﾚｺｰﾄﾞ.

      *----------------------------------------------------------------
      * ｷｬｯｶ-RTN: 却下ファイルに理由を付けて書き、明細を報告する。
      *----------------------------------------------------------------
       ｷｬｯｶ-RTN.
           MOVE NP-ｼﾝｹｲﾔｸ-ﾚｺｰﾄﾞ TO ｷｬｯｶ-ﾃﾞｰﾀ
           MOVE WS-ﾘﾕｳ TO ｷｬｯｶ-ﾘﾕｳ
           MOVE WS-ﾒｯｾｰｼﾞ(WS-ﾘﾕｳ) TO ｷｬｯｶ-ﾒｯｾｰｼﾞ
           WRITE ｷｬｯｶ-ﾚｺｰﾄﾞ
           ADD 1 TO WS-ｷｬｯｶｹﾝｽｳ
           MOVE 4 TO WS-ﾌｯｷｺｰﾄﾞ
           MOVE 0 TO WS-ﾈﾝﾚｲ PM-ﾎｹﾝﾘｮｳ
           MOVE WS-ﾒｯｾｰｼﾞ(WS-ﾘﾕｳ) TO WS-ﾒｲｻｲ-ｹｯｶ
           PERFORM ﾒｲｻｲ-RTN.

      *----------------------------------------------------------------
      * ﾒｲｻｲ-RTN: 1件の明細行。
      *----------------------------------------------------------------
       ﾒｲｻｲ-RTN.
           MOVE NP-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO WS-ﾒｲｻｲ-ｼｮｳｹﾝ
           MOVE NP-ｼｮｳﾋﾝ TO WS-ﾒｲｻｲ-ｼｮｳﾋﾝ
           MOVE NP-ｼﾒｲ TO WS-ﾒｲｻｲ-ｼﾒｲ
           MOVE WS-ﾈﾝﾚｲ TO WS-ﾒｲｻｲ-ﾈﾝﾚｲ
           MOVE NP-ﾎｹﾝｷﾝｶﾞｸ TO WS-ﾒｲｻｲ-ｷﾝｶﾞｸ
           MOVE PM-ﾎｹﾝﾘｮｳ TO WS-ﾒｲｻｲ-ﾎｹﾝﾘｮｳ
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
           MOVE "登録件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾄｳﾛｸｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "却下件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｷｬｯｶｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "登録保険金額" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾄｳﾛｸｷﾝｶﾞｸ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "登録保険料" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾄｳﾛｸﾎｹﾝﾘｮｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
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
           CLOSE NEWPOL-F POLMST-F PRODTB-F RATETB-F GLSCTL-F
                 GLJRNL-F REJECT-F REPORT-F
           DISPLAY "GLD010 登録 " WS-ﾄｳﾛｸｹﾝｽｳ
                   " 却下 " WS-ｷｬｯｶｹﾝｽｳ.

      *----------------------------------------------------------------
      * ｲｼﾞｮｳ-RTN: ファイル異常は 16 で終わる。
      *----------------------------------------------------------------
       ｲｼﾞｮｳ-RTN.
           DISPLAY "GLD010 ファイル異常 " WS-ｴﾗｰﾌｧｲﾙ
                   " 状態 " WS-ｴﾗｰｽﾃｰﾀｽ
           MOVE 16 TO RETURN-CODE
           STOP RUN.

       COPY PRDLOD.

       COPY RATLOD.
