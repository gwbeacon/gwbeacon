package tests

import (
	"context"
	"net"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"

	"errors"

	"github.com/gwbeacon/gwbeacon/lib"
	"github.com/gwbeacon/sdk/v1"
	"google.golang.org/grpc"
)

func initServer(t *testing.T) {
	lis, err := net.Listen("tcp", ":8888")
	if err != nil {
		t.Fatal(err)
		return
	}

	server := grpc.NewServer()
	ch := make(chan int, 1)
	go func() {
		ch <- 1
		go func() {
			lib.LoadAllService(server)
		}()
		err = server.Serve(lis)
		if err != nil {
			t.Fatal(err)
		}
		return
	}()

	t.Log(server.GetServiceInfo())
	<-ch
}

func TestConnector(t *testing.T) {
	initServer(t)
	zapLogger, _ := zap.NewDevelopment(zap.Development(), zap.AddCaller())
	zap.RedirectStdLog(zapLogger)
	conn, err := grpc.DialContext(context.Background(), "localhost:8888", grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_zap.UnaryClientInterceptor(zapLogger),
		)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(grpc_zap.StreamClientInterceptor(zapLogger))),
	)

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
