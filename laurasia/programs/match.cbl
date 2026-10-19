       >>source format is free
*> match: match the presentments to their authorizations
*>
*> reads the presentments and the day's authorizations, both in
*> authorization code order, and marks each presentment: matched to an
*> approval, presented on a declined authorization, for more than the
*> approval allows, or with no authorization at all. authorizations no
*> presentment claims are counted; they lapse.
*>
*> return codes: 0 all matched, 4 some not.
*>
*> history
*> 2011-04-04 rk  lcs-0102 original version.
*> 2014-02-17 rk  lcs-0251 marks instead of rejects.
*> 2018-05-14 ah  lcs-0421 a presentment may exceed its approval by 15%.
identification division.
program-id. match.

environment division.
input-output section.
file-control.
    select pres-in assign to "PRESIN"
        organization is line sequential.
    select auth-in assign to "AUTHIN"
        organization is line sequential.
    select pres-out assign to "PRESOUT"
        organization is line sequential.
    select report-out assign to "MATCHRPT"
        organization is line sequential.

data division.
file section.
fd  pres-in.
01  pres-in-record               pic x(76).
fd  auth-in.
    copy "authlog.cpy".
fd  pres-out.
    copy "presrec.cpy".
fd  report-out.
01  report-line                  pic x(60).

working-storage section.
78  tolerance                    value 1.15.
01  ws-pres-key                  pic x(6).
01  ws-auth-key                  pic x(6).
01  ws-counts.
    05  ws-matched               pic 9(8) value zero.
    05  ws-not-matched           pic 9(8) value zero.
    05  ws-lapsed                pic 9(8) value zero.
01  ws-auth-claimed              pic x value "N".
    88  auth-claimed             value "Y".
01  ws-report-count.
    05  rc-label                 pic x(30).
    05  rc-value                 pic z(7)9.

procedure division.
main section.
    open input pres-in auth-in
         output pres-out report-out
    perform read-presentment
    perform read-authorization
    perform match-next until ws-pres-key = high-values and ws-auth-key = high-values
    perform write-report
    close pres-in auth-in pres-out report-out
    if ws-not-matched > zero
        stop run returning 4
    end-if
    stop run returning 0.

read-presentment section.
    read pres-in into pres-record
        at end move high-values to ws-pres-key
        not at end move ps-auth-code to ws-pres-key
    end-read.

read-authorization section.
    read auth-in
        at end move high-values to ws-auth-key
        not at end move al-auth-code to ws-auth-key
    end-read
    move "N" to ws-auth-claimed.

*> the lower key goes first: a presentment with no authorization, or an
*> authorization no presentment claims. several presentments may claim
*> one authorization.
match-next section.
    evaluate true
        when ws-pres-key < ws-auth-key
            set ps-no-auth to true
            perform write-presentment
            perform read-presentment
        when ws-pres-key > ws-auth-key
            if not auth-claimed
                add 1 to ws-lapsed
            end-if
            perform read-authorization
        when other
            perform mark-presentment
            set auth-claimed to true
            perform write-presentment
            perform read-presentment
    end-evaluate.

mark-presentment section.
    evaluate true
        when al-declined
            set ps-declined to true
        when ps-amount > al-amount * tolerance
            set ps-over-auth to true
        when other
            set ps-matched to true
    end-evaluate.

write-presentment section.
    write pres-record
    if ps-matched
        add 1 to ws-matched
    else
        add 1 to ws-not-matched
    end-if.

write-report section.
    move "MATCH    PRESENTMENTS TO AUTHORIZATIONS" to report-line
    write report-line
    move "MATCHED" to rc-label
    move ws-matched to rc-value
    write report-line from ws-report-count
    move "NOT MATCHED" to rc-label
    move ws-not-matched to rc-value
    write report-line from ws-report-count
    move "AUTHORIZATIONS LAPSED" to rc-label
    move ws-lapsed to rc-value
    write report-line from ws-report-count.
