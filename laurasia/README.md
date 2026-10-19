# Laurasia Card Services

A small COBOL card processor, issuer and acquirer in one, the survey's second specimen: a nightly run that takes in the clearing file, matches it to the day's authorizations, prices the interchange, posts to the cardholders' accounts, moves the disputes on and settles with the merchants, and a billing cycle. Fiction written for this project. Pangaea split into Laurasia and Gondwana; Laurasia is written in another shop's COBOL and run another way, so that what the survey learned from Pangaea alone shows where it does not carry over. It knows nothing of the survey.

## Layout

[docs/](docs/) is the processor's own documentation, in its operations team's voice: the night, the scripts, the files and the return codes.

- [copybooks/](copybooks/): record layouts. The clearing file's copybook has all five of its record layouts, one 01 record each.
- [programs/](programs/): the COBOL programs, free format and lower case, compiled with GnuCOBOL 3.2 in Micro Focus dialect, `-std=mf`. Sections are the unit of work; files are line sequential, relative and indexed.
- [scripts/](scripts/): the batch, as POSIX shell scripts: each step binds its files with `export DD_...`, runs a program or Unix `sort`, and checks the return code with `if`. Files are named by date.
- [testdata/](testdata/): fixtures, each an `input` directory of files to start from, a `RUN` file of the scripts to run on them, in order, and a `golden` directory of what the run leaves behind.

The only Go is [laurasia_test.go](laurasia_test.go), which compiles the programs with [bedrock](../bedrock/run/), runs each fixture's scripts and compares the result with its golden files. `make laurasia-golden` records them again.
