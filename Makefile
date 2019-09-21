.PHONY : default deps test build image docs

export VERSION = 0.1.0-beta1

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
	mkdir -p build/linux && mkdir -p build/darwin && mkdir -p build/windows
	@echo "== Building Last.Backend platform: ${APP}"
	@bash ./hack/build-cross.sh ${APP}

build-plugin:
	@echo "== Pre-building configuration"
	mkdir -p build/linux && mkdir -p build/plugins
	@echo "== Building Last.Backend plugin"
	@bash ./hack/build-plugin.sh

install:
	@echo "== Install binaries"
	@bash ./hack/install-cross.sh ${APP}

image:
	@echo "== Pre-building configuration"
	@sh ./hack/build-images.sh $(app)

image-develop:
	@echo "== Pre-building configuration"
	@sh ./hack/build-images.sh $(app)

run-api:
	@echo "== Run lastbackend rest api"
	@go run ./cmd/api/api.go

run-ctl:
	@echo "== Run lastbackend cluster controller"
	@go run ./cmd/controller/controller.go

run-dns:
	@echo "== Run lastbackend dns daemon"
	@go run ./cmd/discovery/discovery.go

run-exp:
	@echo "== Run lastbackend exporter daemon "
	@go run ./cmd/discovery/discovery.go

run-ing:
	@echo "== Run lastbackend ingress proxy"
	@go run ./cmd/ingress/ingress.go

run-node:
	@echo "== Run node"
	@go run ./cmd/node/node.go

swagger-spec:
	@echo "== Generating Swagger spec for Last.Backend API"
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	swagger generate spec -b ./cmd/kit -m -o ./swagger.json
