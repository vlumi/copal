       >>source format is free
*> presrec: a presentment on its way through the night, a record a line.
*> clrin fills in where it came from, match its authorization's outcome,
*> ixfee the interchange fee.
*>
*> history
*> 2011-04-04 rk  lcs-0102 original version.
*> 2014-02-17 rk  lcs-0251 match status and interchange fee.
01  pres-record.
    05  ps-card-no               pic x(16).
    05  ps-amount                pic s9(9)v99 sign leading separate.
    05  ps-auth-code             pic x(6).
    05  ps-merchant              pic x(10).
    05  ps-mcc                   pic 9(4).
    05  ps-region                pic x.
        88  ps-domestic          value "D".
        88  ps-international     value "I".
    05  ps-product               pic x(2).
        88  ps-classic           value "CL".
        88  ps-gold              value "GD".
        88  ps-business          value "BZ".
    05  ps-txn-date              pic 9(8).
    05  ps-batch-no              pic 9(6).
    05  ps-match                 pic x.
        88  ps-unmatched         value " ".
        88  ps-matched           value "M".
        88  ps-declined          value "D".
        88  ps-no-auth           value "U".
        88  ps-over-auth         value "X".
    05  ps-fee                   pic s9(7)v99 sign leading separate.
