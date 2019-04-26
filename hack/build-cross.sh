#!/bin/bash

go get -u github.com/golang/dep/cmd/dep
dep ensure

mkdir -p build/linux && mkdir -p build/darwin  && mkdir -p build/windows

## declare an array of components variable
declare -a arr=("api" "controller" "node" "ingress" "discovery" "exporter")

if [[ $1 != "" ]]; then
  arr=($1)
fi

## now loop through the components array
for i in "${arr[@]}"
do
 echo "Build '$i' version '$VERSION' for os '$OSTYPE'"
  if [[ "$OSTYPE" == "linux-gnu" || "$OSTYPE" == "linux-musl" ]]; then
    CGO_ENABLED=1 \
    GOOS=linux  go build -ldflags "-X main.Version=$VERSION" -o "build/linux/$i" "cmd/$i/$i.go"
  elif [[ "$OSTYPE" == "darwin"* ]]; then
    GOOS=darwin go build -ldflags "-X main.Version=$VERSION" -o "build/darwin/$i" "cmd/$i/$i.go"
  elif [[ "$OSTYPE" == "windows"* ]]; then
    GOOS=windows go build -ldflags "-X main.Version=$VERSION" -o "build/windows/$i" "cmd/$i/$i.go"
  fi
done
