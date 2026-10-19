       >>source format is free
*> ovlbiz: the over-limit fee of a business card: none; the company is
*> billed for the excess at the cycle.
*>
*> history
*> 2020-02-10 mt  lcs-0466 original version.
identification division.
program-id. ovlbiz.

data division.
linkage section.
01  lk-product                   pic x(2).
01  lk-over-amount               pic s9(7)v99.
01  lk-fee                       pic s9(7)v99.

procedure division using lk-product lk-over-amount lk-fee.
main section.
    move zero to lk-fee
    goback.
