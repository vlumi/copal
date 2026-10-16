#!/bin/sh
# Checks that the tools the build needs are installed, at the versions it expects.
set -eu

status=0

fail() {
    echo "missing: $1" >&2
    status=1
}

check_cobc() {
    if ! command -v cobc >/dev/null; then
        fail "GnuCOBOL 3.2 (brew install gnucobol, or apt install gnucobol3)"
        return
    fi
    version=$(cobc --version | sed -n '1s/.* \([0-9][0-9.]*\)$/\1/p')
    case "$version" in
        3.2|3.2.*) echo "cobc $version" ;;
        *) fail "GnuCOBOL 3.2 (found cobc $version)" ;;
    esac
}

check_command() {
    if command -v "$1" >/dev/null; then
        echo "$1: $(command -v "$1")"
    else
        fail "$1 ($2)"
    fi
}

check_cobc
check_command go "https://go.dev/dl/"

exit $status
