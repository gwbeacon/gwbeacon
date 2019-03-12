#!/bin/sh

./run.sh &

go test -v tests/connector_test.go
