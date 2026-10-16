# ACCRREC

```text
ACCRREC: INTEREST ACCRUAL RECORD, 40 BYTES
ONE PER ACCOUNT, IN ACCOUNT NUMBER ORDER LIKE THE MASTER.
INTEREST IS ACCRUED FOR EVERY CALENDAR DAY BEFORE THE THROUGH
DATE, TO SIX DECIMALS, AND IS NEGATIVE FOR AN OVERDRAWN
BALANCE. WHEN A MONTH CLOSES, ITS INTEREST BECOMES DUE,
ROUNDED TO THE CENT; THE FRACTION STAYS ACCRUED.
THE DUE MONTH IS ZERO ON THE DAYS NO MONTH CLOSES.

MODIFICATION HISTORY
DATE       BY  CHANGE
2001-09-03 TL  CR-0298 ORIGINAL VERSION.
```

Copybook [ACCRREC.cpy](../../../pangaea/copybooks/ACCRREC.cpy), used by [INTACCR](../programs/INTACCR.md), [MONCHG](../programs/MONCHG.md).

## ACCRUAL-RECORD

40 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | ACCRUAL-RECORD | 1 | 40 | | | | [ACCRREC.cpy:14](../../../pangaea/copybooks/ACCRREC.cpy#L14) |
| 05 | AC-ACCT-NO | 1 | 10 | 9(10) | | | [ACCRREC.cpy:15](../../../pangaea/copybooks/ACCRREC.cpy#L15) |
| 05 | AC-THROUGH-DATE | 11 | 8 | 9(8) | | | [ACCRREC.cpy:16](../../../pangaea/copybooks/ACCRREC.cpy#L16) |
| 05 | AC-ACCRUED | 19 | 7 | S9(7)V9(6) | COMP-3 | | [ACCRREC.cpy:17](../../../pangaea/copybooks/ACCRREC.cpy#L17) |
| 05 | AC-DUE | 26 | 5 | S9(7)V99 | COMP-3 | | [ACCRREC.cpy:18](../../../pangaea/copybooks/ACCRREC.cpy#L18) |
| 05 | AC-DUE-MONTH | 31 | 6 | 9(6) | | | [ACCRREC.cpy:19](../../../pangaea/copybooks/ACCRREC.cpy#L19) |
| 05 | FILLER | 37 | 4 | X(4) | | | [ACCRREC.cpy:20](../../../pangaea/copybooks/ACCRREC.cpy#L20) |
