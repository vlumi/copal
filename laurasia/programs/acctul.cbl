       >>source format is free
*> acctul: unload the indexed account file, in card number order
*>
*> history
*> 2011-04-04 rk  lcs-0102 original version.
identification division.
program-id. acctul.

environment division.
input-output section.
file-control.
    select accounts assign to "ACCOUNTS"
        organization is indexed
        access mode is sequential
        record key is ca-card-no
        file status is ws-accounts-status.
    select accounts-out assign to "ACCTOUT"
        organization is line sequential.

data division.
file section.
fd  accounts.
    copy "cardacct.cpy".
fd  accounts-out.
01  account-line                 pic x(61).

working-storage section.
01  ws-accounts-status           pic xx.
    88  accounts-at-end          value "10".

procedure division.
main section.
    open input accounts output accounts-out
    if ws-accounts-status not = "00"
        stop run returning 16
    end-if
    read accounts next record
    perform until accounts-at-end
        write account-line from card-account
        read accounts next record
    end-perform
    close accounts accounts-out
    stop run returning 0.
