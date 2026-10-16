# GLREC

```text
GLREC: GENERAL LEDGER SUMMARY RECORD, 60 BYTES
ONE PER GENERAL LEDGER ACCOUNT, FOR A MONTH: THE DEBITS AND
CREDITS THE MONTH'S POSTINGS MADE TO IT.

MODIFICATION HISTORY
DATE       BY  CHANGE
2004-03-15 TL  CR-0390 ORIGINAL VERSION.
```

Copybook [GLREC.cpy](../../../pangaea/copybooks/GLREC.cpy), used by [GLSUM](../programs/GLSUM.md).

## GL-SUMMARY-RECORD

60 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | GL-SUMMARY-RECORD | 1 | 60 | | | | [GLREC.cpy:10](../../../pangaea/copybooks/GLREC.cpy#L10) |
| 05 | GL-MONTH | 1 | 6 | 9(6) | | | [GLREC.cpy:11](../../../pangaea/copybooks/GLREC.cpy#L11) |
| 05 | GL-ACCOUNT | 7 | 4 | 9(4) | | | [GLREC.cpy:12](../../../pangaea/copybooks/GLREC.cpy#L12) |
| 05 | GL-ACCOUNT-NAME | 11 | 24 | X(24) | | | [GLREC.cpy:13](../../../pangaea/copybooks/GLREC.cpy#L13) |
| 05 | GL-DEBITS | 35 | 8 | S9(13)V99 | COMP-3 | | [GLREC.cpy:14](../../../pangaea/copybooks/GLREC.cpy#L14) |
| 05 | GL-CREDITS | 43 | 8 | S9(13)V99 | COMP-3 | | [GLREC.cpy:15](../../../pangaea/copybooks/GLREC.cpy#L15) |
| 05 | FILLER | 51 | 10 | X(10) | | | [GLREC.cpy:16](../../../pangaea/copybooks/GLREC.cpy#L16) |
