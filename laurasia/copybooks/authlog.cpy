       >>source format is free
*> authlog: the online side's authorizations of the day, a record a line,
*> the ones given in stand-in, when the issuer could not answer, marked.
*>
*> history
*> 2011-04-04 rk  lcs-0102 original version.
01  auth-entry.
    05  al-auth-code             pic x(6).
    05  al-card-no               pic x(16).
    05  al-amount                pic s9(9)v99 sign leading separate.
    05  al-decision              pic x.
        88  al-approved          value "A".
        88  al-declined          value "D".
    05  al-stand-in              pic x.
        88  al-in-stand-in       value "Y".
    05  al-date                  pic 9(8).
