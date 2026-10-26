       IDENTIFICATION DIVISION.
       PROGRAM-ID. TOTALS.
       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT RESULT-IN ASSIGN TO RESULTS
               ORGANIZATION IS SEQUENTIAL.
       DATA DIVISION.
       FILE SECTION.
       FD  RESULT-IN.
       01  RESULT-RECORD.
           05  RS-TYPE                 PIC X.
               88  RS-DETAIL           VALUE 'D'.
               88  RS-TRAILER          VALUE 'T'.
           05  RS-DETAIL-FIELDS.
               10  RS-PAID             PIC X.
                   88  RS-WAS-PAID     VALUE 'Y'.
               10  RS-AMOUNT           PIC 9(7).
           05  RS-TRAILER-FIELDS REDEFINES RS-DETAIL-FIELDS.
               10  RT-COUNT            PIC 9(4).
               10  RT-PAID-COUNT       PIC 9(4).
       WORKING-STORAGE SECTION.
       01  WS-END                      PIC X VALUE 'N'.
           88  AT-END                  VALUE 'Y'.
       01  WS-COUNT                    PIC 9(4) VALUE 0.
       01  WS-PAID-COUNT               PIC 9(4) VALUE 0.
       PROCEDURE DIVISION.
       MAIN-PARA.
           OPEN INPUT RESULT-IN
           PERFORM READ-PARA UNTIL AT-END
           CLOSE RESULT-IN
           STOP RUN.
       READ-PARA.
           READ RESULT-IN
               AT END
                   SET AT-END TO TRUE
               NOT AT END
                   IF RS-DETAIL
                       ADD 1 TO WS-COUNT
                       IF RS-WAS-PAID AND NOT AT-END
                           ADD 1 TO WS-PAID-COUNT
                       END-IF
                   ELSE
                       IF RT-COUNT NOT = WS-COUNT
                         OR RT-PAID-COUNT NOT = WS-PAID-COUNT
                           MOVE 8 TO RETURN-CODE
                       END-IF
                   END-IF
           END-READ.
