#!/bin/sh
set -e

OUT=out
PKG=github.com/X3NOOO/maud
BIN=$OUT/maud
CONFIG=maud.toml

mkdir $OUT 2>/dev/null || true

config() {
    cp $CONFIG.example $OUT/$CONFIG
    "${EDITOR:-nano}" $OUT/$CONFIG
}

release() {
    go mod download
    go build -o $BIN $PKG
}

run() {
    ./$BIN --config $OUT/$CONFIG
}

if [ $# -eq 0 ]; then
    echo "Usage: $0 [function1] [function2] ..."
    exit 1
fi

for func in "$@"; do
    if [ "$(type -t $func)" == "function" ]; then
        $func
    else
        echo "Function '$func' not found."
    fi
done
