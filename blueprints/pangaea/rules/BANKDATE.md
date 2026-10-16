# BANKDATE: decision tables

The business rules of [BANKDATE](../programs/BANKDATE.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

## 0000-MAIN

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `WS-IN-STATUS NOT = '00' OR WS-OUT-STATUS NOT = '00'` | Y | N | N |
| `READ CONTROL-IN AT END` | | Y | N |
| → `OPEN INPUT CONTROL-IN OUTPUT CONTROL-OUT` | 1 | 1 | 1 |
| → `DISPLAY 'BANKDATE: CANNOT OPEN FILES ' WS-IN-STATUS ' ' WS-OUT-STATUS` | 2 | | |
| → `MOVE 16 TO RETURN-CODE` | 3 | 4 | |
| → `STOP RUN` | 4 | 5 | 7 |
| → `READ CONTROL-IN` | | 2 | 2 |
| → `DISPLAY 'BANKDATE: NO BANK CONTROL RECORD'` | | 3 | |
| → `PERFORM 1000-NEXT-BUSINESS-DAY` | | | 3 |
| → `WRITE CONTROL-OUT-RECORD FROM BANK-CONTROL-RECORD` | | | 4 |
| → `CLOSE CONTROL-IN CONTROL-OUT` | | | 5 |
| → `MOVE 0 TO RETURN-CODE` | | | 6 |

Source: [BANKDATE.cbl:52](../../../pangaea/programs/BANKDATE.cbl#L52).
