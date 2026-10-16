# Roadmap

The order in which Copal and the First Bank of Pangaea are built. [brief.md](brief.md) has the reasons and the open decisions; this file has the sequence. Each milestone ends with something that runs, a `make` target that shows it, and `make ci` green.

The two parts grow together: Pangaea gets a little bigger, Copal learns to read what was added, and the blueprints are regenerated as golden files. Somewhere in the middle Pangaea also becomes a bank you can watch: a dataset large enough to feel real, a business calendar, a reset button, and a console in the browser.

## 1. Groundwork

- Evaluate ProLeap and cobol-rekt against a few lines of fixed-format COBOL; decide Copal's language and record it in the brief.
- Toolchain checks in the Makefile: GnuCOBOL 3.2, the chosen language, Graphviz if chosen.
- `make ci` that builds and tests whatever exists, and a GitHub Actions workflow that runs the same target on Ubuntu.

Done when: a pushed branch shows a green build.

## 2. Pangaea: a minimal daily batch

- One copybook for the account master, one for the transaction file (header, detail and trailer records, told apart by a record-type field).
- Three programs: validate (control totals, record order, known accounts), sort, post (new master, rejects, audit trail).
- The job network file in its small JCL-like format, and a runner that executes it and stops on bad return codes.
- A small hand-written fixture: a dozen accounts, one day of transactions with a reject or two. Its outputs are the first golden files.

Done when: `make pangaea-daily` runs the batch on the fixture and `make ci` compares the outputs with the golden ones.

## 3. Copal: layouts and structure

- Fixed-format reading: sequence area, indicator column, comments and continuation lines, `COPY` expansion.
- Record layouts: level, name, offset, length, `PIC`, usage, redefinitions, as data (JSON) and as a Markdown table.
- Program structure: paragraphs, sections, the `PERFORM` graph.
- The first blueprints for Pangaea, checked in.

Done when: `make blueprints` regenerates them and a change to a copybook shows up as a diff in the golden blueprints.

## 4. A living bank

Pangaea stops being a fixture and becomes a running bank.

- A seeded data generator: about two thousand customers and three thousand accounts (current and savings, a few closed or frozen), and for each business day a few thousand transactions with realistic shapes: payday deposits, rent and utility debits, transfers between customers, the occasional overdraft, and a percent or two that the validator has to reject.
- A business calendar. The bank has a current business date of its own, never the wall clock, so every run is reproducible: day N from seed S is always the same.
- The bank's state lives under `pangaea/run/` (ignored by git): the current master, the business date, and one snapshot per closed day.
- Make targets: `bank-reset` regenerates the starting point from the seed, `bank-day` generates the day's transactions and runs the daily batch, `bank-days N` runs several in a row.

The hand-written fixture stays as the CI golden set; the generated bank is for watching, and its determinism is checked by a test that runs a few days twice and compares.

Done when: `make bank-reset bank-days N=5` leaves five closed business days with snapshots, the same every time.

## 5. The bank console

A small web UI that shows the bank's state and how each batch changed it.

- It reads Pangaea's fixed-format files through the record layouts Copal extracted in milestone 3, not through a second hand-written decoder. The console is Copal's first user.
- Views: accounts and balances with search, one account's history across days, the day's transactions and rejects, the audit trail, control totals per run.
- Before and after: pick two business days and see what changed, balances moved, accounts opened or closed, totals reconciled.
- Buttons for the same operations as the make targets (reset, run a day), calling the same scripts. The console is a window onto the bank, not part of it: it never writes Pangaea's files itself.
- Served locally by `make console`; plain server-rendered pages with a little script, no frontend build.

Done when: you can reset the bank, run a day from the browser, and see the master change.

## 6. Pangaea grows, Copal follows

Pangaea reaches the size in the brief, a step at a time, each step with Copal reading it:

- Month-end: interest accrual and posting (`COMP-3` amounts), monthly fees, statements, a general-ledger summary. `bank-days` runs month-end on the last business day of the month.
- The nasties: `REDEFINES`, `OCCURS` and `OCCURS DEPENDING ON`, level-88 conditions, `PERFORM THRU` with a few `GO TO`s, `COPY ... REPLACING`, and a job that only works because an earlier one sorted its input.
- Copal: data flow per paragraph and field lineage; business rules as decision tables; file I/O joined with the job network into one graph, drawn as SVG.

The console gains a page per run: the job network with each step's return code and record counts, and links from steps to the blueprints of their programs.

Done when: Pangaea has its ten or so programs, a month-end in the console, and a blueprint for every program and for the job network.

## 7. The card channel

Worth having, but late and small. A bank that only has a batch misses the coupling that makes real core systems hard to replace: an online day that writes what the night's batch reads, a balance that means different things at different hours, and a batch window during which the online side has to stand still. That coupling is exactly the kind of undocumented contract Copal is meant to surface, and it makes the console worth watching during the day.

- A debit card authorization API over HTTP: authorize, reverse, and a balance enquiry.
- The authorization rule lives in COBOL, as a called program in the style of an online transaction, so Copal reads it like the rest of the bank. A thin server around it does the HTTP and nothing else.
- Authorizations do not touch the master. They append holds to an online log, and the available balance is the master balance minus open holds.
- At end of day the online log becomes a transaction file with header and trailer, the daily batch posts it, and holds are released.
- During the batch window the API answers that the bank is closed.
- A traffic simulator that swipes cards at a configurable rate, and live updates in the console (server-sent events) so balances and holds move while you watch.

Open decisions for this milestone: how the server calls COBOL (a process per request, or libcob in-process), and whether the online log is a sequential file appended under a lock or one of GnuCOBOL's indexed files.

Done when: with the simulator running, the console shows available balances dropping during the day, the API closes for the batch, and the next morning the card transactions appear in the master.

## Later

From the brief, once the above stands: test data derived from the extracted conditions so every branch is reached, run through GnuCOBOL as the oracle a rewrite has to match; data that respects cross-record rules; and only then a model for naming and prose.
