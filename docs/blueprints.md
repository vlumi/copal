# Reading a blueprint

What the survey writes for a system, and what the words in it mean. The blueprints of the First Bank of Pangaea in [blueprints/pangaea/](../blueprints/pangaea/README.md) are the example throughout.

Every blueprint exists twice: as Markdown, for reading, and as JSON next to it, for tools, with the same facts. Every source reference links to the file and line it came from, in the program or in the copybook, so that nothing has to be taken on trust.

## The index

`README.md` lists the programs with their kind, files, paragraphs and decision tables, the copybooks with their records and the programs that use them, the jobs with their steps and findings, and the system's datasets. `datasets.json` next to it has the datasets for tools, `entities.json` the entities, and `schedule.json` the order of the jobs: a job that writes a new generation of a dataset, `NAME(+1)`, runs before a job that reads its current one, an earlier one, or every generation, with the datasets between them.

An entity is a record a copybook describes. For each one, `entities.json` has the datasets its copybook describes, the programs whose job steps write one of them and those that only read them, which is who owns the entity, the programs that use it, and for each field, what each program does with it under its own name for it: the paragraphs that read and write it and the decisions on it. A program's record is the copybook's where it is declared at the same place, so a copy with `REPLACING` counts. For a code field, it also has the state machine: the condition names are the states, and a transition is a rule of a decision table that finds the field in one state, by its condition name or `FIELD = 'value'`, and sets another, by `SET` or `MOVE 'value'`. A rule that sets a state without finding the field in one, or a paragraph that decides nothing and sets one, starts from `any`. A value moved in from another field is listed as a copy, not a transition.

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
6. **amount**: numeric with decimals.
7. **date**: numeric without decimals, of 6, 8 or 14 digits, with `DATE`, `MONTH`, `TIMESTAMP` or `YYMM` in its name. This is the one rule that reads a name, and the only heuristic among the kinds.
8. **number**: any other numeric field.
9. **text**: anything else.

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
- **Field origins:** for every field of every record the program writes to a file, by file and, for a file of several record types such as a report's headings, lines and totals, by record, where its values come from in the end: the fields of the files it reads, literals, initial values, programs called that may change an item passed to them (`call AUTHPUR argument 6 HD-AMOUNT`), the caller of a called program (`parameter 1 RQ-AMOUNT`) and the job's SYSIN (`stream SYSIN RQ-AMOUNT`), followed through every statement and through the bytes that groups and redefinitions share.
- **Records:** the layout of every record, field by field, with kind, position, length, picture, usage and notes. The fields are drawn as the tree the levels declare, condition names under their field, and a field that redefines the one above it hangs from a double branch, `╞`, with the redefined field's notes naming every field that redefines it; a record from a copybook points at the copybook's blueprint instead.

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

One table per paragraph that decides anything, and one per inline loop. Each column is one path through the paragraph; the JSON adds, for each condition, the data items it decides on, a condition name counting as its field. The rows above the line are the conditions, with `Y` or `N` where the path tests them and blank where it does not; the rows below, marked with an arrow, are the actions, numbered in the order the path takes them. A path through an `EVALUATE` is also a path past the `WHEN`s before it, and a path ends at a `GO TO` or `STOP RUN`. A paragraph with more paths than a table can show has a decision shown as an action instead, so the table stays readable.

## Diagrams

Every diagram carries its legend. The same engine draws the program flows, the job networks, the console's run networks, and the diagrams in this documentation, deterministically, so a diagram is a golden file like the text around it.
