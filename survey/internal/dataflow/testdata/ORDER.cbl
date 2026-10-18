       IDENTIFICATION DIVISION.
       PROGRAM-ID. ORDER.
       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT IN-FILE ASSIGN TO INFILE.
           SELECT OUT-FILE ASSIGN TO OUTFILE.
       DATA DIVISION.
       FILE SECTION.
       FD  IN-FILE.
       01  IN-RECORD.
           05  IN-AMOUNT               PIC 9(5).
           05  IN-CODE                 PIC X.
       FD  OUT-FILE.
       01  OUT-RECORD                  PIC X(20).
       WORKING-STORAGE SECTION.
       01  WS-LINE.
           05  WS-AMOUNT               PIC 9(5).
           05  WS-CODE                 PIC X.
           05  WS-COUNT                PIC 9(3)       VALUE 0.
           05  WS-NOTE                 PIC X(11)      VALUE 'NONE'.
       01  WS-TOTAL                    PIC 9(7)       VALUE 0.
       01  WS-EOF                      PIC X          VALUE 'N'.
       PROCEDURE DIVISION.
       MAIN-PARA.
           OPEN INPUT IN-FILE OUTPUT OUT-FILE
           MOVE SPACES TO WS-CODE
           MOVE 99999 TO WS-AMOUNT
           PERFORM READ-PARA
           PERFORM UNTIL WS-EOF = 'Y'
               MOVE IN-AMOUNT TO WS-AMOUNT
               IF IN-CODE = 'A'
                   MOVE 'A' TO WS-CODE
               ELSE
                   MOVE IN-CODE TO WS-CODE
               END-IF
               IF IN-CODE = 'N'
                   MOVE 'NOTED' TO WS-NOTE
               END-IF
               ADD 1 TO WS-COUNT
               ADD IN-AMOUNT TO WS-TOTAL
               WRITE OUT-RECORD FROM WS-LINE
               PERFORM READ-PARA
           END-PERFORM
           MOVE 0 TO WS-TOTAL
           CLOSE IN-FILE OUT-FILE
           STOP RUN.
       READ-PARA.
           READ IN-FILE
               AT END MOVE 'Y' TO WS-EOF
           END-READ.
