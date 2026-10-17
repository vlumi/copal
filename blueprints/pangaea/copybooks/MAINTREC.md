# MAINTREC

```text
MAINTREC: ACCOUNT MAINTENANCE REQUEST, 80 BYTES
THE BRANCHES' REQUESTS TO OPEN AN ACCOUNT, CLOSE ONE, FREEZE
OR UNFREEZE ONE, OR CHANGE ITS OVERDRAFT LIMIT, IN ACCOUNT
NUMBER AND REQUEST ID ORDER. AN OPENING CARRIES THE NEW
ACCOUNT'S CUSTOMER, NAME AND TYPE; A LIMIT CHANGE THE NEW
LIMIT.

MODIFICATION HISTORY
DATE       BY  CHANGE
1988-05-09 JA  CR-0031 ORIGINAL VERSION.
```

Copybook [MAINTREC.cpy](../../../pangaea/copybooks/MAINTREC.cpy), used by [ACCTMNT](../programs/ACCTMNT.md).

## MAINT-REQUEST

80 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | MAINT-REQUEST | group | 1 | 80 | | | | [MAINTREC.cpy:13](../../../pangaea/copybooks/MAINTREC.cpy#L13) |
| 05 | ├ MR-REQUEST-ID | number | 1 | 10 | 9(10) | | | [MAINTREC.cpy:14](../../../pangaea/copybooks/MAINTREC.cpy#L14) |
| 05 | ├ MR-REQUEST-TYPE | code | 11 | 1 | X | | | [MAINTREC.cpy:15](../../../pangaea/copybooks/MAINTREC.cpy#L15) |
| 88 | │ ├ MR-OPEN | | | | | | VALUE 'O' | [MAINTREC.cpy:16](../../../pangaea/copybooks/MAINTREC.cpy#L16) |
| 88 | │ ├ MR-CLOSE | | | | | | VALUE 'C' | [MAINTREC.cpy:17](../../../pangaea/copybooks/MAINTREC.cpy#L17) |
| 88 | │ ├ MR-FREEZE | | | | | | VALUE 'F' | [MAINTREC.cpy:18](../../../pangaea/copybooks/MAINTREC.cpy#L18) |
| 88 | │ ├ MR-UNFREEZE | | | | | | VALUE 'U' | [MAINTREC.cpy:19](../../../pangaea/copybooks/MAINTREC.cpy#L19) |
| 88 | │ └ MR-LIMIT | | | | | | VALUE 'L' | [MAINTREC.cpy:20](../../../pangaea/copybooks/MAINTREC.cpy#L20) |
| 05 | ├ MR-ACCT-NO | number | 12 | 10 | 9(10) | | | [MAINTREC.cpy:21](../../../pangaea/copybooks/MAINTREC.cpy#L21) |
| 05 | ├ MR-CUST-NO | number | 22 | 8 | 9(8) | | | [MAINTREC.cpy:22](../../../pangaea/copybooks/MAINTREC.cpy#L22) |
| 05 | ├ MR-SHORT-NAME | text | 30 | 20 | X(20) | | | [MAINTREC.cpy:23](../../../pangaea/copybooks/MAINTREC.cpy#L23) |
| 05 | ├ MR-ACCT-TYPE | text | 50 | 1 | X | | | [MAINTREC.cpy:24](../../../pangaea/copybooks/MAINTREC.cpy#L24) |
| 05 | ├ MR-OVERDRAFT-LIMIT | amount | 51 | 9 | 9(7)V99 | | | [MAINTREC.cpy:25](../../../pangaea/copybooks/MAINTREC.cpy#L25) |
| 05 | └ FILLER | filler | 60 | 21 | X(21) | | | [MAINTREC.cpy:26](../../../pangaea/copybooks/MAINTREC.cpy#L26) |
