package connector

import (
	"log"

	"sync"

	_ "github.com/gwbeacon/gwbeacon/connector/service"
	"github.com/gwbeacon/gwbeacon/lib"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

type server struct {
	sync.RWMutex
	lib.Server
	sidMaker lib.IDMaker
	midMaker lib.IDMaker
	ss       lib.SessionStore
}

func (s *server) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	return ctx
}
func (s *server) HandleRPC(ctx context.Context, st stats.RPCStats) {
}
func (s *server) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	sid := s.sidMaker.MakeID()
	session := lib.NewSession(sid, info.RemoteAddr.String())
	log.Println("new session:", session)
	return context.WithValue(ctx, lib.ContextSessionKey, session)
}
func (s *server) HandleConn(ctx context.Context, st stats.ConnStats) {
}

func (s *server) onIDChange(id uint32) {
	s.Lock()
	defer s.Unlock()
	if s.sidMaker == nil {
		s.sidMaker = lib.NewIDMaker(id, lib.SessionIDType)
	}
	if s.midMaker == nil {
		s.midMaker = lib.NewIDMaker(id, lib.MessageIDType)
	}
	s.sidMaker.SetConnectorID(id)
	s.midMaker.SetConnectorID(id)
}

func (s *server) Serve(opt ...grpc.ServerOption) error {
	return s.Server.Serve(grpc.StatsHandler(s))
}

func NewServer(port int32, regAddr string) lib.Server {
	info := lib.ServerInfo{
		Port: port,
		Type: lib.FeatureConnector,
	}
	s := &server{
		Server: lib.NewServer(info, lib.GetAllServices()),
		ss:     lib.NewSessionStore(),
	}
	id := s.Register(regAddr, s.onIDChange)
	log.Println(id)
	return s
}
