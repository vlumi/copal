# TRANREC

```text
TRANREC: DAILY TRANSACTION FILE RECORD, 80 BYTES
ONE HEADER, THE DETAILS, ONE TRAILER, TOLD APART BY THE
RECORD TYPE. THE BODY IS REDEFINED FOR EACH RECORD TYPE.
THE TRAILER CARRIES THE SENDER'S CONTROL TOTALS: THE DETAIL
COUNT, THE SUM OF AMOUNTS AND THE SUM OF ACCOUNT NUMBERS.

A TRANSFER (TR) NAMES THE ACCOUNT IT GOES TO. POSTING SPLITS
IT IN TWO: THE DEBIT ON THE SENDING ACCOUNT, THEN A TRANSFER
CREDIT (TC) ON THE RECEIVING ONE. A CREDIT THAT CANNOT POST
GOES BACK AS A RETURNED TRANSFER (RT). TC AND RT ARE MADE BY
POSTING, NEVER RECEIVED IN THE DAILY FILE.

MODIFICATION HISTORY
DATE       BY  CHANGE
1987-03-12 JA  ORIGINAL VERSION.
1991-06-03 PK  CR-0117 ACCOUNT NUMBER HASH IN THE TRAILER.
1996-04-15 PK  CR-0233 TRANSFERS: COUNTER ACCOUNT, FROM THE
               FILLER.
```

Copybook [TRANREC.cpy](../../../pangaea/copybooks/TRANREC.cpy), used by [TXPOST](../programs/TXPOST.md), [TXVALID](../programs/TXVALID.md).

## TRAN-RECORD

80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | TRAN-RECORD | 1 | 80 | | | | [TRANREC.cpy:21](../../../pangaea/copybooks/TRANREC.cpy#L21) |
| 05 | TR-REC-TYPE | 1 | 1 | X | | | [TRANREC.cpy:22](../../../pangaea/copybooks/TRANREC.cpy#L22) |
| 88 | TR-HEADER | | | | | VALUE 'H' | [TRANREC.cpy:23](../../../pangaea/copybooks/TRANREC.cpy#L23) |
| 88 | TR-DETAIL | | | | | VALUE 'D' | [TRANREC.cpy:24](../../../pangaea/copybooks/TRANREC.cpy#L24) |
| 88 | TR-TRAILER | | | | | VALUE 'T' | [TRANREC.cpy:25](../../../pangaea/copybooks/TRANREC.cpy#L25) |
| 05 | TR-BODY | 2 | 79 | X(79) | | | [TRANREC.cpy:26](../../../pangaea/copybooks/TRANREC.cpy#L26) |
| 05 | TR-HEADER-BODY | 2 | 79 | | | REDEFINES TR-BODY | [TRANREC.cpy:27](../../../pangaea/copybooks/TRANREC.cpy#L27) |
| 10 | TH-BUSINESS-DATE | 2 | 8 | 9(8) | | | [TRANREC.cpy:28](../../../pangaea/copybooks/TRANREC.cpy#L28) |
| 10 | TH-SOURCE | 10 | 8 | X(8) | | | [TRANREC.cpy:29](../../../pangaea/copybooks/TRANREC.cpy#L29) |
| 10 | TH-FILE-SEQ | 18 | 6 | 9(6) | | | [TRANREC.cpy:30](../../../pangaea/copybooks/TRANREC.cpy#L30) |
| 10 | FILLER | 24 | 57 | X(57) | | | [TRANREC.cpy:31](../../../pangaea/copybooks/TRANREC.cpy#L31) |
| 05 | TR-DETAIL-BODY | 2 | 79 | | | REDEFINES TR-BODY | [TRANREC.cpy:32](../../../pangaea/copybooks/TRANREC.cpy#L32) |
| 10 | TD-TRAN-ID | 2 | 10 | 9(10) | | | [TRANREC.cpy:33](../../../pangaea/copybooks/TRANREC.cpy#L33) |
| 10 | TD-ACCT-NO | 12 | 10 | 9(10) | | | [TRANREC.cpy:34](../../../pangaea/copybooks/TRANREC.cpy#L34) |
| 10 | TD-TRAN-CODE | 22 | 2 | XX | | | [TRANREC.cpy:35](../../../pangaea/copybooks/TRANREC.cpy#L35) |
| 88 | TD-DEPOSIT | | | | | VALUE 'DP' | [TRANREC.cpy:36](../../../pangaea/copybooks/TRANREC.cpy#L36) |
| 88 | TD-WITHDRAWAL | | | | | VALUE 'WD' | [TRANREC.cpy:37](../../../pangaea/copybooks/TRANREC.cpy#L37) |
| 88 | TD-FEE | | | | | VALUE 'FE' | [TRANREC.cpy:38](../../../pangaea/copybooks/TRANREC.cpy#L38) |
| 88 | TD-TRANSFER | | | | | VALUE 'TR' | [TRANREC.cpy:39](../../../pangaea/copybooks/TRANREC.cpy#L39) |
| 88 | TD-TRANSFER-CREDIT | | | | | VALUE 'TC' | [TRANREC.cpy:40](../../../pangaea/copybooks/TRANREC.cpy#L40) |
| 88 | TD-TRANSFER-RETURN | | | | | VALUE 'RT' | [TRANREC.cpy:41](../../../pangaea/copybooks/TRANREC.cpy#L41) |
| 88 | TD-KNOWN-CODE | | | | | VALUE 'DP', 'WD', 'FE', 'TR' | [TRANREC.cpy:42](../../../pangaea/copybooks/TRANREC.cpy#L42) |
| 10 | TD-AMOUNT | 24 | 11 | 9(9)V99 | | | [TRANREC.cpy:43](../../../pangaea/copybooks/TRANREC.cpy#L43) |
| 10 | TD-DESCRIPTION | 35 | 20 | X(20) | | | [TRANREC.cpy:44](../../../pangaea/copybooks/TRANREC.cpy#L44) |
| 10 | TD-COUNTER-ACCT-NO | 55 | 10 | 9(10) | | | [TRANREC.cpy:45](../../../pangaea/copybooks/TRANREC.cpy#L45) |
| 10 | FILLER | 65 | 16 | X(16) | | | [TRANREC.cpy:46](../../../pangaea/copybooks/TRANREC.cpy#L46) |
| 05 | TR-TRAILER-BODY | 2 | 79 | | | REDEFINES TR-BODY | [TRANREC.cpy:47](../../../pangaea/copybooks/TRANREC.cpy#L47) |
| 10 | TT-DETAIL-COUNT | 2 | 8 | 9(8) | | | [TRANREC.cpy:48](../../../pangaea/copybooks/TRANREC.cpy#L48) |
| 10 | TT-AMOUNT-TOTAL | 10 | 15 | 9(13)V99 | | | [TRANREC.cpy:49](../../../pangaea/copybooks/TRANREC.cpy#L49) |
| 10 | TT-ACCT-HASH | 25 | 15 | 9(15) | | | [TRANREC.cpy:50](../../../pangaea/copybooks/TRANREC.cpy#L50) |
| 10 | FILLER | 40 | 41 | X(41) | | | [TRANREC.cpy:51](../../../pangaea/copybooks/TRANREC.cpy#L51) |
