       IDENTIFICATION DIVISION.
       PROGRAM-ID. CALLER.
       DATA DIVISION.
       WORKING-STORAGE SECTION.
       01  WS-REQUEST                  PIC X(10).
       01  WS-DATE                     PIC 9(8).
       PROCEDURE DIVISION.
       MAIN-PARA.
           CALL 'CALLEE' USING WS-REQUEST BY CONTENT WS-DATE
           CALL 'MISSING'
           GOBACK.
