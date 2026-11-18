      *----------------------------------------------------------------
      * GLD080: 日計表
      *
      * 日次、月次、年次の最後。その日の仕訳 GLJRNL を読み、勘定科目
      * ごとに借方と貸方を集計して日計表を印字する。仕訳1行は勘定
      * 科目を貸借の側に、相手科目を反対の側に同じ金額で記帳する
      * ので、借方合計と貸方合計は必ず合う。同じ日に日次、月次、
      * 年次と続けて実行すると、最後のものがその日の日計表になる。
      * 相手科目のない仕訳、科目表にない科目は数え、4 で終わる。
      *
      * 復帰コード: 0 作成済、4 相手科目のない仕訳か未登録の科目
      *             あり、8 日付が日付でない、16 ファイル異常。
      *
      * 変更履歴
      * 日付       担当 変更
      * 2026-11-19 YN   GL-0828 初版。
      *----------------------------------------------------------------
       IDENTIFICATION DIVISION.
       PROGRAM-ID. GLD080.
       AUTHOR. ゴンドワナ生命 システム部.
       DATE-WRITTEN. 2026-11-19.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT OPTIONAL GLJRNL-F ASSIGN TO "GLJRNL"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-GLJRNL-ST.
           SELECT REPORT-F ASSIGN TO "GLD080R"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-ST.

       DATA DIVISION.
       FILE SECTION.
       FD  GLJRNL-F.
       COPY GLJRNL.
       FD  REPORT-F.
       01  ﾎｳｺｸ-ｷﾞｮｳ                   PIC X(80).

       WORKING-STORAGE SECTION.
       01  WS-ｽﾃｰﾀｽ.
           05  WS-GLJRNL-ST            PIC XX.
               88  WS-GLJRNL-OK                VALUE "00" "05".
               88  WS-GLJRNL-ｵﾜﾘ               VALUE "10".
           05  WS-REPORT-ST            PIC XX.
               88  WS-REPORT-OK                VALUE "00".
       01  WS-ｺﾏﾝﾄﾞｷﾞｮｳ                PIC X(80).
       01  WS-ｴﾗｰﾌｧｲﾙ                  PIC X(8).
       01  WS-ｴﾗｰｽﾃｰﾀｽ                 PIC XX.
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ                PIC 9(8).
       01  WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ            PIC X(7).
       01  WS-ﾎｳｺｳ                     PIC X VALUE "W".
       01  WS-ﾍﾝｶﾝｹｯｶ                  PIC 9.
       01  WS-ｶﾓｸ-ﾃｲｷﾞ.
           05  FILLER PIC X(24) VALUE "1110預金".
           05  FILLER PIC X(24) VALUE "1210契約者貸付金".
           05  FILLER PIC X(24) VALUE "2110責任準備金".
           05  FILLER PIC X(24) VALUE "2310契約者配当準備金".
           05  FILLER PIC X(24) VALUE "2410未払手数料".
           05  FILLER PIC X(24) VALUE "4110保険料収入".
           05  FILLER PIC X(24) VALUE "4310貸付金利息".
           05  FILLER PIC X(24) VALUE "5110死亡保険金".
           05  FILLER PIC X(24) VALUE "5120入院給付金".
           05  FILLER PIC X(24) VALUE "5130解約返戻金".
           05  FILLER PIC X(24) VALUE "5140満期保険金".
           05  FILLER PIC X(24) VALUE "5210責任準備金繰入".
           05  FILLER PIC X(24) VALUE "5310契約者配当金".
           05  FILLER PIC X(24) VALUE "5410代理店手数料".
       01  WS-ｶﾓｸ-TBL REDEFINES WS-ｶﾓｸ-ﾃｲｷﾞ.
           05  WS-ｶﾓｸ OCCURS 14 TIMES INDEXED BY WS-ｶﾓｸ-IX.
               10  WS-ｶﾓｸ-ｺｰﾄﾞ             PIC 9(4).
               10  WS-ｶﾓｸ-ﾒｲ               PIC X(20).
       01  WS-ｼｭｳｹｲ-TBL.
           05  WS-ｼｭｳｹｲ OCCURS 14 TIMES.
               10  WS-ｶﾘｶﾀ                 PIC 9(13).
               10  WS-ｶｼｶﾀ                 PIC 9(13).
       01  WS-N                        PIC 99.
       01  WS-ｷﾁｮｳ-ｶﾓｸ                 PIC 9(4).
       01  WS-ｷﾁｮｳ-ｶﾞﾜ                 PIC X.
           88  WS-ｷﾁｮｳ-ｶﾘｶﾀ                    VALUE "D".
       01  WS-ｹﾝｽｳ.
           05  WS-ｼﾜｹｹﾝｽｳ              PIC 9(7) VALUE 0.
           05  WS-ﾌｾｲｹﾝｽｳ              PIC 9(7) VALUE 0.
           05  WS-ｶﾘｶﾀｺﾞｳｹｲ            PIC 9(13) VALUE 0.
           05  WS-ｶｼｶﾀｺﾞｳｹｲ            PIC 9(13) VALUE 0.
       01  WS-ﾌｯｷｺｰﾄﾞ                  PIC 99 VALUE 0.
       01  WS-ﾎｳｺｸ.
           05  WS-ﾋｮｳﾀﾞｲ-ｷﾞｮｳ.
               10  FILLER                  PIC X(22)
                   VALUE "ゴンドワナ生命  日計表".
               10  FILLER                  PIC X(10) VALUE SPACES.
               10  WS-ﾋｮｳﾀﾞｲ-ﾋﾂﾞｹ          PIC X(7).
           05  WS-ﾐﾀﾞｼ-ｷﾞｮｳ.
               10  FILLER                  PIC X(26) VALUE "科目".
               10  FILLER                  PIC X(17) VALUE "借方".
               10  FILLER                  PIC X(4) VALUE "貸方".
           05  WS-ﾒｲｻｲ-ｷﾞｮｳ.
               10  WS-ﾒｲｻｲ-ｺｰﾄﾞ            PIC 9(4).
               10  FILLER                  PIC X VALUE SPACE.
               10  WS-ﾒｲｻｲ-ﾒｲ              PIC X(20).
               10  WS-ﾒｲｻｲ-ｶﾘｶﾀ            PIC ZZ,ZZZ,ZZZ,ZZ9.
               10  FILLER                  PIC X(3) VALUE SPACES.
               10  WS-ﾒｲｻｲ-ｶｼｶﾀ            PIC ZZ,ZZZ,ZZZ,ZZ9.
           05  WS-ｺﾞｳｹｲ-ｷﾞｮｳ.
               10  FILLER                  PIC X(25) VALUE "合計".
               10  WS-ｺﾞｳｹｲ-ｶﾘｶﾀ           PIC ZZ,ZZZ,ZZZ,ZZ9.
               10  FILLER                  PIC X(3) VALUE SPACES.
               10  WS-ｺﾞｳｹｲ-ｶｼｶﾀ           PIC ZZ,ZZZ,ZZZ,ZZ9.
           05  WS-ｹﾝｽｳ-ｷﾞｮｳ.
               10  WS-ｹﾝｽｳ-ﾗﾍﾞﾙ            PIC X(25).
               10  WS-ｹﾝｽｳ-ｱﾀｲ             PIC ZZ,ZZZ,ZZZ,ZZ9.

       PROCEDURE DIVISION.
      *----------------------------------------------------------------
      * ﾒｲﾝ-RTN: 準備、仕訳の集計、日計表、終了。
      *----------------------------------------------------------------
       ﾒｲﾝ-RTN.
           PERFORM ｼｮｷ-RTN
           PERFORM ｼﾜｹ-RTN UNTIL NOT WS-GLJRNL-OK
           PERFORM ﾎｳｺｸ-RTN
           CLOSE GLJRNL-F REPORT-F
           DISPLAY "GLD080 仕訳 " WS-ｼﾜｹｹﾝｽｳ " 不正 " WS-ﾌｾｲｹﾝｽｳ
           MOVE WS-ﾌｯｷｺｰﾄﾞ TO RETURN-CODE
           STOP RUN.

      *----------------------------------------------------------------
      * ｼｮｷ-RTN: 業務日付、ファイル、最初の仕訳。
      *----------------------------------------------------------------
       ｼｮｷ-RTN.
           ACCEPT WS-ｺﾏﾝﾄﾞｷﾞｮｳ FROM COMMAND-LINE
           MOVE WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8) TO WS-ｷﾞｮｳﾑﾋﾂﾞｹ
           CALL "GLC020" USING WS-ﾎｳｺｳ WS-ｷﾞｮｳﾑﾋﾂﾞｹ
                               WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ WS-ﾍﾝｶﾝｹｯｶ
           IF WS-ﾍﾝｶﾝｹｯｶ NOT = 0
               DISPLAY "GLD080 日付不正: " WS-ｺﾏﾝﾄﾞｷﾞｮｳ(1:8)
               MOVE 8 TO RETURN-CODE
               STOP RUN
           END-IF
           MOVE ZEROS TO WS-ｼｭｳｹｲ-TBL
           OPEN INPUT GLJRNL-F
           OPEN OUTPUT REPORT-F
           EVALUATE TRUE
               WHEN NOT WS-GLJRNL-OK
                   MOVE "GLJRNL" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-GLJRNL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
                   PERFORM ｲｼﾞｮｳ-RTN
               WHEN NOT WS-REPORT-OK
                   MOVE "GLD080R" TO WS-ｴﾗｰﾌｧｲﾙ
                   MOVE WS-REPORT-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
                   PERFORM ｲｼﾞｮｳ-RTN
           END-EVALUATE
           PERFORM GLJRNL-ﾖﾐｺﾐ-RTN.

      *----------------------------------------------------------------
      * GLJRNL-ﾖﾐｺﾐ-RTN: 次の仕訳。
      *----------------------------------------------------------------
       GLJRNL-ﾖﾐｺﾐ-RTN.
           READ GLJRNL-F
           IF NOT WS-GLJRNL-OK AND NOT WS-GLJRNL-ｵﾜﾘ
               MOVE "GLJRNL" TO WS-ｴﾗｰﾌｧｲﾙ
               MOVE WS-GLJRNL-ST TO WS-ｴﾗｰｽﾃｰﾀｽ
               PERFORM ｲｼﾞｮｳ-RTN
           END-IF.

      *----------------------------------------------------------------
      * ｼﾜｹ-RTN: 仕訳1行を、勘定科目は貸借の側に、相手科目は反対の
      * 側に記帳する。相手科目がなければ数えて記帳しない。
      *----------------------------------------------------------------
       ｼﾜｹ-RTN.
           ADD 1 TO WS-ｼﾜｹｹﾝｽｳ
           IF JR-ｱｲﾃｶﾓｸ NOT NUMERIC OR JR-ｱｲﾃｶﾓｸ = 0
               ADD 1 TO WS-ﾌｾｲｹﾝｽｳ
               MOVE 4 TO WS-ﾌｯｷｺｰﾄﾞ
           ELSE
               MOVE JR-ｶﾝｼﾞｮｳｶﾓｸ TO WS-ｷﾁｮｳ-ｶﾓｸ
               MOVE JR-ﾀｲｼｬｸ TO WS-ｷﾁｮｳ-ｶﾞﾜ
               PERFORM ｷﾁｮｳ-RTN
               MOVE JR-ｱｲﾃｶﾓｸ TO WS-ｷﾁｮｳ-ｶﾓｸ
               IF JR-ｶﾘｶﾀ
                   MOVE "C" TO WS-ｷﾁｮｳ-ｶﾞﾜ
               ELSE
                   MOVE "D" TO WS-ｷﾁｮｳ-ｶﾞﾜ
               END-IF
               PERFORM ｷﾁｮｳ-RTN
           END-IF
           PERFORM GLJRNL-ﾖﾐｺﾐ-RTN.

      *----------------------------------------------------------------
      * ｷﾁｮｳ-RTN: 科目表で科目を引き、その側に金額を加える。科目表
      * にない科目は数える。
      *----------------------------------------------------------------
       ｷﾁｮｳ-RTN.
           SET WS-ｶﾓｸ-IX TO 1
           SEARCH WS-ｶﾓｸ
               AT END
                   ADD 1 TO WS-ﾌｾｲｹﾝｽｳ
                   MOVE 4 TO WS-ﾌｯｷｺｰﾄﾞ
               WHEN WS-ｶﾓｸ-ｺｰﾄﾞ(WS-ｶﾓｸ-IX) = WS-ｷﾁｮｳ-ｶﾓｸ
                   SET WS-N TO WS-ｶﾓｸ-IX
                   IF WS-ｷﾁｮｳ-ｶﾘｶﾀ
                       ADD JR-ｷﾝｶﾞｸ TO WS-ｶﾘｶﾀ(WS-N) WS-ｶﾘｶﾀｺﾞｳｹｲ
                   ELSE
                       ADD JR-ｷﾝｶﾞｸ TO WS-ｶｼｶﾀ(WS-N) WS-ｶｼｶﾀｺﾞｳｹｲ
                   END-IF
           END-SEARCH.

      *----------------------------------------------------------------
      * ﾎｳｺｸ-RTN: 動きのあった科目を科目順に、次いで合計と件数。
      *----------------------------------------------------------------
       ﾎｳｺｸ-RTN.
           MOVE WS-ｷﾞｮｳﾑﾋﾂﾞｹ-ﾜﾚｷ TO WS-ﾋｮｳﾀﾞｲ-ﾋﾂﾞｹ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ﾋｮｳﾀﾞｲ-ｷﾞｮｳ
           MOVE SPACES TO ﾎｳｺｸ-ｷﾞｮｳ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ﾐﾀﾞｼ-ｷﾞｮｳ
           PERFORM VARYING WS-N FROM 1 BY 1 UNTIL WS-N > 14
               IF WS-ｶﾘｶﾀ(WS-N) > 0 OR WS-ｶｼｶﾀ(WS-N) > 0
                   MOVE WS-ｶﾓｸ-ｺｰﾄﾞ(WS-N) TO WS-ﾒｲｻｲ-ｺｰﾄﾞ
                   MOVE WS-ｶﾓｸ-ﾒｲ(WS-N) TO WS-ﾒｲｻｲ-ﾒｲ
                   MOVE WS-ｶﾘｶﾀ(WS-N) TO WS-ﾒｲｻｲ-ｶﾘｶﾀ
                   MOVE WS-ｶｼｶﾀ(WS-N) TO WS-ﾒｲｻｲ-ｶｼｶﾀ
                   WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ﾒｲｻｲ-ｷﾞｮｳ
               END-IF
           END-PERFORM
           MOVE WS-ｶﾘｶﾀｺﾞｳｹｲ TO WS-ｺﾞｳｹｲ-ｶﾘｶﾀ
           MOVE WS-ｶｼｶﾀｺﾞｳｹｲ TO WS-ｺﾞｳｹｲ-ｶｼｶﾀ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｺﾞｳｹｲ-ｷﾞｮｳ
           MOVE SPACES TO ﾎｳｺｸ-ｷﾞｮｳ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ
           MOVE "仕訳件数" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ｼﾜｹｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ
           MOVE "相手科目なし・未登録科目" TO WS-ｹﾝｽｳ-ﾗﾍﾞﾙ
           MOVE WS-ﾌｾｲｹﾝｽｳ TO WS-ｹﾝｽｳ-ｱﾀｲ
           WRITE ﾎｳｺｸ-ｷﾞｮｳ FROM WS-ｹﾝｽｳ-ｷﾞｮｳ.

      *----------------------------------------------------------------
      * ｲｼﾞｮｳ-RTN: ファイル異常は 16 で終わる。
      *----------------------------------------------------------------
       ｲｼﾞｮｳ-RTN.
           DISPLAY "GLD080 ファイル異常 " WS-ｴﾗｰﾌｧｲﾙ
                   " 状態 " WS-ｴﾗｰｽﾃｰﾀｽ
           MOVE 16 TO RETURN-CODE
           STOP RUN.
