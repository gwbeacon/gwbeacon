#!/bin/sh
killproc(){
    for pid in `ps aux|grep $1 | grep -v grep |awk '{print $2}'`
	do
		kill -9 $pid
	done
}

quit(){
	echo "quit"
	killproc ./target/bin
}

trap 'quit'  INT

killproc ./target/bin
./build.sh
./target/bin/register > register.log &
sleep 1
./target/bin/session > session.log &
sleep 1
./target/bin/connector > connector.log
