      *================================================================*
      * AUDITREC: AUDIT TRAIL RECORD, 100 BYTES, ONE PER POSTING       *
      * THE AMOUNT IS SIGNED: CREDITS POSITIVE, DEBITS NEGATIVE.       *
      *================================================================*
       01  AUDIT-RECORD.
           05  AU-BUSINESS-DATE        PIC 9(8).
           05  AU-TRAN-ID              PIC 9(10).
           05  AU-ACCT-NO              PIC 9(10).
           05  AU-TRAN-CODE            PIC XX.
           05  AU-AMOUNT               PIC S9(9)V99.
           05  AU-BALANCE-BEFORE       PIC S9(11)V99.
           05  AU-BALANCE-AFTER        PIC S9(11)V99.
           05  FILLER                  PIC X(33).
