       >>source format is free
*> cardacct: a card's account, indexed by card number: its credit limit,
*> its balance, and the balance its last statement showed.
*>
*> history
*> 2011-04-04 rk  lcs-0102 original version.
*> 2013-06-03 rk  lcs-0219 billing cycle day and paid in full flag.
01  card-account.
    05  ca-card-no               pic x(16).
    05  ca-account-no            pic x(10).
    05  ca-product               pic x(2).
        88  ca-classic           value "CL".
        88  ca-gold              value "GD".
        88  ca-business          value "BZ".
    05  ca-status                pic x.
        88  ca-active            value "A".
        88  ca-blocked           value "B".
        88  ca-closed            value "C".
    05  ca-limit                 pic 9(7)v99.
    05  ca-balance               pic s9(7)v99 sign leading separate.
    05  ca-statement-balance     pic s9(7)v99 sign leading separate.
    05  ca-cycle-day             pic 99.
    05  ca-paid-in-full          pic x.
        88  ca-was-paid-in-full  value "Y".
