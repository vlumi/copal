      *================================================================*
      * TRANREC: DAILY TRANSACTION FILE RECORD, 80 BYTES               *
      * ONE HEADER, THE DETAILS, ONE TRAILER, TOLD APART BY THE        *
      * RECORD TYPE. THE BODY IS REDEFINED FOR EACH RECORD TYPE.       *
      * THE TRAILER CARRIES THE SENDER'S CONTROL TOTALS: THE DETAIL    *
      * COUNT, THE SUM OF AMOUNTS AND THE SUM OF ACCOUNT NUMBERS.      *
      *                                                                *
      * A TRANSFER (TR) NAMES THE ACCOUNT IT GOES TO. POSTING SPLITS   *
      * IT IN TWO: THE DEBIT ON THE SENDING ACCOUNT, THEN A TRANSFER   *
      * CREDIT (TC) ON THE RECEIVING ONE. A CREDIT THAT CANNOT POST    *
      * GOES BACK AS A RETURNED TRANSFER (RT). TC AND RT ARE MADE BY   *
      * POSTING, NEVER RECEIVED IN THE DAILY FILE.                     *
      *                                                                *
      * WHEN A MONTH CLOSES, THE BANK CREDITS ITS INTEREST (IN),       *
      * CHARGES OVERDRAFT INTEREST (OD) AND THE MAINTENANCE FEE (MF).  *
      * THESE TOO ARE NEVER RECEIVED IN THE DAILY FILE.                *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 1987-03-12 JA  ORIGINAL VERSION.                               *
      * 1991-06-03 PK  CR-0117 ACCOUNT NUMBER HASH IN THE TRAILER.     *
      * 1996-04-15 PK  CR-0233 TRANSFERS: COUNTER ACCOUNT, FROM THE    *
      *                FILLER.                                         *
      * 2001-09-03 TL  CR-0298 MONTH-END INTEREST AND FEE CODES.       *
      *================================================================*
       01  TRAN-RECORD.
           05  TR-REC-TYPE             PIC X.
               88  TR-HEADER           VALUE 'H'.
               88  TR-DETAIL           VALUE 'D'.
               88  TR-TRAILER          VALUE 'T'.
           05  TR-BODY                 PIC X(79).
           05  TR-HEADER-BODY REDEFINES TR-BODY.
               10  TH-BUSINESS-DATE    PIC 9(8).
               10  TH-SOURCE           PIC X(8).
               10  TH-FILE-SEQ         PIC 9(6).
               10  FILLER              PIC X(57).
           05  TR-DETAIL-BODY REDEFINES TR-BODY.
               10  TD-TRAN-ID          PIC 9(10).
               10  TD-ACCT-NO          PIC 9(10).
               10  TD-TRAN-CODE        PIC XX.
                   88  TD-DEPOSIT      VALUE 'DP'.
                   88  TD-WITHDRAWAL   VALUE 'WD'.
                   88  TD-FEE          VALUE 'FE'.
                   88  TD-TRANSFER     VALUE 'TR'.
                   88  TD-TRANSFER-CREDIT VALUE 'TC'.
                   88  TD-TRANSFER-RETURN VALUE 'RT'.
                   88  TD-INTEREST     VALUE 'IN'.
                   88  TD-OVERDRAFT-INTEREST VALUE 'OD'.
                   88  TD-MAINTENANCE-FEE VALUE 'MF'.
                   88  TD-KNOWN-CODE   VALUE 'DP' 'WD' 'FE' 'TR'.
               10  TD-AMOUNT           PIC 9(9)V99.
               10  TD-DESCRIPTION      PIC X(20).
               10  TD-COUNTER-ACCT-NO  PIC 9(10).
               10  FILLER              PIC X(16).
           05  TR-TRAILER-BODY REDEFINES TR-BODY.
               10  TT-DETAIL-COUNT     PIC 9(8).
               10  TT-AMOUNT-TOTAL     PIC 9(13)V99.
               10  TT-ACCT-HASH        PIC 9(15).
               10  FILLER              PIC X(41).
