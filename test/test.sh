#!/usr/bin/env bash

#echo "===> Test common libs:"
#go test ./libs/...
#echo "\n"
#
#echo "===> Test common packages:"
#go test ./pkg/...
#echo "\n"#

echo "===> Test common packages:"
godep go test -v $(go list ./... | grep -v /vendor/) || { echo 'test failed' ; exit 1; }
echo "\n"