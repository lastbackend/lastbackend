#!/usr/bin/env bash

docker build -t lastbackend/lastbackend -f ./images/lastbackend/Dockerfile .

docker build -t lastbackend/ingress -f ./images/ingress/Dockerfile .

docker build -t lastbackend/discovery -f ./images/discovery/Dockerfile .
