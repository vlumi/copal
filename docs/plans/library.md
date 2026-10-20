# Plan: cobolith, bedrock's COBOL as a library

Bedrock reads COBOL: `cobol` lexes and parses source and copybooks, `layout` lays out records, `record` reads and writes their fields, and `jcl` reads the jobs that run them. Those four become cobolith, a Go library of its own for anyone reading COBOL: first a module within this repository, then a repository of its own once its API stops moving with the survey. The name is COBOL and a stone, in the family of bedrock and the specimens, and has COBOL in it for anyone looking for one.

Read [../roadmap.md](../roadmap.md), milestone 15, first; the ground rules of AGENTS.md apply.

## Ground rules

- **It stands alone.** Cobolith imports nothing of copal, and its tests read only its own `testdata`, written for it, never the specimens. `make ci` checks both.
- **One module in this repository until the API settles.** A change to the library and the change it makes to the survey's blueprints, and to the specimens' oracle, land in one commit, as the goldens rule wants. It moves to a repository of its own when a change to it rarely needs one to the survey alongside: by then most of its work is dialects, which the survey does not wait on.
- **History moves with it.** The packages move with `git mv` now and with `git filter-repo` later, so its history starts where bedrock's did.
- **Honest about coverage.** Below 1.0 and marked pre-production, its README says what it reads, in which dialects, what it has been run against, and what it does not yet read.

## What moves and what stays

| Package | Goes to | Why |
| --- | --- | --- |
| `cobol` | cobolith | Reading COBOL source, copybooks and dialects |
| `layout` | cobolith | Record layouts from data description entries |
| `record` | cobolith | Fields of records through their layouts |
| `jcl` | cobolith | The jobs COBOL batches run by on a mainframe |
| `script` | stays in bedrock | A subset of POSIX sh, Laurasia's batch control, not COBOL's |
| `dataset`, `runner`, `run` | stay in bedrock | This project's harness: running programs, jobs and catalogs |

Bedrock keeps what runs programs, and imports cobolith as survey, Pangaea and Laurasia do.

## Parts

1. **A module of its own.** `cobolith/` at the root, with its own `go.mod` for `github.com/vlumi/cobolith`. Copal's module requires it, with a `replace` to the directory until it moves. The four packages move into it with their tests, imports change across the repository, and the goldens do not. `record`'s and `jcl`'s tests that read Pangaea's copybooks and jobs get fixtures of their own. The make targets test and lint both modules, and the dependency check gains cobolith: it imports nothing of copal, and bedrock may import it.

   As built: the four packages moved with their history, and no blueprint, oracle case or golden output changed. `record`'s tests read an account record and a record of three FILLERs of their own, and `jcl`'s a job of two steps. Bedrock keeps `script`, `dataset`, `runner` and `run`, and the architecture diagram shows every part importing both.
2. **The public API.** Every exported name is either used outside the library, documented and kept, or made unexported. Each package gets a package comment that says what it reads and in which dialects, and runnable examples: parse a program and walk its statements, lay out a copybook's record, read a field of a record. Errors carry the file and line they are about. The API is what the survey, the specimens and the console need, and no more.

   As built: of the exported names, every type is in the signature or a field of another the survey uses, and the constants are the values of exported fields, so the API was nearly what Copal needs already; `ParseStatements`, plumbing the parser alone calls, is now unexported. Each package has an example that runs as a test. Writing them found what the parser does not read yet, for part 3 and the README: statements before a program's first paragraph, `DISP` with a normal or abnormal end such as `(NEW,PASS)`, and a new dataset without its record length.
3. **Tests of its own.** Unit tests for each construct, in each dialect it reads, on fixtures written for the library, with no specimen behind them; their coverage, measured by `go test -cover`, goes in its README. A public conformance suite would show what it does not read: the NIST COBOL-85 suite GnuCOBOL tests with is the candidate, if its license lets it be used and its programs kept as fixtures; that is checked before any of it is copied.
4. **Its README.** What it reads, by dialect and source format; what it has been run against; what it does not read yet; how stable the API is.
5. **A repository of its own.** When the API holds still: `git filter-repo` to its own history, a repository at `github.com/vlumi/cobolith`, tagged `v0.1.0`, and copal requiring that tag in place of the `replace`. From then a change to the library reaches copal as a version bump, and the commit that bumps it carries the goldens it changes.

## Then, in its own repository

The dialects most COBOL is written in, each a pull request with its own tests: GnuCOBOL's own, ACUCOBOL-GT, the 2002 and 2014 standards' additions, `REPLACE` beside `COPY REPLACING`, and `EXEC SQL` and `EXEC CICS` blocks kept whole as opaque statements, rather than failing the parse. Each one the survey can use as it lands.

## Decided

- The name, cobolith: no project, package or repository found by that name.
- In this repository first, as its own module; its own repository when the API settles.
- `script` stays in bedrock with the harness.

Each part one pull request.
