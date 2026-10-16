# TRANREC

Copybook [TRANREC.cpy](../../../pangaea/copybooks/TRANREC.cpy), used by [TXPOST](../programs/TXPOST.md), [TXVALID](../programs/TXVALID.md).

## TRAN-RECORD

80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | TRAN-RECORD | 1 | 80 | | | | [TRANREC.cpy:8](../../../pangaea/copybooks/TRANREC.cpy#L8) |
| 05 | TR-REC-TYPE | 1 | 1 | X | | | [TRANREC.cpy:9](../../../pangaea/copybooks/TRANREC.cpy#L9) |
| 88 | TR-HEADER | | | | | VALUE 'H' | [TRANREC.cpy:10](../../../pangaea/copybooks/TRANREC.cpy#L10) |
| 88 | TR-DETAIL | | | | | VALUE 'D' | [TRANREC.cpy:11](../../../pangaea/copybooks/TRANREC.cpy#L11) |
| 88 | TR-TRAILER | | | | | VALUE 'T' | [TRANREC.cpy:12](../../../pangaea/copybooks/TRANREC.cpy#L12) |
| 05 | TR-BODY | 2 | 79 | X(79) | | | [TRANREC.cpy:13](../../../pangaea/copybooks/TRANREC.cpy#L13) |
| 05 | TR-HEADER-BODY | 2 | 79 | | | REDEFINES TR-BODY | [TRANREC.cpy:14](../../../pangaea/copybooks/TRANREC.cpy#L14) |
| 10 | TH-BUSINESS-DATE | 2 | 8 | 9(8) | | | [TRANREC.cpy:15](../../../pangaea/copybooks/TRANREC.cpy#L15) |
| 10 | TH-SOURCE | 10 | 8 | X(8) | | | [TRANREC.cpy:16](../../../pangaea/copybooks/TRANREC.cpy#L16) |
| 10 | TH-FILE-SEQ | 18 | 6 | 9(6) | | | [TRANREC.cpy:17](../../../pangaea/copybooks/TRANREC.cpy#L17) |
| 10 | FILLER | 24 | 57 | X(57) | | | [TRANREC.cpy:18](../../../pangaea/copybooks/TRANREC.cpy#L18) |
| 05 | TR-DETAIL-BODY | 2 | 79 | | | REDEFINES TR-BODY | [TRANREC.cpy:19](../../../pangaea/copybooks/TRANREC.cpy#L19) |
| 10 | TD-TRAN-ID | 2 | 10 | 9(10) | | | [TRANREC.cpy:20](../../../pangaea/copybooks/TRANREC.cpy#L20) |
| 10 | TD-ACCT-NO | 12 | 10 | 9(10) | | | [TRANREC.cpy:21](../../../pangaea/copybooks/TRANREC.cpy#L21) |
| 10 | TD-TRAN-CODE | 22 | 2 | XX | | | [TRANREC.cpy:22](../../../pangaea/copybooks/TRANREC.cpy#L22) |
| 88 | TD-DEPOSIT | | | | | VALUE 'DP' | [TRANREC.cpy:23](../../../pangaea/copybooks/TRANREC.cpy#L23) |
| 88 | TD-WITHDRAWAL | | | | | VALUE 'WD' | [TRANREC.cpy:24](../../../pangaea/copybooks/TRANREC.cpy#L24) |
| 88 | TD-FEE | | | | | VALUE 'FE' | [TRANREC.cpy:25](../../../pangaea/copybooks/TRANREC.cpy#L25) |
| 88 | TD-KNOWN-CODE | | | | | VALUE 'DP', 'WD', 'FE' | [TRANREC.cpy:26](../../../pangaea/copybooks/TRANREC.cpy#L26) |
| 10 | TD-AMOUNT | 24 | 11 | 9(9)V99 | | | [TRANREC.cpy:27](../../../pangaea/copybooks/TRANREC.cpy#L27) |
| 10 | TD-DESCRIPTION | 35 | 20 | X(20) | | | [TRANREC.cpy:28](../../../pangaea/copybooks/TRANREC.cpy#L28) |
| 10 | FILLER | 55 | 26 | X(26) | | | [TRANREC.cpy:29](../../../pangaea/copybooks/TRANREC.cpy#L29) |
| 05 | TR-TRAILER-BODY | 2 | 79 | | | REDEFINES TR-BODY | [TRANREC.cpy:30](../../../pangaea/copybooks/TRANREC.cpy#L30) |
| 10 | TT-DETAIL-COUNT | 2 | 8 | 9(8) | | | [TRANREC.cpy:31](../../../pangaea/copybooks/TRANREC.cpy#L31) |
| 10 | TT-AMOUNT-TOTAL | 10 | 15 | 9(13)V99 | | | [TRANREC.cpy:32](../../../pangaea/copybooks/TRANREC.cpy#L32) |
| 10 | TT-ACCT-HASH | 25 | 15 | 9(15) | | | [TRANREC.cpy:33](../../../pangaea/copybooks/TRANREC.cpy#L33) |
| 10 | FILLER | 40 | 41 | X(41) | | | [TRANREC.cpy:34](../../../pangaea/copybooks/TRANREC.cpy#L34) |
