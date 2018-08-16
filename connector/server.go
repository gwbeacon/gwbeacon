package main

import (
	"fmt"
	"log"

	"sync"

	_ "github.com/gwbeacon/gwbeacon/connector/service"
	"github.com/gwbeacon/gwbeacon/lib"
	"github.com/gwbeacon/gwbeacon/lib/rpc"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

type server struct {
	sync.RWMutex
	lib.Server
	lib.SessionStore
	streams  map[uint64]interface{}
	sidMaker lib.IDMaker
	midMaker lib.IDMaker
}

func (s *server) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	return ctx
}

func (s *server) HandleRPC(ctx context.Context, st stats.RPCStats) {
}

func (s *server) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	addr := info.RemoteAddr.String()
	log.Println("new connection", addr)

	sid := s.sidMaker.MakeID()
	session := lib.NewSession(uint64(sid), addr)
	ctx = context.WithValue(ctx, lib.ContextSessionKey, session)
	ctx = context.WithValue(ctx, lib.ContextServerKey, lib.Connector(s))
	s.Save(session)
	log.Println("open session:", session)
	return ctx
}

func (s *server) HandleConn(ctx context.Context, st stats.ConnStats) {
	session, ok := ctx.Value(lib.ContextSessionKey).(*rpc.Session)
	if !ok {
		return
	}
	switch st.(type) {
	case *stats.ConnEnd:
		s.Remove(session)
		log.Println("close session:", session)
	default:
		log.Printf("illegal ConnStats type\n")
	}
}

func (s *server) BindStream(sid uint64, stream interface{}) {
	s.Lock()
	defer s.Unlock()
	s.streams[sid] = stream
}

func (s *server) GetStream(sid uint64) interface{} {
	s.Lock()
	defer s.Unlock()
	return s.streams[sid]
}

func (s *server) RemoveStream(sid uint64) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.streams[sid]; ok {
		delete(s.streams, sid)
	}
}

func (s *server) MakeSessionID() uint64 {
	return uint64(s.sidMaker.MakeID())
}

func (s *server) MakeMessageID() uint64 {
	return uint64(s.midMaker.MakeID())
}

func (s *server) onIDChange(id uint32) {
	info := s.GetInfo()
	s.Lock()
	if s.sidMaker == nil {
		s.sidMaker = lib.NewIDMaker(id, lib.SessionIDType, info.TimeBase)
	}
	if s.midMaker == nil {
		s.midMaker = lib.NewIDMaker(id, lib.MessageIDType, info.TimeBase)
	}
	s.Unlock()
	s.sidMaker.SetServerID(id)
	s.midMaker.SetServerID(id)
	servers := s.GetServers(lib.FeatureSession)
	if len(servers) > 0 {
		fmt.Println(servers)
		var sessionServer *rpc.ServerInfo
		for _, server := range servers {
			sessionServer = server
			break
		}
		addr := fmt.Sprintf("%s:%d", lib.ConvertIP(sessionServer.IP), sessionServer.Port)
		sessionStore, err := lib.NewSessionStoreClient(addr, true)
		if err == nil {
			s.SessionStore = sessionStore
		}
	}
	log.Println(servers)
}

func (s *server) Serve(opt ...grpc.ServerOption) error {
	return s.Server.Serve(grpc.StatsHandler(s))
}

func NewServer(port int32, regAddr string) lib.Connector {
	info := lib.ServerInfo{
		Port: port,
		Type: lib.FeatureConnector,
	}
	s := &server{
		Server:  lib.NewServer(info, lib.GetAllServices()),
		streams: make(map[uint64]interface{}),
	}
	id := s.Register(regAddr, s.onIDChange)
	log.Println(id)
	return s
}
