#!/usr/bin/env bash

## declare an array of components variable
declare -a arr=("kit" "agent" "cli")

## now loop through the components array
for i in "${arr[@]}"
do
   echo "Build $i"
   GOOS=linux  go build -ldflags "-X main.Version=$(VERSION)" -o "build/linux/$i" "cmd/$i/$i.go"
   GOOS=darwin go build -ldflags "-X main.Version=$(VERSION)" -o "build/darwin/$i" "cmd/$i/$i.go"
done