# Dataset list

First Bank of Pangaea, Data Center. Revision of 2026-09.

Every dataset the bank keeps, by name. A name with a generation, `NAME(+1)`, is a generation data group: every version is kept, `NAME.G0007V00` being the seventh, and a job names them relative to the newest when it started. `(0)` is the current one, `(+1)` the one the job writes, `(-1)` the one before. A dataset named by its base alone, `TRAN.AUDIT`, is read as every generation it has, newest first. `&&NAME` is a temporary dataset that exists only while its job runs; those are not listed here. Record layouts are in [copybooks/](../copybooks/).

| Dataset | Records | Written by | Read by | Notes |
| --- | --- | --- | --- | --- |
| BANK.CONTROL | BANKCTL, 80 bytes, one record | DAILY NEXTDAY | Every job | The business date. A new generation every day. |
| BANK.RATES | RATEREC, 20 bytes | Nobody since the bank opened | DAILY ACCRUE | The interest rate sheet, by account type and balance tier. Maintained by hand. |
| ACCT.MASTER | ACCTMAST, 100 bytes, in account order | DAILY MAINT | DAILY, MONTHLY, START | The account master. A new generation every day; the day's postings and maintenance applied. |
| ACCT.ACCRUAL | ACCRREC, 40 bytes, in step with the master | DAILY ACCRUE | DAILY, MONTHLY | Each account's accrued interest, and the interest due when a month closed. Packed decimal. |
| ACCT.MAINT | MAINTREC, 80 bytes | The branches | DAILY MAINT | The day's account maintenance requests: open, close, freeze, unfreeze, overdraft limit. |
| TRAN.DAILY | TRANREC, 80 bytes, header, details, trailer | The branches | DAILY VALID | The day's transactions: deposits, withdrawals, fees and transfers. From outside. |
| TRAN.CARD | TRANREC, 80 bytes, header, details, trailer | ONLINE EXTRACT | DAILY VALIDC | The online day's card purchases, ATM withdrawals and ATM fees, from the holds. |
| TRAN.AUDIT | AUDITREC, 100 bytes | DAILY, the four TXPOST steps | MONTHLY | One record per posting: account, transaction, amount, balance before and after. A generation per day, kept forever. |
| TRAN.REJECTS | REJREC, 110 bytes | DAILY VALID, VALIDC and the TXPOST steps | Nobody; returned to the branches | The transactions that did not validate or could not post, with the reason. |
| MAINT.REJECTS | REJREC, 110 bytes | DAILY MAINT | Nobody; returned to the branches | The maintenance requests the rules refused. |
| CARD.MASTER | CARDREC, 80 bytes, in card order | Nobody since the bank opened | DAILY CARDS, START | The cards: account, status, expiry, daily limits. Maintained by hand. |
| ACCT.ONLINE | ACCTONL, 100 bytes, indexed by account | DAILY LOAD, START; updated by ONLINE AUTH | ONLINE | The online account file for one business day: balance, overdraft limit, holds. Indexed: no text form. |
| CARD.ONLINE | CARDREC, 80 bytes, indexed by card | DAILY CARDS, START; updated by ONLINE AUTH | ONLINE | The online card file for one business day, with the day's running totals. Indexed. |
| HOLD.ONLINE | HOLDREC, 100 bytes, indexed by account and hold number | DAILY LOAD, START, empty; written by ONLINE AUTH | ONLINE EXTRACT | The day's holds: what the authorizations approved. Indexed. |
| GL.SUMMARY | GLREC, 60 bytes | MONTHLY LEDGER | Nobody; for the general ledger | Each ledger account's debits and credits for the closed month. Packed decimal. |

## Reject reasons

The rejects carry the program that rejected the record and a four-letter reason.

| Reason | From | Meaning |
| --- | --- | --- |
| ACCT | TXVALID | The account number is not numeric. |
| CODE | TXVALID | The transaction code is not one the file's source sends. |
| AMNT | TXVALID | The amount is not numeric, or zero. |
| CTRA | TXVALID | A transfer names no counter account, or itself. |
| TYPE | TXVALID, ACCTMNT | The record is neither header, detail nor trailer, or the maintenance request is of no known type. |
| UNKN | TXPOST, ACCTMNT | No such account. |
| CLSD | TXPOST, ACCTMNT | The account is closed. |
| FRZN | TXPOST | The account is frozen, and the transaction is a withdrawal or transfer. |
| FUND | TXPOST | The debit would take the balance beyond the overdraft limit. |
| EXST | ACCTMNT | An account to open already exists. |
| BALN | ACCTMNT | An account to close still has a balance. |
| STAT | ACCTMNT | The request does not fit the account's status: freezing a frozen account, unfreezing an open one. |
