# Brief

Written 2026-10-16 for the session that starts this project. Nothing is built yet; this is the plan and the decisions still open.

## Why

The hard part of moving a COBOL core system to anything newer is rarely the translation. It is knowing what the system does: record layouts nobody has documented, programs that hand each other fixed-format files in a fixed order with header and trailer records mixed in, business rules spread across conditions, and a job network that only the operations team can draw. Without that understanding, nobody can prove a rewrite produces the same results, so the rewrite keeps the old shape and its old bottlenecks. Copal is a hobby-scale exploration of that problem: derive the understanding from the source, deterministically.

## Two parts

**Pangaea**, the specimen. A small COBOL batch bank, written for this project (AI-written is fine), that compiles and runs with GnuCOBOL 3.2 (`brew install gnucobol`). Roughly ten programs and a few thousand lines:

- Account master: fixed-length sequential records, layouts in copybooks; indexed once the online channel needs lookups by key (see Decisions).
- Transactions: a daily input file of deposits, withdrawals, transfers and fees, with a header and a trailer record carrying counts and control totals.
- Daily batch: validate transactions, sort, post to the master, write rejects and an audit trail, produce the next day's master.
- Online channel: card purchases and ATM withdrawals authorized during the day by COBOL programs behind a thin HTTP server, writing holds that the night's batch posts. Closed during the batch window.
- Month-end: interest accrual and posting, fees, statements, a general-ledger summary.
- A job network file in a small JCL-like format (steps, programs, input and output datasets, conditions on return codes), and a runner that executes it. GnuCOBOL has no JCL; the network file is also what Copal reads to draw the job graph.
- A data generator for test input, and golden outputs checked in.
- Deliberate nasties, because they are what makes real systems hard to analyze: `REDEFINES`, `OCCURS` (and `OCCURS DEPENDING ON` once), level-88 conditions, `PERFORM THRU`, a few `GO TO`s, `COMP-3` packed decimal, `COPY ... REPLACING`, header, detail and trailer records in one file distinguished by a record-type field, and at least one job that only works because an earlier job sorted its input.

**Copal**, the analyzer. Deterministic, no models at this stage.

1. **Parse.** Fixed-format source (sequence area, indicator column, area A and B), `COPY` expansion with `REPLACING`, and the subset of COBOL that Pangaea uses. Copal has its own parser for this subset; see Decisions for why.
2. **Record layouts.** From the copybooks and data divisions: every field with level, name, offset, length, `PIC`, usage and redefinitions. This is the map of the fixed-format interchange between jobs.
3. **Program structure.** Paragraphs and sections, the `PERFORM` graph (including `THRU` ranges and fall-through), and the control flow within paragraphs.
4. **Data flow.** For each paragraph, the fields it reads and writes; field-to-field lineage through `MOVE`, `COMPUTE` and arithmetic verbs.
5. **Business rules.** Conditions on fields, including level-88 names, extracted as decision tables per program.
6. **File I/O and the job network.** Which program opens which file in which mode, joined with the job network file into one graph of steps and datasets: the batch network, drawn.
7. **Blueprints.** Markdown documents per program and per job network, with diagrams generated as SVG (Graphviz or hand-laid). Checked in as golden files for Pangaea.

## Later, not now

- Test data derived from the extracted conditions so that every branch is reached, run through GnuCOBOL to record real outputs: the oracle a rewrite would have to match.
- Data that respects cross-field and cross-record rules (control totals, header-before-detail, references to existing accounts), which format-only generation gets wrong.
- Only then a model, for naming paragraphs and explaining rules in prose, on top of the facts the analysis established.

## Not the goal

- Translating COBOL into Java or anything else.
- Supporting every dialect. Pangaea's subset first; real-world breadth only if the project ever needs it.

## Decisions

- Storage is files, the way mainframe batch stored it. Everything passed between batch steps is a sequential file of fixed-length records. The account master becomes an indexed file (`ORGANIZATION INDEXED`, GnuCOBOL's counterpart of a VSAM KSDS) when the online channel arrives: the batch keeps reading it sequentially, online authorizations read it by account number. Cards and card holds are indexed files too, keyed by card number and by account and hold number. No embedded SQL: it needs a precompiler and a database server on top of GnuCOBOL, and it moves the layouts out of copybooks into a schema, which is a different analysis problem. No IMS or IDMS: neither has a usable open implementation for GnuCOBOL. Which indexed backend the GnuCOBOL build uses (Berkeley DB or another) is checked when the indexed master is built.
- Copal is written in Go and parses COBOL itself. ProLeap and cobol-rekt were both evaluated on a sample with `COPY ... REPLACING`, `REDEFINES`, `OCCURS DEPENDING ON`, level-88s, `COMP-3`, `PERFORM THRU` and `GO TO`. Both parse it correctly, and neither is enough for Copal. Neither maps positions back to the source file and copybook: ProLeap reports lines of the expanded text, cobol-rekt's JSON has no positions at all. Neither computes record layouts: ProLeap has no offsets or lengths, cobol-rekt gives `COMP-3` fields their display size and fails on `OCCURS 0 TO n`. Neither has data flow Copal could trust: ProLeap has none, cobol-rekt records `ADD A TO B` as modifying A. cobol-rekt's output is also not reproducible as shipped, because it uses random IDs. ProLeap has had no code change since January 2024. Either one would add a JVM and a vendored jar to the toolchain, and Copal would still have to write the layout engine, the `PERFORM` range graph and the data flow itself. Copal needs to cover only the subset Pangaea uses, and Pangaea is written alongside it, so a hand-written parser in Go is the smaller dependency and keeps positions exact. The standard library is enough to start with.
- Diagrams are SVG laid out by Copal itself, not by Graphviz. Blueprints are golden files, so the rendering has to be byte-for-byte reproducible, and Graphviz's layout changes between versions (Homebrew and Ubuntu ship different ones). Copal's graphs are small enough for a simple layered layout.
- Pangaea grows in the order of [roadmap.md](roadmap.md): a minimal daily batch first, Copal's layout and structure extraction on that, then both together.

## Besides the code

- The GitHub Actions workflow runs `make ci` in an Ubuntu 26.04 container, the first Ubuntu with GnuCOBOL 3.2 in apt (24.04 has only 3.1.2), so the history shows green builds and the check exists for branch protection once the repository is public.
- Keep the README's status line current as parts start working.

## Debt

- Fixtures and golden files hold records as text, and a record with packed decimal (`COMP-3`) fields as one line of hex. That round-trips, but a diff of it reads poorly; fields decoded through Copal's record layouts would read better.
- `structure` finds `PERFORM` and `GO TO` statements by scanning a paragraph's tokens, and decides whether a paragraph runs on into the next from its last sentence's tokens: it ends unconditionally only in `GO TO`, `STOP RUN`, `GOBACK` or `EXIT PROGRAM` outside any condition, so an `IF` whose branches both jump still counts as running on. Data flow needs real statement parsing, and `structure` should move onto it then.
- A business day whose batch fails leaves the generations its steps wrote cataloged, and the bank's next run counts from them. The bank's history, and with it the console, assumes every day went through; a failed day needs a restart procedure that removes them.

## Toolchain

GnuCOBOL 3.2 from Homebrew (`cobc -x` for executables, `-free` is not used: Pangaea is fixed-format on purpose). Go, latest stable. `make ci` builds Pangaea, runs its batches against the golden outputs, runs Copal's tests and regenerates the blueprints to compare against the golden ones.
