      *================================================================*
      * BANKDATE: ADVANCE THE BUSINESS DATE                            *
      *                                                                *
      * READS THE BANK CONTROL RECORD AND WRITES THE NEXT ONE, FOR     *
      * THE NEXT BUSINESS DAY: THE NEXT DAY FROM MONDAY TO FRIDAY.     *
      * THE LAST STEP OF THE DAILY BATCH, SO THAT THE NEXT RUN IS FOR  *
      * THE NEXT DAY ONLY IF THIS ONE WENT THROUGH.                    *
      *                                                                *
      * RETURN CODES: 0 ADVANCED, 16 NO CONTROL RECORD OR I/O.         *
      *                                                                *
      * MODIFICATION HISTORY                                           *
      * DATE       BY  CHANGE                                          *
      * 1987-03-12 JA  ORIGINAL VERSION.                               *
      * 1999-01-11 SL  Y2K: 8-DIGIT DATES, INTRINSIC DATE FUNCTIONS    *
      *                INSTEAD OF THE CALENDAR TABLE.                  *
      *================================================================*
       IDENTIFICATION DIVISION.
       PROGRAM-ID. BANKDATE.
       AUTHOR. J. ARKWRIGHT, BATCH SYSTEMS.
       INSTALLATION. FIRST BANK OF PANGAEA, DATA CENTER.
       DATE-WRITTEN. 1987-03-12.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT CONTROL-IN  ASSIGN TO BANKIN
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-IN-STATUS.
           SELECT CONTROL-OUT ASSIGN TO BANKOUT
               ORGANIZATION IS SEQUENTIAL
               FILE STATUS IS WS-OUT-STATUS.

       DATA DIVISION.
       FILE SECTION.
       FD  CONTROL-IN.
           COPY BANKCTL.
       FD  CONTROL-OUT.
       01  CONTROL-OUT-RECORD          PIC X(80).

       WORKING-STORAGE SECTION.
       01  WS-IN-STATUS                PIC XX.
       01  WS-OUT-STATUS               PIC XX.

      * DAY NUMBERS COUNT FROM 1601-01-01, A MONDAY, SO A DAY NUMBER
      * MODULO 7 IS 1 ON MONDAYS, 6 ON SATURDAYS AND 0 ON SUNDAYS.
       01  WS-DAY-NUMBER               PIC 9(7).
       01  WS-WEEKDAY                  PIC 9.
           88  WEEKEND                 VALUE 6 0.

       PROCEDURE DIVISION.
      * READ TODAY'S CONTROL RECORD AND WRITE THE NEXT BUSINESS DAY'S.
       0000-MAIN.
           OPEN INPUT CONTROL-IN
                OUTPUT CONTROL-OUT
           IF WS-IN-STATUS NOT = '00' OR WS-OUT-STATUS NOT = '00'
               DISPLAY 'BANKDATE: CANNOT OPEN FILES '
                       WS-IN-STATUS ' ' WS-OUT-STATUS
               MOVE 16 TO RETURN-CODE
               STOP RUN
           END-IF
           READ CONTROL-IN
               AT END
                   DISPLAY 'BANKDATE: NO BANK CONTROL RECORD'
                   MOVE 16 TO RETURN-CODE
                   STOP RUN
           END-READ
           PERFORM 1000-NEXT-BUSINESS-DAY
           WRITE CONTROL-OUT-RECORD FROM BANK-CONTROL-RECORD
           CLOSE CONTROL-IN CONTROL-OUT
           MOVE 0 TO RETURN-CODE
           STOP RUN.

      * MOVE THE BUSINESS DATE TO THE NEXT DAY THAT IS NOT A WEEKEND.
       1000-NEXT-BUSINESS-DAY.
           COMPUTE WS-DAY-NUMBER =
               FUNCTION INTEGER-OF-DATE(BC-BUSINESS-DATE)
           ADD 1 TO WS-DAY-NUMBER
           COMPUTE WS-WEEKDAY = FUNCTION MOD(WS-DAY-NUMBER, 7)
           PERFORM UNTIL NOT WEEKEND
               ADD 1 TO WS-DAY-NUMBER
               COMPUTE WS-WEEKDAY = FUNCTION MOD(WS-DAY-NUMBER, 7)
           END-PERFORM
           COMPUTE BC-BUSINESS-DATE =
               FUNCTION DATE-OF-INTEGER(WS-DAY-NUMBER).
