#!/bin/bash

rm -rf ./proto/go
mkdir -p ./proto
protoc -I=../protobuf/panel --go_out=./proto ../protobuf/panel/*