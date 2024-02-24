#!/bin/bash

mkdir -p ./proto/go
mkdir -p ./proto/skyhook
protoc --go_out=./proto/ --go_opt=paths=source_relative --go-grpc_out=./proto/ --go-grpc_opt=paths=source_relative -I=../protobuf ../protobuf/skyhook/*