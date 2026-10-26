# Operations guide

Gondwana Life Insurance, Systems Department. Revision of 2026-10. Supersedes the 2022 edition; the anniversaries since the last business date are new.

## The runs

The system holds the policies in an indexed master, `MASTER\POLMST.IDX`, keyed by policy number, and a control record, `MASTER\GLSCTL.BIN`, with the business date it stands at and the next payment and journal numbers. Dates in the master are in the era, `GYYMMDD`: `S` for Showa, `H` for Heisei, `R` for Reiwa, then the year of the era, so `R090301` is March 1, 2027. The paid-to month is the exception, `YYYYMM` in the Western calendar, as the collection agent's files have it. Amounts are whole yen. The business date is the date the batch file is given, as `YYYYMMDD`; no program reads the clock.

The batch files are in `scripts\`, and the scheduler starts them with the programs from `%GLS_BIN%`, the other batch files from `%GLS_SCRIPTS%` and the files from `%GLS_DATA%`:

- `OPEN.BAT DATE`, once, when the system opens: loads the policy extract into the master and writes the control record.
- `DAILY.BAT DATE`, every business day.
- `BILLING.BAT DATE`, once a month before the debit day, the 27th: the transfer request to the collection agent.
- `RESULTS.BAT DATE`, on the day the agent's results arrive, a few business days after the debit, before the daily run.
- `CLOSE.BAT DATE`, on the last business day of the month, after the daily run.
- `ANNUAL.BAT DATE`, on March 31, after the close: the fiscal year ends.

## The daily run

1. **GLD010** reads the sales system's new policies, `IN\NEWPOL.DATE.TXT`, whose dates are Western, and registers them: the product must be sold, the policy number new, the payment mode monthly or annual, the contract and birth dates dates of an era, the sex 1 or 2, the insurance age at the contract date within the product's range, the sum insured within its range, and a rate on the table. The insurance age is the age in full years plus one when six months or more of the next have passed. The premium is the rate per 1,000 yen of the sum, monthly, truncated to the yen; an annual premium is twelve months less five percent. The first premium was collected with the application, so the policy is paid to its contract month and the premium is booked. A policy that fails goes to `OUT\NEWREJ.DATE.TXT` with a reason code.
2. **GLD020** moves the policy year on for every covered policy whose anniversary, the contract's day of the year, has come since the control record's business date, the last the unload moved it to, so that a weekend's anniversaries are taken on the Monday; on the day the system opens, and on a rerun, the day's own. A policy contracted that day has no anniversary yet. When the anniversary is an endowment's maturity date it pays the sum insured plus the accumulated dividends and marks the policy matured; when it is a term policy's it renews the policy at the insurance age attained on the anniversary, with a new premium and ten more years from it, or ends it without payment when that age is past the product's maximum.
3. **GLD030** pays the branches' benefit requests, `IN\CLMREQ.DATE.TXT`, whose event dates are in the era: a death benefit, the sum plus dividends less any policy loan, after which the policy is paid out; a hospitalization benefit, for medical policies only, a daily benefit of a thousandth of the sum for up to 120 days, not within the product's waiting days after the contract date; or a surrender, the cash value per 1,000 yen by product, entry age band and policy year, less the loan, which a lapsed policy may still take. A request for a policy not covered, or whose event is before the contract, is rejected to `OUT\CLMREJ.DATE.TXT`.
4. **GLD040** applies the branches' changes, `IN\MNTREQ.DATE.TXT`: an address, a bank account, a name with its kana, a beneficiary, or a reinstatement of a policy lapsed within the year, which takes the arrears paid at the counter, premiums from the month after the paid-to month to this one, and puts the policy back in force. Rejects go to `OUT\MNTREJ.DATE.TXT`.
5. **GLD060** lapses the policies in grace whose grace has run out: the grace runs from the debit that failed to the end of the month after.
6. **UNLOAD** runs `UNLOAD.BAT`: **GLD090** unloads the master to `OUT\POLMST.DATE.TXT` for the sales system and the accounts, and moves the control record's business date on.

Every program that books anything appends to the day's journal, `OUT\GLJRNL.DATE.TXT`, and every payment to `OUT\PAYOUT.DATE.TXT`, numbered from the control record.

## Billing, results and the close

**GLM010** selects the policies to debit on the 27th: monthly payers paid to before the billing month, and annual payers whose contract month it is; a policy in arrears is debited for two months at most. It writes the request in the banks' fixed-length layout, `OUT\PRMREQ.DATE.DAT`, 120 bytes a record with no delimiters, a header, the debits sorted by bank, branch and account, a trailer with the counts and the amount, and an end record, and notes the billing month in the control record so that a month is not billed twice.

**GLD050** reads the agent's results, `IN\PRMRES.DATE.DAT`, in the same layout with the result filled in. It holds the debits until the trailer and applies nothing if the counts or amounts disagree. A debit that went through moves the paid-to month on, one month or two, or a year, clears the arrears, and puts a policy in grace back in force; one that failed, for lack of funds, a closed account or a refusal, counts an arrear and puts a policy in force into grace from that day. Results for policies no longer paying, matured, paid out or surrendered since the billing, are rejected to `OUT\RESREJ.DATE.TXT`.

**GLM020** works out each covered policy's reserve at the month's end, by product and policy year per 1,000 yen of the sum, books the change since last month, releases the reserve of a policy no longer covered, and prints the reserves by product. **GLM030** works out the agents' commissions on the premiums collected in the month, forty percent in the first policy year, seven in the second and third, two to the tenth, nothing after, prints them by agent with the agent's branch and the month the agent joined, and writes each agent's total to `OUT\COMMIS.DATE.TXT` for payroll. The close ends with the unload.

**GLY010**, at the fiscal year's end, credits a dividend to every covered policy in its third year or later whose product pays one, a rate on the sum, into the policy's accumulated dividends, which the death and maturity benefits pay out.

## Where the output is

Under `%GLS_DATA%`:

- `LOG\`: a line for each step, its name, the program and its errorlevel; `DAILY.DATE.LOG`, `BILLING.DATE.LOG`, `RESULTS.DATE.LOG`, `CLOSE.DATE.LOG`, `ANNUAL.DATE.LOG`, `UNLOAD.DATE.LOG` and `OPEN.LOG`.
- `OUT\`: what goes to the other systems: the rejects, `NEWREJ`, `CLMREJ`, `MNTREJ` and `RESREJ`; the journal; the payments; the transfer request; the commissions; the unloaded master.
- `RPT\`: each program's report, `.RPT`, in Japanese, with its counts and amounts, and anything it displayed, `.OUT`.

## The console

The console, `make life-console`, shows the same without opening the files: every business day run with its batch files, counts and worst errorlevel, each batch file step by step, the files each day read and wrote, and any policy's master record. Its Run button runs the next business days as the scheduler would.

## Return codes

| Code | Meaning | What to do |
| --: | --- | --- |
| 0 | The step did its job. | Nothing. |
| 4 | The step did its job but rejected some records. | Normal. Check the step's report and its rejects. |
| 8 | The step did nothing: the date was not a date, the month was billed already, or the results file does not agree with its trailer. | Check the date given, or ask the agent for the file again. |
| 16 | The step could not do its job: a file error, or a master it cannot update. | A fault. Read the step's `.OUT`. |

A batch file stops at the first step that ends with 8 or more and exits with its code; nothing after it runs.
