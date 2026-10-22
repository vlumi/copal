# Operations guide

Laurasia Card Services, Night Operations. Revision of 2022-09. Supersedes the 2021 edition; the card system's feed is new.

## The night

The card networks send one clearing file a business day, with the presentments of the cards we issue and of the merchants we acquire for; the authorization switch leaves the day's log of what it approved and declined; the card system sends the cardholders' payments and the cards' status changes. The night turns these into postings, merchant payouts and dispute adjustments. The date the night works on is the date it is given, as `YYYYMMDD`; no program reads the clock.

The switch is the network's and runs around the clock; nothing we run closes it. When it cannot reach us it stands in, and the log marks those approvals. The log we are given on a business day holds everything since the last night, a weekend's authorizations included.

The scripts are in `scripts/`, and each takes its programs from `$LCS_BIN`, the other scripts from `$LCS_SCRIPTS` and the files from `$LCS_DATA`.

- `open.sh`, once, when the processor opens: loads the card accounts from the card system's extract into the indexed account file.
- `nightly.sh DATE`, every business day.
- `cycle.sh DATE`, after the night, every day: closes the billing cycles of the cards whose cycle day it is.

## The nightly run

1. **CLRIN** reads the clearing file, `inbound/CLR.DATE.txt`: a file header for the night's date, batches each of a header, presentments and a trailer, and a file trailer. It holds a batch's presentments until its trailer, and lets them go on only if the trailer's count and amount agree; otherwise the whole batch is rejected, `BTOT`. A presentment is rejected on its own for a card number that is not numeric, `CARD`, an amount not above zero, `AMNT`, a merchant category code that is not numeric, `MCC`, a region other than domestic or international, `REGN`, or a product other than classic, gold or business, `PROD`. A record out of place is rejected, `ORDR`.
2. **PRESSORT** and **AUTHSORT** sort the presentments and the authorization log, `inbound/AUTH.DATE.txt`, by authorization code.
3. **MATCH** marks each presentment by its authorization: matched, matched but declined, over the approved amount by more than the tolerance of 15 percent, or unmatched. An authorization no presentment matched has lapsed, and is counted.
4. **FEELD** loads the fee schedule, `inbound/FEES.txt`, into the relative fee table, each fee in the slot its product, merchant group and region give.
5. **IXFEE** prices the interchange of each matched presentment from the fee table: a rate and a fixed amount. A presentment with no fee in the table takes the default, 2 percent and 25 cents.
6. **POST** posts each matched presentment to its card's account, `master/accounts.idx`, and writes the postings. A purchase past the credit limit still posts, as it was approved, and may cost an over-limit fee, by the card's product. A card with no account, `CARD`, a blocked one, `BLKD`, or a closed one, `CLSD`, is rejected for the disputes desk.
7. **PAYMENTS** runs `payments.sh`: **PAYIN** applies the card system's feed, `inbound/PAYMENTS.DATE.txt`, to the accounts. A payment comes off the card's balance and is written as a credit; one that covers the last statement, or clears the balance, marks the statement paid in full, which spares the card interest at the cycle. A status change sets the card active, blocked or closed. A record for a card with no account, `CARD`, a payment of nothing, `AMNT`, a status that is none of the three, `STAT`, or a record of no known kind, `KIND`, is rejected for the card system.
8. **DISPUTES** runs `disputes.sh`: **DSPSORT** sorts the day's dispute events, `inbound/DISPEVT.DATE.txt`, by case; **DISPUTE** moves each case on by its events and writes a new dispute master and the merchants' adjustments; **DSPNEW** puts the new master in place of the old. An event the case's state does not allow is rejected, `STAT`.
9. **SETLSORT** and **ADJSORT** sort the priced presentments and the adjustments by merchant.
10. **SETTLE** works out what each merchant is paid: its presentments less their interchange, and its adjustments. A merchant who ends up owing is not paid, and is listed for collection.
11. **ACCTUL** unloads the account file for the card system.

## The billing cycle

**CYCLE** closes the cycle of each open or blocked card whose cycle day is the day of the month: interest on the last statement's balance, unless it was paid in full, which the grace period allows; the minimum payment, a fiftieth of the balance but at least 25.00, or the whole balance below that; and the statement, which starts unpaid until a payment covers it. It writes the interest as postings. **ACCTUL** unloads the accounts after it.

## Where the output is

Under `$LCS_DATA`:

- `logs/`: a line for each step, its name, the program and its return code; `night.DATE.log`, `cycle.DATE.log` and `open.log`.
- `outbound/`: what goes to the other systems: the rejects, `CLRREJ`, `POSTREJ`, `PAYREJ` and `DISPREJ`; the postings, the credits and the interest; the statements; the merchants' payouts; the accounts.
- `reports/`: each program's report, `.txt`, with its counts and totals, and anything it displayed, `.out`.
- `work/DATE/`: the night's intermediate files. Keep them until the night is checked.

## Return codes

| RC | Meaning | What to do |
| --: | --- | --- |
| 0 | The step did its job. | Nothing. |
| 4 | The step did its job but rejected or set aside some records. | Normal. Check the step's report and its rejects. |
| 8 | CLRIN only: the file trailer does not agree with the batches. | Ask the network for the file again. |
| 12 | CLRIN only: the file is not the night's, or has no header. | Ask the network for the file again. |
| 16 | The step could not do its job: a file error, or an account file it cannot update. | A fault. Read the step's `.out`. |

A script stops at the first step that ends with 8 or more, or a sort that fails, and returns its code; nothing after it runs.
