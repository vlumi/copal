# Roadmap

The order in which Survey and the First Bank of Pangaea are built. [brief.md](brief.md) has the reasons and the open decisions; this file has the sequence. Each milestone ends with something that runs, a `make` target that shows it, and `make ci` green.

The two parts grow together: Pangaea gets a little bigger, the survey learns to read what was added, and the blueprints are regenerated as golden files. Somewhere in the middle Pangaea also becomes a bank you can watch: a dataset large enough to feel real, a business calendar, a reset button, a console in the browser, and an online day of card purchases and ATM withdrawals between the nightly batches.

## 1. Groundwork (done)

- Evaluate ProLeap and cobol-rekt against a few lines of fixed-format COBOL; decide the survey's language and record it in the brief.
- Toolchain checks in the Makefile: GnuCOBOL 3.2, the chosen language, Graphviz if chosen.
- `make ci` that builds and tests whatever exists, and a GitHub Actions workflow that runs the same target on Ubuntu.

Done when: a pushed branch shows a green build.

## 2. Pangaea: a minimal daily batch (done)

- One copybook for the account master, one for the transaction file (header, detail and trailer records, told apart by a record-type field).
- Three programs: validate (control totals, record order, field checks), sort, post (known accounts, account status and overdraft limits; new master, rejects, audit trail). Deposits, withdrawals and fees; transfers come with the living bank, posted as two legs.
- The job network file in its small JCL-like format, and a runner that executes it and stops on bad return codes.
- A small hand-written fixture: a dozen accounts, one day of transactions with a reject or two. Its outputs are the first golden files.

Done when: `make pangaea-daily` runs the batch on the fixture and `make ci` compares the outputs with the golden ones.

## 3. The survey: layouts and structure (done)

- Fixed-format reading: sequence area, indicator column, comments and continuation lines, `COPY` expansion.
- Record layouts: level, name, offset, length, `PIC`, usage, redefinitions, as data (JSON) and as a Markdown table.
- Program structure: paragraphs, sections, the `PERFORM` graph.
- The first blueprints for Pangaea, checked in.

Done when: `make blueprints` regenerates them and a change to a copybook shows up as a diff in the golden blueprints.

## 3b. Blueprints people can read (done)

- The documentation already in the source: the header comment of each program and copybook, the comments before paragraphs and fields, and identification entries such as `AUTHOR` and `DATE-WRITTEN`.
- Pangaea's programs carry what long-lived code carries: identification entries, a modification history in the header banner, a comment before each paragraph.
- A diagram per program: its paragraphs and the `PERFORM`s between them, as SVG laid out by the survey, in the program's blueprint. `GO TO`, fall-through and `THRU` ranges joined it in milestone 6.

Done when: every program blueprint opens with the program's own description and shows its paragraph diagram.

## 4. A living bank (done)

Pangaea stops being a fixture and becomes a running bank.

- Transfers, posted in legs: the debit on the sender, then a second sort and post for the credit on the receiver, and a third for credits that cannot post and go back. The same two programs run three times in one job.
- Generation data groups in the runner, so every day's master, control record, audit trail and rejects stay cataloged as a history, and temporary datasets for everything in between.
- A business calendar. The bank has a business date of its own, never the wall clock, and BANKDATE moves it to the next weekday as the batch's last step, so every run is reproducible: day N from seed S is always the same.
- A seeded data generator: about two thousand customers and three thousand accounts (current and savings, a few closed or frozen), and for each business day the branches' transaction file: salaries on the 25th, rent on the first business day of the month, savings after payday, cash, deposits, transfers between customers, fees, and a few percent that validation or posting rejects.
- The generator reads and writes records through the layouts the survey computes from the copybooks, with no record format written out by hand.
- The bank's state lives under `pangaea/run/bank` (ignored by git). `make bank-reset` regenerates the starting point from a seed, `make bank-days N=5` runs the next five business days.

The hand-written fixtures stay as the CI golden set; the generated bank is for watching, and a test runs three days in two banks from the same seed and compares them byte for byte.

## 5. The bank console (done)

A small web UI that shows the bank's state and how each batch changed it.

- It reads Pangaea's fixed-format files through the record layouts the survey extracted in milestone 3, not through a second hand-written decoder. The console is the survey's first user.
- Views: accounts and balances with search, one account's history across days, the day's transactions and rejects, the audit trail, control totals per run.
- Before and after: pick two business days and see what changed, balances moved, accounts opened or closed, totals reconciled.
- Buttons for the same operations as the make targets (reset, run a day), calling the same scripts. The console is a window onto the bank, not part of it: it never writes Pangaea's files itself.
- Served locally by `make console`; plain server-rendered pages with a little script, no frontend build.

Done when: you can reset the bank, run a day from the browser, and see the master change.

## 6. Pangaea grows, the survey follows (done)

Pangaea reaches the size in the brief, a step at a time, each step with the survey reading it:

- Month-end (done): interest accrual and posting (`COMP-3` amounts), monthly fees, statements, a general-ledger summary. Interest accrues every day; the first business day of a month posts the month before's interest and fees ahead of its own transactions, and a month-end job after the daily batch prints the statements and the general-ledger summary.
- The nasties (done): `REDEFINES`, `OCCURS` and `OCCURS DEPENDING ON`, level-88 conditions, `COPY ... REPLACING`, a job that only works because its input comes in master order, and account maintenance in 1988 style, with `PERFORM THRU` ranges, `GO TO`s and fall-through, which the survey follows into its diagrams.
- the survey (done): file I/O joined with the job network into one graph, drawn as SVG, a statement parser, data flow per paragraph and field lineage within a program, and business rules as decision tables.

The console gains a page per run: the job network with each step's return code and record counts, and links from steps to the blueprints of their programs (done).

Done when: Pangaea has its ten or so programs, a month-end in the console, and a blueprint for every program and for the job network.

## 7. The online channel: cards and ATMs (done)

A bank that only has a batch is not a realistic bank. Customers pay by card and take cash from ATMs all day, and that online day is coupled to the night's batch in ways that make real core systems hard to replace: the online side writes what the batch reads, a balance means different things at different hours, and during the batch window the online side has to stand still. That coupling is exactly the kind of undocumented contract the survey is meant to surface.

- A card file: cards linked to accounts, with status (active, blocked, expired) and daily limits for purchases and for cash.
- An authorization API over HTTP for two kinds of terminal. Card purchases: authorize and reverse. ATMs: cash withdrawal and balance enquiry.
- The authorization rules live in COBOL, as called programs in the style of online transactions, so the survey reads them like the rest of the bank: card status, available balance, overdraft limit, the daily purchase and cash limits, and an ATM fee for withdrawals beyond the free ones each month.
- A thin server around them does the HTTP and nothing else. It passes requests to a COBOL dispatcher that runs for the whole online day and calls the authorization programs, as a transaction monitor would.
- The account master stays sequential. Each night the batch loads the new master into an indexed online account file keyed by account number, which authorizations read single accounts from. The card file is indexed by card number, loaded the same way from a sequential card master.
- Authorizations do not update the master. They write holds to an indexed file keyed by account and hold number and add them to the online account's holds total, and the available balance is the balance plus the overdraft limit less the holds. Each card keeps its running totals for the day.
- At end of day the holds are read in key order into a transaction file with header and trailer, the daily batch posts it and loads the online files afresh, with no holds. The cards' totals are for the date they carry, and start over with a new day's first authorization.
- During the batch window the API answers that the bank is closed.
- A traffic simulator that makes card purchases and ATM withdrawals at a configurable rate, following each card's generated spending habit, with live updates in the console (server-sent events) so balances and holds move while you watch.

Done when: with the simulator running, the console shows available balances dropping during the day and declines for cards over their limits, the API closes for the batch, and the next morning the card and ATM transactions appear in the master.

## 8. The explorer (done)

Static blueprints are good for reading one program or job at a time, but understanding a system means moving through it: from the job network to a step, its program, a paragraph, a record, a field, and following a value across programs through the datasets between them. The explorer is a web app for that, a second reader of the same analysis the blueprints are written from.

- `survey explore` serves the analysis of any system the survey reads, not just Pangaea, locally, in the browser (done).
- Every name is a link: jobs, steps, datasets, programs, paragraphs, records, fields, copybooks. Search finds any of them, and history goes back and forth (done).
- The job networks and program flows are the survey's own diagrams, with pan and zoom (done).
- Lineage across programs: pick a field and see where its values come from, through the statements of each program and through the datasets that pass records from step to step, highlighted on the diagrams (done). Fields are joined across programs by the bytes the writer wrote and the reader reads, which is what a shared copybook comes down to, and a sort passes records on unchanged.
- Server-rendered pages with a little script, like the bank console; no frontend build (done).

Done when: starting from a field in a report, you can follow it back through every program and job step to the input it came from, on the diagrams.

## 9. Design views

The explorer shows the system the way the code is organized. Designing a replacement needs it organized the way the business is. The views, with their plan in [plans/design-views.md](plans/design-views.md):

- Kinds, visible everywhere: amounts, codes, keys, dates and text among fields; main and called programs; external, internal, output and temporary datasets; shown as badges in the explorer, as columns in the blueprints, and as looks of their own in the diagrams (done).
- Rules by field: every decision on a field across every program, and a page per field of a shared record.
- Entities: a page per copybook record with its datasets, its users, who reads and writes each field, and for each code field its state machine, derived from the decision tables.
- Ownership: programs grouped by the entities they write, drawn; the first draft of a service decomposition.
- The boundary: the system's inputs and outputs, files, reports and streams, in one place.
- Forward lineage: from a field to everything its values reach, the mirror of the trail back, on the same diagrams.
- The order of jobs, derived from the generations they write and read.

Done when: starting from the entities page, you can name the system's data model, its state machines, who owns each entity, its contracts with the outside, and what a change to an input field touches, without opening a program.

## Later

The explorer's lineage goes back, not yet forward: where a field's values go. It does not yet follow values into a called program and back through its LINKAGE SECTION, nor from SYSIN through ACCEPT, so the online day's holds stop at AUTHSRV. A program's outputs keep only the first record type it writes to a file.

From the brief, once the above stands: test data derived from the extracted conditions so every branch is reached, run through GnuCOBOL as the oracle a rewrite has to match; data that respects cross-record rules; and only then a model for naming and prose.
