      *================================================================*
      * RATEREC: INTEREST RATE TABLE RECORD, 20 BYTES                  *
      * ONE RECORD PER RATE: CREDIT RATES IN TIERS BY BALANCE, ONE     *
      * DEBIT RATE FOR BALANCES BELOW ZERO, PER ACCOUNT TYPE. A        *
      * BALANCE EARNS THE RATE OF THE HIGHEST TIER IT REACHES, ON THE  *
      * WHOLE BALANCE. RATES ARE ANNUAL FRACTIONS: 0.015 IS 1.5%.      *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * -1754      HAM ORIGINAL VERSION, ON CLAY, IN BABYLON: AT MOST  *
      *                20% A YEAR ON SILVER, 33 1/3% ON GRAIN.         *
      * 2001-09-03 TL  CR-0298 IN COBOL, WITH TIERS. GRAIN DROPPED.    *
      *================================================================*
       01  RATE-RECORD.
           05  RT-ACCT-TYPE            PIC X.
           05  RT-KIND                 PIC X.
               88  RT-CREDIT           VALUE 'C'.
               88  RT-DEBIT            VALUE 'D'.
           05  RT-FLOOR                PIC S9(11)V99  COMP-3.
           05  RT-ANNUAL-RATE          PIC 9V9(6)     COMP-3.
           05  FILLER                  PIC X(7).
