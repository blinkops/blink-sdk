#!/bin/bash

# for older versions that are no longer available on brew
#VERSION=3.15.8
#PROTOC_ZIP=protoc-$VERSION-osx-x86_64.zip
#curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v$VERSION/$PROTOC_ZIP
#sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
#sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
#rm -f $PROTOC_ZIP




# Prerequisits are
# brew install protobuf
# go get -u github.com/golang/protobuf/protoc-gen-go

# gRPC
# protoc --go_out=. --go_opt=paths=source_relative  --go-grpc_out=. --go-grpc_opt=paths=source_relative plugin/proto/integration_pack.proto
# go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

# Twirp
# go install github.com/twitchtv/twirp/protoc-gen-twirp@latest
# protoc --go_out=. --go_opt=paths=source_relative --twirp_out=. plugin/proto/integration_pack.proto
# twirp_out apparently doesn't support "source_relative", so run this as well:
# mv github.com/blinkops/blink-sdk/proto/integration_pack.twirp.go plugin/proto/
