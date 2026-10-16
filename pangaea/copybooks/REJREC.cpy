      *================================================================*
      * REJREC: REJECTED TRANSACTION, 100 BYTES                        *
      * THE TRANSACTION AS RECEIVED, WITH THE STEP THAT REJECTED IT    *
      * AND WHY.                                                       *
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
           05  RJ-TRAN                 PIC X(80).
           05  FILLER                  PIC X(8).
