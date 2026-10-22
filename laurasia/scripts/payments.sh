#!/bin/sh
# payments.sh: take in the card system's feed: the cardholders' payments
# off their balances, and the cards' status changes.
#
# usage: payments.sh DATE, the business date as YYYYMMDD.
# environment: LCS_BIN, the programs; LCS_DATA, the data.
#
# history
# 2022-09-05 ah  lcs-0497 original version.

DATE="$1"
BIN="$LCS_BIN"
DATA="$LCS_DATA"
LOG="$DATA/logs/night.$DATE.log"

# STEP PAYIN: the payments and status changes on the cards' accounts
export DD_FEEDIN="$DATA/inbound/PAYMENTS.$DATE.txt"
export DD_ACCOUNTS="$DATA/master/accounts.idx"
export DD_POSTOUT="$DATA/outbound/CREDITS.$DATE.txt"
export DD_PAYREJ="$DATA/outbound/PAYREJ.$DATE.txt"
export DD_PAYRPT="$DATA/reports/payin.$DATE.txt"
"$BIN/payin" "$DATE" > "$DATA/reports/payin.$DATE.out" 2>&1
rc=$?
echo "PAYIN payin $rc" >> "$LOG"
if [ "$rc" -ge 8 ]; then exit "$rc"; fi
