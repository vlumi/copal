       IDENTIFICATION DIVISION.
       PROGRAM-ID. COPIES.
       PROCEDURE DIVISION.
       COPY-PARA.
           MOVE MASTER-KEY TO ACTIVE-KEY
           IF REQUEST-KEY < MASTER-KEY
               MOVE REQUEST-KEY TO ACTIVE-KEY
           END-IF
           IF MASTER-KEY = ACTIVE-KEY
               DISPLAY 'MASTER'
           END-IF.
       DIGITS-PARA.
           MOVE 0 TO MONTH-NUMBER
           IF CM-MONTH = 12
               DISPLAY 'DECEMBER'
           END-IF.
