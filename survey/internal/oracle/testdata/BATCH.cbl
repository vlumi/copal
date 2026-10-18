       IDENTIFICATION DIVISION.
       PROGRAM-ID. BATCH.
       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT CONTROL-IN ASSIGN TO CTLIN
               ORGANIZATION IS SEQUENTIAL.
           SELECT ACCT-IN ASSIGN TO ACCTIN
               ORGANIZATION IS INDEXED
               ACCESS MODE IS RANDOM
               RECORD KEY IS AC-ACCT-NO.
           SELECT BATCH-IN ASSIGN TO BATCHIN
               ORGANIZATION IS SEQUENTIAL.
           SELECT REPORT-OUT ASSIGN TO RPTOUT
               ORGANIZATION IS LINE SEQUENTIAL.
       DATA DIVISION.
       FILE SECTION.
       FD  CONTROL-IN.
       01  CONTROL-RECORD.
           05  CT-BUSINESS-DATE        PIC 9(8).
       FD  ACCT-IN.
       01  ACCT-RECORD.
           05  AC-ACCT-NO              PIC 9(6).
           05  AC-STATUS               PIC X.
               88  AC-OPEN             VALUE 'O'.
               88  AC-CLOSED           VALUE 'X'.
       FD  BATCH-IN.
           COPY BATCHREC.
       FD  REPORT-OUT.
       01  REPORT-LINE                 PIC X(40).
       WORKING-STORAGE SECTION.
       01  WS-BUSINESS-DATE            PIC 9(8).
       01  WS-COUNT                    PIC 9(6)       VALUE 0.
       01  WS-TOTAL                    PIC 9(9)V99    VALUE 0.
       01  WS-REJECTS                  PIC 9(6)       VALUE 0.
       01  WS-EOF                      PIC X          VALUE 'N'.
           88  END-OF-BATCH            VALUE 'Y'.
       PROCEDURE DIVISION.
       MAIN-PARA.
           OPEN INPUT CONTROL-IN ACCT-IN BATCH-IN OUTPUT REPORT-OUT
           READ CONTROL-IN
               AT END
                   MOVE 16 TO RETURN-CODE
                   STOP RUN
           END-READ
           MOVE CT-BUSINESS-DATE TO WS-BUSINESS-DATE
           PERFORM READ-BATCH
           IF END-OF-BATCH OR NOT BR-HEADER
               MOVE 'NO HEADER' TO REPORT-LINE
               GO TO FINISH-BAD
           END-IF
           IF BH-DATE NOT = WS-BUSINESS-DATE
               MOVE 'WRONG DATE' TO REPORT-LINE
               GO TO FINISH-BAD
           END-IF
           PERFORM READ-BATCH
           PERFORM UNTIL END-OF-BATCH OR BR-TRAILER
               PERFORM CHECK-ITEM
               PERFORM READ-BATCH
           END-PERFORM
           IF END-OF-BATCH
               MOVE 'NO TRAILER' TO REPORT-LINE
               GO TO FINISH-BAD
           END-IF
           IF BT-COUNT NOT = WS-COUNT OR BT-TOTAL NOT = WS-TOTAL
               MOVE 'TOTALS DIFFER' TO REPORT-LINE
               GO TO FINISH-BAD
           END-IF
           IF WS-REJECTS > 0
               MOVE 4 TO RETURN-CODE
           END-IF
           MOVE 'ACCEPTED' TO REPORT-LINE
           WRITE REPORT-LINE
           CLOSE CONTROL-IN ACCT-IN BATCH-IN REPORT-OUT
           STOP RUN.
       READ-BATCH.
           READ BATCH-IN
               AT END
                   SET END-OF-BATCH TO TRUE
           END-READ.
       CHECK-ITEM.
           ADD 1 TO WS-COUNT
           ADD BI-AMOUNT TO WS-TOTAL
           MOVE BI-ACCT-NO TO AC-ACCT-NO
           READ ACCT-IN
               INVALID KEY
                   ADD 1 TO WS-REJECTS
               NOT INVALID KEY
                   IF AC-CLOSED
                       ADD 1 TO WS-REJECTS
                   END-IF
           END-READ
           IF BI-DEBIT AND BI-AMOUNT > 500
               ADD 1 TO WS-REJECTS
           END-IF.
       FINISH-BAD.
           WRITE REPORT-LINE
           CLOSE CONTROL-IN ACCT-IN BATCH-IN REPORT-OUT
           MOVE 12 TO RETURN-CODE
           STOP RUN.
