package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/gwbeacon/gwbeacon/server/connector/service"
	"github.com/gwbeacon/sdk/v1"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println(err)
		return
	}

	server := grpc.NewServer()
	service.LoadAll(server)
	go server.Serve(lis)

	time.Sleep(5 * time.Second)
	conn, err := grpc.DialContext(context.Background(), "localhost:8888", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}
	c1 := v1.NewQueryServiceClient(conn)
	result, err := c1.GetFeatureList(context.Background(), &v1.FeatureQuery{})
	log.Println(result, err)
	time.Sleep(time.Second)
}
