#!/bin/sh
# nightly.sh: the night's run: take in the clearing file, match it to the
# day's authorizations, price the interchange, post to the cards, move the
# disputes on, and settle with the merchants.
#
# usage: nightly.sh DATE, the business date as YYYYMMDD.
# environment: LCS_BIN, the programs; LCS_SCRIPTS, the scripts; LCS_DATA,
# the data.
#
# history
# 2011-04-04 rk  lcs-0102 original version.
# 2012-01-09 ah  lcs-0140 interchange.
# 2015-03-23 mt  lcs-0302 disputes before settlement.
# 2022-09-05 ah  lcs-0497 the card system's feed after posting.

DATE="$1"
export LC_ALL=C
BIN="$LCS_BIN"
DATA="$LCS_DATA"
WORK="$DATA/work/$DATE"
LOG="$DATA/logs/night.$DATE.log"
mkdir -p "$WORK" "$DATA/outbound" "$DATA/reports" "$DATA/logs"

# STEP CLRIN: the clearing file, checked batch by batch
export DD_CLRIN="$DATA/inbound/CLR.$DATE.txt"
export DD_PRESOUT="$WORK/pres.txt"
export DD_CLRREJ="$DATA/outbound/CLRREJ.$DATE.txt"
export DD_CLRRPT="$DATA/reports/clrin.$DATE.txt"
"$BIN/clrin" "$DATE" > "$DATA/reports/clrin.$DATE.out" 2>&1
rc=$?
echo "CLRIN clrin $rc" >> "$LOG"
if [ "$rc" -ge 8 ]; then exit "$rc"; fi

# STEP PRESSORT: the presentments in authorization code order
sort -s -t "|" -k1.29,1.34 "$WORK/pres.txt" > "$WORK/pres.byauth"
rc=$?
echo "PRESSORT sort $rc" >> "$LOG"
if [ "$rc" -ne 0 ]; then exit 16; fi

# STEP AUTHSORT: the day's authorizations in authorization code order
sort -s -t "|" -k1.1,1.6 "$DATA/inbound/AUTH.$DATE.txt" > "$WORK/auth.byauth"
rc=$?
echo "AUTHSORT sort $rc" >> "$LOG"
if [ "$rc" -ne 0 ]; then exit 16; fi

# STEP MATCH: each presentment marked by its authorization
export DD_PRESIN="$WORK/pres.byauth"
export DD_AUTHIN="$WORK/auth.byauth"
export DD_PRESOUT="$WORK/pres.matched"
export DD_MATCHRPT="$DATA/reports/match.$DATE.txt"
"$BIN/match" "$DATE" > "$DATA/reports/match.$DATE.out" 2>&1
rc=$?
echo "MATCH match $rc" >> "$LOG"
if [ "$rc" -ge 8 ]; then exit "$rc"; fi

# STEP FEELD: the fee schedule into the fee table
export DD_FEESCHED="$DATA/inbound/FEES.txt"
export DD_FEETAB="$WORK/fees.rel"
"$BIN/feeld" "$DATE" > "$DATA/reports/feeld.$DATE.out" 2>&1
rc=$?
echo "FEELD feeld $rc" >> "$LOG"
if [ "$rc" -ge 8 ]; then exit "$rc"; fi

# STEP IXFEE: the interchange of each matched presentment
export DD_PRESIN="$WORK/pres.matched"
export DD_FEETAB="$WORK/fees.rel"
export DD_PRESOUT="$WORK/pres.priced"
export DD_FEERPT="$DATA/reports/ixfee.$DATE.txt"
"$BIN/ixfee" "$DATE" > "$DATA/reports/ixfee.$DATE.out" 2>&1
rc=$?
echo "IXFEE ixfee $rc" >> "$LOG"
if [ "$rc" -ge 8 ]; then exit "$rc"; fi

# STEP POST: the purchases on the cards' accounts
export DD_PRESIN="$WORK/pres.priced"
export DD_ACCOUNTS="$DATA/master/accounts.idx"
export DD_POSTOUT="$DATA/outbound/POSTINGS.$DATE.txt"
export DD_POSTREJ="$DATA/outbound/POSTREJ.$DATE.txt"
export DD_POSTRPT="$DATA/reports/post.$DATE.txt"
"$BIN/post" "$DATE" > "$DATA/reports/post.$DATE.out" 2>&1
rc=$?
echo "POST post $rc" >> "$LOG"
if [ "$rc" -ge 8 ]; then exit "$rc"; fi

# STEP PAYMENTS: the card system's feed: payments and status changes
"$LCS_SCRIPTS/payments.sh" "$DATE"
rc=$?
if [ "$rc" -ge 8 ]; then exit "$rc"; fi

# STEP DISPUTES: the disputes moved on, the merchants' adjustments
"$LCS_SCRIPTS/disputes.sh" "$DATE"
rc=$?
if [ "$rc" -ge 8 ]; then exit "$rc"; fi

# STEP SETLSORT: the priced presentments in merchant order
sort -s -t "|" -k1.35,1.44 "$WORK/pres.priced" > "$WORK/pres.bymerchant"
rc=$?
echo "SETLSORT sort $rc" >> "$LOG"
if [ "$rc" -ne 0 ]; then exit 16; fi

# STEP ADJSORT: the adjustments in merchant order
sort -s -t "|" -k1.1,1.10 "$WORK/adjustments.txt" > "$WORK/adjustments.bymerchant"
rc=$?
echo "ADJSORT sort $rc" >> "$LOG"
if [ "$rc" -ne 0 ]; then exit 16; fi

# STEP SETTLE: what each merchant is paid
export DD_PRESIN="$WORK/pres.bymerchant"
export DD_ADJIN="$WORK/adjustments.bymerchant"
export DD_PAYOUT="$DATA/outbound/PAYOUT.$DATE.txt"
export DD_SETLRPT="$DATA/reports/settle.$DATE.txt"
"$BIN/settle" "$DATE" > "$DATA/reports/settle.$DATE.out" 2>&1
rc=$?
echo "SETTLE settle $rc" >> "$LOG"
if [ "$rc" -ge 8 ]; then exit "$rc"; fi

# STEP ACCTUL: the accounts as they stand, for the card system
export DD_ACCOUNTS="$DATA/master/accounts.idx"
export DD_ACCTOUT="$DATA/outbound/ACCOUNTS.$DATE.txt"
"$BIN/acctul" > "$DATA/reports/acctul.$DATE.out" 2>&1
rc=$?
echo "ACCTUL acctul $rc" >> "$LOG"
