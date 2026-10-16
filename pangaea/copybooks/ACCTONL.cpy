      *================================================================*
      * ACCTONL: ONLINE ACCOUNT RECORD, 100 BYTES                      *
      * INDEXED BY ACCOUNT NUMBER. LOADED FROM THE MASTER EACH NIGHT   *
      * BY ACCTLOAD; DURING THE DAY THE AUTHORIZATIONS ADD WHAT THEY   *
      * APPROVE TO THE HOLDS, AND THE AVAILABLE BALANCE IS THE         *
      * BALANCE PLUS THE OVERDRAFT LIMIT LESS THE HOLDS.               *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 2003-05-19 MV  CR-0342 ORIGINAL VERSION.                       *
      *================================================================*
       01  ACCT-ONLINE-RECORD.
           05  AO-ACCT-NO              PIC 9(10).
           05  AO-CUST-NO              PIC 9(8).
           05  AO-SHORT-NAME           PIC X(20).
           05  AO-ACCT-TYPE            PIC X.
               88  AO-CURRENT          VALUE 'C'.
               88  AO-SAVINGS          VALUE 'S'.
           05  AO-STATUS               PIC X.
               88  AO-OPEN             VALUE 'O'.
               88  AO-FROZEN           VALUE 'F'.
      * THE BALANCE AT THE END OF THE BUSINESS DAY IT WAS LOADED ON.
           05  AO-BALANCE              PIC S9(11)V99.
           05  AO-OVERDRAFT-LIMIT      PIC 9(7)V99.
           05  AO-HOLDS                PIC S9(11)V99.
           05  AO-AS-OF-DATE           PIC 9(8).
           05  FILLER                  PIC X(17).
