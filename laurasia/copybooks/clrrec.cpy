       >>source format is free
*> clrrec: the networks' clearing file, a record a line, of five layouts
*> told apart by their first two characters. A file header, then for each
*> acquirer's batch a batch header, its presentments and a batch trailer
*> with their count and amount, then a file trailer with the batches',
*> the presentments' and the amounts' totals.
*>
*> history
*> 2011-04-04 rk  lcs-0102 original version.
*> 2016-09-12 ah  lcs-0377 card product on the presentment.
01  clr-file-header.
    05  clr-record-type          pic x(2).
        88  is-file-header       value "FH".
        88  is-batch-header      value "BH".
        88  is-presentment       value "PR".
        88  is-batch-trailer     value "BT".
        88  is-file-trailer      value "FT".
    05  fh-network               pic x(4).
    05  fh-file-date             pic 9(8).
    05  fh-file-seq              pic 9(4).
01  clr-batch-header.
    05  filler                   pic x(2).
    05  bh-batch-no              pic 9(6).
    05  bh-acquirer              pic x(6).
01  clr-presentment.
    05  filler                   pic x(2).
    05  pr-card-no               pic x(16).
    05  pr-amount                pic s9(9)v99 sign leading separate.
    05  pr-auth-code             pic x(6).
    05  pr-merchant              pic x(10).
    05  pr-mcc                   pic 9(4).
    05  pr-region                pic x.
    05  pr-product               pic x(2).
    05  pr-txn-date              pic 9(8).
01  clr-batch-trailer.
    05  filler                   pic x(2).
    05  bt-count                 pic 9(6).
    05  bt-amount                pic s9(13)v99 sign leading separate.
01  clr-file-trailer.
    05  filler                   pic x(2).
    05  ft-batch-count           pic 9(4).
    05  ft-record-count          pic 9(8).
    05  ft-amount                pic s9(15)v99 sign leading separate.
