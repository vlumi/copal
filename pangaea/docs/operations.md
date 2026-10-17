# Batch operations guide

First Bank of Pangaea, Data Center. Revision of 2026-09. Supersedes the 2011 edition; the online day section is new.

## The business day

The bank runs one business day at a time, Monday to Friday. The date it is working on is the business date in the bank control record, `BANK.CONTROL(0)`, and nothing else: no program reads the clock. A day has two halves.

1. **The online day.** From opening, the authorization service answers the card terminals and the ATMs: purchases, reversals, cash withdrawals and balance enquiries, against the online account and card files loaded the night before. It holds what it approves. See [online.md](online.md).
2. **The batch window.** When the online day closes, the day's holds are extracted as the card file, and the nightly batch runs: `DAILY`, and on the first business day of a month `MONTHLY` after it. The batch posts the branches' transactions and the card file to the master, applies the day's account maintenance, loads the online files for the next day, and moves the business date on. During the window the authorization service answers that the bank is closed.

A new bank is opened once, with `START`, which loads the online files for its first day. After that, every day is the same.

## Running a day

In this installation the scheduler is `make bank-days N=1`, which runs the online day with whatever traffic was sent to it, closes it, and runs the batch. The console's Run button does the same. `make bank-reset` opens a new bank with `START`.

The online day can also be held open for real traffic with `make console` (the service is under `/api/`) or `make online`; interrupting either closes the day and extracts its holds, and the next `make bank-days` runs that day's batch alone. A business day's online day runs once; running the day again does not repeat it.

## Where the output is

Each day's output is under `sysout/DATE/` in the bank's directory: the `DAILY` job's files directly in it, and `ONLINE/` and `MONTHLY/` beside them. `START` writes to `sysout/START/`.

- `JOBLOG`: one line per step with its return code, `BYPASSED` for a step that did not run, or `ABENDED`.
- `JESYSMSG`: for every step that ran, each DD statement and the dataset it was bound to, with the records the dataset held when the step ended; `INDEXED` for an indexed file, whose records are not counted; `SYSOUT`, `DUMMY` or `IN-STREAM` for a DD bound to no dataset.
- `STEP.DDNAME`: the reports, such as `POST.POSTRPT` and `VALID.VALRPT`. Every step writes one. Read them: the counts and totals in them are the first thing to check.
- `STEP.STDOUT`: anything a program displayed, present only when it displayed something. A program that stops with return code 16 says why here.

## Return codes

Every program ends with one of four return codes, and the job log shows them.

| RC | Meaning | What to do |
| --: | --- | --- |
| 0 | The step did its job. | Nothing. |
| 4 | The step did its job but rejected some records: transactions that did not validate or could not post, maintenance requests the rules refuse, or transfers in transit that did not net to zero at month-end. | Normal on most days. Check the step's report and the rejects dataset; the branches get the rejects back. |
| 12 | The step rejected its input as a whole. Only `TXVALID` does this: a transaction file with no header, a header for another date or an unknown source, a trailer whose counts or totals do not match, or records after the trailer. | The file must be corrected and the day run again. Nothing after the step ran. |
| 16 | The step could not do its job: files out of step with each other, input out of order, a new master that does not reconcile, or an I/O error. | This is a fault. Read `STEP.STDOUT`, which says what the program found. Nothing after the step ran; the day is backed out. |

In `DAILY`, every step after the first runs only if no earlier step ended above 4: `COND=(4,LT)` in the JCL. So a 12 or 16 stops the day where it happened, and `BANKDATE` never moves the date on, which is what makes the rerun safe.

## A failed day

A business day's batch goes through whole or not at all. When a step ends with 12 or 16, or abends, the scheduler backs the batch out: every dataset the batch cataloged is deleted again, the generated input files included, and the day's job logs and reports move to `sysout/DATE.failed/`. The online day's own output stays where it was. The next run starts the same business day's batch afresh, from the same control record and master.

Procedure:

1. Find the step in `JOBLOG` under `sysout/DATE.failed/`, and read its report and `STEP.STDOUT`.
2. Correct the cause: a bad input file, a program, a missing dataset.
3. Run the day again. The online day is not run again; its card file is still there.

Nothing in the master, the audit trail or the rejects of earlier days is touched by a failed day. Generations only ever go forward.

## Month-end

On the first business day of a month, `DAILY` posts the closed month's interest and fees before the day's own transactions, and `MONTHLY` runs after it: statements for every account, and the general-ledger summary. `MONTHLY` reads every audit trail at once and picks the closed month's postings itself, so it needs no input of its own. A return code 4 from `GLSUM` means transfers in transit did not net to zero over the month, which means a transfer leg went missing somewhere: find it before the next month-end.

## What is kept

Every generation of every dataset is kept in the catalog, so the history of the master, the audit trail and the rejects since the bank opened is all there, and any day can be reconstructed. The console shows what changed between any two days' openings, reconciled against the postings in between. See [datasets.md](datasets.md).
