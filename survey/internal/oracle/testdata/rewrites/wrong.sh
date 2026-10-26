#!/bin/sh
# A replacement for BANKDATE that gets it wrong: it writes the control
# record back as it is, and returns 0 whatever happened.
cat "$DD_BANKIN" > "$DD_BANKOUT" 2>/dev/null
exit 0
