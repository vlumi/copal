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
       OPEN-PARA.
           EXIT.
       REJECT-PARA.
           EXIT.
