      *================================================================*
      * MAINTREC: ACCOUNT MAINTENANCE REQUEST, 80 BYTES                *
      * THE BRANCHES' REQUESTS TO OPEN AN ACCOUNT, CLOSE ONE, FREEZE   *
      * OR UNFREEZE ONE, OR CHANGE ITS OVERDRAFT LIMIT, IN ACCOUNT     *
      * NUMBER AND REQUEST ID ORDER. AN OPENING CARRIES THE NEW        *
      * ACCOUNT'S CUSTOMER, NAME AND TYPE; A LIMIT CHANGE THE NEW      *
      * LIMIT. THE EFFECTIVE DATE IS THE DAY THE BRANCH TOOK THE       *
      * REQUEST, YYMMDD.                                               *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 1988-05-09 JA  CR-0031 ORIGINAL VERSION.                       *
      * 1991-10-07 PK  CR-0129 EFFECTIVE DATE OF THE REQUEST.          *
      * 1999-01-11 SL  Y2K: LEFT AS IS, WINDOWED ON READ.              *
      *================================================================*
       01  MAINT-REQUEST.
           05  MR-REQUEST-ID           PIC 9(10).
           05  MR-REQUEST-TYPE         PIC X.
               88  MR-OPEN             VALUE 'O'.
               88  MR-CLOSE            VALUE 'C'.
               88  MR-FREEZE           VALUE 'F'.
               88  MR-UNFREEZE         VALUE 'U'.
               88  MR-LIMIT            VALUE 'L'.
           05  MR-ACCT-NO              PIC 9(10).
           05  MR-CUST-NO              PIC 9(8).
           05  MR-SHORT-NAME           PIC X(20).
           05  MR-ACCT-TYPE            PIC X.
           05  MR-OVERDRAFT-LIMIT      PIC 9(7)V99.
           05  MR-EFFECTIVE-DATE       PIC 9(6).
           05  FILLER                  PIC X(15).
