//ONLINE   JOB
//*================================================================*
//* THE ONLINE DAY: AUTHSRV AUTHORIZES CARD PURCHASES, THEIR       *
//* REVERSALS, ATM WITHDRAWALS AND BALANCE ENQUIRIES, ONE REQUEST  *
//* AT A TIME, FROM THE OPENING OF THE BUSINESS DAY TO THE BATCH   *
//* WINDOW. THE REQUESTS ARE THE JOB'S IN-STREAM DATA, WHICH THE   *
//* ONLINE SERVER SUPPLIES FROM THE TERMINALS AS THEY COME, AND    *
//* THE RESPONSES GO TO THE STEP'S OUTPUT, BACK TO THE SERVER.     *
//* WHEN THE REQUESTS END, THE DAY'S HOLDS ARE EXTRACTED AS THE    *
//* CARD FILE, FOR THE NIGHT'S BATCH TO POST.                      *
//*                                                                *
//* READS   BANK.CONTROL(0)                                        *
//* UPDATES ACCT.ONLINE(0)   CARD.ONLINE(0)   HOLD.ONLINE(0)       *
//* WRITES  TRAN.CARD(+1)                                          *
//*                                                                *
//* THE NIGHT'S BATCH LOADS THE ONLINE FILES AND STARTS THE HOLD   *
//* FILE, SO THE JOB CAN BE RUN AGAIN DURING THE DAY.              *
//*================================================================*
//AUTH     EXEC PGM=AUTHSRV
//BANKCTL  DD DSN=BANK.CONTROL(0),DISP=SHR
//ACCTONL  DD DSN=ACCT.ONLINE(0),DISP=OLD
//CARDONL  DD DSN=CARD.ONLINE(0),DISP=OLD
//HOLDS    DD DSN=HOLD.ONLINE(0),DISP=OLD
//AUTHRPT  DD SYSOUT=*
//SYSIN    DD *
//*----------------------------------------------------------------*
//* EXTRACT THE DAY'S HOLDS AS TRANSACTIONS FOR POSTING            *
//*----------------------------------------------------------------*
//EXTRACT  EXEC PGM=HOLDEXT,COND=(4,LT)
//BANKCTL  DD DSN=BANK.CONTROL(0),DISP=SHR
//HOLDS    DD DSN=HOLD.ONLINE(0),DISP=SHR
//TRANOUT  DD DSN=TRAN.CARD(+1),DISP=NEW,
//            DCB=(RECFM=FB,LRECL=80)
//HOLDRPT  DD SYSOUT=*
