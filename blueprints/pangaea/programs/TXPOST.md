# TXPOST

Program [TXPOST.cbl](../../../pangaea/programs/TXPOST.cbl): 7 files, 16 paragraphs.

## Files

| File | Kind | Assigned to | Organization | Status field | Records |
| --- | --- | --- | --- | --- | --- |
| CONTROL-FILE | FD | BANKCTL | SEQUENTIAL | WS-CONTROL-STATUS | BANK-CONTROL-RECORD |
| ACCT-IN | FD | ACCTIN | SEQUENTIAL | WS-ACCT-IN-STATUS | ACCT-MASTER-RECORD |
| TRAN-IN | FD | TRANIN | SEQUENTIAL | WS-TRAN-IN-STATUS | TRAN-RECORD |
| ACCT-OUT | FD | ACCTOUT | SEQUENTIAL | WS-ACCT-OUT-STATUS | ACCT-OUT-RECORD |
| AUDIT-OUT | FD | AUDIT | SEQUENTIAL | WS-AUDIT-STATUS | AUDIT-OUT-RECORD |
| REJECT-OUT | FD | REJECTS | SEQUENTIAL | WS-REJECT-STATUS | REJECT-OUT-RECORD |
| REPORT-OUT | FD | POSTRPT | LINE SEQUENTIAL | WS-REPORT-STATUS | REPORT-LINE |

## Paragraphs

| Paragraph | Source | Performs | Performed by |
| --- | --- | --- | --- |
| 0000-MAIN | [TXPOST.cbl:117](../../../pangaea/programs/TXPOST.cbl#L117) | 1000-INITIALIZE, 2000-MATCH, 5000-RECONCILE, 8000-WRITE-REPORT, 9000-CLOSE-FILES | |
| 1000-INITIALIZE | [TXPOST.cbl:132](../../../pangaea/programs/TXPOST.cbl#L132) | 9900-ABORT, 1200-READ-MASTER, 1300-READ-TRAN | 0000-MAIN |
| 1200-READ-MASTER | [TXPOST.cbl:156](../../../pangaea/programs/TXPOST.cbl#L156) | 9900-ABORT | 1000-INITIALIZE, 2000-MATCH |
| 1300-READ-TRAN | [TXPOST.cbl:172](../../../pangaea/programs/TXPOST.cbl#L172) | 9900-ABORT | 1000-INITIALIZE, 2000-MATCH |
| 2000-MATCH | [TXPOST.cbl:187](../../../pangaea/programs/TXPOST.cbl#L187) | 3900-REJECT-TRAN, 1300-READ-TRAN, 3000-APPLY-TRAN, 4000-WRITE-MASTER, 1200-READ-MASTER | 0000-MAIN |
| 3000-APPLY-TRAN | [TXPOST.cbl:201](../../../pangaea/programs/TXPOST.cbl#L201) | 3900-REJECT-TRAN, 3100-POST-DEPOSIT, 3300-POST-FEE, 3200-POST-WITHDRAWAL | 2000-MATCH |
| 3100-POST-DEPOSIT | [TXPOST.cbl:218](../../../pangaea/programs/TXPOST.cbl#L218) | 3800-WRITE-AUDIT | 3000-APPLY-TRAN |
| 3200-POST-WITHDRAWAL | [TXPOST.cbl:224](../../../pangaea/programs/TXPOST.cbl#L224) | 3900-REJECT-TRAN, 3800-WRITE-AUDIT | 3000-APPLY-TRAN |
| 3300-POST-FEE | [TXPOST.cbl:237](../../../pangaea/programs/TXPOST.cbl#L237) | 3800-WRITE-AUDIT | 3000-APPLY-TRAN |
| 3800-WRITE-AUDIT | [TXPOST.cbl:243](../../../pangaea/programs/TXPOST.cbl#L243) | | 3100-POST-DEPOSIT, 3200-POST-WITHDRAWAL, 3300-POST-FEE |
| 3900-REJECT-TRAN | [TXPOST.cbl:255](../../../pangaea/programs/TXPOST.cbl#L255) | | 2000-MATCH, 3000-APPLY-TRAN, 3200-POST-WITHDRAWAL |
| 4000-WRITE-MASTER | [TXPOST.cbl:261](../../../pangaea/programs/TXPOST.cbl#L261) | | 2000-MATCH |
| 5000-RECONCILE | [TXPOST.cbl:267](../../../pangaea/programs/TXPOST.cbl#L267) | 9900-ABORT | 0000-MAIN |
| 8000-WRITE-REPORT | [TXPOST.cbl:278](../../../pangaea/programs/TXPOST.cbl#L278) | | 0000-MAIN |
| 9000-CLOSE-FILES | [TXPOST.cbl:320](../../../pangaea/programs/TXPOST.cbl#L320) | | 0000-MAIN, 9900-ABORT |
| 9900-ABORT | [TXPOST.cbl:324](../../../pangaea/programs/TXPOST.cbl#L324) | 9000-CLOSE-FILES | 1000-INITIALIZE, 1200-READ-MASTER, 1300-READ-TRAN, 5000-RECONCILE |

## Records

### BANK-CONTROL-RECORD

File CONTROL-FILE. 80 bytes, layout in copybook [BANKCTL](../copybooks/BANKCTL.md).

### ACCT-MASTER-RECORD

File ACCT-IN. 100 bytes, layout in copybook [ACCTMAST](../copybooks/ACCTMAST.md).

### TRAN-RECORD

File TRAN-IN. 80 bytes, layout in copybook [TRANREC](../copybooks/TRANREC.md).

### ACCT-OUT-RECORD

File ACCT-OUT. 100 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | ACCT-OUT-RECORD | 1 | 100 | X(100) | | | [TXPOST.cbl:55](../../../pangaea/programs/TXPOST.cbl#L55) |

### AUDIT-OUT-RECORD

File AUDIT-OUT. 100 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | AUDIT-OUT-RECORD | 1 | 100 | X(100) | | | [TXPOST.cbl:57](../../../pangaea/programs/TXPOST.cbl#L57) |

### REJECT-OUT-RECORD

File REJECT-OUT. 100 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | REJECT-OUT-RECORD | 1 | 100 | X(100) | | | [TXPOST.cbl:59](../../../pangaea/programs/TXPOST.cbl#L59) |

### REPORT-LINE

File REPORT-OUT. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | REPORT-LINE | 1 | 80 | X(80) | | | [TXPOST.cbl:61](../../../pangaea/programs/TXPOST.cbl#L61) |

### WS-FILE-STATUSES

Working storage. 14 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-FILE-STATUSES | 1 | 14 | | | | [TXPOST.cbl:64](../../../pangaea/programs/TXPOST.cbl#L64) |
| 05 | WS-CONTROL-STATUS | 1 | 2 | XX | | | [TXPOST.cbl:65](../../../pangaea/programs/TXPOST.cbl#L65) |
| 05 | WS-ACCT-IN-STATUS | 3 | 2 | XX | | | [TXPOST.cbl:66](../../../pangaea/programs/TXPOST.cbl#L66) |
| 05 | WS-TRAN-IN-STATUS | 5 | 2 | XX | | | [TXPOST.cbl:67](../../../pangaea/programs/TXPOST.cbl#L67) |
| 05 | WS-ACCT-OUT-STATUS | 7 | 2 | XX | | | [TXPOST.cbl:68](../../../pangaea/programs/TXPOST.cbl#L68) |
| 05 | WS-AUDIT-STATUS | 9 | 2 | XX | | | [TXPOST.cbl:69](../../../pangaea/programs/TXPOST.cbl#L69) |
| 05 | WS-REJECT-STATUS | 11 | 2 | XX | | | [TXPOST.cbl:70](../../../pangaea/programs/TXPOST.cbl#L70) |
| 05 | WS-REPORT-STATUS | 13 | 2 | XX | | | [TXPOST.cbl:71](../../../pangaea/programs/TXPOST.cbl#L71) |

### WS-BUSINESS-DATE

Working storage. 8 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-BUSINESS-DATE | 1 | 8 | 9(8) | | | [TXPOST.cbl:73](../../../pangaea/programs/TXPOST.cbl#L73) |

### WS-ABORT-MESSAGE

Working storage. 40 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-ABORT-MESSAGE | 1 | 40 | X(40) | | | [TXPOST.cbl:74](../../../pangaea/programs/TXPOST.cbl#L74) |

### WS-MASTER-KEY

Working storage. 10 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-MASTER-KEY | 1 | 10 | X(10) | | | [TXPOST.cbl:77](../../../pangaea/programs/TXPOST.cbl#L77) |

### WS-TRAN-KEY

Working storage. 10 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TRAN-KEY | 1 | 10 | X(10) | | | [TXPOST.cbl:78](../../../pangaea/programs/TXPOST.cbl#L78) |

### WS-PREV-MASTER-KEY

Working storage. 10 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-PREV-MASTER-KEY | 1 | 10 | X(10) | | VALUE LOW-VALUES | [TXPOST.cbl:79](../../../pangaea/programs/TXPOST.cbl#L79) |

### WS-TRAN-ORDER-KEY

Working storage. 20 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TRAN-ORDER-KEY | 1 | 20 | | | | [TXPOST.cbl:80](../../../pangaea/programs/TXPOST.cbl#L80) |
| 05 | WS-TOK-ACCT-NO | 1 | 10 | X(10) | | | [TXPOST.cbl:81](../../../pangaea/programs/TXPOST.cbl#L81) |
| 05 | WS-TOK-TRAN-ID | 11 | 10 | X(10) | | | [TXPOST.cbl:82](../../../pangaea/programs/TXPOST.cbl#L82) |

### WS-PREV-TRAN-ORDER-KEY

Working storage. 20 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-PREV-TRAN-ORDER-KEY | 1 | 20 | X(20) | | VALUE LOW-VALUES | [TXPOST.cbl:83](../../../pangaea/programs/TXPOST.cbl#L83) |

### WS-POSTING

Working storage. 50 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-POSTING | 1 | 50 | | | | [TXPOST.cbl:85](../../../pangaea/programs/TXPOST.cbl#L85) |
| 05 | WS-BALANCE-BEFORE | 1 | 13 | S9(11)V99 | | | [TXPOST.cbl:86](../../../pangaea/programs/TXPOST.cbl#L86) |
| 05 | WS-NEW-BALANCE | 14 | 13 | S9(11)V99 | | | [TXPOST.cbl:87](../../../pangaea/programs/TXPOST.cbl#L87) |
| 05 | WS-BALANCE-FLOOR | 27 | 13 | S9(11)V99 | | | [TXPOST.cbl:88](../../../pangaea/programs/TXPOST.cbl#L88) |
| 05 | WS-SIGNED-AMOUNT | 40 | 11 | S9(9)V99 | | | [TXPOST.cbl:89](../../../pangaea/programs/TXPOST.cbl#L89) |

### WS-TOTALS

Working storage. 122 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TOTALS | 1 | 122 | | | | [TXPOST.cbl:91](../../../pangaea/programs/TXPOST.cbl#L91) |
| 05 | WS-MASTERS-READ | 1 | 8 | 9(8) | | VALUE 0 | [TXPOST.cbl:92](../../../pangaea/programs/TXPOST.cbl#L92) |
| 05 | WS-MASTERS-WRITTEN | 9 | 8 | 9(8) | | VALUE 0 | [TXPOST.cbl:93](../../../pangaea/programs/TXPOST.cbl#L93) |
| 05 | WS-POSTED-COUNT | 17 | 8 | 9(8) | | VALUE 0 | [TXPOST.cbl:94](../../../pangaea/programs/TXPOST.cbl#L94) |
| 05 | WS-REJECT-COUNT | 25 | 8 | 9(8) | | VALUE 0 | [TXPOST.cbl:95](../../../pangaea/programs/TXPOST.cbl#L95) |
| 05 | WS-OPENING-TOTAL | 33 | 15 | S9(13)V99 | | VALUE 0 | [TXPOST.cbl:96](../../../pangaea/programs/TXPOST.cbl#L96) |
| 05 | WS-DEPOSIT-TOTAL | 48 | 15 | 9(13)V99 | | VALUE 0 | [TXPOST.cbl:97](../../../pangaea/programs/TXPOST.cbl#L97) |
| 05 | WS-WITHDRAWAL-TOTAL | 63 | 15 | 9(13)V99 | | VALUE 0 | [TXPOST.cbl:98](../../../pangaea/programs/TXPOST.cbl#L98) |
| 05 | WS-FEE-TOTAL | 78 | 15 | 9(13)V99 | | VALUE 0 | [TXPOST.cbl:99](../../../pangaea/programs/TXPOST.cbl#L99) |
| 05 | WS-CLOSING-TOTAL | 93 | 15 | S9(13)V99 | | VALUE 0 | [TXPOST.cbl:100](../../../pangaea/programs/TXPOST.cbl#L100) |
| 05 | WS-EXPECTED-TOTAL | 108 | 15 | S9(13)V99 | | VALUE 0 | [TXPOST.cbl:101](../../../pangaea/programs/TXPOST.cbl#L101) |

### AUDIT-RECORD

Working storage. 100 bytes, layout in copybook [AUDITREC](../copybooks/AUDITREC.md).

### REJECT-RECORD

Working storage. 100 bytes, layout in copybook [REJREC](../copybooks/REJREC.md).

### WS-REPORT-COUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-COUNT-LINE | 1 | 80 | | | | [TXPOST.cbl:107](../../../pangaea/programs/TXPOST.cbl#L107) |
| 05 | WS-RC-LABEL | 1 | 30 | X(30) | | | [TXPOST.cbl:108](../../../pangaea/programs/TXPOST.cbl#L108) |
| 05 | WS-RC-VALUE | 31 | 10 | ZZ,ZZZ,ZZ9 | | | [TXPOST.cbl:109](../../../pangaea/programs/TXPOST.cbl#L109) |
| 05 | FILLER | 41 | 40 | X(40) | | VALUE SPACES | [TXPOST.cbl:110](../../../pangaea/programs/TXPOST.cbl#L110) |

### WS-REPORT-AMOUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-AMOUNT-LINE | 1 | 80 | | | | [TXPOST.cbl:111](../../../pangaea/programs/TXPOST.cbl#L111) |
| 05 | WS-RA-LABEL | 1 | 30 | X(30) | | | [TXPOST.cbl:112](../../../pangaea/programs/TXPOST.cbl#L112) |
| 05 | WS-RA-VALUE | 31 | 21 | Z,ZZZ,ZZZ,ZZZ,ZZ9.99- | | | [TXPOST.cbl:113](../../../pangaea/programs/TXPOST.cbl#L113) |
| 05 | FILLER | 52 | 29 | X(29) | | VALUE SPACES | [TXPOST.cbl:114](../../../pangaea/programs/TXPOST.cbl#L114) |
