# Bedrock

The harness that the survey's oracle and the specimens' machine rooms run the specimens with. Nothing here is analysis, and nothing here reads COBOL: that is [cobolith](https://github.com/vlumi/cobolith), the library bedrock stands on.

- [script/](script/): batch scripts in a subset of POSIX sh: variables, exported `DD_` bindings, programs with their output redirected, return codes tested with `if`, `sort`, `mv`, and one script running another.
- [batchfile/](batchfile/): Windows batch files in a subset of their own, read into the same model as the scripts, and run, standing in for `cmd.exe` as `runner` stands in for JES: variables set, `DD_` bindings among them, programs with their output redirected, `echo` into a log, `move`, `call`, and `if errorlevel`.
- [dataset/](dataset/): a system's datasets, files of fixed-length records named by dataset name, with generations, and the catalog that lists them with their record lengths, indexed ones among them.
- [runner/](runner/): running jobs, standing in for the job entry subsystem: each step's program with its DD statements bound to datasets, steps bypassed by their `COND` tests, new datasets and generations cataloged, temporary ones deleted, and a job log.
- [runlog/](runlog/): a system's run record, a line for every step its machine room ran, the business date, the job, the step, the program and the outcome, the same for every system whatever its own logs look like.
- [run/](run/): compiling programs with GnuCOBOL in the dialect they are written in, IBM, Micro Focus or Fujitsu, read in its encoding, as executables or as modules for the programs called with `USING`, and running one at a time with its DD names bound to files, its SYSIN on standard input, and optionally a trace of every statement it executes.

The survey reads scripts with `script` to draw Laurasia's batches, and runs programs, jobs and scripts with the rest for its oracle. Pangaea's machine room runs its jobs with `runner`, each step through `run`; Gondwana's runs its batch files with `batchfile`. All three machine rooms keep their run record with `runlog`, and the survey's explorer reads it. `run` reads a program's header with cobolith's `cobol` to decide how to compile it, and `runner` its jobs with cobolith's `jcl`. Bedrock imports nothing of the survey or the specimens.
