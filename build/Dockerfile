######################################################################
# This Dockerfile is not built, only used to check compilation in CI
######################################################################
FROM golang:bullseye

WORKDIR /go/src/github.com/blinkops/blink-sdk
COPY .. .
RUN go build ./cmd/main
