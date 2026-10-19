       >>source format is free
*> dispevt: an event in a dispute's life, from the cardholder's claim to
*> the network's ruling.
*>
*> history
*> 2015-03-23 mt  lcs-0302 original version.
*> 2019-11-04 mt  lcs-0455 pre-arbitration.
01  dispute-event.
    05  de-case-no               pic 9(8).
    05  de-event                 pic x(2).
        88  de-opened            value "OP".
        88  de-chargeback        value "CB".
        88  de-representment     value "RP".
        88  de-pre-arbitration   value "PA".
        88  de-merchant-accepts  value "AC".
        88  de-issuer-drops      value "DR".
        88  de-ruled-for-issuer  value "RI".
        88  de-ruled-for-merchant value "RM".
    05  de-merchant              pic x(10).
    05  de-amount                pic s9(9)v99 sign leading separate.
    05  de-date                  pic 9(8).
