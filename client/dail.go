package client

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func DefaultOptions() []grpc.DialOption {
	zapLogger, _ := zap.NewDevelopment(zap.Development(), zap.AddCaller())
	zap.RedirectStdLog(zapLogger)
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(grpc_zap.UnaryClientInterceptor(zapLogger))),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(grpc_zap.StreamClientInterceptor(zapLogger))),
	}
	return opts
}

func Dial(addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, DefaultOptions()...)
	return grpc.Dial(addr, opts...)
}
