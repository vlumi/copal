# Bedrock

The language layer that the survey and the bank's machine room both stand on. Nothing here is analysis; it is how COBOL source, copybooks, records and JCL are read.

- [cobol/](cobol/): lexing fixed-format source, expanding `COPY` with `REPLACING`, and parsing the divisions, data entries and statements of the subset Pangaea uses, every token with the file and line it came from.
- [layout/](layout/): record layouts from data description entries: position, length, picture, usage, `OCCURS`, `REDEFINES` and condition names.
- [record/](record/): reading and writing the fields of a record through its layout, including packed decimal and the overpunched signs of display numbers.
- [jcl/](jcl/): jobs, steps and DD statements from job files, with generations, temporary datasets and concatenations.

The survey uses all four to analyze. Pangaea's runner executes the jobs `jcl` reads, its generator and console read and write the bank's records through `layout` and `record`, and its compiler reads a program's header with `cobol` to decide how to compile it. Bedrock imports nothing from either.
