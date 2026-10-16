# CODETAB

```text
CODETAB: TRANSACTION CODE TABLE
FOR EACH TRANSACTION CODE, WHAT A STATEMENT CALLS IT AND THE
GENERAL LEDGER ACCOUNT ITS POSTINGS MEET ON THE OTHER SIDE OF
CUSTOMER DEPOSITS. COPIED WITH THE PROGRAM'S OWN PREFIX FOR
THE :P: TAGS: COPY CODETAB REPLACING ==:P:== BY ==ST==.

MODIFICATION HISTORY
DATE       BY  CHANGE
2004-03-15 TL  CR-0390 ORIGINAL VERSION, FROM STMTPRT, FOR THE
               GENERAL LEDGER SUMMARY.
```

Copybook [CODETAB.cpy](../../../pangaea/copybooks/CODETAB.cpy), used by [GLSUM](../programs/GLSUM.md), [STMTPRT](../programs/STMTPRT.md).

## :P:-CODE-VALUES

312 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | :P:-CODE-VALUES | 1 | 312 | | | | [CODETAB.cpy:13](../../../pangaea/copybooks/CODETAB.cpy#L13) |
| 05 | FILLER | 1 | 2 | XX | | VALUE 'DP' | [CODETAB.cpy:14](../../../pangaea/copybooks/CODETAB.cpy#L14) |
| 05 | FILLER | 3 | 20 | X(20) | | VALUE 'DEPOSIT' | [CODETAB.cpy:15](../../../pangaea/copybooks/CODETAB.cpy#L15) |
| 05 | FILLER | 23 | 4 | 9(4) | | VALUE 1000 | [CODETAB.cpy:17](../../../pangaea/copybooks/CODETAB.cpy#L17) |
| 05 | FILLER | 27 | 2 | XX | | VALUE 'WD' | [CODETAB.cpy:18](../../../pangaea/copybooks/CODETAB.cpy#L18) |
| 05 | FILLER | 29 | 20 | X(20) | | VALUE 'WITHDRAWAL' | [CODETAB.cpy:19](../../../pangaea/copybooks/CODETAB.cpy#L19) |
| 05 | FILLER | 49 | 4 | 9(4) | | VALUE 1000 | [CODETAB.cpy:21](../../../pangaea/copybooks/CODETAB.cpy#L21) |
| 05 | FILLER | 53 | 2 | XX | | VALUE 'FE' | [CODETAB.cpy:22](../../../pangaea/copybooks/CODETAB.cpy#L22) |
| 05 | FILLER | 55 | 20 | X(20) | | VALUE 'FEE' | [CODETAB.cpy:23](../../../pangaea/copybooks/CODETAB.cpy#L23) |
| 05 | FILLER | 75 | 4 | 9(4) | | VALUE 4100 | [CODETAB.cpy:25](../../../pangaea/copybooks/CODETAB.cpy#L25) |
| 05 | FILLER | 79 | 2 | XX | | VALUE 'TR' | [CODETAB.cpy:26](../../../pangaea/copybooks/CODETAB.cpy#L26) |
| 05 | FILLER | 81 | 20 | X(20) | | VALUE 'TRANSFER OUT' | [CODETAB.cpy:27](../../../pangaea/copybooks/CODETAB.cpy#L27) |
| 05 | FILLER | 101 | 4 | 9(4) | | VALUE 2900 | [CODETAB.cpy:29](../../../pangaea/copybooks/CODETAB.cpy#L29) |
| 05 | FILLER | 105 | 2 | XX | | VALUE 'TC' | [CODETAB.cpy:30](../../../pangaea/copybooks/CODETAB.cpy#L30) |
| 05 | FILLER | 107 | 20 | X(20) | | VALUE 'TRANSFER IN' | [CODETAB.cpy:31](../../../pangaea/copybooks/CODETAB.cpy#L31) |
| 05 | FILLER | 127 | 4 | 9(4) | | VALUE 2900 | [CODETAB.cpy:33](../../../pangaea/copybooks/CODETAB.cpy#L33) |
| 05 | FILLER | 131 | 2 | XX | | VALUE 'RT' | [CODETAB.cpy:34](../../../pangaea/copybooks/CODETAB.cpy#L34) |
| 05 | FILLER | 133 | 20 | X(20) | | VALUE 'TRANSFER RETURNED' | [CODETAB.cpy:35](../../../pangaea/copybooks/CODETAB.cpy#L35) |
| 05 | FILLER | 153 | 4 | 9(4) | | VALUE 2900 | [CODETAB.cpy:37](../../../pangaea/copybooks/CODETAB.cpy#L37) |
| 05 | FILLER | 157 | 2 | XX | | VALUE 'IN' | [CODETAB.cpy:38](../../../pangaea/copybooks/CODETAB.cpy#L38) |
| 05 | FILLER | 159 | 20 | X(20) | | VALUE 'CREDIT INTEREST' | [CODETAB.cpy:39](../../../pangaea/copybooks/CODETAB.cpy#L39) |
| 05 | FILLER | 179 | 4 | 9(4) | | VALUE 5000 | [CODETAB.cpy:41](../../../pangaea/copybooks/CODETAB.cpy#L41) |
| 05 | FILLER | 183 | 2 | XX | | VALUE 'OD' | [CODETAB.cpy:42](../../../pangaea/copybooks/CODETAB.cpy#L42) |
| 05 | FILLER | 185 | 20 | X(20) | | VALUE 'OVERDRAFT INTEREST' | [CODETAB.cpy:43](../../../pangaea/copybooks/CODETAB.cpy#L43) |
| 05 | FILLER | 205 | 4 | 9(4) | | VALUE 4000 | [CODETAB.cpy:45](../../../pangaea/copybooks/CODETAB.cpy#L45) |
| 05 | FILLER | 209 | 2 | XX | | VALUE 'MF' | [CODETAB.cpy:46](../../../pangaea/copybooks/CODETAB.cpy#L46) |
| 05 | FILLER | 211 | 20 | X(20) | | VALUE 'MAINTENANCE FEE' | [CODETAB.cpy:47](../../../pangaea/copybooks/CODETAB.cpy#L47) |
| 05 | FILLER | 231 | 4 | 9(4) | | VALUE 4100 | [CODETAB.cpy:49](../../../pangaea/copybooks/CODETAB.cpy#L49) |
| 05 | FILLER | 235 | 2 | XX | | VALUE 'CP' | [CODETAB.cpy:50](../../../pangaea/copybooks/CODETAB.cpy#L50) |
| 05 | FILLER | 237 | 20 | X(20) | | VALUE 'CARD PURCHASE' | [CODETAB.cpy:51](../../../pangaea/copybooks/CODETAB.cpy#L51) |
| 05 | FILLER | 257 | 4 | 9(4) | | VALUE 2100 | [CODETAB.cpy:53](../../../pangaea/copybooks/CODETAB.cpy#L53) |
| 05 | FILLER | 261 | 2 | XX | | VALUE 'AW' | [CODETAB.cpy:54](../../../pangaea/copybooks/CODETAB.cpy#L54) |
| 05 | FILLER | 263 | 20 | X(20) | | VALUE 'ATM WITHDRAWAL' | [CODETAB.cpy:55](../../../pangaea/copybooks/CODETAB.cpy#L55) |
| 05 | FILLER | 283 | 4 | 9(4) | | VALUE 1000 | [CODETAB.cpy:57](../../../pangaea/copybooks/CODETAB.cpy#L57) |
| 05 | FILLER | 287 | 2 | XX | | VALUE 'AF' | [CODETAB.cpy:58](../../../pangaea/copybooks/CODETAB.cpy#L58) |
| 05 | FILLER | 289 | 20 | X(20) | | VALUE 'ATM FEE' | [CODETAB.cpy:59](../../../pangaea/copybooks/CODETAB.cpy#L59) |
| 05 | FILLER | 309 | 4 | 9(4) | | VALUE 4100 | [CODETAB.cpy:61](../../../pangaea/copybooks/CODETAB.cpy#L61) |

## :P:-CODE-TABLE

312 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | :P:-CODE-TABLE | 1 | 312 | | | REDEFINES :P:-CODE-VALUES | [CODETAB.cpy:62](../../../pangaea/copybooks/CODETAB.cpy#L62) |
| 05 | :P:-CODE-ENTRY | 1 | 26 | | | OCCURS 12; INDEXED BY :P:-CX | [CODETAB.cpy:63](../../../pangaea/copybooks/CODETAB.cpy#L63) |
| 10 | :P:-CODE | 1 | 2 | XX | | | [CODETAB.cpy:65](../../../pangaea/copybooks/CODETAB.cpy#L65) |
| 10 | :P:-CODE-TEXT | 3 | 20 | X(20) | | | [CODETAB.cpy:66](../../../pangaea/copybooks/CODETAB.cpy#L66) |
| 10 | :P:-CODE-GL-ACCOUNT | 23 | 4 | 9(4) | | | [CODETAB.cpy:67](../../../pangaea/copybooks/CODETAB.cpy#L67) |
