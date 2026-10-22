#!/bin/sh
# Runs Gondwana's March on the fixture in a fresh gondwana/run/month and
# shows its logs and reports, transcoded from Shift_JIS for the terminal.
# Everything the month wrote is under gondwana/run/month/data.
set -eu

go run ./gondwana/cmd/gondwana month | iconv -f SHIFT_JIS -t UTF-8
