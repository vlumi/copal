# Laurasia Card Services

A small COBOL card processor, issuer and acquirer in one, the survey's second specimen: a nightly run that takes in the clearing file, matches it to the day's authorizations, prices the interchange, posts to the cardholders' accounts, moves the disputes on and settles with the merchants, and a billing cycle. Fiction written for this project. Pangaea split into Laurasia and Gondwana; Laurasia is written in another shop's COBOL and run another way, so that what the survey learned from Pangaea alone shows where it does not carry over. It knows nothing of the survey.

## Layout

[docs/](docs/) is the processor's own documentation, in its operations team's voice: the night, the scripts, the files and the return codes.

- [copybooks/](copybooks/): record layouts. The clearing file's copybook has all five of its record layouts, one 01 record each.
- [programs/](programs/): the COBOL programs, free format and lower case, compiled with GnuCOBOL 3.2 in Micro Focus dialect, `-std=mf`. Sections are the unit of work; files are line sequential, relative and indexed.
- [scripts/](scripts/): the batch, as POSIX shell scripts: each step binds its files with `export DD_...`, runs a program or Unix `sort`, and checks the return code with `if`. Files are named by date.
- [testdata/](testdata/): fixtures, each an `input` directory of files to start from, a `RUN` file of the scripts to run on them, in order, and a `golden` directory of what the run leaves behind.

The Go is the processor's machine room, and small: [night/](night/) compiles the programs with [bedrock](../bedrock/run/) and runs the scripts, a fixture's on a copy of its input or one at a time; [records/](records/) loads the record layouts from the copybooks with cobolith, in the dialect; [generate/](generate/) makes up the portfolio and each day's inbound files from a seed; [processor/](processor/) opens the processor on them and runs it day after day. [laurasia_test.go](laurasia_test.go) runs every fixture and compares the result with its golden files, and `make laurasia-night` runs the one fixture and shows its logs and reports. `make laurasia-golden` records the golden files again.

## The running processor

`make cards-open` opens a processor of about a thousand generated cards, and `make cards-days N=30` runs its next thirty days, a line each. On a business day the switch's log, the networks' clearing file, the card system's feed of payments and status changes and the disputes desk's events are generated from the cards as the last cycle unloaded them, the night runs on them and the billing cycle after it; on a weekend only the cycle runs, and Monday's log holds the weekend's authorizations. Its files are under `laurasia/run/processor/data`, laid out as the [operations guide](docs/operations.md) says. A day that fails is backed out whole, and what it left is kept under `failed/`.

The processor has no online side of its own: the networks' switch authorizes around the clock, standing in for the issuer when it cannot answer, and the night reads only its log. So no hour closes the cards, unlike the bank's own channel. What the generator makes up is what a processor sees: purchases by card product and day of the week, declines past the limit, tips presented over the approval, presentments whose authorization is in an earlier log, a mistyped record now and then, a batch whose trailer is a cent out, cardholders paying in full, the minimum or nothing some days after their statement, cards lost and found, and disputes moving through their states. What months of it show is in [the plan](../docs/plans/living-laurasia.md).
