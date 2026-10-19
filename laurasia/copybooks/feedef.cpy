       >>source format is free
*> feedef: an interchange fee as the fee schedule gives it, a line each:
*> for a card product, a region and a merchant group, a rate on the amount
*> and a fixed fee in cents.
*>
*> history
*> 2012-01-09 ah  lcs-0140 original version.
01  fee-definition.
    05  fd-product               pic x(2).
    05  fd-region                pic x.
    05  fd-group                 pic 9.
    05  fd-rate                  pic 9v9(4).
    05  fd-fixed                 pic 9(5).
