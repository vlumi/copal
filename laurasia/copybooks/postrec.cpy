       >>source format is free
*> postrec: a posting to a card's account, a record a line.
*>
*> history
*> 2011-04-04 rk  lcs-0102 original version.
01  posting.
    05  po-card-no               pic x(16).
    05  po-account-no            pic x(10).
    05  po-kind                  pic x.
        88  po-purchase          value "P".
        88  po-over-limit-fee    value "F".
        88  po-interest          value "I".
    05  po-amount                pic s9(7)v99 sign leading separate.
    05  po-reference             pic x(6).
    05  po-date                  pic 9(8).
