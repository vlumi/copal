# cardacct

```text
cardacct: a card's account, indexed by card number: its credit limit,
its balance, and the balance its last statement showed.

history
2011-04-04 rk  lcs-0102 original version.
2013-06-03 rk  lcs-0219 billing cycle day and paid in full flag.
```

Copybook [cardacct.cpy](../../../laurasia/copybooks/cardacct.cpy), used by [ACCTLD](../programs/ACCTLD.md), [ACCTUL](../programs/ACCTUL.md), [CYCLE](../programs/CYCLE.md), [PAYIN](../programs/PAYIN.md), [POST](../programs/POST.md).

## CARD-ACCOUNT

61 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | CARD-ACCOUNT | group | 1 | 61 | | | | [cardacct.cpy:8](../../../laurasia/copybooks/cardacct.cpy#L8) |
| 05 | ├ CA-CARD-NO | key | 1 | 16 | X(16) | | | [cardacct.cpy:9](../../../laurasia/copybooks/cardacct.cpy#L9) |
| 05 | ├ CA-ACCOUNT-NO | text | 17 | 10 | X(10) | | | [cardacct.cpy:10](../../../laurasia/copybooks/cardacct.cpy#L10) |
| 05 | ├ CA-PRODUCT | code | 27 | 2 | X(2) | | | [cardacct.cpy:11](../../../laurasia/copybooks/cardacct.cpy#L11) |
| 88 | │ ├ CA-CLASSIC | | | | | | VALUE 'CL' | [cardacct.cpy:12](../../../laurasia/copybooks/cardacct.cpy#L12) |
| 88 | │ ├ CA-GOLD | | | | | | VALUE 'GD' | [cardacct.cpy:13](../../../laurasia/copybooks/cardacct.cpy#L13) |
| 88 | │ └ CA-BUSINESS | | | | | | VALUE 'BZ' | [cardacct.cpy:14](../../../laurasia/copybooks/cardacct.cpy#L14) |
| 05 | ├ CA-STATUS | code | 29 | 1 | X | | | [cardacct.cpy:15](../../../laurasia/copybooks/cardacct.cpy#L15) |
| 88 | │ ├ CA-ACTIVE | | | | | | VALUE 'A' | [cardacct.cpy:16](../../../laurasia/copybooks/cardacct.cpy#L16) |
| 88 | │ ├ CA-BLOCKED | | | | | | VALUE 'B' | [cardacct.cpy:17](../../../laurasia/copybooks/cardacct.cpy#L17) |
| 88 | │ └ CA-CLOSED | | | | | | VALUE 'C' | [cardacct.cpy:18](../../../laurasia/copybooks/cardacct.cpy#L18) |
| 05 | ├ CA-LIMIT | amount | 30 | 9 | 9(7)V99 | | | [cardacct.cpy:19](../../../laurasia/copybooks/cardacct.cpy#L19) |
| 05 | ├ CA-BALANCE | amount | 39 | 10 | S9(7)V99 | | SIGN LEADING SEPARATE | [cardacct.cpy:20](../../../laurasia/copybooks/cardacct.cpy#L20) |
| 05 | ├ CA-STATEMENT-BALANCE | amount | 49 | 10 | S9(7)V99 | | SIGN LEADING SEPARATE | [cardacct.cpy:21](../../../laurasia/copybooks/cardacct.cpy#L21) |
| 05 | ├ CA-CYCLE-DAY | number | 59 | 2 | 99 | | | [cardacct.cpy:22](../../../laurasia/copybooks/cardacct.cpy#L22) |
| 05 | └ CA-PAID-IN-FULL | code | 61 | 1 | X | | | [cardacct.cpy:23](../../../laurasia/copybooks/cardacct.cpy#L23) |
| 88 |   └ CA-WAS-PAID-IN-FULL | | | | | | VALUE 'Y' | [cardacct.cpy:24](../../../laurasia/copybooks/cardacct.cpy#L24) |
