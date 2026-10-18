# Bedrock

The language layer that the survey and the bank's machine room both stand on. Nothing here is analysis; it is how COBOL source, copybooks, records and JCL are read, and how COBOL programs are compiled and run.

- [cobol/](cobol/): lexing fixed-format source, expanding `COPY` with `REPLACING`, and parsing the divisions, data entries and statements of the subset Pangaea uses, every token with the file and line it came from.
- [layout/](layout/): record layouts from data description entries: position, length, picture, usage, `OCCURS`, `REDEFINES` and condition names.
- [record/](record/): reading and writing the fields of a record through its layout, including packed decimal and the overpunched signs of display numbers.
- [jcl/](jcl/): jobs, steps and DD statements from job files, with generations, temporary datasets and concatenations.
- [run/](run/): compiling programs with GnuCOBOL, as executables or as modules for the programs called with `USING`, and running one at a time with its DD names bound to files, its SYSIN on standard input, and optionally a trace of every statement it executes.

The survey uses the first four to analyze. Pangaea's runner executes the jobs `jcl` reads, each step through `run`, and its generator and console read and write the bank's records through `layout` and `record`. `run` reads a program's header with `cobol` to decide how to compile it. Bedrock imports nothing from either.
