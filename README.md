# Copal

*Reading the code that still runs the books.*

Copal is a study of one problem: understanding a legacy COBOL batch system well enough to replace it. It has four parts, in four directories.

| Directory | What it is |
| --- | --- |
| [survey/](survey/) | **Survey**, the analyzer. Reads a COBOL system and draws its blueprints: the layout of every record, the structure of every program, which fields each step reads and writes, the business rules buried in conditions, and the job network that ties the programs together through fixed-format files. All of it derived deterministically from the source. `survey explore` serves the same analysis in the browser. |
| [pangaea/](pangaea/) | **The First Bank of Pangaea**, the specimen. A small but working COBOL bank, fiction written for this project: accounts, a daily batch, a month-end run and an online day of card and ATM authorizations, compiled and run with GnuCOBOL. It knows nothing of the survey. |
| [bedrock/](bedrock/) | What both stand on: reading COBOL source and JCL, laying out records from copybooks, and reading and writing records by field. The survey analyzes with it; the bank's machine room runs the bank with it. |
| [console/](console/) | Built for the modernization: a web window that runs the bank, watches its online day, and shows each run's jobs as the survey draws them, with links to the blueprints. The one part that knows both the bank and the survey. |

The survey of Pangaea is checked in as [blueprints/pangaea/](blueprints/pangaea/README.md), golden files that `make ci` regenerates and compares, so a change in the analyzer shows as a diff in what it says about the bank. The dependencies run one way: `bedrock` imports nothing else, `survey` and `pangaea` import `bedrock`, and only `console` imports all three; `make ci` checks it.

Survey does not translate COBOL into anything. It produces the understanding a translation needs, and later the tests a rewrite has to pass.

**Status:** Pangaea runs as a bank: about three thousand generated accounts, a daily batch of validation, interest accrual, sorting and posting in three passes for transfers, month-end interest, fees, statements and a general-ledger summary, and a business calendar, day after day (`make bank-reset bank-days N=20`); and an online day of card purchases and ATM withdrawals, authorized by COBOL programs called from a long-running dispatcher, against indexed files the night's batch loads, and posted by that batch. `make console` opens the console in the browser: accounts, each day's jobs as they ran with their networks, files and reports, what the days changed in the master, and the online channel with its authorizations live; `make simulate` sends it card traffic. Small fixtures check the batch against golden outputs. Survey reads the bank and writes blueprints of every program and copybook: each program's own description, record layouts with positions and lengths, files, its paragraphs with a diagram of how control moves between them, what each paragraph reads and writes, where every data item's values come from, and its business rules as decision tables; and of every job, the network of its steps and the datasets they read and write, with findings where a job and its programs disagree. `make blueprints` regenerates them, and `make explore` serves the explorer, where every job, step, dataset, program, paragraph, field and copybook is a link, and any field's values can be followed back across programs and job steps to the input they came from, on the diagrams. [docs/architecture.md](docs/architecture.md) describes how the parts fit, [docs/blueprints.md](docs/blueprints.md) how to read what the survey writes, [docs/brief.md](docs/brief.md) has the plan and [docs/roadmap.md](docs/roadmap.md) the order of work.

## The names

Copal is young amber: resin that has hardened but not quite fossilized, and the traditional base of varnish for old wood. Said aloud, it is not far from COBOL.

Pangaea is the single supercontinent of the age of the dinosaurs, before it broke into the continents of today. A monolith, waiting to be split.

A survey walks the ground and maps it before anyone builds; bedrock is what everything stands on.
