# clrrec

```text
clrrec: the networks' clearing file, a record a line, of five layouts
told apart by their first two characters. A file header, then for each
acquirer's batch a batch header, its presentments and a batch trailer
with their count and amount, then a file trailer with the batches',
the presentments' and the amounts' totals.

history
2011-04-04 rk  lcs-0102 original version.
2016-09-12 ah  lcs-0377 card product on the presentment.
```

Copybook [clrrec.cpy](../../../laurasia/copybooks/clrrec.cpy), used by [CLRIN](../programs/CLRIN.md).

## Contents

- [CLR-FILE-HEADER](#clr-file-header)
- [CLR-BATCH-HEADER](#clr-batch-header)
- [CLR-PRESENTMENT](#clr-presentment)
- [CLR-BATCH-TRAILER](#clr-batch-trailer)
- [CLR-FILE-TRAILER](#clr-file-trailer)

## CLR-FILE-HEADER

18 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | CLR-FILE-HEADER | group | 1 | 18 | | | | [clrrec.cpy:11](../../../laurasia/copybooks/clrrec.cpy#L11) |
| 05 | ├ CLR-RECORD-TYPE | code | 1 | 2 | X(2) | | | [clrrec.cpy:12](../../../laurasia/copybooks/clrrec.cpy#L12) |
| 88 | │ ├ IS-FILE-HEADER | | | | | | VALUE 'FH' | [clrrec.cpy:13](../../../laurasia/copybooks/clrrec.cpy#L13) |
| 88 | │ ├ IS-BATCH-HEADER | | | | | | VALUE 'BH' | [clrrec.cpy:14](../../../laurasia/copybooks/clrrec.cpy#L14) |
| 88 | │ ├ IS-PRESENTMENT | | | | | | VALUE 'PR' | [clrrec.cpy:15](../../../laurasia/copybooks/clrrec.cpy#L15) |
| 88 | │ ├ IS-BATCH-TRAILER | | | | | | VALUE 'BT' | [clrrec.cpy:16](../../../laurasia/copybooks/clrrec.cpy#L16) |
| 88 | │ └ IS-FILE-TRAILER | | | | | | VALUE 'FT' | [clrrec.cpy:17](../../../laurasia/copybooks/clrrec.cpy#L17) |
| 05 | ├ FH-NETWORK | text | 3 | 4 | X(4) | | | [clrrec.cpy:18](../../../laurasia/copybooks/clrrec.cpy#L18) |
| 05 | ├ FH-FILE-DATE | date | 7 | 8 | 9(8) | | | [clrrec.cpy:19](../../../laurasia/copybooks/clrrec.cpy#L19) |
| 05 | └ FH-FILE-SEQ | number | 15 | 4 | 9(4) | | | [clrrec.cpy:20](../../../laurasia/copybooks/clrrec.cpy#L20) |

## CLR-BATCH-HEADER

14 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | CLR-BATCH-HEADER | group | 1 | 14 | | | | [clrrec.cpy:21](../../../laurasia/copybooks/clrrec.cpy#L21) |
| 05 | ├ FILLER | filler | 1 | 2 | X(2) | | | [clrrec.cpy:22](../../../laurasia/copybooks/clrrec.cpy#L22) |
| 05 | ├ BH-BATCH-NO | number | 3 | 6 | 9(6) | | | [clrrec.cpy:23](../../../laurasia/copybooks/clrrec.cpy#L23) |
| 05 | └ BH-ACQUIRER | text | 9 | 6 | X(6) | | | [clrrec.cpy:24](../../../laurasia/copybooks/clrrec.cpy#L24) |

## CLR-PRESENTMENT

61 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | CLR-PRESENTMENT | group | 1 | 61 | | | | [clrrec.cpy:25](../../../laurasia/copybooks/clrrec.cpy#L25) |
| 05 | ├ FILLER | filler | 1 | 2 | X(2) | | | [clrrec.cpy:26](../../../laurasia/copybooks/clrrec.cpy#L26) |
| 05 | ├ PR-CARD-NO | text | 3 | 16 | X(16) | | | [clrrec.cpy:27](../../../laurasia/copybooks/clrrec.cpy#L27) |
| 05 | ├ PR-AMOUNT | amount | 19 | 12 | S9(9)V99 | | SIGN LEADING SEPARATE | [clrrec.cpy:28](../../../laurasia/copybooks/clrrec.cpy#L28) |
| 05 | ├ PR-AUTH-CODE | text | 31 | 6 | X(6) | | | [clrrec.cpy:29](../../../laurasia/copybooks/clrrec.cpy#L29) |
| 05 | ├ PR-MERCHANT | text | 37 | 10 | X(10) | | | [clrrec.cpy:30](../../../laurasia/copybooks/clrrec.cpy#L30) |
| 05 | ├ PR-MCC | number | 47 | 4 | 9(4) | | | [clrrec.cpy:31](../../../laurasia/copybooks/clrrec.cpy#L31) |
| 05 | ├ PR-REGION | text | 51 | 1 | X | | | [clrrec.cpy:32](../../../laurasia/copybooks/clrrec.cpy#L32) |
| 05 | ├ PR-PRODUCT | text | 52 | 2 | X(2) | | | [clrrec.cpy:33](../../../laurasia/copybooks/clrrec.cpy#L33) |
| 05 | └ PR-TXN-DATE | date | 54 | 8 | 9(8) | | | [clrrec.cpy:34](../../../laurasia/copybooks/clrrec.cpy#L34) |

## CLR-BATCH-TRAILER

24 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | CLR-BATCH-TRAILER | group | 1 | 24 | | | | [clrrec.cpy:35](../../../laurasia/copybooks/clrrec.cpy#L35) |
| 05 | ├ FILLER | filler | 1 | 2 | X(2) | | | [clrrec.cpy:36](../../../laurasia/copybooks/clrrec.cpy#L36) |
| 05 | ├ BT-COUNT | number | 3 | 6 | 9(6) | | | [clrrec.cpy:37](../../../laurasia/copybooks/clrrec.cpy#L37) |
| 05 | └ BT-AMOUNT | amount | 9 | 16 | S9(13)V99 | | SIGN LEADING SEPARATE | [clrrec.cpy:38](../../../laurasia/copybooks/clrrec.cpy#L38) |

## CLR-FILE-TRAILER

32 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | CLR-FILE-TRAILER | group | 1 | 32 | | | | [clrrec.cpy:39](../../../laurasia/copybooks/clrrec.cpy#L39) |
| 05 | ├ FILLER | filler | 1 | 2 | X(2) | | | [clrrec.cpy:40](../../../laurasia/copybooks/clrrec.cpy#L40) |
| 05 | ├ FT-BATCH-COUNT | number | 3 | 4 | 9(4) | | | [clrrec.cpy:41](../../../laurasia/copybooks/clrrec.cpy#L41) |
| 05 | ├ FT-RECORD-COUNT | number | 7 | 8 | 9(8) | | | [clrrec.cpy:42](../../../laurasia/copybooks/clrrec.cpy#L42) |
| 05 | └ FT-AMOUNT | amount | 15 | 18 | S9(15)V99 | | SIGN LEADING SEPARATE | [clrrec.cpy:43](../../../laurasia/copybooks/clrrec.cpy#L43) |
