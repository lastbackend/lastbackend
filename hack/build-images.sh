#!/usr/bin/env bash

if [ $1 == "ingress" ] || [ $1 == ""] ; then
  docker build -t lastbackend/ingress -f ./images/ingress/Dockerfile .
fi

if [ $1 == "lastbackend" ] || [ $1 == "" ] ; then
  docker build -t lastbackend/lastbackend -f ./images/lastbackend/Dockerfile .
fi
