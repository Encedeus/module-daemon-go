#!/bin/bash

rm -rf ./proto/go
mkdir -p ./proto
protoc --go_out=./proto --go_opt=paths=source_relative --go-grpc_out=./proto --go-grpc_opt=paths=source_relative -I=../protobuf/panel ../protobuf/panel/*