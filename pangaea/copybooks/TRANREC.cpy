      *================================================================*
      * TRANREC: DAILY TRANSACTION FILE RECORD, 80 BYTES               *
      * ONE HEADER, THE DETAILS, ONE TRAILER, TOLD APART BY THE        *
      * RECORD TYPE. THE BODY IS REDEFINED FOR EACH RECORD TYPE.       *
      * THE TRAILER CARRIES THE SENDER'S CONTROL TOTALS: THE DETAIL    *
      * COUNT, THE SUM OF AMOUNTS AND THE SUM OF ACCOUNT NUMBERS.      *
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
                   88  TD-KNOWN-CODE   VALUE 'DP' 'WD' 'FE'.
               10  TD-AMOUNT           PIC 9(9)V99.
               10  TD-DESCRIPTION      PIC X(20).
               10  FILLER              PIC X(26).
           05  TR-TRAILER-BODY REDEFINES TR-BODY.
               10  TT-DETAIL-COUNT     PIC 9(8).
               10  TT-AMOUNT-TOTAL     PIC 9(13)V99.
               10  TT-ACCT-HASH        PIC 9(15).
               10  FILLER              PIC X(41).
