#!/bin/sh

TOP_PATH=$(cd `dirname $0`;pwd)
cd $TOP_PATH
mkdir -p target/bin

protoc --go_out=plugins=grpc:./ lib/*.proto

go build  -o target/bin/connector server/connector.go
