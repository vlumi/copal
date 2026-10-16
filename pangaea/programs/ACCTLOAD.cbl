      *================================================================*
      * ACCTLOAD: LOAD THE ONLINE ACCOUNT FILE                         *
      *                                                                *
      * WRITES THE NEW ACCOUNT MASTER'S ACCOUNTS TO A NEW INDEXED      *
      * ONLINE ACCOUNT FILE, FOR THE CARD AND ATM AUTHORIZATIONS OF    *
      * THE NEXT BUSINESS DAY, WITH NO HOLDS, AND STARTS THE DAY'S     *
      * HOLD FILE, EMPTY. CLOSED ACCOUNTS ARE LEFT OUT, SO THAT        *
      * AUTHORIZATIONS FIND NO SUCH ACCOUNT. THE MASTER IS IN ACCOUNT  *
      * NUMBER ORDER, SO THE FILE IS WRITTEN IN KEY ORDER,             *
      * SEQUENTIALLY.                                                  *
      *                                                                *
      * RETURN CODES: 0 LOADED, 16 MASTER OUT OF ORDER OR I/O.         *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 2003-05-19 MV  CR-0342 ORIGINAL VERSION.                       *
      *================================================================*
       IDENTIFICATION DIVISION.
       PROGRAM-ID. ACCTLOAD.
       AUTHOR. M. VARGA, ONLINE SYSTEMS.
       INSTALLATION. FIRST BANK OF PANGAEA, DATA CENTER.
       DATE-WRITTEN. 2003-05-19.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT CONTROL-FILE ASSIGN TO BANKCTL
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-CONTROL-STATUS.
           SELECT ACCT-IN      ASSIGN TO ACCTIN
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-ACCT-IN-STATUS.
           SELECT ONLINE-OUT   ASSIGN TO ACCTOUT
               ORGANIZATION IS INDEXED
               ACCESS MODE IS SEQUENTIAL
               RECORD KEY IS AO-ACCT-NO
               FILE STATUS IS WS-ONLINE-STATUS.
           SELECT HOLD-OUT     ASSIGN TO HOLDOUT
               ORGANIZATION IS INDEXED
               ACCESS MODE IS SEQUENTIAL
               RECORD KEY IS HD-KEY
               FILE STATUS IS WS-HOLD-STATUS.
           SELECT REPORT-OUT   ASSIGN TO LOADRPT
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-STATUS.

       DATA DIVISION.
       FILE SECTION.
       FD  CONTROL-FILE.
           COPY BANKCTL.
       FD  ACCT-IN.
           COPY ACCTMAST.
       FD  ONLINE-OUT.
           COPY ACCTONL.
       FD  HOLD-OUT.
           COPY HOLDREC.
       FD  REPORT-OUT.
       01  REPORT-LINE                 PIC X(80).

       WORKING-STORAGE SECTION.
       01  WS-FILE-STATUSES.
           05  WS-CONTROL-STATUS       PIC XX.
           05  WS-ACCT-IN-STATUS       PIC XX.
           05  WS-ONLINE-STATUS        PIC XX.
           05  WS-HOLD-STATUS          PIC XX.
           05  WS-REPORT-STATUS        PIC XX.

       01  WS-BUSINESS-DATE            PIC 9(8).
       01  WS-ABORT-MESSAGE            PIC X(40).
       01  WS-MASTER-END-SW            PIC X          VALUE 'N'.
           88  MASTER-AT-END           VALUE 'Y'.

       01  WS-TOTALS.
           05  WS-READ-COUNT           PIC 9(8)       VALUE 0.
           05  WS-LOADED-COUNT         PIC 9(8)       VALUE 0.
           05  WS-CLOSED-COUNT         PIC 9(8)       VALUE 0.
           05  WS-LOADED-BALANCE       PIC S9(11)V99  VALUE 0.

       01  WS-REPORT-COUNT-LINE.
           05  WS-RC-LABEL             PIC X(30).
           05  WS-RC-VALUE             PIC ZZ,ZZZ,ZZ9.
           05  FILLER                  PIC X(40)      VALUE SPACES.
       01  WS-REPORT-AMOUNT-LINE.
           05  WS-RA-LABEL             PIC X(30).
           05  WS-RA-VALUE             PIC Z,ZZZ,ZZZ,ZZZ,ZZ9.99-.
           05  FILLER                  PIC X(29)      VALUE SPACES.

       PROCEDURE DIVISION.
      * LOAD EVERY ACCOUNT, THEN REPORT.
       0000-MAIN.
           PERFORM 1000-INITIALIZE
           PERFORM 2000-LOAD-ACCOUNT UNTIL MASTER-AT-END
           PERFORM 8000-WRITE-REPORT
           PERFORM 9000-CLOSE-FILES
           MOVE 0 TO RETURN-CODE
           STOP RUN.

      * OPEN THE FILES AND READ THE FIRST ACCOUNT. THE HOLD FILE IS
      * OPENED ONLY TO START IT.
       1000-INITIALIZE.
           OPEN INPUT  CONTROL-FILE ACCT-IN
                OUTPUT ONLINE-OUT HOLD-OUT REPORT-OUT
           IF WS-FILE-STATUSES NOT = ALL '0'
               DISPLAY 'ACCTLOAD: CANNOT OPEN FILES ' WS-FILE-STATUSES
               MOVE 16 TO RETURN-CODE
               STOP RUN
           END-IF
           READ CONTROL-FILE
               AT END
                   MOVE 'NO BANK CONTROL RECORD' TO WS-ABORT-MESSAGE
                   PERFORM 9900-ABORT
           END-READ
           MOVE BC-BUSINESS-DATE TO WS-BUSINESS-DATE
           PERFORM 1200-READ-ACCOUNT.

       1200-READ-ACCOUNT.
           READ ACCT-IN
               AT END
                   SET MASTER-AT-END TO TRUE
               NOT AT END
                   ADD 1 TO WS-READ-COUNT
           END-READ.

      * WRITE AN ACCOUNT THAT IS NOT CLOSED TO THE ONLINE FILE.
       2000-LOAD-ACCOUNT.
           IF AM-CLOSED
               ADD 1 TO WS-CLOSED-COUNT
           ELSE
               PERFORM 2100-WRITE-ONLINE
           END-IF
           PERFORM 1200-READ-ACCOUNT.

       2100-WRITE-ONLINE.
           MOVE SPACES TO ACCT-ONLINE-RECORD
           MOVE AM-ACCT-NO TO AO-ACCT-NO
           MOVE AM-CUST-NO TO AO-CUST-NO
           MOVE AM-SHORT-NAME TO AO-SHORT-NAME
           MOVE AM-ACCT-TYPE TO AO-ACCT-TYPE
           MOVE AM-STATUS TO AO-STATUS
           MOVE AM-BALANCE TO AO-BALANCE
           MOVE AM-OVERDRAFT-LIMIT TO AO-OVERDRAFT-LIMIT
           MOVE 0 TO AO-HOLDS
           MOVE WS-BUSINESS-DATE TO AO-AS-OF-DATE
           WRITE ACCT-ONLINE-RECORD
               INVALID KEY
                   MOVE 'MASTER NOT IN ACCOUNT NUMBER ORDER'
                       TO WS-ABORT-MESSAGE
                   PERFORM 9900-ABORT
           END-WRITE
           ADD 1 TO WS-LOADED-COUNT
           ADD AM-BALANCE TO WS-LOADED-BALANCE.

      * WRITE THE LOAD REPORT.
       8000-WRITE-REPORT.
           MOVE SPACES TO REPORT-LINE
           STRING 'ACCTLOAD ONLINE ACCOUNT LOAD     BUSINESS DATE '
                  WS-BUSINESS-DATE
               DELIMITED BY SIZE INTO REPORT-LINE
           WRITE REPORT-LINE
           MOVE SPACES TO REPORT-LINE
           WRITE REPORT-LINE
           MOVE 'ACCOUNTS READ' TO WS-RC-LABEL
           MOVE WS-READ-COUNT TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           MOVE 'ACCOUNTS LOADED' TO WS-RC-LABEL
           MOVE WS-LOADED-COUNT TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           MOVE 'CLOSED ACCOUNTS LEFT OUT' TO WS-RC-LABEL
           MOVE WS-CLOSED-COUNT TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           MOVE SPACES TO REPORT-LINE
           WRITE REPORT-LINE
           MOVE 'BALANCES LOADED' TO WS-RA-LABEL
           MOVE WS-LOADED-BALANCE TO WS-RA-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-AMOUNT-LINE.

      * CLOSE ALL FILES.
       9000-CLOSE-FILES.
           CLOSE CONTROL-FILE ACCT-IN ONLINE-OUT HOLD-OUT REPORT-OUT.

      * STOP THE RUN WITH RETURN CODE 16, SAYING WHY.
       9900-ABORT.
           DISPLAY 'ACCTLOAD: ' WS-ABORT-MESSAGE
           PERFORM 9000-CLOSE-FILES
           MOVE 16 TO RETURN-CODE
           STOP RUN.
