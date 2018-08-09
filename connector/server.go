package connector

import (
	"log"

	_ "github.com/gwbeacon/gwbeacon/connector/service"
	"github.com/gwbeacon/gwbeacon/lib"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

type server struct {
	id uint16
	*grpc.Server
	sm *lib.SessionManager
}

func (s *server) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	return ctx
}
func (s *server) HandleRPC(ctx context.Context, st stats.RPCStats) {
}
func (s *server) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	sess := s.sm.OpenSession(s.id, info.RemoteAddr)
	log.Println(sess)
	return context.WithValue(ctx, "session", sess)
}
func (s *server) HandleConn(ctx context.Context, st stats.ConnStats) {
}

func (s *server) Type() lib.ServerType {
	return lib.ConnectorServer
}

func (s *server) SetID(id uint16) {
	s.id = id
}

func NewServer(serverId uint16) lib.Server {
	s := &server{
		id: serverId,
		sm: lib.NewSessionManager(),
	}
	s.Server = grpc.NewServer(grpc.StatsHandler(s))
	if s.Server == nil {
		log.Println("failed to create new grpc server")
		return nil
	}
	lib.LoadAllService(s.Server)
	return s
}
