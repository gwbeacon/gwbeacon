package main

import (
	"log"

	"sync"

	"github.com/gwbeacon/gwbeacon/lib"
	"github.com/gwbeacon/gwbeacon/lib/rpc"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

type sessionServer struct {
	sync.RWMutex
	lib.Server
	cache   lib.SessionStore
	servers map[uint64]uint16
	idMaker lib.IDMaker
}

func (ss *sessionServer) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	return ctx
}

func (ss *sessionServer) HandleRPC(ctx context.Context, st stats.RPCStats) {
}

func (ss *sessionServer) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	addr := info.RemoteAddr.String()
	log.Println("new connection", addr)

	sid := ss.idMaker.MakeID()
	session := lib.NewSession(uint64(sid), addr)
	ctx = context.WithValue(ctx, lib.ContextSessionKey, session)
	ctx = context.WithValue(ctx, lib.ContextServerKey, lib.Server(ss))
	log.Println("open session:", sid)
	return ctx
}

func (ss *sessionServer) HandleConn(ctx context.Context, st stats.ConnStats) {
	sid, ok := ctx.Value(lib.ContextSessionKey).(uint64)
	if !ok {
		return
	}
	switch st.(type) {
	case *stats.ConnEnd:
		log.Println("close session:", sid)
	default:
		log.Printf("illegal ConnStats type\n")
	}
}

func (ss *sessionServer) Save(ctx context.Context, s *rpc.Session) (*rpc.SessionResult, error) {
	sid, _ := ctx.Value(lib.ContextSessionKey).(uint64)
	ss.Lock()
	if serverID, ok := ss.servers[sid]; !ok {
		serverID = lib.ID(s.ID).GetServerID()
		ss.servers[sid] = serverID
	}
	ss.Unlock()
	log.Println("save session to session store server,", s)
	err := ss.cache.Save(s)
	return &rpc.SessionResult{}, err
}

func (ss *sessionServer) Update(ctx context.Context, s *rpc.Session) (*rpc.SessionResult, error) {
	err := ss.cache.Update(s)
	log.Println("update", s, err)
	return &rpc.SessionResult{}, err
}

func (ss *sessionServer) Remove(ctx context.Context, s *rpc.Session) (*rpc.SessionResult, error) {
	log.Println("remove session from session store server,", s)
	res, err := ss.cache.Remove(s)
	if err != nil {
		return &rpc.SessionResult{}, err
	}
	return &rpc.SessionResult{Data: res}, nil
}

func (ss *sessionServer) Replace(ctx context.Context, s *rpc.Session) (*rpc.SessionResult, error) {
	res, err := ss.cache.Replace(s)
	if err != nil {
		return &rpc.SessionResult{}, err
	}
	return &rpc.SessionResult{Data: res}, nil
}

func (ss *sessionServer) Stat(ctx context.Context, s *rpc.SessionStatRequest) (*rpc.SessionStat, error) {
	return ss.cache.Stat()
}

func (ss *sessionServer) Get(ctx context.Context, s *rpc.Session) (*rpc.SessionResult, error) {
	res, err := ss.cache.Get(s)
	if err != nil {
		return &rpc.SessionResult{}, err
	}
	return &rpc.SessionResult{Data: res}, nil
}

type sessionService struct {
	s *sessionServer
}

func (ss *sessionService) Register(gs *grpc.Server) {
	rpc.RegisterSessionStoreServer(gs, ss.s)
}

func (ss *sessionService) GetInfo() *rpc.ServiceInfo {
	return &rpc.ServiceInfo{
		Name:    lib.FeatureSession,
		Version: 1,
	}
}

func (ss *sessionServer) onIDChange(id uint32) {
	ss.Lock()
	if ss.idMaker == nil {
		ss.idMaker = lib.NewIDMaker(id, lib.SessionIDType, 0)
	}
	ss.Unlock()
	ss.idMaker.SetServerID(id)
}

func (ss *sessionServer) Serve(opt ...grpc.ServerOption) error {
	return ss.Server.Serve(grpc.StatsHandler(ss))
}

func NewServer(port int32, regAddr string) lib.Server {
	info := lib.ServerInfo{
		Type: lib.FeatureSession,
		Port: port,
	}
	s := &sessionServer{
		cache:   lib.NewSessionStoreServer(),
		servers: make(map[uint64]uint16),
	}
	service := &sessionService{
		s: s,
	}
	s.Server = lib.NewServer(info, []lib.Service{service})

	id := s.Register(regAddr, s.onIDChange)
	log.Println(id)
	return s
}
