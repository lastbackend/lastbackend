#!/usr/bin/env bash

#echo "===> Test common libs:"
#go test ./libs/...
#echo "\n"
#
#echo "===> Test common packages:"
#go test ./pkg/...
#echo "\n"#

echo "===> Test common packages:"
go test -v $(go list ./... | grep -v /vendor/)
echo "\n"