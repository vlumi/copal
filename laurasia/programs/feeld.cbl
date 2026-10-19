       >>source format is free
*> feeld: load the interchange fee schedule into the fee table
*>
*> reads the fee schedule, a line for each card product, region and
*> merchant group, and writes each fee into the fee table's slot for it,
*> where ixfee finds it without a search.
*>
*> return codes: 0 loaded, 8 a fee for no slot, 16 cannot write.
*>
*> history
*> 2012-01-09 ah  lcs-0140 original version.
identification division.
program-id. feeld.

environment division.
input-output section.
file-control.
    select fee-schedule assign to "FEESCHED"
        organization is line sequential.
    select fee-table assign to "FEETAB"
        organization is relative
        access mode is random
        relative key is ws-slot
        file status is ws-fee-status.

data division.
file section.
fd  fee-schedule.
    copy "feedef.cpy".
fd  fee-table.
    copy "feerec.cpy".

working-storage section.
01  ws-slot                      pic 9(4) comp-5.
01  ws-fee-status                pic xx.
01  ws-at-end                    pic x value "N".
    88  schedule-at-end          value "Y".
01  ws-result                    pic 99 value zero.

procedure division.
main section.
    open input fee-schedule output fee-table
    if ws-fee-status not = "00"
        stop run returning 16
    end-if
    perform until schedule-at-end
        read fee-schedule
            at end set schedule-at-end to true
            not at end perform load-fee
        end-read
    end-perform
    close fee-schedule fee-table
    stop run returning ws-result.

load-fee section.
    call "feeslot" using fd-product fd-region fd-group returning ws-slot
    if ws-slot = zero
        move 8 to ws-result
        exit section
    end-if
    move fd-rate to fs-rate
    move fd-fixed to fs-fixed-cents
    write fee-slot-record
        invalid key move 8 to ws-result
    end-write.
