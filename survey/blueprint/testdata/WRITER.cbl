       IDENTIFICATION DIVISION.
       PROGRAM-ID. WRITER.
       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT IN-FILE ASSIGN TO INFILE.
           SELECT OUT-FILE ASSIGN TO OUTFILE.
       DATA DIVISION.
       FILE SECTION.
       FD  IN-FILE.
       01  IN-RECORD.
           05  IN-DATE                 PIC 9(8).
           05  IN-NAME                 PIC X(20).
       FD  OUT-FILE.
       01  OUT-RECORD.
           05  OUT-DATE                PIC X(7).
           05  OUT-NAME                PIC X(20).
       WORKING-STORAGE SECTION.
       01  WS-DIRECTION                PIC X VALUE 'E'.
       01  WS-ERA-DATE                 PIC X(7).
       PROCEDURE DIVISION.
       MAIN-PARA.
           OPEN INPUT IN-FILE OUTPUT OUT-FILE
           READ IN-FILE
           CALL 'CONVERT' USING WS-DIRECTION IN-DATE WS-ERA-DATE
           MOVE WS-ERA-DATE TO OUT-DATE
           MOVE IN-NAME TO OUT-NAME
           WRITE OUT-RECORD
           CLOSE IN-FILE OUT-FILE
           GOBACK.
