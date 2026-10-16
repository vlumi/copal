//DAILY    JOB
//*----------------------------------------------------------------*
//* DAILY BATCH: VALIDATE THE DAY'S TRANSACTIONS                   *
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
