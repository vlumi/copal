# DISPUTE: decision tables

The business rules of [DISPUTE](../programs/DISPUTE.md), paragraph by paragraph: each column is one path through the paragraph's decisions, with the outcome of each condition it tests, Y or N, and the actions it takes, numbered in the order they run. The WHENs of an EVALUATE are tried in order, and an IF whose THEN jumps away has the rest of its paragraph as its ELSE. The body of an inline PERFORM has a table of its own, for one pass.

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

Source: [dispute.cbl:74](../../../laurasia/programs/dispute.cbl#L74).

## READ-EVENT

| | 1 | 2 |
| --- | :-: | :-: |
| `READ EVENTS-IN AT END` | Y | N |
| → `READ EVENTS-IN` | 1 | 1 |
| → `MOVE ALL '9' TO WS-EVENT-KEY` | 2 | |
| → `MOVE DE-CASE-NO TO WS-EVENT-KEY` | | 2 |

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

Source: [dispute.cbl:91](../../../laurasia/programs/dispute.cbl#L91).

## APPLY-EVENTS

| | 1 | 2 |
| --- | :-: | :-: |
| `HOLDING-CASE` | Y | N |
| → `MOVE WS-EVENT-KEY TO DC-CASE-NO` | 1 | 1 |
| → `PERFORM UNTIL WS-EVENT-KEY NOT = DC-CASE-NO` | 2 | 2 |
| → `MOVE WS-HELD-CASE TO DISPUTE-CASE` | 3 | |
| → `WRITE DISPUTE-CASE` | 4 | |

Source: [dispute.cbl:107](../../../laurasia/programs/dispute.cbl#L107).

## APPLY-EVENTS, in its loop

Each pass of `PERFORM UNTIL WS-EVENT-KEY NOT = DC-CASE-NO`.

| | 1 | 2 | 3 | 4 |
| --- | :-: | :-: | :-: | :-: |
| `HOLDING-CASE` | Y | Y | N | N |
| `HOLDING-CASE, again` | Y | N | Y | N |
| → `MOVE WS-HELD-CASE TO DISPUTE-CASE` | 1 | 1 | | |
| → `PERFORM APPLY-EVENT` | 2 | 2 | 1 | 1 |
| → `MOVE DISPUTE-CASE TO WS-HELD-CASE` | 3 | | 2 | |
| → `PERFORM READ-EVENT` | 4 | 3 | 3 | 2 |

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

Source: [dispute.cbl:125](../../../laurasia/programs/dispute.cbl#L125).
