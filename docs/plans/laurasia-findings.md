# What the survey makes of Laurasia

Part 2 of [second-specimen.md](second-specimen.md): the survey run on Laurasia Card Services as it was after part 1, unchanged, and every place it breaks or comes out wrong. Each is a coupling to Pangaea, to be fixed generically in the part named, never by special-casing Laurasia.

## How it was found

As is, the survey reads none of Laurasia's 14 programs: each stops at its first error, all of them in the parser or the layouts. To see past the first error, a throwaway copy of the analyzer was patched just enough to go on, one break at a time, and so was a copy of Laurasia's sources for the two constructs it could not get past: the nested program was moved to a file of its own and the declaratives were taken out. With those, all 15 programs, the nested one counted, parse and the blueprints are written; what they say was then read against the code. None of the patches is kept. There are no golden blueprints of Laurasia yet, as there is nothing to write until the parser reads it; the first come with part 3.

## Where the parser stops

Each of these stops the survey of every program that has it. They are part 3's.

- **1. Area A in free format.** A paragraph or section header is recognized only in area A, columns 8 to 11 of fixed format. Free format has no areas, so `main section.` at column 1 is taken for a statement, and every program fails with "statements before the first paragraph". All 14 programs.
- **2. Statements of a section.** A section's own statements, before any paragraph in it, are rejected: Laurasia's sections have no paragraphs. A `PERFORM` of a section runs its statements and then its paragraphs, through the last one in the section. All 14.
- **3. `SIGN LEADING SEPARATE`** is rejected, in seven of the nine copybooks and in CYCLE, DISPUTE, INTCALC and SETTLE. Once read, the layout must count the separate sign's byte: without it, CARD-ACCOUNT lays out as 59 bytes and not 61, PRES-RECORD 74 and not 76, CLR-PRESENTMENT 60 and not 61, and every field after the sign is at the wrong position.
- **4. `COMP-5`**, binary in the machine's own byte order, is rejected: CLRIN, CYCLE, FEELD, FEESLOT, INTCALC, IXFEE and the fee table's record. It is not `COMP`, which is big-endian, so reading and writing its values differs too.
- **5. Level 78**, a named constant, is rejected by the layouts: CLRIN, MATCH, IXFEE, OVLFEE, CYCLE. A constant's value is what a rule compares with, as MATCH's `PS-AMOUNT > AL-AMOUNT * TOLERANCE` is 1.15 times the approval.
- **6. `STOP RUN RETURNING`** is rejected, in ten programs. Its value is the program's return code, which the scripts test and the oracle records.
- **7. `RELATIVE KEY`** is rejected: FEELD, IXFEE. A relative file's key is a slot number in working storage, not a field of the record as an indexed file's record key is; the links and the oracle assume the latter.
- **8. `DECLARATIVES`** are taken for a paragraph named DECLARATIVES, and `USE` is rejected: CLRIN. A declarative section is not performed; it runs when an input or output statement on one of the files it names fails, and control comes back after that statement.
- **9. A nested program** is a "second IDENTIFICATION DIVISION": CYCLE, which contains MINPAY. It is a program of its own, callable only from the program around it, and `END PROGRAM` closes it.

## Where it comes out wrong

With the parser patched, these facts are wrong or missing.

### Control flow and rules, part 3

- **10. `EXIT SECTION` does not end the path.** A rule goes on past it to the statements after it. POST's table of `POST-PRESENTMENT` has 24 rules where the code has six, and in most of them a presentment rejected for an unknown, blocked or closed card is posted all the same. DISPUTE's `WHEN OTHER` rejects an event and then stamps the case as updated; FEELD writes a fee it found no slot for, and CLRIN goes on to take a batch from a record it has just rejected as out of order.

### Calls, part 3

- **11. `CALL ... RETURNING`.** `RETURNING` and the item after it are listed among the items passed: CYCLE, FEELD, IXFEE. The item is written by the call, with what the program called returns.
- **12. Lower-case program names.** `call "intcalc"` is not matched with the program INTCALC, so every call in Laurasia is taken for one outside the system: no parameter contracts across the call, no writes through it, and the oracle would not take the callee's tables as its own.
- **13. A call through a data item.** POST's `call ws-fee-program` is not listed as a call at all. The literals moved into `ws-fee-program`, "ovlfee" and "ovlbiz", are the programs it may call.

### Records, part 3

- **14. Several 01 records under one FD are not record types.** CLRIN's inputs list only CLR-FILE-HEADER, the first of the clearing file's five records, and none has types, though the condition names on CLR-RECORD-TYPE say which is which. The oracle also takes a file's first record as its layout. Pangaea's record types are redefinitions within one record.
- **15. Control totals across records.** CLRIN checks the batch trailer's amount against the sum of the presentments' amounts, and the file trailer's likewise; neither is found, as a total is looked for only within the record it is checked in. The counts are found, with no record type.
- **16. `RECORD VARYING ... DEPENDING ON`** is passed over without a word: a file of variable-length records is described as fixed.

### Design facts, part 5

- **17. Links by ordered comparison.** MATCH, SETTLE and DISPUTE merge two files in key order: `<` reads one, `>` the other, and `WHEN OTHER` is the match. Only `=` makes a link today, so none of the three has one; POST's key read is the only link found.
- **18. A state machine in compound conditions.** DISPUTE's EVALUATE tests the case's state and the event together, `HOLDING-CASE ALSO DC-CHARGED-BACK AND DE-REPRESENTMENT`. The decision table is right, but every transition of DC-STATUS is from "any", as the state it is from is inside an `AND`, and each says only that the first `WHEN` was false.

### The batch, part 4

- **19. No jobs.** The survey reads JCL and nothing else, so Laurasia has no jobs, datasets, schedule or lineage across programs: the scripts are its job network.
- **20. The run date on the command line.** CLRIN, POST and CYCLE take the run date with `ACCEPT ... FROM COMMAND-LINE`, the argument the scripts pass, an input the oracle does not give as it gives SYSIN.

### The oracle, part 5

- **21. The dialect.** The oracle compiles every program in IBM dialect, and GnuCOBOL rejects free-format source in it before the search starts.
