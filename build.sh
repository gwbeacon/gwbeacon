#!/bin/sh

TOP_PATH=$(cd `dirname $0`;pwd)
cd $TOP_PATH
mkdir -p target/bin

protoc --go_out=plugins=grpc:./ lib/rpc/*.proto

go build  -o target/bin/connector connector/*.go
go build -o target/bin/register register/*.go
go build -o target/bin/sessionStore session/*.go
