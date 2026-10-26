# BANKDATE: decision tables

The business rules of [BANKDATE](../programs/BANKDATE.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

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

- **Rule 1.** When `WS-IN-STATUS` is not `'00'` or `WS-OUT-STATUS` is not `'00'`, it opens `CONTROL-IN` for input and `CONTROL-OUT` for output; displays `'BANKDATE: CANNOT OPEN FILES '` `WS-IN-STATUS` `' '` `WS-OUT-STATUS`; sets `RETURN-CODE` to 16; and stops the run.
- **Rule 2.** When `WS-IN-STATUS` is `'00'` and `WS-OUT-STATUS` is `'00'`; and reading `CONTROL-IN` finds no more records, it opens `CONTROL-IN` for input and `CONTROL-OUT` for output; reads `CONTROL-IN`; displays `'BANKDATE: NO BANK CONTROL RECORD'`; sets `RETURN-CODE` to 16; and stops the run.
- **Rule 3.** When `WS-IN-STATUS` is `'00'` and `WS-OUT-STATUS` is `'00'`; and reading `CONTROL-IN` finds a record, it opens `CONTROL-IN` for input and `CONTROL-OUT` for output; reads `CONTROL-IN`; performs `1000-NEXT-BUSINESS-DAY`; writes `CONTROL-OUT-RECORD` from `BANK-CONTROL-RECORD`; closes `CONTROL-IN` and `CONTROL-OUT`; sets `RETURN-CODE` to 0; and stops the run.

Source: [BANKDATE.cbl:52](../../../pangaea/programs/BANKDATE.cbl#L52).
