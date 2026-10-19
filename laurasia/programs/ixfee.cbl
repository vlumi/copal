       >>source format is free
*> ixfee: work out each presentment's interchange fee
*>
*> for each matched presentment, finds the merchant group by the
*> merchant's category code, and the fee for the card product, the region
*> and the group in the fee table: a rate on the amount and a fixed fee.
*> a combination the schedule has no fee for gets the default fee. a
*> presentment not matched to an approval pays no fee here; it is not
*> settled.
*>
*> return codes: 0 done, 4 some default fees, 16 no fee table.
*>
*> history
*> 2012-01-09 ah  lcs-0140 original version.
*> 2017-08-28 ah  lcs-0398 travel and fuel as merchant groups of their own.
identification division.
program-id. ixfee.

environment division.
input-output section.
file-control.
    select pres-in assign to "PRESIN"
        organization is line sequential.
    select fee-table assign to "FEETAB"
        organization is relative
        access mode is random
        relative key is ws-slot
        file status is ws-fee-status.
    select pres-out assign to "PRESOUT"
        organization is line sequential.
    select report-out assign to "FEERPT"
        organization is line sequential.

data division.
file section.
fd  pres-in.
01  pres-in-record               pic x(76).
fd  fee-table.
    copy "feerec.cpy".
fd  pres-out.
    copy "presrec.cpy".
fd  report-out.
01  report-line                  pic x(60).

working-storage section.
78  default-rate                 value 0.0200.
78  default-fixed-cents          value 25.
01  ws-slot                      pic 9(4) comp-5.
01  ws-fee-status                pic xx.
01  ws-group                     pic 9.
    88  group-grocery            value 1.
    88  group-restaurants        value 2.
    88  group-travel             value 3.
    88  group-fuel               value 4.
    88  group-other              value 5.
01  ws-at-end                    pic x value "N".
    88  presentments-at-end      value "Y".
01  ws-totals.
    05  ws-fees                  pic s9(9)v99 value zero.
    05  ws-default-count         pic 9(8) value zero.
01  ws-report-amount.
    05  ra-label                 pic x(30).
    05  ra-value                 pic -(9)9.99.
01  ws-report-count.
    05  rc-label                 pic x(30).
    05  rc-value                 pic z(7)9.

procedure division.
main section.
    open input pres-in fee-table
         output pres-out report-out
    if ws-fee-status not = "00"
        stop run returning 16
    end-if
    perform until presentments-at-end
        read pres-in into pres-record
            at end set presentments-at-end to true
            not at end perform price-presentment
        end-read
    end-perform
    perform write-report
    close pres-in fee-table pres-out report-out
    if ws-default-count > zero
        stop run returning 4
    end-if
    stop run returning 0.

price-presentment section.
    move zero to ps-fee
    if ps-matched
        perform find-group
        call "feeslot" using ps-product ps-region ws-group returning ws-slot
        read fee-table
            invalid key
                add 1 to ws-default-count
                compute ps-fee rounded = ps-amount * default-rate + default-fixed-cents / 100
            not invalid key
                compute ps-fee rounded = ps-amount * fs-rate + fs-fixed-cents / 100
        end-read
        add ps-fee to ws-fees
    end-if
    write pres-record.

*> the merchant group by the merchant category code.
find-group section.
    evaluate ps-mcc
        when 5411 thru 5499
            set group-grocery to true
        when 5811 thru 5814
            set group-restaurants to true
        when 3000 thru 3999
        when 4111 thru 4789
            set group-travel to true
        when 5541
        when 5542
        when 5983
            set group-fuel to true
        when other
            set group-other to true
    end-evaluate.

write-report section.
    move "IXFEE    INTERCHANGE" to report-line
    write report-line
    move "INTERCHANGE FEES" to ra-label
    move ws-fees to ra-value
    write report-line from ws-report-amount
    move "DEFAULT FEES" to rc-label
    move ws-default-count to rc-value
    write report-line from ws-report-count.
