       >>source format is free
*> intcalc: a month's interest on a statement balance, in cents, as the
*> return code: a twelfth of the yearly rate of the card's product.
*>
*> history
*> 2013-06-03 rk  lcs-0219 original version.
*> 2021-03-01 mt  lcs-0481 the rate by product.
identification division.
program-id. intcalc.

data division.
working-storage section.
01  ws-yearly-rate               pic v9(4).
01  ws-cents                     pic s9(9) comp-5.
linkage section.
01  lk-product                   pic x(2).
    88  lk-classic               value "CL".
    88  lk-gold                  value "GD".
    88  lk-business              value "BZ".
01  lk-balance                   pic s9(7)v99 sign leading separate.

procedure division using lk-product lk-balance.
main section.
    evaluate true
        when lk-gold
            move 0.1690 to ws-yearly-rate
        when lk-business
            move 0.1490 to ws-yearly-rate
        when other
            move 0.1990 to ws-yearly-rate
    end-evaluate
    compute ws-cents rounded = lk-balance * ws-yearly-rate / 12 * 100
    move ws-cents to return-code
    goback.
