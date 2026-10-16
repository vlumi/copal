# BANKDATE

```text
BANKDATE: ADVANCE THE BUSINESS DATE

READS THE BANK CONTROL RECORD AND WRITES THE NEXT ONE, FOR
THE NEXT BUSINESS DAY: THE NEXT DAY FROM MONDAY TO FRIDAY.
THE LAST STEP OF THE DAILY BATCH, SO THAT THE NEXT RUN IS FOR
THE NEXT DAY ONLY IF THIS ONE WENT THROUGH.

RETURN CODES: 0 ADVANCED, 16 NO CONTROL RECORD OR I/O.

MODIFICATION HISTORY
DATE       BY  CHANGE
1987-03-12 JA  ORIGINAL VERSION.
1999-01-11 SL  Y2K: 8-DIGIT DATES, INTRINSIC DATE FUNCTIONS
               INSTEAD OF THE CALENDAR TABLE.
```

Program [BANKDATE.cbl](../../../pangaea/programs/BANKDATE.cbl): 2 files, 2 paragraphs.

- AUTHOR: J. ARKWRIGHT, BATCH SYSTEMS.
- INSTALLATION: FIRST BANK OF PANGAEA, DATA CENTER.
- DATE-WRITTEN: 1987-03-12.

## Contents

- [Files](#files)
- [Flow](#flow)
- [Paragraphs](#paragraphs)
- [Data flow](#data-flow)
- [Lineage](#lineage)
- [Records](#records)
  - [BANK-CONTROL-RECORD](#bank-control-record)
  - [CONTROL-OUT-RECORD](#control-out-record)
  - [WS-IN-STATUS](#ws-in-status)
  - [WS-OUT-STATUS](#ws-out-status)
  - [WS-DAY-NUMBER](#ws-day-number)
  - [WS-WEEKDAY](#ws-weekday)

## Files

| File | Kind | Assigned to | Organization | Status field | Opened | Records |
| --- | --- | --- | --- | --- | --- | --- |
| CONTROL-IN | FD | BANKIN | SEQUENTIAL | WS-IN-STATUS | INPUT | BANK-CONTROL-RECORD |
| CONTROL-OUT | FD | BANKOUT | SEQUENTIAL | WS-OUT-STATUS | OUTPUT | CONTROL-OUT-RECORD |

## Flow

![The paragraphs of BANKDATE and how control moves between them](BANKDATE.svg)

Solid arrows are PERFORMs, dashed ones GO TOs, and dotted ones control running on into the next paragraph inside a PERFORM THRU range.

## Paragraphs

| Paragraph | Description | Source | Performs | Goes to | Performed by | Reached by |
| --- | --- | --- | --- | --- | --- | --- |
| 0000-MAIN | READ TODAY'S CONTROL RECORD AND WRITE THE NEXT BUSINESS DAY'S. | [BANKDATE.cbl:52](../../../pangaea/programs/BANKDATE.cbl#L52) | 1000-NEXT-BUSINESS-DAY | | | |
| 1000-NEXT-BUSINESS-DAY | MOVE THE BUSINESS DATE TO THE NEXT DAY THAT IS NOT A WEEKEND. | [BANKDATE.cbl:74](../../../pangaea/programs/BANKDATE.cbl#L74) | | | 0000-MAIN | |

## Data flow

| Paragraph | Reads | Writes |
| --- | --- | --- |
| 0000-MAIN | WS-IN-STATUS, WS-OUT-STATUS, BANK-CONTROL-RECORD | RETURN-CODE, BANK-CONTROL-RECORD, CONTROL-OUT-RECORD |
| 1000-NEXT-BUSINESS-DAY | BC-BUSINESS-DATE, WS-DAY-NUMBER, WS-WEEKDAY | WS-DAY-NUMBER, WS-WEEKDAY, BC-BUSINESS-DATE |

## Lineage

Where each data item's values come from, and which records go to which files.

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| BANK-CONTROL-RECORD | file CONTROL-IN | READ | 0000-MAIN | [BANKDATE.cbl:61](../../../pangaea/programs/BANKDATE.cbl#L61) |
| BC-BUSINESS-DATE | WS-DAY-NUMBER | COMPUTE | 1000-NEXT-BUSINESS-DAY | [BANKDATE.cbl:83](../../../pangaea/programs/BANKDATE.cbl#L83) |
| file CONTROL-OUT | CONTROL-OUT-RECORD | WRITE | 0000-MAIN | [BANKDATE.cbl:68](../../../pangaea/programs/BANKDATE.cbl#L68) |
| CONTROL-OUT-RECORD | BANK-CONTROL-RECORD | WRITE | 0000-MAIN | [BANKDATE.cbl:68](../../../pangaea/programs/BANKDATE.cbl#L68) |
| WS-DAY-NUMBER | BC-BUSINESS-DATE | COMPUTE | 1000-NEXT-BUSINESS-DAY | [BANKDATE.cbl:75](../../../pangaea/programs/BANKDATE.cbl#L75) |
| WS-WEEKDAY | WS-DAY-NUMBER | COMPUTE | 1000-NEXT-BUSINESS-DAY | [BANKDATE.cbl:78](../../../pangaea/programs/BANKDATE.cbl#L78) |
| WS-WEEKDAY | WS-DAY-NUMBER | COMPUTE | 1000-NEXT-BUSINESS-DAY | [BANKDATE.cbl:81](../../../pangaea/programs/BANKDATE.cbl#L81) |

## Records

### BANK-CONTROL-RECORD

File CONTROL-IN. 80 bytes, layout in copybook [BANKCTL](../copybooks/BANKCTL.md).

### CONTROL-OUT-RECORD

File CONTROL-OUT. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | CONTROL-OUT-RECORD | 1 | 80 | X(80) | | | [BANKDATE.cbl:38](../../../pangaea/programs/BANKDATE.cbl#L38) |

### WS-IN-STATUS

Working storage. 2 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-IN-STATUS | 1 | 2 | XX | | | [BANKDATE.cbl:41](../../../pangaea/programs/BANKDATE.cbl#L41) |

### WS-OUT-STATUS

Working storage. 2 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-OUT-STATUS | 1 | 2 | XX | | | [BANKDATE.cbl:42](../../../pangaea/programs/BANKDATE.cbl#L42) |

### WS-DAY-NUMBER

Working storage. 7 bytes. DAY NUMBERS COUNT FROM 1601-01-01, A MONDAY, SO A DAY NUMBER MODULO 7 IS 1 ON MONDAYS, 6 ON SATURDAYS AND 0 ON SUNDAYS.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-DAY-NUMBER | 1 | 7 | 9(7) | | | [BANKDATE.cbl:46](../../../pangaea/programs/BANKDATE.cbl#L46) |

### WS-WEEKDAY

Working storage. 1 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-WEEKDAY | 1 | 1 | 9 | | | [BANKDATE.cbl:47](../../../pangaea/programs/BANKDATE.cbl#L47) |
| 88 | WEEKEND | | | | | VALUE 6, 0 | [BANKDATE.cbl:48](../../../pangaea/programs/BANKDATE.cbl#L48) |
