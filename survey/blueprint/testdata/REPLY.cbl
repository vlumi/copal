       IDENTIFICATION DIVISION.
       PROGRAM-ID. REPLY.
       DATA DIVISION.
       WORKING-STORAGE SECTION.
       01  WS-REQUEST                  PIC X(10).
       01  WS-REPLY                    PIC X(10).
       PROCEDURE DIVISION.
       MAIN-PARA.
           ACCEPT WS-REQUEST
           DISPLAY 'REPLY: ' WS-REPLY
           DISPLAY 'FAILED' UPON SYSERR
           GOBACK.
