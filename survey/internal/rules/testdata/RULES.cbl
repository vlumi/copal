       IDENTIFICATION DIVISION.
       PROGRAM-ID. RULES.
       PROCEDURE DIVISION.
       FEE-PARA.
           EVALUATE TRUE
               WHEN AM-CLOSED
                   SET REJECTED TO TRUE
               WHEN AM-CURRENT AND AM-OPEN
                   IF AM-BALANCE < 1000
                       MOVE 2.50 TO FEE
                   END-IF
                   ADD 1 TO FEE-COUNT
           END-EVALUATE.
       CHAIN-PARA.
           IF MR-OPEN
               GO TO OPEN-PARA.
           IF NOT HELD
               SET EXISTS TO TRUE.
           GO TO REJECT-PARA.
       READ-PARA.
           READ IN-FILE
               AT END
                   MOVE HIGH-VALUES TO IN-KEY
               NOT AT END
                   ADD 1 TO READ-COUNT
           END-READ.
       LOOP-PARA.
           PERFORM VARYING WS-I FROM 1 BY 1 UNTIL WS-I > 3
               IF WS-RATE(WS-I) > 0
                   MOVE WS-RATE(WS-I) TO WS-BEST
               END-IF
           END-PERFORM.
       EXIT-PARA.
           IF DONE = 1
               EXIT PARAGRAPH
           END-IF
           DISPLAY 'NOT DONE'.
       AGAIN-PARA.
           IF OK = 1
               PERFORM OPEN-PARA
           END-IF
           IF OK = 1
               DISPLAY 'STILL OK'
           END-IF.
       ABORT-USE-PARA.
           IF BAD = 1
               PERFORM ABORT-PARA
           END-IF
           DISPLAY 'GOOD'.
       ABORT-PARA.
           DISPLAY 'ABORT'
           STOP RUN.
       KNOWN-PARA.
           IF KIND = 'A'
               MOVE 'BAD' TO ERR
           END-IF
           IF ERR = SPACES
               DISPLAY 'FINE'
           END-IF.
       NAMED-PARA.
           SET STATUS-OPEN TO TRUE
           IF STATUS-CLOSED
               DISPLAY 'CLOSED'
           END-IF.
       WIDE-PARA.
           IF A1 = 1
               DISPLAY 'Y1'
           ELSE
               DISPLAY 'N1'
           END-IF
           IF A2 = 1
               DISPLAY 'Y2'
           ELSE
               DISPLAY 'N2'
           END-IF
           IF A3 = 1
               DISPLAY 'Y3'
           ELSE
               DISPLAY 'N3'
           END-IF
           IF A4 = 1
               DISPLAY 'Y4'
           ELSE
               DISPLAY 'N4'
           END-IF
           IF A5 = 1
               DISPLAY 'Y5'
           ELSE
               DISPLAY 'N5'
           END-IF
           IF A6 = 1
               DISPLAY 'Y6'
           ELSE
               DISPLAY 'N6'
           END-IF
           IF B = 1
               GO TO OPEN-PARA
           END-IF
           DISPLAY 'AFTER'.
       JOINED-PARA.
           IF DONE = 1 AND OK = 1
               DISPLAY 'BOTH'
           END-IF
           ADD 1 TO OK
           IF DONE = 1
               DISPLAY 'DONE'
           END-IF.
       CLASS-PARA.
           IF AMT IS NUMERIC
               DISPLAY 'NUMBER'
           END-IF
           IF AMT IS NOT NUMERIC OR AMT = ZERO
               DISPLAY 'BAD'
           END-IF.
       PASS-PARA.
           PERFORM UNTIL DONE = 1 OR BAD = 1
               IF BAD = 1
                   DISPLAY 'NEVER'
               END-IF
               ADD 1 TO DONE
               IF DONE = 1
                   DISPLAY 'DONE'
               END-IF
           END-PERFORM
           PERFORM WITH TEST AFTER UNTIL BAD = 1
               IF BAD = 1
                   DISPLAY 'LAST'
               END-IF
           END-PERFORM.
       OPEN-PARA.
           EXIT.
       REJECT-PARA.
           EXIT.
       LEAVE-SECTION SECTION.
           IF DONE = 1
               EXIT SECTION
           END-IF
           DISPLAY 'STAYED'.
