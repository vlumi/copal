       >>source format is free
*> feeslot: the fee table's slot for a card product, a region and a
*> merchant group, or zero for none: three products, two regions and five
*> groups, thirty slots. the slot is the return code, which a caller
*> receives with call ... returning.
*>
*> history
*> 2012-01-09 ah  lcs-0140 original version.
identification division.
program-id. feeslot.

data division.
working-storage section.
01  ws-product-index             pic 9 comp-5.
01  ws-region-index              pic 9 comp-5.
01  ws-slot                      pic 9(4) comp-5.
linkage section.
01  lk-product                   pic x(2).
01  lk-region                    pic x.
01  lk-group                     pic 9.

procedure division using lk-product lk-region lk-group.
main section.
    move zero to return-code
    evaluate lk-product
        when "CL" move 1 to ws-product-index
        when "GD" move 2 to ws-product-index
        when "BZ" move 3 to ws-product-index
        when other goback
    end-evaluate
    evaluate lk-region
        when "D" move 1 to ws-region-index
        when "I" move 2 to ws-region-index
        when other goback
    end-evaluate
    if lk-group < 1 or lk-group > 5
        goback
    end-if
    compute ws-slot = ((ws-product-index - 1) * 2 + ws-region-index - 1) * 5 + lk-group
    move ws-slot to return-code
    goback.
