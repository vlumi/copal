# DISPUTE: decision tables

The business rules of [DISPUTE](../programs/DISPUTE.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass. Under each table, each rule as a sentence, written from the conditions it tests and the statements it runs, with no more in it than the table has: a condition name with the field it is declared on and its values.

## Contents

- [READ-MASTER](#read-master)
- [READ-EVENT](#read-event)
- [APPLY-NEXT](#apply-next)
- [APPLY-EVENTS](#apply-events)
- [APPLY-EVENTS, in its loop](#apply-events-in-its-loop)
- [APPLY-EVENT](#apply-event)

## READ-MASTER

| | 1 | 2 |
| --- | :-: | :-: |
| `READ MASTER-IN AT END` | Y | N |
| → `READ MASTER-IN` | 1 | 1 |
| → `MOVE ALL '9' TO WS-MASTER-KEY` | 2 | |
| → `MOVE MASTER-IN-RECORD TO DISPUTE-CASE` | | 2 |
| → `MOVE DC-CASE-NO TO WS-MASTER-KEY` | | 3 |

- **Rule 1.** When reading `MASTER-IN` finds no more records, it reads `MASTER-IN` and sets `WS-MASTER-KEY` to `ALL '9'`.
- **Rule 2.** When reading `MASTER-IN` finds a record, it reads `MASTER-IN`, sets `DISPUTE-CASE` to `MASTER-IN-RECORD`, and sets `WS-MASTER-KEY` to `DC-CASE-NO`.

Source: [dispute.cbl:74](../../../laurasia/programs/dispute.cbl#L74).

## READ-EVENT

| | 1 | 2 |
| --- | :-: | :-: |
| `READ EVENTS-IN AT END` | Y | N |
| → `READ EVENTS-IN` | 1 | 1 |
| → `MOVE ALL '9' TO WS-EVENT-KEY` | 2 | |
| → `MOVE DE-CASE-NO TO WS-EVENT-KEY` | | 2 |

- **Rule 1.** When reading `EVENTS-IN` finds no more records, it reads `EVENTS-IN` and sets `WS-EVENT-KEY` to `ALL '9'`.
- **Rule 2.** When reading `EVENTS-IN` finds a record, it reads `EVENTS-IN` and sets `WS-EVENT-KEY` to `DE-CASE-NO`.

Source: [dispute.cbl:82](../../../laurasia/programs/dispute.cbl#L82).

## APPLY-NEXT

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `WS-MASTER-KEY < WS-EVENT-KEY` | Y | N | N |
| `WS-MASTER-KEY = WS-EVENT-KEY` | | Y | N |
| → `MOVE MASTER-IN-RECORD TO DISPUTE-CASE` | 1 | | |
| → `WRITE DISPUTE-CASE` | 2 | | |
| → `PERFORM READ-MASTER` | 3 | 3 | |
| → `MOVE MASTER-IN-RECORD TO WS-HELD-CASE` | | 1 | |
| → `SET HOLDING-CASE TO TRUE` | | 2 | |
| → `PERFORM APPLY-EVENTS` | | 4 | 2 |
| → `SET HOLDING-NOTHING TO TRUE` | | | 1 |

- **Rule 1.** When `WS-MASTER-KEY` is less than `WS-EVENT-KEY`, it sets `DISPUTE-CASE` to `MASTER-IN-RECORD`, writes `DISPUTE-CASE`, and performs `READ-MASTER`.
- **Rule 2.** When `WS-MASTER-KEY` is at least `WS-EVENT-KEY` and `WS-MASTER-KEY` is `WS-EVENT-KEY`, it sets `WS-HELD-CASE` to `MASTER-IN-RECORD`, sets `WS-HOLDING` to `HOLDING-CASE`, performs `READ-MASTER`, and performs `APPLY-EVENTS`.
- **Rule 3.** When `WS-MASTER-KEY` is at least `WS-EVENT-KEY` and `WS-MASTER-KEY` is not `WS-EVENT-KEY`, it sets `WS-HOLDING` to `HOLDING-NOTHING` and performs `APPLY-EVENTS`.

Source: [dispute.cbl:91](../../../laurasia/programs/dispute.cbl#L91).

## APPLY-EVENTS

| | 1 | 2 |
| --- | :-: | :-: |
| `HOLDING-CASE` | Y | N |
| → `MOVE WS-EVENT-KEY TO DC-CASE-NO` | 1 | 1 |
| → `PERFORM UNTIL WS-EVENT-KEY NOT = DC-CASE-NO` | 2 | 2 |
| → `MOVE WS-HELD-CASE TO DISPUTE-CASE` | 3 | |
| → `WRITE DISPUTE-CASE` | 4 | |

- **Rule 1.** When `WS-HOLDING` is `'Y'` (`HOLDING-CASE`), it sets `DC-CASE-NO` to `WS-EVENT-KEY`, repeats its loop until `WS-EVENT-KEY` is not `DC-CASE-NO`, sets `DISPUTE-CASE` to `WS-HELD-CASE`, and writes `DISPUTE-CASE`.
- **Rule 2.** When `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`), it sets `DC-CASE-NO` to `WS-EVENT-KEY` and repeats its loop until `WS-EVENT-KEY` is not `DC-CASE-NO`.

Source: [dispute.cbl:107](../../../laurasia/programs/dispute.cbl#L107).

## APPLY-EVENTS, in its loop

Each pass of `PERFORM UNTIL WS-EVENT-KEY NOT = DC-CASE-NO`.

| | 1 | 2 | 3 |
| --- | :-: | :-: | :-: |
| `HOLDING-CASE` | Y | N | N |
| `HOLDING-CASE, again` | Y | Y | N |
| → `MOVE WS-HELD-CASE TO DISPUTE-CASE` | 1 | | |
| → `PERFORM APPLY-EVENT` | 2 | 1 | 1 |
| → `MOVE DISPUTE-CASE TO WS-HELD-CASE` | 3 | 2 | |
| → `PERFORM READ-EVENT` | 4 | 3 | 2 |

- **Rule 1.** When `WS-HOLDING` is `'Y'` (`HOLDING-CASE`); and `WS-HOLDING` is `'Y'` (`HOLDING-CASE`), tested again, it sets `DISPUTE-CASE` to `WS-HELD-CASE`, performs `APPLY-EVENT`, sets `WS-HELD-CASE` to `DISPUTE-CASE`, and performs `READ-EVENT`.
- **Rule 2.** When `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`); and `WS-HOLDING` is `'Y'` (`HOLDING-CASE`), tested again, it performs `APPLY-EVENT`, sets `WS-HELD-CASE` to `DISPUTE-CASE`, and performs `READ-EVENT`.
- **Rule 3.** When `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`); and `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`), tested again, it performs `APPLY-EVENT` and performs `READ-EVENT`.

Source: [dispute.cbl:109](../../../laurasia/programs/dispute.cbl#L109).

## APPLY-EVENT

| | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 |
| --- | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: |
| `HOLDING-NOTHING AND DE-OPENED` | Y | N | N | N | N | N | N | N | N |
| `HOLDING-CASE AND DC-OPEN AND DE-CHARGEBACK` | | Y | N | N | N | N | N | N | N |
| `HOLDING-CASE AND DC-CHARGED-BACK AND DE-MERCHANT-ACCEPTS` | | | Y | N | N | N | N | N | N |
| `HOLDING-CASE AND DC-CHARGED-BACK AND DE-REPRESENTMENT` | | | | Y | N | N | N | N | N |
| `HOLDING-CASE AND DC-REPRESENTED AND DE-ISSUER-DROPS` | | | | | Y | N | N | N | N |
| `HOLDING-CASE AND DC-REPRESENTED AND DE-PRE-ARBITRATION` | | | | | | Y | N | N | N |
| `HOLDING-CASE AND DC-PRE-ARBITRATION AND DE-RULED-FOR-ISSUER` | | | | | | | Y | N | N |
| `HOLDING-CASE AND DC-PRE-ARBITRATION AND DE-RULED-FOR-MERCHANT` | | | | | | | | Y | N |
| → `MOVE DE-CASE-NO TO DC-CASE-NO` | 1 | | | | | | | | |
| → `MOVE DE-MERCHANT TO DC-MERCHANT` | 2 | | | | | | | | |
| → `MOVE DE-AMOUNT TO DC-AMOUNT` | 3 | | | | | | | | |
| → `MOVE DE-DATE TO DC-OPENED` | 4 | | | | | | | | |
| → `SET DC-OPEN TO TRUE` | 5 | | | | | | | | |
| → `SET HOLDING-CASE TO TRUE` | 6 | | | | | | | | |
| → `MOVE DE-DATE TO DC-UPDATED` | 7 | 4 | 2 | 4 | 2 | 2 | 4 | 2 | |
| → `SET DC-CHARGED-BACK TO TRUE` | | 1 | | | | | | | |
| → `COMPUTE AD-AMOUNT = 0 - DC-AMOUNT` | | 2 | | | | | 2 | | |
| → `PERFORM WRITE-ADJUSTMENT` | | 3 | | 3 | | | 3 | | |
| → `SET DC-WON TO TRUE` | | | 1 | | | | 1 | | |
| → `SET DC-REPRESENTED TO TRUE` | | | | 1 | | | | | |
| → `MOVE DC-AMOUNT TO AD-AMOUNT` | | | | 2 | | | | | |
| → `SET DC-LOST TO TRUE` | | | | | 1 | | | 1 | |
| → `SET DC-PRE-ARBITRATION TO TRUE` | | | | | | 1 | | | |
| → `MOVE SPACES TO REJECT-LINE` | | | | | | | | | 1 |
| → `MOVE 'STAT' TO RJ-REASON` | | | | | | | | | 2 |
| → `MOVE DISPUTE-EVENT TO RJ-EVENT` | | | | | | | | | 3 |
| → `WRITE REJECT-LINE` | | | | | | | | | 4 |
| → `MOVE 4 TO WS-RESULT` | | | | | | | | | 5 |
| → `EXIT SECTION` | | | | | | | | | 6 |

- **Rule 1.** When `WS-HOLDING` is `'N'` (`HOLDING-NOTHING`) and `DE-EVENT` is `'OP'` (`DE-OPENED`), it sets `DC-CASE-NO` to `DE-CASE-NO`, sets `DC-MERCHANT` to `DE-MERCHANT`, sets `DC-AMOUNT` to `DE-AMOUNT`, sets `DC-OPENED` to `DE-DATE`, sets `DC-STATUS` to `DC-OPEN`, sets `WS-HOLDING` to `HOLDING-CASE`, and sets `DC-UPDATED` to `DE-DATE`.
- **Rule 2.** When `WS-HOLDING` is not `'N'` (not `HOLDING-NOTHING`) or `DE-EVENT` is not `'OP'` (not `DE-OPENED`); and `WS-HOLDING` is `'Y'` (`HOLDING-CASE`) and `DC-STATUS` is `'OP'` (`DC-OPEN`) and `DE-EVENT` is `'CB'` (`DE-CHARGEBACK`), it sets `DC-STATUS` to `DC-CHARGED-BACK`, computes `AD-AMOUNT` as `0 - DC-AMOUNT`, performs `WRITE-ADJUSTMENT`, and sets `DC-UPDATED` to `DE-DATE`.
- **Rule 3.** When `WS-HOLDING` is not `'N'` (not `HOLDING-NOTHING`) or `DE-EVENT` is not `'OP'` (not `DE-OPENED`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'OP'` (not `DC-OPEN`) or `DE-EVENT` is not `'CB'` (not `DE-CHARGEBACK`); and `WS-HOLDING` is `'Y'` (`HOLDING-CASE`) and `DC-STATUS` is `'CB'` (`DC-CHARGED-BACK`) and `DE-EVENT` is `'AC'` (`DE-MERCHANT-ACCEPTS`), it sets `DC-STATUS` to `DC-WON` and sets `DC-UPDATED` to `DE-DATE`.
- **Rule 4.** When `WS-HOLDING` is not `'N'` (not `HOLDING-NOTHING`) or `DE-EVENT` is not `'OP'` (not `DE-OPENED`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'OP'` (not `DC-OPEN`) or `DE-EVENT` is not `'CB'` (not `DE-CHARGEBACK`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'CB'` (not `DC-CHARGED-BACK`) or `DE-EVENT` is not `'AC'` (not `DE-MERCHANT-ACCEPTS`); and `WS-HOLDING` is `'Y'` (`HOLDING-CASE`) and `DC-STATUS` is `'CB'` (`DC-CHARGED-BACK`) and `DE-EVENT` is `'RP'` (`DE-REPRESENTMENT`), it sets `DC-STATUS` to `DC-REPRESENTED`, sets `AD-AMOUNT` to `DC-AMOUNT`, performs `WRITE-ADJUSTMENT`, and sets `DC-UPDATED` to `DE-DATE`.
- **Rule 5.** When `WS-HOLDING` is not `'N'` (not `HOLDING-NOTHING`) or `DE-EVENT` is not `'OP'` (not `DE-OPENED`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'OP'` (not `DC-OPEN`) or `DE-EVENT` is not `'CB'` (not `DE-CHARGEBACK`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'CB'` (not `DC-CHARGED-BACK`) or `DE-EVENT` is not `'AC'` (not `DE-MERCHANT-ACCEPTS`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'CB'` (not `DC-CHARGED-BACK`) or `DE-EVENT` is not `'RP'` (not `DE-REPRESENTMENT`); and `WS-HOLDING` is `'Y'` (`HOLDING-CASE`) and `DC-STATUS` is `'RP'` (`DC-REPRESENTED`) and `DE-EVENT` is `'DR'` (`DE-ISSUER-DROPS`), it sets `DC-STATUS` to `DC-LOST` and sets `DC-UPDATED` to `DE-DATE`.
- **Rule 6.** When `WS-HOLDING` is not `'N'` (not `HOLDING-NOTHING`) or `DE-EVENT` is not `'OP'` (not `DE-OPENED`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'OP'` (not `DC-OPEN`) or `DE-EVENT` is not `'CB'` (not `DE-CHARGEBACK`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'CB'` (not `DC-CHARGED-BACK`) or `DE-EVENT` is not `'AC'` (not `DE-MERCHANT-ACCEPTS`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'CB'` (not `DC-CHARGED-BACK`) or `DE-EVENT` is not `'RP'` (not `DE-REPRESENTMENT`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'RP'` (not `DC-REPRESENTED`) or `DE-EVENT` is not `'DR'` (not `DE-ISSUER-DROPS`); and `WS-HOLDING` is `'Y'` (`HOLDING-CASE`) and `DC-STATUS` is `'RP'` (`DC-REPRESENTED`) and `DE-EVENT` is `'PA'` (`DE-PRE-ARBITRATION`), it sets `DC-STATUS` to `DC-PRE-ARBITRATION` and sets `DC-UPDATED` to `DE-DATE`.
- **Rule 7.** When `WS-HOLDING` is not `'N'` (not `HOLDING-NOTHING`) or `DE-EVENT` is not `'OP'` (not `DE-OPENED`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'OP'` (not `DC-OPEN`) or `DE-EVENT` is not `'CB'` (not `DE-CHARGEBACK`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'CB'` (not `DC-CHARGED-BACK`) or `DE-EVENT` is not `'AC'` (not `DE-MERCHANT-ACCEPTS`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'CB'` (not `DC-CHARGED-BACK`) or `DE-EVENT` is not `'RP'` (not `DE-REPRESENTMENT`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'RP'` (not `DC-REPRESENTED`) or `DE-EVENT` is not `'DR'` (not `DE-ISSUER-DROPS`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'RP'` (not `DC-REPRESENTED`) or `DE-EVENT` is not `'PA'` (not `DE-PRE-ARBITRATION`); and `WS-HOLDING` is `'Y'` (`HOLDING-CASE`) and `DC-STATUS` is `'PA'` (`DC-PRE-ARBITRATION`) and `DE-EVENT` is `'RI'` (`DE-RULED-FOR-ISSUER`), it sets `DC-STATUS` to `DC-WON`, computes `AD-AMOUNT` as `0 - DC-AMOUNT`, performs `WRITE-ADJUSTMENT`, and sets `DC-UPDATED` to `DE-DATE`.
- **Rule 8.** When `WS-HOLDING` is not `'N'` (not `HOLDING-NOTHING`) or `DE-EVENT` is not `'OP'` (not `DE-OPENED`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'OP'` (not `DC-OPEN`) or `DE-EVENT` is not `'CB'` (not `DE-CHARGEBACK`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'CB'` (not `DC-CHARGED-BACK`) or `DE-EVENT` is not `'AC'` (not `DE-MERCHANT-ACCEPTS`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'CB'` (not `DC-CHARGED-BACK`) or `DE-EVENT` is not `'RP'` (not `DE-REPRESENTMENT`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'RP'` (not `DC-REPRESENTED`) or `DE-EVENT` is not `'DR'` (not `DE-ISSUER-DROPS`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'RP'` (not `DC-REPRESENTED`) or `DE-EVENT` is not `'PA'` (not `DE-PRE-ARBITRATION`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'PA'` (not `DC-PRE-ARBITRATION`) or `DE-EVENT` is not `'RI'` (not `DE-RULED-FOR-ISSUER`); and `WS-HOLDING` is `'Y'` (`HOLDING-CASE`) and `DC-STATUS` is `'PA'` (`DC-PRE-ARBITRATION`) and `DE-EVENT` is `'RM'` (`DE-RULED-FOR-MERCHANT`), it sets `DC-STATUS` to `DC-LOST` and sets `DC-UPDATED` to `DE-DATE`.
- **Rule 9.** When `WS-HOLDING` is not `'N'` (not `HOLDING-NOTHING`) or `DE-EVENT` is not `'OP'` (not `DE-OPENED`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'OP'` (not `DC-OPEN`) or `DE-EVENT` is not `'CB'` (not `DE-CHARGEBACK`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'CB'` (not `DC-CHARGED-BACK`) or `DE-EVENT` is not `'AC'` (not `DE-MERCHANT-ACCEPTS`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'CB'` (not `DC-CHARGED-BACK`) or `DE-EVENT` is not `'RP'` (not `DE-REPRESENTMENT`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'RP'` (not `DC-REPRESENTED`) or `DE-EVENT` is not `'DR'` (not `DE-ISSUER-DROPS`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'RP'` (not `DC-REPRESENTED`) or `DE-EVENT` is not `'PA'` (not `DE-PRE-ARBITRATION`); `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'PA'` (not `DC-PRE-ARBITRATION`) or `DE-EVENT` is not `'RI'` (not `DE-RULED-FOR-ISSUER`); and `WS-HOLDING` is not `'Y'` (not `HOLDING-CASE`) or `DC-STATUS` is not `'PA'` (not `DC-PRE-ARBITRATION`) or `DE-EVENT` is not `'RM'` (not `DE-RULED-FOR-MERCHANT`), it sets `REJECT-LINE` to spaces, sets `RJ-REASON` to `'STAT'`, sets `RJ-EVENT` to `DISPUTE-EVENT`, writes `REJECT-LINE`, sets `WS-RESULT` to 4, and leaves the section.

Source: [dispute.cbl:125](../../../laurasia/programs/dispute.cbl#L125).
