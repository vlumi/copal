      *----------------------------------------------------------------
      * GLM010: 保険料請求: 口座振替依頼ファイル作成
      *
      * 月次、振替日 (27日) の前。コマンド行の日付の月を請求月とし、
      * 払込中 (有効・猶予中) の契約のうち、月払は払込済年月が請求月
      * より前のもの、年払は契約月が請求月で払込済年月が請求月より
      * 前のものを選ぶ。振替額は保険料に未納の月数を合わせた回数分、
      * 月払は2回分まで、年払は1回分。
      * 収納代行の求める銀行・支店・口座番号順に整列し、銀行間の
      * 固定長レイアウトでヘッダ、データ、トレーラ、エンドを書く。
      * 銀行名と支店名は収納代行が埋めるので空白。
      * 制御レコードの最終請求年月を請求月にし、同じ月に二度請求
      * しない。
      *
      * 復帰コード: 0 作成済、8 日付が日付でないか請求済、
      *             16 ファイル異常。
      *
      * 変更履歴
      * 日付       担当 変更
      * 2003-10-06 MS   GL-0310 初版、口座振替の開始。
      * 2005-02-07 MS   GL-0371 未納分を合わせて2回分まで。
      *----------------------------------------------------------------
       IDENTIFICATION DIVISION.
       PROGRAM-ID. GLM010.
       AUTHOR. ゴンドワナ生命 システム部.
       DATE-WRITTEN. 2003-10-06.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT POLMST-F ASSIGN TO "POLMST"
               ORGANIZATION IS INDEXED
               ACCESS MODE IS SEQUENTIAL
               RECORD KEY IS PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
               FILE STATUS IS WS-POLMST-ST.
           SELECT SORT-F ASSIGN TO "SORTWK".
           SELECT PRMREQ-F ASSIGN TO "PRMREQ"
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-PRMREQ-ST.
           SELECT GLSCTL-F ASSIGN TO "GLSCTL"
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-GLSCTL-ST.
           SELECT REPORT-F ASSIGN TO "GLM010R"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-ST.

       DATA DIVISION.
       FILE SECTION.
       FD  POLMST-F.
       COPY POLMST.
       SD  SORT-F.
       01  SR-ｾｲﾚﾂ-ﾚｺｰﾄﾞ.
           05  SR-ｷﾞﾝｺｳｺｰﾄﾞ            PIC 9(4).
           05  SR-ｼﾃﾝｺｰﾄﾞ              PIC 9(3).
           05  SR-ｺｳｻﾞﾊﾞﾝｺﾞｳ           PIC 9(7).
           05  SR-ﾖｷﾝｼｭﾓｸ              PIC 9.
           05  SR-ｺｳｻﾞﾒｲｷﾞ             PIC X(30).
           05  SR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ          PIC X(10).
           05  SR-ｷﾝｶﾞｸ                PIC 9(10).
       FD  PRMREQ-F
           RECORD CONTAINS 120 CHARACTERS.
       COPY ZENGIN.
       FD  GLSCTL-F
           RECORD CONTAINS 64 CHARACTERS.
       COPY CTLREC.
       FD  REPORT-F.
       01  ﾎｳｺｸ-ｷﾞｮｳ                   PIC X(100).

       WORKING-STORAGE SECTION.
       01  WS-ｽﾃｰﾀｽ.
           05  WS-POLMST-ST            PIC XX.
               88  WS-POLMST-OK                VALUE "00".
               88  WS-POLMST-ｵﾜﾘ               VALUE "10".
           05  WS-PRMREQ-ST            PIC XX.
               88  WS-PRMREQ-OK                VALUE "00".
           05  WS-GLSCTL-ST            PIC XX.
               88  WS-GLSCTL-OK                VALUE "00".
           05  WS-REPORT-ST            PIC XX.
               88  WS-REPORT-OK                VALUE "00".
       01  WS-ｺﾏﾝﾄﾞｷﾞｮｳ                PIC X(80).
       01  WS-ｴﾗｰﾌｧｲﾙ                  PIC X(8).
       01  WS-ｴﾗｰｽﾃｰﾀｽ                 PIC XX.
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ                PIC 9(8).
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾌﾞﾝｶｲ REDEFINES WS-ｷﾞｮｳﾑﾋﾂﾞｹ.
           05  WS-ｾｲｷｭｳ-YYYYMM         PIC 9(6).
           05  WS-ｾｲｷｭｳ-YYYYMM-ﾌﾞﾝｶｲ REDEFINES WS-ｾｲｷｭｳ-YYYYMM.
               10  WS-ｾｲｷｭｳ-YYYY           PIC 9(4).
               10  WS-ｾｲｷｭｳ-MM             PIC 99.
           05  WS-ｷﾞｮｳﾑ-DD             PIC 99.
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ            PIC X(7).
       01  WS-ﾎｳｺｳ                     PIC X VALUE "W".
       01  WS-ﾍﾝｶﾝｹｯｶ                  PIC 9.
       01  WS-ﾋｷｵﾄｼﾋﾞ.
           05  WS-ﾋｷｵﾄｼ-MM             PIC 99.
           05  WS-ﾋｷｵﾄｼ-DD             PIC 99 VALUE 27.
       01  WS-ﾊﾗｲｺﾐｽﾞﾐ.
           05  WS-ﾊﾗｲｺﾐ-YYYY           PIC 9(4).
           05  WS-ﾊﾗｲｺﾐ-MM             PIC 99.
       01  WS-ｶｲｽｳ                     PIC S9(5).
       01  WS-ｿｳｷﾞﾝｺｳ                  PIC 9(4).
       01  WS-ｹﾝｽｳ.
           05  WS-ｾﾝﾀｸｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ｾｲｷｭｳｷﾝｶﾞｸ           PIC 9(13) VALUE 0.
           05  WS-ｷﾞﾝｺｳｹﾝｽｳ            PIC 9(7) VALUE 0.
           05  WS-ｷﾞﾝｺｳｷﾝｶﾞｸ           PIC 9(13) VALUE 0.
       01  WS-ﾌｯｷｺｰﾄﾞ                  PIC 99 VALUE 0.
       01  WS-ﾎｳｺｸ.
           05  WS-ﾋｮｳﾀﾞｲ-ｷﾞｮｳ.
               10  FILLER                  PIC X(28)
                   VALUE "ゴンドワナ生命  保険料請求".
               10  FILLER                  PIC X(4) VALUE SPACES.
               10  WS-ﾋｮｳﾀﾞｲ-ﾋﾂﾞｹ          PIC X(7).
               10  FILLER                  PIC X(6) VALUE "  請求".
               10  WS-ﾋｮｳﾀﾞｲ-ﾈﾝｹﾞﾂ         PIC 9(6).
           05  WS-ｷﾞﾝｺｳ-ｷﾞｮｳ.
               10  FILLER                  PIC X(8) VALUE "銀行".
               10  WS-ｷﾞﾝｺｳ-ｺｰﾄﾞ           PIC 9(4).
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ｷﾞﾝｺｳ-ｹﾝｽｳ           PIC ZZZ,ZZ9.
               10  FILLER                  PIC X(2) VALUE "件".
               10  WS-ｷﾞﾝｺｳ-ｷﾝｶﾞｸ          PIC ZZZ,ZZZ,ZZZ,ZZ9.
               10  FILLER                  PIC X(2) VALUE "円".
           05  WS-ｹﾝｽｳ-ｷﾞｮｳ.
               10  WS-ｹﾝｽｳ-ﾗﾍﾞﾙ            PIC X(16).
               10  WS-ｹﾝｽｳ-ｱﾀｲ             PIC ZZZ,ZZZ,ZZZ,ZZ9.

       PROCEDURE DIVISION.
      *----------------------------------------------------------------
      * ﾒｲﾝ-RTN: 準備、選択と整列と作成、制御レコード、報告、終了。
      *----------------------------------------------------------------
       ﾒｲﾝ-RTN.
           PERFORM ｼｮｷ-RTN
           SORT SORT-F
               ON ASCENDING KEY SR-ｷﾞﾝｺｳｺｰﾄﾞ SR-ｼﾃﾝｺｰﾄﾞ SR-ｺｳｻﾞﾊﾞﾝｺﾞｳ
               INPUT PROCEDURE IS ｾﾝﾀｸ-RTN
               OUTPUT PROCEDURE IS ｻｸｾｲ-RTN
           PERFORM ｾｲｷﾞｮ-RTN
           PERFORM ﾎｳｺｸ-RTN
           PERFORM ｼｭｳﾘｮｳ-RTN
           MOVE WS-ﾌｯｷｺｰﾄﾞ TO RETURN-CODE
           STOP RUN.

      *----------------------------------------------------------------
      * ｼｮｷ-RTN: 業務日付、ファイル、制御レコード、請求済の検査。
      *----------------------------------------------------------------
       ｼｮｷ-RTN.
           ACCEPT WS-ｺﾏﾝﾄﾞｷﾞｮｳ FROM COMMAND-LINE
           MOVE WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8) TO WS-ｷﾞｮｳﾑﾋﾂﾞｹ
           CALL "GLC020" USING WS-ﾎｳｺｳ WS-ｷﾞｮｳﾑﾋﾂﾞｹ
                               WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ WS-ﾍﾝｶﾝｹｯｶ
           IF WS-ﾍﾝｶﾝｹｯｶ NOT = 0
               DISPLAY "GLM010 日付不正: " WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8)
               MOVE 8 TO RETURN-CODE
               STOP RUN
           END-IF
           OPEN INPUT POLMST-F
           OPEN I-O GLSCTL-F
           OPEN OUTPUT PRMREQ-F REPORT-F
           EVALUATE TRUE
               WHEN NOT WS-POLMST-OK
                   MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-GLSCTL-OK
                   MOVE "GLSCTL" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-GLSCTL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-PRMREQ-OK
                   MOVE "PRMREQ" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-PRMREQ-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN NOT WS-REPORT-OK
                   MOVE "GLM010R" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-REPORT-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               WHEN OTHER
                   MOVE SPACES TO WS-ｴﾗｰﾌｧｲﾙ
           END-EVALUATE
           IF WS-ｴﾗｰﾌｧｲﾙ NOT = SPACES
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           READ GLSCTL-F
           IF NOT WS-GLSCTL-OK
               MOVE "GLSCTL" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-GLSCTL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           IF CT-ｻｲｼｭｳｾｲｷｭｳﾈﾝｹﾞﾂ = WS-ｾｲｷｭｳ-YYYYMM
               DISPLAY "GLM010 請求済: " WS-ｾｲｷｭｳ-YYYYMM
               MOVE 8 TO RETURN-CODE
               STOP RUN
           END-IF
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO WS-ﾋｮｳﾀﾞｲ-ﾋﾂﾞｹ
           MOVE WS-ｾｲｷｭｳ-YYYYMM TO WS-ﾋｮｳﾀﾞｲ-ﾈﾝｹﾞﾂ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ﾋｮｳﾀﾞｲ-ｷﾞｮｳ
           MOVE SPACES TO ﾎｳｺｸ-ｷﾞｮｳ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ.

      *----------------------------------------------------------------
      * ｾﾝﾀｸ-RTN: マスタを読み、請求する契約を整列に渡す。
      *----------------------------------------------------------------
       ｾﾝﾀｸ-RTN.
           READ POLMST-F NEXT
           PERFORM UNTIL NOT WS-POLMST-OK
               PERFORM ﾊﾝﾃｲ-RTN
               IF WS-ｶｲｽｳ > 0
                   PERFORM ﾜﾀｼ-RTN
               END-IF
               READ POLMST-F NEXT
           END-PERFORM
           IF NOT WS-POLMST-ｵﾜﾘ
               MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ﾊﾝﾃｲ-RTN: 請求する回数。請求しない契約は 0。
      *----------------------------------------------------------------
       ﾊﾝﾃｲ-RTN.
           MOVE 0 TO WS-ｶｲｽｳ
           IF PM-ﾊﾗｲｺﾐﾁｭｳ AND PM-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ < WS-ｾｲｷｭｳ-YYYYMM
               MOVE PM-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ TO WS-ﾊﾗｲｺﾐｽﾞﾐ
               EVALUATE TRUE
                   WHEN PM-ﾈﾝﾊﾞﾗｲ
                       IF PM-ｹｲﾔｸﾋﾞ(4:2) = WS-ｾｲｷｭｳ-YYYYMM(5:2)
                           MOVE 1 TO WS-ｶｲｽｳ
                       END-IF
                   WHEN OTHER
                       COMPUTE WS-ｶｲｽｳ
                           = (WS-ｾｲｷｭｳ-YYYY - WS-ﾊﾗｲｺﾐ-YYYY) * 12
                           + WS-ｾｲｷｭｳ-MM - WS-ﾊﾗｲｺﾐ-MM
                       IF WS-ｶｲｽｳ > 2
                           MOVE 2 TO WS-ｶｲｽｳ
                       END-IF
               END-EVALUATE
           END-IF.

      *----------------------------------------------------------------
      * ﾜﾀｼ-RTN: 整列レコードを作って RELEASE する。
      *----------------------------------------------------------------
       ﾜﾀｼ-RTN.
           MOVE PM-ｷﾞﾝｺｳｺｰﾄﾞ TO SR-ｷﾞﾝｺｳｺｰﾄﾞ
           MOVE PM-ｼﾃﾝｺｰﾄﾞ TO SR-ｼﾃﾝｺｰﾄﾞ
           MOVE PM-ｺｳｻﾞﾊﾞﾝｺﾞｳ TO SR-ｺｳｻﾞﾊﾞﾝｺﾞｳ
           MOVE PM-ﾖｷﾝｼｭﾓｸ TO SR-ﾖｷﾝｼｭﾓｸ
           MOVE PM-ｺｳｻﾞﾒｲｷﾞ TO SR-ｺｳｻﾞﾒｲｷﾞ
           MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO SR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           COMPUTE SR-ｷﾝｶﾞｸ = PM-ﾎｹﾝﾘｮｳ * WS-ｶｲｽｳ
           RELEASE SR-ｾｲﾚﾂ-ﾚｺｰﾄﾞ
           ADD 1 TO WS-ｾﾝﾀｸｹﾝｽｳ
           ADD SR-ｷﾝｶﾞｸ TO WS-ｾｲｷｭｳｷﾝｶﾞｸ.

      *----------------------------------------------------------------
      * ｻｸｾｲ-RTN: ヘッダ、整列済みのデータ、トレーラ、エンドを書く。
      * 銀行が変わるごとに報告に1行。
      *----------------------------------------------------------------
       ｻｸｾｲ-RTN.
           PERFORM ﾍｯﾀﾞ-RTN
           MOVE 0 TO WS-ｿｳｷﾞﾝｺｳ
           RETURN SORT-F
               AT END
                   CONTINUE
               NOT AT END
                   PERFORM ﾃﾞｰﾀ-RTN
                   RETURN SORT-F
                       AT END
                           CONTINUE
                   END-RETURN
                   PERFORM UNTIL SR-ｷﾞﾝｺｳｺｰﾄﾞ = 0
                       PERFORM ﾃﾞｰﾀ-RTN
                       MOVE 0 TO SR-ｷﾞﾝｺｳｺｰﾄﾞ
                       RETURN SORT-F
                           AT END
                               MOVE 0 TO SR-ｷﾞﾝｺｳｺｰﾄﾞ
                       END-RETURN
                   END-PERFORM
           END-RETURN
           IF WS-ｷﾞﾝｺｳｹﾝｽｳ > 0
               PERFORM ｷﾞﾝｺｳ-ｷﾞｮｳ-RTN
           END-IF
           PERFORM ﾄﾚｰﾗ-RTN.

      *----------------------------------------------------------------
      * ﾍｯﾀﾞ-RTN: 委託者と引落日。
      *----------------------------------------------------------------
       ﾍｯﾀﾞ-RTN.
           MOVE SPACES TO ZG-ﾌﾘｶｴ-ﾚｺｰﾄﾞ
           SET ZG-ﾍｯﾀﾞ TO TRUE
           MOVE 91 TO ZG-H-ｼｭﾍﾞﾂｺｰﾄﾞ
           MOVE 0 TO ZG-H-ｺｰﾄﾞｸﾌﾞﾝ
           MOVE 1234567890 TO ZG-H-ｲﾀｸｼｬｺｰﾄﾞ
           MOVE "ｺﾞﾝﾄﾞﾜﾅｾｲﾒｲﾎｹﾝ" TO ZG-H-ｲﾀｸｼｬﾒｲ
           MOVE WS-ｾｲｷｭｳ-MM TO WS-ﾋｷｵﾄｼ-MM
           MOVE WS-ﾋｷｵﾄｼﾋﾞ TO ZG-H-ﾋｷｵﾄｼﾋﾞ
           MOVE 1 TO ZG-H-ｷﾞﾝｺｳｺｰﾄﾞ
           MOVE SPACES TO ZG-H-ｷﾞﾝｺｳﾒｲ ZG-H-ｼﾃﾝﾒｲ
           MOVE 1 TO ZG-H-ｼﾃﾝｺｰﾄﾞ
           MOVE 1 TO ZG-H-ﾖｷﾝｼｭﾓｸ
           MOVE 1000001 TO ZG-H-ｺｳｻﾞﾊﾞﾝｺﾞｳ
           PERFORM PRMREQ-ｶｷｺﾐ-RTN.

      *----------------------------------------------------------------
      * ﾃﾞｰﾀ-RTN: 1件の振替依頼。銀行が変わったら前の銀行の行。
      *----------------------------------------------------------------
       ﾃﾞｰﾀ-RTN.
           IF SR-ｷﾞﾝｺｳｺｰﾄﾞ NOT = WS-ｿｳｷﾞﾝｺｳ
               IF WS-ｷﾞﾝｺｳｹﾝｽｳ > 0
                   PERFORM ｷﾞﾝｺｳ-ｷﾞｮｳ-RTN
               END-IF
               MOVE SR-ｷﾞﾝｺｳｺｰﾄﾞ TO WS-ｿｳｷﾞﾝｺｳ
           END-IF
           MOVE SPACES TO ZG-ﾌﾘｶｴ-ﾚｺｰﾄﾞ
           SET ZG-ﾃﾞｰﾀ TO TRUE
           MOVE SR-ｷﾞﾝｺｳｺｰﾄﾞ TO ZG-D-ｷﾞﾝｺｳｺｰﾄﾞ
           MOVE SR-ｼﾃﾝｺｰﾄﾞ TO ZG-D-ｼﾃﾝｺｰﾄﾞ
           MOVE SR-ﾖｷﾝｼｭﾓｸ TO ZG-D-ﾖｷﾝｼｭﾓｸ
           MOVE SR-ｺｳｻﾞﾊﾞﾝｺﾞｳ TO ZG-D-ｺｳｻﾞﾊﾞﾝｺﾞｳ
           MOVE SR-ｺｳｻﾞﾒｲｷﾞ TO ZG-D-ｺｳｻﾞﾒｲｷﾞ
           MOVE SR-ｷﾝｶﾞｸ TO ZG-D-ｷﾝｶﾞｸ
           MOVE 0 TO ZG-D-ｼﾝｷｺｰﾄﾞ
           MOVE SR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO ZG-D-ｺｷｬｸﾊﾞﾝｺﾞｳ
           SET ZG-D-ﾐｼｮﾘ TO TRUE
           PERFORM PRMREQ-ｶｷｺﾐ-RTN
           ADD 1 TO WS-ｷﾞﾝｺｳｹﾝｽｳ
           ADD SR-ｷﾝｶﾞｸ TO WS-ｷﾞﾝｺｳｷﾝｶﾞｸ.

      *----------------------------------------------------------------
      * ｷﾞﾝｺｳ-ｷﾞｮｳ-RTN: 銀行ごとの件数と金額を報告し、数え直す。
      *----------------------------------------------------------------
       ｷﾞﾝｺｳ-ｷﾞｮｳ-RTN.
           MOVE WS-ｿｳｷﾞﾝｺｳ TO WS-ｷﾞﾝｺｳ-ｺｰﾄﾞ
           MOVE WS-ｷﾞﾝｺｳｹﾝｽｳ TO WS-ｷﾞﾝｺｳ-ｹﾝｽｳ
           MOVE WS-ｷﾞﾝｺｳｷﾝｶﾞｸ TO WS-ｷﾞﾝｺｳ-ｷﾝｶﾞｸ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｷﾞﾝｺｳ-ｷﾞｮｳ
           MOVE 0 TO WS-ｷﾞﾝｺｳｹﾝｽｳ WS-ｷﾞﾝｺｳｷﾝｶﾞｸ.

      *----------------------------------------------------------------
      * ﾄﾚｰﾗ-RTN: 合計のトレーラとエンド記録。
      *----------------------------------------------------------------
       ﾄﾚｰﾗ-RTN.
           MOVE SPACES TO ZG-ﾌﾘｶｴ-ﾚｺｰﾄﾞ
           SET ZG-ﾄﾚｰﾗ TO TRUE
           MOVE WS-ｾﾝﾀｸｹﾝｽｳ TO ZG-T-ｺﾞｳｹｲｹﾝｽｳ
           MOVE WS-ｾｲｷｭｳｷﾝｶﾞｸ TO ZG-T-ｺﾞｳｹｲｷﾝｶﾞｸ
           MOVE 0 TO ZG-T-ﾌﾘｶｴｽﾞﾐｹﾝｽｳ ZG-T-ﾌﾘｶｴｽﾞﾐｷﾝｶﾞｸ
                     ZG-T-ﾌﾘｶｴﾌﾉｳｹﾝｽｳ ZG-T-ﾌﾘｶｴﾌﾉｳｷﾝｶﾞｸ
           PERFORM PRMREQ-ｶｷｺﾐ-RTN
           MOVE SPACES TO ZG-ﾌﾘｶｴ-ﾚｺｰﾄﾞ
           SET ZG-ｴﾝﾄﾞ TO TRUE
           PERFORM PRMREQ-ｶｷｺﾐ-RTN.

      *----------------------------------------------------------------
      * PRMREQ-ｶｷｺﾐ-RTN: 依頼ファイルに1記録。
      *----------------------------------------------------------------
       PRMREQ-ｶｷｺﾐ-RTN.
           WRITE ZG-ﾌﾘｶｴ-ﾚｺｰﾄﾞ
           IF NOT WS-PRMREQ-OK
               MOVE "PRMREQ" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-PRMREQ-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ｾｲｷﾞｮ-RTN: 最終請求年月を書き戻す。
      *----------------------------------------------------------------
       ｾｲｷﾞｮ-RTN.
           MOVE WS-ｾｲｷｭｳ-YYYYMM TO CT-ｻｲｼｭｳｾｲｷｭｳﾈﾝｹﾞﾂ
           REWRITE CT-ｾｲｷﾞｮ-ﾚｺｰﾄﾞ
           IF NOT WS-GLSCTL-OK
               MOVE "GLSCTL" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-GLSCTL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ﾎｳｺｸ-RTN: 件数と合計。
      *----------------------------------------------------------------
       ﾎｳｺｸ-RTN.
           MOVE SPACES TO ﾎｳｺｸ-ｷﾞｮｳ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ
           MOVE "請求件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｾﾝﾀｸｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "請求金額" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｾｲｷｭｳｷﾝｶﾞｸ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ.

      *----------------------------------------------------------------
      * ｼｭｳﾘｮｳ-RTN: ファイルを閉じ、件数を表示する。
      *----------------------------------------------------------------
       ｼｭｳﾘｮｳ-RTN.
           CLOSE POLMST-F PRMREQ-F GLSCTL-F REPORT-F
           DISPLAY "GLM010 請求 " WS-ｾﾝﾀｸｹﾝｽｳ " 件 "
                   WS-ｾｲｷｭｳｷﾝｶﾞｸ " 円".

      *----------------------------------------------------------------
      * ｲｼﾞｮｳ-RTN: ファイル異常は 16 で終わる。
      *----------------------------------------------------------------
       ｲｼﾞｮｳ-RTN.
           DISPLAY "GLM010 ファイル異常 " WS-ｴﾗｰﾌｧｲﾙ
                   " 状態 " WS-ｴﾗｰｽﾃｰﾀｽ
           MOVE 16 TO RETURN-CODE
           STOP RUN.
