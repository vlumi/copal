# ACCTLOAD

```text
ACCTLOAD: LOAD THE ONLINE ACCOUNT FILE

WRITES THE NEW ACCOUNT MASTER'S ACCOUNTS TO A NEW INDEXED
ONLINE ACCOUNT FILE, FOR THE CARD AND ATM AUTHORIZATIONS OF
THE NEXT BUSINESS DAY, WITH NO HOLDS. CLOSED ACCOUNTS ARE LEFT
OUT, SO THAT AUTHORIZATIONS FIND NO SUCH ACCOUNT. THE MASTER
IS IN ACCOUNT NUMBER ORDER, SO THE FILE IS WRITTEN IN KEY
ORDER, SEQUENTIALLY.

RETURN CODES: 0 LOADED, 16 MASTER OUT OF ORDER OR I/O.

MODIFICATION HISTORY
DATE       BY  CHANGE
2003-05-19 MV  CR-0342 ORIGINAL VERSION.
```

Program [ACCTLOAD.cbl](../../../pangaea/programs/ACCTLOAD.cbl): 4 files, 8 paragraphs. Its business rules: [4 decision tables](../rules/ACCTLOAD.md).

- AUTHOR: M. VARGA, ONLINE SYSTEMS.
- INSTALLATION: FIRST BANK OF PANGAEA, DATA CENTER.
- DATE-WRITTEN: 2003-05-19.

## Contents

- [Files](#files)
- [Flow](#flow)
- [Paragraphs](#paragraphs)
- [Data flow](#data-flow)
- [Lineage](#lineage)
  - [Into BANK-CONTROL-RECORD](#into-bank-control-record)
  - [Into ACCT-MASTER-RECORD](#into-acct-master-record)
  - [Into ACCT-ONLINE-RECORD](#into-acct-online-record)
  - [Into REPORT-LINE](#into-report-line)
  - [Into WS-BUSINESS-DATE](#into-ws-business-date)
  - [Into WS-ABORT-MESSAGE](#into-ws-abort-message)
  - [Into WS-MASTER-END-SW](#into-ws-master-end-sw)
  - [Into WS-TOTALS](#into-ws-totals)
  - [Into WS-REPORT-COUNT-LINE](#into-ws-report-count-line)
  - [Into WS-REPORT-AMOUNT-LINE](#into-ws-report-amount-line)
  - [Into files](#into-files)
  - [Into other data items](#into-other-data-items)
- [Field origins](#field-origins)
  - [File ONLINE-OUT](#file-online-out)
  - [File REPORT-OUT](#file-report-out)
- [Records](#records)
  - [BANK-CONTROL-RECORD](#bank-control-record)
  - [ACCT-MASTER-RECORD](#acct-master-record)
  - [ACCT-ONLINE-RECORD](#acct-online-record)
  - [REPORT-LINE](#report-line)
  - [WS-FILE-STATUSES](#ws-file-statuses)
  - [WS-BUSINESS-DATE](#ws-business-date)
  - [WS-ABORT-MESSAGE](#ws-abort-message)
  - [WS-MASTER-END-SW](#ws-master-end-sw)
  - [WS-TOTALS](#ws-totals)
  - [WS-REPORT-COUNT-LINE](#ws-report-count-line)
  - [WS-REPORT-AMOUNT-LINE](#ws-report-amount-line)

## Files

| File | Kind | Assigned to | Organization | Status field | Opened | Records |
| --- | --- | --- | --- | --- | --- | --- |
| CONTROL-FILE | FD | BANKCTL | SEQUENTIAL | WS-CONTROL-STATUS | INPUT | BANK-CONTROL-RECORD |
| ACCT-IN | FD | ACCTIN | SEQUENTIAL | WS-ACCT-IN-STATUS | INPUT | ACCT-MASTER-RECORD |
| ONLINE-OUT | FD | ACCTOUT | INDEXED, key AO-ACCT-NO | WS-ONLINE-STATUS | OUTPUT | ACCT-ONLINE-RECORD |
| REPORT-OUT | FD | LOADRPT | LINE SEQUENTIAL | WS-REPORT-STATUS | OUTPUT | REPORT-LINE |

## Flow

![The paragraphs of ACCTLOAD and how control moves between them](ACCTLOAD.svg)

## Paragraphs

| Paragraph | Description | Source | Performs | Goes to | Performed by | Reached by |
| --- | --- | --- | --- | --- | --- | --- |
| 0000-MAIN | LOAD EVERY ACCOUNT, THEN REPORT. | [ACCTLOAD.cbl:81](../../../pangaea/programs/ACCTLOAD.cbl#L81) | 1000-INITIALIZE, 2000-LOAD-ACCOUNT, 8000-WRITE-REPORT, 9000-CLOSE-FILES | | | |
| 1000-INITIALIZE | OPEN THE FILES AND READ THE FIRST ACCOUNT. | [ACCTLOAD.cbl:90](../../../pangaea/programs/ACCTLOAD.cbl#L90) | 9900-ABORT, 1200-READ-ACCOUNT | | 0000-MAIN | |
| 1200-READ-ACCOUNT | | [ACCTLOAD.cbl:106](../../../pangaea/programs/ACCTLOAD.cbl#L106) | | | 1000-INITIALIZE, 2000-LOAD-ACCOUNT | |
| 2000-LOAD-ACCOUNT | WRITE AN ACCOUNT THAT IS NOT CLOSED TO THE ONLINE FILE. | [ACCTLOAD.cbl:115](../../../pangaea/programs/ACCTLOAD.cbl#L115) | 2100-WRITE-ONLINE, 1200-READ-ACCOUNT | | 0000-MAIN | |
| 2100-WRITE-ONLINE | | [ACCTLOAD.cbl:123](../../../pangaea/programs/ACCTLOAD.cbl#L123) | 9900-ABORT | | 2000-LOAD-ACCOUNT | |
| 8000-WRITE-REPORT | WRITE THE LOAD REPORT. | [ACCTLOAD.cbl:144](../../../pangaea/programs/ACCTLOAD.cbl#L144) | | | 0000-MAIN | |
| 9000-CLOSE-FILES | CLOSE ALL FILES. | [ACCTLOAD.cbl:168](../../../pangaea/programs/ACCTLOAD.cbl#L168) | | | 0000-MAIN, 9900-ABORT | |
| 9900-ABORT | STOP THE RUN WITH RETURN CODE 16, SAYING WHY. | [ACCTLOAD.cbl:172](../../../pangaea/programs/ACCTLOAD.cbl#L172) | 9000-CLOSE-FILES | | 1000-INITIALIZE, 2100-WRITE-ONLINE | |

## Data flow

| Paragraph | Reads | Writes |
| --- | --- | --- |
| 0000-MAIN | WS-MASTER-END-SW | RETURN-CODE |
| 1000-INITIALIZE | WS-FILE-STATUSES, BC-BUSINESS-DATE | RETURN-CODE, BANK-CONTROL-RECORD, WS-ABORT-MESSAGE, WS-BUSINESS-DATE |
| 1200-READ-ACCOUNT | WS-READ-COUNT | ACCT-MASTER-RECORD, WS-MASTER-END-SW, WS-READ-COUNT |
| 2000-LOAD-ACCOUNT | AM-STATUS, WS-CLOSED-COUNT | WS-CLOSED-COUNT |
| 2100-WRITE-ONLINE | AM-ACCT-NO, AM-CUST-NO, AM-SHORT-NAME, AM-ACCT-TYPE, AM-STATUS, AM-BALANCE, AM-OVERDRAFT-LIMIT, WS-BUSINESS-DATE, WS-LOADED-COUNT, WS-LOADED-BALANCE | ACCT-ONLINE-RECORD, AO-ACCT-NO, AO-CUST-NO, AO-SHORT-NAME, AO-ACCT-TYPE, AO-STATUS, AO-BALANCE, AO-OVERDRAFT-LIMIT, AO-HOLDS, AO-AS-OF-DATE, WS-ABORT-MESSAGE, WS-LOADED-COUNT, WS-LOADED-BALANCE |
| 8000-WRITE-REPORT | WS-BUSINESS-DATE, WS-READ-COUNT, WS-REPORT-COUNT-LINE, WS-LOADED-COUNT, WS-CLOSED-COUNT, WS-LOADED-BALANCE, WS-REPORT-AMOUNT-LINE | REPORT-LINE, WS-RC-LABEL, WS-RC-VALUE, WS-RA-LABEL, WS-RA-VALUE |
| 9000-CLOSE-FILES | | |
| 9900-ABORT | WS-ABORT-MESSAGE | RETURN-CODE |

## Lineage

Where each data item's values come from, grouped by the record they go into. A data item links to the record it belongs to.

### Into BANK-CONTROL-RECORD

Record [BANK-CONTROL-RECORD](#bank-control-record), from copybook [BANKCTL](../copybooks/BANKCTL.md).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| BANK-CONTROL-RECORD | file CONTROL-FILE | READ | 1000-INITIALIZE | [ACCTLOAD.cbl:98](../../../pangaea/programs/ACCTLOAD.cbl#L98) |

### Into ACCT-MASTER-RECORD

Record [ACCT-MASTER-RECORD](#acct-master-record), from copybook [ACCTMAST](../copybooks/ACCTMAST.md).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| ACCT-MASTER-RECORD | file ACCT-IN | READ | 1200-READ-ACCOUNT | [ACCTLOAD.cbl:107](../../../pangaea/programs/ACCTLOAD.cbl#L107) |

### Into ACCT-ONLINE-RECORD

Record [ACCT-ONLINE-RECORD](#acct-online-record), from copybook [ACCTONL](../copybooks/ACCTONL.md).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| ACCT-ONLINE-RECORD | literal SPACES | MOVE | 2100-WRITE-ONLINE | [ACCTLOAD.cbl:124](../../../pangaea/programs/ACCTLOAD.cbl#L124) |
| AO-ACCT-NO | [AM-ACCT-NO](#acct-master-record "in ACCT-MASTER-RECORD") | MOVE | 2100-WRITE-ONLINE | [ACCTLOAD.cbl:125](../../../pangaea/programs/ACCTLOAD.cbl#L125) |
| AO-ACCT-TYPE | [AM-ACCT-TYPE](#acct-master-record "in ACCT-MASTER-RECORD") | MOVE | 2100-WRITE-ONLINE | [ACCTLOAD.cbl:128](../../../pangaea/programs/ACCTLOAD.cbl#L128) |
| AO-AS-OF-DATE | [WS-BUSINESS-DATE](#ws-business-date) | MOVE | 2100-WRITE-ONLINE | [ACCTLOAD.cbl:133](../../../pangaea/programs/ACCTLOAD.cbl#L133) |
| AO-BALANCE | [AM-BALANCE](#acct-master-record "in ACCT-MASTER-RECORD") | MOVE | 2100-WRITE-ONLINE | [ACCTLOAD.cbl:130](../../../pangaea/programs/ACCTLOAD.cbl#L130) |
| AO-CUST-NO | [AM-CUST-NO](#acct-master-record "in ACCT-MASTER-RECORD") | MOVE | 2100-WRITE-ONLINE | [ACCTLOAD.cbl:126](../../../pangaea/programs/ACCTLOAD.cbl#L126) |
| AO-HOLDS | literal 0 | MOVE | 2100-WRITE-ONLINE | [ACCTLOAD.cbl:132](../../../pangaea/programs/ACCTLOAD.cbl#L132) |
| AO-OVERDRAFT-LIMIT | [AM-OVERDRAFT-LIMIT](#acct-master-record "in ACCT-MASTER-RECORD") | MOVE | 2100-WRITE-ONLINE | [ACCTLOAD.cbl:131](../../../pangaea/programs/ACCTLOAD.cbl#L131) |
| AO-SHORT-NAME | [AM-SHORT-NAME](#acct-master-record "in ACCT-MASTER-RECORD") | MOVE | 2100-WRITE-ONLINE | [ACCTLOAD.cbl:127](../../../pangaea/programs/ACCTLOAD.cbl#L127) |
| AO-STATUS | [AM-STATUS](#acct-master-record "in ACCT-MASTER-RECORD") | MOVE | 2100-WRITE-ONLINE | [ACCTLOAD.cbl:129](../../../pangaea/programs/ACCTLOAD.cbl#L129) |

### Into REPORT-LINE

Record [REPORT-LINE](#report-line).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| REPORT-LINE | literal SPACES | MOVE | 8000-WRITE-REPORT | [ACCTLOAD.cbl:145](../../../pangaea/programs/ACCTLOAD.cbl#L145) |
| REPORT-LINE | literal 'ACCTLOAD ONLINE ACCOUNT LOAD     BUSINESS DATE ' | STRING | 8000-WRITE-REPORT | [ACCTLOAD.cbl:146](../../../pangaea/programs/ACCTLOAD.cbl#L146) |
| REPORT-LINE | [WS-BUSINESS-DATE](#ws-business-date) | STRING | 8000-WRITE-REPORT | [ACCTLOAD.cbl:146](../../../pangaea/programs/ACCTLOAD.cbl#L146) |
| REPORT-LINE | literal SPACES | MOVE | 8000-WRITE-REPORT | [ACCTLOAD.cbl:150](../../../pangaea/programs/ACCTLOAD.cbl#L150) |
| REPORT-LINE | [WS-REPORT-COUNT-LINE](#ws-report-count-line) | WRITE | 8000-WRITE-REPORT | [ACCTLOAD.cbl:154](../../../pangaea/programs/ACCTLOAD.cbl#L154) |
| REPORT-LINE | [WS-REPORT-COUNT-LINE](#ws-report-count-line) | WRITE | 8000-WRITE-REPORT | [ACCTLOAD.cbl:157](../../../pangaea/programs/ACCTLOAD.cbl#L157) |
| REPORT-LINE | [WS-REPORT-COUNT-LINE](#ws-report-count-line) | WRITE | 8000-WRITE-REPORT | [ACCTLOAD.cbl:160](../../../pangaea/programs/ACCTLOAD.cbl#L160) |
| REPORT-LINE | literal SPACES | MOVE | 8000-WRITE-REPORT | [ACCTLOAD.cbl:161](../../../pangaea/programs/ACCTLOAD.cbl#L161) |
| REPORT-LINE | [WS-REPORT-AMOUNT-LINE](#ws-report-amount-line) | WRITE | 8000-WRITE-REPORT | [ACCTLOAD.cbl:165](../../../pangaea/programs/ACCTLOAD.cbl#L165) |

### Into WS-BUSINESS-DATE

Record [WS-BUSINESS-DATE](#ws-business-date).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-BUSINESS-DATE | [BC-BUSINESS-DATE](#bank-control-record "in BANK-CONTROL-RECORD") | MOVE | 1000-INITIALIZE | [ACCTLOAD.cbl:103](../../../pangaea/programs/ACCTLOAD.cbl#L103) |

### Into WS-ABORT-MESSAGE

Record [WS-ABORT-MESSAGE](#ws-abort-message).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-ABORT-MESSAGE | literal 'NO BANK CONTROL RECORD' | MOVE | 1000-INITIALIZE | [ACCTLOAD.cbl:100](../../../pangaea/programs/ACCTLOAD.cbl#L100) |
| WS-ABORT-MESSAGE | literal 'MASTER NOT IN ACCOUNT NUMBER ORDER' | MOVE | 2100-WRITE-ONLINE | [ACCTLOAD.cbl:136](../../../pangaea/programs/ACCTLOAD.cbl#L136) |

### Into WS-MASTER-END-SW

Record [WS-MASTER-END-SW](#ws-master-end-sw).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-MASTER-END-SW | literal condition MASTER-AT-END | SET | 1200-READ-ACCOUNT | [ACCTLOAD.cbl:109](../../../pangaea/programs/ACCTLOAD.cbl#L109) |

### Into WS-TOTALS

Record [WS-TOTALS](#ws-totals).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-CLOSED-COUNT | literal 1 | ADD | 2000-LOAD-ACCOUNT | [ACCTLOAD.cbl:117](../../../pangaea/programs/ACCTLOAD.cbl#L117) |
| WS-LOADED-BALANCE | [AM-BALANCE](#acct-master-record "in ACCT-MASTER-RECORD") | ADD | 2100-WRITE-ONLINE | [ACCTLOAD.cbl:141](../../../pangaea/programs/ACCTLOAD.cbl#L141) |
| WS-LOADED-COUNT | literal 1 | ADD | 2100-WRITE-ONLINE | [ACCTLOAD.cbl:140](../../../pangaea/programs/ACCTLOAD.cbl#L140) |
| WS-READ-COUNT | literal 1 | ADD | 1200-READ-ACCOUNT | [ACCTLOAD.cbl:111](../../../pangaea/programs/ACCTLOAD.cbl#L111) |

### Into WS-REPORT-COUNT-LINE

Record [WS-REPORT-COUNT-LINE](#ws-report-count-line).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-RC-LABEL | literal 'ACCOUNTS READ' | MOVE | 8000-WRITE-REPORT | [ACCTLOAD.cbl:152](../../../pangaea/programs/ACCTLOAD.cbl#L152) |
| WS-RC-LABEL | literal 'ACCOUNTS LOADED' | MOVE | 8000-WRITE-REPORT | [ACCTLOAD.cbl:155](../../../pangaea/programs/ACCTLOAD.cbl#L155) |
| WS-RC-LABEL | literal 'CLOSED ACCOUNTS LEFT OUT' | MOVE | 8000-WRITE-REPORT | [ACCTLOAD.cbl:158](../../../pangaea/programs/ACCTLOAD.cbl#L158) |
| WS-RC-VALUE | [WS-READ-COUNT](#ws-totals "in WS-TOTALS") | MOVE | 8000-WRITE-REPORT | [ACCTLOAD.cbl:153](../../../pangaea/programs/ACCTLOAD.cbl#L153) |
| WS-RC-VALUE | [WS-LOADED-COUNT](#ws-totals "in WS-TOTALS") | MOVE | 8000-WRITE-REPORT | [ACCTLOAD.cbl:156](../../../pangaea/programs/ACCTLOAD.cbl#L156) |
| WS-RC-VALUE | [WS-CLOSED-COUNT](#ws-totals "in WS-TOTALS") | MOVE | 8000-WRITE-REPORT | [ACCTLOAD.cbl:159](../../../pangaea/programs/ACCTLOAD.cbl#L159) |

### Into WS-REPORT-AMOUNT-LINE

Record [WS-REPORT-AMOUNT-LINE](#ws-report-amount-line).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-RA-LABEL | literal 'BALANCES LOADED' | MOVE | 8000-WRITE-REPORT | [ACCTLOAD.cbl:163](../../../pangaea/programs/ACCTLOAD.cbl#L163) |
| WS-RA-VALUE | [WS-LOADED-BALANCE](#ws-totals "in WS-TOTALS") | MOVE | 8000-WRITE-REPORT | [ACCTLOAD.cbl:164](../../../pangaea/programs/ACCTLOAD.cbl#L164) |

### Into files

Records written to files.

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| file ONLINE-OUT | [ACCT-ONLINE-RECORD](#acct-online-record) | WRITE | 2100-WRITE-ONLINE | [ACCTLOAD.cbl:134](../../../pangaea/programs/ACCTLOAD.cbl#L134) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [ACCTLOAD.cbl:149](../../../pangaea/programs/ACCTLOAD.cbl#L149) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [ACCTLOAD.cbl:151](../../../pangaea/programs/ACCTLOAD.cbl#L151) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [ACCTLOAD.cbl:154](../../../pangaea/programs/ACCTLOAD.cbl#L154) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [ACCTLOAD.cbl:157](../../../pangaea/programs/ACCTLOAD.cbl#L157) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [ACCTLOAD.cbl:160](../../../pangaea/programs/ACCTLOAD.cbl#L160) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [ACCTLOAD.cbl:162](../../../pangaea/programs/ACCTLOAD.cbl#L162) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | 8000-WRITE-REPORT | [ACCTLOAD.cbl:165](../../../pangaea/programs/ACCTLOAD.cbl#L165) |

### Into other data items

Data items the program does not declare, such as special registers and index names.

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| RETURN-CODE | literal 0 | MOVE | 0000-MAIN | [ACCTLOAD.cbl:86](../../../pangaea/programs/ACCTLOAD.cbl#L86) |
| RETURN-CODE | literal 16 | MOVE | 1000-INITIALIZE | [ACCTLOAD.cbl:95](../../../pangaea/programs/ACCTLOAD.cbl#L95) |
| RETURN-CODE | literal 16 | MOVE | 9900-ABORT | [ACCTLOAD.cbl:175](../../../pangaea/programs/ACCTLOAD.cbl#L175) |

## Field origins

Where the values of each field the program writes to a file come from: the fields of the files it reads, literals, and initial values, followed through every statement that moves a value, and through groups and redefinitions by the bytes they share.

### File ONLINE-OUT

Written from [ACCT-ONLINE-RECORD](#acct-online-record).

| Field | Position | Origins |
| --- | --: | --- |
| AO-ACCT-NO | 1-10 | file ACCT-IN AM-ACCT-NO |
| AO-CUST-NO | 11-18 | file ACCT-IN AM-CUST-NO |
| AO-SHORT-NAME | 19-38 | file ACCT-IN AM-SHORT-NAME |
| AO-ACCT-TYPE | 39 | file ACCT-IN AM-ACCT-TYPE |
| AO-STATUS | 40 | file ACCT-IN AM-STATUS |
| AO-BALANCE | 41-53 | file ACCT-IN AM-BALANCE |
| AO-OVERDRAFT-LIMIT | 54-62 | file ACCT-IN AM-OVERDRAFT-LIMIT |
| AO-HOLDS | 63-75 | literal 0 |
| AO-AS-OF-DATE | 76-83 | file CONTROL-FILE BC-BUSINESS-DATE |

### File REPORT-OUT

Written from [WS-REPORT-AMOUNT-LINE](#ws-report-amount-line).

| Field | Position | Origins |
| --- | --: | --- |
| WS-RA-LABEL | 1-30 | literal 'BALANCES LOADED' |
| WS-RA-VALUE | 31-51 | file ACCT-IN AM-BALANCE; initial value of WS-LOADED-BALANCE |

## Records

### BANK-CONTROL-RECORD

File CONTROL-FILE. 80 bytes, layout in copybook [BANKCTL](../copybooks/BANKCTL.md).

### ACCT-MASTER-RECORD

File ACCT-IN. 100 bytes, layout in copybook [ACCTMAST](../copybooks/ACCTMAST.md).

### ACCT-ONLINE-RECORD

File ONLINE-OUT. 100 bytes, layout in copybook [ACCTONL](../copybooks/ACCTONL.md).

### REPORT-LINE

File REPORT-OUT. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | REPORT-LINE | 1 | 80 | X(80) | | | [ACCTLOAD.cbl:50](../../../pangaea/programs/ACCTLOAD.cbl#L50) |

### WS-FILE-STATUSES

Working storage. 8 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-FILE-STATUSES | 1 | 8 | | | | [ACCTLOAD.cbl:53](../../../pangaea/programs/ACCTLOAD.cbl#L53) |
| 05 | WS-CONTROL-STATUS | 1 | 2 | XX | | | [ACCTLOAD.cbl:54](../../../pangaea/programs/ACCTLOAD.cbl#L54) |
| 05 | WS-ACCT-IN-STATUS | 3 | 2 | XX | | | [ACCTLOAD.cbl:55](../../../pangaea/programs/ACCTLOAD.cbl#L55) |
| 05 | WS-ONLINE-STATUS | 5 | 2 | XX | | | [ACCTLOAD.cbl:56](../../../pangaea/programs/ACCTLOAD.cbl#L56) |
| 05 | WS-REPORT-STATUS | 7 | 2 | XX | | | [ACCTLOAD.cbl:57](../../../pangaea/programs/ACCTLOAD.cbl#L57) |

### WS-BUSINESS-DATE

Working storage. 8 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-BUSINESS-DATE | 1 | 8 | 9(8) | | | [ACCTLOAD.cbl:59](../../../pangaea/programs/ACCTLOAD.cbl#L59) |

### WS-ABORT-MESSAGE

Working storage. 40 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-ABORT-MESSAGE | 1 | 40 | X(40) | | | [ACCTLOAD.cbl:60](../../../pangaea/programs/ACCTLOAD.cbl#L60) |

### WS-MASTER-END-SW

Working storage. 1 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-MASTER-END-SW | 1 | 1 | X | | VALUE 'N' | [ACCTLOAD.cbl:61](../../../pangaea/programs/ACCTLOAD.cbl#L61) |
| 88 | MASTER-AT-END | | | | | VALUE 'Y' | [ACCTLOAD.cbl:62](../../../pangaea/programs/ACCTLOAD.cbl#L62) |

### WS-TOTALS

Working storage. 37 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TOTALS | 1 | 37 | | | | [ACCTLOAD.cbl:64](../../../pangaea/programs/ACCTLOAD.cbl#L64) |
| 05 | WS-READ-COUNT | 1 | 8 | 9(8) | | VALUE 0 | [ACCTLOAD.cbl:65](../../../pangaea/programs/ACCTLOAD.cbl#L65) |
| 05 | WS-LOADED-COUNT | 9 | 8 | 9(8) | | VALUE 0 | [ACCTLOAD.cbl:66](../../../pangaea/programs/ACCTLOAD.cbl#L66) |
| 05 | WS-CLOSED-COUNT | 17 | 8 | 9(8) | | VALUE 0 | [ACCTLOAD.cbl:67](../../../pangaea/programs/ACCTLOAD.cbl#L67) |
| 05 | WS-LOADED-BALANCE | 25 | 13 | S9(11)V99 | | VALUE 0 | [ACCTLOAD.cbl:68](../../../pangaea/programs/ACCTLOAD.cbl#L68) |

### WS-REPORT-COUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-COUNT-LINE | 1 | 80 | | | | [ACCTLOAD.cbl:70](../../../pangaea/programs/ACCTLOAD.cbl#L70) |
| 05 | WS-RC-LABEL | 1 | 30 | X(30) | | | [ACCTLOAD.cbl:71](../../../pangaea/programs/ACCTLOAD.cbl#L71) |
| 05 | WS-RC-VALUE | 31 | 10 | ZZ,ZZZ,ZZ9 | | | [ACCTLOAD.cbl:72](../../../pangaea/programs/ACCTLOAD.cbl#L72) |
| 05 | FILLER | 41 | 40 | X(40) | | VALUE SPACES | [ACCTLOAD.cbl:73](../../../pangaea/programs/ACCTLOAD.cbl#L73) |

### WS-REPORT-AMOUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-AMOUNT-LINE | 1 | 80 | | | | [ACCTLOAD.cbl:74](../../../pangaea/programs/ACCTLOAD.cbl#L74) |
| 05 | WS-RA-LABEL | 1 | 30 | X(30) | | | [ACCTLOAD.cbl:75](../../../pangaea/programs/ACCTLOAD.cbl#L75) |
| 05 | WS-RA-VALUE | 31 | 21 | Z,ZZZ,ZZZ,ZZZ,ZZ9.99- | | | [ACCTLOAD.cbl:76](../../../pangaea/programs/ACCTLOAD.cbl#L76) |
| 05 | FILLER | 52 | 29 | X(29) | | VALUE SPACES | [ACCTLOAD.cbl:77](../../../pangaea/programs/ACCTLOAD.cbl#L77) |
