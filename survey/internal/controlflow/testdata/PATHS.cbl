       IDENTIFICATION DIVISION.
       PROGRAM-ID. PATHS.
       DATA DIVISION.
       WORKING-STORAGE SECTION.
       01  WS-A                        PIC 9.
       01  WS-B                        PIC 9.
       PROCEDURE DIVISION.
       MAIN-PARA.
           MOVE 1 TO WS-A
           PERFORM WORK-PARA
           MOVE 2 TO WS-A
           PERFORM WORK-PARA
           PERFORM CHECK-PARA THRU CHECK-EXIT
           PERFORM UNTIL WS-A = 9
               ADD 1 TO WS-A
           END-PERFORM
           STOP RUN.
       WORK-PARA.
           IF WS-A = 1
               MOVE 3 TO WS-B
           END-IF
           MOVE WS-A TO WS-B.
       CHECK-PARA.
           IF WS-B = 3
               GO TO CHECK-EXIT
           ELSE
               MOVE 4 TO WS-B
           END-IF
           MOVE 5 TO WS-B.
       CHECK-EXIT.
           EXIT.
       ESCAPE-PARA.
           GO TO MAIN-PARA.
