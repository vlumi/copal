      *================================================================*
      * INTACCR: ACCRUE INTEREST                                       *
      *                                                                *
      * ACCRUES EVERY ACCOUNT'S INTEREST FOR THE CALENDAR DAYS FROM    *
      * ITS LAST ACCRUAL UP TO THE BUSINESS DATE, ON THE OPENING       *
      * BALANCE: THE BALANCE THE LAST BUSINESS DAY CLOSED WITH, WHICH  *
      * STOOD OVER THE DAYS IN BETWEEN. THE RATE IS THE ACCOUNT        *
      * TYPE'S HIGHEST CREDIT TIER THE BALANCE REACHES, OR ITS DEBIT   *
      * RATE BELOW ZERO. A YEAR HAS 365 DAYS.                          *
      *                                                                *
      * WHEN THE DAYS CROSS INTO A NEW MONTH, THE DAYS OF THE OLD      *
      * MONTH ARE ACCRUED FIRST AND ITS INTEREST BECOMES DUE, ROUNDED  *
      * TO THE CENT, TO BE CREDITED OR CHARGED. THE ACCRUALS ARE READ  *
      * IN STEP WITH THE MASTER, ONE FOR EACH ACCOUNT.                 *
      *                                                                *
      * RETURN CODES: 0 ACCRUED,                                       *
      *               16 FILES OUT OF STEP, NO RATE, OR I/O.           *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 2001-09-03 TL  CR-0298 ORIGINAL VERSION.                       *
      *================================================================*
       IDENTIFICATION DIVISION.
       PROGRAM-ID. INTACCR.
       AUTHOR. T. LINDQVIST, BATCH SYSTEMS.
       INSTALLATION. FIRST BANK OF PANGAEA, DATA CENTER.
       DATE-WRITTEN. 2001-09-03.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT CONTROL-FILE ASSIGN TO BANKCTL
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-CONTROL-STATUS.
           SELECT RATE-FILE    ASSIGN TO RATES
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-RATE-STATUS.
           SELECT ACCT-IN      ASSIGN TO ACCTIN
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-ACCT-IN-STATUS.
           SELECT ACCR-IN      ASSIGN TO ACCRIN
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-ACCR-IN-STATUS.
           SELECT ACCR-OUT     ASSIGN TO ACCROUT
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-ACCR-OUT-STATUS.
           SELECT REPORT-OUT   ASSIGN TO ACCRRPT
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-STATUS.

       DATA DIVISION.
       FILE SECTION.
       FD  CONTROL-FILE.
           COPY BANKCTL.
       FD  RATE-FILE.
           COPY RATEREC.
       FD  ACCT-IN.
           COPY ACCTMAST.
       FD  ACCR-IN.
           COPY ACCRREC.
       FD  ACCR-OUT.
       01  ACCR-OUT-RECORD             PIC X(40).
       FD  REPORT-OUT.
       01  REPORT-LINE                 PIC X(80).

       WORKING-STORAGE SECTION.
       01  WS-FILE-STATUSES.
           05  WS-CONTROL-STATUS       PIC XX.
           05  WS-RATE-STATUS          PIC XX.
           05  WS-ACCT-IN-STATUS       PIC XX.
           05  WS-ACCR-IN-STATUS       PIC XX.
           05  WS-ACCR-OUT-STATUS      PIC XX.
           05  WS-REPORT-STATUS        PIC XX.

       01  WS-BUSINESS-DATE            PIC 9(8).
       01  WS-ABORT-MESSAGE            PIC X(40).
       01  WS-END-SWITCHES.
           05  WS-MASTER-END-SW        PIC X          VALUE 'N'.
               88  MASTER-AT-END       VALUE 'Y'.
           05  WS-ACCRUAL-END-SW       PIC X          VALUE 'N'.
               88  ACCRUALS-AT-END     VALUE 'Y'.

      * THE RATE TABLE, AS MANY RATES AS THE RATE FILE HOLDS.
       01  WS-RATE-COUNT               PIC 99         VALUE 0.
       01  WS-RATE-TABLE.
           05  WS-RATE OCCURS 1 TO 20 TIMES
                   DEPENDING ON WS-RATE-COUNT.
               10  WS-RT-ACCT-TYPE     PIC X.
               10  WS-RT-KIND          PIC X.
                   88  WS-RT-CREDIT    VALUE 'C'.
                   88  WS-RT-DEBIT     VALUE 'D'.
               10  WS-RT-FLOOR         PIC S9(11)V99  COMP-3.
               10  WS-RT-ANNUAL-RATE   PIC 9V9(6)     COMP-3.
       01  WS-RX                       PIC 99.

      * THE RATE FOR THE ACCOUNT IN HAND, AND THE TIER IT CAME FROM.
       01  WS-ACCOUNT-RATE             PIC 9V9(6)     COMP-3.
       01  WS-RATE-FOUND-SW            PIC X.
           88  RATE-FOUND              VALUE 'Y'.
       01  WS-BEST-FLOOR               PIC S9(11)V99  COMP-3.

      * DAYS COUNT FROM 1601-01-01, AS FUNCTION INTEGER-OF-DATE GIVES.
       01  WS-TODAY                    PIC 9(7).
       01  WS-FROM-DAY                 PIC 9(7).
       01  WS-MONTH-START-DAY          PIC 9(7).
       01  WS-DAYS                     PIC 9(3).
       01  WS-THROUGH-DATE.
           05  WS-TH-YEAR              PIC 9(4).
           05  WS-TH-MONTH             PIC 99.
           05  WS-TH-DAY               PIC 99.
      * THE FIRST DAY OF THE MONTH AFTER THE LAST ACCRUAL.
       01  WS-MONTH-START.
           05  WS-MS-YEAR              PIC 9(4).
           05  WS-MS-MONTH             PIC 99.
           05  WS-MS-DAY               PIC 99.
       01  WS-MONTH-START-DATE REDEFINES WS-MONTH-START
                                       PIC 9(8).

      * THE BALANCE TIMES THE RATE TIMES THE DAYS, EXACT, BEFORE IT IS
      * DIVIDED BY THE DAYS OF A YEAR AND ROUNDED.
       01  WS-INTEREST-PRODUCT         PIC S9(13)V9(8) COMP-3.
       01  WS-INCREMENT                PIC S9(7)V9(6) COMP-3.
       01  WS-DUE-MONTH.
           05  WS-DM-YEAR              PIC 9(4).
           05  WS-DM-MONTH             PIC 99.

       01  WS-TOTALS.
           05  WS-ACCOUNT-COUNT        PIC 9(8)       VALUE 0.
           05  WS-CREDIT-ACCRUED       PIC S9(9)V9(6) COMP-3 VALUE 0.
           05  WS-DEBIT-ACCRUED        PIC S9(9)V9(6) COMP-3 VALUE 0.
           05  WS-CREDIT-DUE           PIC S9(11)V99  COMP-3 VALUE 0.
           05  WS-DEBIT-DUE            PIC S9(11)V99  COMP-3 VALUE 0.
           05  WS-CLOSED-MONTH         PIC 9(6)       VALUE 0.

       01  WS-REPORT-COUNT-LINE.
           05  WS-RC-LABEL             PIC X(30).
           05  WS-RC-VALUE             PIC ZZ,ZZZ,ZZ9.
           05  FILLER                  PIC X(40)      VALUE SPACES.
       01  WS-REPORT-AMOUNT-LINE.
           05  WS-RA-LABEL             PIC X(30).
           05  WS-RA-VALUE             PIC ZZZ,ZZZ,ZZ9.999999-.
           05  FILLER                  PIC X(31)      VALUE SPACES.

       PROCEDURE DIVISION.
      * ACCRUE EVERY ACCOUNT, THEN REPORT.
       0000-MAIN.
           PERFORM 1000-INITIALIZE
           PERFORM 2000-ACCRUE-ACCOUNT
               UNTIL MASTER-AT-END AND ACCRUALS-AT-END
           PERFORM 8000-WRITE-REPORT
           PERFORM 9000-CLOSE-FILES
           MOVE 0 TO RETURN-CODE
           STOP RUN.

      * OPEN THE FILES, LOAD THE RATES AND READ THE FIRST ACCOUNT.
       1000-INITIALIZE.
           OPEN INPUT  CONTROL-FILE RATE-FILE ACCT-IN ACCR-IN
                OUTPUT ACCR-OUT REPORT-OUT
           IF WS-FILE-STATUSES NOT = ALL '0'
               DISPLAY 'INTACCR: CANNOT OPEN FILES ' WS-FILE-STATUSES
               MOVE 16 TO RETURN-CODE
               STOP RUN
           END-IF
           READ CONTROL-FILE
               AT END
                   MOVE 'NO BANK CONTROL RECORD' TO WS-ABORT-MESSAGE
                   PERFORM 9900-ABORT
           END-READ
           MOVE BC-BUSINESS-DATE TO WS-BUSINESS-DATE
           COMPUTE WS-TODAY = FUNCTION INTEGER-OF-DATE(WS-BUSINESS-DATE)
           PERFORM 1100-LOAD-RATES
           PERFORM 1200-READ-ACCOUNT.

      * READ THE RATE FILE INTO THE RATE TABLE.
       1100-LOAD-RATES.
           READ RATE-FILE
           PERFORM UNTIL WS-RATE-STATUS NOT = '00'
               IF WS-RATE-COUNT = 20
                   MOVE 'MORE THAN 20 RATES' TO WS-ABORT-MESSAGE
                   PERFORM 9900-ABORT
               END-IF
               ADD 1 TO WS-RATE-COUNT
               MOVE RT-ACCT-TYPE TO WS-RT-ACCT-TYPE(WS-RATE-COUNT)
               MOVE RT-KIND TO WS-RT-KIND(WS-RATE-COUNT)
               MOVE RT-FLOOR TO WS-RT-FLOOR(WS-RATE-COUNT)
               MOVE RT-ANNUAL-RATE TO WS-RT-ANNUAL-RATE(WS-RATE-COUNT)
               READ RATE-FILE
           END-PERFORM
           IF WS-RATE-STATUS NOT = '10'
               MOVE 'CANNOT READ THE RATES' TO WS-ABORT-MESSAGE
               PERFORM 9900-ABORT
           END-IF.

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

      * ACCRUE ONE ACCOUNT UP TO THE BUSINESS DATE. IF A MONTH CLOSED
      * ON THE WAY, ACCRUE ITS DAYS FIRST AND MAKE ITS INTEREST DUE.
       2000-ACCRUE-ACCOUNT.
           PERFORM 2100-FIND-RATE
           MOVE 0 TO AC-DUE AC-DUE-MONTH
           COMPUTE WS-FROM-DAY =
               FUNCTION INTEGER-OF-DATE(AC-THROUGH-DATE)
           PERFORM 2200-FIND-MONTH-START
           IF WS-MONTH-START-DAY NOT > WS-TODAY
               COMPUTE WS-DAYS = WS-MONTH-START-DAY - WS-FROM-DAY
               PERFORM 2300-ACCRUE-DAYS
               PERFORM 2400-CLOSE-MONTH
               MOVE WS-MONTH-START-DAY TO WS-FROM-DAY
           END-IF
           COMPUTE WS-DAYS = WS-TODAY - WS-FROM-DAY
           PERFORM 2300-ACCRUE-DAYS
           MOVE WS-BUSINESS-DATE TO AC-THROUGH-DATE
           WRITE ACCR-OUT-RECORD FROM ACCRUAL-RECORD
           ADD 1 TO WS-ACCOUNT-COUNT
           PERFORM 1200-READ-ACCOUNT.

      * FIND THE ACCOUNT'S RATE: BELOW ZERO THE DEBIT RATE OF ITS
      * TYPE, OTHERWISE THE CREDIT RATE OF THE HIGHEST TIER IT REACHES.
       2100-FIND-RATE.
           MOVE 'N' TO WS-RATE-FOUND-SW
           PERFORM VARYING WS-RX FROM 1 BY 1
                   UNTIL WS-RX > WS-RATE-COUNT
               IF WS-RT-ACCT-TYPE(WS-RX) = AM-ACCT-TYPE
                   IF AM-BALANCE < 0
                       IF WS-RT-DEBIT(WS-RX)
                           PERFORM 2110-TAKE-RATE
                       END-IF
                   ELSE
                       IF WS-RT-CREDIT(WS-RX)
                          AND WS-RT-FLOOR(WS-RX) NOT > AM-BALANCE
                          AND (NOT RATE-FOUND
                               OR WS-RT-FLOOR(WS-RX) > WS-BEST-FLOOR)
                           PERFORM 2110-TAKE-RATE
                       END-IF
                   END-IF
               END-IF
           END-PERFORM
           IF NOT RATE-FOUND
               MOVE 'NO RATE FOR AN ACCOUNT TYPE' TO WS-ABORT-MESSAGE
               PERFORM 9900-ABORT
           END-IF.

      * TAKE THE RATE IN HAND AS THE ACCOUNT'S, SO FAR.
       2110-TAKE-RATE.
           SET RATE-FOUND TO TRUE
           MOVE WS-RT-ANNUAL-RATE(WS-RX) TO WS-ACCOUNT-RATE
           MOVE WS-RT-FLOOR(WS-RX) TO WS-BEST-FLOOR.

      * THE FIRST DAY OF THE MONTH AFTER THE THROUGH DATE.
       2200-FIND-MONTH-START.
           MOVE AC-THROUGH-DATE TO WS-THROUGH-DATE
           MOVE WS-TH-YEAR TO WS-MS-YEAR
           IF WS-TH-MONTH = 12
               ADD 1 TO WS-MS-YEAR
               MOVE 1 TO WS-MS-MONTH
           ELSE
               COMPUTE WS-MS-MONTH = WS-TH-MONTH + 1
           END-IF
           MOVE 1 TO WS-MS-DAY
           COMPUTE WS-MONTH-START-DAY =
               FUNCTION INTEGER-OF-DATE(WS-MONTH-START-DATE).

      * ACCRUE THE OPENING BALANCE AT THE ACCOUNT'S RATE FOR THE DAYS.
       2300-ACCRUE-DAYS.
           MULTIPLY AM-BALANCE BY WS-ACCOUNT-RATE
               GIVING WS-INTEREST-PRODUCT
           MULTIPLY WS-DAYS BY WS-INTEREST-PRODUCT
           DIVIDE WS-INTEREST-PRODUCT BY 365
               GIVING WS-INCREMENT ROUNDED
           ADD WS-INCREMENT TO AC-ACCRUED
           IF WS-INCREMENT < 0
               ADD WS-INCREMENT TO WS-DEBIT-ACCRUED
           ELSE
               ADD WS-INCREMENT TO WS-CREDIT-ACCRUED
           END-IF.

      * MAKE THE CLOSED MONTH'S INTEREST DUE, TO THE CENT, AND KEEP
      * THE FRACTION ACCRUED.
       2400-CLOSE-MONTH.
           COMPUTE AC-DUE ROUNDED = AC-ACCRUED
           SUBTRACT AC-DUE FROM AC-ACCRUED
           MOVE WS-TH-YEAR TO WS-DM-YEAR
           MOVE WS-TH-MONTH TO WS-DM-MONTH
           MOVE WS-DUE-MONTH TO AC-DUE-MONTH
           MOVE AC-DUE-MONTH TO WS-CLOSED-MONTH
           IF AC-DUE < 0
               ADD AC-DUE TO WS-DEBIT-DUE
           ELSE
               ADD AC-DUE TO WS-CREDIT-DUE
           END-IF.

      * WRITE THE ACCRUAL REPORT.
       8000-WRITE-REPORT.
           MOVE SPACES TO REPORT-LINE
           STRING 'INTACCR  INTEREST ACCRUAL        BUSINESS DATE '
                  WS-BUSINESS-DATE
               DELIMITED BY SIZE INTO REPORT-LINE
           WRITE REPORT-LINE
           MOVE SPACES TO REPORT-LINE
           WRITE REPORT-LINE
           MOVE 'RATES' TO WS-RC-LABEL
           MOVE WS-RATE-COUNT TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           MOVE 'ACCOUNTS ACCRUED' TO WS-RC-LABEL
           MOVE WS-ACCOUNT-COUNT TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           MOVE 'CREDIT INTEREST ACCRUED' TO WS-RA-LABEL
           MOVE WS-CREDIT-ACCRUED TO WS-RA-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-AMOUNT-LINE
           MOVE 'DEBIT INTEREST ACCRUED' TO WS-RA-LABEL
           MOVE WS-DEBIT-ACCRUED TO WS-RA-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-AMOUNT-LINE
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
               MOVE 'CREDIT INTEREST DUE' TO WS-RA-LABEL
               MOVE WS-CREDIT-DUE TO WS-RA-VALUE
               WRITE REPORT-LINE FROM WS-REPORT-AMOUNT-LINE
               MOVE 'DEBIT INTEREST DUE' TO WS-RA-LABEL
               MOVE WS-DEBIT-DUE TO WS-RA-VALUE
               WRITE REPORT-LINE FROM WS-REPORT-AMOUNT-LINE
           END-IF.

      * CLOSE ALL FILES.
       9000-CLOSE-FILES.
           CLOSE CONTROL-FILE RATE-FILE ACCT-IN ACCR-IN
                 ACCR-OUT REPORT-OUT.

      * STOP THE RUN WITH RETURN CODE 16, SAYING WHY.
       9900-ABORT.
           DISPLAY 'INTACCR: ' WS-ABORT-MESSAGE
           PERFORM 9000-CLOSE-FILES
           MOVE 16 TO RETURN-CODE
           STOP RUN.
