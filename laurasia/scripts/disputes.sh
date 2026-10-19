#!/bin/sh
# disputes.sh: move the disputes on by the day's events, and leave the
# merchants' adjustments for the night's settlement.
#
# usage: disputes.sh DATE, the business date as YYYYMMDD.
# environment: LCS_BIN, the programs; LCS_DATA, the data.
#
# history
# 2015-03-23 mt  lcs-0302 original version.

DATE="$1"
export LC_ALL=C
BIN="$LCS_BIN"
DATA="$LCS_DATA"
WORK="$DATA/work/$DATE"
LOG="$DATA/logs/night.$DATE.log"
mkdir -p "$WORK"

# STEP DSPSORT: the events in case order, each case's in the order given
sort -s -t "|" -k1.1,1.8 "$DATA/inbound/DISPEVT.$DATE.txt" > "$WORK/dispevt.sorted"
rc=$?
echo "DSPSORT sort $rc" >> "$LOG"
if [ "$rc" -ne 0 ]; then exit 16; fi

# STEP DISPUTE: the cases moved on, a new master, the adjustments
export DD_DISPIN="$DATA/master/disputes.txt"
export DD_DISPEVT="$WORK/dispevt.sorted"
export DD_DISPOUT="$WORK/disputes.new"
export DD_ADJOUT="$WORK/adjustments.txt"
export DD_DISPREJ="$DATA/outbound/DISPREJ.$DATE.txt"
"$BIN/dispute" "$DATE" > "$DATA/reports/dispute.$DATE.out" 2>&1
rc=$?
echo "DISPUTE dispute $rc" >> "$LOG"
if [ "$rc" -ge 8 ]; then exit "$rc"; fi

# STEP DSPNEW: the new master in place of the old
mv "$WORK/disputes.new" "$DATA/master/disputes.txt"
