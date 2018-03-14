.PHONY : default deps test build install docs

NAME_KIT = lbd
NAME_CLI = lbc

HARDWARE = $(shell uname -m)
OS := $(shell uname)
VERSION ?= 0.9.0

default: deps test build

deps:
	@echo "Configuring Last.Backend Dependencies"
	go get -u github.com/kardianos/govendor
	govendor sync

test:
	@echo "Testing Last.Backend"
	@sh ./hack/run-coverage.sh

docs: docs/*
	@echo "Build Last.Backend Documentation"
	@sh ./hack/build-docs.sh

build:
	@echo "== Pre-building configuration"
	mkdir -p build/linux && mkdir -p build/darwin
	@echo "== Building Last.Backend platform"
	GOOS=linux  go build -ldflags "-X main.Version=$(VERSION)" -o build/linux/$(NAME_KIT) cmd/kit/kit.go
	GOOS=darwin go build -ldflags "-X main.Version=$(VERSION)" -o build/darwin/$(NAME_KIT) cmd/kit/kit.go

install:
	@echo "Install Last.Backend, ${OS} version:= ${VERSION}"
ifeq ($(OS),Linux)
	mv build/linux/$(NAME_KIT) /usr/local/bin/$(NAME_KIT)
endif
ifeq ($(OS) ,Darwin)
	mv build/darwin/$(NAME_KIT) /usr/local/bin/$(NAME_KIT)
endif

image:
	docker build -t lastbackend/lastbackend -f ./images/lastbackend/Dockerfile .

run:
	go run cmd/kit/kit.go --debug=3
