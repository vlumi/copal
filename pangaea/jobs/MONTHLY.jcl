//MONTHLY  JOB
//*================================================================*
//* MONTH-END: THE STATEMENTS OF THE MONTH THAT CLOSED. THE        *
//* SCHEDULER RUNS IT ON THE FIRST BUSINESS DAY OF EVERY MONTH,    *
//* AFTER THE DAILY BATCH HAS POSTED THE MONTH'S INTEREST AND      *
//* FEES. ON THE FIRST DAY OF A NEW BANK NO MONTH HAS CLOSED, AND  *
//* IT PRINTS NOTHING.                                             *
//*                                                                *
//* READS   ACCT.ACCRUAL(0)  ACCT.MASTER(-1)  TRAN.AUDIT           *
//*================================================================*
//*----------------------------------------------------------------*
//* PRINT THE STATEMENTS FROM EVERY AUDIT TRAIL AND THE MASTER THE *
//* CLOSED MONTH ENDED WITH, BEFORE THE DAILY BATCH OF TODAY       *
//*----------------------------------------------------------------*
//STMT     EXEC PGM=STMTPRT
//ACCRIN   DD DSN=ACCT.ACCRUAL(0),DISP=SHR
//ACCTIN   DD DSN=ACCT.MASTER(-1),DISP=SHR
//AUDIT    DD DSN=TRAN.AUDIT,DISP=SHR
//STMTOUT  DD SYSOUT=*
//STMTRPT  DD SYSOUT=*
