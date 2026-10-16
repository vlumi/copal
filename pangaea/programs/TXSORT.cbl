      *================================================================*
      * TXSORT: SORT THE VALID TRANSACTIONS FOR POSTING                *
      *                                                                *
      * POSTING MATCHES TRANSACTIONS AGAINST THE ACCOUNT MASTER, SO    *
      * IT NEEDS THEM IN ACCOUNT NUMBER ORDER. WITHIN AN ACCOUNT THEY  *
      * KEEP THE ORDER OF THEIR TRANSACTION IDS.                       *
      *                                                                *
      * RETURN CODES: 0 SORTED, 16 SORT FAILED.                        *
      *================================================================*
       IDENTIFICATION DIVISION.
       PROGRAM-ID. TXSORT.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT TRAN-IN   ASSIGN TO TRANIN
               ORGANIZATION IS SEQUENTIAL.
           SELECT TRAN-OUT  ASSIGN TO TRANOUT
               ORGANIZATION IS SEQUENTIAL.
           SELECT SORT-WORK ASSIGN TO SORTWK.

       DATA DIVISION.
       FILE SECTION.
       FD  TRAN-IN.
       01  TRAN-IN-RECORD              PIC X(80).
       FD  TRAN-OUT.
       01  TRAN-OUT-RECORD             PIC X(80).
       SD  SORT-WORK.
       01  SORT-RECORD.
           05  SR-REC-TYPE             PIC X.
           05  SR-TRAN-ID              PIC 9(10).
           05  SR-ACCT-NO              PIC 9(10).
           05  FILLER                  PIC X(59).

       PROCEDURE DIVISION.
       0000-MAIN.
           SORT SORT-WORK
               ON ASCENDING KEY SR-ACCT-NO SR-TRAN-ID
               USING TRAN-IN
               GIVING TRAN-OUT
           IF SORT-RETURN NOT = 0
               DISPLAY 'TXSORT: SORT FAILED ' SORT-RETURN
               MOVE 16 TO RETURN-CODE
           ELSE
               MOVE 0 TO RETURN-CODE
           END-IF
           STOP RUN.
