package service

import (
    "errors"
    "fmt"
    "google.golang.org/grpc/grpclog"

    "time"

    "github.com/gwbeacon/gwbeacon/lib"
    "github.com/gwbeacon/gwbeacon/lib/rpc"
    "github.com/gwbeacon/sdk/v1"
    "google.golang.org/grpc"
    "google.golang.org/grpc/peer"
)

type MessageService struct {
	ackStream v1.MessageService_OnAckMessageServer
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
	ctx := stream.Context()
	session, ok := ctx.Value(lib.ContextSessionKey).(*rpc.Session)
	if !ok || session.User == nil || session.User.LoginTime == 0 {
		return errors.New("no session or not login")
	}
	s.ackStream = stream
	for {
		_, err := stream.Recv()
		if err != nil {
			grpclog.Error(err)
			stream.Context().Done()
			return err
		}
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
	connector.BindStream(session.ID, stream)
	for {
		msg, err := stream.Recv()
		if err != nil {
			stream.Context().Done()
			return err
		}
		ack := &v1.AckMessage{
			Domain: msg.Domain,
			To:     msg.From,
			From:   "",
			Id:     msg.Id,
		}
		err = s.ackStream.Send(ack)
		if err != nil {
			stream.Context().Done()
			return err
		}
		msg.Id = connector.MakeMessageID()
		sess := &rpc.Session{
			User: &rpc.UserInfo{
				Name:   msg.To,
				Domain: msg.Domain,
			},
		}
		ret, err := connector.Get(sess)
		if err == nil {
			for _, ss := range ret {
				st := connector.GetStream(ss.ID).(v1.MessageService_OnChatMessageServer)
				err := st.Send(msg)
				if err != nil {
				    grpclog.Error(err)
				    connector.Remove(ss)
                    st.Context().Done()
                }
			}
		}
	}
	return nil
}

func (s *MessageService) OnHeartbeat(stream v1.MessageService_OnHeartbeatServer) error {
	p, _ := peer.FromContext(stream.Context())
	fmt.Println(p.Addr.String())
	for {
		hb, err := stream.Recv()
		if err != nil {
			stream.Context().Done()
			return err
		}
		hb.From = ""
		hb.ServerTime = time.Now().Unix()
		err = stream.Send(hb)
		if err != nil {
			stream.Context().Done()
			return err
		}
	}
	return nil
}
