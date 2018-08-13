package service

import (
	"fmt"

	"errors"

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
	ctx := stream.Context()
	session, ok := ctx.Value(lib.ContextSessionKey).(*rpc.Session)
	if !ok || session.User == nil || session.User.LoginTime == 0 {
		return errors.New("no session or not login")
	}
	connector, ok := ctx.Value(lib.ContextServerKey).(lib.Connector)
	if !ok {
		return errors.New("no connector")
	}
	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}
		msg.Id = connector.MakeMessageID()
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
