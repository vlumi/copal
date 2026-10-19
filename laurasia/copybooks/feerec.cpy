       >>source format is free
*> feerec: an interchange fee in the fee table, a relative file whose slot
*> is worked out from the card product, the region and the merchant group.
*>
*> history
*> 2012-01-09 ah  lcs-0140 original version.
01  fee-slot-record.
    05  fs-rate                  pic 9v9(4).
    05  fs-fixed-cents           pic 9(5) comp-5.
