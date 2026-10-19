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
2. **The survey on it, as it is.** Run it, and write down what breaks and what comes out wrong, program by program: the list of couplings. Nothing is fixed in this part.
3. **The language.** Free format and lower case, sections, `DECLARATIVES`, nested programs, `RETURNING`, level 78, the new usages and clauses, several 01 records under an FD as record types, variable-length and relative files. Each fixed in the parser, the layouts and the analyses, with Pangaea's blueprints unchanged where nothing about Pangaea changed.
4. **The batch from scripts.** A reader for the shell subset the scripts use, in bedrock, beside the JCL reader, and the survey's jobs, datasets, schedule and lineage from either. A step that is not COBOL, such as `sort`, is a step that passes records on, as a sort does in JCL.
5. **The rest of the survey.** Lineage, the design views, the oracle and its batch run on Laurasia, with its own golden cases; and the remaining couplings from part 2.

Each part one pull request or more; part 3 most likely several, by feature.
