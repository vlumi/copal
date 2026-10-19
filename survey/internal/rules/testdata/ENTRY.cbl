       IDENTIFICATION DIVISION.
       PROGRAM-ID. ENTRY.
       PROCEDURE DIVISION.
       MAIN-PARA.
           IF TOTAL = 0
               DISPLAY 'NONE'
           END-IF
           IF ERR = SPACES
               PERFORM CHECK-PARA
               PERFORM JUMPED-PARA
           END-IF
           PERFORM LOOP-PARA UNTIL DONE = 1
           GO TO JUMPED-PARA.
       CHECK-PARA.
           IF ERR = SPACES
               DISPLAY 'GOOD'
           END-IF.
       LOOP-PARA.
           IF DONE = 1
               DISPLAY 'DONE'
           END-IF.
       JUMPED-PARA.
           IF ERR = SPACES
               DISPLAY 'GOOD'
           END-IF
           STOP RUN.
