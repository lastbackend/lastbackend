#!/bin/bash

mkdir -p build/linux && mkdir -p build/darwin  && mkdir -p build/windows

## declare an array of components variable

echo "Build 'lastbackend' version '$VERSION' for os '$OSTYPE'"
if [[ "$OSTYPE" == "linux-gnu" || "$OSTYPE" == "linux-musl" ]]; then
  CGO_ENABLED=1 \
  GOOS=linux  go build -ldflags "-X main.Version=$VERSION" -o "build/linux/lastbackend" "cmd/lastbackend/lastbackend.go"
elif [[ "$OSTYPE" == "darwin"* ]]; then
  GOOS=darwin go build -ldflags "-X main.Version=$VERSION" -o "build/darwin/lastbackend" "cmd/lastbackend/lastbackend.go"
elif [[ "$OSTYPE" == "windows"* ]]; then
  GOOS=windows go build -ldflags "-X main.Version=$VERSION" -o "build/windows/lastbackend" "cmd/lastbackend/lastbackend.go"
fi