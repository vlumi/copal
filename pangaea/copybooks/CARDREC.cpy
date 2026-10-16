      *================================================================*
      * CARDREC: CARD RECORD, 80 BYTES                                 *
      * THE CARD MASTER, IN CARD NUMBER ORDER, AND THE ONLINE CARD     *
      * FILE LOADED FROM IT EACH NIGHT, INDEXED BY CARD NUMBER. A CARD *
      * IS VALID TO THE END OF ITS EXPIRY MONTH. THE AUTHORIZATIONS    *
      * KEEP THE TOTALS OF THE DAY OF CR-TOTALS-DATE, AND THE ATM      *
      * WITHDRAWALS OF ITS MONTH, AND START THEM OVER ON A NEW DAY.    *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 2003-05-19 MV  CR-0342 ORIGINAL VERSION.                       *
      *================================================================*
       01  CARD-RECORD.
           05  CR-CARD-NO              PIC 9(16).
           05  CR-ACCT-NO              PIC 9(10).
           05  CR-CUST-NO              PIC 9(8).
           05  CR-STATUS               PIC X.
               88  CR-ACTIVE           VALUE 'A'.
               88  CR-BLOCKED          VALUE 'B'.
           05  CR-EXPIRY-MONTH         PIC 9(6).
      * DAILY LIMITS.
           05  CR-PURCHASE-LIMIT       PIC 9(5)V99.
           05  CR-CASH-LIMIT           PIC 9(5)V99.
           05  CR-TOTALS-DATE          PIC 9(8).
           05  CR-DAY-PURCHASES        PIC 9(5)V99.
           05  CR-DAY-CASH             PIC 9(5)V99.
           05  CR-MONTH-WITHDRAWALS    PIC 9(3).
