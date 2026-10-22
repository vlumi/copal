# Agent instructions

Copal is a project in four parts, and three specimens. Survey, in `survey/`, reads legacy COBOL batch systems and produces design blueprints from them: record layouts, program structure, data flow, business rules and the job network. The analysis is deterministic; no language model is involved until the deterministic core is solid. The First Bank of Pangaea, a small COBOL bank in `pangaea/`, is the specimen the survey is developed and tested against; Laurasia Card Services, a card processor in `laurasia/`, written in another dialect and run by shell scripts, is the second, which shows where the survey is coupled to the first; Gondwana Life Insurance, a life insurer in `gondwana/`, written in Fujitsu NetCOBOL in Shift_JIS and run by Windows batch files, is the third. None of them knows anything of the survey. Cobolith, a library of its own at github.com/vlumi/cobolith, reads COBOL, lays out its records and reads its jobs; `bedrock/` runs programs, jobs and scripts; survey and the specimens use both, and `console/` is the one part built on all of them. README.md describes the project; docs/brief.md has the plan and the open decisions.

## Rules

- Never mention where this project came from, any company or product it might be compared to as its origin, or any hiring process, in code, documentation, commit messages or anywhere else in the repository. Technical references to well-known tools and systems (GnuCOBOL, ProLeap, IBM mainframes) are fine where they explain a choice.
- Survey is deterministic. Every blueprint fact is derived from the source by parsing and analysis, and is reproducible. Models may come later, only for naming and prose, and only on top of facts the analysis already established.
- The specimens are fiction written for this project. Never copy real proprietary code, layouts or data into them.
- The specimens run: every program compiles with GnuCOBOL, Pangaea's in IBM dialect, Laurasia's in Micro Focus dialect and Gondwana's in Fujitsu dialect, and their batches produce outputs that are checked in as golden files. A change to a specimen updates its golden outputs in the same commit.
- A coupling to one specimen that the other shows is fixed generically, never by special-casing either.
- Survey's blueprints of both specimens are golden files too, so a change in the analyzer shows as a diff in the blueprints.
- The dependencies run one way, and `make ci` checks it: `bedrock` imports nothing of Copal but cobolith, the library; `survey`, `pangaea`, `laurasia` and `gondwana` import `bedrock` and `cobolith`, never each other; only `console` imports the rest. Survey's tests read the specimens as their fixtures, which is data, not an import.
- Code is organized by feature, not by layer.
- All text is English with en-US spelling: code, comments, documentation, commit messages, and the specimens' programs and data. Comments are minimal. The exception is Gondwana, whose point is a Japanese shop's COBOL: its programs, copybooks, batch files and data are in Japanese, comments, data names in half-width katakana, literals and reports alike, in Shift_JIS as its compiler reads them; its documentation and the Go around it are English like the rest.
- Gondwana's sources and data are Shift_JIS in the working tree and UTF-8 in the repository, by `.gitattributes`; edit them with a tool that keeps the encoding, and keep every fixed-format line within 72 bytes, a two-byte character counting two.
- Markdown paragraphs are one line each; the editor wraps. `.markdownlint.jsonc` encodes this.
- Diagrams are generated as SVG or PNG, not Mermaid: the editor in use does not render Mermaid.
- Routine operations are make targets with a `##` description, so plain `make` lists them; scripts do the work, the Makefile names it.

## Commits and pull requests

- The history reads as what actually happened. Commits are granular: each one an atomic change, or a few closely related ones, that builds and passes on its own. No squashing. Something found later, a missed case or a mistake in an earlier commit, gets its own commit saying what it fixes, rather than being folded back into history.
- Commit messages describe the change, not the process.
- A pull request groups the commits of one functional change. Every chunk of work is a branch and a pull request against `main`, started from current `main`. Pushed commits are not rewritten: no amending, squashing or force-pushing. `main` is never rewritten and never committed to directly.
- Run `make ci` before opening or updating a pull request and only push when it passes; a change to Markdown alone needs none.
- Merge commits only, so the by-concern commits stay visible on `main`; branches are deleted on merge. While the repository is private, GitHub's free plan allows neither auto-merge nor branch protection, so merge with `gh pr merge --merge` once the local `make ci` has passed. When the repository goes public, enable auto-merge, protect `main` with the hosted CI job as a required check, and merge with `gh pr merge --auto --merge`, which then waits for it.

## Code

- This is a showcase: the code should read well to a stranger. Idiomatic for its language, small functions with names that say what they do, clear package boundaries, no dead code, no speculative abstractions.
- Build iteratively: working first, then dedicated cleanup commits. Do not gold-plate mid-feature; record debt worth a cleanup pass in docs/brief.md.
- Every analysis step is covered by tests, and by the golden blueprints of the specimens. Linters and formatters run in `make ci` and stay clean.
