      *================================================================*
      * STMTPRT: PRINT THE MONTH'S STATEMENTS                          *
      *                                                                *
      * ON THE FIRST BUSINESS DAY OF A MONTH, AFTER THE DAILY BATCH,   *
      * PRINTS A STATEMENT FOR EVERY ACCOUNT FOR THE MONTH THAT        *
      * CLOSED: THE OPENING BALANCE, EVERY POSTING OF THE MONTH, THE   *
      * INTEREST AND FEES POSTED FOR IT AFTER IT CLOSED, AND THE       *
      * CLOSING BALANCE. A CLOSED ACCOUNT WITH NOTHING TO SHOW GETS    *
      * NONE.                                                          *
      *                                                                *
      * THE MONTH THAT CLOSED IS THE DUE MONTH OF THE DAY'S ACCRUALS.  *
      * THE POSTINGS ARE READ FROM EVERY AUDIT TRAIL, NEWEST FIRST,    *
      * PICKED BY DATE AND SORTED BY ACCOUNT, DATE AND THE ORDER THEY  *
      * WERE READ IN, WHICH WITHIN A DAY IS THE ORDER THEY POSTED IN.  *
      * THEY ARE MATCHED AGAINST THE MASTER THE MONTH CLOSED WITH.     *
      *                                                                *
      * RETURN CODES: 0 PRINTED OR NO MONTH CLOSED,                    *
      *               16 POSTINGS FOR NO ACCOUNT, SORT FAILED OR I/O.  *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 2001-09-03 TL  CR-0298 ORIGINAL VERSION.                       *
      * 2004-03-15 TL  CR-0390 THE CODE TABLE FROM CODETAB, SHARED     *
      *                WITH GLSUM.                                     *
      *================================================================*
       IDENTIFICATION DIVISION.
       PROGRAM-ID. STMTPRT.
       AUTHOR. T. LINDQVIST, BATCH SYSTEMS.
       INSTALLATION. FIRST BANK OF PANGAEA, DATA CENTER.
       DATE-WRITTEN. 2001-09-03.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT ACCR-IN      ASSIGN TO ACCRIN
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-ACCR-IN-STATUS.
           SELECT ACCT-IN      ASSIGN TO ACCTIN
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-ACCT-IN-STATUS.
           SELECT AUDIT-IN     ASSIGN TO AUDIT
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-AUDIT-STATUS.
           SELECT STMT-OUT     ASSIGN TO STMTOUT
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-STMT-STATUS.
           SELECT REPORT-OUT   ASSIGN TO STMTRPT
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-STATUS.
           SELECT STMT-SORT    ASSIGN TO SORTWK.

       DATA DIVISION.
       FILE SECTION.
       FD  ACCR-IN.
           COPY ACCRREC.
       FD  ACCT-IN.
           COPY ACCTMAST.
       FD  AUDIT-IN.
           COPY AUDITREC.
       FD  STMT-OUT.
       01  STMT-LINE                   PIC X(80).
       FD  REPORT-OUT.
       01  REPORT-LINE                 PIC X(80).
       SD  STMT-SORT.
       01  SORT-RECORD.
           05  SR-ACCT-NO              PIC 9(10).
           05  SR-DATE                 PIC 9(8).
           05  SR-SEQ                  PIC 9(8).
           05  SR-TRAN-ID              PIC 9(10).
           05  SR-TRAN-CODE            PIC XX.
           05  SR-AMOUNT               PIC S9(9)V99.
           05  SR-BALANCE-BEFORE       PIC S9(11)V99.
           05  SR-BALANCE-AFTER        PIC S9(11)V99.

       WORKING-STORAGE SECTION.
       01  WS-FILE-STATUSES.
           05  WS-ACCR-IN-STATUS       PIC XX.
           05  WS-ACCT-IN-STATUS       PIC XX.
           05  WS-AUDIT-STATUS         PIC XX.
           05  WS-STMT-STATUS          PIC XX.
           05  WS-REPORT-STATUS        PIC XX.

       01  WS-ABORT-MESSAGE            PIC X(40).
       01  WS-SWITCHES.
           05  WS-MASTER-END-SW        PIC X          VALUE 'N'.
               88  MASTER-AT-END       VALUE 'Y'.
           05  WS-POSTING-END-SW       PIC X          VALUE 'N'.
               88  POSTINGS-AT-END     VALUE 'Y'.

      * THE MONTH THAT CLOSED AND THE ONE AFTER IT, WHOSE FIRST BUSINESS
      * DAY POSTED THE CLOSED MONTH'S INTEREST AND FEES.
       01  WS-CLOSED-MONTH.
           05  WS-CM-YEAR              PIC 9(4).
           05  WS-CM-MONTH             PIC 99.
       01  WS-CLOSED-MONTH-NUMBER REDEFINES WS-CLOSED-MONTH
                                       PIC 9(6).
       01  WS-NEXT-MONTH.
           05  WS-NM-YEAR              PIC 9(4).
           05  WS-NM-MONTH             PIC 99.
       01  WS-POSTING-MONTH            PIC X(6).
       01  WS-POSTING-CODE             PIC XX.
           88  WS-MONTH-END-CODE       VALUE 'IN' 'OD' 'MF'.

      * THE KEY OF THE POSTING IN HAND: HIGH-VALUES AFTER THE LAST.
       01  WS-POSTING-KEY              PIC X(10).
       01  WS-SEQ                      PIC 9(8)       VALUE 0.

       01  WS-MONTH-NAME-VALUES.
           05  FILLER                  PIC X(9)       VALUE 'JANUARY'.
           05  FILLER                  PIC X(9)       VALUE 'FEBRUARY'.
           05  FILLER                  PIC X(9)       VALUE 'MARCH'.
           05  FILLER                  PIC X(9)       VALUE 'APRIL'.
           05  FILLER                  PIC X(9)       VALUE 'MAY'.
           05  FILLER                  PIC X(9)       VALUE 'JUNE'.
           05  FILLER                  PIC X(9)       VALUE 'JULY'.
           05  FILLER                  PIC X(9)       VALUE 'AUGUST'.
           05  FILLER                  PIC X(9)       VALUE 'SEPTEMBER'.
           05  FILLER                  PIC X(9)       VALUE 'OCTOBER'.
           05  FILLER                  PIC X(9)       VALUE 'NOVEMBER'.
           05  FILLER                  PIC X(9)       VALUE 'DECEMBER'.
       01  WS-MONTH-NAMES REDEFINES WS-MONTH-NAME-VALUES.
           05  WS-MONTH-NAME           PIC X(9)       OCCURS 12 TIMES.

      * WHAT EACH TRANSACTION CODE SAYS ON A STATEMENT.
           COPY CODETAB REPLACING ==:P:== BY ==ST==.

       01  WS-STATEMENT.
           05  WS-OPENING              PIC S9(11)V99.
           05  WS-CLOSING              PIC S9(11)V99.
           05  WS-CREDITS              PIC S9(11)V99.
           05  WS-DEBITS               PIC S9(11)V99.

       01  WS-TOTALS.
           05  WS-ACCOUNT-COUNT        PIC 9(8)       VALUE 0.
           05  WS-STATEMENT-COUNT      PIC 9(8)       VALUE 0.
           05  WS-SELECTED-COUNT       PIC 9(8)       VALUE 0.
           05  WS-PRINTED-COUNT        PIC 9(8)       VALUE 0.
           05  WS-CREDIT-TOTAL         PIC S9(13)V99  VALUE 0.
           05  WS-DEBIT-TOTAL          PIC S9(13)V99  VALUE 0.

       01  WS-RULE-LINE                PIC X(80)      VALUE ALL '='.
       01  WS-TITLE-LINE.
           05  FILLER                  PIC X(48)
                                       VALUE 'FIRST BANK OF PANGAEA'.
           05  WS-TL-PERIOD            PIC X(32).
       01  WS-ACCOUNT-LINE.
           05  FILLER                  PIC X(8)       VALUE 'ACCOUNT'.
           05  WS-AL-ACCT-NO           PIC 9(10).
           05  FILLER                  PIC X(3)       VALUE SPACES.
           05  WS-AL-NAME              PIC X(20).
           05  FILLER                  PIC X(16)      VALUE SPACES.
           05  WS-AL-KIND              PIC X(23).
       01  WS-HEADING-LINE.
           05  FILLER                  PIC X(12)      VALUE 'DATE'.
           05  FILLER                  PIC X(12)
                                       VALUE 'TRANSACTION'.
           05  FILLER                  PIC X(20)
                                       VALUE 'DESCRIPTION'.
           05  FILLER                  PIC X(16)
                                       VALUE '          AMOUNT'.
           05  FILLER                  PIC X(19)
                                       VALUE '            BALANCE'.
           05  FILLER                  PIC X(1)       VALUE SPACES.
       01  WS-DATE                     PIC 9(8).
       01  WS-DATE-PARTS REDEFINES WS-DATE.
           05  WS-DATE-YEAR            PIC 9(4).
           05  WS-DATE-MONTH           PIC 99.
           05  WS-DATE-DAY             PIC 99.
       01  WS-POSTING-LINE.
           05  WS-PL-DATE.
               10  WS-PL-YEAR          PIC 9(4).
               10  FILLER              PIC X          VALUE '-'.
               10  WS-PL-MONTH         PIC 99.
               10  FILLER              PIC X          VALUE '-'.
               10  WS-PL-DAY           PIC 99.
           05  FILLER                  PIC X(2)       VALUE SPACES.
           05  WS-PL-TRAN-ID           PIC 9(10).
           05  FILLER                  PIC X(2)       VALUE SPACES.
           05  WS-PL-DESCRIPTION       PIC X(20).
           05  WS-PL-AMOUNT            PIC Z,ZZZ,ZZZ,ZZ9.99-.
           05  FILLER                  PIC X(1)       VALUE SPACES.
           05  WS-PL-BALANCE           PIC ZZ,ZZZ,ZZZ,ZZ9.99-.
       01  WS-BALANCE-LINE.
           05  FILLER                  PIC X(12)      VALUE SPACES.
           05  WS-BL-LABEL             PIC X(50).
           05  WS-BL-BALANCE           PIC ZZ,ZZZ,ZZZ,ZZ9.99-.
       01  WS-SUM-LINE.
           05  FILLER                  PIC X(12)      VALUE SPACES.
           05  FILLER                  PIC X(8)       VALUE 'CREDITS'.
           05  WS-SL-CREDITS           PIC ZZ,ZZZ,ZZZ,ZZ9.99.
           05  FILLER                  PIC X(4)       VALUE SPACES.
           05  FILLER                  PIC X(7)       VALUE 'DEBITS'.
           05  WS-SL-DEBITS            PIC ZZ,ZZZ,ZZZ,ZZ9.99.
           05  FILLER                  PIC X(15)      VALUE SPACES.

       01  WS-REPORT-COUNT-LINE.
           05  WS-RC-LABEL             PIC X(30).
           05  WS-RC-VALUE             PIC ZZ,ZZZ,ZZ9.
           05  FILLER                  PIC X(40)      VALUE SPACES.
       01  WS-REPORT-AMOUNT-LINE.
           05  WS-RA-LABEL             PIC X(30).
           05  WS-RA-VALUE             PIC Z,ZZZ,ZZZ,ZZZ,ZZ9.99-.
           05  FILLER                  PIC X(29)      VALUE SPACES.

       PROCEDURE DIVISION.
      * FIND THE MONTH THAT CLOSED, IF ONE DID, AND PRINT ITS
      * STATEMENTS FROM THE SORTED POSTINGS.
       0000-MAIN.
           PERFORM 1000-INITIALIZE
           IF WS-CLOSED-MONTH-NUMBER NOT = 0
               SORT STMT-SORT
                   ON ASCENDING KEY SR-ACCT-NO SR-DATE SR-SEQ
                   INPUT PROCEDURE IS 2000-SELECT-POSTINGS
                   OUTPUT PROCEDURE IS 3000-PRINT-STATEMENTS
               IF SORT-RETURN NOT = 0
                   MOVE 'SORT FAILED' TO WS-ABORT-MESSAGE
                   PERFORM 9900-ABORT
               END-IF
           END-IF
           PERFORM 8000-WRITE-REPORT
           PERFORM 9000-CLOSE-FILES
           MOVE 0 TO RETURN-CODE
           STOP RUN.

      * OPEN THE FILES AND TAKE THE CLOSED MONTH FROM THE ACCRUALS.
       1000-INITIALIZE.
           OPEN INPUT  ACCR-IN ACCT-IN AUDIT-IN
                OUTPUT STMT-OUT REPORT-OUT
           IF WS-FILE-STATUSES NOT = ALL '0'
               DISPLAY 'STMTPRT: CANNOT OPEN FILES ' WS-FILE-STATUSES
               MOVE 16 TO RETURN-CODE
               STOP RUN
           END-IF
           MOVE 0 TO WS-CLOSED-MONTH-NUMBER
           READ ACCR-IN
               NOT AT END
                   MOVE AC-DUE-MONTH TO WS-CLOSED-MONTH-NUMBER
           END-READ
           MOVE WS-CM-YEAR TO WS-NM-YEAR
           IF WS-CM-MONTH = 12
               ADD 1 TO WS-NM-YEAR
               MOVE 1 TO WS-NM-MONTH
           ELSE
               COMPUTE WS-NM-MONTH = WS-CM-MONTH + 1
           END-IF.

      * RELEASE THE POSTINGS OF THE CLOSED MONTH, AND THE INTEREST AND
      * FEES POSTED FOR IT ON THE FIRST BUSINESS DAY OF THE NEXT.
       2000-SELECT-POSTINGS.
           READ AUDIT-IN
           PERFORM UNTIL WS-AUDIT-STATUS NOT = '00'
               ADD 1 TO WS-SEQ
               MOVE AU-BUSINESS-DATE(1:6) TO WS-POSTING-MONTH
               MOVE AU-TRAN-CODE TO WS-POSTING-CODE
               IF WS-POSTING-MONTH = WS-CLOSED-MONTH
                  OR (WS-POSTING-MONTH = WS-NEXT-MONTH
                      AND WS-MONTH-END-CODE)
                   PERFORM 2100-RELEASE-POSTING
               END-IF
               READ AUDIT-IN
           END-PERFORM
           IF WS-AUDIT-STATUS NOT = '10'
               MOVE 'CANNOT READ THE AUDIT TRAILS' TO WS-ABORT-MESSAGE
               PERFORM 9900-ABORT
           END-IF.

      * RELEASE ONE POSTING TO THE SORT.
       2100-RELEASE-POSTING.
           MOVE AU-ACCT-NO        TO SR-ACCT-NO
           MOVE AU-BUSINESS-DATE  TO SR-DATE
           MOVE WS-SEQ            TO SR-SEQ
           MOVE AU-TRAN-ID        TO SR-TRAN-ID
           MOVE AU-TRAN-CODE      TO SR-TRAN-CODE
           MOVE AU-AMOUNT         TO SR-AMOUNT
           MOVE AU-BALANCE-BEFORE TO SR-BALANCE-BEFORE
           MOVE AU-BALANCE-AFTER  TO SR-BALANCE-AFTER
           RELEASE SORT-RECORD
           ADD 1 TO WS-SELECTED-COUNT.

      * PRINT A STATEMENT FOR EACH ACCOUNT OF THE MASTER, WITH ITS
      * SORTED POSTINGS. EVERY POSTING MUST BELONG TO AN ACCOUNT.
       3000-PRINT-STATEMENTS.
           PERFORM 3100-RETURN-POSTING
           PERFORM 3200-READ-MASTER
           PERFORM UNTIL MASTER-AT-END
               IF WS-POSTING-KEY < AM-ACCT-NO
                   MOVE 'POSTINGS FOR NO ACCOUNT' TO WS-ABORT-MESSAGE
                   PERFORM 9900-ABORT
               END-IF
               PERFORM 3300-PRINT-STATEMENT
               PERFORM 3200-READ-MASTER
           END-PERFORM
           IF NOT POSTINGS-AT-END
               MOVE 'POSTINGS FOR NO ACCOUNT' TO WS-ABORT-MESSAGE
               PERFORM 9900-ABORT
           END-IF.

      * RETURN THE NEXT SORTED POSTING.
       3100-RETURN-POSTING.
           RETURN STMT-SORT
               AT END
                   SET POSTINGS-AT-END TO TRUE
                   MOVE HIGH-VALUES TO WS-POSTING-KEY
               NOT AT END
                   MOVE SR-ACCT-NO TO WS-POSTING-KEY
           END-RETURN.

      * READ THE NEXT ACCOUNT OF THE MASTER.
       3200-READ-MASTER.
           READ ACCT-IN
               AT END
                   SET MASTER-AT-END TO TRUE
               NOT AT END
                   ADD 1 TO WS-ACCOUNT-COUNT
           END-READ.

      * PRINT ONE ACCOUNT'S STATEMENT, UNLESS IT IS CLOSED AND HAS
      * NOTHING TO SHOW.
       3300-PRINT-STATEMENT.
           IF AM-CLOSED AND AM-BALANCE = 0
              AND WS-POSTING-KEY NOT = AM-ACCT-NO
               EXIT PARAGRAPH
           END-IF
           MOVE AM-BALANCE TO WS-OPENING
           IF WS-POSTING-KEY = AM-ACCT-NO
               MOVE SR-BALANCE-BEFORE TO WS-OPENING
           END-IF
           MOVE WS-OPENING TO WS-CLOSING
           MOVE 0 TO WS-CREDITS WS-DEBITS
           PERFORM 3400-PRINT-HEADING
           PERFORM UNTIL WS-POSTING-KEY NOT = AM-ACCT-NO
               PERFORM 3500-PRINT-POSTING
               PERFORM 3100-RETURN-POSTING
           END-PERFORM
           PERFORM 3600-PRINT-CLOSING
           ADD 1 TO WS-STATEMENT-COUNT.

      * PRINT THE STATEMENT'S HEADING AND OPENING BALANCE.
       3400-PRINT-HEADING.
           WRITE STMT-LINE FROM WS-RULE-LINE
           MOVE SPACES TO WS-TL-PERIOD
           STRING 'STATEMENT FOR ' DELIMITED BY SIZE
                  WS-MONTH-NAME(WS-CM-MONTH) DELIMITED BY SPACE
                  ' ' WS-CM-YEAR DELIMITED BY SIZE
               INTO WS-TL-PERIOD
           WRITE STMT-LINE FROM WS-TITLE-LINE
           MOVE AM-ACCT-NO TO WS-AL-ACCT-NO
           MOVE AM-SHORT-NAME TO WS-AL-NAME
           EVALUATE TRUE
               WHEN AM-SAVINGS AND AM-FROZEN
                   MOVE 'SAVINGS ACCOUNT, FROZEN' TO WS-AL-KIND
               WHEN AM-SAVINGS
                   MOVE 'SAVINGS ACCOUNT' TO WS-AL-KIND
               WHEN AM-FROZEN
                   MOVE 'CURRENT ACCOUNT, FROZEN' TO WS-AL-KIND
               WHEN OTHER
                   MOVE 'CURRENT ACCOUNT' TO WS-AL-KIND
           END-EVALUATE
           WRITE STMT-LINE FROM WS-ACCOUNT-LINE
           MOVE SPACES TO STMT-LINE
           WRITE STMT-LINE
           WRITE STMT-LINE FROM WS-HEADING-LINE
           MOVE 'OPENING BALANCE' TO WS-BL-LABEL
           MOVE WS-OPENING TO WS-BL-BALANCE
           WRITE STMT-LINE FROM WS-BALANCE-LINE.

      * PRINT ONE POSTING, DESCRIBED BY ITS TRANSACTION CODE.
       3500-PRINT-POSTING.
           MOVE SR-DATE TO WS-DATE
           MOVE WS-DATE-YEAR TO WS-PL-YEAR
           MOVE WS-DATE-MONTH TO WS-PL-MONTH
           MOVE WS-DATE-DAY TO WS-PL-DAY
           MOVE SR-TRAN-ID TO WS-PL-TRAN-ID
           SET ST-CX TO 1
           SEARCH ST-CODE-ENTRY
               AT END
                   MOVE SR-TRAN-CODE TO WS-PL-DESCRIPTION
               WHEN ST-CODE(ST-CX) = SR-TRAN-CODE
                   MOVE ST-CODE-TEXT(ST-CX) TO WS-PL-DESCRIPTION
           END-SEARCH
           MOVE SR-AMOUNT TO WS-PL-AMOUNT
           MOVE SR-BALANCE-AFTER TO WS-PL-BALANCE
           WRITE STMT-LINE FROM WS-POSTING-LINE
           IF SR-AMOUNT < 0
               SUBTRACT SR-AMOUNT FROM WS-DEBITS
           ELSE
               ADD SR-AMOUNT TO WS-CREDITS
           END-IF
           MOVE SR-BALANCE-AFTER TO WS-CLOSING
           ADD 1 TO WS-PRINTED-COUNT.

      * PRINT THE CLOSING BALANCE AND THE MONTH'S CREDITS AND DEBITS.
       3600-PRINT-CLOSING.
           MOVE 'CLOSING BALANCE' TO WS-BL-LABEL
           MOVE WS-CLOSING TO WS-BL-BALANCE
           WRITE STMT-LINE FROM WS-BALANCE-LINE
           MOVE WS-CREDITS TO WS-SL-CREDITS
           MOVE WS-DEBITS TO WS-SL-DEBITS
           WRITE STMT-LINE FROM WS-SUM-LINE
           MOVE SPACES TO STMT-LINE
           WRITE STMT-LINE
           ADD WS-CREDITS TO WS-CREDIT-TOTAL
           ADD WS-DEBITS TO WS-DEBIT-TOTAL.

      * WRITE THE STATEMENT RUN'S REPORT. EVERY POSTING SELECTED MUST
      * HAVE BEEN PRINTED.
       8000-WRITE-REPORT.
           MOVE 'STMTPRT  MONTHLY STATEMENTS' TO REPORT-LINE
           WRITE REPORT-LINE
           MOVE SPACES TO REPORT-LINE
           WRITE REPORT-LINE
           IF WS-CLOSED-MONTH-NUMBER = 0
               MOVE 'NO MONTH CLOSED' TO REPORT-LINE
               WRITE REPORT-LINE
               EXIT PARAGRAPH
           END-IF
           MOVE SPACES TO REPORT-LINE
           STRING 'MONTH CLOSED                  ' WS-CLOSED-MONTH
               DELIMITED BY SIZE INTO REPORT-LINE
           WRITE REPORT-LINE
           MOVE 'ACCOUNTS READ' TO WS-RC-LABEL
           MOVE WS-ACCOUNT-COUNT TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           MOVE 'STATEMENTS PRINTED' TO WS-RC-LABEL
           MOVE WS-STATEMENT-COUNT TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           MOVE 'POSTINGS SELECTED' TO WS-RC-LABEL
           MOVE WS-SELECTED-COUNT TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           MOVE 'POSTINGS PRINTED' TO WS-RC-LABEL
           MOVE WS-PRINTED-COUNT TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           MOVE 'CREDITS' TO WS-RA-LABEL
           MOVE WS-CREDIT-TOTAL TO WS-RA-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-AMOUNT-LINE
           MOVE 'DEBITS' TO WS-RA-LABEL
           MOVE WS-DEBIT-TOTAL TO WS-RA-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-AMOUNT-LINE
           IF WS-PRINTED-COUNT NOT = WS-SELECTED-COUNT
               MOVE 'POSTINGS SELECTED BUT NOT PRINTED'
                   TO WS-ABORT-MESSAGE
               PERFORM 9900-ABORT
           END-IF.

      * CLOSE ALL FILES.
       9000-CLOSE-FILES.
           CLOSE ACCR-IN ACCT-IN AUDIT-IN STMT-OUT REPORT-OUT.

      * STOP THE RUN WITH RETURN CODE 16, SAYING WHY.
       9900-ABORT.
           DISPLAY 'STMTPRT: ' WS-ABORT-MESSAGE
           PERFORM 9000-CLOSE-FILES
           MOVE 16 TO RETURN-CODE
           STOP RUN.
