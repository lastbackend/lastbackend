.PHONY : build deps

NAME = deployit
HARDWARE = $(shell uname -m)
VERSION ?= 0.1.0

build:
	echo "Building Deploy It"
	go get
	mkdir -p build/linux  && GOOS=linux  go build -ldflags "-X main.Version=$(VERSION)" -o build/linux/$(NAME)
	mkdir -p build/darwin && GOOS=darwin go build -ldflags "-X main.Version=$(VERSION)" -o build/darwin/$(NAME)

deps:
	echo "Installing dependencies"
	go get
