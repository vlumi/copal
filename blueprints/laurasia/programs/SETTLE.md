# SETTLE

```text
settle: work out what each merchant is paid, and what the network is

reads the priced presentments and the dispute adjustments, both in
merchant order. our own merchants, those our acquiring signed, numbered
lcs..., are paid their matched presentments, less their interchange,
plus or minus their adjustments; one that ends up owing is not paid,
and is listed for collection. every other merchant is another
acquirer's: what is due for them is settled with the network, all of it
in one sum.

return codes: 0 all paid, 4 some owe.

history
2011-04-04 rk  lcs-0102 original version.
2015-03-23 mt  lcs-0302 dispute adjustments.
2024-06-03 ah  lcs-0561 our own merchants paid, the rest through the network.
```

Program [settle.cbl](../../../laurasia/programs/settle.cbl): 5 files, 8 paragraphs. Its business rules: [6 decision tables](../rules/SETTLE.md).

## Contents

- [Files](#files)
- [Flow](#flow)
- [Paragraphs](#paragraphs)
- [Data flow](#data-flow)
- [Lineage](#lineage)
  - [Into PRES-RECORD](#into-pres-record)
  - [Into ADJUSTMENT](#into-adjustment)
  - [Into PAYOUT](#into-payout)
  - [Into NETWORK-SETTLEMENT](#into-network-settlement)
  - [Into REPORT-LINE](#into-report-line)
  - [Into WS-PRES-KEY](#into-ws-pres-key)
  - [Into WS-ADJUST-KEY](#into-ws-adjust-key)
  - [Into WS-MERCHANT](#into-ws-merchant)
  - [Into WS-MERCHANT-TOTALS](#into-ws-merchant-totals)
  - [Into WS-TOTALS](#into-ws-totals)
  - [Into WS-NETWORK-TOTALS](#into-ws-network-totals)
  - [Into WS-REPORT-AMOUNT](#into-ws-report-amount)
  - [Into WS-REPORT-COUNT](#into-ws-report-count)
  - [Out of the program](#out-of-the-program)
- [Field origins](#field-origins)
  - [File PAYOUT-OUT](#file-payout-out)
  - [File NETWORK-OUT](#file-network-out)
  - [File REPORT-OUT](#file-report-out)
- [Links between records](#links-between-records)
- [Records](#records)
  - [PRES-RECORD](#pres-record)
  - [ADJUSTMENT](#adjustment)
  - [PAYOUT](#payout)
  - [NETWORK-SETTLEMENT](#network-settlement)
  - [REPORT-LINE](#report-line)
  - [WS-PRES-KEY](#ws-pres-key)
  - [WS-ADJUST-KEY](#ws-adjust-key)
  - [WS-MERCHANT](#ws-merchant)
  - [WS-MERCHANT-TOTALS](#ws-merchant-totals)
  - [WS-TOTALS](#ws-totals)
  - [WS-NETWORK-TOTALS](#ws-network-totals)
  - [WS-REPORT-AMOUNT](#ws-report-amount)
  - [WS-REPORT-COUNT](#ws-report-count)

## Files

| File | Kind | Assigned to | Organization | Status field | Opened | Records |
| --- | --- | --- | --- | --- | --- | --- |
| PRES-IN | FD | PRESIN | LINE SEQUENTIAL | | INPUT | PRES-RECORD |
| ADJUST-IN | FD | ADJIN | LINE SEQUENTIAL | | INPUT | ADJUSTMENT |
| PAYOUT-OUT | FD | PAYOUT | LINE SEQUENTIAL | | OUTPUT | PAYOUT |
| NETWORK-OUT | FD | NETSETL | LINE SEQUENTIAL | | OUTPUT | NETWORK-SETTLEMENT |
| REPORT-OUT | FD | SETLRPT | LINE SEQUENTIAL | | OUTPUT | REPORT-LINE |

## Flow

![The paragraphs of SETTLE and how control moves between them](SETTLE.svg)

## Paragraphs

| Paragraph | Description | Source | Performs | Goes to | Performed by | Reached by |
| --- | --- | --- | --- | --- | --- | --- |
| MAIN | | [settle.cbl:92](../../../laurasia/programs/settle.cbl#L92) | READ-PRESENTMENT, READ-ADJUSTMENT, SETTLE-MERCHANT, WRITE-NETWORK-SETTLEMENT, WRITE-REPORT | | | |
| READ-PRESENTMENT | | [settle.cbl:106](../../../laurasia/programs/settle.cbl#L106) | | | MAIN, SETTLE-MERCHANT | |
| READ-ADJUSTMENT | | [settle.cbl:112](../../../laurasia/programs/settle.cbl#L112) | | | MAIN, SETTLE-MERCHANT | |
| SETTLE-MERCHANT | the next merchant: the lower of the two keys, and everything of both files for it. | [settle.cbl:120](../../../laurasia/programs/settle.cbl#L120) | READ-PRESENTMENT, READ-ADJUSTMENT, WRITE-PAYOUT, ADD-TO-NETWORK | | MAIN | |
| ADD-TO-NETWORK | another acquirer's merchant: what is due goes into the network's sum. | [settle.cbl:144](../../../laurasia/programs/settle.cbl#L144) | | | SETTLE-MERCHANT | |
| WRITE-PAYOUT | | [settle.cbl:150](../../../laurasia/programs/settle.cbl#L150) | | | SETTLE-MERCHANT | |
| WRITE-NETWORK-SETTLEMENT | | [settle.cbl:166](../../../laurasia/programs/settle.cbl#L166) | | | MAIN | |
| WRITE-REPORT | | [settle.cbl:172](../../../laurasia/programs/settle.cbl#L172) | | | MAIN | |

## Data flow

| Paragraph | Reads | Writes |
| --- | --- | --- |
| MAIN | WS-PRES-KEY, WS-ADJUST-KEY, WS-OWING-COUNT | RETURN-CODE |
| READ-PRESENTMENT | PS-MERCHANT | PRES-RECORD, WS-PRES-KEY |
| READ-ADJUSTMENT | AD-MERCHANT | ADJUSTMENT, WS-ADJUST-KEY |
| SETTLE-MERCHANT | WS-PRES-KEY, WS-ADJUST-KEY, WS-MERCHANT, PS-MATCH, PS-AMOUNT, WS-GROSS, PS-FEE, WS-INTERCHANGE, AD-AMOUNT, WS-ADJUSTMENTS, WS-MERCHANT-RANGE | WS-MERCHANT, WS-GROSS, WS-INTERCHANGE, WS-ADJUSTMENTS |
| ADD-TO-NETWORK | WS-NETWORK-COUNT, WS-GROSS, WS-NETWORK-GROSS, WS-INTERCHANGE, WS-NETWORK-INTERCHANGE, WS-ADJUSTMENTS, WS-NETWORK-ADJUSTMENTS | WS-NETWORK-COUNT, WS-NETWORK-GROSS, WS-NETWORK-INTERCHANGE, WS-NETWORK-ADJUSTMENTS |
| WRITE-PAYOUT | WS-MERCHANT, WS-GROSS, WS-INTERCHANGE, WS-ADJUSTMENTS, PA-NET, WS-OWING-COUNT, WS-PAID-COUNT, WS-PAID-TOTAL | PA-MERCHANT, PA-GROSS, PA-INTERCHANGE, PA-ADJUSTMENTS, PA-NET, PA-STATUS, WS-OWING-COUNT, WS-PAID-COUNT, WS-PAID-TOTAL, PAYOUT |
| WRITE-NETWORK-SETTLEMENT | WS-NETWORK-GROSS, WS-NETWORK-INTERCHANGE, WS-NETWORK-ADJUSTMENTS | NS-GROSS, NS-INTERCHANGE, NS-ADJUSTMENTS, NS-NET, NETWORK-SETTLEMENT |
| WRITE-REPORT | WS-PAID-COUNT, WS-REPORT-COUNT, WS-OWING-COUNT, WS-PAID-TOTAL, WS-REPORT-AMOUNT, WS-NETWORK-COUNT, NS-NET | REPORT-LINE, RC-LABEL, RC-VALUE, RA-LABEL, RA-VALUE |

## Lineage

Where each data item's values come from, grouped by the record they go into. A data item links to the record it belongs to.

### Into PRES-RECORD

Record [PRES-RECORD](#pres-record), from copybook [presrec](../copybooks/presrec.md).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| PRES-RECORD | file PRES-IN | READ | READ-PRESENTMENT | [settle.cbl:107](../../../laurasia/programs/settle.cbl#L107) |

### Into ADJUSTMENT

Record [ADJUSTMENT](#adjustment).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| ADJUSTMENT | file ADJUST-IN | READ | READ-ADJUSTMENT | [settle.cbl:113](../../../laurasia/programs/settle.cbl#L113) |

### Into PAYOUT

Record [PAYOUT](#payout).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| PA-ADJUSTMENTS | [WS-ADJUSTMENTS](#ws-merchant-totals "in WS-MERCHANT-TOTALS") | MOVE | WRITE-PAYOUT | [settle.cbl:154](../../../laurasia/programs/settle.cbl#L154) |
| PA-GROSS | [WS-GROSS](#ws-merchant-totals "in WS-MERCHANT-TOTALS") | MOVE | WRITE-PAYOUT | [settle.cbl:152](../../../laurasia/programs/settle.cbl#L152) |
| PA-INTERCHANGE | [WS-INTERCHANGE](#ws-merchant-totals "in WS-MERCHANT-TOTALS") | MOVE | WRITE-PAYOUT | [settle.cbl:153](../../../laurasia/programs/settle.cbl#L153) |
| PA-MERCHANT | [WS-MERCHANT](#ws-merchant) | MOVE | WRITE-PAYOUT | [settle.cbl:151](../../../laurasia/programs/settle.cbl#L151) |
| PA-NET | [WS-GROSS](#ws-merchant-totals "in WS-MERCHANT-TOTALS") | COMPUTE | WRITE-PAYOUT | [settle.cbl:155](../../../laurasia/programs/settle.cbl#L155) |
| PA-NET | [WS-INTERCHANGE](#ws-merchant-totals "in WS-MERCHANT-TOTALS") | COMPUTE | WRITE-PAYOUT | [settle.cbl:155](../../../laurasia/programs/settle.cbl#L155) |
| PA-NET | [WS-ADJUSTMENTS](#ws-merchant-totals "in WS-MERCHANT-TOTALS") | COMPUTE | WRITE-PAYOUT | [settle.cbl:155](../../../laurasia/programs/settle.cbl#L155) |
| PA-STATUS | literal condition PA-OWES | SET | WRITE-PAYOUT | [settle.cbl:157](../../../laurasia/programs/settle.cbl#L157) |
| PA-STATUS | literal condition PA-PAID | SET | WRITE-PAYOUT | [settle.cbl:160](../../../laurasia/programs/settle.cbl#L160) |

### Into NETWORK-SETTLEMENT

Record [NETWORK-SETTLEMENT](#network-settlement).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| NS-ADJUSTMENTS | [WS-NETWORK-ADJUSTMENTS](#ws-network-totals "in WS-NETWORK-TOTALS") | MOVE | WRITE-NETWORK-SETTLEMENT | [settle.cbl:169](../../../laurasia/programs/settle.cbl#L169) |
| NS-GROSS | [WS-NETWORK-GROSS](#ws-network-totals "in WS-NETWORK-TOTALS") | MOVE | WRITE-NETWORK-SETTLEMENT | [settle.cbl:167](../../../laurasia/programs/settle.cbl#L167) |
| NS-INTERCHANGE | [WS-NETWORK-INTERCHANGE](#ws-network-totals "in WS-NETWORK-TOTALS") | MOVE | WRITE-NETWORK-SETTLEMENT | [settle.cbl:168](../../../laurasia/programs/settle.cbl#L168) |
| NS-NET | [WS-NETWORK-GROSS](#ws-network-totals "in WS-NETWORK-TOTALS") | COMPUTE | WRITE-NETWORK-SETTLEMENT | [settle.cbl:170](../../../laurasia/programs/settle.cbl#L170) |
| NS-NET | [WS-NETWORK-INTERCHANGE](#ws-network-totals "in WS-NETWORK-TOTALS") | COMPUTE | WRITE-NETWORK-SETTLEMENT | [settle.cbl:170](../../../laurasia/programs/settle.cbl#L170) |
| NS-NET | [WS-NETWORK-ADJUSTMENTS](#ws-network-totals "in WS-NETWORK-TOTALS") | COMPUTE | WRITE-NETWORK-SETTLEMENT | [settle.cbl:170](../../../laurasia/programs/settle.cbl#L170) |

### Into REPORT-LINE

Record [REPORT-LINE](#report-line).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| REPORT-LINE | literal 'SETTLE   MERCHANT SETTLEMENT' | MOVE | WRITE-REPORT | [settle.cbl:173](../../../laurasia/programs/settle.cbl#L173) |
| REPORT-LINE | [WS-REPORT-COUNT](#ws-report-count) | WRITE | WRITE-REPORT | [settle.cbl:177](../../../laurasia/programs/settle.cbl#L177) |
| REPORT-LINE | [WS-REPORT-COUNT](#ws-report-count) | WRITE | WRITE-REPORT | [settle.cbl:180](../../../laurasia/programs/settle.cbl#L180) |
| REPORT-LINE | [WS-REPORT-AMOUNT](#ws-report-amount) | WRITE | WRITE-REPORT | [settle.cbl:183](../../../laurasia/programs/settle.cbl#L183) |
| REPORT-LINE | [WS-REPORT-COUNT](#ws-report-count) | WRITE | WRITE-REPORT | [settle.cbl:186](../../../laurasia/programs/settle.cbl#L186) |
| REPORT-LINE | [WS-REPORT-AMOUNT](#ws-report-amount) | WRITE | WRITE-REPORT | [settle.cbl:189](../../../laurasia/programs/settle.cbl#L189) |

### Into WS-PRES-KEY

Record [WS-PRES-KEY](#ws-pres-key).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-PRES-KEY | literal HIGH-VALUES | MOVE | READ-PRESENTMENT | [settle.cbl:108](../../../laurasia/programs/settle.cbl#L108) |
| WS-PRES-KEY | [PS-MERCHANT](#pres-record "in PRES-RECORD") | MOVE | READ-PRESENTMENT | [settle.cbl:109](../../../laurasia/programs/settle.cbl#L109) |

### Into WS-ADJUST-KEY

Record [WS-ADJUST-KEY](#ws-adjust-key).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-ADJUST-KEY | literal HIGH-VALUES | MOVE | READ-ADJUSTMENT | [settle.cbl:114](../../../laurasia/programs/settle.cbl#L114) |
| WS-ADJUST-KEY | [AD-MERCHANT](#adjustment "in ADJUSTMENT") | MOVE | READ-ADJUSTMENT | [settle.cbl:115](../../../laurasia/programs/settle.cbl#L115) |

### Into WS-MERCHANT

Record [WS-MERCHANT](#ws-merchant).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-MERCHANT | [WS-PRES-KEY](#ws-pres-key) | MOVE | SETTLE-MERCHANT | [settle.cbl:122](../../../laurasia/programs/settle.cbl#L122) |
| WS-MERCHANT | [WS-ADJUST-KEY](#ws-adjust-key) | MOVE | SETTLE-MERCHANT | [settle.cbl:124](../../../laurasia/programs/settle.cbl#L124) |

### Into WS-MERCHANT-TOTALS

Record [WS-MERCHANT-TOTALS](#ws-merchant-totals).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-ADJUSTMENTS | literal ZERO | MOVE | SETTLE-MERCHANT | [settle.cbl:126](../../../laurasia/programs/settle.cbl#L126) |
| WS-ADJUSTMENTS | [AD-AMOUNT](#adjustment "in ADJUSTMENT") | ADD | SETTLE-MERCHANT | [settle.cbl:135](../../../laurasia/programs/settle.cbl#L135) |
| WS-GROSS | literal ZERO | MOVE | SETTLE-MERCHANT | [settle.cbl:126](../../../laurasia/programs/settle.cbl#L126) |
| WS-GROSS | [PS-AMOUNT](#pres-record "in PRES-RECORD") | ADD | SETTLE-MERCHANT | [settle.cbl:129](../../../laurasia/programs/settle.cbl#L129) |
| WS-INTERCHANGE | literal ZERO | MOVE | SETTLE-MERCHANT | [settle.cbl:126](../../../laurasia/programs/settle.cbl#L126) |
| WS-INTERCHANGE | [PS-FEE](#pres-record "in PRES-RECORD") | ADD | SETTLE-MERCHANT | [settle.cbl:130](../../../laurasia/programs/settle.cbl#L130) |

### Into WS-TOTALS

Record [WS-TOTALS](#ws-totals).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-NETWORK-COUNT | literal 1 | ADD | ADD-TO-NETWORK | [settle.cbl:145](../../../laurasia/programs/settle.cbl#L145) |
| WS-OWING-COUNT | literal 1 | ADD | WRITE-PAYOUT | [settle.cbl:158](../../../laurasia/programs/settle.cbl#L158) |
| WS-PAID-COUNT | literal 1 | ADD | WRITE-PAYOUT | [settle.cbl:161](../../../laurasia/programs/settle.cbl#L161) |
| WS-PAID-TOTAL | [PA-NET](#payout "in PAYOUT") | ADD | WRITE-PAYOUT | [settle.cbl:162](../../../laurasia/programs/settle.cbl#L162) |

### Into WS-NETWORK-TOTALS

Record [WS-NETWORK-TOTALS](#ws-network-totals).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| WS-NETWORK-ADJUSTMENTS | [WS-ADJUSTMENTS](#ws-merchant-totals "in WS-MERCHANT-TOTALS") | ADD | ADD-TO-NETWORK | [settle.cbl:148](../../../laurasia/programs/settle.cbl#L148) |
| WS-NETWORK-GROSS | [WS-GROSS](#ws-merchant-totals "in WS-MERCHANT-TOTALS") | ADD | ADD-TO-NETWORK | [settle.cbl:146](../../../laurasia/programs/settle.cbl#L146) |
| WS-NETWORK-INTERCHANGE | [WS-INTERCHANGE](#ws-merchant-totals "in WS-MERCHANT-TOTALS") | ADD | ADD-TO-NETWORK | [settle.cbl:147](../../../laurasia/programs/settle.cbl#L147) |

### Into WS-REPORT-AMOUNT

Record [WS-REPORT-AMOUNT](#ws-report-amount).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| RA-LABEL | literal 'PAID' | MOVE | WRITE-REPORT | [settle.cbl:181](../../../laurasia/programs/settle.cbl#L181) |
| RA-LABEL | literal 'OWED TO THE NETWORK' | MOVE | WRITE-REPORT | [settle.cbl:187](../../../laurasia/programs/settle.cbl#L187) |
| RA-VALUE | [WS-PAID-TOTAL](#ws-totals "in WS-TOTALS") | MOVE | WRITE-REPORT | [settle.cbl:182](../../../laurasia/programs/settle.cbl#L182) |
| RA-VALUE | [NS-NET](#network-settlement "in NETWORK-SETTLEMENT") | MOVE | WRITE-REPORT | [settle.cbl:188](../../../laurasia/programs/settle.cbl#L188) |

### Into WS-REPORT-COUNT

Record [WS-REPORT-COUNT](#ws-report-count).

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| RC-LABEL | literal 'MERCHANTS PAID' | MOVE | WRITE-REPORT | [settle.cbl:175](../../../laurasia/programs/settle.cbl#L175) |
| RC-LABEL | literal 'MERCHANTS OWING' | MOVE | WRITE-REPORT | [settle.cbl:178](../../../laurasia/programs/settle.cbl#L178) |
| RC-LABEL | literal 'OTHER ACQUIRERS' MERCHANTS' | MOVE | WRITE-REPORT | [settle.cbl:184](../../../laurasia/programs/settle.cbl#L184) |
| RC-VALUE | [WS-PAID-COUNT](#ws-totals "in WS-TOTALS") | MOVE | WRITE-REPORT | [settle.cbl:176](../../../laurasia/programs/settle.cbl#L176) |
| RC-VALUE | [WS-OWING-COUNT](#ws-totals "in WS-TOTALS") | MOVE | WRITE-REPORT | [settle.cbl:179](../../../laurasia/programs/settle.cbl#L179) |
| RC-VALUE | [WS-NETWORK-COUNT](#ws-totals "in WS-TOTALS") | MOVE | WRITE-REPORT | [settle.cbl:185](../../../laurasia/programs/settle.cbl#L185) |

### Out of the program

Records written to files, items passed to the programs called, and items displayed to the job's output.

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| file NETWORK-OUT | [NETWORK-SETTLEMENT](#network-settlement) | WRITE | WRITE-NETWORK-SETTLEMENT | [settle.cbl:171](../../../laurasia/programs/settle.cbl#L171) |
| file PAYOUT-OUT | [PAYOUT](#payout) | WRITE | WRITE-PAYOUT | [settle.cbl:164](../../../laurasia/programs/settle.cbl#L164) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | WRITE-REPORT | [settle.cbl:174](../../../laurasia/programs/settle.cbl#L174) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | WRITE-REPORT | [settle.cbl:177](../../../laurasia/programs/settle.cbl#L177) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | WRITE-REPORT | [settle.cbl:180](../../../laurasia/programs/settle.cbl#L180) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | WRITE-REPORT | [settle.cbl:183](../../../laurasia/programs/settle.cbl#L183) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | WRITE-REPORT | [settle.cbl:186](../../../laurasia/programs/settle.cbl#L186) |
| file REPORT-OUT | [REPORT-LINE](#report-line) | WRITE | WRITE-REPORT | [settle.cbl:189](../../../laurasia/programs/settle.cbl#L189) |

## Field origins

Where the values of each field the program writes to a file come from: the fields of the files it reads, literals, and initial values, followed through every statement that moves a value, and through groups and redefinitions by the bytes they share.

### File PAYOUT-OUT

Written from [PAYOUT](#payout).

| Field | Position | Origins |
| --- | --: | --- |
| PA-MERCHANT | 1-10 | file ADJUST-IN AD-MERCHANT; file PRES-IN PS-MERCHANT; literal HIGH-VALUES |
| PA-GROSS | 11-24 | file PRES-IN PS-AMOUNT; literal ZERO |
| PA-INTERCHANGE | 25-36 | file PRES-IN PS-FEE; literal ZERO |
| PA-ADJUSTMENTS | 37-50 | file ADJUST-IN AD-AMOUNT; literal ZERO |
| PA-NET | 51-64 | file ADJUST-IN AD-AMOUNT; file PRES-IN PS-AMOUNT; file PRES-IN PS-FEE; literal ZERO |
| PA-STATUS | 65 | literal condition PA-OWES; literal condition PA-PAID |

### File NETWORK-OUT

Written from [NETWORK-SETTLEMENT](#network-settlement).

| Field | Position | Origins |
| --- | --: | --- |
| NS-GROSS | 1-16 | file PRES-IN PS-AMOUNT; initial value of WS-NETWORK-GROSS; literal ZERO |
| NS-INTERCHANGE | 17-30 | file PRES-IN PS-FEE; initial value of WS-NETWORK-INTERCHANGE; literal ZERO |
| NS-ADJUSTMENTS | 31-46 | file ADJUST-IN AD-AMOUNT; initial value of WS-NETWORK-ADJUSTMENTS; literal ZERO |
| NS-NET | 47-62 | file ADJUST-IN AD-AMOUNT; file PRES-IN PS-AMOUNT; file PRES-IN PS-FEE; initial value of WS-NETWORK-ADJUSTMENTS; initial value of WS-NETWORK-GROSS; initial value of WS-NETWORK-INTERCHANGE; literal ZERO |

### File REPORT-OUT

Written from [REPORT-LINE](#report-line), [WS-REPORT-COUNT](#ws-report-count) and [WS-REPORT-AMOUNT](#ws-report-amount).

#### Record REPORT-LINE

| Field | Position | Origins |
| --- | --: | --- |
| REPORT-LINE | 1-60 | literal 'SETTLE   MERCHANT SETTLEMENT' |

#### Record WS-REPORT-COUNT

| Field | Position | Origins |
| --- | --: | --- |
| RC-LABEL | 1-30 | literal 'MERCHANTS OWING'; literal 'MERCHANTS PAID'; literal 'OTHER ACQUIRERS' MERCHANTS' |
| RC-VALUE | 31-36 | initial value of WS-NETWORK-COUNT; initial value of WS-OWING-COUNT; initial value of WS-PAID-COUNT; literal 1 |

#### Record WS-REPORT-AMOUNT

| Field | Position | Origins |
| --- | --: | --- |
| RA-LABEL | 1-30 | literal 'OWED TO THE NETWORK'; literal 'PAID' |
| RA-VALUE | 31-45 | file ADJUST-IN AD-AMOUNT; file PRES-IN PS-AMOUNT; file PRES-IN PS-FEE; initial value of WS-NETWORK-ADJUSTMENTS; initial value of WS-NETWORK-GROSS; initial value of WS-NETWORK-INTERCHANGE; initial value of WS-PAID-TOTAL; literal ZERO |

## Links between records

Fields of one input the program relies on to find or match the records of another: moved into an indexed file's record key to read it by, compared for equality, or compared for order to merge the two in that order, each side's values from that field alone.

| Field | Of | Links to | Of | How | Source |
| --- | --- | --- | --- | --- | --- |
| PS-MERCHANT | PRES-IN | AD-MERCHANT | ADJUST-IN | merged in order with | [settle.cbl:121](../../../laurasia/programs/settle.cbl#L121) |

## Records

### PRES-RECORD

File PRES-IN. 76 bytes, layout in copybook [presrec](../copybooks/presrec.md).

### ADJUSTMENT

File ADJUST-IN. 32 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | ADJUSTMENT | group | 1 | 32 | | | | [settle.cbl:40](../../../laurasia/programs/settle.cbl#L40) |
| 05 | ├ AD-MERCHANT | text | 1 | 10 | X(10) | | | [settle.cbl:41](../../../laurasia/programs/settle.cbl#L41) |
| 05 | ├ AD-AMOUNT | amount | 11 | 12 | S9(9)V99 | | SIGN LEADING SEPARATE | [settle.cbl:42](../../../laurasia/programs/settle.cbl#L42) |
| 05 | ├ AD-CASE-NO | number | 23 | 8 | 9(8) | | | [settle.cbl:43](../../../laurasia/programs/settle.cbl#L43) |
| 05 | └ AD-EVENT | text | 31 | 2 | X(2) | | | [settle.cbl:44](../../../laurasia/programs/settle.cbl#L44) |

### PAYOUT

File PAYOUT-OUT. 65 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | PAYOUT | group | 1 | 65 | | | | [settle.cbl:46](../../../laurasia/programs/settle.cbl#L46) |
| 05 | ├ PA-MERCHANT | text | 1 | 10 | X(10) | | | [settle.cbl:47](../../../laurasia/programs/settle.cbl#L47) |
| 05 | ├ PA-GROSS | amount | 11 | 14 | S9(11)V99 | | SIGN LEADING SEPARATE | [settle.cbl:48](../../../laurasia/programs/settle.cbl#L48) |
| 05 | ├ PA-INTERCHANGE | amount | 25 | 12 | S9(9)V99 | | SIGN LEADING SEPARATE | [settle.cbl:49](../../../laurasia/programs/settle.cbl#L49) |
| 05 | ├ PA-ADJUSTMENTS | amount | 37 | 14 | S9(11)V99 | | SIGN LEADING SEPARATE | [settle.cbl:50](../../../laurasia/programs/settle.cbl#L50) |
| 05 | ├ PA-NET | amount | 51 | 14 | S9(11)V99 | | SIGN LEADING SEPARATE | [settle.cbl:51](../../../laurasia/programs/settle.cbl#L51) |
| 05 | └ PA-STATUS | code | 65 | 1 | X | | | [settle.cbl:52](../../../laurasia/programs/settle.cbl#L52) |
| 88 |   ├ PA-PAID | | | | | | VALUE 'P' | [settle.cbl:53](../../../laurasia/programs/settle.cbl#L53) |
| 88 |   └ PA-OWES | | | | | | VALUE 'C' | [settle.cbl:54](../../../laurasia/programs/settle.cbl#L54) |

### NETWORK-SETTLEMENT

File NETWORK-OUT. 62 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | NETWORK-SETTLEMENT | group | 1 | 62 | | | | [settle.cbl:56](../../../laurasia/programs/settle.cbl#L56) |
| 05 | ├ NS-GROSS | amount | 1 | 16 | S9(13)V99 | | SIGN LEADING SEPARATE | [settle.cbl:57](../../../laurasia/programs/settle.cbl#L57) |
| 05 | ├ NS-INTERCHANGE | amount | 17 | 14 | S9(11)V99 | | SIGN LEADING SEPARATE | [settle.cbl:58](../../../laurasia/programs/settle.cbl#L58) |
| 05 | ├ NS-ADJUSTMENTS | amount | 31 | 16 | S9(13)V99 | | SIGN LEADING SEPARATE | [settle.cbl:59](../../../laurasia/programs/settle.cbl#L59) |
| 05 | └ NS-NET | amount | 47 | 16 | S9(13)V99 | | SIGN LEADING SEPARATE | [settle.cbl:60](../../../laurasia/programs/settle.cbl#L60) |

### REPORT-LINE

File REPORT-OUT. 60 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | REPORT-LINE | text | 1 | 60 | X(60) | | | [settle.cbl:62](../../../laurasia/programs/settle.cbl#L62) |

### WS-PRES-KEY

Working storage. 10 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-PRES-KEY | text | 1 | 10 | X(10) | | | [settle.cbl:65](../../../laurasia/programs/settle.cbl#L65) |

### WS-ADJUST-KEY

Working storage. 10 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-ADJUST-KEY | text | 1 | 10 | X(10) | | | [settle.cbl:66](../../../laurasia/programs/settle.cbl#L66) |

### WS-MERCHANT

Working storage. 10 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-MERCHANT | group | 1 | 10 | | | | [settle.cbl:67](../../../laurasia/programs/settle.cbl#L67) |
| 05 | ├ WS-MERCHANT-RANGE | code | 1 | 3 | X(3) | | | [settle.cbl:68](../../../laurasia/programs/settle.cbl#L68) |
| 88 | │ └ WS-OUR-MERCHANT | | | | | | VALUE 'LCS' | [settle.cbl:69](../../../laurasia/programs/settle.cbl#L69) |
| 05 | └ FILLER | filler | 4 | 7 | X(7) | | | [settle.cbl:70](../../../laurasia/programs/settle.cbl#L70) |

### WS-MERCHANT-TOTALS

Working storage. 37 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-MERCHANT-TOTALS | group | 1 | 37 | | | | [settle.cbl:71](../../../laurasia/programs/settle.cbl#L71) |
| 05 | ├ WS-GROSS | amount | 1 | 13 | S9(11)V99 | | | [settle.cbl:72](../../../laurasia/programs/settle.cbl#L72) |
| 05 | ├ WS-INTERCHANGE | amount | 14 | 11 | S9(9)V99 | | | [settle.cbl:73](../../../laurasia/programs/settle.cbl#L73) |
| 05 | └ WS-ADJUSTMENTS | amount | 25 | 13 | S9(11)V99 | | | [settle.cbl:74](../../../laurasia/programs/settle.cbl#L74) |

### WS-TOTALS

Working storage. 33 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TOTALS | group | 1 | 33 | | | | [settle.cbl:75](../../../laurasia/programs/settle.cbl#L75) |
| 05 | ├ WS-PAID-COUNT | number | 1 | 6 | 9(6) | | VALUE ZERO | [settle.cbl:76](../../../laurasia/programs/settle.cbl#L76) |
| 05 | ├ WS-OWING-COUNT | number | 7 | 6 | 9(6) | | VALUE ZERO | [settle.cbl:77](../../../laurasia/programs/settle.cbl#L77) |
| 05 | ├ WS-PAID-TOTAL | amount | 13 | 15 | S9(13)V99 | | VALUE ZERO | [settle.cbl:78](../../../laurasia/programs/settle.cbl#L78) |
| 05 | └ WS-NETWORK-COUNT | number | 28 | 6 | 9(6) | | VALUE ZERO | [settle.cbl:79](../../../laurasia/programs/settle.cbl#L79) |

### WS-NETWORK-TOTALS

Working storage. 43 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-NETWORK-TOTALS | group | 1 | 43 | | | | [settle.cbl:80](../../../laurasia/programs/settle.cbl#L80) |
| 05 | ├ WS-NETWORK-GROSS | amount | 1 | 15 | S9(13)V99 | | VALUE ZERO | [settle.cbl:81](../../../laurasia/programs/settle.cbl#L81) |
| 05 | ├ WS-NETWORK-INTERCHANGE | amount | 16 | 13 | S9(11)V99 | | VALUE ZERO | [settle.cbl:82](../../../laurasia/programs/settle.cbl#L82) |
| 05 | └ WS-NETWORK-ADJUSTMENTS | amount | 29 | 15 | S9(13)V99 | | VALUE ZERO | [settle.cbl:83](../../../laurasia/programs/settle.cbl#L83) |

### WS-REPORT-AMOUNT

Working storage. 45 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-AMOUNT | group | 1 | 45 | | | | [settle.cbl:84](../../../laurasia/programs/settle.cbl#L84) |
| 05 | ├ RA-LABEL | text | 1 | 30 | X(30) | | | [settle.cbl:85](../../../laurasia/programs/settle.cbl#L85) |
| 05 | └ RA-VALUE | edited | 31 | 15 | -(11)9.99 | | | [settle.cbl:86](../../../laurasia/programs/settle.cbl#L86) |

### WS-REPORT-COUNT

Working storage. 36 bytes.

| Level | Field | Kind | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-COUNT | group | 1 | 36 | | | | [settle.cbl:87](../../../laurasia/programs/settle.cbl#L87) |
| 05 | ├ RC-LABEL | text | 1 | 30 | X(30) | | | [settle.cbl:88](../../../laurasia/programs/settle.cbl#L88) |
| 05 | └ RC-VALUE | edited | 31 | 6 | Z(5)9 | | | [settle.cbl:89](../../../laurasia/programs/settle.cbl#L89) |
