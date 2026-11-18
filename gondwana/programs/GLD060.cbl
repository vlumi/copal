      *----------------------------------------------------------------
      * GLD060: 失効処理
      *
      * 日次。猶予中の契約のうち、払込猶予期間を過ぎたものを失効
      * させる。猶予期間は猶予開始日 (振替不能となった日) の属する
      * 月の翌月末まで: 業務月が猶予開始月の翌月より後なら失効。
      * 失効した契約は失効日を業務日付とし、1年以内なら GLD040 の
      * 復活で戻せる。
      * ただし、未納保険料 (月払は払込済年月の翌月から猶予期限の
      * 月まで、年払と半年払は1回分) を貸し付けても貸付残高が解約返戻金
      * (GLC030) 以内に収まる契約は失効させず、未納保険料を自動
      * 振替貸付として残高に加え、払込済年月を進めて有効に戻す。
      * 残高には先に業務日付までの利息 (GLC040) を繰り入れる。
      *
      * 復帰コード: 0 処理済、8 日付が日付でない、16 ファイル異常。
      *
      * 変更履歴
      * 日付       担当 変更
      * 1987-04-01 KT   GL-0001 初版。
      * 2003-10-06 MS   GL-0310 猶予開始日を振替不能日から。
      * 2014-04-07 YN   GL-0598 自動振替貸付。
      * 2026-11-16 YN   GL-0823 半年払。
      * 2026-11-19 YN   GL-0827 収納を収納履歴に残す。
      *----------------------------------------------------------------
       IDENTIFICATION DIVISION.
       PROGRAM-ID. GLD060.
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
           SELECT GLSCTL-F ASSIGN TO "GLSCTL"
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-GLSCTL-ST.
           SELECT OPTIONAL GLJRNL-F ASSIGN TO "GLJRNL"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-GLJRNL-ST.
           SELECT REPORT-F ASSIGN TO "GLD060R"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-ST.

       DATA DIVISION.
       FILE SECTION.
       FD  POLMST-F.
       COPY POLMST.
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
           05  WS-GLSCTL-ST            PIC XX.
               88  WS-GLSCTL-OK                VALUE "00".
           05  WS-GLJRNL-ST            PIC XX.
               88  WS-GLJRNL-OK                VALUE "00" "05".
           05  WS-REPORT-ST            PIC XX.
               88  WS-REPORT-OK                VALUE "00".
           05  WS-CSVTB-ST             PIC XX.
               88  WS-CSVTB-OK                 VALUE "00".
       COPY PRMCAL.
       01  WS-ｺﾏﾝﾄﾞｷﾞｮｳ                PIC X(80).
       01  WS-ｴﾗｰﾌｧｲﾙ                  PIC X(8).
       01  WS-ｴﾗｰｽﾃｰﾀｽ                 PIC XX.
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ                PIC 9(8).
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾌﾞﾝｶｲ REDEFINES WS-ｷﾞｮｳﾑﾋﾂﾞｹ.
           05  WS-ｷﾞｮｳﾑ-YYYYMM         PIC 9(6).
           05  WS-ｷﾞｮｳﾑ-DD             PIC 99.
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ            PIC X(7).
       01  WS-ﾎｳｺｳ                     PIC X.
       01  WS-ﾍﾝｶﾝｹｯｶ                  PIC 9.
       01  WS-ﾕｳﾖｶｲｼﾋﾞ                 PIC 9(8).
       01  WS-ﾕｳﾖ-ﾌﾞﾝｶｲ REDEFINES WS-ﾕｳﾖｶｲｼﾋﾞ.
           05  WS-ﾕｳﾖ-YYYY             PIC 9(4).
           05  WS-ﾕｳﾖ-MM               PIC 99.
           05  WS-ﾕｳﾖ-DD               PIC 99.
       01  WS-ｷｹﾞﾝ.
           05  WS-ｷｹﾞﾝ-YYYY            PIC 9(4).
           05  WS-ｷｹﾞﾝ-MM              PIC 99.
       01  WS-ｷｹﾞﾝ-YYYYMM REDEFINES WS-ｷｹﾞﾝ PIC 9(6).
       01  WS-ﾊﾗｲｺﾐｽﾞﾐ.
           05  WS-ﾊﾗｲｺﾐ-YYYY           PIC 9(4).
           05  WS-ﾊﾗｲｺﾐ-MM             PIC 99.
       01  WS-ﾌﾘｶｴ-ｶｲｽｳ                PIC S9(5).
       01  WS-ﾌﾘｶｴ-ｶﾞｸ                 PIC 9(9).
       01  WS-ﾘｿｸ                      PIC 9(9).
       01  WS-ﾍﾝﾚｲｷﾝ                   PIC 9(9).
       01  WS-ｻﾞﾝﾀﾞｶ                   PIC 9(10).
       01  WS-ﾌﾘｶｴ-ﾊﾝﾃｲ                PIC X.
           88  WS-ﾌﾘｶｴ-ｶ                       VALUE "Y".
           88  WS-ﾌﾘｶｴ-ﾌｶ                      VALUE "N".
       01  WS-ｹﾝｽｳ.
           05  WS-ﾕｳﾖｹﾝｽｳ              PIC 9(7) VALUE 0.
           05  WS-ｼｯｺｳｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ｼｯｺｳｷﾝｶﾞｸ            PIC 9(13) VALUE 0.
           05  WS-ﾌﾘｶｴｹﾝｽｳ             PIC 9(7) VALUE 0.
           05  WS-ﾌﾘｶｴｷﾝｶﾞｸ            PIC 9(13) VALUE 0.
       01  WS-ﾎｳｺｸ.
           05  WS-ﾋｮｳﾀﾞｲ-ｷﾞｮｳ.
               10  FILLER                  PIC X(28)
                   VALUE "ゴンドワナ生命  失効処理".
               10  FILLER                  PIC X(4) VALUE SPACES.
               10  WS-ﾋｮｳﾀﾞｲ-ﾋﾂﾞｹ          PIC X(7).
           05  WS-ﾒｲｻｲ-ｷﾞｮｳ.
               10  WS-ﾒｲｻｲ-ｼｮｳｹﾝ           PIC X(10).
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  WS-ﾒｲｻｲ-ｼﾒｲ             PIC X(20).
               10  FILLER                  PIC X(2) VALUE SPACES.
               10  FILLER                  PIC X(8) VALUE "猶予開始".
               10  WS-ﾒｲｻｲ-ﾕｳﾖ             PIC X(7).
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
      * ｼｮｷ-RTN: 業務日付とファイル。
      *----------------------------------------------------------------
       ｼｮｷ-RTN.
           ACCEPT WS-ｺﾏﾝﾄﾞｷﾞｮｳ FROM COMMAND-LINE
           MOVE WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8) TO WS-ｷﾞｮｳﾑﾋﾂﾞｹ
           MOVE "W" TO WS-ﾎｳｺｳ
           CALL "GLC020" USING WS-ﾎｳｺｳ WS-ｷﾞｮｳﾑﾋﾂﾞｹ
                               WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ WS-ﾍﾝｶﾝｹｯｶ
           IF WS-ﾍﾝｶﾝｹｯｶ NOT = 0
               DISPLAY "GLD060 日付不正: " WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8)
               MOVE 8 TO RETURN-CODE
               STOP RUN
           END-IF
           OPEN I-O POLMST-F GLSCTL-F
           OPEN EXTEND GLJRNL-F
           OPEN OUTPUT REPORT-F
           PERFORM ｵｰﾌﾟﾝ-ｹﾝｻ-RTN
           MOVE SPACES TO JR-ｼﾜｹ-ﾚｺｰﾄﾞ
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
                   MOVE "GLD060R" TO WS-ｴﾗｰﾌｧｲﾙ
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
      * ｼｮﾘ-RTN: 猶予中の契約の猶予期限を求め、過ぎていれば自動
      * 振替貸付、できなければ失効。
      *----------------------------------------------------------------
       ｼｮﾘ-RTN.
           IF PM-ﾕｳﾖﾁｭｳ
               ADD 1 TO WS-ﾕｳﾖｹﾝｽｳ
               MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO WS-ﾒｲｻｲ-ｼｮｳｹﾝ
               MOVE PM-ｼﾒｲ TO WS-ﾒｲｻｲ-ｼﾒｲ
               MOVE PM-ﾕｳﾖｶｲｼﾋﾞ TO WS-ﾒｲｻｲ-ﾕｳﾖ
               MOVE 0 TO WS-ﾒｲｻｲ-ｷﾝｶﾞｸ
               PERFORM ｷｹﾞﾝ-RTN
               IF WS-ｷﾞｮｳﾑ-YYYYMM > WS-ｷｹﾞﾝ-YYYYMM
                   PERFORM ﾌﾘｶｴ-ﾊﾝﾃｲ-RTN
                   IF WS-ﾌﾘｶｴ-ｶ
                       PERFORM ﾌﾘｶｴ-RTN
                   ELSE
                       PERFORM ｼｯｺｳ-RTN
                   END-IF
               ELSE
                   MOVE "猶予中" TO WS-ﾒｲｻｲ-ｹｯｶ
               END-IF
               WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ﾒｲｻｲ-ｷﾞｮｳ
           END-IF
           PERFORM POLMST-ﾖﾐｺﾐ-RTN.

      *----------------------------------------------------------------
      * ｷｹﾞﾝ-RTN: 猶予開始月の翌月が猶予期限の月。
      *----------------------------------------------------------------
       ｷｹﾞﾝ-RTN.
           MOVE "S" TO WS-ﾎｳｺｳ
           CALL "GLC020" USING WS-ﾎｳｺｳ WS-ﾕｳﾖｶｲｼﾋﾞ
                               PM-ﾕｳﾖｶｲｼﾋﾞ WS-ﾍﾝｶﾝｹｯｶ
           MOVE WS-ﾕｳﾖ-YYYY TO WS-ｷｹﾞﾝ-YYYY
           MOVE WS-ﾕｳﾖ-MM TO WS-ｷｹﾞﾝ-MM
           IF WS-ｷｹﾞﾝ-MM = 12
               MOVE 1 TO WS-ｷｹﾞﾝ-MM
               ADD 1 TO WS-ｷｹﾞﾝ-YYYY
           ELSE
               ADD 1 TO WS-ｷｹﾞﾝ-MM
           END-IF.

      *----------------------------------------------------------------
      * ﾌﾘｶｴ-ﾊﾝﾃｲ-RTN: 未納保険料を求め、利息を繰り入れた貸付残高
      * に加えても解約返戻金以内なら自動振替貸付ができる。
      *----------------------------------------------------------------
       ﾌﾘｶｴ-ﾊﾝﾃｲ-RTN.
           MOVE PM-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ TO WS-ﾊﾗｲｺﾐｽﾞﾐ
           IF NOT PM-ﾂｷﾊﾞﾗｲ
               MOVE 1 TO WS-ﾌﾘｶｴ-ｶｲｽｳ
           ELSE
               COMPUTE WS-ﾌﾘｶｴ-ｶｲｽｳ
                   = (WS-ｷｹﾞﾝ-YYYY - WS-ﾊﾗｲｺﾐ-YYYY) * 12
                   + WS-ｷｹﾞﾝ-MM - WS-ﾊﾗｲｺﾐ-MM
           END-IF
           CALL "GLC040" USING PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ PM-ｶｼﾂｹｷｻﾝﾋﾞ
                               WS-ｷﾞｮｳﾑﾋﾂﾞｹ WS-ﾘｿｸ
           CALL "GLC030" USING PM-ｼｮｳﾋﾝ PM-ｶﾆｭｳﾈﾝﾚｲ PM-ｹｲｶﾈﾝｽｳ
                               PM-ﾎｹﾝｷﾝｶﾞｸ WS-ﾍﾝﾚｲｷﾝ WS-CSVTB-ST
           IF NOT WS-CSVTB-OK
               MOVE "CSVTB" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-CSVTB-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           SET WS-ﾌﾘｶｴ-ﾌｶ TO TRUE
           IF WS-ﾌﾘｶｴ-ｶｲｽｳ > 0
               COMPUTE WS-ﾌﾘｶｴ-ｶﾞｸ = PM-ﾎｹﾝﾘｮｳ * WS-ﾌﾘｶｴ-ｶｲｽｳ
               COMPUTE WS-ｻﾞﾝﾀﾞｶ = PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ + WS-ﾘｿｸ
                                 + WS-ﾌﾘｶｴ-ｶﾞｸ
               IF WS-ｻﾞﾝﾀﾞｶ NOT > WS-ﾍﾝﾚｲｷﾝ
                   SET WS-ﾌﾘｶｴ-ｶ TO TRUE
               END-IF
           END-IF.

      *----------------------------------------------------------------
      * ﾌﾘｶｴ-RTN: 未納保険料を貸付残高に加え、払込済年月を月払は
      * 猶予期限の月まで、年払は1年、半年払は半年進め、有効に
      * 戻して書き戻す。
      *----------------------------------------------------------------
       ﾌﾘｶｴ-RTN.
           IF WS-ﾘｿｸ > 0
               MOVE 4310 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
               SET JR-ｶｼｶﾀ TO TRUE
               MOVE WS-ﾘｿｸ TO JR-ｷﾝｶﾞｸ
               MOVE "貸付金利息" TO JR-ﾃｷﾖｳ
               PERFORM ｼﾜｹ-RTN
           END-IF
           MOVE 1210 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
           SET JR-ｶﾘｶﾀ TO TRUE
           MOVE WS-ﾌﾘｶｴ-ｶﾞｸ TO JR-ｷﾝｶﾞｸ
           MOVE "自動振替貸付" TO JR-ﾃｷﾖｳ
           PERFORM ｼﾜｹ-RTN
           MOVE 4110 TO JR-ｶﾝｼﾞｮｳｶﾓｸ
           SET JR-ｶｼｶﾀ TO TRUE
           MOVE WS-ﾌﾘｶｴ-ｶﾞｸ TO JR-ｷﾝｶﾞｸ
           MOVE "振替貸付保険料" TO JR-ﾃｷﾖｳ
           PERFORM ｼﾜｹ-RTN
           MOVE WS-ｻﾞﾝﾀﾞｶ TO PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO PM-ｶｼﾂｹｷｻﾝﾋﾞ PM-ｼﾞｮｳﾀｲﾋﾂﾞｹ
           EVALUATE TRUE
               WHEN PM-ﾈﾝﾊﾞﾗｲ
                   ADD 100 TO PM-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ
               WHEN PM-ﾊﾝﾄｼﾊﾞﾗｲ
                   ADD 6 TO WS-ﾊﾗｲｺﾐ-MM
                   IF WS-ﾊﾗｲｺﾐ-MM > 12
                       SUBTRACT 12 FROM WS-ﾊﾗｲｺﾐ-MM
                       ADD 1 TO WS-ﾊﾗｲｺﾐ-YYYY
                   END-IF
                   MOVE WS-ﾊﾗｲｺﾐｽﾞﾐ TO PM-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ
               WHEN OTHER
                   MOVE WS-ｷｹﾞﾝ-YYYYMM TO PM-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ
           END-EVALUATE
           MOVE 0 TO PM-ﾐﾉｳｶｲｽｳ
           MOVE SPACES TO PM-ﾕｳﾖｶｲｼﾋﾞ
           SET PM-ﾕｳｺｳ TO TRUE
           MOVE WS-ｷﾞｮｳﾑ-YYYYMM TO PM-ｻｲｼｭｳﾆｭｳｷﾝﾈﾝｹﾞﾂ
           MOVE WS-ﾌﾘｶｴ-ｶﾞｸ TO PM-ｻｲｼｭｳﾆｭｳｷﾝｶﾞｸ
           SET WS-ｼｭｳﾉｳ-ｶｼﾂｹ TO TRUE
           MOVE WS-ﾌﾘｶｴ-ｶﾞｸ TO WS-ｼｭｳﾉｳ-ｷﾝｶﾞｸ
           MOVE WS-ﾌﾘｶｴ-ｶｲｽｳ TO WS-ｼｭｳﾉｳ-ｶｲｽｳ
           PERFORM ｼｭｳﾉｳ-ｷﾛｸ-RTN
           PERFORM POLMST-ｶｷﾓﾄﾞｼ-RTN
           ADD 1 TO WS-ﾌﾘｶｴｹﾝｽｳ
           ADD WS-ﾌﾘｶｴ-ｶﾞｸ TO WS-ﾌﾘｶｴｷﾝｶﾞｸ
           MOVE WS-ﾌﾘｶｴ-ｶﾞｸ TO WS-ﾒｲｻｲ-ｷﾝｶﾞｸ
           MOVE "振替貸付" TO WS-ﾒｲｻｲ-ｹｯｶ.

      *----------------------------------------------------------------
      * ｼｯｺｳ-RTN: 失効にして書き戻す。
      *----------------------------------------------------------------
       ｼｯｺｳ-RTN.
           SET PM-ｼｯｺｳ TO TRUE
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO PM-ｼｯｺｳﾋﾞ PM-ｼﾞｮｳﾀｲﾋﾂﾞｹ
           PERFORM POLMST-ｶｷﾓﾄﾞｼ-RTN
           ADD 1 TO WS-ｼｯｺｳｹﾝｽｳ
           ADD PM-ﾎｹﾝｷﾝｶﾞｸ TO WS-ｼｯｺｳｷﾝｶﾞｸ
           MOVE "失効" TO WS-ﾒｲｻｲ-ｹｯｶ.

      *----------------------------------------------------------------
      * POLMST-ｶｷﾓﾄﾞｼ-RTN: 契約を書き戻す。
      *----------------------------------------------------------------
       POLMST-ｶｷﾓﾄﾞｼ-RTN.
           REWRITE PM-ｹｲﾔｸ-ﾚｺｰﾄﾞ
           IF NOT WS-POLMST-OK
               MOVE "POLMST" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-POLMST-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ｼﾜｹ-RTN: 仕訳1行。科目、貸借、金額、摘要は呼ぶ側が入れる。
      *----------------------------------------------------------------
       ｼﾜｹ-RTN.
           MOVE CT-ｼﾞｶｲｼﾜｹﾊﾞﾝｺﾞｳ TO JR-ｼﾜｹﾊﾞﾝｺﾞｳ
           ADD 1 TO CT-ｼﾞｶｲｼﾜｹﾊﾞﾝｺﾞｳ
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO JR-ﾋﾂﾞｹ
           MOVE "GLD060" TO JR-ﾌﾟﾛｸﾞﾗﾑ
           MOVE PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ TO JR-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ
           WRITE JR-ｼﾜｹ-ﾚｺｰﾄﾞ
           IF NOT WS-GLJRNL-OK
               MOVE "GLJRNL" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-GLJRNL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ﾎｳｺｸ-RTN: 件数。
      *----------------------------------------------------------------
       ﾎｳｺｸ-RTN.
           MOVE SPACES TO ﾎｳｺｸ-ｷﾞｮｳ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ
           MOVE "猶予中件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾕｳﾖｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "失効件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｼｯｺｳｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "失効保険金額" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｼｯｺｳｷﾝｶﾞｸ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "振替貸付件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾌﾘｶｴｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "振替貸付金額" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾌﾘｶｴｷﾝｶﾞｸ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ.

      *----------------------------------------------------------------
      * ｼｭｳﾘｮｳ-RTN: 制御レコードを書き戻し、ファイルを閉じ、件数を
      * 表示する。
      *----------------------------------------------------------------
       ｼｭｳﾘｮｳ-RTN.
           REWRITE CT-ｾｲｷﾞｮ-ﾚｺｰﾄﾞ
           IF NOT WS-GLSCTL-OK
               MOVE "GLSCTL" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-GLSCTL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF
           PERFORM ｼｭｳﾉｳ-ｸﾛｰｽﾞ-RTN
           CLOSE POLMST-F GLSCTL-F GLJRNL-F REPORT-F
           DISPLAY "GLD060 猶予中 " WS-ﾕｳﾖｹﾝｽｳ " 失効 " WS-ｼｯｺｳｹﾝｽｳ
                   " 振替貸付 " WS-ﾌﾘｶｴｹﾝｽｳ.

      *----------------------------------------------------------------
      * ｲｼﾞｮｳ-RTN: ファイル異常は 16 で終わる。
      *----------------------------------------------------------------
       ｲｼﾞｮｳ-RTN.
           DISPLAY "GLD060 ファイル異常 " WS-ｴﾗｰﾌｧｲﾙ
                   " 状態 " WS-ｴﾗｰｽﾃｰﾀｽ
           MOVE 16 TO RETURN-CODE
           STOP RUN.

       COPY PRMLOG.
