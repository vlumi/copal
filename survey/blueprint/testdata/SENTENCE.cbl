       IDENTIFICATION DIVISION.
       PROGRAM-ID. SENTENCE.
       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT IN-FILE ASSIGN TO INFILE
               ORGANIZATION IS SEQUENTIAL.
       DATA DIVISION.
       FILE SECTION.
       FD  IN-FILE.
       01  IN-RECORD.
           05  IN-KIND                 PIC X.
               88  IN-OPEN             VALUE 'O' 'R'.
           05  IN-AMOUNT               PIC 9(5).
       WORKING-STORAGE SECTION.
       01  WS-TOTAL                    PIC 9(7).
       01  WS-LIMIT                    PIC 9(7).
       PROCEDURE DIVISION.
       MAIN-PARA.
           OPEN INPUT IN-FILE
           READ IN-FILE
               AT END
                   DISPLAY 'EMPTY'
                   STOP RUN
           END-READ
           IF IN-OPEN AND (IN-AMOUNT > WS-LIMIT OR WS-TOTAL = 0)
               ADD IN-AMOUNT TO WS-TOTAL
           ELSE
               MOVE 0 TO WS-TOTAL WS-LIMIT
               PERFORM CHECK-PARA UNTIL WS-TOTAL NOT < 10
           END-IF
           STOP RUN.
       CHECK-PARA.
           COMPUTE WS-TOTAL = WS-TOTAL + 1.
