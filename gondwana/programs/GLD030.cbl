      *----------------------------------------------------------------
      * GLD030: 給付支払
      *
      * 日次。支社からの給付請求ファイル CLMREQ を読み、死亡保険金、
      * 入院給付金、解約返戻金を支払う。
      * 検査は種別、契約の有無、事由発生日 (和暦で届く)、事由が契約
      * 日以後で業務日付以前であること、契約の状態。死亡・入院は
      * 事由発生日に保障中であったこと: いま保障中 (有効・猶予中)
      * か、失効中でも事由が失効日より前。解約は失効中でもよい。
      * 保障は保障マスタ COVMST から: 事由発生日に有効であった
      * 保障、つまりいま有効か、消滅日が事由より後の保障による。
      * 死亡: 医療以外の保障 (主契約と定期特約) の保険金額の合計。
      * そのような保障のない契約には死亡保障がない。これに配当累計を
      * 加え、貸付残高と未納保険料 (月払は払込済年月の翌月から
      * 事由の月まで、年払は有効でなければ1回分) を引いた額。支払済
      * とする。契約日または復活日から3年以内の自殺は免責で、
      * 保険金に代えて解約返戻金を支払う。
      * 入院: 医療の保障 (医療保険の主契約か入院特約) のある契約に
      * 限る。日数は 1 以上、1入院 120 日まで、通算 1,000 日まで。
      * その保障の契約日から免責日数の間の入院は支払わない。日額は
      * その保障の保険金額に入院日額率を掛けて 1,000 で割った額。
      * 契約はそのまま。
      * 手術: 入院と同じ検査。日数欄の倍率 (10, 20, 40) を日額に
      * 掛けた額。
      * 入院と手術は支払うたびに給付履歴 BENHST に残し、同じ事由
      * 発生日の同じ種別の請求は二重請求として却下する。
      * 解約: 解約返戻金 (GLC030: 商品・加入年齢区分・経過年数の
      * 係数を保険金額 1,000 円あたりに掛けた額) から貸付残高を
      * 引いた額。0 以下なら支払なしで解約とする。特約に返戻金は
      * ない。
      * 死亡と解約では、契約の有効な保障をすべて消滅か解約とする。
      * 貸付残高は業務日付までの利息 (GLC040) を繰り入れた上で、
      * 支払時に回収し、仕訳する。
      *
      * 復帰コード: 0 全件支払、4 却下あり、8 日付が日付でない、
      *             16 ファイル異常。
      *
      * 変更履歴
      * 日付       担当 変更
      * 1987-04-01 KT   GL-0001 初版。
      * 2001-07-02 MS   GL-0271 死亡保険金に配当累計を加える。
      * 2008-06-02 YN   GL-0455 入院給付金、免責期間。
      * 2012-09-03 YN   GL-0561 解約返戻金を GLC030 から。貸付利息。
      * 2026-11-16 YN   GL-0821 事由発生日時点の保障、未納保険料の
      *                         差引、自殺免責、医療保険の死亡請求。
      * 2026-11-16 YN   GL-0822 給付履歴、通算日数、手術給付金。
      * 2026-11-18 YN   GL-0825 保障マスタ: 特約の死亡保障と入院。
      * 2026-11-19 YN   GL-0828 仕訳に相手科目。
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
           SELECT COVMST-F ASSIGN TO "COVMST"
               ORGANIZATION IS INDEXED
               ACCESS MODE IS DYNAMIC
               RECORD KEY IS HS-ｷｰ
               FILE STATUS IS WS-COVMST-ST.
           SELECT PRODTB-F ASSIGN TO "PRODTB"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-PRODTB-ST.
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
           SELECT OPTIONAL BENHST-F ASSIGN TO "BENHST"
               ORGANIZATION IS INDEXED
               ACCESS MODE IS DYNAMIC
               RECORD KEY IS BH-ｷｰ
               FILE STATUS IS WS-BENHST-ST.
           SELECT REPORT-F ASSIGN TO "GLD030R"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-ST.

       DATA DIVISION.
       FILE SECTION.
       FD  CLMREQ-F.
       COPY CLMREQ.
       FD  POLMST-F.
       COPY POLMST.
       FD  COVMST-F.
       COPY COVMST.
       FD  PRODTB-F.
       COPY PRODTB.
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
       FD  BENHST-F.
       COPY BENHST.
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
           05  WS-COVMST-ST            PIC XX.
               88  WS-COVMST-OK                VALUE "00".
               88  WS-COVMST-ｵﾜﾘ               VALUE "10" "23".
           05  WS-PRODTB-ST            PIC XX.
               88  WS-PRODTB-OK                VALUE "00".
               88  WS-PRODTB-ｵﾜﾘ               VALUE "10".
           05  WS-CSVTB-ST             PIC XX.
               88  WS-CSVTB-OK                 VALUE "00".
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
           05  WS-BENHST-ST            PIC XX.
               88  WS-BENHST-OK                VALUE "00" "05".
               88  WS-BENHST-ｵﾜﾘ               VALUE "10" "23".
       01  WS-ｺﾏﾝﾄﾞｷﾞｮｳ                PIC X(80).
       01  WS-ｴﾗｰﾌｧｲﾙ                  PIC X(8).
       01  WS-ｴﾗｰｽﾃｰﾀｽ                 PIC XX.
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ                PIC 9(8).
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ            PIC X(7).
       01  WS-ﾎｳｺｳ                     PIC X.
       01  WS-ﾍﾝｶﾝｹｯｶ                  PIC 9.
       01  WS-ｼﾞﾕｳﾋﾞ                   PIC 9(8).
       01  WS-ｼﾞﾕｳﾋﾞ-ﾌﾞﾝｶｲ REDEFINES WS-ｼﾞﾕｳﾋﾞ.
           05  WS-ｼﾞﾕｳ-YYYY            PIC 9(4).
           05  WS-ｼﾞﾕｳ-MM              PIC 99.
           05  WS-ｼﾞﾕｳ-DD              PIC 99.
       01  WS-ｼｯｺｳﾋﾞ                   PIC 9(8).
       01  WS-ﾒﾝｾｷｷｹﾞﾝ                 PIC 9(8).
       01  WS-ﾒﾝｾｷ-ﾊﾝﾃｲ                PIC X.
           88  WS-ﾒﾝｾｷ-ｱﾘ                      VALUE "Y".
           88  WS-ﾒﾝｾｷ-ﾅｼ                      VALUE "N".
       01  WS-ﾊﾗｲｺﾐｽﾞﾐ.
           05  WS-ﾊﾗｲｺﾐ-YYYY           PIC 9(4).
           05  WS-ﾊﾗｲｺﾐ-MM             PIC 99.
       01  WS-ﾐﾉｳﾂｷｽｳ                  PIC S9(5).
       01  WS-ﾐﾉｳｶﾞｸ                   PIC 9(9).
       01  WS-ﾂｳｻﾝﾆｯｽｳ                 PIC 9(5).
       01  WS-ﾘﾚｷ-ﾚﾝﾊﾞﾝ                PIC 9(3).
       01  WS-ｹｲﾔｸﾋﾞ                   PIC 9(8).
       01  WS-ｹｲｶﾆｯｽｳ                  PIC S9(7).
       01  WS-ﾆｯｽｳ                     PIC 9(3).
       01  WS-ﾆﾁｶﾞｸ                    PIC 9(9).
       01  WS-ｷﾝｶﾞｸ                    PIC S9(11).
       01  WS-ﾍﾝﾚｲｷﾝ                   PIC 9(9).
       01  WS-ﾘｿｸ                      PIC 9(9).
       01  WS-ｼﾎﾞｳｷﾝｶﾞｸ                PIC 9(11).
       01  WS-ｼｮｳﾒﾂﾋﾞ                  PIC 9(8).
       01  WS-ﾎｼｮｳ-ﾊﾝﾃｲ                PIC X.
           88  WS-ﾎｼｮｳ-ｱﾘ                      VALUE "Y".
           88  WS-ﾎｼｮｳ-ﾅｼ                      VALUE "N".
       01  WS-ﾆｭｳｲﾝﾎｼｮｳ.
           05  WS-ﾆｭｳｲﾝﾎｼｮｳ-ﾊﾝﾃｲ       PIC X.
               88  WS-ﾆｭｳｲﾝﾎｼｮｳ-ｱﾘ             VALUE "Y".
               88  WS-ﾆｭｳｲﾝﾎｼｮｳ-ﾅｼ             VALUE "N".
           05  WS-ﾆｭｳｲﾝﾎｼｮｳ-ｷﾝｶﾞｸ      PIC 9(9).
           05  WS-ﾆｭｳｲﾝﾎｼｮｳ-ｹｲﾔｸﾋﾞ     PIC 9(8).
           05  WS-ﾆｭｳｲﾝﾎｼｮｳ-ｿｴｼﾞ       PIC 99.
       COPY TBLWS.
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
           88  WS-ｻｷﾋﾂﾞｹ                       VALUE 9.
           88  WS-ｼﾎﾞｳﾎｼｮｳ-ﾅｼ                  VALUE 10.
           88  WS-ﾊﾞｲﾘﾂ-ﾌｾｲ                    VALUE 11.
           88  WS-ﾆｼﾞｭｳｾｲｷｭｳ                   VALUE 12.
           88  WS-ﾂｳｻﾝ-ﾁｮｳｶ                    VALUE 13.
       01  WS-ﾘﾕｳ-ﾒｯｾｰｼﾞ.
           05  FILLER PIC X(20) VALUE "種別不正".
           05  FILLER PIC X(20) VALUE "契約なし".
           05  FILLER PIC X(20) VALUE "事由発生日不正".
           05  FILLER PIC X(20) VALUE "契約日前の事由".
           05  FILLER PIC X(20) VALUE "保障外".
           05  FILLER PIC X(20) VALUE "入院給付なし".
           05  FILLER PIC X(20) VALUE "日数不正".
           05  FILLER PIC X(20) VALUE "免責期間内".
           05  FILLER PIC X(20) VALUE "事由発生日が先日付".
           05  FILLER PIC X(20) VALUE "死亡保障なし".
           05  FILLER PIC X(20) VALUE "手術倍率不正".
           05  FILLER PIC X(20) VALUE "二重請求".
           05  FILLER PIC X(20) VALUE "通算日数超過".
       01  WS-ﾘﾕｳ-ﾒｯｾｰｼﾞ-TBL REDEFINES WS-ﾘﾕｳ-ﾒｯｾｰｼﾞ.
           05  WS-ﾒｯｾｰｼﾞ OCCURS 13 TIMES PIC X(20).
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
           OPEN INPUT CLMREQ-F PRODTB-F
           OPEN I-O POLMST-F COVMST-F GLSCTL-F BENHST-F
           OPEN EXTEND PAYOUT-F GLJRNL-F
           OPEN OUTPUT REJECT-F REPORT-F
           PERFORM ｵｰﾌﾟﾝ-ｹﾝｻ-RTN
           MOVE SPACES TO PY-ｼﾊﾗｲ-ﾚｺｰﾄﾞ JR-ｼﾜｹ-ﾚｺｰﾄﾞ
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
               WHEN NOT WS-POLMST-OK
                   MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-COVMST-OK
                   MOVE "COVMST" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-COVMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
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
               WHEN NOT WS-BENHST-OK
                   MOVE "BENHST" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-BENHST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
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
               WHEN CR-ｼｭｼﾞｭﾂ
                   MOVE "手術給付金" TO WS-ｼｭﾍﾞﾂﾒｲ
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
           IF WS-ﾘﾕｳ-ﾅｼ AND WS-ｼﾞﾕｳﾋﾞ > WS-ｷﾞｮｳﾑﾋﾂﾞｹ
               SET WS-ｻｷﾋﾂﾞｹ TO TRUE
           END-IF
           IF WS-ﾘﾕｳ-ﾅｼ
               EVALUATE TRUE
                   WHEN CR-ｶｲﾔｸ AND NOT PM-ﾎｼｮｳﾁｭｳ AND NOT PM-ｼｯｺｳ
                       SET WS-ﾎｼｮｳｶﾞｲ TO TRUE
                   WHEN CR-ｶｲﾔｸ
                       CONTINUE
                   WHEN OTHER
                       PERFORM ﾎｼｮｳ-ｹﾝｻ-RTN
               END-EVALUATE
           END-IF
           IF WS-ﾘﾕｳ-ﾅｼ AND NOT CR-ｶｲﾔｸ
               PERFORM ﾎｼｮｳ-ｿｳｻ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ﾎｼｮｳ-ｹﾝｻ-RTN: 事由発生日に保障中であったこと。いま保障中
      * か、失効中で事由が失効日より前。
      *----------------------------------------------------------------
       ﾎｼｮｳ-ｹﾝｻ-RTN.
           EVALUATE TRUE
               WHEN PM-ﾎｼｮｳﾁｭｳ
                   CONTINUE
               WHEN PM-ｼｯｺｳ
                   CALL "GLC020" USING WS-ﾎｳｺｳ WS-ｼｯｺｳﾋﾞ
                                       PM-ｼｯｺｳﾋﾞ WS-ﾍﾝｶﾝｹｯｶ
                   IF WS-ｼﾞﾕｳﾋﾞ NOT < WS-ｼｯｺｳﾋﾞ
                       SET WS-ﾎｼｮｳｶﾞｲ TO TRUE
                   END-IF
               WHEN OTHER
                   SET WS-ﾎｼｮｳｶﾞｲ TO TRUE
           END-EVALUATE.

      *----------------------------------------------------------------
      * ﾎｼｮｳ-ｿｳｻ-RTN: 契約の保障を保障番号順にたどり、事由発生日に
      * 有効であった保障から死亡保険金額の合計と入院の保障を求める。
      * 死亡保障のない契約の死亡請求は却下、入院と手術の請求は入院
      * の検査へ。
      *----------------------------------------------------------------
       ﾎｼｮｳ-ｿｳｻ-RTN.
           MOVE 0 TO WS-ｼﾎﾞｳｷﾝｶﾞｸ
           SET WS-ﾆｭｳｲﾝﾎｼｮｳ-ﾅｼ TO TRUE
           MOVE CR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO HS-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           MOVE 0 TO HS-ﾎｼｮｳﾊﾞﾝｺﾞｳ
           START COVMST-F KEY NOT < HS-ｷｰ
           IF WS-COVMST-OK
               READ COVMST-F NEXT
           END-IF
           PERFORM UNTIL NOT WS-COVMST-OK
                   OR HS-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ NOT = CR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
               PERFORM ﾎｼｮｳ-ﾊﾝﾃｲ-RTN
               READ COVMST-F NEXT
           END-PERFORM
           IF NOT WS-COVMST-OK AND NOT WS-COVMST-ｵﾜﾘ
               MOVE "COVMST" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-COVMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           EVALUATE TRUE
               WHEN CR-ｼﾎﾞｳ AND WS-ｼﾎﾞｳｷﾝｶﾞｸ = 0
                   SET WS-ｼﾎﾞｳﾎｼｮｳ-ﾅｼ TO TRUE
               WHEN CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ
                   PERFORM ﾆｭｳｲﾝ-ｹﾝｻ-RTN
           END-EVALUATE.

      *----------------------------------------------------------------
      * ﾎｼｮｳ-ﾊﾝﾃｲ-RTN: 保障1件。事由発生日に有効であったなら、医療
      * の保障は最初の1件を入院の保障とし、ほかは保険金額を死亡
      * 保険金額に加える。
      *----------------------------------------------------------------
       ﾎｼｮｳ-ﾊﾝﾃｲ-RTN.
           SET WS-ﾎｼｮｳ-ｱﾘ TO TRUE
           IF NOT HS-ﾕｳｺｳ
               CALL "GLC020" USING WS-ﾎｳｺｳ WS-ｼｮｳﾒﾂﾋﾞ
                                   HS-ｼｮｳﾒﾂﾋﾞ WS-ﾍﾝｶﾝｹｯｶ
               IF WS-ｼﾞﾕｳﾋﾞ NOT < WS-ｼｮｳﾒﾂﾋﾞ
                   SET WS-ﾎｼｮｳ-ﾅｼ TO TRUE
               END-IF
           END-IF
           IF WS-ﾎｼｮｳ-ｱﾘ
               SET WS-ｼｮｳﾋﾝ-IX TO 1
               SEARCH WS-ｼｮｳﾋﾝ
                   WHEN WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = HS-ｼｮｳﾋﾝ
                       CONTINUE
               END-SEARCH
               EVALUATE TRUE
                   WHEN NOT WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)
                       ADD HS-ﾎｹﾝｷﾝｶﾞｸ TO WS-ｼﾎﾞｳｷﾝｶﾞｸ
                   WHEN WS-ﾆｭｳｲﾝﾎｼｮｳ-ﾅｼ
                       SET WS-ﾆｭｳｲﾝﾎｼｮｳ-ｱﾘ TO TRUE
                       MOVE HS-ﾎｹﾝｷﾝｶﾞｸ TO WS-ﾆｭｳｲﾝﾎｼｮｳ-ｷﾝｶﾞｸ
                       SET WS-ﾆｭｳｲﾝﾎｼｮｳ-ｿｴｼﾞ TO WS-ｼｮｳﾋﾝ-IX
                       CALL "GLC020" USING WS-ﾎｳｺｳ
                                           WS-ﾆｭｳｲﾝﾎｼｮｳ-ｹｲﾔｸﾋﾞ
                                           HS-ｹｲﾔｸﾋﾞ WS-ﾍﾝｶﾝｹｯｶ
               END-EVALUATE
           END-IF.

      *----------------------------------------------------------------
      * ﾆｭｳｲﾝ-ｹﾝｻ-RTN: 入院の保障があること、日数か倍率、その保障の
      * 免責期間、次いで給付履歴。
      *----------------------------------------------------------------
       ﾆｭｳｲﾝ-ｹﾝｻ-RTN.
           IF WS-ﾆｭｳｲﾝﾎｼｮｳ-ﾅｼ
               SET WS-ﾆｭｳｲﾝ-ﾀｲｼｮｳｶﾞｲ TO TRUE
           ELSE
               COMPUTE WS-ｹｲｶﾆｯｽｳ
                   = FUNCTION INTEGER-OF-DATE(WS-ｼﾞﾕｳﾋﾞ)
                   - FUNCTION INTEGER-OF-DATE(WS-ﾆｭｳｲﾝﾎｼｮｳ-ｹｲﾔｸﾋﾞ)
               EVALUATE TRUE
                   WHEN CR-ﾆｭｳｲﾝ AND CR-ﾆｯｽｳ = 0
                       SET WS-ﾆｯｽｳ-ﾌｾｲ TO TRUE
                   WHEN CR-ｼｭｼﾞｭﾂ AND CR-ﾆｯｽｳ NOT = 10
                           AND CR-ﾆｯｽｳ NOT = 20 AND CR-ﾆｯｽｳ NOT = 40
                       SET WS-ﾊﾞｲﾘﾂ-ﾌｾｲ TO TRUE
                   WHEN WS-ｹｲｶﾆｯｽｳ
                           < WS-ｼｮｳﾋﾝ-ﾒﾝｾｷﾆｯｽｳ(WS-ﾆｭｳｲﾝﾎｼｮｳ-ｿｴｼﾞ)
                       SET WS-ﾒﾝｾｷｷｶﾝ TO TRUE
                   WHEN OTHER
                       PERFORM ﾘﾚｷ-ｹﾝｻ-RTN
               END-EVALUATE
           END-IF.

      *----------------------------------------------------------------
      * ﾘﾚｷ-ｹﾝｻ-RTN: 契約の給付履歴を連番順にたどり、入院の通算
      * 日数と最後の連番を求める。同じ種別で同じ事由発生日の履歴
      * があれば二重請求、入院で通算が 1,000 日に達していれば超過。
      *----------------------------------------------------------------
       ﾘﾚｷ-ｹﾝｻ-RTN.
           MOVE 0 TO WS-ﾂｳｻﾝﾆｯｽｳ WS-ﾘﾚｷ-ﾚﾝﾊﾞﾝ
           MOVE CR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO BH-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           MOVE 0 TO BH-ﾚﾝﾊﾞﾝ
           START BENHST-F KEY NOT < BH-ｷｰ
           IF WS-BENHST-OK
               READ BENHST-F NEXT
           END-IF
           PERFORM UNTIL NOT WS-BENHST-OK
                   OR BH-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ NOT = CR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
               MOVE BH-ﾚﾝﾊﾞﾝ TO WS-ﾘﾚｷ-ﾚﾝﾊﾞﾝ
               IF BH-ﾆｭｳｲﾝ
                   ADD BH-ﾆｯｽｳ TO WS-ﾂｳｻﾝﾆｯｽｳ
               END-IF
               IF BH-ｼｭﾍﾞﾂ = CR-ｼｭﾍﾞﾂ
                       AND BH-ｼﾞﾕｳﾊｯｾｲﾋﾞ = CR-ｼﾞﾕｳﾊｯｾｲﾋﾞ
                   SET WS-ﾆｼﾞｭｳｾｲｷｭｳ TO TRUE
               END-IF
               READ BENHST-F NEXT
           END-PERFORM
           IF NOT WS-BENHST-OK AND NOT WS-BENHST-ｵﾜﾘ
               MOVE "BENHST" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-BENHST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           IF WS-ﾘﾕｳ-ﾅｼ AND CR-ﾆｭｳｲﾝ AND WS-ﾂｳｻﾝﾆｯｽｳ NOT < 1000
               SET WS-ﾂｳｻﾝ-ﾁｮｳｶ TO TRUE
           END-IF.

      *----------------------------------------------------------------
      * ｼﾊﾗｲ-RTN: 種別ごとに金額を求め、支払を書き、契約を更新する。
      *----------------------------------------------------------------
       ｼﾊﾗｲ-RTN.
           EVALUATE TRUE
               WHEN CR-ｼﾎﾞｳ
                   PERFORM ｼﾎﾞｳ-RTN
               WHEN CR-ﾆｭｳｲﾝ
                   PERFORM ﾆｭｳｲﾝ-RTN
               WHEN CR-ｼｭｼﾞｭﾂ
                   PERFORM ｼｭｼﾞｭﾂ-RTN
               WHEN CR-ｶｲﾔｸ
                   PERFORM ｶｲﾔｸ-RTN
           END-EVALUATE
           IF WS-ｷﾝｶﾞｸ > 0
               PERFORM ｼﾊﾗｲ-ｶｷｺﾐ-RTN
           END-IF
           IF CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ
               PERFORM ﾘﾚｷ-ｶｷｺﾐ-RTN
           ELSE
               MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO PM-ｼﾞｮｳﾀｲﾋﾂﾞｹ
               PERFORM ﾎｼｮｳ-ｼｭｳﾘｮｳ-RTN
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
      * ｼﾎﾞｳ-RTN: 免責でなければ、死亡保険金額に配当累計を加え、
      * 貸付残高と未納保険料を引く。
      *----------------------------------------------------------------
       ｼﾎﾞｳ-RTN.
           PERFORM ﾘｿｸ-RTN
           PERFORM ﾒﾝｾｷ-ﾊﾝﾃｲ-RTN
           IF WS-ﾒﾝｾｷ-ｱﾘ
               PERFORM ﾒﾝｾｷ-RTN
           ELSE
               PERFORM ﾐﾉｳ-RTN
               COMPUTE WS-ｷﾝｶﾞｸ = WS-ｼﾎﾞｳｷﾝｶﾞｸ + PM-ﾊｲﾄｳﾙｲｹｲ
                               - PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ - WS-ﾐﾉｳｶﾞｸ
               MOVE 5110 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
               MOVE 1110 TO JR-ｱｲﾃｶﾓｸ
               MOVE WS-ｼﾎﾞｳｷﾝｶﾞｸ TO JR-ｷﾝｶﾞｸ
               SET JR-ｶﾘｶﾀ TO TRUE
               MOVE "死亡保険金" TO JR-ﾃｷﾖｳ
               PERFORM ｼﾜｹ-RTN
               PERFORM ﾊｲﾄｳ-ｶｼﾂｹ-ｼﾜｹ-RTN
               IF WS-ﾐﾉｳｶﾞｸ > 0
                   MOVE 4110 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
                   MOVE 1110 TO JR-ｱｲﾃｶﾓｸ
                   MOVE WS-ﾐﾉｳｶﾞｸ TO JR-ｷﾝｶﾞｸ
                   SET JR-ｶｼｶﾀ TO TRUE
                   MOVE "未納保険料充当" TO JR-ﾃｷﾖｳ
                   PERFORM ｼﾜｹ-RTN
               END-IF
           END-IF
           MOVE 0 TO PM-ﾊｲﾄｳﾙｲｹｲ PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ
           MOVE SPACES TO PM-ｶｼﾂｹｷｻﾝﾋﾞ
           SET PM-ｼﾊﾗｲｽﾞﾐ TO TRUE
           SET PY-ｼﾎﾞｳ TO TRUE.

      *----------------------------------------------------------------
      * ﾒﾝｾｷ-ﾊﾝﾃｲ-RTN: 自殺で、事由が契約日か復活日のうち遅い方から
      * 3年以内なら免責。
      *----------------------------------------------------------------
       ﾒﾝｾｷ-ﾊﾝﾃｲ-RTN.
           SET WS-ﾒﾝｾｷ-ﾅｼ TO TRUE
           IF CR-ｼﾞｻﾂ
               MOVE WS-ｹｲﾔｸﾋﾞ TO WS-ﾒﾝｾｷｷｹﾞﾝ
               IF PM-ﾌｯｶﾂﾋﾞ NOT = SPACES
                   CALL "GLC020" USING WS-ﾎｳｺｳ WS-ﾒﾝｾｷｷｹﾞﾝ
                                       PM-ﾌｯｶﾂﾋﾞ WS-ﾍﾝｶﾝｹｯｶ
               END-IF
               ADD 30000 TO WS-ﾒﾝｾｷｷｹﾞﾝ
               IF WS-ｼﾞﾕｳﾋﾞ < WS-ﾒﾝｾｷｷｹﾞﾝ
                   SET WS-ﾒﾝｾｷ-ｱﾘ TO TRUE
               END-IF
           END-IF.

      *----------------------------------------------------------------
      * ﾒﾝｾｷ-RTN: 保険金に代えて解約返戻金に配当累計を加え、貸付
      * 残高を引いて支払う。0 以下なら支払なし。
      *----------------------------------------------------------------
       ﾒﾝｾｷ-RTN.
           CALL "GLC030" USING PM-ｼｮｳﾋﾝ PM-ｶﾆｭｳﾈﾝﾚｲ PM-ｹｲｶﾈﾝｽｳ
                               PM-ﾎｹﾝｷﾝｶﾞｸ WS-ﾍﾝﾚｲｷﾝ WS-CSVTB-ST
           IF NOT WS-CSVTB-OK
               MOVE "CSVTB" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-CSVTB-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           COMPUTE WS-ｷﾝｶﾞｸ = WS-ﾍﾝﾚｲｷﾝ + PM-ﾊｲﾄｳﾙｲｹｲ
                           - PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ
           IF WS-ｷﾝｶﾞｸ > 0
               MOVE 5130 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
               MOVE 1110 TO JR-ｱｲﾃｶﾓｸ
               MOVE WS-ﾍﾝﾚｲｷﾝ TO JR-ｷﾝｶﾞｸ
               SET JR-ｶﾘｶﾀ TO TRUE
               MOVE "免責時返戻金" TO JR-ﾃｷﾖｳ
               PERFORM ｼﾜｹ-RTN
               PERFORM ﾊｲﾄｳ-ｶｼﾂｹ-ｼﾜｹ-RTN
           ELSE
               MOVE 0 TO WS-ｷﾝｶﾞｸ
           END-IF.

      *----------------------------------------------------------------
      * ﾐﾉｳ-RTN: 死亡保険金から引く未納保険料。月払は払込済年月の
      * 翌月から事由の月まで、年払と半年払は有効でなければ1回分。
      *----------------------------------------------------------------
       ﾐﾉｳ-RTN.
           MOVE PM-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ TO WS-ﾊﾗｲｺﾐｽﾞﾐ
           MOVE 0 TO WS-ﾐﾉｳﾂｷｽｳ
           IF NOT PM-ﾂｷﾊﾞﾗｲ
               IF NOT PM-ﾕｳｺｳ
                   MOVE 1 TO WS-ﾐﾉｳﾂｷｽｳ
               END-IF
           ELSE
               COMPUTE WS-ﾐﾉｳﾂｷｽｳ
                   = (WS-ｼﾞﾕｳ-YYYY - WS-ﾊﾗｲｺﾐ-YYYY) * 12
                   + WS-ｼﾞﾕｳ-MM - WS-ﾊﾗｲｺﾐ-MM
               IF WS-ﾐﾉｳﾂｷｽｳ < 0
                   MOVE 0 TO WS-ﾐﾉｳﾂｷｽｳ
               END-IF
           END-IF
           COMPUTE WS-ﾐﾉｳｶﾞｸ = PM-ﾎｹﾝﾘｮｳ * WS-ﾐﾉｳﾂｷｽｳ.

      *----------------------------------------------------------------
      * ﾘﾚｷ-ｶｷｺﾐ-RTN: 支払った入院か手術を給付履歴の次の連番に
      * 書く。
      *----------------------------------------------------------------
       ﾘﾚｷ-ｶｷｺﾐ-RTN.
           MOVE SPACES TO BH-ﾘﾚｷ-ﾚｺｰﾄﾞ
           MOVE CR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO BH-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           COMPUTE BH-ﾚﾝﾊﾞﾝ = WS-ﾘﾚｷ-ﾚﾝﾊﾞﾝ + 1
           MOVE CR-ｼｭﾍﾞﾂ TO BH-ｼｭﾍﾞﾂ
           MOVE CR-ｼﾞﾕｳﾊｯｾｲﾋﾞ TO BH-ｼﾞﾕｳﾊｯｾｲﾋﾞ
           MOVE WS-ﾆｯｽｳ TO BH-ﾆｯｽｳ
           MOVE WS-ｷﾝｶﾞｸ TO BH-ｷﾝｶﾞｸ
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO BH-ｼﾊﾗｲﾋﾞ
           MOVE CR-ｾｲｷｭｳﾊﾞﾝｺﾞｳ TO BH-ｾｲｷｭｳﾊﾞﾝｺﾞｳ
           WRITE BH-ﾘﾚｷ-ﾚｺｰﾄﾞ
           IF NOT WS-BENHST-OK
               MOVE "BENHST" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-BENHST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ﾆｭｳｲﾝ-RTN: 日額に日数を掛ける。日数は1入院 120 日まで、
      * 通算 1,000 日の残りまで。
      *----------------------------------------------------------------
       ﾆｭｳｲﾝ-RTN.
           IF CR-ﾆｯｽｳ > 120
               MOVE 120 TO WS-ﾆｯｽｳ
           ELSE
               MOVE CR-ﾆｯｽｳ TO WS-ﾆｯｽｳ
           END-IF
           IF WS-ﾂｳｻﾝﾆｯｽｳ + WS-ﾆｯｽｳ > 1000
               COMPUTE WS-ﾆｯｽｳ = 1000 - WS-ﾂｳｻﾝﾆｯｽｳ
           END-IF
           COMPUTE WS-ﾆﾁｶﾞｸ = WS-ﾆｭｳｲﾝﾎｼｮｳ-ｷﾝｶﾞｸ
               * WS-ｼｮｳﾋﾝ-ﾆｭｳｲﾝﾆﾁｶﾞｸﾘﾂ(WS-ﾆｭｳｲﾝﾎｼｮｳ-ｿｴｼﾞ) / 1000
           COMPUTE WS-ｷﾝｶﾞｸ = WS-ﾆﾁｶﾞｸ * WS-ﾆｯｽｳ
           MOVE 5120 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
           MOVE 1110 TO JR-ｱｲﾃｶﾓｸ
           MOVE WS-ｷﾝｶﾞｸ TO JR-ｷﾝｶﾞｸ
           SET JR-ｶﾘｶﾀ TO TRUE
           MOVE "入院給付金" TO JR-ﾃｷﾖｳ
           PERFORM ｼﾜｹ-RTN
           SET PY-ﾆｭｳｲﾝ TO TRUE.

      *----------------------------------------------------------------
      * ｼｭｼﾞｭﾂ-RTN: 日額に日数欄の倍率を掛ける。
      *----------------------------------------------------------------
       ｼｭｼﾞｭﾂ-RTN.
           MOVE CR-ﾆｯｽｳ TO WS-ﾆｯｽｳ
           COMPUTE WS-ﾆﾁｶﾞｸ = WS-ﾆｭｳｲﾝﾎｼｮｳ-ｷﾝｶﾞｸ
               * WS-ｼｮｳﾋﾝ-ﾆｭｳｲﾝﾆﾁｶﾞｸﾘﾂ(WS-ﾆｭｳｲﾝﾎｼｮｳ-ｿｴｼﾞ) / 1000
           COMPUTE WS-ｷﾝｶﾞｸ = WS-ﾆﾁｶﾞｸ * WS-ﾆｯｽｳ
           MOVE 5120 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
           MOVE 1110 TO JR-ｱｲﾃｶﾓｸ
           MOVE WS-ｷﾝｶﾞｸ TO JR-ｷﾝｶﾞｸ
           SET JR-ｶﾘｶﾀ TO TRUE
           MOVE "手術給付金" TO JR-ﾃｷﾖｳ
           PERFORM ｼﾜｹ-RTN
           SET PY-ｼｭｼﾞｭﾂ TO TRUE.

      *----------------------------------------------------------------
      * ｶｲﾔｸ-RTN: 解約返戻金を GLC030 で求め、貸付残高を引く。0 以下
      * なら支払なし。
      *----------------------------------------------------------------
       ｶｲﾔｸ-RTN.
           PERFORM ﾘｿｸ-RTN
           CALL "GLC030" USING PM-ｼｮｳﾋﾝ PM-ｶﾆｭｳﾈﾝﾚｲ PM-ｹｲｶﾈﾝｽｳ
                               PM-ﾎｹﾝｷﾝｶﾞｸ WS-ﾍﾝﾚｲｷﾝ WS-CSVTB-ST
           IF NOT WS-CSVTB-OK
               MOVE "CSVTB" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-CSVTB-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           COMPUTE WS-ｷﾝｶﾞｸ = WS-ﾍﾝﾚｲｷﾝ - PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ
           IF WS-ｷﾝｶﾞｸ > 0
               MOVE 5130 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
               MOVE 1110 TO JR-ｱｲﾃｶﾓｸ
               MOVE WS-ﾍﾝﾚｲｷﾝ TO JR-ｷﾝｶﾞｸ
               SET JR-ｶﾘｶﾀ TO TRUE
               MOVE "解約返戻金" TO JR-ﾃｷﾖｳ
               PERFORM ｼﾜｹ-RTN
               PERFORM ﾊｲﾄｳ-ｶｼﾂｹ-ｼﾜｹ-RTN
           ELSE
               MOVE 0 TO WS-ｷﾝｶﾞｸ
           END-IF
           MOVE 0 TO PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ
           MOVE SPACES TO PM-ｶｼﾂｹｷｻﾝﾋﾞ
           SET PM-ｶｲﾔｸ TO TRUE
           SET PY-ｶｲﾔｸ TO TRUE.

      *----------------------------------------------------------------
      * ﾎｼｮｳ-ｼｭｳﾘｮｳ-RTN: 死亡か解約で終わった契約の有効な保障を、
      * 死亡なら消滅、解約なら解約とする。
      *----------------------------------------------------------------
       ﾎｼｮｳ-ｼｭｳﾘｮｳ-RTN.
           MOVE CR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO HS-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           MOVE 0 TO HS-ﾎｼｮｳﾊﾞﾝｺﾞｳ
           START COVMST-F KEY NOT < HS-ｷｰ
           IF WS-COVMST-OK
               READ COVMST-F NEXT
           END-IF
           PERFORM UNTIL NOT WS-COVMST-OK
                   OR HS-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ NOT = CR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
               IF HS-ﾕｳｺｳ
                   IF CR-ｼﾎﾞｳ
                       SET HS-ｼｮｳﾒﾂ TO TRUE
                   ELSE
                       SET HS-ｶｲﾔｸ TO TRUE
                   END-IF
                   MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO HS-ｼｮｳﾒﾂﾋﾞ
                   REWRITE HS-ﾎｼｮｳ-ﾚｺｰﾄﾞ
               END-IF
               IF WS-COVMST-OK
                   READ COVMST-F NEXT
               END-IF
           END-PERFORM
           IF NOT WS-COVMST-OK AND NOT WS-COVMST-ｵﾜﾘ
               MOVE "COVMST" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-COVMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ﾘｿｸ-RTN: 業務日付までの貸付利息を残高に繰り入れる。
      *----------------------------------------------------------------
       ﾘｿｸ-RTN.
           CALL "GLC040" USING PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ PM-ｶｼﾂｹｷｻﾝﾋﾞ
                               WS-ｷﾞｮｳﾑﾋﾂﾞｹ WS-ﾘｿｸ
           IF WS-ﾘｿｸ > 0
               ADD WS-ﾘｿｸ TO PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ
               MOVE 4310 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
               MOVE 1210 TO JR-ｱｲﾃｶﾓｸ
               MOVE WS-ﾘｿｸ TO JR-ｷﾝｶﾞｸ
               SET JR-ｶｼｶﾀ TO TRUE
               MOVE "貸付金利息" TO JR-ﾃｷﾖｳ
               PERFORM ｼﾜｹ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ﾊｲﾄｳ-ｶｼﾂｹ-ｼﾜｹ-RTN: 配当累計の支払と貸付残高の回収の仕訳。
      *----------------------------------------------------------------
       ﾊｲﾄｳ-ｶｼﾂｹ-ｼﾜｹ-RTN.
           IF CR-ｼﾎﾞｳ AND PM-ﾊｲﾄｳﾙｲｹｲ > 0
               MOVE 2310 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
               MOVE 1110 TO JR-ｱｲﾃｶﾓｸ
               MOVE PM-ﾊｲﾄｳﾙｲｹｲ TO JR-ｷﾝｶﾞｸ
               SET JR-ｶﾘｶﾀ TO TRUE
               MOVE "配当金支払" TO JR-ﾃｷﾖｳ
               PERFORM ｼﾜｹ-RTN
           END-IF
           IF PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0
               MOVE 1210 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
               MOVE 1110 TO JR-ｱｲﾃｶﾓｸ
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
           CLOSE CLMREQ-F POLMST-F COVMST-F PRODTB-F GLSCTL-F
                 BENHST-F PAYOUT-F GLJRNL-F REJECT-F REPORT-F
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
