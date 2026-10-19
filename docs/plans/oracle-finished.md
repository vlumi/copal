# Plan: the oracle, finished

The oracle runs every program of both specimens on cases found from their decision tables, and Pangaea's jobs as batches on those cases. What is left is Laurasia's batch, which is run by shell scripts rather than JCL, and the 95 rules no case reaches. Most of those are rules a run could take, given a file that fails or data no clean case holds, or rules the survey could show no run can take, given what a called program returns.

Read [oracle.md](oracle.md) and [deterministic-sweep.md](deterministic-sweep.md) first; their ground rules apply. Facts only, goldens with the change that causes them, a test for every step, `make ci` before pushing, and no special case for either specimen.

## Where the unreached rules are

| Kind | Rules | Example |
| --- | --: | --- |
| A file status other than success | about 25 | `WS-ACCOUNTS-STATUS = '00'` N, in ACCTLD |
| A combination a called program cannot return | about 26 | AUTHSRV's tests of `RS-RESULT` after the call |
| Totals that agree by construction | a few | `WS-DEBIT-TOTAL = WS-CREDIT-TOTAL` N, in GLSUM |
| Combinations within a program the survey does not rule out | the rest | `WS-REQUEST-KEY < WS-MASTER-KEY` with `WS-MASTER-KEY = WS-ACTIVE-KEY`, in ACCTMNT |

## Parts

1. **A batch of Laurasia's scripts.** `survey oracle -jobs` takes scripts as well as jobs, and runs them as they are run, with sh, on inputs from the programs' cases.

   As built: the scripts run with sh, in the order given, on directories of their environment that the scripts before them wrote to. A variable of the environment that names the directory of what a step runs, `LCS_BIN` or `LCS_SCRIPTS`, names the oracle's stand-ins instead. These run the program, the script, `sort` or `mv` with a program's trace in a file named after its process, and record how each ended, so each step has an outcome without the oracle reading the scripts' own logs. A file a script reads that nothing before it wrote takes the records of the cleanest case of the first program that reads it, through a sort or a move between them. A parameter takes the command line of the cleanest case of the first program it is passed to, the only value that case is consistent with, and the variable that stands for it, `DATE`, has one value across the batch. Laurasia's cases leave the date blank wherever no rule needs it, so the batch runs with a blank date, and its file names show it. `make oracle` runs open, nightly and cycle; every step of the night runs, and open ends with ACCTLD's return code 8, as its one case does.

   Searching Laurasia again showed that CLRIN's checked-in cases were older than the search: a change since had altered what it finds, and the cases had been replayed but not searched again. They are searched again in their own commit, and reach the same rules.
2. **I/O errors.** Cases that fail a file on purpose, the same way every time: a file that is missing, or a record of the wrong length. They reach the rules that only a file status other than success takes.

   As built: three levers for each file, tried after the changes to records, and one failure a case. A file a program reads may be missing, so that opening it fails with status 35; a sequential one may be a directory, so that it opens and reading it fails with status 30; and a file it only writes may be bound to a path in a directory that is not there, so that creating it fails with status 30. GnuCOBOL gives the same statuses on macOS and on Linux, through the same indexed file handler. A record of the wrong length was left out: GnuCOBOL reads a short record of a fixed-length file without an error, and a long line of a line sequential one with status 06, which no program tests. A case records its failure in `input/failed.txt`, and a batch never takes its inputs from one. The rules reached go from 413 to 437 of 508: every rule that only a file status other than success takes is reached, 24 in all, and 71 rules are left.
3. **What a called program returns.** A field the called program sets is carried back into the caller, so that a combination of conditions no call can return is left out of the caller's rules, as one no path can take is.
4. **Totals that agree by construction.** Rules that only bad data reaches, such as totals that disagree when the program computed both, listed apart from the rest, with the reason.

Each part one pull request.
