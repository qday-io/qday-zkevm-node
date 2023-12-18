set -x

init() {
    make build
    make build-docker
    cd test
    make run
    return
}

probeDoc() {
    exec >"$FUNCNAME.log" 2>&1
    cloc .
    find . -iname '*.md'
}

probeCode() {
    exec >"$FUNCNAME.log" 2>&1
    find . -iname 'main.go'
    find . -iname '*main*'
}

probe() {
    exec >"$FUNCNAME.log" 2>&1
    zkevm-node --version
    # zkevm-node version

    zkevm-node --help
    # zkevm-node help

    zkevm-node run --help

    zkevm-node approve --help
    # zkevm-node help run
    return
    # gh repo list -h
    # gh help formatting
    # return
    gh repo list 0xPolygonHermez \
        --limit 200
    # --json name
    # --jq '.'
    b2-zkevm-node --version
    b2-zkevm-node version

    b2-zkevm-node --help
    b2-zkevm-node help

    b2-zkevm-node run --help

    b2-zkevm-node approve --help
    b2-zkevm-node help run
    return
}

txsenderDemo(){
    exec >"$FUNCNAME.log" 2>&1
    txsender --help
}

$@
