       IDENTIFICATION DIVISION.
       PROGRAM-ID. COPYIT.
       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT IN-FILE ASSIGN TO INFILE
               ORGANIZATION IS SEQUENTIAL.
           SELECT OUT-FILE ASSIGN TO OUTFILE
               ORGANIZATION IS SEQUENTIAL.
       DATA DIVISION.
       FILE SECTION.
       FD  IN-FILE.
       01  IN-RECORD                   PIC X(5).
       FD  OUT-FILE.
       01  OUT-RECORD                  PIC X(5).
       WORKING-STORAGE SECTION.
       01  WS-EOF                      PIC X          VALUE 'N'.
       01  WS-COUNT                    PIC 9          VALUE 0.
       PROCEDURE DIVISION.
       MAIN-PARA.
           OPEN INPUT IN-FILE OUTPUT OUT-FILE
           PERFORM UNTIL WS-EOF = 'Y'
               READ IN-FILE
                   AT END MOVE 'Y' TO WS-EOF
                   NOT AT END
                       ADD 1 TO WS-COUNT
                       WRITE OUT-RECORD FROM IN-RECORD
               END-READ
           END-PERFORM
           CLOSE IN-FILE OUT-FILE
           DISPLAY 'COPIED ' WS-COUNT
           MOVE WS-COUNT TO RETURN-CODE
           STOP RUN.
