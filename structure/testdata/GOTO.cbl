       IDENTIFICATION DIVISION.
       PROGRAM-ID. GOTO.
       DATA DIVISION.
       WORKING-STORAGE SECTION.
       01  WS-I                        PIC 9.
       PROCEDURE DIVISION.
       MAIN-PARA.
           PERFORM CHECK-PARA THRU CHECK-EXIT
           STOP RUN.
       CHECK-PARA.
           IF WS-I = 0
               GO TO CHECK-EXIT.
           GO TO CHECK-ONE CHECK-TWO DEPENDING ON WS-I.
       CHECK-ONE.
           MOVE 1 TO WS-I
           GO TO CHECK-EXIT.
       CHECK-TWO.
           MOVE 2 TO WS-I.
       CHECK-EXIT.
           EXIT.
