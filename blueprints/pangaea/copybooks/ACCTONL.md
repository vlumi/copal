# ACCTONL

```text
ACCTONL: ONLINE ACCOUNT RECORD, 100 BYTES
INDEXED BY ACCOUNT NUMBER. LOADED FROM THE MASTER EACH NIGHT
BY ACCTLOAD; DURING THE DAY THE AUTHORIZATIONS ADD WHAT THEY
APPROVE TO THE HOLDS, AND THE AVAILABLE BALANCE IS THE
BALANCE PLUS THE OVERDRAFT LIMIT LESS THE HOLDS.

MODIFICATION HISTORY
DATE       BY  CHANGE
2003-05-19 MV  CR-0342 ORIGINAL VERSION.
```

Copybook [ACCTONL.cpy](../../../pangaea/copybooks/ACCTONL.cpy), used by [ACCTLOAD](../programs/ACCTLOAD.md), [ATMBAL](../programs/ATMBAL.md), [ATMWDL](../programs/ATMWDL.md), [AUTHPUR](../programs/AUTHPUR.md), [AUTHREV](../programs/AUTHREV.md), [AUTHSRV](../programs/AUTHSRV.md).

## ACCT-ONLINE-RECORD

100 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | ACCT-ONLINE-RECORD | group | 1 | 100 | | | | [ACCTONL.cpy:12](../../../pangaea/copybooks/ACCTONL.cpy#L12) |
| 05 | AO-ACCT-NO | key | 1 | 10 | 9(10) | | | [ACCTONL.cpy:13](../../../pangaea/copybooks/ACCTONL.cpy#L13) |
| 05 | AO-CUST-NO | number | 11 | 8 | 9(8) | | | [ACCTONL.cpy:14](../../../pangaea/copybooks/ACCTONL.cpy#L14) |
| 05 | AO-SHORT-NAME | text | 19 | 20 | X(20) | | | [ACCTONL.cpy:15](../../../pangaea/copybooks/ACCTONL.cpy#L15) |
| 05 | AO-ACCT-TYPE | code | 39 | 1 | X | | | [ACCTONL.cpy:16](../../../pangaea/copybooks/ACCTONL.cpy#L16) |
| 88 | AO-CURRENT | | | | | | VALUE 'C' | [ACCTONL.cpy:17](../../../pangaea/copybooks/ACCTONL.cpy#L17) |
| 88 | AO-SAVINGS | | | | | | VALUE 'S' | [ACCTONL.cpy:18](../../../pangaea/copybooks/ACCTONL.cpy#L18) |
| 05 | AO-STATUS | code | 40 | 1 | X | | | [ACCTONL.cpy:19](../../../pangaea/copybooks/ACCTONL.cpy#L19) |
| 88 | AO-OPEN | | | | | | VALUE 'O' | [ACCTONL.cpy:20](../../../pangaea/copybooks/ACCTONL.cpy#L20) |
| 88 | AO-FROZEN | | | | | | VALUE 'F' | [ACCTONL.cpy:21](../../../pangaea/copybooks/ACCTONL.cpy#L21) |
| 05 | AO-BALANCE | amount | 41 | 13 | S9(11)V99 | | THE BALANCE AT THE END OF THE BUSINESS DAY IT WAS LOADED ON. | [ACCTONL.cpy:23](../../../pangaea/copybooks/ACCTONL.cpy#L23) |
| 05 | AO-OVERDRAFT-LIMIT | amount | 54 | 9 | 9(7)V99 | | | [ACCTONL.cpy:24](../../../pangaea/copybooks/ACCTONL.cpy#L24) |
| 05 | AO-HOLDS | amount | 63 | 13 | S9(11)V99 | | | [ACCTONL.cpy:25](../../../pangaea/copybooks/ACCTONL.cpy#L25) |
| 05 | AO-AS-OF-DATE | date | 76 | 8 | 9(8) | | | [ACCTONL.cpy:26](../../../pangaea/copybooks/ACCTONL.cpy#L26) |
| 05 | FILLER | filler | 84 | 17 | X(17) | | | [ACCTONL.cpy:27](../../../pangaea/copybooks/ACCTONL.cpy#L27) |
