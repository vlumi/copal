       IDENTIFICATION DIVISION.
       PROGRAM-ID. FLOW.
       DATA DIVISION.
       WORKING-STORAGE SECTION.
       01  WS-I                        PIC 99.
       PROCEDURE DIVISION.
       MAIN-PARA.
           PERFORM INIT-PARA
           PERFORM WORK-PARA THRU WORK-EXIT
               VARYING WS-I FROM 1 BY 1 UNTIL WS-I > 3
           PERFORM 2 TIMES
               PERFORM INIT-PARA
           END-PERFORM
           STOP RUN.
       INIT-PARA.
           MOVE 0 TO WS-I.
       WORK-PARA.
           CONTINUE.
       WORK-EXIT.
           EXIT.
