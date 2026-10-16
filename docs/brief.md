# Brief

Written 2026-10-16 for the session that starts this project. Nothing is built yet; this is the plan and the decisions still open.

## Why

The hard part of moving a COBOL core system to anything newer is rarely the translation. It is knowing what the system does: record layouts nobody has documented, programs that hand each other fixed-format files in a fixed order with header and trailer records mixed in, business rules spread across conditions, and a job network that only the operations team can draw. Without that understanding, nobody can prove a rewrite produces the same results, so the rewrite keeps the old shape and its old bottlenecks. Copal is a hobby-scale exploration of that problem: derive the understanding from the source, deterministically.

## Two parts

**Pangaea**, the specimen. A small COBOL batch bank, written for this project (AI-written is fine), that compiles and runs with GnuCOBOL 3.2 (`brew install gnucobol`). Roughly ten programs and a few thousand lines:

- Account master: fixed-length sequential records (indexed later, if GnuCOBOL's indexed files behave), layouts in copybooks.
- Transactions: a daily input file of deposits, withdrawals, transfers and fees, with a header and a trailer record carrying counts and control totals.
- Daily batch: validate transactions, sort, post to the master, write rejects and an audit trail, produce the next day's master.
- Month-end: interest accrual and posting, fees, statements, a general-ledger summary.
- A job network file in a small JCL-like format (steps, programs, input and output datasets, conditions on return codes), and a runner that executes it. GnuCOBOL has no JCL; the network file is also what Copal reads to draw the job graph.
- A data generator for test input, and golden outputs checked in.
- Deliberate nasties, because they are what makes real systems hard to analyze: `REDEFINES`, `OCCURS` (and `OCCURS DEPENDING ON` once), level-88 conditions, `PERFORM THRU`, a few `GO TO`s, `COMP-3` packed decimal, `COPY ... REPLACING`, header, detail and trailer records in one file distinguished by a record-type field, and at least one job that only works because an earlier job sorted its input.

**Copal**, the analyzer. Deterministic, no models at this stage.

1. **Parse.** Fixed-format source (sequence area, indicator column, area A and B), `COPY` expansion with `REPLACING`, and the subset of COBOL that Pangaea uses. First decide between writing a subset parser and building on an existing one: evaluate ProLeap (ANTLR-based COBOL 85 parser) and cobol-rekt (an open-source COBOL reverse-engineering toolkit, active as of 2026) before writing anything, and record the decision in this file.
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

## Open decisions for the first session

- Copal's implementation language. Go is the default (Ville's most recent language, simple to ship as one binary); Java would make ProLeap usable directly. Decide after the parser evaluation.
- Graph rendering: Graphviz (`dot`) or hand-laid SVG.
- How far Pangaea goes before Copal starts: suggested order is a minimal daily batch first (two or three programs, one copybook), Copal's layout and structure extraction on that, then grow both together.

## Toolchain

GnuCOBOL 3.2 from Homebrew (`cobc -x` for executables, `-free` is not used: Pangaea is fixed-format on purpose). Go latest stable or the JDK chosen above. `make ci` builds Pangaea, runs its batches against the golden outputs, runs Copal's tests and regenerates the blueprints to compare against the golden ones.
