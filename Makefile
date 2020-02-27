.PHONY : default deps test build image docs

export VERSION = 0.1.0-beta2

HARDWARE = $(shell uname -m)
OS := $(shell uname)

default: test build


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

install:
	@echo "== Install binaries"
	@bash ./hack/install-cross.sh ${APP}

image:
	@echo "== Pre-building configuration"
	@sh ./hack/build-images.sh $(app)

image-develop:
	@echo "== Pre-building configuration"
	@sh ./hack/build-images.sh $(app)

run:
	@echo "== Run lastbackend container platform"
	@go run ./cmd/lastbackend/lastbackend.go -v=3

run-master:
	@echo "== Run lastbackend container platform master"
	@go run ./cmd/lastbackend/lastbackend.go master -v=3

run-minion:
	@echo "== Run lastbackend container platform minion"
	@go run ./cmd/lastbackend/lastbackend.go minion -v=3

swagger-spec:
	@echo "== Generating Swagger spec for Last.Backend API"
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	swagger generate spec -b ./cmd/kit -m -o ./swagger.json
