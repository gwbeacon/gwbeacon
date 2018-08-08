#!/bin/sh

TOP_PATH=$(cd `dirname $0`;pwd)
cd $TOP_PATH
mkdir -p target/bin

go build  -o target/bin/connector connector/connector.go
