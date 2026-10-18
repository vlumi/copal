       IDENTIFICATION DIVISION.
       PROGRAM-ID. CALLEE.
       DATA DIVISION.
       LINKAGE SECTION.
       01  LK-REQUEST.
           05  LK-AMOUNT               PIC 9(5).
           05  LK-HELD                 PIC 9(5).
       01  LK-DATE                     PIC 9(8).
       PROCEDURE DIVISION USING LK-REQUEST LK-DATE.
       MAIN-PARA.
           MOVE LK-AMOUNT TO LK-HELD
           GOBACK.
