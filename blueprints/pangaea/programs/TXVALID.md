# TXVALID

```text
TXVALID: VALIDATE THE DAY'S TRANSACTION FILE

CHECKS THE FILE AS A WHOLE: A HEADER FOR TODAY'S BUSINESS
DATE, THEN DETAILS, THEN A TRAILER WHOSE CONTROL TOTALS MATCH
THE DETAILS. A FILE THAT FAILS THESE CHECKS IS NOT PROCESSED.
CHECKS EACH DETAIL'S FIELDS AND WRITES IT EITHER TO THE VALID
FILE OR TO THE REJECTS. WHETHER THE ACCOUNT EXISTS IS LEFT TO
POSTING, WHICH READS THE MASTER.

RETURN CODES: 0 ALL VALID, 4 SOME DETAILS REJECTED,
              12 FILE REJECTED, 16 I/O ERROR.
```

Program [TXVALID.cbl](../../../pangaea/programs/TXVALID.cbl): 5 files, 10 paragraphs.

## Files

| File | Kind | Assigned to | Organization | Status field | Records |
| --- | --- | --- | --- | --- | --- |
| CONTROL-FILE | FD | BANKCTL | SEQUENTIAL | WS-CONTROL-STATUS | BANK-CONTROL-RECORD |
| TRAN-IN | FD | TRANIN | SEQUENTIAL | WS-TRAN-IN-STATUS | TRAN-RECORD |
| TRAN-OUT | FD | TRANOUT | SEQUENTIAL | WS-TRAN-OUT-STATUS | TRAN-OUT-RECORD |
| REJECT-OUT | FD | REJECTS | SEQUENTIAL | WS-REJECT-STATUS | REJECT-OUT-RECORD |
| REPORT-OUT | FD | VALRPT | LINE SEQUENTIAL | WS-REPORT-STATUS | REPORT-LINE |

## Paragraphs

| Paragraph | Description | Source | Performs | Performed by |
| --- | --- | --- | --- | --- |
| 0000-MAIN | | [TXVALID.cbl:86](../../../pangaea/programs/TXVALID.cbl#L86) | 1000-OPEN-FILES, 2000-READ-HEADER, 3000-READ-DETAILS, 4000-CHECK-TRAILER, 8000-WRITE-REPORT, 9000-CLOSE-FILES | |
| 1000-OPEN-FILES | | [TXVALID.cbl:107](../../../pangaea/programs/TXVALID.cbl#L107) | | 0000-MAIN |
| 2000-READ-HEADER | | [TXVALID.cbl:128](../../../pangaea/programs/TXVALID.cbl#L128) | | 0000-MAIN |
| 3000-READ-DETAILS | | [TXVALID.cbl:143](../../../pangaea/programs/TXVALID.cbl#L143) | 3100-READ-NEXT, 3200-VALIDATE-DETAIL, 3300-REJECT-RECORD | 0000-MAIN |
| 3100-READ-NEXT | | [TXVALID.cbl:160](../../../pangaea/programs/TXVALID.cbl#L160) | | 3000-READ-DETAILS |
| 3200-VALIDATE-DETAIL | | [TXVALID.cbl:166](../../../pangaea/programs/TXVALID.cbl#L166) | 3300-REJECT-RECORD | 3000-READ-DETAILS |
| 3300-REJECT-RECORD | | [TXVALID.cbl:194](../../../pangaea/programs/TXVALID.cbl#L194) | | 3000-READ-DETAILS, 3200-VALIDATE-DETAIL |
| 4000-CHECK-TRAILER | | [TXVALID.cbl:200](../../../pangaea/programs/TXVALID.cbl#L200) | | 0000-MAIN |
| 8000-WRITE-REPORT | | [TXVALID.cbl:220](../../../pangaea/programs/TXVALID.cbl#L220) | | 0000-MAIN |
| 9000-CLOSE-FILES | | [TXVALID.cbl:251](../../../pangaea/programs/TXVALID.cbl#L251) | | 0000-MAIN |

## Records

### BANK-CONTROL-RECORD

File CONTROL-FILE. 80 bytes, layout in copybook [BANKCTL](../copybooks/BANKCTL.md).

### TRAN-RECORD

File TRAN-IN. 80 bytes, layout in copybook [TRANREC](../copybooks/TRANREC.md).

### TRAN-OUT-RECORD

File TRAN-OUT. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | TRAN-OUT-RECORD | 1 | 80 | X(80) | | | [TXVALID.cbl:43](../../../pangaea/programs/TXVALID.cbl#L43) |

### REJECT-OUT-RECORD

File REJECT-OUT. 100 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | REJECT-OUT-RECORD | 1 | 100 | X(100) | | | [TXVALID.cbl:45](../../../pangaea/programs/TXVALID.cbl#L45) |

### REPORT-LINE

File REPORT-OUT. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | REPORT-LINE | 1 | 80 | X(80) | | | [TXVALID.cbl:47](../../../pangaea/programs/TXVALID.cbl#L47) |

### WS-FILE-STATUSES

Working storage. 10 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-FILE-STATUSES | 1 | 10 | | | | [TXVALID.cbl:50](../../../pangaea/programs/TXVALID.cbl#L50) |
| 05 | WS-CONTROL-STATUS | 1 | 2 | XX | | | [TXVALID.cbl:51](../../../pangaea/programs/TXVALID.cbl#L51) |
| 05 | WS-TRAN-IN-STATUS | 3 | 2 | XX | | | [TXVALID.cbl:52](../../../pangaea/programs/TXVALID.cbl#L52) |
| 88 | TRAN-IN-EOF | | | | | VALUE '10' | [TXVALID.cbl:53](../../../pangaea/programs/TXVALID.cbl#L53) |
| 05 | WS-TRAN-OUT-STATUS | 5 | 2 | XX | | | [TXVALID.cbl:54](../../../pangaea/programs/TXVALID.cbl#L54) |
| 05 | WS-REJECT-STATUS | 7 | 2 | XX | | | [TXVALID.cbl:55](../../../pangaea/programs/TXVALID.cbl#L55) |
| 05 | WS-REPORT-STATUS | 9 | 2 | XX | | | [TXVALID.cbl:56](../../../pangaea/programs/TXVALID.cbl#L56) |

### WS-BUSINESS-DATE

Working storage. 8 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-BUSINESS-DATE | 1 | 8 | 9(8) | | | [TXVALID.cbl:58](../../../pangaea/programs/TXVALID.cbl#L58) |

### REJECT-RECORD

Working storage. 100 bytes, layout in copybook [REJREC](../copybooks/REJREC.md).

### WS-FILE-ERROR

Working storage. 40 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-FILE-ERROR | 1 | 40 | X(40) | | VALUE SPACES | [TXVALID.cbl:62](../../../pangaea/programs/TXVALID.cbl#L62) |
| 88 | FILE-IS-GOOD | | | | | VALUE SPACES | [TXVALID.cbl:63](../../../pangaea/programs/TXVALID.cbl#L63) |

### WS-DETAIL-STATE

Working storage. 1 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-DETAIL-STATE | 1 | 1 | X | | | [TXVALID.cbl:65](../../../pangaea/programs/TXVALID.cbl#L65) |
| 88 | DETAIL-IS-VALID | | | | | VALUE 'Y' | [TXVALID.cbl:66](../../../pangaea/programs/TXVALID.cbl#L66) |
| 88 | DETAIL-IS-REJECTED | | | | | VALUE 'N' | [TXVALID.cbl:67](../../../pangaea/programs/TXVALID.cbl#L67) |

### WS-TOTALS

Working storage. 54 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TOTALS | 1 | 54 | | | | [TXVALID.cbl:69](../../../pangaea/programs/TXVALID.cbl#L69) |
| 05 | WS-DETAIL-COUNT | 1 | 8 | 9(8) | | VALUE 0 | [TXVALID.cbl:70](../../../pangaea/programs/TXVALID.cbl#L70) |
| 05 | WS-VALID-COUNT | 9 | 8 | 9(8) | | VALUE 0 | [TXVALID.cbl:71](../../../pangaea/programs/TXVALID.cbl#L71) |
| 05 | WS-REJECT-COUNT | 17 | 8 | 9(8) | | VALUE 0 | [TXVALID.cbl:72](../../../pangaea/programs/TXVALID.cbl#L72) |
| 05 | WS-AMOUNT-TOTAL | 25 | 15 | 9(13)V99 | | VALUE 0 | [TXVALID.cbl:73](../../../pangaea/programs/TXVALID.cbl#L73) |
| 05 | WS-ACCT-HASH | 40 | 15 | 9(15) | | VALUE 0 | [TXVALID.cbl:74](../../../pangaea/programs/TXVALID.cbl#L74) |

### WS-REPORT-COUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-COUNT-LINE | 1 | 80 | | | | [TXVALID.cbl:76](../../../pangaea/programs/TXVALID.cbl#L76) |
| 05 | WS-RC-LABEL | 1 | 30 | X(30) | | | [TXVALID.cbl:77](../../../pangaea/programs/TXVALID.cbl#L77) |
| 05 | WS-RC-VALUE | 31 | 10 | ZZ,ZZZ,ZZ9 | | | [TXVALID.cbl:78](../../../pangaea/programs/TXVALID.cbl#L78) |
| 05 | FILLER | 41 | 40 | X(40) | | VALUE SPACES | [TXVALID.cbl:79](../../../pangaea/programs/TXVALID.cbl#L79) |

### WS-REPORT-AMOUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-AMOUNT-LINE | 1 | 80 | | | | [TXVALID.cbl:80](../../../pangaea/programs/TXVALID.cbl#L80) |
| 05 | WS-RA-LABEL | 1 | 30 | X(30) | | | [TXVALID.cbl:81](../../../pangaea/programs/TXVALID.cbl#L81) |
| 05 | WS-RA-VALUE | 31 | 20 | Z,ZZZ,ZZZ,ZZZ,ZZ9.99 | | | [TXVALID.cbl:82](../../../pangaea/programs/TXVALID.cbl#L82) |
| 05 | FILLER | 51 | 30 | X(30) | | VALUE SPACES | [TXVALID.cbl:83](../../../pangaea/programs/TXVALID.cbl#L83) |
