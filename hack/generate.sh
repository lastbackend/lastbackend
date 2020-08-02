#!/usr/bin/env bash

protoc -I/usr/local/Cellar/protobuf/3.12.4/include  -I. \
  -I${GOPATH}/src \
  -I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.6/third_party/googleapis \
  -I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.6 \
  --grpc-gateway_out=logtostderr=true:./internal \
  --swagger_out=allow_merge=true,merge_file_name=api:./api/openapi \
  --go_out=plugins=grpc:./internal ./api/proto/v1/*.proto