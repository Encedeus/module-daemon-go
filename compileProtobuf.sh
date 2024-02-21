#!/bin/bash

rm -rf ./proto/go
protoc -I=../protobuf/panel --go_out=./proto ../protobuf/panel/*