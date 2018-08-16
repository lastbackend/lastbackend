.PHONY : default deps test build image docs

export VERSION = 0.9.0

HARDWARE = $(shell uname -m)
OS := $(shell uname)

default: deps test build

deps:
	@echo "Configuring Last.Backend Dependencies"
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

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
	@bash ./hack/build-cross.sh

build-cli:
	@echo "== Pre-building cli configuration"
	mkdir -p build/linux && mkdir -p build/darwin
	@echo "== Building Last.Backend CLI"
	@bash ./hack/build-cross.sh cli

install:
	@echo "== Install binaries"
	@bash ./hack/install-cross.sh

image:
	@echo "== Pre-building configuration"
	@sh ./hack/build-images.sh

run-kit:
	@echo "== Run kit daemon all in one"
	@go run ./cmd/kit/kit.go $(app) --config=./contrib/config.yml

run-api:
	@echo "== Run lastbackend rest api daemon all in one"
	@go run ./cmd/kit/kit.go api --config=./contrib/config.yml

run-ctl:
	@echo "== Run lastbackend rest api daemon all in one"
	@go run ./cmd/kit/kit.go ctl --config=./contrib/config.yml


run-node:
	@echo "== Run node"
	@go run ./cmd/node/node.go --config=./contrib/node.yml

swagger-spec:
	@echo "== Generating Swagger spec for Last.Backend API"
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	swagger generate spec -b ./cmd/kit -m -o ./swagger.json
