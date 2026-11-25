       IDENTIFICATION DIVISION.
       PROGRAM-ID. KEEPS.
       PROCEDURE DIVISION.
       LOAN-PARA.
           IF BALANCE > 0
               PERFORM INTEREST-PARA
           END-IF
           PERFORM SETTLE-PARA
           IF IN-FORCE AND BALANCE > 0
               PERFORM EXCESS-PARA
           END-IF
           IF MATURED
               DISPLAY 'MATURED'
           END-IF.
       INTEREST-PARA.
           ADD 1 TO BALANCE.
       SETTLE-PARA.
           IF DUE = 1
               MOVE 0 TO BALANCE
           END-IF.
       EXCESS-PARA.
           SET LAPSED TO TRUE.
