#!/bin/sh
# cycle.sh: close the billing cycles of the cards whose cycle day it is.
#
# usage: cycle.sh DATE, the business date as YYYYMMDD.
# environment: LCS_BIN, the programs; LCS_DATA, the data.
#
# history
# 2013-06-03 rk  lcs-0219 original version.

DATE="$1"
BIN="$LCS_BIN"
DATA="$LCS_DATA"
LOG="$DATA/logs/cycle.$DATE.log"
mkdir -p "$DATA/outbound" "$DATA/reports" "$DATA/logs"

# STEP CYCLE: interest, minimum payments and statements
export DD_ACCOUNTS="$DATA/master/accounts.idx"
export DD_POSTOUT="$DATA/outbound/INTEREST.$DATE.txt"
export DD_STMTOUT="$DATA/outbound/STATEMENTS.$DATE.txt"
"$BIN/cycle" "$DATE" > "$DATA/reports/cycle.$DATE.out" 2>&1
rc=$?
echo "CYCLE cycle $rc" >> "$LOG"
if [ "$rc" -ge 8 ]; then exit "$rc"; fi

# STEP ACCTUL: the accounts as they stand after the cycle
export DD_ACCOUNTS="$DATA/master/accounts.idx"
export DD_ACCTOUT="$DATA/outbound/ACCOUNTS.CYCLE.$DATE.txt"
"$BIN/acctul" > "$DATA/reports/acctul.cycle.$DATE.out" 2>&1
rc=$?
echo "ACCTUL acctul $rc" >> "$LOG"
