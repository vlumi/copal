//DAILY    JOB
//*================================================================*
//* DAILY BATCH: VALIDATE, SORT AND POST THE DAY'S TRANSACTIONS,   *
//* THEN THE TRANSFER CREDITS AND RETURNS THEY LEAD TO, AND MOVE   *
//* TO THE NEXT BUSINESS DAY.                                      *
//*                                                                *
//* READS   BANK.CONTROL(0)  ACCT.MASTER(0)  TRAN.DAILY(0)         *
//* WRITES  BANK.CONTROL(+1) ACCT.MASTER(+1)                       *
//*         TRAN.AUDIT(+1)   TRAN.REJECTS(+1)                      *
//*                                                                *
//* A STEP RUNS ONLY IF NO EARLIER STEP ENDED WITH RC ABOVE 4, SO  *
//* A RUN THAT FAILS LEAVES THE BUSINESS DATE WHERE IT WAS.        *
//*================================================================*
//*----------------------------------------------------------------*
//* VALIDATE THE DAY'S TRANSACTION FILE                            *
//*----------------------------------------------------------------*
//VALID    EXEC PGM=TXVALID
//BANKCTL  DD DSN=BANK.CONTROL(0),DISP=SHR
//TRANIN   DD DSN=TRAN.DAILY(0),DISP=SHR
//TRANOUT  DD DSN=&&VALID,DISP=NEW,
//            DCB=(RECFM=FB,LRECL=80)
//REJECTS  DD DSN=TRAN.REJECTS(+1),DISP=NEW,
//            DCB=(RECFM=FB,LRECL=100)
//VALRPT   DD SYSOUT=*
//*----------------------------------------------------------------*
//* SORT THEM INTO ACCOUNT ORDER AND POST THEM                     *
//*----------------------------------------------------------------*
//SORT     EXEC PGM=TXSORT,COND=(4,LT)
//TRANIN   DD DSN=&&VALID,DISP=OLD
//TRANOUT  DD DSN=&&SORTED,DISP=NEW,
//            DCB=(RECFM=FB,LRECL=80)
//POST     EXEC PGM=TXPOST,COND=(4,LT)
//BANKCTL  DD DSN=BANK.CONTROL(0),DISP=SHR
//ACCTIN   DD DSN=ACCT.MASTER(0),DISP=SHR
//TRANIN   DD DSN=&&SORTED,DISP=OLD
//ACCTOUT  DD DSN=&&POSTED,DISP=NEW,
//            DCB=(RECFM=FB,LRECL=100)
//AUDIT    DD DSN=TRAN.AUDIT(+1),DISP=MOD,
//            DCB=(RECFM=FB,LRECL=100)
//REJECTS  DD DSN=TRAN.REJECTS(+1),DISP=MOD
//LEGS     DD DSN=&&CREDITS,DISP=NEW,
//            DCB=(RECFM=FB,LRECL=80)
//POSTRPT  DD SYSOUT=*
//*----------------------------------------------------------------*
//* SORT AND POST THE CREDITS OF THE TRANSFERS THAT POSTED         *
//*----------------------------------------------------------------*
//SORTCR   EXEC PGM=TXSORT,COND=(4,LT)
//TRANIN   DD DSN=&&CREDITS,DISP=OLD
//TRANOUT  DD DSN=&&CREDSORT,DISP=NEW,
//            DCB=(RECFM=FB,LRECL=80)
//POSTCR   EXEC PGM=TXPOST,COND=(4,LT)
//BANKCTL  DD DSN=BANK.CONTROL(0),DISP=SHR
//ACCTIN   DD DSN=&&POSTED,DISP=OLD
//TRANIN   DD DSN=&&CREDSORT,DISP=OLD
//ACCTOUT  DD DSN=&&CREDITED,DISP=NEW,
//            DCB=(RECFM=FB,LRECL=100)
//AUDIT    DD DSN=TRAN.AUDIT(+1),DISP=MOD
//REJECTS  DD DSN=TRAN.REJECTS(+1),DISP=MOD
//LEGS     DD DSN=&&RETURNS,DISP=NEW,
//            DCB=(RECFM=FB,LRECL=80)
//POSTRPT  DD SYSOUT=*
//*----------------------------------------------------------------*
//* SORT AND POST THE RETURNS OF THE CREDITS THAT COULD NOT POST,  *
//* WRITING THE NEW MASTER                                         *
//*----------------------------------------------------------------*
//SORTRT   EXEC PGM=TXSORT,COND=(4,LT)
//TRANIN   DD DSN=&&RETURNS,DISP=OLD
//TRANOUT  DD DSN=&&RETSORT,DISP=NEW,
//            DCB=(RECFM=FB,LRECL=80)
//POSTRT   EXEC PGM=TXPOST,COND=(4,LT)
//BANKCTL  DD DSN=BANK.CONTROL(0),DISP=SHR
//ACCTIN   DD DSN=&&CREDITED,DISP=OLD
//TRANIN   DD DSN=&&RETSORT,DISP=OLD
//ACCTOUT  DD DSN=ACCT.MASTER(+1),DISP=NEW,
//            DCB=(RECFM=FB,LRECL=100)
//AUDIT    DD DSN=TRAN.AUDIT(+1),DISP=MOD
//REJECTS  DD DSN=TRAN.REJECTS(+1),DISP=MOD
//LEGS     DD DUMMY
//POSTRPT  DD SYSOUT=*
//*----------------------------------------------------------------*
//* MOVE TO THE NEXT BUSINESS DAY                                  *
//*----------------------------------------------------------------*
//NEXTDAY  EXEC PGM=BANKDATE,COND=(4,LT)
//BANKIN   DD DSN=BANK.CONTROL(0),DISP=SHR
//BANKOUT  DD DSN=BANK.CONTROL(+1),DISP=NEW,
//            DCB=(RECFM=FB,LRECL=80)
