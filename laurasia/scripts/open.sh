#!/bin/sh
# open.sh: open the card processor: load the card accounts into the indexed
# account file the nights post to.
#
# usage: open.sh
# environment: LCS_BIN, the programs; LCS_DATA, the data.
#
# history
# 2011-04-04 rk  lcs-0102 original version.

BIN="$LCS_BIN"
DATA="$LCS_DATA"
LOG="$DATA/logs/open.log"
mkdir -p "$DATA/master" "$DATA/logs" "$DATA/reports"

# STEP ACCTLD: the accounts, from the card system's extract
export DD_ACCTIN="$DATA/inbound/accounts.txt"
export DD_ACCOUNTS="$DATA/master/accounts.idx"
"$BIN/acctld" > "$DATA/reports/acctld.out" 2>&1
rc=$?
echo "ACCTLD acctld $rc" >> "$LOG"
if [ "$rc" -ge 8 ]; then exit "$rc"; fi
