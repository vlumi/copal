       IDENTIFICATION DIVISION.
       PROGRAM-ID. ORIGINS.
       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT IN-FILE  ASSIGN TO INDD.
           SELECT OUT-FILE ASSIGN TO OUTDD.
       DATA DIVISION.
       FILE SECTION.
       FD  IN-FILE.
       01  IN-RECORD.
           05  IN-KIND                 PIC X.
           05  IN-BODY                 PIC X(9).
           05  IN-DETAIL REDEFINES IN-BODY.
               10  IN-ACCT             PIC 9(5).
               10  IN-AMOUNT           PIC 9(4).
       FD  OUT-FILE.
       01  OUT-RECORD                  PIC X(8).
       WORKING-STORAGE SECTION.
       01  WS-HOLD.
           05  WS-KIND                 PIC X.
           05  WS-BODY                 PIC X(9).
       01  WS-OUT.
           05  WS-OUT-ACCT             PIC 9(5).
           05  WS-OUT-TAG              PIC X(3)       VALUE 'ABC'.
       PROCEDURE DIVISION.
       MAIN-PARA.
           READ IN-FILE
           MOVE 5 TO IN-AMOUNT
           MOVE IN-RECORD TO WS-HOLD
           MOVE IN-ACCT TO WS-OUT-ACCT
           WRITE OUT-RECORD FROM WS-OUT
           STOP RUN.
