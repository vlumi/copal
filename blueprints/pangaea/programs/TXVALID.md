# TXVALID

```text
TXVALID: VALIDATE THE DAY'S TRANSACTION FILE

CHECKS THE FILE AS A WHOLE: A HEADER FOR TODAY'S BUSINESS
DATE FROM A KNOWN SOURCE, THE BRANCHES OR THE CARDS, THEN
DETAILS, THEN A TRAILER WHOSE CONTROL TOTALS MATCH THE
DETAILS. A FILE THAT FAILS THESE CHECKS IS NOT PROCESSED.
CHECKS EACH DETAIL'S FIELDS AND WRITES IT EITHER TO THE VALID
FILE OR TO THE REJECTS: THE ACCOUNT NUMBER, THE CODE, WHICH
MUST BE ONE ITS SOURCE SENDS, THE AMOUNT, AND FOR A TRANSFER
THE ACCOUNT IT GOES TO. WHETHER THE ACCOUNTS EXIST IS LEFT TO
POSTING, WHICH READS THE MASTER. THE REJECTS ARE APPENDED TO,
AS THE DAILY BATCH VALIDATES THE CARD FILE AFTER THE
BRANCHES'.

RETURN CODES: 0 ALL VALID, 4 SOME DETAILS REJECTED,
              12 FILE REJECTED, 16 I/O ERROR.

MODIFICATION HISTORY
DATE       BY  CHANGE
1987-03-12 JA  ORIGINAL VERSION.
1991-06-03 PK  CR-0117 ACCOUNT NUMBER HASH IN THE TRAILER.
1996-04-15 PK  CR-0233 TRANSFERS MUST NAME ANOTHER ACCOUNT.
1999-01-11 SL  Y2K: 8-DIGIT BUSINESS DATE FROM THE BANK
               CONTROL RECORD.
2003-05-19 MV  CR-0342 THE CARD FILE: ITS SOURCE AND CODES.
               REJECTS APPENDED TO.
2004-09-20 PK  CR-0388 REJECT RECORDS OF UNKNOWN TYPE
               INSTEAD OF STOPPING THE RUN.
2009-03-02 MV  CR-0444 THE DATE OF THE REJECTION ON EVERY
               REJECT.
```

Program [TXVALID.cbl](../../../pangaea/programs/TXVALID.cbl): 5 files, 10 paragraphs. Its business rules: [8 decision tables](../rules/TXVALID.md).

- AUTHOR: J. ARKWRIGHT, BATCH SYSTEMS.
- INSTALLATION: FIRST BANK OF PANGAEA, DATA CENTER.
- DATE-WRITTEN: 1987-03-12.

## Contents

- [Files](#files)
- [Flow](#flow)
- [Paragraphs](#paragraphs)
- [Data flow](#data-flow)
- [Lineage](#lineage)
  - [Into BANK-CONTROL-RECORD](#into-bank-control-record)
  - [Into TRAN-RECORD](#into-tran-record)
  - [Into TRAN-OUT-RECORD](#into-tran-out-record)
  - [Into REJECT-OUT-RECORD](#into-reject-out-record)
  - [Into REPORT-LINE](#into-report-line)
  - [Into WS-BUSINESS-DATE](#into-ws-business-date)
  - [Into WS-SOURCE](#into-ws-source)
  - [Into REJECT-RECORD](#into-reject-record)
  - [Into WS-FILE-ERROR](#into-ws-file-error)
  - [Into WS-DETAIL-STATE](#into-ws-detail-state)
  - [Into WS-TOTALS](#into-ws-totals)
  - [Into WS-REPORT-COUNT-LINE](#into-ws-report-count-line)
  - [Into WS-REPORT-AMOUNT-LINE](#into-ws-report-amount-line)
  - [Out of the program](#out-of-the-program)
  - [Into other data items](#into-other-data-items)
- [Field origins](#field-origins)
  - [File TRAN-OUT](#file-tran-out)
  - [File REJECT-OUT](#file-reject-out)
  - [File REPORT-OUT](#file-report-out)
- [Records](#records)
  - [BANK-CONTROL-RECORD](#bank-control-record)
  - [TRAN-RECORD](#tran-record)
  - [TRAN-OUT-RECORD](#tran-out-record)
  - [REJECT-OUT-RECORD](#reject-out-record)
  - [REPORT-LINE](#report-line)
  - [WS-FILE-STATUSES](#ws-file-statuses)
  - [WS-BUSINESS-DATE](#ws-business-date)
  - [WS-SOURCE](#ws-source)
  - [REJECT-RECORD](#reject-record)
  - [WS-FILE-ERROR](#ws-file-error)
  - [WS-DETAIL-STATE](#ws-detail-state)
  - [WS-TOTALS](#ws-totals)
  - [WS-REPORT-COUNT-LINE](#ws-report-count-line)
  - [WS-REPORT-AMOUNT-LINE](#ws-report-amount-line)

## Files

| File | Kind | Assigned to | Organization | Status field | Opened | Records |
| --- | --- | --- | --- | --- | --- | --- |
| CONTROL-FILE | FD | BANKCTL | SEQUENTIAL | WS-CONTROL-STATUS | INPUT | BANK-CONTROL-RECORD |
| TRAN-IN | FD | TRANIN | SEQUENTIAL | WS-TRAN-IN-STATUS | INPUT | TRAN-RECORD |
| TRAN-OUT | FD | TRANOUT | SEQUENTIAL | WS-TRAN-OUT-STATUS | OUTPUT | TRAN-OUT-RECORD |
| REJECT-OUT | FD | REJECTS | SEQUENTIAL | WS-REJECT-STATUS | EXTEND | REJECT-OUT-RECORD |
| REPORT-OUT | FD | VALRPT | LINE SEQUENTIAL | WS-REPORT-STATUS | OUTPUT | REPORT-LINE |

## Flow

![The paragraphs of TXVALID and how control moves between them](TXVALID.svg)

## Paragraphs

| Paragraph | Description | Source | Performs | Goes to | Performed by | Reached by |
| --- | --- | --- | --- | --- | --- | --- |
| 0000-MAIN | VALIDATE THE FILE AND SET THE RETURN CODE. | [TXVALID.cbl:113](../../../pangaea/programs/TXVALID.cbl#L113) | 1000-OPEN-FILES, 2000-READ-HEADER, 3000-READ-DETAILS, 4000-CHECK-TRAILER, 8000-WRITE-REPORT, 9000-CLOSE-FILES | | | |
| 1000-OPEN-FILES | OPEN THE FILES AND READ THE BUSINESS DATE. THE REJECTS ARE EXTENDED, AND CREATED IF NEW. | [TXVALID.cbl:136](../../../pangaea/programs/TXVALID.cbl#L136) | | | 0000-MAIN | |
| 2000-READ-HEADER | THE FIRST RECORD MUST BE THE HEADER FOR THE BUSINESS DATE, FROM A KNOWN SOURCE. | [TXVALID.cbl:163](../../../pangaea/programs/TXVALID.cbl#L163) | | | 0000-MAIN | |
| 3000-READ-DETAILS | VALIDATE EACH RECORD UNTIL THE TRAILER. | [TXVALID.cbl:185](../../../pangaea/programs/TXVALID.cbl#L185) | 3100-READ-NEXT, 3200-VALIDATE-DETAIL, 3300-REJECT-RECORD | | 0000-MAIN | |
| 3100-READ-NEXT | READ ONE RECORD. THE FILE MUST NOT END BEFORE ITS TRAILER. | [TXVALID.cbl:203](../../../pangaea/programs/TXVALID.cbl#L203) | | | 3000-READ-DETAILS | |
| 3200-VALIDATE-DETAIL | ADD A DETAIL TO THE TOTALS, THEN CHECK ITS FIELDS. | [TXVALID.cbl:210](../../../pangaea/programs/TXVALID.cbl#L210) | 3300-REJECT-RECORD | | 3000-READ-DETAILS | |
| 3300-REJECT-RECORD | WRITE THE RECORD TO THE REJECTS WITH ITS REASON. | [TXVALID.cbl:245](../../../pangaea/programs/TXVALID.cbl#L245) | | | 3000-READ-DETAILS, 3200-VALIDATE-DETAIL | |
| 4000-CHECK-TRAILER | THE TRAILER MUST MATCH THE TOTALS, AND END THE FILE. | [TXVALID.cbl:252](../../../pangaea/programs/TXVALID.cbl#L252) | | | 0000-MAIN | |
| 8000-WRITE-REPORT | WRITE THE VALIDATION REPORT. | [TXVALID.cbl:273](../../../pangaea/programs/TXVALID.cbl#L273) | | | 0000-MAIN | |
| 9000-CLOSE-FILES | CLOSE ALL FILES. | [TXVALID.cbl:305](../../../pangaea/programs/TXVALID.cbl#L305) | | | 0000-MAIN | |

## Data flow

| Paragraph | Reads | Writes |
| --- | --- | --- |
| 0000-MAIN | WS-FILE-ERROR, WS-REJECT-COUNT | RETURN-CODE |
| 1000-OPEN-FILES | WS-CONTROL-STATUS, WS-TRAN-IN-STATUS, WS-TRAN-OUT-STATUS, WS-REJECT-STATUS, WS-REPORT-STATUS, WS-FILE-STATUSES, BC-BUSINESS-DATE, WS-BUSINESS-DATE | RETURN-CODE, BANK-CONTROL-RECORD, WS-BUSINESS-DATE, REJECT-RECORD, RJ-REJECT-DATE |
| 2000-READ-HEADER | WS-FILE-ERROR, TR-REC-TYPE, TH-BUSINESS-DATE, WS-BUSINESS-DATE, TH-SOURCE | TRAN-RECORD, WS-FILE-ERROR, WS-SOURCE |
| 3000-READ-DETAILS | TR-REC-TYPE, WS-FILE-ERROR | WS-FILE-ERROR, RJ-REASON |
| 3100-READ-NEXT | | TRAN-RECORD, WS-FILE-ERROR |
| 3200-VALIDATE-DETAIL | WS-DETAIL-COUNT, TD-AMOUNT, WS-AMOUNT-TOTAL, TD-ACCT-NO, WS-ACCT-HASH, WS-SOURCE, TD-TRAN-CODE, TD-COUNTER-ACCT-NO, WS-DETAIL-STATE, TRAN-RECORD, WS-VALID-COUNT | WS-DETAIL-COUNT, WS-AMOUNT-TOTAL, WS-ACCT-HASH, WS-DETAIL-STATE, RJ-REASON, TRAN-OUT-RECORD, WS-VALID-COUNT |
| 3300-REJECT-RECORD | TRAN-RECORD, REJECT-RECORD, WS-REJECT-COUNT | RJ-PROGRAM, RJ-TRAN, REJECT-OUT-RECORD, WS-REJECT-COUNT |
| 4000-CHECK-TRAILER | TT-DETAIL-COUNT, WS-DETAIL-COUNT, TT-AMOUNT-TOTAL, WS-AMOUNT-TOTAL, TT-ACCT-HASH, WS-ACCT-HASH, WS-FILE-ERROR | WS-FILE-ERROR, TRAN-RECORD |
| 8000-WRITE-REPORT | WS-BUSINESS-DATE, WS-DETAIL-COUNT, WS-REPORT-COUNT-LINE, WS-VALID-COUNT, WS-REJECT-COUNT, WS-AMOUNT-TOTAL, WS-REPORT-AMOUNT-LINE, WS-FILE-ERROR | REPORT-LINE, WS-RC-LABEL, WS-RC-VALUE, WS-RA-LABEL, WS-RA-VALUE |
| 9000-CLOSE-FILES | | |

## Lineage

Where each data item's values come from, grouped by the record they go into. A data item links to the record it belongs to.

### Into BANK-CONTROL-RECORD

Record [BANK-CONTROL-RECORD](#bank-control-record), from copybook [BANKCTL](../copybooks/BANKCTL.md).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| BANK-CONTROL-RECORD | file CONTROL-FILE | READ | 1000-OPEN-FILES | [TXVALID.cbl:149](../../../pangaea/programs/TXVALID.cbl#L149) |

### Into TRAN-RECORD

Record [TRAN-RECORD](#tran-record), from copybook [TRANREC](../copybooks/TRANREC.md).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| TRAN-RECORD | file TRAN-IN | READ | 2000-READ-HEADER | [TXVALID.cbl:164](../../../pangaea/programs/TXVALID.cbl#L164) |
| TRAN-RECORD | file TRAN-IN | READ | 3100-READ-NEXT | [TXVALID.cbl:204](../../../pangaea/programs/TXVALID.cbl#L204) |
| TRAN-RECORD | file TRAN-IN | READ | 4000-CHECK-TRAILER | [TXVALID.cbl:265](../../../pangaea/programs/TXVALID.cbl#L265) |

### Into TRAN-OUT-RECORD

Record [TRAN-OUT-RECORD](#tran-out-record).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| TRAN-OUT-RECORD | [TRAN-RECORD](#tran-record) | WRITE | 3200-VALIDATE-DETAIL | [TXVALID.cbl:238](../../../pangaea/programs/TXVALID.cbl#L238) |

### Into REJECT-OUT-RECORD

Record [REJECT-OUT-RECORD](#reject-out-record).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| REJECT-OUT-RECORD | [REJECT-RECORD](#reject-record) | WRITE | 3300-REJECT-RECORD | [TXVALID.cbl:248](../../../pangaea/programs/TXVALID.cbl#L248) |

### Into REPORT-LINE

Record [REPORT-LINE](#report-line).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| REPORT-LINE | literal SPACES | MOVE | 8000-WRITE-REPORT | [TXVALID.cbl:274](../../../pangaea/programs/TXVALID.cbl#L274) |
| REPORT-LINE | literal 'TXVALID  TRANSACTION VALIDATION  BUSINESS DATE ' | STRING | 8000-WRITE-REPORT | [TXVALID.cbl:275](../../../pangaea/programs/TXVALID.cbl#L275) |
| REPORT-LINE | [WS-BUSINESS-DATE](#ws-business-date) | STRING | 8000-WRITE-REPORT | [TXVALID.cbl:275](../../../pangaea/programs/TXVALID.cbl#L275) |
| REPORT-LINE | literal SPACES | MOVE | 8000-WRITE-REPORT | [TXVALID.cbl:279](../../../pangaea/programs/TXVALID.cbl#L279) |
| REPORT-LINE | [WS-REPORT-COUNT-LINE](#ws-report-count-line) | WRITE | 8000-WRITE-REPORT | [TXVALID.cbl:283](../../../pangaea/programs/TXVALID.cbl#L283) |
| REPORT-LINE | [WS-REPORT-COUNT-LINE](#ws-report-count-line) | WRITE | 8000-WRITE-REPORT | [TXVALID.cbl:286](../../../pangaea/programs/TXVALID.cbl#L286) |
| REPORT-LINE | [WS-REPORT-COUNT-LINE](#ws-report-count-line) | WRITE | 8000-WRITE-REPORT | [TXVALID.cbl:289](../../../pangaea/programs/TXVALID.cbl#L289) |
| REPORT-LINE | [WS-REPORT-AMOUNT-LINE](#ws-report-amount-line) | WRITE | 8000-WRITE-REPORT | [TXVALID.cbl:292](../../../pangaea/programs/TXVALID.cbl#L292) |
| REPORT-LINE | literal SPACES | MOVE | 8000-WRITE-REPORT | [TXVALID.cbl:293](../../../pangaea/programs/TXVALID.cbl#L293) |
| REPORT-LINE | literal SPACES | MOVE | 8000-WRITE-REPORT | [TXVALID.cbl:295](../../../pangaea/programs/TXVALID.cbl#L295) |
| REPORT-LINE | literal 'FILE ACCEPTED' | MOVE | 8000-WRITE-REPORT | [TXVALID.cbl:297](../../../pangaea/programs/TXVALID.cbl#L297) |
| REPORT-LINE | literal 'FILE REJECTED: ' | STRING | 8000-WRITE-REPORT | [TXVALID.cbl:299](../../../pangaea/programs/TXVALID.cbl#L299) |
| REPORT-LINE | [WS-FILE-ERROR](#ws-file-error) | STRING | 8000-WRITE-REPORT | [TXVALID.cbl:299](../../../pangaea/programs/TXVALID.cbl#L299) |

### Into WS-BUSINESS-DATE

Record [WS-BUSINESS-DATE](#ws-business-date).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-BUSINESS-DATE | [BC-BUSINESS-DATE](#bank-control-record "in BANK-CONTROL-RECORD") | MOVE | 1000-OPEN-FILES | [TXVALID.cbl:155](../../../pangaea/programs/TXVALID.cbl#L155) |

### Into WS-SOURCE

Record [WS-SOURCE](#ws-source).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-SOURCE | [TH-SOURCE](#tran-record "in TRAN-RECORD") | MOVE | 2000-READ-HEADER | [TXVALID.cbl:181](../../../pangaea/programs/TXVALID.cbl#L181) |

### Into REJECT-RECORD

Record [REJECT-RECORD](#reject-record), from copybook [REJREC](../copybooks/REJREC.md).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| REJECT-RECORD | literal SPACES | MOVE | 1000-OPEN-FILES | [TXVALID.cbl:156](../../../pangaea/programs/TXVALID.cbl#L156) |
| RJ-PROGRAM | literal 'TXVALID' | MOVE | 3300-REJECT-RECORD | [TXVALID.cbl:246](../../../pangaea/programs/TXVALID.cbl#L246) |
| RJ-REASON | literal 'TYPE' | MOVE | 3000-READ-DETAILS | [TXVALID.cbl:194](../../../pangaea/programs/TXVALID.cbl#L194) |
| RJ-REASON | literal condition RJ-BAD-ACCOUNT | SET | 3200-VALIDATE-DETAIL | [TXVALID.cbl:221](../../../pangaea/programs/TXVALID.cbl#L221) |
| RJ-REASON | literal condition RJ-BAD-CODE | SET | 3200-VALIDATE-DETAIL | [TXVALID.cbl:225](../../../pangaea/programs/TXVALID.cbl#L225) |
| RJ-REASON | literal condition RJ-BAD-AMOUNT | SET | 3200-VALIDATE-DETAIL | [TXVALID.cbl:229](../../../pangaea/programs/TXVALID.cbl#L229) |
| RJ-REASON | literal condition RJ-BAD-COUNTER | SET | 3200-VALIDATE-DETAIL | [TXVALID.cbl:234](../../../pangaea/programs/TXVALID.cbl#L234) |
| RJ-REJECT-DATE | [WS-BUSINESS-DATE](#ws-business-date) | STRING | 1000-OPEN-FILES | [TXVALID.cbl:157](../../../pangaea/programs/TXVALID.cbl#L157) |
| RJ-REJECT-DATE | literal '-' | STRING | 1000-OPEN-FILES | [TXVALID.cbl:157](../../../pangaea/programs/TXVALID.cbl#L157) |
| RJ-TRAN | [TRAN-RECORD](#tran-record) | MOVE | 3300-REJECT-RECORD | [TXVALID.cbl:247](../../../pangaea/programs/TXVALID.cbl#L247) |

### Into WS-FILE-ERROR

Record [WS-FILE-ERROR](#ws-file-error).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-FILE-ERROR | literal 'EMPTY FILE' | MOVE | 2000-READ-HEADER | [TXVALID.cbl:166](../../../pangaea/programs/TXVALID.cbl#L166) |
| WS-FILE-ERROR | literal 'FIRST RECORD IS NOT A HEADER' | MOVE | 2000-READ-HEADER | [TXVALID.cbl:172](../../../pangaea/programs/TXVALID.cbl#L172) |
| WS-FILE-ERROR | literal 'HEADER IS NOT FOR THE BUSINESS DATE' | MOVE | 2000-READ-HEADER | [TXVALID.cbl:174](../../../pangaea/programs/TXVALID.cbl#L174) |
| WS-FILE-ERROR | literal 'HEADER IS FROM AN UNKNOWN SOURCE' | MOVE | 2000-READ-HEADER | [TXVALID.cbl:177](../../../pangaea/programs/TXVALID.cbl#L177) |
| WS-FILE-ERROR | literal 'SECOND HEADER' | MOVE | 3000-READ-DETAILS | [TXVALID.cbl:192](../../../pangaea/programs/TXVALID.cbl#L192) |
| WS-FILE-ERROR | literal 'NO TRAILER' | MOVE | 3100-READ-NEXT | [TXVALID.cbl:206](../../../pangaea/programs/TXVALID.cbl#L206) |
| WS-FILE-ERROR | literal 'TRAILER DETAIL COUNT DIFFERS' | MOVE | 4000-CHECK-TRAILER | [TXVALID.cbl:255](../../../pangaea/programs/TXVALID.cbl#L255) |
| WS-FILE-ERROR | literal 'TRAILER AMOUNT TOTAL DIFFERS' | MOVE | 4000-CHECK-TRAILER | [TXVALID.cbl:258](../../../pangaea/programs/TXVALID.cbl#L258) |
| WS-FILE-ERROR | literal 'TRAILER ACCOUNT HASH DIFFERS' | MOVE | 4000-CHECK-TRAILER | [TXVALID.cbl:261](../../../pangaea/programs/TXVALID.cbl#L261) |
| WS-FILE-ERROR | literal 'RECORDS AFTER THE TRAILER' | MOVE | 4000-CHECK-TRAILER | [TXVALID.cbl:267](../../../pangaea/programs/TXVALID.cbl#L267) |

### Into WS-DETAIL-STATE

Record [WS-DETAIL-STATE](#ws-detail-state).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-DETAIL-STATE | literal condition DETAIL-IS-VALID | SET | 3200-VALIDATE-DETAIL | [TXVALID.cbl:218](../../../pangaea/programs/TXVALID.cbl#L218) |
| WS-DETAIL-STATE | literal condition DETAIL-IS-REJECTED | SET | 3200-VALIDATE-DETAIL | [TXVALID.cbl:222](../../../pangaea/programs/TXVALID.cbl#L222) |
| WS-DETAIL-STATE | literal condition DETAIL-IS-REJECTED | SET | 3200-VALIDATE-DETAIL | [TXVALID.cbl:226](../../../pangaea/programs/TXVALID.cbl#L226) |
| WS-DETAIL-STATE | literal condition DETAIL-IS-REJECTED | SET | 3200-VALIDATE-DETAIL | [TXVALID.cbl:230](../../../pangaea/programs/TXVALID.cbl#L230) |
| WS-DETAIL-STATE | literal condition DETAIL-IS-REJECTED | SET | 3200-VALIDATE-DETAIL | [TXVALID.cbl:235](../../../pangaea/programs/TXVALID.cbl#L235) |

### Into WS-TOTALS

Record [WS-TOTALS](#ws-totals).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-ACCT-HASH | [TD-ACCT-NO](#tran-record "in TRAN-RECORD") | ADD | 3200-VALIDATE-DETAIL | [TXVALID.cbl:216](../../../pangaea/programs/TXVALID.cbl#L216) |
| WS-AMOUNT-TOTAL | [TD-AMOUNT](#tran-record "in TRAN-RECORD") | ADD | 3200-VALIDATE-DETAIL | [TXVALID.cbl:213](../../../pangaea/programs/TXVALID.cbl#L213) |
| WS-DETAIL-COUNT | literal 1 | ADD | 3200-VALIDATE-DETAIL | [TXVALID.cbl:211](../../../pangaea/programs/TXVALID.cbl#L211) |
| WS-REJECT-COUNT | literal 1 | ADD | 3300-REJECT-RECORD | [TXVALID.cbl:249](../../../pangaea/programs/TXVALID.cbl#L249) |
| WS-VALID-COUNT | literal 1 | ADD | 3200-VALIDATE-DETAIL | [TXVALID.cbl:239](../../../pangaea/programs/TXVALID.cbl#L239) |

### Into WS-REPORT-COUNT-LINE

Record [WS-REPORT-COUNT-LINE](#ws-report-count-line).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-RC-LABEL | literal 'DETAILS READ' | MOVE | 8000-WRITE-REPORT | [TXVALID.cbl:281](../../../pangaea/programs/TXVALID.cbl#L281) |
| WS-RC-LABEL | literal 'DETAILS VALID' | MOVE | 8000-WRITE-REPORT | [TXVALID.cbl:284](../../../pangaea/programs/TXVALID.cbl#L284) |
| WS-RC-LABEL | literal 'RECORDS REJECTED' | MOVE | 8000-WRITE-REPORT | [TXVALID.cbl:287](../../../pangaea/programs/TXVALID.cbl#L287) |
| WS-RC-VALUE | [WS-DETAIL-COUNT](#ws-totals "in WS-TOTALS") | MOVE | 8000-WRITE-REPORT | [TXVALID.cbl:282](../../../pangaea/programs/TXVALID.cbl#L282) |
| WS-RC-VALUE | [WS-VALID-COUNT](#ws-totals "in WS-TOTALS") | MOVE | 8000-WRITE-REPORT | [TXVALID.cbl:285](../../../pangaea/programs/TXVALID.cbl#L285) |
| WS-RC-VALUE | [WS-REJECT-COUNT](#ws-totals "in WS-TOTALS") | MOVE | 8000-WRITE-REPORT | [TXVALID.cbl:288](../../../pangaea/programs/TXVALID.cbl#L288) |

### Into WS-REPORT-AMOUNT-LINE

Record [WS-REPORT-AMOUNT-LINE](#ws-report-amount-line).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-RA-LABEL | literal 'AMOUNT TOTAL' | MOVE | 8000-WRITE-REPORT | [TXVALID.cbl:290](../../../pangaea/programs/TXVALID.cbl#L290) |
| WS-RA-VALUE | [WS-AMOUNT-TOTAL](#ws-totals "in WS-TOTALS") | MOVE | 8000-WRITE-REPORT | [TXVALID.cbl:291](../../../pangaea/programs/TXVALID.cbl#L291) |

### Out of the program

Records written to files, items passed to the programs called, and items displayed to the job's output.

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| file REJECT-OUT | [REJECT-OUT-RECORD](#reject-out-record) | WRITE | 3300-REJECT-RECORD | [TXVALID.cbl:248](../../../pangaea/programs/TXVALID.cbl#L248) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [TXVALID.cbl:278](../../../pangaea/programs/TXVALID.cbl#L278) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [TXVALID.cbl:280](../../../pangaea/programs/TXVALID.cbl#L280) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [TXVALID.cbl:283](../../../pangaea/programs/TXVALID.cbl#L283) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [TXVALID.cbl:286](../../../pangaea/programs/TXVALID.cbl#L286) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [TXVALID.cbl:289](../../../pangaea/programs/TXVALID.cbl#L289) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [TXVALID.cbl:292](../../../pangaea/programs/TXVALID.cbl#L292) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [TXVALID.cbl:294](../../../pangaea/programs/TXVALID.cbl#L294) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [TXVALID.cbl:302](../../../pangaea/programs/TXVALID.cbl#L302) |
| SYSOUT | [WS-FILE-STATUSES](#ws-file-statuses) | DISPLAY | 1000-OPEN-FILES | [TXVALID.cbl:145](../../../pangaea/programs/TXVALID.cbl#L145) |
| file TRAN-OUT | [TRAN-OUT-RECORD](#tran-out-record) | WRITE | 3200-VALIDATE-DETAIL | [TXVALID.cbl:238](../../../pangaea/programs/TXVALID.cbl#L238) |

### Into other data items

Data items the program does not declare, such as special registers and index names.

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| RETURN-CODE | literal 12 | MOVE | 0000-MAIN | [TXVALID.cbl:126](../../../pangaea/programs/TXVALID.cbl#L126) |
| RETURN-CODE | literal 4 | MOVE | 0000-MAIN | [TXVALID.cbl:128](../../../pangaea/programs/TXVALID.cbl#L128) |
| RETURN-CODE | literal 0 | MOVE | 0000-MAIN | [TXVALID.cbl:130](../../../pangaea/programs/TXVALID.cbl#L130) |
| RETURN-CODE | literal 16 | MOVE | 1000-OPEN-FILES | [TXVALID.cbl:146](../../../pangaea/programs/TXVALID.cbl#L146) |
| RETURN-CODE | literal 16 | MOVE | 1000-OPEN-FILES | [TXVALID.cbl:152](../../../pangaea/programs/TXVALID.cbl#L152) |

## Field origins

Where the values of each field the program writes to a file come from: the fields of the files it reads, literals, and initial values, followed through every statement that moves a value, and through groups and redefinitions by the bytes they share.

### File TRAN-OUT

Written from [TRAN-RECORD](#tran-record).

| Field | Position | Origins |
| --- | --: | --- |
| TR-REC-TYPE | 1 | file TRAN-IN TR-REC-TYPE |
| TR-BODY | 2-80 | file TRAN-IN TR-BODY |
| TH-BUSINESS-DATE | 2-9 | file TRAN-IN TH-BUSINESS-DATE |
| TH-SOURCE | 10-17 | file TRAN-IN TH-SOURCE |
| TH-FILE-SEQ | 18-23 | file TRAN-IN TH-FILE-SEQ |
| TD-TRAN-ID | 2-11 | file TRAN-IN TD-TRAN-ID |
| TD-ACCT-NO | 12-21 | file TRAN-IN TD-ACCT-NO |
| TD-TRAN-CODE | 22-23 | file TRAN-IN TD-TRAN-CODE |
| TD-AMOUNT | 24-34 | file TRAN-IN TD-AMOUNT |
| TD-DESCRIPTION | 35-54 | file TRAN-IN TD-DESCRIPTION |
| TD-COUNTER-ACCT-NO | 55-64 | file TRAN-IN TD-COUNTER-ACCT-NO |
| TT-DETAIL-COUNT | 2-9 | file TRAN-IN TH-BUSINESS-DATE |
| TT-AMOUNT-TOTAL | 10-24 | file TRAN-IN TT-AMOUNT-TOTAL |
| TT-ACCT-HASH | 25-39 | file TRAN-IN TT-ACCT-HASH |

### File REJECT-OUT

Written from [REJECT-RECORD](#reject-record).

| Field | Position | Origins |
| --- | --: | --- |
| RJ-PROGRAM | 1-8 | literal 'TXVALID' |
| RJ-REASON | 9-12 | literal 'TYPE'; literal condition RJ-BAD-ACCOUNT; literal condition RJ-BAD-AMOUNT; literal condition RJ-BAD-CODE; literal condition RJ-BAD-COUNTER |
| RJ-TRAN | 13-92 | file TRAN-IN TRAN-RECORD |
| RJ-REJECT-DATE | 101-110 | file CONTROL-FILE BC-BUSINESS-DATE; literal '-' |

### File REPORT-OUT

Written from [REPORT-LINE](#report-line), [WS-REPORT-COUNT-LINE](#ws-report-count-line) and [WS-REPORT-AMOUNT-LINE](#ws-report-amount-line).

#### Record REPORT-LINE

| Field | Position | Origins |
| --- | --: | --- |
| REPORT-LINE | 1-80 | file CONTROL-FILE BC-BUSINESS-DATE; initial value of WS-FILE-ERROR; literal 'EMPTY FILE'; literal 'FILE ACCEPTED'; literal 'FILE REJECTED: '; literal 'FIRST RECORD IS NOT A HEADER'; literal 'HEADER IS FROM AN UNKNOWN SOURCE'; literal 'HEADER IS NOT FOR THE BUSINESS DATE'; literal 'NO TRAILER'; literal 'RECORDS AFTER THE TRAILER'; literal 'SECOND HEADER'; literal 'TRAILER ACCOUNT HASH DIFFERS'; literal 'TRAILER AMOUNT TOTAL DIFFERS'; literal 'TRAILER DETAIL COUNT DIFFERS'; literal 'TXVALID  TRANSACTION VALIDATION  BUSINESS DATE '; literal SPACES |

#### Record WS-REPORT-COUNT-LINE

| Field | Position | Origins |
| --- | --: | --- |
| WS-RC-LABEL | 1-30 | literal 'DETAILS READ'; literal 'DETAILS VALID'; literal 'RECORDS REJECTED' |
| WS-RC-VALUE | 31-40 | initial value of WS-DETAIL-COUNT; initial value of WS-REJECT-COUNT; initial value of WS-VALID-COUNT; literal 1 |

#### Record WS-REPORT-AMOUNT-LINE

| Field | Position | Origins |
| --- | --: | --- |
| WS-RA-LABEL | 1-30 | literal 'AMOUNT TOTAL' |
| WS-RA-VALUE | 31-50 | file TRAN-IN TD-AMOUNT; initial value of WS-AMOUNT-TOTAL |

## Records

### BANK-CONTROL-RECORD

File CONTROL-FILE. 80 bytes, layout in copybook [BANKCTL](../copybooks/BANKCTL.md).

### TRAN-RECORD

File TRAN-IN. 80 bytes, layout in copybook [TRANREC](../copybooks/TRANREC.md).

### TRAN-OUT-RECORD

File TRAN-OUT. 80 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | TRAN-OUT-RECORD | text | 1 | 80 | X(80) | | | [TXVALID.cbl:65](../../../pangaea/programs/TXVALID.cbl#L65) |

### REJECT-OUT-RECORD

File REJECT-OUT. 110 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | REJECT-OUT-RECORD | text | 1 | 110 | X(110) | | | [TXVALID.cbl:67](../../../pangaea/programs/TXVALID.cbl#L67) |

### REPORT-LINE

File REPORT-OUT. 80 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | REPORT-LINE | text | 1 | 80 | X(80) | | | [TXVALID.cbl:69](../../../pangaea/programs/TXVALID.cbl#L69) |

### WS-FILE-STATUSES

Working storage. 10 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-FILE-STATUSES | group | 1 | 10 | | | | [TXVALID.cbl:72](../../../pangaea/programs/TXVALID.cbl#L72) |
| 05 | ├ WS-CONTROL-STATUS | text | 1 | 2 | XX | | | [TXVALID.cbl:73](../../../pangaea/programs/TXVALID.cbl#L73) |
| 05 | ├ WS-TRAN-IN-STATUS | code | 3 | 2 | XX | | | [TXVALID.cbl:74](../../../pangaea/programs/TXVALID.cbl#L74) |
| 88 | │ └ TRAN-IN-EOF | | | | | | VALUE '10' | [TXVALID.cbl:75](../../../pangaea/programs/TXVALID.cbl#L75) |
| 05 | ├ WS-TRAN-OUT-STATUS | text | 5 | 2 | XX | | | [TXVALID.cbl:76](../../../pangaea/programs/TXVALID.cbl#L76) |
| 05 | ├ WS-REJECT-STATUS | code | 7 | 2 | XX | | | [TXVALID.cbl:77](../../../pangaea/programs/TXVALID.cbl#L77) |
| 88 | │ └ REJECTS-OPENED | | | | | | VALUE '00', '05' | [TXVALID.cbl:78](../../../pangaea/programs/TXVALID.cbl#L78) |
| 05 | └ WS-REPORT-STATUS | text | 9 | 2 | XX | | | [TXVALID.cbl:79](../../../pangaea/programs/TXVALID.cbl#L79) |

### WS-BUSINESS-DATE

Working storage. 8 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-BUSINESS-DATE | date | 1 | 8 | 9(8) | | | [TXVALID.cbl:81](../../../pangaea/programs/TXVALID.cbl#L81) |

### WS-SOURCE

Working storage. 8 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-SOURCE | code | 1 | 8 | X(8) | | | [TXVALID.cbl:82](../../../pangaea/programs/TXVALID.cbl#L82) |
| 88 | ├ FROM-BRANCHES | | | | | | VALUE 'BRANCHES' | [TXVALID.cbl:83](../../../pangaea/programs/TXVALID.cbl#L83) |
| 88 | └ FROM-CARDS | | | | | | VALUE 'CARDS' | [TXVALID.cbl:84](../../../pangaea/programs/TXVALID.cbl#L84) |

### REJECT-RECORD

Working storage. 110 bytes, layout in copybook [REJREC](../copybooks/REJREC.md).

### WS-FILE-ERROR

Working storage. 40 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-FILE-ERROR | code | 1 | 40 | X(40) | | VALUE SPACES | [TXVALID.cbl:88](../../../pangaea/programs/TXVALID.cbl#L88) |
| 88 | └ FILE-IS-GOOD | | | | | | VALUE SPACES | [TXVALID.cbl:89](../../../pangaea/programs/TXVALID.cbl#L89) |

### WS-DETAIL-STATE

Working storage. 1 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-DETAIL-STATE | code | 1 | 1 | X | | | [TXVALID.cbl:91](../../../pangaea/programs/TXVALID.cbl#L91) |
| 88 | ├ DETAIL-IS-VALID | | | | | | VALUE 'Y' | [TXVALID.cbl:92](../../../pangaea/programs/TXVALID.cbl#L92) |
| 88 | └ DETAIL-IS-REJECTED | | | | | | VALUE 'N' | [TXVALID.cbl:93](../../../pangaea/programs/TXVALID.cbl#L93) |

### WS-TOTALS

Working storage. 54 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TOTALS | group | 1 | 54 | | | | [TXVALID.cbl:95](../../../pangaea/programs/TXVALID.cbl#L95) |
| 05 | ├ WS-DETAIL-COUNT | number | 1 | 8 | 9(8) | | VALUE 0 | [TXVALID.cbl:96](../../../pangaea/programs/TXVALID.cbl#L96) |
| 05 | ├ WS-VALID-COUNT | number | 9 | 8 | 9(8) | | VALUE 0 | [TXVALID.cbl:97](../../../pangaea/programs/TXVALID.cbl#L97) |
| 05 | ├ WS-REJECT-COUNT | number | 17 | 8 | 9(8) | | VALUE 0 | [TXVALID.cbl:98](../../../pangaea/programs/TXVALID.cbl#L98) |
| 05 | ├ WS-AMOUNT-TOTAL | amount | 25 | 15 | 9(13)V99 | | VALUE 0 | [TXVALID.cbl:99](../../../pangaea/programs/TXVALID.cbl#L99) |
| 05 | └ WS-ACCT-HASH | number | 40 | 15 | 9(15) | | VALUE 0 | [TXVALID.cbl:100](../../../pangaea/programs/TXVALID.cbl#L100) |

### WS-REPORT-COUNT-LINE

Working storage. 80 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-COUNT-LINE | group | 1 | 80 | | | | [TXVALID.cbl:102](../../../pangaea/programs/TXVALID.cbl#L102) |
| 05 | ├ WS-RC-LABEL | text | 1 | 30 | X(30) | | | [TXVALID.cbl:103](../../../pangaea/programs/TXVALID.cbl#L103) |
| 05 | ├ WS-RC-VALUE | edited | 31 | 10 | ZZ,ZZZ,ZZ9 | | | [TXVALID.cbl:104](../../../pangaea/programs/TXVALID.cbl#L104) |
| 05 | └ FILLER | filler | 41 | 40 | X(40) | | VALUE SPACES | [TXVALID.cbl:105](../../../pangaea/programs/TXVALID.cbl#L105) |

### WS-REPORT-AMOUNT-LINE

Working storage. 80 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-AMOUNT-LINE | group | 1 | 80 | | | | [TXVALID.cbl:106](../../../pangaea/programs/TXVALID.cbl#L106) |
| 05 | ├ WS-RA-LABEL | text | 1 | 30 | X(30) | | | [TXVALID.cbl:107](../../../pangaea/programs/TXVALID.cbl#L107) |
| 05 | ├ WS-RA-VALUE | edited | 31 | 20 | Z,ZZZ,ZZZ,ZZZ,ZZ9.99 | | | [TXVALID.cbl:108](../../../pangaea/programs/TXVALID.cbl#L108) |
| 05 | └ FILLER | filler | 51 | 30 | X(30) | | VALUE SPACES | [TXVALID.cbl:109](../../../pangaea/programs/TXVALID.cbl#L109) |
