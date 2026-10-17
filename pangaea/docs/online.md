# The authorization service

First Bank of Pangaea, Online Systems. Revision of 2026-09.

The authorization service answers the card terminals and the ATMs during the online day. It is `AUTHSRV`, running as the `AUTH` step of the `ONLINE` job from opening until the batch window, reading one request at a time and writing one response. The programs that decide, `AUTHPUR`, `AUTHREV`, `ATMWDL` and `ATMBAL`, are called for each request; only `AUTHSRV` touches the files.

## Requests

| Function | From | Does |
| --- | --- | --- |
| PURC | Card terminal | A purchase: holds the amount on the account. |
| REVS | Card terminal | Reverses a purchase of the same card by its authorization code; releases the hold. Only purchases are reversed, and only once. |
| WDRL | ATM | A cash withdrawal: holds the amount and any fee. |
| BALQ | ATM | A balance enquiry: the balance and the available balance. Changes nothing. |

Over HTTP the terminals send JSON to `/api/purchase`, `/api/reversal`, `/api/withdrawal` and `/api/balance`, for instance `{"card": "9100000000000018", "amount": "12.30", "merchant": "CAFE"}`, and get the response as JSON. `GET /api/status` says whether the service is open.

## The rules, in plain terms

- A card must be active. A blocked card is declined; an expired card is not even loaded for the day, so it is unknown.
- The account must be open. A frozen account is declined.
- **Available balance** is the balance the day opened with, plus the overdraft limit, less what is already held today.
- A purchase must fit the card's **daily purchase limit**, counting the day's earlier purchases, and the available balance.
- A withdrawal must be in **notes of 20**, fit the card's **daily cash limit** counting the day's earlier withdrawals, and fit the available balance together with its fee.
- The first **four withdrawals in a calendar month** are free; each one after that carries the **ATM fee of 2.00**, held with the amount.
- A card's totals are for the business date they carry. The first authorization of a new day starts them over, and a new month starts the withdrawal count over.

## Result codes

| Code | Meaning |
| --: | --- |
| 00 | Approved. The response carries the authorization code, the amount, the fee, and the account's balance and available balance after it. |
| 12 | Invalid request: an unknown function, or a reversal of something that is not a purchase. |
| 13 | Invalid amount: zero, or a withdrawal not in notes of 20. |
| 14 | No such card. |
| 25 | No such authorization: a reversal names an authorization code the card does not have today, or one already reversed. |
| 51 | Insufficient funds: the amount, with any fee, is more than the available balance. |
| 61 | Over the daily purchase or cash limit. |
| 62 | Card blocked. |
| 78 | Account not open: frozen, or no longer in the online file. |
| 91 | Bank closed: the batch window. Answered by the server, not by `AUTHSRV`. |

## The day's end

When the service closes, `HOLDEXT` extracts the holds still open as the card file: one transaction per purchase or withdrawal, and one per ATM fee, with the control totals the batch checks. Reversed purchases are left out. The batch posts the card file with the branches' transactions, whatever the balance by then: the money was authorized when it was spent. The next morning's online files are loaded from the new master, with no holds.

The service can be stopped and started during the day. It continues numbering the day's authorizations from the highest hold on file. The day's authorizations are counted in `AUTH.AUTHRPT`, and the extraction in `EXTRACT.HOLDRPT`, under `sysout/DATE/ONLINE/`.
