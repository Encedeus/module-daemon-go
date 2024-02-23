#!/bin/bash

mkdir -p ./proto/go
protoc --go_out=./proto/go --go_opt=paths=source_relative --go-grpc_out=./proto/go --go-grpc_opt=paths=source_relative -I=../protobuf ../protobuf/skyhook/*