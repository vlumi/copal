# Plan: a deterministic sweep

Before a model is let near the blueprints, what the code already says should all be in them, and right. A sweep over the analysis found facts that are wrong or too coarse, and facts the survey derives for its own use but does not publish. This plan fixes and publishes them, in five parts, one pull request each, in this order.

The ground rules of [design-views.md](design-views.md) apply: facts only, new facts in the model, goldens with the change that causes them, a test for every step, screenshots of every changed page before merging, `make ci` before pushing. The oracle's cases are goldens too: a change that moves the tables moves the cases' coverage, and `make oracle` searches for them again.

## Part 1: wrong facts

- **A paragraph that never returns ends a path.** Pangaea performs an abort paragraph, which stops the run, 36 times; 41 of the 384 rules take statements after it, in 22 tables, and every decision after it doubles the paths. A paragraph never returns when every path through it ends in `STOP RUN` or `GOBACK`, or in a `PERFORM` of one that never returns; a rule's path ends at a `PERFORM` of it.
- **A call changes only what the program called writes.** Within one program, every program called by reference may change what it is passed, so AUTHSRV's authorization table grew past what a table shows; none of the programs it calls writes the request. The blueprint, which has every program, tells data flow which of the items a `CALL` passes the program called writes, through its `USING` items, by offset; a program outside the system still may change all it is passed by reference. The oracle takes its tables from the blueprint, so that the cases and the blueprints agree.

Checks: no rule takes a statement after a `PERFORM` of an abort paragraph; AUTHSRV's `2100-AUTHORIZE` is a table again, with `RS-APPROVED` among its conditions.

As built: both checks hold. The tables went from 384 rules to 368 with the paragraphs that never return, and to 408 with AUTHSRV's table back; the oracle's search reaches 266 of them, against 245 of 384 before. What AUTHSRV's table gained is mostly combinations of request types no run takes, which part 2 is for.

## Part 2: rules no run can take

Most rules no case reaches are rules no run can take: a table does not weigh what one condition's outcome says about another's, or what a statement on the path has just moved. Along each path, keep the values the path knows: a literal moved into a field, or set by a condition name, and a condition's outcome tested before. A condition on a field whose value the path knows has a known outcome, and a path that would give it the other is no rule. A value is forgotten where a statement may change the field, as part 1 tells.

Checks: TXVALID's header check loses the rules that find the file good right after moving an error into it; the oracle's unreached rules shrink by those, and the rest are listed with their conditions as now.

## Part 3: record types

The oracle tells a file's record types apart, by the field the program tests to tell them, and the redefinition each type has. Make that a fact of the blueprint: each record's types, their values and their fields. Then:

- Field origins and inputs by record type: a field of the trailer gets the trailer's origins, not the detail's over the same bytes, as a statement that sets the type on a path says which type the path writes.
- A program's inputs have every record type it reads, not the first.
- A program that updates a dataset in place, reading and rewriting it, counts among its writers in lineage.

## Part 4: design facts between records

- **Keys and links:** a field moved into an indexed file's key to read it, or compared with another file's key to match it, refers to that file's records; with the record order a file's types are tested in, and its control totals, a field compared with an item the program adds another field up into. The oracle infers these for its own use; publish them, and show them on the entity pages as the data model's links.
- **Interfaces:** for a program called with `USING`, which of its parameters' fields it reads, writes, or both: the contract a replacement has to keep.

## Part 5: smaller facts

- **Rejection reasons:** for each code field, the rules across the programs that set each of its values, with their conditions: the system's validations in one table.
- **Kinds from lineage:** a text field that holds a date, by its origins, is a date in text.
- **Decision tables in the explorer**, beside their coverage by the oracle's cases.
