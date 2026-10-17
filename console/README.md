# The console

A web window onto the running First Bank of Pangaea, built for the modernization: it runs the bank, watches its online day, and shows each run's jobs as the survey draws them. It is the one part of Copal that knows both the bank and the survey, so it lives beside them rather than in either.

`make console` serves it at <http://127.0.0.1:8080/> (`ADDR` to change it):

- the bank today: its business date, accounts and balances, and every day it has run, with each day's worst return code;
- accounts, searched by account number, customer number or name, and for each account its master record field by field, every posting and every reject since the reset;
- each business day: every job it ran, the online day, the daily batch and on the first business day of a month the month-end job, each drawn as the survey draws its network, with every step's return code and every dataset's records, with links from steps to the blueprints of their programs, and with its reports; then the transaction file's control totals against its details, the rejects of transactions and of maintenance requests, the audit trail and the transaction file;
- what changed between any two days' openings: balances, status changes and new accounts, reconciled against the postings in the audit trails in between;
- the online day as it happens: every authorization as it is answered, approved or declined, and the balances and available balances of the accounts as they move, followed as server-sent events.

It reads every record through the layouts computed from the copybooks, the same way the bank's generator writes them. Its buttons reset the bank and run business days by calling the same code as `make bank-reset` and `make bank-days`; the console itself never writes the bank's files.

The console also serves the bank's online channel, under `/api/`: it opens the business day's online day as it starts, closes it for the batch when you run days, during which the channel answers that the bank is closed, and opens the next day's after. `make simulate` sends it card traffic, five requests a second by default (`RATE`): purchases, withdrawals, balance enquiries and the odd reversal, for the bank's cards, each with a spending habit of its own, so that some run into their daily limits and some accounts out of money. Interrupting the console closes the online day and extracts its holds.
