       IDENTIFICATION DIVISION.
       PROGRAM-ID. CALLER.
       DATA DIVISION.
       WORKING-STORAGE SECTION.
       01  WS-REQUEST.
           05  WS-AMOUNT               PIC 9(5).
           05  WS-HELD                 PIC 9(5).
       01  WS-DATE                     PIC 9(8).
       01  WS-OUT                      PIC 9(5).
       PROCEDURE DIVISION.
       MAIN-PARA.
           ACCEPT WS-REQUEST
           CALL 'CALLEE' USING WS-REQUEST BY CONTENT WS-DATE
           MOVE WS-HELD TO WS-OUT
           DISPLAY WS-OUT
           GOBACK.
