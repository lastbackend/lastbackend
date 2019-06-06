#!/bin/bash

declare -a arr=("lastbackend")

if [[ $1 != "" ]]; then
  arr=($1)
fi

## now loop through the components array
for i in "${arr[@]}"
do
 echo "Build '$i' version '$VERSION'"
 docker build -t "index.lstbknd.net/lastbackend/$i" -f "./images/$i/Dockerfile" .
done
