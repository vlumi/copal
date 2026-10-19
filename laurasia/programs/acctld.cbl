       >>source format is free
*> acctld: load the card accounts into the indexed account file
*>
*> return codes: 0 loaded, 8 some cards twice, 16 cannot write.
*>
*> history
*> 2011-04-04 rk  lcs-0102 original version.
identification division.
program-id. acctld.

environment division.
input-output section.
file-control.
    select accounts-in assign to "ACCTIN"
        organization is line sequential.
    select accounts assign to "ACCOUNTS"
        organization is indexed
        access mode is random
        record key is ca-card-no
        file status is ws-accounts-status.

data division.
file section.
fd  accounts-in.
01  account-line                 pic x(61).
fd  accounts.
    copy "cardacct.cpy".

working-storage section.
01  ws-accounts-status           pic xx.
01  ws-at-end                    pic x value "N".
    88  accounts-at-end          value "Y".
01  ws-result                    pic 99 value zero.

procedure division.
main section.
    open input accounts-in output accounts
    if ws-accounts-status not = "00"
        stop run returning 16
    end-if
    perform until accounts-at-end
        read accounts-in
            at end set accounts-at-end to true
            not at end
                write card-account from account-line
                    invalid key move 8 to ws-result
                end-write
        end-read
    end-perform
    close accounts-in accounts
    stop run returning ws-result.
