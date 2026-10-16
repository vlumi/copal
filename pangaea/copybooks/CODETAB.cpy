      *================================================================*
      * CODETAB: TRANSACTION CODE TABLE                                *
      * FOR EACH TRANSACTION CODE, WHAT A STATEMENT CALLS IT AND THE   *
      * GENERAL LEDGER ACCOUNT ITS POSTINGS MEET ON THE OTHER SIDE OF  *
      * CUSTOMER DEPOSITS. COPIED WITH THE PROGRAM'S OWN PREFIX FOR    *
      * THE :P: TAGS: COPY CODETAB REPLACING ==:P:== BY ==ST==.        *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 2004-03-15 TL  CR-0390 ORIGINAL VERSION, FROM STMTPRT, FOR THE *
      *                GENERAL LEDGER SUMMARY.                         *
      *================================================================*
       01  :P:-CODE-VALUES.
           05  FILLER                  PIC XX         VALUE 'DP'.
           05  FILLER                  PIC X(20)
                                       VALUE 'DEPOSIT'.
           05  FILLER                  PIC 9(4)       VALUE 1000.
           05  FILLER                  PIC XX         VALUE 'WD'.
           05  FILLER                  PIC X(20)
                                       VALUE 'WITHDRAWAL'.
           05  FILLER                  PIC 9(4)       VALUE 1000.
           05  FILLER                  PIC XX         VALUE 'FE'.
           05  FILLER                  PIC X(20)
                                       VALUE 'FEE'.
           05  FILLER                  PIC 9(4)       VALUE 4100.
           05  FILLER                  PIC XX         VALUE 'TR'.
           05  FILLER                  PIC X(20)
                                       VALUE 'TRANSFER OUT'.
           05  FILLER                  PIC 9(4)       VALUE 2900.
           05  FILLER                  PIC XX         VALUE 'TC'.
           05  FILLER                  PIC X(20)
                                       VALUE 'TRANSFER IN'.
           05  FILLER                  PIC 9(4)       VALUE 2900.
           05  FILLER                  PIC XX         VALUE 'RT'.
           05  FILLER                  PIC X(20)
                                       VALUE 'TRANSFER RETURNED'.
           05  FILLER                  PIC 9(4)       VALUE 2900.
           05  FILLER                  PIC XX         VALUE 'IN'.
           05  FILLER                  PIC X(20)
                                       VALUE 'CREDIT INTEREST'.
           05  FILLER                  PIC 9(4)       VALUE 5000.
           05  FILLER                  PIC XX         VALUE 'OD'.
           05  FILLER                  PIC X(20)
                                       VALUE 'OVERDRAFT INTEREST'.
           05  FILLER                  PIC 9(4)       VALUE 4000.
           05  FILLER                  PIC XX         VALUE 'MF'.
           05  FILLER                  PIC X(20)
                                       VALUE 'MAINTENANCE FEE'.
           05  FILLER                  PIC 9(4)       VALUE 4100.
           05  FILLER                  PIC XX         VALUE 'CP'.
           05  FILLER                  PIC X(20)
                                       VALUE 'CARD PURCHASE'.
           05  FILLER                  PIC 9(4)       VALUE 2100.
           05  FILLER                  PIC XX         VALUE 'AW'.
           05  FILLER                  PIC X(20)
                                       VALUE 'ATM WITHDRAWAL'.
           05  FILLER                  PIC 9(4)       VALUE 1000.
           05  FILLER                  PIC XX         VALUE 'AF'.
           05  FILLER                  PIC X(20)
                                       VALUE 'ATM FEE'.
           05  FILLER                  PIC 9(4)       VALUE 4100.
       01  :P:-CODE-TABLE REDEFINES :P:-CODE-VALUES.
           05  :P:-CODE-ENTRY          OCCURS 12 TIMES
                                       INDEXED BY :P:-CX.
               10  :P:-CODE            PIC XX.
               10  :P:-CODE-TEXT       PIC X(20).
               10  :P:-CODE-GL-ACCOUNT PIC 9(4).
