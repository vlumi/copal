      *================================================================*
      * ATMWDL: AUTHORIZE AN ATM CASH WITHDRAWAL                       *
      *                                                                *
      * CALLED BY AUTHSRV WITH THE REQUEST, THE CARD AND ITS ACCOUNT.  *
      * APPROVES THE WITHDRAWAL IF THE CARD IS ACTIVE, THE ACCOUNT     *
      * OPEN, THE AMOUNT IN NOTES OF 20, THE DAY'S CASH WITHIN THE     *
      * CARD'S DAILY CASH LIMIT, AND THE AMOUNT AND ANY FEE WITHIN THE *
      * AVAILABLE BALANCE. THE FIRST FOUR WITHDRAWALS OF A MONTH ARE   *
      * FREE; EACH ONE AFTER THAT CARRIES THE ATM FEE. AN APPROVED     *
      * WITHDRAWAL AND ITS FEE ARE HELD ON THE ACCOUNT AND COUNTED ON  *
      * THE CARD, AND THE HOLD RECORD FILLED IN FOR AUTHSRV TO WRITE.  *
      * NO FILES.                                                      *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 2003-05-19 MV  CR-0342 ORIGINAL VERSION.                       *
      * 2007-02-05 MV  CR-0415 ATM FEE AFTER FOUR FREE WITHDRAWALS A   *
      *                MONTH.                                          *
      *================================================================*
       IDENTIFICATION DIVISION.
       PROGRAM-ID. ATMWDL.
       AUTHOR. M. VARGA, ONLINE SYSTEMS.
       INSTALLATION. FIRST BANK OF PANGAEA, DATA CENTER.
       DATE-WRITTEN. 2003-05-19.

       DATA DIVISION.
       WORKING-STORAGE SECTION.
       01  WS-NOTE                     PIC 99         VALUE 20.
       01  WS-FREE-WITHDRAWALS         PIC 9(3)       VALUE 4.
       01  WS-ATM-FEE                  PIC 9V99       VALUE 2.00.
       01  WS-AVAILABLE                PIC S9(11)V99.
       01  WS-DAY-CASH                 PIC 9(7)V99.
       01  WS-FEE                      PIC 9V99.
       01  WS-TO-HOLD                  PIC 9(7)V99.

       LINKAGE SECTION.
           COPY AUTHMSG.
           COPY BANKCTL.
           COPY ACCTONL.
           COPY CARDREC.
           COPY HOLDREC.

       PROCEDURE DIVISION USING AUTH-REQUEST AUTH-RESPONSE
               BANK-CONTROL-RECORD ACCT-ONLINE-RECORD CARD-RECORD
               HOLD-RECORD.
      * DECIDE, AND HOLD THE AMOUNT AND FEE IF APPROVED.
       0000-MAIN.
           PERFORM 1000-START-CARD-TOTALS
           ADD AO-BALANCE AO-OVERDRAFT-LIMIT GIVING WS-AVAILABLE
           SUBTRACT AO-HOLDS FROM WS-AVAILABLE
           COMPUTE WS-DAY-CASH = CR-DAY-CASH + RQ-AMOUNT
           IF CR-MONTH-WITHDRAWALS < WS-FREE-WITHDRAWALS
               MOVE 0 TO WS-FEE
           ELSE
               MOVE WS-ATM-FEE TO WS-FEE
           END-IF
           COMPUTE WS-TO-HOLD = RQ-AMOUNT + WS-FEE
           EVALUATE TRUE
               WHEN CR-BLOCKED
                   SET RS-CARD-BLOCKED TO TRUE
                   MOVE 'CARD BLOCKED' TO RS-MESSAGE
               WHEN NOT AO-OPEN
                   SET RS-ACCOUNT-NOT-OPEN TO TRUE
                   MOVE 'ACCOUNT NOT OPEN' TO RS-MESSAGE
               WHEN RQ-AMOUNT = 0
               WHEN FUNCTION MOD(RQ-AMOUNT, WS-NOTE) NOT = 0
                   SET RS-INVALID-AMOUNT TO TRUE
                   MOVE 'AMOUNT NOT IN NOTES OF 20' TO RS-MESSAGE
               WHEN WS-DAY-CASH > CR-CASH-LIMIT
                   SET RS-OVER-LIMIT TO TRUE
                   MOVE 'OVER DAILY CASH LIMIT' TO RS-MESSAGE
               WHEN WS-TO-HOLD > WS-AVAILABLE
                   SET RS-INSUFFICIENT-FUNDS TO TRUE
                   MOVE 'INSUFFICIENT FUNDS' TO RS-MESSAGE
               WHEN OTHER
                   PERFORM 2000-APPROVE
           END-EVALUATE
           GOBACK.

      * THE CARD'S TOTALS ARE FOR ONE DAY, AND ITS COUNT OF
      * WITHDRAWALS FOR ONE MONTH: START THEM OVER ON A NEW ONE.
       1000-START-CARD-TOTALS.
           IF CR-TOTALS-DATE NOT = BC-BUSINESS-DATE
               IF CR-TOTALS-DATE(1:6) NOT = BC-BUSINESS-DATE(1:6)
                   MOVE 0 TO CR-MONTH-WITHDRAWALS
               END-IF
               MOVE 0 TO CR-DAY-PURCHASES CR-DAY-CASH
               MOVE BC-BUSINESS-DATE TO CR-TOTALS-DATE
           END-IF.

      * HOLD THE AMOUNT AND THE FEE, AND COUNT THE WITHDRAWAL.
       2000-APPROVE.
           ADD WS-TO-HOLD TO AO-HOLDS
           MOVE WS-DAY-CASH TO CR-DAY-CASH
           ADD 1 TO CR-MONTH-WITHDRAWALS
           MOVE AO-ACCT-NO TO HD-ACCT-NO
           MOVE CR-CARD-NO TO HD-CARD-NO
           SET HD-WITHDRAWAL TO TRUE
           SET HD-OPEN TO TRUE
           MOVE RQ-AMOUNT TO HD-AMOUNT
           MOVE WS-FEE TO HD-FEE
           MOVE RQ-TERMINAL TO HD-TERMINAL
           MOVE RQ-TIMESTAMP TO HD-TIMESTAMP
           MOVE BC-BUSINESS-DATE TO HD-BUSINESS-DATE
           SET RS-APPROVED TO TRUE
           MOVE 'APPROVED' TO RS-MESSAGE
           MOVE RQ-AMOUNT TO RS-AMOUNT
           MOVE WS-FEE TO RS-FEE
           MOVE AO-BALANCE TO RS-BALANCE
           SUBTRACT WS-TO-HOLD FROM WS-AVAILABLE
           MOVE WS-AVAILABLE TO RS-AVAILABLE.
