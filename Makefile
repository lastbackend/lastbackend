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

generate:
	@mkdir -p "api_pb"
    protoc -I/usr/local/Cellar/protobuf/3.12.4/include  -I. \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway \
		--grpc-gateway_out=logtostderr=true:./api/proto \
		--swagger_out=allow_merge=true,merge_file_name=api:. \
		--go_out=plugins=grpc:./api_pb ./api/proto/*.proto

image:
	@echo "== Pre-building configuration"
	@sh ./hack/build-images.sh $(app)

image-develop:
	@echo "== Pre-building configuration"
	@sh ./hack/build-images.sh $(app)

run:
	@echo "== Run lastbackend container platform"
	@go run ./cmd/lastbackend/lastbackend.go daemon 

run-master:
	@echo "== Run lastbackend container platform master"
	@go run ./cmd/lastbackend/lastbackend.go daemon --no-schedule 

run-minion:
	@echo "== Run lastbackend container platform minion"
	@go run ./cmd/lastbackend/lastbackend.go minion -c config/linux/minion.yml -v=3

swagger-spec:
	@echo "== Generating Swagger spec for Last.Backend API"
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	swagger generate spec -b ./cmd/kit -m -o ./swagger.json
