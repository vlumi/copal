# Plan: a second specimen

Every heuristic in the survey was tuned against one codebase, written by one hand, in one style. Before a model is let near the blueprints, the facts should hold on a system the survey was not grown on. Laurasia Card Services is that system: a card processor, issuer and acquirer in one, written in another shop's COBOL, run by shell scripts instead of JCL. Pangaea split into Laurasia and Gondwana; the names stay in the family.

The deliverable is not the specimen but what it shows: every place the survey breaks on it, or comes out wrong, is a coupling to Pangaea, and is fixed generically, never by special-casing Laurasia.

## Ground rules

- Laurasia is fiction, like Pangaea: no real proprietary code, layouts, network formats or data. It looks like card processing, not like any network's actual specification.
- It runs: every program compiles with GnuCOBOL in its own dialect, its batches produce outputs that are checked in as golden files, and a change to it updates them in the same commit.
- It knows nothing of the survey, as Pangaea does not. It imports bedrock at most, and its batch runs from shell scripts, not from Go.
- The survey's blueprints of it are golden files too, from the start: the first ones are wrong, and the history shows them getting right.

## What it varies, on purpose

Each difference breaks an assumption the survey makes today.

| Pangaea | Laurasia |
| --- | --- |
| IBM dialect, `-std=ibm`, fixed format, upper case | Micro Focus dialect, `-std=mf`, free format, lower case, `*>` comments |
| Paragraphs, `PERFORM` of paragraphs and `THRU` ranges | Sections as the unit of work, paragraphs within them, `PERFORM` of sections |
| No `DECLARATIVES` | `DECLARATIVES` with `USE AFTER ERROR` for file errors |
| Record types as `REDEFINES` of one record, a type field | Several 01 records under one FD, and variable-length records, `RECORD VARYING ... DEPENDING ON` |
| Sequential and indexed files | Relative files too, and line sequential reports |
| Display and `COMP-3` numbers | `COMP-5` and `BINARY`, `SIGN SEPARATE`, level-78 constants |
| Static calls with `USING` | Calls with `RETURNING`, a call through a data item, a nested program |
| JCL: steps, DD statements, generations, `COND` | Shell scripts: `cobcrun` with `export DD_...` bindings, return codes checked with `if`, Unix `sort` as a step, files named by date instead of generations |
| Abort paragraphs named `9900-ABORT`, switches named `...-SW` | Other conventions: `stop run returning`, flags with `88` levels on `pic x` named otherwise |

## The domain

A card processor sees both sides of a card payment:

- **Authorization**, as a log the online side leaves behind: approvals, declines, stand-in decisions.
- **Clearing:** the networks' files of presentments, in variable-length records of several layouts, a file header and trailer, batch headers and trailers with counts and amounts.
- **Matching** presentments to authorizations, and the ones that match nothing.
- **Interchange:** a fee for each presentment from a table by card product, merchant category and region, the widest decision tables of the two specimens.
- **Posting** to the cardholders' accounts, with credit limits and over-limit fees.
- **The billing cycle:** interest on the average daily balance, a grace period, a minimum payment, statements.
- **Disputes:** a chargeback's life, from first chargeback through representment and pre-arbitration, a state machine the survey should find.
- **Merchant settlement:** what each merchant is paid, net of fees and chargebacks.

## Parts

1. **The first slice of Laurasia.** About seven programs: clearing ingest, matching, interchange, posting, the billing cycle, disputes and merchant settlement; their copybooks; three scripts, the nightly run, the cycle and the disputes; a generator for its inputs and a fixture with golden outputs; its own README and operations guide. The repo's rules and docs stop saying Pangaea where they mean a specimen. Bedrock's build takes the dialect and the source format as options, as `-std=ibm` is hard-wired today.

   As built: fourteen programs rather than seven, as interchange needs its fee table loaded, posting its account file loaded and unloaded, and posting, the cycle and the fee table each call a module; four scripts, with `open.sh` to load the accounts, and the night calling the disputes' script. The fixture, one night and its cycle, was made by hand rather than by a generator, which is left for when the survey needs more data than one night. GnuCOBOL does not implement `RETURNING` on the called program's side, so the modules set `RETURN-CODE` and the callers use `CALL ... RETURNING`. The source format is the `>>SOURCE FORMAT IS FREE` directive in each file, which bedrock's lexer reads, so the build needs only the dialect: `run.IBM` or `run.MicroFocus`. Bedrock also reads `COPY` of a literal name with its extension. The oracle and the console still compile in `run.IBM`; the survey has not yet read Laurasia.
2. **The survey on it, as it is.** Run it, and write down what breaks and what comes out wrong, program by program: the list of couplings. Nothing is fixed in this part.

   As built: [laurasia-findings.md](laurasia-findings.md). The survey reads none of the programs as they are; with the parser patched in a throwaway copy, 21 findings, 9 where the parser stops and 12 where the facts are wrong or missing, each marked with the part that fixes it.
3. **The language.** Free format and lower case, sections, `DECLARATIVES`, nested programs, `RETURNING`, level 78, the new usages and clauses, several 01 records under an FD as record types, variable-length and relative files. Each fixed in the parser, the layouts and the analyses, with Pangaea's blueprints unchanged where nothing about Pangaea changed.

   As built, in #94 to #98: findings 1 to 16 of [laurasia-findings.md](laurasia-findings.md) are fixed, and Laurasia's blueprints are golden since every program reads. The dialect turned out to be the language's, not only the compiler's: Micro Focus lays binary items out in as few bytes as their digits need, so `cobol.Dialect` carries it to the layouts, and `survey -dialect` reads a system in it. A count's record type was guessed from any sum of the same file, which Laurasia showed wrong; now only from a sum the same paragraph adds. Left for later: a level-78 constant's value is shown but not used as a value by the rules; lineage stops at what a call returns rather than following it into the program called; `record` reads overpunched signs EBCDIC style only, which Laurasia does not use; and the oracle has no shape for a file whose types are records of their own, which part 5 needs.
4. **The batch from scripts.** A reader for the shell subset the scripts use, in bedrock, beside the JCL reader, and the survey's jobs, datasets, schedule and lineage from either. A step that is not COBOL, such as `sort`, is a step that passes records on, as a sort does in JCL.

   As built: `bedrock/script` reads the subset Laurasia's scripts use, and a script is a job in the survey. Its datasets are the files its `DD_` variables and redirections name, with the variables that only stand for a parameter kept as names, so `$DATA/work/$DATE/pres.txt` is one dataset in every script that names it. A sort or `mv` passes records on, in the copybooks of the datasets either side and in lineage; a step may run another script; the test after a step says when the script stops. Scripts name files, not generations, so a job that creates a file runs before the jobs that read it. Finding 19 is fixed; finding 20, the run date on the command line, is the oracle's, for part 5. Reading whole records into others turned out to blur lineage in Pangaea too: an origin now names the bytes of the record it comes from.
5. **The rest of the survey.** Lineage, the design views, the oracle and its batch run on Laurasia, with its own golden cases; and the remaining couplings from part 2.

   As built: findings 17, 18, 20 and 21 are fixed, and with them every finding of part 2. Merge links come from keys compared for order where one takes the highest value at its file's end, and state machines read a conjunction's parts. The oracle runs Laurasia's programs in their dialect and reaches every rule a run can take but those of I/O errors and of combinations the code rules out: through traces of lower-case source and sections, the command line as an input stream, files of several records laid out as one record area, relative files loaded by record number, records read whole laid out as the record read into, constants and digits as values to try, and arguments that are fields. Its cases are golden in `oracle/laurasia`, and `make explore-laurasia` serves the explorer. Left for later: the oracle runs Laurasia's programs one by one, not its scripts as a batch, as it runs Pangaea's jobs; Laurasia's own test runs the night. Several fixes changed Pangaea's facts for the better, where Laurasia showed a coupling Pangaea hid: byte-exact origins of records read whole, ACCTMNT's merge link.

Each part one pull request or more; part 3 most likely several, by feature.
