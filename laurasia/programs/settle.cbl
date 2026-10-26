       >>source format is free
*> settle: work out what each merchant is paid, and what the network is
*>
*> reads the priced presentments and the dispute adjustments, both in
*> merchant order. our own merchants, those our acquiring signed, numbered
*> lcs..., are paid their matched presentments, less their interchange,
*> plus or minus their adjustments; one that ends up owing is not paid,
*> and is listed for collection. every other merchant is another
*> acquirer's: what is due for them is settled with the network, all of it
*> in one sum.
*>
*> return codes: 0 all paid, 4 some owe.
*>
*> history
*> 2011-04-04 rk  lcs-0102 original version.
*> 2015-03-23 mt  lcs-0302 dispute adjustments.
*> 2024-06-03 ah  lcs-0561 our own merchants paid, the rest through the network.
identification division.
program-id. settle.

environment division.
input-output section.
file-control.
    select pres-in assign to "PRESIN"
        organization is line sequential.
    select adjust-in assign to "ADJIN"
        organization is line sequential.
    select payout-out assign to "PAYOUT"
        organization is line sequential.
    select network-out assign to "NETSETL"
        organization is line sequential.
    select report-out assign to "SETLRPT"
        organization is line sequential.

data division.
file section.
fd  pres-in.
    copy "presrec.cpy".
fd  adjust-in.
01  adjustment.
    05  ad-merchant              pic x(10).
    05  ad-amount                pic s9(9)v99 sign leading separate.
    05  ad-case-no               pic 9(8).
    05  ad-event                 pic x(2).
fd  payout-out.
01  payout.
    05  pa-merchant              pic x(10).
    05  pa-gross                 pic s9(11)v99 sign leading separate.
    05  pa-interchange           pic s9(9)v99 sign leading separate.
    05  pa-adjustments           pic s9(11)v99 sign leading separate.
    05  pa-net                   pic s9(11)v99 sign leading separate.
    05  pa-status                pic x.
        88  pa-paid              value "P".
        88  pa-owes              value "C".
fd  network-out.
01  network-settlement.
    05  ns-gross                 pic s9(13)v99 sign leading separate.
    05  ns-interchange           pic s9(11)v99 sign leading separate.
    05  ns-adjustments           pic s9(13)v99 sign leading separate.
    05  ns-net                   pic s9(13)v99 sign leading separate.
fd  report-out.
01  report-line                  pic x(60).

working-storage section.
01  ws-pres-key                  pic x(10).
01  ws-adjust-key                pic x(10).
01  ws-merchant.
    05  ws-merchant-range        pic x(3).
        88  ws-our-merchant      value "LCS".
    05  filler                   pic x(7).
01  ws-merchant-totals.
    05  ws-gross                 pic s9(11)v99.
    05  ws-interchange           pic s9(9)v99.
    05  ws-adjustments           pic s9(11)v99.
01  ws-totals.
    05  ws-paid-count            pic 9(6) value zero.
    05  ws-owing-count           pic 9(6) value zero.
    05  ws-paid-total            pic s9(13)v99 value zero.
    05  ws-network-count         pic 9(6) value zero.
01  ws-network-totals.
    05  ws-network-gross         pic s9(13)v99 value zero.
    05  ws-network-interchange   pic s9(11)v99 value zero.
    05  ws-network-adjustments   pic s9(13)v99 value zero.
01  ws-report-amount.
    05  ra-label                 pic x(30).
    05  ra-value                 pic -(11)9.99.
01  ws-report-count.
    05  rc-label                 pic x(30).
    05  rc-value                 pic z(5)9.

procedure division.
main section.
    open input pres-in adjust-in
         output payout-out network-out report-out
    perform read-presentment
    perform read-adjustment
    perform settle-merchant until ws-pres-key = high-values and ws-adjust-key = high-values
    perform write-network-settlement
    perform write-report
    close pres-in adjust-in payout-out network-out report-out
    if ws-owing-count > zero
        stop run returning 4
    end-if
    stop run returning 0.

read-presentment section.
    read pres-in
        at end move high-values to ws-pres-key
        not at end move ps-merchant to ws-pres-key
    end-read.

read-adjustment section.
    read adjust-in
        at end move high-values to ws-adjust-key
        not at end move ad-merchant to ws-adjust-key
    end-read.

*> the next merchant: the lower of the two keys, and everything of both
*> files for it.
settle-merchant section.
    if ws-pres-key < ws-adjust-key
        move ws-pres-key to ws-merchant
    else
        move ws-adjust-key to ws-merchant
    end-if
    move zero to ws-gross ws-interchange ws-adjustments
    perform until ws-pres-key not = ws-merchant
        if ps-matched
            add ps-amount to ws-gross
            add ps-fee to ws-interchange
        end-if
        perform read-presentment
    end-perform
    perform until ws-adjust-key not = ws-merchant
        add ad-amount to ws-adjustments
        perform read-adjustment
    end-perform
    if ws-our-merchant
        perform write-payout
    else
        perform add-to-network
    end-if.
*> another acquirer's merchant: what is due goes into the network's sum.
add-to-network section.
    add 1 to ws-network-count
    add ws-gross to ws-network-gross
    add ws-interchange to ws-network-interchange
    add ws-adjustments to ws-network-adjustments.

write-payout section.
    move ws-merchant to pa-merchant
    move ws-gross to pa-gross
    move ws-interchange to pa-interchange
    move ws-adjustments to pa-adjustments
    compute pa-net = ws-gross - ws-interchange + ws-adjustments
    if pa-net < zero
        set pa-owes to true
        add 1 to ws-owing-count
    else
        set pa-paid to true
        add 1 to ws-paid-count
        add pa-net to ws-paid-total
    end-if
    write payout.

write-network-settlement section.
    move ws-network-gross to ns-gross
    move ws-network-interchange to ns-interchange
    move ws-network-adjustments to ns-adjustments
    compute ns-net = ws-network-gross - ws-network-interchange + ws-network-adjustments
    write network-settlement.
write-report section.
    move "SETTLE   MERCHANT SETTLEMENT" to report-line
    write report-line
    move "MERCHANTS PAID" to rc-label
    move ws-paid-count to rc-value
    write report-line from ws-report-count
    move "MERCHANTS OWING" to rc-label
    move ws-owing-count to rc-value
    write report-line from ws-report-count
    move "PAID" to ra-label
    move ws-paid-total to ra-value
    write report-line from ws-report-amount
    move "OTHER ACQUIRERS' MERCHANTS" to rc-label
    move ws-network-count to rc-value
    write report-line from ws-report-count
    move "OWED TO THE NETWORK" to ra-label
    move ns-net to ra-value
    write report-line from ws-report-amount.
