# TXVALID

```text
TXVALID: VALIDATE THE DAY'S TRANSACTION FILE

CHECKS THE FILE AS A WHOLE: A HEADER FOR TODAY'S BUSINESS
DATE, THEN DETAILS, THEN A TRAILER WHOSE CONTROL TOTALS MATCH
THE DETAILS. A FILE THAT FAILS THESE CHECKS IS NOT PROCESSED.
CHECKS EACH DETAIL'S FIELDS AND WRITES IT EITHER TO THE VALID
FILE OR TO THE REJECTS: THE ACCOUNT NUMBER, THE CODE, THE
AMOUNT, AND FOR A TRANSFER THE ACCOUNT IT GOES TO. WHETHER
THE ACCOUNTS EXIST IS LEFT TO POSTING, WHICH READS THE MASTER.

RETURN CODES: 0 ALL VALID, 4 SOME DETAILS REJECTED,
              12 FILE REJECTED, 16 I/O ERROR.

MODIFICATION HISTORY
DATE       BY  CHANGE
1987-03-12 JA  ORIGINAL VERSION.
1991-06-03 PK  CR-0117 ACCOUNT NUMBER HASH IN THE TRAILER.
1996-04-15 PK  CR-0233 TRANSFERS MUST NAME ANOTHER ACCOUNT.
1999-01-11 SL  Y2K: 8-DIGIT BUSINESS DATE FROM THE BANK
               CONTROL RECORD.
2004-09-20 PK  CR-0388 REJECT RECORDS OF UNKNOWN TYPE
               INSTEAD OF STOPPING THE RUN.
```

Program [TXVALID.cbl](../../../pangaea/programs/TXVALID.cbl): 5 files, 10 paragraphs.

- AUTHOR: J. ARKWRIGHT, BATCH SYSTEMS.
- INSTALLATION: FIRST BANK OF PANGAEA, DATA CENTER.
- DATE-WRITTEN: 1987-03-12.

## Files

| File | Kind | Assigned to | Organization | Status field | Records |
| --- | --- | --- | --- | --- | --- |
| CONTROL-FILE | FD | BANKCTL | SEQUENTIAL | WS-CONTROL-STATUS | BANK-CONTROL-RECORD |
| TRAN-IN | FD | TRANIN | SEQUENTIAL | WS-TRAN-IN-STATUS | TRAN-RECORD |
| TRAN-OUT | FD | TRANOUT | SEQUENTIAL | WS-TRAN-OUT-STATUS | TRAN-OUT-RECORD |
| REJECT-OUT | FD | REJECTS | SEQUENTIAL | WS-REJECT-STATUS | REJECT-OUT-RECORD |
| REPORT-OUT | FD | VALRPT | LINE SEQUENTIAL | WS-REPORT-STATUS | REPORT-LINE |

## Flow

![The paragraphs of TXVALID and how control moves between them](TXVALID.svg)

Solid arrows are PERFORMs, dashed ones GO TOs, and dotted ones control running on into the next paragraph inside a PERFORM THRU range.

## Paragraphs

| Paragraph | Description | Source | Performs | Goes to | Performed by | Reached by |
| --- | --- | --- | --- | --- | --- | --- |
| 0000-MAIN | VALIDATE THE FILE AND SET THE RETURN CODE. | [TXVALID.cbl:101](../../../pangaea/programs/TXVALID.cbl#L101) | 1000-OPEN-FILES, 2000-READ-HEADER, 3000-READ-DETAILS, 4000-CHECK-TRAILER, 8000-WRITE-REPORT, 9000-CLOSE-FILES | | | |
| 1000-OPEN-FILES | OPEN THE FILES AND READ THE BUSINESS DATE. | [TXVALID.cbl:123](../../../pangaea/programs/TXVALID.cbl#L123) | | | 0000-MAIN | |
| 2000-READ-HEADER | THE FIRST RECORD MUST BE THE HEADER FOR THE BUSINESS DATE. | [TXVALID.cbl:145](../../../pangaea/programs/TXVALID.cbl#L145) | | | 0000-MAIN | |
| 3000-READ-DETAILS | VALIDATE EACH RECORD UNTIL THE TRAILER. | [TXVALID.cbl:161](../../../pangaea/programs/TXVALID.cbl#L161) | 3100-READ-NEXT, 3200-VALIDATE-DETAIL, 3300-REJECT-RECORD | | 0000-MAIN | |
| 3100-READ-NEXT | READ ONE RECORD. THE FILE MUST NOT END BEFORE ITS TRAILER. | [TXVALID.cbl:179](../../../pangaea/programs/TXVALID.cbl#L179) | | | 3000-READ-DETAILS | |
| 3200-VALIDATE-DETAIL | ADD A DETAIL TO THE TOTALS, THEN CHECK ITS FIELDS. | [TXVALID.cbl:186](../../../pangaea/programs/TXVALID.cbl#L186) | 3300-REJECT-RECORD | | 3000-READ-DETAILS | |
| 3300-REJECT-RECORD | WRITE THE RECORD TO THE REJECTS WITH ITS REASON. | [TXVALID.cbl:220](../../../pangaea/programs/TXVALID.cbl#L220) | | | 3000-READ-DETAILS, 3200-VALIDATE-DETAIL | |
| 4000-CHECK-TRAILER | THE TRAILER MUST MATCH THE TOTALS, AND END THE FILE. | [TXVALID.cbl:227](../../../pangaea/programs/TXVALID.cbl#L227) | | | 0000-MAIN | |
| 8000-WRITE-REPORT | WRITE THE VALIDATION REPORT. | [TXVALID.cbl:248](../../../pangaea/programs/TXVALID.cbl#L248) | | | 0000-MAIN | |
| 9000-CLOSE-FILES | CLOSE ALL FILES. | [TXVALID.cbl:280](../../../pangaea/programs/TXVALID.cbl#L280) | | | 0000-MAIN | |

## Records

### BANK-CONTROL-RECORD

File CONTROL-FILE. 80 bytes, layout in copybook [BANKCTL](../copybooks/BANKCTL.md).

### TRAN-RECORD

File TRAN-IN. 80 bytes, layout in copybook [TRANREC](../copybooks/TRANREC.md).

### TRAN-OUT-RECORD

File TRAN-OUT. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | TRAN-OUT-RECORD | 1 | 80 | X(80) | | | [TXVALID.cbl:57](../../../pangaea/programs/TXVALID.cbl#L57) |

### REJECT-OUT-RECORD

File REJECT-OUT. 100 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | REJECT-OUT-RECORD | 1 | 100 | X(100) | | | [TXVALID.cbl:59](../../../pangaea/programs/TXVALID.cbl#L59) |

### REPORT-LINE

File REPORT-OUT. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | REPORT-LINE | 1 | 80 | X(80) | | | [TXVALID.cbl:61](../../../pangaea/programs/TXVALID.cbl#L61) |

### WS-FILE-STATUSES

Working storage. 10 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-FILE-STATUSES | 1 | 10 | | | | [TXVALID.cbl:64](../../../pangaea/programs/TXVALID.cbl#L64) |
| 05 | WS-CONTROL-STATUS | 1 | 2 | XX | | | [TXVALID.cbl:65](../../../pangaea/programs/TXVALID.cbl#L65) |
| 05 | WS-TRAN-IN-STATUS | 3 | 2 | XX | | | [TXVALID.cbl:66](../../../pangaea/programs/TXVALID.cbl#L66) |
| 88 | TRAN-IN-EOF | | | | | VALUE '10' | [TXVALID.cbl:67](../../../pangaea/programs/TXVALID.cbl#L67) |
| 05 | WS-TRAN-OUT-STATUS | 5 | 2 | XX | | | [TXVALID.cbl:68](../../../pangaea/programs/TXVALID.cbl#L68) |
| 05 | WS-REJECT-STATUS | 7 | 2 | XX | | | [TXVALID.cbl:69](../../../pangaea/programs/TXVALID.cbl#L69) |
| 05 | WS-REPORT-STATUS | 9 | 2 | XX | | | [TXVALID.cbl:70](../../../pangaea/programs/TXVALID.cbl#L70) |

### WS-BUSINESS-DATE

Working storage. 8 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-BUSINESS-DATE | 1 | 8 | 9(8) | | | [TXVALID.cbl:72](../../../pangaea/programs/TXVALID.cbl#L72) |

### REJECT-RECORD

Working storage. 100 bytes, layout in copybook [REJREC](../copybooks/REJREC.md).

### WS-FILE-ERROR

Working storage. 40 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-FILE-ERROR | 1 | 40 | X(40) | | VALUE SPACES | [TXVALID.cbl:76](../../../pangaea/programs/TXVALID.cbl#L76) |
| 88 | FILE-IS-GOOD | | | | | VALUE SPACES | [TXVALID.cbl:77](../../../pangaea/programs/TXVALID.cbl#L77) |

### WS-DETAIL-STATE

Working storage. 1 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-DETAIL-STATE | 1 | 1 | X | | | [TXVALID.cbl:79](../../../pangaea/programs/TXVALID.cbl#L79) |
| 88 | DETAIL-IS-VALID | | | | | VALUE 'Y' | [TXVALID.cbl:80](../../../pangaea/programs/TXVALID.cbl#L80) |
| 88 | DETAIL-IS-REJECTED | | | | | VALUE 'N' | [TXVALID.cbl:81](../../../pangaea/programs/TXVALID.cbl#L81) |

### WS-TOTALS

Working storage. 54 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TOTALS | 1 | 54 | | | | [TXVALID.cbl:83](../../../pangaea/programs/TXVALID.cbl#L83) |
| 05 | WS-DETAIL-COUNT | 1 | 8 | 9(8) | | VALUE 0 | [TXVALID.cbl:84](../../../pangaea/programs/TXVALID.cbl#L84) |
| 05 | WS-VALID-COUNT | 9 | 8 | 9(8) | | VALUE 0 | [TXVALID.cbl:85](../../../pangaea/programs/TXVALID.cbl#L85) |
| 05 | WS-REJECT-COUNT | 17 | 8 | 9(8) | | VALUE 0 | [TXVALID.cbl:86](../../../pangaea/programs/TXVALID.cbl#L86) |
| 05 | WS-AMOUNT-TOTAL | 25 | 15 | 9(13)V99 | | VALUE 0 | [TXVALID.cbl:87](../../../pangaea/programs/TXVALID.cbl#L87) |
| 05 | WS-ACCT-HASH | 40 | 15 | 9(15) | | VALUE 0 | [TXVALID.cbl:88](../../../pangaea/programs/TXVALID.cbl#L88) |

### WS-REPORT-COUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-COUNT-LINE | 1 | 80 | | | | [TXVALID.cbl:90](../../../pangaea/programs/TXVALID.cbl#L90) |
| 05 | WS-RC-LABEL | 1 | 30 | X(30) | | | [TXVALID.cbl:91](../../../pangaea/programs/TXVALID.cbl#L91) |
| 05 | WS-RC-VALUE | 31 | 10 | ZZ,ZZZ,ZZ9 | | | [TXVALID.cbl:92](../../../pangaea/programs/TXVALID.cbl#L92) |
| 05 | FILLER | 41 | 40 | X(40) | | VALUE SPACES | [TXVALID.cbl:93](../../../pangaea/programs/TXVALID.cbl#L93) |

### WS-REPORT-AMOUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-AMOUNT-LINE | 1 | 80 | | | | [TXVALID.cbl:94](../../../pangaea/programs/TXVALID.cbl#L94) |
| 05 | WS-RA-LABEL | 1 | 30 | X(30) | | | [TXVALID.cbl:95](../../../pangaea/programs/TXVALID.cbl#L95) |
| 05 | WS-RA-VALUE | 31 | 20 | Z,ZZZ,ZZZ,ZZZ,ZZ9.99 | | | [TXVALID.cbl:96](../../../pangaea/programs/TXVALID.cbl#L96) |
| 05 | FILLER | 51 | 30 | X(30) | | VALUE SPACES | [TXVALID.cbl:97](../../../pangaea/programs/TXVALID.cbl#L97) |
