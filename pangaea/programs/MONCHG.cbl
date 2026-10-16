      *================================================================*
      * MONCHG: MONTH-END INTEREST AND FEES                            *
      *                                                                *
      * ON THE DAY INTACCR CLOSES A MONTH, WRITES THE MONTH'S CHARGES  *
      * AND CREDITS AS TRANSACTIONS FOR TXPOST: THE INTEREST DUE,      *
      * CREDITED (IN) OR CHARGED AS OVERDRAFT INTEREST (OD), AND THE   *
      * MAINTENANCE FEE (MF) OF AN OPEN CURRENT ACCOUNT WHOSE BALANCE  *
      * CLOSED THE MONTH BELOW THE WAIVER BALANCE. ON OTHER DAYS IT    *
      * WRITES NOTHING.                                                *
      *                                                                *
      * THE MASTER AND THE ACCRUALS ARE READ IN STEP, SO THE           *
      * TRANSACTIONS COME OUT IN ACCOUNT NUMBER ORDER, AND TXPOST      *
      * POSTS THEM WITHOUT A SORT. TRANSACTION IDS ARE 00, THE CLOSED  *
      * MONTH AS YYMM, AND A SEQUENCE NUMBER.                          *
      *                                                                *
      * RETURN CODES: 0 WRITTEN, 16 FILES OUT OF STEP OR I/O.          *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 2001-09-03 TL  CR-0298 ORIGINAL VERSION.                       *
      *================================================================*
       IDENTIFICATION DIVISION.
       PROGRAM-ID. MONCHG.
       AUTHOR. T. LINDQVIST, BATCH SYSTEMS.
       INSTALLATION. FIRST BANK OF PANGAEA, DATA CENTER.
       DATE-WRITTEN. 2001-09-03.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT CONTROL-FILE ASSIGN TO BANKCTL
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-CONTROL-STATUS.
           SELECT ACCT-IN      ASSIGN TO ACCTIN
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-ACCT-IN-STATUS.
           SELECT ACCR-IN      ASSIGN TO ACCRIN
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-ACCR-IN-STATUS.
           SELECT TRAN-OUT     ASSIGN TO TRANOUT
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-TRAN-OUT-STATUS.
           SELECT REPORT-OUT   ASSIGN TO MONRPT
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-STATUS.

       DATA DIVISION.
       FILE SECTION.
       FD  CONTROL-FILE.
           COPY BANKCTL.
       FD  ACCT-IN.
           COPY ACCTMAST.
       FD  ACCR-IN.
           COPY ACCRREC.
       FD  TRAN-OUT.
       01  TRAN-OUT-RECORD             PIC X(80).
       FD  REPORT-OUT.
       01  REPORT-LINE                 PIC X(80).

       WORKING-STORAGE SECTION.
       01  WS-FILE-STATUSES.
           05  WS-CONTROL-STATUS       PIC XX.
           05  WS-ACCT-IN-STATUS       PIC XX.
           05  WS-ACCR-IN-STATUS       PIC XX.
           05  WS-TRAN-OUT-STATUS      PIC XX.
           05  WS-REPORT-STATUS        PIC XX.

       01  WS-BUSINESS-DATE            PIC 9(8).
       01  WS-ABORT-MESSAGE            PIC X(40).
       01  WS-END-SWITCHES.
           05  WS-MASTER-END-SW        PIC X          VALUE 'N'.
               88  MASTER-AT-END       VALUE 'Y'.
           05  WS-ACCRUAL-END-SW       PIC X          VALUE 'N'.
               88  ACCRUALS-AT-END     VALUE 'Y'.

      * THE MAINTENANCE FEE, AND THE BALANCE THAT WAIVES IT.
       01  WS-MAINTENANCE-FEE          PIC 9(3)V99    VALUE 2.50.
       01  WS-FEE-WAIVER-BALANCE       PIC 9(7)V99    VALUE 1000.00.

       01  WS-CLOSED-MONTH             PIC 9(6)       VALUE 0.
       01  WS-CLOSED-MONTH-PARTS REDEFINES WS-CLOSED-MONTH.
           05  WS-CM-CENTURY           PIC 99.
           05  WS-CM-YYMM              PIC 9(4).
       01  WS-TRAN-ID.
           05  FILLER                  PIC 99         VALUE 0.
           05  WS-TI-YYMM              PIC 9(4).
           05  WS-TI-SEQ               PIC 9(4)       VALUE 0.

           COPY TRANREC.

       01  WS-TOTALS.
           05  WS-INTEREST-COUNT       PIC 9(8)       VALUE 0.
           05  WS-INTEREST-TOTAL       PIC 9(11)V99   VALUE 0.
           05  WS-OVERDRAFT-COUNT      PIC 9(8)       VALUE 0.
           05  WS-OVERDRAFT-TOTAL      PIC 9(11)V99   VALUE 0.
           05  WS-FEE-COUNT            PIC 9(8)       VALUE 0.
           05  WS-FEE-TOTAL            PIC 9(11)V99   VALUE 0.
           05  WS-WAIVED-COUNT         PIC 9(8)       VALUE 0.

       01  WS-REPORT-COUNT-LINE.
           05  WS-RC-LABEL             PIC X(30).
           05  WS-RC-VALUE             PIC ZZ,ZZZ,ZZ9.
           05  FILLER                  PIC X(40)      VALUE SPACES.
       01  WS-REPORT-AMOUNT-LINE.
           05  WS-RA-LABEL             PIC X(30).
           05  WS-RA-VALUE             PIC Z,ZZZ,ZZZ,ZZZ,ZZ9.99-.
           05  FILLER                  PIC X(29)      VALUE SPACES.

       PROCEDURE DIVISION.
      * CHARGE EVERY ACCOUNT, THEN REPORT.
       0000-MAIN.
           PERFORM 1000-INITIALIZE
           PERFORM 2000-CHARGE-ACCOUNT
               UNTIL MASTER-AT-END AND ACCRUALS-AT-END
           PERFORM 8000-WRITE-REPORT
           PERFORM 9000-CLOSE-FILES
           MOVE 0 TO RETURN-CODE
           STOP RUN.

      * OPEN THE FILES AND READ THE FIRST ACCOUNT. THE TRANSACTION
      * RECORD STARTS BLANK: NO COUNTER ACCOUNT.
       1000-INITIALIZE.
           OPEN INPUT  CONTROL-FILE ACCT-IN ACCR-IN
                OUTPUT TRAN-OUT REPORT-OUT
           IF WS-FILE-STATUSES NOT = ALL '0'
               DISPLAY 'MONCHG: CANNOT OPEN FILES ' WS-FILE-STATUSES
               MOVE 16 TO RETURN-CODE
               STOP RUN
           END-IF
           READ CONTROL-FILE
               AT END
                   MOVE 'NO BANK CONTROL RECORD' TO WS-ABORT-MESSAGE
                   PERFORM 9900-ABORT
           END-READ
           MOVE BC-BUSINESS-DATE TO WS-BUSINESS-DATE
           MOVE SPACES TO TRAN-RECORD
           PERFORM 1200-READ-ACCOUNT.

      * READ THE NEXT ACCOUNT AND ITS ACCRUAL, WHICH MUST MATCH.
       1200-READ-ACCOUNT.
           READ ACCT-IN
               AT END
                   SET MASTER-AT-END TO TRUE
           END-READ
           READ ACCR-IN
               AT END
                   SET ACCRUALS-AT-END TO TRUE
           END-READ
           IF MASTER-AT-END AND ACCRUALS-AT-END
               CONTINUE
           ELSE
               IF MASTER-AT-END OR ACCRUALS-AT-END
                  OR AC-ACCT-NO NOT = AM-ACCT-NO
                   MOVE 'ACCRUALS OUT OF STEP WITH THE MASTER'
                       TO WS-ABORT-MESSAGE
                   PERFORM 9900-ABORT
               END-IF
           END-IF.

      * WRITE AN ACCOUNT'S INTEREST AND FEE, IF ITS MONTH CLOSED.
       2000-CHARGE-ACCOUNT.
           IF AC-DUE-MONTH NOT = 0
               MOVE AC-DUE-MONTH TO WS-CLOSED-MONTH
               MOVE WS-CM-YYMM TO WS-TI-YYMM
               PERFORM 2100-SETTLE-INTEREST
               PERFORM 2200-CHARGE-FEE
           END-IF
           PERFORM 1200-READ-ACCOUNT.

      * CREDIT THE INTEREST DUE, OR CHARGE IT IF THE ACCOUNT WAS
      * OVERDRAWN.
       2100-SETTLE-INTEREST.
           EVALUATE TRUE
               WHEN AC-DUE > 0
                   SET TD-INTEREST TO TRUE
                   MOVE AC-DUE TO TD-AMOUNT
                   MOVE 'CREDIT INTEREST' TO TD-DESCRIPTION
                   PERFORM 2900-WRITE-TRAN
                   ADD 1 TO WS-INTEREST-COUNT
                   ADD TD-AMOUNT TO WS-INTEREST-TOTAL
               WHEN AC-DUE < 0
                   SET TD-OVERDRAFT-INTEREST TO TRUE
                   COMPUTE TD-AMOUNT = 0 - AC-DUE
                   MOVE 'OVERDRAFT INTEREST' TO TD-DESCRIPTION
                   PERFORM 2900-WRITE-TRAN
                   ADD 1 TO WS-OVERDRAFT-COUNT
                   ADD TD-AMOUNT TO WS-OVERDRAFT-TOTAL
               WHEN OTHER
                   CONTINUE
           END-EVALUATE.

      * CHARGE THE MAINTENANCE FEE: OPEN CURRENT ACCOUNTS ONLY, AND
      * NOT IF THE MONTH CLOSED AT OR ABOVE THE WAIVER BALANCE.
       2200-CHARGE-FEE.
           EVALUATE AM-ACCT-TYPE ALSO AM-STATUS
               WHEN 'C' ALSO 'O'
                   IF AM-BALANCE < WS-FEE-WAIVER-BALANCE
                       SET TD-MAINTENANCE-FEE TO TRUE
                       MOVE WS-MAINTENANCE-FEE TO TD-AMOUNT
                       MOVE 'MAINTENANCE FEE' TO TD-DESCRIPTION
                       PERFORM 2900-WRITE-TRAN
                       ADD 1 TO WS-FEE-COUNT
                       ADD TD-AMOUNT TO WS-FEE-TOTAL
                   ELSE
                       ADD 1 TO WS-WAIVED-COUNT
                   END-IF
               WHEN OTHER
                   CONTINUE
           END-EVALUATE.

      * WRITE THE TRANSACTION FOR THE ACCOUNT, WITH THE NEXT ID.
       2900-WRITE-TRAN.
           SET TR-DETAIL TO TRUE
           ADD 1 TO WS-TI-SEQ
           MOVE WS-TRAN-ID TO TD-TRAN-ID
           MOVE AM-ACCT-NO TO TD-ACCT-NO
           WRITE TRAN-OUT-RECORD FROM TRAN-RECORD.

      * WRITE THE MONTH-END CHARGES REPORT.
       8000-WRITE-REPORT.
           MOVE SPACES TO REPORT-LINE
           STRING 'MONCHG   MONTH-END CHARGES       BUSINESS DATE '
                  WS-BUSINESS-DATE
               DELIMITED BY SIZE INTO REPORT-LINE
           WRITE REPORT-LINE
           MOVE SPACES TO REPORT-LINE
           WRITE REPORT-LINE
           IF WS-CLOSED-MONTH = 0
               MOVE 'NO MONTH CLOSED' TO REPORT-LINE
               WRITE REPORT-LINE
           ELSE
               MOVE SPACES TO REPORT-LINE
               STRING 'MONTH CLOSED                  ' WS-CLOSED-MONTH
                   DELIMITED BY SIZE INTO REPORT-LINE
               WRITE REPORT-LINE
               MOVE SPACES TO REPORT-LINE
               WRITE REPORT-LINE
               MOVE 'INTEREST CREDITS' TO WS-RC-LABEL
               MOVE WS-INTEREST-COUNT TO WS-RC-VALUE
               WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
               MOVE 'INTEREST CREDITED' TO WS-RA-LABEL
               MOVE WS-INTEREST-TOTAL TO WS-RA-VALUE
               WRITE REPORT-LINE FROM WS-REPORT-AMOUNT-LINE
               MOVE 'OVERDRAFT INTEREST CHARGES' TO WS-RC-LABEL
               MOVE WS-OVERDRAFT-COUNT TO WS-RC-VALUE
               WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
               MOVE 'OVERDRAFT INTEREST CHARGED' TO WS-RA-LABEL
               MOVE WS-OVERDRAFT-TOTAL TO WS-RA-VALUE
               WRITE REPORT-LINE FROM WS-REPORT-AMOUNT-LINE
               MOVE 'MAINTENANCE FEES' TO WS-RC-LABEL
               MOVE WS-FEE-COUNT TO WS-RC-VALUE
               WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
               MOVE 'MAINTENANCE FEES CHARGED' TO WS-RA-LABEL
               MOVE WS-FEE-TOTAL TO WS-RA-VALUE
               WRITE REPORT-LINE FROM WS-REPORT-AMOUNT-LINE
               MOVE 'MAINTENANCE FEES WAIVED' TO WS-RC-LABEL
               MOVE WS-WAIVED-COUNT TO WS-RC-VALUE
               WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           END-IF.

      * CLOSE ALL FILES.
       9000-CLOSE-FILES.
           CLOSE CONTROL-FILE ACCT-IN ACCR-IN TRAN-OUT REPORT-OUT.

      * STOP THE RUN WITH RETURN CODE 16, SAYING WHY.
       9900-ABORT.
           DISPLAY 'MONCHG: ' WS-ABORT-MESSAGE
           PERFORM 9000-CLOSE-FILES
           MOVE 16 TO RETURN-CODE
           STOP RUN.
