# Copal

*Reading the code that still runs the books.*

Copal is a study of one problem: understanding a legacy COBOL batch system well enough to replace it. The hard part of such a move is rarely the translation; it is knowing what the system does, when nobody can draw the job network or say which rules are in the code. Copal derives that understanding from the source, deterministically, and tests it on three COBOL systems written for the purpose. It does not translate COBOL into anything.

## The parts

| Directory | Part | What it is |
| --- | --- | --- |
| [survey/](survey/) | **Survey** | The analyzer. Reads a COBOL system and writes its blueprints: record layouts, program structure, data flow, business rules as decision tables, and the job network. Serves the same analysis in the browser, and finds test cases that reach the rules. |
| [pangaea/](pangaea/) | **The First Bank of Pangaea** | The first specimen: a small working bank in IBM COBOL and JCL. A daily batch, a month-end, and an online day of card and ATM authorizations. |
| [laurasia/](laurasia/) | **Laurasia Card Services** | The second specimen: a card processor in Micro Focus COBOL, free format, run by shell scripts. It shows where the survey assumed Pangaea. |
| [gondwana/](gondwana/) | **Gondwana Life Insurance** | The third specimen: a life insurer in Fujitsu NetCOBOL, in Shift_JIS with Japanese names and comments, yen and era dates, run by Windows batch files. It shows where the survey assumed ASCII, cents and the Western calendar. |
| [bedrock/](bedrock/) | **Bedrock** | What runs the specimens: compiling and running COBOL programs, jobs and scripts, and keeping datasets and their catalog. |
| [console/](console/) | **The console** | A web window that runs the bank, watches its online day, and shows each run's jobs as the survey draws them. The one part that knows both the bank and the survey. |
| [cobolith](https://github.com/vlumi/cobolith) | **Cobolith** | The library everything reads COBOL with: source, copybooks, record layouts and fields, and JCL, in the IBM, Micro Focus, Hitachi and Fujitsu dialects. Written here, now a repository of its own. |

![The parts of Copal and what imports what](docs/diagrams/architecture.svg)

The dependencies run one way, and `make ci` checks it: the specimens know nothing of the survey, and only the console knows both. All three specimens compile and run with GnuCOBOL, and their outputs are golden files. So are the survey's blueprints of them: a change to the analyzer shows as a diff in what it says about a real system.

## What you can run

Everything is a make target, and plain `make` lists them. You need GnuCOBOL 3.2 (`brew install gnucobol`), Go, Node for the Markdown linter, and `golangci-lint` for `make ci`.

### The bank

| Target | What it does | Where to look |
| --- | --- | --- |
| `make pangaea-daily` | Runs the daily batch on a small fixture: the job log and every step's report | the terminal |
| `make bank-reset`, `make bank-days N=5` | Opens a bank of about 3,000 generated accounts and runs business days, one line each | `pangaea/run/bank/` |
| `make console` | The bank in the browser: accounts, each day's jobs drawn as networks with return codes and record counts, what the days changed | <http://127.0.0.1:4510/> |
| `make simulate` | Sends card and ATM traffic to the console's online channel | the console's **Online** page |

### The card processor

| Target | What it does | Where to look |
| --- | --- | --- |
| `make laurasia-night` | Runs Laurasia's night on its fixture: the clearing file in, matching, pricing, posting, disputes, settlement and the billing cycle, with each step's log and every program's report | the terminal, and `laurasia/run/night/data/` |
| `make cards-open`, `make cards-days N=30` | Opens a card processor of about 1,000 generated cards and runs days, one line each: the night on business days, the billing cycle every day | `laurasia/run/processor/data/` |

### The life insurer

| Target | What it does | Where to look |
| --- | --- | --- |
| `make gondwana-month` | Runs Gondwana's March on its fixture: the system opened, three business days, the premium billing and the collection agent's results, the month-end close and the year-end dividends, with each batch file's log and every program's report, transcoded for the terminal | the terminal, and `gondwana/run/month/data/` |
| `make life-open`, `make life-days N=261` | Opens a life insurer of about 2,000 generated policies and runs business days, one line each: the daily run, and the billing, the agent's results, the close and the year end on their days | `gondwana/run/insurer/data/` |

### The survey

| Target | What it does | Where to look |
| --- | --- | --- |
| `make blueprints` | Regenerates the survey of the three specimens | [blueprints/pangaea/](blueprints/pangaea/README.md), [blueprints/laurasia/](blueprints/laurasia/README.md), [blueprints/gondwana/](blueprints/gondwana/README.md) |
| `make explore-pangaea`, `make explore-laurasia`, `make explore-gondwana` | The same analysis as a web app: every job, step, dataset, program, paragraph and field a page, and a field's values followed across programs on the diagrams | <http://127.0.0.1:1192/> |
| `make oracle` | Finds, for each program of the specimens, input data that reaches the rules of its decision tables, by running it; replays the batches on those cases | [oracle/pangaea/](oracle/pangaea/README.md), [oracle/laurasia/](oracle/laurasia/README.md), [oracle/gondwana/](oracle/gondwana/README.md) |
| `make conform-pangaea PROGRAM=TXPOST WITH=bin/txpost` | Runs a replacement for a program, any executable, on the program's recorded cases, and names the first field it gets wrong; `conform-laurasia` and `conform-gondwana` likewise | the terminal |
| `make conform-demo` | The same on BANKDATE rewritten in Go, which conforms on all its cases | the terminal |

### Checks

| Target | What it does |
| --- | --- |
| `make ci` | What the hosted build runs: linters, the dependency check, every test and golden file, the oracle's cases replayed |
| `make pangaea-golden`, `make laurasia-golden`, `make gondwana-golden` | Rewrite a specimen's golden outputs after a deliberate change to it |

## What the survey writes

| For each | The blueprint has |
| --- | --- |
| Program | Its own header comment; files; the flow of its paragraphs, drawn; what each paragraph reads and writes; where every value comes from; its records field by field |
| Copybook | Its records, with each field's kind, and the programs that use it |
| Job | Its steps and the datasets between them, drawn; each step's DD statements bound to the program's files; findings where the job and its programs disagree |
| Paragraph | A decision table: each path through its conditions, and the actions it takes |
| Entity | Who reads, writes and decides on each field across the programs; the state machine of each code field; who owns it |
| Field | Its trail back to the inputs its values come from, and forward to everything they reach, through programs and job steps |
| Rule | A recorded case: input data that reaches it, and what the program wrote and returned; what a replacement for the program has to match |

[docs/blueprints.md](docs/blueprints.md) explains how to read them. Every fact links to the source line it came from, and no language model is involved.

## By the numbers

| | |
| --- | --- |
| Specimens | 3, in three dialects: 46 programs, about 10,000 lines of COBOL, 41 copybooks, 4 jobs in JCL, 5 shell scripts and 7 batch files |
| Blueprints | 153 pages and 63 diagrams, regenerated and compared by `make ci` |
| Decision rules | 497 in the bank and the card processor, of which 224 recorded cases reach 468; 14 are marked with the reason the search cannot reach them, 15 are left. The insurer's 677, counting a called program's by each caller, of which 162 cases reach 318 |
| Cobolith | 0.8.3: four dialects, the Report Writer, and 439 of the 459 programs of NIST's COBOL-85 test suite read |

## A quick tour

1. `make toolchain`: check that the tools are there.
2. `make pangaea-daily`: the bank in one screen; `make laurasia-night`: the card processor's night in one; `make gondwana-month`: the life insurer's March, in Japanese.
3. `make bank-reset` and `make bank-days N=5`: five business days of the generated bank. `make cards-open` and `make cards-days N=30`: a month of the card processor, weekends included. `make life-open` and `make life-days N=261`: a fiscal year of the life insurer, its dividends on the last day.
4. `make console`, then open an account, a day's jobs, and **Compare**.
5. `make simulate` from another terminal: the **Online** page fills with authorizations. Run a day from the console: the channel closes for the batch, and the next morning the purchases and withdrawals are in the master.
6. Read a blueprint: [TXPOST](blueprints/pangaea/programs/TXPOST.md), its [decision tables](blueprints/pangaea/rules/TXPOST.md), and the [DAILY job](blueprints/pangaea/jobs/DAILY.md).
7. `make explore-pangaea`: search for `AM-BALANCE`, open it in TXPOST, and follow its values back to the files that come from outside. `make explore-gondwana`: the same, in Japanese, with the Report Writer's lines among the records.
8. Read an oracle case: [TXPOST's](oracle/pangaea/TXPOST/README.md), rule by rule.
9. `make conform-demo`: a rewrite of BANKDATE in Go checked against its cases. Break the rewrite, in `survey/internal/oracle/testdata/rewrites/bankdate/`, and run it again: the report names the field.
10. `make ci`.

## Reading more

| Document | What it covers |
| --- | --- |
| [docs/architecture.md](docs/architecture.md) | How the parts fit, the dependency rule, and the golden files |
| [docs/blueprints.md](docs/blueprints.md) | How to read what the survey writes |
| [survey/README.md](survey/README.md) | The survey's pipeline, analysis by analysis, and what it does not do |
| [pangaea/README.md](pangaea/README.md), [pangaea/docs/](pangaea/docs/README.md) | The bank, and its own operations guide in its operators' voice |
| [laurasia/README.md](laurasia/README.md), [laurasia/docs/](laurasia/docs/operations.md) | The card processor and its night |
| [gondwana/README.md](gondwana/README.md), [gondwana/docs/](gondwana/docs/operations.md) | The life insurer and its month |
| [docs/roadmap.md](docs/roadmap.md), [docs/brief.md](docs/brief.md) | The order the work was done in, and the plan and decisions it started from |

## The names

Copal is young amber: resin that has hardened but not quite fossilized, and the traditional base of varnish for old wood. Said aloud, it is not far from COBOL. Pangaea is the single supercontinent before it broke into the continents of today: a monolith, waiting to be split; Laurasia and Gondwana are the two it split into. A survey walks the ground and maps it before anyone builds; bedrock is what everything stands on; a cobolith is the stone it is cut from.

## License

MIT, see [LICENSE](LICENSE). The specimens, their programs, data and documents included, are fiction written for this project and covered by the same license.
