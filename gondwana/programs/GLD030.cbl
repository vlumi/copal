      *----------------------------------------------------------------
      * GLD030: 給付支払
      *
      * 日次。支社からの給付請求ファイル CLMREQ を読み、死亡保険金、
      * 入院給付金、解約返戻金を支払う。
      * 検査は種別、契約の有無、事由発生日 (和暦で届く)、事由が契約
      * 日以後であること、契約の状態: 死亡・入院は保障中 (有効・
      * 猶予中)、解約は失効中でもよい。
      * 死亡: 保険金額に配当累計を加え、貸付残高を引いた額。支払済
      * とする。
      * 入院: 医療保険に限る。日数は 1 以上、120 日まで。契約日から
      * 免責日数の間の入院は支払わない。日額は保険金額に入院日額率
      * を掛けて 1,000 で割った額。契約はそのまま。
      * 解約: 解約返戻金係数 (商品・加入年齢区分・経過年数) を保険
      * 金額 1,000 円あたりに掛け、貸付残高を引いた額。0 以下なら
      * 支払なしで解約とする。
      * 貸付残高は支払時に回収し、仕訳する。
      *
      * 復帰コード: 0 全件支払、4 却下あり、8 日付が日付でない、
      *             16 ファイル異常。
      *
      * 変更履歴
      * 日付       担当 変更
      * 1987-04-01 KT   GL-0001 初版。
      * 2001-07-02 MS   GL-0271 死亡保険金に配当累計を加える。
      * 2008-06-02 YN   GL-0455 入院給付金、免責期間。
      *----------------------------------------------------------------
       IDENTIFICATION DIVISION.
       PROGRAM-ID. GLD030.
       AUTHOR. ゴンドワナ生命 システム部.
       DATE-WRITTEN. 1987-04-01.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT CLMREQ-F ASSIGN TO "CLMREQ"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-CLMREQ-ST.
           SELECT POLMST-F ASSIGN TO "POLMST"
               ORGANIZATION IS INDEXED
               ACCESS MODE IS RANDOM
               RECORD KEY IS PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
               FILE STATUS IS WS-POLMST-ST.
           SELECT PRODTB-F ASSIGN TO "PRODTB"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-PRODTB-ST.
           SELECT CSVTB-F ASSIGN TO "CSVTB"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-CSVTB-ST.
           SELECT GLSCTL-F ASSIGN TO "GLSCTL"
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-GLSCTL-ST.
           SELECT OPTIONAL PAYOUT-F ASSIGN TO "PAYOUT"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-PAYOUT-ST.
           SELECT OPTIONAL GLJRNL-F ASSIGN TO "GLJRNL"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-GLJRNL-ST.
           SELECT REJECT-F ASSIGN TO "CLMREJ"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REJECT-ST.
           SELECT REPORT-F ASSIGN TO "GLD030R"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-ST.

       DATA DIVISION.
       FILE SECTION.
       FD  CLMREQ-F.
       COPY CLMREQ.
       FD  POLMST-F.
       COPY POLMST.
       FD  PRODTB-F.
       COPY PRODTB.
       FD  CSVTB-F.
       COPY CSVTB.
       FD  GLSCTL-F
           RECORD CONTAINS 64 CHARACTERS.
       COPY CTLREC.
       FD  PAYOUT-F.
       COPY PAYOUT.
       FD  GLJRNL-F.
       COPY GLJRNL.
       FD  REJECT-F.
       01  ｷｬｯｶ-ﾚｺｰﾄﾞ.
           05  ｷｬｯｶ-ﾃﾞｰﾀ               PIC X(80).
           05  ｷｬｯｶ-ﾘﾕｳ                PIC 99.
           05  ｷｬｯｶ-ﾒｯｾｰｼﾞ             PIC X(20).
       FD  REPORT-F.
       01  ﾎｳｺｸ-ｷﾞｮｳ                   PIC X(100).

       WORKING-STORAGE SECTION.
       01  WS-ｽﾃｰﾀｽ.
           05  WS-CLMREQ-ST            PIC XX.
               88  WS-CLMREQ-OK                VALUE "00".
               88  WS-CLMREQ-ｵﾜﾘ               VALUE "10".
           05  WS-POLMST-ST            PIC XX.
               88  WS-POLMST-OK                VALUE "00".
               88  WS-POLMST-ﾅｼ                VALUE "23".
           05  WS-PRODTB-ST            PIC XX.
               88  WS-PRODTB-OK                VALUE "00".
               88  WS-PRODTB-ｵﾜﾘ               VALUE "10".
           05  WS-CSVTB-ST             PIC XX.
               88  WS-CSVTB-OK                 VALUE "00".
               88  WS-CSVTB-ｵﾜﾘ                VALUE "10".
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
       01  WS-ｺﾏﾝﾄﾞｷﾞｮｳ                PIC X(80).
       01  WS-ｴﾗｰﾌｧｲﾙ                  PIC X(8).
       01  WS-ｴﾗｰｽﾃｰﾀｽ                 PIC XX.
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ                PIC 9(8).
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ            PIC X(7).
       01  WS-ﾎｳｺｳ                     PIC X.
       01  WS-ﾍﾝｶﾝｹｯｶ                  PIC 9.
       01  WS-ｼﾞﾕｳﾋﾞ                   PIC 9(8).
       01  WS-ｹｲﾔｸﾋﾞ                   PIC 9(8).
       01  WS-ｹｲｶﾆｯｽｳ                  PIC S9(7).
       01  WS-ﾆｯｽｳ                     PIC 9(3).
       01  WS-ﾆﾁｶﾞｸ                    PIC 9(9).
       01  WS-ｷﾝｶﾞｸ                    PIC S9(11).
       01  WS-ﾈﾝﾚｲｸﾌﾞﾝ                 PIC 9.
       01  WS-ｼﾊﾗｲｼｭﾍﾞﾂ                PIC X.
       COPY TBLWS.
       01  WS-ｶｲﾔｸ-TBL.
           05  WS-ｶｲﾔｸ-ｼｮｳﾋﾝ OCCURS 10 TIMES.
               10  WS-ｶｲﾔｸ-ｸﾌﾞﾝ OCCURS 6 TIMES.
                   15  WS-ｶｲﾔｸ-ｹｲｽｳ OCCURS 40 TIMES
                                           PIC 9(4)V9(3).
       01  WS-ｹﾝｽｳ.
           05  WS-ﾖﾐｺﾐｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ｼﾊﾗｲｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ｷｬｯｶｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ｼﾊﾗｲｷﾝｶﾞｸ            PIC 9(13) VALUE 0.
       01  WS-ﾘﾕｳ                      PIC 99.
           88  WS-ﾘﾕｳ-ﾅｼ                       VALUE 0.
           88  WS-ｼｭﾍﾞﾂ-ﾌｾｲ                    VALUE 1.
           88  WS-ｹｲﾔｸ-ﾅｼ                      VALUE 2.
           88  WS-ｼﾞﾕｳﾋﾞ-ﾌｾｲ                   VALUE 3.
           88  WS-ｹｲﾔｸﾋﾞ-ﾏｴ                    VALUE 4.
           88  WS-ﾎｼｮｳｶﾞｲ                      VALUE 5.
           88  WS-ﾆｭｳｲﾝ-ﾀｲｼｮｳｶﾞｲ               VALUE 6.
           88  WS-ﾆｯｽｳ-ﾌｾｲ                     VALUE 7.
           88  WS-ﾒﾝｾｷｷｶﾝ                      VALUE 8.
       01  WS-ﾘﾕｳ-ﾒｯｾｰｼﾞ.
           05  FILLER PIC X(20) VALUE "種別不正".
           05  FILLER PIC X(20) VALUE "契約なし".
           05  FILLER PIC X(20) VALUE "事由発生日不正".
           05  FILLER PIC X(20) VALUE "契約日前の事由".
           05  FILLER PIC X(20) VALUE "保障外".
           05  FILLER PIC X(20) VALUE "入院給付なし".
           05  FILLER PIC X(20) VALUE "日数不正".
           05  FILLER PIC X(20) VALUE "免責期間内".
       01  WS-ﾘﾕｳ-ﾒｯｾｰｼﾞ-TBL REDEFINES WS-ﾘﾕｳ-ﾒｯｾｰｼﾞ.
           05  WS-ﾒｯｾｰｼﾞ OCCURS 8 TIMES PIC X(20).
       01  WS-ｼｭﾍﾞﾂﾒｲ                  PIC X(12).
       01  WS-ﾌｯｷｺｰﾄﾞ                  PIC 99 VALUE 0.
       01  WS-ﾎｳｺｸ.
           05  WS-ﾋｮｳﾀﾞｲ-ｷﾞｮｳ.
               10  FILLER                  PIC X(28)
                   VALUE "ゴンドワナ生命  給付支払".
               10  FILLER                  PIC X(4) VALUE SPACES.
               10  WS-ﾋｮｳﾀﾞｲ-ﾋﾂﾞｹ          PIC X(7).
           05  WS-ﾒｲｻｲ-ｷﾞｮｳ.
               10  WS-ﾒｲｻｲ-ｾｲｷｭｳ           PIC X(8).
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
      * ﾒｲﾝ-RTN: 準備、1件ずつの支払、報告、終了。
      *----------------------------------------------------------------
       ﾒｲﾝ-RTN.
           PERFORM ｼｮｷ-RTN
           PERFORM ｼｮﾘ-RTN UNTIL WS-CLMREQ-ｵﾜﾘ
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
               DISPLAY "GLD030 日付不正: " WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8)
               MOVE 8 TO RETURN-CODE
               STOP RUN
           END-IF
           OPEN INPUT CLMREQ-F PRODTB-F CSVTB-F
           OPEN I-O POLMST-F GLSCTL-F
           OPEN EXTEND PAYOUT-F GLJRNL-F
           OPEN OUTPUT REJECT-F REPORT-F
           PERFORM ｵｰﾌﾟﾝ-ｹﾝｻ-RTN
           MOVE SPACES TO PY-ｼﾊﾗｲ-ﾚｺｰﾄﾞ JR-ｼﾜｹ-ﾚｺｰﾄﾞ
           PERFORM ｼｮｳﾋﾝ-ﾖﾐｺﾐ-RTN
           PERFORM ｶｲﾔｸ-ﾖﾐｺﾐ-RTN
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
           PERFORM CLMREQ-ﾖﾐｺﾐ-RTN.

      *----------------------------------------------------------------
      * ｵｰﾌﾟﾝ-ｹﾝｻ-RTN: 開けなかったファイルがあれば異常終了。
      *----------------------------------------------------------------
       ｵｰﾌﾟﾝ-ｹﾝｻ-RTN.
           EVALUATE TRUE
               WHEN NOT WS-CLMREQ-OK
                   MOVE "CLMREQ" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-CLMREQ-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-PRODTB-OK
                   MOVE "PRODTB" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-PRODTB-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-CSVTB-OK
                   MOVE "CSVTB" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-CSVTB-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
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
                   MOVE "CLMREJ" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-REJECT-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-REPORT-OK
                   MOVE "GLD030R" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-REPORT-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN OTHER
                   MOVE SPACES TO WS-ｴﾗｰﾌｧｲﾙ
           END-EVALUATE
           IF WS-ｴﾗｰﾌｧｲﾙ NOT = SPACES
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ｶｲﾔｸ-ﾖﾐｺﾐ-RTN: 解約返戻金係数を商品・年齢区分・経過年数の
      * テーブルに読み込む。商品テーブルにない商品は読み飛ばす。
      *----------------------------------------------------------------
       ｶｲﾔｸ-ﾖﾐｺﾐ-RTN.
           MOVE ZEROS TO WS-ｶｲﾔｸ-TBL
           READ CSVTB-F
           PERFORM UNTIL NOT WS-CSVTB-OK
               SET WS-ｼｮｳﾋﾝ-IX TO 1
               SEARCH WS-ｼｮｳﾋﾝ
                   WHEN WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = CV-ｼｮｳﾋﾝ
                       MOVE CV-ｹｲｽｳ TO WS-ｶｲﾔｸ-ｹｲｽｳ(WS-ｼｮｳﾋﾝ-IX,
                                                    CV-ﾈﾝﾚｲｸﾌﾞﾝ,
                                                    CV-ｹｲｶﾈﾝｽｳ)
               END-SEARCH
               READ CSVTB-F
           END-PERFORM
           IF NOT WS-CSVTB-ｵﾜﾘ
               MOVE "CSVTB" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-CSVTB-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * CLMREQ-ﾖﾐｺﾐ-RTN: 請求ファイルの次の行。
      *----------------------------------------------------------------
       CLMREQ-ﾖﾐｺﾐ-RTN.
           READ CLMREQ-F
           IF WS-CLMREQ-OK
               ADD 1 TO WS-ﾖﾐｺﾐｹﾝｽｳ
           ELSE
               IF NOT WS-CLMREQ-ｵﾜﾘ
                   MOVE "CLMREQ" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-CLMREQ-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
                   PERFORM ｲｼﾞｮｳ-RTN
               END-IF
           END-IF.

      *----------------------------------------------------------------
      * ｼｮﾘ-RTN: 1件を検査し、通れば支払、通らなければ却下。
      *----------------------------------------------------------------
       ｼｮﾘ-RTN.
           PERFORM ｹﾝｻ-RTN
           IF WS-ﾘﾕｳ-ﾅｼ
               PERFORM ｼﾊﾗｲ-RTN
           ELSE
               PERFORM ｷｬｯｶ-RTN
           END-IF
           PERFORM CLMREQ-ﾖﾐｺﾐ-RTN.

      *----------------------------------------------------------------
      * ｹﾝｻ-RTN: 請求に共通の検査、次いで種別ごとの検査。
      *----------------------------------------------------------------
       ｹﾝｻ-RTN.
           SET WS-ﾘﾕｳ-ﾅｼ TO TRUE
           MOVE 0 TO WS-ｷﾝｶﾞｸ
           EVALUATE TRUE
               WHEN CR-ｼﾎﾞｳ
                   MOVE "死亡保険金" TO WS-ｼｭﾍﾞﾂﾒｲ
               WHEN CR-ﾆｭｳｲﾝ
                   MOVE "入院給付金" TO WS-ｼｭﾍﾞﾂﾒｲ
               WHEN CR-ｶｲﾔｸ
                   MOVE "解約返戻金" TO WS-ｼｭﾍﾞﾂﾒｲ
               WHEN OTHER
                   MOVE "不明" TO WS-ｼｭﾍﾞﾂﾒｲ
                   SET WS-ｼｭﾍﾞﾂ-ﾌｾｲ TO TRUE
           END-EVALUATE
           IF WS-ﾘﾕｳ-ﾅｼ
               MOVE CR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
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
           END-IF
           IF WS-ﾘﾕｳ-ﾅｼ
               MOVE "S" TO WS-ﾎｳｺｳ
               CALL "GLC020" USING WS-ﾎｳｺｳ WS-ｼﾞﾕｳﾋﾞ
                                   CR-ｼﾞﾕｳﾊｯｾｲﾋﾞ WS-ﾍﾝｶﾝｹｯｶ
               IF WS-ﾍﾝｶﾝｹｯｶ NOT = 0
                   SET WS-ｼﾞﾕｳﾋﾞ-ﾌｾｲ TO TRUE
               END-IF
           END-IF
           IF WS-ﾘﾕｳ-ﾅｼ
               CALL "GLC020" USING WS-ﾎｳｺｳ WS-ｹｲﾔｸﾋﾞ
                                   PM-ｹｲﾔｸﾋﾞ WS-ﾍﾝｶﾝｹｯｶ
               IF WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ
                   SET WS-ｹｲﾔｸﾋﾞ-ﾏｴ TO TRUE
               END-IF
           END-IF
           IF WS-ﾘﾕｳ-ﾅｼ
               EVALUATE TRUE
                   WHEN CR-ｶｲﾔｸ AND NOT PM-ﾎｼｮｳﾁｭｳ AND NOT PM-ｼｯｺｳ
                       SET WS-ﾎｼｮｳｶﾞｲ TO TRUE
                   WHEN NOT CR-ｶｲﾔｸ AND NOT PM-ﾎｼｮｳﾁｭｳ
                       SET WS-ﾎｼｮｳｶﾞｲ TO TRUE
                   WHEN CR-ﾆｭｳｲﾝ
                       PERFORM ﾆｭｳｲﾝ-ｹﾝｻ-RTN
               END-EVALUATE
           END-IF.

      *----------------------------------------------------------------
      * ﾆｭｳｲﾝ-ｹﾝｻ-RTN: 医療保険であること、日数、免責期間。
      *----------------------------------------------------------------
       ﾆｭｳｲﾝ-ｹﾝｻ-RTN.
           SET WS-ｼｮｳﾋﾝ-IX TO 1
           SEARCH WS-ｼｮｳﾋﾝ
               WHEN WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = PM-ｼｮｳﾋﾝ
                   CONTINUE
           END-SEARCH
           COMPUTE WS-ｹｲｶﾆｯｽｳ = FUNCTION INTEGER-OF-DATE(WS-ｼﾞﾕｳﾋﾞ)
                             - FUNCTION INTEGER-OF-DATE(WS-ｹｲﾔｸﾋﾞ)
           EVALUATE TRUE
               WHEN NOT WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)
                   SET WS-ﾆｭｳｲﾝ-ﾀｲｼｮｳｶﾞｲ TO TRUE
               WHEN CR-ﾆｯｽｳ = 0
                   SET WS-ﾆｯｽｳ-ﾌｾｲ TO TRUE
               WHEN WS-ｹｲｶﾆｯｽｳ < WS-ｼｮｳﾋﾝ-ﾒﾝｾｷﾆｯｽｳ(WS-ｼｮｳﾋﾝ-IX)
                   SET WS-ﾒﾝｾｷｷｶﾝ TO TRUE
           END-EVALUATE.

      *----------------------------------------------------------------
      * ｼﾊﾗｲ-RTN: 種別ごとに金額を求め、支払を書き、契約を更新する。
      *----------------------------------------------------------------
       ｼﾊﾗｲ-RTN.
           EVALUATE TRUE
               WHEN CR-ｼﾎﾞｳ
                   PERFORM ｼﾎﾞｳ-RTN
               WHEN CR-ﾆｭｳｲﾝ
                   PERFORM ﾆｭｳｲﾝ-RTN
               WHEN CR-ｶｲﾔｸ
                   PERFORM ｶｲﾔｸ-RTN
           END-EVALUATE
           IF WS-ｷﾝｶﾞｸ > 0
               PERFORM ｼﾊﾗｲ-ｶｷｺﾐ-RTN
           END-IF
           IF NOT CR-ﾆｭｳｲﾝ
               MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO PM-ｼﾞｮｳﾀｲﾋﾂﾞｹ
           END-IF
           REWRITE PM-ｹｲﾔｸ-ﾚｺｰﾄﾞ
           IF NOT WS-POLMST-OK
               MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           ADD 1 TO WS-ｼﾊﾗｲｹﾝｽｳ
           MOVE "支払" TO WS-ﾒｲｻｲ-ｹｯｶ
           PERFORM ﾒｲｻｲ-RTN.

      *----------------------------------------------------------------
      * ｼﾎﾞｳ-RTN: 保険金額に配当累計を加え、貸付残高を引く。
      *----------------------------------------------------------------
       ｼﾎﾞｳ-RTN.
           COMPUTE WS-ｷﾝｶﾞｸ = PM-ﾎｹﾝｷﾝｶﾞｸ + PM-ﾊｲﾄｳﾙｲｹｲ
                           - PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ
           MOVE 5110 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
           MOVE PM-ﾎｹﾝｷﾝｶﾞｸ TO JR-ｷﾝｶﾞｸ
           SET JR-ｶﾘｶﾀ TO TRUE
           MOVE "死亡保険金" TO JR-ﾃｷﾖｳ
           PERFORM ｼﾜｹ-RTN
           PERFORM ﾊｲﾄｳ-ｶｼﾂｹ-ｼﾜｹ-RTN
           MOVE 0 TO PM-ﾊｲﾄｳﾙｲｹｲ PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ
           SET PM-ｼﾊﾗｲｽﾞﾐ TO TRUE
           SET PY-ｼﾎﾞｳ TO TRUE.

      *----------------------------------------------------------------
      * ﾆｭｳｲﾝ-RTN: 日額に日数を掛ける。日数は 120 日まで。
      *----------------------------------------------------------------
       ﾆｭｳｲﾝ-RTN.
           IF CR-ﾆｯｽｳ > 120
               MOVE 120 TO WS-ﾆｯｽｳ
           ELSE
               MOVE CR-ﾆｯｽｳ TO WS-ﾆｯｽｳ
           END-IF
           COMPUTE WS-ﾆﾁｶﾞｸ = PM-ﾎｹﾝｷﾝｶﾞｸ
               * WS-ｼｮｳﾋﾝ-ﾆｭｳｲﾝﾆﾁｶﾞｸﾘﾂ(WS-ｼｮｳﾋﾝ-IX) / 1000
           COMPUTE WS-ｷﾝｶﾞｸ = WS-ﾆﾁｶﾞｸ * WS-ﾆｯｽｳ
           MOVE 5120 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
           MOVE WS-ｷﾝｶﾞｸ TO JR-ｷﾝｶﾞｸ
           SET JR-ｶﾘｶﾀ TO TRUE
           MOVE "入院給付金" TO JR-ﾃｷﾖｳ
           PERFORM ｼﾜｹ-RTN
           SET PY-ﾆｭｳｲﾝ TO TRUE.

      *----------------------------------------------------------------
      * ｶｲﾔｸ-RTN: 係数から解約返戻金を求め、貸付残高を引く。0 以下
      * なら支払なし。
      *----------------------------------------------------------------
       ｶｲﾔｸ-RTN.
           SET WS-ｼｮｳﾋﾝ-IX TO 1
           SEARCH WS-ｼｮｳﾋﾝ
               WHEN WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = PM-ｼｮｳﾋﾝ
                   CONTINUE
           END-SEARCH
           EVALUATE TRUE
               WHEN PM-ｶﾆｭｳﾈﾝﾚｲ < 20
                   MOVE 1 TO WS-ﾈﾝﾚｲｸﾌﾞﾝ
               WHEN PM-ｶﾆｭｳﾈﾝﾚｲ < 60
                   COMPUTE WS-ﾈﾝﾚｲｸﾌﾞﾝ = PM-ｶﾆｭｳﾈﾝﾚｲ / 10
               WHEN OTHER
                   MOVE 6 TO WS-ﾈﾝﾚｲｸﾌﾞﾝ
           END-EVALUATE
           IF PM-ｹｲｶﾈﾝｽｳ > 40
               MOVE 40 TO WS-ﾆｯｽｳ
           ELSE
               MOVE PM-ｹｲｶﾈﾝｽｳ TO WS-ﾆｯｽｳ
           END-IF
           COMPUTE WS-ｷﾝｶﾞｸ = PM-ﾎｹﾝｷﾝｶﾞｸ / 1000
               * WS-ｶｲﾔｸ-ｹｲｽｳ(WS-ｼｮｳﾋﾝ-IX, WS-ﾈﾝﾚｲｸﾌﾞﾝ, WS-ﾆｯｽｳ)
               - PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ
           IF WS-ｷﾝｶﾞｸ > 0
               MOVE 5130 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
               COMPUTE JR-ｷﾝｶﾞｸ = WS-ｷﾝｶﾞｸ + PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ
               SET JR-ｶﾘｶﾀ TO TRUE
               MOVE "解約返戻金" TO JR-ﾃｷﾖｳ
               PERFORM ｼﾜｹ-RTN
               PERFORM ﾊｲﾄｳ-ｶｼﾂｹ-ｼﾜｹ-RTN
           ELSE
               MOVE 0 TO WS-ｷﾝｶﾞｸ
           END-IF
           MOVE 0 TO PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ
           SET PM-ｶｲﾔｸ TO TRUE
           SET PY-ｶｲﾔｸ TO TRUE.

      *----------------------------------------------------------------
      * ﾊｲﾄｳ-ｶｼﾂｹ-ｼﾜｹ-RTN: 配当累計の支払と貸付残高の回収の仕訳。
      *----------------------------------------------------------------
       ﾊｲﾄｳ-ｶｼﾂｹ-ｼﾜｹ-RTN.
           IF CR-ｼﾎﾞｳ AND PM-ﾊｲﾄｳﾙｲｹｲ > 0
               MOVE 2310 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
               MOVE PM-ﾊｲﾄｳﾙｲｹｲ TO JR-ｷﾝｶﾞｸ
               SET JR-ｶﾘｶﾀ TO TRUE
               MOVE "配当金支払" TO JR-ﾃｷﾖｳ
               PERFORM ｼﾜｹ-RTN
           END-IF
           IF PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0
               MOVE 1210 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
               MOVE PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ TO JR-ｷﾝｶﾞｸ
               SET JR-ｶｼｶﾀ TO TRUE
               MOVE "貸付金回収" TO JR-ﾃｷﾖｳ
               PERFORM ｼﾜｹ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ｼﾜｹ-RTN: 仕訳1行。科目、貸借、金額、摘要は呼ぶ側が入れる。
      *----------------------------------------------------------------
       ｼﾜｹ-RTN.
           MOVE CT-ｼﾞｶｲｼﾜｹﾊﾞﾝｺﾞｳ TO JR-ｼﾜｹﾊﾞﾝｺﾞｳ
           ADD 1 TO CT-ｼﾞｶｲｼﾜｹﾊﾞﾝｺﾞｳ
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO JR-ﾋﾂﾞｹ
           MOVE "GLD030" TO JR-ﾌﾟﾛｸﾞﾗﾑ
           MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO JR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           WRITE JR-ｼﾜｹ-ﾚｺｰﾄﾞ
           IF NOT WS-GLJRNL-OK
               MOVE "GLJRNL" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-GLJRNL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ｼﾊﾗｲ-ｶｷｺﾐ-RTN: 支払レコードを採番して書く。
      *----------------------------------------------------------------
       ｼﾊﾗｲ-ｶｷｺﾐ-RTN.
           MOVE CT-ｼﾞｶｲｼﾊﾗｲﾊﾞﾝｺﾞｳ TO PY-ｼﾊﾗｲﾊﾞﾝｺﾞｳ
           ADD 1 TO CT-ｼﾞｶｲｼﾊﾗｲﾊﾞﾝｺﾞｳ
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO PY-ｼﾊﾗｲﾋﾞ
           MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO PY-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           MOVE WS-ｷﾝｶﾞｸ TO PY-ｷﾝｶﾞｸ
           MOVE CR-ｳｹﾄﾘﾆﾝｶﾅ TO PY-ｳｹﾄﾘﾆﾝｶﾅ
           MOVE CR-ｾｲｷｭｳﾊﾞﾝｺﾞｳ TO PY-ｾｲｷｭｳﾊﾞﾝｺﾞｳ
           WRITE PY-ｼﾊﾗｲ-ﾚｺｰﾄﾞ
           ADD WS-ｷﾝｶﾞｸ TO WS-ｼﾊﾗｲｷﾝｶﾞｸ.

      *----------------------------------------------------------------
      * ｷｬｯｶ-RTN: 却下ファイルに理由を付けて書き、明細を報告する。
      *----------------------------------------------------------------
       ｷｬｯｶ-RTN.
           MOVE CR-ｾｲｷｭｳ-ﾚｺｰﾄﾞ TO ｷｬｯｶ-ﾃﾞｰﾀ
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
           MOVE CR-ｾｲｷｭｳﾊﾞﾝｺﾞｳ TO WS-ﾒｲｻｲ-ｾｲｷｭｳ
           MOVE CR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO WS-ﾒｲｻｲ-ｼｮｳｹﾝ
           MOVE WS-ｼｭﾍﾞﾂﾒｲ TO WS-ﾒｲｻｲ-ｼｭﾍﾞﾂ
           MOVE WS-ｷﾝｶﾞｸ TO WS-ﾒｲｻｲ-ｷﾝｶﾞｸ
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
           MOVE "支払件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｼﾊﾗｲｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "却下件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｷｬｯｶｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "支払金額" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｼﾊﾗｲｷﾝｶﾞｸ TO WS-ｹﾝｽｳ-ｱﾀｲ
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
           CLOSE CLMREQ-F POLMST-F PRODTB-F CSVTB-F GLSCTL-F
                 PAYOUT-F GLJRNL-F REJECT-F REPORT-F
           DISPLAY "GLD030 支払 " WS-ｼﾊﾗｲｹﾝｽｳ
                   " 却下 " WS-ｷｬｯｶｹﾝｽｳ.

      *----------------------------------------------------------------
      * ｲｼﾞｮｳ-RTN: ファイル異常は 16 で終わる。
      *----------------------------------------------------------------
       ｲｼﾞｮｳ-RTN.
           DISPLAY "GLD030 ファイル異常 " WS-ｴﾗｰﾌｧｲﾙ
                   " 状態 " WS-ｴﾗｰｽﾃｰﾀｽ
           MOVE 16 TO RETURN-CODE
           STOP RUN.

       COPY PRDLOD.
