#!/usr/bin/env bash

echo "===> Test common libs:"
go test ./libs/...
echo "\n"

echo "===> Test common packages:"
go test ./pkg/...
echo "\n"