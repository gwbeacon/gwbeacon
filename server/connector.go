package main

import (
	"context"
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/gwbeacon/gwbeacon/connector"
	"github.com/gwbeacon/gwbeacon/lib"
)

func register(registerAddr string) {

	info := &lib.ServerInfo{
		Type:     lib.ConnectorServer,
		Services: lib.GetAllServiceInfo(),
	}
	conn, err := grpc.DialContext(context.Background(), registerAddr, grpc.WithInsecure())

	if err != nil {
		log.Println("1", err)
		return
	}
	c1 := lib.NewClusterClient(conn)
	stream, err := c1.Register(context.Background(), info)
	if err != nil {
		log.Println("2", err)
		return
	}
	result, err := stream.Recv()
	if err != nil {
		log.Println("3", err)
	}
	log.Println(result)
}

func main() {
	var port = ""
	var registerAddr = ""
	flag.StringVar(&port, "port", "8888", "-port 8888")
	flag.StringVar(&registerAddr, "register", "localhost:9999", "-register localhost:9999")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println(err)
		return
	}

	server := connector.NewServer(0)
	go register(registerAddr)
	err = server.Serve(lis)
}
