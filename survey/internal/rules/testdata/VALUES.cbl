       IDENTIFICATION DIVISION.
       PROGRAM-ID. VALUES.
       PROCEDURE DIVISION.
       GIVE-PARA.
           MOVE SPACES TO RESULT-GROUP
           IF KIND = 'A'
               PERFORM DECLINE-PARA
           END-IF
           IF RESULT = SPACES
               DISPLAY 'NONE'
           ELSE
               IF RESULT = 'OK'
                   DISPLAY 'OK'
               END-IF
           END-IF
           STOP RUN.
       DECLINE-PARA.
           MOVE 'NO' TO RESULT.
