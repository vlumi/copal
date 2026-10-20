#!/bin/sh
# Runs Laurasia's night on the fixture in a fresh laurasia/run/night and
# shows its logs and reports. Everything the night wrote is under
# laurasia/run/night/data.
set -eu

go run ./laurasia/cmd/laurasia night
