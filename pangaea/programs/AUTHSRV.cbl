      *================================================================*
      * AUTHSRV: SERVE THE DAY'S CARD AND ATM AUTHORIZATIONS           *
      *                                                                *
      * RUNS FOR THE WHOLE ONLINE DAY, AS A TRANSACTION MONITOR DOES.  *
      * READS ONE AUTHORIZATION REQUEST AT A TIME FROM SYSIN, READS    *
      * THE CARD AND ITS ACCOUNT FROM THE ONLINE FILES, AND FOR A      *
      * REVERSAL THE HOLD, CALLS THE PROGRAM FOR THE REQUEST, AND      *
      * WRITES THE RESPONSE TO SYSOUT. WHEN A REQUEST IS APPROVED IT   *
      * REWRITES THE ACCOUNT AND THE CARD AND WRITES OR REWRITES THE   *
      * HOLD. THE CALLED PROGRAMS DECIDE; ONLY AUTHSRV OPENS FILES, SO *
      * THAT EACH FILE HAS ONE READER AND WRITER. ENDS WITH THE        *
      * REQUESTS, WHEN THE BATCH WINDOW OPENS, AND REPORTS THE DAY.    *
      *                                                                *
      * HOLD NUMBERS, THE AUTHORIZATION CODES, NUMBER THE DAY'S HOLDS  *
      * FROM 1, AND GO ON FROM THE HIGHEST IN THE HOLD FILE WHEN THE   *
      * SERVER IS STARTED AGAIN DURING THE DAY.                        *
      *                                                                *
      * RETURN CODES: 0 SERVED, 16 I/O.                                *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 2003-05-19 MV  CR-0342 ORIGINAL VERSION.                       *
      *================================================================*
       IDENTIFICATION DIVISION.
       PROGRAM-ID. AUTHSRV.
       AUTHOR. M. VARGA, ONLINE SYSTEMS.
       INSTALLATION. FIRST BANK OF PANGAEA, DATA CENTER.
       DATE-WRITTEN. 2003-05-19.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT CONTROL-FILE ASSIGN TO BANKCTL
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-CONTROL-STATUS.
           SELECT ACCT-FILE    ASSIGN TO ACCTONL
               ORGANIZATION IS INDEXED
               ACCESS MODE IS RANDOM
               RECORD KEY IS AO-ACCT-NO
               FILE STATUS IS WS-ACCT-STATUS.
           SELECT CARD-FILE    ASSIGN TO CARDONL
               ORGANIZATION IS INDEXED
               ACCESS MODE IS RANDOM
               RECORD KEY IS CR-CARD-NO
               FILE STATUS IS WS-CARD-STATUS.
           SELECT HOLD-FILE    ASSIGN TO HOLDS
               ORGANIZATION IS INDEXED
               ACCESS MODE IS DYNAMIC
               RECORD KEY IS HD-KEY
               FILE STATUS IS WS-HOLD-STATUS.
           SELECT REPORT-OUT   ASSIGN TO AUTHRPT
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-STATUS.

       DATA DIVISION.
       FILE SECTION.
       FD  CONTROL-FILE.
           COPY BANKCTL.
       FD  ACCT-FILE.
           COPY ACCTONL.
       FD  CARD-FILE.
           COPY CARDREC.
       FD  HOLD-FILE.
           COPY HOLDREC.
       FD  REPORT-OUT.
       01  REPORT-LINE                 PIC X(80).

       WORKING-STORAGE SECTION.
       01  WS-FILE-STATUSES.
           05  WS-CONTROL-STATUS       PIC XX.
           05  WS-ACCT-STATUS          PIC XX.
           05  WS-CARD-STATUS          PIC XX.
           05  WS-HOLD-STATUS          PIC XX.
           05  WS-REPORT-STATUS        PIC XX.

       01  WS-ABORT-MESSAGE            PIC X(40).
       01  WS-INPUT-END-SW             PIC X          VALUE 'N'.
           88  END-OF-INPUT            VALUE 'Y'.
       01  WS-HOLDS-END-SW             PIC X          VALUE 'N'.
           88  END-OF-HOLDS            VALUE 'Y'.
       01  WS-LAST-HOLD-NO             PIC 9(6)       VALUE 0.

           COPY AUTHMSG.

       01  WS-TOTALS.
           05  WS-REQUEST-COUNT        PIC 9(8)       VALUE 0.
           05  WS-APPROVED-COUNT       PIC 9(8)       VALUE 0.
           05  WS-DECLINED-COUNT       PIC 9(8)       VALUE 0.
           05  WS-HELD-TOTAL           PIC 9(11)V99   VALUE 0.
           05  WS-RELEASED-TOTAL       PIC 9(11)V99   VALUE 0.

       01  WS-REPORT-COUNT-LINE.
           05  WS-RC-LABEL             PIC X(30).
           05  WS-RC-VALUE             PIC ZZ,ZZZ,ZZ9.
           05  FILLER                  PIC X(40)      VALUE SPACES.
       01  WS-REPORT-AMOUNT-LINE.
           05  WS-RA-LABEL             PIC X(30).
           05  WS-RA-VALUE             PIC Z,ZZZ,ZZZ,ZZZ,ZZ9.99-.
           05  FILLER                  PIC X(29)      VALUE SPACES.

       PROCEDURE DIVISION.
      * SERVE REQUESTS UNTIL THERE ARE NO MORE, THEN REPORT.
       0000-MAIN.
           PERFORM 1000-INITIALIZE
           PERFORM 2000-SERVE-REQUEST UNTIL END-OF-INPUT
           PERFORM 8000-WRITE-REPORT
           PERFORM 9000-CLOSE-FILES
           MOVE 0 TO RETURN-CODE
           STOP RUN.

      * OPEN THE FILES, READ THE BUSINESS DATE AND FIND THE LAST HOLD
      * NUMBER GIVEN TODAY.
       1000-INITIALIZE.
           OPEN INPUT  CONTROL-FILE
                I-O    ACCT-FILE CARD-FILE HOLD-FILE
                OUTPUT REPORT-OUT
           IF WS-FILE-STATUSES NOT = ALL '0'
               DISPLAY 'AUTHSRV: CANNOT OPEN FILES ' WS-FILE-STATUSES
                   UPON SYSERR
               MOVE 16 TO RETURN-CODE
               STOP RUN
           END-IF
           READ CONTROL-FILE
               AT END
                   MOVE 'NO BANK CONTROL RECORD' TO WS-ABORT-MESSAGE
                   PERFORM 9900-ABORT
           END-READ
           PERFORM 1100-FIND-LAST-HOLD.

      * READ THROUGH THE HOLDS FOR THE HIGHEST HOLD NUMBER.
       1100-FIND-LAST-HOLD.
           MOVE LOW-VALUES TO HD-KEY
           START HOLD-FILE KEY IS NOT LESS THAN HD-KEY
               INVALID KEY
                   SET END-OF-HOLDS TO TRUE
           END-START
           PERFORM UNTIL END-OF-HOLDS
               READ HOLD-FILE NEXT RECORD
                   AT END
                       SET END-OF-HOLDS TO TRUE
                   NOT AT END
                       IF HD-HOLD-NO > WS-LAST-HOLD-NO
                           MOVE HD-HOLD-NO TO WS-LAST-HOLD-NO
                       END-IF
               END-READ
           END-PERFORM.

      * READ THE NEXT REQUEST, AND ANSWER IT.
       2000-SERVE-REQUEST.
           ACCEPT AUTH-REQUEST
               ON EXCEPTION
                   SET END-OF-INPUT TO TRUE
               NOT ON EXCEPTION
                   PERFORM 2100-AUTHORIZE
                   DISPLAY AUTH-RESPONSE
           END-ACCEPT.

      * FIND THE RECORDS THE REQUEST IS ABOUT, CALL THE PROGRAM FOR IT,
      * AND KEEP WHAT IT CHANGED IF IT APPROVED.
       2100-AUTHORIZE.
           ADD 1 TO WS-REQUEST-COUNT
           MOVE SPACES TO AUTH-RESPONSE HOLD-RECORD
           IF RQ-PURCHASE OR RQ-REVERSAL OR RQ-WITHDRAWAL
              OR RQ-BALANCE-ENQUIRY
               PERFORM 2200-READ-CARD-AND-ACCOUNT
           ELSE
               SET RS-INVALID-REQUEST TO TRUE
               MOVE 'UNKNOWN REQUEST' TO RS-MESSAGE
           END-IF
           IF RS-RESULT = SPACES AND RQ-REVERSAL
               PERFORM 2300-READ-HOLD
           END-IF
           IF RS-RESULT = SPACES
               EVALUATE TRUE
                   WHEN RQ-PURCHASE
                       CALL 'AUTHPUR' USING AUTH-REQUEST AUTH-RESPONSE
                           BANK-CONTROL-RECORD ACCT-ONLINE-RECORD
                           CARD-RECORD HOLD-RECORD
                   WHEN RQ-REVERSAL
                       CALL 'AUTHREV' USING AUTH-REQUEST AUTH-RESPONSE
                           BANK-CONTROL-RECORD ACCT-ONLINE-RECORD
                           CARD-RECORD HOLD-RECORD
                   WHEN RQ-WITHDRAWAL
                       CALL 'ATMWDL' USING AUTH-REQUEST AUTH-RESPONSE
                           BANK-CONTROL-RECORD ACCT-ONLINE-RECORD
                           CARD-RECORD HOLD-RECORD
                   WHEN RQ-BALANCE-ENQUIRY
                       CALL 'ATMBAL' USING AUTH-REQUEST AUTH-RESPONSE
                           BANK-CONTROL-RECORD ACCT-ONLINE-RECORD
                           CARD-RECORD HOLD-RECORD
               END-EVALUATE
           END-IF
           IF RS-APPROVED
               ADD 1 TO WS-APPROVED-COUNT
               IF NOT RQ-BALANCE-ENQUIRY
                   PERFORM 2400-SAVE-CHANGES
               END-IF
           ELSE
               ADD 1 TO WS-DECLINED-COUNT
           END-IF.

      * READ THE CARD, AND THE ACCOUNT IT DRAWS ON.
       2200-READ-CARD-AND-ACCOUNT.
           MOVE RQ-CARD-NO TO CR-CARD-NO
           READ CARD-FILE
               INVALID KEY
                   SET RS-NO-SUCH-CARD TO TRUE
                   MOVE 'NO SUCH CARD' TO RS-MESSAGE
               NOT INVALID KEY
                   MOVE CR-ACCT-NO TO AO-ACCT-NO
                   READ ACCT-FILE
                       INVALID KEY
                           SET RS-ACCOUNT-NOT-OPEN TO TRUE
                           MOVE 'NO SUCH ACCOUNT' TO RS-MESSAGE
                   END-READ
           END-READ.

      * READ THE HOLD OF THE AUTHORIZATION TO REVERSE.
       2300-READ-HOLD.
           MOVE CR-ACCT-NO TO HD-ACCT-NO
           MOVE RQ-AUTH-CODE TO HD-HOLD-NO
           READ HOLD-FILE
               INVALID KEY
                   SET RS-NO-SUCH-AUTH TO TRUE
                   MOVE 'NO SUCH AUTHORIZATION' TO RS-MESSAGE
           END-READ.

      * REWRITE THE ACCOUNT AND THE CARD, AND WRITE A NEW HOLD UNDER THE
      * NEXT HOLD NUMBER OR REWRITE THE REVERSED ONE.
       2400-SAVE-CHANGES.
           REWRITE ACCT-ONLINE-RECORD
               INVALID KEY
                   MOVE 'CANNOT REWRITE THE ACCOUNT' TO WS-ABORT-MESSAGE
                   PERFORM 9900-ABORT
           END-REWRITE
           REWRITE CARD-RECORD
               INVALID KEY
                   MOVE 'CANNOT REWRITE THE CARD' TO WS-ABORT-MESSAGE
                   PERFORM 9900-ABORT
           END-REWRITE
           IF RQ-REVERSAL
               REWRITE HOLD-RECORD
                   INVALID KEY
                       MOVE 'CANNOT REWRITE THE HOLD'
                           TO WS-ABORT-MESSAGE
                       PERFORM 9900-ABORT
               END-REWRITE
               ADD HD-AMOUNT TO WS-RELEASED-TOTAL
           ELSE
               ADD 1 TO WS-LAST-HOLD-NO
               MOVE WS-LAST-HOLD-NO TO HD-HOLD-NO RS-AUTH-CODE
               WRITE HOLD-RECORD
                   INVALID KEY
                       MOVE 'CANNOT WRITE THE HOLD' TO WS-ABORT-MESSAGE
                       PERFORM 9900-ABORT
               END-WRITE
               ADD HD-AMOUNT TO WS-HELD-TOTAL
               ADD HD-FEE TO WS-HELD-TOTAL
           END-IF.

      * WRITE THE DAY'S REPORT.
       8000-WRITE-REPORT.
           MOVE SPACES TO REPORT-LINE
           STRING 'AUTHSRV  ONLINE AUTHORIZATIONS   BUSINESS DATE '
                  BC-BUSINESS-DATE
               DELIMITED BY SIZE INTO REPORT-LINE
           WRITE REPORT-LINE
           MOVE SPACES TO REPORT-LINE
           WRITE REPORT-LINE
           MOVE 'REQUESTS' TO WS-RC-LABEL
           MOVE WS-REQUEST-COUNT TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           MOVE 'APPROVED' TO WS-RC-LABEL
           MOVE WS-APPROVED-COUNT TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           MOVE 'DECLINED' TO WS-RC-LABEL
           MOVE WS-DECLINED-COUNT TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           MOVE 'LAST AUTHORIZATION CODE' TO WS-RC-LABEL
           MOVE WS-LAST-HOLD-NO TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           MOVE SPACES TO REPORT-LINE
           WRITE REPORT-LINE
           MOVE 'HELD WITH FEES' TO WS-RA-LABEL
           MOVE WS-HELD-TOTAL TO WS-RA-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-AMOUNT-LINE
           MOVE 'RELEASED BY REVERSALS' TO WS-RA-LABEL
           MOVE WS-RELEASED-TOTAL TO WS-RA-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-AMOUNT-LINE.

      * CLOSE ALL FILES.
       9000-CLOSE-FILES.
           CLOSE CONTROL-FILE ACCT-FILE CARD-FILE HOLD-FILE REPORT-OUT.

      * STOP THE RUN WITH RETURN CODE 16, SAYING WHY.
       9900-ABORT.
           DISPLAY 'AUTHSRV: ' WS-ABORT-MESSAGE UPON SYSERR
           PERFORM 9000-CLOSE-FILES
           MOVE 16 TO RETURN-CODE
           STOP RUN.
