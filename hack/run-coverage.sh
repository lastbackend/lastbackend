#!/usr/bin/env bash

echo "===> Test common packages:"
set -e
echo "mode: atomic" > coverage.txt

for d in $(go list ./... | grep -v vendor); do
    go test -race -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out | grep -v "mode: atomic" >> coverage.txt
        rm profile.out
    fi
done
