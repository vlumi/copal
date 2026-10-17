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

WHEN A MONTH CLOSES, THE BANK CREDITS ITS INTEREST (IN),
CHARGES OVERDRAFT INTEREST (OD) AND THE MAINTENANCE FEE (MF).
THESE TOO ARE NEVER RECEIVED IN THE DAILY FILE.

THE CARD FILE, FROM THE ONLINE DAY'S HOLDS, HAS THE SAME FORM
AND ITS OWN SOURCE: CARD PURCHASES (CP), ATM WITHDRAWALS (AW)
AND ATM FEES (AF), WHICH WERE AUTHORIZED WHEN THEY WERE MADE.

MODIFICATION HISTORY
DATE       BY  CHANGE
1987-03-12 JA  ORIGINAL VERSION.
1991-06-03 PK  CR-0117 ACCOUNT NUMBER HASH IN THE TRAILER.
1996-04-15 PK  CR-0233 TRANSFERS: COUNTER ACCOUNT, FROM THE
               FILLER.
2001-09-03 TL  CR-0298 MONTH-END INTEREST AND FEE CODES.
2003-05-19 MV  CR-0342 THE CARD FILE AND ITS CODES.
```

Copybook [TRANREC.cpy](../../../pangaea/copybooks/TRANREC.cpy), used by [HOLDEXT](../programs/HOLDEXT.md), [MONCHG](../programs/MONCHG.md), [TXPOST](../programs/TXPOST.md), [TXVALID](../programs/TXVALID.md).

## TRAN-RECORD

80 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | TRAN-RECORD | group | 1 | 80 | | | | [TRANREC.cpy:31](../../../pangaea/copybooks/TRANREC.cpy#L31) |
| 05 | TR-REC-TYPE | code | 1 | 1 | X | | | [TRANREC.cpy:32](../../../pangaea/copybooks/TRANREC.cpy#L32) |
| 88 | TR-HEADER | | | | | | VALUE 'H' | [TRANREC.cpy:33](../../../pangaea/copybooks/TRANREC.cpy#L33) |
| 88 | TR-DETAIL | | | | | | VALUE 'D' | [TRANREC.cpy:34](../../../pangaea/copybooks/TRANREC.cpy#L34) |
| 88 | TR-TRAILER | | | | | | VALUE 'T' | [TRANREC.cpy:35](../../../pangaea/copybooks/TRANREC.cpy#L35) |
| 05 | TR-BODY | text | 2 | 79 | X(79) | | | [TRANREC.cpy:36](../../../pangaea/copybooks/TRANREC.cpy#L36) |
| 05 | TR-HEADER-BODY | group | 2 | 79 | | | REDEFINES TR-BODY | [TRANREC.cpy:37](../../../pangaea/copybooks/TRANREC.cpy#L37) |
| 10 | TH-BUSINESS-DATE | date | 2 | 8 | 9(8) | | | [TRANREC.cpy:38](../../../pangaea/copybooks/TRANREC.cpy#L38) |
| 10 | TH-SOURCE | code | 10 | 8 | X(8) | | | [TRANREC.cpy:39](../../../pangaea/copybooks/TRANREC.cpy#L39) |
| 88 | TH-FROM-BRANCHES | | | | | | VALUE 'BRANCHES' | [TRANREC.cpy:40](../../../pangaea/copybooks/TRANREC.cpy#L40) |
| 88 | TH-FROM-CARDS | | | | | | VALUE 'CARDS' | [TRANREC.cpy:41](../../../pangaea/copybooks/TRANREC.cpy#L41) |
| 10 | TH-FILE-SEQ | number | 18 | 6 | 9(6) | | | [TRANREC.cpy:42](../../../pangaea/copybooks/TRANREC.cpy#L42) |
| 10 | FILLER | filler | 24 | 57 | X(57) | | | [TRANREC.cpy:43](../../../pangaea/copybooks/TRANREC.cpy#L43) |
| 05 | TR-DETAIL-BODY | group | 2 | 79 | | | REDEFINES TR-BODY | [TRANREC.cpy:44](../../../pangaea/copybooks/TRANREC.cpy#L44) |
| 10 | TD-TRAN-ID | number | 2 | 10 | 9(10) | | | [TRANREC.cpy:45](../../../pangaea/copybooks/TRANREC.cpy#L45) |
| 10 | TD-ACCT-NO | number | 12 | 10 | 9(10) | | | [TRANREC.cpy:46](../../../pangaea/copybooks/TRANREC.cpy#L46) |
| 10 | TD-TRAN-CODE | code | 22 | 2 | XX | | | [TRANREC.cpy:47](../../../pangaea/copybooks/TRANREC.cpy#L47) |
| 88 | TD-DEPOSIT | | | | | | VALUE 'DP' | [TRANREC.cpy:48](../../../pangaea/copybooks/TRANREC.cpy#L48) |
| 88 | TD-WITHDRAWAL | | | | | | VALUE 'WD' | [TRANREC.cpy:49](../../../pangaea/copybooks/TRANREC.cpy#L49) |
| 88 | TD-FEE | | | | | | VALUE 'FE' | [TRANREC.cpy:50](../../../pangaea/copybooks/TRANREC.cpy#L50) |
| 88 | TD-TRANSFER | | | | | | VALUE 'TR' | [TRANREC.cpy:51](../../../pangaea/copybooks/TRANREC.cpy#L51) |
| 88 | TD-TRANSFER-CREDIT | | | | | | VALUE 'TC' | [TRANREC.cpy:52](../../../pangaea/copybooks/TRANREC.cpy#L52) |
| 88 | TD-TRANSFER-RETURN | | | | | | VALUE 'RT' | [TRANREC.cpy:53](../../../pangaea/copybooks/TRANREC.cpy#L53) |
| 88 | TD-INTEREST | | | | | | VALUE 'IN' | [TRANREC.cpy:54](../../../pangaea/copybooks/TRANREC.cpy#L54) |
| 88 | TD-OVERDRAFT-INTEREST | | | | | | VALUE 'OD' | [TRANREC.cpy:55](../../../pangaea/copybooks/TRANREC.cpy#L55) |
| 88 | TD-MAINTENANCE-FEE | | | | | | VALUE 'MF' | [TRANREC.cpy:56](../../../pangaea/copybooks/TRANREC.cpy#L56) |
| 88 | TD-CARD-PURCHASE | | | | | | VALUE 'CP' | [TRANREC.cpy:57](../../../pangaea/copybooks/TRANREC.cpy#L57) |
| 88 | TD-ATM-WITHDRAWAL | | | | | | VALUE 'AW' | [TRANREC.cpy:58](../../../pangaea/copybooks/TRANREC.cpy#L58) |
| 88 | TD-ATM-FEE | | | | | | VALUE 'AF' | [TRANREC.cpy:59](../../../pangaea/copybooks/TRANREC.cpy#L59) |
| 88 | TD-KNOWN-CODE | | | | | | VALUE 'DP', 'WD', 'FE', 'TR' | [TRANREC.cpy:60](../../../pangaea/copybooks/TRANREC.cpy#L60) |
| 88 | TD-CARD-CODE | | | | | | VALUE 'CP', 'AW', 'AF' | [TRANREC.cpy:61](../../../pangaea/copybooks/TRANREC.cpy#L61) |
| 10 | TD-AMOUNT | amount | 24 | 11 | 9(9)V99 | | | [TRANREC.cpy:62](../../../pangaea/copybooks/TRANREC.cpy#L62) |
| 10 | TD-DESCRIPTION | text | 35 | 20 | X(20) | | | [TRANREC.cpy:63](../../../pangaea/copybooks/TRANREC.cpy#L63) |
| 10 | TD-COUNTER-ACCT-NO | number | 55 | 10 | 9(10) | | | [TRANREC.cpy:64](../../../pangaea/copybooks/TRANREC.cpy#L64) |
| 10 | FILLER | filler | 65 | 16 | X(16) | | | [TRANREC.cpy:65](../../../pangaea/copybooks/TRANREC.cpy#L65) |
| 05 | TR-TRAILER-BODY | group | 2 | 79 | | | REDEFINES TR-BODY | [TRANREC.cpy:66](../../../pangaea/copybooks/TRANREC.cpy#L66) |
| 10 | TT-DETAIL-COUNT | number | 2 | 8 | 9(8) | | | [TRANREC.cpy:67](../../../pangaea/copybooks/TRANREC.cpy#L67) |
| 10 | TT-AMOUNT-TOTAL | amount | 10 | 15 | 9(13)V99 | | | [TRANREC.cpy:68](../../../pangaea/copybooks/TRANREC.cpy#L68) |
| 10 | TT-ACCT-HASH | number | 25 | 15 | 9(15) | | | [TRANREC.cpy:69](../../../pangaea/copybooks/TRANREC.cpy#L69) |
| 10 | FILLER | filler | 40 | 41 | X(41) | | | [TRANREC.cpy:70](../../../pangaea/copybooks/TRANREC.cpy#L70) |
