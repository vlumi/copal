# Copal

*Reading the code that still runs the books.*

Copal is a study of one problem: understanding a legacy COBOL batch system well enough to replace it. It has four parts, in four directories, and a second specimen beside the first.

| Directory | What it is |
| --- | --- |
| [survey/](survey/) | **Survey**, the analyzer. Reads a COBOL system and draws its blueprints: the layout of every record, the structure of every program, which fields each step reads and writes, the business rules buried in conditions, and the job network that ties the programs together through fixed-format files. All of it derived deterministically from the source. `survey explore` serves the same analysis in the browser. |
| [pangaea/](pangaea/) | **The First Bank of Pangaea**, the specimen. A small but working COBOL bank, fiction written for this project: accounts, a daily batch, a month-end run and an online day of card and ATM authorizations, compiled and run with GnuCOBOL. It knows nothing of the survey. |
| [laurasia/](laurasia/) | **Laurasia Card Services**, the second specimen. A card processor, issuer and acquirer in one, in another shop's COBOL: Micro Focus dialect, free format, sections, run by shell scripts instead of JCL. Clearing, matching, interchange, posting, disputes, settlement and a billing cycle. It is there to show where the survey assumes Pangaea. |
| [bedrock/](bedrock/) | What both stand on: reading COBOL source and JCL, laying out records from copybooks, and reading and writing records by field. The survey analyzes with it; the bank's machine room runs the bank with it. |
| [console/](console/) | Built for the modernization: a web window that runs the bank, watches its online day, and shows each run's jobs as the survey draws them, with links to the blueprints. The one part that knows both the bank and the survey. |

The survey of Pangaea is checked in as [blueprints/pangaea/](blueprints/pangaea/README.md), and of Laurasia as [blueprints/laurasia/](blueprints/laurasia/README.md), golden files that `make ci` regenerates and compares, so a change in the analyzer shows as a diff in what it says about either. The dependencies run one way: `bedrock` imports nothing else, `survey`, `pangaea` and `laurasia` import `bedrock`, and only `console` imports the rest; `make ci` checks it.

Survey does not translate COBOL into anything. It produces the understanding a translation needs, and later the tests a rewrite has to pass.

**Status:** Pangaea runs as a bank: about three thousand generated accounts, a daily batch of validation, interest accrual, sorting and posting in three passes for transfers, month-end interest, fees, statements and a general-ledger summary, and a business calendar, day after day (`make bank-reset bank-days N=20`); and an online day of card purchases and ATM withdrawals, authorized by COBOL programs called from a long-running dispatcher, against indexed files the night's batch loads, and posted by that batch. `make console` opens the console in the browser: accounts, each day's jobs as they ran with their networks, files and reports, what the days changed in the master, and the online channel with its authorizations live; `make simulate` sends it card traffic. Small fixtures check the batch against golden outputs. Survey reads the bank and writes blueprints of every program and copybook: each program's own description, record layouts with positions and lengths, files, its paragraphs with a diagram of how control moves between them, what each paragraph reads and writes, where every data item's values come from, and its business rules as decision tables; and of every job, the network of its steps and the datasets they read and write, with findings where a job and its programs disagree. `make blueprints` regenerates them, and `make explore` serves the explorer, where every job, step, dataset, program, paragraph, field and copybook is a link, and any field's values can be followed back across programs and job steps to the input they came from, on the diagrams. `make oracle` runs each program, compiled with GnuCOBOL's trace code, on input data made from its own conditions until the cases reach as many rules of its decision tables as they can; the cases, with what the programs wrote, are checked in and replayed by the tests. Laurasia Card Services, the second specimen, in Micro Focus COBOL and run by shell scripts, is surveyed the same way, with nothing special-cased for either: [blueprints/laurasia/](blueprints/laurasia/README.md), [oracle/laurasia/](oracle/laurasia/README.md), and `make explore-laurasia`. [docs/architecture.md](docs/architecture.md) describes how the parts fit, [docs/blueprints.md](docs/blueprints.md) how to read what the survey writes, [docs/brief.md](docs/brief.md) has the plan and [docs/roadmap.md](docs/roadmap.md) the order of work.

## A quick tour

Everything is a make target, and plain `make` lists them. You need GnuCOBOL 3.2 (`brew install gnucobol`), Go, Node for the Markdown linter, and `golangci-lint` for `make ci`.

1. Check that the tools are there.

   ```sh
   make toolchain
   ```

2. Run the daily batch on the small fixture. It prints the job log and every step's report: the bank in one screen.

   ```sh
   make pangaea-daily
   ```

3. Open a bank of about three thousand generated accounts on 2026-09-01, and run five business days, one line each: transactions, rejects, the worst return code, and the month-end job when a month closes. Everything it made is under `pangaea/run/bank/`.

   ```sh
   make bank-reset
   make bank-days N=5
   ```

4. Serve the console and open <http://127.0.0.1:4510/>: the bank today, with buttons to run days and reset. Open a day to see its jobs drawn as networks with each step's return code and each dataset's records, and the reports; open an account for its master record field by field and its postings; **Compare** shows what the days changed, reconciled against the postings.

   ```sh
   make console
   ```

5. With the console running, send it card and ATM traffic from another terminal. The console's **Online** page shows each authorization as it is answered, approvals and declines, and balances moving. Then run a day from the console: the channel answers that the bank is closed for the batch, and the next morning the card purchases and withdrawals are in the master, under the codes CP, AW and AF.

   ```sh
   make simulate
   ```

6. Read the survey of the bank, which is checked in: start at [blueprints/pangaea/README.md](blueprints/pangaea/README.md), then a program such as [TXPOST](blueprints/pangaea/programs/TXPOST.md), its [decision tables](blueprints/pangaea/rules/TXPOST.md), and the [DAILY job](blueprints/pangaea/jobs/DAILY.md). Regenerating them shows a change to the bank or the analyzer as a diff.

   ```sh
   make blueprints
   ```

7. Serve the explorer and open <http://127.0.0.1:1192/>: the same analysis as a web app. Search for `AM-BALANCE`, open it in TXPOST, and follow its values back: the trail runs through every program and job step to the files that come from outside, highlighted on the diagrams.

   ```sh
   make explore
   ```

8. Look at the oracle's cases, also checked in: input data for each program that reaches the rules of its decision tables, with what the program wrote and returned, the evidence a rewrite has to match. Start at [oracle/pangaea/README.md](oracle/pangaea/README.md), then [TXPOST's](oracle/pangaea/TXPOST/README.md), rule by rule; the explorer shows them on each program's page. Searching for them again takes a couple of minutes.

   ```sh
   make oracle
   ```

9. Read the rest: the bank's own operations guide in [pangaea/docs/](pangaea/docs/README.md), in its operators' voice; [docs/architecture.md](docs/architecture.md) on how the four parts fit; [docs/blueprints.md](docs/blueprints.md) on how to read what the survey writes.

10. Run what the hosted build runs: linters, the dependency check, every test and the golden files, the oracle's cases replayed among them.

   ```sh
   make ci
   ```

## The names

Copal is young amber: resin that has hardened but not quite fossilized, and the traditional base of varnish for old wood. Said aloud, it is not far from COBOL.

Pangaea is the single supercontinent of the age of the dinosaurs, before it broke into the continents of today. A monolith, waiting to be split.

A survey walks the ground and maps it before anyone builds; bedrock is what everything stands on.

## License

MIT, see [LICENSE](LICENSE). The First Bank of Pangaea, its programs, data and documents included, is fiction written for this project and is covered by the same license.
