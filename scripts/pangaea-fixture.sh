#!/bin/sh
# Runs the daily batch on the fixture in a fresh pangaea/run/daily and shows
# the job log and reports. The datasets it created are unloaded as text into
# pangaea/run/daily/out.
set -eu

run=pangaea/run/daily
rm -rf "$run"
go run ./pangaea/cmd/pangaea build
go run ./pangaea/cmd/pangaea load -datasets "$run/datasets" pangaea/testdata/daily/input/*.txt
go run ./pangaea/cmd/pangaea run -datasets "$run/datasets" -sysout "$run/sysout" pangaea/jobs/DAILY.jcl
go run ./pangaea/cmd/pangaea unload -datasets "$run/datasets" "$run/out"
for report in "$run"/sysout/*.*; do
    printf '\n== %s\n' "$(basename "$report")"
    cat "$report"
done
