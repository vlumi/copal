# Architecture

Copal is one Go module in four parts, each a directory at the top of the repository, and the files where two of them meet. What reads COBOL is a library of its own, [cobolith](https://github.com/vlumi/cobolith), which every part imports.

![The parts of Copal: the console imports the survey, the bank, bedrock and cobolith; the survey, the bank and Laurasia import bedrock and cobolith, and bedrock imports cobolith; the survey writes the blueprints and reads the bank](diagrams/architecture.svg)

## The parts

- **Survey**, `survey/`, is the analyzer: it reads a COBOL system and writes its blueprints, and serves the explorer. Only `blueprint`, its model of a system, is exported; every analysis is under `survey/internal/`. [survey/README.md](../survey/README.md) describes the pipeline.
- **Pangaea**, `pangaea/`, is the specimen: a legacy COBOL bank with the machine room that runs it in Go. It is written as a legacy system would be, documented as one in [pangaea/docs/](../pangaea/docs/README.md), and it knows nothing of the survey.
- **Laurasia**, `laurasia/`, is the second specimen: a card processor in another shop's COBOL, Micro Focus dialect and free format, run by shell scripts instead of JCL, with no machine room beyond its test. It is there to show where the survey assumes Pangaea. Its operations guide is in [laurasia/docs/](../laurasia/docs/operations.md).
- **Cobolith**, [cobolith](https://github.com/vlumi/cobolith), is the language layer everything stands on: reading COBOL source and copybooks in the IBM and Micro Focus dialects, laying out records, reading and writing them by field, and reading JCL. It was written here and is now a library in a repository of its own, which Copal requires by version; how it got there is in [plans/library.md](plans/library.md).
- **Bedrock**, `bedrock/`, is the harness built on it: keeping datasets and their catalog, reading batch scripts, and compiling COBOL programs with GnuCOBOL and running them and the jobs that run them. The survey's oracle runs programs with it, the bank's machine room runs the bank with it. None of it is analysis.
- **The console**, `console/`, is the one part built for the modernization rather than inherited from the past: it runs the bank, watches its online day, and shows each run's jobs as the survey draws them, with links into the blueprints. It is the only part that knows both the bank and the survey.

## The rule

Dependencies run one way, and `make lint-deps` checks it in CI with `go list`, tests included: `bedrock` imports nothing of Copal, and of cobolith, the library; `survey`, `pangaea` and `laurasia` import `bedrock` and `cobolith` and never each other; only `console` imports the rest. The survey's tests read `pangaea/` as their fixture and compare with `blueprints/pangaea/`, which is data, not an import.

The rule exists because a bank must not depend on the tool that studies it. Everything the bank's Go needed from the analyzer turned out to be either language (now cobolith and bedrock) or presentation (now the console), and once those were moved the import went away.

## Golden files, twice

Two things are checked in that a build regenerates and compares:

- **Pangaea's outputs.** Each fixture under `pangaea/testdata/` runs its jobs through GnuCOBOL, and the datasets, job logs and reports they produce are golden. A change to the bank updates them in the same commit, so the history of the bank's behavior is in the history of its goldens.
- **Laurasia's outputs.** Each fixture under `laurasia/testdata/` runs its scripts on its input files, and what they leave behind is golden, the same way.
- **The survey of Pangaea.** `blueprints/pangaea/` is what the survey says about the bank. A change in the analyzer shows as a diff there, which is how its behavior is reviewed: not by reading the code change alone, but by reading what it changed in the blueprints of a real system.

Both depend on everything being deterministic: the analysis, the diagram layout, and the bank itself, whose days follow from a seed and never read the clock.

## One module

Two modules with a Go workspace would make the sibling relation formal, but add toolchain friction for nothing a reader gains, and the survey's tests reading `pangaea/` across a module boundary would be awkward. One module, four trees, one rule.

## Commands and make

Each part has its own command: `survey/cmd/survey`, `pangaea/cmd/pangaea`, `console/cmd/console`. Routine operations are make targets with a description, so plain `make` lists them; the Makefile names the work and scripts and commands do it. `make ci` is what the hosted build runs: the toolchain check, Markdown and Go linting, the dependency check, and every test, the goldens included.

The diagrams on this page are drawn by `make docs-diagrams` with the survey's own diagram package, and a test keeps them current.
