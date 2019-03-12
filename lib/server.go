package lib

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"

	"github.com/gwbeacon/gwbeacon/lib/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func ConvertIP(ip string) string {
	if ip == "::1" {
		return "[::1]"
	}
	return ip
}

type Server interface {
	GetInfo() *rpc.ServerInfo
	GetServers(tp string) map[uint32]*rpc.ServerInfo
	Register(regAddr string, onIDChange func(uint32)) uint32
	Serve(opt ...grpc.ServerOption) error
	Stop()
}

type Connector interface {
	Server
	SessionStore
	BindStream(sid uint64, ctx interface{})
	GetStream(sid uint64) interface{}
	RemoveStream(sid uint64)
	MakeSessionID() uint64
	MakeMessageID() uint64
}

type server struct {
	sync.RWMutex
	info        rpc.ServerInfo
	regAddr     string
	serversInfo map[string]map[uint32]*rpc.ServerInfo
	services    []Service
	grpcServer  *grpc.Server
	onIDChange  func(uint32)
	registered  bool
}

type ServerInfo struct {
	Type string
	Port int
}

func DefaultOptions() []grpc.ServerOption {
	zapLogger, _ := zap.NewDevelopment(zap.Development(), zap.AddCaller())
	zap.RedirectStdLog(zapLogger)
	options := []grpc.ServerOption{
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(zapLogger),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(zapLogger),
		)),
	}
	return options
}

func NewServer(info ServerInfo, services []Service) Server {
	s := &server{
		info: rpc.ServerInfo{
			Type:     info.Type,
			Port:     int32(info.Port),
			Services: make([]*rpc.ServiceInfo, 0),
		},
		serversInfo: make(map[string]map[uint32]*rpc.ServerInfo),
		services:    services,
	}
	for _, service := range services {
		s.info.Services = append(s.info.Services, service.GetInfo())
	}
	return s
}

func (s *server) GetInfo() *rpc.ServerInfo {
	s.RLock()
	defer s.RUnlock()
	return &s.info
}

func (s *server) GetServers(tp string) map[uint32]*rpc.ServerInfo {
	s.RLock()
	defer s.RUnlock()
	if res, ok := s.serversInfo[tp]; ok {
		return res
	}
	return map[uint32]*rpc.ServerInfo{}
}

func (s *server) Register(regAddr string, onIDChange func(id uint32)) uint32 {
	s.Lock()
	if s.registered {
		s.Unlock()
		return s.info.ID
	}
	s.regAddr = regAddr
	idNotify := make(chan uint32, 1)
	s.onIDChange = func(id uint32) {
		idNotify <- id
		close(idNotify)
	}
	s.Unlock()
	go func() {
		for {
			err := s.register()
			if err != nil {
				log.Println(err)
				time.Sleep(5 * time.Second)
			}
		}
	}()
	id := <-idNotify
	onIDChange(id)
	s.Lock()
	s.onIDChange = onIDChange
	s.Unlock()
	return id
}

func (s *server) Serve(opts ...grpc.ServerOption) error {
	addr := fmt.Sprintf(":%d", s.info.Port)
	fmt.Println(addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer func() {
		_ = lis.Close()
	}()
	opts = append(opts, DefaultOptions()...)
	server := grpc.NewServer(opts...)
	if server == nil {
		return errors.New("create grpc server failed")
	}
	s.Lock()
	for _, service := range s.services {
		service.Register(server)
	}
	s.grpcServer = server
	s.Unlock()
	log.Println("start server", s.info)
	err = server.Serve(lis)
	return err
}

func (s *server) Stop() {
	s.Lock()
	defer s.Unlock()
	s.grpcServer.GracefulStop()
}

func (s *server) register() error {
	s.RLock()
	regAddr := s.regAddr
	serverInfo := &s.info
	s.RUnlock()
	conn, err := grpc.DialContext(context.Background(), regAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	fmt.Println("begin to register server", s.info)
	client := rpc.NewClusterClient(conn)
	stream, err := client.DoRegister(context.Background())
	if err != nil {
		_ = conn.Close()
		return err
	}
	err = stream.Send(serverInfo)
	if err != nil {
		_ = conn.Close()
		return err
	}
	for {
		result, err1 := stream.Recv()
		if err1 != nil {
			_ = conn.Close()
			err = err1
			break
		}
		s.Lock()
		if result.ID > 0 {
			s.info.ID = result.ID
			s.onIDChange(result.ID)
		}
		for tp, idMap := range result.DownServers {
			if idMap == nil || len(idMap.Servers) == 0 {
				continue
			}
			for _, info := range idMap.Servers {
				if _, ok := s.serversInfo[tp]; ok {
					if _, ok := s.serversInfo[tp][info.ID]; ok {
						delete(s.serversInfo[tp], info.ID)
					}
				}
			}
		}

		for tp, idMap := range result.NewServers {
			s.serversInfo[tp] = idMap.Servers
		}
		s.Unlock()
	}
	return err
}
