package service

import (
	"fmt"

	"github.com/gwbeacon/gwbeacon/lib"
	"github.com/gwbeacon/gwbeacon/lib/rpc"
	"github.com/gwbeacon/sdk/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type MessageService struct {
}

func init() {
	lib.RegisterService(&MessageService{})
}

func (s *MessageService) Register(gs *grpc.Server) {
	v1.RegisterMessageServiceServer(gs, s)
}

func (s *MessageService) GetInfo() *rpc.ServiceInfo {
	return &rpc.ServiceInfo{
		Version: int32(v1.SdkVersion_V1),
		Name:    lib.FeatureMessageService,
	}
}

func (s *MessageService) OnAckMessage(stream v1.MessageService_OnAckMessageServer) error {
	p, _ := peer.FromContext(stream.Context())
	fmt.Println(p.Addr.String())
	for {
		stream.Recv()
	}
	return nil
}

func (s *MessageService) OnChatMessage(stream v1.MessageService_OnChatMessageServer) error {
	//ctx := stream.Context()
	//session := ctx.Value(lib.ContextSessionKey)
	for {
		_, err := stream.Recv()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *MessageService) OnHeartbeat(stream v1.MessageService_OnHeartbeatServer) error {
	p, _ := peer.FromContext(stream.Context())
	fmt.Println(p.Addr.String())
	for {
		stream.Recv()
	}
	return nil
}
