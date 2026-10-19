# authlog

```text
authlog: the online side's authorizations of the day, a record a line,
the ones given in stand-in, when the issuer could not answer, marked.

history
2011-04-04 rk  lcs-0102 original version.
```

Copybook [authlog.cpy](../../../laurasia/copybooks/authlog.cpy), used by [MATCH](../programs/MATCH.md).

## AUTH-ENTRY

44 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | AUTH-ENTRY | group | 1 | 44 | | | | [authlog.cpy:7](../../../laurasia/copybooks/authlog.cpy#L7) |
| 05 | ├ AL-AUTH-CODE | text | 1 | 6 | X(6) | | | [authlog.cpy:8](../../../laurasia/copybooks/authlog.cpy#L8) |
| 05 | ├ AL-CARD-NO | text | 7 | 16 | X(16) | | | [authlog.cpy:9](../../../laurasia/copybooks/authlog.cpy#L9) |
| 05 | ├ AL-AMOUNT | amount | 23 | 12 | S9(9)V99 | | SIGN LEADING SEPARATE | [authlog.cpy:10](../../../laurasia/copybooks/authlog.cpy#L10) |
| 05 | ├ AL-DECISION | code | 35 | 1 | X | | | [authlog.cpy:11](../../../laurasia/copybooks/authlog.cpy#L11) |
| 88 | │ ├ AL-APPROVED | | | | | | VALUE 'A' | [authlog.cpy:12](../../../laurasia/copybooks/authlog.cpy#L12) |
| 88 | │ └ AL-DECLINED | | | | | | VALUE 'D' | [authlog.cpy:13](../../../laurasia/copybooks/authlog.cpy#L13) |
| 05 | ├ AL-STAND-IN | code | 36 | 1 | X | | | [authlog.cpy:14](../../../laurasia/copybooks/authlog.cpy#L14) |
| 88 | │ └ AL-IN-STAND-IN | | | | | | VALUE 'Y' | [authlog.cpy:15](../../../laurasia/copybooks/authlog.cpy#L15) |
| 05 | └ AL-DATE | date | 37 | 8 | 9(8) | | | [authlog.cpy:16](../../../laurasia/copybooks/authlog.cpy#L16) |
