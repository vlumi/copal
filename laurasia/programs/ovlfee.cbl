       >>source format is free
*> ovlfee: the over-limit fee of a consumer card: a tenth of what the
*> balance is over the limit, at least 5.00 and at most 35.00; a gold
*> card's first 50.00 over is free.
*>
*> history
*> 2013-06-03 rk  lcs-0219 original version.
identification division.
program-id. ovlfee.

data division.
working-storage section.
78  least-fee                    value 5.00.
78  most-fee                     value 35.00.
78  gold-allowance               value 50.00.
linkage section.
01  lk-product                   pic x(2).
    88  lk-gold                  value "GD".
01  lk-over-amount               pic s9(7)v99.
01  lk-fee                       pic s9(7)v99.

procedure division using lk-product lk-over-amount lk-fee.
main section.
    if lk-gold and lk-over-amount not > gold-allowance
        move zero to lk-fee
        goback
    end-if
    compute lk-fee rounded = lk-over-amount / 10
    evaluate true
        when lk-fee < least-fee
            move least-fee to lk-fee
        when lk-fee > most-fee
            move most-fee to lk-fee
    end-evaluate
    goback.
