       >>source format is free
*> cycle: close the billing cycle of the cards whose cycle day it is
*>
*> for each open or blocked card whose cycle day is the run's day of the
*> month: interest on the last statement's balance, unless it was paid in
*> full, as the grace period allows; the minimum payment; and the
*> statement, with the balance it closes on, which the next cycle charges
*> interest on.
*>
*> return codes: 0 done, 16 cannot update the accounts.
*>
*> history
*> 2013-06-03 rk  lcs-0219 original version.
*> 2021-03-01 mt  lcs-0481 interest by the card's product.
identification division.
program-id. cycle.

environment division.
input-output section.
file-control.
    select accounts assign to "ACCOUNTS"
        organization is indexed
        access mode is sequential
        record key is ca-card-no
        file status is ws-accounts-status.
    select postings-out assign to "POSTOUT"
        organization is line sequential.
    select statements-out assign to "STMTOUT"
        organization is line sequential.

data division.
file section.
fd  accounts.
    copy "cardacct.cpy".
fd  postings-out.
    copy "postrec.cpy".
fd  statements-out.
01  statement-line               pic x(70).

working-storage section.
01  ws-run-date.
    05  ws-run-year-month        pic 9(6).
    05  ws-run-day               pic 99.
01  ws-accounts-status           pic xx.
    88  accounts-at-end          value "10".
01  ws-interest-cents            pic s9(9) comp-5.
01  ws-interest                  pic s9(7)v99.
01  ws-minimum                   pic s9(7)v99.
01  ws-closed-count              pic 9(6) value zero.
01  ws-statement-head.
    05  sh-card-no               pic x(16).
    05  filler                   pic x value space.
    05  sh-account-no            pic x(10).
    05  filler                   pic x(4) value " ON ".
    05  sh-date                  pic 9(8).
01  ws-statement-amount.
    05  filler                   pic x(4) value spaces.
    05  sa-label                 pic x(30).
    05  sa-value                 pic -(7)9.99.

procedure division.
main section.
    accept ws-run-date from command-line
    open i-o accounts
         output postings-out statements-out
    if ws-accounts-status not = "00"
        stop run returning 16
    end-if
    read accounts next record
    perform until accounts-at-end
        if ca-cycle-day = ws-run-day and not ca-closed
            perform close-cycle
        end-if
        read accounts next record
    end-perform
    close accounts postings-out statements-out
    stop run returning 0.

close-cycle section.
    move zero to ws-interest
    if not ca-was-paid-in-full and ca-statement-balance > zero
        call "intcalc" using ca-product ca-statement-balance returning ws-interest-cents
        compute ws-interest = ws-interest-cents / 100
    end-if
    if ws-interest > zero
        add ws-interest to ca-balance
        perform post-interest
    end-if
    call "minpay" using ca-balance ws-minimum
    perform write-statement
    move ca-balance to ca-statement-balance
    rewrite card-account
        invalid key stop run returning 16
    end-rewrite
    add 1 to ws-closed-count.

post-interest section.
    move ca-card-no to po-card-no
    move ca-account-no to po-account-no
    set po-interest to true
    move ws-interest to po-amount
    move spaces to po-reference
    move ws-run-date to po-date
    write posting.

write-statement section.
    move ca-card-no to sh-card-no
    move ca-account-no to sh-account-no
    move ws-run-date to sh-date
    write statement-line from ws-statement-head
    move "PREVIOUS BALANCE" to sa-label
    move ca-statement-balance to sa-value
    write statement-line from ws-statement-amount
    move "INTEREST" to sa-label
    move ws-interest to sa-value
    write statement-line from ws-statement-amount
    move "NEW BALANCE" to sa-label
    move ca-balance to sa-value
    write statement-line from ws-statement-amount
    move "MINIMUM PAYMENT" to sa-label
    move ws-minimum to sa-value
    write statement-line from ws-statement-amount.

*> minpay: the minimum payment of a balance: the whole of a small one,
*> otherwise a fiftieth of it, but at least the floor.
identification division.
program-id. minpay.

data division.
working-storage section.
78  payment-floor                value 25.00.
linkage section.
01  lk-balance                   pic s9(7)v99 sign leading separate.
01  lk-minimum                   pic s9(7)v99.

procedure division using lk-balance lk-minimum.
main section.
    evaluate true
        when lk-balance not > zero
            move zero to lk-minimum
        when lk-balance not > payment-floor
            move lk-balance to lk-minimum
        when other
            compute lk-minimum rounded = lk-balance / 50
            if lk-minimum < payment-floor
                move payment-floor to lk-minimum
            end-if
    end-evaluate
    goback.
end program minpay.
end program cycle.
