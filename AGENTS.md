# Agent instructions

Copal is a project in four parts. Survey, in `survey/`, reads legacy COBOL batch systems and produces design blueprints from them: record layouts, program structure, data flow, business rules and the job network. The analysis is deterministic; no language model is involved until the deterministic core is solid. The First Bank of Pangaea, a small COBOL bank in `pangaea/`, is the specimen the survey is developed and tested against; it knows nothing of the survey. `bedrock/` is the language layer both use, and `console/` the one part built on both. README.md describes the project; docs/brief.md has the plan and the open decisions.

## Rules

- Never mention where this project came from, any company or product it might be compared to as its origin, or any hiring process, in code, documentation, commit messages or anywhere else in the repository. Technical references to well-known tools and systems (GnuCOBOL, ProLeap, IBM mainframes) are fine where they explain a choice.
- Survey is deterministic. Every blueprint fact is derived from the source by parsing and analysis, and is reproducible. Models may come later, only for naming and prose, and only on top of facts the analysis already established.
- Pangaea is fiction written for this project. Never copy real proprietary code, layouts or data into it.
- Pangaea runs: every program compiles with GnuCOBOL and the daily and month-end batches produce outputs that are checked in as golden files. A change to Pangaea updates its golden outputs in the same commit.
- Survey's blueprints for Pangaea are golden files too, so a change in the analyzer shows as a diff in the blueprints.
- The dependencies run one way, and `make ci` checks it: `bedrock` imports nothing else; `survey` and `pangaea` import `bedrock`, never each other; only `console` imports all three. Survey's tests read `pangaea/` as their fixture, which is data, not an import.
- Code is organized by feature, not by layer.
- All text is English with en-US spelling: code, comments, documentation, commit messages, and Pangaea's programs and data. Comments are minimal.
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
- Every analysis step is covered by tests, and by the golden blueprints for Pangaea. Linters and formatters run in `make ci` and stay clean.
