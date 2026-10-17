       IDENTIFICATION DIVISION.
       PROGRAM-ID. REPLY.
       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT IN-FILE ASSIGN TO INDD.
       DATA DIVISION.
       FILE SECTION.
       FD  IN-FILE.
       01  IN-RECORD                   PIC X(10).
       WORKING-STORAGE SECTION.
       01  WS-IN.
           05  WS-IN-KEY               PIC X(4).
           05  WS-IN-REST              PIC X(6).
       01  WS-REQUEST                  PIC X(10).
       01  WS-REPLY                    PIC X(10).
       PROCEDURE DIVISION.
       MAIN-PARA.
           ACCEPT WS-REQUEST
           READ IN-FILE INTO WS-IN
           MOVE WS-IN-KEY TO WS-REPLY
           DISPLAY 'REPLY: ' WS-REPLY
           DISPLAY 'FAILED' UPON SYSERR
           GOBACK.
