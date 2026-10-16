      *================================================================*
      * HOLDEXT: EXTRACT THE DAY'S HOLDS FOR POSTING                   *
      *                                                                *
      * AT THE END OF THE ONLINE DAY, READS THE HOLDS IN KEY ORDER,    *
      * WHICH IS ACCOUNT ORDER, AND WRITES THE CARD FILE FOR THE       *
      * NIGHT'S BATCH IN THE FORM OF THE DAILY TRANSACTION FILE: A     *
      * HEADER FROM THE CARDS FOR THE BUSINESS DATE, A DETAIL FOR      *
      * EVERY CARD PURCHASE AND ATM WITHDRAWAL STILL HELD, ANOTHER     *
      * FOR THE FEE OF A WITHDRAWAL THAT HAD ONE, AND A TRAILER WITH   *
      * THE CONTROL TOTALS. REVERSED PURCHASES ARE LEFT OUT.           *
      * TRANSACTION IDS ARE 9, THE HOLD NUMBER, 00, AND 1 FOR THE      *
      * AMOUNT OR 2 FOR THE FEE.                                       *
      *                                                                *
      * RETURN CODES: 0 EXTRACTED, 16 A HOLD OF ANOTHER DAY OR I/O.    *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 2003-05-19 MV  CR-0342 ORIGINAL VERSION.                       *
      * 2007-02-05 MV  CR-0415 ATM FEES.                               *
      *================================================================*
       IDENTIFICATION DIVISION.
       PROGRAM-ID. HOLDEXT.
       AUTHOR. M. VARGA, ONLINE SYSTEMS.
       INSTALLATION. FIRST BANK OF PANGAEA, DATA CENTER.
       DATE-WRITTEN. 2003-05-19.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT CONTROL-FILE ASSIGN TO BANKCTL
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-CONTROL-STATUS.
           SELECT HOLD-FILE    ASSIGN TO HOLDS
               ORGANIZATION IS INDEXED
               ACCESS MODE IS SEQUENTIAL
               RECORD KEY IS HD-KEY
               FILE STATUS IS WS-HOLD-STATUS.
           SELECT TRAN-OUT     ASSIGN TO TRANOUT
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-TRAN-OUT-STATUS.
           SELECT REPORT-OUT   ASSIGN TO HOLDRPT
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-STATUS.

       DATA DIVISION.
       FILE SECTION.
       FD  CONTROL-FILE.
           COPY BANKCTL.
       FD  HOLD-FILE.
           COPY HOLDREC.
       FD  TRAN-OUT.
       01  TRAN-OUT-RECORD             PIC X(80).
       FD  REPORT-OUT.
       01  REPORT-LINE                 PIC X(80).

       WORKING-STORAGE SECTION.
       01  WS-FILE-STATUSES.
           05  WS-CONTROL-STATUS       PIC XX.
           05  WS-HOLD-STATUS          PIC XX.
           05  WS-TRAN-OUT-STATUS      PIC XX.
           05  WS-REPORT-STATUS        PIC XX.

       01  WS-ABORT-MESSAGE            PIC X(40).
       01  WS-HOLDS-END-SW             PIC X          VALUE 'N'.
           88  END-OF-HOLDS            VALUE 'Y'.

           COPY TRANREC.

       01  WS-TRAN-ID.
           05  FILLER                  PIC 9          VALUE 9.
           05  WS-TI-HOLD-NO           PIC 9(6).
           05  FILLER                  PIC 99         VALUE 0.
           05  WS-TI-PART              PIC 9.

       01  WS-TOTALS.
           05  WS-HOLD-COUNT           PIC 9(8)       VALUE 0.
           05  WS-REVERSED-COUNT       PIC 9(8)       VALUE 0.
           05  WS-DETAIL-COUNT         PIC 9(8)       VALUE 0.
           05  WS-AMOUNT-TOTAL         PIC 9(13)V99   VALUE 0.
           05  WS-ACCT-HASH            PIC 9(15)      VALUE 0.

       01  WS-REPORT-COUNT-LINE.
           05  WS-RC-LABEL             PIC X(30).
           05  WS-RC-VALUE             PIC ZZ,ZZZ,ZZ9.
           05  FILLER                  PIC X(40)      VALUE SPACES.
       01  WS-REPORT-AMOUNT-LINE.
           05  WS-RA-LABEL             PIC X(30).
           05  WS-RA-VALUE             PIC Z,ZZZ,ZZZ,ZZZ,ZZ9.99.
           05  FILLER                  PIC X(30)      VALUE SPACES.

       PROCEDURE DIVISION.
      * WRITE THE HEADER, A DETAIL OR TWO PER HOLD, AND THE TRAILER.
       0000-MAIN.
           PERFORM 1000-INITIALIZE
           PERFORM 2000-EXTRACT-HOLD UNTIL END-OF-HOLDS
           PERFORM 3000-WRITE-TRAILER
           PERFORM 8000-WRITE-REPORT
           PERFORM 9000-CLOSE-FILES
           MOVE 0 TO RETURN-CODE
           STOP RUN.

      * OPEN THE FILES, WRITE THE HEADER AND READ THE FIRST HOLD.
       1000-INITIALIZE.
           OPEN INPUT  CONTROL-FILE HOLD-FILE
                OUTPUT TRAN-OUT REPORT-OUT
           IF WS-FILE-STATUSES NOT = ALL '0'
               DISPLAY 'HOLDEXT: CANNOT OPEN FILES ' WS-FILE-STATUSES
               MOVE 16 TO RETURN-CODE
               STOP RUN
           END-IF
           READ CONTROL-FILE
               AT END
                   MOVE 'NO BANK CONTROL RECORD' TO WS-ABORT-MESSAGE
                   PERFORM 9900-ABORT
           END-READ
           MOVE SPACES TO TRAN-RECORD
           SET TR-HEADER TO TRUE
           MOVE BC-BUSINESS-DATE TO TH-BUSINESS-DATE
           SET TH-FROM-CARDS TO TRUE
           MOVE 1 TO TH-FILE-SEQ
           WRITE TRAN-OUT-RECORD FROM TRAN-RECORD
           PERFORM 1200-READ-HOLD.

       1200-READ-HOLD.
           READ HOLD-FILE
               AT END
                   SET END-OF-HOLDS TO TRUE
               NOT AT END
                   ADD 1 TO WS-HOLD-COUNT
           END-READ.

      * WRITE A HOLD STILL OPEN AS A TRANSACTION, AND ITS FEE AS
      * ANOTHER. EVERY HOLD MUST BE OF THE BUSINESS DAY.
       2000-EXTRACT-HOLD.
           IF HD-BUSINESS-DATE NOT = BC-BUSINESS-DATE
               MOVE 'A HOLD OF ANOTHER BUSINESS DAY' TO WS-ABORT-MESSAGE
               PERFORM 9900-ABORT
           END-IF
           IF HD-REVERSED
               ADD 1 TO WS-REVERSED-COUNT
           ELSE
               MOVE SPACES TO TRAN-RECORD
               SET TR-DETAIL TO TRUE
               MOVE HD-HOLD-NO TO WS-TI-HOLD-NO
               MOVE HD-ACCT-NO TO TD-ACCT-NO
               MOVE HD-TERMINAL TO TD-DESCRIPTION
               MOVE 0 TO TD-COUNTER-ACCT-NO
               IF HD-PURCHASE
                   SET TD-CARD-PURCHASE TO TRUE
               ELSE
                   SET TD-ATM-WITHDRAWAL TO TRUE
               END-IF
               MOVE HD-AMOUNT TO TD-AMOUNT
               MOVE 1 TO WS-TI-PART
               PERFORM 2900-WRITE-DETAIL
               IF HD-FEE > 0
                   SET TD-ATM-FEE TO TRUE
                   MOVE HD-FEE TO TD-AMOUNT
                   MOVE 'ATM FEE' TO TD-DESCRIPTION
                   MOVE 2 TO WS-TI-PART
                   PERFORM 2900-WRITE-DETAIL
               END-IF
           END-IF
           PERFORM 1200-READ-HOLD.

      * WRITE A DETAIL AND ADD IT TO THE CONTROL TOTALS.
       2900-WRITE-DETAIL.
           MOVE WS-TRAN-ID TO TD-TRAN-ID
           WRITE TRAN-OUT-RECORD FROM TRAN-RECORD
           ADD 1 TO WS-DETAIL-COUNT
           ADD TD-AMOUNT TO WS-AMOUNT-TOTAL
           ADD TD-ACCT-NO TO WS-ACCT-HASH.

       3000-WRITE-TRAILER.
           MOVE SPACES TO TRAN-RECORD
           SET TR-TRAILER TO TRUE
           MOVE WS-DETAIL-COUNT TO TT-DETAIL-COUNT
           MOVE WS-AMOUNT-TOTAL TO TT-AMOUNT-TOTAL
           MOVE WS-ACCT-HASH TO TT-ACCT-HASH
           WRITE TRAN-OUT-RECORD FROM TRAN-RECORD.

      * WRITE THE EXTRACTION REPORT.
       8000-WRITE-REPORT.
           MOVE SPACES TO REPORT-LINE
           STRING 'HOLDEXT  HOLDS FOR POSTING       BUSINESS DATE '
                  BC-BUSINESS-DATE
               DELIMITED BY SIZE INTO REPORT-LINE
           WRITE REPORT-LINE
           MOVE SPACES TO REPORT-LINE
           WRITE REPORT-LINE
           MOVE 'HOLDS READ' TO WS-RC-LABEL
           MOVE WS-HOLD-COUNT TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           MOVE 'REVERSED, LEFT OUT' TO WS-RC-LABEL
           MOVE WS-REVERSED-COUNT TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           MOVE 'TRANSACTIONS WRITTEN' TO WS-RC-LABEL
           MOVE WS-DETAIL-COUNT TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           MOVE 'AMOUNT TOTAL' TO WS-RA-LABEL
           MOVE WS-AMOUNT-TOTAL TO WS-RA-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-AMOUNT-LINE.

      * CLOSE ALL FILES.
       9000-CLOSE-FILES.
           CLOSE CONTROL-FILE HOLD-FILE TRAN-OUT REPORT-OUT.

      * STOP THE RUN WITH RETURN CODE 16, SAYING WHY.
       9900-ABORT.
           DISPLAY 'HOLDEXT: ' WS-ABORT-MESSAGE
           PERFORM 9000-CLOSE-FILES
           MOVE 16 TO RETURN-CODE
           STOP RUN.
