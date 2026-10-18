# Plan: the oracle

A rewrite has to do what the old system does, rule for rule. The survey already knows the rules, as decision tables; what it lacks is evidence of what the programs do when each rule applies. The oracle is that evidence: for each program, input data that reaches every rule of every table, run through GnuCOBOL, with the outputs, reports and return codes recorded as golden cases a rewrite can be tested against. Program by program first, then whole batches.

Read [AGENTS.md](../../AGENTS.md), [architecture.md](../architecture.md), [survey/README.md](../../survey/README.md), [bedrock/README.md](../../bedrock/README.md) and the package comments of `survey/internal/rules`, `pangaea/runner` and `pangaea/compile` first. The ground rules of [design-views.md](design-views.md) apply: facts only, new facts in the model, goldens with the change that causes them, a test for every step, `make ci` before pushing.

## Decisions

- **Running COBOL moves into bedrock.** Compiling with `cobc` and running a program with its files bound are what Pangaea's machine room and the survey's oracle both need, so they belong in bedrock, where survey and pangaea may both import them. Pangaea's job runner stays in pangaea, built on top. The oracle stays a survey feature that works on any system the survey reads, not only Pangaea.
- **Coverage-guided search, not symbolic solving.** Many conditions test working-storage values computed from the inputs: a new balance against the overdraft floor. Propagating inputs symbolically through COBOL arithmetic, redefinitions and files is a large project before the first result. Instead, candidate values come from the decision tables: each condition's literals and condition-name values, with their boundaries, put into the input fields the lineage says the condition depends on. Each candidate is run with tracing, and the records that reach a rule no earlier record reached are kept. The search is deterministic: a fixed order, a fixed seed, and the same cases every time. The rules it cannot reach are reported, with the conditions that block them.
- **Programs first, then the batch.** A case is one main program run alone, with synthesized files bound to its DDs. Called programs are reached through their callers. Batch cases, built from the program cases and run through the jobs, come last.

## How coverage is measured

GnuCOBOL 3.2 traces every paragraph and statement it executes, with its source line, when a program is compiled with `-ftraceall` and run with `COB_SET_TRACE=Y` and `COB_TRACE_FILE`. A rule of a decision table is a path through a paragraph: the outcome of each condition it tests and the statements it runs. A pass through the paragraph in the trace reached the rule when the statements it executed, between entering the paragraph and leaving it, are the rule's. So the rules gain the source positions of the statements on each path, and coverage is the replay of a trace against them.

## Part 1: bedrock runs COBOL

`bedrock/run`:

- `Compile(sources, copybookDirs, outDir, flags)`: what `pangaea/compile` does today, an executable for a main program and a module for a called one, plus the trace flag when asked.
- `Program(binary, Binding)`: run one program with its DD names bound to files through the environment, its in-stream data on stdin, `COB_LIBRARY_PATH` for its modules, and optionally a trace file. Returns the return code, whether it abended, and stdout.

`pangaea/compile` and the step runner in `pangaea/runner` are rebuilt on it; the job runner's catalog, generations, COND tests and logs stay in pangaea. The dependency check still passes, and Pangaea's golden outputs do not change.

One pull request: the package with tests on a small fixture program; pangaea on top of it; docs.

## Part 2: rule coverage from a trace

- `survey/internal/rules`: each rule gains the positions of the statements on its path, in order.
- `survey/internal/coverage`, new: read a GnuCOBOL trace into passes through paragraphs, by program and source line, and replay each pass against the program's tables. The result is, for each table, the rules reached and how many passes reached each.
- Tests compile fixture programs with tracing and run them: CI has GnuCOBOL already. Also check the replay on Pangaea's daily fixture: run it traced and report the coverage of today's goldens. That is the baseline, and the first number worth a look.

One pull request.

## Part 3: synthesizing records

`survey/internal/oracle`, the search, for one main program:

- **Inputs:** the files it reads, bound by the job step that runs it, with their record layouts; its SYSIN, if it accepts; and for a program it calls, the caller's arguments come through the caller.
- **The base case:** one valid record of each input file, filled field by field by kind: keys and numbers increasing, dates from a fixed business date, amounts small and positive, codes at their first condition name, text by name. The bank control record carries the same business date.
- **Candidates:** for each condition of each table, the input fields it depends on, from the lineage back to the files. Values come from the condition's literals and the values of its condition names, each one, one above and one below for numbers, plus spaces, zero and the field's largest value. Each candidate is the base case with one field changed, then two for conditions on two fields.
- **Records that belong together:** some rules hold between records, and format-only generation breaks them. The survey's facts give them:
  - **Record order:** a field compared with condition names that tell header, detail and trailer apart gives the file's order: header first, details, then trailer.
  - **Control totals:** a trailer field compared with an item that accumulates a detail field, by ADD from the file, becomes that sum, or that count for an item that adds 1.
  - **References:** a field the program matches against a key of another input file draws its values from that file's keys, and a candidate breaks the reference on purpose to reach the rule for an unknown key.
- **The loop:** run each candidate traced, replay the trace, and keep the candidates that reach a new rule. Stop when a round adds nothing. Merge the kept candidates into as few cases as keep the same coverage: one file of many details reaches more rules in one run than many files of one.
- **Unreached rules:** listed per table, with the condition outcomes the search never produced.

Indexed input files, such as the online account and hold files, need loading. Load them through `bedrock/run` with a small generated COBOL loader, so that the oracle needs nothing outside GnuCOBOL.

One pull request, possibly two: the base case and candidates first, then the rules between records.

## Part 4: the program cases

- For each main program of the system, its cases under `blueprints/<system>/oracle/PROGRAM/NN/`: the input files as text, as Pangaea's fixtures are, the outputs, the reports, the return code, and the rules each case reaches.
- `oracle.md` per program: the coverage, rule by rule, linked to the decision tables in the program's blueprint, and the unreached rules with why.
- `survey oracle -copybooks dir -out dir ...` writes them, and `make oracle` regenerates Pangaea's. They are golden: CI runs the cases and compares the outputs, which also catches a change in Pangaea's behavior, and regenerates the search to compare the cases themselves, if the time allows; otherwise only on demand.
- The explorer's program page gets the coverage, and each decision table shows which rules a case reaches, with a link to the case.

Checks against Pangaea: every rule of `TXVALID`'s header, detail and trailer checks is reached, including each trailer total that differs; `TXPOST` reaches the overdraft floor exactly at the limit and one cent past it, and a fee past the limit; `ACCTMNT` reaches every rejection reason; `AUTHSRV` reaches `AUTHPUR`'s and `ATMWDL`'s approvals and declines through SYSIN requests.

One pull request.

## Part 5: the batch

Cases for the daily and month-end jobs: the program cases' inputs for the jobs' first steps, the rest produced by the steps themselves, run through the job runner, with the datasets each step writes recorded. Running a job needs the job runner from Pangaea's machine room; whether that also moves into bedrock, or the batch cases are run from the console, which imports both, is decided when this part starts, once parts 1 to 4 show how much of the runner is generic.

One pull request.

## Later

A model for naming paragraphs and explaining rules in prose, on top of the facts and now the evidence, stays last, as the brief says.
