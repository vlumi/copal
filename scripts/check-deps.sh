#!/bin/sh
# Checks that the dependencies between Copal's parts run one way: bedrock
# imports nothing else, survey and pangaea import bedrock and never each
# other, and only console imports all three. Tests count too.
set -eu

module=github.com/vlumi/copal
status=0

forbid() {
    part=$1
    shift
    for forbidden in "$@"; do
        offenders=$(go list -f '{{$p := .ImportPath}}{{range .Imports}}{{$p}} {{.}}{{"\n"}}{{end}}{{range .TestImports}}{{$p}} {{.}}{{"\n"}}{{end}}{{range .XTestImports}}{{$p}} {{.}}{{"\n"}}{{end}}' "./$part/..." |
            grep " $module/$forbidden" || true)
        if [ -n "$offenders" ]; then
            echo "$part must not import $forbidden:"
            echo "$offenders" | sed 's/^/  /'
            status=1
        fi
    done
}

forbid bedrock survey pangaea console
forbid survey pangaea console
forbid pangaea survey console
exit $status
