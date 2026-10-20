# Bedrock

The harness that the survey's oracle and the bank's machine room run the specimens with. Nothing here is analysis, and nothing here reads COBOL: that is [cobolith](https://github.com/vlumi/cobolith), the library bedrock stands on.

- [script/](script/): batch scripts in a subset of POSIX sh: variables, exported `DD_` bindings, programs with their output redirected, return codes tested with `if`, `sort`, `mv`, and one script running another.
- [dataset/](dataset/): a system's datasets, files of fixed-length records named by dataset name, with generations, and the catalog that lists them with their record lengths, indexed ones among them.
- [runner/](runner/): running jobs, standing in for the job entry subsystem: each step's program with its DD statements bound to datasets, steps bypassed by their `COND` tests, new datasets and generations cataloged, temporary ones deleted, and a job log.
- [run/](run/): compiling programs with GnuCOBOL in the dialect they are written in, IBM or Micro Focus, as executables or as modules for the programs called with `USING`, and running one at a time with its DD names bound to files, its SYSIN on standard input, and optionally a trace of every statement it executes.

The survey reads scripts with `script` to draw Laurasia's batches, and runs programs, jobs and scripts with the rest for its oracle. Pangaea's machine room runs its jobs with `runner`, each step through `run`. `run` reads a program's header with cobolith's `cobol` to decide how to compile it, and `runner` its jobs with cobolith's `jcl`. Bedrock imports nothing of the survey or the specimens.
