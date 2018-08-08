package tests

import (
	"context"
	"net"
	"testing"

	"errors"

	"github.com/gwbeacon/gwbeacon/lib"
	"github.com/gwbeacon/sdk/v1"
	"google.golang.org/grpc"
)

func TestConnector(t *testing.T) {
	lis, err := net.Listen("tcp", ":8888")
	if err != nil {
		t.Fatal(err)
		return
	}

	server := grpc.NewServer()
	lib.LoadAllService(server)
	ch := make(chan int, 1)
	go func() {
		ch <- 1
		err = server.Serve(lis)
		if err != nil {
			t.Fatal(err)
		}
		return
	}()

	t.Log(server.GetServiceInfo())
	<-ch
	conn, err := grpc.DialContext(context.Background(), "localhost:8888", grpc.WithInsecure())
	if err != nil {
		t.Fatal()
		return
	}
	c1 := v1.NewQueryServiceClient(conn)
	result, err := c1.GetFeatureList(context.Background(), &v1.FeatureQuery{})
	if err != nil {
		t.Fatal(err)
		return
	}
	if result == nil || len(result.List) == 0 {
		t.Error(errors.New("empty result"))
	} else {
		t.Log(result)
	}
}
