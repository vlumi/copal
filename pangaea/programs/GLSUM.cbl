      *================================================================*
      * GLSUM: SUMMARIZE THE MONTH FOR THE GENERAL LEDGER              *
      *                                                                *
      * ON THE FIRST BUSINESS DAY OF A MONTH, AFTER THE DAILY BATCH,   *
      * POSTS EVERY POSTING OF THE MONTH THAT CLOSED TO THE GENERAL    *
      * LEDGER, AS STMTPRT PICKS THEM: ONE SIDE TO CUSTOMER DEPOSITS,  *
      * THE OTHER TO THE ACCOUNT THE CODE TABLE NAMES FOR ITS          *
      * TRANSACTION CODE. A CREDIT TO A CUSTOMER CREDITS DEPOSITS, A   *
      * DEBIT DEBITS THEM. WRITES EACH LEDGER ACCOUNT'S DEBITS AND     *
      * CREDITS FOR THE MONTH, AND REPORTS THEM.                       *
      *                                                                *
      * EVERY TRANSFER'S LEGS POST ON ONE DAY, SO TRANSFERS IN TRANSIT *
      * MUST NET TO ZERO OVER A MONTH. IF THEY DO NOT, A LEG WENT      *
      * MISSING, AND THE RUN ENDS WITH RETURN CODE 4.                  *
      *                                                                *
      * RETURN CODES: 0 SUMMARIZED OR NO MONTH CLOSED,                 *
      *               4 TRANSFERS IN TRANSIT DO NOT NET TO ZERO,       *
      *               16 AN UNKNOWN CODE OR ACCOUNT, OUT OF BALANCE,   *
      *               OR I/O.                                          *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 1494-11-10 LP  EVERY DEBIT HAS ITS CREDIT: THE SUMMA, VENICE.  *
      * 2004-03-15 TL  CR-0390 IN COBOL.                               *
      *================================================================*
       IDENTIFICATION DIVISION.
       PROGRAM-ID. GLSUM.
       AUTHOR. T. LINDQVIST, BATCH SYSTEMS.
       INSTALLATION. FIRST BANK OF PANGAEA, DATA CENTER.
       DATE-WRITTEN. 2004-03-15.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT ACCR-IN      ASSIGN TO ACCRIN
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-ACCR-IN-STATUS.
           SELECT AUDIT-IN     ASSIGN TO AUDIT
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-AUDIT-STATUS.
           SELECT GL-OUT       ASSIGN TO GLOUT
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-GL-STATUS.
           SELECT REPORT-OUT   ASSIGN TO GLRPT
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-REPORT-STATUS.

       DATA DIVISION.
       FILE SECTION.
       FD  ACCR-IN.
           COPY ACCRREC.
       FD  AUDIT-IN.
           COPY AUDITREC.
       FD  GL-OUT.
           COPY GLREC.
       FD  REPORT-OUT.
       01  REPORT-LINE                 PIC X(80).

       WORKING-STORAGE SECTION.
       01  WS-FILE-STATUSES.
           05  WS-ACCR-IN-STATUS       PIC XX.
           05  WS-AUDIT-STATUS         PIC XX.
           05  WS-GL-STATUS            PIC XX.
           05  WS-REPORT-STATUS        PIC XX.

       01  WS-ABORT-MESSAGE            PIC X(40).

      * THE MONTH THAT CLOSED AND THE ONE AFTER IT, AS IN STMTPRT.
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

      * WHICH LEDGER ACCOUNT EACH TRANSACTION CODE POSTS AGAINST.
           COPY CODETAB REPLACING ==:P:== BY ==GS==.

      * THE LEDGER ACCOUNTS, AND WHAT THE MONTH POSTED TO EACH.
       01  WS-LEDGER-COUNT             PIC 9          VALUE 7.
       01  WS-LEDGER-VALUES.
           05  FILLER                  PIC 9(4)       VALUE 1000.
           05  FILLER                  PIC X(24)      VALUE 'CASH'.
           05  FILLER                  PIC 9(4)       VALUE 2000.
           05  FILLER                  PIC X(24)
                                       VALUE 'CUSTOMER DEPOSITS'.
           05  FILLER                  PIC 9(4)       VALUE 2100.
           05  FILLER                  PIC X(24)
                                       VALUE 'CARD SETTLEMENT'.
           05  FILLER                  PIC 9(4)       VALUE 2900.
           05  FILLER                  PIC X(24)
                                       VALUE 'TRANSFERS IN TRANSIT'.
           05  FILLER                  PIC 9(4)       VALUE 4000.
           05  FILLER                  PIC X(24)
                                       VALUE 'INTEREST INCOME'.
           05  FILLER                  PIC 9(4)       VALUE 4100.
           05  FILLER                  PIC X(24)
                                       VALUE 'FEE INCOME'.
           05  FILLER                  PIC 9(4)       VALUE 5000.
           05  FILLER                  PIC X(24)
                                       VALUE 'INTEREST EXPENSE'.
       01  WS-LEDGER REDEFINES WS-LEDGER-VALUES.
           05  WS-LEDGER-ENTRY         OCCURS 7 TIMES
                                       INDEXED BY WS-LX.
               10  WS-LE-ACCOUNT       PIC 9(4).
               10  WS-LE-NAME          PIC X(24).
       01  WS-LEDGER-TOTALS.
           05  WS-LEDGER-TOTAL         OCCURS 7 TIMES
                                       INDEXED BY WS-TX.
               10  WS-LT-DEBITS        PIC S9(13)V99  COMP-3.
               10  WS-LT-CREDITS       PIC S9(13)V99  COMP-3.

       01  WS-DEPOSITS-ACCOUNT         PIC 9(4)       VALUE 2000.
       01  WS-TRANSIT-ACCOUNT          PIC 9(4)       VALUE 2900.
       01  WS-POSTING.
           05  WS-AMOUNT               PIC S9(9)V99.
           05  WS-DEBIT-ACCOUNT        PIC 9(4).
           05  WS-CREDIT-ACCOUNT       PIC 9(4).
           05  WS-WANTED-ACCOUNT       PIC 9(4).

       01  WS-TOTALS.
           05  WS-POSTING-COUNT        PIC 9(8)       VALUE 0.
           05  WS-DEBIT-TOTAL          PIC S9(13)V99  VALUE 0.
           05  WS-CREDIT-TOTAL         PIC S9(13)V99  VALUE 0.
           05  WS-DEPOSITS-CHANGE      PIC S9(13)V99  VALUE 0.
           05  WS-TRANSIT-NET          PIC S9(13)V99  VALUE 0.

       01  WS-HEADING-LINE.
           05  FILLER                  PIC X(30)      VALUE 'ACCOUNT'.
           05  FILLER                  PIC X(18)
                                       VALUE '            DEBITS'.
           05  FILLER                  PIC X(2)       VALUE SPACES.
           05  FILLER                  PIC X(18)
                                       VALUE '           CREDITS'.
           05  FILLER                  PIC X(12)      VALUE SPACES.
       01  WS-LEDGER-LINE.
           05  WS-LL-ACCOUNT           PIC X(4).
           05  FILLER                  PIC X(1)       VALUE SPACES.
           05  WS-LL-NAME              PIC X(25).
           05  WS-LL-DEBITS            PIC Z,ZZZ,ZZZ,ZZZ,ZZ9.99.
           05  FILLER                  PIC X(2)       VALUE SPACES.
           05  WS-LL-CREDITS           PIC Z,ZZZ,ZZZ,ZZZ,ZZ9.99.
           05  FILLER                  PIC X(10)      VALUE SPACES.
       01  WS-REPORT-COUNT-LINE.
           05  WS-RC-LABEL             PIC X(30).
           05  WS-RC-VALUE             PIC ZZ,ZZZ,ZZ9.
           05  FILLER                  PIC X(40)      VALUE SPACES.
       01  WS-REPORT-AMOUNT-LINE.
           05  WS-RA-LABEL             PIC X(30).
           05  WS-RA-VALUE             PIC Z,ZZZ,ZZZ,ZZZ,ZZ9.99-.
           05  FILLER                  PIC X(29)      VALUE SPACES.

       PROCEDURE DIVISION.
      * FIND THE MONTH THAT CLOSED, IF ONE DID, POST ITS POSTINGS TO
      * THE LEDGER AND WRITE THE SUMMARY.
       0000-MAIN.
           PERFORM 1000-INITIALIZE
           IF WS-CLOSED-MONTH-NUMBER NOT = 0
               PERFORM 2000-POST-MONTH
               PERFORM 3000-WRITE-LEDGER
           END-IF
           PERFORM 8000-WRITE-REPORT
           PERFORM 9000-CLOSE-FILES
           IF WS-TRANSIT-NET NOT = 0
               MOVE 4 TO RETURN-CODE
           ELSE
               MOVE 0 TO RETURN-CODE
           END-IF
           STOP RUN.

      * OPEN THE FILES, TAKE THE CLOSED MONTH FROM THE ACCRUALS AND
      * START THE LEDGER AT ZERO.
       1000-INITIALIZE.
           OPEN INPUT  ACCR-IN AUDIT-IN
                OUTPUT GL-OUT REPORT-OUT
           IF WS-FILE-STATUSES NOT = ALL '0'
               DISPLAY 'GLSUM: CANNOT OPEN FILES ' WS-FILE-STATUSES
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
           END-IF
           PERFORM VARYING WS-TX FROM 1 BY 1
                   UNTIL WS-TX > WS-LEDGER-COUNT
               MOVE 0 TO WS-LT-DEBITS(WS-TX) WS-LT-CREDITS(WS-TX)
           END-PERFORM.

      * POST THE POSTINGS OF THE CLOSED MONTH, AND THE INTEREST AND
      * FEES POSTED FOR IT ON THE FIRST BUSINESS DAY OF THE NEXT.
       2000-POST-MONTH.
           READ AUDIT-IN
           PERFORM UNTIL WS-AUDIT-STATUS NOT = '00'
               MOVE AU-BUSINESS-DATE(1:6) TO WS-POSTING-MONTH
               MOVE AU-TRAN-CODE TO WS-POSTING-CODE
               IF WS-POSTING-MONTH = WS-CLOSED-MONTH
                  OR (WS-POSTING-MONTH = WS-NEXT-MONTH
                      AND WS-MONTH-END-CODE)
                   PERFORM 2100-POST-TO-LEDGER
               END-IF
               READ AUDIT-IN
           END-PERFORM
           IF WS-AUDIT-STATUS NOT = '10'
               MOVE 'CANNOT READ THE AUDIT TRAILS' TO WS-ABORT-MESSAGE
               PERFORM 9900-ABORT
           END-IF.

      * POST ONE POSTING: DEPOSITS ON THE CUSTOMER'S SIDE, THE CODE'S
      * LEDGER ACCOUNT ON THE OTHER.
       2100-POST-TO-LEDGER.
           SET GS-CX TO 1
           SEARCH GS-CODE-ENTRY
               AT END
                   MOVE 'UNKNOWN TRANSACTION CODE' TO WS-ABORT-MESSAGE
                   PERFORM 9900-ABORT
               WHEN GS-CODE(GS-CX) = AU-TRAN-CODE
                   CONTINUE
           END-SEARCH
           IF AU-AMOUNT < 0
               COMPUTE WS-AMOUNT = 0 - AU-AMOUNT
               MOVE WS-DEPOSITS-ACCOUNT TO WS-DEBIT-ACCOUNT
               MOVE GS-CODE-GL-ACCOUNT(GS-CX) TO WS-CREDIT-ACCOUNT
           ELSE
               MOVE AU-AMOUNT TO WS-AMOUNT
               MOVE GS-CODE-GL-ACCOUNT(GS-CX) TO WS-DEBIT-ACCOUNT
               MOVE WS-DEPOSITS-ACCOUNT TO WS-CREDIT-ACCOUNT
           END-IF
           MOVE WS-DEBIT-ACCOUNT TO WS-WANTED-ACCOUNT
           PERFORM 2200-FIND-ACCOUNT
           ADD WS-AMOUNT TO WS-LT-DEBITS(WS-TX)
           MOVE WS-CREDIT-ACCOUNT TO WS-WANTED-ACCOUNT
           PERFORM 2200-FIND-ACCOUNT
           ADD WS-AMOUNT TO WS-LT-CREDITS(WS-TX)
           ADD 1 TO WS-POSTING-COUNT.

      * FIND THE WANTED LEDGER ACCOUNT, AND POINT AT ITS TOTALS.
       2200-FIND-ACCOUNT.
           SET WS-LX TO 1
           SEARCH WS-LEDGER-ENTRY
               AT END
                   MOVE 'UNKNOWN LEDGER ACCOUNT' TO WS-ABORT-MESSAGE
                   PERFORM 9900-ABORT
               WHEN WS-LE-ACCOUNT(WS-LX) = WS-WANTED-ACCOUNT
                   SET WS-TX TO WS-LX
           END-SEARCH.

      * WRITE A SUMMARY RECORD FOR EACH LEDGER ACCOUNT, AND CHECK THAT
      * THE LEDGER BALANCES.
       3000-WRITE-LEDGER.
           PERFORM VARYING WS-LX FROM 1 BY 1
                   UNTIL WS-LX > WS-LEDGER-COUNT
               SET WS-TX TO WS-LX
               MOVE SPACES TO GL-SUMMARY-RECORD
               MOVE WS-CLOSED-MONTH-NUMBER TO GL-MONTH
               MOVE WS-LE-ACCOUNT(WS-LX) TO GL-ACCOUNT
               MOVE WS-LE-NAME(WS-LX) TO GL-ACCOUNT-NAME
               MOVE WS-LT-DEBITS(WS-TX) TO GL-DEBITS
               MOVE WS-LT-CREDITS(WS-TX) TO GL-CREDITS
               WRITE GL-SUMMARY-RECORD
               ADD WS-LT-DEBITS(WS-TX) TO WS-DEBIT-TOTAL
               ADD WS-LT-CREDITS(WS-TX) TO WS-CREDIT-TOTAL
               EVALUATE WS-LE-ACCOUNT(WS-LX)
                   WHEN WS-DEPOSITS-ACCOUNT
                       COMPUTE WS-DEPOSITS-CHANGE =
                           WS-LT-CREDITS(WS-TX) - WS-LT-DEBITS(WS-TX)
                   WHEN WS-TRANSIT-ACCOUNT
                       COMPUTE WS-TRANSIT-NET =
                           WS-LT-DEBITS(WS-TX) - WS-LT-CREDITS(WS-TX)
               END-EVALUATE
           END-PERFORM
           IF WS-DEBIT-TOTAL NOT = WS-CREDIT-TOTAL
               MOVE 'LEDGER OUT OF BALANCE' TO WS-ABORT-MESSAGE
               PERFORM 9900-ABORT
           END-IF.

      * WRITE THE GENERAL LEDGER REPORT.
       8000-WRITE-REPORT.
           MOVE 'GLSUM    GENERAL LEDGER SUMMARY' TO REPORT-LINE
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
           MOVE 'POSTINGS' TO WS-RC-LABEL
           MOVE WS-POSTING-COUNT TO WS-RC-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-COUNT-LINE
           MOVE SPACES TO REPORT-LINE
           WRITE REPORT-LINE
           WRITE REPORT-LINE FROM WS-HEADING-LINE
           PERFORM VARYING WS-LX FROM 1 BY 1
                   UNTIL WS-LX > WS-LEDGER-COUNT
               SET WS-TX TO WS-LX
               MOVE WS-LE-ACCOUNT(WS-LX) TO WS-LL-ACCOUNT
               MOVE WS-LE-NAME(WS-LX) TO WS-LL-NAME
               MOVE WS-LT-DEBITS(WS-TX) TO WS-LL-DEBITS
               MOVE WS-LT-CREDITS(WS-TX) TO WS-LL-CREDITS
               WRITE REPORT-LINE FROM WS-LEDGER-LINE
           END-PERFORM
           MOVE SPACES TO WS-LL-ACCOUNT
           MOVE 'TOTAL' TO WS-LL-NAME
           MOVE WS-DEBIT-TOTAL TO WS-LL-DEBITS
           MOVE WS-CREDIT-TOTAL TO WS-LL-CREDITS
           WRITE REPORT-LINE FROM WS-LEDGER-LINE
           MOVE SPACES TO REPORT-LINE
           WRITE REPORT-LINE
           MOVE 'CUSTOMER DEPOSITS CHANGE' TO WS-RA-LABEL
           MOVE WS-DEPOSITS-CHANGE TO WS-RA-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-AMOUNT-LINE
           MOVE 'TRANSFERS IN TRANSIT NET' TO WS-RA-LABEL
           MOVE WS-TRANSIT-NET TO WS-RA-VALUE
           WRITE REPORT-LINE FROM WS-REPORT-AMOUNT-LINE
           MOVE SPACES TO REPORT-LINE
           WRITE REPORT-LINE
           IF WS-TRANSIT-NET = 0
               MOVE 'LEDGER IN BALANCE' TO REPORT-LINE
           ELSE
               MOVE 'TRANSFERS IN TRANSIT DO NOT NET TO ZERO'
                   TO REPORT-LINE
           END-IF
           WRITE REPORT-LINE.

      * CLOSE ALL FILES.
       9000-CLOSE-FILES.
           CLOSE ACCR-IN AUDIT-IN GL-OUT REPORT-OUT.

      * STOP THE RUN WITH RETURN CODE 16, SAYING WHY.
       9900-ABORT.
           DISPLAY 'GLSUM: ' WS-ABORT-MESSAGE
           PERFORM 9000-CLOSE-FILES
           MOVE 16 TO RETURN-CODE
           STOP RUN.
