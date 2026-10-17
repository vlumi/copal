# AUTHMSG

```text
AUTHMSG: AUTHORIZATION REQUEST, 80 BYTES, AND RESPONSE, 90
A CARD TERMINAL OR ATM SENDS A REQUEST FOR A CARD: A PURCHASE,
THE REVERSAL OF ONE BY ITS AUTHORIZATION CODE, A CASH
WITHDRAWAL OR A BALANCE ENQUIRY. THE RESPONSE GIVES THE
RESULT, AND WHEN APPROVED THE AUTHORIZATION CODE, THE AMOUNT
HELD AND THE FEE, AND THE ACCOUNT'S BALANCE AND AVAILABLE
BALANCE AFTERWARDS.

MODIFICATION HISTORY
DATE       BY  CHANGE
2003-05-19 MV  CR-0342 ORIGINAL VERSION.
2007-02-05 MV  CR-0415 FEE IN THE RESPONSE.
```

Copybook [AUTHMSG.cpy](../../../pangaea/copybooks/AUTHMSG.cpy), used by [ATMBAL](../programs/ATMBAL.md), [ATMWDL](../programs/ATMWDL.md), [AUTHPUR](../programs/AUTHPUR.md), [AUTHREV](../programs/AUTHREV.md), [AUTHSRV](../programs/AUTHSRV.md).

## AUTH-REQUEST

80 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | AUTH-REQUEST | group | 1 | 80 | | | | [AUTHMSG.cpy:15](../../../pangaea/copybooks/AUTHMSG.cpy#L15) |
| 05 | RQ-FUNCTION | code | 1 | 4 | X(4) | | | [AUTHMSG.cpy:16](../../../pangaea/copybooks/AUTHMSG.cpy#L16) |
| 88 | RQ-PURCHASE | | | | | | VALUE 'PURC' | [AUTHMSG.cpy:17](../../../pangaea/copybooks/AUTHMSG.cpy#L17) |
| 88 | RQ-REVERSAL | | | | | | VALUE 'REVS' | [AUTHMSG.cpy:18](../../../pangaea/copybooks/AUTHMSG.cpy#L18) |
| 88 | RQ-WITHDRAWAL | | | | | | VALUE 'WDRL' | [AUTHMSG.cpy:19](../../../pangaea/copybooks/AUTHMSG.cpy#L19) |
| 88 | RQ-BALANCE-ENQUIRY | | | | | | VALUE 'BALQ' | [AUTHMSG.cpy:20](../../../pangaea/copybooks/AUTHMSG.cpy#L20) |
| 05 | RQ-CARD-NO | number | 5 | 16 | 9(16) | | | [AUTHMSG.cpy:21](../../../pangaea/copybooks/AUTHMSG.cpy#L21) |
| 05 | RQ-AMOUNT | amount | 21 | 9 | 9(7)V99 | | | [AUTHMSG.cpy:22](../../../pangaea/copybooks/AUTHMSG.cpy#L22) |
| 05 | RQ-AUTH-CODE | number | 30 | 6 | 9(6) | | | [AUTHMSG.cpy:23](../../../pangaea/copybooks/AUTHMSG.cpy#L23) |
| 05 | RQ-TERMINAL | text | 36 | 20 | X(20) | | | [AUTHMSG.cpy:24](../../../pangaea/copybooks/AUTHMSG.cpy#L24) |
| 05 | RQ-TIMESTAMP | date | 56 | 14 | 9(14) | | | [AUTHMSG.cpy:25](../../../pangaea/copybooks/AUTHMSG.cpy#L25) |
| 05 | FILLER | filler | 70 | 11 | X(11) | | | [AUTHMSG.cpy:26](../../../pangaea/copybooks/AUTHMSG.cpy#L26) |

## AUTH-RESPONSE

90 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | AUTH-RESPONSE | group | 1 | 90 | | | | [AUTHMSG.cpy:27](../../../pangaea/copybooks/AUTHMSG.cpy#L27) |
| 05 | RS-RESULT | code | 1 | 2 | XX | | | [AUTHMSG.cpy:28](../../../pangaea/copybooks/AUTHMSG.cpy#L28) |
| 88 | RS-APPROVED | | | | | | VALUE '00' | [AUTHMSG.cpy:29](../../../pangaea/copybooks/AUTHMSG.cpy#L29) |
| 88 | RS-INVALID-REQUEST | | | | | | VALUE '12' | [AUTHMSG.cpy:30](../../../pangaea/copybooks/AUTHMSG.cpy#L30) |
| 88 | RS-INVALID-AMOUNT | | | | | | VALUE '13' | [AUTHMSG.cpy:31](../../../pangaea/copybooks/AUTHMSG.cpy#L31) |
| 88 | RS-NO-SUCH-CARD | | | | | | VALUE '14' | [AUTHMSG.cpy:32](../../../pangaea/copybooks/AUTHMSG.cpy#L32) |
| 88 | RS-NO-SUCH-AUTH | | | | | | VALUE '25' | [AUTHMSG.cpy:33](../../../pangaea/copybooks/AUTHMSG.cpy#L33) |
| 88 | RS-INSUFFICIENT-FUNDS | | | | | | VALUE '51' | [AUTHMSG.cpy:34](../../../pangaea/copybooks/AUTHMSG.cpy#L34) |
| 88 | RS-OVER-LIMIT | | | | | | VALUE '61' | [AUTHMSG.cpy:35](../../../pangaea/copybooks/AUTHMSG.cpy#L35) |
| 88 | RS-CARD-BLOCKED | | | | | | VALUE '62' | [AUTHMSG.cpy:36](../../../pangaea/copybooks/AUTHMSG.cpy#L36) |
| 88 | RS-ACCOUNT-NOT-OPEN | | | | | | VALUE '78' | [AUTHMSG.cpy:37](../../../pangaea/copybooks/AUTHMSG.cpy#L37) |
| 05 | FILLER | filler | 3 | 1 | X | | | [AUTHMSG.cpy:38](../../../pangaea/copybooks/AUTHMSG.cpy#L38) |
| 05 | RS-AUTH-CODE | number | 4 | 6 | 9(6) | | | [AUTHMSG.cpy:39](../../../pangaea/copybooks/AUTHMSG.cpy#L39) |
| 05 | FILLER | filler | 10 | 1 | X | | | [AUTHMSG.cpy:40](../../../pangaea/copybooks/AUTHMSG.cpy#L40) |
| 05 | RS-AMOUNT | edited | 11 | 10 | Z(6)9.99 | | | [AUTHMSG.cpy:41](../../../pangaea/copybooks/AUTHMSG.cpy#L41) |
| 05 | FILLER | filler | 21 | 1 | X | | | [AUTHMSG.cpy:42](../../../pangaea/copybooks/AUTHMSG.cpy#L42) |
| 05 | RS-FEE | edited | 22 | 6 | ZZ9.99 | | | [AUTHMSG.cpy:43](../../../pangaea/copybooks/AUTHMSG.cpy#L43) |
| 05 | FILLER | filler | 28 | 1 | X | | | [AUTHMSG.cpy:44](../../../pangaea/copybooks/AUTHMSG.cpy#L44) |
| 05 | RS-BALANCE | edited | 29 | 15 | -9(11).99 | | | [AUTHMSG.cpy:45](../../../pangaea/copybooks/AUTHMSG.cpy#L45) |
| 05 | FILLER | filler | 44 | 1 | X | | | [AUTHMSG.cpy:46](../../../pangaea/copybooks/AUTHMSG.cpy#L46) |
| 05 | RS-AVAILABLE | edited | 45 | 15 | -9(11).99 | | | [AUTHMSG.cpy:47](../../../pangaea/copybooks/AUTHMSG.cpy#L47) |
| 05 | FILLER | filler | 60 | 1 | X | | | [AUTHMSG.cpy:48](../../../pangaea/copybooks/AUTHMSG.cpy#L48) |
| 05 | RS-MESSAGE | text | 61 | 30 | X(30) | | | [AUTHMSG.cpy:49](../../../pangaea/copybooks/AUTHMSG.cpy#L49) |
