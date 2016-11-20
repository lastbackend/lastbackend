.PHONY : default deps test build install

NAME = lastbackend
HARDWARE = $(shell uname -m)
OS := $(shell uname)
VERSION ?= 0.1.0

default: deps; test; build;

deps:
	echo "Configuring Last.Backend"
	go get -u github.com/tools/godep
	godep restore

test:
	@echo "Testing Last.Backend"
	@sh ./test/test.sh

build:
	echo "Building Last.Backend"
	mkdir -p build/linux  && GOOS=linux  go build -ldflags "-X main.Version=$(VERSION)" -o build/linux/$(NAME)
	mkdir -p build/darwin && GOOS=darwin go build -ldflags "-X main.Version=$(VERSION)" -o build/darwin/$(NAME)

install:
	echo "Install Last.Backend, ${OS} version:= ${VERSION}"
ifeq ($(OS),Linux)
	mv build/linux/$(NAME) /usr/local/bin/lb
endif
ifeq ($(OS) ,Darwin)
	mv build/darwin/$(NAME) /usr/local/bin/lb
endif
	chmod +x /usr/local/bin/$(NAME)


