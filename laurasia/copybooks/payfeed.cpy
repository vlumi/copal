       >>source format is free
*> payfeed: the card system's daily feed, a record a line: a cardholder's
*> payment, with the bank's reference, or a change of a card's status.
*>
*> history
*> 2022-09-05 ah  lcs-0497 original version.
01  feed-record.
    05  pf-card-no               pic x(16).
    05  pf-kind                  pic x.
        88  pf-payment           value "P".
        88  pf-status-change     value "S".
    05  pf-amount                pic s9(7)v99 sign leading separate.
    05  pf-status                pic x.
    05  pf-reference             pic x(6).
    05  pf-date                  pic 9(8).
