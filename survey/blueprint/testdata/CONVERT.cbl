       IDENTIFICATION DIVISION.
       PROGRAM-ID. CONVERT.
       DATA DIVISION.
       WORKING-STORAGE SECTION.
       01  WS-YEAR                     PIC 9(4).
       LINKAGE SECTION.
       01  LK-DIRECTION                PIC X.
       01  LK-WESTERN.
           05  LK-W-YEAR               PIC 9(4).
           05  LK-W-MMDD               PIC 9(4).
       01  LK-ERA.
           05  LK-E-ERA                PIC X.
           05  LK-E-YEAR               PIC 99.
           05  LK-E-MMDD               PIC 9(4).
       PROCEDURE DIVISION USING LK-DIRECTION LK-WESTERN LK-ERA.
       MAIN-PARA.
           IF LK-DIRECTION = 'E'
               MOVE 'R' TO LK-E-ERA
               COMPUTE LK-E-YEAR = LK-W-YEAR - 2018
               MOVE LK-W-MMDD TO LK-E-MMDD
           ELSE
               COMPUTE LK-W-YEAR = LK-E-YEAR + 2018
               MOVE LK-E-MMDD TO LK-W-MMDD
           END-IF
           GOBACK.
