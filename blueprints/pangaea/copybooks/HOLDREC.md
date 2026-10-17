# HOLDREC

```text
HOLDREC: HOLD RECORD, 100 BYTES
AN AMOUNT AUTHORIZED ONLINE AND NOT YET POSTED: A CARD
PURCHASE OR AN ATM WITHDRAWAL AND ITS FEE. INDEXED BY ACCOUNT
AND HOLD NUMBER, THE AUTHORIZATION CODE, WHICH NUMBERS THE
DAY'S HOLDS FROM 1. A REVERSED PURCHASE STAYS, MARKED.

MODIFICATION HISTORY
DATE       BY  CHANGE
2003-05-19 MV  CR-0342 ORIGINAL VERSION.
2007-02-05 MV  CR-0415 ATM FEE.
```

Copybook [HOLDREC.cpy](../../../pangaea/copybooks/HOLDREC.cpy), used by [ACCTLOAD](../programs/ACCTLOAD.md), [ATMBAL](../programs/ATMBAL.md), [ATMWDL](../programs/ATMWDL.md), [AUTHPUR](../programs/AUTHPUR.md), [AUTHREV](../programs/AUTHREV.md), [AUTHSRV](../programs/AUTHSRV.md), [HOLDEXT](../programs/HOLDEXT.md).

## HOLD-RECORD

100 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | HOLD-RECORD | group | 1 | 100 | | | | [HOLDREC.cpy:13](../../../pangaea/copybooks/HOLDREC.cpy#L13) |
| 05 | ├ HD-KEY | group | 1 | 16 | | | | [HOLDREC.cpy:14](../../../pangaea/copybooks/HOLDREC.cpy#L14) |
| 10 | │ ├ HD-ACCT-NO | number | 1 | 10 | 9(10) | | | [HOLDREC.cpy:15](../../../pangaea/copybooks/HOLDREC.cpy#L15) |
| 10 | │ └ HD-HOLD-NO | number | 11 | 6 | 9(6) | | | [HOLDREC.cpy:16](../../../pangaea/copybooks/HOLDREC.cpy#L16) |
| 05 | ├ HD-CARD-NO | number | 17 | 16 | 9(16) | | | [HOLDREC.cpy:17](../../../pangaea/copybooks/HOLDREC.cpy#L17) |
| 05 | ├ HD-KIND | code | 33 | 1 | X | | | [HOLDREC.cpy:18](../../../pangaea/copybooks/HOLDREC.cpy#L18) |
| 88 | │ ├ HD-PURCHASE | | | | | | VALUE 'P' | [HOLDREC.cpy:19](../../../pangaea/copybooks/HOLDREC.cpy#L19) |
| 88 | │ └ HD-WITHDRAWAL | | | | | | VALUE 'W' | [HOLDREC.cpy:20](../../../pangaea/copybooks/HOLDREC.cpy#L20) |
| 05 | ├ HD-STATUS | code | 34 | 1 | X | | | [HOLDREC.cpy:21](../../../pangaea/copybooks/HOLDREC.cpy#L21) |
| 88 | │ ├ HD-OPEN | | | | | | VALUE 'O' | [HOLDREC.cpy:22](../../../pangaea/copybooks/HOLDREC.cpy#L22) |
| 88 | │ └ HD-REVERSED | | | | | | VALUE 'R' | [HOLDREC.cpy:23](../../../pangaea/copybooks/HOLDREC.cpy#L23) |
| 05 | ├ HD-AMOUNT | amount | 35 | 9 | 9(7)V99 | | | [HOLDREC.cpy:24](../../../pangaea/copybooks/HOLDREC.cpy#L24) |
| 05 | ├ HD-FEE | amount | 44 | 5 | 9(3)V99 | | | [HOLDREC.cpy:25](../../../pangaea/copybooks/HOLDREC.cpy#L25) |
| 05 | ├ HD-TERMINAL | text | 49 | 20 | X(20) | | | [HOLDREC.cpy:26](../../../pangaea/copybooks/HOLDREC.cpy#L26) |
| 05 | ├ HD-TIMESTAMP | date | 69 | 14 | 9(14) | | | [HOLDREC.cpy:27](../../../pangaea/copybooks/HOLDREC.cpy#L27) |
| 05 | ├ HD-BUSINESS-DATE | date | 83 | 8 | 9(8) | | | [HOLDREC.cpy:28](../../../pangaea/copybooks/HOLDREC.cpy#L28) |
| 05 | └ FILLER | filler | 91 | 10 | X(10) | | | [HOLDREC.cpy:29](../../../pangaea/copybooks/HOLDREC.cpy#L29) |
