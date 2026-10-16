       IDENTIFICATION DIVISION.
       PROGRAM-ID. FLOW.
       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT IN-FILE  ASSIGN TO INDD.
           SELECT OUT-FILE ASSIGN TO OUTDD.
       DATA DIVISION.
       FILE SECTION.
       FD  IN-FILE.
       01  IN-RECORD.
           05  IN-AMOUNT               PIC 9(5).
           05  IN-CODE                 PIC X.
               88  IN-DEPOSIT          VALUE 'D'.
       FD  OUT-FILE.
       01  OUT-RECORD                  PIC X(6).
       WORKING-STORAGE SECTION.
       01  WS-TOTAL                    PIC 9(7).
       01  WS-COUNT                    PIC 9(5).
       01  WS-LINE.
           05  WS-LINE-TOTAL           PIC 9(7).
       PROCEDURE DIVISION.
       MAIN-PARA.
           READ IN-FILE
           IF IN-DEPOSIT
               ADD IN-AMOUNT TO WS-TOTAL
               ADD 1 TO WS-COUNT
           END-IF
           PERFORM WRITE-PARA
           STOP RUN.
       WRITE-PARA.
           COMPUTE WS-LINE-TOTAL = WS-TOTAL / WS-COUNT
           SET IN-DEPOSIT TO TRUE
           WRITE OUT-RECORD FROM WS-LINE.
