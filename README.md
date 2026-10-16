# Copal

*Reading the code that still runs the books.*

Copal reads a legacy COBOL batch system and draws its blueprints: the layout of every record, the structure of every program, which fields each step reads and writes, the business rules buried in conditions, and the job network that ties the programs together through fixed-format files. All of it is derived deterministically from the source.

It is developed against **the First Bank of Pangaea**, a small but working COBOL batch bank in [pangaea/](pangaea/): accounts, transactions, a daily batch and a month-end run, compiled and executed with GnuCOBOL.

Copal does not translate COBOL into anything. It produces the understanding a translation needs, and later the tests a rewrite has to pass.

**Status:** Pangaea runs as a bank: about three thousand generated accounts, a daily batch of validation, interest accrual, sorting and posting in three passes for transfers, month-end interest, fees, statements and a general-ledger summary, and a business calendar, day after day (`make bank-reset bank-days N=20`). `make console` opens it in the browser: accounts, each day's files and reports, and what the days changed in the master, read through the record layouts Copal computes, with buttons to reset the bank and run days. Small fixtures check the batch against golden outputs. Copal reads it and writes [blueprints](blueprints/pangaea/README.md) of every program and copybook: each program's own description, record layouts with positions and lengths, files, and its paragraphs with a diagram of how control moves between them; and of every job, the network of its steps and the datasets they read and write, with findings where a job and its programs disagree. `make blueprints` regenerates them. [docs/brief.md](docs/brief.md) has the plan and [docs/roadmap.md](docs/roadmap.md) the order of work.

## The names

Copal is young amber: resin that has hardened but not quite fossilized, and the traditional base of varnish for old wood. Said aloud, it is not far from COBOL.

Pangaea is the single supercontinent of the age of the dinosaurs, before it broke into the continents of today. A monolith, waiting to be split.
