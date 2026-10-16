# TXSORT

```text
TXSORT: SORT TRANSACTIONS FOR POSTING

POSTING MATCHES TRANSACTIONS AGAINST THE ACCOUNT MASTER, SO
IT NEEDS THEM IN ACCOUNT NUMBER ORDER. WITHIN AN ACCOUNT THEY
KEEP THE ORDER OF THEIR TRANSACTION IDS. THE DAILY BATCH SORTS
THE VALID TRANSACTIONS, THEN THE TRANSFER CREDITS, THEN THE
RETURNED TRANSFERS.

RETURN CODES: 0 SORTED, 16 SORT FAILED.

MODIFICATION HISTORY
DATE       BY  CHANGE
1987-03-12 JA  ORIGINAL VERSION.
1993-02-08 PK  CR-0142 KEEP TRANSACTION ID ORDER WITHIN AN
               ACCOUNT.
```

Program [TXSORT.cbl](../../../pangaea/programs/TXSORT.cbl): 3 files, 1 paragraphs.

- AUTHOR: J. ARKWRIGHT, BATCH SYSTEMS.
- INSTALLATION: FIRST BANK OF PANGAEA, DATA CENTER.
- DATE-WRITTEN: 1987-03-12.

## Files

| File | Kind | Assigned to | Organization | Status field | Opened | Records |
| --- | --- | --- | --- | --- | --- | --- |
| TRAN-IN | FD | TRANIN | SEQUENTIAL | | INPUT | TRAN-IN-RECORD |
| TRAN-OUT | FD | TRANOUT | SEQUENTIAL | | OUTPUT | TRAN-OUT-RECORD |
| SORT-WORK | SD | SORTWK | SEQUENTIAL | | | SORT-RECORD |

## Flow

![The paragraphs of TXSORT and how control moves between them](TXSORT.svg)

Solid arrows are PERFORMs, dashed ones GO TOs, and dotted ones control running on into the next paragraph inside a PERFORM THRU range.

## Paragraphs

| Paragraph | Description | Source | Performs | Goes to | Performed by | Reached by |
| --- | --- | --- | --- | --- | --- | --- |
| 0000-MAIN | SORT BY ACCOUNT NUMBER, THEN TRANSACTION ID. | [TXSORT.cbl:48](../../../pangaea/programs/TXSORT.cbl#L48) | | | | |

## Records

### TRAN-IN-RECORD

File TRAN-IN. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | TRAN-IN-RECORD | 1 | 80 | X(80) | | | [TXSORT.cbl:36](../../../pangaea/programs/TXSORT.cbl#L36) |

### TRAN-OUT-RECORD

File TRAN-OUT. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | TRAN-OUT-RECORD | 1 | 80 | X(80) | | | [TXSORT.cbl:38](../../../pangaea/programs/TXSORT.cbl#L38) |

### SORT-RECORD

File SORT-WORK. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | SORT-RECORD | 1 | 80 | | | | [TXSORT.cbl:40](../../../pangaea/programs/TXSORT.cbl#L40) |
| 05 | SR-REC-TYPE | 1 | 1 | X | | | [TXSORT.cbl:41](../../../pangaea/programs/TXSORT.cbl#L41) |
| 05 | SR-TRAN-ID | 2 | 10 | 9(10) | | | [TXSORT.cbl:42](../../../pangaea/programs/TXSORT.cbl#L42) |
| 05 | SR-ACCT-NO | 12 | 10 | 9(10) | | | [TXSORT.cbl:43](../../../pangaea/programs/TXSORT.cbl#L43) |
| 05 | FILLER | 22 | 59 | X(59) | | | [TXSORT.cbl:44](../../../pangaea/programs/TXSORT.cbl#L44) |
