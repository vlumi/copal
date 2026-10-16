# The First Bank of Pangaea

A small COBOL batch bank, the specimen Copal is developed against. Fiction written for this project.

## Layout

- [copybooks/](copybooks/): record layouts, one copybook per record.
- [programs/](programs/): the COBOL programs, fixed-format, compiled with GnuCOBOL 3.2 in IBM dialect.
- [jobs/](jobs/): job files in a small subset of JCL, run by the Go runner in [runner/](runner/).
- [testdata/](testdata/): fixtures, each an `input` directory of datasets to start from and a `golden` directory of the datasets, job log and reports a run produces.

The Go packages around the COBOL are the bank's machine room, not the bank: [jcl/](jcl/) parses job files, [dataset/](dataset/) keeps the datasets and their catalog, [runner/](runner/) runs jobs and [compile/](compile/) compiles the programs. `cmd/pangaea` puts them on the command line.

## The daily batch

[DAILY.jcl](jobs/DAILY.jcl) runs three steps, each only if no earlier step ended with a return code (RC) above 4:

1. **TXVALID** checks the day's transaction file: a header for the business date in the bank control record, details, and a trailer whose detail count, amount total and account number hash match the details. A file that fails is rejected as a whole (RC 12). Details with a bad account number, transaction code or amount go to the rejects (RC 4).
2. **TXSORT** sorts the valid details into account number order.
3. **TXPOST** matches them against the account master and writes the new master, an audit trail record per posting and the transactions it could not post: unknown or closed accounts, withdrawals from frozen accounts or beyond the overdraft limit. It checks that the new master's balances equal the old ones plus what it posted.

`make pangaea-daily` runs it on the fixture in `testdata/daily` and shows the reports. `make test` runs it on every fixture and compares with the golden files; `make pangaea-golden` rewrites them after a deliberate change.

## Conventions

- Datasets are files of fixed-length records without delimiters, like a mainframe's `RECFM=FB`. The catalog in each dataset directory lists their record lengths. Fixtures and golden files hold the same records as text, one per line.
- Signed numbers are zoned decimal with an overpunched sign, EBCDIC style: the last digit is replaced by one character that carries both the digit and the sign, so the field is as long as its digits. `{ABCDEFGHI` stand for +0 to +9, `}JKLMNOPQR` for -0 to -9. In a 13-digit balance, `PIC S9(11)V99`, 55,000.01 is `0000005500001` with its last digit 1 written as `A`: `000000550000A`. -150.05 is `000000001500N`.
- A program knows its files by DD name (`ASSIGN TO TRANIN`); the runner binds each DD to its dataset through the `DD_TRANIN` environment variable.
- Programs never read the clock. The business date comes from the bank control record.
- Return codes: 0 fine, 4 warnings such as rejected records, 12 input rejected, 16 the program could not do its job.
