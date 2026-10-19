       >>source format is free
*> dispcase: a cardholder's dispute of a presentment, a case in the
*> dispute master, with where it stands.
*>
*> history
*> 2015-03-23 mt  lcs-0302 original version.
*> 2019-11-04 mt  lcs-0455 pre-arbitration.
01  dispute-case.
    05  dc-case-no               pic 9(8).
    05  dc-status                pic x(2).
        88  dc-open              value "OP".
        88  dc-charged-back      value "CB".
        88  dc-represented       value "RP".
        88  dc-pre-arbitration   value "PA".
        88  dc-won               value "WN".
        88  dc-lost              value "LS".
    05  dc-merchant              pic x(10).
    05  dc-amount                pic s9(9)v99 sign leading separate.
    05  dc-opened                pic 9(8).
    05  dc-updated               pic 9(8).
