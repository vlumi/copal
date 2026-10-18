       IDENTIFICATION DIVISION.
       PROGRAM-ID. ESCAPE.
       DATA DIVISION.
       WORKING-STORAGE SECTION.
       01  WS-A                        PIC 9.
       01  WS-B                        PIC 9.
       PROCEDURE DIVISION.
       MAIN-PARA.
           MOVE 1 TO WS-A
           PERFORM WORK-PARA
           STOP RUN.
       WORK-PARA.
           MOVE 2 TO WS-A
           GO TO DONE-PARA.
       DONE-PARA.
           MOVE WS-A TO WS-B
           STOP RUN.
