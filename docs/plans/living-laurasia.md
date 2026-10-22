# Plan: a living Laurasia

Laurasia Card Services ran one fixture: seven cards, nineteen clearing lines, ten dispute events, one night and its cycle, all written by hand. That is enough for the survey, which reads the code, and for the golden files, which check the programs. It is not enough to show that the second specimen is a system rather than an example. The First Bank of Pangaea earned that with `make bank-reset` and `make bank-days`: a generated portfolio, and business days that chain, each starting from what the last one left. Laurasia gets the same, and running it for months shows what one night cannot.

Read [second-specimen.md](second-specimen.md) first; its ground rules apply. Laurasia is fiction, it runs, it knows nothing of the survey, and it imports bedrock and cobolith but never Pangaea, so it has a small calendar and a generator of its own rather than the bank's.

## Ground rules

- The scripts run as the operators run them: `open.sh` once, `nightly.sh DATE` every business day, `cycle.sh DATE` every day, each through `sh` with the programs, the scripts and the data from the environment they expect. The Go around them writes the inbound files and reads the outbound ones, as the card system, the switch and the networks would, and never touches the masters itself.
- Everything follows from a seed. The same seed and the same days give the same files, byte for byte, which a test checks by running two processors side by side.
- A day goes through whole or not at all. The masters are copied before the day; if any step fails, they are put back and everything the day wrote moves aside, so the next run starts the day afresh. A night that half posted and ran again would post twice.

## The switch

Laurasia has no online side of its own, and no hour at which its cards stop working. The networks' authorization switch approves and declines around the clock, standing in for the issuer when the issuer cannot answer, and the night reads nothing of it but its log. The log of a business day holds every authorization since the last night, a weekend's included, with the ones given in stand-in marked. That is how a card issuer's batch window is invisible to its cardholders, and it is a deliberate difference from the bank, whose own ATM and card channel closes for its batch. The generator keeps it so: the switch's view of a card is the balance the last night unloaded, plus what it approved since.

## Parts

1. **The portfolio and the days.** A generator for the card system's extract, the fee schedule and the merchants, and for each day the switch's log, the networks' clearing file and the disputes desk's events; a processor that opens on the portfolio and runs day after day; a command and make targets to drive it; a test that two processors from one seed agree, and one that a failed day is backed out.

   As built: `laurasia/records` loads the record layouts from the copybooks with cobolith, in Micro Focus dialect, and `laurasia/generate` writes and reads the files by field through them. The portfolio is about a thousand cards over the three products, limits and balances by product, a few blocked or closed, cycle days spread over the month, and sixty merchants under five acquirers, each in a line of business with its category code and its tickets. A day's purchases are drawn card by card, by product and by the day of the week; the switch approves against the last unload, declines a card past its limit or not active, and stands in now and then, for a card the processor has not been told of too. The networks present most approvals, more than approved where there are tips, a declined purchase now and then, and a few presentments whose authorization was in an earlier log; a presentment in a few thousand is mistyped, and a batch in a few dozen files has a trailer a cent out. Dispute events move the cases in flight through their state machine, with an event a case cannot take once in a while, and open new cases on last night's sales. `laurasia/processor` runs it: `make cards-open` and `make cards-days N=30` print a line a day, the night's counts on a business day and the cycle's statements on every day, and `laurasia/run/processor/data` is the processor's `LCS_DATA`.

2. **The payments feed.** Found by running it. Nothing in Laurasia reduces a balance: POST adds purchases and over-limit fees, CYCLE adds interest, and no program credits a payment or sets a card's paid-in-full flag. A real processor receives the cardholders' payments from the card system every day. Over a hundred and fifty days from one seed, with the switch declining as cards fill up:

   | Night | Presentments | Declines in the day's log | Over-limit fees |
   | --- | --: | --: | --: |
   | 2026-10-08 | 588 | 13 | 5.90 |
   | 2026-12-03 | 367 | 246 | 77.56 |
   | 2027-02-26 | 222 | 383 | 246.05 |

   By the end, 670 of the 1,000 cards stand at their limit. The specimen needs a payments file from the card system, `PAYMENTS.DATE.txt`, a program that credits each payment to its card and marks the statement paid in full when it clears it, run before the cycle, and in the same feed the card system's status changes, since nothing blocks, unblocks or closes a card after the processor opens either. That is a change to the specimen: its golden outputs, the operations guide, the blueprints and the oracle's cases move with it, and the survey reading a changed system is a test in itself.

Done when: the processor runs for a year from a seed with its cards paying, declining and disputing at rates that hold steady, and the survey's blueprints of it show the payments feed as they show the rest.
