#!/usr/bin/env bash

echo "===> Test common packages:"
go test -v $(go list ./... | grep -v /vendor/) || { echo 'test failed' ; exit 1; }
echo "\n"