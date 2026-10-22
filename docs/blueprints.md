# Reading a blueprint

What the survey writes for a system, and what the words in it mean. The blueprints of the First Bank of Pangaea in [blueprints/pangaea/](../blueprints/pangaea/README.md) are the example throughout.

Every blueprint exists twice: as Markdown, for reading, and as JSON next to it, for tools, with the same facts. Every source reference links to the file and line it came from, in the program or in the copybook, so that nothing has to be taken on trust.

## The index

`README.md` lists the programs with their kind, files, paragraphs and decision tables, the copybooks with their records, or the paragraphs of a copybook of procedure, and the programs that use them, the jobs with their steps and findings, and the system's datasets. `datasets.json` next to it has the datasets for tools, `entities.json` the entities, and `schedule.json` the order of the jobs: a job that writes a new generation of a dataset, `NAME(+1)`, runs before a job that reads its current one, an earlier one, or every generation, with the datasets between them.

An entity is a record a copybook describes. For each one, `entities.json` has the datasets its copybook describes, the programs whose job steps write one of them and those that only read them, which is who owns the entity, the programs that use it, and for each field, what each program does with it under its own name for it: the paragraphs that read and write it and the decisions on it. A program's record is the copybook's where it is declared at the same place, so a copy with `REPLACING` counts. For a code field, it also has the state machine: the condition names are the states, and a transition is a rule of a decision table that finds the field in one state, by its condition name or `FIELD = 'value'`, and sets another, by `SET` or `MOVE 'value'`. A rule that sets a state without finding the field in one starts from the states the rules that perform its paragraph hold the field in, and theirs in turn, or from `any` where some way in, a `PERFORM` that knows no state, a `GO TO` or a fall-through, knows none; so does a paragraph that decides nothing and sets one. A value moved in from another field is listed as a copy, not a transition. Each transition carries the outcomes of the conditions that decide it, those of the `PERFORM`s on the way in and then its rule's: why it happens, so that a field of reasons, such as a rejection's, lists the system's validations. A condition decides a rule's action where the rules part on it, one that found the conditions before it the same finding it otherwise and not taking the action; a condition every rule finds the same way, or one some rule finds otherwise and takes the action all the same, is not a reason, so a transition every path of a paragraph takes is listed once, with none. The entity's links with the other entities and the control totals of its records are gathered from the programs. Each field has its layout, its picture and usage, and a date field, or a text field that is not a redefinition, lists the date fields in another layout its values come from in the programs that write it: one business date in eight digits, six, and text. The explorer's entity pages group the dates of all the entities by layout.

A program is **main**, or **called** when it is passed items with `PROCEDURE DIVISION USING`.

The datasets are every dataset the jobs name, by the name without its generation, each with its copybook, whether it is indexed, the job steps that write and read it, and its kind:

- **external**: no step writes it, so it comes from outside the system, as the day's transactions do;
- **output**: no step reads it, so it goes outside, as the rejects and the ledger do;
- **temporary**: `&&NAME`, it lives within its job;
- **internal**: anything else, passed from step to step and from job to job.

## Kinds of field

Every field of a layout, in a program or a copybook, has a kind, the first of these that fits:

1. **group**: it has fields of its own.
2. **filler**: it is named `FILLER`.
3. **key**: it is the `RECORD KEY` of one of the program's files, or a key one of its sorts sorts on. A copybook's field is a key where a program that uses the copybook says so. A key is only what the code says is one: the account master's account number is read in sequence, so it is a number.
4. **code**: it has condition names, level 88.
5. **edited**: its picture is numeric-edited, for printing.
6. **rate**: numeric, named for a rate: `RATE`, `FACTOR` or `PERCENT` in its name, or `ｹｲｽｳ`, or ending in `ﾘﾂ`.
7. **amount**: numeric with decimals; or without, named for an amount, `AMOUNT`, `AMT` or `BALANCE`, or `ｷﾝｶﾞｸ` or `ｻﾞﾝﾀﾞｶ`, or ending in `ｶﾞｸ`, `ﾘｮｳ` or `ｷﾝ`, as a currency without a minor unit leaves it.
8. **date**: numeric without decimals, of 4, 6, 8 or 14 digits, or text of 4 to 8 bytes, named for a date: `DATE`, `MONTH`, `TIMESTAMP` or `YYMM` in its name, or `ﾋﾂﾞｹ`, `ﾈﾝｹﾞﾂ`, `ｶﾞｯﾋﾟ`, `ﾜﾚｷ` or `ｾｲﾚｷ`, or ending in `ﾋﾞ`.
9. **number**: any other numeric field.
10. **text**: anything else.

Telling a rate, an amount or a date by its name is the one rule that reads a name, and the only heuristic among the kinds. What the code does with a field then carries the kind along: a numeric field an amount is moved or computed into is an amount; a field a date is moved into, of a size a date takes, is a date, as is an argument of `INTEGER-OF-DATE`; and a field of an entity is of the kind any program's record of it shows. A date's arithmetic is not followed: the days between two dates are a number.

## A program, `programs/NAME.md`

In order:

- **The program's own description:** its header comment, as written, modification history included. The survey adds nothing to it.
- **Files:** every file the program declares, the DD name it is assigned to, its organization and record key, its status field, the modes it is opened in, and its records.
- **Flow:** the paragraphs and how control moves between them, `PERFORM` solid, `GO TO` dashed, falling into the next paragraph dotted, with the first paragraph emphasized. Each paragraph links to its source.
- **Paragraphs:** each with its comment, what it performs and goes to, and what performs or reaches it.
- **Displays:** in the JSON, the program's `DISPLAY` statements to the job's output.
- **Calls:** for a program that is called, what it is passed and by whom; for one that calls, what it calls and what it passes.
- **Data flow:** what each paragraph reads and writes.
- **Lineage:** every statement that moves a value into a data item, grouped by the record the item belongs to: from which item, literal, file, program called or the job's SYSIN, by which verb, in which paragraph; and under **Out of the program**, the records written to files, the items passed to programs called, by argument, and the items displayed to SYSOUT.
- **Inputs:** in the JSON, for every record the program reads from a file, where the values of each field go: the data items they reach and the files they are written to. The mirror of the field origins, which the explorer's forward lineage follows from.
- **Field origins:** for every field of every record the program writes to a file, by file and, for a file of several records, such as a report's headings, lines and totals, by record, and for a record of types, such as a header, details and a trailer laid out in redefinitions of one record, by type, from the writes of that type alone, where its values come from in the end: the fields of the files it reads, literals, initial values, programs called that may change an item passed to them (`call AUTHPUR argument 6 HD-AMOUNT`), the caller of a called program (`parameter 1 RQ-AMOUNT`) and the job's SYSIN (`stream SYSIN RQ-AMOUNT`), followed back from the statements that write the record, along the paths control can take to them, through the bytes that groups and redefinitions share, and no further than a statement that overwrites those bytes.
- **Links between records:** the fields of one input the program relies on to find or match the records of another: moved into the record key of an indexed file it reads, or compared for equality, each side's values from that field alone. Derived, not declared: the relations of the data model.
- **Control totals:** the fields of an input the program checks for equality against what it adds up from the input's other records, a field's sum or their count, over the records of a type.
- **Parameters:** for a program called with `USING`, which fields of each item it is passed it reads and writes: the contract its callers rely on.
- **Records:** for a record laid out in redefinitions, its types: the field whose condition names tell them apart, as the program tests or sets them, and each type's redefinition, paired in order where there are as many of each. Then the layout of every record, field by field, with kind, position, length, picture, usage and notes. The fields are drawn as the tree the levels declare, condition names under their field, and a field that redefines the one above it hangs from a double branch, `╞`, with the redefined field's notes naming every field that redefines it; a record from a copybook points at the copybook's blueprint instead.

## A copybook, `copybooks/NAME.md`

Its description, the programs that use it, and the layout of each record in it, with each field's kind.

## A job, `jobs/NAME.md`

- **The job's own description:** its comment banner.
- **Steps:** each step, its program, the datasets it reads and writes, what bypasses it, and its description from the JCL's comments.
- **Network:** the steps and the datasets between them. A dataset one step writes and another reads is drawn once, between them; one only read or only written is drawn next to its step. Temporary datasets are dashed, external ones drawn in a warm fill of their own, and each dataset names the copybook that describes its records, where a program says so, and its kind and whether it is indexed, unless it is internal and sequential.
- **The DD statements of each step:** the dataset, its disposition, the program's file bound to it, that file's copybook, and how the program uses the file.
- **Datasets:** every dataset the job names, the copybook that describes it, its kind across the system, and the steps of the job that write and read it.
- **Findings:** where the job and its programs disagree. A DD no file is assigned to; a file no DD is bound to; a program that writes a dataset it shares with `DISP=SHR`, appends to one that is not `DISP=MOD`, or only reads one it creates; a new dataset organized otherwise than its file; a temporary dataset read before any step writes it, or written and never read; a step running a program that was not analyzed; a program that ACCEPTs with no SYSIN; a dataset two programs describe with different copybooks. A job with none says so.

## Decision tables, `rules/NAME.md`

One table per paragraph that decides anything, and one per inline loop. Each column is one path through the paragraph; the JSON adds, for each condition, the data items it decides on, a condition name counting as its field. The rows above the line are the conditions, with `Y` or `N` where the path tests them and blank where it does not; the rows below, marked with an arrow, are the actions, numbered in the order the path takes them. A path through an `EVALUATE` is also a path past the `WHEN`s before it, and a path ends at a `GO TO`, `STOP RUN` or `EXIT PARAGRAPH`, or at a `PERFORM` of a paragraph that never returns, such as one that aborts the run. A condition a path tests twice has the same outcome both times, unless a statement between the tests may change what it refers to: a move into its bytes, a call to a program that writes them, by what the program called writes of what it is passed, or a paragraph performed that does either; a program outside the system may write all it is passed by reference. Then the second test is a row of its own, `CONDITION, again`. A path also knows the values it has given fields, by moving a literal or setting a condition name, and those its tests have found, until a statement may change them; a path that would find a condition other than they decide is no rule, so a table does not move an error message into a field and then find the field's condition name for no error true. A paragraph entered by `PERFORM` alone starts its paths knowing what every `PERFORM` of it knows alike: the outcomes of the conditions it tests, and of those the paragraphs it performs test, which carry on down; a `READ`, `WRITE` or other file statement changes the file's status, so a test of the status after it is a new test. A paragraph with more paths than a table can show has a decision shown as an action instead, so the table stays readable.

## The oracle's cases, `oracle/NAME/`

Beside the blueprints, `survey oracle` writes the cases it found for each main program: the evidence of what the programs do when each rule applies.

- **`README.md`** lists the programs, with the rules their cases reach.
- **`PROGRAM/README.md`** lists every rule of the program's decision tables, and of the programs it calls, with the first case that reaches it. A rule no case reaches shows the outcomes of its conditions; most of those are rules no run can take, which the tables do not rule out. Two kinds are marked with the reason the search cannot reach them, and `coverage.json` lists them as `ownTotals` and `failedWrites`:
  - a rule that finds two totals the program adds up itself to differ, items of its own working storage it only adds to from zero, since no input field sets either;
  - a rule that finds a WRITE or REWRITE by key to fail, since the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile.
- **A case** is a directory:
  - `README.md`: how the run ended, and the rules it is the first to reach;
  - `input/FILE.txt`: each input file, a record a block, a field a line, `NAME VALUE`, a number as a decimal and text quoted; `SYSIN.txt` is what the program accepts;
  - `input/failed.txt`, for a case that fails a file on purpose: the file and how. `missing`, a file read that is not there, so opening it fails; `unreadable`, a directory in place of a sequential file read, so it opens but reading it fails; `uncreatable`, a file written whose directory is not there, so creating it fails. GnuCOBOL reports each the same way on every platform;
  - `output/FILE.txt`: what it wrote to each file it writes or updates, a record a line, quoted with escapes where it holds bytes that do not print; `STDOUT.txt` is what it displayed.
- **`jobs/JOB/`** is a batch: the jobs run in the order the operator gives, Pangaea's being START, ONLINE, DAILY and MONTHLY, since the JCL does not say when a job runs. Each job runs on one catalog, on what the jobs before it wrote; a dataset none of them wrote takes the records of the cleanest case, the first with the lowest return code, of the first program that reads it. `README.md` has how each step ended and the rules each program reached; `input/DSN.txt` the datasets taken from the cases, field by field; `output/DSN.txt` each dataset the job created or changed, a record a line; and `sysout/` the job's output, its log among it.
- **A batch of scripts**, Laurasia's open, nightly and cycle, runs them with sh, as they are run: on a directory for each variable of their environment that names data, with the oracle's own stand-ins for the programs, scripts, sorts and moves their steps run, which record how each ended. A parameter, such as the business date, takes the command line of the cleanest case of the first program it is passed to. Its files are named by that variable and their path in its directory, `input/LCS_DATA/inbound/FEES.txt`, and `output/` has every file the scripts created or changed, a program's as the program reads it, an indexed or relative one unloaded.

The tests replay every case and the batches, and compare what they reach and write, so a change in a program's behavior, or in the survey's tables, shows as a difference.

## Diagrams

Every diagram carries its legend. The same engine draws the program flows, the job networks, the console's run networks, and the diagrams in this documentation, deterministically, so a diagram is a golden file like the text around it.
