package main

import (
	"fmt"
	"net"

	_ "github.com/gwbeacon/gwbeacon/connector/service"
	"github.com/gwbeacon/gwbeacon/lib"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println(err)
		return
	}

	server := grpc.NewServer()
	lib.LoadAllService(server)
	err = server.Serve(lis)
}
