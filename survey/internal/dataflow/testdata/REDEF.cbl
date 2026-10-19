       IDENTIFICATION DIVISION.
       PROGRAM-ID. REDEF.
       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT IN-FILE ASSIGN TO INFILE.
       DATA DIVISION.
       FILE SECTION.
       FD  IN-FILE.
       01  IN-RECORD.
           05  IN-MONTH                PIC 9(6).
       WORKING-STORAGE SECTION.
       01  WS-MONTH.
           05  WS-YEAR                 PIC 9(4).
           05  WS-MM                   PIC 99.
       01  WS-MONTH-NUMBER REDEFINES WS-MONTH PIC 9(6).
       PROCEDURE DIVISION.
       MAIN-PARA.
           OPEN INPUT IN-FILE
           READ IN-FILE
           MOVE IN-MONTH TO WS-MONTH-NUMBER
           IF WS-MM = 12
               DISPLAY 'DECEMBER'
           END-IF
           STOP RUN.
