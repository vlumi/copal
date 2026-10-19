       >>source format is free
*> post: post the night's matched presentments to the cards' accounts
*>
*> each presentment matched to an approval is a purchase on its card's
*> account. a purchase that takes the balance over the credit limit still
*> posts, as it was approved, but may cost an over-limit fee, which the fee
*> program the account's product names works out. a presentment for a card
*> with no account, or a blocked or closed one, is rejected for the
*> disputes desk; a presentment not matched is not posted.
*>
*> return codes: 0 all posted, 4 some rejected, 16 cannot update.
*>
*> history
*> 2011-04-04 rk  lcs-0102 original version.
*> 2013-06-03 rk  lcs-0219 over-limit fees.
*> 2020-02-10 mt  lcs-0466 the fee program by name, for the business cards'.
identification division.
program-id. post.

environment division.
input-output section.
file-control.
    select pres-in assign to "PRESIN"
        organization is line sequential.
    select accounts assign to "ACCOUNTS"
        organization is indexed
        access mode is random
        record key is ca-card-no
        file status is ws-accounts-status.
    select postings-out assign to "POSTOUT"
        organization is line sequential.
    select reject-out assign to "POSTREJ"
        organization is line sequential.
    select report-out assign to "POSTRPT"
        organization is line sequential.

data division.
file section.
fd  pres-in.
    copy "presrec.cpy".
fd  accounts.
    copy "cardacct.cpy".
fd  postings-out.
    copy "postrec.cpy".
fd  reject-out.
01  reject-line.
    05  rj-reason                pic x(4).
    05  filler                   pic x.
    05  rj-presentment           pic x(76).
fd  report-out.
01  report-line                  pic x(60).

working-storage section.
01  ws-run-date                  pic 9(8).
01  ws-accounts-status           pic xx.
01  ws-at-end                    pic x value "N".
    88  presentments-at-end      value "Y".
01  ws-fee-program               pic x(8).
01  ws-reject-reason             pic x(4).
01  ws-over-amount               pic s9(7)v99.
01  ws-fee                       pic s9(7)v99.
01  ws-totals.
    05  ws-posted                pic 9(8) value zero.
    05  ws-rejected              pic 9(8) value zero.
    05  ws-purchases             pic s9(9)v99 value zero.
    05  ws-fees                  pic s9(9)v99 value zero.
01  ws-report-count.
    05  rc-label                 pic x(30).
    05  rc-value                 pic z(7)9.
01  ws-report-amount.
    05  ra-label                 pic x(30).
    05  ra-value                 pic -(9)9.99.

procedure division.
main section.
    accept ws-run-date from command-line
    open input pres-in
         i-o accounts
         output postings-out reject-out report-out
    if ws-accounts-status not = "00"
        stop run returning 16
    end-if
    perform until presentments-at-end
        read pres-in
            at end set presentments-at-end to true
            not at end
                if ps-matched
                    perform post-presentment
                end-if
        end-read
    end-perform
    perform write-report
    close pres-in accounts postings-out reject-out report-out
    if ws-rejected > zero
        stop run returning 4
    end-if
    stop run returning 0.

post-presentment section.
    move ps-card-no to ca-card-no
    read accounts
        invalid key
            move "CARD" to rj-reason
            perform reject-presentment
            exit section
    end-read
    evaluate true
        when ca-blocked
            move "BLKD" to rj-reason
            perform reject-presentment
            exit section
        when ca-closed
            move "CLSD" to rj-reason
            perform reject-presentment
            exit section
    end-evaluate
    add ps-amount to ca-balance
    move "P" to po-kind
    move ps-amount to po-amount
    perform write-posting
    add ps-amount to ws-purchases
    if ca-balance > ca-limit
        perform charge-over-limit
    end-if
    rewrite card-account
        invalid key stop run returning 16
    end-rewrite.

*> the fee for going over the limit, by the program the card's product
*> names; one that charges none returns zero.
charge-over-limit section.
    evaluate true
        when ca-business
            move "ovlbiz" to ws-fee-program
        when other
            move "ovlfee" to ws-fee-program
    end-evaluate
    compute ws-over-amount = ca-balance - ca-limit
    call ws-fee-program using ca-product ws-over-amount ws-fee
    if ws-fee > zero
        add ws-fee to ca-balance ws-fees
        move "F" to po-kind
        move ws-fee to po-amount
        perform write-posting
    end-if.

write-posting section.
    move ca-card-no to po-card-no
    move ca-account-no to po-account-no
    move ps-auth-code to po-reference
    move ws-run-date to po-date
    write posting
    add 1 to ws-posted.

reject-presentment section.
    move rj-reason to ws-reject-reason
    move spaces to reject-line
    move ws-reject-reason to rj-reason
    move pres-record to rj-presentment
    write reject-line
    add 1 to ws-rejected.

write-report section.
    move "POST     CARD ACCOUNTS" to report-line
    write report-line
    move "POSTINGS" to rc-label
    move ws-posted to rc-value
    write report-line from ws-report-count
    move "PRESENTMENTS REJECTED" to rc-label
    move ws-rejected to rc-value
    write report-line from ws-report-count
    move "PURCHASES" to ra-label
    move ws-purchases to ra-value
    write report-line from ws-report-amount
    move "OVER-LIMIT FEES" to ra-label
    move ws-fees to ra-value
    write report-line from ws-report-amount.
