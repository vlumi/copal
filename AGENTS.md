# Agent instructions

Copal reads legacy COBOL batch systems and produces design blueprints from them: record layouts, program structure, data flow, business rules and the job network. The analysis is deterministic; no language model is involved until the deterministic core is solid. The First Bank of Pangaea, a small COBOL batch bank in `pangaea/`, is the specimen it is developed and tested against. README.md describes the project; docs/brief.md has the plan and the open decisions.

## Rules

- Never mention where this project came from, any company or product it might be compared to as its origin, or any hiring process, in code, documentation, commit messages or anywhere else in the repository. Technical references to well-known tools and systems (GnuCOBOL, ProLeap, IBM mainframes) are fine where they explain a choice.
- Copal is deterministic. Every blueprint fact is derived from the source by parsing and analysis, and is reproducible. Models may come later, only for naming and prose, and only on top of facts the analysis already established.
- Pangaea is fiction written for this project. Never copy real proprietary code, layouts or data into it.
- Pangaea runs: every program compiles with GnuCOBOL and the daily and month-end batches produce outputs that are checked in as golden files. A change to Pangaea updates its golden outputs in the same commit.
- Copal's blueprints for Pangaea are golden files too, so a change in the analyzer shows as a diff in the blueprints.
- Code is organized by feature, not by layer.
- Comments: minimal, English, en-US spelling.
- Markdown paragraphs are one line each; the editor wraps. `.markdownlint.jsonc` encodes this.
- Diagrams are generated as SVG or PNG, not Mermaid: the editor in use does not render Mermaid.
- Routine operations are make targets with a `##` description, so plain `make` lists them. Run `make ci` before every push to `main` and only push when it passes; a push that changes only Markdown needs none.
- Commit and push to `main` directly when something works; frequent small commits, no pull requests.
