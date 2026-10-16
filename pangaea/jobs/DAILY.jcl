//DAILY    JOB
//*================================================================*
//* DAILY BATCH: VALIDATE THE DAY'S TRANSACTIONS FROM THE BRANCHES *
//* AND THE CARDS, ACCRUE INTEREST AND POST WHAT A CLOSED MONTH    *
//* OWES, THEN SORT AND POST THE DAY'S TRANSACTIONS, THE TRANSFER  *
//* CREDITS AND RETURNS THEY LEAD TO, APPLY THE DAY'S ACCOUNT      *
//* MAINTENANCE, LOAD THE ONLINE ACCOUNT FILE, MOVE TO THE NEXT    *
//* BUSINESS DAY, AND LOAD THE ONLINE CARD FILE FOR IT.            *
//*                                                                *
//* READS   BANK.CONTROL(0)  ACCT.MASTER(0)  TRAN.DAILY(0)         *
//*         TRAN.CARD(0)     ACCT.ACCRUAL(0)  ACCT.MAINT(0)        *
//*         CARD.MASTER(0)   BANK.RATES                            *
//* WRITES  BANK.CONTROL(+1) ACCT.MASTER(+1) ACCT.ACCRUAL(+1)      *
//*         TRAN.AUDIT(+1)   TRAN.REJECTS(+1) MAINT.REJECTS(+1)    *
//*         ACCT.ONLINE(+1)  CARD.ONLINE(+1)  HOLD.ONLINE(+1)      *
//*                                                                *
//* A STEP RUNS ONLY IF NO EARLIER STEP ENDED WITH RC ABOVE 4, SO  *
//* A RUN THAT FAILS LEAVES THE BUSINESS DATE WHERE IT WAS.        *
//*================================================================*
//*----------------------------------------------------------------*
//* VALIDATE THE BRANCHES' TRANSACTION FILE                        *
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
//* VALIDATE THE CARD FILE, THE ONLINE DAY'S AUTHORIZED CARD AND   *
//* ATM TRANSACTIONS                                               *
//*----------------------------------------------------------------*
//VALIDC   EXEC PGM=TXVALID,COND=(4,LT)
//BANKCTL  DD DSN=BANK.CONTROL(0),DISP=SHR
//TRANIN   DD DSN=TRAN.CARD(0),DISP=SHR
//TRANOUT  DD DSN=&&CARDVAL,DISP=NEW,
//            DCB=(RECFM=FB,LRECL=80)
//REJECTS  DD DSN=TRAN.REJECTS(+1),DISP=MOD
//VALRPT   DD SYSOUT=*
//*----------------------------------------------------------------*
//* ACCRUE INTEREST UP TO TODAY ON THE OPENING BALANCES. WHEN A    *
//* MONTH HAS CLOSED, POST ITS INTEREST AND FEES. THE CHARGES COME *
//* IN MASTER ORDER, SO THEY ARE POSTED WITHOUT A SORT.            *
//*----------------------------------------------------------------*
//ACCRUE   EXEC PGM=INTACCR,COND=(4,LT)
//BANKCTL  DD DSN=BANK.CONTROL(0),DISP=SHR
//RATES    DD DSN=BANK.RATES,DISP=SHR
//ACCTIN   DD DSN=ACCT.MASTER(0),DISP=SHR
//ACCRIN   DD DSN=ACCT.ACCRUAL(0),DISP=SHR
//ACCROUT  DD DSN=ACCT.ACCRUAL(+1),DISP=NEW,
//            DCB=(RECFM=FB,LRECL=40)
//ACCRRPT  DD SYSOUT=*
//CHARGE   EXEC PGM=MONCHG,COND=(4,LT)
//BANKCTL  DD DSN=BANK.CONTROL(0),DISP=SHR
//ACCTIN   DD DSN=ACCT.MASTER(0),DISP=SHR
//ACCRIN   DD DSN=ACCT.ACCRUAL(+1),DISP=SHR
//TRANOUT  DD DSN=&&MONTHLY,DISP=NEW,
//            DCB=(RECFM=FB,LRECL=80)
//MONRPT   DD SYSOUT=*
//POSTMO   EXEC PGM=TXPOST,COND=(4,LT)
//BANKCTL  DD DSN=BANK.CONTROL(0),DISP=SHR
//ACCTIN   DD DSN=ACCT.MASTER(0),DISP=SHR
//TRANIN   DD DSN=&&MONTHLY,DISP=OLD
//ACCTOUT  DD DSN=&&CHARGED,DISP=NEW,
//            DCB=(RECFM=FB,LRECL=100)
//AUDIT    DD DSN=TRAN.AUDIT(+1),DISP=MOD,
//            DCB=(RECFM=FB,LRECL=100)
//REJECTS  DD DSN=TRAN.REJECTS(+1),DISP=MOD
//LEGS     DD DUMMY
//POSTRPT  DD SYSOUT=*
//*----------------------------------------------------------------*
//* SORT THE DAY'S TRANSACTIONS, FROM THE BRANCHES AND THE CARDS,  *
//* INTO ACCOUNT ORDER AND POST THEM                               *
//*----------------------------------------------------------------*
//SORT     EXEC PGM=TXSORT,COND=(4,LT)
//TRANIN   DD DSN=&&VALID,DISP=OLD
//         DD DSN=&&CARDVAL,DISP=OLD
//TRANOUT  DD DSN=&&SORTED,DISP=NEW,
//            DCB=(RECFM=FB,LRECL=80)
//POST     EXEC PGM=TXPOST,COND=(4,LT)
//BANKCTL  DD DSN=BANK.CONTROL(0),DISP=SHR
//ACCTIN   DD DSN=&&CHARGED,DISP=OLD
//TRANIN   DD DSN=&&SORTED,DISP=OLD
//ACCTOUT  DD DSN=&&POSTED,DISP=NEW,
//            DCB=(RECFM=FB,LRECL=100)
//AUDIT    DD DSN=TRAN.AUDIT(+1),DISP=MOD
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
//* SORT AND POST THE RETURNS OF THE CREDITS THAT COULD NOT POST   *
//*----------------------------------------------------------------*
//SORTRT   EXEC PGM=TXSORT,COND=(4,LT)
//TRANIN   DD DSN=&&RETURNS,DISP=OLD
//TRANOUT  DD DSN=&&RETSORT,DISP=NEW,
//            DCB=(RECFM=FB,LRECL=80)
//POSTRT   EXEC PGM=TXPOST,COND=(4,LT)
//BANKCTL  DD DSN=BANK.CONTROL(0),DISP=SHR
//ACCTIN   DD DSN=&&CREDITED,DISP=OLD
//TRANIN   DD DSN=&&RETSORT,DISP=OLD
//ACCTOUT  DD DSN=&&RETURNED,DISP=NEW,
//            DCB=(RECFM=FB,LRECL=100)
//AUDIT    DD DSN=TRAN.AUDIT(+1),DISP=MOD
//REJECTS  DD DSN=TRAN.REJECTS(+1),DISP=MOD
//LEGS     DD DUMMY
//POSTRPT  DD SYSOUT=*
//*----------------------------------------------------------------*
//* APPLY THE DAY'S ACCOUNT MAINTENANCE, WRITING THE NEW MASTER    *
//*----------------------------------------------------------------*
//MAINT    EXEC PGM=ACCTMNT,COND=(4,LT)
//BANKCTL  DD DSN=BANK.CONTROL(0),DISP=SHR
//ACCTIN   DD DSN=&&RETURNED,DISP=OLD
//MAINTIN  DD DSN=ACCT.MAINT(0),DISP=SHR
//ACCTOUT  DD DSN=ACCT.MASTER(+1),DISP=NEW,
//            DCB=(RECFM=FB,LRECL=100)
//REJECTS  DD DSN=MAINT.REJECTS(+1),DISP=NEW,
//            DCB=(RECFM=FB,LRECL=100)
//MNTRPT   DD SYSOUT=*
//*----------------------------------------------------------------*
//* LOAD THE NEW MASTER INTO THE ONLINE ACCOUNT FILE, WHICH THE    *
//* CARD AND ATM AUTHORIZATIONS READ DURING THE NEXT BUSINESS DAY, *
//* AND START THAT DAY'S HOLD FILE                                 *
//*----------------------------------------------------------------*
//LOAD     EXEC PGM=ACCTLOAD,COND=(4,LT)
//BANKCTL  DD DSN=BANK.CONTROL(0),DISP=SHR
//ACCTIN   DD DSN=ACCT.MASTER(+1),DISP=SHR
//ACCTOUT  DD DSN=ACCT.ONLINE(+1),DISP=NEW,
//            RECORG=KS,LRECL=100
//HOLDOUT  DD DSN=HOLD.ONLINE(+1),DISP=NEW,
//            RECORG=KS,LRECL=100
//LOADRPT  DD SYSOUT=*
//*----------------------------------------------------------------*
//* MOVE TO THE NEXT BUSINESS DAY                                  *
//*----------------------------------------------------------------*
//NEXTDAY  EXEC PGM=BANKDATE,COND=(4,LT)
//BANKIN   DD DSN=BANK.CONTROL(0),DISP=SHR
//BANKOUT  DD DSN=BANK.CONTROL(+1),DISP=NEW,
//            DCB=(RECFM=FB,LRECL=80)
//*----------------------------------------------------------------*
//* LOAD THE CARDS INTO THE ONLINE CARD FILE FOR THE NEXT BUSINESS *
//* DAY, LEAVING OUT THOSE THAT HAVE EXPIRED BY THEN               *
//*----------------------------------------------------------------*
//CARDS    EXEC PGM=CARDLOAD,COND=(4,LT)
//BANKCTL  DD DSN=BANK.CONTROL(+1),DISP=SHR
//CARDIN   DD DSN=CARD.MASTER(0),DISP=SHR
//CARDOUT  DD DSN=CARD.ONLINE(+1),DISP=NEW,
//            RECORG=KS,LRECL=80
//LOADRPT  DD SYSOUT=*
