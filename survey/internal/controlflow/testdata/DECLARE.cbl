       >>source format is free
identification division.
program-id. declare.
environment division.
input-output section.
file-control.
    select in-file assign to "IN"
        organization is line sequential.
data division.
file section.
fd  in-file.
01  in-record                    pic x(10).
working-storage section.
01  ws-failed                    pic x value "n".
procedure division.
declaratives.
in-errors section.
    use after error procedure on in-file.
note-failure.
    move "y" to ws-failed.
end declaratives.
main section.
    open input in-file
    move "z" to ws-failed
    goback.
