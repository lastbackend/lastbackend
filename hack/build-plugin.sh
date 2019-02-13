#!/bin/bash
docker build -f ./images/plugins/docker/Dockerfile.dev -t log-driver .

id=$(docker create log-driver true)
mkdir -p ./build/plugins/docker/rootfs
docker export "$id" | tar -x -C ./build/plugins/docker/rootfs
cp ./images/plugins/docker/config.json ./build/plugins/docker/config.json
docker plugin rm lastbackend
docker plugin create lastbackend ./build/plugins/docker
docker plugin enable lastbackend