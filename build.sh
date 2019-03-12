#!/bin/sh

export GO111MODULE=on

TOP_PATH=$(cd `dirname $0`;pwd)
cd $TOP_PATH
mkdir -p target/bin

protoc --go_out=plugins=grpc:./ lib/rpc/*.proto

go build  -o target/bin/connector cmd/connector.go
go build -o target/bin/register cmd/register.go
go build -o target/bin/session cmd/session.go
go build -o target/bin/client cmd/client.go
