       >>source format is free
*> clrin: take in the networks' clearing file
*>
*> reads the night's clearing file, checks it is for the business date the
*> run is given, and checks each acquirer's batch against its trailer: the
*> presentments of a batch whose count and amount agree go on to matching,
*> those of a batch that does not are rejected whole. a presentment that
*> does not validate is rejected alone, but counts towards its batch. the
*> file trailer's totals are checked last.
*>
*> return codes: 0 all taken, 4 some rejected, 8 file totals differ,
*>               12 not the night's file, 16 cannot read or write.
*>
*> history
*> 2011-04-04 rk  lcs-0102 original version.
*> 2014-02-17 rk  lcs-0251 presentments held until their batch checks.
*> 2016-09-12 ah  lcs-0377 card product validated.
identification division.
program-id. clrin.

environment division.
input-output section.
file-control.
    select clearing-in assign to "CLRIN"
        organization is line sequential
        file status is ws-clearing-status.
    select pres-out assign to "PRESOUT"
        organization is line sequential
        file status is ws-pres-status.
    select reject-out assign to "CLRREJ"
        organization is line sequential
        file status is ws-reject-status.
    select report-out assign to "CLRRPT"
        organization is line sequential.

data division.
file section.
fd  clearing-in
    record varying in size from 14 to 61 characters depending on ws-line-length.
    copy "clrrec.cpy".
fd  pres-out.
    copy "presrec.cpy".
fd  reject-out.
01  reject-line.
    05  rj-reason                pic x(4).
    05  filler                   pic x.
    05  rj-batch-no              pic 9(6).
    05  filler                   pic x.
    05  rj-record                pic x(61).
fd  report-out.
01  report-line                  pic x(60).

working-storage section.
78  max-batch                    value 500.
01  ws-run-date                  pic 9(8).
01  ws-line-length               pic 9(4) comp-5.
01  ws-statuses.
    05  ws-clearing-status       pic xx.
        88  clearing-at-end      value "10".
    05  ws-pres-status           pic xx.
    05  ws-reject-status         pic xx.
01  ws-io-failed                 pic x value "N".
    88  io-failed                value "Y".
01  ws-batch.
    05  ws-batch-no              pic 9(6).
    05  ws-batch-count           pic 9(6) comp-5.
    05  ws-batch-amount          pic s9(13)v99 comp-3.
    05  ws-held-count            pic 9(4) comp-5.
    05  ws-held occurs 0 to 500 times depending on ws-held-count.
        10  ws-held-presentment  pic x(61).
01  ws-file-totals.
    05  ws-batches               pic 9(4) comp-5 value zero.
    05  ws-records               pic 9(8) comp-5 value zero.
    05  ws-amount                pic s9(15)v99 comp-3 value zero.
01  ws-counts.
    05  ws-taken                 pic 9(8) value zero.
    05  ws-rejected              pic 9(8) value zero.
01  ws-record-read               pic x(61).
01  ws-reason                    pic x(4).
    88  presentment-is-good      value spaces.
01  ws-result                    pic 99 value zero.
01  ws-report-count.
    05  rc-label                 pic x(30).
    05  rc-value                 pic z(7)9.

procedure division.
declaratives.
file-errors section.
    use after error procedure on clearing-in pres-out reject-out.
note-failure.
    set io-failed to true.
end declaratives.

main section.
    accept ws-run-date from command-line
    open input clearing-in
         output pres-out reject-out report-out
    if io-failed
        stop run returning 16
    end-if
    perform read-clearing
    if clearing-at-end or not is-file-header or fh-file-date not = ws-run-date
        move 12 to ws-result
    else
        perform read-clearing
        perform take-batch until clearing-at-end or is-file-trailer or io-failed
        perform check-file-trailer
    end-if
    perform write-report
    close clearing-in pres-out reject-out report-out
    if io-failed
        move 16 to ws-result
    end-if
    stop run returning ws-result.

read-clearing section.
    read clearing-in
        at end set clearing-at-end to true
    end-read.

*> a batch: its header, its presentments held until its trailer, which
*> lets them go on or rejects them all. the held presentments pass through
*> the file's record, so the trailer read is set aside meanwhile.
take-batch section.
    if not is-batch-header
        move "ORDR" to ws-reason
        perform reject-record
        perform read-clearing
        exit section
    end-if
    move bh-batch-no to ws-batch-no
    move zero to ws-batch-count ws-batch-amount ws-held-count
    add 1 to ws-batches
    perform read-clearing
    perform take-presentment until clearing-at-end or not is-presentment
    move clr-presentment to ws-record-read
    if is-batch-trailer and bt-count = ws-batch-count and bt-amount = ws-batch-amount
        perform release-batch
    else
        perform reject-batch
    end-if
    move ws-record-read to clr-presentment
    if is-batch-trailer
        perform read-clearing
    end-if.

take-presentment section.
    add 1 to ws-batch-count ws-records
    add pr-amount to ws-batch-amount ws-amount
    perform validate-presentment
    if presentment-is-good and ws-held-count < max-batch
        add 1 to ws-held-count
        move clr-presentment to ws-held-presentment(ws-held-count)
    else
        perform reject-record
    end-if
    perform read-clearing.

validate-presentment section.
    evaluate true
        when pr-card-no is not numeric
            move "CARD" to ws-reason
        when pr-amount not > zero
            move "AMNT" to ws-reason
        when pr-mcc is not numeric
            move "MCC " to ws-reason
        when pr-region not = "D" and pr-region not = "I"
            move "REGN" to ws-reason
        when pr-product not = "CL" and pr-product not = "GD" and pr-product not = "BZ"
            move "PROD" to ws-reason
        when other
            move spaces to ws-reason
    end-evaluate.

release-batch section.
    perform varying ws-line-length from 1 by 1 until ws-line-length > ws-held-count
        move ws-held-presentment(ws-line-length) to clr-presentment
        move spaces to pres-record
        move pr-card-no to ps-card-no
        move pr-amount to ps-amount
        move pr-auth-code to ps-auth-code
        move pr-merchant to ps-merchant
        move pr-mcc to ps-mcc
        move pr-region to ps-region
        move pr-product to ps-product
        move pr-txn-date to ps-txn-date
        move ws-batch-no to ps-batch-no
        move zero to ps-fee
        write pres-record
        add 1 to ws-taken
    end-perform.

reject-batch section.
    move "BTOT" to ws-reason
    perform varying ws-line-length from 1 by 1 until ws-line-length > ws-held-count
        move ws-held-presentment(ws-line-length) to clr-presentment
        perform reject-record
    end-perform.

reject-record section.
    move spaces to reject-line
    move ws-reason to rj-reason
    move ws-batch-no to rj-batch-no
    move clr-presentment to rj-record
    write reject-line
    add 1 to ws-rejected
    if ws-result < 4
        move 4 to ws-result
    end-if.

check-file-trailer section.
    if not is-file-trailer
        or ft-batch-count not = ws-batches
        or ft-record-count not = ws-records
        or ft-amount not = ws-amount
        move 8 to ws-result
    end-if.

write-report section.
    move spaces to report-line
    string "CLRIN    CLEARING FILE OF " ws-run-date delimited by size into report-line
    write report-line
    move "PRESENTMENTS TAKEN" to rc-label
    move ws-taken to rc-value
    write report-line from ws-report-count
    move "RECORDS REJECTED" to rc-label
    move ws-rejected to rc-value
    write report-line from ws-report-count
    move "BATCHES" to rc-label
    move ws-batches to rc-value
    write report-line from ws-report-count.
