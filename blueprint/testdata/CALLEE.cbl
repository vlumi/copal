       IDENTIFICATION DIVISION.
       PROGRAM-ID. CALLEE.
       DATA DIVISION.
       LINKAGE SECTION.
       01  LK-REQUEST                  PIC X(10).
       01  LK-DATE                     PIC 9(8).
       PROCEDURE DIVISION USING LK-REQUEST LK-DATE.
       MAIN-PARA.
           MOVE LK-DATE TO LK-REQUEST
           GOBACK.
