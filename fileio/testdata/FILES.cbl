       IDENTIFICATION DIVISION.
       PROGRAM-ID. FILES.
       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT IN-FILE   ASSIGN TO INDD.
           SELECT OUT-FILE  ASSIGN TO OUTDD.
           SELECT LOG-FILE  ASSIGN TO LOGDD.
           SELECT WORK-SORT ASSIGN TO SORTWK.
           SELECT SORTED    ASSIGN TO SORTDD.
       DATA DIVISION.
       FILE SECTION.
       FD  IN-FILE.
       01  IN-RECORD                   PIC X(10).
       FD  OUT-FILE.
       01  OUT-RECORD                  PIC X(10).
       FD  LOG-FILE.
       01  LOG-RECORD                  PIC X(10).
       SD  WORK-SORT.
       01  WORK-RECORD                 PIC X(10).
       FD  SORTED.
       01  SORTED-RECORD               PIC X(10).
       PROCEDURE DIVISION.
       MAIN-PARA.
           OPEN INPUT IN-FILE
                OUTPUT OUT-FILE
                EXTEND LOG-FILE
           SORT WORK-SORT ON ASCENDING KEY WORK-RECORD
               USING IN-FILE
               GIVING SORTED
           CLOSE IN-FILE OUT-FILE LOG-FILE
           OPEN INPUT IN-FILE
           STOP RUN.
