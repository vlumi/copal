# Bedrock

The language layer that the survey and the bank's machine room both stand on. Nothing here is analysis; it is how COBOL source, copybooks, records and JCL are read, and how COBOL programs are compiled and run.

- [cobol/](cobol/): lexing fixed-format and free-format source, expanding `COPY` with `REPLACING`, and parsing the divisions, data entries and statements of the subset the specimens use, every token with the file and line it came from; and the dialects, IBM and Micro Focus, and what they lay out differently.
- [layout/](layout/): record layouts from data description entries, in a dialect: position, length, picture, usage, `OCCURS`, `REDEFINES` and condition names.
- [record/](record/): reading and writing the fields of a record through its layout, including packed decimal and the overpunched signs of display numbers.
- [jcl/](jcl/): jobs, steps and DD statements from job files, with generations, temporary datasets and concatenations.
- [script/](script/): batch scripts in a subset of POSIX sh: variables, exported `DD_` bindings, programs with their output redirected, return codes tested with `if`, `sort`, `mv`, and one script running another.
- [dataset/](dataset/): a system's datasets, files of fixed-length records named by dataset name, with generations, and the catalog that lists them with their record lengths, indexed ones among them.
- [runner/](runner/): running jobs, standing in for the job entry subsystem: each step's program with its DD statements bound to datasets, steps bypassed by their `COND` tests, new datasets and generations cataloged, temporary ones deleted, and a job log.
- [run/](run/): compiling programs with GnuCOBOL in the dialect they are written in, IBM or Micro Focus, as executables or as modules for the programs called with `USING`, and running one at a time with its DD names bound to files, its SYSIN on standard input, and optionally a trace of every statement it executes.

The survey uses the first four to analyze, and the rest for its oracle, which runs the programs it reads. Pangaea's machine room runs its jobs with `runner`, each step through `run`, and its generator and console read and write the bank's records through `layout` and `record`. `run` reads a program's header with `cobol` to decide how to compile it. Bedrock imports nothing from either.
