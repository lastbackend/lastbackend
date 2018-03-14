#!/bin/bash

## declare an array of components variable
declare -a arr=("kit" "node")

## now loop through the components array
for i in "${arr[@]}"
do
    echo "Install '$i'"
    if [[ "$OSTYPE" == "linux-gnu" || "$OSTYPE" == "linux-musl" ]]; then
        mv  build/linux/$i /usr/bin/$i
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        mv  build/darwin/$i /usr/bin/$i
    fi
done
