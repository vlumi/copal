      *================================================================*
      * AUTHMSG: AUTHORIZATION REQUEST, 80 BYTES, AND RESPONSE, 90     *
      * A CARD TERMINAL OR ATM SENDS A REQUEST FOR A CARD: A PURCHASE, *
      * THE REVERSAL OF ONE BY ITS AUTHORIZATION CODE, A CASH          *
      * WITHDRAWAL OR A BALANCE ENQUIRY. THE RESPONSE GIVES THE        *
      * RESULT, AND WHEN APPROVED THE AUTHORIZATION CODE, THE AMOUNT   *
      * HELD AND THE FEE, AND THE ACCOUNT'S BALANCE AND AVAILABLE      *
      * BALANCE AFTERWARDS.                                            *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 2003-05-19 MV  CR-0342 ORIGINAL VERSION.                       *
      * 2007-02-05 MV  CR-0415 FEE IN THE RESPONSE.                    *
      *================================================================*
       01  AUTH-REQUEST.
           05  RQ-FUNCTION             PIC X(4).
               88  RQ-PURCHASE         VALUE 'PURC'.
               88  RQ-REVERSAL         VALUE 'REVS'.
               88  RQ-WITHDRAWAL       VALUE 'WDRL'.
               88  RQ-BALANCE-ENQUIRY  VALUE 'BALQ'.
           05  RQ-CARD-NO              PIC 9(16).
           05  RQ-AMOUNT               PIC 9(7)V99.
           05  RQ-AUTH-CODE            PIC 9(6).
           05  RQ-TERMINAL             PIC X(20).
           05  RQ-TIMESTAMP            PIC 9(14).
           05  FILLER                  PIC X(11).
       01  AUTH-RESPONSE.
           05  RS-RESULT               PIC XX.
               88  RS-APPROVED         VALUE '00'.
               88  RS-INVALID-REQUEST  VALUE '12'.
               88  RS-INVALID-AMOUNT   VALUE '13'.
               88  RS-NO-SUCH-CARD     VALUE '14'.
               88  RS-NO-SUCH-AUTH     VALUE '25'.
               88  RS-INSUFFICIENT-FUNDS VALUE '51'.
               88  RS-OVER-LIMIT       VALUE '61'.
               88  RS-CARD-BLOCKED     VALUE '62'.
               88  RS-ACCOUNT-NOT-OPEN VALUE '78'.
           05  FILLER                  PIC X.
           05  RS-AUTH-CODE            PIC 9(6).
           05  FILLER                  PIC X.
           05  RS-AMOUNT               PIC Z(6)9.99.
           05  FILLER                  PIC X.
           05  RS-FEE                  PIC ZZ9.99.
           05  FILLER                  PIC X.
           05  RS-BALANCE              PIC -9(11).99.
           05  FILLER                  PIC X.
           05  RS-AVAILABLE            PIC -9(11).99.
           05  FILLER                  PIC X.
           05  RS-MESSAGE              PIC X(30).
