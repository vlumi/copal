      *================================================================*
      * ATMBAL: ANSWER AN ATM BALANCE ENQUIRY                          *
      *                                                                *
      * CALLED BY AUTHSRV WITH THE REQUEST, THE CARD AND ITS ACCOUNT.  *
      * FOR AN ACTIVE CARD OF AN OPEN ACCOUNT, GIVES THE BALANCE AND   *
      * THE AVAILABLE BALANCE: THE BALANCE PLUS THE OVERDRAFT LIMIT    *
      * LESS WHAT IS HELD. CHANGES NOTHING. NO FILES.                  *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 2003-05-19 MV  CR-0342 ORIGINAL VERSION.                       *
      *================================================================*
       IDENTIFICATION DIVISION.
       PROGRAM-ID. ATMBAL.
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
      * GIVE THE BALANCES TO AN ACTIVE CARD OF AN OPEN ACCOUNT.
       0000-MAIN.
           EVALUATE TRUE
               WHEN CR-BLOCKED
                   SET RS-CARD-BLOCKED TO TRUE
                   MOVE 'CARD BLOCKED' TO RS-MESSAGE
               WHEN NOT AO-OPEN
                   SET RS-ACCOUNT-NOT-OPEN TO TRUE
                   MOVE 'ACCOUNT NOT OPEN' TO RS-MESSAGE
               WHEN OTHER
                   SET RS-APPROVED TO TRUE
                   MOVE 'BALANCE' TO RS-MESSAGE
                   MOVE AO-BALANCE TO RS-BALANCE
                   ADD AO-BALANCE AO-OVERDRAFT-LIMIT GIVING WS-AVAILABLE
                   SUBTRACT AO-HOLDS FROM WS-AVAILABLE
                   MOVE WS-AVAILABLE TO RS-AVAILABLE
           END-EVALUATE
           GOBACK.
