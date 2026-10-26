# Plan: a replacement checked

The brief opens with the reason a rewrite keeps the old shape: nobody can prove it produces the same results. The oracle now holds, for every main program of the three specimens, the cases that reach its rules and what the program wrote and returned on each. Those cases are the acceptance test a replacement has to pass, and nothing yet runs one against them. This plan closes that loop with no translation and no model: a harness that runs any executable the way the oracle ran the program, on the recorded cases, and says where it differs.

## Ground rules

- The survey still translates nothing. The replacement is somebody else's, in any language; the harness only runs it and compares.
- A replacement gets what the program got, no more: the case's command line as its argument, its SYSIN on its standard input, each file's path as `DD_ASSIGN` in its environment, and its exit status is its return code. That is the contract bedrock runs the COBOL by, so a replacement that passes here runs in the same job.
- Files are records. A sequential file comes as the program's SELECT describes it, back to back or a line each. An indexed file, which GnuCOBOL keeps as a Berkeley DB file no other runtime reads, comes as a sequential file of its records in key order, the first of a key kept as the oracle's loader keeps it, and is read back the same way: a replacement that updates one writes it back whole, in key order.
- What the runtime itself displayed in the recorded run, GnuCOBOL's warnings of an implicit CLOSE, is no replacement's to match, and is left out of the comparison on both sides.
- The report names the difference the way the blueprints name things: the file, the record, and the fields of the record whose bytes differ, with both values.

## Parts

1. **A program's cases.** `survey conform` takes a system, a program and a replacement, runs the replacement on the program's recorded cases and prints a report: a line for the whole, and a table of the cases with their differences in return code, in each file written, record by record and field by field, and in what was displayed. It fails unless the replacement conforms on every case. `make conform-pangaea PROGRAM=NAME WITH=executable`, and the same for the other two specimens. To show it, BANKDATE, the bank's last daily step, rewritten in Go, and a script that gets it wrong: the test checks that the rewrite conforms on all three cases and that the report names what the script gets wrong, the date it leaves as it was, the return code and the line it does not display. `make conform-demo` builds the rewrite and runs it.

   As built: `oracle.Conform` in the oracle's package, since it reads the cases as the oracle wrote them and stages a run as the oracle stages one. Writing the rewrite showed the harness earning its keep before the harness was finished: the first version counted days with a `time.Duration`, which overflows between 1601 and 2026, and the report named the field and the two dates.

2. **A step of a batch.** The batches in `oracle/NAME/jobs/` run the jobs on the cases; a replacement should be able to take one step of a job, with the datasets the steps before it wrote and the ones after it read, so that a rewrite is proven in its place in the network and not only on its own files.

3. **A replacement that keeps the indexed files.** A refactored COBOL program, or one compiled otherwise, reads the indexed files as they are; the harness should load and unload them through the oracle's loaders for such a replacement, so that the same cases check a change within COBOL too.

Done when: a replacement for any main program of the three specimens is checked against its cases in one command, with a report that names the first field it gets wrong, and a replaced step runs in its job's batch the same way.
