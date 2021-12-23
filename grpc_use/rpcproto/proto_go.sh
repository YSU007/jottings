#!/usr/bin/env bash

protoc *.proto --go_out=plugins=grpc:./
go run ../protogen/protogen.go $GOPATH