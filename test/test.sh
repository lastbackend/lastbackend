#!/usr/bin/env bash

echo "===> Test client package:"
go test ./cmd/client/...
echo "\n"

echo "===> Test daemon package:"
go test ./cmd/daemon/...
echo "\n"

echo "===> Test builder package:"
go test ./cmd/builder/...
echo "\n"

echo "===> Test common libs:"
go test ./libs/...
echo "\n"

echo "===> Test common packages:"
go test ./pkg/...
echo "\n"

