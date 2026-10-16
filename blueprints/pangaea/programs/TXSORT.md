# TXSORT

```text
TXSORT: SORT THE VALID TRANSACTIONS FOR POSTING

POSTING MATCHES TRANSACTIONS AGAINST THE ACCOUNT MASTER, SO
IT NEEDS THEM IN ACCOUNT NUMBER ORDER. WITHIN AN ACCOUNT THEY
KEEP THE ORDER OF THEIR TRANSACTION IDS.

RETURN CODES: 0 SORTED, 16 SORT FAILED.
```

Program [TXSORT.cbl](../../../pangaea/programs/TXSORT.cbl): 3 files, 1 paragraphs.

## Files

| File | Kind | Assigned to | Organization | Status field | Records |
| --- | --- | --- | --- | --- | --- |
| TRAN-IN | FD | TRANIN | SEQUENTIAL | | TRAN-IN-RECORD |
| TRAN-OUT | FD | TRANOUT | SEQUENTIAL | | TRAN-OUT-RECORD |
| SORT-WORK | SD | SORTWK | SEQUENTIAL | | SORT-RECORD |

## Paragraphs

| Paragraph | Description | Source | Performs | Performed by |
| --- | --- | --- | --- | --- |
| 0000-MAIN | | [TXSORT.cbl:36](../../../pangaea/programs/TXSORT.cbl#L36) | | |

## Records

### TRAN-IN-RECORD

File TRAN-IN. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | TRAN-IN-RECORD | 1 | 80 | X(80) | | | [TXSORT.cbl:25](../../../pangaea/programs/TXSORT.cbl#L25) |

### TRAN-OUT-RECORD

File TRAN-OUT. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | TRAN-OUT-RECORD | 1 | 80 | X(80) | | | [TXSORT.cbl:27](../../../pangaea/programs/TXSORT.cbl#L27) |

### SORT-RECORD

File SORT-WORK. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | SORT-RECORD | 1 | 80 | | | | [TXSORT.cbl:29](../../../pangaea/programs/TXSORT.cbl#L29) |
| 05 | SR-REC-TYPE | 1 | 1 | X | | | [TXSORT.cbl:30](../../../pangaea/programs/TXSORT.cbl#L30) |
| 05 | SR-TRAN-ID | 2 | 10 | 9(10) | | | [TXSORT.cbl:31](../../../pangaea/programs/TXSORT.cbl#L31) |
| 05 | SR-ACCT-NO | 12 | 10 | 9(10) | | | [TXSORT.cbl:32](../../../pangaea/programs/TXSORT.cbl#L32) |
| 05 | FILLER | 22 | 59 | X(59) | | | [TXSORT.cbl:33](../../../pangaea/programs/TXSORT.cbl#L33) |
