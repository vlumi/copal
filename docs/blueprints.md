# Reading a blueprint

What the survey writes for a system, and what the words in it mean. The blueprints of the First Bank of Pangaea in [blueprints/pangaea/](../blueprints/pangaea/README.md) are the example throughout.

Every blueprint exists twice: as Markdown, for reading, and as JSON next to it, for tools, with the same facts. Every source reference links to the file and line it came from, in the program or in the copybook, so that nothing has to be taken on trust.

## The index

`README.md` lists the programs with their files, paragraphs and decision tables, the copybooks with their records and the programs that use them, and the jobs with their steps and findings.

## A program, `programs/NAME.md`

In order:

- **The program's own description:** its header comment, as written, modification history included. The survey adds nothing to it.
- **Files:** every file the program declares, the DD name it is assigned to, its organization and record key, its status field, the modes it is opened in, and its records.
- **Flow:** the paragraphs and how control moves between them, `PERFORM` solid, `GO TO` dashed, falling into the next paragraph dotted, with the first paragraph emphasized. Each paragraph links to its source.
- **Paragraphs:** each with its comment, what it performs and goes to, and what performs or reaches it.
- **Calls:** for a program that is called, what it is passed and by whom; for one that calls, what it calls and what it passes.
- **Data flow:** what each paragraph reads and writes.
- **Lineage:** every statement that moves a value into a data item, grouped by the record the item belongs to: from which item, literal or file, by which verb, in which paragraph.
- **Field origins:** for every field of every record the program writes to a file, where its values come from in the end: the fields of the files it reads, literals and initial values, followed through every statement and through the bytes that groups and redefinitions share.
- **Records:** the layout of every record, field by field, with position, length, picture, usage and notes; a record from a copybook points at the copybook's blueprint instead.

## A copybook, `copybooks/NAME.md`

Its description, the programs that use it, and the layout of each record in it.

## A job, `jobs/NAME.md`

- **The job's own description:** its comment banner.
- **Steps:** each step, its program, the datasets it reads and writes, what bypasses it, and its description from the JCL's comments.
- **Network:** the steps and the datasets between them. A dataset one step writes and another reads is drawn once, between them; one only read or only written is drawn next to its step. Temporary datasets are dashed, and each dataset names the copybook that describes its records, where a program says so.
- **The DD statements of each step:** the dataset, its disposition, the program's file bound to it, that file's copybook, and how the program uses the file.
- **Datasets:** every dataset the job names, the copybook that describes it, and the steps that write and read it.
- **Findings:** where the job and its programs disagree. A DD no file is assigned to; a file no DD is bound to; a program that writes a dataset it shares with `DISP=SHR`, appends to one that is not `DISP=MOD`, or only reads one it creates; a new dataset organized otherwise than its file; a temporary dataset read before any step writes it, or written and never read; a step running a program that was not analyzed; a program that ACCEPTs with no SYSIN; a dataset two programs describe with different copybooks. A job with none says so.

## Decision tables, `rules/NAME.md`

One table per paragraph that decides anything, and one per inline loop. Each column is one path through the paragraph. The rows above the line are the conditions, with `Y` or `N` where the path tests them and blank where it does not; the rows below, marked with an arrow, are the actions, numbered in the order the path takes them. A path through an `EVALUATE` is also a path past the `WHEN`s before it, and a path ends at a `GO TO` or `STOP RUN`. A paragraph with more paths than a table can show has a decision shown as an action instead, so the table stays readable.

## Diagrams

Every diagram carries its legend. The same engine draws the program flows, the job networks, the console's run networks, and the diagrams in this documentation, deterministically, so a diagram is a golden file like the text around it.
