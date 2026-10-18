       IDENTIFICATION DIVISION.
       PROGRAM-ID. ANSWERER.
       DATA DIVISION.
       LINKAGE SECTION.
       01  LK-QUESTION.
           05  LK-KIND                 PIC X.
           05  LK-ANSWER               PIC X.
       PROCEDURE DIVISION USING LK-QUESTION.
       MAIN-PARA.
           IF LK-KIND = 'A'
               MOVE 'Y' TO LK-ANSWER
           END-IF
           GOBACK.
