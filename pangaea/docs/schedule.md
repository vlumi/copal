# Job schedule

First Bank of Pangaea, Data Center. Revision of 2026-09.

The jobs, in the order they run on a business day, with their steps. `COND=(4,LT)` means the step runs only if no earlier step ended with a return code above 4. The JCL in [jobs/](../jobs/) is the authority; this is the reading of it.

## START, once, when the bank opens

| Step | Program | Does |
| --- | --- | --- |
| LOAD | ACCTLOAD | Loads the master into the online account file for the first business day, and starts an empty hold file. |
| CARDS | CARDLOAD | Loads the card master into the online card file for that day. |

## ONLINE, the online day

| Step | Program | Does |
| --- | --- | --- |
| AUTH | AUTHSRV | Answers authorization requests, one at a time, from opening until the batch window; updates the online account, card and hold files. |
| EXTRACT | HOLDEXT | At close, writes the holds still open as the card file, `TRAN.CARD(+1)`, for the batch to post. |

## DAILY, the nightly batch

| Step | Program | Does |
| --- | --- | --- |
| VALID | TXVALID | Validates the branches' transaction file, `TRAN.DAILY(0)`: header, details, trailer. Writes the rejects. |
| VALIDC | TXVALID | Validates the card file, `TRAN.CARD(0)`, the same way; accepts only card and ATM codes from it. |
| ACCRUE | INTACCR | Accrues interest up to today on every account; when a month has closed, makes its interest due. |
| CHARGE | MONCHG | On the day a month closes, writes its interest credits, overdraft interest and maintenance fees as transactions; otherwise nothing. |
| POSTMO | TXPOST | Posts the month-end charges, which come in account order. |
| SORT | TXSORT | Sorts the valid transactions of both files into account order. |
| POST | TXPOST | Posts them. A posted transfer writes its credit for the receiving account. |
| SORTCR | TXSORT | Sorts the transfer credits. |
| POSTCR | TXPOST | Posts them. A credit to an unknown or closed account writes a return to its sender. |
| SORTRT | TXSORT | Sorts the returns. |
| POSTRT | TXPOST | Posts them. |
| MAINT | ACCTMNT | Applies the branches' account maintenance requests and writes the day's new master, `ACCT.MASTER(+1)`. |
| LOAD | ACCTLOAD | Loads the new master into the online account file for the next business day, and starts its hold file. |
| NEXTDAY | BANKDATE | Writes the bank control record for the next business day. |
| CARDS | CARDLOAD | Loads the card master into the online card file for that day, leaving out the cards expired by then. |

Every step from VALIDC on is `COND=(4,LT)`. The four TXPOST steps append to one audit trail, `TRAN.AUDIT(+1)`, and one rejects file, `TRAN.REJECTS(+1)`, which VALID and VALIDC started.

## MONTHLY, after DAILY on the first business day of a month

| Step | Program | Does |
| --- | --- | --- |
| STMT | STMTPRT | Prints a statement for every account for the month that closed, from every audit trail. |
| LEDGER | GLSUM | Summarizes the month for the general ledger, `GL.SUMMARY(+1)`, in double entry. `COND=(4,LT)`. |
