      *================================================================*
      * TXVALID: VALIDATE THE DAY'S TRANSACTION FILE                   *
      *                                                                *
      * CHECKS THE FILE AS A WHOLE: A HEADER FOR TODAY'S BUSINESS      *
      * DATE, THEN DETAILS, THEN A TRAILER WHOSE CONTROL TOTALS MATCH  *
      * THE DETAILS. A FILE THAT FAILS THESE CHECKS IS NOT PROCESSED.  *
      * CHECKS EACH DETAIL'S FIELDS AND WRITES IT EITHER TO THE VALID  *
      * FILE OR TO THE REJECTS. WHETHER THE ACCOUNT EXISTS IS LEFT TO  *
      * POSTING, WHICH READS THE MASTER.                               *
      *                                                                *
      * RETURN CODES: 0 ALL VALID, 4 SOME DETAILS REJECTED,            *
      *               12 FILE REJECTED, 16 I/O ERROR.                  *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 1987-03-12 JA  ORIGINAL VERSION.                               *
      * 1991-06-03 PK  CR-0117 ACCOUNT NUMBER HASH IN THE TRAILER.     *
      * 1999-01-11 SL  Y2K: 8-DIGIT BUSINESS DATE FROM THE BANK        *
      *                CONTROL RECORD.                                 *
      * 2004-09-20 PK  CR-0388 REJECT RECORDS OF UNKNOWN TYPE          *
      *                INSTEAD OF STOPPING THE RUN.                    *
      *================================================================*
       IDENTIFICATION DIVISION.
       PROGRAM-ID. TXVALID.
       AUTHOR. J. ARKWRIGHT, BATCH SYSTEMS.
       INSTALLATION. FIRST BANK OF PANGAEA, DATA CENTER.
       DATE-WRITTEN. 1987-03-12.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT CONTROL-FILE ASSIGN TO BANKCTL
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-CONTROL-STATUS.
           SELECT TRAN-IN      ASSIGN TO TRANIN
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-TRAN-IN-STATUS.
           SELECT TRAN-OUT     ASSIGN TO TRANOUT
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-TRAN-OUT-STATUS.
           SELECT REJECT-OUT   ASSIGN TO REJECTS
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-REJECT-STATUS.
           SELECT REPORT-OUT   ASSIGN TO VALRPT
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-STATUS.

       DATA DIVISION.
       FILE SECTION.
       FD  CONTROL-FILE.
           COPY BANKCTL.
       FD  TRAN-IN.
           COPY TRANREC.
       FD  TRAN-OUT.
       01  TRAN-OUT-RECORD             PIC X(80).
       FD  REJECT-OUT.
       01  REJECT-OUT-RECORD           PIC X(100).
       FD  REPORT-OUT.
       01  REPORT-LINE                 PIC X(80).

       WORKING-STORAGE SECTION.
       01  WS-FILE-STATUSES.
           05  WS-CONTROL-STATUS       PIC XX.
           05  WS-TRAN-IN-STATUS       PIC XX.
               88  TRAN-IN-EOF         VALUE '10'.
           05  WS-TRAN-OUT-STATUS      PIC XX.
           05  WS-REJECT-STATUS        PIC XX.
           05  WS-REPORT-STATUS        PIC XX.

       01  WS-BUSINESS-DATE            PIC 9(8).

           COPY REJREC.

       01  WS-FILE-ERROR               PIC X(40) VALUE SPACES.
           88  FILE-IS-GOOD            VALUE SPACES.

       01  WS-DETAIL-STATE             PIC X.
           88  DETAIL-IS-VALID         VALUE 'Y'.
           88  DETAIL-IS-REJECTED      VALUE 'N'.

       01  WS-TOTALS.
           05  WS-DETAIL-COUNT         PIC 9(8)       VALUE 0.
           05  WS-VALID-COUNT          PIC 9(8)       VALUE 0.
           05  WS-REJECT-COUNT         PIC 9(8)       VALUE 0.
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
      * VALIDATE THE FILE AND SET THE RETURN CODE.
       0000-MAIN.
           PERFORM 1000-OPEN-FILES
           PERFORM 2000-READ-HEADER
           IF FILE-IS-GOOD
               PERFORM 3000-READ-DETAILS
           END-IF
           IF FILE-IS-GOOD
               PERFORM 4000-CHECK-TRAILER
           END-IF
           PERFORM 8000-WRITE-REPORT
           PERFORM 9000-CLOSE-FILES
           EVALUATE TRUE
               WHEN NOT FILE-IS-GOOD
                   MOVE 12 TO RETURN-CODE
               WHEN WS-REJECT-COUNT > 0
                   MOVE 4 TO RETURN-CODE
               WHEN OTHER
                   MOVE 0 TO RETURN-CODE
           END-EVALUATE
           STOP RUN.

      * OPEN THE FILES AND READ THE BUSINESS DATE.
       1000-OPEN-FILES.
           OPEN INPUT  CONTROL-FILE TRAN-IN
                OUTPUT TRAN-OUT REJECT-OUT REPORT-OUT
           IF WS-CONTROL-STATUS NOT = '00'
              OR WS-TRAN-IN-STATUS NOT = '00'
              OR WS-TRAN-OUT-STATUS NOT = '00'
              OR WS-REJECT-STATUS NOT = '00'
              OR WS-REPORT-STATUS NOT = '00'
               DISPLAY 'TXVALID: CANNOT OPEN FILES ' WS-FILE-STATUSES
               MOVE 16 TO RETURN-CODE
               STOP RUN
           END-IF
           READ CONTROL-FILE
               AT END
                   DISPLAY 'TXVALID: NO BANK CONTROL RECORD'
                   MOVE 16 TO RETURN-CODE
                   STOP RUN
           END-READ
           MOVE BC-BUSINESS-DATE TO WS-BUSINESS-DATE
           MOVE SPACES TO REJECT-RECORD.

      * THE FIRST RECORD MUST BE THE HEADER FOR THE BUSINESS DATE.
       2000-READ-HEADER.
           READ TRAN-IN
               AT END
                   MOVE 'EMPTY FILE' TO WS-FILE-ERROR
           END-READ
           EVALUATE TRUE
               WHEN NOT FILE-IS-GOOD
                   CONTINUE
               WHEN NOT TR-HEADER
                   MOVE 'FIRST RECORD IS NOT A HEADER' TO WS-FILE-ERROR
               WHEN TH-BUSINESS-DATE NOT = WS-BUSINESS-DATE
                   MOVE 'HEADER IS NOT FOR THE BUSINESS DATE'
                       TO WS-FILE-ERROR
           END-EVALUATE.

      * VALIDATE EACH RECORD UNTIL THE TRAILER.
       3000-READ-DETAILS.
           PERFORM 3100-READ-NEXT
           PERFORM UNTIL TR-TRAILER OR NOT FILE-IS-GOOD
               EVALUATE TRUE
                   WHEN TR-DETAIL
                       PERFORM 3200-VALIDATE-DETAIL
                   WHEN TR-HEADER
                       MOVE 'SECOND HEADER' TO WS-FILE-ERROR
                   WHEN OTHER
                       MOVE 'TYPE' TO RJ-REASON
                       PERFORM 3300-REJECT-RECORD
               END-EVALUATE
               IF FILE-IS-GOOD
                   PERFORM 3100-READ-NEXT
               END-IF
           END-PERFORM.

      * READ ONE RECORD. THE FILE MUST NOT END BEFORE ITS TRAILER.
       3100-READ-NEXT.
           READ TRAN-IN
               AT END
                   MOVE 'NO TRAILER' TO WS-FILE-ERROR
           END-READ.

      * ADD A DETAIL TO THE TOTALS, THEN CHECK ITS FIELDS.
       3200-VALIDATE-DETAIL.
           ADD 1 TO WS-DETAIL-COUNT
           IF TD-AMOUNT IS NUMERIC
               ADD TD-AMOUNT TO WS-AMOUNT-TOTAL
           END-IF
           IF TD-ACCT-NO IS NUMERIC
               ADD TD-ACCT-NO TO WS-ACCT-HASH
           END-IF
           SET DETAIL-IS-VALID TO TRUE
           EVALUATE TRUE
               WHEN TD-ACCT-NO IS NOT NUMERIC
                   SET RJ-BAD-ACCOUNT TO TRUE
                   SET DETAIL-IS-REJECTED TO TRUE
               WHEN NOT TD-KNOWN-CODE
                   SET RJ-BAD-CODE TO TRUE
                   SET DETAIL-IS-REJECTED TO TRUE
               WHEN TD-AMOUNT IS NOT NUMERIC
                 OR TD-AMOUNT = ZERO
                   SET RJ-BAD-AMOUNT TO TRUE
                   SET DETAIL-IS-REJECTED TO TRUE
           END-EVALUATE
           IF DETAIL-IS-VALID
               WRITE TRAN-OUT-RECORD FROM TRAN-RECORD
               ADD 1 TO WS-VALID-COUNT
           ELSE
               PERFORM 3300-REJECT-RECORD
           END-IF.

      * WRITE THE RECORD TO THE REJECTS WITH ITS REASON.
       3300-REJECT-RECORD.
           MOVE 'TXVALID' TO RJ-PROGRAM
           MOVE TRAN-RECORD TO RJ-TRAN
           WRITE REJECT-OUT-RECORD FROM REJECT-RECORD
           ADD 1 TO WS-REJECT-COUNT.

      * THE TRAILER MUST MATCH THE TOTALS, AND END THE FILE.
       4000-CHECK-TRAILER.
           EVALUATE TRUE
               WHEN TT-DETAIL-COUNT NOT = WS-DETAIL-COUNT
                   MOVE 'TRAILER DETAIL COUNT DIFFERS'
                       TO WS-FILE-ERROR
               WHEN TT-AMOUNT-TOTAL NOT = WS-AMOUNT-TOTAL
                   MOVE 'TRAILER AMOUNT TOTAL DIFFERS'
                       TO WS-FILE-ERROR
               WHEN TT-ACCT-HASH NOT = WS-ACCT-HASH
                   MOVE 'TRAILER ACCOUNT HASH DIFFERS'
                       TO WS-FILE-ERROR
           END-EVALUATE
           IF FILE-IS-GOOD
               READ TRAN-IN
                   NOT AT END
                       MOVE 'RECORDS AFTER THE TRAILER'
                           TO WS-FILE-ERROR
               END-READ
           END-IF.

      * WRITE THE VALIDATION REPORT.
       8000-WRITE-REPORT.
           MOVE SPACES TO REPORT-LINE
           STRING 'TXVALID  TRANSACTION VALIDATION  BUSINESS DATE '
                  WS-BUSINESS-DATE
               DELIMITED BY SIZE INTO REPORT-LINE
           WRITE REPORT-LINE
           MOVE SPACES TO REPORT-LINE
           WRITE REPORT-LINE
           MOVE 'DETAILS READ' TO WS-RC-LABEL
           MOVE WS-DETAIL-COUNT TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           MOVE 'DETAILS VALID' TO WS-RC-LABEL
           MOVE WS-VALID-COUNT TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           MOVE 'RECORDS REJECTED' TO WS-RC-LABEL
           MOVE WS-REJECT-COUNT TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           MOVE 'AMOUNT TOTAL' TO WS-RA-LABEL
           MOVE WS-AMOUNT-TOTAL TO WS-RA-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-AMOUNT-LINE
           MOVE SPACES TO REPORT-LINE
           WRITE REPORT-LINE
           MOVE SPACES TO REPORT-LINE
           IF FILE-IS-GOOD
               MOVE 'FILE ACCEPTED' TO REPORT-LINE
           ELSE
               STRING 'FILE REJECTED: ' WS-FILE-ERROR
                   DELIMITED BY SIZE INTO REPORT-LINE
           END-IF
           WRITE REPORT-LINE.

      * CLOSE ALL FILES.
       9000-CLOSE-FILES.
           CLOSE CONTROL-FILE TRAN-IN TRAN-OUT REJECT-OUT REPORT-OUT.
