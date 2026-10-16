      *================================================================*
      * HOLDREC: HOLD RECORD, 100 BYTES                                *
      * AN AMOUNT AUTHORIZED ONLINE AND NOT YET POSTED: A CARD         *
      * PURCHASE OR AN ATM WITHDRAWAL AND ITS FEE. INDEXED BY ACCOUNT  *
      * AND HOLD NUMBER, THE AUTHORIZATION CODE, WHICH NUMBERS THE     *
      * DAY'S HOLDS FROM 1. A REVERSED PURCHASE STAYS, MARKED.         *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 2003-05-19 MV  CR-0342 ORIGINAL VERSION.                       *
      * 2007-02-05 MV  CR-0415 ATM FEE.                                *
      *================================================================*
       01  HOLD-RECORD.
           05  HD-KEY.
               10  HD-ACCT-NO          PIC 9(10).
               10  HD-HOLD-NO          PIC 9(6).
           05  HD-CARD-NO              PIC 9(16).
           05  HD-KIND                 PIC X.
               88  HD-PURCHASE         VALUE 'P'.
               88  HD-WITHDRAWAL       VALUE 'W'.
           05  HD-STATUS               PIC X.
               88  HD-OPEN             VALUE 'O'.
               88  HD-REVERSED         VALUE 'R'.
           05  HD-AMOUNT               PIC 9(7)V99.
           05  HD-FEE                  PIC 9(3)V99.
           05  HD-TERMINAL             PIC X(20).
           05  HD-TIMESTAMP            PIC 9(14).
           05  HD-BUSINESS-DATE        PIC 9(8).
           05  FILLER                  PIC X(10).
