       >>source format is free
*> dispute: move the disputes on by the day's events
*>
*> applies the day's dispute events, in case order, to the dispute master,
*> and writes the new master. a dispute opens when the cardholder claims;
*> the chargeback takes the amount back from the merchant; the merchant
*> may accept it, or represent, which gives the amount back; the issuer
*> may then drop the dispute, or take it to pre-arbitration, where the
*> network rules for one side. an event a case cannot take where it stands
*> is rejected. each move of money is an adjustment for the merchant's
*> settlement.
*>
*> return codes: 0 all applied, 4 some rejected, 16 cases out of order.
*>
*> history
*> 2015-03-23 mt  lcs-0302 original version.
*> 2019-11-04 mt  lcs-0455 pre-arbitration.
identification division.
program-id. dispute.

environment division.
input-output section.
file-control.
    select master-in assign to "DISPIN"
        organization is line sequential.
    select events-in assign to "DISPEVT"
        organization is line sequential.
    select master-out assign to "DISPOUT"
        organization is line sequential.
    select adjust-out assign to "ADJOUT"
        organization is line sequential.
    select reject-out assign to "DISPREJ"
        organization is line sequential.

data division.
file section.
fd  master-in.
01  master-in-record             pic x(48).
fd  events-in.
    copy "dispevt.cpy".
fd  master-out.
    copy "dispcase.cpy".
fd  adjust-out.
01  adjustment.
    05  ad-merchant              pic x(10).
    05  ad-amount                pic s9(9)v99 sign leading separate.
    05  ad-case-no               pic 9(8).
    05  ad-event                 pic x(2).
fd  reject-out.
01  reject-line.
    05  rj-reason                pic x(4).
    05  filler                   pic x.
    05  rj-event                 pic x(40).

working-storage section.
01  ws-master-key                pic 9(8).
01  ws-event-key                 pic 9(8).
01  ws-held-case                 pic x(48).
01  ws-holding                   pic x value "N".
    88  holding-case             value "Y".
    88  holding-nothing          value "N".
01  ws-result                    pic 99 value zero.

procedure division.
main section.
    open input master-in events-in
         output master-out adjust-out reject-out
    perform read-master
    perform read-event
    perform apply-next until ws-master-key = all "9" and ws-event-key = all "9"
    close master-in events-in master-out adjust-out reject-out
    stop run returning ws-result.

read-master section.
    read master-in
        at end move all "9" to ws-master-key
        not at end
            move master-in-record to dispute-case
            move dc-case-no to ws-master-key
    end-read.

read-event section.
    read events-in
        at end move all "9" to ws-event-key
        not at end move de-case-no to ws-event-key
    end-read.

*> the lower case number first: a case with no events today is copied,
*> the events of a case are applied to it in turn, and an event of no
*> case opens one, if it opens.
apply-next section.
    evaluate true
        when ws-master-key < ws-event-key
            move master-in-record to dispute-case
            write dispute-case
            perform read-master
        when ws-master-key = ws-event-key
            move master-in-record to ws-held-case
            set holding-case to true
            perform read-master
            perform apply-events
        when other
            set holding-nothing to true
            perform apply-events
    end-evaluate.

apply-events section.
    move ws-event-key to dc-case-no
    perform until ws-event-key not = dc-case-no
        if holding-case
            move ws-held-case to dispute-case
        end-if
        perform apply-event
        if holding-case
            move dispute-case to ws-held-case
        end-if
        perform read-event
    end-perform
    if holding-case
        move ws-held-case to dispute-case
        write dispute-case
    end-if.

*> the state machine of a dispute.
apply-event section.
    evaluate true also true
        when holding-nothing also de-opened
            move de-case-no to dc-case-no
            move de-merchant to dc-merchant
            move de-amount to dc-amount
            move de-date to dc-opened
            set dc-open to true
            set holding-case to true
        when holding-case also dc-open and de-chargeback
            set dc-charged-back to true
            compute ad-amount = 0 - dc-amount
            perform write-adjustment
        when holding-case also dc-charged-back and de-merchant-accepts
            set dc-won to true
        when holding-case also dc-charged-back and de-representment
            set dc-represented to true
            move dc-amount to ad-amount
            perform write-adjustment
        when holding-case also dc-represented and de-issuer-drops
            set dc-lost to true
        when holding-case also dc-represented and de-pre-arbitration
            set dc-pre-arbitration to true
        when holding-case also dc-pre-arbitration and de-ruled-for-issuer
            set dc-won to true
            compute ad-amount = 0 - dc-amount
            perform write-adjustment
        when holding-case also dc-pre-arbitration and de-ruled-for-merchant
            set dc-lost to true
        when other
            move spaces to reject-line
            move "STAT" to rj-reason
            move dispute-event to rj-event
            write reject-line
            move 4 to ws-result
            exit section
    end-evaluate
    move de-date to dc-updated.

write-adjustment section.
    move dc-merchant to ad-merchant
    move dc-case-no to ad-case-no
    move de-event to ad-event
    write adjustment.
