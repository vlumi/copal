#!/bin/sh
# Checks that the dependencies between Copal's parts run one way: bedrock
# imports nothing else; survey and the specimens, pangaea and laurasia,
# import bedrock and never each other; and only console imports survey,
# pangaea and bedrock. Tests count too. Cobolith, which reads COBOL, is a
# library of its own that any of them may import.
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

forbid bedrock survey pangaea laurasia console
forbid survey pangaea laurasia console
forbid pangaea survey laurasia console
forbid laurasia survey pangaea console
exit $status
