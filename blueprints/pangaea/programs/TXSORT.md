# TXSORT

```text
TXSORT: SORT THE VALID TRANSACTIONS FOR POSTING

POSTING MATCHES TRANSACTIONS AGAINST THE ACCOUNT MASTER, SO
IT NEEDS THEM IN ACCOUNT NUMBER ORDER. WITHIN AN ACCOUNT THEY
KEEP THE ORDER OF THEIR TRANSACTION IDS.

RETURN CODES: 0 SORTED, 16 SORT FAILED.

MODIFICATION HISTORY
DATE       BY  CHANGE
1987-03-12 JA  ORIGINAL VERSION.
1993-02-08 PK  CR-0142 KEEP TRANSACTION ID ORDER WITHIN AN
               ACCOUNT.
```

Program [TXSORT.cbl](../../../pangaea/programs/TXSORT.cbl): 3 files, 1 paragraphs.

- AUTHOR: J. ARKWRIGHT, BATCH SYSTEMS.
- INSTALLATION: FIRST BANK OF PANGAEA, DATA CENTRE.
- DATE-WRITTEN: 1987-03-12.

## Files

| File | Kind | Assigned to | Organization | Status field | Records |
| --- | --- | --- | --- | --- | --- |
| TRAN-IN | FD | TRANIN | SEQUENTIAL | | TRAN-IN-RECORD |
| TRAN-OUT | FD | TRANOUT | SEQUENTIAL | | TRAN-OUT-RECORD |
| SORT-WORK | SD | SORTWK | SEQUENTIAL | | SORT-RECORD |

## Paragraphs

| Paragraph | Description | Source | Performs | Performed by |
| --- | --- | --- | --- | --- |
| 0000-MAIN | SORT BY ACCOUNT NUMBER, THEN TRANSACTION ID. | [TXSORT.cbl:46](../../../pangaea/programs/TXSORT.cbl#L46) | | |

## Records

### TRAN-IN-RECORD

File TRAN-IN. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | TRAN-IN-RECORD | 1 | 80 | X(80) | | | [TXSORT.cbl:34](../../../pangaea/programs/TXSORT.cbl#L34) |

### TRAN-OUT-RECORD

File TRAN-OUT. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | TRAN-OUT-RECORD | 1 | 80 | X(80) | | | [TXSORT.cbl:36](../../../pangaea/programs/TXSORT.cbl#L36) |

### SORT-RECORD

File SORT-WORK. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | SORT-RECORD | 1 | 80 | | | | [TXSORT.cbl:38](../../../pangaea/programs/TXSORT.cbl#L38) |
| 05 | SR-REC-TYPE | 1 | 1 | X | | | [TXSORT.cbl:39](../../../pangaea/programs/TXSORT.cbl#L39) |
| 05 | SR-TRAN-ID | 2 | 10 | 9(10) | | | [TXSORT.cbl:40](../../../pangaea/programs/TXSORT.cbl#L40) |
| 05 | SR-ACCT-NO | 12 | 10 | 9(10) | | | [TXSORT.cbl:41](../../../pangaea/programs/TXSORT.cbl#L41) |
| 05 | FILLER | 22 | 59 | X(59) | | | [TXSORT.cbl:42](../../../pangaea/programs/TXSORT.cbl#L42) |
