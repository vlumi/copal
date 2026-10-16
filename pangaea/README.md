# The First Bank of Pangaea

A small COBOL batch bank, the specimen Copal is developed against. Fiction written for this project.

## Layout

- [copybooks/](copybooks/): record layouts, one copybook per record.
- [programs/](programs/): the COBOL programs, fixed-format, compiled with GnuCOBOL 3.2 in IBM dialect.
- [jobs/](jobs/): job files in a small subset of JCL, run by the Go runner in [runner/](runner/): the daily batch and the month-end job.
- [testdata/](testdata/): fixtures, each an `input` directory of datasets to start from and a `golden` directory of the datasets, job log and reports a run produces.

The Go packages around the COBOL are the bank's machine room, not the bank: [jcl/](jcl/) parses job files, [dataset/](dataset/) keeps the datasets and their catalog, [runner/](runner/) runs jobs and [compile/](compile/) compiles the programs. [generate/](generate/) makes up customers and their business, and [bank/](bank/) runs the bank day after day, and [console/](console/) shows it in the browser. `cmd/pangaea` puts them on the command line.

## The daily batch

[DAILY.jcl](jobs/DAILY.jcl) runs one business day. A step runs only if no earlier step ended with a return code (RC) above 4.

1. **TXVALID** checks the day's transaction file: a header for the business date in the bank control record, details, and a trailer whose detail count, amount total and account number hash match the details. A file that fails is rejected as a whole (RC 12). Details with a bad account number, transaction code, amount or transfer counter account go to the rejects (RC 4).
2. **INTACCR** accrues interest for the calendar days since the last business day, on the balances that stood over them, at the rate table's rate for the account type and balance tier, or its debit rate below zero. Accruals are packed decimal to six decimals. When the days cross into a new month, the old month's interest becomes due, rounded to the cent.
3. **MONCHG**, on the day a month closes, writes its interest credits, overdraft interest and the maintenance fee of open current accounts below 1,000.00 as transactions, in account order, and **TXPOST** posts them without a sort. On other days there are none, and TXPOST copies the master.
4. **TXSORT** sorts the valid details into account number order.
5. **TXPOST** matches them against the account master and writes a new master, an audit trail record per posting, and the transactions it could not post: unknown or closed accounts, withdrawals and transfers from frozen accounts or beyond the overdraft limit. It checks that the new master's balances equal the old ones plus what it posted.
6. **TXSORT** and **TXPOST** again, for the transfer credits: a transfer that posted debits its sender in step 5 and writes a credit for the receiving account, which posts here.
7. **TXSORT** and **TXPOST** a third time, for the returns: a credit to an unknown or closed account goes back to its sender.
8. **ACCTMNT** applies the branches' account maintenance requests, in account order, and writes the day's new master: it opens accounts, closes empty ones, freezes and unfreezes them, and changes overdraft limits, and rejects what the rules do not allow. It is the bank's oldest style of code, from 1988: `PERFORM ... THRU` ranges that jump to their exit paragraph with `GO TO` and run on from one paragraph into the next.
9. **BANKDATE** writes the bank control record for the next business day, Monday to Friday.

The day's inputs are the current generations, (0), of the bank control record, the account master, the accruals, the transaction file and the maintenance requests, and the rate table; it writes the next generations, (+1), of the control record, the master and the accruals, and the day's audit trail and the rejects of both transactions and maintenance. Everything in between is temporary. So a month's interest and fees are posted on the first business day of the next month, as the batch's first postings.

`make pangaea-daily` runs it on the fixture in `testdata/daily` and shows the reports; the fixture in `testdata/month-end` closes October and runs the month-end job after it.

## The month-end job

[MONTHLY.jcl](jobs/MONTHLY.jcl) runs on the first business day of every month, after the daily batch has posted the closed month's interest and fees; the bank's scheduler starts it, as an operations team's would. Both of its programs read every audit trail at once, the whole generation data group, and pick the closed month's postings, and the interest and fees posted for it, by the same rule.

1. **STMTPRT** prints a statement for every account: the opening balance, the month's postings in the order they posted, the interest and fees posted for the month, and the closing balance. It picks the postings in a `SORT` input procedure and prints in the output procedure against the master the month closed with.
2. **GLSUM** summarizes the month for the general ledger in double entry: every posting moves customer deposits on one side and, on the other, the ledger account the code table names for its transaction code: cash, fee or interest income, interest expense, or transfers in transit. It writes the month's debits and credits per ledger account to `GL.SUMMARY`, and ends with RC 4 if transfers in transit do not net to zero, which would mean a transfer leg went missing.

The two programs share the transaction code table, [CODETAB](copybooks/CODETAB.cpy), each copying it with a prefix of its own: `COPY CODETAB REPLACING ==:P:== BY ==ST==`. `make test` runs it on every fixture and compares with the golden files; `make pangaea-golden` rewrites them after a deliberate change.

## The running bank

`make bank-reset` makes a new bank in `pangaea/run/bank`: about three thousand accounts for two thousand customers, generated from a seed (`SEED=1`, `ACCOUNTS=3000` by default), on business day 2026-09-01. `make bank-days N=5` runs the next five business days: for each it generates the branches' transaction and maintenance files for the day and runs the daily batch on them, and on the first business day of a month the month-end job too.

Every generation stays cataloged in `pangaea/run/bank/datasets`, so the history of the master, the audit trail and the rejects is all there; each day's job log and reports are in `pangaea/run/bank/sysout/DATE`. The same seed always gives the same bank, day after day.

The generator ([generate/](generate/)) reads and writes records through the layouts Copal computes from the copybooks, so it never hand-codes a record format. Customers have habits that follow from the seed: a salary on the 25th, rent on the first business day of the month, savings after payday, cash, deposits and transfers to other customers. Branches make the odd mistake, and customers sometimes try to spend money they do not have, so a few percent of each day ends up in the rejects.

## The console

`make console` serves a window onto the running bank at <http://localhost:8080/> (`ADDR` to change it):

- the bank today: its business date, accounts and balances, and every day it has run, with each day's worst return code;
- accounts, searched by account number, customer number or name, and for each account its master record field by field, every posting and every reject since the reset;
- each business day: the job's steps, the transaction file's control totals against its details, the rejects, the reports, the audit trail and the transaction file;
- what changed between any two days' openings: balances, status changes and new accounts, reconciled against the postings in the audit trails in between.

It reads every record through the layouts Copal computes from the copybooks, the same way the generator writes them. Its buttons reset the bank and run business days by calling the same code as `make bank-reset` and `make bank-days`; the console itself never writes the bank's files.

## Conventions

- Datasets are files of fixed-length records without delimiters, like a mainframe's record format FB (fixed, blocked). The catalog in each dataset directory lists their record lengths (LRECL). Fixtures and golden files hold the same records as text, one per line.
- A generation data group (GDG) keeps every version of a dataset: `ACCT.MASTER.G0007V00` is the seventh master. A job names them relative to the newest when it started: `ACCT.MASTER(0)` is the current one, `ACCT.MASTER(+1)` the one it writes. A temporary dataset, `&&SORTED`, lives only while its job runs. Named by its base alone, `TRAN.AUDIT`, a group is every generation it had when the job started, newest first, read as one file. So is a concatenation: a DD statement without a name adds its dataset to the one above.
- A record with bytes that are not printable is written in hex in fixtures and golden files, the way COBOL writes a hex literal: `X'00012C'`.
- Packed decimal (`COMP-3`, `PACKED-DECIMAL`) holds two digits a byte and the sign in the last half-byte: C positive, D negative, F unsigned. In `PIC S9(3)V99 COMP-3`, 12.00 is the three bytes `X'01200C'`.
- Signed display numbers are zoned decimal with an overpunched sign, EBCDIC style: the last digit is replaced by one character that carries both the digit and the sign, so the field is as long as its digits. `{ABCDEFGHI` stand for +0 to +9, `}JKLMNOPQR` for -0 to -9. In a 13-digit balance, `PIC S9(11)V99`, 55,000.01 is `0000005500001` with its last digit 1 written as `A`: `000000550000A`. -150.05 is `000000001500N`.
- A program knows its files by DD name (`ASSIGN TO TRANIN`); the runner binds each DD to its dataset through the `DD_TRANIN` environment variable.
- Programs never read the clock. The business date comes from the bank control record.
- Return codes: 0 fine, 4 warnings such as rejected records, 12 input rejected, 16 the program could not do its job.
