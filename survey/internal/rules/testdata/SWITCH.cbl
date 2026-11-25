       IDENTIFICATION DIVISION.
       PROGRAM-ID. SWITCH.
       PROCEDURE DIVISION.
       WRITE-PARA.
           IF NOT-OPEN
               PERFORM OPEN-PARA
           END-IF
           IF IS-OPEN
               DISPLAY 'WRITE'
           END-IF
           GOBACK.
       OPEN-PARA.
           IF AVAILABLE = 'Y'
               SET IS-OPEN TO TRUE
           END-IF.
       CLOSE-PARA.
           SET NOT-OPEN TO TRUE.
