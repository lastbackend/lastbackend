#!/usr/bin/env bash
echo "install dep package manager"
go get -u github.com/golang/dep/cmd/dep
echo "install dependencies"
dep ensure