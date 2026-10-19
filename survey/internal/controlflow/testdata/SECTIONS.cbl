       >>source format is free
identification division.
program-id. sections.
data division.
working-storage section.
01  ws-a                         pic 9.
procedure division.
main section.
    perform check
    move 2 to ws-a
    goback.
check section.
    if ws-a = 1
        exit section
    end-if
    move 3 to ws-a.
check-more.
    move 4 to ws-a.
after-check section.
    move 5 to ws-a.
