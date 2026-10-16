# TXPOST

```text
TXPOST: POST TRANSACTIONS TO THE ACCOUNT MASTER

MATCHES THE SORTED TRANSACTIONS AGAINST THE OLD MASTER, BOTH
IN ACCOUNT NUMBER ORDER, AND WRITES THE NEW MASTER: EVERY OLD
RECORD, WITH ITS TRANSACTIONS APPLIED. EACH POSTING GOES TO
THE AUDIT TRAIL; A TRANSACTION THAT CANNOT BE POSTED GOES TO
THE REJECTS. BOTH ARE APPENDED TO, AS THE DAILY BATCH RUNS
TXPOST FOUR TIMES:
  0. THE MONTH-END INTEREST AND FEES FROM MONCHG, ON THE DAY
     A MONTH CLOSES; ON OTHER DAYS NONE.
  1. THE DAY'S TRANSACTIONS. A POSTED TRANSFER WRITES A
     TRANSFER CREDIT FOR THE RECEIVING ACCOUNT TO THE LEGS.
  2. THE TRANSFER CREDITS, SORTED. ONE FOR AN UNKNOWN OR
     CLOSED ACCOUNT WRITES A RETURNED TRANSFER TO THE LEGS.
  3. THE RETURNED TRANSFERS, SORTED.

RULES: NOTHING POSTS TO A CLOSED ACCOUNT. CREDITS AND FEES
POST TO A FROZEN ACCOUNT, WITHDRAWALS AND TRANSFERS DO NOT. A
DEBIT MAY TAKE THE BALANCE DOWN TO THE OVERDRAFT LIMIT, NO
FURTHER. A FEE ALWAYS POSTS, EVEN PAST THE LIMIT. INTEREST
POSTS AS A CREDIT, OVERDRAFT INTEREST AND MAINTENANCE FEES AS
FEES.

RETURN CODES: 0 ALL POSTED, 4 SOME REJECTED,
              16 INPUT OUT OF ORDER, OUT OF BALANCE OR I/O.

MODIFICATION HISTORY
DATE       BY  CHANGE
1987-04-02 JA  ORIGINAL VERSION.
1994-11-14 PK  CR-0201 OVERDRAFT LIMITS.
1996-04-15 PK  CR-0233 TRANSFERS, POSTED IN LEGS OVER THREE
               RUNS. AUDIT AND REJECTS ARE APPENDED TO.
1999-01-11 SL  Y2K: 8-DIGIT DATES IN THE AUDIT TRAIL.
2003-05-26 MV  CR-0352 FEES POST EVEN PAST THE OVERDRAFT
               LIMIT.
2001-09-03 TL  CR-0298 MONTH-END INTEREST AND FEES.
2011-08-15 MV  CR-0471 RECONCILE THE NEW MASTER BEFORE
               ENDING THE RUN.
```

Program [TXPOST.cbl](../../../pangaea/programs/TXPOST.cbl): 8 files, 17 paragraphs.

- AUTHOR: J. ARKWRIGHT, BATCH SYSTEMS.
- INSTALLATION: FIRST BANK OF PANGAEA, DATA CENTER.
- DATE-WRITTEN: 1987-04-02.

## Contents

- [Files](#files)
- [Flow](#flow)
- [Paragraphs](#paragraphs)
- [Data flow](#data-flow)
- [Lineage](#lineage)
- [Records](#records)
  - [BANK-CONTROL-RECORD](#bank-control-record)
  - [ACCT-MASTER-RECORD](#acct-master-record)
  - [TRAN-RECORD](#tran-record)
  - [ACCT-OUT-RECORD](#acct-out-record)
  - [AUDIT-OUT-RECORD](#audit-out-record)
  - [REJECT-OUT-RECORD](#reject-out-record)
  - [LEGS-OUT-RECORD](#legs-out-record)
  - [REPORT-LINE](#report-line)
  - [WS-FILE-STATUSES](#ws-file-statuses)
  - [WS-BUSINESS-DATE](#ws-business-date)
  - [WS-ABORT-MESSAGE](#ws-abort-message)
  - [WS-MASTER-KEY](#ws-master-key)
  - [WS-TRAN-KEY](#ws-tran-key)
  - [WS-PREV-MASTER-KEY](#ws-prev-master-key)
  - [WS-TRAN-ORDER-KEY](#ws-tran-order-key)
  - [WS-PREV-TRAN-ORDER-KEY](#ws-prev-tran-order-key)
  - [WS-POSTING](#ws-posting)
  - [WS-SAVED-TRAN](#ws-saved-tran)
  - [WS-LEG-ACCT-NO](#ws-leg-acct-no)
  - [WS-TOTALS](#ws-totals)
  - [AUDIT-RECORD](#audit-record)
  - [REJECT-RECORD](#reject-record)
  - [WS-REPORT-COUNT-LINE](#ws-report-count-line)
  - [WS-REPORT-AMOUNT-LINE](#ws-report-amount-line)

## Files

| File | Kind | Assigned to | Organization | Status field | Opened | Records |
| --- | --- | --- | --- | --- | --- | --- |
| CONTROL-FILE | FD | BANKCTL | SEQUENTIAL | WS-CONTROL-STATUS | INPUT | BANK-CONTROL-RECORD |
| ACCT-IN | FD | ACCTIN | SEQUENTIAL | WS-ACCT-IN-STATUS | INPUT | ACCT-MASTER-RECORD |
| TRAN-IN | FD | TRANIN | SEQUENTIAL | WS-TRAN-IN-STATUS | INPUT | TRAN-RECORD |
| ACCT-OUT | FD | ACCTOUT | SEQUENTIAL | WS-ACCT-OUT-STATUS | OUTPUT | ACCT-OUT-RECORD |
| AUDIT-OUT | FD | AUDIT | SEQUENTIAL | WS-AUDIT-STATUS | EXTEND | AUDIT-OUT-RECORD |
| REJECT-OUT | FD | REJECTS | SEQUENTIAL | WS-REJECT-STATUS | EXTEND | REJECT-OUT-RECORD |
| LEGS-OUT | FD | LEGS | SEQUENTIAL | WS-LEGS-STATUS | OUTPUT | LEGS-OUT-RECORD |
| REPORT-OUT | FD | POSTRPT | LINE SEQUENTIAL | WS-REPORT-STATUS | OUTPUT | REPORT-LINE |

## Flow

![The paragraphs of TXPOST and how control moves between them](TXPOST.svg)

Solid arrows are PERFORMs, dashed ones GO TOs, and dotted ones control running on into the next paragraph inside a PERFORM THRU range.

## Paragraphs

| Paragraph | Description | Source | Performs | Goes to | Performed by | Reached by |
| --- | --- | --- | --- | --- | --- | --- |
| 0000-MAIN | MATCH UNTIL BOTH FILES END, THEN RECONCILE AND REPORT. | [TXPOST.cbl:158](../../../pangaea/programs/TXPOST.cbl#L158) | 1000-INITIALIZE, 2000-MATCH, 5000-RECONCILE, 8000-WRITE-REPORT, 9000-CLOSE-FILES | | | |
| 1000-INITIALIZE | OPEN THE FILES AND READ THE FIRST RECORD OF EACH. THE AUDIT TRAIL AND THE REJECTS ARE EXTENDED, AND CREATED IF NEW. | [TXPOST.cbl:175](../../../pangaea/programs/TXPOST.cbl#L175) | 9900-ABORT, 1200-READ-MASTER, 1300-READ-TRAN | | 0000-MAIN | |
| 1200-READ-MASTER | READ THE NEXT OLD MASTER RECORD, CHECKING ACCOUNT ORDER. | [TXPOST.cbl:202](../../../pangaea/programs/TXPOST.cbl#L202) | 9900-ABORT | | 1000-INITIALIZE, 2000-MATCH | |
| 1300-READ-TRAN | READ THE NEXT TRANSACTION, CHECKING ITS ORDER. | [TXPOST.cbl:219](../../../pangaea/programs/TXPOST.cbl#L219) | 9900-ABORT | | 1000-INITIALIZE, 2000-MATCH | |
| 2000-MATCH | ONE STEP OF THE MATCH, BY WHICHEVER KEY IS LOWER: REJECT A TRANSACTION FOR A MISSING ACCOUNT, APPLY ONE TO ITS ACCOUNT, OR WRITE A FINISHED ACCOUNT. | [TXPOST.cbl:237](../../../pangaea/programs/TXPOST.cbl#L237) | 3900-REJECT-TRAN, 1300-READ-TRAN, 3000-APPLY-TRAN, 4000-WRITE-MASTER, 1200-READ-MASTER | | 0000-MAIN | |
| 3000-APPLY-TRAN | APPLY ONE TRANSACTION TO ITS ACCOUNT BY THE POSTING RULES. | [TXPOST.cbl:252](../../../pangaea/programs/TXPOST.cbl#L252) | 3900-REJECT-TRAN, 3100-POST-CREDIT, 3300-POST-FEE, 3200-POST-DEBIT | | 2000-MATCH | |
| 3100-POST-CREDIT | CREDIT A DEPOSIT, A TRANSFER CREDIT, A RETURNED TRANSFER OR INTEREST. | [TXPOST.cbl:277](../../../pangaea/programs/TXPOST.cbl#L277) | 3800-WRITE-AUDIT | | 3000-APPLY-TRAN | |
| 3200-POST-DEBIT | DEBIT A WITHDRAWAL OR TRANSFER IF IT STAYS WITHIN THE OVERDRAFT LIMIT. A POSTED TRANSFER SENDS ITS CREDIT ON. | [TXPOST.cbl:285](../../../pangaea/programs/TXPOST.cbl#L285) | 3900-REJECT-TRAN, 3800-WRITE-AUDIT, 3500-WRITE-LEG | | 3000-APPLY-TRAN | |
| 3300-POST-FEE | DEBIT A FEE, OVERDRAFT INTEREST OR A MAINTENANCE FEE, WHATEVER THE BALANCE. | [TXPOST.cbl:303](../../../pangaea/programs/TXPOST.cbl#L303) | 3800-WRITE-AUDIT | | 3000-APPLY-TRAN | |
| 3500-WRITE-LEG | WRITE THE NEXT LEG OF A TRANSFER FOR THE COUNTER ACCOUNT: A TRANSFER CREDIT AFTER A TRANSFER, A RETURN AFTER A CREDIT THAT COULD NOT POST. | [TXPOST.cbl:312](../../../pangaea/programs/TXPOST.cbl#L312) | | | 3200-POST-DEBIT, 3900-REJECT-TRAN | |
| 3800-WRITE-AUDIT | RECORD A POSTING IN THE AUDIT TRAIL. | [TXPOST.cbl:327](../../../pangaea/programs/TXPOST.cbl#L327) | | | 3100-POST-CREDIT, 3200-POST-DEBIT, 3300-POST-FEE | |
| 3900-REJECT-TRAN | WRITE THE TRANSACTION TO THE REJECTS WITH ITS REASON. A TRANSFER CREDIT THAT CANNOT POST GOES BACK TO ITS SENDER. | [TXPOST.cbl:341](../../../pangaea/programs/TXPOST.cbl#L341) | 3500-WRITE-LEG | | 2000-MATCH, 3000-APPLY-TRAN, 3200-POST-DEBIT | |
| 4000-WRITE-MASTER | WRITE AN ACCOUNT, ITS TRANSACTIONS APPLIED, TO THE NEW MASTER. | [TXPOST.cbl:351](../../../pangaea/programs/TXPOST.cbl#L351) | | | 2000-MATCH | |
| 5000-RECONCILE | THE NEW MASTER MUST HOLD THE OLD BALANCES PLUS WHAT WAS POSTED. | [TXPOST.cbl:357](../../../pangaea/programs/TXPOST.cbl#L357) | 9900-ABORT | | 0000-MAIN | |
| 8000-WRITE-REPORT | WRITE THE POSTING REPORT. | [TXPOST.cbl:369](../../../pangaea/programs/TXPOST.cbl#L369) | | | 0000-MAIN | |
| 9000-CLOSE-FILES | CLOSE ALL FILES. | [TXPOST.cbl:415](../../../pangaea/programs/TXPOST.cbl#L415) | | | 0000-MAIN, 9900-ABORT | |
| 9900-ABORT | STOP THE RUN WITH RETURN CODE 16, SAYING WHY. | [TXPOST.cbl:420](../../../pangaea/programs/TXPOST.cbl#L420) | 9000-CLOSE-FILES | | 1000-INITIALIZE, 1200-READ-MASTER, 1300-READ-TRAN, 5000-RECONCILE | |

## Data flow

| Paragraph | Reads | Writes |
| --- | --- | --- |
| 0000-MAIN | WS-MASTER-KEY, WS-TRAN-KEY, WS-REJECT-COUNT | RETURN-CODE |
| 1000-INITIALIZE | WS-CONTROL-STATUS, WS-ACCT-IN-STATUS, WS-TRAN-IN-STATUS, WS-ACCT-OUT-STATUS, WS-AUDIT-STATUS, WS-REJECT-STATUS, WS-LEGS-STATUS, WS-REPORT-STATUS, WS-FILE-STATUSES, BC-BUSINESS-DATE | RETURN-CODE, BANK-CONTROL-RECORD, WS-ABORT-MESSAGE, WS-BUSINESS-DATE, AUDIT-RECORD, REJECT-RECORD |
| 1200-READ-MASTER | AM-ACCT-NO, WS-MASTER-KEY, WS-PREV-MASTER-KEY, WS-MASTERS-READ, AM-BALANCE, WS-OPENING-TOTAL | ACCT-MASTER-RECORD, WS-MASTER-KEY, WS-ABORT-MESSAGE, WS-PREV-MASTER-KEY, WS-MASTERS-READ, WS-OPENING-TOTAL |
| 1300-READ-TRAN | TD-ACCT-NO, TD-TRAN-ID, WS-TRAN-ORDER-KEY, WS-PREV-TRAN-ORDER-KEY | TRAN-RECORD, WS-TRAN-KEY, WS-TOK-ACCT-NO, WS-TOK-TRAN-ID, WS-ABORT-MESSAGE, WS-PREV-TRAN-ORDER-KEY |
| 2000-MATCH | WS-TRAN-KEY, WS-MASTER-KEY | RJ-REASON |
| 3000-APPLY-TRAN | AM-BALANCE, AM-STATUS, TD-TRAN-CODE | WS-BALANCE-BEFORE, RJ-REASON |
| 3100-POST-CREDIT | TD-AMOUNT, AM-BALANCE, WS-CREDIT-TOTAL | AM-BALANCE, WS-CREDIT-TOTAL, WS-SIGNED-AMOUNT |
| 3200-POST-DEBIT | AM-BALANCE, TD-AMOUNT, AM-OVERDRAFT-LIMIT, WS-NEW-BALANCE, WS-BALANCE-FLOOR, WS-DEBIT-TOTAL, TD-TRAN-CODE | WS-NEW-BALANCE, WS-BALANCE-FLOOR, RJ-REASON, AM-BALANCE, WS-DEBIT-TOTAL, WS-SIGNED-AMOUNT |
| 3300-POST-FEE | TD-AMOUNT, AM-BALANCE, WS-FEE-TOTAL | AM-BALANCE, WS-FEE-TOTAL, WS-SIGNED-AMOUNT |
| 3500-WRITE-LEG | TRAN-RECORD, TD-TRAN-CODE, TD-ACCT-NO, TD-COUNTER-ACCT-NO, WS-LEG-ACCT-NO, WS-LEG-COUNT, WS-SAVED-TRAN | WS-SAVED-TRAN, TD-TRAN-CODE, WS-LEG-ACCT-NO, TD-ACCT-NO, TD-COUNTER-ACCT-NO, LEGS-OUT-RECORD, WS-LEG-COUNT, TRAN-RECORD |
| 3800-WRITE-AUDIT | WS-BUSINESS-DATE, TD-TRAN-ID, TD-ACCT-NO, TD-TRAN-CODE, WS-SIGNED-AMOUNT, WS-BALANCE-BEFORE, AM-BALANCE, AUDIT-RECORD, WS-POSTED-COUNT | AU-BUSINESS-DATE, AM-LAST-ACTIVITY-DATE, AU-TRAN-ID, AU-ACCT-NO, AU-TRAN-CODE, AU-AMOUNT, AU-BALANCE-BEFORE, AU-BALANCE-AFTER, AUDIT-OUT-RECORD, WS-POSTED-COUNT |
| 3900-REJECT-TRAN | TRAN-RECORD, REJECT-RECORD, WS-REJECT-COUNT, TD-TRAN-CODE | RJ-PROGRAM, RJ-TRAN, REJECT-OUT-RECORD, WS-REJECT-COUNT |
| 4000-WRITE-MASTER | ACCT-MASTER-RECORD, WS-MASTERS-WRITTEN, AM-BALANCE, WS-CLOSING-TOTAL | ACCT-OUT-RECORD, WS-MASTERS-WRITTEN, WS-CLOSING-TOTAL |
| 5000-RECONCILE | WS-OPENING-TOTAL, WS-CREDIT-TOTAL, WS-EXPECTED-TOTAL, WS-DEBIT-TOTAL, WS-FEE-TOTAL, WS-CLOSING-TOTAL, WS-MASTERS-WRITTEN, WS-MASTERS-READ | WS-EXPECTED-TOTAL, WS-ABORT-MESSAGE |
| 8000-WRITE-REPORT | WS-BUSINESS-DATE, WS-MASTERS-READ, WS-REPORT-COUNT-LINE, WS-MASTERS-WRITTEN, WS-POSTED-COUNT, WS-REJECT-COUNT, WS-LEG-COUNT, WS-OPENING-TOTAL, WS-REPORT-AMOUNT-LINE, WS-CREDIT-TOTAL, WS-DEBIT-TOTAL, WS-FEE-TOTAL, WS-CLOSING-TOTAL | REPORT-LINE, WS-RC-LABEL, WS-RC-VALUE, WS-RA-LABEL, WS-RA-VALUE |
| 9000-CLOSE-FILES | | |
| 9900-ABORT | WS-ABORT-MESSAGE | RETURN-CODE |

## Lineage

Where each data item's values come from, and which records go to which files.

| Into | From | Statement | Paragraph | Source |
| --- | --- | --- | --- | --- |
| ACCT-MASTER-RECORD | file ACCT-IN | READ | 1200-READ-MASTER | [TXPOST.cbl:203](../../../pangaea/programs/TXPOST.cbl#L203) |
| file ACCT-OUT | ACCT-OUT-RECORD | WRITE | 4000-WRITE-MASTER | [TXPOST.cbl:352](../../../pangaea/programs/TXPOST.cbl#L352) |
| ACCT-OUT-RECORD | ACCT-MASTER-RECORD | WRITE | 4000-WRITE-MASTER | [TXPOST.cbl:352](../../../pangaea/programs/TXPOST.cbl#L352) |
| AM-BALANCE | TD-AMOUNT | ADD | 3100-POST-CREDIT | [TXPOST.cbl:278](../../../pangaea/programs/TXPOST.cbl#L278) |
| AM-BALANCE | WS-NEW-BALANCE | MOVE | 3200-POST-DEBIT | [TXPOST.cbl:292](../../../pangaea/programs/TXPOST.cbl#L292) |
| AM-BALANCE | TD-AMOUNT | SUBTRACT | 3300-POST-FEE | [TXPOST.cbl:304](../../../pangaea/programs/TXPOST.cbl#L304) |
| AM-LAST-ACTIVITY-DATE | WS-BUSINESS-DATE | MOVE | 3800-WRITE-AUDIT | [TXPOST.cbl:328](../../../pangaea/programs/TXPOST.cbl#L328) |
| AU-ACCT-NO | TD-ACCT-NO | MOVE | 3800-WRITE-AUDIT | [TXPOST.cbl:331](../../../pangaea/programs/TXPOST.cbl#L331) |
| AU-AMOUNT | WS-SIGNED-AMOUNT | MOVE | 3800-WRITE-AUDIT | [TXPOST.cbl:333](../../../pangaea/programs/TXPOST.cbl#L333) |
| AU-BALANCE-AFTER | AM-BALANCE | MOVE | 3800-WRITE-AUDIT | [TXPOST.cbl:335](../../../pangaea/programs/TXPOST.cbl#L335) |
| AU-BALANCE-BEFORE | WS-BALANCE-BEFORE | MOVE | 3800-WRITE-AUDIT | [TXPOST.cbl:334](../../../pangaea/programs/TXPOST.cbl#L334) |
| AU-BUSINESS-DATE | WS-BUSINESS-DATE | MOVE | 3800-WRITE-AUDIT | [TXPOST.cbl:328](../../../pangaea/programs/TXPOST.cbl#L328) |
| AU-TRAN-CODE | TD-TRAN-CODE | MOVE | 3800-WRITE-AUDIT | [TXPOST.cbl:332](../../../pangaea/programs/TXPOST.cbl#L332) |
| AU-TRAN-ID | TD-TRAN-ID | MOVE | 3800-WRITE-AUDIT | [TXPOST.cbl:330](../../../pangaea/programs/TXPOST.cbl#L330) |
| file AUDIT-OUT | AUDIT-OUT-RECORD | WRITE | 3800-WRITE-AUDIT | [TXPOST.cbl:336](../../../pangaea/programs/TXPOST.cbl#L336) |
| AUDIT-OUT-RECORD | AUDIT-RECORD | WRITE | 3800-WRITE-AUDIT | [TXPOST.cbl:336](../../../pangaea/programs/TXPOST.cbl#L336) |
| BANK-CONTROL-RECORD | file CONTROL-FILE | READ | 1000-INITIALIZE | [TXPOST.cbl:191](../../../pangaea/programs/TXPOST.cbl#L191) |
| file LEGS-OUT | LEGS-OUT-RECORD | WRITE | 3500-WRITE-LEG | [TXPOST.cbl:322](../../../pangaea/programs/TXPOST.cbl#L322) |
| LEGS-OUT-RECORD | TRAN-RECORD | WRITE | 3500-WRITE-LEG | [TXPOST.cbl:322](../../../pangaea/programs/TXPOST.cbl#L322) |
| file REJECT-OUT | REJECT-OUT-RECORD | WRITE | 3900-REJECT-TRAN | [TXPOST.cbl:344](../../../pangaea/programs/TXPOST.cbl#L344) |
| REJECT-OUT-RECORD | REJECT-RECORD | WRITE | 3900-REJECT-TRAN | [TXPOST.cbl:344](../../../pangaea/programs/TXPOST.cbl#L344) |
| REPORT-LINE | WS-BUSINESS-DATE | STRING | 8000-WRITE-REPORT | [TXPOST.cbl:371](../../../pangaea/programs/TXPOST.cbl#L371) |
| REPORT-LINE | WS-REPORT-COUNT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:379](../../../pangaea/programs/TXPOST.cbl#L379) |
| REPORT-LINE | WS-REPORT-COUNT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:382](../../../pangaea/programs/TXPOST.cbl#L382) |
| REPORT-LINE | WS-REPORT-COUNT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:385](../../../pangaea/programs/TXPOST.cbl#L385) |
| REPORT-LINE | WS-REPORT-COUNT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:388](../../../pangaea/programs/TXPOST.cbl#L388) |
| REPORT-LINE | WS-REPORT-COUNT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:391](../../../pangaea/programs/TXPOST.cbl#L391) |
| REPORT-LINE | WS-REPORT-AMOUNT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:396](../../../pangaea/programs/TXPOST.cbl#L396) |
| REPORT-LINE | WS-REPORT-AMOUNT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:399](../../../pangaea/programs/TXPOST.cbl#L399) |
| REPORT-LINE | WS-REPORT-AMOUNT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:402](../../../pangaea/programs/TXPOST.cbl#L402) |
| REPORT-LINE | WS-REPORT-AMOUNT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:405](../../../pangaea/programs/TXPOST.cbl#L405) |
| REPORT-LINE | WS-REPORT-AMOUNT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:408](../../../pangaea/programs/TXPOST.cbl#L408) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:374](../../../pangaea/programs/TXPOST.cbl#L374) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:376](../../../pangaea/programs/TXPOST.cbl#L376) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:379](../../../pangaea/programs/TXPOST.cbl#L379) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:382](../../../pangaea/programs/TXPOST.cbl#L382) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:385](../../../pangaea/programs/TXPOST.cbl#L385) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:388](../../../pangaea/programs/TXPOST.cbl#L388) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:391](../../../pangaea/programs/TXPOST.cbl#L391) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:393](../../../pangaea/programs/TXPOST.cbl#L393) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:396](../../../pangaea/programs/TXPOST.cbl#L396) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:399](../../../pangaea/programs/TXPOST.cbl#L399) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:402](../../../pangaea/programs/TXPOST.cbl#L402) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:405](../../../pangaea/programs/TXPOST.cbl#L405) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:408](../../../pangaea/programs/TXPOST.cbl#L408) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:410](../../../pangaea/programs/TXPOST.cbl#L410) |
| file REPORT-OUT | REPORT-LINE | WRITE | 8000-WRITE-REPORT | [TXPOST.cbl:412](../../../pangaea/programs/TXPOST.cbl#L412) |
| RJ-TRAN | TRAN-RECORD | MOVE | 3900-REJECT-TRAN | [TXPOST.cbl:343](../../../pangaea/programs/TXPOST.cbl#L343) |
| TD-ACCT-NO | TD-COUNTER-ACCT-NO | MOVE | 3500-WRITE-LEG | [TXPOST.cbl:320](../../../pangaea/programs/TXPOST.cbl#L320) |
| TD-COUNTER-ACCT-NO | WS-LEG-ACCT-NO | MOVE | 3500-WRITE-LEG | [TXPOST.cbl:321](../../../pangaea/programs/TXPOST.cbl#L321) |
| TRAN-RECORD | file TRAN-IN | READ | 1300-READ-TRAN | [TXPOST.cbl:220](../../../pangaea/programs/TXPOST.cbl#L220) |
| TRAN-RECORD | WS-SAVED-TRAN | MOVE | 3500-WRITE-LEG | [TXPOST.cbl:324](../../../pangaea/programs/TXPOST.cbl#L324) |
| WS-BALANCE-BEFORE | AM-BALANCE | MOVE | 3000-APPLY-TRAN | [TXPOST.cbl:253](../../../pangaea/programs/TXPOST.cbl#L253) |
| WS-BALANCE-FLOOR | AM-OVERDRAFT-LIMIT | COMPUTE | 3200-POST-DEBIT | [TXPOST.cbl:287](../../../pangaea/programs/TXPOST.cbl#L287) |
| WS-BUSINESS-DATE | BC-BUSINESS-DATE | MOVE | 1000-INITIALIZE | [TXPOST.cbl:196](../../../pangaea/programs/TXPOST.cbl#L196) |
| WS-CLOSING-TOTAL | AM-BALANCE | ADD | 4000-WRITE-MASTER | [TXPOST.cbl:354](../../../pangaea/programs/TXPOST.cbl#L354) |
| WS-CREDIT-TOTAL | TD-AMOUNT | ADD | 3100-POST-CREDIT | [TXPOST.cbl:279](../../../pangaea/programs/TXPOST.cbl#L279) |
| WS-DEBIT-TOTAL | TD-AMOUNT | ADD | 3200-POST-DEBIT | [TXPOST.cbl:293](../../../pangaea/programs/TXPOST.cbl#L293) |
| WS-EXPECTED-TOTAL | WS-OPENING-TOTAL | MOVE | 5000-RECONCILE | [TXPOST.cbl:358](../../../pangaea/programs/TXPOST.cbl#L358) |
| WS-EXPECTED-TOTAL | WS-CREDIT-TOTAL | ADD | 5000-RECONCILE | [TXPOST.cbl:359](../../../pangaea/programs/TXPOST.cbl#L359) |
| WS-EXPECTED-TOTAL | WS-DEBIT-TOTAL | SUBTRACT | 5000-RECONCILE | [TXPOST.cbl:360](../../../pangaea/programs/TXPOST.cbl#L360) |
| WS-EXPECTED-TOTAL | WS-FEE-TOTAL | SUBTRACT | 5000-RECONCILE | [TXPOST.cbl:361](../../../pangaea/programs/TXPOST.cbl#L361) |
| WS-FEE-TOTAL | TD-AMOUNT | ADD | 3300-POST-FEE | [TXPOST.cbl:305](../../../pangaea/programs/TXPOST.cbl#L305) |
| WS-LEG-ACCT-NO | TD-ACCT-NO | MOVE | 3500-WRITE-LEG | [TXPOST.cbl:319](../../../pangaea/programs/TXPOST.cbl#L319) |
| WS-MASTER-KEY | AM-ACCT-NO | MOVE | 1200-READ-MASTER | [TXPOST.cbl:207](../../../pangaea/programs/TXPOST.cbl#L207) |
| WS-NEW-BALANCE | AM-BALANCE | COMPUTE | 3200-POST-DEBIT | [TXPOST.cbl:286](../../../pangaea/programs/TXPOST.cbl#L286) |
| WS-NEW-BALANCE | TD-AMOUNT | COMPUTE | 3200-POST-DEBIT | [TXPOST.cbl:286](../../../pangaea/programs/TXPOST.cbl#L286) |
| WS-OPENING-TOTAL | AM-BALANCE | ADD | 1200-READ-MASTER | [TXPOST.cbl:215](../../../pangaea/programs/TXPOST.cbl#L215) |
| WS-PREV-MASTER-KEY | WS-MASTER-KEY | MOVE | 1200-READ-MASTER | [TXPOST.cbl:213](../../../pangaea/programs/TXPOST.cbl#L213) |
| WS-PREV-TRAN-ORDER-KEY | WS-TRAN-ORDER-KEY | MOVE | 1300-READ-TRAN | [TXPOST.cbl:231](../../../pangaea/programs/TXPOST.cbl#L231) |
| WS-RA-VALUE | WS-OPENING-TOTAL | MOVE | 8000-WRITE-REPORT | [TXPOST.cbl:395](../../../pangaea/programs/TXPOST.cbl#L395) |
| WS-RA-VALUE | WS-CREDIT-TOTAL | MOVE | 8000-WRITE-REPORT | [TXPOST.cbl:398](../../../pangaea/programs/TXPOST.cbl#L398) |
| WS-RA-VALUE | WS-DEBIT-TOTAL | MOVE | 8000-WRITE-REPORT | [TXPOST.cbl:401](../../../pangaea/programs/TXPOST.cbl#L401) |
| WS-RA-VALUE | WS-FEE-TOTAL | MOVE | 8000-WRITE-REPORT | [TXPOST.cbl:404](../../../pangaea/programs/TXPOST.cbl#L404) |
| WS-RA-VALUE | WS-CLOSING-TOTAL | MOVE | 8000-WRITE-REPORT | [TXPOST.cbl:407](../../../pangaea/programs/TXPOST.cbl#L407) |
| WS-RC-VALUE | WS-MASTERS-READ | MOVE | 8000-WRITE-REPORT | [TXPOST.cbl:378](../../../pangaea/programs/TXPOST.cbl#L378) |
| WS-RC-VALUE | WS-MASTERS-WRITTEN | MOVE | 8000-WRITE-REPORT | [TXPOST.cbl:381](../../../pangaea/programs/TXPOST.cbl#L381) |
| WS-RC-VALUE | WS-POSTED-COUNT | MOVE | 8000-WRITE-REPORT | [TXPOST.cbl:384](../../../pangaea/programs/TXPOST.cbl#L384) |
| WS-RC-VALUE | WS-REJECT-COUNT | MOVE | 8000-WRITE-REPORT | [TXPOST.cbl:387](../../../pangaea/programs/TXPOST.cbl#L387) |
| WS-RC-VALUE | WS-LEG-COUNT | MOVE | 8000-WRITE-REPORT | [TXPOST.cbl:390](../../../pangaea/programs/TXPOST.cbl#L390) |
| WS-SAVED-TRAN | TRAN-RECORD | MOVE | 3500-WRITE-LEG | [TXPOST.cbl:313](../../../pangaea/programs/TXPOST.cbl#L313) |
| WS-SIGNED-AMOUNT | TD-AMOUNT | MOVE | 3100-POST-CREDIT | [TXPOST.cbl:280](../../../pangaea/programs/TXPOST.cbl#L280) |
| WS-SIGNED-AMOUNT | TD-AMOUNT | COMPUTE | 3200-POST-DEBIT | [TXPOST.cbl:294](../../../pangaea/programs/TXPOST.cbl#L294) |
| WS-SIGNED-AMOUNT | TD-AMOUNT | COMPUTE | 3300-POST-FEE | [TXPOST.cbl:306](../../../pangaea/programs/TXPOST.cbl#L306) |
| WS-TOK-ACCT-NO | TD-ACCT-NO | MOVE | 1300-READ-TRAN | [TXPOST.cbl:224](../../../pangaea/programs/TXPOST.cbl#L224) |
| WS-TOK-TRAN-ID | TD-TRAN-ID | MOVE | 1300-READ-TRAN | [TXPOST.cbl:225](../../../pangaea/programs/TXPOST.cbl#L225) |
| WS-TRAN-KEY | TD-ACCT-NO | MOVE | 1300-READ-TRAN | [TXPOST.cbl:224](../../../pangaea/programs/TXPOST.cbl#L224) |

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
| 01 | ACCT-OUT-RECORD | 1 | 100 | X(100) | | | [TXPOST.cbl:84](../../../pangaea/programs/TXPOST.cbl#L84) |

### AUDIT-OUT-RECORD

File AUDIT-OUT. 100 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | AUDIT-OUT-RECORD | 1 | 100 | X(100) | | | [TXPOST.cbl:86](../../../pangaea/programs/TXPOST.cbl#L86) |

### REJECT-OUT-RECORD

File REJECT-OUT. 100 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | REJECT-OUT-RECORD | 1 | 100 | X(100) | | | [TXPOST.cbl:88](../../../pangaea/programs/TXPOST.cbl#L88) |

### LEGS-OUT-RECORD

File LEGS-OUT. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | LEGS-OUT-RECORD | 1 | 80 | X(80) | | | [TXPOST.cbl:90](../../../pangaea/programs/TXPOST.cbl#L90) |

### REPORT-LINE

File REPORT-OUT. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | REPORT-LINE | 1 | 80 | X(80) | | | [TXPOST.cbl:92](../../../pangaea/programs/TXPOST.cbl#L92) |

### WS-FILE-STATUSES

Working storage. 16 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-FILE-STATUSES | 1 | 16 | | | | [TXPOST.cbl:95](../../../pangaea/programs/TXPOST.cbl#L95) |
| 05 | WS-CONTROL-STATUS | 1 | 2 | XX | | | [TXPOST.cbl:96](../../../pangaea/programs/TXPOST.cbl#L96) |
| 05 | WS-ACCT-IN-STATUS | 3 | 2 | XX | | | [TXPOST.cbl:97](../../../pangaea/programs/TXPOST.cbl#L97) |
| 05 | WS-TRAN-IN-STATUS | 5 | 2 | XX | | | [TXPOST.cbl:98](../../../pangaea/programs/TXPOST.cbl#L98) |
| 05 | WS-ACCT-OUT-STATUS | 7 | 2 | XX | | | [TXPOST.cbl:99](../../../pangaea/programs/TXPOST.cbl#L99) |
| 05 | WS-AUDIT-STATUS | 9 | 2 | XX | | | [TXPOST.cbl:100](../../../pangaea/programs/TXPOST.cbl#L100) |
| 88 | AUDIT-OPENED | | | | | VALUE '00', '05' | [TXPOST.cbl:101](../../../pangaea/programs/TXPOST.cbl#L101) |
| 05 | WS-REJECT-STATUS | 11 | 2 | XX | | | [TXPOST.cbl:102](../../../pangaea/programs/TXPOST.cbl#L102) |
| 88 | REJECTS-OPENED | | | | | VALUE '00', '05' | [TXPOST.cbl:103](../../../pangaea/programs/TXPOST.cbl#L103) |
| 05 | WS-LEGS-STATUS | 13 | 2 | XX | | | [TXPOST.cbl:104](../../../pangaea/programs/TXPOST.cbl#L104) |
| 05 | WS-REPORT-STATUS | 15 | 2 | XX | | | [TXPOST.cbl:105](../../../pangaea/programs/TXPOST.cbl#L105) |

### WS-BUSINESS-DATE

Working storage. 8 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-BUSINESS-DATE | 1 | 8 | 9(8) | | | [TXPOST.cbl:107](../../../pangaea/programs/TXPOST.cbl#L107) |

### WS-ABORT-MESSAGE

Working storage. 40 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-ABORT-MESSAGE | 1 | 40 | X(40) | | | [TXPOST.cbl:108](../../../pangaea/programs/TXPOST.cbl#L108) |

### WS-MASTER-KEY

Working storage. 10 bytes. MATCH KEYS: HIGH-VALUES ONCE A FILE IS AT ITS END.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-MASTER-KEY | 1 | 10 | X(10) | | | [TXPOST.cbl:111](../../../pangaea/programs/TXPOST.cbl#L111) |

### WS-TRAN-KEY

Working storage. 10 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TRAN-KEY | 1 | 10 | X(10) | | | [TXPOST.cbl:112](../../../pangaea/programs/TXPOST.cbl#L112) |

### WS-PREV-MASTER-KEY

Working storage. 10 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-PREV-MASTER-KEY | 1 | 10 | X(10) | | VALUE LOW-VALUES | [TXPOST.cbl:113](../../../pangaea/programs/TXPOST.cbl#L113) |

### WS-TRAN-ORDER-KEY

Working storage. 20 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TRAN-ORDER-KEY | 1 | 20 | | | | [TXPOST.cbl:114](../../../pangaea/programs/TXPOST.cbl#L114) |
| 05 | WS-TOK-ACCT-NO | 1 | 10 | X(10) | | | [TXPOST.cbl:115](../../../pangaea/programs/TXPOST.cbl#L115) |
| 05 | WS-TOK-TRAN-ID | 11 | 10 | X(10) | | | [TXPOST.cbl:116](../../../pangaea/programs/TXPOST.cbl#L116) |

### WS-PREV-TRAN-ORDER-KEY

Working storage. 20 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-PREV-TRAN-ORDER-KEY | 1 | 20 | X(20) | | VALUE LOW-VALUES | [TXPOST.cbl:117](../../../pangaea/programs/TXPOST.cbl#L117) |

### WS-POSTING

Working storage. 50 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-POSTING | 1 | 50 | | | | [TXPOST.cbl:119](../../../pangaea/programs/TXPOST.cbl#L119) |
| 05 | WS-BALANCE-BEFORE | 1 | 13 | S9(11)V99 | | | [TXPOST.cbl:120](../../../pangaea/programs/TXPOST.cbl#L120) |
| 05 | WS-NEW-BALANCE | 14 | 13 | S9(11)V99 | | | [TXPOST.cbl:121](../../../pangaea/programs/TXPOST.cbl#L121) |
| 05 | WS-BALANCE-FLOOR | 27 | 13 | S9(11)V99 | | | [TXPOST.cbl:122](../../../pangaea/programs/TXPOST.cbl#L122) |
| 05 | WS-SIGNED-AMOUNT | 40 | 11 | S9(9)V99 | | | [TXPOST.cbl:123](../../../pangaea/programs/TXPOST.cbl#L123) |

### WS-SAVED-TRAN

Working storage. 80 bytes. A LEG IS WRITTEN FROM THE TRANSACTION RECORD ITSELF, TURNED AROUND; THE TRANSACTION IS SAVED HERE AND PUT BACK AFTER.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-SAVED-TRAN | 1 | 80 | X(80) | | | [TXPOST.cbl:127](../../../pangaea/programs/TXPOST.cbl#L127) |

### WS-LEG-ACCT-NO

Working storage. 10 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-LEG-ACCT-NO | 1 | 10 | 9(10) | | | [TXPOST.cbl:128](../../../pangaea/programs/TXPOST.cbl#L128) |

### WS-TOTALS

Working storage. 130 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-TOTALS | 1 | 130 | | | | [TXPOST.cbl:130](../../../pangaea/programs/TXPOST.cbl#L130) |
| 05 | WS-MASTERS-READ | 1 | 8 | 9(8) | | VALUE 0 | [TXPOST.cbl:131](../../../pangaea/programs/TXPOST.cbl#L131) |
| 05 | WS-MASTERS-WRITTEN | 9 | 8 | 9(8) | | VALUE 0 | [TXPOST.cbl:132](../../../pangaea/programs/TXPOST.cbl#L132) |
| 05 | WS-POSTED-COUNT | 17 | 8 | 9(8) | | VALUE 0 | [TXPOST.cbl:133](../../../pangaea/programs/TXPOST.cbl#L133) |
| 05 | WS-REJECT-COUNT | 25 | 8 | 9(8) | | VALUE 0 | [TXPOST.cbl:134](../../../pangaea/programs/TXPOST.cbl#L134) |
| 05 | WS-LEG-COUNT | 33 | 8 | 9(8) | | VALUE 0 | [TXPOST.cbl:135](../../../pangaea/programs/TXPOST.cbl#L135) |
| 05 | WS-OPENING-TOTAL | 41 | 15 | S9(13)V99 | | VALUE 0 | [TXPOST.cbl:136](../../../pangaea/programs/TXPOST.cbl#L136) |
| 05 | WS-CREDIT-TOTAL | 56 | 15 | 9(13)V99 | | VALUE 0 | [TXPOST.cbl:137](../../../pangaea/programs/TXPOST.cbl#L137) |
| 05 | WS-DEBIT-TOTAL | 71 | 15 | 9(13)V99 | | VALUE 0 | [TXPOST.cbl:138](../../../pangaea/programs/TXPOST.cbl#L138) |
| 05 | WS-FEE-TOTAL | 86 | 15 | 9(13)V99 | | VALUE 0 | [TXPOST.cbl:139](../../../pangaea/programs/TXPOST.cbl#L139) |
| 05 | WS-CLOSING-TOTAL | 101 | 15 | S9(13)V99 | | VALUE 0 | [TXPOST.cbl:140](../../../pangaea/programs/TXPOST.cbl#L140) |
| 05 | WS-EXPECTED-TOTAL | 116 | 15 | S9(13)V99 | | VALUE 0 | [TXPOST.cbl:141](../../../pangaea/programs/TXPOST.cbl#L141) |

### AUDIT-RECORD

Working storage. 100 bytes, layout in copybook [AUDITREC](../copybooks/AUDITREC.md).

### REJECT-RECORD

Working storage. 100 bytes, layout in copybook [REJREC](../copybooks/REJREC.md).

### WS-REPORT-COUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-COUNT-LINE | 1 | 80 | | | | [TXPOST.cbl:147](../../../pangaea/programs/TXPOST.cbl#L147) |
| 05 | WS-RC-LABEL | 1 | 30 | X(30) | | | [TXPOST.cbl:148](../../../pangaea/programs/TXPOST.cbl#L148) |
| 05 | WS-RC-VALUE | 31 | 10 | ZZ,ZZZ,ZZ9 | | | [TXPOST.cbl:149](../../../pangaea/programs/TXPOST.cbl#L149) |
| 05 | FILLER | 41 | 40 | X(40) | | VALUE SPACES | [TXPOST.cbl:150](../../../pangaea/programs/TXPOST.cbl#L150) |

### WS-REPORT-AMOUNT-LINE

Working storage. 80 bytes.

| Level | Field | Start | Length | Picture | Usage | Notes | Source |
| --: | --- | --: | --: | --- | --- | --- | --- |
| 01 | WS-REPORT-AMOUNT-LINE | 1 | 80 | | | | [TXPOST.cbl:151](../../../pangaea/programs/TXPOST.cbl#L151) |
| 05 | WS-RA-LABEL | 1 | 30 | X(30) | | | [TXPOST.cbl:152](../../../pangaea/programs/TXPOST.cbl#L152) |
| 05 | WS-RA-VALUE | 31 | 21 | Z,ZZZ,ZZZ,ZZZ,ZZ9.99- | | | [TXPOST.cbl:153](../../../pangaea/programs/TXPOST.cbl#L153) |
| 05 | FILLER | 52 | 29 | X(29) | | VALUE SPACES | [TXPOST.cbl:154](../../../pangaea/programs/TXPOST.cbl#L154) |
