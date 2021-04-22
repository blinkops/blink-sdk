#!/bin/bash

# Prerequisits are
# brew install protobuf
# go get -u github.com/golang/protobuf/protoc-gen-go
# go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

protoc --go_out=. --go_opt=paths=source_relative  --go-grpc_out=. --go-grpc_opt=paths=source_relative plugin/proto/plugin.proto
