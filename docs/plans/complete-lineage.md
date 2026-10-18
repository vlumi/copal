# Plan: complete lineage

The trail of a field, back or forward, is the survey's most useful view and its least complete. It stops at a `CALL`, so the online day's authorizations are opaque: a hold's amount stops at `AUTHSRV` instead of reaching the rule in `AUTHPUR` that approved it and the terminal's request that asked for it. It stops at `ACCEPT` and starts at `DISPLAY` only as notes. It sees one record type per file. And it ignores the order of statements, so a value every path overwrites still counts. This plan closes those gaps in four parts, one pull request each, in this order.

Read [AGENTS.md](../../AGENTS.md), [architecture.md](../architecture.md), [survey/README.md](../../survey/README.md) and the package comments of `survey/internal/dataflow` and `survey/internal/lineage` first. The ground rules of [design-views.md](design-views.md) apply: facts only, new facts in the model, goldens with the change that causes them, a test for every step, screenshots of every changed page before merging, `make ci` before pushing.

## Part 1: across CALL, and from ACCEPT to DISPLAY

### Within a program

`survey/internal/dataflow` gains four kinds of edge, each with the verb that makes it:

- `CALL 'P' USING a, b` moves each item into the program called: an edge from the item to `P`, `ToCall`, with the argument's position, 1-based, in `Arg`. Passed `BY CONTENT` or by reference alike.
- Each item passed by reference may come back changed: an edge from `P`, `FromCall`, with its `Arg`, to the item. A `CALL` through a data item names no program and makes no edges.
- `ACCEPT x`, with no `FROM`, reads the job's `SYSIN`: an edge from `SYSIN`, `FromStream`, to `x`.
- `DISPLAY x`, with no `UPON`, writes the job's output: an edge from `x` to `SYSOUT`, `ToStream`.

`Trace` treats `FromCall` and `FromStream` as origins, as it treats files: `Origin` gains `Call` and `Arg` for the first, and `Stream` for the second, each with the `Item` at the bytes, as a file origin has. A program called with `USING` also gets an origin for every item of its linkage records: the caller passed it in, `Origin{Parameter: n, Item}`, where `n` is the position of the `USING` item the bytes are in. It is an origin whether or not the program also sets the item, as with every origin while statements are not ordered (part 4).

`Affects` treats `ToCall` and `ToStream` as it treats files: `Effect` gains `Calls`, the bytes passed to each argument of each program called, and `Streams`. An item of a linkage record that the program writes also reaches the caller: `Effect.Parameters`, the bytes of each `USING` item written.

Strings: `call AUTHPUR argument 6 HD-AMOUNT`, `parameter 6 HD-AMOUNT`, `stream SYSIN AUTH-REQUEST`. The field origins of the blueprints show them; the goldens change for the programs that call, are called or ACCEPT.

### Across programs

`survey/internal/lineage` follows them:

- Back from a `Call` origin to the program called, in the same job step: the bytes of the item within the argument passed, at the same offset within the called program's `USING` item at that position, as a node in the called program. The bytes are found with the blueprint's records, by the argument's and the parameter's `Start`.
- Back from a `Parameter` origin to every program that calls this one, each `CALL` of it, in the same job step: the bytes at the same offset within the item the caller passed at that position.
- Back from a `Stream` origin: a `Source` with `Stream: "SYSIN"` and the step's SYSIN, which is outside the system, as an external dataset is, and is shown on the boundary.
- Forward, the mirror: from `Calls` into the called program's parameters, from `Parameters` back to each caller's argument, and to `SYSOUT` as a target, as `Displays` are today, which `Effect.Streams` replaces.

The job step of a called program is its caller's: `Node.Job` and `Node.Step` carry over, as they do through a sort.

Checks against Pangaea:

- Back from `HOLDEXT`'s `TD-AMOUNT` reaches `AUTHSRV`'s `HD-AMOUNT`, then `AUTHPUR`'s and `ATMWDL`'s `HD-AMOUNT`, then their `RQ-AMOUNT`, then `AUTHSRV`'s `RQ-AMOUNT`, then the `SYSIN` of `ONLINE`'s `AUTH` step.
- Back from `AUTHSRV`'s `RS-AVAILABLE`, the response, reaches `AUTHPUR`'s `WS-AVAILABLE`, `AO-BALANCE`, `AO-OVERDRAFT-LIMIT` and `AO-HOLDS`, and through `ACCT-FILE` the `ACCTLOAD` step that loaded them.
- Forward from `AUTHSRV`'s `RQ-AMOUNT` reaches the hold file, `HOLDEXT`, the card file, the master's balance, and `SYSOUT`.
- The trace page shows calls and streams as hops of their own, and highlights the called programs' paragraphs on their flows.

One pull request: the edges; the origins and effects; the field origins in the blueprints, with goldens; lineage across calls; across streams; the trace page; docs.

As built: all checks hold. Two refinements. Within the caller every program called by reference is an origin of what it is passed, since the caller cannot know which callee changes it; across programs, `lineage` keeps only a called program with a statement that moves a value into the bytes, so `ATMBAL` and `AUTHREV` drop out of the hold amount's trail. And a called program no step runs is placed in the steps that run its callers, so a trail can start in `AUTHPUR`. Forward, `DISPLAY` stays a note on the hop rather than a target, as it was; the plan's `Effect.Streams` was not needed. The blueprints' lineage tables put files, programs called and SYSOUT under one heading, "Out of the program", which replaces "Into files".

## Part 2: every record type a program writes

`blueprint.Program.Outputs` keeps the first record type a program writes to a file; a file of header, details and trailer, or a report of headings, lines and totals, has more. Make `Output` one per record written, not per file: the file, the record written, and its fields' origins. The Markdown's field origins get a heading per record under the file's. `lineage.outputs` matches a reader's bytes against every record type the writer wrote to the file, so a reader's field over a header's bytes finds the header's field as well as the detail's.

Checks: `STMTPRT`'s `STMT-OUT` shows its statement lines as well as `WS-SUM-LINE`; `HOLDEXT`'s `TRAN-OUT` has the header, the detail and the trailer, each written from `TRAN-RECORD`, which is the same record redefined, so one entry; `TXVALID`'s reports have their lines. Back from a TXVALID field over a header's bytes finds `HOLDEXT`'s `TH-BUSINESS-DATE`.

One pull request; many goldens change.

As built: the checks hold, with one correction to the plan. `HOLDEXT`'s header and `TXVALID`'s back to it already worked, since `TRAN-RECORD` is one record redefined and its leaves include every record type's fields. The change shows where a file is written `FROM` several working-storage records: every program's report, and `STMTPRT`'s statements, which also writes `STMT-LINE` itself for the blank lines. No dataset in Pangaea that another step reads is written that way, so a unit test covers the reader's side. A record written bare and `FROM` others is listed for both.

## Part 3: order-aware origins

The [brief](../brief.md)'s debt: origins ignore the order of statements. A value that every path to the item's use overwrites is not an origin.

Follow values back along the program's control flow rather than over all its statements. `survey/internal/structure` knows the paragraphs and how control moves between them, and `bedrock/cobol` the statements in order. A reaching-definitions analysis over statements: for the use of an item at a statement, the definitions of its bytes that can reach it along some path with no definition of the same bytes in between. Start, for a field written to a file, at the `WRITE`; for a field displayed, at the `DISPLAY`; and for a field asked about without a statement, at every use. A definition that covers the bytes only in part, such as a field of a group, kills only those bytes.

This is the largest part. Do it in steps, each its own commit: the control-flow graph of statements, with `PERFORM` ranges, `GO TO`, fall-through and the branches of `IF`, `EVALUATE` and the conditional phrases; reaching definitions over it, by bytes; `Trace` and `Affects` on top, keeping the order-free analysis as the fallback for what the graph does not cover, say a `GO TO DEPENDING ON` with a computed target; then the blueprints.

Checks: before implementing, read the programs and list the origins in today's goldens that no path can reach, with the paths that kill them; the implementation must remove exactly those. Two cases to keep apart: a value overwritten on every path, which this part removes, and a field of a record redefined for another record type, such as `MONCHG`'s trailer fields over the description's bytes, which order does not explain and which stays until the survey knows which record type a path writes; say so in survey/README.md. `TXPOST`'s `AM-BALANCE` written by `4000-WRITE-MASTER` must keep every posting as an origin. Record in the plan, as built, every origin that disappeared from the goldens and why.

One pull request. Remove the debt from the brief.

As built: the checks were made against an audit done separately, before the implementation's goldens were compared, of every origin in the goldens that no path can reach. The implementation removes exactly the audit's list, except where the two rules differ on purpose, and `TXPOST`'s `AM-BALANCE` keeps every posting.

- `survey/internal/controlflow` lays out the paths through the statements, a performed range anew for each `PERFORM`, so that control returns only to the statement after the one that entered it. Every Pangaea program is covered; a `GO TO` out of a performed range would leave a program to the order-free analysis.
- A statement overwrites the items it writes by plain name and does not also read: `MOVE`, `COMPUTE`, the arithmetic verbs into an item they do not update, `SET`, `ACCEPT`, `READ` and `RETURN` into their records, and a `WRITE FROM` into the file's record. A bare `WRITE` overwrites nothing.
- Field origins start at the `WRITE`s of each record; `Trace` without a statement starts at every use of the item, and `Affects` at every statement that sets it and at the entry.

What disappeared from the goldens, and why:

- Every program's report line, written bare, no longer collects the count and amount lines written `FROM` other records, nor what those carry: the counters' initial values, their literals, the totals' inputs. A `MOVE SPACES` or a heading literal overwrites the line before each bare `WRITE`.
- `STMTPRT`'s `STMT-LINE`, written bare only for blank lines, keeps only `SPACES`. Its sort record's fields lose the values the output procedure's `RETURN` reads back, which come after the input procedure's `RELEASE`.
- `HOLDEXT`'s header, detail and trailer fields lose the feedback through `WS-ACCT-HASH` and `WS-AMOUNT-TOTAL`: the detail's amount is always the hold's amount or fee, set after the `MOVE SPACES`, and the trailer is written after the last detail. The report's amount line loses the same.
- `TXPOST`'s transfer legs lose the input's account numbers and code, which `3500-WRITE-LEG` swaps and sets on every path; its audit and reject records lose the leg's values, which the restore from `WS-SAVED-TRAN` overwrites first.
- `AUTHSRV`'s last authorization code loses `RQ-AUTH-CODE` and the held numbers passed to the programs called: the `READ` of the hold file overwrites the record first.
- `GLSUM`'s ledger account and name lose `SPACES`, overwritten from the ledger table; `INTACCR`'s through date loses the accrual read in, overwritten by the business date.

Where the audit and the implementation differ: a `READ` with `INVALID KEY` or a `RETURN` with `AT END` counts as leaving its record as it was, which it does on that path, so `AUTHSRV`'s keys moved in before a keyed read, and `STMTPRT`'s audit values in memory, stay. `MONCHG`'s transaction id gains `WS-TI-SEQ`'s initial value, which the audit found missing: `ADD 1` keeps the old value, and the bytes left once a statement overwrites part of a group are followed on their own. For the same reason a few file origins of `TXPOST`'s leg trailer are named by a narrower field over the same bytes. Values that only the outcome of a condition keeps from the write stay, as do a redefined record type's fields over another type's bytes; survey/README.md says so.

## Part 4: Pangaea's honest mess

The bank is tidier than a real one. Give the survey real inconsistencies to find, each a change to Pangaea with its goldens:

- `MAINTREC` keeps a six-digit effective date, `YYMMDD`, never widened in the Y2K project; `ACCTMNT` windows it on read, at 50, with the history entry `1999-01-11 SL Y2K: LEFT AS IS, WINDOWED ON READ`. Generated requests fill it.
- `REJREC` carries the date it was rejected as text, `PIC X(10)`, `YYYY-MM-DD`, moved from the business date with `STRING`.
- The survey then has one business date in three layouts. Add to the entity pages a section on fields of the same kind with different layouts across the entities, by the kinds and the lineage: a field whose values come from a date field of another layout is listed with both.

Checks: the trail of `REJREC`'s date reaches `BANK.CONTROL`'s `BC-BUSINESS-DATE` through the `STRING`; the entity pages list the three layouts.

One pull request.

As built: the checks hold. `MR-EFFECTIVE-DATE` came with CR-0129 in 1991, which dates an opening from it; generated and fixture requests carry the business date, so the master is unchanged and only the rejects' golden outputs change. The rejection date did not fit the reject record's FILLER, so CR-0444 widened the record to 110 bytes and appended `RJ-REJECT-DATE` at the end. Each entity field has its layout, and a date field or a text field that is not a redefinition lists the date fields of another layout its values come from; amounts computed from dates and a body redefined for several record types are left out. The entity pages show four layouts, since the online requests' and holds' timestamps are fourteen digits, and find two more dates from another layout besides the plan's: `ACCT-MASTER-RECORD`'s open date from the request's, and `ACCRUAL-RECORD`'s due month from the business date. A six-digit picture holds both `YYMMDD` and `YYYYMM`; the survey cannot tell them apart.

## Next: the oracle

Planned on its own, in [oracle.md](oracle.md).
