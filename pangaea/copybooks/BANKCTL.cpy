      *================================================================*
      * BANKCTL: BANK CONTROL RECORD, ONE PER FILE                     *
      * THE BUSINESS DATE THE BATCH RUNS FOR, NEVER THE CLOCK.         *
      *================================================================*
       01  BANK-CONTROL-RECORD.
           05  BC-BUSINESS-DATE        PIC 9(8).
           05  FILLER                  PIC X(72).
