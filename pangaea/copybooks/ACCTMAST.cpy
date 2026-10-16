      *================================================================*
      * ACCTMAST: ACCOUNT MASTER RECORD, 100 BYTES                     *
      * IN ACCOUNT NUMBER ORDER. THE BALANCE MAY GO BELOW ZERO DOWN    *
      * TO THE OVERDRAFT LIMIT, AND BELOW THAT ONLY BY FEES.           *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 1987-03-12 JA  ORIGINAL VERSION.                               *
      * 1994-11-14 PK  CR-0201 OVERDRAFT LIMIT, FROM THE FILLER.       *
      * 1999-01-11 SL  Y2K: 8-DIGIT DATES.                             *
      *================================================================*
       01  ACCT-MASTER-RECORD.
           05  AM-ACCT-NO              PIC 9(10).
           05  AM-CUST-NO              PIC 9(8).
           05  AM-SHORT-NAME           PIC X(20).
           05  AM-ACCT-TYPE            PIC X.
               88  AM-CURRENT          VALUE 'C'.
               88  AM-SAVINGS          VALUE 'S'.
           05  AM-STATUS               PIC X.
               88  AM-OPEN             VALUE 'O'.
               88  AM-FROZEN           VALUE 'F'.
               88  AM-CLOSED           VALUE 'X'.
           05  AM-OPEN-DATE            PIC 9(8).
           05  AM-LAST-ACTIVITY-DATE   PIC 9(8).
           05  AM-BALANCE              PIC S9(11)V99.
           05  AM-OVERDRAFT-LIMIT      PIC 9(7)V99.
           05  FILLER                  PIC X(22).
