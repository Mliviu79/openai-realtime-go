#!/bin/bash

mods=(
    .
    ./contrib/ws-gorilla
)

for mod in "${mods[@]}"; do
    if [ -d "$mod" ]; then
        (cd "$mod" && "$@") || exit $?
    fi
done
