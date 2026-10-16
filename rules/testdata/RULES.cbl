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
       OPEN-PARA.
           EXIT.
       REJECT-PARA.
           EXIT.
