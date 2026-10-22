# Gondwana Life Insurance

A small COBOL life insurer, the survey's third specimen: policies that run for decades, a daily run of new business, anniversaries, benefits, changes and lapses, a monthly premium collection through the banks' transfer file, a month-end close of reserves and commissions, and a yearly dividend. Fiction written for this project. Pangaea split into Laurasia and Gondwana; Gondwana is written in a Japanese shop's COBOL, Fujitsu NetCOBOL, in Shift_JIS, with its names and comments in Japanese, its amounts in whole yen and its dates in the Japanese eras, and it is run by Windows batch files. It knows nothing of the survey.

## Layout

[docs/](docs/) is the insurer's own documentation, in its operations team's voice: the runs, the batch files, the files and the return codes.

- [copybooks/](copybooks/): record layouts, and the loading of the product and rate tables as procedure copybooks. The data names are half-width katakana with an ASCII record prefix, `PM-ｼｮｳｹﾝﾊﾞﾝｺﾞｳ` for the policy number of the master, as the code of that world was written when compilers had only single-byte kana.
- [programs/](programs/): the COBOL programs, fixed format, compiled with GnuCOBOL 3.2 under its default configuration, which takes what NetCOBOL and it share. Programs are named the shop's way, `GLD010` for the first daily step; the header comment says what each does. Paragraphs are named `ﾒｲﾝ-RTN`, `ｼｮｷ-RTN` and the like; files, reports and programs keep ASCII names, which the compiler turns into C identifiers.
- [scripts/](scripts/): the batch, as Windows batch files: each step binds its files with `set DD_...=`, runs a program with its output redirected, logs the errorlevel, and stops the file with `if errorlevel 8`. Files are named by date. Bedrock runs them with an interpreter of its own, since no `cmd.exe` runs where the tests do.
- [testdata/](testdata/): fixtures, each an `input` directory of files to start from, a `RUN` file of the batch files to run on them, in order, and a `golden` directory of what the run leaves behind.

Everything the programs read and write is Shift_JIS, and the working tree holds it that way; the repository stores these files as UTF-8 through `.gitattributes`, so that diffs and the web read them. A fixed-format line is at most 72 bytes, a two-byte character counting two.

The Go is the insurer's machine room, and small: [batch/](batch/) compiles the programs with [bedrock](../bedrock/run/) and runs the batch files through its [interpreter](../bedrock/batchfile/), a fixture's on a copy of its input or one at a time; [cmd/gondwana](cmd/gondwana/) puts them on the command line. [gondwana_test.go](gondwana_test.go) runs every fixture and compares the result with its golden files, and `make gondwana-month` runs the one fixture and shows its logs and reports, transcoded for the terminal. `make gondwana-golden` records the golden files again.

## The month

The fixture is March 2027, the ninth year of Reiwa, with policies from Showa, Heisei and Reiwa in one master. The system opens on March 1 from the extract, and the month runs three business days, the billing on the 10th for the debit on the 27th, the collection agent's results on the 30th, the close on the 31st and, March being the end of the fiscal year, the dividends. What happens in it is in the [operations guide](docs/operations.md).
