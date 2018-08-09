package service

import (
	"fmt"

	"github.com/gwbeacon/gwbeacon/lib"
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

func (s *MessageService) Version() int32 {
	return int32(v1.SdkVersion_V1)
}

func (s *MessageService) Type() int32 {
	return int32(v1.FeatureType_FeatureTypeMessage)
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
	p, _ := peer.FromContext(stream.Context())
	fmt.Println(p.Addr.String())
	for {
		stream.Recv()
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
