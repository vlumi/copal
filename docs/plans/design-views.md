# Plan: design views in the explorer

The explorer is organized the way the code is: by program, by job, by field. Designing a replacement needs the same facts organized the way the business is: by entity, by rule, by boundary, by who owns what, and by what a change touches. This plan adds those views, and makes the kinds of things visibly different throughout. It is written for whoever implements it, agent or person, to work through without the conversation that produced it.

Read [AGENTS.md](../../AGENTS.md) first, then [architecture.md](../architecture.md), [blueprints.md](../blueprints.md) and [survey/README.md](../../survey/README.md). Every rule there applies: deterministic facts only, feature packages, granular commits, goldens updated in the same commit, `make ci` before pushing, one branch and pull request per part below, merged with `gh pr merge --merge --delete-branch` once local `make ci` and the hosted check pass.

## Ground rules for this work

- A view shows facts the analysis established. Where a classification needs a heuristic, such as reading a field as a date from its name, the rule is stated in the code's comment and in [blueprints.md](../blueprints.md), it is applied the same way everywhere, and the blueprint JSON carries the result so that a change in the rule shows as a diff in `blueprints/pangaea/`.
- New facts go into the blueprint model in `survey/blueprint`, not into the explorer, so that the Markdown blueprints, the JSON and the explorer say the same thing. The explorer only arranges and links.
- Every new page gets a route in `survey/internal/explore/server.go`, a template, links to it from the pages that mention the things it shows, and a check in `survey/internal/explore/explore_test.go`. The crawl test there follows every link from the overview, so a new page is tested for being there as soon as something links to it; add content checks for what it must show.
- Diagrams are drawn with `survey/internal/diagram` and must stay deterministic. A new node kind gets a legend key and a look of its own, not only a different corner radius.
- Pangaea is the test. Each view is checked against what we know of the bank: the examples given under each part are the acceptance checks.
- Record what a view does not attempt in its page's own note and in [survey/README.md](../../survey/README.md) under what the survey does not do.

Useful mechanics, learned the hard way: `goimports` is at `~/go/bin/goimports` and wants `-local github.com/vlumi/copal`; `rsvg-convert` previews SVG but does not render text on paths, headless Chrome does (`"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" --headless=new --screenshot=out.png URL`); `make blueprints` regenerates the goldens after a model change and `make pangaea-golden` the bank's; chain shell commands with `&&`.

## Part 1: kinds, visible everywhere

What it is for: at a glance, a reader should tell an amount from a code from a key, a called program from a main one, an external input from a scratch file, without reading the pictures and the JCL.

### Field kinds

Add `Kind string` to `blueprint.Field`, JSON `kind`, computed in `newField` from the layout, by these rules in this order, first match wins:

1. `group`: the field has fields.
2. `filler`: the name is `FILLER`.
3. `key`: the field is a file's `RECORD KEY`, or a field of a SORT's work record named as a sort key. Record keys are in `File.RecordKey`; sort keys need `fileio.Sort` to carry them, so add `Keys []string` to it from the SORT statement's `Reads`.
4. `code`: the field has level-88 condition names.
5. `edited`: the picture is numeric-edited, that is, contains any of `Z`, `*`, `,`, `.`, `-`, `+`, `$`, `CR` or `DB` outside `V` and `S`; `layout` already distinguishes the category, use it.
6. `amount`: numeric with a scale above zero.
7. `date`: numeric without scale, of 6, 8 or 14 digits, whose name contains `DATE`, `MONTH`, `TIMESTAMP` or `YYMM`. A name heuristic, say so in the comment.
8. `number`: any other numeric field.
9. `text`: anything else.

Show it: a `Kind` column in the layout tables of the Markdown blueprints (programs and copybooks), the `kind` in the JSON, and in the explorer a small badge next to the field name on the field page and in layout tables, with a CSS class per kind, `kind-amount` and so on, each with its own muted color and the same color in light and dark modes. Keep a legend of the badges on the copybook page, once, at the top of the layouts.

Checks against Pangaea: `AM-BALANCE` and `TD-AMOUNT` are amounts; `AM-STATUS`, `TD-TRAN-CODE` and `RS-RESULT` are codes; `AO-ACCT-NO` is a key as a record key and `SR-ACCT-NO` as a sort key, while `AM-ACCT-NO` comes out as a `number`, since the master is read sequentially and nothing names it a key. State that limit in the page's note rather than guess from the name: a rule by name is not exact, and the view must not say more than the code does. `BC-BUSINESS-DATE` and `AM-OPEN-DATE` are dates; `RS-BALANCE` is edited; `CR-DAY-PURCHASES` is an amount.

### Program kinds

Add `Kind string` to `blueprint.Program`, JSON `kind`: `called` if `Using` is not empty, otherwise `main`. In the explorer, the index and the search show a badge, and the program page says "called program" or "main program" in its crumbs. The job network already shows which steps run which program; no change there.

### Datasets of the system, with kinds

Add `Datasets []SystemDataset` to `blueprint.System`, built after the jobs are analyzed, one per dataset name without generation, temporary ones included with their `&&`:

```go
type SystemDataset struct {
    Name      string    `json:"dataset"`
    Copybook  string    `json:"copybook,omitempty"`
    Indexed   bool      `json:"indexed,omitempty"`
    Temporary bool      `json:"temporary,omitempty"`
    Kind      string    `json:"kind"`
    WrittenBy []StepRef `json:"writtenBy,omitempty"`
    ReadBy    []StepRef `json:"readBy,omitempty"`
}

type StepRef struct {
    Job  string `json:"job"`
    Step string `json:"step"`
}
```

`Kind` is `external` if no step writes it, `output` if no step reads it, `temporary` for `&&` names, otherwise `internal`. A step writes a dataset when its DD's uses include write, append or update, as `DD.writes` says, and reads it when they include read or update. Add the datasets table to the blueprints' index README, with the kinds, and replace the explorer's dataset index with it. Badges for `external`, `output`, `temporary`, `indexed` on the dataset pages and in the job network's labels, where the node for an external dataset gets a look of its own in the diagram, with a legend key, so that the system's inputs stand out on the DAILY network.

Checks: `TRAN.DAILY`, `ACCT.MAINT`, `BANK.RATES` and `CARD.MASTER` are external; `GL.SUMMARY`, `MAINT.REJECTS` and `TRAN.REJECTS` are outputs; `ACCT.ONLINE` is internal and indexed; `&&SORTED` is temporary.

### Records

`Record.Section` already says `file`, `working-storage` or `linkage`. Show it as a badge on the field page's crumbs and in the program page's records table; no new fact.

One pull request, commits: sort keys in `fileio`; field kinds in the model and the blueprints; program kinds; system datasets; the explorer's badges and the diagram's external datasets; docs.

## Part 2: rules by field, and entities with their state machines

### Conditions know their fields

`rules.Table.Conditions` are strings. Add `Refs [][]string`, parallel to `Conditions`, the data names each condition refers to, from `cobol.Condition.Refs()`; carry it into `blueprint.DecisionTable` as `conditionRefs`. In `blueprint`, resolve a level-88 name to its field with the records' `Conditions`, the way `dataflow` does with its `names`, so that a condition `AM-CLOSED` counts as a condition on `AM-STATUS`.

### Rules by field

On the explorer's field page, a section **Rules**: every decision table in the program with a condition on the field, with the conditions' text and a link to the paragraph. Then a new page, `/entities/{record}/fields/{field}`, for a field of a copybook record across every program that uses the copybook: the programs, and in each the paragraphs that read it, write it, and decide on it, with the conditions. The program-level field page links to it when the record comes from a copybook. The Markdown rules blueprints stay as they are; this view is the explorer's.

Checks: `AM-STATUS` across the system shows `TXPOST` deciding in `3000-APPLY-TRAN`, `ACCTMNT` in its maintenance paragraphs, `ACCTLOAD` in `2000-LOAD-ACCOUNT`; `AM-OVERDRAFT-LIMIT` shows `TXPOST`, `ACCTMNT`, `AUTHPUR`, `ATMWDL` and `ATMBAL`.

As built: the `AM-STATUS` check holds, with `MONCHG` and `STMTPRT` deciding on it too. `AM-OVERDRAFT-LIMIT` shows `TXPOST`, `ACCTMNT` and `ACCTLOAD`; the authorization programs use `AO-OVERDRAFT-LIMIT` of the online account, a record of its own copybook, `ACCTONL`, whose page shows `AUTHPUR`, `ATMWDL`, `ATMBAL` and `AUTHREV`. The two are joined by `ACCTLOAD`'s move, which is lineage, not this view.

### Entities

A new page per copybook record, `/entities/{record}`, and an index at `/entities`. An entity is a record a copybook describes. The page shows: the datasets described by it, with their kinds; the programs that use it, and for each field, who reads and who writes it, as compact lists; the fields with their kinds; and for each `code` field, its **state machine**.

The state machine is derived from the decision tables and the lineage, deterministically:

- The states are the field's level-88 values, by condition name. A field whose 88s have several values each is still drawn by condition name.
- A transition from state A to state B exists where a decision table has a rule in which a condition on the field holds for A (the condition is A's 88 name with entry `Y`, or `NOT` A with `N`, or `FIELD = 'literal'` matching A's value) and one of the rule's actions sets B: `SET B TO TRUE`, or `MOVE 'literal' TO FIELD` with B's value. The label is the program and paragraph.
- A rule that sets B without testing the field is an entry from anywhere: draw it from a node named `any`.
- Values moved from another field are not transitions and are not drawn; the page says so.

Draw it with `diagram`: states as nodes and transitions as edges. Edges need a label for the program and paragraph, which `diagram.Edge` does not have, so give it an optional `Label string`, drawn at the edge's midpoint, deterministic like the rest, with a test. Put the state machine on the entity page and, for the field itself, on its cross-program page.

Checks: `AM-STATUS` shows `O` to `F` and `F` to `O` from `ACCTMNT`'s freeze and unfreeze, `O` to `X` from close, and an entry to `O` from open; `HD-STATUS` shows `O` to `R` from `AUTHREV`; `CR-STATUS` has no transitions, which is itself a finding worth stating: nothing in the system blocks or unblocks a card.

As built, by the rules above: freeze, unfreeze, open and `CR-STATUS` hold. Close is an entry to `X` from any state: `2220-CLOSE` tests only the balance, and the test that the account is not closed is in `2200-APPLY-REQUEST`, which goes to it; a frozen account can be closed, so `O` to `X` would say less than the code. `HD-STATUS` to `R` is from any state for the same reason: `AUTHREV` tests `HD-REVERSED` in `0000-MAIN`, false, and sets it in `2000-APPROVE`, which it performs. Following decisions across paragraphs would need the paths between them, which the rules do not follow; the page states the limit.

One pull request, commits: condition refs in `rules` and the model; the field page's rules section; the cross-program field page; the entity pages; edge labels in `diagram`; state machines; docs.

## Part 3: ownership and the boundary

### Ownership

A page, `/ownership`, grouping programs by the entities they write: for each entity, the programs that write a dataset described by it, and the programs that only read one. Draw it as one diagram, entities as data nodes and programs as boxes, a solid edge from a program to an entity it writes and a dotted one to an entity it reads, legend included. Under the diagram, a table of the groups, and a note that this is the first draft of a service decomposition, derived from who writes what, not a design.

Checks: the master is written by `TXPOST` and `ACCTMNT` and read by nine programs; accruals by `INTACCR` alone; holds, cards and the online accounts by `AUTHSRV`, with `ACCTLOAD` and `CARDLOAD` loading them; the ledger by `GLSUM` alone.

### The boundary

A page, `/boundary`, with the system's contracts with the outside: the external datasets with their copybooks, linked to the entity and the layout; the outputs, datasets no step reads; the reports, every SYSOUT DD by job and step, with the program's file and record; and the in-stream data, every SYSIN of a program that ACCEPTs, with what the program DISPLAYs as its reply. State plainly that the survey sees the files and the streams, not what is on the other end: the branches' systems and the HTTP server around `AUTHSRV` are outside it.

Checks: inputs `TRAN.DAILY` (`TRANREC`), `ACCT.MAINT`, `BANK.RATES`, `CARD.MASTER`, and `AUTHSRV`'s SYSIN; outputs `GL.SUMMARY`, the rejects, every `*RPT` report, `STMTPRT`'s statements, and `AUTHSRV`'s responses.

One pull request, commits: ownership; the boundary; docs.

## Part 4: forward lineage

Today `lineage.Back` follows a field's values to where they came from. Add the other direction, the one design needs for impact: change this, and here is what it touches.

- `dataflow.Flow.Affects(item)`: the mirror of `Trace`: through every statement that moves a value out of the item or out of bytes it shares, to the items it reaches, on to the files it is written to, with the statements followed. Same byte-position following through groups and redefinitions, same tests on `ORIGINS.cbl` turned around.
- `blueprint.Program.Affects(item)` and, for every record a program reads from a file, the fields of it, so that a reader's field can be found over the same bytes a writer wrote: this is `outputs` in reverse, call it `inputs`, built the same way from the lineage's `FromFile` edges.
- `lineage.Forward`: from a field, to the files its program writes it to, to the datasets the step binds, to every step that reads them, to the reader's field over the same bytes, on through sorts, until it reaches an output of the boundary: a report, a dataset no step reads, or a `DISPLAY`. Share the dataset and sort resolution with `Back`; refactor what is common rather than copy it.
- The trace page takes `direction=forward`, defaulting to back, with the same hop list and the same highlighting on the job networks and the flows. The field page links to both: "Follow its values back" and "Follow where its values go".

Checks: forward from `TD-AMOUNT` in `TXVALID` reaches `TXPOST`'s `AM-BALANCE`, the audit trail, and from there `STMTPRT`'s statement lines and `GLSUM`'s ledger amounts; forward from `BC-BUSINESS-DATE` in `BANKDATE` reaches nearly everything, which the page should survive without timing out, so cap the hops the way `Back` caps sort depth and say when it did.

One pull request, commits: `Affects` in `dataflow`; inputs in the model; `Forward` with the shared resolution; the trace page's direction; docs.

## Part 5: the order of jobs, derived

The JCL does not say which job runs when; the operations guide does. But the generations say it too: a job that writes `NAME(+1)` runs before a job that reads `NAME(0)`, and a job that reads `NAME(-1)` runs after the one that wrote the generation before. Derive the order from the jobs' datasets and draw it: jobs as boxes, an edge from the job that writes a generation to the one that reads it, labeled with the dataset. A cycle is the business day's loop and should be drawn as one: `DAILY` writes `BANK.CONTROL(+1)` that `ONLINE` and the next `DAILY` read. Put it on a page, `/schedule`, with the derived order in words under the diagram, and the note that the survey derives the order from the data and cannot know the calendar: that `MONTHLY` runs on the first business day of a month is in the operations guide, not in the files.

Checks: `START` before `ONLINE`; `ONLINE` before `DAILY` through `TRAN.CARD`; `DAILY` before `MONTHLY` through `TRAN.AUDIT` and `ACCT.ACCRUAL`; `DAILY` before the next `ONLINE` through `BANK.CONTROL`, `ACCT.ONLINE`, `CARD.ONLINE` and `HOLD.ONLINE`.

One pull request.

## Finishing

- The explorer's overview links to the new pages in a section of their own, **Design views**, with one line each on what it is for.
- [blueprints.md](../blueprints.md) gains the kinds and the datasets table; [survey/README.md](../../survey/README.md) gains the views and their limits; [roadmap.md](../roadmap.md) marks milestone 9 done, item by item.
- Take screenshots of each new page on Pangaea with headless Chrome and look at them before merging. A page that is correct and unreadable is not done.
- The field, entity and state-machine views are the ones that most change how the data reads. If the work has to stop early, stop after part 2.
