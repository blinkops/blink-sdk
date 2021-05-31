######################################################################
# This Dockerfile is not built, only used to check compilation in CI
######################################################################
FROM golang:1.16.3

WORKDIR /go/src/github.com/blinkops/plugin-sdk
COPY .. .
RUN go build ./cmd/main
