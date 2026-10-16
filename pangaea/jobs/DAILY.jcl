//DAILY    JOB
//*================================================================*
//* DAILY BATCH: VALIDATE, SORT AND POST THE DAY'S TRANSACTIONS    *
//* A STEP RUNS ONLY IF NO EARLIER STEP ENDED WITH RC ABOVE 4.     *
//*================================================================*
//*----------------------------------------------------------------*
//* VALIDATE THE DAY'S TRANSACTION FILE                            *
//*----------------------------------------------------------------*
//VALID    EXEC PGM=TXVALID
//BANKCTL  DD DSN=BANK.CONTROL,DISP=SHR
//TRANIN   DD DSN=TRAN.DAILY,DISP=SHR
//TRANOUT  DD DSN=TRAN.VALID,DISP=NEW,
//            DCB=(RECFM=FB,LRECL=80)
//REJECTS  DD DSN=TRAN.VALID.REJECTS,DISP=NEW,
//            DCB=(RECFM=FB,LRECL=100)
//VALRPT   DD SYSOUT=*
//*----------------------------------------------------------------*
//* SORT THE VALID TRANSACTIONS INTO ACCOUNT ORDER                 *
//*----------------------------------------------------------------*
//SORT     EXEC PGM=TXSORT,COND=(4,LT)
//TRANIN   DD DSN=TRAN.VALID,DISP=OLD
//TRANOUT  DD DSN=TRAN.SORTED,DISP=NEW,
//            DCB=(RECFM=FB,LRECL=80)
//*----------------------------------------------------------------*
//* POST THEM TO THE ACCOUNT MASTER, WRITING THE NEW MASTER        *
//*----------------------------------------------------------------*
//POST     EXEC PGM=TXPOST,COND=(4,LT)
//BANKCTL  DD DSN=BANK.CONTROL,DISP=SHR
//ACCTIN   DD DSN=ACCT.MASTER,DISP=OLD
//TRANIN   DD DSN=TRAN.SORTED,DISP=OLD
//ACCTOUT  DD DSN=ACCT.MASTER.NEW,DISP=NEW,
//            DCB=(RECFM=FB,LRECL=100)
//AUDIT    DD DSN=TRAN.AUDIT,DISP=NEW,
//            DCB=(RECFM=FB,LRECL=100)
//REJECTS  DD DSN=TRAN.POST.REJECTS,DISP=NEW,
//            DCB=(RECFM=FB,LRECL=100)
//POSTRPT  DD SYSOUT=*
