# CARDREC

```text
CARDREC: CARD RECORD, 80 BYTES
THE CARD MASTER, IN CARD NUMBER ORDER, AND THE ONLINE CARD
FILE LOADED FROM IT EACH NIGHT, INDEXED BY CARD NUMBER. A CARD
IS VALID TO THE END OF ITS EXPIRY MONTH. THE AUTHORIZATIONS
KEEP THE TOTALS OF THE DAY OF CR-TOTALS-DATE, AND THE ATM
WITHDRAWALS OF ITS MONTH, AND START THEM OVER ON A NEW DAY.

MODIFICATION HISTORY
DATE       BY  CHANGE
2003-05-19 MV  CR-0342 ORIGINAL VERSION.
```

Copybook [CARDREC.cpy](../../../pangaea/copybooks/CARDREC.cpy), used by [ATMBAL](../programs/ATMBAL.md), [ATMWDL](../programs/ATMWDL.md), [AUTHPUR](../programs/AUTHPUR.md), [AUTHREV](../programs/AUTHREV.md), [AUTHSRV](../programs/AUTHSRV.md), [CARDLOAD](../programs/CARDLOAD.md).

## CARD-RECORD

80 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | CARD-RECORD | group | 1 | 80 | | | | [CARDREC.cpy:13](../../../pangaea/copybooks/CARDREC.cpy#L13) |
| 05 | ├ CR-CARD-NO | key | 1 | 16 | 9(16) | | | [CARDREC.cpy:14](../../../pangaea/copybooks/CARDREC.cpy#L14) |
| 05 | ├ CR-ACCT-NO | number | 17 | 10 | 9(10) | | | [CARDREC.cpy:15](../../../pangaea/copybooks/CARDREC.cpy#L15) |
| 05 | ├ CR-CUST-NO | number | 27 | 8 | 9(8) | | | [CARDREC.cpy:16](../../../pangaea/copybooks/CARDREC.cpy#L16) |
| 05 | ├ CR-STATUS | code | 35 | 1 | X | | | [CARDREC.cpy:17](../../../pangaea/copybooks/CARDREC.cpy#L17) |
| 88 | │ ├ CR-ACTIVE | | | | | | VALUE 'A' | [CARDREC.cpy:18](../../../pangaea/copybooks/CARDREC.cpy#L18) |
| 88 | │ └ CR-BLOCKED | | | | | | VALUE 'B' | [CARDREC.cpy:19](../../../pangaea/copybooks/CARDREC.cpy#L19) |
| 05 | ├ CR-EXPIRY-MONTH | date | 36 | 6 | 9(6) | | | [CARDREC.cpy:20](../../../pangaea/copybooks/CARDREC.cpy#L20) |
| 05 | ├ CR-PURCHASE-LIMIT | amount | 42 | 7 | 9(5)V99 | | DAILY LIMITS. | [CARDREC.cpy:22](../../../pangaea/copybooks/CARDREC.cpy#L22) |
| 05 | ├ CR-CASH-LIMIT | amount | 49 | 7 | 9(5)V99 | | | [CARDREC.cpy:23](../../../pangaea/copybooks/CARDREC.cpy#L23) |
| 05 | ├ CR-TOTALS-DATE | date | 56 | 8 | 9(8) | | | [CARDREC.cpy:24](../../../pangaea/copybooks/CARDREC.cpy#L24) |
| 05 | ├ CR-DAY-PURCHASES | amount | 64 | 7 | 9(5)V99 | | | [CARDREC.cpy:25](../../../pangaea/copybooks/CARDREC.cpy#L25) |
| 05 | ├ CR-DAY-CASH | amount | 71 | 7 | 9(5)V99 | | | [CARDREC.cpy:26](../../../pangaea/copybooks/CARDREC.cpy#L26) |
| 05 | └ CR-MONTH-WITHDRAWALS | number | 78 | 3 | 9(3) | | | [CARDREC.cpy:27](../../../pangaea/copybooks/CARDREC.cpy#L27) |
