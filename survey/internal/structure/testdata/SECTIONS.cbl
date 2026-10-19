       >>source format is free
identification division.
program-id. sections.
data division.
working-storage section.
01  ws-a                         pic 9.
procedure division.
main section.
    perform check
    perform tail-only
    goback.
check section.
    move 3 to ws-a.
check-more.
    move 4 to ws-a.
tail-only section.
tail-step.
    move 5 to ws-a.
