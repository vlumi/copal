      *================================================================*
      * ACCRREC: INTEREST ACCRUAL RECORD, 40 BYTES                     *
      * ONE PER ACCOUNT, IN ACCOUNT NUMBER ORDER LIKE THE MASTER.      *
      * INTEREST IS ACCRUED FOR EVERY CALENDAR DAY BEFORE THE THROUGH  *
      * DATE, TO SIX DECIMALS, AND IS NEGATIVE FOR AN OVERDRAWN        *
      * BALANCE. WHEN A MONTH CLOSES, ITS INTEREST BECOMES DUE,        *
      * ROUNDED TO THE CENT; THE FRACTION STAYS ACCRUED.               *
      * THE DUE MONTH IS ZERO ON THE DAYS NO MONTH CLOSES.             *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 2001-09-03 TL  CR-0298 ORIGINAL VERSION.                       *
      *================================================================*
       01  ACCRUAL-RECORD.
           05  AC-ACCT-NO              PIC 9(10).
           05  AC-THROUGH-DATE         PIC 9(8).
           05  AC-ACCRUED              PIC S9(7)V9(6) COMP-3.
           05  AC-DUE                  PIC S9(7)V99   COMP-3.
           05  AC-DUE-MONTH            PIC 9(6).
           05  FILLER                  PIC X(4).
