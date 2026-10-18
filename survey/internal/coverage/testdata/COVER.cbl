       IDENTIFICATION DIVISION.
       PROGRAM-ID. COVER.
       DATA DIVISION.
       WORKING-STORAGE SECTION.
       01  WS-N                        PIC 9          VALUE 0.
       01  WS-KIND                     PIC X.
       PROCEDURE DIVISION.
       MAIN-PARA.
           MOVE 'A' TO WS-KIND
           PERFORM CHECK-PARA
           MOVE 'B' TO WS-KIND
           PERFORM CHECK-PARA
           PERFORM UNTIL WS-N = 3
               ADD 1 TO WS-N
               IF WS-N = 2
                   DISPLAY 'TWO'
               END-IF
           END-PERFORM
           STOP RUN.
       CHECK-PARA.
           EVALUATE WS-KIND
               WHEN 'A'
                   DISPLAY 'KIND A'
               WHEN 'B'
                   DISPLAY 'KIND B'
               WHEN OTHER
                   DISPLAY 'OTHER'
           END-EVALUATE.
