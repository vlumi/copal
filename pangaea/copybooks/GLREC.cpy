      *================================================================*
      * GLREC: GENERAL LEDGER SUMMARY RECORD, 60 BYTES                 *
      * ONE PER GENERAL LEDGER ACCOUNT, FOR A MONTH: THE DEBITS AND    *
      * CREDITS THE MONTH'S POSTINGS MADE TO IT.                       *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 2004-03-15 TL  CR-0390 ORIGINAL VERSION.                       *
      *================================================================*
       01  GL-SUMMARY-RECORD.
           05  GL-MONTH                PIC 9(6).
           05  GL-ACCOUNT              PIC 9(4).
           05  GL-ACCOUNT-NAME         PIC X(24).
           05  GL-DEBITS               PIC S9(13)V99  COMP-3.
           05  GL-CREDITS              PIC S9(13)V99  COMP-3.
           05  FILLER                  PIC X(10).
