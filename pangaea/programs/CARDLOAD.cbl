      *================================================================*
      * CARDLOAD: LOAD THE ONLINE CARD FILE                            *
      *                                                                *
      * WRITES THE CARD MASTER'S CARDS TO A NEW INDEXED ONLINE CARD    *
      * FILE FOR THE NEXT BUSINESS DAY, THE DATE OF THE NEW BANK       *
      * CONTROL RECORD. CARDS THAT HAVE EXPIRED BY THAT DAY ARE LEFT   *
      * OUT, SO THAT AUTHORIZATIONS FIND NO SUCH CARD; BLOCKED CARDS   *
      * ARE LOADED, FOR THE AUTHORIZATIONS TO DECLINE. THE MASTER IS   *
      * IN CARD NUMBER ORDER, SO THE FILE IS WRITTEN IN KEY ORDER,     *
      * SEQUENTIALLY.                                                  *
      *                                                                *
      * RETURN CODES: 0 LOADED, 16 MASTER OUT OF ORDER OR I/O.         *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 2003-05-19 MV  CR-0342 ORIGINAL VERSION.                       *
      *================================================================*
       IDENTIFICATION DIVISION.
       PROGRAM-ID. CARDLOAD.
       AUTHOR. M. VARGA, ONLINE SYSTEMS.
       INSTALLATION. FIRST BANK OF PANGAEA, DATA CENTER.
       DATE-WRITTEN. 2003-05-19.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT CONTROL-FILE ASSIGN TO BANKCTL
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-CONTROL-STATUS.
           SELECT CARD-IN      ASSIGN TO CARDIN
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-CARD-IN-STATUS.
           SELECT ONLINE-OUT   ASSIGN TO CARDOUT
               ORGANIZATION IS INDEXED
               ACCESS MODE IS SEQUENTIAL
               RECORD KEY IS OC-CARD-NO
               FILE STATUS IS WS-ONLINE-STATUS.
           SELECT REPORT-OUT   ASSIGN TO LOADRPT
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-STATUS.

       DATA DIVISION.
       FILE SECTION.
       FD  CONTROL-FILE.
           COPY BANKCTL.
       FD  CARD-IN.
           COPY CARDREC.
       FD  ONLINE-OUT.
       01  ONLINE-CARD.
           05  OC-CARD-NO              PIC 9(16).
           05  FILLER                  PIC X(64).
       FD  REPORT-OUT.
       01  REPORT-LINE                 PIC X(80).

       WORKING-STORAGE SECTION.
       01  WS-FILE-STATUSES.
           05  WS-CONTROL-STATUS       PIC XX.
           05  WS-CARD-IN-STATUS       PIC XX.
           05  WS-ONLINE-STATUS        PIC XX.
           05  WS-REPORT-STATUS        PIC XX.

       01  WS-BUSINESS-DATE            PIC 9(8).
       01  WS-BUSINESS-DATE-PARTS REDEFINES WS-BUSINESS-DATE.
           05  WS-BUSINESS-MONTH       PIC 9(6).
           05  FILLER                  PIC 99.
       01  WS-ABORT-MESSAGE            PIC X(40).
       01  WS-MASTER-END-SW            PIC X          VALUE 'N'.
           88  MASTER-AT-END           VALUE 'Y'.

       01  WS-TOTALS.
           05  WS-READ-COUNT           PIC 9(8)       VALUE 0.
           05  WS-LOADED-COUNT         PIC 9(8)       VALUE 0.
           05  WS-BLOCKED-COUNT        PIC 9(8)       VALUE 0.
           05  WS-EXPIRED-COUNT        PIC 9(8)       VALUE 0.

       01  WS-REPORT-COUNT-LINE.
           05  WS-RC-LABEL             PIC X(30).
           05  WS-RC-VALUE             PIC ZZ,ZZZ,ZZ9.
           05  FILLER                  PIC X(40)      VALUE SPACES.

       PROCEDURE DIVISION.
      * LOAD EVERY CARD, THEN REPORT.
       0000-MAIN.
           PERFORM 1000-INITIALIZE
           PERFORM 2000-LOAD-CARD UNTIL MASTER-AT-END
           PERFORM 8000-WRITE-REPORT
           PERFORM 9000-CLOSE-FILES
           MOVE 0 TO RETURN-CODE
           STOP RUN.

      * OPEN THE FILES, READ THE DAY TO LOAD FOR AND THE FIRST CARD.
       1000-INITIALIZE.
           OPEN INPUT  CONTROL-FILE CARD-IN
                OUTPUT ONLINE-OUT REPORT-OUT
           IF WS-FILE-STATUSES NOT = ALL '0'
               DISPLAY 'CARDLOAD: CANNOT OPEN FILES ' WS-FILE-STATUSES
               MOVE 16 TO RETURN-CODE
               STOP RUN
           END-IF
           READ CONTROL-FILE
               AT END
                   MOVE 'NO BANK CONTROL RECORD' TO WS-ABORT-MESSAGE
                   PERFORM 9900-ABORT
           END-READ
           MOVE BC-BUSINESS-DATE TO WS-BUSINESS-DATE
           PERFORM 1200-READ-CARD.

       1200-READ-CARD.
           READ CARD-IN
               AT END
                   SET MASTER-AT-END TO TRUE
               NOT AT END
                   ADD 1 TO WS-READ-COUNT
           END-READ.

      * WRITE A CARD THAT HAS NOT EXPIRED TO THE ONLINE FILE.
       2000-LOAD-CARD.
           IF CR-EXPIRY-MONTH < WS-BUSINESS-MONTH
               ADD 1 TO WS-EXPIRED-COUNT
           ELSE
               IF CR-BLOCKED
                   ADD 1 TO WS-BLOCKED-COUNT
               END-IF
               PERFORM 2100-WRITE-ONLINE
           END-IF
           PERFORM 1200-READ-CARD.

       2100-WRITE-ONLINE.
           WRITE ONLINE-CARD FROM CARD-RECORD
               INVALID KEY
                   MOVE 'MASTER NOT IN CARD NUMBER ORDER'
                       TO WS-ABORT-MESSAGE
                   PERFORM 9900-ABORT
           END-WRITE
           ADD 1 TO WS-LOADED-COUNT.

      * WRITE THE LOAD REPORT.
       8000-WRITE-REPORT.
           MOVE SPACES TO REPORT-LINE
           STRING 'CARDLOAD ONLINE CARD LOAD        BUSINESS DATE '
                  WS-BUSINESS-DATE
               DELIMITED BY SIZE INTO REPORT-LINE
           WRITE REPORT-LINE
           MOVE SPACES TO REPORT-LINE
           WRITE REPORT-LINE
           MOVE 'CARDS READ' TO WS-RC-LABEL
           MOVE WS-READ-COUNT TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           MOVE 'CARDS LOADED' TO WS-RC-LABEL
           MOVE WS-LOADED-COUNT TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           MOVE '  OF WHICH BLOCKED' TO WS-RC-LABEL
           MOVE WS-BLOCKED-COUNT TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           MOVE 'EXPIRED CARDS LEFT OUT' TO WS-RC-LABEL
           MOVE WS-EXPIRED-COUNT TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE.

      * CLOSE ALL FILES.
       9000-CLOSE-FILES.
           CLOSE CONTROL-FILE CARD-IN ONLINE-OUT REPORT-OUT.

      * STOP THE RUN WITH RETURN CODE 16, SAYING WHY.
       9900-ABORT.
           DISPLAY 'CARDLOAD: ' WS-ABORT-MESSAGE
           PERFORM 9000-CLOSE-FILES
           MOVE 16 TO RETURN-CODE
           STOP RUN.
