# REJREC

```text
REJREC: REJECTED TRANSACTION, 110 BYTES
THE TRANSACTION AS RECEIVED, WITH THE STEP THAT REJECTED IT
AND WHY. ACCOUNT MAINTENANCE REJECTS ITS REQUESTS IN THE SAME
FORM, THE REQUEST IN PLACE OF THE TRANSACTION. THE DATE IT
WAS REJECTED IS TEXT, YYYY-MM-DD, AS THE BRANCHES READ IT.

MODIFICATION HISTORY
DATE       BY  CHANGE
1987-03-12 JA  ORIGINAL VERSION.
1988-05-09 JA  CR-0031 ACCOUNT MAINTENANCE REASONS.
2009-03-02 MV  CR-0444 THE DATE OF THE REJECTION. RECORD
               WIDENED TO 110 BYTES.
```

Copybook [REJREC.cpy](../../../pangaea/copybooks/REJREC.cpy), used by [ACCTMNT](../programs/ACCTMNT.md), [TXPOST](../programs/TXPOST.md), [TXVALID](../programs/TXVALID.md).

## REJECT-RECORD

110 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | REJECT-RECORD | group | 1 | 110 | | | | [REJREC.cpy:15](../../../pangaea/copybooks/REJREC.cpy#L15) |
| 05 | ├ RJ-PROGRAM | text | 1 | 8 | X(8) | | | [REJREC.cpy:16](../../../pangaea/copybooks/REJREC.cpy#L16) |
| 05 | ├ RJ-REASON | code | 9 | 4 | X(4) | | | [REJREC.cpy:17](../../../pangaea/copybooks/REJREC.cpy#L17) |
| 88 | │ ├ RJ-BAD-ACCOUNT | | | | | | VALUE 'ACCT' | [REJREC.cpy:18](../../../pangaea/copybooks/REJREC.cpy#L18) |
| 88 | │ ├ RJ-BAD-COUNTER | | | | | | VALUE 'CTRA' | [REJREC.cpy:19](../../../pangaea/copybooks/REJREC.cpy#L19) |
| 88 | │ ├ RJ-BAD-CODE | | | | | | VALUE 'CODE' | [REJREC.cpy:20](../../../pangaea/copybooks/REJREC.cpy#L20) |
| 88 | │ ├ RJ-BAD-AMOUNT | | | | | | VALUE 'AMNT' | [REJREC.cpy:21](../../../pangaea/copybooks/REJREC.cpy#L21) |
| 88 | │ ├ RJ-BAD-TYPE | | | | | | VALUE 'TYPE' | [REJREC.cpy:22](../../../pangaea/copybooks/REJREC.cpy#L22) |
| 88 | │ ├ RJ-NO-ACCOUNT | | | | | | VALUE 'UNKN' | [REJREC.cpy:23](../../../pangaea/copybooks/REJREC.cpy#L23) |
| 88 | │ ├ RJ-CLOSED | | | | | | VALUE 'CLSD' | [REJREC.cpy:24](../../../pangaea/copybooks/REJREC.cpy#L24) |
| 88 | │ ├ RJ-FROZEN | | | | | | VALUE 'FRZN' | [REJREC.cpy:25](../../../pangaea/copybooks/REJREC.cpy#L25) |
| 88 | │ ├ RJ-NO-FUNDS | | | | | | VALUE 'FUND' | [REJREC.cpy:26](../../../pangaea/copybooks/REJREC.cpy#L26) |
| 88 | │ ├ RJ-EXISTS | | | | | | VALUE 'EXST' | [REJREC.cpy:27](../../../pangaea/copybooks/REJREC.cpy#L27) |
| 88 | │ ├ RJ-HAS-BALANCE | | | | | | VALUE 'BALN' | [REJREC.cpy:28](../../../pangaea/copybooks/REJREC.cpy#L28) |
| 88 | │ └ RJ-BAD-STATUS | | | | | | VALUE 'STAT' | [REJREC.cpy:29](../../../pangaea/copybooks/REJREC.cpy#L29) |
| 05 | ├ RJ-TRAN | text | 13 | 80 | X(80) | | | [REJREC.cpy:30](../../../pangaea/copybooks/REJREC.cpy#L30) |
| 05 | ├ FILLER | filler | 93 | 8 | X(8) | | | [REJREC.cpy:31](../../../pangaea/copybooks/REJREC.cpy#L31) |
| 05 | └ RJ-REJECT-DATE | date | 101 | 10 | X(10) | | | [REJREC.cpy:32](../../../pangaea/copybooks/REJREC.cpy#L32) |
