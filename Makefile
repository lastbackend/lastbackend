.PHONY : default deps test build install

NAME_DAEMON = lbd
NAME_AGENT = lba
NAME_CLI = lbc

HARDWARE = $(shell uname -m)
OS := $(shell uname)
VERSION ?= 0.1.0

default: deps test build

deps:
	echo "Configuring Last.Backend"
	go get -u github.com/tools/godep
	godep restore

test:
	@echo "Testing Last.Backend"
	@sh ./test/test.sh

build:
	echo "Pre-building configuration"
	mkdir -p build/linux && mkdir -p build/darwin
	echo "Building Last.Backend daemon"
	GOOS=linux  go build -ldflags "-X main.Version=$(VERSION)" -o build/linux/$(NAME_DAEMON) cmd/daemon/daemon.go
	GOOS=darwin go build -ldflags "-X main.Version=$(VERSION)" -o build/darwin/$(NAME_DAEMON) cmd/daemon/daemon.go
	echo "Building Last.Backend node agent"
	GOOS=linux  go build -ldflags "-X main.Version=$(VERSION)" -o build/linux/$(NAME_DAEMON) cmd/agent/agent.go
	GOOS=darwin go build -ldflags "-X main.Version=$(VERSION)" -o build/darwin/$(NAME_DAEMON) cmd/agent/agent.go
	echo "Building Last.Backend CLI"
	GOOS=linux  go build -ldflags "-X main.Version=$(VERSION)" -o build/linux/$(NAME_CLI) cmd/client/client.go
	GOOS=darwin go build -ldflags "-X main.Version=$(VERSION)" -o build/darwin/$(NAME_CLI) cmd/client/client.go

install:
	echo "Install Last.Backend, ${OS} version:= ${VERSION}"
ifeq ($(OS),Linux)
	mv build/linux/$(NAME_CLI) /usr/local/bin/$(NAME_CLI)
	mv build/linux/$(NAME_DAEMON) /usr/local/bin/$(NAME_DAEMON)
	mv build/linux/$(NAME_AGENT) /usr/local/bin/$(NAME_AGENT)
endif
ifeq ($(OS) ,Darwin)
	mv build/darwin/$(NAME_CLI) /usr/local/bin/$(NAME_CLI)
	mv build/darwin/$(NAME_DAEMON) /usr/local/bin/$(NAME_DAEMON)
	mv build/darwin/$(NAME_AGENT) /usr/local/bin/$(NAME_AGENT)
endif

run-agent:
	go run cmd/agent/agent.go daemon -d


