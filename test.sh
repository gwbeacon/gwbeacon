#!/bin/sh
killproc(){
    ps aux|grep $1 | grep -v grep |awk '{print $2}' |xargs kill -9
}

killproc ./target/bin
./build.sh
./target/bin/register > register.log &
sleep 1
./target/bin/sessionStore > session.log &
sleep 1
./target/bin/connector > connector.log &
sleep 1

go test -v tests/connector_test.go

killproc ./target/bin
