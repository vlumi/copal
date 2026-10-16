      *================================================================*
      * AUTHREV: REVERSE A CARD PURCHASE                               *
      *                                                                *
      * CALLED BY AUTHSRV WITH THE REQUEST, THE CARD, ITS ACCOUNT AND  *
      * THE HOLD OF THE AUTHORIZATION CODE TO REVERSE. A MERCHANT      *
      * REVERSES A PURCHASE OF THE SAME CARD THAT IS STILL HELD, EVEN  *
      * IF THE CARD HAS BEEN BLOCKED SINCE: THE HOLD IS RELEASED FROM  *
      * THE ACCOUNT, TAKEN OFF THE CARD'S PURCHASES OF THE DAY, AND    *
      * MARKED REVERSED FOR AUTHSRV TO REWRITE. CASH WITHDRAWALS ARE   *
      * NOT REVERSED. NO FILES.                                        *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 2003-05-19 MV  CR-0342 ORIGINAL VERSION.                       *
      *================================================================*
       IDENTIFICATION DIVISION.
       PROGRAM-ID. AUTHREV.
       AUTHOR. M. VARGA, ONLINE SYSTEMS.
       INSTALLATION. FIRST BANK OF PANGAEA, DATA CENTER.
       DATE-WRITTEN. 2003-05-19.

       DATA DIVISION.
       WORKING-STORAGE SECTION.
       01  WS-AVAILABLE                PIC S9(11)V99.

       LINKAGE SECTION.
           COPY AUTHMSG.
           COPY BANKCTL.
           COPY ACCTONL.
           COPY CARDREC.
           COPY HOLDREC.

       PROCEDURE DIVISION USING AUTH-REQUEST AUTH-RESPONSE
               BANK-CONTROL-RECORD ACCT-ONLINE-RECORD CARD-RECORD
               HOLD-RECORD.
      * DECIDE, AND RELEASE THE HOLD IF THE REVERSAL IS APPROVED.
       0000-MAIN.
           EVALUATE TRUE
               WHEN HD-CARD-NO NOT = CR-CARD-NO
                   SET RS-NO-SUCH-AUTH TO TRUE
                   MOVE 'AUTHORIZATION OF ANOTHER CARD' TO RS-MESSAGE
               WHEN NOT HD-PURCHASE
                   SET RS-INVALID-REQUEST TO TRUE
                   MOVE 'ONLY PURCHASES ARE REVERSED' TO RS-MESSAGE
               WHEN HD-REVERSED
                   SET RS-NO-SUCH-AUTH TO TRUE
                   MOVE 'ALREADY REVERSED' TO RS-MESSAGE
               WHEN OTHER
                   PERFORM 2000-APPROVE
           END-EVALUATE
           GOBACK.

      * RELEASE THE HOLD. THE CARD'S TOTALS COUNT IT ONLY IF THEY ARE
      * STILL FOR THE DAY IT WAS HELD ON.
       2000-APPROVE.
           SUBTRACT HD-AMOUNT FROM AO-HOLDS
           IF CR-TOTALS-DATE = HD-BUSINESS-DATE
               SUBTRACT HD-AMOUNT FROM CR-DAY-PURCHASES
           END-IF
           SET HD-REVERSED TO TRUE
           SET RS-APPROVED TO TRUE
           MOVE 'REVERSED' TO RS-MESSAGE
           MOVE HD-HOLD-NO TO RS-AUTH-CODE
           MOVE HD-AMOUNT TO RS-AMOUNT
           MOVE 0 TO RS-FEE
           MOVE AO-BALANCE TO RS-BALANCE
           ADD AO-BALANCE AO-OVERDRAFT-LIMIT GIVING WS-AVAILABLE
           SUBTRACT AO-HOLDS FROM WS-AVAILABLE
           MOVE WS-AVAILABLE TO RS-AVAILABLE.
