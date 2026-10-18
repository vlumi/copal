      *================================================================*
      * REJREC: REJECTED TRANSACTION, 110 BYTES                        *
      * THE TRANSACTION AS RECEIVED, WITH THE STEP THAT REJECTED IT    *
      * AND WHY. ACCOUNT MAINTENANCE REJECTS ITS REQUESTS IN THE SAME  *
      * FORM, THE REQUEST IN PLACE OF THE TRANSACTION. THE DATE IT     *
      * WAS REJECTED IS TEXT, YYYY-MM-DD, AS THE BRANCHES READ IT.     *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 1987-03-12 JA  ORIGINAL VERSION.                               *
      * 1988-05-09 JA  CR-0031 ACCOUNT MAINTENANCE REASONS.            *
      * 2009-03-02 MV  CR-0444 THE DATE OF THE REJECTION. RECORD       *
      *                WIDENED TO 110 BYTES.                           *
      *================================================================*
       01  REJECT-RECORD.
           05  RJ-PROGRAM              PIC X(8).
           05  RJ-REASON               PIC X(4).
               88  RJ-BAD-ACCOUNT      VALUE 'ACCT'.
               88  RJ-BAD-COUNTER      VALUE 'CTRA'.
               88  RJ-BAD-CODE         VALUE 'CODE'.
               88  RJ-BAD-AMOUNT       VALUE 'AMNT'.
               88  RJ-BAD-TYPE         VALUE 'TYPE'.
               88  RJ-NO-ACCOUNT       VALUE 'UNKN'.
               88  RJ-CLOSED           VALUE 'CLSD'.
               88  RJ-FROZEN           VALUE 'FRZN'.
               88  RJ-NO-FUNDS         VALUE 'FUND'.
               88  RJ-EXISTS           VALUE 'EXST'.
               88  RJ-HAS-BALANCE      VALUE 'BALN'.
               88  RJ-BAD-STATUS       VALUE 'STAT'.
           05  RJ-TRAN                 PIC X(80).
           05  FILLER                  PIC X(8).
           05  RJ-REJECT-DATE          PIC X(10).
