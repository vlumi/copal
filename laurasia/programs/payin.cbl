       >>source format is free
*> payin: take in the card system's feed
*>
*> applies the day's feed from the card system to the cards' accounts: a
*> payment comes off the card's balance, and marks the last statement
*> paid in full when it covers it or clears the balance, which spares the
*> card interest at the cycle; a status change sets the card active,
*> blocked or closed. a record for a card with no account, a payment of
*> nothing, a status that is none of the three, or a record of no known
*> kind is rejected for the card system.
*>
*> return codes: 0 all applied, 4 some rejected, 16 cannot update.
*>
*> history
*> 2022-09-05 ah  lcs-0497 original version.
identification division.
program-id. payin.

environment division.
input-output section.
file-control.
    select feed-in assign to "FEEDIN"
        organization is line sequential.
    select accounts assign to "ACCOUNTS"
        organization is indexed
        access mode is random
        record key is ca-card-no
        file status is ws-accounts-status.
    select postings-out assign to "POSTOUT"
        organization is line sequential.
    select reject-out assign to "PAYREJ"
        organization is line sequential.
    select report-out assign to "PAYRPT"
        organization is line sequential.

data division.
file section.
fd  feed-in.
    copy "payfeed.cpy".
fd  accounts.
    copy "cardacct.cpy".
fd  postings-out.
    copy "postrec.cpy".
fd  reject-out.
01  reject-line.
    05  rj-reason                pic x(4).
    05  filler                   pic x.
    05  rj-record                pic x(42).
fd  report-out.
01  report-line                  pic x(60).

working-storage section.
01  ws-run-date                  pic 9(8).
01  ws-accounts-status           pic xx.
01  ws-at-end                    pic x value "N".
    88  feed-at-end              value "Y".
01  ws-reason                    pic x(4).
    88  record-is-good           value spaces.
01  ws-totals.
    05  ws-payments              pic 9(8) value zero.
    05  ws-status-changes        pic 9(8) value zero.
    05  ws-rejected              pic 9(8) value zero.
    05  ws-paid                  pic s9(9)v99 value zero.
01  ws-report-count.
    05  rc-label                 pic x(30).
    05  rc-value                 pic z(7)9.
01  ws-report-amount.
    05  ra-label                 pic x(30).
    05  ra-value                 pic -(9)9.99.

procedure division.
main section.
    accept ws-run-date from command-line
    open input feed-in
         i-o accounts
         output postings-out reject-out report-out
    if ws-accounts-status not = "00"
        stop run returning 16
    end-if
    perform until feed-at-end
        read feed-in
            at end set feed-at-end to true
            not at end perform apply-record
        end-read
    end-perform
    perform write-report
    close feed-in accounts postings-out reject-out report-out
    if ws-rejected > zero
        stop run returning 4
    end-if
    stop run returning 0.

apply-record section.
    perform validate-record
    if not record-is-good
        perform reject-record
        exit section
    end-if
    move pf-card-no to ca-card-no
    read accounts
        invalid key
            move "CARD" to ws-reason
            perform reject-record
            exit section
    end-read
    if pf-payment
        perform apply-payment
    else
        perform apply-status
    end-if
    rewrite card-account
        invalid key stop run returning 16
    end-rewrite.

validate-record section.
    evaluate true
        when not pf-payment and not pf-status-change
            move "KIND" to ws-reason
        when pf-payment and pf-amount not > zero
            move "AMNT" to ws-reason
        when pf-status-change and pf-status not = "A" and pf-status not = "B" and pf-status not = "C"
            move "STAT" to ws-reason
        when other
            move spaces to ws-reason
    end-evaluate.

*> the payment off the balance; one that covers the last statement, or
*> clears the balance, pays the statement in full.
apply-payment section.
    subtract pf-amount from ca-balance
    if ca-balance not > zero or pf-amount not < ca-statement-balance
        set ca-was-paid-in-full to true
    end-if
    move ca-card-no to po-card-no
    move ca-account-no to po-account-no
    set po-payment to true
    compute po-amount = 0 - pf-amount
    move pf-reference to po-reference
    move ws-run-date to po-date
    write posting
    add 1 to ws-payments
    add pf-amount to ws-paid.

apply-status section.
    move pf-status to ca-status
    add 1 to ws-status-changes.

reject-record section.
    move spaces to reject-line
    move ws-reason to rj-reason
    move feed-record to rj-record
    write reject-line
    add 1 to ws-rejected.

write-report section.
    move "PAYIN    CARD SYSTEM FEED" to report-line
    write report-line
    move "PAYMENTS" to rc-label
    move ws-payments to rc-value
    write report-line from ws-report-count
    move "STATUS CHANGES" to rc-label
    move ws-status-changes to rc-value
    write report-line from ws-report-count
    move "RECORDS REJECTED" to rc-label
    move ws-rejected to rc-value
    write report-line from ws-report-count
    move "PAID" to ra-label
    move ws-paid to ra-value
    write report-line from ws-report-amount.
