.PHONY : build

NAME = deployit
HARDWARE = $(shell uname -m)
OS := $(shell uname)
VERSION ?= 0.1.0

build:
	echo "Building Deploy It"
	go get -u github.com/tools/godep
	godep restore
	mkdir -p build/linux  && GOOS=linux  go build -ldflags "-X main.Version=$(VERSION)" -o build/linux/$(NAME)
	mkdir -p build/darwin && GOOS=darwin go build -ldflags "-X main.Version=$(VERSION)" -o build/darwin/$(NAME)

install:
	echo "Install Deploy IT, ${OS} version:= ${VERSION}"
ifeq ($(OS),Linux)
	mv build/linux/$(NAME) /usr/local/bin/
endif
ifeq ($(OS) ,Darwin)
	mv build/darwin/$(NAME) /usr/local/bin/
endif
	chmod +x /usr/local/bin/$(NAME)
